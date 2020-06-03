
###############################################################################
# Initial
###############################################################################
.PHONY: setup-mac
setup-mac:
	brew install jq mkcert

.PHONY: goget
goget:
	go get -u -d -v ./...


.PHONY: install-ssl
install-ssl:
	mkcert -install
	mkcert localhost 127.0.0.1

.PHONY: install-sqlboiler
install-sqlboiler: SQLBOILER_VERSION=3.7.1
install-sqlboiler:
	echo SQLBOILER_VERSION is $(SQLBOILER_VERSION)
	go get github.com/volatiletech/sqlboiler@v$(SQLBOILER_VERSION)
	go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql@v$(SQLBOILER_VERSION)

.PHONY: install-sqlboiler2
install-sqlboiler2:
	cd ${GOPATH}/src/github.com/volatiletech/sqlboiler
	#git checkout refs/tags/v3.7.1
	git pull
	go get ./...
	go build -i -v -o ${GOPATH}/bin/sqlboiler .

# https://github.com/volatiletech/sqlboiler/issues/633
# https://github.com/volatiletech/sqlboiler/issues/607
# sqlboiler 3.6.1 cannot convert type: types.Decimal => named tag: v3.3.1 in github.com/ericlagergren/decimal works
# https://github.com/golang/go/issues/35732
# https://forum.golangbridge.org/t/solved-error-when-using-go-modules-in-existing-project/15908/9
.PHONY: update-decimal
update-decimal:
	#go get -u github.com/ericlagergren/decimal@v3.3.1 => error
	#go mod edit -require github.com/ericlagergren/decimal@v3.3.1
	#GONOSUMDB=github.com/ericlagergren/decimal go install github.com/volatiletech/sqlboiler
	go get github.com/ericlagergren/decimal@v0.0.0-20181231230500-73749d4874d5

.PHONY: imports
imports:
	./scripts/imports.sh

# FIXME: just after updating package outside from this repository, `go get` doesn't update that package for a while
#goget:
#	go get github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi@ca80219

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lintfix
lintfix:
	golangci-lint run --fix

# FIXME: file is not generated with --templates option if files are existing
# As workaround, modify files in ./templates/..
.PHONY: generate-db-definition
generate-db-definition:
	sqlboiler --wipe \
	--templates ${GOPATH}/src/github.com/volatiletech/sqlboiler/templates,\
	${GOPATH}/src/github.com/volatiletech/sqlboiler/templates_test,\
	${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/templates \
	mysql

.PHONY: sqlboiler
sqlboiler:
	sqlboiler --wipe mysql

# ifacemaker is tool which makes interface from struct
#.PHONY: generate-go-interface
#generate-go-interface:
#	#walletrepo
#	ifacemaker -f pkg/model/rdb/walletrepo/account_pubkey_repo.go -s WalletRepository -i WalletStorager -p rdb
#	ifacemaker -f pkg/model/rdb/walletrepo/payment_request_repo.go -s WalletRepository -i WalletStorager -p rdb
#	ifacemaker -f pkg/model/rdb/walletrepo/tx_input_repo.go -s WalletRepository -i WalletStorager -p rdb
#	ifacemaker -f pkg/model/rdb/walletrepo/tx_output_repo.go -s WalletRepository -i WalletStorager -p rdb
#	ifacemaker -f pkg/model/rdb/walletrepo/tx_repo.go -s WalletRepository -i WalletStorager -p rdb
#	#keygenrepo
#	ifacemaker -f pkg/model/rdb/keygenrepo/account_key_repo.go -s KeygenRepository -i KeygenStorager -p rdb
#	ifacemaker -f pkg/model/rdb/keygenrepo/added_pubkey_history_repo.go -s KeygenRepository -i KeygenStorager -p rdb
#	ifacemaker -f pkg/model/rdb/keygenrepo/seed_repo.go -s KeygenRepository -i KeygenStorager -p rdb

# git tag
#git tag v2.0.0 cfeca390b781af79321fb644c056bf6e755fdc7e
#git push origin v2.0.0

###############################################################################
# From inside docker container
###############################################################################
.PHONY: bld-linux
bld-linux: update-decimal
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/watch ./cmd/watch/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.authName=auth1" -o /go/bin/sign ./cmd/sign/main.go

###############################################################################
# Build on local
###############################################################################
.PHONY: bld
bld: update-decimal
	go build -i -v -o ${GOPATH}/bin/watch ./cmd/watch/
	go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -ldflags "-X main.authName=auth1" -i -v -o ${GOPATH}/bin/sign ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -i -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
	go build -ldflags "-X main.authName=auth3" -i -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
	go build -ldflags "-X main.authName=auth4" -i -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
	go build -ldflags "-X main.authName=auth5" -i -v -o ${GOPATH}/bin/sign5 ./cmd/sign/

.PHONY: bldw
bldw:
	go build -i -v -o ${GOPATH}/bin/watch ./cmd/watch/

.PHONY: bldk
bldk:
	go build -i -v -o ${GOPATH}/bin/keygen ./cmd/keygen/

.PHONY: blds
blds:
	go build -ldflags "-X main.authName=auth1" -i -v -o ${GOPATH}/bin/sign ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -i -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
	go build -ldflags "-X main.authName=auth3" -i -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
	go build -ldflags "-X main.authName=auth4" -i -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
	go build -ldflags "-X main.authName=auth5" -i -v -o ${GOPATH}/bin/sign5 ./cmd/sign/


run:
	go run ./cmd/watch/ -conf ./data/config/watch.toml


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
#// +build integration


###############################################################################
# Docker and compose
###############################################################################
# build docker images
.PHONY: bld-docker-all
bld-docker-all:
	docker-compose build

# build golang image
.PHONY: bld-docker-go
bld-docker-go:
	docker-compose build base-golang

# build ubuntu image
.PHONY: bld-docker-ubuntu
bld-docker-ubuntu:
	docker-compose build base-ubuntu

# build bitcoin core server
.PHONY: bld-docker-btc
bld-docker-btc:
	docker-compose build btc-watch

# build bitcoin cash core server
.PHONY: bld-docker-bch
bld-docker-bch:
	docker-compose -f docker-compose.bch.yml build bch-watch


# run bitcoin core server
.PHONY: up-docker-btc
up-docker-btc:
	docker-compose up btc-watch btc-keygen btc-sign

# run bitcoin cash core server
.PHONY: up-docker-bch
up-docker-bch:
	docker-compose -f docker-compose.bch.yml up bch-watch

# run ethereum node server
.PHONY: up-docker-eth
up-docker-eth:
	docker-compose -f docker-compose.eth.yml up eth-node

# run all databases
.PHONY: up-docker-db
up-docker-db:
	docker-compose up btc-watch-db btc-keygen-db btc-sign-db

# run logging middleware
# logging and monitoring
.PHONY: up-docker-logger
up-docker-logger:
	docker-compose up fluentd elasticsearch grafana

# remove database volumes
.PHONY: rm-db-volumes
rm-db-volumes:
	#docker rm -f $(docker ps -a --format "{{.Names}}")
	#docker volume rm -f $(docker volume ls --format "{{.Name}}")
	#docker-compose down -v
	#docker-compose down
	docker volume rm -f go-crypto-wallet_btc-keygen-db
	docker volume rm -f go-crypto-wallet_btc-sign-db
	docker volume rm -f go-crypto-wallet_btc-watch-db

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
# Grafana
###############################################################################
# http://localhost:3000


###############################################################################
# auto key generator
###############################################################################
.PHONY: generate-btc-key-local
generate-btc-key-local:
	./scripts/operation/generate-btc-key-local.sh btc

.PHONY: generate-bch-key-local
generate-bch-key-local:
	./scripts/operation/generate-btc-key-local.sh bch

.PHONY: generate-eth-key-local
generate-eth-key-local:
	./scripts/operation/generate-eth-key-local.sh eth

###############################################################################
# payment request
###############################################################################
.PHONY: reset-payment-request
reset-payment-request:
	mysql -h 127.0.0.1 -u root -p${MYSQL_ROOT_PASSWORD} -P 3307 < ./docker/mysql/watch/init.d/payment_request.sql

.PHONY: reset-payment-request-docker
reset-payment-request-docker:
	docker-compose exec btc-watch-db mysql -u root -proot  -e "$(cat ./docker/mysql/watch/init.d/payment_request.sql)"

###############################################################################
# Operation
###############################################################################
include ./Makefile_watch_op.mk
include ./Makefile_keygen_op.mk
include ./Makefile_sign_op.mk


###############################################################################
# wallet
###############################################################################
# run only once, even if wallet.dat is removed
.PHONY: create-wallets
create-wallets:
	bitcoin-cli createwallet watch
	bitcoin-cli createwallet keygen
	bitcoin-cli createwallet sign1
	bitcoin-cli createwallet sign2
	bitcoin-cli createwallet sign3
	bitcoin-cli createwallet sign4
	bitcoin-cli createwallet sign5
	bitcoin-cli listwallets

# list loaded wallets (listed wallet is not needed to load, these wallet can be unloaded
.PHONY: list-wallets
list-wallets:
	bitcoin-cli listwallets

# required after bitcoind restarted
# however, it takes much time in bitcoin ABC (for BCH)
#  so, removing any wallet.dat from server before restarting in BCH, then create wallet again.
.PHONY: load-wallet
load-wallets:
	bitcoin-cli loadwallet watch
	bitcoin-cli loadwallet keygen
	bitcoin-cli loadwallet sign1
	bitcoin-cli loadwallet sign2
	bitcoin-cli loadwallet sign3
	bitcoin-cli loadwallet sign4
	bitcoin-cli loadwallet sign5

#.PHONY: unload-wallet
#unload-wallet:
#	bitcoin-cli -rpcwallet=watch unloadwallet
#	bitcoin-cli -rpcwallet=keygen unloadwallet
#	bitcoin-cli -rpcwallet=sign1 unloadwallet

.PHONY: encrypt-wallets
encrypt-wallets:
	keygen api encryptwallet -passphrase test
	sign -wallet sign1 api encryptwallet -passphrase test
	sign2 -wallet sign2 api encryptwallet -passphrase test
	sign3 -wallet sign3 api encryptwallet -passphrase test
	sign4 -wallet sign4 api encryptwallet -passphrase test
	sign5 -wallet sign5 api encryptwallet -passphrase test

#.PHONY: dump-wallet
dump-wallet:
	keygen api walletpassphrase -passphrase test
	keygen api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/keygen.bk
	sign -wallet sign1 api walletpassphrase -passphrase test
	sign -wallet sign1 api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/sign1.bk
	sign2 -wallet sign2 api walletpassphrase -passphrase test
	sign2 -wallet sign2 api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/sign2.bk
	sign3 -wallet sign3 api walletpassphrase -passphrase test
	sign3 -wallet sign3 api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/sign3.bk
	sign4 -wallet sign4 api walletpassphrase -passphrase test
	sign4 -wallet sign4 api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/sign4.bk
	sign5 -wallet sign4 api walletpassphrase -passphrase test
	sign5 -wallet sign4 api dumpwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/sign5.bk
	#bitcoin-cli -rpcwallet=watch dumpwallet "watch"
	#bitcoin-cli -rpcwallet=keygen dumpwallet "keygen"
	#bitcoin-cli -rpcwallet=sign dumpwallet "sign"

.PHONY: wallet-info
wallet-info:
	bitcoin-cli -rpcwallet=watch getwalletinfo
	bitcoin-cli -rpcwallet=keygen getwalletinfo
	bitcoin-cli -rpcwallet=sign1 getwalletinfo

.PHONY: import-wallet
import-wallet:
	keygen api walletpassphrase -passphrase test
	keygen api importwallet -file ${GOPATH}/src/github.com/hiromaily/go-crypto-wallet/data/dump/keygen.bk

###############################################################################
# Utility
###############################################################################
.PHONY: rm-local-wallet-dat
rm-local-wallet-dat:
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/wallet.dat
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/watch
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/keygen
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/sign1
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/sign2
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/sign3
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/sign4
	rm -rf ~/Library/Application\ Support/Bitcoin/testnet3/wallets/sign5

.PHONY: rm-docker-wallet-dat
rm-docker-wallet-dat:
	# BTC
	rm -rf ./docker/btc/data/testnet3/wallets/wallet.data
	# BCH
	rm -rf ./docker/bch/data/testnet3/wallets/wallet.dat
	rm -rf ./docker/bch/data/testnet3/wallets/watch
	rm -rf ./docker/bch/data/testnet3/wallets/keygen
	rm -rf ./docker/bch/data/testnet3/wallets/sign1
	rm -rf ./docker/bch/data/testnet3/wallets/sign2
	rm -rf ./docker/bch/data/testnet3/wallets/sign3
	rm -rf ./docker/bch/data/testnet3/wallets/sign4
	rm -rf ./docker/bch/data/testnet3/wallets/sign5


.PHONY: rm-files
rm-files:
	rm -rf ./data/btc/address/*.csv
	rm -rf ./data/btc/pubkey/*.csv
	rm -rf ./data/btc/tx/deposit/*
	rm -rf ./data/btc/tx/payment/*
	rm -rf ./data/btc/tx/transfer/*
	touch ./data/btc/tx/deposit/.gitkeep
	touch ./data/btc/tx/payment/.gitkeep
	touch ./data/btc/tx/transfer/.gitkeep

.PHONY: clean
clean: rm-db-volumes rm-local-wallet-dat

#after that, run `make up-docker-db`

# bitcoin-cli
#bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo

