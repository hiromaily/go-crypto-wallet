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
