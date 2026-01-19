###############################################################################
# Code Generator Targets
###############################################################################

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
# Protocol Buffer Code Generation
#------------------------------------------------------------------------------
# This project uses Protobuf Edition 2024.
# See docs/proto.md for details on edition support and tooling.
#
# Two generation methods are available:
# 1. protoc (native) - Recommended for edition 2024 support
# 2. buf - For when buf CLI adds edition 2024 support
#------------------------------------------------------------------------------

# Output directory for generated proto files
PROTO_OUT_DIR := internal/infrastructure/api/xrp/protogen

# Minimum protoc version required for edition 2024
PROTOC_MIN_VERSION := 33

#------------------------------------------------------------------------------
# protoc-based generation (Recommended)
#------------------------------------------------------------------------------
# Direct protoc usage for edition 2024 support.
# Requires: protoc >= 33.0, protoc-gen-go, protoc-gen-go-grpc
#------------------------------------------------------------------------------

# Check protoc version (edition 2024 requires >= 33.0)
.PHONY: proto-version-check
proto-version-check:
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed"; exit 1; }
	@PROTOC_VERSION=$$(protoc --version | sed 's/libprotoc //' | cut -d. -f1); \
	if [ "$$PROTOC_VERSION" -lt $(PROTOC_MIN_VERSION) ]; then \
		echo "Error: protoc version $$PROTOC_VERSION is too old. Edition 2024 requires protoc >= $(PROTOC_MIN_VERSION).0"; \
		echo "Current version: $$(protoc --version)"; \
		echo "Please upgrade protoc: https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi; \
	echo "protoc version check passed: $$(protoc --version)"

# Generate Go code from proto files using protoc (recommended for edition 2024)
.PHONY: proto
proto: proto-version-check clean-pb
	@echo "Generating proto files with protoc (edition 2024)..."
	protoc \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		-I proto/rippleapi \
		proto/rippleapi/*.proto

#------------------------------------------------------------------------------
# buf-based generation (for future use)
#------------------------------------------------------------------------------
# buf provides additional features:
# - Linting with buf lint
# - Breaking change detection with buf breaking
# - Managed mode for go_package options
#
# NOTE: As of buf CLI v1.63.0, edition 2024 is not yet supported.
# Use 'make proto' (protoc) until buf adds support.
#------------------------------------------------------------------------------

# Generate Go code from proto files using buf
# NOTE: Currently fails with edition 2024 - use 'make proto' instead
.PHONY: proto-buf
proto-buf: clean-pb
	@echo "Generating proto files with buf..."
	@echo "WARNING: buf CLI may not support edition 2024 yet"
	buf generate

# Format proto files with buf
.PHONY: proto-fmt
proto-fmt:
	buf format -w

# Check proto file formatting (for CI)
.PHONY: proto-fmt-check
proto-fmt-check:
	buf format -d --exit-code

# Lint proto files with buf
# NOTE: May fail with edition 2024 until buf adds support
.PHONY: proto-lint
proto-lint:
	buf lint

# Check for breaking changes in proto files
.PHONY: breaking-proto
breaking-proto:
	buf breaking --against '.git#branch=main'

#------------------------------------------------------------------------------
# Common targets
#------------------------------------------------------------------------------

# Clean generated protobuf files
.PHONY: clean-pb
clean-pb:
	rm -rf $(PROTO_OUT_DIR)/*.pb.go

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
