###############################################################################
# Ethereum Targets
###############################################################################

# E2E defaults (can be overridden on the command line, e.g. NODE_TYPE=geth DB=mysql)
NODE_TYPE ?= anvil
DB ?= sqlite

###############################################################################
# Docker Compose Targets
###############################################################################

# run ethereum node server (Geth + Lodestar)
.PHONY: up-docker-eth
up-docker-eth:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml --profile geth up

.PHONY: up-docker-eth-d
up-docker-eth-d:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml --profile geth up -d

.PHONY: stop-docker-eth
stop-docker-eth:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml --profile geth stop

# run ethereum lodestar (requires geth profile)
.PHONY: up-docker-lodestar
up-docker-lodestar:
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml --profile geth up lodestar

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

###############################################################################
# Anvil (Foundry Local Node)
###############################################################################
.PHONY: up-docker-anvil
up-docker-anvil:
	docker compose -f compose.eth.yaml --profile anvil up -d anvil

.PHONY: stop-docker-anvil
stop-docker-anvil:
	docker compose -f compose.eth.yaml --profile anvil stop anvil

###############################################################################
# E2E Tests
#
# Pattern 1: Single-sig EIP-1559
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Usage:
#   make eth-e2e-p1-reset             # Run with Anvil + SQLite
#   make eth-e2e-p1-reset NODE_TYPE=geth DB=mysql  # Run with Geth + MySQL
###############################################################################

.PHONY: eth-e2e-p1-reset
eth-e2e-p1-reset: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh --reset

.PHONY: eth-e2e-p1
eth-e2e-p1: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh

.PHONY: eth-e2e-p1-verbose
eth-e2e-p1-verbose: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh --verbose

.PHONY: eth-e2e-p1-ci
eth-e2e-p1-ci: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh --non-interactive

.PHONY: eth-e2e-p1-cleanup
eth-e2e-p1-cleanup:
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh --cleanup
