# Fix BTC E2E Pattern 5 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 5: P2WPKH Native SegWit Single-sig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 5 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 5 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2WPKH (BIP84 Native SegWit) |
| **Script Type** | Single-sig |
| **Address Format** | `bc1q...` (Mainnet), `tb1q...` (Testnet), `bcrt1q...` (regtest) |
| **Signature Requirement** | Single-sig (1 signature) |
| **Descriptor** | `wpkh([fingerprint/84'/0'/0']xpub.../0/*)` |
| **Required Wallets** | watch, keygen |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="bech32"` |

### Comparison with Other Single-sig Patterns

| Item | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) | Pattern 5 (P2WPKH) | Pattern 9 (P2TR) |
|------|-------------------|-------------------------|--------------------|--------------------|
| BIP | BIP44 | BIP49 | **BIP84** | BIP86 |
| Address Prefix | `1...`/`m...` | `3...`/`2...` | **`bc1q...`/`bcrt1q...`** | `bc1p...` |
| Descriptor | `pkh(...)` | `sh(wpkh(...))` | **`wpkh(...)`** | `tr(...)` |
| SegWit | No | Yes (wrapped) | **Yes (native)** | Yes (Taproot) |
| Transaction Size | Largest | Medium | **Smaller** | Smallest |
| Legacy Compatible | Yes | Yes | **No** | No |
| address_type | `legacy` | `p2sh-segwit` | **`bech32`** | `taproot` |

### Why Use P2WPKH (Native SegWit)?

| Aspect | Description |
|--------|-------------|
| ✅ Most efficient SegWit | Smallest transaction size among SegWit v0 types |
| ✅ Lower fees | ~40% smaller than P2PKH, ~15% smaller than P2SH-P2WPKH |
| ✅ Bech32 encoding | Error detection and correction capabilities |
| ✅ Widely adopted | Supported by most modern wallets and exchanges |
| ❌ Legacy incompatible | Very old wallets may not recognize `bc1q...` addresses |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p5`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 5

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 5 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 script (Single-sig base)
- @scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh - Pattern 3 script (SegWit Single-sig reference)
- @config/wallet/account/account.yaml - Single-sig account config

## Pre-check: Environment Variables

**Pattern 5 requires `WALLET_ADDRESS_TYPE="bech32"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Implementation Steps

### Step 1: Create Script

Base on Pattern 1 (`e2e-p1-p2pkh-singlesig.sh`) or Pattern 3 (`e2e-p3-p2sh-p2wpkh-singlesig.sh`) with these changes:

1. Filename: `e2e-p5-p2wpkh-singlesig.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="bech32"`
3. Header comments: Update to Pattern 5 specs
4. Address validation logic: Check for `bcrt1q...` format (regtest bech32)

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 5: P2WPKH Native SegWit Single-sig
###############################################################################
.PHONY: btc-e2e-p5-reset
btc-e2e-p5-reset:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --reset

.PHONY: btc-e2e-p5
btc-e2e-p5:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh

.PHONY: btc-e2e-p5-verbose
btc-e2e-p5-verbose:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --verbose

.PHONY: btc-e2e-p5-ci
btc-e2e-p5-ci:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --non-interactive

.PHONY: btc-e2e-p5-cleanup
btc-e2e-p5-cleanup:
	./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --cleanup
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p5-reset

# With debug output
./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --verbose
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
| Descriptor Export | BIP84 Descriptor | `internal/infrastructure/wallet/key/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

## Technical Specification: P2WPKH (Native SegWit)

### Address Structure

```
P2WPKH Address (Native SegWit v0):
┌─────────────────────────────────────────────────────────────┐
│  Bech32 Encoding                                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 0                                      │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Public Key Hash (20 bytes)                          │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1q + 38 characters (mainnet)
        bcrt1q + 38 characters (regtest)
```

### Descriptor Format

```
wpkh([fingerprint/84'/0'/0']xpub.../0/*)
     └─ BIP84 derivation path (Native SegWit)
```

### Key Derivation Path (BIP84)

```
m / 84' / coin_type' / account' / change / address_index
         └─ 0' for mainnet, 1' for testnet/regtest
```

### Signing Flow (Single-sig)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (sign with single key)
    ↓
Watch Wallet (broadcast)
```

## Pattern 5 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 5 specific errors:

### address_type Mismatch

**Symptoms**: `m...`, `n...`, `2...`, or `3...` addresses generated instead of `bcrt1q...`

**Cause**: `address_type` is not `bech32`

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "bech32"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh
```

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using wrong format (e.g., `pkh(...)` or `sh(wpkh(...))`)

**Check**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Expected format
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "wpkh([...]xpub.../0/*)"  (NOT pkh or sh(wpkh))
```

### key_type Auto-derivation Check

**Check**: Verify `key_type` correctly derived from `address_type`

| address_type | Expected key_type |
|--------------|-------------------|
| `bech32` | `bip84` |

Related code: `AddrType.ToKeyType()` in `internal/domain/address/types.go`

### Bech32 Address Validation

**Symptoms**: Address format validation error

**Cause**: Regex or validation logic doesn't recognize `bcrt1q...` format

**Check**:

```bash
# Verify address format in regtest
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "address": "bcrt1q...", "iswitness": true, "witness_version": 0
```

### Bitcoin Core Version Compatibility

**Note**: Native SegWit requires Bitcoin Core 0.16.0+ (2018).
This should not be an issue with modern versions but check if using older images.

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2WPKH)
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": false, "iswitness": true, "witness_version": 0, "witness_program": "<20-byte-hash>"
```

### Descriptor Check

```bash
# Keygen descriptor list
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch descriptor list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors

# Verify wpkh descriptor format
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep wpkh
```

### Address Derivation Test

```bash
# Derive addresses from descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  deriveaddresses "wpkh([fingerprint/84'/1'/0']xpub.../0/*)" "[0,5]"
# Should return bcrt1q... addresses
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
2. `docs/crypto/btc/operations/e2e-transaction-patterns.md` - Update implementation status (Pattern 5: 🔶 → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 5 specific fixes to `P2WPKH Single-sig` related code
- When modifying common code, verify impact on other patterns (especially 1, 3, 9)
- Confirm regression with unit tests when modifying common functions

### Bech32 vs Bech32m

- Pattern 5 (P2WPKH): Uses **Bech32** encoding (SegWit v0)
- Pattern 9 (P2TR): Uses **Bech32m** encoding (SegWit v1/Taproot)
- Do not confuse these formats

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh --cleanup

# Full reset (including data)
make btc-e2e-p5-reset
```
