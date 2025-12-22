###############################################################################
# Utility Targets
###############################################################################

.PHONY: timestamp
timestamp:
	@echo $(timestamp)

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
