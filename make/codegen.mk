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

# Minimum buf version required for edition 2024 (expected support)
BUF_MIN_VERSION := 1.64

#------------------------------------------------------------------------------
# protoc-based generation (Recommended)
#------------------------------------------------------------------------------
# Direct protoc usage for edition 2024 support.
# Requires: protoc >= 33.0, protoc-gen-go, protoc-gen-go-grpc
#------------------------------------------------------------------------------

# Check protoc version (edition 2024 requires >= 33.0)
# Uses sort -V for robust semantic version comparison
.PHONY: proto-version-check
proto-version-check:
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed"; exit 1; }
	@PROTOC_VERSION=$$(protoc --version | sed 's/libprotoc //'); \
	if [ "$$(printf '%s\n' '$(PROTOC_MIN_VERSION)' "$$PROTOC_VERSION" | sort -V | head -n 1)" != "$(PROTOC_MIN_VERSION)" ]; then \
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
# Use 'make proto' (protoc) until buf adds support (expected in v1.64.0+).
#------------------------------------------------------------------------------

# Check buf version (edition 2024 expected to require >= 1.64.0)
# Uses sort -V for robust semantic version comparison
.PHONY: buf-version-check
buf-version-check:
	@command -v buf >/dev/null 2>&1 || { echo "Error: buf is not installed"; exit 1; }
	@BUF_VERSION=$$(buf --version); \
	if [ "$$(printf '%s\n' '$(BUF_MIN_VERSION)' "$$BUF_VERSION" | sort -V | head -n 1)" != "$(BUF_MIN_VERSION)" ]; then \
		echo "Error: buf version $$BUF_VERSION is too old. Edition 2024 requires buf >= $(BUF_MIN_VERSION).0"; \
		echo "Current version: $$BUF_VERSION"; \
		echo "Please upgrade buf: https://buf.build/docs/cli/installation/"; \
		exit 1; \
	fi; \
	echo "buf version check passed: $$BUF_VERSION"

# Generate Go code from proto files using buf
# NOTE: Requires buf >= 1.64.0 for edition 2024 support
.PHONY: proto-buf
proto-buf: buf-version-check clean-pb
	@echo "Generating proto files with buf..."
	buf generate

# Format proto files with buf
# NOTE: Fails with edition 2024 - use 'make proto-clang-fmt' instead
.PHONY: proto-fmt-buf
proto-fmt-buf: buf-version-check
	buf format -w

# Check proto file formatting with buf (for CI)
# NOTE: Fails with edition 2024 - use 'make proto-clang-fmt-check' instead
.PHONY: proto-fmt-buf-check
proto-fmt-buf-check: buf-version-check
	buf format -d --exit-code

# Lint proto files with buf
# NOTE: Fails with edition 2024 until buf adds support
.PHONY: proto-lint
proto-lint: buf-version-check
	buf lint

# Check for breaking changes in proto files
.PHONY: breaking-proto
breaking-proto:
	buf breaking --against '.git#branch=main'

#------------------------------------------------------------------------------
# clang-format based formatting (for edition 2024)
#------------------------------------------------------------------------------
# clang-format supports proto files and works with edition 2024.
# This is the recommended formatting method until buf adds edition 2024 support.
#------------------------------------------------------------------------------

# Format proto files with clang-format (recommended for edition 2024)
.PHONY: proto-fmt
proto-fmt:
	@command -v clang-format >/dev/null 2>&1 || { echo "Error: clang-format is not installed"; exit 1; }
	@echo "Formatting proto files with clang-format..."
	@find proto -name "*.proto" -exec clang-format -i {} \;
	@echo "Proto files formatted successfully"

# Check proto file formatting with clang-format (for CI)
.PHONY: proto-fmt-check
proto-fmt-check:
	@command -v clang-format >/dev/null 2>&1 || { echo "Error: clang-format is not installed"; exit 1; }
	@echo "Checking proto file formatting with clang-format..."
	@DIFF=$$(find proto -name "*.proto" -exec clang-format --dry-run -Werror {} \; 2>&1); \
	if [ -n "$$DIFF" ]; then \
		echo "Proto files are not formatted. Run 'make proto-fmt' to fix."; \
		echo "$$DIFF"; \
		exit 1; \
	fi; \
	echo "Proto file formatting check passed"

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
