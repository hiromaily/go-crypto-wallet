###############################################################################
# XRP/Ripple Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run ripple node server (foreground)
.PHONY: up-docker-xrp
up-docker-xrp:
	docker compose -f compose.xrp.yaml up xrp-node

# run ripple node server (background)
.PHONY: up-docker-xrp-d
up-docker-xrp-d:
	docker compose -f compose.xrp.yaml up -d xrp-node

# run all XRP services (xrp-node + xrpl-grpc-server)
.PHONY: up-docker-xrp-all
up-docker-xrp-all:
	docker compose -f compose.xrp.yaml up -d

# stop XRP services
.PHONY: down-docker-xrp
down-docker-xrp:
	docker compose -f compose.xrp.yaml down

# view xrp-node logs
.PHONY: logs-xrp-node
logs-xrp-node:
	docker compose -f compose.xrp.yaml logs -f xrp-node

# check xrp-node server info
.PHONY: xrp-node-info
xrp-node-info:
	docker compose -f compose.xrp.yaml exec xrp-node server_info
