###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run all databases
.PHONY: up-docker-db
up-docker-db:
	docker compose up watch-db keygen-db sign-db

# remove database volumes
.PHONY: rm-db-volumes
rm-db-volumes:
	docker volume rm -f watch-db
	docker volume rm -f keygen-db
	docker volume rm -f sign-db

###############################################################################
# Schema Export Targets
###############################################################################

# Export schema from watch-db container
.PHONY: dump-schema-watch
dump-schema-watch:
	mkdir -p data/dump/sql
	docker exec watch-db mysqldump -u root -proot --no-data --skip-triggers watch > data/dump/sql/dump_watch.sql

# Export schema from keygen-db container
.PHONY: dump-schema-keygen
dump-schema-keygen:
	mkdir -p data/dump/sql
	docker exec keygen-db mysqldump -u root -proot --no-data --skip-triggers keygen > data/dump/sql/dump_keygen.sql

# Export schema from sign-db container
.PHONY: dump-schema-sign
dump-schema-sign:
	mkdir -p data/dump/sql
	docker exec sign-db mysqldump -u root -proot --no-data --skip-triggers sign > data/dump/sql/dump_sign.sql

# Export all schemas from all database containers
.PHONY: dump-schema-all
dump-schema-all: dump-schema-watch dump-schema-keygen dump-schema-sign
