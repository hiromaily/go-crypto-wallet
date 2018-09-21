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

# 初回実行用、seed及び、アカウント各種のキーを生成する
#1. $ ./tools/integration_on_docker.sh 6
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

# 指定したアカウントのkeyを追加で生成
# generate_additional_key accountName number
# e.g. generate_additional_key client 5
# $ ./tools/integration_on_docker.sh 7
function generate_additional_key() {
    if test "$1" = "" ; then
        echo argument1 is required as account
    fi
    if test "$2" = "" ; then
        echo argument2 is required as key number
    fi
    #TODO:account名の厳密なチェックが必要かも


    #keyを生成
    coldwallet1 -k -m 10 -n $2 -a $1

    #作成したAccountのPrivateKeyをColdWalletにimportする
    coldwallet1 -k -m 20 -a $1

    #作成したAccountのPublicアドレスをcsvファイルとしてexportする"
    file_name=$(coldwallet1 -k -m 30 -a "$1")

    #client, authorization は除外
    if [ $1 != client ] && [ $1 != authorization ]; then
        #coldwallet1からexportしたAccountのpublicアドレスをcoldWallet2にimportする
        coldwallet2 -k -m 30 -i ${file_name##*\[fileName\]: } -a $1

        #`addmultisigaddress`を実行し、multisigアドレスを生成する。パラメータは、accountのアドレス、authorizationのアドレス
        coldwallet2 -k -m 40 -a $1

        #作成したAccountのMultisigアドレスをcsvファイルとしてexportする
        file_name=$(coldwallet2 -k -m 50 -a "$1")

        #coldwallet2からexportしたAccountのmultisigアドレスをcoldWallet1にimportする
        coldwallet1 -k -m 40 -i ${file_name##*\[fileName\]: } -a $1

        #multisigのimport後、AccountのMultisigをcsvファイルとしてexportする
        file_name=$(coldwallet1 -k -m 50 -a "$1")
    fi

    #authorization は除外
    if [ $1 != authorization ]; then
        #coldwalletで生成したAccountのアドレスをwalletにimportする
        wallet -k -m 1 -i ${file_name##*\[fileName\]: } -a $1
    fi
}

# 自動テスト(全ロジックをこのtestにて検証する)
# TODO:送金後の、6confirmationチェックロジックの実装
# $ ./tools/integration_on_docker.sh 8
function auto_testing() {
    #check
    #docker-compose exec btc-wallet bitcoin-cli getnetworkinfo
    #docker-compose exec btc-wallet bitcoin-cli listunspent

    #test用にkeyを新たに追加(TODO:アカウント別にkeyを追記するshが必要)
    generate_additional_key receipt 1
    generate_additional_key payment 2

    #quoineアカウントからpaymentアカウントに出金する
    file_unsigned=$(wallet -t -m 1 -a quoine -z payment)
    #署名1
    file_signed1=`coldwallet1 -s -m 1 -i "${file_unsigned##*\[fileName\]: }"`
    #署名2
    file_signed2=`coldwallet2 -s -m 1 -i "${file_signed1##*\[fileName\]: }"`
    #送信
    wallet -s -m 1 -i ${file_signed2##*\[fileName\]: }

    #run after 6confirmation, so monitoring is required
    #check_confirmation payment

    #検証用出金データをリセット
    wallet -d -m 2
    #出金 + 未署名トランザクション作成
    file_unsigned=$(wallet -p -m 1)
    #署名1
    file_signed1=`coldwallet1 -s -m 1 -i "${file_unsigned##*\[fileName\]: }"`
    #署名2
    file_signed2=`coldwallet2 -s -m 1 -i "${file_signed1##*\[fileName\]: }"`
    #送信
    wallet -s -m 1 -i ${file_signed2##*\[fileName\]: }

    #run after 6confirmation, so monitoring is required

    #入金 + 未署名トランザクション作成
    file_unsigned=$(wallet -r -m 1)
    #署名1
    file_signed1=`coldwallet1 -s -m 1 -i "${file_unsigned##*\[fileName\]: }"`
    #署名2
    file_signed2=`coldwallet2 -s -m 1 -i "${file_signed1##*\[fileName\]: }"`
    #送信
    wallet -s -m 1 -i ${file_signed2##*\[fileName\]: }

    #run after 6confirmation, so monitoring is required

    #receiptアカウントからpaymentアカウントに出金する
    file_unsigned=$(wallet -t -m 1 -a receipt -z payment)
    #署名1
    file_signed1=`coldwallet1 -s -m 1 -i "${file_unsigned##*\[fileName\]: }"`
    #署名2
    file_signed2=`coldwallet2 -s -m 1 -i "${file_signed1##*\[fileName\]: }"`
    #送信
    wallet -s -m 1 -i ${file_signed2##*\[fileName\]: }

    #run after 6confirmation, so monitoring is required

}

# $1 account
function check_confirmation() {
    echo check 6 confirmation
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
elif [ $1 -eq 6 ]; then
    #generate all keys
    generate_all
elif [ $1 -eq 7 ]; then
    #add key
    #generate_additional_key client 5
    #generate_additional_key receipt 5
    #generate_additional_key payment 5
    generate_additional_key quoine 5
elif [ $1 -eq 8 ]; then
    #auto testing
    auto_testing
else
    echo argument is out of range
fi
