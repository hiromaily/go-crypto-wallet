###############################################################################
# Code Generator Targets
###############################################################################
# proto related targets
include make/codegen_proto.mk

###############################################################################
# generate sqlite schema from postgres scheme
###############################################################################
# WIP
.PHONY: generate-sqlite-schema
generate-sqlite-schema:
	pg2sqlite -i tools/sqlc/schemas/postgres/01_watch.sql \
	-o tools/sqlc/schemas/sqlite/01_watch.sql
	pg2sqlite -i tools/sqlc/schemas/postgres/02_keygen.sql \
	-o tools/sqlc/schemas/sqlite/02_keygen.sql
	pg2sqlite -i tools/sqlc/schemas/postgres/03_sign.sql \
	-o tools/sqlc/schemas/sqlite/03_sign.sql

###############################################################################
# Generate composite SQLite schemas for E2E testing
###############################################################################
# Creates per-wallet schemas in tools/sqlc/schemas/sqlite/e2e/ that include
# all tables each wallet needs at runtime. The sign wallet gets seed and
# musig2_nonces tables from the keygen schema (shared tables that can't be
# duplicated in SQLC schemas due to glob-based merging).
#------------------------------------------------------------------------------
.PHONY: generate-e2e-sqlite-schemas
generate-e2e-sqlite-schemas:
	./scripts/db/generate-e2e-sqlite-schemas.sh

########################################################################
# sqlc
#------------------------------------------------------------------------------
# Generate Go code from SQL queries for PostgreSQL using sqlc
# Schemas: tools/sqlc/schemas/postgres/*.sql
# Queries: tools/sqlc/queries/postgres/*.sql
# Output: internal/infrastructure/database/postgres/sqlcgen/
#------------------------------------------------------------------------------
.PHONY: sqlc
sqlc:
	cd tools/sqlc && sqlc generate -f sqlc_postgres.yml

# Alias for postgres (default config)
.PHONY: sqlc-postgres
sqlc-postgres: sqlc

# Generate Go code from SQL queries for MySQL using sqlc
# Schemas: tools/sqlc/schemas/mysql/*.sql
# Queries: tools/sqlc/queries/mysql/*.sql
# Output: internal/infrastructure/database/mysql/sqlcgen/
.PHONY: sqlc-mysql
sqlc-mysql:
	cd tools/sqlc && sqlc generate -f sqlc_mysql.yml

# Generate Go code from SQL queries for SQLite using sqlc
# Schemas: tools/sqlc/schemas/sqlite/*.sql
# Queries: tools/sqlc/queries/mysql/*.sql
# Output: internal/infrastructure/database/sqlite/sqlcgen/
.PHONY: sqlc-sqlite
sqlc-sqlite:
	cd tools/sqlc && sqlc generate -f sqlc_sqlite.yml

# Generate sqlc code for all database backends (MySQL, SQLite, and PostgreSQL)
.PHONY: sqlc-all
sqlc-all: sqlc sqlc-sqlite sqlc-postgres

###############################################################################
# mockery
#------------------------------------------------------------------------------
# Generate mock implementations from Go interfaces using mockery
# Configuration: .mockery.yaml
# Output: Internal mocks directories (next to interfaces)
#
# Usage:
#   make mockery        - Generate all mocks defined in .mockery.yaml (cleans old mocks first)
#   make clean-mocks    - Remove all generated mock files
#------------------------------------------------------------------------------
.PHONY: mockery
mockery: clean-mocks
	go tool github.com/vektra/mockery/v3

.PHONY: clean-mocks
clean-mocks:
	find . -type d -name "mocks" -exec rm -rf {} + 2>/dev/null || true


###############################################################################
# ABI
#------------------------------------------------------------------------------
# ABI code generation using abigen
# abigen is a tool for generating Go code from Ethereum smart contract ABIs
#------------------------------------------------------------------------------
# Generate ABI code from token.abi file using abigen
# Source: contracts/token.abi
# Output: internal/infrastructure/contract/token-abi.go
#------------------------------------------------------------------------------
.PHONY: generate-abi
generate-abi:
	abigen --abi ./contracts/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go
