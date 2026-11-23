#!/bin/bash
# Assembles a local preview of the published site in public/. It copies the
# cf-pages branch from origin, which holds every published version, and
# builds this checkout over its nightly version.
#
# Serve it with `task preview:serve` on http://localhost:1314.

set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

echo "Fetching cf-pages from origin..."
git fetch --quiet origin cf-pages

rm -rf public
mkdir public
git archive FETCH_HEAD | tar -C public -xf -

rm -rf public/nightly
HUGO_PARAMS_VERSIONING_VERSION=nightly hugo --gc --minify \
    --baseURL http://localhost:1314/nightly/ --destination public/nightly
