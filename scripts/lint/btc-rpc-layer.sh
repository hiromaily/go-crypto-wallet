#!/usr/bin/env bash
# btc-rpc-layer.sh - BTC RPC layer boundary checker
#
# Rule: client.RawRequest must only be used in pkg/chains/btc/rpc/.
# The infrastructure layer (internal/infrastructure/api/) must NOT call
# RawRequest directly; it must delegate to the pkg layer via pkgrpc.
#
# See: .claude/rules/btc/pkg-or-internal.md
#
# Auto-generated files (containing "Code generated" or "DO NOT EDIT") are skipped.
# Commented-out lines are also skipped.

set -euo pipefail

VIOLATIONS=0

is_generated() {
	grep -qE "(Code generated|DO NOT EDIT)" "$1"
}

# ---------------------------------------------------------------------------
# 1. infrastructure/api must not call client.RawRequest directly
#    All RawRequest calls belong in pkg/chains/btc/rpc/
# ---------------------------------------------------------------------------
echo ""
echo "=== [1] Checking for RawRequest usage in internal/infrastructure/api ==="

header_shown=0
while IFS= read -r file; do
	is_generated "$file" && continue
	# grep for RawRequest actual method calls, ignoring:
	#   - comment lines (lines starting with optional whitespace then //)
	#   - string literals (a " before .RawRequest means it's inside a quoted string)
	if grep -qE "^[^\"]*\.RawRequest\(" "$file" && grep -E "^[^\"]*\.RawRequest\(" "$file" | grep -qvE "^[[:space:]]*//"; then
		if [ "$header_shown" -eq 0 ]; then
			echo "WARNING: RawRequest called directly in infrastructure layer (move to pkg/chains/btc/rpc/):" >&2
			echo ""
			header_shown=1
		fi
		VIOLATIONS=$((VIOLATIONS + 1))
		echo "  $file"
		grep -nE "^[^\"]*\.RawRequest\(" "$file" | grep -vE "^[0-9]+:[[:space:]]*//" | sed 's/^/       /'
		echo ""
	fi
done < <(find "internal/infrastructure/api" -name "*.go" ! -name "*_test.go")

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
if [ "$VIOLATIONS" -eq 0 ]; then
	echo "No BTC RPC layer violations found."
	exit 0
else
	echo "Found $VIOLATIONS violation(s). See warnings above." >&2
	exit 1
fi
