###############################################################################
# BTC Related Targets
###############################################################################

# E2E Testing
include make/btc_e2e.mk


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
# Docker Compose Targets
###############################################################################

# run all bitcoin core servers
.PHONY: up-docker-btc
up-docker-btc:
	docker compose -f compose.btc.yaml up btc-watch btc-keygen btc-sign1 btc-sign2

# run bitcoin core servers in detached mode
.PHONY: up-docker-btc-d
up-docker-btc-d:
	docker compose -f compose.btc.yaml up -d btc-watch btc-keygen btc-sign1 btc-sign2

# stop bitcoin core servers
.PHONY: down-docker-btc
down-docker-btc:
	docker compose -f compose.btc.yaml down

# stop bitcoin core servers and remove volumes
.PHONY: down-docker-btc-v
down-docker-btc-v:
	docker compose -f compose.btc.yaml down -v

###############################################################################
# auto key generator
###############################################################################
.PHONY: generate-btc-key-local
generate-btc-key-local:
	./scripts/operation/generate-btc-key.sh btc false 5

.PHONY: generate-bch-key-local
generate-bch-key-local:
	./scripts/operation/generate-btc-key.sh bch false 5

