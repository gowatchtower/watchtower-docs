#!/bin/bash
# Builds the notification template preview for the docs site.
#
# It compiles the tplprev Go module to tplprev.wasm, and copies wasm_exec.js,
# the Go runtime that loads the module in the browser. wasm_exec.js comes from
# the installed Go toolchain, so it matches the compiler that built the module.
#
# The build stamps the module's metadata package with the last commit that
# changed the module, that commit's date, and its `git describe` version. The
# output therefore stays byte-identical until the module changes, which keeps
# an unchanged preview from being committed again with each published build.
#
# Usage:
#   build.sh --module DIR --output DIR
#
#   --module  Directory of the tplprev Go module.
#   --output  Directory that receives tplprev.wasm and wasm_exec.js.
#
# Requires git and go.

set -euo pipefail

die() {
    if [ "${GITHUB_ACTIONS:-}" = true ]; then
        echo "::error::$*"
    fi
    echo "Error: $*" >&2
    exit 1
}

usage() {
    sed -n '/^# Usage:/,/^#   --output/p' "$0" | sed 's/^# \{0,1\}//' >&2
    exit 2
}

# Prints $1 with every character that is unsafe in a linker flag replaced.
sanitize() {
    printf '%s' "$1" | tr -d '\n\r' | tr -c 'A-Za-z0-9._+:-' '_'
}

# Copies wasm_exec.js from the Go toolchain into the directory $1. Go keeps the
# file in lib/wasm, and kept it in misc/wasm before Go 1.24.
copy_runtime() {
    local -r output="$1"
    local goroot dir

    goroot="$(go env GOROOT)"
    for dir in lib/wasm misc/wasm; do
        if [ -f "$goroot/$dir/wasm_exec.js" ]; then
            cp "$goroot/$dir/wasm_exec.js" "$output/wasm_exec.js"
            return 0
        fi
    done
    die "wasm_exec.js was not found in $goroot"
}

# Compiles the module in $1 to $2/tplprev.wasm. The build leaves out local
# paths and Go's own version control stamp, which would differ between builds
# of the same source.
build_module() {
    local -r module="$1" output="$2"
    local package version commit date flags

    commit="$(git log -1 --format=%H -- "$module")"
    [ -n "$commit" ] || die "$module has no commits"

    package="$(go -C "$module" list -m)/internal/metadata"
    version="$(sanitize "$(git describe --tags --always "$commit")")"
    date="$(sanitize "$(git log -1 --format=%cI "$commit")")"

    flags="-X $package.Version=$version -X $package.Commit=$commit -X $package.Date=$date"
    GOOS=js GOARCH=wasm go -C "$module" build -trimpath -buildvcs=false \
        -ldflags "$flags" -o "$output/tplprev.wasm" .
}

main() {
    local module="" output="" tool

    while [ $# -gt 0 ]; do
        case "$1" in
        --module) module="${2:-}"; shift 2 || usage ;;
        --output) output="${2:-}"; shift 2 || usage ;;
        *) usage ;;
        esac
    done
    [ -n "$module" ] && [ -n "$output" ] || usage

    for tool in git go; do
        command -v "$tool" >/dev/null || die "$tool is required"
    done
    [ -f "$module/go.mod" ] || die "$module is not a Go module"

    mkdir -p "$output"
    output="$(cd "$output" && pwd)"

    build_module "$module" "$output"
    copy_runtime "$output"
    echo "Built $output/tplprev.wasm and $output/wasm_exec.js."
}

main "$@"
