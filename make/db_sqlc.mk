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

# Clean old sqlc schema files before extracting new ones
# This removes manually created schema files to prevent duplication errors
.PHONY: clean-sqlc-schemas
clean-sqlc-schemas:
	@echo "Cleaning old sqlc schema files..."
	@rm -f tools/sqlc/schemas/*.sql
	@echo "Done."

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
		tools/sqlc/schemas/01_watch.sql

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
		tools/sqlc/schemas/02_keygen.sql

# Extract sqlc schema from sign dump file
# This extracts CREATE TABLE statements from MySQL dump and formats them for sqlc:
# - Excludes atlas_schema_revisions table (not needed for sqlc)
# - Excludes seed table (exists in keygen schema, avoid duplication)
# - Excludes musig2_nonces table (exists in keygen schema, avoid duplication)
# - Extracts CREATE TABLE statements from dump file
# - Removes backticks (sqlc prefers without them)
# - Removes MySQL conditional comments
# - Adds blank line after each CREATE TABLE statement
.PHONY: extract-sqlc-schema-sign
extract-sqlc-schema-sign: dump-schema-sign
	@scripts/db/extract-sqlc-schema.sh \
		sign \
		data/dump/sql/dump_sign.sql \
		tools/sqlc/schemas/03_sign.sql

# Extract sqlc schemas from all dump files
# This cleans old schema files first, then extracts new ones from database dumps
.PHONY: extract-sqlc-schema-all
extract-sqlc-schema-all: clean-sqlc-schemas
	@$(MAKE) extract-sqlc-schema-watch
	@$(MAKE) extract-sqlc-schema-keygen
	@$(MAKE) extract-sqlc-schema-sign

###############################################################################
# Generate sqlc code
###############################################################################
# need to run after database is reset and migrations are applied
# This target waits for migration services to complete before extracting schemas
.PHONY: regenerate-sqlc-from-current-db
regenerate-sqlc-from-current-db:
	make extract-sqlc-schema-all
	make sqlc
