modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)
GOLANGCI_VERSION=v1.46.2
#PROTOC_BIN=protoc
PROTOC_BIN=buf protoc

###############################################################################
# Initial Settings
###############################################################################
.PHONY: check-ver
check-ver:
	#echo $(modVer)
	#echo $(currentVer)
	@if [ ${currentVer} -lt ${modVer} ]; then\
		echo go version ${modVer}++ is required but your go version is ${currentVer};\
	fi

.PHONY: install-mac-tools
install-mac-tools:
	brew install jq mkcert go-task/tap/go-task

.PHONY: install-protobuf
install-protobuf:
	brew install protobuf prototool
	brew tap bufbuild/buf                                                                                                                                                                      (git)-[master]
	brew install buf

.PHONY: install-ssl
install-ssl:
	mkcert -install
	mkcert localhost 127.0.0.1

.PHONY: install-tools
install-tools:
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/icholy/gomajor@latest

.PHONY: install-proto-plugin
install-proto-plugin:
	# refer to https://developers.google.com/protocol-buffers/docs/reference/go-generated
	# provides a protoc-gen-go binary which protoc uses when invoked with the --go_out command-line flag.
	# The --go_out flag
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	# The go-grpc_out flag
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# gogo protobuf
	go install github.com/gogo/protobuf/protoc-gen-gogo@latest
	#go get github.com/gogo/protobuf/proto
	#go get github.com/gogo/protobuf/jsonpb
	#go get github.com/gogo/protobuf/gogoproto

.PHONY: goget
goget:
	go get ./...

###############################################################################
# Code Generator
###############################################################################
# sqlboiler
#------------------------------------------------------------------------------
# To generate all schema, modify `docker/mysql/watch/init.d/watch.sql` according to comments
# Then recreate database
# ```
# $ docker compose rm -f -s watch-db
# $ docker volume rm -f go-crypto-wallet_watch-db
# $ docker compose up watch-db
# ```
# Make sure `watch-db` includes tables of keygen-db/sign-db
# Then, run `make sqlboiler`
# Make sure `make build` works
# Revert `docker/mysql/watch/init.d/watch.sql`
#------------------------------------------------------------------------------
.PHONY: sqlboiler
sqlboiler:
	sqlboiler --wipe mysql

.PHONY: sqlboiler-with-template
sqlboiler-with-template:
	sqlboiler --wipe mysql \
	--templates \
	${GOPATH}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.8.6/templates/main, \
	${GOPATH}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.8.6/templates/test, \
	templates

.PHONY: generate-abi
generate-abi:
	abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./pkg/contract/token-abi.go

###############################################################################
# Protocol Buffer
#------------------------------------------------------------------------------
# run `make install-proto-plugin` in advance
###############################################################################
.PHONY: get-third-proto
get-third-proto:
	./scripts/get_third_proto.sh

.PHONY: lint-proto
lint-proto:
	prototool lint

.PHONY:protoc-go
protoc-go: clean-pb
	$(PROTOC_BIN) \
	--go_out=./pkg/wallet/api/xrpgrp/xrp/ --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/wallet/api/xrpgrp/xrp/ --go-grpc_opt=paths=source_relative  \
	--proto_path=./data/proto/rippleapi \
	--proto_path=./data/proto/third_party \
	data/proto/**/*.proto

.PHONY: clean-pb
clean-pb:
	rm -rf pkg/wallet/api/xrpgrp/xrp/*.pb.go

###############################################################################
# Linter
###############################################################################
.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run

# Bug: format doesn't work on files which has tags
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: check-upgrade
check-upgrade:
	gomajor list

###############################################################################
# From inside docker container
###############################################################################
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/watch ./cmd/watch/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.authName=auth1" -o /go/bin/sign ./cmd/sign/main.go

###############################################################################
# Build on local
# - authName on sign works as account name
###############################################################################
.PHONY: tidy
tidy:
	go mod tidy -compat=1.17

.PHONY: build
build: tidy
	go build -v -o ${GOPATH}/bin/watch ./cmd/watch/
	go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
	go build -ldflags "-X main.authName=auth3" -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
	go build -ldflags "-X main.authName=auth4" -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
	go build -ldflags "-X main.authName=auth5" -v -o ${GOPATH}/bin/sign5 ./cmd/sign/

.PHONY: build-watch
build-watch:
	go build -v -o ${GOPATH}/bin/watch ./cmd/watch/

.PHONY: build-keygen
build-keygen:
	go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/

.PHONY: build-sign
build-sign:
	go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
	go build -ldflags "-X main.authName=auth3" -v -o ${GOPATH}/bin/sign3 ./cmd/sign/
	go build -ldflags "-X main.authName=auth4" -v -o ${GOPATH}/bin/sign4 ./cmd/sign/
	go build -ldflags "-X main.authName=auth5" -v -o ${GOPATH}/bin/sign5 ./cmd/sign/

run:
	go run ./cmd/watch/ -conf ./data/config/watch.toml

###############################################################################
# Test on local
###############################################################################
.PHONY: gotest
gotest:
	go test -v ./...

.PHONY: gotest-addr
gotest-addr:
	go test -tags=integration -v -run pkg/wallets/api/btc/...
	go test -tags=integration -v -run GetAddressInfo pkg/wallets/api/btc/...
	go test -v pkg/wallets/api/btc/... -run GetAddressInfo

.PHONY: gotest-integration
gotest-integration:
	go test -v -tags=integration ./...


###############################################################################
# Docker and compose
###############################################################################

# run bitcoin core server
.PHONY: up-docker-btc
up-docker-btc:
	docker compose up btc-watch btc-keygen btc-sign

# run bitcoin cash core server
.PHONY: up-docker-bch
up-docker-bch:
	docker compose -f docker-compose.bch.yml up bch-watch

# run ethereum node server
.PHONY: up-docker-eth
up-docker-eth:
	docker compose -f docker-compose.eth.yml up eth-node

# run ripple node server
.PHONY: up-docker-xrp
up-docker-xrp:
	docker compose -f docker-compose.xrp.yml up xrp-node

# run all databases
.PHONY: up-docker-db
up-docker-db:
	docker compose up watch-db keygen-db sign-db

# run logging middleware
# logging and monitoring
#.PHONY: up-docker-logger
#up-docker-logger:
#	docker compose up fluentd elasticsearch grafana

# remove database volumes
.PHONY: rm-db-volumes
rm-db-volumes:
	#docker rm -f $(docker ps -a --format "{{.Names}}")
	#docker volume rm -f $(docker volume ls --format "{{.Name}}")
	#docker compose down -v
	#docker compose down
	docker volume rm -f watch-db
	docker volume rm -f keygen-db
	docker volume rm -f sign-db

###############################################################################
# Bitcoin core on local
###############################################################################
.PHONY: bitcoin-run
bitcoin-run:
	bitcoind -daemon

.PHONY: bitcoin-stop
bitcoin-stop:
	bitcoin-cli stop

# MacOS only
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
	./scripts/operation/generate-btc-key.sh btc false 5

.PHONY: generate-bch-key-local
generate-bch-key-local:
	./scripts/operation/generate-btc-key.sh bch false 5

.PHONY: generate-eth-key-local
generate-eth-key-local:
	./scripts/operation/generate-eth-key.sh eth

###############################################################################
# payment request
###############################################################################
.PHONY: reset-payment-request
reset-payment-request:
	mysql -h 127.0.0.1 -u root -p${MYSQL_ROOT_PASSWORD} -P 3307 < ./docker/mysql/watch/init.d/payment_request.sql

.PHONY: reset-payment-request-docker
reset-payment-request-docker:
	docker compose exec watch-db mysql -u root -proot  -e "$(cat ./docker/mysql/watch/init.d/payment_request.sql)"

###############################################################################
# Operation
###############################################################################
include ./Makefile_watch_op.mk
include ./Makefile_keygen_op.mk
include ./Makefile_sign_op.mk


###############################################################################
# wallet
###############################################################################
.PHONY: create-wallets
create-wallets:
	bitcoin-cli-watch createwallet watch
	bitcoin-cli-keygen createwallet keygen
	bitcoin-cli-sign createwallet sign1
	bitcoin-cli-sign createwallet sign2
	bitcoin-cli-sign createwallet sign3
	bitcoin-cli-sign createwallet sign4
	bitcoin-cli-sign createwallet sign5
	bitcoin-cli-sign listwallets

# run only once, even if wallet.dat is removed
.PHONY: create-wallets-one-bitcoind
create-wallets-one-bitcoind:
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
	rm -rf ./docker/btc/data/testnet3/wallets/watch
	rm -rf ./docker/btc/data/testnet3/wallets/keygen
	rm -rf ./docker/btc/data/testnet3/wallets/sign1
	rm -rf ./docker/btc/data/testnet3/wallets/sign2
	rm -rf ./docker/btc/data/testnet3/wallets/sign3
	rm -rf ./docker/btc/data/testnet3/wallets/sign4
	rm -rf ./docker/btc/data/testnet3/wallets/sign5
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
# - using arguments
# $ bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo
# - check sync information
# $ bitcoin-cli getblockchaininfo