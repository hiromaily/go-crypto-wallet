# PSBT User Guide

This guide explains how to use Partially Signed Bitcoin Transactions (PSBT) in go-crypto-wallet for creating, signing, and broadcasting Bitcoin transactions.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [PSBT Basics](#psbt-basics)
4. [Transaction Workflows](#transaction-workflows)
5. [File Management](#file-management)
6. [Address Types](#address-types)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)

---

## Overview

### What is PSBT?

PSBT (Partially Signed Bitcoin Transaction) is a Bitcoin standard (BIP 174) that provides a structured format for unsigned and partially signed transactions. It allows multiple parties to collaborate on signing a transaction without exposing private keys.

### Benefits Over Legacy CSV Format

| Feature | Legacy CSV | PSBT |
|---------|-----------|------|
| **Format** | Custom CSV | Standardized (BIP 174) |
| **Encoding** | Plain text | Base64 |
| **Compatibility** | go-crypto-wallet only | Bitcoin Core, hardware wallets, other tools |
| **Metadata** | Limited | Rich (UTXO info, derivation paths, etc.) |
| **Security** | Basic | Enhanced with structured validation |
| **Error Handling** | Limited | Comprehensive |
| **File Extension** | Various | `.psbt` |

### When to Use PSBT

- ✅ **All new transactions** - PSBT is the recommended format
- ✅ **Multisig transactions** - Better tracking of signatures
- ✅ **Offline signing** - Secure air-gapped wallet operations
- ✅ **Hardware wallet integration** - Future compatibility
- ⚠️ **Legacy transactions** - CSV format still supported for backward compatibility

---

## Prerequisites

### System Requirements

- **Watch Wallet** (online): Bitcoin Core node or compatible RPC endpoint
- **Keygen Wallet** (offline): Isolated system for key generation and first signature
- **Sign Wallet** (offline): Isolated system for additional signatures
- **Go version**: 1.21 or higher
- **Bitcoin Core**: v22.0 or higher (for PSBT support)

### Wallet Configuration

Ensure your wallets are properly configured:

```bash
# Watch wallet configuration
data/config/btc_watch.toml

# Keygen wallet configuration
data/config/btc_keygen.toml

# Sign wallet configuration (for multisig)
data/config/btc_sign.toml
```

See `docs/crypto/btc/operation_example.md` for configuration details.

---

## PSBT Basics

### PSBT File Format

#### File Naming Convention

PSBT files follow this naming pattern:

```
{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
```

**Components:**
- `actionType`: Transaction type (`deposit`, `payment`, `transfer`)
- `txID`: Database transaction ID
- `txType`: Status (`unsigned`, `signed`)
- `signedCount`: Number of signatures collected (0, 1, 2, ...)
- `timestamp`: Unix timestamp in nanoseconds

**Examples:**
```
deposit_8_unsigned_0_1534744535097796209.psbt    # Unsigned deposit (0 signatures)
deposit_8_unsigned_1_1534744535097796210.psbt    # Partially signed (1 signature)
deposit_8_signed_2_1534744535097796211.psbt      # Fully signed (2 signatures)
```

#### File Content

PSBT files contain base64-encoded binary data:

```
cHNidP8BAHECAAAAAZt/TvyKa6hVH3n8FwUPKA...
```

**Do not edit PSBT files manually!** Use wallet commands to create and sign PSBTs.

### PSBT States

A PSBT progresses through these states:

```
1. Unsigned (0 signatures) ───> Watch Wallet creates
                    │
2. Partially Signed (1 sig) ──> Keygen Wallet signs
                    │
3. Partially Signed (2 sig) ──> Sign Wallet signs (multisig)
                    │
4. Fully Signed ──────────────> Watch Wallet finalizes & broadcasts
```

---

## Transaction Workflows

### Workflow 1: Deposit Transaction (Single-Signature)

**Scenario**: Receiving funds from users (client → deposit account)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
cd /path/to/go-crypto-wallet
./watch create deposit --fee 0.0001
```

**Output:**
```
Created unsigned transaction: deposit_8_unsigned_0_1534744535097796209.psbt
Transaction ID: 8
Inputs: 5
Outputs: 2
Amount: 1.5 BTC
Fee: 0.0001 BTC
```

**File Generated:**
- `data/tx/btc/deposit_8_unsigned_0_1534744535097796209.psbt`

#### Step 2: Transfer PSBT to Keygen Wallet

Transfer the unsigned PSBT file to the offline Keygen wallet system:

```bash
# Using USB drive or secure file transfer
cp data/tx/btc/deposit_8_unsigned_0_*.psbt /media/usb/
```

#### Step 3: Sign Transaction (Keygen Wallet)

```bash
# On the offline Keygen wallet system
cd /path/to/go-crypto-wallet
./keygen sign --file /media/usb/deposit_8_unsigned_0_1534744535097796209.psbt
```

**Output:**
```
Signed transaction successfully
Input file: deposit_8_unsigned_0_1534744535097796209.psbt
Output file: deposit_8_signed_1_1534744535097796210.psbt
Is fully signed: true
Transaction ready for broadcasting
```

**File Generated:**
- `data/tx/btc/deposit_8_signed_1_1534744535097796210.psbt`

#### Step 4: Transfer Signed PSBT Back to Watch Wallet

```bash
# Copy signed PSBT back to Watch wallet
cp data/tx/btc/deposit_8_signed_1_*.psbt /media/usb/
```

#### Step 5: Broadcast Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
./watch send --file /media/usb/deposit_8_signed_1_1534744535097796210.psbt
```

**Output:**
```
Transaction broadcast successfully
Transaction hash: a1b2c3d4e5f6...
Status: Sent
```

---

### Workflow 2: Payment Transaction (Multisig 2-of-2)

**Scenario**: Sending funds to external addresses (payment account → external)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
./watch create payment --fee 0.0002
```

**Output:**
```
Created unsigned transaction: payment_12_unsigned_0_1534744600000000000.psbt
Transaction ID: 12
Inputs: 3
Outputs: 2 (1 recipient + 1 change)
Amount: 0.5 BTC
Fee: 0.0002 BTC
Multisig: 2-of-2 (requires 2 signatures)
```

#### Step 2: First Signature (Keygen Wallet)

```bash
# Transfer to Keygen wallet
./keygen sign --file payment_12_unsigned_0_1534744600000000000.psbt
```

**Output:**
```
Signed transaction successfully
Output file: payment_12_unsigned_1_1534744600000000001.psbt
Is fully signed: false
Signatures: 1/2
Next: Transfer to Sign wallet for second signature
```

**Note:** Transaction is **not** fully signed yet (1 of 2 signatures).

#### Step 3: Second Signature (Sign Wallet)

```bash
# Transfer to Sign wallet
./sign sign --file payment_12_unsigned_1_1534744600000000001.psbt
```

**Output:**
```
Signed transaction successfully
Output file: payment_12_signed_2_1534744600000000002.psbt
Is fully signed: true
Signatures: 2/2
Transaction ready for broadcasting
```

#### Step 4: Broadcast Transaction (Watch Wallet)

```bash
./watch send --file payment_12_signed_2_1534744600000000002.psbt
```

**Output:**
```
Transaction broadcast successfully
Transaction hash: f6e5d4c3b2a1...
Status: Sent
```

---

### Workflow 3: Transfer Transaction (Multisig 2-of-2)

**Scenario**: Moving funds between internal accounts (stored → payment)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
./watch create transfer --sender stored --receiver payment --amount 10.0 --fee 0.0003
```

**Output:**
```
Created unsigned transaction: transfer_15_unsigned_0_1534744700000000000.psbt
Transaction ID: 15
Sender: stored account
Receiver: payment account
Amount: 10.0 BTC
Fee: 0.0003 BTC
Multisig: 2-of-2
```

#### Steps 2-4: Same as Payment Workflow

Follow the same signing and broadcasting steps as the payment transaction:

1. Keygen wallet signs (first signature)
2. Sign wallet signs (second signature)
3. Watch wallet broadcasts (finalization)

---

## File Management

### File Locations

#### Watch Wallet

```
data/tx/btc/
├── deposit_8_unsigned_0_*.psbt      # Created here
├── deposit_8_signed_1_*.psbt        # Receives from Keygen
├── payment_12_unsigned_0_*.psbt     # Created here
└── payment_12_signed_2_*.psbt       # Receives from Sign
```

#### Keygen Wallet

```
data/tx/btc/
├── deposit_8_unsigned_0_*.psbt      # Receives from Watch
├── deposit_8_signed_1_*.psbt        # Creates here
├── payment_12_unsigned_0_*.psbt     # Receives from Watch
└── payment_12_unsigned_1_*.psbt     # Creates here
```

#### Sign Wallet

```
data/tx/btc/
├── payment_12_unsigned_1_*.psbt     # Receives from Keygen
└── payment_12_signed_2_*.psbt       # Creates here
```

### File Transfer Best Practices

#### For Air-Gapped Systems

1. **Use USB Drives**
   ```bash
   # Mount USB
   mount /dev/sdb1 /media/usb

   # Copy PSBT
   cp data/tx/btc/payment_12_unsigned_0_*.psbt /media/usb/

   # Safely unmount
   umount /media/usb
   ```

2. **Use QR Codes** (for smaller transactions)
   ```bash
   # Generate QR code
   qrencode -o psbt.png < payment_12_unsigned_0_*.psbt

   # Scan and decode on offline system
   zbarimg psbt.png > payment_12_unsigned_0_*.psbt
   ```

3. **Use Optical Data Transfer** (most secure)
   - Print PSBT as QR code or text
   - Manually type or scan on offline system

#### Security Considerations

- ✅ **Virus scan USB drives** before use on offline systems
- ✅ **Use dedicated USB drives** for wallet operations only
- ✅ **Verify file integrity** with checksums (sha256sum)
- ❌ **Never** connect offline wallets to the internet
- ❌ **Never** use the same USB drive for other purposes

### File Cleanup

#### Automatic Cleanup

Wallets automatically manage PSBT files:

- Unsigned PSBTs are kept until signed
- Partially signed PSBTs are kept until fully signed
- Fully signed PSBTs are kept until broadcasted
- Broadcasted transaction PSBTs can be archived

#### Manual Cleanup

```bash
# Archive old PSBT files
mkdir -p data/tx/btc/archive/2025-01
mv data/tx/btc/*_signed_*_1735*.psbt data/tx/btc/archive/2025-01/

# Compress archives
tar -czf psbt-archive-2025-01.tar.gz data/tx/btc/archive/2025-01/
```

---

## Address Types

PSBT supports all Bitcoin address types:

### Supported Address Types

| Type | Format | Example | Description |
|------|--------|---------|-------------|
| **P2PKH** | Legacy | `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa` | Original Bitcoin address |
| **P2SH-SegWit** | Nested SegWit | `3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy` | SegWit wrapped in P2SH |
| **P2WPKH** | Bech32 | `bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4` | Native SegWit |
| **P2TR** | Taproot | `bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr` | Taproot (BIP86) |

### Address Type Features

#### P2PKH (Legacy)

```bash
# Create transaction with P2PKH addresses
./watch create deposit --address-type p2pkh
```

- **Pros**: Universal compatibility, well-tested
- **Cons**: Larger transaction size, higher fees
- **Signature**: ECDSA
- **Use Case**: Maximum compatibility with old wallets

#### P2SH-SegWit (Nested SegWit)

```bash
./watch create deposit --address-type p2sh-segwit
```

- **Pros**: SegWit benefits with legacy compatibility
- **Cons**: Larger than native SegWit
- **Signature**: ECDSA
- **Use Case**: Transition period, maximum compatibility

#### P2WPKH (Bech32 Native SegWit)

```bash
./watch create deposit --address-type p2wpkh
```

- **Pros**: Smallest size, lowest fees, error detection
- **Cons**: Not universally supported by old wallets
- **Signature**: ECDSA
- **Use Case**: Modern wallets, cost optimization

#### P2TR (Taproot)

```bash
./watch create deposit --address-type taproot
```

- **Pros**: Best privacy, Schnorr signatures, script flexibility, lowest fees
- **Cons**: Requires Bitcoin Core 22.0+, newest standard
- **Signature**: Schnorr (BIP340)
- **Use Case**: Maximum privacy and efficiency

### Mixed Address Types

PSBT handles transactions with mixed input types seamlessly:

```
Transaction with mixed inputs:
├── Input 1: P2PKH (legacy)       → ECDSA signature
├── Input 2: P2WPKH (SegWit)      → ECDSA signature
└── Input 3: P2TR (Taproot)       → Schnorr signature

PSBT automatically:
✅ Uses correct signature algorithm per input
✅ Validates each input independently
✅ Combines all signatures correctly
```

---

## Troubleshooting

### Common Errors

#### Error: "PSBT is not fully signed"

**Symptom:**
```
Error: PSBT is not fully signed - cannot finalize incomplete PSBT
```

**Cause:** Trying to broadcast a PSBT that doesn't have all required signatures.

**Solution:**
```bash
# Check signature status
./keygen sign --file payment_12_unsigned_1_*.psbt --dry-run

# Ensure all required signatures are collected:
# - Single-sig: 1 signature (Keygen)
# - 2-of-2 multisig: 2 signatures (Keygen + Sign)
```

#### Error: "Invalid PSBT format"

**Symptom:**
```
Error: failed to validate PSBT: invalid PSBT format
```

**Cause:** PSBT file is corrupted or invalid base64.

**Solution:**
```bash
# Verify file integrity
sha256sum payment_12_unsigned_0_*.psbt

# Re-create transaction if file is corrupted
./watch create payment --fee 0.0002

# Do not edit PSBT files manually
```

#### Error: "Transaction already broadcast"

**Symptom:**
```
Error: transaction already sent
Transaction ID: (empty)
```

**Cause:** Transaction was already broadcast to the network.

**Solution:**
- This is **not an error** - the transaction was successfully sent previously
- Check blockchain explorer to confirm transaction status
- No action needed

#### Error: "Failed to read PSBT file"

**Symptom:**
```
Error: failed to read PSBT file: invalid file extension (expected .psbt)
```

**Cause:** File doesn't have `.psbt` extension.

**Solution:**
```bash
# Ensure file has correct extension
mv payment_12_signed_2_1534744600000000002 payment_12_signed_2_1534744600000000002.psbt
```

#### Error: "Missing private key"

**Symptom:**
```
Error: failed to sign PSBT: private key not found for address
```

**Cause:** Wallet doesn't have the required private key.

**Solution:**
```bash
# Verify key exists in database
# Keygen wallet:
./keygen list-keys --account client

# Sign wallet:
./sign list-keys --auth auth1

# If key is missing, import or regenerate keys
./keygen import-privkey --file keys.txt
```

### Debugging Tips

#### Inspect PSBT Contents

Use Bitcoin Core to inspect PSBT:

```bash
# Decode PSBT (Bitcoin Core)
bitcoin-cli decodepsbt "$(cat payment_12_unsigned_0_*.psbt)"
```

**Output shows:**
- Inputs and their UTXOs
- Outputs and amounts
- Fee
- Signatures status
- Missing information

#### Validate PSBT

```bash
# Validate PSBT format (Bitcoin Core)
bitcoin-cli analyzepsbt "$(cat payment_12_unsigned_0_*.psbt)"
```

**Output shows:**
- Number of inputs
- Required signatures per input
- Current signature count
- Missing signatures
- Estimated fee
- Next steps

#### Check Transaction Status

```bash
# Check transaction in database
sqlite3 data/db/btc_watch.db "SELECT * FROM btc_tx WHERE id = 12;"

# Check if transaction is on blockchain
bitcoin-cli getrawtransaction <txid> 1
```

---

## Best Practices

### Security

1. **Air-Gapped Signing**
   - ✅ Keep Keygen and Sign wallets offline at all times
   - ✅ Use dedicated, isolated computers for offline wallets
   - ✅ Never connect offline wallets to networks

2. **File Transfer Security**
   - ✅ Use dedicated USB drives for PSBT transfer
   - ✅ Virus scan USB drives before use
   - ✅ Verify file checksums after transfer

3. **Private Key Protection**
   - ✅ Store seeds and private keys in secure offline storage
   - ✅ Use hardware security modules (HSMs) for production
   - ✅ Implement proper access controls

4. **Transaction Verification**
   - ✅ Always verify transaction amounts before signing
   - ✅ Check recipient addresses carefully
   - ✅ Verify fees are reasonable

### Operations

1. **File Naming**
   - ✅ Use auto-generated filenames (don't rename PSBT files)
   - ✅ Keep files organized by transaction type
   - ✅ Archive old PSBTs regularly

2. **Workflow Documentation**
   - ✅ Document your signing workflow
   - ✅ Create checklists for each transaction type
   - ✅ Train operators on PSBT procedures

3. **Testing**
   - ✅ Test on testnet before mainnet
   - ✅ Test with small amounts first
   - ✅ Verify full workflow end-to-end

4. **Monitoring**
   - ✅ Monitor transaction confirmations
   - ✅ Track fee rates and adjust as needed
   - ✅ Set up alerts for transaction failures

### Performance

1. **Fee Management**
   - ✅ Monitor mempool for optimal fee rates
   - ✅ Use higher fees for urgent transactions
   - ✅ Consider consolidating UTXOs during low fee periods

2. **UTXO Management**
   - ✅ Avoid creating dust outputs
   - ✅ Consolidate UTXOs when fees are low
   - ✅ Monitor UTXO set size

3. **Transaction Batching**
   - ✅ Batch multiple payments into single transaction
   - ✅ Reduces overall fees
   - ✅ Increases efficiency

### Backup and Recovery

1. **Regular Backups**
   - ✅ Backup seeds and private keys securely
   - ✅ Backup wallet databases regularly
   - ✅ Test recovery procedures

2. **Disaster Recovery**
   - ✅ Document recovery procedures
   - ✅ Store backups in multiple secure locations
   - ✅ Test recovery periodically

3. **Key Management**
   - ✅ Use BIP39 mnemonics for seed backup
   - ✅ Store seed backups in secure, offline locations
   - ✅ Consider multi-signature setup for critical keys

---

## Additional Resources

### Documentation

- [PSBT Implementation Details](psbt_implementation.md)
- [PSBT Migration Guide](psbt_migration.md)
- [PSBT Developer Guide](psbt_developer_guide.md)
- [Operation Examples](operation_example.md)

### Standards

- [BIP 174: PSBT Specification](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 341: Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP 86: Key Derivation for Taproot](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)

### Tools

- [Bitcoin Core](https://bitcoincore.org/) - Reference implementation with PSBT support
- [btcd](https://github.com/btcsuite/btcd) - Go Bitcoin implementation (used by go-crypto-wallet)
- [Electrum](https://electrum.org/) - Desktop wallet with PSBT support
- [Sparrow Wallet](https://sparrowwallet.com/) - Modern wallet with excellent PSBT features

### Support

- GitHub Issues: [https://github.com/hiromaily/go-crypto-wallet/issues](https://github.com/hiromaily/go-crypto-wallet/issues)
- Bitcoin Stack Exchange: [https://bitcoin.stackexchange.com/](https://bitcoin.stackexchange.com/)

---

## Glossary

- **PSBT**: Partially Signed Bitcoin Transaction (BIP 174)
- **UTXO**: Unspent Transaction Output
- **Multisig**: Multi-signature address requiring multiple signatures
- **Air-gapped**: Computer never connected to networks
- **Schnorr**: Signature algorithm used by Taproot (BIP 340)
- **ECDSA**: Elliptic Curve Digital Signature Algorithm (legacy)
- **Bech32**: Address encoding for SegWit (BIP 173)
- **Taproot**: Bitcoin upgrade enabling Schnorr signatures (BIP 341)
- **BIP**: Bitcoin Improvement Proposal
- **RBF**: Replace-By-Fee transaction replacement mechanism

---

**Last Updated**: 2025-01-27
**Version**: 1.0 (PSBT Phase 2 Complete)
