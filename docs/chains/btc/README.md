# Bitcoin (BTC) Technical Reference

This document provides a comprehensive technical reference for Bitcoin implementation in the go-crypto-wallet system. It covers specifications, protocol details, and links to official documentation to help AI agents and developers understand Bitcoin's architecture and implement features correctly.

## Documentation Structure

This directory is organized into the following categories:

| File / Directory | Description | Audience |
|------------------|-------------|----------|
| [architecture.md](../../../template/sections/misc/architecture.md) | **Wallet architecture** — wallet roles, use case boundary map, Keygen vs Sign signing | Developers |
| [overview/](../../../template/sections/misc/overview/README) | Fundamental technical references and Bitcoin basics | All |
| [operations/](../../../template/sections/misc/operations/README) | Wallet operation guides and transaction flows | Operators |
| [keygen/](../../../template/sections/misc/keygen/README) | Key generation design and improvements | Developers |
| [psbt/](../../../template/sections/misc/psbt/README) | PSBT implementation and usage guides | All |
| [descriptor/](../../../template/sections/misc/descriptor/README) | Output Descriptor implementation | Developers |
| [taproot/](../../../template/sections/misc/taproot/README) | Taproot (BIP341/BIP86) guides | All |
| [musig2/](../../../template/sections/misc/musig2/README) | MuSig2 multisignature implementation | All |
| archive/ | Outdated documentation (reference only) | - |

## Quick Start

### For Operators

1. Start with [operations/wallet-flow.md](../../../template/sections/misc/operations/wallet-flow.md) for wallet setup and transaction flows
2. Review [operations/e2e-transaction-patterns.md](../../../template/sections/misc/operations/e2e-transaction-patterns.md) for transaction types
3. See [psbt/user-guide.md](../../../template/sections/misc/psbt/user-guide.md) for offline signing workflows

### For Developers

1. Read [architecture.md](../../../template/sections/misc/architecture.md) for the wallet boundary map and use case assignments
2. Read [overview/technical-reference.md](../../../template/sections/misc/overview/technical-reference.md) for Bitcoin protocol fundamentals
3. Review feature-specific architecture docs: [descriptor/architecture.md](../../../template/sections/misc/descriptor/architecture.md), [musig2/architecture.md](../../../template/sections/misc/musig2/architecture.md)
4. Check [psbt/developer-guide.md](../../../template/sections/misc/psbt/developer-guide.md) for PSBT implementation details

---

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

See [overview/address-types.md](../../../template/sections/misc/overview/address-types.md) for detailed comparison.

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

**Sighash Types:**

| Type | Value | Description |
|------|-------|-------------|
| SIGHASH_ALL | 0x01 | Sign all inputs and outputs (default) |
| SIGHASH_NONE | 0x02 | Sign all inputs, no outputs |
| SIGHASH_SINGLE | 0x03 | Sign all inputs, matching output only |
| SIGHASH_ANYONECANPAY | 0x80 | Modifier: sign only current input |

### Schnorr Signatures (Taproot)

Used for P2TR transactions. Introduced with Taproot (BIP340).

**Advantages:**

- Fixed 64-byte size (vs variable ECDSA)
- Linear: enables signature aggregation (MuSig2)
- Provably secure under standard assumptions
- Batch verification is faster

See [taproot/user-guide.md](../../../template/sections/misc/taproot/user-guide.md) for details.

---

## Multisig & MuSig2

### Traditional Multisig (P2SH/P2WSH)

**M-of-N Redeem Script:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

### MuSig2 (Schnorr Signature Aggregation)

MuSig2 enables N-of-N multisig that appears as single-sig on-chain.

**Benefits:**

- 30-50% smaller transactions
- Maximum privacy (looks like single-sig)
- Lower fees
- No on-chain multisig indicator

**Critical Security: NONCE MANAGEMENT**

- **NEVER reuse nonces** - reusing leaks private key
- Generate fresh nonces for every transaction
- Delete nonces after signing

See [musig2/](../../../template/sections/misc/musig2/README) for detailed documentation.

---

## PSBT (Partially Signed Bitcoin Transactions)

PSBT (BIP174) is the standard format for offline/multi-party signing workflows.

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

See [psbt/](../../../template/sections/misc/psbt/README) for detailed documentation.

---

## Network & Consensus

### Networks

| Network | Purpose | Port | RPC Port | Magic Bytes |
|---------|---------|------|----------|-------------|
| **Mainnet** | Production | 8333 | 8332 | 0xF9BEB4D9 |
| **Testnet3** | Public testing | 18333 | 18332 | 0x0B110907 |
| **Signet** | Controlled testing | 38333 | 38332 | 0x0A03CF40 |
| **Regtest** | Local development | 18444 | 18443 | 0xFABFB5DA |

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

---

## Wallet Implementation

### Wallet Types in This System

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, broadcast, monitor | Online |
| **Keygen** | Generate keys, first signature | Offline (air-gapped) |
| **Sign** | Additional signatures (multisig) | Offline (air-gapped) |

### Account Types

| Account | Purpose | Multisig |
|---------|---------|----------|
| **client** | Customer deposit addresses | No |
| **deposit** | Aggregate client funds | No |
| **payment** | Outgoing payments | Yes (2-of-3 or 3-of-3) |
| **stored** | Cold storage | Yes |

For the common 3-wallet transaction flow (chain-agnostic), see [docs/transaction-flow.md](../../../template/transaction-flow.md).
For BTC-specific procedures and Mermaid diagrams, see [operations/wallet-flow.md](../../../template/sections/misc/operations/wallet-flow.md).

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

See [musig2/security.md](../../../template/sections/misc/musig2/security.md) for details.

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

### By Category

| Category | Documents |
|----------|-----------|
| **Overview** | [technical-reference.md](../../../template/sections/misc/overview/technical-reference.md), [address-types.md](../../../template/sections/misc/overview/address-types.md) |
| **Operations** | [wallet-flow.md](../../../template/sections/misc/operations/wallet-flow.md), [e2e-transaction-patterns.md](../../../template/sections/misc/operations/e2e-transaction-patterns.md), [wallet-flow-improvements-2025.md](../../../template/sections/misc/operations/wallet-flow-improvements-2025.md) |
| **Key Generation** | [improvements-2025.md](../../../template/sections/misc/keygen/improvements-2025.md), [interface-design.md](../../../template/sections/misc/keygen/interface-design.md) |
| **PSBT** | [user-guide.md](../../../template/sections/misc/psbt/user-guide.md), [developer-guide.md](../../../template/sections/misc/psbt/developer-guide.md), [implementation.md](../../../template/sections/misc/psbt/implementation.md) |
| **Descriptor** | [user-guide.md](../../../template/sections/misc/descriptor/user-guide.md), [architecture.md](../../../template/sections/misc/descriptor/architecture.md), [api.md](../../../template/sections/misc/descriptor/api.md) |
| **Taproot** | [user-guide.md](../../../template/sections/misc/taproot/user-guide.md), [testing.md](../../../template/sections/misc/taproot/testing.md) |
| **MuSig2** | [user-guide.md](../../../template/sections/misc/musig2/user-guide.md), [architecture.md](../../../template/sections/misc/musig2/architecture.md), [security.md](../../../template/sections/misc/musig2/security.md) |
| **Testing** | [pattern3-verification.md](../../../template/sections/misc/operations/pattern3-verification.md) |

### Related Resources

| Resource | Location |
|----------|----------|
| E2E Test Scripts | [scripts/operation/btc/e2e/](https://github.com/hiromaily/go-crypto-wallet/tree/main/scripts/operation/btc/e2e) |
| Project Testing Standards | [docs/guidelines/testing.md](../../../template/guidelines/testing.md) |
| Security Standards | [docs/guidelines/security.md](../../../template/guidelines/security.md) |

---

## Version Information

| Component | Minimum Version | Recommended |
|-----------|-----------------|-------------|
| **Bitcoin Core** | v22.0 (Taproot) | v26.0+ |
| **btcd** | v0.24.0 | Latest |
| **Go** | 1.25 | 1.25+ |

---

**Document Version:** 3.0
**Last Updated:** 2026-01-16
**Maintainer:** go-crypto-wallet team
