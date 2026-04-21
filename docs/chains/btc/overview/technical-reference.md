<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/btc/overview/technical-reference.tpl.md · Run `make docs` to regenerate.
-->

# Bitcoin and Bitcoin Cash Technical Guide for AI Agents (2025)

This document provides a comprehensive technical reference for AI Agents working with Bitcoin (BTC) and Bitcoin Cash (BCH) transaction systems, covering address types, multisig mechanisms, and transaction sending workflows.

## Table of Contents

1. [Overview](#overview)
2. [Address Types](#address-types)
3. [Multisig Mechanisms](#multisig-mechanisms)
4. [Transaction Sending Mechanisms](#transaction-sending-mechanisms)
5. [Key Differences Between BTC and BCH](#key-differences-between-btc-and-bch)
6. [Quick Reference Tables](#quick-reference-tables)

---

## Overview

### Bitcoin (BTC) vs Bitcoin Cash (BCH)

| Feature | Bitcoin (BTC) | Bitcoin Cash (BCH) |
|---------|---------------|-------------------|
| **Fork Date** | Original (2009) | August 2017 (hard fork from BTC) |
| **Block Size** | 1-4 MB (with SegWit) | 32 MB |
| **SegWit Support** | ✅ Yes (activated 2017) | ❌ No |
| **Taproot Support** | ✅ Yes (activated 2021) | ❌ No |
| **Primary Focus** | Store of value, Layer 2 solutions | Peer-to-peer electronic cash |
| **Core Implementation** | Bitcoin Core | Bitcoin ABC / BCHN |
| **Required Version** | v22.0+ (for Taproot) | v0.24.0+ |

### Shared Characteristics

Both BTC and BCH share the following foundational properties:

- **UTXO Model**: Unspent Transaction Output model for tracking ownership
- **ECDSA Signatures**: secp256k1 elliptic curve for cryptographic operations
- **SHA-256 Mining**: Same proof-of-work algorithm
- **Base58Check Encoding**: For legacy address format encoding
- **Script Language**: Bitcoin Script for transaction validation

---

## Address Types

### Bitcoin (BTC) Address Types

Bitcoin supports four main address types, each corresponding to different script types and BIP standards:

#### 1. Legacy (P2PKH) - BIP44

```
Purpose: 44'
Derivation Path: m/44'/0'/account'/change/index
```

| Property | Value |
|----------|-------|
| **Prefix** | `1` (mainnet), `m` or `n` (testnet) |
| **Length** | 25-34 characters |
| **Encoding** | Base58Check |
| **Script Type** | Pay-to-Public-Key-Hash |
| **Example** | `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa` |

**ScriptPubKey Structure:**

```
OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
```

**Pros:**

- Maximum compatibility with all wallets/exchanges
- Widely understood and documented

**Cons:**

- Largest transaction size
- Highest fees
- No SegWit benefits

#### 2. P2SH-SegWit (P2SH-P2WPKH) - BIP49

```
Purpose: 49'
Derivation Path: m/49'/0'/account'/change/index
```

| Property | Value |
|----------|-------|
| **Prefix** | `3` (mainnet), `2` (testnet) |
| **Length** | 34 characters |
| **Encoding** | Base58Check |
| **Script Type** | Pay-to-Script-Hash wrapping SegWit |
| **Example** | `3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy` |

**ScriptPubKey Structure:**

```
OP_HASH160 <scriptHash> OP_EQUAL
```

**Pros:**

- SegWit fee savings (~26% reduction)
- Backward compatible with legacy wallets
- Improved transaction malleability protection

**Cons:**

- Not as efficient as native SegWit
- More complex script structure

#### 3. Native SegWit (P2WPKH) - BIP84

```
Purpose: 84'
Derivation Path: m/84'/0'/account'/change/index
```

| Property | Value |
|----------|-------|
| **Prefix** | `bc1q` (mainnet), `tb1q` (testnet) |
| **Length** | 42 characters |
| **Encoding** | Bech32 |
| **Script Type** | Pay-to-Witness-Public-Key-Hash |
| **Example** | `bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq` |

**Witness Script Structure:**

```
0 <20-byte pubKeyHash>
```

**Pros:**

- Best pre-Taproot fee efficiency (~38% reduction)
- Full SegWit benefits
- Error detection with Bech32 checksum

**Cons:**

- Some older wallets don't support sending to bc1q

#### 4. Taproot (P2TR) - BIP86

```
Purpose: 86'
Derivation Path: m/86'/0'/account'/change/index
```

| Property | Value |
|----------|-------|
| **Prefix** | `bc1p` (mainnet), `tb1p` (testnet) |
| **Length** | 62 characters |
| **Encoding** | Bech32m |
| **Script Type** | Pay-to-Taproot |
| **Signature** | Schnorr (BIP340) |
| **Example** | `bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr` |

**Witness Script Structure:**

```
1 <32-byte x-only pubKey>
```

**Pros:**

- Best fee efficiency (~50% reduction for multisig)
- Enhanced privacy (all spends look identical)
- Schnorr signature aggregation (MuSig2)
- Script path spending for complex contracts

**Cons:**

- Requires Bitcoin Core v22.0+
- Some exchanges still don't support withdrawal to bc1p

### Bitcoin Cash (BCH) Address Types

Bitcoin Cash uses the **CashAddr** format (introduced 2018) with its own prefix and encoding:

#### 1. CashAddr P2PKH

| Property | Value |
|----------|-------|
| **Prefix** | `bitcoincash:q` (mainnet), `bchtest:q` (testnet) |
| **Length** | ~50 characters (with prefix) |
| **Encoding** | CashAddr (modified Bech32) |
| **Script Type** | Pay-to-Public-Key-Hash |
| **Example** | `bitcoincash:qp3wjpa3tjlj042z2wv7hahsldgwhwy0rq9sywjpyy` |

**ScriptPubKey Structure:**

```
OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
```

#### 2. CashAddr P2SH

| Property | Value |
|----------|-------|
| **Prefix** | `bitcoincash:p` (mainnet), `bchtest:p` (testnet) |
| **Length** | ~50 characters (with prefix) |
| **Encoding** | CashAddr (modified Bech32) |
| **Script Type** | Pay-to-Script-Hash |
| **Example** | `bitcoincash:pr0662zpd7vr936d83f64u629v886aan7c77r3j5v5` |

**CashAddr Characteristics:**

- Case-insensitive (always lowercase recommended)
- Built-in checksum for error detection
- Network prefix prevents cross-chain sending errors
- Can omit prefix when network is unambiguous

#### 3. Legacy Format (Backward Compatible)

BCH can also use legacy Base58Check addresses (starts with `1` for P2PKH, `3` for P2SH), but CashAddr is strongly recommended to prevent accidental BTC/BCH cross-sends.

---

## Multisig Mechanisms

### Multisig Concepts

**Definition:** Multi-signature (multisig) requires M-of-N signatures to authorize a transaction.

**Common Configurations:**

- **2-of-2**: Both parties must sign (joint control)
- **2-of-3**: Any 2 of 3 parties must sign (recovery/escrow)
- **3-of-5**: Corporate treasury management

### Bitcoin Multisig Types

#### 1. Traditional P2SH Multisig (Legacy)

```
Address Prefix: 3 (mainnet)
Script Type: Pay-to-Script-Hash
BIP: BIP11, BIP16
```

**Redeem Script Structure:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

**Example 2-of-3:**

```
2 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

**Characteristics:**

| Property | Value |
|----------|-------|
| **Maximum Keys** | 15 (due to OP_CHECKMULTISIG limit) |
| **Transaction Size** | ~370-400 bytes (2-of-3) |
| **Privacy** | Low (multisig is visible on-chain) |
| **Signature Type** | ECDSA |

#### 2. P2WSH Multisig (Native SegWit)

```
Address Prefix: bc1q (longer than P2WPKH)
Script Type: Pay-to-Witness-Script-Hash
BIP: BIP141, BIP143
```

**Witness Script Structure:**

```
0 <32-byte scriptHash>
```

**Characteristics:**

| Property | Value |
|----------|-------|
| **Maximum Keys** | 15 |
| **Transaction Size** | ~270-300 bytes (2-of-3) |
| **Fee Savings** | ~27% vs P2SH |
| **Privacy** | Low (multisig visible) |
| **Signature Type** | ECDSA |

#### 3. Taproot Multisig (Key Path)

With Taproot, multisig can be implemented in two ways:

**A. Key Path Spending (MuSig2)**

```
Address Prefix: bc1p
Script Type: Pay-to-Taproot (key path)
BIP: BIP340, BIP341, BIP327
```

**MuSig2 Protocol (Two-Round):**

```
Round 1: Nonce Generation (Parallel)
├── Signer 1 generates nonce
├── Signer 2 generates nonce
└── Signer N generates nonce
    │
    ▼
  Exchange public nonces
    │
    ▼
Round 2: Partial Signing (Sequential)
├── Signer 1 creates partial signature
├── Signer 2 creates partial signature
└── Signer N creates partial signature
    │
    ▼
  Aggregate into single signature
```

**MuSig2 Characteristics:**

| Property | Value |
|----------|-------|
| **On-Chain Appearance** | Single signature (indistinguishable from single-sig) |
| **Transaction Size** | ~200-250 bytes |
| **Fee Savings** | 30-50% vs P2WSH multisig |
| **Privacy** | Maximum (looks like single-sig) |
| **Signature Type** | Schnorr (64 bytes) |
| **Flexibility** | N-of-N only (all signers required) |

**CRITICAL: Nonce Security**

- **NEVER reuse nonces** - Reusing a nonce leaks the private key
- Generate fresh nonces for every transaction
- Delete nonces after signing

**B. Script Path Spending**

For M-of-N (where M < N), Taproot uses script path:

```
Taproot Output Key = Internal Key + TapTweak(Merkle Root)
```

**Script Tree Structure:**

```
              Root
             /    \
        Branch    Leaf3
        /    \     (2-of-3 script)
    Leaf1   Leaf2
    (2-of-3) (2-of-3)
```

**Characteristics:**

- Multiple spending conditions organized in Merkle tree
- Only the used script path is revealed on-chain
- Unused branches remain private

### Bitcoin Cash Multisig

BCH supports only traditional P2SH multisig (no SegWit/Taproot):

```
Address Prefix: p (CashAddr)
Script Type: Pay-to-Script-Hash
```

**Characteristics:**

| Property | Value |
|----------|-------|
| **Maximum Keys** | 15 |
| **Transaction Size** | ~370-400 bytes (2-of-3) |
| **Privacy** | Low |
| **Signature Type** | ECDSA |

---

## Transaction Sending Mechanisms

### UTXO Model

Both BTC and BCH use the UTXO (Unspent Transaction Output) model:

```
Transaction Structure:
┌──────────────────────────────────────────────┐
│ Version (4 bytes)                            │
├──────────────────────────────────────────────┤
│ Input Count (VarInt)                         │
├──────────────────────────────────────────────┤
│ Inputs[] (variable)                          │
│   ├── Previous TX Hash (32 bytes)            │
│   ├── Previous Output Index (4 bytes)        │
│   ├── ScriptSig Length (VarInt)              │
│   ├── ScriptSig (variable)                   │
│   └── Sequence (4 bytes)                     │
├──────────────────────────────────────────────┤
│ Output Count (VarInt)                        │
├──────────────────────────────────────────────┤
│ Outputs[] (variable)                         │
│   ├── Value (8 bytes)                        │
│   ├── ScriptPubKey Length (VarInt)           │
│   └── ScriptPubKey (variable)                │
├──────────────────────────────────────────────┤
│ Locktime (4 bytes)                           │
└──────────────────────────────────────────────┘
```

**SegWit Transaction (BTC only):**

```
┌──────────────────────────────────────────────┐
│ Version (4 bytes)                            │
├──────────────────────────────────────────────┤
│ Marker (1 byte) = 0x00                       │
│ Flag (1 byte) = 0x01                         │
├──────────────────────────────────────────────┤
│ Input Count (VarInt)                         │
├──────────────────────────────────────────────┤
│ Inputs[] (scriptSig is empty for SegWit)     │
├──────────────────────────────────────────────┤
│ Output Count (VarInt)                        │
├──────────────────────────────────────────────┤
│ Outputs[]                                    │
├──────────────────────────────────────────────┤
│ Witness[] (SegWit data)                      │
│   └── Per input: signatures and pubkeys      │
├──────────────────────────────────────────────┤
│ Locktime (4 bytes)                           │
└──────────────────────────────────────────────┘
```

### Transaction Workflow

#### Standard Transaction Flow

```
1. Create Transaction
   └── Select UTXOs (inputs)
   └── Define outputs (recipients + change)
   └── Calculate fee

2. Sign Transaction
   └── For each input:
       └── Generate signature using private key
       └── Signature covers: outputs, input data, sighash

3. Broadcast Transaction
   └── Send to node/network
   └── Node validates and propagates
   └── Miners include in block

4. Confirmation
   └── Block mined containing transaction
   └── Each subsequent block adds confirmation
```

### PSBT (Partially Signed Bitcoin Transaction)

**BIP174 - Standard for offline signing and multi-party transactions**

PSBT enables:

- Offline/air-gapped signing
- Hardware wallet integration
- Multi-party signing coordination
- Watch-only wallet workflows

**PSBT Structure:**

```
┌──────────────────────────────────────────────┐
│ Magic Bytes: "psbt" + 0xFF (5 bytes)         │
├──────────────────────────────────────────────┤
│ Global Map                                   │
│   ├── PSBT_GLOBAL_UNSIGNED_TX               │
│   ├── PSBT_GLOBAL_XPUB (optional)           │
│   └── PSBT_GLOBAL_VERSION                   │
├──────────────────────────────────────────────┤
│ Input Maps[] (one per input)                 │
│   ├── PSBT_IN_WITNESS_UTXO                  │
│   ├── PSBT_IN_PARTIAL_SIG                   │
│   ├── PSBT_IN_SIGHASH_TYPE                  │
│   ├── PSBT_IN_BIP32_DERIVATION              │
│   └── PSBT_IN_TAP_* (Taproot fields)        │
├──────────────────────────────────────────────┤
│ Output Maps[] (one per output)               │
│   ├── PSBT_OUT_BIP32_DERIVATION             │
│   └── PSBT_OUT_TAP_* (Taproot fields)       │
└──────────────────────────────────────────────┘
```

**PSBT Workflow:**

```
Creator (Watch Wallet - Online)
    │
    │ 1. Create unsigned PSBT
    │    - Select UTXOs
    │    - Add witness data
    │    - Save as .psbt file
    ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_0_*.psbt    │
└─────────────────────────────────┘
    │
    │ Transfer via USB/QR
    ▼
Signer 1 (Keygen Wallet - Offline)
    │
    │ 2. Sign with first key
    │    - Add partial signature
    ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_1_*.psbt    │
└─────────────────────────────────┘
    │
    │ Transfer via USB/QR
    ▼
Signer 2 (Sign Wallet - Offline)
    │
    │ 3. Sign with second key
    │    - Add partial signature
    ▼
┌─────────────────────────────────┐
│ payment_15_signed_2_*.psbt      │
└─────────────────────────────────┘
    │
    │ Transfer via USB/QR
    ▼
Finalizer (Watch Wallet - Online)
    │
    │ 4. Finalize PSBT
    │    - Combine signatures
    │    - Create scriptSig/witness
    │    - Extract final transaction
    │
    │ 5. Broadcast
    │    - Send to network
    ▼
    Confirmed on blockchain
```

### MuSig2 Transaction Workflow

For Taproot multisig with MuSig2:

```
Creator (Watch Wallet)
    │
    │ 1. Create unsigned PSBT with MuSig2 metadata
    ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_0_*.psbt    │
└─────────────────────────────────┘
    │
    ▼
═══════════════════════════════════════
    ROUND 1: NONCE GENERATION (PARALLEL)
═══════════════════════════════════════
    │
    ├──────┬──────┬──────┐
    ▼      ▼      ▼      ▼
  Signer1 Signer2 Signer3
    │      │      │
    │ Generate nonces (can run in parallel)
    │      │      │
    └──────┴──────┴──────┐
                        ▼
┌─────────────────────────────────┐
│ payment_15_nonce_0_*.psbt       │
│ (contains all 3 public nonces)  │
└─────────────────────────────────┘
    │
    ▼
═══════════════════════════════════════
    ROUND 2: PARTIAL SIGNING (SEQUENTIAL)
═══════════════════════════════════════
    │
    │ All nonces must be collected first!
    ▼
  Signer1 ─── Partial Sig 1 ───┐
                              ▼
  Signer2 ─── Partial Sig 2 ───┐
                              ▼
  Signer3 ─── Partial Sig 3 ───┐
                              ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_3_*.psbt    │
│ (contains all partial sigs)     │
└─────────────────────────────────┘
    │
    ▼
═══════════════════════════════════════
    AGGREGATION & BROADCAST
═══════════════════════════════════════
    │
Aggregator (Watch Wallet)
    │
    │ Aggregate partial signatures → single Schnorr sig
    ▼
┌─────────────────────────────────┐
│ payment_15_signed_3_*.psbt      │
│ (final 64-byte Schnorr sig)     │
└─────────────────────────────────┘
    │
    │ Broadcast
    ▼
  Blockchain (looks like single-sig tx!)
```

### Bitcoin Cash Transaction Workflow

BCH uses a simpler workflow without PSBT (legacy raw transactions):

```
Creator (Watch Wallet)
    │
    │ 1. Create raw transaction
    │    - Hex-encoded transaction
    ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_0_*.tx      │
│ (hex,prev_addresses format)     │
└─────────────────────────────────┘
    │
    ▼
Signer 1 (Keygen Wallet)
    │
    │ 2. Sign with ECDSA
    │    - Create scriptSig
    ▼
┌─────────────────────────────────┐
│ payment_15_unsigned_1_*.tx      │
└─────────────────────────────────┘
    │
    ▼
Signer 2 (Sign Wallet) [if multisig]
    │
    │ 3. Add second signature
    ▼
┌─────────────────────────────────┐
│ payment_15_signed_2_*.tx        │
└─────────────────────────────────┘
    │
    ▼
Broadcaster (Watch Wallet)
    │
    │ 4. Send to network
    ▼
  Confirmed on blockchain
```

### Signature Types Comparison

| Signature | Algorithm | Size | Used By |
|-----------|-----------|------|---------|
| **ECDSA** | secp256k1 | ~71-72 bytes (DER encoded) | BTC (legacy/SegWit), BCH |
| **Schnorr** | BIP340 | 64 bytes (fixed) | BTC (Taproot only) |

**Schnorr Benefits:**

- Fixed 64-byte size (vs variable ECDSA)
- Linearity enables signature aggregation
- Batch verification is faster
- Provably secure under standard assumptions

---

## Key Differences Between BTC and BCH

### Feature Comparison Matrix

| Feature | Bitcoin (BTC) | Bitcoin Cash (BCH) |
|---------|---------------|-------------------|
| **SegWit** | ✅ Full support | ❌ Not implemented |
| **Taproot** | ✅ Full support | ❌ Not implemented |
| **Schnorr Signatures** | ✅ BIP340 | ❌ Not available |
| **MuSig2** | ✅ BIP327 | ❌ Not available |
| **PSBT** | ✅ BIP174 | ⚠️ Limited support |
| **Address Format** | Bech32/Bech32m | CashAddr |
| **Block Size** | 1-4 MB (SegWit) | 32 MB |
| **Transaction Malleability** | ✅ Fixed (SegWit) | ⚠️ Still present |
| **RBF** | ✅ Replace-by-fee | ⚠️ Child-pays-for-parent |

### Implementation Differences in go-crypto-wallet

The codebase handles BTC and BCH through inheritance:

```go
// Bitcoin Cash embeds and extends Bitcoin
type BitcoinCash struct {
    btc.Bitcoin  // Inherits all Bitcoin functionality
}

// BCH overrides chain parameters
func (b *BitcoinCash) initChainParams() {
    switch conf.Name {
    case chaincfg.TestNet3Params.Name:
        conf.Net = TestnetMagic  // BCH-specific magic
    case chaincfg.MainNetParams.Name:
        conf.Net = MainnetMagic
    }
}
```

**Shared Code Path:**

- Transaction creation
- UTXO selection
- Fee calculation
- Basic signing (ECDSA)

**BTC-Only Features:**

- SegWit witness data handling
- Taproot/Schnorr signing
- MuSig2 nonce management
- PSBT creation/parsing
- Bech32/Bech32m encoding

**BCH-Specific:**

- CashAddr encoding/decoding
- Different network magic bytes
- No witness data in transactions

---

## Quick Reference Tables

### Address Type Selection Guide

| Use Case | BTC Recommendation | BCH Recommendation |
|----------|-------------------|-------------------|
| **New wallet setup** | Taproot (bc1p) | CashAddr P2PKH |
| **Maximum compatibility** | P2SH-SegWit (3) | CashAddr P2PKH |
| **Best fees** | Taproot (bc1p) | CashAddr P2PKH |
| **Simple multisig** | P2WSH | P2SH |
| **Private multisig** | MuSig2 Taproot | N/A |
| **Hardware wallet** | Native SegWit (bc1q) | CashAddr P2PKH |

### Multisig Selection Guide

| Requirement | BTC Recommendation | BCH Recommendation |
|-------------|-------------------|-------------------|
| **2-of-3 corporate** | P2WSH or Taproot script path | P2SH |
| **N-of-N with privacy** | MuSig2 (Taproot) | N/A |
| **Maximum compatibility** | P2SH | P2SH |
| **Best fee efficiency** | MuSig2 (Taproot) | N/A (use P2SH) |
| **Threshold (M < N)** | Taproot script path | P2SH |

### Transaction Size Comparison (2-of-3 Multisig)

| Type | Size | Fee @ 10 sat/vB | Notes |
|------|------|-----------------|-------|
| **P2SH (BTC/BCH)** | ~370 bytes | ~3,700 sats | Legacy, most compatible |
| **P2WSH (BTC)** | ~270 bytes | ~2,700 sats | SegWit, ~27% savings |
| **MuSig2 (BTC)** | ~215 bytes | ~2,150 sats | 42% savings, max privacy |

### Network Parameters

| Parameter | BTC Mainnet | BTC Testnet | BCH Mainnet | BCH Testnet |
|-----------|-------------|-------------|-------------|-------------|
| **P2PKH Prefix** | `1` | `m`/`n` | `q` (CashAddr) | `q` (bchtest:) |
| **P2SH Prefix** | `3` | `2` | `p` (CashAddr) | `p` (bchtest:) |
| **Bech32 Prefix** | `bc1` | `tb1` | N/A | N/A |
| **Bech32m Prefix** | `bc1p` | `tb1p` | N/A | N/A |
| **Default Port** | 8333 | 18333 | 8333 | 18333 |
| **RPC Port** | 8332 | 18332 | 8332 | 18332 |

---

## Related Documentation

### Project-Specific Guides

- [Taproot User Guide](/chains/btc/taproot/user-guide) - Detailed Taproot implementation
- [MuSig2 User Guide](/chains/btc/musig2/user-guide) - MuSig2 multisig operations
- [PSBT Developer Guide](/chains/btc/psbt/developer-guide) - PSBT implementation details
- [Wallet Flow](/chains/btc/operations/wallet-flow) - Wallet configuration and usage

### External References

- [BIP32 - HD Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [BIP44 - Multi-Account HD](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
- [BIP84 - Native SegWit](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki)
- [BIP86 - Taproot Derivation](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)
- [BIP141 - SegWit](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki)
- [BIP174 - PSBT](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP327 - MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [BIP340 - Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP341 - Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BCH CashAddr Spec](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md)

---

**Document Version:** 1.0
**Last Updated:** 2025-01-07
**Target Audience:** AI Agents, Developers
**Applicable Versions:** BTC Core v22.0+, BCH ABC v0.24.0+
