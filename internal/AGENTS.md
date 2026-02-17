# `internal/` Directory Guidelines

This document provides guidelines for AI agents working on packages in the `internal/` directory.

## Overview

The `internal/` directory contains **internal application code** following Clean Architecture principles.
These packages are application-specific and are **NOT** intended for external use.

**Key Characteristics:**

- Application-specific code (not reusable libraries)
- Clean Architecture with strict layer separation
- Implements dependency inversion (domain defines interfaces, infrastructure implements)

## Architecture Layers

| Layer | Directory | Purpose |
|-------|-----------|---------|
| Domain | `internal/domain/` | Pure business logic, ZERO infrastructure dependencies |
| Application | `internal/application/` | Use case orchestration, port interfaces, DTOs |
| Infrastructure | `internal/infrastructure/` | External dependencies, implementations only |
| Interface Adapters | `internal/interface-adapters/` | CLI, HTTP, wallet adapters |

### Layer Rules (See detailed rules in `.claude/rules/internal/`)

- **Domain**: No imports from other layers. Pure business logic.
- **Application**: Imports domain. Defines port interfaces in `ports/`.
- **Infrastructure**: Implements port interfaces. No interface definitions.
- **Interface Adapters**: Depends only on application layer.

## Directory Structure

```
internal/
├── domain/              # Business logic, entities, validators
│   ├── account/         # Account types and rules
│   ├── address/         # Address entity
│   ├── bitcoin/         # Bitcoin-specific types
│   ├── coin/            # Cryptocurrency types
│   ├── ethereum/        # Ethereum-specific types
│   ├── key/             # Key value objects
│   ├── multisig/        # MuSig2, multisig validators
│   ├── transaction/     # Transaction entity
│   ├── wallet/          # Wallet types, descriptors
│   └── xrp/             # XRP-specific types
├── application/         # Use cases and port interfaces
│   ├── dto/             # Data Transfer Objects
│   ├── ports/           # Interface definitions
│   └── usecase/         # keygen/, sign/, watch/
├── infrastructure/      # Implementations
│   ├── api/             # Blockchain API clients (btc, eth, xrp)
│   ├── database/        # SQLC-generated code (DO NOT EDIT)
│   ├── repository/      # Repository implementations
│   ├── storage/         # File storage implementations
│   └── wallet/key/      # Key generation (BIP32-86, HD wallet)
├── interface-adapters/  # External interfaces
│   ├── cli/             # CLI commands (keygen, sign, watch)
│   ├── http/            # HTTP handlers
│   └── wallet/          # Wallet adapters
├── di/                  # Dependency injection (panic allowed)
└── integration_test/    # Integration tests
```

## Critical Rules

### Dependency Direction

```
Interface Adapters → Application → Domain ← Infrastructure
```

### Never Do

- Import `infrastructure/` from `domain/` or `application/usecase/`
- Import `interface-adapters/` from `application/` or `infrastructure/`
- Define interfaces in `infrastructure/` (define in `application/ports/`)
- Edit files with `DO NOT EDIT` comments (auto-generated)

### Auto-Generated Files (DO NOT EDIT)

- `internal/infrastructure/database/*/sqlcgen/*.go`
- `internal/infrastructure/api/xrp/xrp/*.pb.go`

## Common Commands

```bash
make go-lint      # Lint and auto-fix
make check-build  # Build verification
make gotest       # Run unit tests
make tidy         # go mod tidy
```

## `internal/` vs `pkg/`

| Directory | Purpose | Can Import |
|-----------|---------|------------|
| `internal/` | Application-specific, Clean Architecture | `pkg/` |
| `pkg/` | Shared utilities, external use | NEVER `internal/` |

## References

### Rules (Auto-triggered by file paths)

- `.claude/rules/internal/clean-architecture.md` - Layer dependencies
- `.claude/rules/internal/domain-layer.md` - Domain layer rules
- `.claude/rules/internal/application-layer.md` - Application layer rules
- `.claude/rules/internal/infrastructure-layer.md` - Infrastructure layer rules
- `.claude/rules/internal/interface-adapters-layer.md` - Interface adapters layer rules
- `.claude/rules/go/conventions.md` - Go conventions
- `.claude/rules/go/repository.md` - Repository pattern
- `.claude/rules/go/usecase.md` - Use case pattern

### Guidelines

- [Coding Conventions](../docs/guidelines/coding-conventions.md)
- [Security](../docs/guidelines/security.md)
- [Testing](../docs/guidelines/testing.md)
- [Workflow](../docs/guidelines/workflow.md)
- [Architecture](../docs/guidelines/architecture.md)
- [Database](../docs/database/db-management.md)
- [Code Generation](../docs/guidelines/code-generation.md)

### Other

- [AGENTS.md](../AGENTS.md) - Root agent guidelines
- [pkg/AGENTS.md](../pkg/AGENTS.md) - Public packages
- [ARCHITECTURE.md](../ARCHITECTURE.md) - System architecture
