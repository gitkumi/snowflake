#!/usr/bin/env bash
#
# Validate a release tag against the previous one. Used by scripts/release.sh
# before tagging, and by the Release workflow before publishing.
#
#   scripts/validate-tag.sh v1.105.0
#
# Requires full tag history (actions/checkout needs fetch-depth: 0).
#
set -euo pipefail

die() { printf '\033[31mtag validation failed:\033[0m %s\n' "$*" >&2; exit 1; }

[ $# -eq 1 ] || { echo "usage: $0 <vX.Y.Z>" >&2; exit 2; }
tag="$1"

[[ $tag =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]] || die "'$tag' is not a vX.Y.Z version"
major=${BASH_REMATCH[1]}

# --- Go module major-version suffix ----------------------------------------
#
# Go requires v2+ module paths to end in /vN. Tagging v2.0.0 on a module named
# ".../snowflake" produces a tag the proxy refuses to serve: the release looks
# green locally but 'go install ...@latest' silently keeps resolving to v1.x.

module=$(awk '/^module /{print $2; exit}' go.mod)
if [ "$major" -ge 2 ]; then
  [ "${module##*/}" = "v${major}" ] || die \
"go.mod declares 'module $module' but the tag is $tag.

Go requires v2+ modules to end in /v${major}. To release v${major}:
  1. set 'module $module/v${major}' in go.mod
  2. rewrite internal imports to $module/v${major}/...
  3. commit, push, then tag"
else
  case "$module" in
    */v[0-9]*) die "go.mod module path '$module' carries a major-version suffix but the tag is $tag" ;;
  esac
fi

# --- legal successor --------------------------------------------------------
#
# The greatest tag strictly below this one. -v:refname is a real semver sort;
# plain lexicographic sorting ranks v0.99.0 above v0.103.0.

prev=$(git tag --list 'v*' --sort=v:refname | grep -vFx "$tag" | awk -v t="$tag" '
  { tags[NR] = $0 }
  END {
    for (i = NR; i >= 1; i--) {
      split(substr(tags[i], 2), a, ".")
      split(substr(t, 2), b, ".")
      if (a[1] < b[1] || (a[1] == b[1] && (a[2] < b[2] || (a[2] == b[2] && a[3] < b[3])))) {
        print tags[i]; exit
      }
    }
  }')

if [ -z "$prev" ]; then
  echo "no earlier tag than $tag; skipping successor check"
  exit 0
fi

[[ $prev =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]] || { echo "previous tag '$prev' is not vX.Y.Z; skipping successor check"; exit 0; }
p_major=${BASH_REMATCH[1]}
p_minor=${BASH_REMATCH[2]}
p_patch=${BASH_REMATCH[3]}

legal="v${p_major}.${p_minor}.$((p_patch + 1)) v${p_major}.$((p_minor + 1)).0 v$((p_major + 1)).0.0"
case " $legal " in
  *" $tag "*) ;;
  *) die "$tag is not a legal successor to $prev.

A major bump resets minor and patch; a minor bump resets patch. Valid next
versions after $prev are:
$(printf '%s\n' $legal | sed 's/^/  /')" ;;
esac

echo "ok: $tag is a valid successor to $prev (module $module)"
