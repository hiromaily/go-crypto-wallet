#!/bin/sh

# - prepare coin for payment address
#  https://coinfaucet.eu/en/btc-testnet/

# - data in payment_request table is required
# - receiver_address should be replaced from client address to retrieve coin
# - and payment_id in payment_request should be NULL

set -eu

# reset payment_request
echo 'reset payment_request'
docker-compose exec btc-wallet-db mysql -u root -proot  -e "$(cat ./docker/mysql/wallet/init.d/payment_request.sql)"

# create unsigned tx
echo 'create payment tx'
tx_file=$(wallet create payment)
if [ "`echo $tx_file | grep 'No utxo'`" ]; then
  echo 'No utxo'
  exit 0
fi

# sign on keygen wallet for 1st signature
echo 'sign on 1st '${tx_file##*\[fileName\]: }
tx_file_signed=`keygen sign -file "${tx_file##*\[fileName\]: }"`

# sign on sign wallet for 2nd signature
echo 'sign on 2nd '${tx_file_signed##*\[fileName\]: }
tx_file_signed2=`sign sign -file "${tx_file_signed##*\[fileName\]: }"`

# send signed tx
echo 'send tx '${tx_file_signed2##*\[fileName\]: }
tx_id=`wallet send -file "${tx_file_signed2##*\[fileName\]: }"`
echo 'txID:'${tx_id}
