
###############################################################################
# Initial
###############################################################################
goget:
	go get -u -d -v ./...


###############################################################################
# Build
###############################################################################
bld:
	go build -o wallet ./cmd/wallet/main.go
	go build -o coldwallet1 ./cmd/coldwallet1/main.go
	go build -o coldwallet2 ./cmd/coldwallet2/main.go

bld-windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/wallet.exe ./cmd/wallet/main.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/coldwallet1.exe ./cmd/coldwallet1/main.go
	#GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/coldwallet2.exe ./cmd/coldwallet2/main.go
	zip -r ./bin/windows_amd64/wallet.zip ./bin/windows_amd64/wallet.exe
	zip -r ./bin/windows_amd64/coldwallet1.zip ./bin/windows_amd64/coldwallet1.exe
	#zip -r ./bin/windows_amd64/coldwallet2.zip ./bin/windows_amd64/coldwallet2.exe
	rm -f ./bin/windows_amd64/wallet.exe
	rm -f ./bin/windows_amd64/coldwallet1.exe
	#rm -f ./bin/windows_amd64/coldwallet2.exe


###############################################################################
# Run 入金
###############################################################################
# TODO:定期的に実行して、動作を確認すること(これを自動化しておきたい)

# 入金データを集約し、未署名のトランザクションを作成する
create-unsigned: bld
	./wallet -m 1

# 入金データを集約し、未署名のトランザクションを作成する(更に手数料を調整したい場合)
create-unsigned-fee: bld
	./wallet -m 1 -f 1.5

# 未署名のトランザクションに署名する
sign: bld
	./coldwallet1 -m 5 -i ./data/tx/receipt/receipt_8_unsigned_1534832793024491932

# 署名済トランザクションを送信する
send: bld
	./wallet -m 3 -i ./data/tx/receipt/receipt_8_signed_1534832879778945174

# 送金ステータスを監視し、6confirmationsになったら、statusをdoneに更新する
	./wallet -m 10


# テストデータ作成のために入金の一連の流れをまとめて実行する
create-receipt-all: bld
	./wallet -m 20


###############################################################################
# Run 出金
###############################################################################
# TODO:定期的に実行して、動作を確認すること(これを自動化しておきたい)

# 出金データから出金トランザクションを作成する
create-payment: bld
	./wallet -m 2

# 出金データから出金トランザクションを作成する(更に手数料を調整したい場合)
create-payment-fee: bld
	./wallet -m 2 -f 1.5

# 出金用に未署名のトランザクションに署名する
sign-payment: bld
	./coldwallet1 -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772

# 出金用に署名済トランザクションを送信する
send-payment: bld
	./wallet -m 3 -i ./data/tx/payment/payment_3_signed_1534833088943126101


# テストデータ作成のために出金の一連の流れをまとめて実行する
create-payment-all: bld
	./wallet -m 21


###############################################################################
# Run 送金監視
###############################################################################
detect-sent-transaction:
	./wallet -m 10


###############################################################################
# Run 各種Debug機能
###############################################################################
# 出金依頼データの再利用のため、DBを書き換える
run-reset:
	./wallet -d -m 11

# 現在の手数料算出(estimatesmartfee)
run-fee:
	./wallet -d -m 2

# ネットワーク情報取得(getnetworkinfo)
run-info:
	./wallet -d -m 4


###############################################################################
# Run Watch only wallet: アドレス機能
###############################################################################
import-pub:
	./wallet -m 11 -i ./data/pubkey/client_1535423628425011000.csv
	#追加されたアドレスを確認するには、`getaddressesbyaccount ""`コマンド


###############################################################################
# Run coldwallet1: Key生成 機能
###############################################################################
# 出金依頼データの再利用のため、DBを書き換える
gen-seed:
	./coldwallet1 -d -m 2

gen-client:
	./coldwallet1 -d -m 3

gen-receipt:
	./coldwallet1 -d -m 4

gen-payment:
	./coldwallet1 -d -m 5

import-priv:
	./coldwallet1 -d -m 10

export-pub:
	./coldwallet1 -d -m 11


###############################################################################
# Test
###############################################################################
gotest:
	go test -v ./...


###############################################################################
# Docker and compose
###############################################################################
bld-docker-go:
	docker build --no-cache -t cayenne-wallet-go:1.10.3 -f ./docker/golang/Dockerfile .


###############################################################################
# Bitcoin core
###############################################################################
bitcoin-run:
	bitcoind -daemon

bitcoin-stop:
	bitcoin-cli stop


###############################################################################
# Utility
###############################################################################
.PHONY: clean
clean:
	rm -rf wallet coldwallet1