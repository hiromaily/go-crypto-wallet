###############################################################################
# E2E Tests
#
# Pattern 1: Single-sig Payment Transfer
#   DB:           sqlite (default) or mysql
#   XRP_WS_HOST:  rippled WebSocket host (default: 127.0.0.1)
#   XRP_WS_PORT:  rippled WebSocket port (default: 6006)
#
# Usage:
#   make xrp-e2e-p1-reset             # Run with SQLite (full reset)
#   make xrp-e2e-p1                   # Run without reset
#   make xrp-e2e-p1-verbose           # Run with verbose output
#   make xrp-e2e-p1-ci                # Run non-interactive (CI mode)
#   make xrp-e2e-p1-cleanup           # Stop containers and cleanup
###############################################################################

.PHONY: xrp-e2e-p1-reset
xrp-e2e-p1-reset: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p1.sh --reset

.PHONY: xrp-e2e-p1
xrp-e2e-p1: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p1.sh

.PHONY: xrp-e2e-p1-verbose
xrp-e2e-p1-verbose: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p1.sh --verbose

.PHONY: xrp-e2e-p1-ci
xrp-e2e-p1-ci: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p1.sh --non-interactive

.PHONY: xrp-e2e-p1-cleanup
xrp-e2e-p1-cleanup:
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p1.sh --cleanup

###############################################################################
# E2E Tests
#
# Pattern 2: 2-of-2 Multisig Payment Transfer
#   DB:           sqlite (default) or mysql
#   XRP_WS_HOST:  rippled WebSocket host (default: 127.0.0.1)
#   XRP_WS_PORT:  rippled WebSocket port (default: 6006)
#
# Usage:
#   make xrp-e2e-p2-reset             # Run with SQLite (full reset)
#   make xrp-e2e-p2                   # Run without reset
#   make xrp-e2e-p2-ci                # Run non-interactive (CI mode)
#   make xrp-e2e-p2-cleanup           # Stop containers and cleanup
#
# Note: Sign 1 and Sign 2 phases are currently documented (pending implementation).
#       Phases 1-7 (setup, SignerListSet, multisig TX creation) run end-to-end.
###############################################################################

.PHONY: xrp-e2e-p2-reset
xrp-e2e-p2-reset: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p2.sh --reset

.PHONY: xrp-e2e-p2
xrp-e2e-p2: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p2.sh

.PHONY: xrp-e2e-p2-ci
xrp-e2e-p2-ci: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p2.sh --non-interactive

.PHONY: xrp-e2e-p2-cleanup
xrp-e2e-p2-cleanup:
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-p2.sh --cleanup

###############################################################################
# Parallel E2E Testing - Run Multiple Patterns Concurrently
###############################################################################
# Usage:
#   make xrp-e2e-parallel                          # Run all patterns (P1) in parallel
#   make xrp-e2e-parallel PATTERNS=1               # Run specific pattern
#   make xrp-e2e-parallel MAX_PARALLEL=1           # Limit concurrent processes
#   make xrp-e2e-ci-all                            # Run all patterns in CI mode
#
# Parameters:
#   PATTERNS      - Comma-separated list or range (e.g., "1")
#                   Default: "1" (all patterns)
#   MAX_PARALLEL  - Maximum number of concurrent processes
#                   Default: 1
#   VERBOSE       - Show real-time output (true/false)
#                   Default: false
#
# Notes:
#   - Uses SQLite backend for isolated database per pattern
#   - Each pattern runs with --non-interactive flag
#   - Logs are saved to data/logs/e2e-parallel-xrp/
#   - A single shared rippled node is used across all patterns
###############################################################################

# Default parameters
XRP_PATTERNS ?= 1
XRP_MAX_PARALLEL ?= 1
XRP_VERBOSE ?= false

# Parallel runner script path
E2E_XRP_PARALLEL_SCRIPT := ./scripts/operation/xrp/e2e/e2e-parallel-runner.sh

# Build verbose flag
XRP_VERBOSE_FLAG :=
ifeq ($(XRP_VERBOSE),true)
	XRP_VERBOSE_FLAG := --verbose
endif

# Run E2E tests in parallel
# Usage: make xrp-e2e-parallel [XRP_PATTERNS=1] [XRP_MAX_PARALLEL=1] [XRP_VERBOSE=true]
.PHONY: xrp-e2e-parallel
xrp-e2e-parallel: build-all
	$(E2E_XRP_PARALLEL_SCRIPT) --patterns "$(XRP_PATTERNS)" --max-parallel "$(XRP_MAX_PARALLEL)" $(XRP_VERBOSE_FLAG)

# Run all E2E tests in parallel for CI (non-interactive mode)
# Usage: make xrp-e2e-ci-all [XRP_PATTERNS=1] [XRP_MAX_PARALLEL=1] [XRP_VERBOSE=true]
.PHONY: xrp-e2e-ci-all
xrp-e2e-ci-all: build-all
	$(E2E_XRP_PARALLEL_SCRIPT) --patterns "$(XRP_PATTERNS)" --max-parallel "$(XRP_MAX_PARALLEL)" $(XRP_VERBOSE_FLAG) --ci
