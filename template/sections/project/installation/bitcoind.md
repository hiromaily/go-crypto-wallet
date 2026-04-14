### Bitcoind Setup

At least, one bitcoin core server and 1 database (with 3 schemas: watch, keygen, sign) are required.

1. copy `bitcoin.conf` from ./config/blockchain/bitcoind/ to ./docker/nodes/btc/watch, keygen, sign1 directory respectively.

- I recommend signet network.

1. run bitcoind node containers

```
docker compose -f compose.btc.yaml up btc-watch btc-keygen btc-sign
```

1. setup `bitcoin-cli` using docker
    - after running `btc-watch` container, set alias on shell

   ```zsh
   alias bitcoin-cli-watch='docker exec -it btc-watch bitcoin-cli'
   alias bitcoin-cli-keygen='docker exec -it btc-keygen bitcoin-cli'
   alias bitcoin-cli-sign='docker exec -it btc-sign bitcoin-cli'
   ```

2. create wallets on bitcoind respectively

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

3. load wallet (required if btc containers restarted)

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

4. operation

- see [Operation Example](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/btc/OperationExample.md)
