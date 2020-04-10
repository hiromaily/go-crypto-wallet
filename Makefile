
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
	go build -i -v -o ${GOPATH}/bin/wallet ./cmd/wallet/
	#go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen-wallet/
	#go build -i -v -o ${GOPATH}/bin/sign ./cmd/signature-wallet/

run: bld
	wallet -conf ./data/toml/btc/local_watch_only.toml

# docker-compose up db-btc-wallet

# bitcoin-cli
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo

# wallet -conf ./data/toml/btc/local_watch_only.toml api estimatefee

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
	docker-compose up btc-wallet btc-keygen btc-signature db-btc-wallet db-btc-keygen db-btc-signature

#bitcoin coreのみ起動
.PHONY: up-docker-core
up-docker-core:
	docker-compose up btc-wallet btc-keygen btc-signature

#データベースのみ起動
.PHONY: up-docker-dbs
up-docker-dbs:
	docker-compose up db-btc-wallet db-btc-keygen db-btc-signature

.PHONY: up-docker-apps
up-docker-apps:
	docker-compose up watch-only-wallet

#ログ系システムのみ起動
.PHONY: up-docker-logger
up-docker-logger:
	docker-compose up fluentd elasticsearch grafana

.PHONY: up-docker-only-watch-wallet
up-docker-only-watch-wallet:
	docker-compose up btc-wallet db-btc-wallet watch-only-wallet

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
# Operation
###############################################################################
include ./Makefile_operation


###############################################################################
# Utility
###############################################################################
.PHONY: clean
clean:
	rm -rf wallet keygen sign