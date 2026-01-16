# Fix BTC E2E Pattern 9 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 9: P2TR Taproot Single-sig) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)
- @.claude/skills/btc-terminology/SKILL.md - **CRITICAL**: Understand `bech32m` (encoding) vs `taproot` (address_type)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 9 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 9 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2TR (BIP86 Taproot) |
| **Script Type** | Single-sig (Key Path Spend) |
| **Address Format** | `bc1p...` (Mainnet), `tb1p...` (Testnet), `bcrt1p...` (regtest) |
| **Signature Requirement** | Single-sig (1 Schnorr signature) |
| **Descriptor** | `tr([fingerprint/86'/0'/0']xpub.../0/*)` |
| **Required Wallets** | watch, keygen |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="taproot"` |
| **Bitcoin Core Version** | **v22.0+** (Required for Taproot/Schnorr) |

### Comparison with Other Single-sig Patterns

| Item | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) | Pattern 5 (P2WPKH) | **Pattern 9 (P2TR)** |
|------|-------------------|-------------------------|--------------------|--------------------|
| BIP | BIP44 | BIP49 | BIP84 | **BIP86** |
| Address Prefix | `1...`/`m...` | `3...`/`2...` | `bc1q...`/`bcrt1q...` | **`bc1p...`/`bcrt1p...`** |
| Descriptor | `pkh(...)` | `sh(wpkh(...))` | `wpkh(...)` | **`tr(...)`** |
| SegWit Version | N/A | v0 (wrapped) | v0 (native) | **v1 (Taproot)** |
| Signature Algorithm | ECDSA | ECDSA | ECDSA | **Schnorr (BIP340)** |
| Transaction Size | Largest | Medium | Smaller | **Smallest** |
| Legacy Compatible | Yes | Yes | No | **No** |
| address_type | `legacy` | `p2sh-segwit` | `bech32` | **`taproot`** |
| Encoding | Base58Check | Base58Check | Bech32 | **Bech32m** |

### Why Use P2TR (Taproot)?

| Aspect | Description |
|--------|-------------|
| ✅ Smallest transactions | ~30-40% smaller than P2PKH, ~15% smaller than P2WPKH |
| ✅ Lowest fees | Most efficient single-sig format |
| ✅ Enhanced privacy | All Taproot spends look identical on-chain |
| ✅ Schnorr signatures | Faster verification, enables signature aggregation |
| ✅ Future-proof | Latest Bitcoin address standard (activated Nov 2021) |
| ❌ Newer standard | Some older services may not support `bc1p...` addresses |
| ❌ Requires Bitcoin Core v22.0+ | Older nodes cannot validate Taproot transactions |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p9`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 9

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 9 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @docs/crypto/btc/TAPROOT_GUIDE.md - Taproot user guide
- @docs/crypto/btc/e2e_transaction_patterns.md - Pattern 9 details
- @scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh - Pattern 5 script (Native SegWit Single-sig reference)
- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 script (Single-sig base)
- @config/wallet/account.yaml - Single-sig account config

## Pre-check: Environment Variables

**Pattern 9 requires `WALLET_ADDRESS_TYPE="taproot"`.**

> ⚠️ **CRITICAL**: Use `"taproot"` (address type), NOT `"bech32m"` (encoding format).
> See `btc-terminology` skill for details.

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "taproot"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Pre-check: Bitcoin Core Version

**Pattern 9 requires Bitcoin Core v22.0 or later for Taproot/Schnorr support.**

```bash
# Check version
docker exec btc-watch bitcoin-cli --version
# Bitcoin Core RPC client version v22.0.0 or higher required

# Verify Taproot is active (on mainnet/testnet)
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo | grep -A 5 taproot
```

## Implementation Steps

### Step 1: Create Script

Base on Pattern 5 (`e2e-p5-p2wpkh-singlesig.sh`) with these changes:

1. Filename: `e2e-p9-p2tr-singlesig.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="taproot"` (NOT `"bech32m"`)
3. Header comments: Update to Pattern 9 specs
4. Address validation logic: Check for `bcrt1p...` format (regtest Taproot)
5. Descriptor format: Use `tr(...)` instead of `wpkh(...)`

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 9: P2TR Taproot Single-sig
###############################################################################
.PHONY: btc-e2e-p9-reset
btc-e2e-p9-reset:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --reset

.PHONY: btc-e2e-p9
btc-e2e-p9:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh

.PHONY: btc-e2e-p9-verbose
btc-e2e-p9-verbose:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --verbose

.PHONY: btc-e2e-p9-ci
btc-e2e-p9-ci:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --non-interactive

.PHONY: btc-e2e-p9-cleanup
btc-e2e-p9-cleanup:
	./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --cleanup
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p9-reset

# With debug output
./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --verbose
```

> **Note**: For build and verification commands, see common rules.

### Step 4: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Bitcoin Core Version | Bitcoin Core | Must be v22.0+ for Taproot |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | BIP86 HD Key derivation | `internal/application/usecase/keygen/` |
| Descriptor Export | Taproot Descriptor | `internal/infrastructure/wallet/key/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT + Schnorr signing | `internal/infrastructure/wallet/api/btc/` |

## Technical Specification: P2TR (Taproot)

### What is Taproot?

Taproot is Bitcoin's major upgrade (activated November 2021) introducing:

| Component | BIP | Description |
|-----------|-----|-------------|
| **Taproot** | BIP341 | New output type combining Schnorr + MAST |
| **Tapscript** | BIP342 | Script semantics for Taproot spends |
| **Schnorr Signatures** | BIP340 | New signature scheme (64 bytes vs 71-72) |
| **Key Derivation** | BIP86 | Derivation path for Taproot key path spend |

### Address Structure

```
P2TR Address (SegWit v1 - Taproot):
┌─────────────────────────────────────────────────────────────┐
│  Bech32m Encoding                                            │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 1                                      │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Tweaked Public Key (32 bytes, x-only)               │ │ │
│  │  │  = internal_key + tweak * G                          │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1p + 58 characters (mainnet)
        tb1p + 58 characters (testnet/signet)
        bcrt1p + 58 characters (regtest)
```

### Key Path vs Script Path Spending

Pattern 9 uses **Key Path Spend** (the simpler path):

| Spend Type | Description | Use Case |
|------------|-------------|----------|
| **Key Path** (Pattern 9) | Direct signature with tweaked key | Single-sig, simple multisig |
| Script Path | Reveal script + satisfy conditions | Complex conditions, timelocks |

### Descriptor Format (BIP86)

```
tr([fingerprint/86'/coin'/account']xpub.../0/*)
   └─ BIP86 derivation path (Taproot)

Example (regtest):
tr([a1b2c3d4/86'/1'/0']tpub.../0/*)
```

### Key Derivation Path (BIP86)

```
m / 86' / coin_type' / account' / change / address_index
         └─ 0' for mainnet, 1' for testnet/regtest
```

### Schnorr Signature (BIP340)

```
Signature: 64 bytes (r || s)
  - r: 32 bytes (x-coordinate of R point)
  - s: 32 bytes (scalar)

ECDSA Signature: 71-72 bytes (DER encoded)
  - 30 (type) + length + 02 + r_len + r + 02 + s_len + s
```

### Signing Flow (Single-sig)

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (sign with Schnorr signature)
    ↓
Watch Wallet (broadcast)
```

## Pattern 9 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 9 specific errors:

### Bitcoin Core Version Error

**Symptoms**: `Unknown address type` or `Taproot not supported`

**Cause**: Bitcoin Core version is older than v22.0

**Solution**:

```bash
# Check version
docker exec btc-watch bitcoin-cli --version

# If older than v22.0, update Docker image
# Edit compose.btc.yaml to use newer Bitcoin Core image
```

### address_type Mismatch

**Symptoms**: `bcrt1q...` (bech32) addresses generated instead of `bcrt1p...` (bech32m encoded)

**Cause**: `address_type` is not `taproot`

> ⚠️ **CRITICAL**: `address_type` must be `"taproot"`, NOT `"bech32m"`.
> `bech32m` is the encoding format, not the address type. See `btc-terminology` skill.

**Solution**:

```bash
# Check environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "taproot"

# Check script setting
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh
```

### Bech32 vs Bech32m Confusion

**Critical**: P2TR uses **Bech32m** encoding (different from P2WPKH's Bech32)

| Address Type | Encoding | Witness Version | Prefix (regtest) |
|--------------|----------|-----------------|------------------|
| P2WPKH (Pattern 5) | Bech32 | v0 | `bcrt1q...` |
| **P2TR (Pattern 9)** | **Bech32m** | **v1** | **`bcrt1p...`** |

**Symptoms**: Checksum error or invalid address

**Solution**: Verify Bech32m encoding is used for Taproot addresses

### Descriptor Format Error

**Symptoms**: Error during Descriptor export/import

**Cause**: Using wrong format (e.g., `wpkh(...)` instead of `tr(...)`)

**Check**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Expected format
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "tr([...]xpub.../0/*)"  (NOT wpkh or pkh)
```

### key_type Auto-derivation Check

**Check**: Verify `key_type` correctly derived from `address_type`

| address_type | Expected key_type |
|--------------|-------------------|
| `taproot` | `bip86` |

Related code: `AddrType.ToKeyType()` in `internal/domain/address/types.go`

### Schnorr Signing Error

**Symptoms**: Signature verification failed or invalid signature length

**Cause**: ECDSA signature used instead of Schnorr, or wrong signing method

**Check**:

```bash
# Verify signing method in PSBT
docker exec btc-watch bitcoin-cli -regtest analyzepsbt "${psbt_hex}"
```

**Solution**: Ensure Taproot signing path uses BIP340 Schnorr signatures

### Tweaked Key Error

**Symptoms**: `tweaked public key mismatch` or similar

**Cause**: Internal key not properly tweaked for key path spend

**Note**: BIP86 defines key path spend without script tree:

```
output_key = internal_key + H_TapTweak(internal_key) * G
```

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Check Taproot activation (should be active in regtest)
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo | jq '.softforks.taproot'

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2TR)
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": false, "iswitness": true, "witness_version": 1
```

### Descriptor Check

```bash
# Keygen descriptor list
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch descriptor list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors

# Verify tr descriptor format
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep "tr("
```

### Address Derivation Test

```bash
# Derive addresses from descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  deriveaddresses "tr([fingerprint/86'/1'/0']xpub.../0/*)" "[0,5]"
# Should return bcrt1p... addresses
```

### Signature Verification

```bash
# Check transaction details
docker exec btc-watch bitcoin-cli -regtest \
  decoderawtransaction "<raw_tx_hex>"

# Verify witness data (should be 64 bytes for Schnorr)
# witness: ["<64-byte-schnorr-signature>"]
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

### Taproot-Specific Code Locations

| Path | Role |
|------|------|
| Taproot address generation | `btcutil.NewAddressTaproot()` |
| BIP86 key derivation | `m/86'/...` path handling |
| Schnorr signing | BIP340 signature creation |
| Tweaked key calculation | Internal key tweaking for key path |

## Documentation Updates

After creating script, update these documents:

1. `scripts/operation/btc/e2e/README.md` - Add to script list
2. `docs/crypto/btc/e2e_transaction_patterns.md` - Update implementation status (Pattern 9: ❌ → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list
4. `make/btc_e2e.mk` - Add Makefile targets

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 9 specific fixes to `P2TR Taproot` related code
- When modifying common code, verify impact on other patterns (especially 1, 3, 5)
- Confirm regression with unit tests when modifying common functions

### Bech32 vs Bech32m (Critical)

> See `btc-terminology` skill for complete explanation.

- **Pattern 5 (P2WPKH)**: Uses **Bech32** encoding (SegWit v0, prefix `bcrt1q`)
- **Pattern 9 (P2TR)**: Uses **Bech32m** encoding (SegWit v1, prefix `bcrt1p`)
- These are **different encodings** - do not confuse them
- Bech32m has a different checksum constant than Bech32
- **CRITICAL**: `address_type` config uses `"taproot"`, NOT `"bech32m"`

### Bitcoin Core Version Requirement

- Pattern 9 **requires Bitcoin Core v22.0+**
- Older versions cannot create or validate Taproot transactions
- Verify Docker image version before starting E2E test

> **Note**: For build rules, security, see common rules.

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh --cleanup

# Full reset (including data)
make btc-e2e-p9-reset
```
