<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/btc/keygen/README.tpl.md · Run `make docs` to regenerate.
-->

# Key Generation

This directory contains documentation for Bitcoin key generation design and improvements.

## Contents

| Document | Description |
|----------|-------------|
| [improvements-2025.md](./improvements-2025.md) | 2025 key generation modernization improvements |
| [improvements-2025-ja.md](./improvements-2025-ja.md) | 2025 improvements (Japanese version) |
| [interface-design.md](./interface-design.md) | Key generator interface design documentation |

## Audience

- Developers working on key generation features
- Security reviewers

## Key Topics

- Taproot (BIP341/BIP86) support
- BIP49 (P2WPKH-P2SH) implementation
- BIP85 (Deterministic Entropy) consideration
- Descriptor wallets support
- MuSig2 improvements
- Random number generation enhancement
- Security enhancements

## Security Notice

⚠️ Key generation is a security-critical operation. Changes to this area require careful review and testing.

## Related Documentation

- [../taproot/](../taproot/README.md) - Taproot implementation guides
- [../musig2/](../musig2/README.md) - MuSig2 architecture
- [../descriptor/](../descriptor/README.md) - Output descriptor guides
- [../README.md](../README.md) - Main BTC documentation index
