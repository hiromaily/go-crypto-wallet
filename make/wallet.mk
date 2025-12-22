###############################################################################
# Wallet Management Targets
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
# payment request
###############################################################################
.PHONY: reset-payment-request
reset-payment-request:
	mysql -h 127.0.0.1 -u root -p${MYSQL_ROOT_PASSWORD} -P 3307 < ./docker/mysql/watch/init.d/payment_request.sql

.PHONY: reset-payment-request-docker
reset-payment-request-docker:
	docker compose exec watch-db mysql -u root -proot  -e "$(cat ./docker/mysql/watch/init.d/payment_request.sql)"
