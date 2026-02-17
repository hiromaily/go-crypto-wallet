###############################################################################
# sqlc Format, Validation, and Lint Targets
# Validates MySQL, SQLite, and PostgreSQL configurations
###############################################################################

# Compile SQL queries and schemas to check for syntax and type errors
# `sqlc compile` - Statically check SQL for syntax and type errors
# Runs for MySQL (default), SQLite, and PostgreSQL configurations
.PHONY: sqlc-compile
sqlc-compile:
	@echo "Compiling MySQL SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile
	@echo "✓ MySQL SQL compilation successful"
	@echo "Compiling SQLite SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL compilation successful"
	@echo "Compiling PostgreSQL SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile -f sqlc_postgresql.yml
	@echo "✓ PostgreSQL SQL compilation successful"

# Vet SQL queries for potential issues
# `sqlc vet` - Examines queries for potential problems
# Runs for MySQL (default), SQLite, and PostgreSQL configurations
.PHONY: sqlc-vet
sqlc-vet:
	@echo "Vetting MySQL SQL queries..."
	@cd tools/sqlc && sqlc vet
	@echo "✓ MySQL SQL queries passed vetting"
	@echo "Vetting SQLite SQL queries..."
	@cd tools/sqlc && sqlc vet -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL queries passed vetting"
	@echo "Vetting PostgreSQL SQL queries..."
	@cd tools/sqlc && sqlc vet -f sqlc_postgresql.yml
	@echo "✓ PostgreSQL SQL queries passed vetting"

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
# Validates MySQL, SQLite, and PostgreSQL configurations
.PHONY: sqlc-validate
sqlc-validate: sqlc-compile sqlc-vet
	@echo "✓ All sqlc validation checks passed (MySQL + SQLite + PostgreSQL)"

# Lint SQL queries (alias for vet)
# This target runs vet to check queries for potential issues
.PHONY: sqlc-lint
sqlc-lint: sqlc-vet

###############################################################################
# MySQL Schema Export Targets
###############################################################################

# Export watch schema from wallet-mysql container
.PHONY: dump-schema-watch
dump-schema-watch:
	mkdir -p data/dump/sql
	docker exec wallet-mysql mysqldump -u root -proot --no-data --skip-triggers watch > data/dump/sql/dump_watch.sql

# Export keygen schema from wallet-mysql container
.PHONY: dump-schema-keygen
dump-schema-keygen:
	mkdir -p data/dump/sql
	docker exec wallet-mysql mysqldump -u root -proot --no-data --skip-triggers keygen > data/dump/sql/dump_keygen.sql

# Export sign schema from wallet-mysql container
.PHONY: dump-schema-sign
dump-schema-sign:
	mkdir -p data/dump/sql
	docker exec wallet-mysql mysqldump -u root -proot --no-data --skip-triggers sign > data/dump/sql/dump_sign.sql

# Export all schemas from wallet-mysql container
.PHONY: dump-schema-all
dump-schema-all: dump-schema-watch dump-schema-keygen dump-schema-sign

###############################################################################
# MySQL Schema Extraction Targets (for sqlc)
###############################################################################

# Clean old sqlc schema files before extracting new ones
# This removes manually created schema files to prevent duplication errors
.PHONY: clean-sqlc-schemas
clean-sqlc-schemas:
	@echo "Cleaning old sqlc schema files..."
	@rm -f tools/sqlc/schemas/mysql/*.sql
	@echo "Done."

# Extract sqlc schema from watch dump file
# This extracts CREATE TABLE statements from MySQL dump and formats them for sqlc
.PHONY: extract-sqlc-schema-watch
extract-sqlc-schema-watch: dump-schema-watch
	@scripts/db/extract-sqlc-schema.sh \
		watch \
		data/dump/sql/dump_watch.sql \
		tools/sqlc/schemas/mysql/01_watch.sql

# Extract sqlc schema from keygen dump file
.PHONY: extract-sqlc-schema-keygen
extract-sqlc-schema-keygen: dump-schema-keygen
	@scripts/db/extract-sqlc-schema.sh \
		keygen \
		data/dump/sql/dump_keygen.sql \
		tools/sqlc/schemas/mysql/02_keygen.sql

# Extract sqlc schema from sign dump file
# Excludes seed and musig2_nonces tables (exist in keygen schema, avoid duplication)
.PHONY: extract-sqlc-schema-sign
extract-sqlc-schema-sign: dump-schema-sign
	@scripts/db/extract-sqlc-schema.sh \
		sign \
		data/dump/sql/dump_sign.sql \
		tools/sqlc/schemas/mysql/03_sign.sql

# Extract sqlc schemas from all dump files
# This cleans old schema files first, then extracts new ones from database dumps
.PHONY: extract-sqlc-schema-all
extract-sqlc-schema-all: clean-sqlc-schemas
	@$(MAKE) extract-sqlc-schema-watch
	@$(MAKE) extract-sqlc-schema-keygen
	@$(MAKE) extract-sqlc-schema-sign

###############################################################################
# MySQL: Generate sqlc code from current DB
###############################################################################
# need to run after database is reset and migrations are applied
.PHONY: regenerate-sqlc-from-current-db
regenerate-sqlc-from-current-db:
	make extract-sqlc-schema-all
	make sqlc

###############################################################################
# PostgreSQL Schema Export Targets
###############################################################################

# Export watch schema from wallet-postgres container
.PHONY: dump-schema-postgresql-watch
dump-schema-postgresql-watch:
	mkdir -p data/dump/sql
	docker exec wallet-postgres pg_dump -U postgres --schema-only --no-owner --no-privileges watch > data/dump/sql/dump_watch.pg.sql

# Export keygen schema from wallet-postgres container
.PHONY: dump-schema-postgresql-keygen
dump-schema-postgresql-keygen:
	mkdir -p data/dump/sql
	docker exec wallet-postgres pg_dump -U postgres --schema-only --no-owner --no-privileges keygen > data/dump/sql/dump_keygen.pg.sql

# Export sign schema from wallet-postgres container
.PHONY: dump-schema-postgresql-sign
dump-schema-postgresql-sign:
	mkdir -p data/dump/sql
	docker exec wallet-postgres pg_dump -U postgres --schema-only --no-owner --no-privileges sign > data/dump/sql/dump_sign.pg.sql

# Export all schemas from wallet-postgres container
.PHONY: dump-schema-postgresql-all
dump-schema-postgresql-all: dump-schema-postgresql-watch dump-schema-postgresql-keygen dump-schema-postgresql-sign

###############################################################################
# PostgreSQL Schema Extraction Targets (for sqlc)
###############################################################################

# Clean old PostgreSQL sqlc schema files before extracting new ones
.PHONY: clean-sqlc-schemas-postgresql
clean-sqlc-schemas-postgresql:
	@echo "Cleaning old PostgreSQL sqlc schema files..."
	@rm -f tools/sqlc/schemas/postgresql/*.sql
	@echo "Done."

# Extract PostgreSQL sqlc schema from watch dump file
.PHONY: extract-sqlc-schema-postgresql-watch
extract-sqlc-schema-postgresql-watch: dump-schema-postgresql-watch
	@scripts/db/extract-sqlc-schema-postgresql.sh \
		watch \
		data/dump/sql/dump_watch.pg.sql \
		tools/sqlc/schemas/postgresql/01_watch.sql

# Extract PostgreSQL sqlc schema from keygen dump file
.PHONY: extract-sqlc-schema-postgresql-keygen
extract-sqlc-schema-postgresql-keygen: dump-schema-postgresql-keygen
	@scripts/db/extract-sqlc-schema-postgresql.sh \
		keygen \
		data/dump/sql/dump_keygen.pg.sql \
		tools/sqlc/schemas/postgresql/02_keygen.sql

# Extract PostgreSQL sqlc schema from sign dump file
# Excludes seed and musig2_nonces tables (exist in keygen schema, avoid duplication)
.PHONY: extract-sqlc-schema-postgresql-sign
extract-sqlc-schema-postgresql-sign: dump-schema-postgresql-sign
	@scripts/db/extract-sqlc-schema-postgresql.sh \
		sign \
		data/dump/sql/dump_sign.pg.sql \
		tools/sqlc/schemas/postgresql/03_sign.sql

# Extract PostgreSQL sqlc schemas from all dump files
.PHONY: extract-sqlc-schema-postgresql-all
extract-sqlc-schema-postgresql-all: clean-sqlc-schemas-postgresql
	@$(MAKE) extract-sqlc-schema-postgresql-watch
	@$(MAKE) extract-sqlc-schema-postgresql-keygen
	@$(MAKE) extract-sqlc-schema-postgresql-sign

###############################################################################
# PostgreSQL: Generate sqlc code from current DB
###############################################################################
# need to run after database is reset and migrations are applied
.PHONY: regenerate-sqlc-postgresql-from-current-db
regenerate-sqlc-postgresql-from-current-db:
	make extract-sqlc-schema-postgresql-all
	make sqlc-postgresql
