# Fix BTC E2E Pattern 11 Test #{issue_number}

Implement and fix BTC E2E test (Pattern 11: P2TR Tapscript M-of-N) in **regtest environment**.

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)
- @.claude/skills/btc-terminology/SKILL.md - **CRITICAL**: Understand `bech32m` (encoding) vs `taproot` (address_type)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command creates/runs `scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh` and analyzes/fixes any errors.

> **Note**: This E2E test runs in local regtest (Regression Test) environment.
> It does not connect to actual Bitcoin network (mainnet/testnet).

### Pattern 11 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 11 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2TR (BIP86 Taproot) |
| **Script Type** | Tapscript (M-of-N Script Path Spend) |
| **Address Format** | `bc1p...` (Mainnet), `tb1p...` (Testnet), `bcrt1p...` (regtest) |
| **Signature Requirement** | M-of-N threshold (e.g., 2-of-3) |
| **Descriptor** | `tr(internal_key,sortedmulti_a(2,[fp1]xpub1,[fp2]xpub2,[fp3]xpub3))` |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="taproot"` |
| **Bitcoin Core Version** | **v22.0+** (Required for Taproot/Tapscript) |
| **Spending Path** | **Script Path** (not Key Path) |

### Key Path vs Script Path Spending

| Spend Type | Description | Use Case | Pattern |
|------------|-------------|----------|---------|
| **Key Path** | Direct signature with tweaked key | Single-sig, efficient N-of-N | Pattern 9 |
| **MuSig2** | N-of-N aggregated signature | All-party required | Pattern 10 |
| **Script Path** | Reveal Merkle proof + satisfy script | **M-of-N threshold, complex conditions** | **Pattern 11** |

### Comparison with Related Patterns

| Item | Pattern 6 (P2WSH 2-of-3) | Pattern 10 (MuSig2) | **Pattern 11 (Tapscript)** |
|------|--------------------------|---------------------|----------------------------|
| BIP | BIP84 | BIP86 + BIP327 | **BIP86 + BIP341 + BIP342** |
| Address | `bcrt1q...` (62 chars) | `bcrt1p...` | **`bcrt1p...`** |
| Descriptor | `wsh(sortedmulti(2,...))` | `tr(musig(...))` | **`tr(internal_key,sortedmulti_a(M,...))`** |
| Threshold | M-of-N | N-of-N only | **M-of-N (flexible)** |
| On-Chain | Multiple ECDSA sigs | 1 Schnorr sig | **M Schnorr sigs + script** |
| Privacy | Script visible | Maximum | **Good (unused paths hidden)** |
| Transaction Size | ~370-400 vBytes | ~99 vBytes | **~150-200 vBytes** |

### Why Use Tapscript (Script Path)?

| Aspect | Description |
|--------|-------------|
| ✅ M-of-N threshold | True threshold signing (2-of-3, 3-of-5, etc.) |
| ✅ Script flexibility | Multiple spending conditions in Merkle tree |
| ✅ Unused paths hidden | Only reveal the script path you use |
| ✅ Schnorr signatures | Each signature is 64 bytes (vs 71-72 ECDSA) |
| ✅ Smaller than P2WSH | More efficient than legacy multisig |
| ❌ Larger than MuSig2 | Script path requires script reveal |
| ❌ Requires MAST | More complex implementation |
| ❌ Multiple signatures | M signatures on-chain (not aggregated) |

### Tapscript Architecture

```
Taproot Output (P2TR):
┌─────────────────────────────────────────────────────────────┐
│  Output Key = Internal Key + H_TapTweak(internal_key, merkle_root) * G │
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  Internal Key (can be used for Key Path spend)          ││
│  └─────────────────────────────────────────────────────────┘│
│                           +                                 │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  Merkle Root (Script Tree)                               ││
│  │  ┌─────────────────────────────────────────────────────┐││
│  │  │         Root                                        │││
│  │  │        /    \                                       │││
│  │  │    Leaf1   Leaf2                                    │││
│  │  │   (2-of-3) (timelock)                               │││
│  │  └─────────────────────────────────────────────────────┘││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘

Script Path Spend:
1. Reveal internal key
2. Reveal script (e.g., 2-of-3 multisig)
3. Provide Merkle proof to the script
4. Satisfy script conditions (M Schnorr signatures)
```

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p11`
- **Commit type**: `feat(btc)` (for new script) / `fix(btc)` (for fixes)
- **Scope**: BTC E2E Pattern 11

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Implement/fix locally without creating branch or PR.

## Pattern 11 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @docs/crypto/btc/TAPROOT_GUIDE.md - Taproot user guide
- @docs/crypto/btc/btc_bch_technical_guide.md - Script Path details
- @docs/crypto/btc/e2e_transaction_patterns.md - Pattern 11 details
- @scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh - Pattern 9 script (Taproot Key Path reference)
- @scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh - Pattern 10 script (MuSig2 reference)
- @scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh - Pattern 4 script (2-of-3 multisig reference)
- @config/wallet/account_2of3.yaml - 2-of-3 account config

## Pre-check: Environment Variables

**Pattern 11 requires `WALLET_ADDRESS_TYPE="taproot"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "taproot"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Pre-check: Bitcoin Core Version

**Pattern 11 requires Bitcoin Core v22.0 or later for Taproot/Tapscript support.**

```bash
# Check version
docker exec btc-watch bitcoin-cli --version
# Bitcoin Core RPC client version v22.0.0 or higher required

# Verify Taproot is active
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo | grep -A 5 taproot
```

## Implementation Steps

### Step 1: Create Script

Base on Pattern 9 (`e2e-p9-p2tr-singlesig.sh`) and Pattern 4 (`e2e-p4-p2sh-p2wsh-2of3.sh`) with these changes:

1. Filename: `e2e-p11-p2tr-tapscript.sh`
2. Environment variable: `WALLET_ADDRESS_TYPE="taproot"`
3. Header comments: Update to Pattern 11 specs
4. Address validation logic: Check for `bcrt1p...` format (regtest Taproot)
5. Descriptor format: Use `tr(internal_key,sortedmulti_a(M,...))` for Script Path
6. Script tree: Build Merkle tree with M-of-N script
7. Signing flow: M signatures required (Schnorr)
8. Required wallets: watch, keygen, sign1, sign2
9. Account config: Use `account_2of3.yaml`

### Step 2: Add Makefile Targets

Add to `make/btc_e2e.mk`:

```makefile
###############################################################################
# E2E Testing - Pattern 11: P2TR Tapscript M-of-N
###############################################################################
.PHONY: btc-e2e-p11-reset
btc-e2e-p11-reset:
 ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --reset

.PHONY: btc-e2e-p11
btc-e2e-p11:
 ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh

.PHONY: btc-e2e-p11-verbose
btc-e2e-p11-verbose:
 ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --verbose

.PHONY: btc-e2e-p11-ci
btc-e2e-p11-ci:
 ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --non-interactive

.PHONY: btc-e2e-p11-cleanup
btc-e2e-p11-cleanup:
 ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --cleanup
```

### Step 3: Run E2E Test

```bash
# Full reset and run (recommended)
make btc-e2e-p11-reset

# With debug output
./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --verbose
```

> **Note**: For build and verification commands, see common rules.

### Step 4: Error Analysis

Identify the phase where error occurred and investigate related code:

| Phase | Related Code | Description |
|-------|--------------|-------------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Bitcoin Core Version | Bitcoin Core | Must be v22.0+ for Taproot/Tapscript |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | BIP86 HD Key derivation | `internal/application/usecase/keygen/` |
| **Script Tree** | Merkle tree construction | Tapscript leaf building |
| **Internal Key** | Tweak calculation | Internal key + Merkle root tweak |
| Descriptor Export | Tapscript Descriptor | `internal/infrastructure/wallet/key/descriptor/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| **Script Path Signing** | M Schnorr signatures | Script reveal + signatures |
| Broadcast | Transaction broadcast | Final transaction broadcast |

## Technical Specification: Tapscript

### What is Tapscript?

Tapscript (BIP342) defines the script semantics for Taproot script path spending. It allows embedding multiple spending conditions in a Merkle tree, with only the used condition revealed on-chain.

| Component | BIP | Description |
|-----------|-----|-------------|
| **Taproot** | BIP341 | Output structure, spending rules |
| **Tapscript** | BIP342 | Script semantics for script path |
| **Schnorr Signatures** | BIP340 | 64-byte signatures |
| **Key Derivation** | BIP86 | Derivation path for Taproot |

### Tapscript vs Traditional Multisig

| Feature | P2WSH Multisig | **Tapscript M-of-N** |
|---------|----------------|---------------------|
| Script Visibility | Always visible | **Only used path visible** |
| Signature Type | ECDSA (71-72 bytes) | **Schnorr (64 bytes)** |
| Address Type | `bc1q...` | **`bc1p...`** |
| Unused Scripts | Visible on-chain | **Hidden in Merkle tree** |
| Multiple Conditions | Single script | **Multiple scripts in tree** |

### Address Structure

```
Tapscript P2TR Address:
┌─────────────────────────────────────────────────────────────┐
│  Bech32m Encoding                                            │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 1 (Taproot)                            │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Output Key (32 bytes, x-only)                       │ │ │
│  │  │  = internal_key + H_TapTweak(P, merkle_root) * G     │ │ │
│  │  │  (Commits to both internal key AND script tree)      │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1p + 58 characters (mainnet)
        bcrt1p + 58 characters (regtest)
```

### Descriptor Format

```
tr(internal_key,sortedmulti_a(2,[fp1/86'/1'/1']xpub1/0/*,[fp2/86'/1'/1']xpub2/0/*,[fp3/86'/1'/1']xpub3/0/*))
   │            └─ Tapscript leaf: 2-of-3 multisig with Schnorr
   └─ Internal key (can be NUMS point for script-only)

NUMS (Nothing Up My Sleeve) point:
- Use when you want ONLY script path spending
- No one can spend via key path
- Example: H (hash of "Nothing Up My Sleeve")
```

### Merkle Tree Structure

```
Script Tree Example (2-of-3 + timelock fallback):
           ┌───────────────────┐
           │    Merkle Root     │
           └───────────────────┘
                  /    \
                 /      \
    ┌───────────────┐  ┌───────────────┐
    │   Leaf A      │  │   Leaf B      │
    │ 2-of-3 Multi  │  │   Timelock    │
    │ (Schnorr sigs)│  │ + 1-of-3      │
    └───────────────┘  └───────────────┘

Spending via Leaf A:
1. Reveal internal key
2. Reveal Leaf A script (2-of-3 multisig)
3. Provide Merkle proof (hash of Leaf B)
4. Provide 2 Schnorr signatures
```

### Key Derivation Path (BIP86)

```
m / 86' / coin_type' / account' / change / address_index
         └─ 1' for testnet/regtest
                      └─ 1' for multisig accounts
```

### Signing Flow (2-of-3 Tapscript)

```
Watch Wallet (create unsigned PSBT)
    ↓
┌───────────────────────────────────────────────────┐
│ Construct Script Path Spend:                      │
│ - Include control block (internal key + path)     │
│ - Include script (2-of-3 multisig)                │
│ - Include Merkle proof (sibling hashes)           │
└───────────────────────────────────────────────────┘
    ↓
Keygen Wallet (1st Schnorr signature)
    ↓
Sign1 Wallet (2nd Schnorr signature) ← Complete here! (2-of-3)
    ↓
Watch Wallet (finalize + broadcast)

※ Sign2 not required for 2-of-3
※ All signatures are Schnorr (BIP340)
```

## Pattern 11 Specific Errors

For common errors (No utxo, RPC connection, etc.), see common rules. Below are Pattern 11 specific errors:

### Bitcoin Core Version Error

**Symptoms**: `Unknown address type` or `Taproot not supported`

**Cause**: Bitcoin Core version is older than v22.0

**Solution**:

```bash
# Check version
docker exec btc-watch bitcoin-cli --version

# If older than v22.0, update Docker image
```

### Invalid Tapscript Descriptor

**Symptoms**: `Invalid descriptor` or `Unsupported script type`

**Cause**: Wrong descriptor format for Tapscript

**Solution**:

```bash
# Check descriptor format
# Correct: tr(internal_key,sortedmulti_a(2,...))
# Wrong:   tr(sortedmulti(2,...))  # missing internal key
# Wrong:   wsh(sortedmulti(2,...)) # wrong output type

# Verify descriptor
docker exec btc-watch bitcoin-cli -regtest \
  getdescriptorinfo "tr(internal_key,sortedmulti_a(2,...))"
```

### Control Block Error

**Symptoms**: `Invalid control block` or `Merkle proof verification failed`

**Cause**: Incorrect control block construction

**Solution**:

```bash
# Control block structure:
# - 1 byte: leaf version + parity bit
# - 32 bytes: internal key (x-only)
# - 32 bytes * N: Merkle path (sibling hashes)

# Verify in PSBT
docker exec btc-watch bitcoin-cli -regtest \
  decodepsbt "<psbt_hex>" | jq '.inputs[].tap_internal_key'
```

### Script Execution Error

**Symptoms**: `Script evaluation failed` or `OP_CHECKSIGADD failed`

**Cause**: Incorrect script or signature format

**Solution**:

```bash
# Tapscript uses different opcodes:
# - OP_CHECKSIGADD (replaces OP_CHECKMULTISIG)
# - Each signature is validated individually

# Check script format
# Correct: <pk1> OP_CHECKSIG <pk2> OP_CHECKSIGADD <pk3> OP_CHECKSIGADD 2 OP_EQUAL
# Wrong:   OP_2 <pk1> <pk2> <pk3> OP_3 OP_CHECKMULTISIG (legacy format)
```

### Signature Format Error

**Symptoms**: `Invalid signature` or `Signature must be 64 bytes`

**Cause**: Using ECDSA signature instead of Schnorr

**Solution**:

```bash
# Tapscript requires Schnorr signatures (BIP340)
# - Exactly 64 bytes (r || s)
# - No SIGHASH byte appended by default

# Check signature in witness
docker exec btc-watch bitcoin-cli -regtest \
  decoderawtransaction "<raw_tx>" | jq '.vin[0].txinwitness'
# Each signature should be 128 hex chars (64 bytes)
```

### Bech32 vs Bech32m Confusion

**Symptoms**: Checksum error or invalid address

**Cause**: Using Bech32 encoding instead of Bech32m for Taproot addresses

| Address Type | Encoding | Witness Version | Prefix (regtest) |
|--------------|----------|-----------------|------------------|
| P2WPKH/P2WSH | Bech32 | v0 | `bcrt1q...` |
| **P2TR (Tapscript)** | **Bech32m** | **v1** | **`bcrt1p...`** |

### Threshold Mismatch

**Symptoms**: Transaction broadcasts after wrong number of signatures

**Cause**: Descriptor has wrong threshold (e.g., M=3 instead of M=2)

**Solution**:

```bash
# Check threshold in descriptor
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep sortedmulti_a
# Should contain "sortedmulti_a(2," for 2-of-3
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

# Check address info (verify P2TR with script)
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

# Verify Tapscript descriptor format
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors | jq '.descriptors[].desc' | grep "tr("
```

### Script Tree Analysis

```bash
# Analyze Taproot output
docker exec btc-watch bitcoin-cli -regtest \
  getaddressinfo "<bc1p_address>" | jq '.witness_program'

# Decode PSBT to see Taproot fields
docker exec btc-watch bitcoin-cli -regtest \
  decodepsbt "<psbt_hex>" | jq '.inputs[0] | {
    tap_internal_key,
    tap_merkle_root,
    tap_leaf_script,
    tap_bip32_derivs
  }'
```

### PSBT Taproot Fields

```bash
# Check Taproot-specific PSBT fields
docker exec btc-watch bitcoin-cli -regtest \
  decodepsbt "<psbt_hex>"

# Expected Taproot fields:
# - PSBT_IN_TAP_KEY_SIG: Key path signature (not used in script path)
# - PSBT_IN_TAP_SCRIPT_SIG: Script path signatures
# - PSBT_IN_TAP_LEAF_SCRIPT: Tapscript leaf being spent
# - PSBT_IN_TAP_BIP32_DERIVATION: BIP32 derivation info
# - PSBT_IN_TAP_INTERNAL_KEY: Internal key
# - PSBT_IN_TAP_MERKLE_ROOT: Merkle root of script tree
```

### Transaction Analysis

```bash
# Decode final transaction
docker exec btc-watch bitcoin-cli -regtest \
  decoderawtransaction "<raw_tx_hex>"

# Check witness structure for Script Path spend:
# witness: [
#   <signature_1>,    # 64 bytes Schnorr
#   <signature_2>,    # 64 bytes Schnorr
#   <script>,         # The Tapscript being executed
#   <control_block>   # Internal key + Merkle path
# ]
```

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/infrastructure/wallet/key/descriptor/` | Descriptor processing |
| `internal/infrastructure/wallet/key/tapscript/` | Tapscript processing |
| `internal/domain/address/types.go` | address_type → key_type conversion |
| `pkg/config/loader.go` | Config loader |

### Tapscript-Specific Code Locations

| Path | Role |
|------|------|
| Script tree building | `TapscriptTree.Build()` |
| Merkle root calculation | `TapscriptTree.MerkleRoot()` |
| Control block construction | `ControlBlock.Serialize()` |
| Internal key tweaking | `ComputeTweakedOutputKey()` |
| Leaf hash calculation | `TapLeafHash()` |
| Schnorr signing | BIP340 signature creation |

## Security Considerations

### Script Path Privacy

| Aspect | Benefit |
|--------|---------|
| ✅ Hidden conditions | Unused script paths remain secret |
| ✅ Smaller witness | Only reveal one path |
| ⚠️ Used path revealed | The spent script is visible on-chain |

### Key Management

| Rule | Description |
|------|-------------|
| ✅ Secure internal key | Protect internal key for key path fallback |
| ✅ NUMS point option | Use if you want script-only spending |
| ✅ Merkle tree privacy | Only hash of unused scripts revealed |

### Signature Security

| Rule | Description |
|------|-------------|
| ✅ Schnorr signatures | All signatures must be BIP340 Schnorr |
| ✅ No signature malleability | Schnorr signatures are non-malleable |
| ✅ Proper SIGHASH | Use appropriate sighash flags |

## Documentation Updates

After creating script, update these documents:

1. `scripts/operation/btc/e2e/README.md` - Add to script list
2. `docs/crypto/btc/e2e_transaction_patterns.md` - Update implementation status (Pattern 11: 🔜 → ✅)
3. `.claude/rules/btc/e2e-script.md` - Add to pattern list
4. `make/btc_e2e.mk` - Add Makefile targets

## Cautions

### Avoid Impact on Other Patterns

- Limit Pattern 11 specific fixes to `Tapscript` related code
- When modifying common code, verify impact on other patterns (especially 9, 10)
- Confirm regression with unit tests when modifying common functions

### Tapscript vs Key Path vs MuSig2 (Critical Differences)

| Item | Pattern 9 (Key Path) | Pattern 10 (MuSig2) | **Pattern 11 (Tapscript)** |
|------|---------------------|---------------------|---------------------------|
| Signers | 1 | N (all required) | **M-of-N (flexible)** |
| Protocol | Simple Schnorr sign | 2-Round MuSig2 | **Script path reveal + M sigs** |
| On-chain sigs | 1 (64 bytes) | 1 (64 bytes) | **M (64 bytes each)** |
| Script revealed | None | None | **Yes (used script)** |
| Descriptor | `tr([fp]xpub)` | `tr(musig(...))` | **`tr(internal_key,sortedmulti_a(M,...))`** |

### Script Tree Complexity

| Scenario | Recommendation |
|----------|----------------|
| Simple 2-of-3 | Single leaf with `sortedmulti_a(2,...)` |
| 2-of-3 + timelock | Two leaves (normal + timelock fallback) |
| Complex conditions | Multiple leaves with different scripts |

### Bitcoin Core Version Requirement

- Pattern 11 **requires Bitcoin Core v22.0+**
- Older versions cannot create or validate Taproot/Tapscript transactions
- Verify Docker image version before starting E2E test

> **Note**: For build rules, security, see common rules.

## Performance Comparison

| Metric | P2WSH 2-of-3 | MuSig2 3-of-3 | **Tapscript 2-of-3** |
|--------|--------------|---------------|---------------------|
| Transaction Size | ~370 vBytes | ~99 vBytes | **~150-200 vBytes** |
| Signatures | 2 × ECDSA (~142 bytes) | 1 × Schnorr (64 bytes) | **2 × Schnorr (128 bytes)** |
| Script Visible | Yes | No | **Partial (used path only)** |
| Threshold | M-of-N | N-of-N only | **M-of-N** |
| Fees (@ 10 sat/vB) | ~3,700 sats | ~990 sats | **~1,500-2,000 sats** |

## Cleanup

```bash
# Stop containers only
./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh --cleanup

# Full reset (including data)
make btc-e2e-p11-reset
```
