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
	bun run lint

# Fix lint issues
.PHONY: lint-fix
lint-fix:
	bun run lint:fix

# Format code with Biome
.PHONY: format
format:
	bun run format

# TypeScript type checking
.PHONY: typecheck
typecheck:
	bun run typecheck

# Generate proto files (TypeScript only via root Makefile)
.PHONY: proto
proto:
	bun run proto

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf dist node_modules src/gen
