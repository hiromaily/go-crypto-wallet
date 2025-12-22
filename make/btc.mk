###############################################################################
# Bitcoin Core Targets
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
# auto key generator
###############################################################################
.PHONY: generate-btc-key-local
generate-btc-key-local:
	./scripts/operation/generate-btc-key.sh btc false 5

.PHONY: generate-bch-key-local
generate-bch-key-local:
	./scripts/operation/generate-btc-key.sh bch false 5
