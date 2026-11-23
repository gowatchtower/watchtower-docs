#!/bin/bash
# Commits a docs version built by the build-versioned-docs action and pushes it
# to the site branch. A build that changed nothing is not committed.
#
# Usage:
#   publish.sh --site DIR --version VERSION --branch BRANCH
#
#   --site     Worktree of the site branch that holds the build.
#   --version  Version that was built, named in the commit message.
#   --branch   Branch that holds the site.
#
# Requires git. Reads the token for the push from $TOKEN. Pushes to the
# repository $GITHUB_REPOSITORY on $GITHUB_SERVER_URL, and names the commit
# $GITHUB_SHA in the message.

set -euo pipefail

readonly COMMITTER_NAME='github-actions[bot]'
readonly COMMITTER_EMAIL='41898282+github-actions[bot]@users.noreply.github.com'

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

main() {
    local site="" version="" branch="" name remote

    while [ $# -gt 0 ]; do
        case "$1" in
        --site) site="${2:-}"; shift 2 || usage ;;
        --version) version="${2:-}"; shift 2 || usage ;;
        --branch) branch="${2:-}"; shift 2 || usage ;;
        *) usage ;;
        esac
    done
    [ -n "$site" ] && [ -n "$version" ] && [ -n "$branch" ] || usage

    for name in TOKEN GITHUB_REPOSITORY GITHUB_SERVER_URL GITHUB_SHA; do
        [ -n "${!name:-}" ] || die "$name is not set"
    done
    git -C "$site" rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "$site is not a git worktree"

    git -C "$site" add --all
    if git -C "$site" diff --cached --quiet; then
        echo "The build of $version changed nothing on $branch."
        return 0
    fi

    git -C "$site" -c user.name="$COMMITTER_NAME" -c user.email="$COMMITTER_EMAIL" \
        commit --quiet -m "docs: publish $version from ${GITHUB_SHA::7}"

    remote="${GITHUB_SERVER_URL/:\/\//://x-access-token:$TOKEN@}/$GITHUB_REPOSITORY.git"
    git -C "$site" push --quiet "$remote" "HEAD:refs/heads/$branch"
    echo "Published $version to $branch."
}

main "$@"
