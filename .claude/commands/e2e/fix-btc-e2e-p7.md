# Fix BTC E2E Pattern 7 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 7: P2WSH Native SegWit 3-of-3 Multisig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 7 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 7 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2WSH (BIP84 Native SegWit Multisig) |
| **Script Type** | 3-of-3 Multisig |
| **Address Format** | `bc1q...` (Mainnet), `tb1q...` (Testnet), `bcrt1q...` (regtest) |
| **Signature Requirement** | 3-of-3 (all 3 signatures required) |
| **Descriptor** | `wsh(sortedmulti(3,[fp/84'/1'/1]xpub1/0/*,[fp/84'/1'/1]xpub2/0/*,[fp/84'/1'/1]xpub3/0/*))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="bech32"` |

### Comparison with Related Patterns

| Item | Pattern 6 (P2WSH 2-of-3) | Pattern 7 (P2WSH 3-of-3) | Pattern 8 (P2SH-P2WSH 3-of-3) |
|------|--------------------------|--------------------------|-------------------------------|
| BIP | BIP84 | **BIP84** | BIP49 |
| Address | `bcrt1q...` | **`bcrt1q...`** | `2...` |
| Descriptor | `wsh(sortedmulti(2,...))` | **`wsh(sortedmulti(3,...))`** | `sh(wsh(sortedmulti(3,...)))` |
| Threshold | 2-of-3 | **3-of-3** | 3-of-3 |
| P2SH Wrapper | No | **No** | Yes |
| Legacy Compatible | No | **No** | Yes |
| Transaction Size | Smallest | **Smallest** | Medium |

### Why Use P2WSH 3-of-3 (Native SegWit)?

| Aspect | Description |
|--------|-------------|
| ✅ Highest security | All 3 keys must sign - no single point of failure |
| ✅ Most efficient 3-of-3 | Smallest transaction size for 3-of-3 multisig |
| ✅ Modern standard | The recommended format for high-security multisig |
| ✅ Bech32 encoding | Error detection and correction capabilities |
| ❌ All must sign | Requires coordination of all 3 signers |
| ❌ Legacy incompatible | Very old wallets may not recognize `bc1q...` addresses |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p7`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 7

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 7 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh - Pattern 6 script (P2WSH 2-of-3 reference)
- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Pattern 8 script (P2SH-P2WSH 3-of-3 reference)
- @config/wallet/account_3of3.yaml - 3-of-3 multisig config

## Pre-check: Environment Variables

**Pattern 7 requires `WALLET_ADDRESS_TYPE="bech32"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Implementation Steps

### Step 1: Create Script

Base on Pattern 6 (`e2e-p6-p2wsh-2of3.sh`) or Pattern 8 (`e2e-p8-p2sh-p2wsh-3of3.sh`) with these changes:

1. Filename: `e2e-p7-p2wsh-3of3.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="bech32"`
3. Header comments: Update to Pattern 7 specs
4. Address validation logic: Check for `bcrt1q...` format (regtest bech32)
5. Descriptor format: `wsh(sortedmulti(3,...))` (3-of-3 threshold, no sh() wrapper)
6. Signing flow: All 3 signatures required (Keygen + Sign1 + Sign2)
7. Account config: Use `account_3of3.yaml`

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 7: P2WSH Native SegWit 3-of-3 Multisig
###############################################################################
.PHONY: btc-e2e-p7-reset
btc-e2e-p7-reset:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --reset

.PHONY: btc-e2e-p7
btc-e2e-p7:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh

.PHONY: btc-e2e-p7-verbose
btc-e2e-p7-verbose:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --verbose

.PHONY: btc-e2e-p7-ci
btc-e2e-p7-ci:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --non-interactive

.PHONY: btc-e2e-p7-cleanup
btc-e2e-p7-cleanup:
	./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --cleanup
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p7-reset

# With debug output
./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --verbose
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

## Technical Specification: P2WSH Native SegWit 3-of-3

### Address Structure

```
P2WSH 3-of-3 Address (Native SegWit):
┌─────────────────────────────────────────────────────────────┐
│  Bech32 Encoding (NO P2SH wrapper)                           │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 0                                      │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Witness Script Hash (32 bytes)                      │ │ │
│  │  │  ┌─────────────────────────────────────────────────┐ │ │ │
│  │  │  │  sortedmulti(3, pubkey1, pubkey2, pubkey3)       │ │ │ │
│  │  │  │  → 3-of-3 multisig script (ALL must sign)        │ │ │ │
│  │  │  └─────────────────────────────────────────────────┘ │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1q + 58 characters (mainnet, 32-byte witness program)
        bcrt1q + 58 characters (regtest)
```

### Descriptor Format

```
wsh(sortedmulti(3,[fingerprint/84'/1'/1]xpub1/0/*,[fingerprint/84'/1'/1]xpub2/0/*,[fingerprint/84'/1'/1]xpub3/0/*))
                └─ 3-of-3 threshold (all must sign)
```

### Key Derivation Path (BIP84 for Multisig)

```
m / 84' / coin_type' / account' / change / address_index
         └─ 1' for testnet/regtest
                      └─ 1' for multisig accounts
```

### Signing Flow (3-of-3)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature)
    ↓
Sign2 Wallet (3rd signature) ← Complete here!
    ↓
Watch Wallet (broadcast)

※ ALL 3 signatures required - no early completion
```

## Pattern 7 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 7 specific errors:

### address_type Mismatch

**Symptoms**: `2...` or `3...` addresses generated instead of `bcrt1q...`

**Cause**: `address_type` is `p2sh-segwit` instead of `bech32`

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh
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
# → "wsh(sortedmulti(3,[...]xpub.../0/*)...)"  (NOT sh(wsh(...)))
```

### Threshold Mismatch (2-of-3 vs 3-of-3)

**Symptoms**: Transaction broadcasts after only 2 signatures

**Cause**: Descriptor has `sortedmulti(2,...)` instead of `sortedmulti(3,...)`

**Check**:

```bash
# Check multisig threshold in descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep sortedmulti
# Should contain "sortedmulti(3," not "sortedmulti(2,"
```

### Incomplete Signature Error

**Symptoms**: "Incomplete signature" error after 2 signatures

**Cause**: Correct behavior for 3-of-3! Need to add Sign2's signature.

**Solution**: Ensure signing flow includes all 3 wallets (Keygen → Sign1 → Sign2)

### fullpubkey Import Error

**Symptoms**: Error during Multisig setup

**Cause**: fullpubkey format mismatch or ordering issue

**Solution**:

1. Verify correctly exported from `sign1`, `sign2`
2. Check import order to `keygen`
3. Related code: `internal/infrastructure/wallet/key/fullpubkey/`

### Account Config Mismatch

**Symptoms**: Wrong threshold in generated descriptor

**Cause**: Using `account_2of3.yaml` instead of `account_3of3.yaml`

**Solution**:

```bash
# Check account config in script
grep "account.*yaml" scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh
# Should reference account_3of3.yaml
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
```

### Descriptor Check

```bash
# Keygen descriptor list
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch descriptor list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors

# Verify wsh descriptor format with 3-of-3
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep "sortedmulti(3"
```

### PSBT Analysis

```bash
# Check signature status in PSBT
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  analyzepsbt "${psbt_hex}"
# → Check "next" field progression:
#   After 1st sig: "signer" (needs more)
#   After 2nd sig: "signer" (needs more)
#   After 3rd sig: "finalizer" (ready to broadcast)
```

### Signing Progress Check

```bash
# Decode PSBT to see partial signatures
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  decodepsbt "${psbt_hex}"
# Check "partial_signatures" count in each input
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
2. `docs/crypto/btc/e2e_transaction_patterns.md` - Update implementation status (Pattern 7: ❌ → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 7 specific fixes to `P2WSH 3-of-3` related code
- When modifying common code, verify impact on other patterns (especially 6, 8)
- Confirm regression with unit tests when modifying common functions

### 3-of-3 vs 2-of-3 Signing Logic

| Item | 2-of-3 (Pattern 6) | 3-of-3 (Pattern 7) |
|------|--------------------|--------------------|
| Threshold | `sortedmulti(2,...)` | `sortedmulti(3,...)` |
| Required sigs | Any 2 of 3 | All 3 |
| Signing wallets | Keygen + Sign1 | Keygen + Sign1 + Sign2 |
| Early completion | Yes (after 2nd sig) | No (must complete all) |

Ensure script logic correctly handles the 3-of-3 requirement and doesn't skip Sign2.

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh --cleanup

# Full reset (including data)
make btc-e2e-p7-reset
```
