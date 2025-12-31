###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database
.PHONY: reset-docker-db
reset-docker-db:
	docker compose down -v
	docker compose up


###############################################################################
# DB Related Targets
###############################################################################
# Atlas-related targets
include make/db_atlas.mk
# sqlc-related targets
include make/db_sqlc.mk
