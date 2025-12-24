# Agent Guidelines for go-crypto-wallet

This document provides guidelines for AI agents working on this project.

## Project Overview

- This project is a cryptocurrency wallet implementation in Go supporting BTC, BCH, ETH, XRP, and ERC-20 tokens
- The project is currently under refactoring based on Clean Architecture and Clean Code principles
- Security is of utmost importance (private key management, offline wallets)
- The project follows the `pkg` layout pattern

## Architecture Principles

- Follow Clean Architecture principles
- Maintain clear layer separation (domain, application, infrastructure)
- Use dependency injection and abstract with interfaces
- Follow the `pkg` layout pattern

### Domain Layer Guidelines

The `pkg/domain/` package contains pure business logic with **ZERO infrastructure dependencies**.

**Key Principles:**

- Domain layer has NO dependencies on infrastructure (no database, no API clients, no file I/O)
- Domain defines interfaces; infrastructure implements them (Dependency Inversion Principle)
- All domain logic must be testable without mocks (pure functions preferred)
- Domain is the most stable layer - changes here affect all other layers

**Domain Layer Structure:**

- **Types & Value Objects**: Immutable objects defined by values (AccountType, TxType, CoinTypeCode)
- **Entities**: Objects with unique identity and lifecycle (not yet fully implemented)
- **Validators**: Business rule validation functions
- **Domain Services**: Stateless services with business logic

**Important:**

- When adding new business logic, first consider if it belongs in the domain layer
- Use domain validators for input validation before infrastructure operations
- Business rules should be in domain, not scattered across services
- For backward compatibility, old packages (`pkg/wallet/types.go`, `pkg/account/types.go`, etc.)
  now provide type aliases to domain types

### Application Layer (Use Case) Guidelines

The `pkg/application/usecase/` package implements the use case layer following Clean Architecture principles.

**Key Principles:**

- Use cases orchestrate business logic by coordinating domain objects and infrastructure services
- Each use case represents a single business operation with clear input and output
- Use cases act as thin wrappers that transform DTOs, delegate to services, and wrap errors with context
- Use cases depend on domain layer and infrastructure layer through interfaces (Dependency Inversion)
- Organized by wallet type (watch, keygen, sign) and cryptocurrency (btc, eth, xrp, shared)

**Use Case Structure:**

```go
// Use case interface definition
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

// Input/Output DTOs
type XxxInput struct {
    Param1 string
    Param2 int
}

type XxxOutput struct {
    Result string
}

// Implementation
type xxxUseCase struct {
    service ServiceInterface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    result, err := u.service.SomeMethod(input.Param1, input.Param2)
    if err != nil {
        return nil, fmt.Errorf("failed to execute xxx: %w", err)
    }
    return &XxxOutput{Result: result}, nil
}
```

**DTO Conventions:**

- **Input DTOs**: Contain all parameters needed for the use case operation
- **Output DTOs**: Contain all results returned by the use case
- DTOs use domain types (not primitive types when domain types exist)
- DTOs are passed by value for inputs, returned as pointers for outputs

**Error Handling:**

- Wrap service errors with context using `fmt.Errorf` with `%w`
- Error messages should describe the use case operation that failed
- Return domain errors when business rule violations occur
- Let infrastructure errors propagate with added context

**Organization Structure:**

```
pkg/application/usecase/
├── keygen/
│   ├── interfaces.go              # Use case interfaces
│   ├── btc/                       # Bitcoin-specific use cases
│   ├── eth/                       # Ethereum-specific use cases
│   ├── xrp/                       # XRP-specific use cases
│   └── shared/                    # Shared use cases (all coins)
├── sign/
│   ├── interfaces.go
│   ├── btc/
│   ├── eth/
│   ├── xrp/
│   └── shared/
└── watch/
    ├── interfaces.go
    ├── btc/
    ├── eth/
    ├── xrp/
    └── shared/
```

**Testing Approach:**

Use cases currently have constructor tests that verify:
- Use case can be instantiated with dependencies
- Correct interface implementation

For comprehensive testing strategy, see `docs/TESTING_STRATEGY.md`.

**When to Create a New Use Case:**

- New command functionality is added (commands should use use cases, not services directly)
- Existing service logic needs to be exposed to commands with different DTO structure
- Business logic needs to coordinate multiple services
- Transaction boundaries need to be defined

**Important:**

- Commands in `pkg/command/` should ONLY depend on use cases, NOT services directly
- Use cases should be small and focused on a single operation
- Avoid business logic in use cases; delegate to domain or services
- Use cases are the entry point to application logic from command layer

## Coding Standards

- Follow `golangci-lint` configuration (`.golangci.yml`)
- Format code with `make format` (uses `gofumpt` and `goimports` via golangci-lint)
  - Import order: standard → third-party → local
- Use `make lint-fix` to run linting and formatting together (executes lint checks and format fixes)
- Maintain consistent naming conventions (lowercase package names, exported functions start with uppercase)

## Common Commands

After making code changes, use these commands to verify code correctness:

- `make lint-fix`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

**Command Constraints**:

- **DO NOT** use `go build -v` directly; use `make check-build` instead
- **DO NOT** use `go tool golangci-lint` directly; use `make lint-fix` instead

## Error Handling

- Wrap errors with `fmt.Errorf` + `%w`
- Use `errors.Is`/`errors.As` for error checking
- Include context information in error messages

## Panic Usage

Following the principle of separation of concerns, the project clearly separates the instance construction phase
from the instance usage phase.
Therefore, `panic` is only allowed during instance construction.
Specifically, `panic` is acceptable in:

- `main.go` files (application entry points)
- `pkg/di` package (dependency injection container)

**Important:**

- `panic` should **NOT** be used in business logic, service layers, or infrastructure layers
- Use proper error handling with error returns in all other packages
- The separation allows for fail-fast behavior during initialization while maintaining robust error handling during runtime

## Context Management

- Add `context.Context` to all API calls
- Implement timeouts and cancellation
- Implement graceful shutdown

## Security

- **NEVER** log private keys or sensitive information
- Encrypt or zero-clear private keys in memory when possible
- Do not pass passwords via command-line arguments; use secure input methods
- Conduct security review when making changes involving sensitive information

## Wallet Types Understanding

- **Watch Wallet**: Online, public keys only, creates and sends transactions
- **Keygen Wallet**: Offline, generates keys, first signature for multisig
- **Sign Wallet**: Offline, second and subsequent signatures for multisig

## Directory Structure

- `cmd/`: Application entry points (keygen, sign, watch)
- `pkg/`: Package code
  - `domain/`: **Domain layer** - Pure business logic (ZERO infrastructure dependencies)
    - `account/`: Account types, validators, and business rules
    - `transaction/`: Transaction types, state machine, validators
    - `wallet/`: Wallet types and definitions
    - `key/`: Key value objects and validators
    - `multisig/`: Multisig validators and business rules
    - `coin/`: Cryptocurrency type definitions
  - `application/`: **Application layer** - Use case layer (Clean Architecture)
    - `usecase/`: Use case implementations organized by wallet type
      - `keygen/`: Key generation use cases (btc, eth, xrp, shared)
      - `sign/`: Signing use cases (btc, eth, xrp, shared)
      - `watch/`: Watch wallet use cases (btc, eth, xrp, shared)
  - `infrastructure/`: **Infrastructure layer** - External dependencies and implementations
    - `api/`: External API clients
      - `bitcoin/`: Bitcoin/BCH Core RPC API clients (btc, bch)
      - `ethereum/`: Ethereum JSON-RPC API clients (eth, erc20)
      - `ripple/`: Ripple gRPC API clients (xrp)
    - `database/`: Database connections and generated code
      - `mysql/`: MySQL connection management
      - `sqlc/`: SQLC generated database code
    - `repository/`: Data persistence implementations
      - `cold/`: Cold wallet repository (keygen, sign)
      - `watch/`: Watch wallet repository
    - `storage/`: File storage implementations
      - `file/`: File-based storage (address, transaction)
    - `network/`: Network communication
      - `websocket/`: WebSocket client implementations
  - `wallet/service/`: **Application layer** - Business logic orchestration (legacy/transitional)
    - `keygen/`: Key generation services (btc, eth, xrp, shared)
    - `sign/`: Signing services (btc, eth, xrp, shared)
    - `watch/`: Watch wallet services (btc, eth, xrp, shared)
  - `wallet/key/`: Key generation logic - Infrastructure layer
  - `wallet/wallets/`: Wallet implementations (btcwallet, ethwallet, xrpwallet)
  - `command/`: Command implementations (keygen, sign, watch)
    - `keygen/`: Keygen command implementations (api, create, export, imports, sign)
    - `sign/`: Sign command implementations (create, export, imports, sign)
    - `watch/`: Watch command implementations (api, create, imports, monitor, send)
  - `di/`: Dependency injection container
  - `config/`: Configuration management
  - `logger/`: Logging utilities
  - `address/`: Address formatting and utilities (bch, xrp)
  - `account/`: Account-related utilities (backward compatibility type aliases)
  - `contract/`: Smart contract utilities (ERC-20 token ABI)
  - `converter/`: Data conversion utilities
  - `debug/`: Debug utilities
  - `fullpubkey/`: Full public key formatting utilities
  - `models/`: Data models (rdb)
  - `serial/`: Serialization utilities
  - `testutil/`: Test utilities (btc, eth, xrp, repository, suite)
  - `uuid/`: UUID generation utilities
- `data/`: Generated files, configuration files
  - `address/`: Address data files (bch, btc, eth, xrp)
  - `config/`: Configuration files (account, wallet configs, node configs)
  - `contract/`: Contract ABI files
  - `keystore/`: Keystore files
  - `proto/`: Protocol buffer definitions (rippleapi)
  - `tx/`: Transaction data files (bch, btc, eth, xrp)
- `scripts/`: Operation scripts
  - `operation/`: Wallet operation scripts
  - `setup/`: Setup scripts for blockchain nodes

**Architecture Dependency Direction:**

```text
Application Layer (application/usecase, wallet/service, command) → Domain Layer (domain/*) ← Infrastructure Layer (infrastructure/*, wallet/key)
```

## Refactoring Status

- Refer to `REFACTORING_CHECKLIST.md` for current refactoring tasks
- Make changes incrementally without breaking existing functionality
- Run tests before and after refactoring

## Testing

- Use `//go:build integration` tag for integration tests
- Separate unit tests and integration tests
- Measure and improve test coverage

## Dependency Management

- Use `go mod tidy` to organize dependencies
- Follow procedures in `REFACTORING_CHECKLIST.md` when updating dependencies
- Run security scans (`govulncheck`)

## Logging

- Use structured logging
- Set appropriate log levels
- **NEVER** log sensitive information (private keys, passwords, etc.)

## Patterns to Avoid

- Using `log.Fatal` (except in `main`)
- Using `panic` outside of instance construction (i.e., outside of `main.go` and `pkg/di` package)
- Leaving commented-out code
- Unused imports, variables, or functions
- Ignoring errors (detected by `errcheck`)
- Using `go build -v` directly (use `make check-build` instead)
- Using `go tool golangci-lint` directly (use `make lint-fix` instead)
- For `unused-receiver` lint errors: **Remove the receiver entirely** instead of renaming it to `_`.
  Renaming to `_` will only cause the same error to appear for other receivers,
  so convert the method to a function from the start.

## Recommended Patterns

- Abstraction through interfaces
- Dependency injection
- Proper error wrapping with context
- Use of `context.Context`

## Documentation

- Add godoc comments to exported functions and methods
- Add package-level comments
- Include usage examples for complex logic

## Multi-Chain Support

- **BTC/BCH**: Bitcoin Core RPC API
- **ETH**: Ethereum JSON-RPC API, ERC-20 token support
- **XRP**: Communication via gRPC with ripple-lib-server

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments
  (typically auto-generated files from tools like sqlc, protoc, or go generate)
- **Git Operations**:
  - Allowed: `git add`, `git commit`, and `git push` to GitHub
  - **NOT allowed**: `git merge` operations
  - **NOT allowed**: `gh` command merge operations (e.g., `gh pr merge`)
  - **NOT allowed**: `git commit` and `git push` to `main` or `master` branches
