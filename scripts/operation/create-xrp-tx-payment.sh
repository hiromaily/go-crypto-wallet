#!/bin/sh

# prepare coin for client address
# https://xrpl.org/xrp-testnet-faucet.html

# get conf from faucet address
# watch -coin xrp api sendcoin -address raWG2eo1tEXwN4HtGFJCagvukC2nBuiHxC

set -eu

#watch -coin xrp create db

# create unsigned tx
echo 'create deposit tx'
tx_file=$(watch -coin xrp create payment)

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen -coin xrp sign -file "${tx_file##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
watch -coin xrp send -file "${tx_file_signed##*\[fileName\]: }"
