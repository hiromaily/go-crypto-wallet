#!/bin/sh

set -eu

# reset DB
#docker compose stop wallet-mysql
docker compose rm -f -s wallet-mysql
docker volume rm -f go-crypto-wallet_wallet-mysql

# reset bitcoind dat
docker compose -f compose.btc.yaml stop btc-watch btc-keygen btc-sign1 btc-sign2
rm -rf ./docker/nodes/btc/watch/signet/wallets/*
rm -rf ./docker/nodes/btc/keygen/signet/wallets/*
rm -rf ./docker/nodes/btc/sign1/signet/wallets/*
rm -rf ./docker/nodes/btc/sign2/signet/wallets/*
rm -rf ./docker/nodes/btc/watch/regtest/wallets/*
rm -rf ./docker/nodes/btc/keygen/regtest/wallets/*
rm -rf ./docker/nodes/btc/sign1/regtest/wallets/*
rm -rf ./docker/nodes/btc/sign2/regtest/wallets/*
