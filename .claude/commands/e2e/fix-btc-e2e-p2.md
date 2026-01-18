# Fix BTC E2E Pattern 2 Test #{issue_number}

Run and fix BTC E2E test (Pattern 2: P2PKH 2-of-3 Multisig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command runs `scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh` in **Bitcoin Core regtest environment** and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 2 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 2 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2PKH (BIP44 Legacy) |
| **Script Type** | 2-of-3 Multisig (P2SH wrapped) |
| **Address Format** | `2...` (regtest/testnet P2SH) |
| **Signature Requirement** | 2-of-3 (any 2 signatures complete) |
| **Descriptor** | `sh(multi(2, [fp/44'/0'/0']xpub1/0/*, xpub2/0/*, xpub3/0/*))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="legacy"` |

### Differences from Pattern 1 (Single-sig)

| Item | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | 2-of-3 Multisig |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| Address Format | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey Exchange | Not required | Required |
| Account Config | `account_singlesig.yaml` | `account_2of3.yaml` |

### Differences from Pattern 8 (3-of-3)

| Item | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| Key Type | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| Signature Requirement | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |
| Signing Flow | Complete with 2 signatures | Requires all 3 signatures |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p2`
- **Commit type**: `fix(btc)`
- **Scope**: BTC E2E Pattern 2

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 2 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Target script
- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 script (Single-sig reference)
- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Pattern 8 script (Multisig reference)
- @config/wallet/account/account_2of3.yaml - 2-of-3 multisig config

## Pre-check: Environment Variables

**Pattern 2 requires `WALLET_ADDRESS_TYPE="legacy"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Execution Steps

### Step 1: Run E2E Test

**Always use Makefile targets** (build is automatic):

```bash
# Full reset and run (recommended)
make btc-e2e-reset P=2

# Or run from existing state
make btc-e2e P=2

# With debug output
make btc-e2e-verbose P=2
```

> **Note**: Do NOT run scripts directly. Makefile includes build dependency.

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

## Technical Specification: Signing Flow (2-of-3)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature) ← Complete here
    ↓
Watch Wallet (broadcast)

※ Sign2 not required - 2 signatures satisfy 2-of-3
```

### Key Technical Points

1. **2-of-3 vs 3-of-3**
   - Pattern 2: Complete with any 2 signatures
   - Pattern 8: All 3 signatures required
   - Signing flow control differs

2. **Descriptor Format**
   - `sh(multi(2, ...))` - P2SH wrapper + 2-of-3 multisig
   - Can use `multi` instead of `sortedmulti`
   - Note key order

3. **HD Key Derivation Path**
   - BIP44: `m/44'/0'/account'/change/index`
   - Different from Pattern 8 (BIP49)

## Pattern 2 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 2 specific errors:

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using BIP49 format (P2SH-SegWit)

**Solution**: Verify `address_type` is `legacy`

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
```

### Address Format Differs

**Symptoms**: `m...` or `n...` addresses generated (Single-sig P2PKH address)

**Cause**: Generating Single-sig addresses

**Solution**:

- Check `multisig` section in `account_2of3.yaml`
- Verify fullpubkey import succeeded
- Check `required: 2` is set

### Insufficient/Too Many Signatures

**Symptoms**: "Incomplete signature" or "Too many signatures" error on transaction send

**Cause**:

- 2nd signature not applied correctly
- Applied 3rd signature (unnecessary)

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

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/key/fullpubkey/` | fullpubkey processing |
| `internal/domain/wallet/key/` | Key domain model |

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 2 specific fixes to `P2PKH 2-of-3` related code
- When modifying common code, verify impact on other patterns (especially 1, 8)
- Confirm regression with unit tests when modifying common functions

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
make btc-e2e-cleanup P=2

# Full reset (including data)
make btc-e2e-reset P=2
```
