
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

.PHONY: generate-db-definition
generate-db-definition:
	sqlboiler --wipe mysql
	sqlboiler --config sqlboiler.keygen.toml --wipe mysql


# ifacemaker is tool which makes interface from struct
.PHONY: generate-go-interface
generate-go-interface:
	#walletrepo
	ifacemaker -f pkg/model/rdb/walletrepo/account_pubkey_repo.go -s WalletRepository -i WalletStorager -p rdb
	ifacemaker -f pkg/model/rdb/walletrepo/payment_request_repo.go -s WalletRepository -i WalletStorager -p rdb
	ifacemaker -f pkg/model/rdb/walletrepo/tx_input_repo.go -s WalletRepository -i WalletStorager -p rdb
	ifacemaker -f pkg/model/rdb/walletrepo/tx_output_repo.go -s WalletRepository -i WalletStorager -p rdb
	ifacemaker -f pkg/model/rdb/walletrepo/tx_repo.go -s WalletRepository -i WalletStorager -p rdb
	#keygenrepo
	ifacemaker -f pkg/model/rdb/keygenrepo/account_key_repo.go -s KeygenRepository -i KeygenStorager -p rdb
	ifacemaker -f pkg/model/rdb/keygenrepo/added_pubkey_history_repo.go -s KeygenRepository -i KeygenStorager -p rdb
	ifacemaker -f pkg/model/rdb/keygenrepo/seed_repo.go -s KeygenRepository -i KeygenStorager -p rdb

# git tag
#git tag v2.0.0 cfeca390b781af79321fb644c056bf6e755fdc7e
#git push origin v2.0.0

###############################################################################
# From inside docker container
###############################################################################
.PHONY: bld-linux
bld-linux:
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/wallet ./cmd/wallet/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/sign ./cmd/signature/main.go

###############################################################################
# Build on local
###############################################################################
.PHONY: bld
bld:
	go build -i -v -o ${GOPATH}/bin/wallet ./cmd/wallet/
	go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -i -v -o ${GOPATH}/bin/sign ./cmd/signature/

.PHONY: bldw
bldw:
	go build -i -v -o ${GOPATH}/bin/wallet ./cmd/wallet/

.PHONY: bldk
bldk:
	go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen/

.PHONY: blds
blds:
	go build -i -v -o ${GOPATH}/bin/sign ./cmd/signature/


run:
	go run ./cmd/wallet/ -conf ./data/config/btc/wallet.toml


###############################################################################
# Test on local
###############################################################################
.PHONY: gotest
gotest:
	go test -v ./...

addr-test:
	go test -tags=integration -v -run pkg/wallets/api/btc/...
	go test -tags=integration -v -run GetAddressInfo pkg/wallets/api/btc/...
	go test -v pkg/wallets/api/btc/... -run GetAddressInfo
	go test -v yaml/yaml_test.go -run TestYAMLTable -log 1
#// +build integration


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

.PHONY: up-docker-btc
up-docker-btc:
	docker-compose up btc-wallet btc-keygen btc-signature

.PHONY: up-docker-db
up-docker-db:
	docker-compose up btc-wallet-db btc-keygen-db btc-signature-db

.PHONY: up-docker-only-watch-wallet
up-docker-only-watch-wallet:
	docker-compose up btc-wallet btc-wallet-db

.PHONY: up-docker-btc-all
up-docker-btc-all: up-docker-btc up-docker-db

#.PHONY: up-docker-apps
#up-docker-apps:
#	docker-compose up watch-only-wallet

# logging and monitoring
.PHONY: up-docker-logger
up-docker-logger:
	docker-compose up fluentd elasticsearch grafana


.PHONY: rm-db-volumes
rm-db-volumes:
	docker rm -f $(docker ps -a --format "{{.Names}}")
	docker volume rm -f $(docker volume ls --format "{{.Name}}")


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


###############################################################################
# wallet.dat
###############################################################################
.PHONY: rm-local-wallet-dat
rm-local-wallet-dat:
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/wallet.dat

.PHONY: rm-docker-wallet-dat-all
rm-docker-wallet-dat-all:
	rm -rf ./docker/btc/data1/testnet3/wallets/wallet.data
	rm -rf ./docker/btc/data2/testnet3/wallets/wallet.data
	rm -rf ./docker/btc/data3/testnet3/wallets/wallet.data


###############################################################################
# Grafana
###############################################################################
# http://localhost:3000


###############################################################################
# auto key generator
###############################################################################
.PHONY: generate-key-local
generate-key-local:
	./scripts/operation/generate-key-local.sh

# preparation
# make clean


###############################################################################
# auto key generator
###############################################################################
.PHONY: reset-payment-request
reset-payment-request:
	mysql -h 127.0.0.1 -u root -p${MYSQL_ROOT_PASSWORD} -P 3307 < ./docker/mysql/wallet/init.d/payment_request.sql

.PHONY: reset-payment-request-docker
reset-payment-request-docker:
	docker-compose exec btc-wallet-db mysql -u root -proot  -e "$(cat ./docker/mysql/wallet/init.d/payment_request.sql)"

###############################################################################
# Operation
###############################################################################
include ./Makefile_wallet_op.mk
include ./Makefile_keygen_op.mk
include ./Makefile_signature_op.mk


###############################################################################
# Utility
###############################################################################
.PHONY: clean
clean: rm-db-volumes rm-local-wallet-dat

#after that, run `make up-docker-db`

# bitcoin-cli
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo
