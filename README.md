## Overview

### What is MuSig2?

MuSig2 is a cryptographic protocol that enables multiple parties to create a **single aggregated Schnorr signature** that looks identical to a single-signature transaction on the blockchain. Unlike traditional multisig (P2SH, P2WSH) where multiple signatures are stored on-chain, MuSig2 aggregates multiple signatures into one, providing significant benefits.

### Benefits Over Traditional Multisig

| Feature | Traditional P2WSH Multisig | MuSig2 |
|---------|---------------------------|--------|
| **On-Chain Appearance** | Multiple signatures visible | Single signature (looks like single-sig) |
| **Transaction Size** | ~370-400 bytes (2-of-3) | ~200-250 bytes (30-50% smaller) |
| **Privacy** | Multisig is visible | Indistinguishable from single-sig |
| **Fees** | Higher (proportional to size) | 30-50% lower |
| **Signature Algorithm** | ECDSA | Schnorr (BIP340) |
| **Address Type** | P2WSH (bc1q...) | P2TR Taproot (bc1p...) |
| **Compatibility** | Older standard | Modern (Bitcoin Core 22.0+) |

### When to Use MuSig2

- ✅ **New multisig setups** - Best privacy and efficiency
- ✅ **High-volume operations** - Significant fee savings over time
- ✅ **Privacy-focused applications** - Transactions look like single-sig
- ✅ **Modern infrastructure** - Requires Bitcoin Core 22.0+ and Taproot support
- ⚠️ **Legacy multisig** - Traditional P2WSH still supported for backward compatibility

### How MuSig2 Works (Two-Round Protocol)

```
Round 1: Nonce Generation (Parallel)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Generate Nonce 1                   │
│  Sign Wallet 1 → Generate Nonce 2  (can run in     │
│  Sign Wallet 2 → Generate Nonce 3   parallel)      │
└─────────────────────────────────────────────────────┘
                        ↓
            Exchange nonces via PSBT files
                        ↓
Round 2: Signing (Sequential)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Create Partial Signature 1         │
│  Sign Wallet 1 → Create Partial Signature 2         │
│  Sign Wallet 2 → Create Partial Signature 3         │
└─────────────────────────────────────────────────────┘
                        ↓
            Collect partial signatures
                        ↓
Aggregation (Watch Wallet)
┌─────────────────────────────────────────────────────┐
│  Watch Wallet → Aggregate Partial Signatures        │
│              → Verify Final Signature               │
│              → Broadcast Transaction                │
└─────────────────────────────────────────────────────┘
```

**Key Security Feature:**

- Each wallet generates a **nonce** (random value) in Round 1
- Nonces must be **unique per transaction** and **never reused**
- Reusing nonces can leak private keys - this is critical!

---

## Supported Chains

| Chain | Address Types | Highlights | E2E Patterns |
|-------|--------------|------------|--------------|
| **[BTC](./docs/chains/btc/README.md)** | P2PKH → P2TR | Taproot (BIP341), MuSig2 (BIP327), Descriptor Wallets (BIP380), PSBT | [11 patterns](./docs/chains/btc/operations/e2e-transaction-patterns.md) |
| **[BCH](./docs/chains/bch/README.md)** | P2PKH, P2SH | CashAddr format, SIGHASH\_FORKID replay protection | 3 patterns |
| **[ETH](./docs/chains/eth/README.md)** | EOA `0x...` | ERC-20 tokens, Safe multisig (v1.4.1), MPC-TSS threshold signing | 4 patterns |
| **[XRP](./docs/chains/xrp/README.md)** | Classic `r...` | Ed25519 / secp256k1, multisig, offline keygen | 2 patterns |

Each pattern is implemented in Go and verified through real transactions on regtest or a local node — not mocks.

## Quick Start

Requires Docker. See the [Installation Guide](./docs/getting-started/installation.md) for full setup.

```bash
# Bitcoin (BTC) — 11 patterns from P2PKH to Taproot
make btc-e2e P=1    # P2PKH single-sig
make btc-e2e P=9    # P2TR Taproot (Schnorr)

# Bitcoin Cash (BCH)
make bch-e2e P=1    # P2PKH single-sig
make bch-e2e P=2    # P2SH 2-of-3 multisig

# Ethereum (ETH)
make eth-e2e-p1     # Single-sig EIP-1559
make eth-e2e-p3     # Safe 2-of-2 multisig
make eth-e2e-p4     # MPC-TSS 2-of-3 threshold signing

# XRP Ledger
make xrp-e2e-p1     # Single-sig payment
make xrp-e2e-p2     # 2-of-2 multisig payment
```

See the [E2E Transaction Patterns Guide](./docs/chains/btc/operations/e2e-transaction-patterns.md) for all patterns, CI mode (`make btc-e2e-ci P=1`), and reset options.

## Documentation

| Topic | Link |
|-------|------|
| Getting Started | [Installation & Setup](./docs/getting-started/installation.md) |
| Architecture | [System Design](https://hiromaily.github.io/go-crypto-wallet/architecture/) |
| CLI Reference | [Command × Chain × UseCase matrix](./internal/interface-adapters/cli/README.md) |
| BTC Guide | [Bitcoin](./docs/chains/btc/README.md) |
| BCH Guide | [Bitcoin Cash](./docs/chains/bch/README.md) |
| ETH Guide | [Ethereum](./docs/chains/eth/README.md) |
| XRP Guide | [XRP Ledger](./docs/chains/xrp/README.md) |
| E2E Patterns | [Transaction Pattern Guide](./docs/chains/btc/operations/e2e-transaction-patterns.md) |
| Full Docs | [Documentation Site](https://hiromaily.github.io/go-crypto-wallet/) |
