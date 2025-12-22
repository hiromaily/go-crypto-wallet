###############################################################################
# Docker and Compose Targets
###############################################################################

# run bitcoin core server
.PHONY: up-docker-btc
up-docker-btc:
	docker compose up btc-watch btc-keygen btc-sign

# run bitcoin cash core server
.PHONY: up-docker-bch
up-docker-bch:
	docker compose -f compose.bch.yaml up bch-watch

# run ethereum node server
.PHONY: up-docker-eth
up-docker-eth:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml up geth lodestar

.PHONY: up-docker-eth-d
up-docker-eth-d:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml up -d geth lodestar

.PHONY: stop-docker-eth
stop-docker-eth:
	GETH_VERSION=$(GETH_VERSION) GETH_HTTP_PORT=$(GETH_HTTP_PORT) \
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml stop

# run ethereum lodestar
.PHONY: up-docker-lodestar
up-docker-lodestar:
	LODESTAR_VERSION=$(LODESTAR_VERSION) BEACON_HTTP_PORT=$(BEACON_HTTP_PORT) TARGET_NETWORK=$(TARGET_NETWORK) \
	CHECKPOINT_SYNC_URL=$(CHECKPOINT_SYNC_URL) \
	docker compose -f compose.eth.yaml up lodestar

# run ripple node server
.PHONY: up-docker-xrp
up-docker-xrp:
	docker compose -f compose.xrp.yaml up xrp-node

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
