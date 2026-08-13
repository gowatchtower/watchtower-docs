#!/bin/bash
# Builds tplprev.wasm for the notification template preview and copies
# wasm_exec.js into static/. If wasm_exec.js is not found in
# $GOROOT/lib/wasm/, it is downloaded from the Go repository.

set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

mkdir -p ./static

echo "Copying wasm_exec.js..."
WASM_EXEC_SRC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
WASM_EXEC_DEST="./static/wasm_exec.js"
WASM_EXEC_URL="https://raw.githubusercontent.com/golang/go/master/lib/wasm/wasm_exec.js"

if [ -f "$WASM_EXEC_SRC" ]; then
    cp "$WASM_EXEC_SRC" "$WASM_EXEC_DEST"
    echo "Copied wasm_exec.js to ./static/"
else
    echo "wasm_exec.js not found at $WASM_EXEC_SRC. Downloading from $WASM_EXEC_URL..."
    if ! curl -L -o "$WASM_EXEC_DEST" "$WASM_EXEC_URL"; then
        echo "Error: Failed to download wasm_exec.js from $WASM_EXEC_URL" >&2
        exit 1
    fi
    echo "Downloaded wasm_exec.js to ./static/"
fi

echo "Building tplprev.wasm..."
sanitize() {
    printf '%s' "$1" | tr -d '\n\r' | tr -c 'A-Za-z0-9._+-:' '_'
}

VERSION="$(sanitize "$(git describe --tags --always --dirty 2>/dev/null || true)")"
COMMIT="$(sanitize "$(git rev-parse --verify HEAD 2>/dev/null || true)")"
DATE="$(sanitize "$(git log -1 --format=%cI 2>/dev/null || true)")"
[ -n "$VERSION" ] || VERSION=dev
[ -n "$COMMIT" ] || COMMIT=none
[ -n "$DATE" ] || DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

LDFLAGS="-X github.com/gowatchtower/tplprev/internal/metadata.Version=${VERSION}"
LDFLAGS="${LDFLAGS} -X github.com/gowatchtower/tplprev/internal/metadata.Commit=${COMMIT}"
LDFLAGS="${LDFLAGS} -X github.com/gowatchtower/tplprev/internal/metadata.Date=${DATE}"
GOARCH=wasm GOOS=js go -C ./tools/tplprev build -ldflags "${LDFLAGS}" -o ../../static/tplprev.wasm .

echo "Files in ./static:"
ls -l ./static/tplprev.wasm ./static/wasm_exec.js
