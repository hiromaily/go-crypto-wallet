# Cosmos Hub
token of the Cosmos Hub is the `ATOM`

## References
- [Introduction](https://hub.cosmos.network/main/hub-overview/overview.html)
- [github CosmosHub/Gaia](https://github.com/cosmos/gaia)
- [Cosmos Hub Testnets](https://github.com/cosmos/testnets)
- [Join the Public Testnet(github)](https://github.com/cosmos/gaia/blob/main/docs/hub-tutorials/join-testnet.md)
- [Join the Public Testnet](https://hub.cosmos.network/main/hub-tutorials/join-testnet.html)
- [vega-testnet](https://vega-explorer.hypha.coop/) ... testnet explorer
- 
## Install Gaiad
- [Installation](https://github.com/cosmos/gaia/blob/main/docs/getting-started/installation.md)

#### Install
```
$ git clone -b v6.0.0 https://github.com/cosmos/gaia
$ cd gaia
$ make install
# cd ../ && rm -rf gaia
```

#### gaiad sub command list
```
$ gaiad --help                                                                                                                                                                               (git)-[master]
Stargate Cosmos Hub App

Usage:
gaiad [command]

Available Commands:
add-genesis-account Add a genesis account to genesis.json
collect-gentxs      Collect genesis txs and output a genesis.json file
config              Create or query an application CLI configuration file
debug               Tool for helping with debugging your application
export              Export state to JSON
gentx               Generate a genesis tx carrying a self delegation
help                Help about any command
init                Initialize private validator, p2p, genesis, and application configuration files
keys                Manage your application's keys
query               Querying subcommands
start               Run the full node
status              Query remote node for status
tendermint          Tendermint subcommands
testnet             Initialize files for a simapp testnet
tx                  Transactions subcommands
unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
validate-genesis    validates the genesis file at the default location or at the location passed as an arg
version             Print the application binary version information

Flags:
-h, --help                help for gaiad
--home string         directory for config and data (default "/Users/hiroki.yasui/.gaia")
--log_format string   The logging format (json|plain) (default "plain")
--log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
--trace               print out full stack trace on errors

Use "gaiad [command] --help" for more information about a command.
```

## Setup node for testnet
```
gaiad init mymoniker

# Prepare genesis file
wget https://github.com/cosmos/vega-test/raw/master/public-testnet/modified_genesis_public_testnet/genesis.json.gz
gzip -d genesis.json.gz
mv genesis.json $HOME/.gaia/config/genesis.json

# Set minimum gas price & peers
cd $HOME/.gaia/config
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001uatom"/' app.toml
sed -i 's/persistent_peers = ""/persistent_peers = "<persistent_peer_node_id_1@persistent_peer_address_1:p2p_port>,<persistent_peer_node_id_2@persistent_peer_address_2:p2p_port>"/' config.toml
```

## Run node
```
gaiad start --x-crisis-skip-assert-invariants
```