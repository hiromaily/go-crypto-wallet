###############################################################################
# Watch Only wallet
###############################################################################
###############################################################################
# Run アドレスのImport機能
###############################################################################
# keygen-walletでexportしたpublicアドレスをimportする
import-pub1:
	wallet -k -m 1 -i ./data/pubkey/client_1535423628425011000.csv

import-pub2:
	wallet -k -m 2 -i ./data/pubkey/client_1535423628425011000.csv

import-pub3:
	wallet -k -m 3 -i ./data/pubkey/client_1535423628425011000.csv


###############################################################################
# Run 入金
###############################################################################
# TODO:定期的に実行して、動作を確認すること(これを自動化しておきたい)

# 入金データを集約し、未署名のトランザクションを作成する
create-unsigned: bld
	wallet -r -m 1

# 入金データを集約し、未署名のトランザクションを作成する(更に手数料を調整したい場合)
create-unsigned-fee: bld
	wallet -r -m 1 -f 1.5

# 入金確認のみ[WIP]
check-unsigned: bld
	wallet -r -m 2

# [coldwallet] 未署名のトランザクションに署名する
sign: bld
	coldwallet1 -w 1 -s -m 1 -i ./data/tx/receipt/receipt_8_unsigned_1534832793024491932

# 署名済トランザクションを送信する
send: bld
	wallet -s -m 1 -i ./data/tx/receipt/receipt_8_signed_1534832879778945174

# 送金ステータスを監視し、6confirmationsになったら、statusをdoneに更新する
	wallet -n -m 1


# Debug用
# テストデータ作成のために入金の一連の流れをまとめて実行する
create-receipt-all: bld
	wallet -r -m 10


###############################################################################
# Run 出金
###############################################################################
# TODO:定期的に実行して、動作を確認すること(これを自動化しておきたい)

# 出金データから出金トランザクションを作成する
create-payment: bld
	wallet -p -m 1

# 出金データから出金トランザクションを作成する(更に手数料を調整したい場合)
create-payment-fee: bld
	wallet -p -m 1 -f 1.5


# [coldwallet]出金用に未署名のトランザクションに署名する #出金時の署名は2回
sign-payment1: bld
	coldwallet1 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772

sign-payment2: bld
	coldwallet2 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772


# 出金用に署名済トランザクションを送信する
send-payment: bld
	wallet -s -m 3 -i ./data/tx/payment/payment_3_signed_1534833088943126101


# Debug用
# テストデータ作成のために出金の一連の流れをまとめて実行する
create-payment-all: bld
	wallet -p -m 1


###############################################################################
# Run 送金監視
###############################################################################
detect-sent-transaction:
	wallet -n -m 1


###############################################################################
# Run 各種Debug機能
###############################################################################
# 出金依頼データの作成を行う (coldwallet側で生成したデータをwalletにimport後)
run-create-testdata:
	wallet -d -m 1

# 出金依頼データの再利用のため、DBを書き換える
run-db-reset:
	wallet -d -m 2


###############################################################################
# Run Bitcoin API
###############################################################################
# 現在の手数料算出(estimatesmartfee)
run-fee:
	wallet -d -m 2
	#wallet -c ./data/toml/dev1-btccore01.toml -d -m 2

# ネットワーク情報取得(getnetworkinfo)
run-info:
	wallet -d -m 4

