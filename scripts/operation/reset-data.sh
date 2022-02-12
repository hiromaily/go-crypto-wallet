#!/bin/sh

set -eu

# reset DB
#docker compose stop watch-db keygen-db sign-db
docker compose rm -f -s watch-db keygen-db sign-db
docker volume rm -f go-crypto-wallet_watch-db
docker volume rm -f go-crypto-wallet_keygen-db
docker volume rm -f go-crypto-wallet_sign-db

# reset bitcoind dat
docker compose stop btc-watch btc-keygen btc-sign
rm -rf ./docker/nodes/btc/data1/signet/wallets/watch
rm -rf ./docker/nodes/btc/data2/signet/wallets/keygen
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign1
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign2
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign3
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign4
rm -rf ./docker/nodes/btc/data3/signet/wallets/sign5
