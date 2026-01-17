###############################################################################
# sqlc Format, Validation, and Lint Targets
# Validates both MySQL and SQLite configurations
###############################################################################

# Compile SQL queries and schemas to check for syntax and type errors
# `sqlc compile` - Statically check SQL for syntax and type errors
# Runs for both MySQL (default) and SQLite configurations
.PHONY: sqlc-compile
sqlc-compile:
	@echo "Compiling MySQL SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile
	@echo "✓ MySQL SQL compilation successful"
	@echo "Compiling SQLite SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL compilation successful"

# Vet SQL queries for potential issues
# `sqlc vet` - Examines queries for potential problems
# Runs for both MySQL (default) and SQLite configurations
.PHONY: sqlc-vet
sqlc-vet:
	@echo "Vetting MySQL SQL queries..."
	@cd tools/sqlc && sqlc vet
	@echo "✓ MySQL SQL queries passed vetting"
	@echo "Vetting SQLite SQL queries..."
	@cd tools/sqlc && sqlc vet -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL queries passed vetting"

# Verify schema, queries, and configuration
# `sqlc verify` - Verify schema, queries, and configuration for this project
# Note: This command requires sqlc Cloud connection. If you don't use sqlc Cloud,
# you can skip this command and use sqlc-compile and sqlc-vet instead.
.PHONY: sqlc-verify
sqlc-verify:
	@echo "Verifying sqlc configuration, schemas, and queries..."
	@echo "Note: This command requires sqlc Cloud. Use sqlc-validate for local validation."
	@cd tools/sqlc && sqlc verify --no-remote || (echo "⚠️  sqlc verify requires sqlc Cloud connection. Use sqlc-validate for local validation." && exit 0)

# Validate SQL queries and schemas (combines compile and vet)
# This target runs compile and vet checks for local validation
# Validates both MySQL and SQLite configurations
.PHONY: sqlc-validate
sqlc-validate: sqlc-compile sqlc-vet
	@echo "✓ All sqlc validation checks passed (MySQL + SQLite)"

# Lint SQL queries (alias for vet)
# This target runs vet to check queries for potential issues
.PHONY: sqlc-lint
sqlc-lint: sqlc-vet

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
