# go-crypto-wallet

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/xrp-img.jpg?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/ethereum-img.png?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![Test](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml/badge.svg)](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml)
[![GitHub release](https://img.shields.io/badge/release-v6.2.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://hiromaily.github.io/go-crypto-wallet/)

**[📖 Documentation Site](https://hiromaily.github.io/go-crypto-wallet/)**

Wallet functionalities to create raw transaction, to sign on unsigned transaction,
to send signed transaction for BTC, BCH, ETH, XRP and so on.

# Wallet Type

This is explained for BTC/BCH for now.
There are mainly 3 wallets separately and these wallets are expected to be installed in each different devices.

## 1.Watch only wallet

- Only this wallet run online to access to BTC/BCH Nodes.
- Only pubkey address is stored. Private key is NOT stored for security reason. That's why this is called `watch only wallet`.
- Major functionalities are
  - creating unsigned transaction
  - sending signed transaction
  - monitoring transaction status.

## 2.Keygen wallet as cold wallet

- Key management functionalities for accounts.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts
  - generating keys based on `HD Wallet`
  - generating multisig addressed according to account setting
  - exporting pubkey addresses as csv file which is imported from `Watch only wallet`
  - signing on unsigned transaction as first sign. However, multisig addresses could not be completed by only this wallet.

## 3.Sign wallet as cold wallet (Auth wallet)

- The internal authorization operators would use this wallet to sign on unsigned transaction for multisig addresses.
- Each of operators would be given own authorization account and Sing wallet apps.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts for own auth account
  - generating keys based on `HD Wallet` for own auth account
  - exporting full-pubkey addresses as csv file which is imported from `Keygen wallet` to generate multisig address
  - signing on unsigned transaction as second or more signs for multisig addresses.

# Workflow diagram

## BTC

### 1. Generate keys

![generate keys](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/0_key%20generation%20diagram.png?raw=true)

### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address

![create tx](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address

![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)

# Wallet Architecture

## Three Wallet Types

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Watch Wallet  │     │  Keygen Wallet  │     │   Sign Wallet   │
│    (Online)     │     │   (Offline)     │     │   (Offline)     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ • Monitor txs   │     │ • Generate keys │     │ • Auth signing  │
│ • Create unsig  │     │ • Create multis │     │ • Second+ sign  │
│ • Send signed   │     │ • First sign    │     │ • Export pubkey │
│ • Import pubkey │     │ • Export pubkey │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │    CSV/File Export    │    CSV/File Export    │
        └───────────────────────┴───────────────────────┘
```

## Security Model

1. **Keygen Wallet** (Offline): Generates HD wallet seeds and keys. Never connects to network.
2. **Sign Wallet** (Offline): Provides authorization signatures. Each operator has own instance.
3. **Watch Wallet** (Online): Only stores public keys. Cannot sign transactions.
