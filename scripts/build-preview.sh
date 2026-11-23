#!/bin/bash
# Assembles the experimental multi-version preview in public/. Cloudflare
# Pages runs it as the build command for the preview project.
#
# The legacy MkDocs versions come from the gh-pages branch of
# nicholas-fedor/watchtower, copied as they are. Each tag in this repository
# names the Watchtower release it documents, and its Hugo build replaces the
# legacy build of that release. The newest tag takes the latest alias unless
# a newer legacy release holds it. Releases without a tag keep their legacy
# build. This checkout is published as nightly and replaces the legacy rolling
# dev build, so both version selectors list nightly in place of dev.
#
# The tree is served from CF_PAGES_URL on Cloudflare Pages, and from
# http://localhost:1314 otherwise (see `task preview:serve`).

set -euo pipefail

LEGACY_REPO="https://github.com/nicholas-fedor/watchtower.git"
LEGACY_BRANCH="gh-pages"

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

legacy="$(mktemp -d)"
trap 'rm -rf "$legacy"' EXIT

echo "Fetching legacy docs from $LEGACY_REPO ($LEGACY_BRANCH)..."
git clone --quiet --depth 1 --single-branch --branch "$LEGACY_BRANCH" "$LEGACY_REPO" "$legacy"
echo "Legacy docs at commit $(git -C "$legacy" rev-parse --short HEAD)"

rm -rf public
mkdir public
tar -C "$legacy" --exclude=./.git -cf - . | tar -C public -xf -

base_url="${CF_PAGES_URL:-http://localhost:1314}"

# A shallow clone, such as the one Cloudflare Pages builds from, may lack the
# tags, and Hugo dates pages from their Git history. Both need the full
# history, which is a few megabytes.
if [ "$(git rev-parse --is-shallow-repository)" = true ]; then
    echo "Fetching the full history and tags into the shallow clone..."
    git fetch --quiet --unshallow --tags origin
fi

mapfile -t tags < <(git tag --list 'v[0-9]*' --sort=version:refname)

newest_legacy="$(
    python3 -I -c 'import json, sys; print("\n".join(e["version"] for e in json.load(open(sys.argv[1]))))' public/versions.json |
        grep -E '^v[0-9]' | sort -V | tail -n 1
)"

for tag in "${tags[@]}"; do
    alias_args=()
    if [ "$tag" = "${tags[-1]}" ] &&
        [ "$(printf '%s\n' "$tag" "$newest_legacy" | sort -V | tail -n 1)" = "$tag" ]; then
        alias_args=(--alias latest)
    fi
    ./scripts/publish-version.sh \
        --tree public \
        --version "$tag" \
        --ref "$tag" \
        "${alias_args[@]}" \
        --base-url "$base_url"
done

./scripts/publish-version.sh \
    --tree public \
    --version nightly \
    --replaces dev \
    --base-url "$base_url"
