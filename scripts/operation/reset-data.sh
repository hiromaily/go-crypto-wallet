#!/bin/sh

set -eu

# reset DB
#docker compose stop wallet-db
docker compose rm -f -s wallet-db
docker volume rm -f go-crypto-wallet_wallet-db

# reset bitcoind dat
docker compose -f compose.btc.yaml stop btc-watch btc-keygen btc-sign1 btc-sign2
rm -rf ./docker/nodes/btc/data1/signet/wallets/watch
rm -rf ./docker/nodes/btc/data2/signet/wallets/keygen
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign1
rm -rf ./docker/nodes/btc/data4/signet/wallets/sign2
rm -rf ./docker/nodes/btc/data1/regtest/wallets/watch
rm -rf ./docker/nodes/btc/data2/regtest/wallets/keygen
rm -rf ./docker/nodes/btc/data3/regtest/wallets/sign1
rm -rf ./docker/nodes/btc/data4/regtest/wallets/sign2
