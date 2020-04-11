
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

# wallet -conf ./data/toml/btc/local_watch_only.toml api estimatefee

# docker-compose up db-btc-wallet

# bitcoin-cli
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo

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
	docker volume rm btc-wallet-db btc-keygen-db btc-signature-db

.PHONY: rm-docker-wallet-dat-all
rm-docker-wallet-dat-all:
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

.PHONY: rm-local-wallet-dat
rm-local-wallet-dat:
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
	./scripts/operation/integration_on_docker.sh 99


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