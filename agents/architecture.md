# Architecture Guidelines

This document describes the Clean Architecture principles and layer guidelines for the go-crypto-wallet project.

## Architecture Principles

- Follow Clean Architecture principles
- Maintain clear layer separation (domain, application, infrastructure)
- Use dependency injection and abstract with interfaces
- Follow the `pkg` layout pattern

## Domain Layer Guidelines

The `internal/domain/` package contains pure business logic with **ZERO infrastructure dependencies**.

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

## Application Layer (Use Case) Guidelines

The `internal/application/usecase/` package implements the use case layer following Clean Architecture principles.

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

```text
internal/application/usecase/
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

For comprehensive testing strategy, see [Testing Guidelines](testing.md).

**When to Create a New Use Case:**

- New command functionality is added (commands should use use cases, not services directly)
- Existing service logic needs to be exposed to commands with different DTO structure
- Business logic needs to coordinate multiple services
- Transaction boundaries need to be defined

**Important:**

- Commands in `internal/interface-adapters/cli/` should ONLY depend on use cases, NOT services directly
- Use cases should be small and focused on a single operation
- Avoid business logic in use cases; delegate to domain or services
- Use cases are the entry point to application logic from command layer

## Directory Structure

- `cmd/`: Application entry points (keygen, sign, watch)
- `internal/`: Internal packages (application-specific, not for external use)
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
    - `contract/`: Smart contract utilities (ERC-20 token ABI generated code)
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
    - `wallet/key/`: Key generation logic - Infrastructure layer
  - `interface-adapters/`: **Interface Adapters layer** - Adapters between use cases and external interfaces
    - `cli/`: CLI command adapters (keygen, sign, watch)
      - `keygen/`: Keygen command implementations (api, create, export, imports, sign)
      - `sign/`: Sign command implementations (create, export, imports, sign)
      - `watch/`: Watch command implementations (api, create, imports, monitor, send)
    - `wallet/`: Wallet adapter interfaces and implementations
      - `interfaces.go`: Wallet interfaces (Keygener, Signer, Watcher)
      - `btc/`: Bitcoin wallet implementations
      - `eth/`: Ethereum wallet implementations
      - `xrp/`: XRP wallet implementations
  - `wallet/service/`: **Application layer** - Business logic orchestration (legacy/transitional)
    - `keygen/`: Key generation services (btc, eth, xrp, shared)
    - `sign/`: Signing services (btc, eth, xrp, shared)
    - `watch/`: Watch wallet services (btc, eth, xrp, shared)
  - `di/`: Dependency injection container
- `pkg/`: Shared packages (reusable, for external use)
  - `config/`: Configuration management utilities
    - `testutil/`: Test utilities for configuration
  - `logger/`: Logging utilities (structured logging, noop logger, slog support)
  - `converter/`: Data conversion utilities
  - `debug/`: Debug utilities
  - `serial/`: Serialization utilities
  - `testutil/`: Test utilities for various components (btc, eth, xrp, repository, suite)
  - `uuid/`: UUID generation utilities
  - `db/mysql/`: MySQL database connection utilities
  - `decimal/`: Decimal number utilities
  - `grpc/`: gRPC client utilities
  - `websocket/`: WebSocket client utilities
  - `di/`: Legacy dependency injection container (for backward compatibility)

  **Important**: See `pkg/AGENTS.md` for detailed guidelines on working with `pkg/` packages.
  **Critical Rule**: Packages in `pkg/` MUST NOT import or depend on any packages in `internal/` directory.
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
Interface Adapters (interface-adapters/*) → Application Layer (application/usecase, wallet/service) → Domain Layer (domain/*) ← Infrastructure Layer (infrastructure/*)
```

## See Also

- [Core Principles](core.md) - Security, error handling, and core patterns
- [Testing Guidelines](testing.md) - Testing strategy for each layer
- [`internal/AGENTS.md`](../internal/AGENTS.md) - Detailed guidelines for `internal/` directory
- [`pkg/AGENTS.md`](../pkg/AGENTS.md) - Guidelines for `pkg/` directory
