# Fix BTC E2E Pattern 10 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 10: P2TR MuSig2 N-of-N) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)
- @.claude/skills/btc-terminology/SKILL.md - **CRITICAL**: Understand `bech32m` (encoding) vs `taproot` (address_type)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 10 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 10 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2TR (BIP86 Taproot) |
| **Script Type** | MuSig2 (N-of-N Signature Aggregation) |
| **Address Format** | `bc1p...` (Mainnet), `tb1p...` (Testnet), `bcrt1p...` (regtest) |
| **Signature Requirement** | N-of-N (all signers required, aggregated into single signature) |
| **Descriptor** | `tr(musig([fp1/86'/1'/1']xpub1,[fp2/86'/1'/1']xpub2,[fp3/86'/1'/1']xpub3)/0/*)` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="taproot"` |
| **Bitcoin Core Version** | **v22.0+** (Required for Taproot/Schnorr) |
| **Protocol** | **2-Round MuSig2** (BIP327) |

### Comparison with Related Patterns

| Item | Pattern 7 (P2WSH 3-of-3) | Pattern 9 (P2TR Single-sig) | **Pattern 10 (P2TR MuSig2)** |
|------|--------------------------|----------------------------|------------------------------|
| BIP | BIP84 | BIP86 | **BIP86 + BIP327** |
| Address | `bcrt1q...` | `bcrt1p...` | **`bcrt1p...`** |
| Descriptor | `wsh(sortedmulti(3,...))` | `tr(...)` | **`tr(musig(...))`** |
| Signers | 3 | 1 | **N (all required)** |
| On-Chain Signatures | 3 ECDSA | 1 Schnorr | **1 Schnorr (aggregated)** |
| Transaction Size | ~370-400 vBytes | ~99 vBytes | **~99 vBytes** |
| Privacy | Multisig visible | Single-sig | **Single-sig appearance** |
| Signature Algorithm | ECDSA | Schnorr (BIP340) | **Schnorr (BIP340)** |

### Why Use MuSig2?

| Aspect | Description |
|--------|-------------|
| ✅ Maximum privacy | N-of-N multisig looks identical to single-sig on-chain |
| ✅ Smallest size | Same size as single-sig (~30-50% smaller than P2WSH multisig) |
| ✅ Lowest fees | Aggregated signature = minimal witness data |
| ✅ Schnorr benefits | Batch verification, linear signature aggregation |
| ✅ Future-proof | Modern cryptographic standard |
| ❌ Interactive | Requires 2-round communication between all signers |
| ❌ All must sign | N-of-N only (no threshold like 2-of-3) |
| ❌ Nonce management | Critical: nonce reuse leaks private keys |
| ❌ Complex implementation | More complex than traditional multisig |

### MuSig2 Protocol Overview

```
Round 1: Nonce Generation (can be parallelized)
┌─────────────────────────────────────────────────────────┐
│  Keygen Wallet → Generate Nonce 1 (66 bytes)            │
│  Sign1 Wallet  → Generate Nonce 2 (66 bytes)            │
│  Sign2 Wallet  → Generate Nonce 3 (66 bytes)            │
└─────────────────────────────────────────────────────────┘
                        ↓
            Exchange nonces via PSBT files
                        ↓
Round 2: Signing (sequential, after all nonces collected)
┌─────────────────────────────────────────────────────────┐
│  Keygen Wallet → Create Partial Signature 1 (32 bytes)  │
│  Sign1 Wallet  → Create Partial Signature 2 (32 bytes)  │
│  Sign2 Wallet  → Create Partial Signature 3 (32 bytes)  │
└─────────────────────────────────────────────────────────┘
                        ↓
            Collect partial signatures
                        ↓
Aggregation (Watch Wallet)
┌─────────────────────────────────────────────────────────┐
│  Watch Wallet → Aggregate Partial Signatures            │
│              → Final Schnorr Signature (64 bytes)       │
│              → Broadcast Transaction                    │
└─────────────────────────────────────────────────────────┘
```

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p10`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 10

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 10 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @docs/crypto/btc/musig2/user-guide.md - MuSig2 user guide (essential reading)
- @docs/crypto/btc/taproot/user-guide.md - Taproot user guide
- @docs/crypto/btc/operations/e2e-transaction-patterns.md - Pattern 10 details
- @scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh - Pattern 9 script (Taproot reference)
- @scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh - Pattern 7 script (3-of-3 multisig reference)
- @config/wallet/account/account_3of3.yaml - 3-of-3 account config

## Pre-check: Environment Variables

**Pattern 10 requires `WALLET_ADDRESS_TYPE="taproot"`.**

> ⚠️ **CRITICAL**: Use `"taproot"` (address type), NOT `"bech32m"` (encoding format).
> See `btc-terminology` skill for details.

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "taproot"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Pre-check: Bitcoin Core Version

**Pattern 10 requires Bitcoin Core v22.0 or later for Taproot/Schnorr support.**

```bash
# Check version
docker exec btc-watch bitcoin-cli --version
# Bitcoin Core RPC client version v22.0.0 or higher required

# Verify Taproot is active
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo | grep -A 5 taproot
```

## Implementation Steps

### Step 1: Create Script

Base on Pattern 9 (`e2e-p9-p2tr-singlesig.sh`) and Pattern 7 (`e2e-p7-p2wsh-3of3.sh`) with these changes:

1. Filename: `e2e-p10-p2tr-musig2.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="taproot"` (NOT `"bech32m"`)
3. Header comments: Update to Pattern 10 specs
4. Address validation logic: Check for `bcrt1p...` format (regtest Taproot)
5. Descriptor format: Use `tr(musig(...))` for aggregated key
6. **Round 1**: Implement nonce generation for all wallets
7. **Round 2**: Implement partial signature creation
8. **Aggregation**: Implement signature aggregation in Watch wallet
9. Required wallets: watch, keygen, sign1, sign2
10. Account config: Use `account_3of3.yaml` (N-of-N)

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 10: P2TR MuSig2 N-of-N
###############################################################################
.PHONY: btc-e2e-p10-reset
btc-e2e-p10-reset:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --reset

.PHONY: btc-e2e-p10
btc-e2e-p10:
	make btc-e2e P=10

.PHONY: btc-e2e-p10-verbose
btc-e2e-p10-verbose:
	make btc-e2e-verbose P=10

.PHONY: btc-e2e-p10-ci
btc-e2e-p10-ci:
	./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh --non-interactive

.PHONY: btc-e2e-p10-cleanup
btc-e2e-p10-cleanup:
	make btc-e2e-cleanup P=10
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-reset P=10

# With debug output
make btc-e2e-verbose P=10
```

> **Note**: For build and verification commands, see common rules.

### Step 4: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Bitcoin Core Version | Bitcoin Core | Must be v22.0+ for Taproot |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | BIP86 HD Key derivation | `internal/application/usecase/keygen/` |
| MuSig2 Key Aggregation | Aggregated public key | `internal/infrastructure/wallet/key/musig2/` |
| Descriptor Export | MuSig2 Taproot Descriptor | `internal/infrastructure/wallet/key/descriptor/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| **Round 1 - Nonce** | Nonce generation | MuSig2 nonce generation (critical!) |
| **Round 2 - Sign** | Partial signatures | MuSig2 partial signature creation |
| **Aggregation** | Signature aggregation | MuSig2 signature aggregation |
| Broadcast | Transaction broadcast | Final transaction broadcast |

## Technical Specification: MuSig2

### What is MuSig2?

MuSig2 (Simple Two-Round Schnorr Multisignatures) is defined in **BIP327**. It enables N-of-N multisig with a **single aggregated Schnorr signature** that is indistinguishable from single-sig on-chain.

| Component | BIP | Description |
|-----------|-----|-------------|
| **MuSig2 Protocol** | BIP327 | Two-round signature aggregation protocol |
| **Schnorr Signatures** | BIP340 | 64-byte signatures (vs 71-72 for ECDSA) |
| **Taproot** | BIP341 | Output type for Schnorr signatures |
| **Key Derivation** | BIP86 | Derivation path for Taproot key path spend |

### MuSig2 vs Traditional Multisig

| Feature | Traditional P2WSH Multisig | **MuSig2** |
|---------|---------------------------|-----------|
| On-Chain Appearance | Multiple signatures visible | **Single signature** |
| Transaction Size | ~370-400 bytes (2-of-3) | **~200-250 bytes** |
| Privacy | Multisig is visible | **Indistinguishable from single-sig** |
| Fees | Higher | **30-50% lower** |
| Signature Algorithm | ECDSA | **Schnorr (BIP340)** |
| Address Type | P2WSH (`bc1q...`) | **P2TR (`bc1p...`)** |
| Threshold | M-of-N (flexible) | **N-of-N only** |

### Address Structure

```
MuSig2 P2TR Address:
┌─────────────────────────────────────────────────────────────┐
│  Bech32m Encoding                                            │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 1 (Taproot)                            │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Aggregated Public Key (32 bytes, x-only)            │ │ │
│  │  │  = MuSig2.KeyAgg(pk1, pk2, pk3)                      │ │ │
│  │  │  (Looks identical to single-sig Taproot address)     │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1p + 58 characters (mainnet)
        bcrt1p + 58 characters (regtest)
```

### Descriptor Format

```
tr(musig([fp1/86'/1'/1']xpub1,[fp2/86'/1'/1']xpub2,[fp3/86'/1'/1']xpub3)/0/*)
   └─ MuSig2 aggregated key from 3 signers
```

### Key Derivation Path (BIP86 for Multisig)

```
m / 86' / coin_type' / account' / change / address_index
         └─ 1' for testnet/regtest
                      └─ 1' for multisig accounts
```

### Two-Round Protocol Details

#### Round 1: Nonce Generation

```
Each signer generates:
- Secret nonce: (secnonce1, secnonce2) - 64 bytes total, MUST be kept secret
- Public nonce: (R1, R2) - 66 bytes total, shared with others

CRITICAL: Nonces MUST be:
- Generated fresh for EACH transaction
- NEVER reused (reuse = private key leak)
- Securely deleted after signing
```

#### Round 2: Partial Signing

```
After collecting all public nonces:
- Aggregate nonces: R = R1_agg + R2_agg
- Create partial signature: s_i = k_i + e * a_i * x_i
- Partial signature size: 32 bytes per signer
```

#### Aggregation

```
Combine partial signatures:
- s = s_1 + s_2 + s_3
- Final signature: (R, s) - 64 bytes total
- Verify: s*G == R + e*P (where P is aggregated pubkey)
```

### Signing Flow (3-of-3 MuSig2)

```
Watch Wallet (create unsigned PSBT)
    ↓
┌───────────────────────────────────────────────────┐
│ ROUND 1: Nonce Generation (parallel)              │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│ │   Keygen    │ │   Sign1     │ │   Sign2     │   │
│ │ Gen Nonce 1 │ │ Gen Nonce 2 │ │ Gen Nonce 3 │   │
│ └─────────────┘ └─────────────┘ └─────────────┘   │
└───────────────────────────────────────────────────┘
    ↓ Exchange nonces via PSBT
┌───────────────────────────────────────────────────┐
│ ROUND 2: Partial Signing (sequential)             │
│   Keygen → Sign1 → Sign2                          │
│   (each adds partial signature to PSBT)           │
└───────────────────────────────────────────────────┘
    ↓
Watch Wallet (aggregate signatures → broadcast)

※ ALL signers required in both rounds
```

## Pattern 10 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 10 specific errors:

### Bitcoin Core Version Error

**Symptoms**: `Unknown address type` or `Taproot not supported`

**Cause**: Bitcoin Core version is older than v22.0

**Solution**:

```bash
# Check version
docker exec btc-watch bitcoin-cli --version

# If older than v22.0, update Docker image
```

### Nonce Reuse Error (CRITICAL)

**Symptoms**: `nonce reuse detected - cannot sign with same nonce`

**Cause**: Attempting to reuse a nonce (CRITICAL SECURITY VIOLATION)

**Solution**:

```bash
# NEVER reuse nonces - this WILL leak private keys
# Generate fresh nonces for each transaction
# If nonces were already used, create a NEW transaction

# Check nonce status
./keygen musig2 status --file <psbt_file>
```

> **WARNING**: Nonce reuse in MuSig2 is a **critical security vulnerability** that can leak private keys. Always generate fresh nonces for each transaction.

### Missing Nonces Error

**Symptoms**: `cannot proceed to Round 2 - missing nonces`

**Cause**: Trying to sign before all nonces are generated

**Solution**:

```bash
# Check PSBT status
./keygen musig2 status --file <psbt_file>

# Expected output:
# Nonces collected: 3/3 ✓

# Ensure ALL wallets have generated nonces:
# - Keygen wallet
# - Sign1 wallet
# - Sign2 wallet
```

### Invalid Partial Signature Error

**Symptoms**: `partial signature verification failed`

**Cause**: Partial signature doesn't match expected format or is corrupted

**Solution**:

```bash
# Verify PSBT integrity
sha256sum <psbt_file>

# Verify partial signatures
./watch musig2 verify --file <psbt_file>
```

### Signature Aggregation Failed

**Symptoms**: `failed to aggregate signatures - verification failed`

**Cause**: One or more partial signatures are invalid or missing

**Solution**:

```bash
# Check all partial signatures are present
./watch musig2 status --file <psbt_file>

# Expected:
# - Nonces: 3/3 ✓
# - Partial signatures: 3/3 ✓
```

### Bech32 vs Bech32m Confusion

**Symptoms**: Checksum error or invalid address

**Cause**: Using Bech32 encoding instead of Bech32m for Taproot addresses

**Solution**: Verify Bech32m encoding is used:

| Address Type | Encoding | Witness Version | Prefix (regtest) |
|--------------|----------|-----------------|------------------|
| P2WPKH/P2WSH | Bech32 | v0 | `bcrt1q...` |
| **P2TR (MuSig2)** | **Bech32m** | **v1** | **`bcrt1p...`** |

### MuSig2 Key Aggregation Error

**Symptoms**: `failed to aggregate public keys` or invalid aggregated key

**Cause**: Public key format mismatch or ordering issue

**Solution**:

```bash
# Verify public keys are in correct format (32-byte x-only)
# Check key ordering is consistent across all wallets
```

### PSBT State Confusion

**Symptoms**: Commands fail due to unexpected PSBT state

**Cause**: Running commands out of order

**Solution**:

```bash
# Check PSBT state
./keygen musig2 status --file <psbt_file>

# Expected progression:
# 1. Unsigned (no nonces) → Watch creates PSBT
# 2. Nonces Generated → Round 1 complete
# 3. Partially Signed (1 sig) → Keygen signed
# 4. Partially Signed (2 sig) → Sign1 signed
# 5. Partially Signed (3 sig) → Sign2 signed
# 6. Aggregated & Finalized → Watch aggregated
```

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Check Taproot activation
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo | jq '.softforks.taproot'

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Check address info (verify P2TR MuSig2)
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

# Verify MuSig2 descriptor format
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep "musig"
```

### MuSig2 Status Check

```bash
# Check MuSig2 PSBT status
./keygen musig2 status --file <psbt_file>

# Expected output:
# PSBT Status:
#   Transaction ID: 15
#   Type: payment
#   State: Round 2 - Signing
#   Nonces collected: 3/3 ✓
#   Partial signatures: 2/3
#   Next step: Sign with Sign2 wallet
```

### Nonce Verification

```bash
# Check nonce uniqueness (critical!)
./keygen musig2 verify-nonces --file <psbt_file>

# Expected:
# ✓ All nonces are unique
# ✓ No nonce reuse detected
# ✓ Safe to proceed to Round 2
```

### Signature Verification

```bash
# Verify partial signatures
./watch musig2 verify --file <psbt_file>

# Expected:
# Partial Signature Verification:
#   ✓ Keygen signature: valid
#   ✓ Sign1 signature: valid
#   ✓ Sign2 signature: valid
#   ✓ All signatures compatible for aggregation
```

### Transaction Analysis

```bash
# Check final transaction (after aggregation)
docker exec btc-watch bitcoin-cli -regtest \
  decoderawtransaction "<raw_tx_hex>"

# Verify witness data (should be single 64-byte Schnorr signature)
# witness: ["<64-byte-schnorr-signature>"]
```

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/key/descriptor/` | Descriptor processing |
| `internal/infrastructure/wallet/key/musig2/` | MuSig2 key aggregation |
| `internal/domain/address/types.go` | address_type → key_type conversion |
| `pkg/config/loader.go` | Config loader |

### MuSig2-Specific Code Locations

| Path | Role |
|------|------|
| MuSig2 key aggregation | `MuSig2.KeyAgg()` |
| Nonce generation | `MuSig2.NonceGen()` |
| Partial signing | `MuSig2.Sign()` |
| Signature aggregation | `MuSig2.PartialSigAgg()` |
| BIP86 key derivation | `m/86'/...` path handling |
| Schnorr verification | BIP340 signature verification |

## Security Considerations

### Nonce Management (CRITICAL)

| Rule | Description |
|------|-------------|
| ✅ Fresh nonces | Generate new nonces for EVERY transaction |
| ✅ Secure storage | Store secret nonces securely during protocol |
| ✅ Immediate deletion | Delete secret nonces after signing |
| ❌ **NEVER reuse** | Nonce reuse = **private key leak** |
| ❌ **NEVER share secrets** | Only share public nonces |

### Air-Gapped Operations

| Rule | Description |
|------|-------------|
| ✅ Offline signing | Keep Keygen, Sign1, Sign2 offline |
| ✅ USB transfer | Use dedicated USB drives for PSBT files |
| ✅ Verify checksums | Verify file integrity after transfer |
| ❌ **NEVER** online | Never connect signing wallets to network |

## Documentation Updates

After creating script, update these documents:

1. `scripts/operation/btc/e2e/README.md` - Add to script list
2. `docs/crypto/btc/operations/e2e-transaction-patterns.md` - Update implementation status (Pattern 10: 🔜 → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list
4. `make/btc_e2e.mk` - Add Makefile targets

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 10 specific fixes to `MuSig2` related code
- When modifying common code, verify impact on other patterns (especially 9)
- Confirm regression with unit tests when modifying common functions

### MuSig2 vs Single-sig Taproot (Critical Difference)

| Item | Pattern 9 (P2TR Single-sig) | **Pattern 10 (P2TR MuSig2)** |
|------|----------------------------|------------------------------|
| Signers | 1 | **N (all required)** |
| Protocol | Simple Schnorr sign | **2-Round MuSig2** |
| Nonces | Not applicable | **Critical (fresh per tx)** |
| Key | Single public key | **Aggregated public key** |
| Descriptor | `tr([fp/86'/...]xpub)` | **`tr(musig([fp1]xpub1,[fp2]xpub2,...))`** |

### 2-Round Protocol Requirements

| Round | Action | Parallelizable | Required Wallets |
|-------|--------|----------------|------------------|
| **Round 1** | Nonce generation | **Yes** | All (Keygen, Sign1, Sign2) |
| **Round 2** | Partial signing | **No** (sequential) | All (Keygen, Sign1, Sign2) |
| **Aggregation** | Signature combination | N/A | Watch only |

> **Important**: Round 2 can only start AFTER all nonces from Round 1 are collected.

### Bitcoin Core Version Requirement

- Pattern 10 **requires Bitcoin Core v22.0+**
- Older versions cannot create or validate Taproot transactions
- Verify Docker image version before starting E2E test

> **Note**: For build rules, security, see common rules.

## Performance Comparison

| Metric | P2WSH 3-of-3 Multisig | **MuSig2 3-of-3** | Improvement |
|--------|----------------------|-------------------|-------------|
| Transaction Size | ~550 bytes | **~200 bytes** | **~64% smaller** |
| Signature Data | 3 × ECDSA (~213 bytes) | **1 × Schnorr (64 bytes)** | **~70% smaller** |
| On-Chain Privacy | Multisig visible | **Single-sig appearance** | **Maximum** |
| Fees (@ 10 sat/vB) | ~5,500 sats | **~2,000 sats** | **~64% lower** |

## Cleanup

```bash
# Stop containers only
make btc-e2e-cleanup P=10

# Full reset (including data)
make btc-e2e-reset P=10
```
