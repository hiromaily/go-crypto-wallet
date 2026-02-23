###############################################################################
# sqlc Format, Validation, and Lint Targets
# Validates MySQL, SQLite, and PostgreSQL configurations
###############################################################################

# Compile SQL queries and schemas to check for syntax and type errors
# `sqlc compile` - Statically check SQL for syntax and type errors
# Runs for PostgreSQL (default), MySQL configurations (SQLite skipped)
# Note: **SQLite config references `./queries/mysql/*.sql` MySQL queries.**
# So, sqlite must be ignored
.PHONY: sqlc-compile
sqlc-compile:
	@for dialect in postgres mysql; do \
		echo "Compiling $$dialect SQL queries and schemas..."; \
		cd tools/sqlc && sqlc compile -f sqlc_$${dialect}.yml && cd ../..; \
		echo "✓ $$dialect SQL compilation successful"; \
	done

# Vet SQL queries for potential issues
# `sqlc vet` - Examines queries for potential problems
# Runs for PostgreSQL (default), MySQL configurations (SQLite skipped)
# Note: **SQLite config references `./queries/mysql/*.sql` MySQL queries.**
.PHONY: sqlc-vet
sqlc-vet:
	@for dialect in postgres mysql; do \
		echo "Vetting $$dialect SQL queries..."; \
		cd tools/sqlc && sqlc vet -f sqlc_$${dialect}.yml && cd ../..; \
		echo "✓ $$dialect SQL queries passed vetting"; \
	done

# Validate SQL queries and schemas (combines compile and vet)
# This target runs compile and vet checks for local validation
# Validates MySQL, PostgreSQL configurations (SQLite skipped)
.PHONY: sqlc-validate
sqlc-validate: sqlc-compile sqlc-vet
	@echo "✓ All sqlc validation checks passed (MySQL + SQLite + PostgreSQL)"

# Lint SQL queries (alias for vet)
# This target runs vet to check queries for potential issues
.PHONY: sqlc-lint
sqlc-lint: sqlc-vet

# Verify schema, queries, and configuration
# `sqlc verify` - Verify schema, queries, and configuration for this project
# Note: **This command requires sqlc Cloud connection.** If you don't use sqlc Cloud,
# you can skip this command and use sqlc-compile and sqlc-vet instead.
.PHONY: sqlc-verify
sqlc-verify:
	@echo "Verifying sqlc configuration, schemas, and queries..."
	@echo "Note: This command requires sqlc Cloud. Use sqlc-validate for local validation."
	@for dialect in postgres mysql sqlite; do \
		cd tools/sqlc && sqlc verify -f sqlc_$${dialect}.yml --no-remote && cd ../.. || \
		(echo "⚠️  sqlc verify requires sqlc Cloud connection. Use sqlc-validate for local validation." && cd ../.. && exit 0); \
	done

###############################################################################
# Dump database schema from running Database.
# At this time, dumped files are in ./data/dump/.
# So `make extract-sqlc-schema-all` must be run later
# Usage: make dump-schema-<name> [DB_DIALECT=mysql|postgres]
###############################################################################

# Export watch schema from database container
# Usage: make dump-schema-watch [DB_DIALECT=mysql|postgres]
.PHONY: dump-schema-watch
dump-schema-watch:
	mkdir -p data/dump/sql
ifeq ($(DB_DIALECT),postgres)
	docker exec wallet-$(DB_DIALECT) pg_dump -U postgres --schema-only --no-owner --no-privileges watch > data/dump/sql/dump_watch.$(DB_DIALECT).sql
else
	docker exec wallet-$(DB_DIALECT) mysqldump -u root -proot --no-data --skip-triggers watch > data/dump/sql/dump_watch.$(DB_DIALECT).sql
endif

# Export keygen schema from database container
# Usage: make dump-schema-keygen [DB_DIALECT=mysql|postgres]
.PHONY: dump-schema-keygen
dump-schema-keygen:
	mkdir -p data/dump/sql
ifeq ($(DB_DIALECT),postgres)
	docker exec wallet-$(DB_DIALECT) pg_dump -U postgres --schema-only --no-owner --no-privileges keygen > data/dump/sql/dump_keygen.$(DB_DIALECT).sql
else
	docker exec wallet-$(DB_DIALECT) mysqldump -u root -proot --no-data --skip-triggers keygen > data/dump/sql/dump_keygen.$(DB_DIALECT).sql
endif

# Export sign schema from database container
# Usage: make dump-schema-sign [DB_DIALECT=mysql|postgres]
.PHONY: dump-schema-sign
dump-schema-sign:
	mkdir -p data/dump/sql
ifeq ($(DB_DIALECT),postgres)
	docker exec wallet-$(DB_DIALECT) pg_dump -U postgres --schema-only --no-owner --no-privileges sign > data/dump/sql/dump_sign.$(DB_DIALECT).sql
else
	docker exec wallet-$(DB_DIALECT) mysqldump -u root -proot --no-data --skip-triggers sign > data/dump/sql/dump_sign.$(DB_DIALECT).sql
endif

# Export all schemas from database container
# Usage: make dump-schema-all [DB_DIALECT=mysql|postgres]
.PHONY: dump-schema-all
dump-schema-all: dump-schema-watch dump-schema-keygen dump-schema-sign

# Convenience aliases for dialect-specific dump
.PHONY: dump-schema-all-mysql
dump-schema-all-mysql:
	@$(MAKE) dump-schema-all DB_DIALECT=mysql

###############################################################################
# Schema Extraction for sqlc
# Run after `make dump-schema-all`
# Usage: make extract-sqlc-schema-<name> [DB_DIALECT=postgres|mysql]
###############################################################################

# Clean old sqlc schema files before extracting new ones
# Usage: make clean-sqlc-schemas [DB_DIALECT=postgres|mysql]
.PHONY: clean-sqlc-schemas
clean-sqlc-schemas:
	@echo "Cleaning old sqlc schema files ($(DB_DIALECT))..."
	@rm -f tools/sqlc/schemas/$(DB_DIALECT)/*.sql
	@echo "Done."

# Extract sqlc schema from watch dump file
# Usage: make extract-sqlc-schema-watch [DB_DIALECT=mysql|postgres]
.PHONY: extract-sqlc-schema-watch
extract-sqlc-schema-watch: dump-schema-watch
	@scripts/db/extract-sqlc-schema-$(DB_DIALECT).sh \
		watch \
		data/dump/sql/dump_watch.$(DB_DIALECT).sql \
		tools/sqlc/schemas/$(DB_DIALECT)/01_watch.sql

# Extract sqlc schema from keygen dump file
# Usage: make extract-sqlc-schema-keygen [DB_DIALECT=mysql|postgres]
.PHONY: extract-sqlc-schema-keygen
extract-sqlc-schema-keygen: dump-schema-keygen
	@scripts/db/extract-sqlc-schema-$(DB_DIALECT).sh \
		keygen \
		data/dump/sql/dump_keygen.$(DB_DIALECT).sql \
		tools/sqlc/schemas/$(DB_DIALECT)/02_keygen.sql

# Extract sqlc schema from sign dump file
# Excludes seed and musig2_nonces tables (exist in keygen schema, avoid duplication)
# Usage: make extract-sqlc-schema-sign [DB_DIALECT=mysql|postgres]
.PHONY: extract-sqlc-schema-sign
extract-sqlc-schema-sign: dump-schema-sign
	@scripts/db/extract-sqlc-schema-$(DB_DIALECT).sh \
		sign \
		data/dump/sql/dump_sign.$(DB_DIALECT).sql \
		tools/sqlc/schemas/$(DB_DIALECT)/03_sign.sql

# Extract sqlc schemas from all dump files
# Usage: make extract-sqlc-schema-all [DB_DIALECT=mysql|postgres]
.PHONY: extract-sqlc-schema-all
extract-sqlc-schema-all: clean-sqlc-schemas
	@$(MAKE) extract-sqlc-schema-watch DB_DIALECT=$(DB_DIALECT)
	@$(MAKE) extract-sqlc-schema-keygen DB_DIALECT=$(DB_DIALECT)
	@$(MAKE) extract-sqlc-schema-sign DB_DIALECT=$(DB_DIALECT)

# Convenience aliases for dialect-specific extraction
.PHONY: extract-sqlc-schema-all-mysql
extract-sqlc-schema-all-mysql:
	@$(MAKE) extract-sqlc-schema-all DB_DIALECT=mysql

###############################################################################
# Generate sqlc code from current DB
# Usage: make regenerate-sqlc-from-current-db [DB_DIALECT=mysql|postgres]
###############################################################################
# need to run after database is reset and migrations are applied
.PHONY: regenerate-sqlc-from-current-db
regenerate-sqlc-from-current-db:
	$(MAKE) extract-sqlc-schema-all DB_DIALECT=$(DB_DIALECT)
	$(MAKE) sqlc-$(DB_DIALECT)

# Convenience aliases for dialect-specific regeneration
.PHONY: regenerate-sqlc-from-current-db-mysql
regenerate-sqlc-from-current-db-mysql:
	@$(MAKE) regenerate-sqlc-from-current-db DB_DIALECT=mysql

# sqlite queries must be run separately
# `make sqlc-sqlite`
