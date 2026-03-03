#!/usr/bin/env bash
# eth-rpc-layer.sh - ETH RPC layer boundary checker
#
# Rule: CallContext must only be used in pkg/chains/eth/rpc/.
# The infrastructure layer (internal/infrastructure/api/) must NOT call
# CallContext directly; it must delegate to the pkg layer via pkgrpc.
#
# See: .claude/rules/chains/eth/pkg-or-internal.md
#
# Auto-generated files (containing "Code generated" or "DO NOT EDIT") are skipped.
# Commented-out lines are also skipped.

set -euo pipefail

VIOLATIONS=0

is_generated() {
	grep -qE "(Code generated|DO NOT EDIT)" "$1"
}

# ---------------------------------------------------------------------------
# 1. infrastructure/api must not call CallContext directly
#    All CallContext calls belong in pkg/chains/eth/rpc/
# ---------------------------------------------------------------------------
echo ""
echo "=== [1] Checking for CallContext usage in internal/infrastructure/api ==="

header_shown=0
while IFS= read -r file; do
	is_generated "$file" && continue
	# grep for CallContext actual method calls, ignoring:
	#   - comment lines (lines starting with optional whitespace then //)
	#   - string literals (a " before .CallContext means it's inside a quoted string)
	if grep -qE "^[^\"]*\.CallContext\(" "$file" && grep -E "^[^\"]*\.CallContext\(" "$file" | grep -qvE "^[[:space:]]*//"; then
		if [ "$header_shown" -eq 0 ]; then
			echo "WARNING: CallContext called directly in infrastructure layer (move to pkg/chains/eth/rpc/):" >&2
			echo ""
			header_shown=1
		fi
		VIOLATIONS=$((VIOLATIONS + 1))
		echo "  $file"
		grep -nE "^[^\"]*\.CallContext\(" "$file" | grep -vE "^[0-9]+:[[:space:]]*///" | sed 's/^/       /'
		echo ""
	fi
done < <(find "internal/infrastructure/api" -name "*.go" ! -name "*_test.go")

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
if [ "$VIOLATIONS" -eq 0 ]; then
	echo "No ETH RPC layer violations found."
	exit 0
else
	echo "Found $VIOLATIONS violation(s). See warnings above." >&2
	exit 1
fi
