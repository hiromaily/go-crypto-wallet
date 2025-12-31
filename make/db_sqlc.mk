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

###############################################################################
# Schema Extraction Targets (for sqlc)
###############################################################################

# Extract sqlc schema from watch dump file
# This extracts CREATE TABLE statements from MySQL dump and formats them for sqlc:
# - Excludes atlas_schema_revisions table (not needed for sqlc)
# - Extracts CREATE TABLE statements from dump file
# - Removes backticks (sqlc prefers without them)
# - Removes MySQL conditional comments
# - Adds blank line after each CREATE TABLE statement
.PHONY: extract-sqlc-schema-watch
extract-sqlc-schema-watch: dump-schema-watch
	@scripts/db/extract-sqlc-schema.sh \
		watch \
		data/dump/sql/dump_watch.sql \
		tools/sqlc/schemas/extracted_watch.sql

# Extract sqlc schema from keygen dump file
# This extracts CREATE TABLE statements from MySQL dump and formats them for sqlc:
# - Excludes atlas_schema_revisions table (not needed for sqlc)
# - Extracts CREATE TABLE statements from dump file
# - Removes backticks (sqlc prefers without them)
# - Removes MySQL conditional comments
# - Adds blank line after each CREATE TABLE statement
.PHONY: extract-sqlc-schema-keygen
extract-sqlc-schema-keygen: dump-schema-keygen
	@scripts/db/extract-sqlc-schema.sh \
		keygen \
		data/dump/sql/dump_keygen.sql \
		tools/sqlc/schemas/extracted_keygen.sql

# Extract sqlc schema from sign dump file
# This extracts CREATE TABLE statements from MySQL dump and formats them for sqlc:
# - Excludes atlas_schema_revisions table (not needed for sqlc)
# - Excludes seed table (exists in keygen schema, avoid duplication)
# - Extracts CREATE TABLE statements from dump file
# - Removes backticks (sqlc prefers without them)
# - Removes MySQL conditional comments
# - Adds blank line after each CREATE TABLE statement
.PHONY: extract-sqlc-schema-sign
extract-sqlc-schema-sign: dump-schema-sign
	@scripts/db/extract-sqlc-schema.sh \
		sign \
		data/dump/sql/dump_sign.sql \
		tools/sqlc/schemas/extracted_sign.sql

# Extract sqlc schemas from all dump files
.PHONY: extract-sqlc-schema-all
extract-sqlc-schema-all: extract-sqlc-schema-watch extract-sqlc-schema-keygen extract-sqlc-schema-sign

###############################################################################
# Generate sqlc code
###############################################################################
# need to run after database is reset and migrations are applied
# This target waits for migration services to complete before extracting schemas
.PHONY: regenerate-sqlc
regenerate-sqlc: wait-migrations
	make extract-sqlc-schema-all
	make sqlc

# Wait for all migration services to complete
# Migration services exit after applying migrations (restart: "no" in compose.yaml)
# This checks that all migration containers have exited successfully (Exited (0))
.PHONY: wait-migrations
wait-migrations:
	@echo "Waiting for migration services to complete..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		status_output=$$(docker compose ps -a wallet-db-migrate-watch wallet-db-migrate-keygen wallet-db-migrate-sign 2>/dev/null || echo ""); \
		if [ -n "$$status_output" ]; then \
			if echo "$$status_output" | grep -q "wallet-db-migrate-watch.*Exited (0)" && \
			   echo "$$status_output" | grep -q "wallet-db-migrate-keygen.*Exited (0)" && \
			   echo "$$status_output" | grep -q "wallet-db-migrate-sign.*Exited (0)"; then \
				echo "✓ All migration services completed"; \
				exit 0; \
			fi; \
		fi; \
		sleep 1; \
		timeout=$$((timeout - 1)); \
	done; \
	echo "⚠ Warning: Migration services may not have completed within timeout."; \
	echo "Current status:"; \
	docker compose ps -a wallet-db-migrate-watch wallet-db-migrate-keygen wallet-db-migrate-sign 2>/dev/null || true; \
	exit 1
