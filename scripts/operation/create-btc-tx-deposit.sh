#!/bin/sh

# prepare coin for client address
# https://coinfaucet.eu/en/btc-testnet/

set -eu

# create unsigned tx
echo 'create deposit tx'
tx_file=$(watch create deposit)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
keygen api walletpassphrase -passphrase test
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`
keygen api walletlock

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
tx_id=`watch send -file "${tx_file_signed##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }

# check confirmation
bitcoin-cli -rpcuser=xyz -rpcpassword=xyz -rpcwallet=watch gettransaction ${tx_id##*txID: } | jq .confirmations
