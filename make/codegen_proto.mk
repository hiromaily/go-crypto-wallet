###############################################################################
# Protocol Buffer Code Generation
#------------------------------------------------------------------------------
# This project uses Protobuf Edition 2024.
# See docs/proto.md for details on edition support and tooling.
#
# Two generation methods are available:
# 1. protoc (native) - Currently used (edition 2024 support)
# 2. buf - For future use when buf CLI adds edition 2024 support
#
# Usage:
#   make proto          - Generate Go proto files (uses PROTO_TOOL)
#   make proto-ts       - Generate TypeScript proto files (uses PROTO_TOOL)
#   make proto-all      - Generate all proto files
#   make proto-fmt      - Format proto files
#   make proto-lint     - Lint proto files (buf only)
#
# Override tool selection:
#   make proto PROTO_TOOL=buf      - Use buf for generation
#   make proto PROTO_TOOL=protoc   - Use protoc for generation
#------------------------------------------------------------------------------

#------------------------------------------------------------------------------
# Configuration
#------------------------------------------------------------------------------

# Tool selection: protoc (default) or buf
# Change to 'buf' when buf CLI >= 1.64.0 adds edition 2024 support
PROTO_TOOL ?= protoc

# Output directories
PROTO_GO_OUT_DIR := internal/infrastructure/api/xrp/protogen
PROTO_TS_OUT_DIR := apps/xrpl-grpc-server/src/gen

# Minimum versions required for edition 2024
PROTOC_MIN_VERSION := 33
BUF_MIN_VERSION := 1.64

#------------------------------------------------------------------------------
# Wrapper Targets (use PROTO_TOOL flag)
#------------------------------------------------------------------------------

# Generate Go code from proto files
.PHONY: proto
proto:
ifeq ($(PROTO_TOOL),buf)
	@$(MAKE) proto-go-buf
else
	@$(MAKE) proto-go-protoc
endif

# Generate TypeScript code from proto files
.PHONY: proto-ts
proto-ts:
ifeq ($(PROTO_TOOL),buf)
	@$(MAKE) proto-ts-buf
else
	@$(MAKE) proto-ts-protoc
endif

# Generate all proto files (Go + TypeScript)
.PHONY: proto-all
proto-all: proto proto-ts

# Format proto files
.PHONY: proto-fmt
proto-fmt:
ifeq ($(PROTO_TOOL),buf)
	@$(MAKE) proto-fmt-buf
else
	@$(MAKE) proto-fmt-clang
endif

# Check proto file formatting (for CI)
.PHONY: proto-fmt-check
proto-fmt-check:
ifeq ($(PROTO_TOOL),buf)
	@$(MAKE) proto-fmt-buf-check
else
	@$(MAKE) proto-fmt-clang-check
endif

# Lint proto files (buf only)
.PHONY: proto-lint
proto-lint:
ifeq ($(PROTO_TOOL),buf)
	@$(MAKE) proto-lint-buf
else
	@echo "Proto linting requires buf. Run: make proto-lint PROTO_TOOL=buf"
	@echo "Note: buf >= $(BUF_MIN_VERSION).0 required for edition 2024"
endif

#------------------------------------------------------------------------------
# protoc-based targets
#------------------------------------------------------------------------------
# Direct protoc usage for edition 2024 support.
# Requires: protoc >= 33.0, protoc-gen-go, protoc-gen-go-grpc
#------------------------------------------------------------------------------

# Check protoc version (edition 2024 requires >= 33.0)
.PHONY: protoc-version-check
protoc-version-check:
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed"; exit 1; }
	@PROTOC_VERSION=$$(protoc --version | sed 's/libprotoc //'); \
	if [ "$$(printf '%s\n' '$(PROTOC_MIN_VERSION)' "$$PROTOC_VERSION" | sort -V | head -n 1)" != "$(PROTOC_MIN_VERSION)" ]; then \
		echo "Error: protoc version $$PROTOC_VERSION is too old. Edition 2024 requires protoc >= $(PROTOC_MIN_VERSION).0"; \
		echo "Current version: $$(protoc --version)"; \
		echo "Please upgrade protoc: https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi; \
	echo "protoc version check passed: $$(protoc --version)"

# Generate Go code using protoc
.PHONY: proto-go-protoc
proto-go-protoc: protoc-version-check clean-pb-go
	@echo "Generating Go proto files with protoc (edition 2024)..."
	protoc \
		--go_out=$(PROTO_GO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		-I proto/rippleapi \
		proto/rippleapi/*.proto

# Generate TypeScript code using protoc
# Requires: protoc-gen-es >= 2.x (install via: cd apps/xrpl-grpc-server && bun install)
# Note: protoc-gen-es v2 generates both messages and services in *_pb.ts files
# Note: protoc-gen-connect-es is NOT used because it doesn't support Edition 2024
.PHONY: proto-ts-protoc
proto-ts-protoc: protoc-version-check clean-pb-ts
	@echo "Generating TypeScript proto files with protoc (edition 2024)..."
	@test -f apps/xrpl-grpc-server/node_modules/.bin/protoc-gen-es || { \
		echo "Error: protoc-gen-es is not installed"; \
		echo "Run: cd apps/xrpl-grpc-server && bun install"; \
		exit 1; \
	}
	@mkdir -p $(PROTO_TS_OUT_DIR)
	protoc \
		--plugin=protoc-gen-es=apps/xrpl-grpc-server/node_modules/.bin/protoc-gen-es \
		--es_out=$(PROTO_TS_OUT_DIR) \
		--es_opt=target=ts \
		-I proto/rippleapi \
		proto/rippleapi/*.proto

#------------------------------------------------------------------------------
# buf-based targets (for future use)
#------------------------------------------------------------------------------
# buf provides additional features: linting, breaking change detection, managed mode
# NOTE: As of buf CLI v1.63.0, edition 2024 is not yet supported.
# These targets will work when buf >= 1.64.0 is released.
#------------------------------------------------------------------------------

# Check buf version (edition 2024 requires >= 1.64.0)
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

# Generate Go code using buf
.PHONY: proto-go-buf
proto-go-buf: buf-version-check clean-pb-go
	@echo "Generating Go proto files with buf..."
	buf generate

# Generate TypeScript code using buf
# Note: buf.gen.yaml must include TypeScript plugins
.PHONY: proto-ts-buf
proto-ts-buf: buf-version-check clean-pb-ts
	@echo "Generating TypeScript proto files with buf..."
	buf generate

# Format proto files with buf
.PHONY: proto-fmt-buf
proto-fmt-buf: buf-version-check
	@echo "Formatting proto files with buf..."
	buf format -w

# Check proto file formatting with buf (for CI)
.PHONY: proto-fmt-buf-check
proto-fmt-buf-check: buf-version-check
	buf format -d --exit-code

# Lint proto files with buf
.PHONY: proto-lint-buf
proto-lint-buf: buf-version-check
	buf lint

# Check for breaking changes in proto files
.PHONY: proto-breaking
proto-breaking: buf-version-check
	buf breaking --against '.git#branch=main'

#------------------------------------------------------------------------------
# clang-format based formatting (for edition 2024)
#------------------------------------------------------------------------------
# clang-format supports proto files and works with edition 2024.
# Used as default formatter until buf adds edition 2024 support.
#------------------------------------------------------------------------------

# Format proto files with clang-format
.PHONY: proto-fmt-clang
proto-fmt-clang:
	@command -v clang-format >/dev/null 2>&1 || { echo "Error: clang-format is not installed"; exit 1; }
	@echo "Formatting proto files with clang-format..."
	@find proto -name "*.proto" -exec clang-format -i {} \;
	@echo "Proto files formatted successfully"

# Check proto file formatting with clang-format (for CI)
.PHONY: proto-fmt-clang-check
proto-fmt-clang-check:
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
# Clean targets
#------------------------------------------------------------------------------

# Clean generated Go protobuf files
.PHONY: clean-pb-go
clean-pb-go:
	rm -rf $(PROTO_GO_OUT_DIR)/*.pb.go

# Clean generated TypeScript protobuf files
.PHONY: clean-pb-ts
clean-pb-ts:
	rm -rf $(PROTO_TS_OUT_DIR)/*

# Clean all generated protobuf files
.PHONY: clean-pb
clean-pb: clean-pb-go clean-pb-ts
