###############################################################################
# E2E Testing - Bitcoin Transaction Patterns
###############################################################################
# Usage:
#   make btc-e2e P=<pattern>           # Run E2E test (SQLite default)
#   make btc-e2e-reset P=<pattern>     # Run with --reset flag (recommended)
#   make btc-e2e-verbose P=<pattern>   # Run with --verbose flag
#   make btc-e2e-ci P=<pattern>        # Run with --non-interactive flag
#   make btc-e2e-cleanup P=<pattern>   # Run with --cleanup flag
#
# Database Backend:
#   DB=sqlite (default) - Use local SQLite files with SQLC schemas
#                         Faster, no Docker MySQL container needed
#   DB=mysql            - Use Docker MySQL container with Atlas migrations
#
# Patterns:
#   1:  P2PKH Single-sig
#   2:  P2PKH 2-of-3 Multisig
#   3:  P2SH-P2WPKH Single-sig
#   4:  P2SH-P2WSH 2-of-3 Multisig
#   5:  P2WPKH Native SegWit Single-sig
#   6:  P2WSH Native SegWit 2-of-3 Multisig
#   7:  P2WSH Native SegWit 3-of-3 Multisig
#   8:  P2SH-P2WSH 3-of-3 Multisig
#   9:  P2TR Taproot Single-sig
#   10: P2TR MuSig2 N-of-N
#   11: P2TR Tapscript M-of-N
#
# Examples:
#   make btc-e2e-reset P=1           # SQLite (default)
#   make btc-e2e-reset P=1 DB=mysql  # MySQL
###############################################################################

# Default pattern
P ?= 1

# Default database type (sqlite or mysql)
DB ?= sqlite

# Script mapping for each pattern
E2E_SCRIPT_1  := e2e-p1-p2pkh-singlesig.sh
E2E_SCRIPT_2  := e2e-p2-p2pkh-2of3.sh
E2E_SCRIPT_3  := e2e-p3-p2sh-p2wpkh-singlesig.sh
E2E_SCRIPT_4  := e2e-p4-p2sh-p2wsh-2of3.sh
E2E_SCRIPT_5  := e2e-p5-p2wpkh-singlesig.sh
E2E_SCRIPT_6  := e2e-p6-p2wsh-2of3.sh
E2E_SCRIPT_7  := e2e-p7-p2wsh-3of3.sh
E2E_SCRIPT_8  := e2e-p8-p2sh-p2wsh-3of3.sh
E2E_SCRIPT_9  := e2e-p9-p2tr-singlesig.sh
E2E_SCRIPT_10 := e2e-p10-p2tr-musig2.sh
E2E_SCRIPT_11 := e2e-p11-p2tr-tapscript.sh

# Select script based on pattern number
E2E_SCRIPT = $(E2E_SCRIPT_$(P))
E2E_SCRIPT_PATH = ./scripts/operation/btc/e2e/$(E2E_SCRIPT)

# Database type environment variable
E2E_DB_TYPE = DB_TYPE=$(DB)

# Validate pattern number
.PHONY: _btc-e2e-validate
_btc-e2e-validate:
	@if [ -z "$(E2E_SCRIPT)" ]; then \
		echo "Error: Invalid pattern P=$(P). Valid patterns are 1-11."; \
		exit 1; \
	fi

# Run E2E test with --reset flag (recommended for fresh start)
# Note: build-all ensures binaries are up-to-date before running E2E tests
.PHONY: btc-e2e-reset
btc-e2e-reset: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --reset

# Run complete E2E test
.PHONY: btc-e2e
btc-e2e: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH)

# Run E2E test with verbose output
.PHONY: btc-e2e-verbose
btc-e2e-verbose: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --verbose

# Run E2E test in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-ci
btc-e2e-ci: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --non-interactive

# Cleanup E2E test environment (no build needed for cleanup)
.PHONY: btc-e2e-cleanup
btc-e2e-cleanup: _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --cleanup

# Show available patterns
.PHONY: btc-e2e-help
btc-e2e-help:
	@echo "Bitcoin E2E Test Patterns:"
	@echo "  P=1  : P2PKH Single-sig"
	@echo "  P=2  : P2PKH 2-of-3 Multisig"
	@echo "  P=3  : P2SH-P2WPKH Single-sig"
	@echo "  P=4  : P2SH-P2WSH 2-of-3 Multisig"
	@echo "  P=5  : P2WPKH Native SegWit Single-sig"
	@echo "  P=6  : P2WSH Native SegWit 2-of-3 Multisig"
	@echo "  P=7  : P2WSH Native SegWit 3-of-3 Multisig"
	@echo "  P=8  : P2SH-P2WSH 3-of-3 Multisig"
	@echo "  P=9  : P2TR Taproot Single-sig"
	@echo "  P=10 : P2TR MuSig2 N-of-N"
	@echo "  P=11 : P2TR Tapscript M-of-N"
	@echo ""
	@echo "Database Backend:"
	@echo "  DB=sqlite (default) - Local SQLite files with SQLC schemas"
	@echo "                        No Docker MySQL container needed (faster)"
	@echo "  DB=mysql            - Docker MySQL container with Atlas migrations"
	@echo ""
	@echo "Usage:"
	@echo "  make btc-e2e P=<pattern>              # Run E2E test (SQLite default)"
	@echo "  make btc-e2e-reset P=<pattern>        # Run with --reset flag (recommended)"
	@echo "  make btc-e2e-reset P=1 DB=mysql       # Run with MySQL"
	@echo "  make btc-e2e-verbose P=<pattern>      # Run with --verbose flag"
	@echo "  make btc-e2e-ci P=<pattern>           # Run with --non-interactive flag"
	@echo "  make btc-e2e-cleanup P=<pattern>      # Run with --cleanup flag"
