#!/bin/sh

#After `docker-compose up`

###############################################################################
#coldwallet1
###############################################################################
#seedを生成
coldwallet1 -k -m 1

#keyを生成
coldwallet1 -k -m 10 -n 10 -a client  #client
coldwallet1 -k -m 10 -n 5  -a receipt #receipt
coldwallet1 -k -m 10 -n 5  -a payment #payment
coldwallet1 -k -m 10 -n 5  -a quoine  #quoine
coldwallet1 -k -m 10 -n 5  -a fee     #fee
coldwallet1 -k -m 10 -n 5  -a store   #stored

#作成したAccountのPrivateKeyをColdWalletにimportする
coldwallet1 -k -m 20 -a client
coldwallet1 -k -m 20 -a receipt
coldwallet1 -k -m 20 -a payment
coldwallet1 -k -m 20 -a quoine
coldwallet1 -k -m 20 -a fee
coldwallet1 -k -m 20 -a stored

#作成したAccountのPublicアドレスをcsvファイルとしてexportする"
coldwallet1 -k -m 30 -a client
coldwallet1 -k -m 30 -a receipt
coldwallet1 -k -m 30 -a payment
coldwallet1 -k -m 30 -a quoine
coldwallet1 -k -m 30 -a fee
coldwallet1 -k -m 30 -a stored


###############################################################################
#coldwallet2
###############################################################################
#seedを生成
coldwallet2 -k -m 1

#AuthのKeyを生成
coldwallet2 -k -m 10

#作成したAuthのPrivateKeyをColdWalletにimportする
coldwallet2 -k -m 20

#coldwallet1からexportしたAccountのpublicアドレスをcoldWallet2にimportする
coldwallet2 -k -m 30 -i ./data/pubkey/xxx.csv -a receipt
coldwallet2 -k -m 30 -i ./data/pubkey/xxx.csv -a payment
coldwallet2 -k -m 30 -i ./data/pubkey/xxx.csv -a quoine
coldwallet2 -k -m 30 -i ./data/pubkey/xxx.csv -a fee
coldwallet2 -k -m 30 -i ./data/pubkey/xxx.csv -a stored

#`addmultisigaddress`を実行する。パラメータは、accountのアドレス、authorizationのアドレス
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
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv -a receipt
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv -a payment
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv -a quoine
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv -a fee
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv -a stored

#multisigのimport後、AccountのMultisigをcsvファイルとしてexportする
coldwallet1 -k -m 50 -a receipt
coldwallet1 -k -m 50 -a payment
coldwallet1 -k -m 50 -a quoine
coldwallet1 -k -m 50 -a fee
coldwallet1 -k -m 50 -a stored


###############################################################################
#watch only wallet
###############################################################################
wallet -k -m 1 -x -i ./data/pubkey/xxx.csv
wallet -k -m 2 -x -i ./data/pubkey/xxx.csv
wallet -k -m 3 -x -i ./data/pubkey/xxx.csv

wallet -d -m 1
