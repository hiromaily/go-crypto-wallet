## Core Specifications

### Cryptographic Primitives

#### Elliptic Curve (secp256k1)

Ethereum uses the same secp256k1 elliptic curve as Bitcoin for signing:

```
Private Key:  32 bytes (256-bit scalar)
Public Key:   64 bytes (uncompressed, without 04 prefix) or 65 bytes with prefix
Address:      20 bytes = Keccak256(pubkey)[12:]
```

#### Hash Functions

| Function | Usage |
|----------|-------|
| **Keccak256** | Address derivation, transaction hashing, function selectors |
| **RLP** | Transaction serialization encoding |
| **SHA3** | Generic hashing (note: Ethereum uses Keccak256, not NIST SHA3) |

#### Address Derivation

```
1. Generate secp256k1 key pair
2. Take uncompressed public key (64 bytes, no prefix)
3. Compute Keccak256 hash (32 bytes)
4. Take last 20 bytes
5. Apply EIP-55 checksum encoding
```

### Data Encoding

| Format | Usage |
|--------|-------|
| **RLP** | Transaction serialization |
| **Hex** | Addresses, transaction data (`0x` prefix) |
| **EIP-55** | Mixed-case checksum address encoding |
| **ABI** | Smart contract call data encoding |

---
