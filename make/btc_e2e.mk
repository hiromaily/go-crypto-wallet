###############################################################################
# E2E Testing - Pattern 1: P2PKH Single-sig
###############################################################################
# Run Bitcoin E2E workflow Pattern 1 from completely fresh state (recommended)
.PHONY: btc-e2e-p1-reset
btc-e2e-p1-reset:
	./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 1
.PHONY: btc-e2e-p1
btc-e2e-p1:
	./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh

# Run Bitcoin E2E workflow Pattern 1 with verbose output
.PHONY: btc-e2e-p1-verbose
btc-e2e-p1-verbose:
	./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --verbose

# Run Bitcoin E2E workflow Pattern 1 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p1-ci
btc-e2e-p1-ci:
	./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 1
.PHONY: btc-e2e-p1-cleanup
btc-e2e-p1-cleanup:
	./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --cleanup

###############################################################################
# E2E Testing - Pattern 2: P2PKH 2-of-3 Multisig
###############################################################################
# Run Bitcoin E2E workflow Pattern 2 from completely fresh state (recommended)
.PHONY: btc-e2e-p2-reset
btc-e2e-p2-reset:
	./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 2
.PHONY: btc-e2e-p2
btc-e2e-p2:
	./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh

# Run Bitcoin E2E workflow Pattern 2 with verbose output
.PHONY: btc-e2e-p2-verbose
btc-e2e-p2-verbose:
	./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh --verbose

# Run Bitcoin E2E workflow Pattern 2 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p2-ci
btc-e2e-p2-ci:
	./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 2
.PHONY: btc-e2e-p2-cleanup
btc-e2e-p2-cleanup:
	./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh --cleanup

###############################################################################
# E2E Testing - Pattern 3: P2SH-P2WPKH Single-sig
###############################################################################
# Run Bitcoin E2E workflow Pattern 3 from completely fresh state (recommended)
.PHONY: btc-e2e-p3-reset
btc-e2e-p3-reset:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 3
.PHONY: btc-e2e-p3
btc-e2e-p3:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh

# Run Bitcoin E2E workflow Pattern 3 with verbose output
.PHONY: btc-e2e-p3-verbose
btc-e2e-p3-verbose:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --verbose

# Run Bitcoin E2E workflow Pattern 3 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p3-ci
btc-e2e-p3-ci:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 3
.PHONY: btc-e2e-p3-cleanup
btc-e2e-p3-cleanup:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --cleanup

###############################################################################
# E2E Testing - Pattern 8: P2SH-P2WSH 3-of-3 Multisig
###############################################################################
# Run Bitcoin E2E workflow Pattern 8 from completely fresh state (recommended)
.PHONY: btc-e2e-p8-reset
btc-e2e-p8-reset:
	./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 8
.PHONY: btc-e2e-p8
btc-e2e-p8:
	./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh

# Run Bitcoin E2E workflow Pattern 8 with verbose output
.PHONY: btc-e2e-p8-verbose
btc-e2e-p8-verbose:
	./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --verbose

# Run Bitcoin E2E workflow Pattern 8 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p8-ci
btc-e2e-p8-ci:
	./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 8
.PHONY: btc-e2e-p8-cleanup
btc-e2e-p8-cleanup:
	./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --cleanup

