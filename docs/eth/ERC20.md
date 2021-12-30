# ERC20 Token
## How to try ETH20 token
1. create eth account and register by running `generate-eth-key.sh` in ./scripts/operation/
2. change value in `coin` column to `hyc` as example in `address` table on watch-db.
```
UPDATE `watch`.`address` SET `coin` = 'hyt';
```
3. deploy contract to your ethereum network from [erc20-token](https://github.com/hiromaily/erc20-token)
```
yarn run deploy-dev2
```
4. copy contract address and set into `contract_address` in eth_watch.toml
5. set `master_address` as well
6. setup [erc20-token](https://github.com/hiromaily/erc20-token)
    - `.envrc` includes `CONTRACT_ADDRESS` and `OWNER_ADDRESS` which need to be changed.
7. transfer token from master address to specific address
```
yarn ts-node src/web3.ts --mode transfer --address 0x5c2415367A9558Cb95926619337859aD64beA345 --amount 100
```
8. check balance
```
yarn ts-node src/web3.ts --mode balance --address 0x5c2415367A9558Cb95926619337859aD64beA345
```
9. run command `watch -coin hyt create deposit`
