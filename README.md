# go-bitcoin
Wallet functionalities handling BTC, BCH, ETH and so on. Currencies would be added step by step.

## Structures
This is explained for BTC for now.
There are mainly 3 wallets separately.

### Watch only wallet
- This wallet could access to BTC Network
- Only Bitcoin address is stored. Private key is NOT stored here. That's why this is called watch only wallet.
- It works as detection coin received, creation of unsigned transaction and client to call Bitcoin APIs.

### Cold wallet1
- This wallet is key management functionalities. It generates seed and private keys as HD wallet and exports address for watch only wallet.
- Sign unsigned transaction by certain keys.
- Outside network is not used at all.

### Cold wallet2
- This wallet is key management for authorization by multi-signature address. It also generates seed and private keys for authorization accounts.
- Sign unsigned transaction by certain keys.
- Outside network is not used at all.


## Install
- This project is ongoing. Until project done to some extent, I don't use package management tool like dep. So you can get packages as below command.
```
go get -u -d -v ./...'
```