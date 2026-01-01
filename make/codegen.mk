###############################################################################
# Code Generator Targets
###############################################################################

###############################################################################
# sqlc
#------------------------------------------------------------------------------
# Generate Go code from SQL queries using sqlc
# Schemas: tools/sqlc/schemas/*.sql
# Queries: tools/sqlc/queries/*.sql
# Output: internal/infrastructure/database/sqlc/
#------------------------------------------------------------------------------
.PHONY: sqlc
sqlc:
	cd tools/sqlc && sqlc generate

# ABI
.PHONY: generate-abi
generate-abi:
	abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go

###############################################################################
# mockery
#------------------------------------------------------------------------------
# Generate mock implementations from Go interfaces using mockery
# Configuration: .mockery.yaml
# Output: Internal mocks directories (next to interfaces)
#
# Usage:
#   make mockery        - Generate all mocks defined in .mockery.yaml
#   make clean-mocks    - Remove all generated mock files
#------------------------------------------------------------------------------
.PHONY: mockery
mockery:
	go tool github.com/vektra/mockery/v3

.PHONY: clean-mocks
clean-mocks:
	find . -type d -name "mocks" -exec rm -rf {} + 2>/dev/null || true

###############################################################################
# Protocol Buffer (buf-based generation)
#------------------------------------------------------------------------------
# Protocol Buffer code generation using buf
# buf replaces direct protoc usage and provides:
# - Linting with buf lint
# - Breaking change detection with buf breaking
# - Code generation with buf generate
#------------------------------------------------------------------------------

# Lint proto files with buf
.PHONY: lint-proto
lint-proto:
	buf lint

# Check for breaking changes in proto files
.PHONY: breaking-proto
breaking-proto:
	buf breaking --against '.git#branch=master'

# Generate Go code from proto files using buf
.PHONY: protoc-go
protoc-go: clean-pb
	buf generate

# Clean generated protobuf files
.PHONY: clean-pb
clean-pb:
	rm -rf internal/infrastructure/api/ripple/xrp/*.pb.go
