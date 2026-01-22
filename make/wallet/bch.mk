###############################################################################
# Bitcoin Cash Node Targets
###############################################################################

# E2E Testing
include make/wallet/bch_e2e.mk

###############################################################################
# Docker Compose Targets
###############################################################################

# run all bitcoin cash node servers
.PHONY: up-docker-bch
up-docker-bch:
	docker compose -f compose.bch.yaml up bch-watch bch-keygen bch-sign1 bch-sign2

# run bitcoin cash node servers in detached mode
.PHONY: up-docker-bch-d
up-docker-bch-d:
	docker compose -f compose.bch.yaml up -d bch-watch bch-keygen bch-sign1 bch-sign2

# stop bitcoin cash node servers
.PHONY: down-docker-bch
down-docker-bch:
	docker compose -f compose.bch.yaml down

# stop bitcoin cash node servers and remove volumes
.PHONY: down-docker-bch-v
down-docker-bch-v:
	docker compose -f compose.bch.yaml down -v


