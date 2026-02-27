#!/usr/bin/env bash
# deps.sh - Clean Architecture dependency violation checker
#
# Checks:
#   1. application/usecase -> infrastructure (direct import violation)
#   2. Exported interface declarations inside infrastructure (should live in application/ports)
#   3. pkg -> internal (forbidden reverse dependency)
#
# Auto-generated files (containing "Code generated" or "DO NOT EDIT") are skipped.

set -euo pipefail

MODULE="github.com/hiromaily/go-crypto-wallet"
VIOLATIONS=0

is_generated() {
  grep -qE "(Code generated|DO NOT EDIT)" "$1"
}

# check_violations <description> <find_path> <grep_pattern> <warning_message>
check_violations() {
  local description="$1"
  local find_path="$2"
  local grep_pattern="$3"
  local warning_message="$4"

  echo ""
  echo "=== ${description} ==="

  local header_shown=0
  while IFS= read -r file; do
    is_generated "$file" && continue
    if grep -qE "${grep_pattern}" "$file"; then
      if [ "$header_shown" -eq 0 ]; then
        echo "WARNING: ${warning_message}" >&2
        header_shown=1
      fi
      VIOLATIONS=$((VIOLATIONS + 1))
      echo "  $file"
      grep -nE "${grep_pattern}" "$file" | sed 's/^/       /'
    fi
  done < <(find "${find_path}" -name "*.go" ! -name "*_test.go")
}

# ---------------------------------------------------------------------------
# 1. usecase must not import infrastructure directly
# ---------------------------------------------------------------------------
check_violations \
  "[1] Checking usecase -> infrastructure direct imports" \
  "internal/application/usecase" \
  "\"${MODULE}/internal/infrastructure" \
  "Direct infrastructure import in usecase layer:"

# ---------------------------------------------------------------------------
# 2. infrastructure must not define exported interfaces
#    (exported = starts with uppercase; they belong in application/ports)
# ---------------------------------------------------------------------------
check_violations \
  "[2] Checking for exported interface definitions inside infrastructure" \
  "internal/infrastructure" \
  "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{" \
  "Exported interface defined in infrastructure layer (move to application/ports):"

# ---------------------------------------------------------------------------
# 3. pkg must not import internal
# ---------------------------------------------------------------------------
check_violations \
  "[3] Checking pkg -> internal dependencies" \
  "pkg" \
  "\"${MODULE}/internal" \
  "pkg imports internal package:"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
if [ "$VIOLATIONS" -eq 0 ]; then
  echo "No dependency violations found."
  exit 0
else
  echo "Found $VIOLATIONS violation(s). See warnings above." >&2
  exit 1
fi
