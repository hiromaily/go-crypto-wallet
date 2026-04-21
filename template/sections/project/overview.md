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
