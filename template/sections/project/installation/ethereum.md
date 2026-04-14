### Ethereum Setup

It depends on which node you choose

#### A. go-ethereum

- run node by docker compose

```
make up-docker-geth
 or
docker compose -f compose.eth.yaml up geth
```

- If you have exported data, run `make import-geth-data` after tweaking parameters before running `make up-docker-geth`.

##### [WIP] Call API => move to operation example

1. `watch -coin eth api clientversion`

```
client version: Geth/v1.10.15-stable-8be800ff/linux-amd64/go1.17.5
```

1. `watch -coin eth api nodeinfo`
2. `watch -coin eth api syncing`
3. `watch -coin eth api netversion`

#### B. Ganache

- run node by docker compose

```
docker compose -f compose.eth.yaml up ganache
```

- prepare sql file if you choose Ganache.
  But, first account(index[0]) must not be used. See more instruction [here](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/eth/Ganache.md)
