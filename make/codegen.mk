###############################################################################
# Code Generator Targets
###############################################################################
# proto related targets
include make/codegen_proto.mk

###############################################################################
# sqlc
#------------------------------------------------------------------------------
# Generate Go code from SQL queries using sqlc
# Schemas: tools/sqlc/schemas/*.sql
# Queries: tools/sqlc/queries/*.sql
# Output: internal/infrastructure/database/mysql/sqlcgen/
#------------------------------------------------------------------------------
.PHONY: sqlc
sqlc:
	cd tools/sqlc && sqlc generate

# Generate Go code from SQL queries for SQLite using sqlc
# Schemas: tools/sqlc/schemas_sqlite/*.sql
# Queries: tools/sqlc/queries/*.sql
# Output: internal/infrastructure/database/sqlite/sqlcgen/
#------------------------------------------------------------------------------
.PHONY: sqlc-sqlite
sqlc-sqlite:
	cd tools/sqlc && sqlc generate -f sqlc_sqlite.yml

# Generate sqlc code for all database backends (MySQL and SQLite)
#------------------------------------------------------------------------------
.PHONY: sqlc-all
sqlc-all: sqlc sqlc-sqlite

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
