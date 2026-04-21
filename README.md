<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/README.tpl.md · Run `make docs` to regenerate.
-->

<p align="center">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/bitcoin-img.svg?sanitize=true" alt="Bitcoin" width="140px">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/ethereum-img.png?raw=true" alt="Ethereum" width="140px">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/xrp-img.jpg?raw=true" alt="XRP" width="140px">
</p>

# go-crypto-wallet

**A production-grade multi-chain cold wallet system built in Go.**

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![Test](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml/badge.svg)](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml)
[![GitHub release](https://img.shields.io/badge/release-v6.2.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://hiromaily.github.io/go-crypto-wallet/)

**[📖 Documentation Site](https://hiromaily.github.io/go-crypto-wallet/)**

---

Managing cryptocurrency at scale means two things in conflict: keys must stay offline to be secure, and wallets must stay online to be useful. Most implementations compromise one for the other.

go-crypto-wallet resolves this with a **three-wallet architecture**: an offline **Keygen** wallet that generates and holds private keys, an offline **Sign** wallet held by independent authorizers for multisig, and an online **Watch** wallet that creates and broadcasts transactions without ever touching a private key.

```text
┌──────────────────┐      ┌──────────────────┐
│  Keygen Wallet   │      │   Sign Wallet    │
│   (OFFLINE)      │      │   (OFFLINE)      │
│                  │      │                  │
│  HD key gen      │      │  Auth signing    │
│  Multisig addrs  │      │  2nd+ signature  │
└────────┬─────────┘      └────────┬─────────┘
         │  export pubkeys / sign  │
         └────────────┬────────────┘
                      │
             ┌────────▼─────────┐
             │   Watch Wallet   │
             │    (ONLINE)      │
             │                  │
             │  Create tx       │
             │  Broadcast tx    │
             │  Monitor         │
             └──────────────────┘
```

Every transaction pattern — from legacy P2PKH to Taproot MuSig2 and MPC-TSS — is implemented in Go and verified end-to-end against real local nodes.

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
