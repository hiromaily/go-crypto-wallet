# Ethereum

- [go-ethereum](https://github.com/ethereum/go-ethereum)
- [Getting Started with Geth](https://geth.ethereum.org/docs/getting-started)
- [parity](https://www.parity.io/ethereum/)

## Free ebooks
- [Ethereum Development with Go](https://goethereumbook.org/en/)

## Development
- [Dev mode](https://geth.ethereum.org/getting-started/dev-mode)  
Geth has a development mode which sets up a single node Ethereum test network with a number of options optimized for developing on local machines. You enable it with the --dev argument.

- [JSON RPC](https://github.com/ethereum/wiki/wiki/JSON-RPC)
- [Management APIs](https://github.com/ethereum/go-ethereum/wiki/Management-APIs)
- [Setting up private network or local cluster](https://github.com/ethereum/go-ethereum/wiki/Setting-up-private-network-or-local-cluster)
- [Etherscan for goerli](https://goerli.etherscan.io/)
- [Goerli Testnet Faucet](https://goerli-faucet.slock.it/)
- [Status](https://stats.goerli.net/)

## Parity
- [parity](https://www.parity.io/ethereum/)
- [goerli/testnet](https://github.com/goerli/testnet)
- [homebrew-paritytech](https://github.com/paritytech/homebrew-paritytech)

## Install ethereum and Run
### Install ethereum on MacOS
```
$ brew tap ethereum/ethereum
$ brew install ethereum
```

### Run on testnet
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
- to use UnlockAccount()
```
# --allow-insecure-unlock
```

# Rest API using [HTTPie](https://httpie.org/)
$ http http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
 or
$ http --auth USERNAME:PASSWORD http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
```

## geth useful option
```
$ geth --help                                                                                                                                                         (git)-[master]

USAGE:
   geth [options] command [command options] [arguments...]

VERSION:
   1.9.13-stable

COMMANDS:
   account                            Manage accounts
   attach                             Start an interactive JavaScript environment (connect to node)
   console                            Start an interactive JavaScript environment
   copydb                             Create a local chain from a target chaindata folder
   dump                               Dump a specific block from storage
   dumpconfig                         Show configuration values
   dumpgenesis                        Dumps genesis block JSON configuration to stdout
   export                             Export blockchain into file
   export-preimages                   Export the preimage database into an RLP stream
   import                             Import a blockchain file
   import-preimages                   Import the preimage database from an RLP stream
   init                               Bootstrap and initialize a new genesis block
   inspect                            Inspect the storage size for each type of data in the database
   js                                 Execute the specified JavaScript files
   license                            Display license information
   makecache                          Generate ethash verification cache (for testing)
   makedag                            Generate ethash mining DAG (for testing)
   removedb                           Remove blockchain and state databases
   retesteth                          Launches geth in retesteth mode
   version                            Print version numbers
   wallet                             Manage Ethereum presale wallets
   help, h                            Shows a list of commands or help for one command
```
```
ETHEREUM OPTIONS:
  --config value                      TOML configuration file
  --datadir value                     Data directory for the databases and keystore (default: "/Users/hy/Library/Ethereum")
  --keystore value                    Directory for the keystore (default = inside the datadir)
  --networkid value                   Network identifier (integer, 1=Frontier, 3=Ropsten, 4=Rinkeby, 5=Görli) (default: 1)
  --goerli                            Görli network: pre-configured proof-of-authority test network
  --syncmode value                    Blockchain sync mode ("fast", "full", or "light") (default: fast)
DEVELOPER CHAIN OPTIONS:
  --dev                               Ephemeral proof-of-authority network with a pre-funded developer account, mining enabled
API AND CONSOLE OPTIONS:
  --ipcdisable                        Disable the IPC-RPC server
  --ipcpath value                     Filename for IPC socket/pipe within the datadir (explicit paths escape it)
  --rpc                               Enable the HTTP-RPC server
  --rpcaddr value                     HTTP-RPC server listening interface (default: "localhost")
  --rpcport value                     HTTP-RPC server listening port (default: 8545)
  --rpcapi value                      API's offered over the HTTP-RPC interface
```

## Install parity and Run
```
# Install parity on MacOS
$ brew tap paritytech/paritytech
$ brew install parity

# Run on testnet
$ parity --chain=goerli --geth
 or
$ cp ./data/config/parity/config.toml ${HOME}/Library/Application Support/io.parity.ethereum/
$ parity
```

## geth Rest API
- eth_syncing
```
http http://127.0.0.1:8545 method=eth_syncing params:='[]' id=1
```
- web3_clientVersion
```
http http://127.0.0.1:8545 method=web3_clientVersion params:='[]' id=67
```



## How to implement multisig on Ethereum
- [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
- [corvus-ch/shamir](https://github.com/corvus-ch/shamir)
