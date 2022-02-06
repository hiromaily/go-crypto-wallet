# go-crypto-wallet

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/xrp-img.jpg?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/ethereum-img.png?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![codebeat badge](https://codebeat.co/badges/792a7c07-2352-4b7e-8083-0a323368b26f)](https://codebeat.co/projects/github-com-hiromaily-go-crypto-wallet-master)
[![GitHub release](https://img.shields.io/badge/release-v5.0.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/LICENSE)

Wallet functionalities to create raw transaction, to sing on unsigned transaction, to send signed transaction for BTC, BCH, ETH, XRP and so on.  

## What kind of coin can be used?
- Bitcoin
- Bitcoin Cash
- Ethereum
- ERC-20 Token
- Ripple


## Current development
- This project is under refactoring
  - based on `Clean Code`, `Clean Architecture`, [`Refactoring`](https://martinfowler.com/articles/refactoring-2nd-ed.html)
- Bitcoin Core version 22.0 is released. Must be cought up.

## Expected use cases
### 1.Deposit functionality
- Pubkey addresses are given to our users first.
- Users would want to deposit coins on our system.
- After users sent coins to their given addresses, these all amount of coins are sent to our safe addresses managed offline by cold wallet

### 2.Payment functionality
- Users would want to withdraw their coins to specific addresses.
- Transaction is created and sent after payment is requested by users.

### 3.Transfer functionality
- Internal use. Each accounts can transfer coins among internal accounts.


## Wallet Type
This is explained for BTC/BCH for now.  
There are mainly 3 wallets separately and these wallets are expected to be installed in each different devices.

### 1.Watch only wallet
- Only this wallet run online to access to BTC/BCH Nodes.
- Only pubkey address is stored. Private key is NOT stored for security reason. That's why this is called `watch only wallet`.
- Major functionalities are
    - creating unsigned transaction
    - sending signed transaction
    - monitoring transaction status.

### 2.Keygen wallet as cold wallet
- Key management functionalities for accounts.  
- This wallet is expected to work offline.
- Major functionalities are
    - generating seed for accounts
    - generating keys based on `HD Wallet`
    - generating multisig addressed according to account setting
    - exporting pubkey addresses as csv file which is imported from `Watch only wallet`
    - signing on unsigned transaction as first sign. However, multisig addresses could not be completed by only this wallet.

### 3.Sign wallet as cold wallet (Auth wallet)
- The internal authorization operators would use this wallet to sign on unsigned transaction for multisig addresses.
- Each of operators would be given own authorization account and Sing wallet apps.
- This wallet is expected to work offline.
- Major functionalities are
    - generating seed for accounts for own auth account
    - generating keys based on `HD Wallet` for own auth account
    - exporting full-pubkey addresses as csv file which is imported from `Keygen wallet` to generate multisig address
    - signing on unsigned transaction as second or more signs for multisig addresses.


## Workflow diagram
### BTC
#### 1. Generate keys
![generate keys](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/0_key%20generation%20diagram.png?raw=true)

#### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address.
![create tx](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

#### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address.
![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)


## Requirements
- Golang 1.16+
- [golangci-lint](https://github.com/golangci/golangci-lint) 1.42+ (for development)
- [direnv](https://direnv.net/)
- [Docker](https://www.docker.com/get-started)
    - MySQL 5.7
    - Node Server
        - BTC: [Bitcoin Core 0.18+ for Bitcoin node](https://bitcoin.org/en/bitcoin-core/)
        - BCH: [Bitcoin ABC 0.21+ for Bitcoin cash node](https://www.bitcoinabc.org/)
        - ETH: [go-ethereum](https://github.com/ethereum/go-ethereum), [Ganache](https://www.trufflesuite.com/ganache), [erc20-token](https://github.com/hiromaily/go-crypto-wallet/tree/master/web/erc20-token)
        - XRP: [rippled](https://xrpl.org/manage-the-rippled-server.html), [ripple-lib-server](https://github.com/hiromaily/go-crypto-wallet/tree/master/web/ripple-lib-server)

## Components inside repository
- ripple-lib-server
  - see web/ripple-lib-server
- erc20-token
  - see web/erc20-token

## Installation
see [Installation](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/Installation.md)

### Command example
- [see Makefile](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile)
    - [Makefile](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile)
    - [Makefile for watch wallet operation](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile_watch_op.mk)
    - [Makefile for keygen wallet operation](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile_keygen_op.mk)
    - [Makefile for sign wallet operation](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile_sign_op.mk)
- [see operation scripts](https://github.com/hiromaily/go-crypto-wallet/tree/master/scripts/operation)

### Setup keys for BTC
- [see generate-btc-key.sh](https://github.com/hiromaily/go-crypto-wallet/blob/master/scripts/operation/generate-btc-key.sh)

Keygen wallet
```
# create seed
keygen create seed

# create hdkey for client, deposit, payment account
keygen create hdkey -account client -keynum 10
keygen create hdkey -account deposit -keynum 10
keygen create hdkey -account payment -keynum 10
keygen create hdkey -account stored -keynum 10

# import generated private key into keygen wallet
keygen import privkey -account client
keygen import privkey -account deposit
keygen import privkey -account payment
keygen import privkey -account stored
```

Sign wallet
```
# create seed
sign create seed

# create hdkey for authorization
sign -wallet sign1 create hdkey
sign2 -wallet sign2 create hdkey
sign3 -wallet sign3 create hdkey
sign4 -wallet sign4 create hdkey
sign5 -wallet sign5 create hdkey

# import generated private key into sign wallet
sign -wallet sign1 import privkey
sign2 -wallet sign2 import privkey
sign3 -wallet sign3 import privkey
sign4 -wallet sign4 import privkey
sign5 -wallet sign5 import privkey

# export full-pubkey as csv file
sign -wallet sign1 export fullpubkey
sign2 -wallet sign2 export fullpubkey
sign3 -wallet sign3 export fullpubkey
sign4 -wallet sign4 export fullpubkey
sign5 -wallet sign5 export fullpubkey
```

Keygen wallet
```
# import full-pubkey
keygen import fullpubkey -file auth1-fullpubkey-file
keygen import fullpubkey -file auth2-fullpubkey-file
keygen import fullpubkey -file auth3-fullpubkey-file
keygen import fullpubkey -file auth4-fullpubkey-file
keygen import fullpubkey -file auth5-fullpubkey-file

# create multisig address
keygen create multisig -account deposit
keygen create multisig -account payment
keygen create multisig -account stored

# export address
keygen export address -account client
keygen export address -account deposit
keygen export address -account payment
keygen export address -account stored
```

Watch wallet
```
# import addresses generated by keygen wallet
watch import address -account client -file client-address-file
watch import address -account deposit -file deposit-address-file
watch import address -account payment -file payment-address-file
watch import address -account stored -file stored-address-file
```

### Operation for deposit action
```
# check client addresses if it receives coin
watch create deposit

# sign on keygen wallet
keygen sign -file xxx.file

# send signed tx
watch send -file xxx.csv

```

### Operation for payment action
```
# check payment_request if there are requests
wallet create payment

# sign on keygen wallet for first sigunature
keygen sign -file xxx.file

# sign on sign wallet for second sigunature
sign sign -file xxx.file

# send signed tx
watch send -file xxx.csv

```

## TODO
### Basics
- [ ] Add ATOM tokens on [Cosmos Hub](https://hub.cosmos.network/main/hub-overview/overview.html)
- [ ] Add [Polkadot](https://polkadot.network/technology/)
- [ ] Various monitoring patterns to detect suspicious operations.
- [ ] Add Github Action as CI

### For BTC/BCH
- [ ] Multisig-address is used only once because of security reason, so after tx is sent, related receiver addresses should be updated by is_allocated=true.
- [ ] Sent tx is not proceeded in bitcoin network if fee is not enough comparatively. So re-sending tx functionality is required adding more fee.

### For ERC20 token
- [ ] Add any useful APIs using contract equivalent to ETH APIs
- [ ] Monitoring for ERC20 token

### For ETH
- [ ] Make sure that `quantity-tag` is used properly. e.g. when getting balance, which quantity-tag should be used, latest or pending.
- [ ] Handling secret of private key properly. Password could be passed from command line argument.

### For XRP
- [ ] Handling secret of private key properly. Password could be passed from command line argument.

## Project layout patterns
- The `pkg` layout pattern, refer to the [linked](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2) URLs for details.