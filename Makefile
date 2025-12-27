###############################################################################
# go-crypto-wallet Makefile
###############################################################################
# This Makefile is organized into modules for better maintainability.
# All modules are located in the make/ directory.
#
# Module organization:
# - vars.mk:         Variable definitions and constants
# - install.mk:      Installation targets
# - build.mk:        Build-related targets
# - test.mk:         Testing targets
# - btc.mk:          Bitcoin-specific targets
# - eth.mk:          Ethereum-specific targets
# - xrp.mk:          XRP/Ripple-specific targets
# - codegen.mk:      Code generation targets
# - lint.mk:         Linting and code quality checks
# - wallet.mk:       Wallet management operations
# - utils.mk:        Utility functions and cleanup
# - watch_op.mk:     Watch wallet operations
# - keygen_op.mk:    Keygen wallet operations
# - sign_op.mk:      Sign wallet operations
###############################################################################

# Include modules in logical order based on dependencies
# Dependencies:
# - vars.mk: Base variables (no dependencies, must be first)
# - install.mk: Depends on vars.mk (modVer, currentVer, GOLANGCI_VERSION)
# - eth.mk: Depends on vars.mk (GETH_VERSION, LODESTAR_VERSION, etc.)
# - utils.mk: Depends on vars.mk (timestamp) and db.mk (rm-db-volumes target)
#
# Logical grouping:
# 1. Base: Variables and configuration
# 2. Development tools: Install, build, test, codegen, lint
# 3. Blockchain-specific: BTC, ETH, XRP
# 4. Infrastructure: Database
# 5. Wallet operations: Wallet management, utilities, watch/keygen/sign operations

# 1. Base: Variables (must be first)
include make/vars.mk

# 2. Development tools
include make/install.mk
include make/build.mk
include make/test.mk
include make/codegen.mk
include make/lint.mk

# 3. Blockchain-specific
include make/btc.mk
include make/eth.mk
include make/xrp.mk

# 4. Infrastructure
include make/db.mk

# 5. Wallet operations
include make/wallet.mk
include make/utils.mk
include make/watch_op.mk
include make/keygen_op.mk
include make/sign_op.mk

###############################################################################
# Standard Targets (required by checkmake)
###############################################################################

.PHONY: all
all: check-build

.PHONY: clean
clean:

.PHONY: test
test: gotest
