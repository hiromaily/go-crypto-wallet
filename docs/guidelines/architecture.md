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

## Application Ports Layer (Interface Definitions)

The `internal/application/ports/` package contains **interface definitions** (abstractions) for infrastructure dependencies, following the Dependency Inversion Principle.

**Key Principles:**

- **Interfaces are defined in the layer that uses them** (application layer), NOT in the infrastructure layer
- **Infrastructure layer contains only implementations**, never interface definitions
- This inverts the dependency direction: Infrastructure depends on Application, not vice versa
- Ports define contracts that infrastructure implementations must fulfill
- **Interfaces MUST use application-layer DTOs**, NOT infrastructure types (Clean Architecture Dependency Rule)

**Current Port Packages:**

```text
internal/application/ports/
├── bitcoin/
│   └── interface.go              # Bitcoiner interface (Bitcoin/BCH API abstraction)
├── persistence/
│   └── repository.go             # Repository interfaces (database abstractions)
└── storage/
    └── interface.go              # TransactionFileRepositorier interface (file storage abstraction)
```

**DTOs for Port Interfaces:**

Port interfaces must use **application-layer DTOs** (defined in `internal/application/dto/`) instead of infrastructure types. This ensures the application layer has zero dependencies on infrastructure.

**DTO Location:**

- DTOs are defined in `internal/application/dto/{coin}/` (e.g., `internal/application/dto/btc/dto.go`)
- DTOs contain only domain types, standard library types, or external library types (e.g., `btcutil.Amount`)
- DTOs MUST NOT import infrastructure packages

**Example: Bitcoin API Abstraction with DTOs**

```go
// DTOs in internal/application/dto/btc/dto.go
package btc

import "github.com/btcsuite/btcd/btcutil"

type AddressInfo struct {
    Address      string
    ScriptPubKey string
    IsWitness    bool
    Labels       []string
    // ... other fields
}

type UnspentOutput struct {
    TxID          string
    Vout          uint32
    Address       string
    Amount        btcutil.Amount
    Confirmations int64
    // ... other fields
}

// Interface definition in application/ports/btc/interface.go
package btc

import btcdto "internal/application/dto/btc"

type Bitcoiner interface {
    GetBalance() (btcutil.Amount, error)
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)
    ListUnspent(confirmationNum uint64) ([]btcdto.UnspentOutput, error)
    // ... other methods using DTOs
}

// Implementation in infrastructure/api/btc/btc/bitcoin.go
package btc

import (
    portsBtc "internal/application/ports/btc"
    btcdto "internal/application/dto/btc"
)

type Bitcoin struct {
    client *rpcclient.Client
}

// Bitcoin implements portsBtc.Bitcoiner interface
// Maps infrastructure types to application DTOs
func (b *Bitcoin) GetAddressInfo(addr string) (*btcdto.AddressInfo, error) {
    // Call Bitcoin Core RPC (returns infrastructure type)
    result, err := b.client.GetAddressInfo(addr)
    if err != nil {
        return nil, err
    }
    
    // Map infrastructure type to application DTO
    return &btcdto.AddressInfo{
        Address:      result.Address,
        ScriptPubKey: result.ScriptPubKey,
        IsWitness:    result.IsWitness,
        Labels:       result.Labels,
        // ... map all fields
    }, nil
}
```

**Why Interfaces Belong in Application Layer:**

1. **Dependency Inversion Principle**: High-level modules (application) should not depend on low-level modules (infrastructure). Both should depend on abstractions.
2. **Stable Abstractions**: The application defines what it needs; infrastructure provides implementations.
3. **Testability**: Application layer can be tested with mock implementations without infrastructure dependencies.
4. **Clean Architecture**: Core business logic (application) is independent of external frameworks and libraries (infrastructure).

**Important:**

- **NEVER** define interfaces in the infrastructure layer
- **ALWAYS** define interfaces in `application/ports/` when infrastructure needs abstraction
- **NEVER** use infrastructure types in port interface method signatures (use application DTOs instead)
- **ALWAYS** create DTOs in `internal/application/dto/` for data structures returned by infrastructure
- Infrastructure packages import and implement these port interfaces
- Infrastructure implementations map between infrastructure types and application DTOs
- Use cases depend on port interfaces, not concrete implementations

**DTO Creation Guidelines:**

When creating port interfaces that return data from infrastructure:

1. **Identify infrastructure types** used in interface methods
2. **Create corresponding DTOs** in `internal/application/dto/{coin}/dto.go` (e.g., `internal/application/dto/btc/dto.go`)
3. **Update interface methods** to use DTOs instead of infrastructure types
4. **Update infrastructure implementations** to map infrastructure types → DTOs
5. **Update use cases** to work with new DTOs

**Example Pattern:**

```go
// ❌ WRONG: Interface depends on infrastructure type
type Bitcoiner interface {
    GetAddressInfo(addr string) (*btc.GetAddressInfoResult, error)  // btc is infrastructure
}

// ✅ CORRECT: Interface uses application DTO
type Bitcoiner interface {
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)  // btcdto is application DTO
}
```

## Application Layer (Use Case) Guidelines

The `internal/application/usecase/` package implements the use case layer following Clean Architecture principles.

**Key Principles:**

- Use cases orchestrate business logic by coordinating domain objects and infrastructure services
- Each use case represents a single business operation with clear input and output
- Use cases act as thin wrappers that transform DTOs, delegate to services, and wrap errors with context
- **Use cases depend on interfaces defined in `application/ports/`** (Dependency Inversion)
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
    - `dto/`: **Data Transfer Objects** - Application-layer DTOs for port interfaces
      - `btc/`: Bitcoin DTOs (AddressInfo, UnspentOutput, TransactionResult, etc.)
      - Other coin DTOs (eth, xrp) as needed
    - `ports/`: **Interface definitions (abstractions)** - Contracts for infrastructure implementations
      - `btc/`: Bitcoiner interface (Bitcoin/BCH API abstraction)
      - `persistence/`: Repository interfaces (database abstractions)
      - `storage/`: File storage interfaces (TransactionFileRepositorier)
    - `usecase/`: Use case implementations organized by wallet type
      - `keygen/`: Key generation use cases (btc, eth, xrp, shared)
      - `sign/`: Signing use cases (btc, eth, xrp, shared)
      - `watch/`: Watch wallet use cases (btc, eth, xrp, shared)
  - `infrastructure/`: **Infrastructure layer** - **Implementation only** (NO interface definitions)
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
┌─────────────────────────────────────────────────────────────────────┐
│                    Clean Architecture Layers                         │
│                                                                      │
│  Interface Adapters (interface-adapters/*)                          │
│           ↓                                                          │
│  Application Layer (application/usecase/) ← depends on              │
│           ↓                                  ↓                       │
│  Application Ports (application/ports/) ← implemented by            │
│           ↓                                  ↓                       │
│  Domain Layer (domain/*)              Infrastructure Layer          │
│                                        (infrastructure/*)            │
│                                        - Implementations ONLY        │
│                                        - NO interface definitions    │
└─────────────────────────────────────────────────────────────────────┘

Key Principles:
1. Infrastructure implements interfaces defined in application/ports/
2. Application layer depends on application/ports/ (abstractions), not infrastructure (concrete)
3. This follows Dependency Inversion Principle (DIP)
4. Dependency flow: Interface Adapters → Application (Use Cases + Ports) ← Infrastructure
```

## See Also

- [Core Principles](core.md) - Security, error handling, and core patterns
- [Testing Guidelines](testing.md) - Testing strategy for each layer
- `internal/AGENTS.md` - Detailed guidelines for `internal/` directory
- `pkg/AGENTS.md` - Guidelines for `pkg/` directory
