# MuSig2

This directory contains documentation for MuSig2 (BIP327) multisignature implementation.

## Contents

| Document | Description | Audience |
|----------|-------------|----------|
| [user-guide.md](../../../../docs/chains/btc/musig2/user-guide.md) | How to use MuSig2 for multisig operations | Operators |
| [architecture.md](../../../../docs/chains/btc/musig2/architecture.md) | MuSig2 system architecture and data flow | Developers |
| [security.md](../../../../docs/chains/btc/musig2/security.md) | Security considerations and nonce management | All |
| [migration-from-traditional.md](../../../../docs/chains/btc/musig2/migration-from-traditional.md) | Migration guide from traditional P2WSH multisig | All |

## What is MuSig2?

MuSig2 is a two-round Schnorr multisignature protocol (BIP327) that enables N-of-N multisig appearing as single-sig on-chain:

- **Smaller Transactions**: 30-50% size reduction vs traditional P2WSH multisig
- **Lower Fees**: Proportional to transaction size reduction
- **Better Privacy**: Multisig transactions look like single-sig on-chain
- **Schnorr Signatures**: Uses BIP340 via Taproot (P2TR)

## Two-Round Protocol

```
Round 1: Nonce Generation (Parallel)
├── Each signer generates random nonce
├── Public nonces are exchanged
└── Aggregate nonce computed

Round 2: Partial Signing (Sequential)
├── Each signer creates partial signature
├── Partial signatures are collected
└── Aggregated into single 64-byte Schnorr signature
```

## ⚠️ Critical Security: Nonce Management

**NEVER reuse nonces** - Reusing nonces in MuSig2 will leak the private key!

- Generate fresh nonces for every transaction
- Delete nonces immediately after signing
- See [security.md](../../../../docs/chains/btc/musig2/security.md) for details

## CLI Commands

```bash
# Create MuSig2 address
keygen --coin btc create musig2-address

# Generate nonces
keygen --coin btc musig2 nonce
sign --coin btc musig2 nonce

# Create partial signatures
keygen --coin btc musig2 sign
sign --coin btc musig2 sign

# Aggregate signatures
watch --coin btc musig2 aggregate
```

## Related Documentation

- [../taproot/](../../../../docs/chains/btc/taproot/README.md) - MuSig2 requires Taproot
- [../psbt/](../../../../docs/chains/btc/psbt/README.md) - PSBT used for coordination
- [../operations/wallet-flow.md](../../../../docs/chains/btc/operations/wallet-flow.md) - Transaction flows
- [../README.md](../../../../docs/chains/btc/README.md) - Main BTC documentation index

## References

- [BIP327 - MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [MuSig2 Paper](https://eprint.iacr.org/2020/1261)
