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
