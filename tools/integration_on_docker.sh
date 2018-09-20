#!/bin/sh

# $ ./tools/integration_on_docker.sh 1
function cold1_generate_key() {
    #seedを生成
    coldwallet1 -k -m 1

    #keyを生成
    coldwallet1 -k -m 10 -n 10 -a client
    coldwallet1 -k -m 10 -n 5  -a receipt
    coldwallet1 -k -m 10 -n 5  -a payment
    coldwallet1 -k -m 10 -n 5  -a quoine
    coldwallet1 -k -m 10 -n 5  -a fee
    coldwallet1 -k -m 10 -n 5  -a stored

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
}

function cold2_generate_key() {
    #seedを生成
    coldwallet2 -k -m 1

    #AuthのKeyを生成
    coldwallet2 -k -m 10

    #作成したAuthのPrivateKeyをColdWalletにimportする
    coldwallet2 -k -m 20
}

function cold2_import_export_keys() {
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
}

function cold1_import_export_keys() {
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
}

function watch_only_import_keys() {
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
}

# $ ./tools/integration_on_docker.sh 99
function generate_all() {
   #seedを生成
    coldwallet1 -k -m 1
    coldwallet2 -k -m 1

    #keyを生成
    coldwallet1 -k -m 10 -n 10 -a client
    coldwallet1 -k -m 10 -n 5  -a receipt
    coldwallet1 -k -m 10 -n 5  -a payment
    coldwallet1 -k -m 10 -n 5  -a quoine
    coldwallet1 -k -m 10 -n 5  -a fee
    coldwallet1 -k -m 10 -n 5  -a stored
    #AuthのKeyを生成
    coldwallet2 -k -m 10

    #作成したAccountのPrivateKeyをColdWalletにimportする
    coldwallet1 -k -m 20 -a client
    coldwallet1 -k -m 20 -a receipt
    coldwallet1 -k -m 20 -a payment
    coldwallet1 -k -m 20 -a quoine
    coldwallet1 -k -m 20 -a fee
    coldwallet1 -k -m 20 -a stored
    #作成したAuthのPrivateKeyをColdWalletにimportする
    coldwallet2 -k -m 20

    #作成したAccountのPublicアドレスをcsvファイルとしてexportする"
    file_client=$(coldwallet1 -k -m 30 -a client)
    file_receipt=$(coldwallet1 -k -m 30 -a receipt)
    file_payment=$(coldwallet1 -k -m 30 -a payment)
    file_quoine=$(coldwallet1 -k -m 30 -a quoine)
    file_fee=$(coldwallet1 -k -m 30 -a fee)
    file_stored=$(coldwallet1 -k -m 30 -a stored)

    #coldwallet1からexportしたAccountのpublicアドレスをcoldWallet2にimportする
    coldwallet2 -k -m 30 -i ${file_receipt##*\[fileName\]: } -a receipt
    coldwallet2 -k -m 30 -i ${file_payment##*\[fileName\]: } -a payment
    coldwallet2 -k -m 30 -i ${file_quoine##*\[fileName\]: } -a quoine
    coldwallet2 -k -m 30 -i ${file_fee##*\[fileName\]: } -a fee
    coldwallet2 -k -m 30 -i ${file_stored##*\[fileName\]: } -a stored

    #`addmultisigaddress`を実行し、multisigアドレスを生成する。パラメータは、accountのアドレス、authorizationのアドレス
    coldwallet2 -k -m 40 -a receipt
    coldwallet2 -k -m 40 -a payment
    coldwallet2 -k -m 40 -a quoine
    coldwallet2 -k -m 40 -a fee
    coldwallet2 -k -m 40 -a stored

    #作成したAccountのMultisigアドレスをcsvファイルとしてexportする
    file_receipt=$(coldwallet2 -k -m 50 -a receipt)
    file_payment=$(coldwallet2 -k -m 50 -a payment)
    file_quoine=$(coldwallet2 -k -m 50 -a quoine)
    file_fee=$(coldwallet2 -k -m 50 -a fee)
    file_stored=$(coldwallet2 -k -m 50 -a stored)

    #coldwallet2からexportしたAccountのmultisigアドレスをcoldWallet1にimportする
    coldwallet1 -k -m 40 -i ${file_receipt##*\[fileName\]: } -a receipt
    coldwallet1 -k -m 40 -i ${file_payment##*\[fileName\]: } -a payment
    coldwallet1 -k -m 40 -i ${file_quoine##*\[fileName\]: } -a quoine
    coldwallet1 -k -m 40 -i ${file_fee##*\[fileName\]: } -a fee
    coldwallet1 -k -m 40 -i ${file_stored##*\[fileName\]: } -a stored

    #multisigのimport後、AccountのMultisigをcsvファイルとしてexportする
    file_receipt=$(coldwallet1 -k -m 50 -a receipt)
    file_payment=$(coldwallet1 -k -m 50 -a payment)
    file_quoine=$(coldwallet1 -k -m 50 -a quoine)
    file_fee=$(coldwallet1 -k -m 50 -a fee)
    file_stored=$(coldwallet1 -k -m 50 -a stored)

    #coldwalletで生成したAccountのアドレスをwalletにimportする
    wallet -k -m 1 -i ${file_client##*\[fileName\]: } -a client
    wallet -k -m 1 -i ${file_receipt##*\[fileName\]: } -a receipt
    wallet -k -m 1 -i ${file_payment##*\[fileName\]: } -a payment
    wallet -k -m 1 -i ${file_quoine##*\[fileName\]: } -a quoine
    wallet -k -m 1 -i ${file_fee##*\[fileName\]: } -a fee
    wallet -k -m 1 -i ${file_stored##*\[fileName\]: } -a stored

    #検証用に出金データを作成する
    wallet -d -m 1
}

set -eu

# make sure parameter
echo prameter1 is $1

#debug
#ret=$(wallet -d -m 4)
##ファイル名取得
#echo ${ret##*\[fileName\]: }


# main
if [ $1 -eq 1 ]; then
    cold1_generate_key
elif [ $1 -eq 2 ]; then
    cold2_generate_key
elif [ $1 -eq 3 ]; then
    cold2_import_export_keys
elif [ $1 -eq 4 ]; then
    cold1_import_export_keys
elif [ $1 -eq 5 ]; then
    watch_only_import_keys
else
    generate_all
fi
