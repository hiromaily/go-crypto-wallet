# go-bitcoin

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-bitcoin/master/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-bitcoin)](https://goreportcard.com/report/github.com/hiromaily/go-bitcoin)
[![codebeat badge](https://codebeat.co/badges/792a7c07-2352-4b7e-8083-0a323368b26f)](https://codebeat.co/projects/github-com-hiromaily-go-bitcoin-master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/LICENSE)

Wallet functionalities handling BTC, BCH and so on. Currencies would be added step by step.

## Wallet Type
This is explained for BTC for now.
There are mainly 3 wallets separately and these wallets are expected to be installed each diffrent devices.

### Watch only wallet
- This wallet could access to BTC Network
- Only Bitcoin public address is stored. Private key is NOT stored here for security reason. That's why this is called watch only wallet.
- It works for creating unsigned transaction, sending signed transaction and monitoring trasaction status.

### Keygen wallet as cold wallet
- This wallet is key management functionalities. It generates seed and private keys as HD wallet and exports address for watch only wallet.
- Sign unsigned transaction as first signature. Multisig address can not be completed by only this wallet.
- Outside network is not used at all because of cold wallet.

### Signature wallet as cold wallet
- This wallet is signature management for authorization by multi-signature address. It also generates seed and private keys for authorization accounts.
- Sign unsigned transaction as second sigunature. Mustisig address must require this wallet.
- Outside network is not used at all because of cold wallet.

## Workflow diagram
### 1. Generate keys
![generate keys](https://raw.githubusercontent.com/hiromaily/go-bitcoin/master/images/0_key%20generation%20diagram.png?raw=true)

### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address.
![create tx](https://raw.githubusercontent.com/hiromaily/go-bitcoin/master/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address.
![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-bitcoin/master/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)


## Requirements
- Bitcoin Core 0.18+
- MySQL 5.7
- Golang 1.13+
- Docker


## Install on local
At least, one bitcoin core server and 3 different databases are required.  
After [bitcoin core installation](https://github.com/bitcoin/bitcoin/blob/master/doc/build-osx.md) is done
```
# run database on docker
$ docker-compose up btc-watch-db btc-keygen-db btc-signature-db

# run bitcoind
$ bitcoind

# create wallets separately
$ bitcoin-cli createwallet watch
$ bitcoin-cli createwallet keygen
$ bitcoin-cli createwallet sign
$ bitcoin-cli listwallets
[
  "",
  "watch",
  "keygen",
  "sign"
]

# build source
go build -i -v -o ${GOPATH}/bin/watch ./cmd/watch/
go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -i -v -o ${GOPATH}/bin/sign ./cmd/sign/
```

### Configuration
- [watch.toml](https://github.com/hiromaily/go-bitcoin/blob/master/data/config/btc/watch.toml)
- [keygen.toml](https://github.com/hiromaily/go-bitcoin/blob/master/data/config/btc/keygen.toml)
- [sign.toml](https://github.com/hiromaily/go-bitcoin/blob/master/data/config/btc/sign.toml)

## Example
- [see scripts](https://github.com/hiromaily/go-bitcoin/tree/master/scripts/operation)
- [see Makefile](https://github.com/hiromaily/go-bitcoin/blob/master/Makefile)

### Setup for any keys
- [see scripts](https://github.com/hiromaily/go-bitcoin/blob/master/scripts/operation/generate-key-local.sh)

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
- [x] Change ORM to sqlboiler
- [x] Implement proper database transaction.
- [x] Re-design procedure for creating Multisig address.
- [x] Flexible multisig proportion M:N. For now only 2:2 fixed proportion is available.
- [ ] Implement `encryptwallet` to keygen/sign wallets
- [ ] Lock/Unlock wallet for security reason, authorized operator is available.
- [ ] Backup/Restore for wallet.dat. If wallt.dat is broken, it's not easy to recover.
- [ ] Refactoring logic for unsigned transaction creation.
- [ ] Multisig-address is used only once because of security reason, so after tx is sent, related receiver addresses should be updated by is_allocated=true.
- [ ] Tweak key generation logic for multisig address for BCH.
- [ ] Various monitoring patterns to detect suspicious operations.
- [ ] Sent tx is not proceeded in bitcoin network if fee is not enough comparatively. So re-sending tx functionality is required adding more fee.
- [ ] Scaling Bitcoin Core server and synchronize wallet.dat among Bitcoin Core cluster.
- [ ] High coverage of UnitTest.
- [ ] Add CircleCI or other CI service
- [ ] Add new coins like Ethereum.


## Project layout patterns
- The `pkg` layout pattern, refer to the [linked](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2) URLs for details.