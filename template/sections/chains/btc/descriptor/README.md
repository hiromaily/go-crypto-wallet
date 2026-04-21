# Output Descriptors

This directory contains documentation for Bitcoin Output Descriptor implementation.

## Contents

| Document | Description | Audience |
|----------|-------------|----------|
| [user-guide.md](../../../../../docs/chains/btc/descriptor/user-guide.md) | How to use descriptors for address generation | Operators |
| [architecture.md](../../../../../docs/chains/btc/descriptor/architecture.md) | Descriptor architecture following Clean Architecture | Developers |
| [development.md](../../../../../docs/chains/btc/descriptor/development.md) | Development guide for descriptor features | Developers |
| [api.md](../../../../../docs/chains/btc/descriptor/api.md) | API reference for descriptor operations | Developers |
| [examples.md](../../../../../docs/chains/btc/descriptor/examples.md) | Descriptor format examples | All |
| [compatibility.md](../../../../../docs/chains/btc/descriptor/compatibility.md) | Bitcoin Core compatibility notes | Developers |
| [migration.md](../../../../../docs/chains/btc/descriptor/migration.md) | Migration guide from legacy address handling | All |

## What are Output Descriptors?

Output Descriptors are a language for describing collections of output scripts. They provide:

- **Precise Address Specification**: Unambiguous description of how to derive addresses
- **Wallet Interoperability**: Standard format supported by Bitcoin Core and other wallets
- **Backup Simplicity**: Single string captures all information needed to recover addresses

## Supported Descriptor Types

| Type | Description | Example |
|------|-------------|---------|
| `pkh()` | P2PKH (Legacy) | `pkh([fingerprint/44'/0'/0']xpub.../0/*)` |
| `sh(wpkh())` | P2SH-P2WPKH (Nested SegWit) | `sh(wpkh([...]))` |
| `wpkh()` | P2WPKH (Native SegWit) | `wpkh([fingerprint/84'/0'/0']xpub.../0/*)` |
| `tr()` | P2TR (Taproot) | `tr([fingerprint/86'/0'/0']xpub.../0/*)` |
| `wsh()` | P2WSH (Multisig) | `wsh(multi(2,...))` |

## Related Documentation

- [../keygen/](../../../../../docs/chains/btc/keygen/README.md) - Key generation that uses descriptors
- [../taproot/](../../../../../docs/chains/btc/taproot/README.md) - Taproot descriptors
- [../README.md](../../../../../docs/chains/btc/README.md) - Main BTC documentation index

## References

- [Bitcoin Core Descriptors Documentation](https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md)
- [BIP380 - Output Script Descriptors General Operation](https://github.com/bitcoin/bips/blob/master/bip-0380.mediawiki)
