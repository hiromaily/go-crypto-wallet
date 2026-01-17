# Fix BTC E2E Pattern 8 Test #{issue_number}

Run and fix BTC E2E test (Pattern 8: P2SH-P2WSH 3-of-3 Multisig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command runs `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh` in **Bitcoin Core regtest environment** and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 8 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 8 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2SH-P2WSH (BIP49 wrapped SegWit) |
| **Script Type** | 3-of-3 Multisig |
| **Address Format** | `2...` (regtest/testnet P2SH) |
| **Signature Requirement** | 3-of-3 (all 3 signatures required) |
| **Descriptor** | `sh(wsh(sortedmulti(3, xpub1, xpub2, xpub3)))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="p2sh-segwit"` |

### Differences from Pattern 2 (2-of-3)

| Item | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| Key Type | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| Signature Requirement | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |
| Signing Flow | Complete with 2 signatures | Requires all 3 signatures |
| address_type | `legacy` | `p2sh-segwit` |

### Differences from Pattern 1 (Single-sig)

| Item | Pattern 1 | Pattern 8 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | 3-of-3 Multisig |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| Address Format | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey Exchange | Not required | Required |
| Account Config | `account_singlesig.yaml` | `account_3of3.yaml` |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p8`
- **Commit type**: `fix(btc)`
- **Scope**: BTC E2E Pattern 8

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 8 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Target script
- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Pattern 2 script (2-of-3 reference)
- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 script (Single-sig reference)
- @config/wallet/account/account_3of3.yaml - 3-of-3 multisig config

## Pre-check: Environment Variables

**Pattern 8 requires `WALLET_ADDRESS_TYPE="p2sh-segwit"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Execution Steps

### Step 1: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p8-reset

# Or run from existing state
./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh

# With debug output
./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --verbose
```

> **Note**: For build and verification commands, see common rules.

### Step 2: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Multisig Setup | Descriptor export/import | `internal/application/usecase/watch/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

### Step 3: Fix Code

Load appropriate skill based on error type (see Related Skills in common rules).

## Technical Specification: Signing Flow (3-of-3)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature)
    ↓
Sign2 Wallet (3rd signature) ← Complete here
    ↓
Watch Wallet (broadcast)

※ All 3 signatures required
```

### Key Technical Points

1. **Descriptor Solvability**
   - P2SH-P2WSH addresses require Descriptor import
   - Traditional address import makes UTXOs "unsolvable"
   - Solve with `importdescriptors` RPC

2. **3-of-3 Multisig**
   - Uses addresses derived from 3 xpubs
   - Private keys from all signers required
   - Partial signatures transmitted via PSBT format

3. **HD Key Derivation Path**
   - BIP49: `m/49'/0'/account'/change/index`
   - Same index keys used across wallets
   - Different from Pattern 2 (BIP44)

## Pattern 8 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 8 specific errors:

### "Unsolvable" UTXO Error

**Cause**: P2SH-P2WSH address imported without Descriptor

**Solution**:

1. Export correct format with `descriptor export`
2. Import with `descriptor import --account`
3. Related code: `internal/infrastructure/wallet/api/btc/descriptor.go`

```bash
# Debug: Check address info
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# solvable: true required
```

### Signing Error (3-of-3 specific)

**Cause**: Private key not imported or key mismatch

**Solution**:

1. Verify `import privkey` executed in each wallet
2. Check fullpubkey import order
3. Verify all 3 signatures applied
4. Related code: `internal/infrastructure/wallet/api/btc/sign.go`

**Check**:

```bash
# Check PSBT signature status
btc_cli "btc-watch" analyzepsbt "${psbt_hex}"
```

### fullpubkey Import Error

**Symptoms**: Error during Multisig setup

**Cause**: fullpubkey format mismatch or ordering issue

**Solution**:

1. Verify correctly exported from `sign1`, `sign2`
2. Check import order to `keygen`
3. Related code: `internal/infrastructure/wallet/key/fullpubkey/`

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using BIP44 format (Legacy)

**Solution**: Verify `address_type` is `p2sh-segwit`

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "p2sh-segwit"
```

**Expected format** (Pattern 8):

```
sh(wsh(sortedmulti(3, [fp/49'/0'/0']xpub1/0/*, xpub2/0/*, xpub3/0/*)))
```

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/api/btc/descriptor.go` | Descriptor processing |
| `internal/infrastructure/wallet/api/btc/sign.go` | PSBT signing |
| `internal/infrastructure/wallet/key/fullpubkey/` | fullpubkey processing |
| `internal/domain/wallet/key/` | Key domain model |

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 8 specific fixes to `P2SH-P2WSH` related code
- When modifying common code, verify impact on other patterns (especially 1, 2)
- Confirm regression with unit tests when modifying common functions

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh --cleanup

# Full reset (including data)
make btc-e2e-p8-reset
```
