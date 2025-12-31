# ERC20-Token

erc20 token `HYT`

## Deploy to testnet using Geth client

1. generate address/private key on Geth for token manager who can deploy contract and mint token
2. run geth with token manager address

```
geth --goerli --rpc --rpcaddr 0.0.0.0 --rpcapi admin,debug,web3,eth,txpool,net,personal --unlock 0xXXXXXXXXXXXXXXXX --password pw --allow-insecure-unlock
```

3. deploy contract

- token manager address is required in `migrations/2_all_contracts.js` before running

```
truffle migrate --network deploy-dev2 --reset
```

4. mint token to address

- environment variable `NODE_URL`, `CONTRACT_ADDRESS`, `OWNER_ADDRESS` are required. see `.envrc`.

```
# blance
yarn ts-node src/web3.ts --mode balance --address 0xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# transfer
yarn ts-node src/web3.ts --mode transfer --address 0xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --amount 100
```

## Initial settings for development

```
## install
yarn init
yarn global add truffle
yarn add --dev eslint
yarn add --dev prettier
yarn add --dev prettier-plugin-solidity
yarn add --dev solhint
yarn add web3
yarn add @openzeppelin/contracts

## typescrript
yarn add --dev typescript
yarn add --dev ts-node
yarn add --dev @types/node
yarn add --dev @typescript-eslint/parser
yarn add --dev @typescript-eslint/eslint-plugin
yarn add --dev eslint-plugin-compat
yarn add --dev eslint-plugin-eslint-comments
yarn add --dev eslint-plugin-prettier
yarn add --dev eslint-config-prettier

tsc --init


## truffle
truffle init
truffle create contract HyToken
truffle create test testHyToken
```

## Requirements for development

- [Ganache](https://www.trufflesuite.com/ganache)
- [truffle](https://www.trufflesuite.com/docs/truffle/getting-started/installation)
- [openzeppelin-contracts](https://github.com/OpenZeppelin/openzeppelin-contracts)
