###############################################################################
# Ethereum Targets
###############################################################################

###############################################################################
# Geth specific
###############################################################################
.PHONY:geth-help
geth-help:
	docker run --rm ethereum/client-go:$(GETH_VERSION) --help

# geth image based on ethereum/client-go:v1.10.26 with curl commnad
.PHONY:build-geth-image
build-geth-image:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	docker compose -f compose.eth.yaml build --no-cache geth

.PHONY:import-geth-data
import-geth-data:
	docker run -v $(CURDIR)/docker/nodes/eth/backup:/backup -v $(CURDIR)/docker/nodes/eth/$(TARGET_NETWORK):/data ethereum/client-go:$(GETH_VERSION) import --datadir=/data /backup/exported-file

# run after geth stopped
.PHONY:export-geth-data
export-geth-data:
	docker run -v $(CURDIR)/docker/nodes/eth/backup:/backup -v $(CURDIR)/docker/nodes/eth/$(TARGET_NETWORK):/data ethereum/client-go:$(GETH_VERSION) export --datadir=/data /backup/exported-file-$(timestamp)

.PHONY:check-execution-block
check-execution-block:
	curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}' localhost:$(GETH_HTTP_PORT)

.PHONY:check-execution-syncing
check-execution-syncing:
	curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' localhost:$(GETH_HTTP_PORT)

###############################################################################
# auto key generator
###############################################################################
.PHONY: generate-eth-key-local
generate-eth-key-local:
	./scripts/operation/generate-eth-key.sh eth

###############################################################################
# Grafana
###############################################################################
# http://localhost:3000
