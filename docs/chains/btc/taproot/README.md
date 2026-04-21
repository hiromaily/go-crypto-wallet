# Taproot

This directory contains documentation for Taproot (BIP341/BIP86) implementation.

## Contents

| Document | Description | Audience |
|----------|-------------|----------|
| [user-guide.md](./user-guide.md) | How to use Taproot addresses in the wallet system | All |
| [testing.md](./testing.md) | Taproot test procedures and verification | Developers |

## What is Taproot?

Taproot is Bitcoin's major upgrade (activated November 2021) introducing:

- **P2TR Addresses**: `bc1p...` (mainnet), `tb1p...` (testnet)
- **Schnorr Signatures**: BIP340 - more efficient than ECDSA
- **BIP86 Derivation**: `m/86'/0'/account'/0/index`

## Benefits

| Feature | Benefit |
|---------|---------|
| **Lower Fees** | 30-50% smaller transaction sizes |
| **Enhanced Privacy** | All spends look identical on-chain |
| **Signature Aggregation** | Enables MuSig2 for efficient multisig |
| **Faster Validation** | Schnorr signatures verify faster |

## When to Use Taproot

✅ New wallets and addresses
✅ When transaction fee savings matter
✅ When privacy is important
✅ For multisig setups (via MuSig2)

## Quick Start

```yaml
# config/wallet/btc/keygen.yaml
address_type: taproot  # Use Taproot (P2TR)
```

## Related Documentation

- [../musig2/](../musig2/README) - MuSig2 leverages Taproot
- [../keygen/improvements-2025.md](../keygen/improvements-2025.md) - Taproot key generation
- [../operations/wallet-flow.md](../operations/wallet-flow.md) - Transaction flows
- [../README.md](../README.md) - Main BTC documentation index

## References

- [BIP340 - Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP341 - Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP342 - Tapscript](https://github.com/bitcoin/bips/blob/master/bip-0342.mediawiki)
- [BIP350 - Bech32m](https://github.com/bitcoin/bips/blob/master/bip-0350.mediawiki)
- [BIP86 - Taproot Derivation](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)
