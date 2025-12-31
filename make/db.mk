###############################################################################
# Database Targets
###############################################################################

###############################################################################
# include other db related makefiles
###############################################################################

# SQLC-related targets are defined in make/db_sqlc.mk
include make/db_sqlc.mk

# Atlas-related targets are defined in make/db_atlas.mk
include make/db_atlas.mk

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database and migration services
# remove all containers and volumes before starting
.PHONY: up-docker-db
up-docker-db:
	docker compose down -v || true
	docker compose up
