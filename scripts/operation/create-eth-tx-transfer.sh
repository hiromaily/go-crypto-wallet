#!/bin/sh

# prepare coin for client address
# https://goerli-faucet.slock.it/

set -eu

# create unsigned tx
echo 'create deposit tx'
tx_file=$(watch -coin eth create transfer -account1 deposit -account2 payment)

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen -coin eth sign -file "${tx_file##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
watch -coin eth send -file "${tx_file_signed##*\[fileName\]: }"
