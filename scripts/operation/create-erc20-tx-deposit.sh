#!/bin/sh

# prepare token for client address
# https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/Installation.md#ethereum-setup
# WIP

set -eu

# create unsigned tx
echo 'create deposit tx'
tx_file=$(watch -coin hyt create deposit)

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen -coin hyt sign -file "${tx_file##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
watch -coin eth send -file "${tx_file_signed##*\[fileName\]: }"
