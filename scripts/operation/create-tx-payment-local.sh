#!/bin/sh

# - prepare coin for payment address
#  https://coinfaucet.eu/en/btc-testnet/

# - data in payment_request table is required
# - receiver_address should be replaced from client address to retrieve coin
# - and payment_id in payment_request should be NULL

set -eu

# reset payment_request
echo 'reset payment_request'
#docker-compose exec btc-watch-db mysql -u root -proot  -e "$(cat ./docker/mysql/sqls/payment_request.sql)"
watch db create

# create unsigned tx
echo 'create payment tx'
tx_file=$(watch create payment)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet for 1st signature
echo 'sign on 1st '${tx_file##*\[fileName\]: }
keygen api walletpassphrase -passphrase test
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`
keygen api walletlock

# sign on sign wallet for 2nd signature
# FIXME: somehow passphrase is not required because wif is used
echo 'sign on 2nd '${tx_file_signed##*\[fileName\]: }
#sign -wallet sign1 api walletpassphrase -passphrase test
tx_file_signed2=`sign -wallet sign1 sign -file "${tx_file_signed##*\[fileName\]: }"`
#sign -wallet sign1 api walletlock

# sign on sign wallet for 3rd signature
# FIXME: somehow passphrase is not required because wif is used
echo 'sign on 3rd '${tx_file_signed##*\[fileName\]: }
#sign2 -wallet sign2 api walletpassphrase -passphrase test
tx_file_signed3=`sign2 -wallet sign2 sign -file "${tx_file_signed##*\[fileName\]: }"`
#sign -wallet sign1 api walletlock

# send signed tx
echo 'send tx '${tx_file_signed3##*\[fileName\]: }
tx_id=`watch send -file "${tx_file_signed3##*\[fileName\]: }"`
echo 'txID:'${tx_id##*txID: }
