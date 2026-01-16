---
paths: "scripts/operation/btc/e2e/**"
---

# BTC E2E Script Development Rules

Rules applied when creating or modifying Bitcoin E2E scripts.

## Required Documentation

Read the following documents before creating or modifying scripts:

| Document | Contents |
|----------|----------|
| `docs/crypto/btc/e2e_transaction_patterns.md` | Detailed specifications for all 11 patterns |
| `scripts/operation/common.sh` | Common utility functions |
| `pkg/config/README.md` | Configuration override via environment variables |
| `config/wallet/README.md` | Wallet configuration file policies |

## Script Structure Conventions

### Header Comments

Each script must include header comments in the following format:

```bash
#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern N: [Pattern Name]
# This script automates the complete Bitcoin workflow for [description]
# Usage: ./scripts/operation/btc/e2e/e2e-pN-{description}.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/btc/e2e_transaction_patterns.md - E2E transaction patterns
#
# Transaction Pattern:
#   Pattern N: BTC [Address Type] [Signature Requirement]
#   - Address Type: [P2PKH/P2SH-P2WPKH/etc.]
#   - Address Format: `...` (Mainnet), `...` (Testnet/Regtest)
#   - Signature Requirement: [Single-sig/2-of-3/3-of-3]
#   - Descriptor: [descriptor format]
#
# Required Config Settings:
#   - config/wallet/btc_watch.yaml:  address_type: "[type]"
#   - config/wallet/btc_keygen.yaml: address_type: "[type]"
```

### Environment Variable Section

```bash
###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# These environment variables override config file values.
# Priority: Environment Variables > Config File > Default Values
#
# Pattern N requires:
#   - address_type: "[type]" (derives key_type automatically)
export WALLET_ADDRESS_TYPE="[type]"
```

### Account Configuration Files

Use the appropriate configuration file based on the pattern:

| Pattern | Account Config |
|---------|----------------|
| Single-sig | `config/wallet/account.yaml` |
| 2-of-3 Multisig | `config/wallet/account_2of3.yaml` |
| 3-of-3 Multisig | `config/wallet/account_3of3.yaml` |

## Configuration File Policy (Important)

### ❌ Do NOT Edit Config Files Directly

Do **not** edit configuration files (`btc_watch.yaml`, `btc_keygen.yaml`, etc.) directly.
Use **environment variables** to override settings when different values are needed.

### ✅ Override via Environment Variables

```bash
# Export within the script
export WALLET_ADDRESS_TYPE="legacy"

# Priority order:
# 1. Environment Variables (highest priority)
# 2. Config File
# 3. Default Values (lowest priority)
```

### Automatic key_type Derivation

`key_type` is **automatically derived** from `address_type`. The `WALLET_KEY_TYPE` environment variable is not needed.

| address_type | Derived key_type | Use Case |
|--------------|------------------|----------|
| `legacy` | `bip44` | P2PKH (Pattern 1, 2) |
| `p2sh-segwit` | `bip49` | P2SH-P2WPKH/P2SH-P2WSH (Pattern 3, 4, 8) |
| `bech32` | `bip84` | Native SegWit (Pattern 5, 6, 7) |
| `taproot` / `bech32m` | `bip86` | Taproot (Pattern 9, 10, 11) |

Reference: `AddrType.ToKeyType()` in `internal/domain/address/types.go`

## Pattern-Specific Settings

| Pattern | Description | address_type | Address Format | Signature |
|---------|-------------|--------------|----------------|-----------|
| 1 | P2PKH Single-sig | `legacy` | `m.../n...` | Single |
| 2 | P2PKH 2-of-3 | `legacy` | `2...` (P2SH) | 2-of-3 |
| 3 | P2SH-P2WPKH Single-sig | `p2sh-segwit` | `2...` | Single |
| 4 | P2SH-P2WSH 2-of-3 | `p2sh-segwit` | `2...` | 2-of-3 |
| 5 | P2WPKH Native SegWit Single-sig | `bech32` | `bcrt1q...` | Single |
| 6 | P2WSH 2-of-3 | `bech32` | `bcrt1q...` (62 chars) | 2-of-3 |
| 7 | P2WSH 3-of-3 | `bech32` | `bcrt1q...` (62 chars) | 3-of-3 |
| 8 | P2SH-P2WSH 3-of-3 | `p2sh-segwit` | `2...` | 3-of-3 |
| 9 | P2TR Taproot Single-sig | `taproot` | `bcrt1p...` | Single |

## Build and Verification Rules

### Build Before E2E Execution (Required)

```bash
# Build all wallet binaries
make build-all

# Output: ${GOPATH}/bin/watch, keygen, sign1, sign2
```

**Important**: Do not run `go build` directly.

### Verification After Go Code Changes

```bash
# 1. Lint, build, and test
make go-lint && make check-build && make gotest

# 2. Rebuild binaries (required)
make build-all

# 3. Run E2E test (e.g., make btc-e2e-p1-reset for Pattern 1)
make btc-e2e-pN-reset
```

### Verification After Shell Script Changes

```bash
make shfmt
```

## Retry Limit

**If the fix-test cycle exceeds 5 iterations, organize progress and report.**

### Escalation Conditions

- Same error occurs repeatedly
- Deep understanding of Bitcoin specifications required
- Large-scale Go code changes needed

### Progress Report Format

```markdown
## Progress Report

### Error Details
[Error message that occurred]

### Attempted Fixes
1. [Fix attempt 1]
2. [Fix attempt 2]

### Current State
[Description of current state]

### Environment Variable Status
- WALLET_ADDRESS_TYPE: [value]

### Next Steps
[Required next actions]
```

## Security Rules

- ❌ Do NOT log private keys
- ❌ Do NOT use test passphrases/RPC credentials in production
- Reference: `docs/standards/security.md`

## Avoiding Impact on Other Patterns

### When Modifying common.sh

- Always verify impact on other E2E patterns
- Do not break existing pattern behavior

### Pattern-Specific Changes

- Set environment variables locally within each script
- Confirm regression with unit tests when modifying shared code

## Common Errors

### "No utxo" Error

1. Verify Descriptor is correctly imported
2. Confirm block generation (101+) is complete
3. Check `address_type` is correct

```bash
# Debug
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
```

### address_type Mismatch

When generated address differs from expected:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE
```

### duplicate key Error

When data from previous execution remains:

```bash
make btc-e2e-pN-reset  # e.g., make btc-e2e-p1-reset
```

## Makefile Targets

Add targets to `make/btc_e2e.mk` when creating new scripts.

**Naming Convention**: `btc-e2e-pN` where N is the pattern number.

```makefile
# Example for Pattern 1 (e2e-p1-p2pkh-singlesig.sh)
.PHONY: btc-e2e-p1-reset
btc-e2e-p1-reset:
  ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --reset

.PHONY: btc-e2e-p1
btc-e2e-p1:
  ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh

.PHONY: btc-e2e-p1-verbose
btc-e2e-p1-verbose:
  ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --verbose

.PHONY: btc-e2e-p1-ci
btc-e2e-p1-ci:
  ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --non-interactive

.PHONY: btc-e2e-p1-cleanup
btc-e2e-p1-cleanup:
  ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --cleanup
```

**Current Scripts**:

| Pattern | Script | Make Target | Status |
|---------|--------|-------------|--------|
| 1 | `e2e-p1-p2pkh-singlesig.sh` | `btc-e2e-p1` | ✅ |
| 2 | `e2e-p2-p2pkh-2of3.sh` | `btc-e2e-p2` | ✅ |
| 3 | `e2e-p3-p2sh-p2wpkh-singlesig.sh` | `btc-e2e-p3` | ✅ |
| 4 | `e2e-p4-p2sh-p2wsh-2of3.sh` | `btc-e2e-p4` | ✅ |
| 5 | `e2e-p5-p2wpkh-singlesig.sh` | `btc-e2e-p5` | ✅ |
| 6 | `e2e-p6-p2wsh-2of3.sh` | `btc-e2e-p6` | ✅ |
| 7 | `e2e-p7-p2wsh-3of3.sh` | `btc-e2e-p7` | ✅ |
| 8 | `e2e-p8-p2sh-p2wsh-3of3.sh` | `btc-e2e-p8` | ✅ |
| 9 | `e2e-p9-p2tr-singlesig.sh` | `btc-e2e-p9` | ✅ |
| 10 | `e2e-p10-p2tr-musig2.sh` | `btc-e2e-p10` | ✅ |
| 11 | `e2e-p11-p2tr-tapscript.sh` | `btc-e2e-p11` | ✅ |

## Related Skills

| Skill | Use Case |
|-------|----------|
| `shell-scripts` | Script creation/modification |
| `go-development` | Go code changes |
| `makefile-update` | Makefile updates |
| `git-workflow` | Branch management |
