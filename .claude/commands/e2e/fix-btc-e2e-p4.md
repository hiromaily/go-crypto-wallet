# Fix BTC E2E Pattern 4 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 4: P2SH-P2WSH 2-of-3 Multisig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 4 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 4 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2SH-P2WSH (BIP49 Nested SegWit Multisig) |
| **Script Type** | 2-of-3 Multisig |
| **Address Format** | `3...` (Mainnet), `2...` (regtest/testnet) |
| **Signature Requirement** | 2-of-3 (any 2 signatures complete) |
| **Descriptor** | `sh(wsh(sortedmulti(2,[fp/49'/1'/1]xpub1/0/*,[fp/49'/1'/1]xpub2/0/*,[fp/49'/1'/1]xpub3/0/*)))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="p2sh-segwit"` |

### Comparison with Other Patterns

| Item | Pattern 2 | Pattern 4 | Pattern 8 |
|------|-----------|-----------|-----------|
| Key Type | BIP44 (Legacy) | **BIP49 (P2SH-SegWit)** | BIP49 (P2SH-SegWit) |
| Signature Requirement | 2-of-3 | **2-of-3** | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | **`sh(wsh(sortedmulti(2,...)))`** | `sh(wsh(sortedmulti(3,...)))` |
| SegWit | ❌ No | **✅ Yes (wrapped)** | ✅ Yes (wrapped) |
| Transaction Size | Larger | **Smaller** | Similar to P4 |
| address_type | `legacy` | **`p2sh-segwit`** | `p2sh-segwit` |

### Differences from Pattern 3 (Single-sig)

| Item | Pattern 3 | Pattern 4 |
|------|-----------|-----------|
| Signature Requirement | Single-sig | 2-of-3 Multisig |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| Descriptor | `sh(wpkh(...))` | `sh(wsh(sortedmulti(2,...)))` |
| fullpubkey Exchange | Not required | Required |

### Differences from Pattern 8 (3-of-3)

| Item | Pattern 4 | Pattern 8 |
|------|-----------|-----------|
| Signature Requirement | 2-of-3 | 3-of-3 |
| Descriptor | `sh(wsh(sortedmulti(2,...)))` | `sh(wsh(sortedmulti(3,...)))` |
| Signing Flow | Complete with 2 signatures | Requires all 3 signatures |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p4`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 4

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 4 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Pattern 2 script (2-of-3 reference)
- @scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh - Pattern 3 script (P2SH-SegWit reference)
- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Pattern 8 script (P2SH-P2WSH 3-of-3 reference)
- @config/wallet/account/account_2of3.yaml - 2-of-3 multisig config

## Pre-check: Environment Variables

**Pattern 4 requires `WALLET_ADDRESS_TYPE="p2sh-segwit"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Implementation Steps

### Step 1: Create Script

Base on Pattern 8 (`e2e-p8-p2sh-p2wsh-3of3.sh`) with 2-of-3 modifications:

1. Filename: `e2e-p4-p2sh-p2wsh-2of3.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="p2sh-segwit"`
3. Header comments: Update to Pattern 4 specs
4. Signature requirement: Change from 3-of-3 to 2-of-3
5. Signing flow: Skip Sign2 (complete after Keygen + Sign1)

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 4: P2SH-P2WSH 2-of-3 Multisig
###############################################################################
.PHONY: btc-e2e-p4-reset
btc-e2e-p4-reset:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --reset

.PHONY: btc-e2e-p4
btc-e2e-p4:
	make btc-e2e P=4

.PHONY: btc-e2e-p4-verbose
btc-e2e-p4-verbose:
	make btc-e2e-verbose P=4

.PHONY: btc-e2e-p4-ci
btc-e2e-p4-ci:
	./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh --non-interactive

.PHONY: btc-e2e-p4-cleanup
btc-e2e-p4-cleanup:
	make btc-e2e-cleanup P=4
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-reset P=4

# With debug output
make btc-e2e-verbose P=4
```

> **Note**: For build and verification commands, see common rules.

### Step 4: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Multisig Setup | fullpubkey exchange | `internal/infrastructure/wallet/key/fullpubkey/` |
| Descriptor Export | BIP49 Descriptor | `internal/infrastructure/wallet/key/descriptor/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

## Technical Specification: P2SH-P2WSH 2-of-3

### Address Structure

```
P2SH-P2WSH 2-of-3 Address:
┌─────────────────────────────────────────────────────────────┐
│  P2SH wrapper (for legacy wallet compatibility)              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  P2WSH (SegWit Script Hash)                              │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  sortedmulti(2, pubkey1, pubkey2, pubkey3)           │ │ │
│  │  │  → 2-of-3 multisig script                            │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Descriptor Format

```
sh(wsh(sortedmulti(2,[fingerprint/49'/1'/1]xpub1/0/*,[fingerprint/49'/1'/1]xpub2/0/*,[fingerprint/49'/1'/1]xpub3/0/*)))
                   └─ 2-of-3 threshold
                              └─ BIP49 derivation path for multisig
```

### Signing Flow (2-of-3)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature) ← Complete here!
    ↓
Watch Wallet (broadcast)

※ Sign2 not required - 2 signatures satisfy 2-of-3
```

### Key Technical Points

1. **2-of-3 vs 3-of-3 Threshold**
   - Pattern 4: `sortedmulti(2, ...)` - Complete with any 2 signatures
   - Pattern 8: `sortedmulti(3, ...)` - All 3 signatures required

2. **sortedmulti vs multi**
   - `sortedmulti`: Keys sorted by pubkey (deterministic order)
   - `multi`: Keys in specified order
   - Prefer `sortedmulti` for consistency

3. **SegWit Benefits**
   - Smaller transaction size compared to Pattern 2 (legacy multisig)
   - Witness data separated from transaction
   - Lower fees

## Pattern 4 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 4 specific errors:

### address_type Mismatch

**Symptoms**: `m...` or `n...` addresses generated (Legacy P2PKH address)

**Cause**: `address_type` is `legacy` instead of `p2sh-segwit`

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh
```

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using wrong descriptor format (e.g., `sh(multi(...))` instead of `sh(wsh(sortedmulti(...)))`)

**Check**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Expected format
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "sh(wsh(sortedmulti(2,[...]xpub.../0/*)...))"
```

### Threshold Mismatch

**Symptoms**: "Incomplete signature" error after 2 signatures

**Cause**: Descriptor has `sortedmulti(3,...)` instead of `sortedmulti(2,...)`

**Check**:

```bash
# Check multisig threshold in descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep sortedmulti
# Should contain "sortedmulti(2," not "sortedmulti(3,"
```

### fullpubkey Import Error

**Symptoms**: Error during Multisig setup

**Cause**: fullpubkey format mismatch or ordering issue

**Solution**:

1. Verify correctly exported from `sign1`, `sign2`
2. Check import order to `keygen`
3. Related code: `internal/infrastructure/wallet/key/fullpubkey/`

### Too Many Signatures Error

**Symptoms**: Error when Sign2 also signs

**Cause**: Trying to apply 3rd signature when only 2 are needed

**Solution**:

- Verify signing flow stops after Sign1 (2nd signature)
- Check script logic for 2-of-3 termination condition

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2SH-P2WSH)
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": true, "iswitness": false, "embedded": { "isscript": true, "iswitness": true, "witness_version": 0, ... }
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

### PSBT Analysis

```bash
# Check signature status in PSBT
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  analyzepsbt "${psbt_hex}"
# → Check "next" field: should show "finalizer" after 2 signatures
```

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/key/descriptor/` | Descriptor processing |
| `internal/infrastructure/wallet/key/fullpubkey/` | fullpubkey processing |
| `internal/domain/address/types.go` | address_type → key_type conversion |
| `pkg/config/loader.go` | Config loader |

## Documentation Updates

After creating script, update these documents:

1. `scripts/operation/btc/e2e/README.md` - Add to script list
2. `docs/crypto/btc/operations/e2e-transaction-patterns.md` - Update implementation status (Pattern 4: ❌ → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 4 specific fixes to `P2SH-P2WSH 2-of-3` related code
- When modifying common code, verify impact on other patterns (especially 2, 3, 8)
- Confirm regression with unit tests when modifying common functions

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
make btc-e2e-cleanup P=4

# Full reset (including data)
make btc-e2e-reset P=4
```
