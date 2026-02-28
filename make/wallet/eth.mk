###############################################################################
# Ethereum Targets
###############################################################################

# E2E defaults (can be overridden on the command line, e.g. NODE_TYPE=geth DB=mysql)
NODE_TYPE ?= anvil # anvil, geth
DB ?= sqlite

# ETH Variables
GETH_HTTP_PORT=8546
BEACON_HTTP_PORT=9596
GETH_VERSION=v1.17.0
LODESTAR_VERSION=v1.40.0
#ETH_CHAIN_ID=11155111 # used in docker-compose.eth.yml.
TARGET_NETWORK=sepolia
# https://eth-clients.github.io/checkpoint-sync-endpoints/
CHECKPOINT_SYNC_URL=https://beaconstate-${TARGET_NETWORK}.chainsafe.io

###############################################################################
# Docker Compose Targets
###############################################################################

# Common environment variables for ETH docker compose
ETH_COMPOSE_ENV := \
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL)

# run ethereum node server (Geth + Lodestar)
.PHONY: up-docker-geth
up-docker-geth:
	$(ETH_COMPOSE_ENV) docker compose -f compose.eth.yaml --profile geth up

.PHONY: up-docker-geth-d
up-docker-geth-d:
	$(ETH_COMPOSE_ENV) docker compose -f compose.eth.yaml --profile geth up -d

.PHONY: stop-docker-geth
stop-docker-geth:
	$(ETH_COMPOSE_ENV) docker compose -f compose.eth.yaml --profile geth stop

# run ethereum lodestar (requires geth profile)
.PHONY: up-docker-lodestar
up-docker-lodestar:
	$(ETH_COMPOSE_ENV) docker compose -f compose.eth.yaml --profile geth up lodestar

###############################################################################
# Geth specific
###############################################################################
.PHONY:geth-help
geth-help:
	docker run --rm ethereum/client-go:$(GETH_VERSION) --help

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
# Anvil (Foundry Local Node)
###############################################################################

# Anvil RPC port (host-side; container always listens on 8545)
ANVIL_RPC_PORT ?= 8546

# forge binary — use Foundry toolchain installation path
FORGE ?= $(HOME)/.foundry/bin/forge

# Private key for Anvil account 0, derived from the project mnemonic.
# Override via environment variable for a different deployer account.
ANVIL_PRIVATE_KEY ?= 0xd3f1327df33c2a67ac03877f37410f33ed3a37085a84c739a7f5a137e5152bf2

.PHONY: up-docker-anvil
up-docker-anvil:
	docker compose -f compose.eth.yaml --profile anvil up -d anvil

.PHONY: stop-docker-anvil
stop-docker-anvil:
	docker compose -f compose.eth.yaml --profile anvil stop anvil

## deploy-hyt: Deploy HYT ERC-20 contract to the local Anvil node.
##   Requires: Anvil running (make up-docker-anvil)
##   Outputs:  Contract address printed and saved to broadcast/ artifacts.
##   Override: ANVIL_RPC_PORT=8546 ANVIL_PRIVATE_KEY=0x... make deploy-hyt
.PHONY: deploy-hyt
deploy-hyt:
	cd apps/eth-contracts && \
		PRIVATE_KEY=$(ANVIL_PRIVATE_KEY) \
		$(FORGE) script script/DeployHYT.s.sol \
			--rpc-url http://localhost:$(ANVIL_RPC_PORT) \
			--broadcast
	@echo ""
	@echo "Contract address:"
	@grep -oE '"contractAddress"\s*:\s*"0x[a-fA-F0-9]+"' \
		apps/eth-contracts/broadcast/DeployHYT.s.sol/31337/run-latest.json \
		| head -1 | grep -oE '0x[a-fA-F0-9]+'

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
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1.sh --reset

.PHONY: eth-e2e-p1
eth-e2e-p1: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1.sh

.PHONY: eth-e2e-p1-verbose
eth-e2e-p1-verbose: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1.sh --verbose

.PHONY: eth-e2e-p1-ci
eth-e2e-p1-ci: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1.sh --non-interactive

.PHONY: eth-e2e-p1-cleanup
eth-e2e-p1-cleanup:
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p1.sh --cleanup

###############################################################################
# E2E Tests
#
# Pattern 2: ERC-20 HYT token transfer (EIP-1559)
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Requires the HYT contract to be deployed first:
#   make deploy-hyt
#
# Usage:
#   make eth-e2e-p2-reset             # Run with Anvil + SQLite (full reset)
#   make eth-e2e-p2                   # Run without reset
#   make eth-e2e-p2-reset HYT_CONTRACT_ADDRESS=0x...  # Override contract address
###############################################################################

.PHONY: eth-e2e-p2-reset
eth-e2e-p2-reset: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p2.sh --reset

.PHONY: eth-e2e-p2
eth-e2e-p2: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p2.sh

.PHONY: eth-e2e-p2-verbose
eth-e2e-p2-verbose: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p2.sh --verbose

.PHONY: eth-e2e-p2-ci
eth-e2e-p2-ci: build-all
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p2.sh --non-interactive

.PHONY: eth-e2e-p2-cleanup
eth-e2e-p2-cleanup:
	NODE_TYPE=$(NODE_TYPE) DB_TYPE=$(DB) ./scripts/operation/eth/e2e/e2e-p2.sh --cleanup
