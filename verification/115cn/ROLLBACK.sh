#!/usr/bin/env bash
set -euo pipefail
if [[ $# -ne 1 ]]; then
  echo "usage: $0 TARGET_DIRECTORY" >&2
  exit 2
fi
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
TARGET_ROOT="$(cd "$1" && pwd)"
if [[ "$TARGET_ROOT" == "$SOURCE_ROOT" ]]; then
  echo "refusing to roll back the active source checkout" >&2
  exit 3
fi
if command -v git.exe >/dev/null 2>&1 && command -v wslpath >/dev/null 2>&1; then
  TARGET_ARG="$(wslpath -w "$TARGET_ROOT")"
  PATCH_ARG="$(wslpath -w "$SCRIPT_DIR/DIFF_FILE.patch")"
  GIT_CMD=(git.exe -C "$TARGET_ARG")
else
  PATCH_ARG="$SCRIPT_DIR/DIFF_FILE.patch"
  GIT_CMD=(git -C "$TARGET_ROOT")
fi
"${GIT_CMD[@]}" apply --no-index -R --ignore-space-change --check "$PATCH_ARG"
"${GIT_CMD[@]}" apply --no-index -R --ignore-space-change "$PATCH_ARG"
echo "ROLLBACK_OK target=$TARGET_ROOT restored=100CN-baseline"