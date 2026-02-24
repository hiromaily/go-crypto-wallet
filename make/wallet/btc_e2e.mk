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
#                         Faster, no Docker container needed
#   DB=postgres         - Use Docker PostgreSQL container with Atlas migrations
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
#   make btc-e2e-reset P=1              # SQLite (default)
#   make btc-e2e-reset P=1 DB=postgres  # PostgreSQL
#   make btc-e2e-reset P=1 DB=mysql     # MySQL
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
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e-reset P=<1-11> [DB=sqlite|postgres|mysql]
.PHONY: btc-e2e-reset
btc-e2e-reset: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --reset

# Run complete E2E test
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e P=<1-11> [DB=sqlite|postgres|mysql]
.PHONY: btc-e2e
btc-e2e: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH)

# Run E2E test with verbose output
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e-verbose P=<1-11> [DB=sqlite|postgres|mysql]
.PHONY: btc-e2e-verbose
btc-e2e-verbose: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --verbose

# Run E2E test in non-interactive mode (for CI/CD)
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e-ci P=<1-11> [DB=sqlite|postgres|mysql]
.PHONY: btc-e2e-ci
btc-e2e-ci: build-all _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --non-interactive

# Cleanup E2E test environment (no build needed for cleanup)
# Usage: make btc-e2e-cleanup P=<1-11> [DB=sqlite|postgres|mysql]
.PHONY: btc-e2e-cleanup
btc-e2e-cleanup: _btc-e2e-validate
	$(E2E_DB_TYPE) $(E2E_SCRIPT_PATH) --cleanup

# Show available patterns
# Usage: make btc-e2e-help
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
	@echo "                        No Docker container needed (faster)"
	@echo "  DB=postgres         - Docker PostgreSQL container with Atlas migrations"
	@echo "  DB=mysql            - Docker MySQL container with Atlas migrations"
	@echo ""
	@echo "Usage:"
	@echo "  make btc-e2e P=<pattern>              # Run E2E test (SQLite default)"
	@echo "  make btc-e2e-reset P=<pattern>        # Run with --reset flag (recommended)"
	@echo "  make btc-e2e-reset P=1 DB=postgres    # Run with PostgreSQL"
	@echo "  make btc-e2e-reset P=1 DB=mysql       # Run with MySQL"
	@echo "  make btc-e2e-verbose P=<pattern>      # Run with --verbose flag"
	@echo "  make btc-e2e-ci P=<pattern>           # Run with --non-interactive flag"
	@echo "  make btc-e2e-cleanup P=<pattern>      # Run with --cleanup flag"

###############################################################################
# Parallel E2E Testing - Run Multiple Patterns Concurrently
###############################################################################
# Usage:
#   make btc-e2e-parallel                          # Run all patterns (P1-P11) in parallel
#   make btc-e2e-parallel PATTERNS=1,2,3           # Run specific patterns
#   make btc-e2e-parallel PATTERNS=1-5             # Run pattern range
#   make btc-e2e-parallel MAX_PARALLEL=3           # Limit concurrent processes
#   make btc-e2e-ci-all                            # Run all patterns in CI mode
#
# Parameters:
#   PATTERNS      - Comma-separated list or range (e.g., "1,2,3" or "1-11")
#                   Default: "1-11" (all patterns)
#   MAX_PARALLEL  - Maximum number of concurrent processes
#                   Default: 11 (all patterns run in parallel)
#   VERBOSE       - Show real-time output (true/false)
#                   Default: false
#
# Notes:
#   - Uses SQLite backend for isolated database per pattern
#   - Each pattern runs with --reset --non-interactive flags
#   - Logs are saved to data/logs/e2e-parallel/
#   - Ideal for CI/CD environments to reduce execution time
#
# Examples:
#   make btc-e2e-parallel                          # All patterns, full parallelism
#   make btc-e2e-parallel PATTERNS=1-5 MAX_PARALLEL=3
#   make btc-e2e-ci-all                            # CI mode (all patterns)
###############################################################################

# Default parameters
PATTERNS ?= 1-11
MAX_PARALLEL ?= 11
VERBOSE ?= false

# Parallel runner script path
E2E_PARALLEL_SCRIPT := ./scripts/operation/btc/e2e/e2e-parallel-runner.sh

# Build verbose flag
VERBOSE_FLAG :=
ifeq ($(VERBOSE),true)
	VERBOSE_FLAG := --verbose
endif

# Run E2E tests in parallel
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e-parallel [PATTERNS=1-11] [MAX_PARALLEL=11] [VERBOSE=true]
# e.g. make btc-e2e-parallel PATTERNS=1-9 MAX_PARALLEL=9
.PHONY: btc-e2e-parallel
btc-e2e-parallel: build-all
	$(E2E_PARALLEL_SCRIPT) --patterns $(PATTERNS) --max-parallel $(MAX_PARALLEL) $(VERBOSE_FLAG)

# Run all E2E tests in parallel for CI (non-interactive mode)
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make btc-e2e-ci-all [PATTERNS=1-11] [MAX_PARALLEL=11] [VERBOSE=true]
.PHONY: btc-e2e-ci-all
btc-e2e-ci-all: build-all
	$(E2E_PARALLEL_SCRIPT) --patterns $(PATTERNS) --max-parallel $(MAX_PARALLEL) $(VERBOSE_FLAG) --ci
