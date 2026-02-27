#!/usr/bin/env bash
# visibility.sh - Encapsulation and visibility convention checker
#
# Checks:
#   1. Exported struct definitions inside internal/infrastructure should be unexported.
#      Infrastructure structs are used only through application/ports interfaces,
#      so their concrete types must not be visible outside the DI layer.
#      NOTE: Structs embedded by a child package (e.g. Bitcoin embedded by BitcoinCash)
#      are a known exception and must remain exported.
#
#   2. Interface definitions in pkg/ must reside in interface.go or {package_name}.go.
#      Each pkg package may contain both abstract and concrete types; keeping all
#      interface declarations in a single, predictable file makes the public API obvious.
#
# Auto-generated files (containing "Code generated" or "DO NOT EDIT") are skipped.

set -euo pipefail

VIOLATIONS=0

is_generated() {
  grep -qE "(Code generated|DO NOT EDIT)" "$1"
}

# ---------------------------------------------------------------------------
# 1. infrastructure must use unexported (private) structs
#    Exported struct types in infrastructure can leak concrete types to callers
#    who are supposed to depend only on application/ports interfaces.
# ---------------------------------------------------------------------------
echo ""
echo "=== [1] Checking for exported struct definitions inside infrastructure ==="

header_shown=0
while IFS= read -r file; do
  is_generated "$file" && continue
  if grep -qE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+struct[[:space:]]*\{" "$file"; then
    if [ "$header_shown" -eq 0 ]; then
      echo "WARNING: Exported struct in infrastructure (should be unexported; callers must use ports interfaces):" >&2
      echo ""
      header_shown=1
    fi
    VIOLATIONS=$((VIOLATIONS + 1))
    echo "  $file"
    grep -nE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+struct[[:space:]]*\{" "$file" | sed 's/^/       /'
    echo ""
  fi
done < <(find "internal/infrastructure" -name "*.go" ! -name "*_test.go" ! -path "*/testutil/*")

# ---------------------------------------------------------------------------
# 2. Interface declarations in pkg/ must be in interface.go or {pkg_name}.go
#    Concrete implementations may live in any other file.
# ---------------------------------------------------------------------------
echo "=== [2] Checking interface file placement in pkg/ ==="

header_shown=0
while IFS= read -r file; do
  is_generated "$file" && continue
  if grep -qE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{" "$file"; then
    filename=$(basename "$file")
    pkg_dir=$(basename "$(dirname "$file")")
    if [[ "$filename" != "interface.go" && "$filename" != "${pkg_dir}.go" ]]; then
      if [ "$header_shown" -eq 0 ]; then
        echo "WARNING: Interface defined in unexpected file in pkg/ (move to interface.go or ${pkg_dir}.go):" >&2
        echo ""
        header_shown=1
      fi
      VIOLATIONS=$((VIOLATIONS + 1))
      echo "  $file"
      grep -nE "^[[:space:]]*type[[:space:]]+[A-Z][A-Za-z0-9_]*[[:space:]]+interface[[:space:]]*\{" "$file" | sed 's/^/       /'
      echo ""
    fi
  fi
done < <(find "pkg" -name "*.go" ! -name "*_test.go")

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
if [ "$VIOLATIONS" -eq 0 ]; then
  echo "No visibility violations found."
  exit 0
else
  echo "Found $VIOLATIONS violation(s). See warnings above." >&2
  exit 1
fi
