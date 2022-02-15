# Installation

This installation expects MacOS environment.

## Requirements
- Golang 1.16+
- [golangci-lint](https://github.com/golangci/golangci-lint) 1.43+ (for development)
- [direnv](https://direnv.net/)
- [Docker](https://www.docker.com/get-started)

## Common Setup
1. install Golang, Docker
2. build `watch`, `keygen`, `auth` wallets
- only each sign wallet includes corresponding account name as `authName` into binary
```
make build
 or
go build -v -o ${GOPATH}/bin/watch ./cmd/watch/
go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
go build -ldflags "-X main.authName=auth3" -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
go build -ldflags "-X main.authName=auth4" -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
go build -ldflags "-X main.authName=auth5" -v -o ${GOPATH}/bin/sign5 ./cmd/sign/
```
3. configure config files in [./data/config/*.toml](https://github.com/hiromaily/go-crypto-wallet/tree/master/data/config)
4. set environment variables
   - install [direnv](https://direnv.net/)
   - see `.envrc`
   - modify `.envrc` if needed
   - execute `direnv allow` on terminal
5. run Database containers
```
docker compose up watch-db keygen-db sign-db
```

## Bitcoind Setup
At least, one bitcoin core server and 3 different databases are required.

1. copy `bitcoin.conf` from ./data/config/bitcoind/ to ./docker/nodes/btc/data1, data2, data3 directory respectively.
  - I recommend signet network.
2. run bitcoind node containers
```
docker compose up btc-watch btc-keygen btc-sign
```
3. setup `bitcoin-cli` using docker
    - after running `btc-watch` container, set alias on shell
   ```zsh
   alias bitcoin-cli-watch='docker exec -it btc-watch bitcoin-cli'
   alias bitcoin-cli-keygen='docker exec -it btc-keygen bitcoin-cli'
   alias bitcoin-cli-sign='docker exec -it btc-sign bitcoin-cli'
   ```
4. create wallets on bitcoind respectively 
   ```
   ./scripts/operation/create-bitcoind-wallet.sh
     or
   bitcoin-cli-watch createwallet watch
   bitcoin-cli-keygen createwallet keygen
   bitcoin-cli-sign createwallet sign1
   bitcoin-cli-sign createwallet sign2
   bitcoin-cli-sign createwallet sign3
   bitcoin-cli-sign createwallet sign4
   bitcoin-cli-sign createwallet sign5
   ```
5. load wallet (required if btc containers restarted)
   ```
   ./scripts/operation/load-bitcoind-wallet.sh
     or
   bitcoin-cli-watch loadwallet watch
   bitcoin-cli-keygen loadwallet keygen
   bitcoin-cli-sign loadwallet sign1
   bitcoin-cli-sign loadwallet sign2
   bitcoin-cli-sign loadwallet sign3
   bitcoin-cli-sign loadwallet sign4
   bitcoin-cli-sign loadwallet sign5
   ```
6. operation
  - see [Operation Example](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/btc/OperationExample.md)

## Bitcoind Setup without container 
1. install `bitcoind` on macOS directly if needed
  - see [bitcoin core installation](https://github.com/bitcoin/bitcoin/blob/master/doc/build-osx.md)
2. run bitcoind `$ bitcoind`
3. create wallets separately (if only one node used)
    ```
    $ bitcoin-cli createwallet watch
    $ bitcoin-cli createwallet keygen
    $ bitcoin-cli createwallet sign1
    $ bitcoin-cli createwallet sign2
    $ bitcoin-cli createwallet sign3
    $ bitcoin-cli createwallet sign4
    $ bitcoin-cli createwallet sign5
    $ bitcoin-cli listwallets
    [
      "",
      "watch",
      "keygen",
      "sign1",
      "sign2",
      "sign3",
      "sign4",
      "sign5"
    ]
    ```

## Ethereum Setup
It depends on which node you choose

### A. go-ethereum
- run node by docker compose
```
make up-docker-eth
 or
docker compose -f docker-compose.eth.yml up eth-node
```

### B. Ganache
- run node by docker compose
```
docker compose -f docker-compose.eth.yml up ganache
```
- prepare sql file if you choose Ganache.  
  But, first account(index[0]) must not be used. See more instruction [here](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/eth/Ganache.md)

## ERC20 Token Setup
- deploy ERC-20 token contract if needed
- Original ERC-20 token is [here](https://github.com/hiromaily/go-crypto-wallet/tree/master/web/erc20-token). See [`scripts/operation/deploy-token.sh`](https://github.com/hiromaily/go-crypto-wallet/blob/master/scripts/operation/deploy-token.sh)
```
cd ./web/erc20-token
yarn install

# deploy contract to current network
yarn run deploy       # using 7545 port
 or
yarn run deploy-dev2  # using 8545 port
```

- copy `contract address` in console and modify `contract_address` at `ethereum.erc20s.hyt` section in ./config/eth_watch.toml
- copy `account` in console and modify `master_address` at `ethereum.erc20s.hyt` section in ./config/eth_watch.toml

```
# check balance
yarn ts-node src/web3.ts --mode balance --address 0xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
# transfer to specific address
yarn ts-node src/web3.ts --mode transfer --address 0xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --amount 100
```
   
## Ripple Setup
