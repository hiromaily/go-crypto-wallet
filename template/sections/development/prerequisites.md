## Prerequisites

### System Requirements

- **Watch Wallet** (online): Bitcoin Core node v22.0+ with Taproot support
- **Keygen Wallet** (offline): Isolated system for key generation and first signature
- **Sign Wallet** (offline): Isolated system for additional signatures
- **Go version**: 1.21 or higher
- **Bitcoin Core**: v22.0 or higher (for Taproot/Schnorr support)

### Required Features

MuSig2 builds on top of existing wallet features:

- ✅ **Phase 1**: Taproot support (BIP340 Schnorr signatures)
- ✅ **Phase 2**: PSBT support (BIP174)
- ✅ **Phase 3**: MuSig2 implementation (current)

### Wallet Configuration

Ensure your wallets are properly configured:

```bash
# Watch wallet configuration
config/wallet/btc/watch.yaml

# Keygen wallet configuration
config/wallet/btc/keygen.yaml

# Sign wallet configuration (for multisig)
config/wallet/btc/sign1.yaml
```

**Important**: All wallet commands require the `--config` flag to specify the configuration file:

```bash
# Example usage
./watch --config config/wallet/btc/watch.yaml --coin btc create payment
./keygen --config config/wallet/btc/keygen.yaml --coin btc sign --file tx.psbt
./sign1 --config config/wallet/btc/sign1.yaml --coin btc sign --file tx.psbt
```

See `docs/chains/btc/operation_example.md` for configuration details.

---
