# `pkg/` Directory Guidelines

Guidelines for AI agents working on packages in the `pkg/` directory.

## Overview

The `pkg/` directory contains **public packages** that can be imported by external code.
These packages provide shared utilities, configuration management, logging, and common functionality.

## Directory Structure

```
pkg/
├── config/         # Configuration management (account, wallet, loader)
├── cryptocurrency/ # Crypto utilities (btc, bch, eth, xrp)
├── db/             # Database connections (mysql, sqlite)
├── debug/          # Debug utilities
├── decimal/        # Decimal number utilities
├── di/             # Legacy DI container (panic allowed)
├── grpc/           # gRPC client utilities
├── logger/         # Logging (global, slog, noop)
├── null/           # Null value converters
├── retry/          # Retry with exponential backoff
├── serializer/     # Serialization utilities
├── uuid/           # UUID generation
└── websocket/      # WebSocket client
```

## Critical Rule

**Packages in `pkg/` MUST NOT import from `internal/` directory.**

This is non-negotiable:

- `pkg/` = public APIs for external code
- `internal/` = internal implementation details
- Mixing breaks encapsulation and creates circular dependencies

## Common Commands

```bash
make go-lint      # Lint and auto-fix
make check-build  # Build verification
make gotest       # Run unit tests
make tidy         # go mod tidy
```

## Patterns

### Avoid

- Importing from `internal/`
- Using `panic` (except `pkg/di` legacy)
- Commented-out code
- Ignoring errors
- Logging sensitive information

### Recommended

- Self-contained utility functions
- Proper error wrapping with context
- Use `context.Context` for cancellation
- Unit tests for exported functions

## Important Notes

- Financial project - make changes carefully
- Public APIs - consider backward compatibility
- **DO NOT** edit files with `DO NOT EDIT` comments

## References

### Rules

- `.claude/rules/go/conventions.md` - Go conventions

### Guidelines

- [Coding Conventions](../docs/guidelines/coding-conventions.md)
- [Security](../docs/guidelines/security.md)
- [Testing](../docs/guidelines/testing.md)

### Other

- [AGENTS.md](../AGENTS.md) - Root agent guidelines
- [internal/AGENTS.md](../internal/AGENTS.md) - Internal packages
