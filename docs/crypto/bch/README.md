# Bitcoin Cash (BCH) Technical Reference

This document provides a comprehensive technical reference for Bitcoin Cash implementation in the go-crypto-wallet system. It covers specifications, protocol details, and links to official documentation to help AI agents and developers understand BCH's architecture and implement features correctly.

## Table of Contents

1. [Overview](#overview)
2. [Core Specifications](#core-specifications)
3. [Address Types & Key Derivation](#address-types--key-derivation)
4. [Transaction Architecture](#transaction-architecture)
5. [Signing Mechanisms](#signing-mechanisms)
6. [Multisig Implementation](#multisig-implementation)
7. [Network & Consensus](#network--consensus)
8. [Fee Management](#fee-management)
9. [Wallet Implementation](#wallet-implementation)
10. [RPC & API Reference](#rpc--api-reference)
11. [Differences from Bitcoin](#differences-from-bitcoin)
12. [Security Considerations](#security-considerations)
13. [Testing Resources](#testing-resources)
14. [Official References](#official-references)
15. [Project Documentation](#project-documentation)

---

## Overview

### What is Bitcoin Cash?

Bitcoin Cash (BCH) is a cryptocurrency that forked from Bitcoin (BTC) on August 1, 2017 (block 478,558). The fork was primarily motivated by the desire for larger block sizes to increase transaction throughput. BCH focuses on being a peer-to-peer electronic cash system, prioritizing low fees and fast confirmations.

### Key Characteristics (2026)

| Property | Value |
|----------|-------|
| **Fork Date** | August 1, 2017 (Block 478,558) |
| **Block Time** | ~10 minutes |
| **Block Size** | 32 MB |
| **Total Supply** | 21,000,000 BCH |
| **Current Block Reward** | 3.125 BCH (post-2024 halving) |
| **Next Halving** | ~2028 |
| **Consensus Algorithm** | SHA-256 Proof of Work |
| **Cryptographic Curve** | secp256k1 |
| **Signature Algorithm** | ECDSA only (no Schnorr) |
| **SegWit Support** | ❌ No |
| **Taproot Support** | ❌ No |

### Fork History

| Date | Event | Block Height |
|------|-------|--------------|
| 2017-08-01 | Bitcoin Cash fork from Bitcoin | 478,558 |
| 2017-11-13 | DAA (Difficulty Adjustment Algorithm) upgrade | 504,031 |
| 2018-05-15 | 32 MB block size, OP_RETURN 220 bytes | 530,356 |
| 2018-11-15 | BSV fork (contentious hard fork) | 556,766 |
| 2020-05-15 | IFP (Infrastructure Funding Plan) controversy | 635,258 |
| 2020-11-15 | BCHN becomes dominant implementation | 661,647 |
| 2023-05-15 | CashTokens upgrade | - |

### BCH vs BTC Comparison

| Feature | Bitcoin Cash (BCH) | Bitcoin (BTC) |
|---------|-------------------|---------------|
| **Block Size** | 32 MB | 1-4 MB (with SegWit) |
| **SegWit** | ❌ Not supported | ✅ Activated 2017 |
| **Taproot** | ❌ Not supported | ✅ Activated 2021 |
| **Schnorr Signatures** | ❌ Not available | ✅ BIP340 |
| **MuSig2** | ❌ Not available | ✅ BIP327 |
| **Address Format** | CashAddr | Bech32/Bech32m |
| **Transaction Malleability** | ⚠️ Still present | ✅ Fixed by SegWit |
| **Primary Focus** | Peer-to-peer cash | Store of value |
| **Average Fee** | Very low (<$0.01) | Variable (can be high) |

---

## Core Specifications

### Cryptographic Primitives

#### Elliptic Curve (secp256k1)

BCH uses the same secp256k1 elliptic curve as Bitcoin:

```
Curve Parameters:
- p = 2^256 - 2^32 - 977
- a = 0
- b = 7
- G = (0x79BE667E..., 0x483ADA77...)
- n = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFE BAAEDCE6AF48A03BBFD25E8CD0364141
```

#### Hash Functions

| Function | Usage |
|----------|-------|
| **SHA-256** | Block hashing, TXID calculation, PoW |
| **RIPEMD-160** | Address generation (hash160) |
| **HASH160** | SHA256 + RIPEMD160 for pubkey hashing |
| **HASH256** | Double SHA256 for transaction/block hashing |

### Data Encoding

| Format | Usage | Reference |
|--------|-------|-----------|
| **CashAddr** | Primary address format (recommended) | [CashAddr Spec](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md) |
| **Base58Check** | Legacy addresses (backward compatible) | [Base58Check](https://en.bitcoin.it/wiki/Base58Check_encoding) |
| **WIF** | Private key encoding | [Wallet Import Format](https://en.bitcoin.it/wiki/Wallet_import_format) |
| **Hex** | Raw transaction data | Standard hexadecimal |

### Network Magic Bytes

BCH uses different network magic bytes to distinguish from BTC:

| Network | BCH Magic | BTC Magic |
|---------|-----------|-----------|
| **Mainnet** | `0xe8f3e1e3` | `0xf9beb4d9` |
| **Testnet** | `0xf4f3e5f4` | `0x0b110907` |
| **Regtest** | `0xfabfb5da` | `0xfabfb5da` |

---

## Address Types & Key Derivation

### CashAddr Format (Recommended)

CashAddr was introduced in January 2018 to prevent accidental cross-chain sends between BCH and BTC.

#### CashAddr P2PKH

| Property | Value |
|----------|-------|
| **Prefix (Mainnet)** | `bitcoincash:q` |
| **Prefix (Testnet)** | `bchtest:q` |
| **Length** | ~42 characters (without prefix) |
| **Encoding** | Modified Bech32 |
| **Version Byte** | `0x00` |
| **Example** | `bitcoincash:qp3wjpa3tjlj042z2wv7hahsldgwhwy0rq9sywjpyy` |

**ScriptPubKey Structure:**

```
OP_DUP OP_HASH160 <20-byte pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
```

#### CashAddr P2SH

| Property | Value |
|----------|-------|
| **Prefix (Mainnet)** | `bitcoincash:p` |
| **Prefix (Testnet)** | `bchtest:p` |
| **Length** | ~42 characters (without prefix) |
| **Encoding** | Modified Bech32 |
| **Version Byte** | `0x08` |
| **Example** | `bitcoincash:pr0662zpd7vr936d83f64u629v886aan7c77r3j5v5` |

**ScriptPubKey Structure:**

```
OP_HASH160 <20-byte scriptHash> OP_EQUAL
```

### CashAddr Characteristics

- **Case-insensitive**: Always use lowercase (recommended)
- **Checksum**: Built-in error detection (Polymod checksum)
- **Network prefix**: Prevents cross-chain address confusion
- **Omittable prefix**: Can omit `bitcoincash:` when context is clear
- **Human-readable**: Clearly distinguishes BCH from BTC

### CashAddr Encoding

```
CashAddr = Prefix + ":" + Payload

Payload = Base32Encode(VersionByte + Hash + Checksum)

Base32 Alphabet: qpzry9x8gf2tvdw0s3jn54khce6mua7l
```

### Legacy Address Format (Backward Compatible)

BCH also supports legacy Base58Check addresses for backward compatibility:

| Type | Prefix (Mainnet) | Prefix (Testnet) |
|------|------------------|------------------|
| **P2PKH** | `1` | `m` or `n` |
| **P2SH** | `3` | `2` |

**⚠️ Warning**: Using legacy addresses risks accidental BTC/BCH cross-sends. Always prefer CashAddr.

### HD Wallet Derivation (BIP44)

BCH uses BIP44 with coin type 145:

```
Derivation Path: m/44'/145'/account'/change/index

- Purpose: 44' (BIP44)
- Coin Type: 145' (Bitcoin Cash per SLIP-0044)
- Account: 0' (first account)
- Change: 0 (external) or 1 (internal/change)
- Index: 0, 1, 2, ... (address index)
```

**Testnet uses coin type 1:**

```
Derivation Path: m/44'/1'/account'/change/index
```

**References:**

- [BIP32 - HD Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [BIP39 - Mnemonic Codes](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- [BIP44 - Multi-Account Hierarchy](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
- [SLIP-0044 - Coin Types](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)

---

## Transaction Architecture

### UTXO Model

BCH uses the same UTXO (Unspent Transaction Output) model as Bitcoin:

```
UTXO = {
    txid:      32-byte transaction hash
    vout:      output index (uint32)
    value:     satoshi amount (int64)
    scriptPubKey: locking script
}
```

### Transaction Structure

BCH uses legacy transaction format (no SegWit):

```
+--------------------+
| Version (4 bytes)  |  Usually 1 or 2
+--------------------+
| Input Count        |  VarInt
+--------------------+
| Inputs[]           |
|   - prevTxHash     |  32 bytes
|   - prevVout       |  4 bytes
|   - scriptSigLen   |  VarInt
|   - scriptSig      |  variable (contains signature)
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

### Transaction Size

Since BCH doesn't support SegWit, transaction sizes are larger than BTC SegWit transactions:

| Transaction Type | Size | Fee @ 1 sat/byte |
|------------------|------|------------------|
| P2PKH (1-in, 2-out) | ~226 bytes | ~226 sats |
| P2PKH (2-in, 2-out) | ~374 bytes | ~374 sats |
| 2-of-3 Multisig (1-in, 2-out) | ~370-400 bytes | ~400 sats |

### Sighash Types

BCH uses the same sighash types as legacy Bitcoin, with an additional fork ID:

| Type | Value | Description |
|------|-------|-------------|
| SIGHASH_ALL | 0x01 | Sign all inputs and outputs |
| SIGHASH_NONE | 0x02 | Sign all inputs, no outputs |
| SIGHASH_SINGLE | 0x03 | Sign all inputs, matching output only |
| SIGHASH_ANYONECANPAY | 0x80 | Modifier: sign only current input |
| SIGHASH_FORKID | 0x40 | BCH fork ID flag (added post-fork) |

**BCH Sighash Example:**

```
SIGHASH_ALL + SIGHASH_FORKID = 0x41
```

### Replay Protection

BCH implements replay protection using SIGHASH_FORKID to prevent transactions from being valid on both BCH and BTC chains.

---

## Signing Mechanisms

### ECDSA Signatures

BCH uses ECDSA (secp256k1) for all signatures. Unlike Bitcoin, BCH does NOT support Schnorr signatures.

**Signature Format (DER Encoded):**

```
0x30 [total-length]
  0x02 [r-length] [r]
  0x02 [s-length] [s]
[sighash-type]  (includes FORKID: 0x41)
```

**Size:** 71-73 bytes (variable due to DER encoding)

### BCH Signature Digest Algorithm

Post-fork BCH uses a modified signature digest algorithm (similar to BIP143) for replay protection:

```
Signature Digest:
1. nVersion (4 bytes)
2. hashPrevouts (32 bytes)
3. hashSequence (32 bytes)
4. outpoint (36 bytes)
5. scriptCode (variable)
6. value (8 bytes)
7. nSequence (4 bytes)
8. hashOutputs (32 bytes)
9. nLocktime (4 bytes)
10. sighash type (4 bytes) - includes FORKID
```

---

## Multisig Implementation

### Traditional P2SH Multisig

BCH supports only traditional P2SH multisig (no P2WSH or MuSig2):

**Redeem Script Structure:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

**Example 2-of-3:**

```
2 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

### Multisig Characteristics

| Property | Value |
|----------|-------|
| **Maximum Keys** | 15 (OP_CHECKMULTISIG limit) |
| **Address Type** | P2SH (CashAddr `p` prefix) |
| **Transaction Size** | ~370-400 bytes (2-of-3) |
| **Privacy** | Low (multisig visible on-chain) |
| **Signature Type** | ECDSA |

### Multisig Workflow

```
1. Create Redeem Script
   └── <M> <PubKey1> ... <PubKeyN> <N> OP_CHECKMULTISIG

2. Create P2SH Address
   └── Hash160(redeemScript) → CashAddr P2SH

3. Create Transaction
   └── Reference UTXO with P2SH scriptPubKey

4. Sign Transaction (Sequential)
   └── Signer 1 → Signer 2 → ... → Signer M

5. Assemble scriptSig
   └── OP_0 <Sig1> <Sig2> ... <SigM> <redeemScript>

6. Broadcast
   └── Send to network
```

### Limitations (vs Bitcoin)

| Feature | BCH | BTC |
|---------|-----|-----|
| P2SH Multisig | ✅ | ✅ |
| P2WSH Multisig | ❌ | ✅ |
| Taproot Multisig | ❌ | ✅ |
| MuSig2 | ❌ | ✅ |
| Privacy | Low | High (with MuSig2) |
| Fee Efficiency | Standard | 30-50% lower with MuSig2 |

---

## Network & Consensus

### Networks

| Network | Purpose | P2P Port | RPC Port |
|---------|---------|----------|----------|
| **Mainnet** | Production | 8333 | 8332 |
| **Testnet3** | Public testing | 18333 | 18332 |
| **Testnet4** | Public testing (newer) | 28333 | 28332 |
| **Regtest** | Local development | 18444 | 18443 |
| **Chipnet** | Contract testing | 48333 | 48332 |

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
| Transactions[]       |  Up to 32 MB
+----------------------+
```

### Difficulty Adjustment Algorithm (DAA)

BCH uses ASERT (Absolutely Scheduled Exponentially Rising Targets) DAA since November 2020:

- Adjusts difficulty every block
- Targets 10-minute block times
- Responds smoothly to hashrate changes

### Node Implementations

| Implementation | Description | Repository |
|----------------|-------------|------------|
| **Bitcoin Cash Node (BCHN)** | Primary implementation | [bitcoincashnode/bitcoincashnode](https://gitlab.com/bitcoin-cash-node/bitcoin-cash-node) |
| **Bitcoin Unlimited** | Alternative implementation | [BitcoinUnlimited/BitcoinUnlimited](https://github.com/BitcoinUnlimited/BitcoinUnlimited) |
| **BCHD** | Go implementation | [gcash/bchd](https://github.com/gcash/bchd) |
| **Flowee** | C++ implementation | [flowee-org/thehub](https://gitlab.com/FloweeTheHub/thehub) |

### Recommended Node Version

- **Bitcoin Cash Node**: v27.0.0+ (2026)
- Includes latest consensus rules and CashTokens support

---

## Fee Management

### Fee Characteristics

BCH typically has very low fees due to larger block size:

| Metric | Typical Value |
|--------|---------------|
| **Minimum Relay Fee** | 1 sat/byte |
| **Average Fee** | 1-2 sat/byte |
| **Fee per Transaction** | < 500 satoshis (~$0.002) |

### Fee Estimation

```bash
# Bitcoin Cash Node RPC
bitcoin-cli estimatefee <conf_target>
```

### Low Fee Benefits

- Suitable for micropayments
- Ideal for frequent, small transactions
- No significant fee market competition

---

## Wallet Implementation

### Implementation Architecture

In this codebase, BCH implementation **embeds and extends** Bitcoin:

```go
// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
    btc.Bitcoin
}

// Overrides chain parameters for BCH
func (b *BitcoinCash) initChainParams() {
    conf := b.GetChainConf()
    switch conf.Name {
    case chaincfg.TestNet3Params.Name:
        conf.Net = TestnetMagic
    case chaincfg.MainNetParams.Name:
        conf.Net = MainnetMagic
    }
}
```

### Shared Code with Bitcoin

The following functionality is shared:

- UTXO selection
- Transaction creation (non-SegWit)
- ECDSA signing
- Fee calculation
- Basic RPC communication

### BCH-Specific Overrides

The following methods are overridden for BCH:

- `GetAddressInfo()` - Different RPC response format
- `initChainParams()` - BCH network magic bytes
- Address encoding/decoding (CashAddr)

### Wallet Types

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, broadcast, monitor | Online |
| **Keygen** | Generate keys, first signature | Offline (air-gapped) |
| **Sign** | Additional signatures (multisig) | Offline (air-gapped) |

### Key Operations

| Operation | Description |
|-----------|-------------|
| `create seed` | Generate BIP39 mnemonic |
| `create hdkey` | Derive HD keys (BIP44 m/44'/145'/...) |
| `export address` | Export CashAddr addresses |
| `create transaction` | Create unsigned transaction |
| `sign` | Sign transaction with ECDSA |
| `send transaction` | Broadcast to BCH network |

---

## RPC & API Reference

### Bitcoin Cash Node RPC

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
| `getaddressinfo` | Get address information |
| `validateaddress` | Validate address format |

### BCH-Specific RPC Differences

| Command | BCH Behavior | Notes |
|---------|--------------|-------|
| `getaddressinfo` | Different response structure | No SegWit fields |
| `decodescript` | No witness fields | Legacy scripts only |

### Go Libraries

| Library | Purpose | Repository |
|---------|---------|------------|
| **gcash/bchd** | BCH full node in Go | [github.com/gcash/bchd](https://github.com/gcash/bchd) |
| **gcash/bchutil** | BCH utilities | [github.com/gcash/bchutil](https://github.com/gcash/bchutil) |
| **cpacia/bchutil** | CashAddr encoding | [github.com/cpacia/bchutil](https://github.com/cpacia/bchutil) |

### CashAddr Library Usage

```go
import "github.com/cpacia/bchutil"

// Encode to CashAddr
cashaddr, err := bchutil.EncodeAddress(hash160, &chaincfg.MainNetParams)

// Decode from CashAddr
decoded, err := bchutil.DecodeAddress(cashaddr, &chaincfg.MainNetParams)
```

---

## Differences from Bitcoin

### Feature Comparison

| Feature | Bitcoin (BTC) | Bitcoin Cash (BCH) |
|---------|---------------|-------------------|
| **Block Size** | 1-4 MB (with SegWit) | 32 MB |
| **SegWit** | ✅ Activated 2017 | ❌ Not implemented |
| **Taproot** | ✅ Activated 2021 | ❌ Not implemented |
| **Schnorr Signatures** | ✅ BIP340 | ❌ Not available |
| **MuSig2** | ✅ BIP327 | ❌ Not available |
| **Address Format** | Bech32/Bech32m | CashAddr |
| **PSBT** | ✅ BIP174 | ⚠️ Limited support |
| **RBF** | ✅ BIP125 | ⚠️ Different approach |
| **Transaction Malleability** | ✅ Fixed by SegWit | ⚠️ Still present |

### Implementation Differences in go-crypto-wallet

| Aspect | BTC | BCH |
|--------|-----|-----|
| **Transaction Format** | SegWit (with witness) | Legacy (no witness) |
| **Signature** | ECDSA or Schnorr | ECDSA only |
| **Multisig** | P2SH, P2WSH, MuSig2 | P2SH only |
| **Address Encoding** | Base58, Bech32, Bech32m | CashAddr, Base58 |
| **PSBT Support** | Full | Limited |

### Code Organization

```
internal/infrastructure/api/btc/
├── btc/          # Bitcoin implementation (base)
│   ├── bitcoin.go
│   ├── address.go
│   ├── psbt.go
│   └── ...
└── bch/          # Bitcoin Cash (extends btc)
    ├── bitcoin_cash.go   # Embeds btc.Bitcoin
    ├── address.go        # Override GetAddressInfo
    └── account.go
```

---

## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Use air-gapped systems for key generation and signing
- Implement proper entropy for key generation

### Address Security

- **Always use CashAddr** to prevent BTC/BCH cross-sends
- Validate address format before sending
- Double-check network prefix (`bitcoincash:` vs `bchtest:`)

### Transaction Security

- Verify all transaction details before signing
- Implement multi-signature for high-value accounts
- Validate change addresses

### Replay Protection

- BCH transactions include SIGHASH_FORKID
- Transactions are not valid on BTC chain
- Always verify fork ID in signatures

### 51% Attack Considerations

BCH has lower hashrate than BTC, making it theoretically more vulnerable to 51% attacks:

- Wait for more confirmations for large transactions
- Consider 10+ confirmations for high-value transfers

---

## Testing Resources

### Testnet Faucets

| Network | Faucet URL |
|---------|------------|
| **Testnet3** | [tbch.googol.cash](https://tbch.googol.cash/) |
| **Testnet4** | [tbch4.googol.cash](https://tbch4.googol.cash/) |

### Block Explorers

| Network | Explorer |
|---------|----------|
| **Mainnet** | [blockchair.com/bitcoin-cash](https://blockchair.com/bitcoin-cash) |
| **Mainnet** | [explorer.bitcoin.com](https://explorer.bitcoin.com/bch) |
| **Mainnet** | [bch.loping.net](https://bch.loping.net/) |
| **Testnet** | [tbch.loping.net](https://tbch.loping.net/) |
| **Chipnet** | [chipnet.chaingraph.cash](https://chipnet.chaingraph.cash/) |

### Development Tools

| Tool | Purpose |
|------|---------|
| **Bitcoin Cash Node** | Full node implementation |
| **BCHD** | Go full node |
| **Electron Cash** | Desktop wallet |
| **Bitbox SDK** | JavaScript development kit |

---

## Official References

### Specifications

| Document | Description | Link |
|----------|-------------|------|
| **CashAddr Spec** | Address format specification | [cashaddr.md](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md) |
| **Transaction Format** | BCH transaction specification | [transaction.md](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/transaction.md) |
| **OP_RETURN** | Data carrier outputs | [op_return.md](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/op_return.md) |
| **CashTokens** | Token specification | [cashtokens.md](https://github.com/cashtokens/cashtokens) |

### CHIPs (Cash Improvement Proposals)

| CHIP | Title | Status |
|------|-------|--------|
| [CHIP-2021-01](https://gitlab.com/GeneralProtocols/research/chips) | Restrict Transaction Version | Deployed |
| [CHIP-2021-05](https://gitlab.com/GeneralProtocols/research/chips) | VM Limits | Deployed |
| [CHIP-2022-02](https://gitlab.com/GeneralProtocols/research/chips) | CashTokens | Deployed |
| [CHIP-2022-05](https://gitlab.com/GeneralProtocols/research/chips) | P2SH32 | Deployed |

### Documentation Resources

- [Bitcoin Cash Documentation](https://documentation.cash/)
- [Bitcoin Cash Research](https://bitcoincashresearch.org/)
- [CashScript Documentation](https://cashscript.org/)

### Relevant BIPs (Shared with Bitcoin)

| BIP | Title | Applicability |
|-----|-------|---------------|
| [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) | HD Wallets | ✅ Used |
| [BIP39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) | Mnemonic Codes | ✅ Used |
| [BIP44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) | Multi-Account HD | ✅ Used (coin type 145) |
| [BIP11](https://github.com/bitcoin/bips/blob/master/bip-0011.mediawiki) | M-of-N Multisig | ✅ Used |
| [BIP16](https://github.com/bitcoin/bips/blob/master/bip-0016.mediawiki) | P2SH | ✅ Used |
| [BIP141](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki) | SegWit | ❌ Not used |
| [BIP174](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki) | PSBT | ⚠️ Limited |
| [BIP340](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki) | Schnorr | ❌ Not used |

---

## Project Documentation

### Related Guides

| Document | Description |
|----------|-------------|
| [BTC/BCH Technical Guide](../btc/btc_bch_technical_guide.md) | Comprehensive BTC/BCH comparison |
| [Bitcoin README](../btc/README.md) | Bitcoin technical reference |

### Implementation Files

| File | Description |
|------|-------------|
| `internal/infrastructure/api/btc/bch/bitcoin_cash.go` | Main BCH implementation |
| `internal/infrastructure/api/btc/bch/address.go` | Address handling override |
| `internal/domain/coin/types.go` | Coin type definitions (BCH = 145) |

---

## Quick Reference

### Address Format Summary

| Type | Mainnet Prefix | Testnet Prefix | Example |
|------|----------------|----------------|---------|
| P2PKH (CashAddr) | `bitcoincash:q` | `bchtest:q` | `bitcoincash:qp3wjpa3tjlj042z2wv7hahsldgwhwy0rq9sywjpyy` |
| P2SH (CashAddr) | `bitcoincash:p` | `bchtest:p` | `bitcoincash:pr0662zpd7vr936d83f64u629v886aan7c77r3j5v5` |
| P2PKH (Legacy) | `1` | `m`/`n` | `1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2` |
| P2SH (Legacy) | `3` | `2` | `3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy` |

### Network Parameters

| Parameter | Mainnet | Testnet |
|-----------|---------|---------|
| Network Magic | `0xe8f3e1e3` | `0xf4f3e5f4` |
| P2P Port | 8333 | 18333 |
| RPC Port | 8332 | 18332 |
| Coin Type (SLIP44) | 145 | 1 |
| Address Prefix | `bitcoincash:` | `bchtest:` |

### Transaction Checklist

- [ ] Use CashAddr format for addresses
- [ ] Include SIGHASH_FORKID (0x40) in signatures
- [ ] Set appropriate fee (typically 1 sat/byte)
- [ ] Verify replay protection is active
- [ ] Wait for sufficient confirmations

---

**Document Version:** 1.0
**Last Updated:** 2026-01-07
**Maintainer:** go-crypto-wallet team
