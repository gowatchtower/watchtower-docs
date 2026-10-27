#!/bin/bash
# Assembles the experimental multi-version preview in public/. Cloudflare
# Pages runs it as the build command for the preview project.
#
# The legacy MkDocs versions come from the gh-pages branch of
# nicholas-fedor/watchtower, copied as they are. This checkout is published as
# nightly and replaces the legacy rolling dev build, so both version selectors
# list nightly in place of dev.
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

./scripts/publish-version.sh \
    --tree public \
    --version nightly \
    --replaces dev \
    --base-url "${CF_PAGES_URL:-http://localhost:1314}"
