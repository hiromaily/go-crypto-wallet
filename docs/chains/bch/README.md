<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/bch/README.tpl.md · Run `make docs` to regenerate.
-->

## Table of Contents

1. [Overview](#overview)
2. [Architecture Layers](#architecture-layers)
3. [Component Interactions](#component-interactions)
4. [Data Flow](#data-flow)
5. [Security Architecture](#security-architecture)
6. [Database Schema](#database-schema)
7. [API Reference](#api-reference)
8. [Implementation Notes](#implementation-notes)

---

## Overview

### What is MuSig2?

MuSig2 is a two-round Schnorr multisignature protocol (BIP327) that enables multiple parties to create a single aggregated signature that is indistinguishable from a standard single-signature transaction on the blockchain. This provides:

- **Smaller transactions**: 30-50% size reduction compared to traditional P2WSH multisig
- **Lower fees**: Proportional to transaction size reduction
- **Better privacy**: Multisig transactions look like single-sig on-chain
- **Schnorr signatures**: Uses BIP340 Schnorr signatures via Taproot (P2TR)

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│   Interface Adapters (CLI Commands)                     │
│   - keygen create musig2-address                        │
│   - keygen musig2 nonce                                 │
│   - keygen musig2 sign                                  │
│   - sign musig2 nonce                                   │
│   - sign musig2 sign                                    │
│   - watch musig2 aggregate                              │
└────────────────────┬────────────────────────────────────┘
                     │ depends on
┌────────────────────▼────────────────────────────────────┐
│   Application Layer (Use Cases)                         │
│   Keygen:  CreateMuSig2AddressUseCase                   │
│            GenerateMuSig2NonceUseCase                   │
│            MuSig2SignUseCase                            │
│   Sign:    GenerateMuSig2NonceUseCase                   │
│            MuSig2SignUseCase                            │
│   Watch:   AggregateMuSig2SignaturesUseCase             │
└────────────────────┬────────────────────────────────────┘
                     │ depends on
┌────────────────────▼────────────────────────────────────┐
│   Domain Layer (Business Logic)                         │
│   - MuSig2 Types (domain/musig2/)                       │
│   - Validators                                          │
│   - Business Rules                                      │
└─────────────────────────────────────────────────────────┘
                     ▲ implements
┌────────────────────┴────────────────────────────────────┐
│   Infrastructure Layer (External Dependencies)          │
│   - MuSig2Service (btcd/btcec/v2/schnorr/musig2)       │
│   - AccountKeyRepository (MySQL)                        │
│   - AuthFullPubkeyRepository (MySQL)                    │
│   - FileStorage (PSBT files)                            │
└─────────────────────────────────────────────────────────┘
```

### Design Principles

1. **Clean Architecture**: Strict layer separation with dependency inversion
2. **Security First**: Nonce uniqueness enforced at multiple levels
3. **Type Safety**: Domain types for all MuSig2 operations
4. **Testability**: All components have clear interfaces
5. **Offline Support**: Keygen and Sign wallets work completely offline
6. **PSBT Integration**: MuSig2 data stored in PSBT proprietary fields

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

## Known Issues and Workarounds

### BCH Node RPC Bug: Multisig `complete` Flag

#### Issue Description

Bitcoin Cash Node has a bug in the `signrawtransactionwithkey` RPC method where it incorrectly reports `complete=true` after adding a partial signature to multisig transactions:

| Multisig Type | Expected Behavior | Actual BCH Behavior | Impact |
|---------------|-------------------|---------------------|--------|
| 2-of-3 | `complete=false` after 1st sig, `complete=true` after 2nd sig | Same (correct) | None |
| 3-of-3 | `complete=false` after 1st and 2nd sig, `complete=true` after 3rd sig | `complete=true` after 2nd sig (incorrect) | Subsequent signers cannot add their signatures |

#### Root Cause

The BCH Node's `signrawtransactionwithkey` RPC method incorrectly evaluates transaction completeness for multisig transactions. When adding the 2nd signature to a 3-of-3 multisig transaction, it reports `complete=true` even though a 3rd signature is still required.

#### Symptoms

1. **Transaction File Type**: After 2nd signature in 3-of-3 multisig, transaction file is marked as `signed` instead of remaining `unsigned`
2. **Missing prevTx Metadata**: When marked as `signed`, prevTx metadata (TXID, vout, scriptPubKey, redeemScript, amount) is omitted from the transaction file
3. **Sign2 Wallet Failure**: Without prevTx metadata, Sign2 wallet cannot add the 3rd signature, resulting in error:

   ```
   Error: fail to sign transaction: fail to sign raw transaction with auth key:
   result of sign raw transaction includes error: Input not found or already spent
   ```

#### Workaround Implementation

**Issue #485** implemented the following workarounds:

##### 1. Accept Both `unsigned` and `signed` Transaction Types (Sign Wallet)

Modified `internal/application/usecase/sign/bch/sign_transaction.go` to accept both transaction types:

```go
// For BCH multisig, accept both "unsigned" and "signed" files since BCH may mark
// a transaction as "signed" even when it needs additional signatures for multisig
if fileType.TxType != domainTx.TxTypeUnsigned && fileType.TxType != domainTx.TxTypeSigned {
    return signusecase.SignTransactionOutput{}, fmt.Errorf(
        "invalid txType: %s (expected unsigned or signed)", fileType.TxType)
}
```

##### 2. Always Include prevTx Metadata for Multisig (Keygen and Sign Wallets)

Modified both keygen and sign wallets to always include prevTx metadata for multisig transactions, regardless of the `complete` flag:

**Keygen Wallet** (`internal/application/usecase/keygen/bch/sign_transaction.go`):

```go
var generatedFileName string
// For multisig transactions, always include prevTx metadata for subsequent signers,
// even if BCH incorrectly reports isSigned=true for partial signatures
isMultisig := u.multisigAccount != nil
if isSigned && !isMultisig {
    // For fully signed single-sig transactions, just write the hex
    generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, signedHex)
} else {
    // For partially signed or multisig, include prevTx metadata for next signer
    content := u.formatSignedTxContent(signedHex, prevTxs)
    generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, content)
}
```

**Sign Wallet** (`internal/application/usecase/sign/bch/sign_transaction.go`):

```go
var generatedFileName string
// For multisig transactions, always include prevTx metadata for subsequent signers,
// even if BCH incorrectly reports isSigned=true for partial signatures
isMultisig := u.multisigAccount != nil
if isSigned && !isMultisig {
    // For fully signed single-sig transactions, just write the hex
    generatedFileName, err = u.txFileRepo.WriteFile(path, signedHex)
} else {
    // For partially signed or multisig, include prevTx metadata for next signer
    content := u.formatSignedTxContent(signedHex, prevTxs)
    generatedFileName, err = u.txFileRepo.WriteFile(path, content)
}
```

##### 3. E2E Script Fixes

Fixed duplicate `--wallet` flags in E2E scripts:

- **Issue**: RPC host already included wallet name (e.g., `127.0.0.1:30332/wallet/sign1-p2`)
- **Additional `--wallet` flag** created invalid path (e.g., `127.0.0.1:30332/wallet/sign1-p2/wallet/sign1`)
- **Solution**: Removed all redundant `--wallet` flags from `e2e-p2-p2sh-2of3.sh` and `e2e-p3-p2sh-3of3.sh`

#### Testing Results

After implementing the workarounds:

| Pattern | Description | Status | Transaction ID |
|---------|-------------|--------|----------------|
| Pattern 2 | 2-of-3 Multisig | ✅ Pass | `62c4331beddc2279c1567e41a23cca16d0ed4e25b376578ee4a3980c3f6cb240` |
| Pattern 3 | 3-of-3 Multisig | ✅ Pass | `4dee6d3978176f02ce01a9b4dacbfc01e84f662d1a4ab8ab8e7bb00c3a413a0c` |

#### Impact on Development

- **Multisig Detection**: Use non-nil `multisigAccount` field to determine if transaction requires multisig handling
- **prevTx Metadata**: Always include for multisig, regardless of `complete` flag
- **Transaction Type Validation**: Sign wallets must accept both `unsigned` and `signed` transaction files
- **E2E Scripts**: Avoid redundant `--wallet` flags when RPC host already includes wallet name

#### Future Considerations

This workaround should be maintained until Bitcoin Cash Node fixes the `complete` flag behavior. Monitor BCH Node releases for:

- Fixes to `signrawtransactionwithkey` multisig evaluation
- Changes to transaction completeness logic
- Updates to RPC response format

#### Related Issues

- GitHub Issue #485: "Bitcoin Core wallet not loaded when importing private keys"
- Related commits:
  - `69a2df37` - Initial Pattern 2 fix
  - `1902e9b2` - Extended fix to Pattern 3

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

## E2E Transaction Patterns

This section describes the E2E (End-to-End) transaction patterns available for Bitcoin Cash in the go-crypto-wallet system.

> **Common flow reference**: The 3-wallet setup, signing, and monitoring flows are defined in
> [docs/transaction-flow.md](../../transaction-flow.md). The patterns below describe
> BCH-specific address formats, signing algorithms, and protocol constraints.

### BCH Pattern Limitations vs BTC

Due to protocol differences, BCH supports **fewer patterns** than BTC:

| Feature | BTC Support | BCH Support | Impact |
|---------|-------------|-------------|--------|
| **SegWit** | ✅ Activated 2017 | ❌ Not implemented | No P2WPKH, P2WSH patterns |
| **Taproot** | ✅ Activated 2021 | ❌ Not implemented | No P2TR patterns |
| **Schnorr Signatures** | ✅ BIP340 | ❌ Not available | ECDSA only |
| **MuSig2** | ✅ BIP327 | ❌ Not available | No signature aggregation |
| **PSBT** | ✅ BIP174 | ⚠️ Limited | Raw transaction format |

### BCH E2E Pattern Matrix

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Status |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **CashAddr P2PKH** | **Single-sig** | **`bitcoincash:q...`** | **🔶 Manual testing** |
| **2** | **CashAddr P2SH** | **2-of-3 Multisig** | **`bitcoincash:p...`** | **❌ Not implemented** |
| **3** | **CashAddr P2SH** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e-workflow.sh** |

### Pattern 1: BCH CashAddr P2PKH Single-sig

```
Address Type: CashAddr P2PKH (BIP44)
Signing Requirements: Single-sig (Keygen only)
Address Format: bitcoincash:q... (P2PKH), bchtest:q... (Testnet)
```

**Workflow:**

Follows the [common single-sig flow](../../transaction-flow.md#single-sig-flow).
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- Simple and fast (completed with single signature)
- Sign1/Sign2 wallets not required
- Uses BIP44 key derivation path (`m/44'/145'/account'/change/index` for mainnet)
- CashAddr P2PKH format (`bitcoincash:q...` prefix)
- Suitable for `client` account type

**Key Derivation:**

```
Mainnet: m/44'/145'/account'/change/index
Testnet/Regtest: m/44'/1'/account'/change/index
```

**Implementation Status:** 🔶 Manual testing (no dedicated E2E script)

### Pattern 2: BCH CashAddr P2SH 2-of-3 Multisig

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 2-of-3 Multisig (any 2 of Keygen, Sign1, Sign2)
Address Format: bitcoincash:p... (P2SH), bchtest:p... (Testnet)
```

**Workflow:**

Follows the [common multisig flow](../../transaction-flow.md#multisig-flow-m-of-n) with M=2, N=3.
Signing stops after 2 signatures (`isCompleted: true`); Sign2 is not required.
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- Requires 2 of 3 possible signatures
- Provides redundancy - losing one key doesn't prevent access
- Uses BIP44 key derivation with BIP11 multisig
- CashAddr P2SH format (`bitcoincash:p...` prefix)
- Suitable for `deposit`, `payment` accounts

**Redeem Script:**

```
2 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

**Implementation Status:** ❌ Not implemented (planned)

### Pattern 3: BCH CashAddr P2SH 3-of-3 Multisig (Current E2E)

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 3-of-3 Multisig (Keygen + Sign1 + Sign2)
Address Format: bitcoincash:p... (P2SH), bchtest:p... (Testnet)
E2E Script: scripts/operation/bch/e2e-workflow.sh
```

**Workflow:**

Follows the [common multisig flow](../../transaction-flow.md#multisig-flow-m-of-n) with M=3, N=3.
All 3 signatures required (Keygen + Sign1 + Sign2).
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- All 3 signatures required (maximum security)
- Uses BIP44 key derivation with BIP11 multisig
- CashAddr P2SH format (`bitcoincash:p...` prefix)
- Suitable for `stored` account type (cold storage)

**Redeem Script:**

```
3 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

**Implementation Status:** ✅ Implemented in `scripts/operation/bch/e2e-workflow.sh`

### BCH vs BTC Pattern Comparison

| BTC Pattern | BCH Equivalent | Notes |
|-------------|----------------|-------|
| Pattern 1 (P2PKH Single-sig) | **BCH Pattern 1** | Same concept, CashAddr format |
| Pattern 2 (P2PKH 2-of-3) | **BCH Pattern 2** | Same concept, CashAddr format |
| Pattern 3 (P2SH-P2WPKH) | ❌ N/A | BCH has no SegWit |
| Pattern 4 (P2SH-P2WSH 2-of-3) | ❌ N/A | BCH has no SegWit |
| Pattern 5 (P2WPKH) | ❌ N/A | BCH has no SegWit |
| Pattern 6 (P2WSH 2-of-3) | ❌ N/A | BCH has no SegWit |
| Pattern 7 (P2WSH 3-of-3) | **BCH Pattern 3** | Similar but no SegWit benefits |
| Pattern 8 (P2SH-P2WSH 3-of-3) | **BCH Pattern 3** | Similar but no SegWit benefits |
| Pattern 9 (P2TR Single-sig) | ❌ N/A | BCH has no Taproot |
| Pattern 10 (P2TR MuSig2) | ❌ N/A | BCH has no Taproot/Schnorr |
| Pattern 11 (P2TR Tapscript) | ❌ N/A | BCH has no Taproot |

### Account Types and Recommended Patterns

| Account | Purpose | Recommended BCH Pattern | Reason |
|---------|---------|-------------------------|--------|
| **client** | Customer deposit addresses | Pattern 1 (Single-sig) | Simple, customer-side operations |
| **deposit** | Deposit aggregation | Pattern 2 or 3 (Multisig) | Enhanced security |
| **payment** | Payments | Pattern 2 or 3 (Multisig) | Approval workflow |
| **stored** | Long-term storage | Pattern 3 (3-of-3) | Highest security |

### E2E Script Reference

| Script | Pattern | Description |
|--------|---------|-------------|
| `scripts/operation/bch/e2e-workflow.sh` | Pattern 3 | 3-of-3 Multisig workflow |

**Usage:**

```bash
# Run complete E2E workflow
./scripts/operation/bch/e2e-workflow.sh

# Run with reset (fresh state)
./scripts/operation/bch/e2e-workflow.sh --reset

# Run with verbose output
./scripts/operation/bch/e2e-workflow.sh --verbose

# Cleanup only
./scripts/operation/bch/e2e-workflow.sh --cleanup
```

### Future BCH Upgrades (2025-2026)

While not currently implemented in go-crypto-wallet, BCH has planned upgrades:

| Upgrade | CHIP | Description | go-crypto-wallet Impact |
|---------|------|-------------|-------------------------|
| **CashTokens** | CHIP-2022-02 | Native fungible/NFT tokens | Not in scope |
| **VM Limits** | CHIP-2021-05 | Extended script capabilities | Not in scope |
| **BigInt** | CHIP-2024-07 | High-precision arithmetic | Not in scope |
| **Loops** | CHIP-2021-05 | Script loops (OP_BEGIN/OP_UNTIL) | Not in scope |
| **Functions** | CHIP-2025-05 | Reusable script functions | Not in scope |
| **Pay-2-Script** | CHIP-2024-12 | Direct P2S outputs | Not in scope |

**Note:** go-crypto-wallet focuses on standard transaction patterns. Advanced contract features are out of scope.

### Implementation Priority

| Priority | Pattern | Reason |
|----------|---------|--------|
| ✅ Completed | Pattern 3 (3-of-3 Multisig) | Security-first approach |
| 🔶 Manual | Pattern 1 (Single-sig) | Basic functionality available |
| 🔜 Planned | Pattern 2 (2-of-3 Multisig) | Operational flexibility |

---

## Project Documentation

### Related Guides

| Document | Description |
|----------|-------------|
| [BTC E2E Transaction Patterns](../btc/operations/e2e-transaction-patterns.md) | Comprehensive BTC patterns guide |
| BTC/BCH Technical Guide (internal document) | Comprehensive BTC/BCH comparison |
| [Bitcoin README](../btc/README.md) | Bitcoin technical reference |
| [BCH E2E Workflow](https://github.com/hiromaily/go-crypto-wallet/blob/main/scripts/operation/bch/README.md) | BCH E2E script documentation |

### Implementation Files

| File | Description |
|------|-------------|
| `internal/infrastructure/api/btc/bch/bitcoin_cash.go` | Main BCH implementation |
| `internal/infrastructure/api/btc/bch/address.go` | Address handling override |
| `internal/infrastructure/api/btc/bch/account.go` | Account handling override |
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

## Changelog

### Version 1.10 (2026-01-17)

- ✅ Fixed BCH Pattern 3 (3-of-3 Multisig) UTXO retrieval issue (Closes #423, PR #426)
- Resolved CashAddr format mismatch in `ListUnspentByAccount` address comparison
- Updated BCH Pattern Matrix to show all 3 patterns with correct E2E script paths
- Updated E2E Script Reference with proper BCH script locations
- Added implementation status details for BCH Pattern 3
- All BCH patterns now have dedicated E2E scripts in `scripts/operation/bch/e2e/`

### Version 1.9 (2026-01-17)

- Updated BCH Pattern Matrix with accurate pattern numbering (1, 2, 3)
- Added BCH limitations section (no SegWit, Taproot, Schnorr, MuSig2)
- Updated BCH Pattern 3 (3-of-3 Multisig) with detailed workflow and signing flow
- Added cross-reference link to BCH Technical Reference
- Renamed "BCH Pattern 2" to "BCH Pattern 3" to reflect correct pattern numbering

### Version 1.8 (2026-01-16)

- ✅ Pattern 11 (P2TR Tapscript M-of-N) framework implemented
- E2E script `e2e-p11-p2tr-tapscript.sh` created (Closes #381)
- Implements Tapscript Script Path spending framework with 2-of-3 threshold
- Uses BIP86 key derivation + BIP342 Tapscript semantics
- Script tree with Merkle proof and control block structure
- M × Schnorr signatures for Script Path spend
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)
- ~50% smaller than P2WSH 2-of-3 multisig
- Enhanced privacy: unused script paths hidden in Merkle tree
- Note: Full Tapscript implementation pending (currently uses placeholder)

### Version 1.7 (2026-01-16)

- ✅ Pattern 9 (P2TR Taproot Single-sig) is now fully working
- E2E script `e2e-p9-p2tr-singlesig.sh` completed and verified (Closes #377)
- Fixed Taproot address derivation to use x-only public keys (32 bytes)
- Implements BIP86 key derivation for Taproot key path spending
- Uses Schnorr signatures (BIP340, 64 bytes)
- Most efficient single-sig transaction format
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)

### Version 1.6 (2026-01-16)

- ✅ Pattern 8 (P2SH-P2WSH 3-of-3 Multisig) is now fully working
- E2E script `e2e-p8-p2sh-p2wsh-3of3.sh` completed and verified
- Fixed receiver address generation to use P2SH-SegWit format (Closes #374)
- P2SH-wrapped SegWit multisig with legacy compatibility (`2...` addresses in regtest)
- Implements BIP49 key derivation for 3-of-3 multisig
- All 3 signatures required (Keygen + Sign1 + Sign2)

### Version 1.5 (2026-01-16)

- ✅ Pattern 6 (P2WSH Native SegWit 2-of-3 Multisig) is now fully working
- E2E script `e2e-p6-p2wsh-2of3.sh` completed and verified
- Native SegWit multisig with Bech32 encoding (`bcrt1q...` 62-char addresses)
- Added native SegWit descriptor support (`wsh`, `wpkh`) in PSBT infrastructure
- Most efficient multisig format - no P2SH wrapper overhead
- Implements BIP84 key derivation and BIP67 sorted multisig keys

### Version 1.4 (2026-01-15)

- ✅ Pattern 5 (P2WPKH Native SegWit Single-sig) is now fully working
- E2E script `e2e-p5-p2wpkh-singlesig.sh` completed and verified
- Native SegWit with Bech32 encoding (`bcrt1q...` addresses)

### Version 1.3 (2026-01-15)

- ✅ Pattern 3 (P2SH-P2WPKH Single-sig) is now fully working
- E2E script `e2e-p3-p2sh-p2wpkh-singlesig.sh` completed and verified

### Version 1.2 (2026-01-15)

- ✅ Pattern 2 (P2PKH 2-of-3 Multisig) is now fully working
- Fixed key derivation path mismatch for multisig accounts (PR #357)
- Added detailed explanation of the fix and root cause
- Updated implementation status for 2-of-3 multisig

### Version 1.1 (2026-01-14)

- Initial comprehensive documentation of all patterns

# Bitcoin (BTC) Technical Reference

This document provides a comprehensive technical reference for Bitcoin implementation in the go-crypto-wallet system. It covers specifications, protocol details, and links to official documentation to help AI agents and developers understand Bitcoin's architecture and implement features correctly.
