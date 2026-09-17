#!/usr/bin/env bash
#
# Regenerate THIRD_PARTY_LICENSES/ from the Go module cache.
#
# Collects every module that is actually linked into a released binary — the
# union across all target platforms, since cobra pulls in extra dependencies on
# Windows — and copies its license text out of the module cache.
#
# Usage:
#   scripts/update-third-party-licenses.sh           # regenerate
#   scripts/update-third-party-licenses.sh --check   # verify, change nothing
#
# --check exits non-zero if a linked module has no license text committed, which
# is what CI runs. Release archives must carry these texts.

set -euo pipefail

cd "$(dirname "$0")/.."

OUT_DIR="THIRD_PARTY_LICENSES"
CHECK_ONLY=false
[ "${1:-}" = "--check" ] && CHECK_ONLY=true

# Platforms GoReleaser builds for.
PLATFORMS="linux darwin windows"

# Modules linked into at least one released binary.
linked_modules() {
  for os in $PLATFORMS; do
    GOOS="$os" GOARCH=amd64 go list -deps ./... 2>/dev/null
  done |
    grep -E '^(github\.com|golang\.org|gopkg\.in)/' |
    sed -E 's#^(github\.com/[^/]+/[^/]+).*#\1#; s#^(golang\.org/x/[^/]+).*#\1#; s#^(gopkg\.in/[^/]+).*#\1#' |
    grep -v '^github\.com/torreirow/soltty$' |
    sort -u
}

missing=0
found=0

for mod in $(linked_modules); do
  dir="$(go list -m -f '{{.Dir}}' "$mod" 2>/dev/null || true)"
  if [ -z "$dir" ] || [ ! -d "$dir" ]; then
    echo "ERROR: module $mod is linked but not in the module cache" >&2
    echo "       run 'go mod download' first" >&2
    missing=$((missing + 1))
    continue
  fi

  src="$(find "$dir" -maxdepth 1 -type f \( -iname 'LICENSE*' -o -iname 'LICENCE*' -o -iname 'COPYING*' \) | head -1)"
  if [ -z "$src" ]; then
    echo "ERROR: no license file found for $mod in $dir" >&2
    missing=$((missing + 1))
    continue
  fi

  dest="$OUT_DIR/$(echo "$mod" | tr '/' '_').LICENSE"

  if $CHECK_ONLY; then
    if [ ! -f "$dest" ]; then
      echo "ERROR: missing license text: $dest ($mod)" >&2
      missing=$((missing + 1))
    elif ! cmp -s "$src" "$dest"; then
      echo "ERROR: license text out of date: $dest ($mod)" >&2
      echo "       run scripts/update-third-party-licenses.sh" >&2
      missing=$((missing + 1))
    else
      found=$((found + 1))
    fi
  else
    mkdir -p "$OUT_DIR"
    cp "$src" "$dest"
    chmod 644 "$dest"
    echo "  $mod -> $dest"
    found=$((found + 1))
  fi
done

# A NOTICE file in an Apache-2.0 dependency must be redistributed too (§4(d)).
for mod in $(linked_modules); do
  dir="$(go list -m -f '{{.Dir}}' "$mod" 2>/dev/null || true)"
  [ -n "$dir" ] || continue
  if find "$dir" -maxdepth 1 -type f -iname 'NOTICE*' | grep -q .; then
    echo "WARNING: $mod ships a NOTICE file; it must be redistributed as well" >&2
    missing=$((missing + 1))
  fi
done

if [ "$missing" -gt 0 ]; then
  echo >&2
  echo "$missing problem(s) found across $((found + missing)) linked module(s)." >&2
  exit 1
fi

if $CHECK_ONLY; then
  echo "All $found linked module license texts are present and current."
else
  echo "Wrote license texts for $found linked module(s) to $OUT_DIR/"
fi
