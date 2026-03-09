# MuSig2 User Guide

This guide explains how to use MuSig2 (Simple Two-Round Schnorr Multisignatures) in go-crypto-wallet for creating efficient, private multisig transactions with reduced fees and improved privacy.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [MuSig2 Basics](#musig2-basics)
4. [Transaction Workflows](#transaction-workflows)
5. [File Management](#file-management)
6. [Address Creation](#address-creation)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)
9. [Performance Comparison](#performance-comparison)

---

## Overview

### What is MuSig2?

MuSig2 is a cryptographic protocol that enables multiple parties to create a **single aggregated Schnorr signature** that looks identical to a single-signature transaction on the blockchain. Unlike traditional multisig (P2SH, P2WSH) where multiple signatures are stored on-chain, MuSig2 aggregates multiple signatures into one, providing significant benefits.

### Benefits Over Traditional Multisig

| Feature | Traditional P2WSH Multisig | MuSig2 |
|---------|---------------------------|--------|
| **On-Chain Appearance** | Multiple signatures visible | Single signature (looks like single-sig) |
| **Transaction Size** | ~370-400 bytes (2-of-3) | ~200-250 bytes (30-50% smaller) |
| **Privacy** | Multisig is visible | Indistinguishable from single-sig |
| **Fees** | Higher (proportional to size) | 30-50% lower |
| **Signature Algorithm** | ECDSA | Schnorr (BIP340) |
| **Address Type** | P2WSH (bc1q...) | P2TR Taproot (bc1p...) |
| **Compatibility** | Older standard | Modern (Bitcoin Core 22.0+) |

### When to Use MuSig2

- ✅ **New multisig setups** - Best privacy and efficiency
- ✅ **High-volume operations** - Significant fee savings over time
- ✅ **Privacy-focused applications** - Transactions look like single-sig
- ✅ **Modern infrastructure** - Requires Bitcoin Core 22.0+ and Taproot support
- ⚠️ **Legacy multisig** - Traditional P2WSH still supported for backward compatibility

### How MuSig2 Works (Two-Round Protocol)

```
Round 1: Nonce Generation (Parallel)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Generate Nonce 1                   │
│  Sign Wallet 1 → Generate Nonce 2  (can run in     │
│  Sign Wallet 2 → Generate Nonce 3   parallel)      │
└─────────────────────────────────────────────────────┘
                        ↓
            Exchange nonces via PSBT files
                        ↓
Round 2: Signing (Sequential)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Create Partial Signature 1         │
│  Sign Wallet 1 → Create Partial Signature 2         │
│  Sign Wallet 2 → Create Partial Signature 3         │
└─────────────────────────────────────────────────────┘
                        ↓
            Collect partial signatures
                        ↓
Aggregation (Watch Wallet)
┌─────────────────────────────────────────────────────┐
│  Watch Wallet → Aggregate Partial Signatures        │
│              → Verify Final Signature               │
│              → Broadcast Transaction                │
└─────────────────────────────────────────────────────┘
```

**Key Security Feature:**

- Each wallet generates a **nonce** (random value) in Round 1
- Nonces must be **unique per transaction** and **never reused**
- Reusing nonces can leak private keys - this is critical!

---

## Prerequisites

### System Requirements

- **Watch Wallet** (online): Bitcoin Core node v22.0+ with Taproot support
- **Keygen Wallet** (offline): Isolated system for key generation and first signature
- **Sign Wallet** (offline): Isolated system for additional signatures
- **Go version**: 1.21 or higher
- **Bitcoin Core**: v22.0 or higher (for Taproot/Schnorr support)

### Required Features

MuSig2 builds on top of existing wallet features:

- ✅ **Phase 1**: Taproot support (BIP340 Schnorr signatures)
- ✅ **Phase 2**: PSBT support (BIP174)
- ✅ **Phase 3**: MuSig2 implementation (current)

### Wallet Configuration

Ensure your wallets are properly configured:

```bash
# Watch wallet configuration
config/wallet/btc/watch.yaml

# Keygen wallet configuration
config/wallet/btc/keygen.yaml

# Sign wallet configuration (for multisig)
config/wallet/btc/sign1.yaml
```

**Important**: All wallet commands require the `--config` flag to specify the configuration file:

```bash
# Example usage
./watch --config config/wallet/btc/watch.yaml --coin btc create payment
./keygen --config config/wallet/btc/keygen.yaml --coin btc sign --file tx.psbt
./sign1 --config config/wallet/btc/sign1.yaml --coin btc sign --file tx.psbt
```

See `docs/chains/btc/operation_example.md` for configuration details.

---

## MuSig2 Basics

### MuSig2 Address Format

MuSig2 uses **Taproot (P2TR)** addresses:

```
bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr
└─┬┘└──────────────────────────────────────────────────────────┘
  │                         Taproot address
bc1p = Bech32m Taproot address prefix
```

**Key Properties:**

- Starts with `bc1p` on mainnet, `tb1p` on testnet
- 62 characters long (Bech32m encoding)
- Contains aggregated public key (looks like single-sig)
- Supports Schnorr signatures only

### PSBT Files with MuSig2 Data

MuSig2 transactions use PSBT format with extension fields:

#### File Naming Convention

```
{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
```

**Examples:**

```
payment_15_unsigned_0_1735680000000000000.psbt       # Unsigned (no nonces)
payment_15_nonce_0_1735680000000000001.psbt          # After Round 1 (nonces)
payment_15_unsigned_1_1735680000000000002.psbt       # After partial sig 1
payment_15_unsigned_2_1735680000000000003.psbt       # After partial sig 2
payment_15_signed_3_1735680000000000004.psbt         # Fully signed (aggregated)
```

#### PSBT Extension Fields for MuSig2

MuSig2 data is stored in PSBT proprietary fields:

- **Public Nonces**: Stored per-input (66 bytes per signer)
- **Partial Signatures**: Stored per-input (32 bytes per signer)
- **Aggregated Signature**: Final signature after aggregation (64 bytes)

**Note**: These fields are automatically managed by wallet commands.

### MuSig2 States

A MuSig2 transaction progresses through these states:

```
1. Unsigned (no nonces) ──────> Watch Wallet creates PSBT
                    │
2. Nonces Generated ──────────> Round 1: All wallets generate nonces
                    │
3. Partially Signed (1 sig) ──> Keygen Wallet signs (Round 2)
                    │
4. Partially Signed (2 sig) ──> Sign Wallet 1 signs (Round 2)
                    │
5. Partially Signed (3 sig) ──> Sign Wallet 2 signs (Round 2)
                    │
6. Aggregated & Finalized ────> Watch Wallet aggregates & broadcasts
```

---

## Transaction Workflows

### Workflow 1: Payment Transaction with MuSig2 (3-of-3)

**Scenario**: Sending funds to external addresses using MuSig2 multisig

#### Step 1: Create MuSig2 Address (One-Time Setup)

```bash
# On Keygen wallet - create MuSig2 Taproot addresses
./keygen create musig2-address --account payment
```

**Output:**

```
✓ MuSig2 Taproot addresses created
Account: payment
Addresses created: 10
Address type: P2TR (Taproot)
Multisig type: MuSig2 3-of-3
Example address: tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxq...
```

#### Step 2: Create Unsigned Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
./watch create payment --multisig-type musig2 --fee 0.0001
```

**Output:**

```
Created unsigned MuSig2 transaction: payment_15_unsigned_0_1735680000000000000.psbt
Transaction ID: 15
Inputs: 2
Outputs: 2 (1 recipient + 1 change)
Amount: 0.5 BTC
Fee: 0.0001 BTC
Multisig: MuSig2 3-of-3
Round 1: Ready for nonce generation
```

**File Generated:**

- `data/tx/btc/payment_15_unsigned_0_1735680000000000000.psbt`

#### Step 3: Round 1 - Generate Nonces (All Wallets, Parallel)

**Important**: Nonce generation can be done **in parallel** for all three wallets.

##### Keygen Wallet

```bash
# Transfer PSBT to Keygen wallet
# On offline Keygen wallet system
./keygen musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Keygen
Nonce size: 66 bytes
Output file: payment_15_unsigned_0_1735680000000000001.psbt
Status: Nonce added to PSBT
Next: Transfer to Sign wallets for their nonce generation
```

##### Sign Wallet 1

```bash
# Transfer updated PSBT to Sign wallet 1
# On offline Sign wallet #1 system
./sign musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Sign (auth1)
Nonce size: 66 bytes
Output file: payment_15_unsigned_0_1735680000000000002.psbt
Status: Nonce added to PSBT
Next: Transfer to Sign wallet #2 for nonce generation
```

##### Sign Wallet 2

```bash
# Transfer updated PSBT to Sign wallet 2
# On offline Sign wallet #2 system
./sign musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Sign (auth2)
Nonce size: 66 bytes
Output file: payment_15_nonce_0_1735680000000000003.psbt
Status: All nonces collected (3/3)
Next: Round 2 - Partial signature creation
```

**Note**: After Step 3, all three nonces are stored in the PSBT file.

#### Step 4: Round 2 - Create Partial Signatures (All Wallets, Sequential)

**Important**: Signing must be done **sequentially** after all nonces are collected.

##### Keygen Wallet Signs

```bash
# Transfer PSBT with all nonces to Keygen wallet
# On offline Keygen wallet system
./keygen musig2 sign --file /media/usb/payment_15_nonce_0_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Keygen
Signature size: 32 bytes
Output file: payment_15_unsigned_1_1735680000000000004.psbt
Status: Partial signature 1/3
Next: Transfer to Sign wallet #1 for signing
```

##### Sign Wallet 1 Signs

```bash
# Transfer PSBT to Sign wallet 1
# On offline Sign wallet #1 system
./sign musig2 sign --file /media/usb/payment_15_unsigned_1_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Sign (auth1)
Signature size: 32 bytes
Output file: payment_15_unsigned_2_1735680000000000005.psbt
Status: Partial signature 2/3
Next: Transfer to Sign wallet #2 for signing
```

##### Sign Wallet 2 Signs

```bash
# Transfer PSBT to Sign wallet 2
# On offline Sign wallet #2 system
./sign musig2 sign --file /media/usb/payment_15_unsigned_2_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Sign (auth2)
Signature size: 32 bytes
Output file: payment_15_unsigned_3_1735680000000000006.psbt
Status: All partial signatures collected (3/3)
Next: Transfer to Watch wallet for signature aggregation
```

#### Step 5: Aggregate Signatures (Watch Wallet)

```bash
# Transfer PSBT with all partial signatures to Watch wallet
# On online Watch wallet system
./watch musig2 aggregate --file /media/usb/payment_15_unsigned_3_*.psbt
```

**Output:**

```
✓ MuSig2 signatures aggregated successfully
Final signature size: 64 bytes (Schnorr)
Verification: PASSED
Output file: payment_15_signed_3_1735680000000000007.psbt
Status: Ready for broadcasting
Transaction size: 215 bytes
Traditional multisig size: ~370 bytes
Size reduction: 41.9%
Fee reduction: 41.9%
```

#### Step 6: Broadcast Transaction (Watch Wallet)

```bash
# On online Watch wallet system
./watch send --file payment_15_signed_3_*.psbt
```

**Output:**

```
✓ Transaction broadcast successfully
Transaction hash: a1b2c3d4e5f6789...
Status: Sent
Confirmations: 0 (pending)
On-chain appearance: Single-signature transaction (MuSig2)
Privacy: Maximum (indistinguishable from single-sig)
```

---

### Workflow 2: Deposit Transaction (Single-Signature Taproot Spend)

**Scenario**: Receiving funds from users (client → deposit account) using a single-signature Taproot address.

For deposit transactions using single-signature Taproot addresses:

#### Step 1: Create Unsigned Transaction

```bash
# On Watch wallet
./watch create deposit --multisig-type musig2 --fee 0.0001
```

#### Step 2: Round 1 - Generate Nonce (Keygen Wallet Only)

```bash
# On Keygen wallet
./keygen musig2 nonce --file deposit_8_unsigned_0_*.psbt
```

#### Step 3: Round 2 - Sign (Keygen Wallet Only)

```bash
# On Keygen wallet
./keygen musig2 sign --file deposit_8_nonce_0_*.psbt
```

**Note**: For single-signature, only Keygen wallet participates.

#### Step 4: Finalize and Broadcast

```bash
# On Watch wallet
./watch musig2 aggregate --file deposit_8_unsigned_1_*.psbt
./watch send --file deposit_8_signed_1_*.psbt
```

---

## File Management

### File Locations

#### Watch Wallet

```
data/tx/btc/
├── payment_15_unsigned_0_*.psbt     # Created here (unsigned)
├── payment_15_nonce_0_*.psbt        # Receives after Round 1
├── payment_15_unsigned_3_*.psbt     # Receives after Round 2
└── payment_15_signed_3_*.psbt       # Creates after aggregation
```

#### Keygen Wallet

```
data/tx/btc/
├── payment_15_unsigned_0_...0.psbt     # Receives from Watch
├── payment_15_unsigned_0_...1.psbt     # Creates after nonce gen
├── payment_15_nonce_0_*.psbt           # Receives from Sign wallets
└── payment_15_unsigned_1_*.psbt        # Creates after signing
```

#### Sign Wallet 1

```
data/tx/btc/
├── payment_15_unsigned_0_...1.psbt     # Receives from Keygen
├── payment_15_unsigned_0_...2.psbt     # Creates after nonce gen
├── payment_15_unsigned_1_*.psbt        # Receives from Keygen
└── payment_15_unsigned_2_*.psbt        # Creates after signing
```

#### Sign Wallet 2

```
data/tx/btc/
├── payment_15_unsigned_0_...2.psbt     # Receives from Sign 1
├── payment_15_unsigned_0_...3.psbt     # Creates after nonce gen
├── payment_15_unsigned_2_*.psbt        # Receives from Sign 1
└── payment_15_unsigned_3_*.psbt        # Creates after signing
```

### File Transfer Best Practices

#### For Air-Gapped Systems

1. **Use USB Drives**

   ```bash
   # Mount USB
   mount /dev/sdb1 /media/usb

   # Copy PSBT
   cp data/tx/btc/payment_15_unsigned_0_*.psbt /media/usb/

   # Safely unmount
   umount /media/usb
   ```

2. **Verify File Integrity**

   ```bash
   # Generate checksum
   sha256sum payment_15_unsigned_0_*.psbt > checksum.txt

   # Verify on destination
   sha256sum -c checksum.txt
   ```

3. **Track PSBT State**

   ```bash
   # Check PSBT status
   ./keygen musig2 status --file payment_15_unsigned_0_*.psbt
   ```

#### Security Considerations

- ✅ **Virus scan USB drives** before use on offline systems
- ✅ **Use dedicated USB drives** for wallet operations only
- ✅ **Verify file checksums** after transfer
- ✅ **Track PSBT state** to avoid confusion
- ❌ **Never** connect offline wallets to the internet
- ❌ **Never** reuse nonces (critical security issue)

### Nonce Management

#### Nonce Lifecycle

```
1. Generate → Store in PSBT → Collect all nonces → Sign → Delete
```

**Critical Rules:**

- Each transaction needs fresh nonces
- Nonces are stored securely during Round 1
- Nonces are deleted after signing (Round 2)
- **Never reuse nonces** (will leak private key)

#### Nonce Storage

Nonces are stored in:

- **PSBT proprietary fields** (during transaction flow)
- **Wallet database** (temporary, for Round 2 signing)

**Cleanup**: Nonces are automatically deleted after signing.

---

## Address Creation

### Creating MuSig2 Addresses

#### Command

```bash
# On Keygen wallet
./keygen create musig2-address --account <account> [--count <number>]
```

#### Parameters

- `--account`: Account type (payment, stored, etc.)
- `--count`: Number of addresses to create (default: 10)

#### Example

```bash
# Create 20 MuSig2 addresses for payment account
./keygen create musig2-address --account payment --count 20
```

#### Output

```
✓ MuSig2 Taproot addresses created successfully
Account: payment
Addresses created: 20
Address type: P2TR (Taproot)
Multisig type: MuSig2 3-of-3
Public keys: 3 (Keygen + Sign1 + Sign2)
Sample addresses:
  tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxq...
  tb1pq8r3t5ys7whg3vz2nq4x6kj7h5pq9c8r3t5ys7whg3vz2nq4x6k...
Next: Export addresses to Watch wallet
```

### Address Properties

MuSig2 addresses have these characteristics:

- **Address Type**: P2TR (Taproot)
- **Public Key**: Aggregated from all signers
- **On-Chain**: Looks like single-sig address
- **Privacy**: Maximum (indistinguishable from single-sig)
- **Efficiency**: Smallest possible address type

### Exporting Addresses

```bash
# On Keygen wallet
./keygen export address --account payment --file payment_musig2_addresses.txt

# Transfer to Watch wallet
# On Watch wallet
./watch import address --file payment_musig2_addresses.txt
```

---

## Troubleshooting

### Common Errors

#### Error: "Nonce already used for this transaction"

**Symptom:**

```
Error: nonce reuse detected - cannot sign with same nonce
Transaction ID: 15
```

**Cause:** Attempting to reuse a nonce (critical security violation).

**Solution:**

- **Never reuse nonces** - this will leak your private key
- Generate fresh nonces for each transaction
- If nonces were already generated, create a new transaction

#### Error: "Not all nonces collected"

**Symptom:**

```
Error: cannot proceed to Round 2 - missing nonces
Expected: 3 nonces
Received: 2 nonces
```

**Cause:** Trying to sign before all nonces are generated.

**Solution:**

```bash
# Check PSBT status
./keygen musig2 status --file payment_15_unsigned_0_*.psbt

# Ensure all wallets have generated nonces:
# - Keygen wallet
# - Sign wallet 1
# - Sign wallet 2
```

#### Error: "Invalid partial signature"

**Symptom:**

```
Error: partial signature verification failed
Signer: auth1
```

**Cause:** Partial signature doesn't match expected format or is corrupted.

**Solution:**

```bash
# Verify PSBT integrity
sha256sum payment_15_unsigned_2_*.psbt

# Re-sign with the problematic wallet
./sign musig2 sign --file payment_15_unsigned_1_*.psbt
```

#### Error: "Signature aggregation failed"

**Symptom:**

```
Error: failed to aggregate signatures - verification failed
```

**Cause:** One or more partial signatures are invalid or missing.

**Solution:**

```bash
# Check all partial signatures are present
./watch musig2 status --file payment_15_unsigned_3_*.psbt

# Expected output should show:
# - Nonces: 3/3 ✓
# - Partial signatures: 3/3 ✓

# If any are missing, re-run signing for that wallet
```

#### Error: "MuSig2 not supported by Bitcoin Core"

**Symptom:**

```
Error: Bitcoin Core does not support Taproot/MuSig2
Current version: v21.0
```

**Cause:** Bitcoin Core version is too old.

**Solution:**

- Upgrade Bitcoin Core to v22.0 or higher
- Taproot and Schnorr signatures require Bitcoin Core 22.0+

### Debugging Tips

#### Inspect PSBT Status

```bash
# Check MuSig2 PSBT status
./keygen musig2 status --file payment_15_unsigned_0_*.psbt
```

**Output:**

```
PSBT Status:
  Transaction ID: 15
  Type: payment
  State: Round 1 - Nonce generation
  Nonces collected: 2/3
    ✓ Keygen: [66 bytes]
    ✓ Sign (auth1): [66 bytes]
    ✗ Sign (auth2): missing
  Partial signatures: 0/3
  Next step: Generate nonce for Sign (auth2)
```

#### Verify Nonce Uniqueness

```bash
# Check nonce uniqueness
./keygen musig2 verify-nonces --file payment_15_nonce_0_*.psbt
```

**Output:**

```
✓ All nonces are unique
✓ No nonce reuse detected
✓ Safe to proceed to Round 2
```

#### Check Signature Validity

```bash
# Verify partial signatures
./watch musig2 verify --file payment_15_unsigned_3_*.psbt
```

**Output:**

```
Partial Signature Verification:
  ✓ Keygen signature: valid
  ✓ Sign (auth1) signature: valid
  ✓ Sign (auth2) signature: valid
  ✓ All signatures compatible for aggregation
```

---

## Best Practices

### Security

1. **Nonce Management** (Critical)
   - ✅ Generate fresh nonces for every transaction
   - ✅ Store nonces securely during Round 1-2
   - ✅ Delete nonces immediately after signing
   - ❌ **NEVER reuse nonces** (will leak private key)
   - ❌ Never share secret nonces with anyone

2. **Air-Gapped Signing**
   - ✅ Keep Keygen and Sign wallets offline at all times
   - ✅ Use dedicated, isolated computers for offline wallets
   - ✅ Never connect offline wallets to networks

3. **File Transfer Security**
   - ✅ Use dedicated USB drives for PSBT transfer
   - ✅ Virus scan USB drives before use
   - ✅ Verify file checksums after transfer

4. **Transaction Verification**
   - ✅ Verify transaction amounts before signing
   - ✅ Check recipient addresses carefully
   - ✅ Verify aggregated signature before broadcasting

### Operations

1. **Workflow Management**
   - ✅ Track PSBT state at each step
   - ✅ Use file checksums for integrity verification
   - ✅ Document workflow for each transaction type
   - ✅ Create checklists for Round 1 and Round 2

2. **Parallel vs Sequential Operations**
   - ✅ **Round 1** (nonce generation): Can be done **in parallel**
   - ✅ **Round 2** (signing): Must be done **sequentially** after nonces collected
   - ✅ Collect all nonces before proceeding to Round 2

3. **Testing**
   - ✅ Test on testnet/signet before mainnet
   - ✅ Test with small amounts first
   - ✅ Verify full workflow end-to-end
   - ✅ Test nonce uniqueness enforcement

4. **Monitoring**
   - ✅ Monitor transaction confirmations
   - ✅ Track fee savings compared to traditional multisig
   - ✅ Set up alerts for transaction failures

### Performance

1. **Fee Optimization**
   - ✅ MuSig2 transactions are 30-50% smaller
   - ✅ Fees are proportionally lower
   - ✅ Monitor mempool for optimal fee rates
   - ✅ Consider consolidating UTXOs during low fee periods

2. **Transaction Batching**
   - ✅ Batch multiple payments into single transaction
   - ✅ Further reduces fees per payment
   - ✅ Increases privacy (harder to link payments)

### Backup and Recovery

1. **Regular Backups**
   - ✅ Backup seeds and private keys securely
   - ✅ Backup wallet databases regularly
   - ✅ Store backups in multiple secure locations
   - ✅ Test recovery procedures

2. **Key Management**
   - ✅ Use BIP39 mnemonics for seed backup
   - ✅ Store seed backups offline in secure locations
   - ✅ Consider hardware security modules (HSMs) for production

---

## Performance Comparison

### Transaction Size Comparison

| Multisig Type | Transaction Size | Signature Data | Fee (@ 10 sat/vB) |
|---------------|------------------|----------------|-------------------|
| **Traditional 2-of-3 P2WSH** | ~370 bytes | 2x ECDSA signatures (~142 bytes) | ~3,700 sats |
| **MuSig2 3-of-3 P2TR** | ~215 bytes | 1x Schnorr signature (64 bytes) | ~2,150 sats |
| **Reduction** | **41.9%** | **54.9%** | **41.9%** |

### Real-World Example

**Scenario**: 1000 payment transactions per month

| Metric | Traditional Multisig | MuSig2 | Savings |
|--------|---------------------|--------|---------|
| **Total Size** | 370,000 bytes (~361 KB) | 215,000 bytes (~210 KB) | 155 KB |
| **Total Fees** | 3,700,000 sats (~0.037 BTC) | 2,150,000 sats (~0.0215 BTC) | 0.0155 BTC |
| **Monthly Savings** | - | - | **$775** @ $50k BTC |
| **Annual Savings** | - | - | **$9,300** @ $50k BTC |

### Privacy Benefits

**Traditional Multisig:**

```
On-Chain: Clearly visible as multisig
- Multiple signatures visible
- Number of signers visible
- Redeem script visible
- Easy to track and analyze
```

**MuSig2:**

```
On-Chain: Looks like single-sig transaction
- Single aggregated signature
- Indistinguishable from single-sig
- No visible multisig indicators
- Maximum privacy
```

### Performance Metrics

From integration testing:

- **MuSig2 Signature Size**: 64 bytes (Schnorr)
- **Traditional 2-of-2 Signature Size**: ~142 bytes (2x ECDSA)
- **Size Reduction**: 54.9%
- **Fee Reduction**: 30-50% (proportional to size)
- **Privacy**: Indistinguishable from single-sig on-chain
- **Confirmation Time**: Same as single-sig (no difference)

---

## Additional Resources

### Documentation

- [PSBT User Guide](/chains/btc/psbt/user-guide) - Prerequisite for MuSig2
- [Taproot Guide](/chains/btc/taproot/user-guide) - Understanding Taproot addresses
- [Wallet Flow](/chains/btc/operations/wallet-flow) - Wallet setup and configuration
- [MuSig2 Architecture](/chains/btc/musig2/architecture) - For developers

### Standards and Specifications

- [BIP 327: MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 341: Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP 174: PSBT Specification](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [MuSig2 Paper](https://eprint.iacr.org/2020/1261) - Academic research paper

### Tools

- [Bitcoin Core](https://bitcoincore.org/) - Reference implementation (v22.0+ required)
- [btcd](https://github.com/btcsuite/btcd) - Go Bitcoin implementation (used by go-crypto-wallet)
- [Sparrow Wallet](https://sparrowwallet.com/) - Desktop wallet with Taproot support

### Support

- GitHub Issues: [https://github.com/hiromaily/go-crypto-wallet/issues](https://github.com/hiromaily/go-crypto-wallet/issues)
- Bitcoin Stack Exchange: [https://bitcoin.stackexchange.com/](https://bitcoin.stackexchange.com/)

---

## Glossary

- **MuSig2**: Simple Two-Round Schnorr Multisignatures protocol
- **Schnorr**: Signature algorithm used by Taproot (BIP340)
- **ECDSA**: Elliptic Curve Digital Signature Algorithm (legacy)
- **P2TR**: Pay-to-Taproot address type
- **P2WSH**: Pay-to-Witness-Script-Hash (traditional multisig)
- **Nonce**: Random value used in MuSig2 protocol (must be unique per transaction)
- **Partial Signature**: Individual signature created by each signer in Round 2
- **Signature Aggregation**: Combining partial signatures into single signature
- **Taproot**: Bitcoin upgrade enabling Schnorr signatures (BIP 341)
- **PSBT**: Partially Signed Bitcoin Transaction (BIP 174)
- **Air-gapped**: Computer never connected to networks
- **BIP**: Bitcoin Improvement Proposal

---

**Last Updated**: 2025-12-30
**Version**: 1.0 (MuSig2 Phase 3)
