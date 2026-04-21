# Wallet Operations

This directory contains guides for operating the Bitcoin wallet system.

## Contents

| Document | Description |
|----------|-------------|
| [wallet-flow.md](../../../../docs/chains/btc/operations/wallet-flow.md) | Setup procedures and transaction flow for Watch/Keygen/Sign wallets |
| [e2e-transaction-patterns.md](../../../../docs/chains/btc/operations/e2e-transaction-patterns.md) | E2E test patterns for various key types and signature schemes |
| [wallet-flow-improvements-2025.md](../../../../docs/chains/btc/operations/wallet-flow-improvements-2025.md) | 2025 workflow enhancement plans |

## Audience

- Wallet operators
- System administrators
- Developers implementing new features

## Key Concepts

1. **Three-Wallet Architecture**
   - **Watch Wallet** (Online): Creates transactions, broadcasts, monitors
   - **Keygen Wallet** (Offline): Generates keys, provides first signature
   - **Sign Wallet** (Offline): Provides additional signatures for multisig

2. **Transaction Types**
   - Deposit: Aggregate client funds to cold storage
   - Payment: Process withdrawal requests
   - Transfer: Internal account movements

## Related Documentation

- [../psbt/](../../../../docs/chains/btc/psbt/README) - PSBT guides for offline signing
- [../musig2/](../../../../docs/chains/btc/musig2/README) - MuSig2 multisig operations
- [../README.md](../../../../docs/chains/btc/README.md) - Main BTC documentation index
