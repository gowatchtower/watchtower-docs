#!/bin/bash
# Builds this checkout as one version of the multi-version docs tree and
# records it in the tree's versions.json, replacing `mike deploy`.
#
# The tree uses mike's layout, so legacy MkDocs builds can live beside Hugo
# builds: each version in its own directory, plus a versions.json listing
# {version, title, aliases} entries that both the legacy MkDocs Material
# selector and the Hugo selector read.
#
# Usage:
#   publish-version.sh --tree DIR --version VERSION --base-url URL [--ref REF] [--alias NAME] [--replaces OLD]
#
#   --tree      Directory holding the assembled site. Created if missing.
#   --version   Directory and versions.json entry to publish, such as nightly or v1.24.0.
#   --base-url  URL the tree is served from, such as https://watchtower.nickfedor.com.
#   --ref       Git ref whose pages to publish, such as the tag v1.23.2. Its
#               content/ and static/ are built with this checkout's layouts,
#               so every version gets the current version selector. Without
#               it, the working tree is published as it is.
#   --alias     Alias to move to this version, such as latest. The version is
#               also built a second time into DIR/NAME/, so alias URLs serve
#               full pages rather than redirects.
#   --replaces  Version this one takes over, such as dev for nightly. Its
#               directory and versions.json entry are removed, and its aliases
#               move to this version.
#
# The script only changes DIR. Committing or uploading the tree is up to the caller.

set -euo pipefail

usage() {
    sed -n '/^# Usage:/,/^# The script/p' "$0" | sed 's/^# \{0,1\}//' >&2
    exit 2
}

tree=""
version=""
base_url=""
ref=""
alias_name=""
replaces=""

while [ $# -gt 0 ]; do
    case "$1" in
    --tree) tree="${2:-}"; shift 2 ;;
    --version) version="${2:-}"; shift 2 ;;
    --base-url) base_url="${2:-}"; shift 2 ;;
    --ref) ref="${2:-}"; shift 2 ;;
    --alias) alias_name="${2:-}"; shift 2 ;;
    --replaces) replaces="${2:-}"; shift 2 ;;
    *) usage ;;
    esac
done

[ -n "$tree" ] && [ -n "$version" ] && [ -n "$base_url" ] || usage

# Version, alias, and replaced names are directory names that get deleted, so
# each is limited to one plain path segment and they must all differ.
names=("$version" ${alias_name:+"$alias_name"} ${replaces:+"$replaces"})
for name in "${names[@]}"; do
    if ! [[ "$name" =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]]; then
        echo "Error: invalid version, alias, or replaced name: $name" >&2
        exit 2
    fi
done

if [ "$(printf '%s\n' "${names[@]}" | sort -u | wc -l)" -ne "${#names[@]}" ]; then
    echo "Error: --version, --alias, and --replaces must all differ" >&2
    exit 2
fi

for tool in hugo python3; do
    command -v "$tool" >/dev/null || { echo "Error: $tool is required" >&2; exit 1; }
done

repo_root="$(cd "$(dirname "$0")/.." && pwd)"

if [ -n "$ref" ] && ! git -C "$repo_root" rev-parse --verify --quiet "$ref^{commit}" >/dev/null; then
    echo "Error: unknown ref: $ref" >&2
    exit 2
fi

mkdir -p "$tree"
tree="$(cd "$tree" && pwd)"
base_url="${base_url%/}"

# Hugo builds from the working tree, or with --ref from a temporary worktree of
# the ref. That worktree keeps the ref's content/ and static/, so each version
# keeps its own pages, images, and template preview WASM, and Hugo dates its
# pages from the ref's history. Everything else comes from the working tree.
# That includes static/template-preview/, because the template-preview
# shortcode in layouts/ loads it.
source="$repo_root"
if [ -n "$ref" ]; then
    source="$(mktemp -d)"
    trap 'git -C "$repo_root" worktree remove --force "$source" 2>/dev/null || rm -rf "$source"' EXIT
    git -C "$repo_root" worktree add --quiet --detach "$source" "$ref"

    find "$source" -mindepth 1 -maxdepth 1 ! -name .git ! -name content ! -name static -exec rm -rf {} +
    rm -rf "$source/static/template-preview"

    # Copies the files git tracks or would track outside content/ and static/,
    # skipping tracked files deleted from the working tree.
    (
        cd "$repo_root"
        {
            git ls-files -z --cached --others --exclude-standard --deduplicate -- . ':!content' ':!static'
            git ls-files -z --cached --others --exclude-standard --deduplicate -- static/template-preview
        } |
            while IFS= read -r -d '' file; do
                [ -e "$file" ] && printf '%s\0' "$file"
            done |
            tar --null -T - -cf -
    ) | tar -C "$source" -xf -

    # Homepages tagged before versioning link from the host root, such as
    # link="/quickstart/", which skips the version's directory. Dropping the
    # leading slash makes them relative to the version's base URL.
    sed -i -E 's#(link|src)="/([^/])#\1="\2#g' "$source/content/_index.md"
fi

build() {
    local dir="$1"
    echo "Building $version${ref:+ from $ref} into $tree/$dir/ (base URL $base_url/$dir/)..."
    rm -rf "${tree:?}/$dir"
    HUGO_PARAMS_VERSIONING_VERSION="$version" hugo \
        --source "$source" \
        --gc \
        --minify \
        --baseURL "$base_url/$dir/" \
        --destination "$tree/$dir"
}

(cd "$source" && hugo mod get)

build "$version"
if [ -n "$alias_name" ]; then
    build "$alias_name"
fi
if [ -n "$replaces" ]; then
    echo "Removing $tree/$replaces/, replaced by $version..."
    rm -rf "${tree:?}/$replaces"
fi

# Writes versions.json in mike's format and order: entries that are not
# releases (such as nightly) first, then releases newest first. Rebuilding a
# version keeps its existing title, aliases, and properties.
python3 -I - "$tree/versions.json" "$version" "$alias_name" "$replaces" <<'EOF'
import json
import os
import re
import sys

path, version, alias, replaces = sys.argv[1:5]
release = re.compile(r"^v?[0-9]+(\.[0-9]+)*$")

try:
    with open(path) as f:
        entries = json.load(f)
except FileNotFoundError:
    entries = []

existing = next((e for e in entries if e["version"] == version), {})
replaced = next((e for e in entries if e["version"] == replaces), {})
entries = [e for e in entries if e["version"] not in (version, replaces)]
for entry in entries:
    entry["aliases"] = [a for a in entry["aliases"] if a != alias]

aliases = set(existing.get("aliases", [])) | set(replaced.get("aliases", []))
if alias:
    aliases.add(alias)
entries.append({
    **existing,
    "version": version,
    "title": existing.get("title", version),
    "aliases": sorted(aliases),
})

others = [e for e in entries if not release.match(e["version"])]
releases = sorted(
    (e for e in entries if release.match(e["version"])),
    key=lambda e: [int(part) for part in e["version"].lstrip("v").split(".")],
    reverse=True,
)

with open(path + ".tmp", "w") as f:
    json.dump(others + releases, f, indent=2)
    f.write("\n")
os.replace(path + ".tmp", path)
EOF

echo "Published $version${alias_name:+ (alias $alias_name)}${replaces:+, replacing $replaces,} to $tree"
