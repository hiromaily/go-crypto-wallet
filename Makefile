###############################################################################
# go-crypto-wallet Makefile
###############################################################################
# This Makefile is organized into modules for better maintainability.
# All modules are located in the make/ directory.
#
# Module organization:
# - utils.mk:        Utility functions and cleanup
# - vars.mk:         Variable definitions and constants
# - install.mk:      Installation targets
# - codegen.mk:      Code generation targets
# - lint.mk:         Linting and code quality checks
# - build.mk:        Build-related targets
# - test.mk:         Testing targets
# - release.mk:      Release and versioning targets
# - db.mk:           Database-specific targets
# - ai.mk:           AI-related targets
# - docs.mk:         docs-ssot tool related targets
# - wallet/btc.mk:   Bitcoin-specific targets
# - wallet/bch.mk:   Bitcoin Cash-specific targets
# - wallet/eth.mk:   Ethereum-specific targets
# - wallet/xrp.mk:   XRP/Ripple-specific targets
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
include make/utils.mk

# 2. Development tools
include make/install.mk
include make/ai.mk
include make/build.mk
include make/test.mk
include make/codegen.mk
include make/lint.mk
include make/release.mk

# 3. Blockchain-specific
include make/wallet/btc.mk
include make/wallet/bch.mk
include make/wallet/eth.mk
include make/wallet/xrp.mk

# 4. Infrastructure
include make/db.mk

# 5. Wallet operations
#include make/wallet.mk
#include make/watch_op.mk
#include make/keygen_op.mk
#include make/sign_op.mk

###############################################################################
# Standard Targets (required by checkmake)
###############################################################################

.PHONY: all
all: check-build

.PHONY: clean
clean:

.PHONY: test
test: gotest
