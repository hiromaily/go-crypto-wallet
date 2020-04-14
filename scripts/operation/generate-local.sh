#!/bin/sh

#After reset database and `docker-compose up`

# reset wallet.dat
rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/wallet.dat
sleep 5

###############################################################################
# keygen wallet
###############################################################################
# create seed
keygen create seed

# create hdkey for client, receipt, payment account
keygen create hdkey -account client -keynum 10
keygen create hdkey -account receipt -keynum 10
keygen create hdkey -account payment -keynum 10
keygen create hdkey -account stored -keynum 10

# import generated private key into keygen wallet
keygen import privkey -account client
keygen import privkey -account receipt
keygen import privkey -account payment
keygen import privkey -account stored

# export created public address as csv
keygen export address -account client
keygen export address -account receipt
keygen export address -account payment
keygen export address -account stored


###############################################################################
# sign wallet
###############################################################################
# create seed
sign create seed

# create hdkey for authorization
sign create hdkey

# import generated private key into sign wallet
sign import privkey
# done

# import pubkey exported from keygen wallet into sign wallet
sign import address -account receipt -file receipt_1_1586831083436291000.csv
sign import address -account payment -file payment_1_1586831473462845000.csv
sign import address -account stored -file stored_1_1586834862891724000.csv


#coldwallet1からexportしたAccountのpublicアドレスをcoldWallet2にimportする
coldwallet2 -k -m 30 -i ./data/pubkey/receipt_1_xxx.csv -a receipt
coldwallet2 -k -m 30 -i ./data/pubkey/payment_1_xxx.csv -a payment
coldwallet2 -k -m 30 -i ./data/pubkey/quoine_1_xxx.csv -a quoine
coldwallet2 -k -m 30 -i ./data/pubkey/fee_1_xxx.csv -a fee
coldwallet2 -k -m 30 -i ./data/pubkey/stored_1_xxx.csv -a stored

#`addmultisigaddress`を実行し、multisigアドレスを生成する。パラメータは、accountのアドレス、authorizationのアドレス
coldwallet2 -k -m 40 -a receipt
coldwallet2 -k -m 40 -a payment
coldwallet2 -k -m 40 -a quoine
coldwallet2 -k -m 40 -a fee
coldwallet2 -k -m 40 -a stored

#作成したAccountのMultisigアドレスをcsvファイルとしてexportする
coldwallet2 -k -m 50 -a receipt
coldwallet2 -k -m 50 -a payment
coldwallet2 -k -m 50 -a quoine
coldwallet2 -k -m 50 -a fee
coldwallet2 -k -m 50 -a stored


###############################################################################
#coldwallet1
###############################################################################
#coldwallet2からexportしたAccountのmultisigアドレスをcoldWallet1にimportする
coldwallet1 -k -m 40 -i ./data/pubkey/receipt_2_xxx.csv -a receipt
coldwallet1 -k -m 40 -i ./data/pubkey/payment_2_xxx.csv -a payment
coldwallet1 -k -m 40 -i ./data/pubkey/quoine_2_xxx.csv -a quoine
coldwallet1 -k -m 40 -i ./data/pubkey/fee_2_xxx.csv -a fee
coldwallet1 -k -m 40 -i ./data/pubkey/stored_2xxx.csv -a stored

#multisigのimport後、AccountのMultisigをcsvファイルとしてexportする
coldwallet1 -k -m 50 -a receipt
coldwallet1 -k -m 50 -a payment
coldwallet1 -k -m 50 -a quoine
coldwallet1 -k -m 50 -a fee
coldwallet1 -k -m 50 -a stored


###############################################################################
#watch only wallet
###############################################################################
#coldwalletで生成したAccountのアドレスをwalletにimportする
#wallet -k -m 1 -x -i ./data/pubkey/client_1_xxx.csv -a client #-x rescan=true(coreのwallet.datをリセットした場合)
wallet -k -m 1 -i ./data/pubkey/client_1_xxx.csv -a client
wallet -k -m 1 -i ./data/pubkey/receipt_3_xxx.csv -a receipt
wallet -k -m 1 -i ./data/pubkey/payment_3_xxx.csv -a payment
wallet -k -m 1 -i ./data/pubkey/quoine_3_xxx.csv -a quoine
wallet -k -m 1 -i ./data/pubkey/fee_3_xxx.csv -a fee
wallet -k -m 1 -i ./data/pubkey/stored_3_xxx.csv -a stored

#検証用に出金データを作成する
wallet -d -m 1
