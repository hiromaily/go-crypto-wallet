#!/bin/sh

# prepare coin for client address
# https://goerli-faucet.slock.it/

set -eu

watch -coin eth create db

# create unsigned tx
echo 'create payment tx'
tx_file=$(watch -coin eth create payment)

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen -coin eth sign -file "${tx_file##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
watch -coin eth send -file "${tx_file_signed##*\[fileName\]: }"
