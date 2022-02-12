# Operation Example
## Scenario 1
- 3 nodes are running `btc-watch`,`btc-keygen`,`btc-sign` respectively.
- To use multisig functionality, multiple auth account is required.
- Set 5 auth accounts on [account.toml](https://github.com/hiromaily/go-crypto-wallet/blob/master/data/config/account.toml)

### Prerequisite
See [Installation](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/Installation.md#bitcoind-setup)
1. `bitcoin-cli` need to be ready
2. [account.toml](https://github.com/hiromaily/go-crypto-wallet/blob/master/data/config/account.toml) configuration
3. create wallets

### 1. Generate Key for Bitcoin
```
$ ./scripts/operation/generate-btc-key.sh btc false 5
```