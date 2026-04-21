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
