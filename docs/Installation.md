# Installation

This installation expects macOS environment.

## Requirements
- Golang 1.16+
- [golangci-lint](https://github.com/golangci/golangci-lint) 1.42+ (for development)
- [direnv](https://direnv.net/)
- [Docker](https://www.docker.com/get-started)

## Common Setup
1. install Golang, Docker
2. run Database by docker compose
```
make up-docker-db
 or
docker compose up watch-db keygen-db sign-db
```
3. build `watch`, `keygen`, `auth` wallets
```
make bld
 or
go build -v -o ${GOPATH}/bin/watch ./cmd/watch/
go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign ./cmd/sign/
go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
go build -ldflags "-X main.authName=auth3" -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
go build -ldflags "-X main.authName=auth4" -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
go build -ldflags "-X main.authName=auth5" -v -o ${GOPATH}/bin/sign5 ./cmd/sign/
```
4. configure config files in [./data/config/*.toml](https://github.com/hiromaily/go-crypto-wallet/tree/master/data/config) (better after node setup)
5. set environment variables
   - install [direnv](https://direnv.net/)
   - see `.envrc`
   - modify `.envrc` if needed
   - execute `direnv allow` on terminal

## Bitcoin Setup
At least, one bitcoin core server and 3 different databases are required.

1. run bitcoin node by docker-compose
```
make bld-docker-btc
 or
docker compose build btc-watch

make up-docker-btc
 or
docker compose up btc-watch btc-keygen btc-sign
```

2. install `bitcoind` on macOS directly if needed
    - see [bitcoin core installation](https://github.com/bitcoin/bitcoin/blob/master/doc/build-osx.md)
    - run bitcoind
    ```
    $ bitcoind
    ```
    - create wallets separately (if only one node used)
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
It depends on what node you choose

1. run node by docker compose
    - run go-ethereum
    ```
    make up-docker-eth
     or
    docker compose -f docker-compose.eth.yml up eth-node
    ```
    - run Ganache
    ```
    docker compose -f docker-compose.eth.yml up ganache
    ```

2. prepare sql file for Ganache. see more instruction [here](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/eth/Ganache.md)

3. deploy ERC-20 token if needded
Original ERC-20 token is [erc20-token](https://github.com/hiromaily/erc20-token). See [`scripts/operation/deploy-token.sh`](https://github.com/hiromaily/go-crypto-wallet/blob/master/scripts/operation/deploy-token.sh)
    - run the below
    ```
    git clone https://github.com/hiromaily/erc20-token.git
    cd erc20-token

    yarn install
    yarn run deploy-dev2  # using 8545 port
     or
    yarn run deploy       # using 7545 port
    ```
    - copy `contract address` in console and modify `contract_address` at `ethereum.erc20s.hyt` section in ./config/eth_watch.toml
    - copy `account` in console and modify `master_address` at `ethereum.erc20s.hyt` section in ./config/eth_watch.toml

## Ripple Setup
