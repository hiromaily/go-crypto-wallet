# Install dependencies
.PHONY: install
install:
	bun install

# Development server with hot reload
.PHONY: dev
dev:
	bun --hot src/index.ts

# Build for production
.PHONY: build
build:
	bun build src/index.ts --outdir=dist --target=bun

# Lint code with Biome
.PHONY: lint
lint:
	bun run biome check .

# Fix lint issues
.PHONY: lint-fix
lint-fix:
	bun run biome check --write .

# Format code with Biome
.PHONY: format
format:
	bun run biome format --write .

# TypeScript type checking
.PHONY: typecheck
typecheck:
	bun run tsc --noEmit

# Generate proto files
.PHONY: proto
proto:
	bun run buf generate ../../proto/rippleapi

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf dist node_modules
