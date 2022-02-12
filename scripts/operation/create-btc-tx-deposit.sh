#!/bin/sh

# prepare coin for client address
# https://coinfaucet.eu/en/btc-testnet/

set -eu

ENCRYPTED="${1:?false}"

CLI_WATCH="docker exec -it btc-watch bitcoin-cli"

# create unsigned tx
echo "------------------------------------------------"
echo "create unsigned tx"
echo "------------------------------------------------"
tx_file=$(watch create deposit)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet
echo "------------------------------------------------"
echo "sign on unsigned tx by keygen wallet "${tx_file##*\[fileName\]: }
echo "------------------------------------------------"
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletpassphrase -passphrase test
fi
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletlock
fi

# send signed tx
echo "------------------------------------------------"
echo "send signed tx "${tx_file_signed##*\[fileName\]: }
echo "------------------------------------------------"
tx_id=`watch send -file "${tx_file_signed##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }

# check confirmation
#$CLI_WATCH -rpcuser=xyz -rpcpassword=xyz -rpcwallet=watch gettransaction ${tx_id##*txID: } | jq .confirmations
