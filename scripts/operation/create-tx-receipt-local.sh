#!/bin/sh

# prepare coin for client address
# https://coinfaucet.eu/en/btc-testnet/

set -eu

# create unsigned tx
echo 'create receipt tx'
tx_file=$(wallet create receipt)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet
echo 'sign on '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed##*\[fileName\]: }
tx_id=`wallet send -file "${tx_file_signed##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }

# check confirmation
bitcoin-cli -rpcuser=xyz -rpcpassword=xyz -rpcwallet=watch gettransaction ${tx_id##*txID: } | jq .confirmations
