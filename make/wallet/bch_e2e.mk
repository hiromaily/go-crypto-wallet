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
# Note: build-all uses incremental build - only rebuilds when Go sources change
.PHONY: bch-e2e-reset
bch-e2e-reset: build-all _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --reset

# Run complete E2E test
# Note: build-all uses incremental build - only rebuilds when Go sources change
.PHONY: bch-e2e
bch-e2e: build-all _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH)

# Run E2E test with verbose output
# Note: build-all uses incremental build - only rebuilds when Go sources change
.PHONY: bch-e2e-verbose
bch-e2e-verbose: build-all _bch-e2e-validate
	$(BCH_E2E_SCRIPT_PATH) --verbose

# Run E2E test in non-interactive mode (for CI/CD)
# Note: build-all uses incremental build - only rebuilds when Go sources change
.PHONY: bch-e2e-ci
bch-e2e-ci: build-all _bch-e2e-validate
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

###############################################################################
# Parallel E2E Testing - Run Multiple Patterns Concurrently
###############################################################################
# Usage:
#   make bch-e2e-parallel                          # Run all patterns (P1-P3) in parallel
#   make bch-e2e-parallel PATTERNS=1,2             # Run specific patterns
#   make bch-e2e-parallel PATTERNS=1-3             # Run pattern range
#   make bch-e2e-parallel MAX_PARALLEL=2           # Limit concurrent processes
#   make bch-e2e-ci-all                            # Run all patterns in CI mode
#
# Parameters:
#   PATTERNS      - Comma-separated list or range (e.g., "1,2,3" or "1-3")
#                   Default: "1-3" (all patterns)
#   MAX_PARALLEL  - Maximum number of concurrent processes
#                   Default: 3 (all patterns run in parallel)
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
#   make bch-e2e-parallel                          # All patterns, full parallelism
#   make bch-e2e-parallel PATTERNS=1,2 MAX_PARALLEL=2
#   make bch-e2e-ci-all                            # CI mode (all patterns)
###############################################################################

# Default parameters
PATTERNS ?= 1-3
MAX_PARALLEL ?= 3
VERBOSE ?= false

# Parallel runner script path
BCH_E2E_PARALLEL_SCRIPT := ./scripts/operation/bch/e2e/e2e-parallel-runner.sh

# Build verbose flag
BCH_VERBOSE_FLAG :=
ifeq ($(VERBOSE),true)
	BCH_VERBOSE_FLAG := --verbose
endif

# Run E2E tests in parallel
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make bch-e2e-parallel [PATTERNS=1-3] [MAX_PARALLEL=3] [VERBOSE=true]
# e.g. make bch-e2e-parallel PATTERNS=1-3 MAX_PARALLEL=3
.PHONY: bch-e2e-parallel
bch-e2e-parallel: build-all
	$(BCH_E2E_PARALLEL_SCRIPT) --patterns $(PATTERNS) --max-parallel $(MAX_PARALLEL) $(BCH_VERBOSE_FLAG)

# Run all E2E tests in parallel for CI (non-interactive mode)
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make bch-e2e-ci-all [PATTERNS=1-3] [MAX_PARALLEL=3] [VERBOSE=true]
.PHONY: bch-e2e-ci-all
bch-e2e-ci-all: build-all
	$(BCH_E2E_PARALLEL_SCRIPT) --patterns $(PATTERNS) --max-parallel $(MAX_PARALLEL) $(BCH_VERBOSE_FLAG) --ci

# Check List
# For PostgreSQL
# make bch-e2e-reset P=1 DB=postgres
# make bch-e2e-reset P=2 DB=postgres
# make bch-e2e-reset P=3 DB=postgres
# For SQLite
# make bch-e2e-parallel PATTERNS=1-3 MAX_PARALLEL=3
