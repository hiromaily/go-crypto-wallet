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
