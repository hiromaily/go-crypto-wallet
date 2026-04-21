# PSBT (Partially Signed Bitcoin Transactions)

This directory contains documentation for PSBT implementation and usage in the go-crypto-wallet system.

## Contents

| Document | Description | Audience |
|----------|-------------|----------|
| [user-guide.md](../../../../../docs/chains/btc/psbt/user-guide.md) | How to use PSBT for offline signing workflows | Operators |
| [developer-guide.md](../../../../../docs/chains/btc/psbt/developer-guide.md) | Technical implementation details for developers | Developers |
| [implementation.md](../../../../../docs/chains/btc/psbt/implementation.md) | PSBT implementation architecture | Developers |
| [migration.md](../../../../../docs/chains/btc/psbt/migration.md) | Migration guide from legacy to PSBT | All |
| [poc-example.md](../../../../../docs/chains/btc/psbt/poc-example.md) | Proof of concept examples | Developers |

## What is PSBT?

PSBT (BIP174) is the standard format for Bitcoin transactions that require offline or multi-party signing. It enables:

- **Offline Signing**: Sign transactions on air-gapped devices
- **Multi-party Signing**: Coordinate signatures from multiple parties
- **Hardware Wallet Support**: Compatible with hardware security modules

## PSBT Workflow

```
1. Creator (Watch Wallet) → Create unsigned PSBT
2. Updater → Add metadata (optional)
3. Signer(s) (Keygen/Sign Wallet) → Add signatures
4. Combiner → Combine PSBTs (optional)
5. Finalizer → Create final transaction
6. Extractor → Extract broadcastable transaction
```

## Related Documentation

- [../operations/wallet-flow.md](../../../../../docs/chains/btc/operations/wallet-flow.md) - Transaction flow using PSBT
- [../musig2/](../../../../../docs/chains/btc/musig2/README.md) - MuSig2 uses PSBT for coordination
- [../README.md](../../../../../docs/chains/btc/README.md) - Main BTC documentation index

## References

- [BIP174 - PSBT](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP370 - PSBT Version 2](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki)
- [BIP371 - Taproot PSBT Fields](https://github.com/bitcoin/bips/blob/master/bip-0371.mediawiki)
