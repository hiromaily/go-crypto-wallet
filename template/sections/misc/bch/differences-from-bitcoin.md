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
