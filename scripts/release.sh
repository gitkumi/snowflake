#!/usr/bin/env bash
#
# Cut a release tag. Never hand-type a version again.
#
#   scripts/release.sh patch     # v1.104.0 -> v1.104.1
#   scripts/release.sh minor     # v1.104.0 -> v1.105.0
#   scripts/release.sh major     # v1.104.0 -> v2.0.0   (requires /v2 module path)
#   scripts/release.sh v1.105.0  # explicit, still validated
#
set -euo pipefail

BRANCH="master"
REMOTE="origin"

die() { printf '\033[31merror:\033[0m %s\n' "$*" >&2; exit 1; }
info() { printf '\033[36m==>\033[0m %s\n' "$*"; }

[ $# -eq 1 ] || die "usage: $0 <major|minor|patch|vX.Y.Z>"

# --- preflight -------------------------------------------------------------

current_branch=$(git rev-parse --abbrev-ref HEAD)
[ "$current_branch" = "$BRANCH" ] || die "on branch '$current_branch', expected '$BRANCH'"

[ -z "$(git status --porcelain)" ] || die "working tree is dirty; commit or stash first"

info "fetching $REMOTE"
git fetch --quiet "$REMOTE" "$BRANCH" --tags

local_head=$(git rev-parse HEAD)
remote_head=$(git rev-parse "$REMOTE/$BRANCH")
[ "$local_head" = "$remote_head" ] || die "HEAD ($(git rev-parse --short HEAD)) != $REMOTE/$BRANCH ($(git rev-parse --short "$REMOTE/$BRANCH")); push or pull first"

# --- work out the current and next version ---------------------------------

# -v:refname is a real semver sort. Plain `git tag -l` is lexicographic and
# will tell you v0.99.0 is newer than v0.103.0. That is how v1.104.0 happened.
latest=$(git tag --list 'v*' --sort=v:refname | tail -n 1)
[ -n "$latest" ] || die "no existing v* tags found"

[[ $latest =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]] || die "latest tag '$latest' is not vX.Y.Z"
cur_major=${BASH_REMATCH[1]}
cur_minor=${BASH_REMATCH[2]}
cur_patch=${BASH_REMATCH[3]}

case "$1" in
  major) next="v$((cur_major + 1)).0.0" ;;
  minor) next="v${cur_major}.$((cur_minor + 1)).0" ;;
  patch) next="v${cur_major}.${cur_minor}.$((cur_patch + 1))" ;;
  v*)    next="$1" ;;
  *)     die "unknown bump '$1'; use major, minor, patch, or an explicit vX.Y.Z" ;;
esac

[[ $next =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "'$next' is not a valid vX.Y.Z version"

# --- validate ---------------------------------------------------------------

if git rev-parse -q --verify "refs/tags/$next" >/dev/null; then
  die "tag $next already exists locally"
fi
if git ls-remote --exit-code --tags "$REMOTE" "refs/tags/$next" >/dev/null 2>&1; then
  die "tag $next already exists on $REMOTE"
fi

# Same rules CI enforces: legal successor, and the /vN module path for v2+.
"$(dirname "$0")/validate-tag.sh" "$next"

module=$(awk '/^module /{print $2; exit}' go.mod)

# A tag is immutable once the proxy has served it, so gate on the tests here
# rather than discovering the failure after the tag is public. SKIP_TESTS=1
# to bypass when CI has already run green on this exact commit.
if [ "${SKIP_TESTS:-}" = "1" ]; then
  info "SKIP_TESTS=1, skipping test run"
else
  info "running make test"
  make test
fi

# --- go ---------------------------------------------------------------------

echo
info "module:  $module"
info "latest:  $latest"
info "next:    $next"
info "commit:  $(git log -1 --format='%h %s')"
echo
read -r -p "Tag and push $next? [y/N] " reply
[ "$reply" = "y" ] || [ "$reply" = "Y" ] || die "aborted"

git tag -a "$next" -m "Release $next"
git push "$REMOTE" "refs/tags/$next"

info "pushed $next — https://github.com/gitkumi/snowflake/actions"
