# `internal/` Directory Guidelines

This document provides guidelines for AI agents working on packages in the `internal/` directory.

## Overview

The `internal/` directory contains **internal application code** that follows Clean Architecture principles.
These packages are application-specific and are **NOT** intended for external use (unlike `pkg/` packages).

**Key Characteristics:**

- Application-specific code (not reusable libraries)
- Follows Clean Architecture with clear layer separation
- Organized by architectural layers (domain, application, infrastructure, interface-adapters)
- Implements dependency inversion (domain defines interfaces, infrastructure implements them)

## Architecture Layers

The `internal/` directory is organized into four main layers following Clean Architecture:

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Pure business logic with **ZERO infrastructure dependencies**

**Key Principles:**

- **NO dependencies** on infrastructure (no database, no API clients, no file I/O)
- Defines interfaces that infrastructure must implement (Dependency Inversion Principle)
- Contains pure business logic that can be tested without mocks
- Most stable layer - changes here affect all other layers

**Structure:**

- `account/`: Account types, validators, and business rules
- `transaction/`: Transaction types, state machine, validators
- `wallet/`: Wallet types and definitions
- `key/`: Key value objects and validators
- `multisig/`: Multisig validators and business rules
- `coin/`: Cryptocurrency type definitions

**When Working in Domain Layer:**

- ✅ Add business rules and validators
- ✅ Define value objects and entities
- ✅ Create domain interfaces for infrastructure to implement
- ✅ Write pure functions (no side effects)
- ❌ **NEVER** import from `infrastructure/`, `application/`, or `interface-adapters/`
- ❌ **NEVER** use database, API clients, or file I/O
- ❌ **NEVER** add logging (use `pkg/logger/` if absolutely necessary, but prefer no logging)

**Testing:**

- Use pure unit tests without mocks
- Test business logic in isolation
- No infrastructure dependencies required

### 2. Application Layer (`internal/application/`)

**Purpose**: Use case orchestration, interface definitions (ports), and business logic coordination

**Key Principles:**

- Orchestrates business logic by coordinating domain objects and infrastructure services
- **Defines interfaces (ports) for infrastructure dependencies** (Dependency Inversion Principle)
- Each use case represents a single business operation
- Organized by wallet type (watch, keygen, sign) and cryptocurrency (btc, eth, xrp, shared)

**Structure:**

- `dto/`: **Data Transfer Objects** - Application-layer DTOs for port interfaces
- `btc/`: Bitcoin DTOs (AddressInfo, UnspentOutput, TransactionResult, etc.)
- Other coin DTOs (eth, xrp) as needed
- `ports/`: **Interface definitions (abstractions)** - Contracts that infrastructure must implement
  - `btc/`: Bitcoiner interface (Bitcoin/BCH API abstraction)
  - `persistence/`: Repository interfaces (database abstractions)
  - `storage/`: File storage interfaces (TransactionFileRepositorier)
- `usecase/`: Use case implementations
  - `keygen/`: Key generation use cases
  - `sign/`: Signing use cases
  - `watch/`: Watch wallet use cases

**Application Ports (`application/ports/`):**

**CRITICAL PRINCIPLE**: **Interfaces are defined in the layer that uses them, NOT in the infrastructure layer.**

- ✅ Define interfaces in `application/ports/` when infrastructure needs abstraction
- ✅ Infrastructure packages import and implement these port interfaces
- ✅ Use cases depend on port interfaces, not concrete implementations
- ✅ **ALWAYS** use application-layer DTOs in interface method signatures (NOT infrastructure types)
- ❌ **NEVER** define interfaces in the infrastructure layer
- ❌ **NEVER** import concrete infrastructure types in use cases
- ❌ **NEVER** use infrastructure types in port interface method signatures

**Example:**

```go
// DTOs in internal/application/dto/btc/dto.go
package btc

import "github.com/btcsuite/btcd/btcutil"

type AddressInfo struct {
    Address      string
    ScriptPubKey string
    IsWitness    bool
    Labels       []string
}

// Interface definition in application/ports/btc/interface.go
package btc

import btcdto "internal/application/dto/btc"

type Bitcoiner interface {
    GetBalance() (btcutil.Amount, error)
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)  // Uses DTO, not infrastructure type
}

// Implementation in infrastructure/api/bitcoin/btc/bitcoin.go
package btc

import (
    portsBtc "internal/application/ports/btc"
    btcdto "internal/application/dto/btc"
)

type Bitcoin struct {
    client *rpcclient.Client
}

// Bitcoin implements portsBtc.Bitcoiner
// Maps infrastructure type to application DTO
func (b *Bitcoin) GetAddressInfo(addr string) (*btcdto.AddressInfo, error) {
    // Call infrastructure API (returns infrastructure type)
    result, err := b.client.GetAddressInfo(addr)
    if err != nil {
        return nil, err
    }

    // Map to application DTO
    return &btcdto.AddressInfo{
        Address:      result.Address,
        ScriptPubKey: result.ScriptPubKey,
        IsWitness:    result.IsWitness,
        Labels:       result.Labels,
    }, nil
}

// Use case depends on interface, not implementation
type xxxUseCase struct {
    btcClient portsBtc.Bitcoiner  // Interface from ports
}
```

**When Working in Application Layer:**

- ✅ Create use cases that orchestrate domain and infrastructure
- ✅ Define interfaces in `application/ports/` for infrastructure abstractions
- ✅ Create DTOs in `internal/application/dto/` for data structures returned by infrastructure
- ✅ Use application DTOs in port interface method signatures (NOT infrastructure types)
- ✅ Transform DTOs between layers
- ✅ Wrap errors with context
- ✅ Depend on domain layer and `application/ports/` interfaces
- ❌ **NEVER** import from `interface-adapters/` (dependency direction violation)
- ❌ **NEVER** contain business logic (delegate to domain)
- ❌ **NEVER** directly access infrastructure implementations (use port interfaces)
- ❌ **NEVER** define interfaces in infrastructure layer
- ❌ **NEVER** use infrastructure types in port interface method signatures

**Use Case Pattern:**

```go
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

type xxxUseCase struct {
    service ServiceInterface  // Domain or infrastructure interface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    // Validate input using domain validators
    if err := domain.ValidateXxx(input); err != nil {
        return nil, fmt.Errorf("invalid input: %w", err)
    }

    // Delegate to service
    result, err := u.service.SomeMethod(input.Param1)
    if err != nil {
        return nil, fmt.Errorf("failed to execute xxx: %w", err)
    }

    return &XxxOutput{Result: result}, nil
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: External dependencies and technical implementations - **IMPLEMENTATIONS ONLY, NO INTERFACE DEFINITIONS**

**Key Principles:**

- **Implements interfaces defined in `application/ports/`** (Dependency Inversion Principle)
- **Contains ONLY concrete implementations, NEVER interface definitions**
- Contains NO business logic (only technical implementation)
- Easily replaceable and mockable for testing
- Organized by technical concern (database, api, repository, storage, network)

**CRITICAL PRINCIPLE**: **Infrastructure layer contains NO abstraction layer. All interfaces are defined in `application/ports/`.**

**Structure:**

- `api/`: External API clients (Bitcoin, Ethereum, Ripple) - **implementations only**
- `database/`: Database connections and generated code (MySQL, sqlc)
- `repository/`: Data persistence implementations - **implements `application/ports/persistence` interfaces**
- `storage/`: File storage implementations - **implements `application/ports/storage` interfaces**
- `network/`: Network communication (WebSocket, gRPC)
- `contract/`: Smart contract utilities
- `wallet/key/`: Key generation logic (infrastructure concern)
- `config/`: Configuration management

**When Working in Infrastructure Layer:**

- ✅ Implement interfaces from `application/ports/`
- ✅ Import from `application/ports/` to implement port interfaces
- ✅ Import from `application/dto/` to map infrastructure types to DTOs
- ✅ Map infrastructure types to application DTOs in interface implementations
- ✅ Handle external system communication
- ✅ Convert between domain entities and external formats
- ✅ Manage technical concerns (database, network, file I/O)
- ❌ **NEVER** contain business logic (delegate to domain)
- ❌ **NEVER** validate business rules (use domain validators)
- ❌ **NEVER** import from `application/usecase/` or `interface-adapters/`
- ❌ **NEVER** make domain decisions
- ❌ **NEVER** define interfaces in infrastructure (define them in `application/ports/` instead)

**Infrastructure Component Guidelines:**

- **Database**: Connection management, query execution, transaction management
- **Repository**: CRUD operations, convert between domain entities and database models
- **API Clients**: Communicate with external blockchain APIs, handle network errors
- **Storage**: File I/O for transaction and address data
- **Network**: Connection management (WebSocket, gRPC)

#### Infrastructure Layer Structure and Responsibilities

The infrastructure layer is organized into three categories based on their relationship with domain concepts:

##### 1. Domain-Agnostic I/O (Low-Level Technical Implementation)

These packages provide low-level I/O operations without domain awareness:

- **`infrastructure/database/`**: Database connections, SQLC-generated code, query execution
- **`infrastructure/storage/`**: File I/O operations (transaction files, address CSV files)
- **`infrastructure/api/`**: External API clients (Bitcoin, Ethereum, Ripple)

**Characteristics:**

- ✅ **NO interface definitions** - These are concrete implementations only
- ✅ Handle raw data (bytes, strings, database rows)
- ✅ No domain entity conversion
- ✅ Technical concerns only (connection management, error handling, data serialization)

**Example:**

```go
// infrastructure/storage/file/transaction/transaction.go
type TransactionFileRepository struct {
    filePath string
}

func (r *TransactionFileRepository) ReadFile(path string) (string, error) {
    // Raw file I/O - returns string, no domain conversion
    return os.ReadFile(path)
}
```

##### 2. Domain-Aware I/O (Repository Layer)

These packages provide domain-aware data access with conversion responsibilities:

- **`infrastructure/repository/*`**: Repository implementations that convert between infrastructure types and domain entities

**Characteristics:**

- ✅ **NO interface definitions** - Interfaces are defined in `application/ports/`
- ✅ **Depend on** domain-agnostic I/O packages (`database`, `storage`, `api`)
- ✅ **Convert** between infrastructure types (sqlcgen, API responses) and domain entities
- ✅ **Implement** interfaces defined in `application/ports/persistence`
- ✅ Handle domain-specific error mapping

**Dependency Flow:**

```text
infrastructure/repository/*
    ↓ depends on
infrastructure/database/  (or storage/, api/)
    ↓ uses
sqlcgen.* (or file I/O, API responses)
    ↓ converts to
domain.* entities
```

**Example:**

```go
// infrastructure/repository/cold/nonce_repository_sqlc.go
type NonceRepositorySqlc struct {
    queries *sqlcgen.Queries  // Depends on database layer
}

// Implements application/ports/persistence interface
// Converts sqlcgen type to domain entity
func (r *NonceRepositorySqlc) GetNonce(...) (*multisig.NonceCommitment, error) {
    nonce, err := r.queries.GetNonceBySignerAndTx(...)  // Infrastructure type
    if err != nil {
        return nil, err
    }
    return convertToNonceCommitment(&nonce)  // Convert to domain entity
}
```

##### 3. Interface Definitions (Application Ports)

All repository and storage interfaces are defined in the application layer:

- **`application/ports/persistence/`**: Repository interfaces for data persistence
- **`application/ports/storage/`**: File storage interfaces
- **`application/ports/btc/`**: Blockchain API interfaces

**Characteristics:**

- ✅ **All interfaces** are defined here, NOT in infrastructure
- ✅ Infrastructure packages implement these interfaces
- ✅ Use cases depend on these interfaces, not concrete implementations
- ✅ May use domain entities or DTOs in method signatures

**Example:**

```go
// application/ports/persistence/repository.go
type AddressRepositorier interface {
    GetAll(accountType domainAccount.AccountType) ([]*sqlcgen.Address, error)
}

// application/usecase/watch/address.go
type addressUseCase struct {
    addressRepo persistence.AddressRepositorier  // Depends on interface
}
```

##### Design Principles

1. **Separation of Concerns:**
   - Domain-agnostic I/O: Raw data handling
   - Domain-aware I/O: Domain entity conversion
   - Interface definitions: Application layer abstraction

2. **Dependency Direction:**

   ```text
   application/ports/* (interfaces)
       ↑ implements
   infrastructure/repository/* (domain-aware)
       ↑ depends on
   infrastructure/database/ (domain-agnostic)
   ```

3. **No Interface Definitions in Infrastructure:**
   - ❌ **NEVER** define interfaces in `infrastructure/`
   - ✅ **ALWAYS** define interfaces in `application/ports/`
   - ✅ Infrastructure packages only contain implementations

##### Migration Notes

- `infrastructure/repository/cold/interfaces.go` and
  `infrastructure/repository/watch/interfaces.go` contain type aliases for backward compatibility
- These are aliases to `application/ports/persistence` interfaces
- New code should import interfaces directly from `application/ports/`

##### Repository Pattern: Converting Between Infrastructure and Domain Types

The repository layer follows a consistent pattern for converting between infrastructure types (sqlcgen, API responses) and domain entities:

**Key Principles:**

1. **Private conversion functions**: Each repository has `convertToXxx` and `convertFromXxx` helper functions
2. **Bidirectional conversion**: Repository can convert infrastructure→domain and domain→infrastructure
3. **Validation during conversion**: Domain constraints are enforced during conversion
4. **Public methods expose only domain entities**: Repository interface methods work exclusively with domain types

**Example Pattern:**

```go
// infrastructure/repository/watch/address_sqlc.go
type AddressRepositorySqlc struct {
    queries      *sqlcgen.Queries  // Infrastructure database access
    coinTypeCode domainCoin.CoinTypeCode
}

// Private conversion: sqlcgen → domain entity
func convertToAddress(sqlcAddr *sqlcgen.Address) (*domainAddress.Address, error) {
    addr := &domainAddress.Address{
        ID:            sqlcAddr.ID,
        CoinTypeCode:  domainCoin.CoinTypeCode(sqlcAddr.Coin),
        AccountType:   domainAccount.AccountType(sqlcAddr.Account),
        WalletAddress: sqlcAddr.WalletAddress,
        IsAllocated:   sqlcAddr.IsAllocated,
    }

    // Parse timestamps if present
    if sqlcAddr.UpdatedAt.Valid {
        addr.UpdatedAt = &sqlcAddr.UpdatedAt.Time
    }

    return addr, nil
}

// Private conversion: domain entity → sqlcgen
func convertFromAddress(addr *domainAddress.Address) *sqlcgen.Address {
    sqlcAddr := &sqlcgen.Address{
        ID:            addr.ID,
        Coin:          sqlcgen.AddressCoin(addr.CoinTypeCode.String()),
        Account:       sqlcgen.AddressAccount(addr.AccountType.String()),
        WalletAddress: addr.WalletAddress,
        IsAllocated:   addr.IsAllocated,
    }

    if addr.UpdatedAt != nil {
        sqlcAddr.UpdatedAt = sql.NullTime{Time: *addr.UpdatedAt, Valid: true}
    }

    return sqlcAddr
}

// Public method - uses domain entities only
func (r *AddressRepositorySqlc) GetOneUnAllocated(
    accountType domainAccount.AccountType,
) (*domainAddress.Address, error) {
    // Call infrastructure (returns sqlcgen type)
    addr, err := r.queries.GetOneUnallocatedAddress(
        context.Background(),
        sqlcgen.GetOneUnallocatedAddressParams{
            Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
            Account: sqlcgen.AddressAccount(accountType.String()),
        },
    )
    if err != nil {
        return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
    }

    // Convert to domain entity before returning
    return convertToAddress(&addr)
}
```

**Benefits of This Pattern:**

- ✅ **Clean separation**: Infrastructure details don't leak into application layer
- ✅ **Type safety**: Domain constraints are enforced at repository boundary
- ✅ **Testability**: Domain entities are easy to construct in tests
- ✅ **Maintainability**: Changes to database schema only affect repository layer
- ✅ **Flexibility**: Can easily switch database implementations without changing use cases

**Common Conversion Patterns:**

1. **Enum types**: Convert strings to domain enum types with validation
   ```go
   accountType := domainAccount.AccountType(sqlcRow.Account)
   ```

2. **Nullable types**: Handle SQL NULL values properly
   ```go
   if sqlcRow.UpdatedAt.Valid {
       entity.UpdatedAt = &sqlcRow.UpdatedAt.Time
   }
   ```

3. **Numeric types**: Convert between int/int64/uint64 as needed
   ```go
   sequence := uint64(sqlcRow.Sequence)
   ```

4. **Validation**: Use domain validators during conversion
   ```go
   txType, err := domainTx.TxTypeFromInt8(sqlcRow.TxType)
   if err != nil {
       return nil, fmt.Errorf("invalid tx type: %w", err)
   }
   ```

**Anti-Patterns to Avoid:**

- ❌ **DON'T** expose sqlcgen types in repository interface methods
- ❌ **DON'T** skip validation during conversion
- ❌ **DON'T** use infrastructure types in use cases
- ❌ **DON'T** create public conversion functions (keep them private to repository)

### 4. Interface Adapters Layer (`internal/interface-adapters/`)

**Purpose**: Adapters between use cases and external interfaces

**Key Principles:**

- Converts between external formats and application DTOs
- Depends ONLY on application layer (use cases)
- Provides adapters for CLI, HTTP, and wallet interfaces

**Structure:**

- `cli/`: CLI command adapters (keygen, sign, watch)
- `wallet/`: Wallet adapter interfaces and implementations
- `http/`: HTTP handlers (if applicable)

**When Working in Interface Adapters Layer:**

- ✅ Call use cases from commands
- ✅ Convert CLI arguments to use case DTOs
- ✅ Convert use case outputs to CLI output
- ✅ Handle user-facing errors and formatting
- ❌ **NEVER** import from `infrastructure/` or `domain/` directly
- ❌ **NEVER** contain business logic (delegate to use cases)
- ❌ **NEVER** call services directly (use use cases)

**Command Pattern:**

```go
type CreateCommand struct {
    useCase CreateUseCase  // Application layer use case
}

func (c *CreateCommand) Run(args []string) error {
    // Convert CLI args to use case input
    input := CreateInput{
        Param1: args[0],
        Param2: args[1],
    }

    // Call use case
    output, err := c.useCase.Execute(context.Background(), input)
    if err != nil {
        return fmt.Errorf("failed to create: %w", err)
    }

    // Format output for CLI
    fmt.Printf("Created: %s\n", output.Result)
    return nil
}
```

## Dependency Direction Rules

**CRITICAL**: Dependencies must flow in ONE direction only:

```text
Interface Adapters → Application Layer → Domain Layer ← Infrastructure Layer
```

**Rules:**

1. **Domain Layer**: NO dependencies on other layers
2. **Application Layer**: Can depend on Domain Layer only
3. **Infrastructure Layer**: Can depend on Domain Layer only (implements domain interfaces)
4. **Interface Adapters**: Can depend on Application Layer only

**Violations to Avoid:**

- ❌ Domain importing from infrastructure
- ❌ Application importing from infrastructure directly (use interfaces)
- ❌ Application importing from interface-adapters
- ❌ Infrastructure importing from application
- ❌ Interface-adapters importing from infrastructure or domain directly

## `internal/` vs `pkg/` Directory

**`internal/` Directory:**

- Application-specific code
- Follows Clean Architecture layers
- NOT intended for external use
- Can import from `pkg/` (shared utilities)

**`pkg/` Directory:**

- Shared, reusable packages
- Can be imported by external code
- **MUST NOT** import from `internal/` (critical rule)
- Contains utilities (logger, config, converter, etc.)

**When to Use Each:**

- Use `internal/` for application-specific business logic
- Use `pkg/` for shared utilities that could be used by external code
- If unsure, prefer `internal/` (can be refactored to `pkg/` later if needed)

## Legacy Code (`internal/wallet/service/`)

**Status**: Legacy/transitional code being refactored to use case layer

**Guidelines:**

- New code should use `internal/application/usecase/` instead
- Legacy services are being gradually replaced by use cases
- When refactoring, move logic to appropriate layer:
  - Business logic → `domain/`
  - Orchestration → `application/usecase/`
  - Technical implementation → `infrastructure/`

## Dependency Injection (`internal/di/`)

**Purpose**: Wire up dependencies between layers

**Guidelines:**

- `panic` is allowed here (instance construction phase)
- Wire use cases, services, repositories, and API clients
- Follow dependency direction rules
- Keep wiring logic separate from business logic

## Common Patterns

### Adding New Business Logic

1. **Check if it belongs in Domain Layer** (business rules, validators)
2. **If orchestration needed**, create a use case in Application Layer
3. **If external system interaction**, implement in Infrastructure Layer
4. **If user interface**, create adapter in Interface Adapters Layer

### Adding New Cryptocurrency Support

1. **Domain**: Add coin type, validators, business rules
2. **Application**: Add use cases for new coin
3. **Infrastructure**: Add API client, repository implementations
4. **Interface Adapters**: Add CLI commands, wallet adapters

### Refactoring Legacy Code

1. Identify the layer the code belongs to
2. Move business logic to domain layer
3. Create use cases in application layer
4. Move technical implementation to infrastructure layer
5. Update interface adapters to use new use cases
6. Update dependency injection

## Testing Strategy

**Domain Layer:**

- Pure unit tests without mocks
- Test business logic in isolation
- Fast, deterministic tests

**Application Layer:**

- Test use cases with mocked infrastructure
- Verify orchestration logic
- Test error handling and DTO transformation

**Infrastructure Layer:**

- Unit tests with mocked external dependencies
- Integration tests with real external systems
- Test error handling and retries

**Interface Adapters Layer:**

- Test command parsing and output formatting
- Test use case integration (with mocked use cases)
- Test error handling and user-facing messages

## Common Commands

After making code changes in `internal/`, use these commands:

- `make go-lint`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

## Important Notes

- This is a financial-related project; make changes carefully
- Follow Clean Architecture principles strictly
- Maintain dependency direction rules
- Security is critical (private key management, offline wallets)
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)

## References

### Root Documentation

- **[AGENTS.md](../AGENTS.md)**: Agent behavior guidelines and core values
- **[llms.txt](../llms.txt)**: AI-friendly project sitemap
- **[ARCHITECTURE.md](../ARCHITECTURE.md)**: System architecture overview
- **[Core Principles](../docs/ai-agents/guidelines/core.md)**: Security, error handling, panic usage, and core patterns
- **[Architecture Guidelines](../docs/ai-agents/guidelines/architecture.md)**: Clean Architecture principles and layer guidelines
- **[Coding Standards](../docs/ai-agents/guidelines/coding-standards.md)**: Linting, formatting, and code style
- **[Database Management](../docs/ai-agents/guidelines/database.md)**: Database schema changes and SQLC workflow
- **[Code Generation](../docs/ai-agents/guidelines/code-generation.md)**: Auto-generated files and code generation tools
- **[Workflow Guidelines](../docs/ai-agents/guidelines/workflow.md)**: Git operations and dependency management
- **[Testing Guidelines](../docs/ai-agents/guidelines/testing.md)**: Testing strategy and requirements

### Directory Documentation

- **[Pkg Guidelines](../pkg/AGENTS.md)**: Guidelines for `pkg/` directory
- `internal/domain/doc.go`: Domain layer documentation
- `internal/infrastructure/doc.go`: Infrastructure layer documentation
- `internal/application/doc.go`: Application layer documentation
