###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database
.PHONY: up-docker-db
up-docker-db:
	docker compose up wallet-db

# remove database volume
.PHONY: rm-db-volumes
rm-db-volumes:
	docker volume rm -f go-crypto-wallet_wallet-db

###############################################################################
# Schema Export Targets
###############################################################################

# Export watch schema from wallet-db container
.PHONY: dump-schema-watch
dump-schema-watch:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers watch > data/dump/sql/dump_watch.sql

# Export keygen schema from wallet-db container
.PHONY: dump-schema-keygen
dump-schema-keygen:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers keygen > data/dump/sql/dump_keygen.sql

# Export sign schema from wallet-db container
.PHONY: dump-schema-sign
dump-schema-sign:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers sign > data/dump/sql/dump_sign.sql

# Export all schemas from wallet-db container
.PHONY: dump-schema-all
dump-schema-all: dump-schema-watch dump-schema-keygen dump-schema-sign
