# Bitcoin (BTC) Technical Reference

This document provides a comprehensive technical reference for Bitcoin implementation in the go-crypto-wallet system. It covers specifications, protocol details, and links to official documentation to help AI agents and developers understand Bitcoin's architecture and implement features correctly.

## Table of Contents

1. [Overview](#overview)
2. [Core Specifications](#core-specifications)
3. [Address Types & Key Derivation](#address-types--key-derivation)
4. [Transaction Architecture](#transaction-architecture)
5. [Signing Mechanisms](#signing-mechanisms)
6. [Multisig & MuSig2](#multisig--musig2)
7. [PSBT (Partially Signed Bitcoin Transactions)](#psbt-partially-signed-bitcoin-transactions)
8. [Network & Consensus](#network--consensus)
9. [Fee Management](#fee-management)
10. [Wallet Implementation](#wallet-implementation)
11. [RPC & API Reference](#rpc--api-reference)
12. [Security Considerations](#security-considerations)
13. [Testing Resources](#testing-resources)
14. [Official References](#official-references)
15. [Project Documentation](#project-documentation)

---

## Overview

### What is Bitcoin?

Bitcoin is a decentralized digital currency that operates on a peer-to-peer network without central authority. It uses a UTXO (Unspent Transaction Output) model for tracking ownership and proof-of-work consensus for block validation.

### Key Characteristics (2026)

| Property | Value |
|----------|-------|
| **Launch Date** | January 3, 2009 |
| **Block Time** | ~10 minutes |
| **Block Size** | 1-4 MB (with SegWit) |
| **Total Supply** | 21,000,000 BTC |
| **Current Block Reward** | 3.125 BTC (post-2024 halving) |
| **Next Halving** | ~2028 (Block 1,050,000) |
| **Consensus Algorithm** | SHA-256 Proof of Work |
| **Cryptographic Curve** | secp256k1 |
| **Signature Algorithms** | ECDSA (legacy), Schnorr (Taproot) |

### Protocol Upgrades Timeline

| Year | Upgrade | Key Features |
|------|---------|--------------|
| 2017 | **SegWit (BIP141)** | Transaction malleability fix, increased capacity |
| 2021 | **Taproot (BIP340/341)** | Schnorr signatures, MAST, privacy improvements |
| 2023 | **BIP327 MuSig2** | Standardized multi-signature aggregation |
| 2024+ | **OP_CAT (Proposed)** | Enhanced scripting capabilities |

---

## Core Specifications

### Cryptographic Primitives

#### Elliptic Curve (secp256k1)

Bitcoin uses the secp256k1 elliptic curve for all cryptographic operations:

```
Curve Parameters:
- p = 2^256 - 2^32 - 977
- a = 0
- b = 7
- G = (0x79BE667E..., 0x483ADA77...)
- n = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFE BAAEDCE6AF48A03BBFD25E8CD0364141
```

**Reference:**

- [SEC 2: Recommended Elliptic Curve Domain Parameters](https://www.secg.org/sec2-v2.pdf)
- [Bitcoin secp256k1 Library](https://github.com/bitcoin-core/secp256k1)

#### Hash Functions

| Function | Usage |
|----------|-------|
| **SHA-256** | Block hashing, TXID calculation, PoW |
| **RIPEMD-160** | Address generation (hash160 = RIPEMD160(SHA256(x))) |
| **HASH160** | SHA256 + RIPEMD160 for pubkey hashing |
| **HASH256** | Double SHA256 for transaction/block hashing |
| **Tagged Hashes** | BIP340 Schnorr signatures (SHA256 with tag) |

### Data Encoding

| Format | Usage | Reference |
|--------|-------|-----------|
| **Base58Check** | Legacy addresses (P2PKH, P2SH) | [Base58Check](https://en.bitcoin.it/wiki/Base58Check_encoding) |
| **Bech32** | Native SegWit addresses (P2WPKH, P2WSH) | [BIP173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki) |
| **Bech32m** | Taproot addresses (P2TR) | [BIP350](https://github.com/bitcoin/bips/blob/master/bip-0350.mediawiki) |
| **WIF** | Private key encoding | [Wallet Import Format](https://en.bitcoin.it/wiki/Wallet_import_format) |
| **Hex** | Raw transaction data | Standard hexadecimal |

---

## Address Types & Key Derivation

### Address Types Supported

| Type | BIP | Prefix (Mainnet) | Prefix (Testnet) | Description |
|------|-----|------------------|------------------|-------------|
| **P2PKH** | BIP44 | `1` | `m`/`n` | Legacy Pay-to-Public-Key-Hash |
| **P2SH** | BIP16 | `3` | `2` | Pay-to-Script-Hash |
| **P2SH-P2WPKH** | BIP49 | `3` | `2` | SegWit wrapped in P2SH |
| **P2WPKH** | BIP84 | `bc1q` | `tb1q` | Native SegWit |
| **P2WSH** | BIP141 | `bc1q` | `tb1q` | SegWit Script Hash |
| **P2TR** | BIP86 | `bc1p` | `tb1p` | Taproot (recommended) |

### HD Wallet Derivation Paths

| Standard | Path | Address Type |
|----------|------|--------------|
| **BIP44** | `m/44'/0'/account'/change/index` | P2PKH (Legacy) |
| **BIP49** | `m/49'/0'/account'/change/index` | P2SH-P2WPKH |
| **BIP84** | `m/84'/0'/account'/change/index` | P2WPKH (Native SegWit) |
| **BIP86** | `m/86'/0'/account'/change/index` | P2TR (Taproot) |

**Coin Types:**

- `0'` = Bitcoin Mainnet
- `1'` = Bitcoin Testnet/Signet

**References:**

- [BIP32 - HD Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [BIP39 - Mnemonic Codes](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- [BIP43 - Purpose Field](https://github.com/bitcoin/bips/blob/master/bip-0043.mediawiki)
- [BIP44 - Multi-Account Hierarchy](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

### ScriptPubKey Formats

```
P2PKH:      OP_DUP OP_HASH160 <20-byte pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
P2SH:       OP_HASH160 <20-byte scriptHash> OP_EQUAL
P2WPKH:     0x00 <20-byte pubKeyHash>
P2WSH:      0x00 <32-byte witnessScriptHash>
P2TR:       0x51 <32-byte x-only pubKey>
```

---

## Transaction Architecture

### UTXO Model

Bitcoin uses the Unspent Transaction Output (UTXO) model:

```
UTXO = {
    txid:      32-byte transaction hash
    vout:      output index (uint32)
    value:     satoshi amount (int64)
    scriptPubKey: locking script
}
```

**Key Concepts:**

- Each transaction consumes UTXOs (inputs) and creates new UTXOs (outputs)
- Total inputs must equal outputs + transaction fee
- UTXOs can only be spent once (double-spend protection)

### Transaction Structure

#### Legacy Transaction Format

```
+--------------------+
| Version (4 bytes)  |
+--------------------+
| Input Count        |  VarInt
+--------------------+
| Inputs[]           |
|   - prevTxHash     |  32 bytes
|   - prevVout       |  4 bytes
|   - scriptSigLen   |  VarInt
|   - scriptSig      |  variable
|   - sequence       |  4 bytes
+--------------------+
| Output Count       |  VarInt
+--------------------+
| Outputs[]          |
|   - value          |  8 bytes (satoshis)
|   - scriptPubKeyLen|  VarInt
|   - scriptPubKey   |  variable
+--------------------+
| Locktime (4 bytes) |
+--------------------+
```

#### SegWit Transaction Format (BIP141)

```
+--------------------+
| Version (4 bytes)  |
+--------------------+
| Marker (0x00)      |  1 byte (SegWit indicator)
| Flag (0x01)        |  1 byte
+--------------------+
| Input Count        |
+--------------------+
| Inputs[]           |  (scriptSig empty for SegWit)
+--------------------+
| Output Count       |
+--------------------+
| Outputs[]          |
+--------------------+
| Witness[]          |  SegWit data (signatures, pubkeys)
+--------------------+
| Locktime (4 bytes) |
+--------------------+
```

### Transaction Weight & Virtual Size

SegWit introduced weight units for fee calculation:

```
Weight = (Non-witness data × 4) + Witness data
Virtual Size (vBytes) = Weight ÷ 4

Fee = Virtual Size × Fee Rate (sat/vB)
```

**Typical Sizes:**

| Transaction Type | Weight | vBytes | Fee @ 10 sat/vB |
|------------------|--------|--------|-----------------|
| P2PKH (1-in, 2-out) | ~680 | ~170 | ~1,700 sats |
| P2WPKH (1-in, 2-out) | ~440 | ~110 | ~1,100 sats |
| P2TR (1-in, 2-out) | ~396 | ~99 | ~990 sats |
| 2-of-3 Multisig (P2WSH) | ~1,100 | ~275 | ~2,750 sats |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | ~1,400 sats |

**Reference:**

- [BIP141 - SegWit](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki)
- [BIP144 - SegWit Peer Services](https://github.com/bitcoin/bips/blob/master/bip-0144.mediawiki)

---

## Signing Mechanisms

### ECDSA Signatures (Legacy/SegWit)

Used for P2PKH, P2SH, P2WPKH, and P2WSH transactions.

**Signature Format (DER Encoded):**

```
0x30 [total-length]
  0x02 [r-length] [r]
  0x02 [s-length] [s]
[sighash-type]
```

**Size:** 71-73 bytes (variable due to DER encoding)

**Sighash Types:**

| Type | Value | Description |
|------|-------|-------------|
| SIGHASH_ALL | 0x01 | Sign all inputs and outputs (default) |
| SIGHASH_NONE | 0x02 | Sign all inputs, no outputs |
| SIGHASH_SINGLE | 0x03 | Sign all inputs, matching output only |
| SIGHASH_ANYONECANPAY | 0x80 | Modifier: sign only current input |

**Reference:**

- [BIP143 - SegWit Signature Verification](https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki)

### Schnorr Signatures (Taproot)

Used for P2TR transactions. Introduced with Taproot (BIP340).

**Advantages:**

- Fixed 64-byte size (vs variable ECDSA)
- Linear: enables signature aggregation (MuSig2)
- Provably secure under standard assumptions
- Batch verification is faster

**Signature Format:**

```
[32-byte R] [32-byte s]  (64 bytes total, no sighash byte appended)
```

**Key Format (x-only):**

- Taproot uses 32-byte x-only public keys
- Y-coordinate is implicitly even (BIP340 convention)

**References:**

- [BIP340 - Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP341 - Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP342 - Tapscript](https://github.com/bitcoin/bips/blob/master/bip-0342.mediawiki)

---

## Multisig & MuSig2

### Traditional Multisig (P2SH/P2WSH)

**M-of-N Redeem Script:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

**Characteristics:**

- Maximum 15 keys (OP_CHECKMULTISIG limit)
- Multiple signatures visible on-chain
- Higher fees due to larger transaction size

**Reference:**

- [BIP11 - M-of-N Standard Transactions](https://github.com/bitcoin/bips/blob/master/bip-0011.mediawiki)
- [BIP16 - P2SH](https://github.com/bitcoin/bips/blob/master/bip-0016.mediawiki)

### MuSig2 (Schnorr Signature Aggregation)

MuSig2 enables N-of-N multisig that appears as single-sig on-chain.

**Two-Round Protocol:**

```
Round 1: Nonce Generation (Parallel)
├── Each signer generates random nonce
├── Public nonces are exchanged
└── Aggregate nonce computed

Round 2: Partial Signing (Sequential)
├── Each signer creates partial signature
├── Partial signatures are collected
└── Aggregate into single 64-byte Schnorr signature
```

**Benefits:**

- 30-50% smaller transactions
- Maximum privacy (looks like single-sig)
- Lower fees
- No on-chain multisig indicator

**Critical Security: NONCE MANAGEMENT**

- **NEVER reuse nonces** - reusing leaks private key
- Generate fresh nonces for every transaction
- Delete nonces after signing

**References:**

- [BIP327 - MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [MuSig2 Paper (Cryptology ePrint)](https://eprint.iacr.org/2020/1261)

### Taproot Script Path (M-of-N where M < N)

For threshold signatures (e.g., 2-of-3), Taproot uses Merkle script trees:

```
Taproot Output Key = Internal Key + TapTweak(Merkle Root)

              Root
             /    \
        Branch    Leaf3
        /    \     (2-of-3 script)
    Leaf1   Leaf2
    (2-of-3) (2-of-3)
```

- Only the used script path is revealed on-chain
- Unused branches remain private

---

## PSBT (Partially Signed Bitcoin Transactions)

PSBT (BIP174) is the standard format for offline/multi-party signing workflows.

### PSBT Structure

```
+----------------------+
| Magic: "psbt" + 0xFF |  5 bytes
+----------------------+
| Global Map           |
|   - UNSIGNED_TX      |
|   - XPUB (optional)  |
|   - VERSION          |
+----------------------+
| Input Maps[]         |
|   - WITNESS_UTXO     |
|   - PARTIAL_SIG[]    |
|   - SIGHASH_TYPE     |
|   - BIP32_DERIVATION |
|   - TAP_* fields     |
+----------------------+
| Output Maps[]        |
|   - BIP32_DERIVATION |
|   - TAP_* fields     |
+----------------------+
```

### PSBT Workflow

```
1. Creator (Watch Wallet - Online)
   └── Create unsigned PSBT with UTXO data

2. Updater (Optional)
   └── Add metadata (derivation paths, etc.)

3. Signer(s) (Offline Wallets)
   └── Add partial signatures

4. Combiner (Optional)
   └── Combine multiple PSBTs

5. Finalizer (Watch Wallet)
   └── Create final scriptSig/witness

6. Extractor
   └── Extract broadcastable transaction
```

### PSBT Extensions for Taproot

| Field | Description |
|-------|-------------|
| PSBT_IN_TAP_KEY_SIG | Key path Schnorr signature |
| PSBT_IN_TAP_SCRIPT_SIG | Script path signature |
| PSBT_IN_TAP_LEAF_SCRIPT | Tapscript leaf |
| PSBT_IN_TAP_BIP32_DERIVATION | Taproot derivation info |
| PSBT_IN_TAP_INTERNAL_KEY | Internal key |
| PSBT_IN_TAP_MERKLE_ROOT | Script tree Merkle root |

**References:**

- [BIP174 - PSBT](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP370 - PSBT Version 2](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki)
- [BIP371 - Taproot PSBT Fields](https://github.com/bitcoin/bips/blob/master/bip-0371.mediawiki)

---

## Network & Consensus

### Networks

| Network | Purpose | Port | RPC Port | Magic Bytes |
|---------|---------|------|----------|-------------|
| **Mainnet** | Production | 8333 | 8332 | 0xF9BEB4D9 |
| **Testnet3** | Public testing | 18333 | 18332 | 0x0B110907 |
| **Signet** | Controlled testing | 38333 | 38332 | 0x0A03CF40 |
| **Regtest** | Local development | 18444 | 18443 | 0xFABFB5DA |

### Block Structure

```
+----------------------+
| Block Header (80 B)  |
|   - Version (4 B)    |
|   - PrevBlockHash    |  32 bytes
|   - MerkleRoot       |  32 bytes
|   - Timestamp (4 B)  |
|   - Bits (4 B)       |  Difficulty target
|   - Nonce (4 B)      |
+----------------------+
| Transaction Count    |  VarInt
+----------------------+
| Transactions[]       |
+----------------------+
```

### Confirmation Guidelines

| Confirmations | Risk Level | Typical Use Case |
|---------------|------------|------------------|
| 0 (unconfirmed) | High | Very small amounts, trusted parties |
| 1 | Medium | Small retail transactions |
| 3 | Low | Most commerce |
| 6 | Very Low | Large transactions |
| 100+ | None | Coinbase maturity |

---

## Fee Management

### Fee Estimation

Bitcoin Core provides fee estimation via RPC:

```bash
# Estimate fee for confirmation in N blocks
bitcoin-cli estimatesmartfee <conf_target> [estimate_mode]

# Modes: UNSET, ECONOMICAL, CONSERVATIVE
```

### Fee Rate Sources

| Source | Endpoint/Method |
|--------|-----------------|
| **Bitcoin Core** | `estimatesmartfee` RPC |
| **Mempool.space** | `https://mempool.space/api/v1/fees/recommended` |
| **Blockstream** | `https://blockstream.info/api/fee-estimates` |

### Fee Optimization Strategies

1. **SegWit/Taproot** - Use native SegWit or Taproot for smaller transactions
2. **UTXO Consolidation** - Consolidate UTXOs during low-fee periods
3. **Batching** - Combine multiple payments in single transaction
4. **RBF** - Use Replace-by-Fee for fee bumping if needed

**Reference:**

- [BIP125 - Replace-by-Fee](https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki)

---

## Wallet Implementation

### Wallet Types in This System

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, broadcast, monitor | Online |
| **Keygen** | Generate keys, first signature | Offline (air-gapped) |
| **Sign** | Additional signatures (multisig) | Offline (air-gapped) |

### Key Operations

| Operation | Wallet | Description |
|-----------|--------|-------------|
| `create seed` | Keygen | Generate BIP39 mnemonic |
| `create hdkey` | Keygen | Derive HD keys |
| `export address` | Keygen | Export addresses for import |
| `import address` | Watch | Import addresses to monitor |
| `create transaction` | Watch | Create unsigned PSBT |
| `sign` | Keygen/Sign | Add signature to PSBT |
| `send transaction` | Watch | Broadcast signed transaction |
| `monitor transaction` | Watch | Track confirmations |

### Account Types

| Account | Purpose | Multisig |
|---------|---------|----------|
| **client** | Customer deposit addresses | No |
| **deposit** | Aggregate client funds | No |
| **payment** | Outgoing payments | Yes (2-of-3 or 3-of-3) |
| **stored** | Cold storage | Yes |

---

## RPC & API Reference

### Bitcoin Core RPC

**Essential Commands:**

| Command | Description |
|---------|-------------|
| `getblockchaininfo` | Network and sync status |
| `getbalance` | Wallet balance |
| `listunspent` | List UTXOs |
| `createrawtransaction` | Create raw transaction |
| `signrawtransactionwithkey` | Sign with provided keys |
| `sendrawtransaction` | Broadcast transaction |
| `gettransaction` | Get transaction details |
| `walletprocesspsbt` | Process PSBT |
| `finalizepsbt` | Finalize PSBT |
| `decodepsbt` | Decode/analyze PSBT |

**Reference:**

- [Bitcoin Core RPC Documentation](https://developer.bitcoin.org/reference/rpc/)
- [Bitcoin Core JSON-RPC API Reference](https://bitcoincore.org/en/doc/)

### Go Libraries

| Library | Purpose | Repository |
|---------|---------|------------|
| **btcd** | Full node implementation | [github.com/btcsuite/btcd](https://github.com/btcsuite/btcd) |
| **btcutil** | Address/transaction utilities | [github.com/btcsuite/btcd/btcutil](https://github.com/btcsuite/btcd) |
| **btcec** | secp256k1 cryptography | [github.com/btcsuite/btcd/btcec](https://github.com/btcsuite/btcd) |
| **psbt** | PSBT implementation | [github.com/btcsuite/btcd/btcutil/psbt](https://github.com/btcsuite/btcd) |
| **txscript** | Script parsing/building | [github.com/btcsuite/btcd/txscript](https://github.com/btcsuite/btcd) |

---

## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Use air-gapped systems for key generation and signing
- Implement proper entropy for key generation
- Use hardware security modules (HSMs) for production

### Transaction Security

- Verify all transaction details before signing
- Implement multi-signature for high-value accounts
- Use PSBT for offline signing workflows
- Validate change addresses

### Nonce Security (MuSig2)

- **CRITICAL:** Never reuse nonces in MuSig2
- Generate cryptographically secure random nonces
- Delete nonces immediately after signing

### Network Security

- Use authenticated RPC connections
- Implement TLS for all network communication
- Validate transaction broadcasts

---

## Testing Resources

### Testnet Faucets

| Network | Faucet URL |
|---------|------------|
| **Testnet3** | [testnet-faucet.com](https://testnet-faucet.com/btc-testnet/) |
| **Signet** | [signetfaucet.com](https://signetfaucet.com/) |
| **Signet (Alt)** | [alt.signetfaucet.com](https://alt.signetfaucet.com/) |

### Block Explorers

| Network | Explorer |
|---------|----------|
| **Mainnet** | [mempool.space](https://mempool.space/) |
| **Mainnet** | [blockstream.info](https://blockstream.info/) |
| **Testnet3** | [mempool.space/testnet](https://mempool.space/testnet) |
| **Signet** | [mempool.space/signet](https://mempool.space/signet) |
| **Signet** | [explorer.bc-2.jp](https://explorer.bc-2.jp/) |

### Development Tools

| Tool | Purpose |
|------|---------|
| **Bitcoin Core** | Full node reference implementation |
| **btcdeb** | Bitcoin script debugger |
| **Sparrow Wallet** | Desktop wallet with PSBT support |
| **Electrum** | Lightweight wallet |

---

## Official References

### Bitcoin Improvement Proposals (BIPs)

#### Key Management & Addresses

| BIP | Title | Status |
|-----|-------|--------|
| [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) | HD Wallets | Final |
| [BIP39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) | Mnemonic Seed | Final |
| [BIP44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) | Multi-Account HD | Final |
| [BIP49](https://github.com/bitcoin/bips/blob/master/bip-0049.mediawiki) | P2SH-P2WPKH Derivation | Final |
| [BIP84](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki) | Native SegWit Derivation | Final |
| [BIP86](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki) | Taproot Derivation | Final |

#### SegWit & Taproot

| BIP | Title | Status |
|-----|-------|--------|
| [BIP141](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki) | SegWit Consensus | Final |
| [BIP143](https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki) | SegWit Signature Verification | Final |
| [BIP173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki) | Bech32 Addresses | Final |
| [BIP340](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki) | Schnorr Signatures | Final |
| [BIP341](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki) | Taproot | Final |
| [BIP342](https://github.com/bitcoin/bips/blob/master/bip-0342.mediawiki) | Tapscript | Final |
| [BIP350](https://github.com/bitcoin/bips/blob/master/bip-0350.mediawiki) | Bech32m Addresses | Final |

#### PSBT & Transactions

| BIP | Title | Status |
|-----|-------|--------|
| [BIP174](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki) | PSBT | Final |
| [BIP370](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki) | PSBT Version 2 | Draft |
| [BIP371](https://github.com/bitcoin/bips/blob/master/bip-0371.mediawiki) | Taproot PSBT Fields | Draft |
| [BIP125](https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki) | Replace-by-Fee | Final |

#### Multisig

| BIP | Title | Status |
|-----|-------|--------|
| [BIP11](https://github.com/bitcoin/bips/blob/master/bip-0011.mediawiki) | M-of-N Standard | Final |
| [BIP16](https://github.com/bitcoin/bips/blob/master/bip-0016.mediawiki) | P2SH | Final |
| [BIP327](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki) | MuSig2 | Draft |

### Official Documentation

- [Bitcoin Developer Documentation](https://developer.bitcoin.org/)
- [Bitcoin Core Documentation](https://bitcoincore.org/en/doc/)
- [Bitcoin Wiki](https://en.bitcoin.it/wiki/Main_Page)
- [Learn Me a Bitcoin](https://learnmeabitcoin.com/)
- [Bitcoin Optech](https://bitcoinops.org/)

### Academic Papers

- [Bitcoin Whitepaper](https://bitcoin.org/bitcoin.pdf) - Satoshi Nakamoto (2008)
- [MuSig2 Paper](https://eprint.iacr.org/2020/1261) - Nick, Ruffing, Seurin, Wuille (2020)
- [Schnorr Signatures for secp256k1](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)

---

## Project Documentation

### Implementation Guides

| Document | Description |
|----------|-------------|
| [BTC/BCH Technical Guide](btc_bch_technical_guide.md) | Comprehensive technical reference |
| [Taproot User Guide](TAPROOT_GUIDE.md) | Taproot address and transaction guide |
| [MuSig2 User Guide](musig2_guide.md) | MuSig2 multisig operations |
| [PSBT Developer Guide](psbt_developer_guide.md) | PSBT implementation details |
| [PSBT User Guide](psbt_user_guide.md) | PSBT usage guide |
| [Operation Examples](operation_example.md) | Wallet setup and operation |

### Architecture & Design

| Document | Description |
|----------|-------------|
| [Descriptor Architecture](./descriptor_architecture.md) | Output descriptor design |
| [MuSig2 Architecture](./musig2_architecture.md) | MuSig2 system architecture |
| [Key Generation Improvements](key_generation_improvements_2025.md) | 2025 key generation updates |
| [Wallet Flow Improvements](wallet_flow_improvements_2025.md) | 2025 workflow enhancements |

### Testing

| Document | Description |
|----------|-------------|
| [Taproot Testing](../testing/TAPROOT_TESTING.md) | Taproot test procedures |

---

## Version Information

| Component | Minimum Version | Recommended |
|-----------|-----------------|-------------|
| **Bitcoin Core** | v22.0 (Taproot) | v26.0+ |
| **btcd** | v0.24.0 | Latest |
| **Go** | 1.25 | 1.25+ |

---

**Document Version:** 2.0
**Last Updated:** 2026-01-07
**Maintainer:** go-crypto-wallet team
