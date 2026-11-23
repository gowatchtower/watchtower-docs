#!/bin/bash
# Builds the Hugo site in the current directory as one version of a
# multi-version docs site kept on a branch.
#
# The build lands in a detached worktree of that branch and is not committed.
# A failed build leaves the worktree as the branch has it.
#
# The branch uses mike's layout. Each version has its own directory, and a
# versions.json beside them lists the versions. Versions that are not
# releases, such as nightly, are listed first, in the order they were added.
# Releases follow, newest first. A release is a name such as v1.24.0. The
# newest release holds the latest alias. A version is built into its own
# directory, and again into a directory for each alias it holds. Alias URLs
# then serve full pages rather than redirects.
#
# Usage:
#   build.sh --version VERSION --base-url URL --branch BRANCH
#
#   --version   Directory and versions.json entry to build, such as nightly or v1.24.0.
#   --base-url  URL the site branch is served from, such as https://example.org.
#   --branch    Branch that holds the site. It must already be fetched from origin.
#
# Requires git and hugo. The site reads its version from
# HUGO_PARAMS_VERSIONING_VERSION. The worktree's location is written to
# $GITHUB_OUTPUT as path.

set -euo pipefail

# A version or alias names a directory inside the site, so it is limited to
# one plain path segment.
readonly NAME_PATTERN='^[A-Za-z0-9][A-Za-z0-9._-]*$'
readonly RELEASE_PATTERN='^v[0-9]+(\.[0-9]+)*$'
readonly LATEST='latest'

# The versions listed in versions.json, in file order. Each version's title and
# space-separated aliases are keyed by its name.
versions=()
declare -A titles=()
declare -A aliases=()

# Directory the versions are built into before they replace the site's copies.
staging=""

die() {
    if [ "${GITHUB_ACTIONS:-}" = true ]; then
        echo "::error::$*"
    fi
    echo "Error: $*" >&2
    exit 1
}

usage() {
    sed -n '/^# Usage:/,/^#   --branch/p' "$0" | sed 's/^# \{0,1\}//' >&2
    exit 2
}

is_release() {
    [[ "$1" =~ $RELEASE_PATTERN ]]
}

check_name() {
    [[ "$1" =~ $NAME_PATTERN ]] || die "invalid version or alias name: $1"
}

# Reads the versions.json file $1 into versions, titles, and aliases. A missing
# file lists no versions.
#
# The file must be laid out with one field per line, as mike and this script
# write it. Anything else stops the build, including a field other than
# version, title, and aliases. The file is rewritten later, and a field this
# function skipped would be lost.
read_versions() {
    local -r file="$1"
    local -r field='^[[:space:]]*"(version|title)":[[:space:]]*"([^"\\]*)",?[[:space:]]*$'
    local -r alias_list='^[[:space:]]*"aliases":[[:space:]]*\[(.*)$'
    local -r alias_item='^[^"]*"([^"\\]*)"(.*)$'
    local -r punctuation='^[][{},[:space:]]*$'
    local line text item version="" title="" entry_aliases="" in_aliases=false number=0

    [ -f "$file" ] || return 0

    while IFS= read -r line || [ -n "$line" ]; do
        number=$((number + 1))
        text=""

        if [ "$in_aliases" = true ]; then
            text="$line"
        elif [[ "$line" =~ $alias_list ]]; then
            in_aliases=true
            text="${BASH_REMATCH[1]}"
        elif [[ "$line" =~ $field ]]; then
            if [ "${BASH_REMATCH[1]}" = version ]; then
                version="${BASH_REMATCH[2]}"
            else
                title="${BASH_REMATCH[2]}"
            fi
            continue
        elif [[ "$line" =~ $punctuation ]]; then
            # A closing brace ends the entry.
            if [[ "$line" == *"}"* ]]; then
                check_name "$version"
                [ -z "${titles[$version]+listed}" ] || die "$file lists $version twice"
                versions+=("$version")
                titles[$version]="${title:-$version}"
                aliases[$version]="${entry_aliases# }"
                version="" title="" entry_aliases=""
            fi
            continue
        else
            die "$file line $number is not supported: $line"
        fi

        # Collects the alias names on this line. The closing bracket ends the list.
        while [[ "$text" =~ $alias_item ]]; do
            item="${BASH_REMATCH[1]}"
            text="${BASH_REMATCH[2]}"
            check_name "$item"
            entry_aliases+=" $item"
        done
        if [[ "$text" == *"]"* ]]; then
            in_aliases=false
        fi
    done <"$file"

    [ -z "$version$title$entry_aliases" ] && [ "$in_aliases" = false ] ||
        die "$file ends inside an entry"
}

# Lists version $1 if it is not listed yet. A release that no listed release is
# newer than takes the latest alias from whichever version holds it.
record_version() {
    local -r version="$1"
    local name item newest kept

    if [ -z "${titles[$version]+listed}" ]; then
        versions+=("$version")
        titles[$version]="$version"
        aliases[$version]=""
    fi

    is_release "$version" || return 0
    newest="$(printf '%s\n' "${versions[@]}" | grep -E "$RELEASE_PATTERN" | sort --version-sort | tail -n 1)"
    [ "$newest" = "$version" ] || return 0

    for name in "${versions[@]}"; do
        kept=()
        for item in ${aliases[$name]}; do
            [ "$item" = "$LATEST" ] || kept+=("$item")
        done
        aliases[$name]="${kept[*]}"
    done
    aliases[$version]="${aliases[$version]:+${aliases[$version]} }$LATEST"
}

# Prints the listed versions in the order versions.json holds them.
ordered_versions() {
    local name
    for name in "${versions[@]}"; do
        is_release "$name" || echo "$name"
    done
    printf '%s\n' "${versions[@]}" | grep -E "$RELEASE_PATTERN" | sort --version-sort --reverse || true
}

# Writes the listed versions to the versions.json file $1.
write_versions() {
    local -r file="$1"
    local name item separator="" alias_separator

    {
        printf '['
        while IFS= read -r name; do
            printf '%s\n  {\n    "version": "%s",\n    "title": "%s",\n    "aliases": [' \
                "$separator" "$name" "${titles[$name]}"
            alias_separator=""
            for item in ${aliases[$name]}; do
                printf '%s\n      "%s"' "$alias_separator" "$item"
                alias_separator=","
            done
            if [ -n "$alias_separator" ]; then
                printf '\n    '
            fi
            printf ']\n  }'
            separator=","
        done < <(ordered_versions)
        printf '\n]\n'
    } >"$file.tmp"
    mv "$file.tmp" "$file"
}

# Builds version $1 for the base URL $2 into the site worktree $3. Every
# directory is built in a staging directory first. The site's copies are
# replaced only after all of them have built.
build_version() {
    local -r version="$1" base_url="$2" site="$3"
    local dir

    staging="$(mktemp -d "$site.build.XXXXXX")"

    for dir in "$version" ${aliases[$version]}; do
        echo "Building $version for $base_url/$dir/"
        HUGO_PARAMS_VERSIONING_VERSION="$version" hugo --gc --minify \
            --baseURL "$base_url/$dir/" --destination "$staging/$dir"
    done

    for dir in "$version" ${aliases[$version]}; do
        rm -rf "${site:?}/$dir"
        mv "$staging/$dir" "$site/$dir"
    done
}

main() {
    local version="" base_url="" branch="" site tool

    while [ $# -gt 0 ]; do
        case "$1" in
        --version) version="${2:-}"; shift 2 || usage ;;
        --base-url) base_url="${2:-}"; shift 2 || usage ;;
        --branch) branch="${2:-}"; shift 2 || usage ;;
        *) usage ;;
        esac
    done
    [ -n "$version" ] && [ -n "$base_url" ] && [ -n "$branch" ] || usage

    check_name "$version"
    [ "$version" != "$LATEST" ] || die "$LATEST is an alias and cannot be built as a version"
    [[ "$base_url" =~ ^https?://[^[:space:]]+$ ]] || die "invalid base URL: $base_url"
    base_url="${base_url%/}"

    for tool in git hugo; do
        command -v "$tool" >/dev/null || die "$tool is required"
    done

    git rev-parse --verify --quiet "refs/remotes/origin/$branch^{commit}" >/dev/null ||
        die "origin/$branch is not fetched. Check out the repository with fetch-depth: 0."

    site="${RUNNER_TEMP:-${TMPDIR:-/tmp}}/versioned-docs"
    [ ! -e "$site" ] || die "$site already exists"
    git worktree add --quiet --detach "$site" "origin/$branch"
    trap 'rm -rf "$staging"' EXIT

    read_versions "$site/versions.json"
    record_version "$version"
    build_version "$version" "$base_url" "$site"
    write_versions "$site/versions.json"

    echo "path=$site" >>"${GITHUB_OUTPUT:-/dev/null}"
}

# Runs only when executed, so the functions can be loaded on their own.
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
