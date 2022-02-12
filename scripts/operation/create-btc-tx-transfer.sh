#!/bin/sh

set -eu

ENCRYPTED="${1:?false}"
AMOUNT=${1:?0.0001}

# create unsigned tx
echo "------------------------------------------------"
echo 'create transfer tx from deposit to payment'
echo "------------------------------------------------"
tx_file=$(watch create transfer -account1 deposit -account2 payment -amount 0.0001)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet for 1st signature
echo "------------------------------------------------"
echo 'sign on 1st '${tx_file##*\[fileName\]: }
echo "------------------------------------------------"
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletpassphrase -passphrase test
fi
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletlock
fi

# sign on sign wallet for 2nd signature
echo "------------------------------------------------"
echo 'sign on 2nd '${tx_file_signed##*\[fileName\]: }
echo "------------------------------------------------"
tx_file_signed2=`sign1 -wallet sign1 sign -file "${tx_file_signed##*\[fileName\]: }"`

# sign on sign wallet for 3rd signature
#echo 'sign on 3rd '${tx_file_signed##*\[fileName\]: }
#tx_file_signed3=`sign2 -wallet sign2 sign -file "${tx_file_signed2##*\[fileName\]: }"`

# send signed tx
echo "------------------------------------------------"
echo 'send tx '${tx_file_signed2##*\[fileName\]: }
echo "------------------------------------------------"
tx_id=`watch send -file "${tx_file_signed2##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }
