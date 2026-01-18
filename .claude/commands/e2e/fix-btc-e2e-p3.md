# Fix BTC E2E Pattern 3 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 3: P2SH-P2WPKH Single-sig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 3 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 3 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2SH-P2WPKH (BIP49 Nested SegWit) |
| **Script Type** | Single-sig |
| **Address Format** | `3...` (Mainnet), `2...` (regtest/testnet) |
| **Signature Requirement** | Single-sig (1 signature) |
| **Descriptor** | `sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))` |
| **Required Wallets** | watch, keygen |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="p2sh-segwit"` |

### Differences from Pattern 1 (P2PKH Single-sig)

| Item | Pattern 1 | Pattern 3 |
|------|-----------|-----------|
| Key Type | BIP44 (Legacy) | BIP49 (Nested SegWit) |
| Address Format | `m.../n...` (P2PKH) | `2...` (P2SH) |
| Descriptor | `pkh(...)` | `sh(wpkh(...))` |
| Environment Variable | `legacy` | `p2sh-segwit` |
| Transaction Size | Larger | Smaller (SegWit discount) |

### Differences from Pattern 8 (P2SH-P2WSH 3-of-3)

| Item | Pattern 3 | Pattern 8 |
|------|-----------|-----------|
| Signature Requirement | Single-sig | 3-of-3 Multisig |
| Descriptor | `sh(wpkh(...))` | `sh(wsh(sortedmulti(3,...)))` |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| fullpubkey Exchange | Not required | Required |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p3`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 3

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 3 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 script (Single-sig base)
- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Pattern 8 script (P2SH-SegWit reference)
- @config/wallet/account/account.yaml - Single-sig account config

## Pre-check: Environment Variables

**Pattern 3 requires `WALLET_ADDRESS_TYPE="p2sh-segwit"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Implementation Steps

### Step 1: Create Script

Base on Pattern 1 (`e2e-p1-p2pkh-singlesig.sh`) with these changes:

1. Filename: `e2e-p3-p2sh-p2wpkh-singlesig.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="p2sh-segwit"`
3. Header comments: Update to Pattern 3 specs
4. Address validation logic: Check for `2...` format

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 3: P2SH-P2WPKH Single-sig
###############################################################################
.PHONY: btc-e2e-p3-reset
btc-e2e-p3-reset:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --reset

.PHONY: btc-e2e-p3
btc-e2e-p3:
	make btc-e2e P=3

.PHONY: btc-e2e-p3-verbose
btc-e2e-p3-verbose:
	make btc-e2e-verbose P=3

.PHONY: btc-e2e-p3-ci
btc-e2e-p3-ci:
	./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --non-interactive

.PHONY: btc-e2e-p3-cleanup
btc-e2e-p3-cleanup:
	make btc-e2e-cleanup P=3
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-reset P=3

# With debug output
make btc-e2e-verbose P=3
```

> **Note**: For build and verification commands, see common rules.

### Step 4: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Descriptor Export | BIP49 Descriptor | `internal/infrastructure/wallet/key/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

## Technical Specification: P2SH-P2WPKH (Nested SegWit)

### Address Structure

```
P2SH-P2WPKH Address:
┌─────────────────────────────────────────────────────────────┐
│  P2SH wrapper                                                │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  P2WPKH (Native SegWit)                                  │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Public Key Hash                                     │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Descriptor Format

```
sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
           └─ BIP49 derivation path
```

### Signing Flow (Single-sig)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (sign with single key)
    ↓
Watch Wallet (broadcast)
```

## Pattern 3 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 3 specific errors:

### address_type Mismatch

**Symptoms**: `m...` or `n...` addresses generated (P2PKH address)

**Cause**: `address_type` is `legacy`

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh
```

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using BIP44 format (`pkh(...)`)

**Check**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Expected format
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "sh(wpkh([...]xpub.../0/*))"
```

### key_type Auto-derivation Check

**Check**: Verify `key_type` correctly derived from `address_type`

| address_type | Expected key_type |
|--------------|-------------------|
| `p2sh-segwit` | `bip49` |

Related code: `AddrType.ToKeyType()` in `internal/domain/address/types.go`

### Transaction Signing Error

**Symptoms**: Witness-related error during signing

**Cause**: Issue with SegWit transaction witness data processing

**Check**:

```bash
# Analyze PSBT
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  analyzepsbt "${psbt_hex}"
```

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2SH-P2WPKH)
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": true, "iswitness": false, "embedded": { "isscript": false, "iswitness": true, "witness_version": 0, ... }
```

### Descriptor Check

```bash
# Keygen descriptor list
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch descriptor list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors
```

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/key/descriptor/` | Descriptor processing |
| `internal/domain/address/types.go` | address_type → key_type conversion |
| `pkg/config/loader.go` | Config loader |

## Documentation Updates

After creating script, update these documents:

1. `scripts/operation/btc/e2e/README.md` - Add to script list
2. `docs/crypto/btc/operations/e2e-transaction-patterns.md` - Update implementation status
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 3 specific fixes to `P2SH-P2WPKH Single-sig` related code
- When modifying common code, verify impact on other patterns (especially 1, 2, 8)
- Confirm regression with unit tests when modifying common functions

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
make btc-e2e-cleanup P=3

# Full reset (including data)
make btc-e2e-reset P=3
```
