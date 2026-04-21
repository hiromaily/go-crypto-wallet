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
