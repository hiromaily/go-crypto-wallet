# Fix BTC E2E Pattern 6 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 6: P2WSH Native SegWit 2-of-3 Multisig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 6 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 6 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2WSH (BIP84 Native SegWit Multisig) |
| **Script Type** | 2-of-3 Multisig |
| **Address Format** | `bc1q...` (Mainnet), `tb1q...` (Testnet), `bcrt1q...` (regtest) |
| **Signature Requirement** | 2-of-3 (any 2 signatures complete) |
| **Descriptor** | `wsh(sortedmulti(2,[fp/84'/1'/1]xpub1/0/*,[fp/84'/1'/1]xpub2/0/*,[fp/84'/1'/1]xpub3/0/*))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="bech32"` |

### Comparison with Related Patterns

| Item | Pattern 4 (P2SH-P2WSH) | Pattern 5 (P2WPKH) | Pattern 6 (P2WSH) | Pattern 7 (P2WSH 3-of-3) |
|------|------------------------|--------------------|--------------------|--------------------------|
| BIP | BIP49 | BIP84 | **BIP84** | BIP84 |
| Type | 2-of-3 Multisig | Single-sig | **2-of-3 Multisig** | 3-of-3 Multisig |
| Address | `2...` (P2SH) | `bcrt1q...` | **`bcrt1q...`** | `bcrt1q...` |
| Descriptor | `sh(wsh(sortedmulti(2,...)))` | `wpkh(...)` | **`wsh(sortedmulti(2,...))`** | `wsh(sortedmulti(3,...))` |
| P2SH Wrapper | Yes | No | **No** | No |
| Legacy Compatible | Yes | No | **No** | No |
| Transaction Size | Medium | Small | **Smallest (multisig)** | Similar to P6 |

### Why Use P2WSH (Native SegWit Multisig)?

| Aspect | Description |
|--------|-------------|
| ✅ Most efficient multisig | Smallest transaction size for multisig (no P2SH wrapper overhead) |
| ✅ Lower fees | ~15% smaller than P2SH-P2WSH (Pattern 4) |
| ✅ Bech32 encoding | Error detection and correction capabilities |
| ✅ Modern standard | The recommended format for new multisig wallets |
| ❌ Legacy incompatible | Very old wallets may not recognize `bc1q...` addresses |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p6`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 6

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 6 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh - Pattern 4 script (P2SH-P2WSH 2-of-3 reference)
- @scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh - Pattern 5 script (Native SegWit reference)
- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Pattern 2 script (2-of-3 multisig reference)
- @config/wallet/account_2of3.yaml - 2-of-3 multisig config

## Pre-check: Environment Variables

**Pattern 6 requires `WALLET_ADDRESS_TYPE="bech32"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Implementation Steps

### Step 1: Create Script

Base on Pattern 4 (`e2e-p4-p2sh-p2wsh-2of3.sh`) with these changes:

1. Filename: `e2e-p6-p2wsh-2of3.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="bech32"`
3. Header comments: Update to Pattern 6 specs
4. Address validation logic: Check for `bcrt1q...` format (regtest bech32)
5. Descriptor format: `wsh(sortedmulti(2,...))` instead of `sh(wsh(sortedmulti(2,...)))`

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 6: P2WSH Native SegWit 2-of-3 Multisig
###############################################################################
.PHONY: btc-e2e-p6-reset
btc-e2e-p6-reset:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --reset

.PHONY: btc-e2e-p6
btc-e2e-p6:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh

.PHONY: btc-e2e-p6-verbose
btc-e2e-p6-verbose:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --verbose

.PHONY: btc-e2e-p6-ci
btc-e2e-p6-ci:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --non-interactive

.PHONY: btc-e2e-p6-cleanup
btc-e2e-p6-cleanup:
	./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --cleanup
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p6-reset

# With debug output
./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --verbose
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
| Descriptor Export | BIP84 Descriptor | `internal/infrastructure/wallet/key/descriptor/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

## Technical Specification: P2WSH Native SegWit 2-of-3

### Address Structure

```
P2WSH 2-of-3 Address (Native SegWit):
┌─────────────────────────────────────────────────────────────┐
│  Bech32 Encoding (NO P2SH wrapper)                           │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 0                                      │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Witness Script Hash (32 bytes)                      │ │ │
│  │  │  ┌─────────────────────────────────────────────────┐ │ │ │
│  │  │  │  sortedmulti(2, pubkey1, pubkey2, pubkey3)       │ │ │ │
│  │  │  │  → 2-of-3 multisig script                        │ │ │ │
│  │  │  └─────────────────────────────────────────────────┘ │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1q + 58 characters (mainnet, 32-byte witness program)
        bcrt1q + 58 characters (regtest)

Note: P2WSH addresses are longer than P2WPKH (62 vs 42 chars)
      due to 32-byte script hash vs 20-byte pubkey hash
```

### Descriptor Format

```
wsh(sortedmulti(2,[fingerprint/84'/1'/1]xpub1/0/*,[fingerprint/84'/1'/1]xpub2/0/*,[fingerprint/84'/1'/1]xpub3/0/*))
    └─ No sh() wrapper!           └─ BIP84 derivation path for multisig
```

### Key Derivation Path (BIP84 for Multisig)

```
m / 84' / coin_type' / account' / change / address_index
         └─ 1' for testnet/regtest
                      └─ 1' for multisig accounts (non-hardened in descriptor)
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

## Pattern 6 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 6 specific errors:

### address_type Mismatch

**Symptoms**: `2...` or `3...` addresses generated instead of `bcrt1q...`

**Cause**: `address_type` is `p2sh-segwit` instead of `bech32`

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh
```

### Descriptor Format Error (P2SH Wrapper)

**Symptoms**: Error during Descriptor export/import, or `2...` addresses generated

**Cause**: Using `sh(wsh(...))` instead of `wsh(...)`

**Check**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Expected format (NO sh() wrapper)
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "wsh(sortedmulti(2,[...]xpub.../0/*)...)"  (NOT sh(wsh(...)))
```

### Address Length Difference

**Note**: P2WSH addresses are longer than P2WPKH addresses:

| Type | Length | Example |
|------|--------|---------|
| P2WPKH | 42 chars | `bcrt1q...` (20-byte pubkey hash) |
| P2WSH | 62 chars | `bcrt1q...` (32-byte script hash) |

Both start with `bcrt1q` but P2WSH is 20 characters longer.

### fullpubkey Import Error

**Symptoms**: Error during Multisig setup

**Cause**: fullpubkey format mismatch or ordering issue

**Solution**:

1. Verify correctly exported from `sign1`, `sign2`
2. Check import order to `keygen`
3. Related code: `internal/infrastructure/wallet/key/fullpubkey/`

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

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2WSH Native SegWit)
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": true, "iswitness": true, "witness_version": 0, "witness_program": "<32-byte-hash>"
# Note: isscript=true (it's a script), iswitness=true (native segwit)
```

### Descriptor Check

```bash
# Keygen descriptor list
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch descriptor list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors

# Verify wsh descriptor format (NOT sh(wsh))
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep -E "^\"wsh\("
```

### Address Derivation Test

```bash
# Derive addresses from descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  deriveaddresses "wsh(sortedmulti(2,[fp1]xpub1/0/*,[fp2]xpub2/0/*,[fp3]xpub3/0/*))" "[0,5]"
# Should return bcrt1q... addresses (62 chars each)
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
2. `docs/crypto/btc/e2e_transaction_patterns.md` - Update implementation status (Pattern 6: ❌ → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 6 specific fixes to `P2WSH 2-of-3` related code
- When modifying common code, verify impact on other patterns (especially 4, 5, 7)
- Confirm regression with unit tests when modifying common functions

### P2WSH vs P2SH-P2WSH

| Item | P2WSH (Pattern 6) | P2SH-P2WSH (Pattern 4) |
|------|-------------------|------------------------|
| Descriptor | `wsh(...)` | `sh(wsh(...))` |
| Address | `bcrt1q...` (bech32) | `2...` (base58) |
| Size | Smaller | Larger (P2SH overhead) |
| Compatibility | Modern wallets | All wallets |

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh --cleanup

# Full reset (including data)
make btc-e2e-p6-reset
```
