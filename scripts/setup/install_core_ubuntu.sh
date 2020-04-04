#!/bin/sh

sudo apt-add-repository ppa:bitcoin/bitcoin
sudo apt-get update
sudo apt-get install bitcoind

mkdir ~/.bitcoin

cat <<EOF >> ~/.bitcoin/bitcoin.conf
testnet=1
server=1
rpcuser=hiromaily
rpcpassword=hiromaily
txindex=1
zmqpubrawblock=tcp://127.0.0.1:29000
zmqpubrawtx=tcp://127.0.0.1:2900

rpcport=18332
allowip=127.0.0.1
rpcallowip=111.98.254.212
EOF
