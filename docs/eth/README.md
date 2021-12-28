# Ethereum

## Free ebooks
- [Ethereum Development with Go](https://goethereumbook.org/en/)

## Development
- [JSON RPC](https://eth.wiki/json-rpc/API)
- [Ethereum Network](https://ethereum.org/en/developers/docs/networks/)
    - [Goerli Testnet Faucet](https://goerli-faucet.slock.it/)
    - [Etherscan for goerli](https://goerli.etherscan.io/)
    - [Status for Goerli](https://stats.goerli.net/)
    - [github goerli/testnet](https://github.com/goerli/testnet)

### go-ethereum
- [go-ethereum](https://github.com/ethereum/go-ethereum)
- [Getting Started with Geth](https://geth.ethereum.org/docs/getting-started)
- [Private Network Tutorial](https://geth.ethereum.org/docs/getting-started/private-net)
- [Dev mode](https://geth.ethereum.org/docs/getting-started/dev-mode)
  - Geth has a development mode which sets up a single node Ethereum test network with a number of options optimized for developing on local machines. You enable it with the --dev argument.

### OpenEthereum
- [github](https://github.com/openethereum/openethereum)
- [docs](https://openethereum.github.io/)

### Ganache
- [Ganache](https://www.trufflesuite.com/ganache)


## Install ethereum
### Install ethereum on MacOS
```
$ brew install ethereum
```

## Run geth on testnet
- with console
```
$ geth --goerli --rpc console
```
- allow any rpcapi
```
$ geth --goerli --rpc --rpcapi admin,debug,web3,eth,txpool,net,personal
# available="[admin debug web3 eth txpool personal clique miner net]
```
- acccess from any hosts
```
$ geth --goerli --rpc --rpcaddr 0.0.0.0 --rpcapi admin,debug,web3,eth,txpool,net,personal
```
- set keystore directory
```
#  --keystore ${HOME}/work/go/src/github.com/hiromaily/go-crypto-wallet/data/keystore
```
- unlock account
```
geth --goerli --rpc --rpcaddr 0.0.0.0 --rpcapi admin,debug,web3,eth,txpool,net,personal --unlock 0xF512F9E94c7B97916ec69cd80F3750F4410EaA63 --password pw --allow-insecure-unlock
```

### Using IPC
```
$ geth --goerli --ipcapi admin,debug,web3,eth,txpool,net,personal
```

## Rest API using [HTTPie](https://httpie.org/)
```
$ http http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
 or
$ http --auth USERNAME:PASSWORD http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
```

- eth_syncing
```
http http://127.0.0.1:8545 method=eth_syncing params:='[]' id=1
```
- web3_clientVersion
```
http http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
```


## Install OpenEthereum
```
# Install OpenEthereum on MacOS
$ brew tap openethereum/openethereum
$ brew install openethereum

# Run OpenEthereum on testnet
$ openethereum --chain goerli --jsonrpc-apis personal
```

## Go Contract Bindings
- [Go Contract Bindings](https://geth.ethereum.org/docs/dapp/native-bindings)
### install Go binding generator
```
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```
### generate go code
- create abi json file from built contract json file. only abi element is required. then
```
abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./pkg/contract/token-abi.go
```

## How to implement multisig on Ethereum
- [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
- [corvus-ch/shamir](https://github.com/corvus-ch/shamir)
