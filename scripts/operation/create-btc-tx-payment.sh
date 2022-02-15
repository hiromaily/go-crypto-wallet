#!/bin/sh

# - prepare coin for payment address

# - data in payment_request table is required
# - receiver_address should be replaced from client address to retrieve coin
# - and payment_id in payment_request should be NULL

set -eu

ENCRYPTED="${1:?false}"

# reset payment_request
echo "------------------------------------------------"
echo 'reset payment_request'
echo "------------------------------------------------"
#watch create db
# rerun the below command to reset
#```
# docker compose exec watch-db mysql -u root -proot  -e "$(cat ./docker/mysql/sqls/payment_request.sql)"
#```

# create unsigned tx
echo "------------------------------------------------"
echo 'create payment tx'
echo "------------------------------------------------"
tx_file=$(watch create payment)
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
# FIXME: somehow passphrase is not required because wif is used
#sign -wallet sign1 api walletpassphrase -passphrase test
tx_file_signed2=`sign1 -wallet sign1 sign -file "${tx_file_signed##*\[fileName\]: }"`
#sign -wallet sign1 api walletlock

# sign on sign wallet for 3rd signature
echo "------------------------------------------------"
echo 'sign on 3rd '${tx_file_signed##*\[fileName\]: }
echo "------------------------------------------------"
# FIXME: somehow passphrase is not required because wif is used
#sign2 -wallet sign2 api walletpassphrase -passphrase test
tx_file_signed3=`sign2 -wallet sign2 sign -file "${tx_file_signed##*\[fileName\]: }"`
#sign -wallet sign1 api walletlock

# send signed tx
echo "------------------------------------------------"
echo 'send tx '${tx_file_signed3##*\[fileName\]: }
echo "------------------------------------------------"
tx_id=`watch send -file "${tx_file_signed3##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }
