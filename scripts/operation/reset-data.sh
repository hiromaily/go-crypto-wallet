#!/bin/sh

set -eu

# reset DB
#docker compose stop wallet-db
docker compose rm -f -s wallet-db
docker volume rm -f go-crypto-wallet_wallet-db

# reset bitcoind dat
docker compose -f compose.btc.yaml stop btc-watch btc-keygen btc-sign1 btc-sign2
rm -rf ./docker/nodes/btc/watch/signet/wallets/watch
rm -rf ./docker/nodes/btc/keygen/signet/wallets/keygen
rm -rf ./docker/nodes/btc/sign1/signet/wallets/sign1
rm -rf ./docker/nodes/btc/sign2/signet/wallets/sign2
rm -rf ./docker/nodes/btc/watch/regtest/wallets/watch
rm -rf ./docker/nodes/btc/keygen/regtest/wallets/keygen
rm -rf ./docker/nodes/btc/sign1/regtest/wallets/sign1
rm -rf ./docker/nodes/btc/sign2/regtest/wallets/sign2
