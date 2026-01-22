###############################################################################
# Wallet Utility Targets
###############################################################################

# remove local wallet data
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

# remove docker wallet data
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


# bitcoin-cli
# - using arguments
# $ bitcoin-cli -rpcuser=xyz -rpcpassword=xyz getnetworkinfo
# - check sync information
# $ bitcoin-cli getblockchaininfo
