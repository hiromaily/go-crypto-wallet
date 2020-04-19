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
#tx is sent!! txID: a473a8bd3424e2b1e81155bf55e407b06342795cc3779a1cb194a3f532551030

# check confirmation
#${tx_id##*txID: }
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz gettransaction a473a8bd3424e2b1e81155bf55e407b06342795cc3779a1cb194a3f532551030
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz gettransaction a473a8bd3424e2b1e81155bf55e407b06342795cc3779a1cb194a3f532551030 | jq .confirmations

# check confirmation
bitcoin-cli -rpcuser=xyz -rpcpassword=xyz gettransaction ${tx_id##*txID: } | jq .confirmations
