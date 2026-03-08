###############################################################################
# E2E Tests
#
# Pattern 1: Single-sig EIP-1559
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Usage:
#   make eth-e2e-p1-reset             # Run with Anvil + SQLite
#   make eth-e2e-p1-reset NODE_TYPE=geth DB=mysql  # Run with Geth + MySQL
###############################################################################

.PHONY: eth-e2e-p1-reset
eth-e2e-p1-reset: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p1.sh --reset

.PHONY: eth-e2e-p1
eth-e2e-p1: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p1.sh

.PHONY: eth-e2e-p1-verbose
eth-e2e-p1-verbose: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p1.sh --verbose

.PHONY: eth-e2e-p1-ci
eth-e2e-p1-ci: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p1.sh --non-interactive

.PHONY: eth-e2e-p1-cleanup
eth-e2e-p1-cleanup:
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p1.sh --cleanup

###############################################################################
# E2E Tests
#
# Pattern 2: ERC-20 HYT token transfer (EIP-1559)
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Requires the HYT contract to be deployed first:
#   make deploy-hyt
#
# Usage:
#   make eth-e2e-p2-reset             # Run with Anvil + SQLite (full reset)
#   make eth-e2e-p2                   # Run without reset
#   make eth-e2e-p2-reset HYT_CONTRACT_ADDRESS=0x...  # Override contract address
###############################################################################

.PHONY: eth-e2e-p2-reset
eth-e2e-p2-reset: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p2.sh --reset

.PHONY: eth-e2e-p2
eth-e2e-p2: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p2.sh

.PHONY: eth-e2e-p2-verbose
eth-e2e-p2-verbose: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p2.sh --verbose

.PHONY: eth-e2e-p2-ci
eth-e2e-p2-ci: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p2.sh --non-interactive

.PHONY: eth-e2e-p2-cleanup
eth-e2e-p2-cleanup:
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p2.sh --cleanup

###############################################################################
# E2E Tests
#
# Pattern 3: Safe 2-of-2 multisig payment
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Requires:
#   - Anvil running (started automatically)
#   - apps/eth-contracts/node_modules installed (bun install in eth-contracts)
#
# Usage:
#   make eth-e2e-p3-reset             # Run with Anvil + SQLite (full reset)
#   make eth-e2e-p3                   # Run without reset
#   make eth-e2e-p3-ci                # Non-interactive (CI mode)
###############################################################################

.PHONY: eth-e2e-p3-reset
eth-e2e-p3-reset: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p3.sh --reset

.PHONY: eth-e2e-p3
eth-e2e-p3: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p3.sh

.PHONY: eth-e2e-p3-verbose
eth-e2e-p3-verbose: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p3.sh --verbose

.PHONY: eth-e2e-p3-ci
eth-e2e-p3-ci: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p3.sh --non-interactive

.PHONY: eth-e2e-p3-cleanup
eth-e2e-p3-cleanup:
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p3.sh --cleanup

###############################################################################
# E2E Tests
#
# Pattern 4: MPC-TSS 2-of-3 threshold signing
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Requires:
#   - Anvil running (started automatically)
#   - tss-lib/v2 pre-computation parameters (generated automatically)
#   - Ports 9001-9003 available for MPC gRPC daemons
#
# Usage:
#   make eth-e2e-p4-reset             # Run with Anvil + SQLite (full reset)
#   make eth-e2e-p4                   # Run without reset
#   make eth-e2e-p4-ci                # Non-interactive (CI mode)
###############################################################################

.PHONY: eth-e2e-p4-reset
eth-e2e-p4-reset: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p4.sh --reset

.PHONY: eth-e2e-p4
eth-e2e-p4: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p4.sh

.PHONY: eth-e2e-p4-verbose
eth-e2e-p4-verbose: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p4.sh --verbose

.PHONY: eth-e2e-p4-ci
eth-e2e-p4-ci: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p4.sh --non-interactive

.PHONY: eth-e2e-p4-cleanup
eth-e2e-p4-cleanup:
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-p4.sh --cleanup

###############################################################################
# Parallel E2E Testing - Run Multiple Patterns Concurrently
###############################################################################
# Usage:
#   make eth-e2e-parallel                          # Run all patterns (P1-P4) in parallel
#   make eth-e2e-parallel PATTERNS=1,2             # Run specific patterns
#   make eth-e2e-parallel MAX_PARALLEL=2           # Limit concurrent processes
#   make eth-e2e-ci-all                            # Run all patterns in CI mode
#
# Parameters:
#   PATTERNS      - Comma-separated list or range (e.g., "1,2,3,4" or "1-4")
#                   Default: "1-4" (all patterns)
#   MAX_PARALLEL  - Maximum number of concurrent processes
#                   Default: 3 (all patterns run in parallel)
#   VERBOSE       - Show real-time output (true/false)
#                   Default: false
#
# Notes:
#   - Uses SQLite backend for isolated database per pattern
#   - Each pattern runs with --non-interactive flag
#   - Logs are saved to data/logs/e2e-parallel-eth/
#   - Ideal for CI/CD environments to reduce execution time
###############################################################################

# Default parameters
PATTERNS ?= 1-4
MAX_PARALLEL ?= 4
VERBOSE ?= false

# Parallel runner script path
E2E_ETH_PARALLEL_SCRIPT := ./scripts/operation/eth/e2e/e2e-parallel-runner.sh

# Build verbose flag
ETH_VERBOSE_FLAG :=
ifeq ($(VERBOSE),true)
	ETH_VERBOSE_FLAG := --verbose
endif

# Run E2E tests in parallel
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make eth-e2e-parallel [PATTERNS=1-3] [MAX_PARALLEL=3] [VERBOSE=true]
#  e.g. make eth-e2e-parallel PATTERNS=1-3 MAX_PARALLEL=3
.PHONY: eth-e2e-parallel
eth-e2e-parallel: build-all
	NODE_TYPE="$(NODE_TYPE)" $(E2E_ETH_PARALLEL_SCRIPT) --patterns "$(PATTERNS)" --max-parallel "$(MAX_PARALLEL)" $(ETH_VERBOSE_FLAG)

# Run all E2E tests in parallel for CI (non-interactive mode)
# Note: build-all uses incremental build - only rebuilds when Go sources change
# Usage: make eth-e2e-ci-all [PATTERNS=1-2] [MAX_PARALLEL=2] [VERBOSE=true]
.PHONY: eth-e2e-ci-all
eth-e2e-ci-all: build-all
	NODE_TYPE="$(NODE_TYPE)" $(E2E_ETH_PARALLEL_SCRIPT) --patterns "$(PATTERNS)" --max-parallel "$(MAX_PARALLEL)" $(ETH_VERBOSE_FLAG) --ci
