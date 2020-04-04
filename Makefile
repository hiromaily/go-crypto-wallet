
###############################################################################
# Initial
###############################################################################
.PHONY: setup-mac
setup-mac:
	brew install jq

.PHONY: goget
goget:
	go get -u -d -v ./...

.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lintfix
lintfix:
	golangci-lint run --fix

###############################################################################
# From inside docker container
###############################################################################
.PHONY: bld-linux
bld-linux:
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/wallet ./cmd/wallet/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/keygen ./cmd/keygen-wallet/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/sign ./cmd/signature-wallet/main.go

###############################################################################
# Build on local
###############################################################################
.PHONY: bld
bld:
	go build -i -v -o ${GOPATH}/bin/wallet ./cmd/wallet/main.go
	go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen-wallet/main.go
	go build -i -v -o ${GOPATH}/bin/sign ./cmd/signature-wallet/main.go

###############################################################################
# Test on local
###############################################################################
.PHONY: gotest
gotest:
	go test -v ./...


###############################################################################
# Docker and compose
###############################################################################
# build docker images
.PHONY: bld-docker-all
bld-docker-all:
	docker-compose build

.PHONY: bld-docker-go
bld-docker-go:
	docker-compose build base-golang

.PHONY: bld-docker-ubuntu
bld-docker-ubuntu:
	docker-compose build base-ubuntu

.PHONY: bld-docker-btc
bld-docker-btc:
	docker-compose build btc-wallet

#bld-docker-bch:
#	docker-compose -f docker-compose.bch.yml build bch-wallet
#up-docker-bch:
#	docker-compose -f docker-compose.bch.yml up bch-wallet

#bitcoin coreとdbをまとめて起動(基本的にこれを使うことになるはず)
.PHONY: up-local-dev-btc
up-local-dev-btc:
	docker-compose up btc-wallet btc-cold-wallet1 btc-cold-wallet2 btc-wallet-db btc-cold-wallet1-db btc-cold-wallet2-db

#bitcoin coreのみ起動
.PHONY: up-docker-core
up-docker-core:
	docker-compose up btc-wallet btc-cold-wallet1 btc-cold-wallet2

#データベースのみ起動
.PHONY: up-docker-dbs
up-docker-dbs:
	docker-compose up btc-wallet-db btc-cold-wallet1-db btc-cold-wallet2-db

.PHONY: up-docker-apps
up-docker-apps:
	docker-compose up watch-only-wallet

#ログ系システムのみ起動
.PHONY: up-docker-logger
up-docker-logger:
	docker-compose up fluentd elasticsearch grafana

.PHONY: up-docker-only-watch-wallet
up-docker-only-watch-wallet:
	docker-compose up btc-wallet btc-wallet-db watch-only-wallet

.PHONY: clear-db-volumes
clear-db-volumes:
	docker rm -f $(docker ps -a --format "{{.Names}}")
	docker volume rm go-bitcoin_db1 go-bitcoin_db2 go-bitcoin_db3

.PHONY: remove-wallet-dat
remove-wallet-dat:
	rm -rf ./docker/btc/data1/testnet3/wallets/wallet.data
	rm -rf ./docker/btc/data2/testnet3/wallets/wallet.data
	rm -rf ./docker/btc/data3/testnet3/wallets/wallet.data


###############################################################################
# Bitcoin core on local
###############################################################################
.PHONY: bitcoin-run
bitcoin-run:
	bitcoind -daemon

.PHONY: bitcoin-stop
bitcoin-stop:
	bitcoin-cli stop

.PHONY: cd-btc-dir
cd-btc-dir:
	cd ~/Library/Application\ Support/Bitcoin

.PHONY: reset-wallet-dat
reset-wallet-dat:
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/wallet.dat


###############################################################################
# Grafana
###############################################################################
# http://localhost:3000


###############################################################################
# Automation on docker
###############################################################################
.PHONY: auto-generation
auto-generation:
	./tools/integration_on_docker.sh 99


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


###############################################################################
# Cold wallet1 (Client/Receipt/PaymentのKey管理)
###############################################################################
###############################################################################
# Run coldwallet1: Key生成 機能
###############################################################################
# development
develop:
	coldwallet1 -d

# seedを生成する
gen-seed:
	coldwallet1 -k -m 1


# Clientのkeyを生成する
gen-client-key:
	coldwallet1 -k -m 10

# Receiptのkeyを生成する
gen-receipt-key:
	coldwallet1 -k -m 11

# Paymentのkeyを生成する
gen-payment:
	coldwallet1 -k -m 12


# Clientのprivate keyをcoldwalletに登録する
add-client-priv-key:
	coldwallet1 -k -m 20

# Receiptのprivate keyをcoldwalletに登録する
add-receipt-priv-key:
	coldwallet1 -k -m 21

# Paymentのprivate keyをcoldwalletに登録する
add-payment-priv-key:
	coldwallet1 -k -m 22


# Clientのpubアドレスをexportする
export-client-pub-key:
	coldwallet1 -k -m 30

# Receiptのpubアドレスをexportする
export-receipt-pub-key:
	coldwallet1 -k -m 31

# Paymentのpubアドレスをexportする
export-payment-pub-key:
	coldwallet1 -k -m 32


# Receiptのmultisigアドレスをimportする
import-receipt-multisig-address:
	coldwallet1 -k -m 40

# Paymentのmultisigアドレスをimportする
import-payment-multisig-address:
	coldwallet1 -k -m 41



###############################################################################
# Cold wallet2 (Authorizationのキー/ Receipt/PaymentのMultisigアドレス管理)
###############################################################################
###############################################################################
# Run coldwallet2: Key生成 機能
###############################################################################
# seedを生成する
gen-seed2:
	coldwallet2 -k -m 1


# Authorizationのkeyを生成する
gen-authorization-key:
	coldwallet2 -k -m 13


# Authorizationのprivate keyをcoldwalletに登録する
add-authorization-priv-key:
	coldwallet2 -k -m 23


# ReceiptのPublicアドレス(full public key)をimportする
import-receipt-pub-key:
	coldwallet2 -k -m 33 -i ./data/pubkey/receipt_1535613888391656000.csv

# PaymentのPublicアドレス(full public key)をimportする
import-payment-pub-key:
	coldwallet2 -k -m 34 -i ./data/pubkey/payment_1535613934762230000.csv


# Receiptのmultisigアドレスを生成し、登録する
add-multisig-receipt:
	coldwallet2 -k -m 50

# Paymentのmultisigアドレスを生成し、登録する
add-multisig-payment:
	coldwallet2 -k -m 51


# Receiptのmultisigアドレスをexportする
export-multisig-receipt:
	coldwallet2 -k -m 60

# Paymentのmultisigアドレスをexportする
export-multisig-payment:
	coldwallet2 -k -m 61


###############################################################################
# Utility
###############################################################################
.PHONY: clean
clean:
	rm -rf wallet keygen sign