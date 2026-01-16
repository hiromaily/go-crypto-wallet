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
# Protocol Buffer (buf-based generation)
#------------------------------------------------------------------------------
# Protocol Buffer code generation using buf
# buf replaces direct protoc usage and provides:
# - Linting with buf lint
# - Breaking change detection with buf breaking
# - Code generation with buf generate
#------------------------------------------------------------------------------

# Format proto files with buf
.PHONY: proto-fmt
proto-fmt:
	buf format -w

# Check proto file formatting (for CI)
.PHONY: proto-fmt-check
proto-fmt-check:
	buf format -d --exit-code

# Lint proto files with buf
.PHONY: proto-lint
proto-lint:
	buf lint

# Check for breaking changes in proto files
.PHONY: breaking-proto
breaking-proto:
	buf breaking --against '.git#branch=main'

# Generate Go code from proto files using buf
.PHONY: protoc-go
protoc-go: clean-pb
	buf generate

# Clean generated protobuf files
.PHONY: clean-pb
clean-pb:
	rm -rf internal/infrastructure/api/ripple/xrp/*.pb.go

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
