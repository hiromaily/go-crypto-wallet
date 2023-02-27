#!/bin/sh

# [WIP] prepare coin for client address

set -eu

# create unsigned tx
echo 'create transfer tx'
tx_file=$(watch create transfer -account1 deposit -account2 payment)
if [ "$(echo $tx_file | grep 'No utxo')" ]; then
	echo 'No utxo'
	exit 0
fi

# sign on keygen wallet for 1st signature
echo 'sign on 1st '${tx_file##*\[fileName\]: }
tx_file_signed=$(keygen sign -file "${tx_file##*\[fileName\]: }")

# sign on sign wallet for 2nd signature
echo 'sign on 2nd '${tx_file_signed##*\[fileName\]: }
tx_file_signed2=$(sign -wallet sign1 sign -file "${tx_file_signed##*\[fileName\]: }")

# sign on sign wallet for 3rd signature
#echo 'sign on 3rd '${tx_file_signed##*\[fileName\]: }
#tx_file_signed3=`sign2 -wallet sign2 sign -file "${tx_file_signe2d##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed2##*\[fileName\]: }
tx_id=$(watch send -file "${tx_file_signed2##*\[fileName\]: }")
echo 'txID:'${tx_id##*txID: }
