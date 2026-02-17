# Output Descriptors

This directory contains documentation for Bitcoin Output Descriptor implementation.

## Contents

| Document | Description | Audience |
|----------|-------------|----------|
| [user-guide.md](user-guide.md) | How to use descriptors for address generation | Operators |
| [architecture.md](architecture.md) | Descriptor architecture following Clean Architecture | Developers |
| [development.md](development.md) | Development guide for descriptor features | Developers |
| [api.md](api.md) | API reference for descriptor operations | Developers |
| [examples.md](examples.md) | Descriptor format examples | All |
| [compatibility.md](compatibility.md) | Bitcoin Core compatibility notes | Developers |
| [migration.md](migration.md) | Migration guide from legacy address handling | All |

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

- [../keygen/](../keygen/) - Key generation that uses descriptors
- [../taproot/](../taproot/) - Taproot descriptors
- [../README.md](../README.md) - Main BTC documentation index

## References

- [Bitcoin Core Descriptors Documentation](https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md)
- [BIP380 - Output Script Descriptors General Operation](https://github.com/bitcoin/bips/blob/master/bip-0380.mediawiki)
