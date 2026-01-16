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
# E2E Testing - Pattern 4: P2SH-P2WSH 2-of-3 Multisig
###############################################################################
# Run Bitcoin E2E workflow Pattern 4 from completely fresh state (recommended)
.PHONY: btc-e2e-p4-reset
btc-e2e-p4-reset:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 4
.PHONY: btc-e2e-p4
btc-e2e-p4:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh

# Run Bitcoin E2E workflow Pattern 4 with verbose output
.PHONY: btc-e2e-p4-verbose
btc-e2e-p4-verbose:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --verbose

# Run Bitcoin E2E workflow Pattern 4 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p4-ci
btc-e2e-p4-ci:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 4
.PHONY: btc-e2e-p4-cleanup
btc-e2e-p4-cleanup:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --cleanup

###############################################################################
# E2E Testing - Pattern 5: P2WPKH Native SegWit Single-sig
###############################################################################
# Run Bitcoin E2E workflow Pattern 5 from completely fresh state (recommended)
.PHONY: btc-e2e-p5-reset
btc-e2e-p5-reset:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 5
.PHONY: btc-e2e-p5
btc-e2e-p5:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh

# Run Bitcoin E2E workflow Pattern 5 with verbose output
.PHONY: btc-e2e-p5-verbose
btc-e2e-p5-verbose:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --verbose

# Run Bitcoin E2E workflow Pattern 5 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p5-ci
btc-e2e-p5-ci:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 5
.PHONY: btc-e2e-p5-cleanup
btc-e2e-p5-cleanup:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --cleanup

###############################################################################
# E2E Testing - Pattern 6: P2WSH Native SegWit 2-of-3 Multisig
###############################################################################
# Run Bitcoin E2E workflow Pattern 6 from completely fresh state (recommended)
.PHONY: btc-e2e-p6-reset
btc-e2e-p6-reset:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 6
.PHONY: btc-e2e-p6
btc-e2e-p6:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh

# Run Bitcoin E2E workflow Pattern 6 with verbose output
.PHONY: btc-e2e-p6-verbose
btc-e2e-p6-verbose:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --verbose

# Run Bitcoin E2E workflow Pattern 6 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p6-ci
btc-e2e-p6-ci:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 6
.PHONY: btc-e2e-p6-cleanup
btc-e2e-p6-cleanup:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --cleanup

###############################################################################
# E2E Testing - Pattern 7: P2WSH Native SegWit 3-of-3 Multisig
###############################################################################
# Run Bitcoin E2E workflow Pattern 7 from completely fresh state (recommended)
.PHONY: btc-e2e-p7-reset
btc-e2e-p7-reset:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 7
.PHONY: btc-e2e-p7
btc-e2e-p7:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh

# Run Bitcoin E2E workflow Pattern 7 with verbose output
.PHONY: btc-e2e-p7-verbose
btc-e2e-p7-verbose:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --verbose

# Run Bitcoin E2E workflow Pattern 7 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p7-ci
btc-e2e-p7-ci:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 7
.PHONY: btc-e2e-p7-cleanup
btc-e2e-p7-cleanup:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --cleanup

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

###############################################################################
# E2E Testing - Pattern 9: P2TR Taproot Single-sig
###############################################################################
# Run Bitcoin E2E workflow Pattern 9 from completely fresh state (recommended)
.PHONY: btc-e2e-p9-reset
btc-e2e-p9-reset:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 9
.PHONY: btc-e2e-p9
btc-e2e-p9:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh

# Run Bitcoin E2E workflow Pattern 9 with verbose output
.PHONY: btc-e2e-p9-verbose
btc-e2e-p9-verbose:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --verbose

# Run Bitcoin E2E workflow Pattern 9 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p9-ci
btc-e2e-p9-ci:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 9
.PHONY: btc-e2e-p9-cleanup
btc-e2e-p9-cleanup:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --cleanup


###############################################################################
# E2E Testing - Pattern 10: P2TR MuSig2 N-of-N
###############################################################################
# Run Bitcoin E2E workflow Pattern 10 from completely fresh state (recommended)
.PHONY: btc-e2e-p10-reset
btc-e2e-p10-reset:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 10
.PHONY: btc-e2e-p10
btc-e2e-p10:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh

# Run Bitcoin E2E workflow Pattern 10 with verbose output
.PHONY: btc-e2e-p10-verbose
btc-e2e-p10-verbose:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --verbose

# Run Bitcoin E2E workflow Pattern 10 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p10-ci
btc-e2e-p10-ci:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 10
.PHONY: btc-e2e-p10-cleanup
btc-e2e-p10-cleanup:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --cleanup
###############################################################################
# E2E Testing - Pattern 11: P2TR Tapscript M-of-N
###############################################################################
# Run Bitcoin E2E workflow Pattern 11 from completely fresh state (recommended)
.PHONY: btc-e2e-p11-reset
btc-e2e-p11-reset:
	./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 11
.PHONY: btc-e2e-p11
btc-e2e-p11:
	./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh

# Run Bitcoin E2E workflow Pattern 11 with verbose output
.PHONY: btc-e2e-p11-verbose
btc-e2e-p11-verbose:
	./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --verbose

# Run Bitcoin E2E workflow Pattern 11 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p11-ci
btc-e2e-p11-ci:
	./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 11
.PHONY: btc-e2e-p11-cleanup
btc-e2e-p11-cleanup:
	./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --cleanup
