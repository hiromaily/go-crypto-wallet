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
# - docker.mk:       Docker and compose operations
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

# Include modules in logical order
# Variables must come first as they're used by other modules
include make/vars.mk
include make/install.mk
include make/build.mk
include make/test.mk
include make/docker.mk
include make/btc.mk
include make/eth.mk
include make/xrp.mk
include make/codegen.mk
include make/lint.mk
include make/wallet.mk
include make/utils.mk
include make/watch_op.mk
include make/keygen_op.mk
include make/sign_op.mk
