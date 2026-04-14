### ERC20 Token Setup

- deploy ERC-20 token contract if needed
- Original ERC-20 token is [here](https://github.com/hiromaily/go-crypto-wallet/tree/main/web/erc20-token). See [`scripts/operation/deploy-token.sh`](https://github.com/hiromaily/go-crypto-wallet/blob/main/scripts/operation/deploy-token.sh)

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
