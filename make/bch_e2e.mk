###############################################################################
# E2E Testing - Bitcoin Cash Transaction Patterns
###############################################################################
# Usage:
#   make bch-e2e P=<pattern>           # Run E2E test
#   make bch-e2e-reset P=<pattern>     # Run with --reset flag (recommended)
#   make bch-e2e-verbose P=<pattern>   # Run with --verbose flag
#   make bch-e2e-ci P=<pattern>        # Run with --non-interactive flag
#   make bch-e2e-cleanup P=<pattern>   # Run with --cleanup flag
#
# Patterns:
#   1: P2PKH Single-sig (CashAddr)
#   2: P2SH 2-of-3 Multisig (CashAddr)
#   3: P2SH 3-of-3 Multisig (CashAddr)
#
# Note: BCH does NOT support SegWit or Taproot patterns
###############################################################################

# Default pattern
P ?= 1

# Script mapping for each pattern
BCH_E2E_SCRIPT_1 := e2e-p1-p2pkh-singlesig.sh
BCH_E2E_SCRIPT_2 := e2e-p2-p2sh-2of3.sh
BCH_E2E_SCRIPT_3 := e2e-p3-p2sh-3of3.sh

# Select script based on pattern number
BCH_E2E_SCRIPT = $(BCH_E2E_SCRIPT_$(P))
BCH_E2E_SCRIPT_PATH = ./scripts/operation/bch/e2e/$(BCH_E2E_SCRIPT)

# Validate pattern number
.PHONY: _bch-e2e-validate
_bch-e2e-validate:
	@if [ -z "$(BCH_E2E_SCRIPT)" ]; then \
		echo "Error: Invalid pattern P=$(P). Valid patterns are 1-3."; \
		echo "  P=1: P2PKH Single-sig"; \
		echo "  P=2: P2SH 2-of-3 Multisig"; \
		echo "  P=3: P2SH 3-of-3 Multisig"; \
		exit 1; \
	fi

# Run E2E test with --reset flag (recommended for fresh start)
.PHONY: bch-e2e-reset
bch-e2e-reset: _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --reset

# Run complete E2E test
.PHONY: bch-e2e
bch-e2e: _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH)

# Run E2E test with verbose output
.PHONY: bch-e2e-verbose
bch-e2e-verbose: _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --verbose

# Run E2E test in non-interactive mode (for CI/CD)
.PHONY: bch-e2e-ci
bch-e2e-ci: _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --non-interactive

# Cleanup E2E test environment
.PHONY: bch-e2e-cleanup
bch-e2e-cleanup: _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --cleanup

# Show available patterns
.PHONY: bch-e2e-help
bch-e2e-help:
	@echo "Bitcoin Cash E2E Test Patterns:"
	@echo "  P=1: P2PKH Single-sig (CashAddr)"
	@echo "  P=2: P2SH 2-of-3 Multisig (CashAddr)"
	@echo "  P=3: P2SH 3-of-3 Multisig (CashAddr)"
	@echo ""
	@echo "Note: BCH does NOT support SegWit or Taproot patterns"
	@echo ""
	@echo "Usage:"
	@echo "  make bch-e2e P=<pattern>         # Run E2E test"
	@echo "  make bch-e2e-reset P=<pattern>   # Run with --reset flag (recommended)"
	@echo "  make bch-e2e-verbose P=<pattern> # Run with --verbose flag"
	@echo "  make bch-e2e-ci P=<pattern>      # Run with --non-interactive flag"
	@echo "  make bch-e2e-cleanup P=<pattern> # Run with --cleanup flag"
