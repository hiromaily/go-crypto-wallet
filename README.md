# go-bitcoin
Wallet functionalities handling BTC, BCH and so on. Currencies would be added step by step.

### Note
Drastic refactoring is ongoing

- handling bitcoin core version 0.19
- any comments should be English, not Japanese
- apply for domain driven design

## Requirements
Bitcoin Core 0.19 minimum

## Wallet Type
This is explained for BTC for now.
There are mainly 3 wallets separately.

### Watch only wallet
- This wallet could access to BTC Network
- Only Bitcoin address is stored. Private key is NOT stored here. That's why this is called watch only wallet.
- It works as detection coin received, creation of unsigned transaction and client to call Bitcoin APIs.

### Keygen wallet as cold wallet
- This wallet is key management functionalities. It generates seed and private keys as HD wallet and exports address for watch only wallet.
- Sign unsigned transaction by certain keys.
- Outside network is not used at all.

### Signature wallet as cold wallet
- This wallet is key management for authorization by multi-signature address. It also generates seed and private keys for authorization accounts.
- Sign unsigned transaction by certain keys.
- Outside network is not used at all.


## Install
- This project is ongoing. Until project done to some extent, I don't use package management tool like dep. So you can get packages as below command.
```
go get -u -d -v ./...'
```

## Project layout patterns
- The `pkg` layout pattern, refer to the [linked](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2) URLs for details.