#!/bin/sh

# Prerequisite
# create wallet
# $ make create-wallets
#
# encrypt wallet
#  $ make encrypt-wallets

# reset database and wallet if you want to start from the beginning
# $ make rm-local-wallet-dat
# $ make create-wallets

set -eu
# debug
#set -eux


COIN="${1:?btc}"
ENCRYPTED="${2:?false}"
SIGN_WALLET_NUM=${3:?5}

###############################################################################
# keygen wallet
###############################################################################
# create seed
echo "------------------------------------------------"
echo "create seed"
echo "------------------------------------------------"
keygen create seed

# create hdkey for client, deposit, payment account
echo "------------------------------------------------"
echo "create hdkey for client, deposit, payment account"
echo "------------------------------------------------"
for account in client deposit payment stored; do
  keygen -coin ${COIN} create hdkey -account $account -keynum 10
done

# import generated private key into keygen wallet
echo "------------------------------------------------"
echo "import generated private key into keygen wallet"
echo "------------------------------------------------"
# if wallet is encrypted, walletpassphrase is required before
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletpassphrase -passphrase test
fi
for account in client deposit payment stored; do
  # FIXME: error occurred => done
  # fail to call ImportPrivKeyRescan(): -18: No wallet is loaded.
  # Load a wallet using loadwallet or create a new one with createwallet.
  # (Note: A default wallet is no longer automatically created)
  keygen -coin ${COIN} import privkey -account $account
done
if [ "$ENCRYPTED" = "true" ]; then
  keygen api walletlock
fi

###############################################################################
# sign wallet
###############################################################################
# create seed
echo "------------------------------------------------"
echo "create seed"
echo "------------------------------------------------"
sign create seed

# create hdkey for authorization
echo "------------------------------------------------"
echo "create hdkey for authorization"
echo "------------------------------------------------"
for i in $(seq 1 $SIGN_WALLET_NUM); do
  echo $i
  sign$i -coin ${COIN} -wallet sign$i create hdkey
done

# import generated private key into sign wallet
echo "------------------------------------------------"
echo "import generated private key into sign wallet"
echo "------------------------------------------------"
# if wallet is encrypted, walletpassphrase is required before
for i in $(seq 1 $SIGN_WALLET_NUM); do
  if [ "$ENCRYPTED" = "true" ]; then
    sign$i -coin ${COIN} -wallet sign$i api walletpassphrase -passphrase test
  fi
  sign$i -coin ${COIN} -wallet sign$i import privkey
  if [ "$ENCRYPTED" = "true" ]; then
    sign$i -coin ${COIN} -wallet sign$i api walletlock
  fi
done

# export full-pubkey as csv file
echo "------------------------------------------------"
echo "export full-pubkey as csv file"
echo "------------------------------------------------"
# sign -wallet sign1 export fullpubkey
file_fullpubkey_auth1=$(sign1 -coin "${COIN}" -wallet sign1 export fullpubkey)
file_fullpubkey_auth2=$(sign2 -coin "${COIN}" -wallet sign2 export fullpubkey)
file_fullpubkey_auth3=$(sign3 -coin "${COIN}" -wallet sign3 export fullpubkey)
file_fullpubkey_auth4=$(sign4 -coin "${COIN}" -wallet sign4 export fullpubkey)
file_fullpubkey_auth5=$(sign5 -coin "${COIN}" -wallet sign5 export fullpubkey)


###############################################################################
# keygen wallet
###############################################################################
# import full-pubkey
echo "------------------------------------------------"
echo "import full-pubkey"
echo "------------------------------------------------"
# keygen import fullpubkey -file ./data/pubkey/auth1_1588399093997165000.csv
keygen -coin ${COIN} import fullpubkey -file ${file_fullpubkey_auth1##*\[fileName\]: }
keygen -coin ${COIN} import fullpubkey -file ${file_fullpubkey_auth2##*\[fileName\]: }
keygen -coin ${COIN} import fullpubkey -file ${file_fullpubkey_auth3##*\[fileName\]: }
keygen -coin ${COIN} import fullpubkey -file ${file_fullpubkey_auth4##*\[fileName\]: }
keygen -coin ${COIN} import fullpubkey -file ${file_fullpubkey_auth5##*\[fileName\]: }

# create multisig address
echo "------------------------------------------------"
echo "create multisig address"
echo "------------------------------------------------"
keygen -coin ${COIN} create multisig -account deposit
keygen -coin ${COIN} create multisig -account payment
keygen -coin ${COIN} create multisig -account stored

# export address
echo "------------------------------------------------"
echo "export address"
echo "------------------------------------------------"
file_address_client=$(keygen -coin "${COIN}" export address -account client)
file_address_deposit=$(keygen -coin "${COIN}" export address -account deposit)
file_address_payment=$(keygen -coin "${COIN}" export address -account payment)
file_address_stored=$(keygen -coin "${COIN}" export address -account stored)


###############################################################################
# watch only wallet
###############################################################################
# import addresses generated by keygen wallet
# if wallet.dat is deleted, rescan is required by `-rescan`
echo "------------------------------------------------"
echo "import addresses generated by keygen wallet"
echo "------------------------------------------------"
watch -coin ${COIN} import address -file ${file_address_client##*\[fileName\]: }
watch -coin ${COIN} import address -file ${file_address_deposit##*\[fileName\]: }
watch -coin ${COIN} import address -file ${file_address_payment##*\[fileName\]: }
watch -coin ${COIN} import address -file ${file_address_stored##*\[fileName\]: }
