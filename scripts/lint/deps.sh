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

warn() {
  echo "WARNING: $*" >&2
  VIOLATIONS=$((VIOLATIONS + 1))
}

is_generated() {
  grep -qE "(Code generated|DO NOT EDIT)" "$1"
}

# ---------------------------------------------------------------------------
# 1. usecase must not import infrastructure directly
# ---------------------------------------------------------------------------
echo "=== [1] Checking usecase -> infrastructure direct imports ==="

HEADER_SHOWN=0
while IFS= read -r file; do
  is_generated "$file" && continue
  if grep -qE "\"${MODULE}/internal/infrastructure" "$file"; then
    if [ "$HEADER_SHOWN" -eq 0 ]; then
      echo "WARNING: Direct infrastructure import in usecase layer:" >&2
      HEADER_SHOWN=1
    fi
    VIOLATIONS=$((VIOLATIONS + 1))
    echo "  $file"
    grep -nE "\"${MODULE}/internal/infrastructure" "$file" | while IFS= read -r line; do
      echo "       $line"
    done
  fi
done < <(find internal/application/usecase -name "*.go" ! -name "*_test.go")

# ---------------------------------------------------------------------------
# 2. infrastructure must not define exported interfaces
#    (exported = starts with uppercase; they belong in application/ports)
# ---------------------------------------------------------------------------
echo ""
echo "=== [2] Checking for exported interface definitions inside infrastructure ==="

HEADER_SHOWN=0
while IFS= read -r file; do
  is_generated "$file" && continue
  # Match exported interfaces: type UppercaseName interface {
  if grep -qE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{" "$file"; then
    if [ "$HEADER_SHOWN" -eq 0 ]; then
      echo "WARNING: Exported interface defined in infrastructure layer (move to application/ports):" >&2
      HEADER_SHOWN=1
    fi
    VIOLATIONS=$((VIOLATIONS + 1))
    echo "  $file"
    grep -nE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{" "$file" | while IFS= read -r line; do
      echo "       $line"
    done
  fi
done < <(find internal/infrastructure -name "*.go" ! -name "*_test.go")

# ---------------------------------------------------------------------------
# 3. pkg must not import internal
# ---------------------------------------------------------------------------
echo ""
echo "=== [3] Checking pkg -> internal dependencies ==="

HEADER_SHOWN=0
while IFS= read -r file; do
  is_generated "$file" && continue
  if grep -qE "\"${MODULE}/internal" "$file"; then
    if [ "$HEADER_SHOWN" -eq 0 ]; then
      echo "WARNING: pkg imports internal package:" >&2
      HEADER_SHOWN=1
    fi
    VIOLATIONS=$((VIOLATIONS + 1))
    echo "  $file"
    grep -nE "\"${MODULE}/internal" "$file" | while IFS= read -r line; do
      echo "       $line"
    done
  fi
done < <(find pkg -name "*.go" ! -name "*_test.go")

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
