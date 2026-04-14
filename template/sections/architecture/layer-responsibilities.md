## Layer Responsibilities

### Domain Layer (`internal/domain/`)

**Purpose**: Pure business logic with **ZERO infrastructure dependencies**

**Contains**:

- Value Objects (AccountType, TxType, CoinTypeCode)
- Entities (objects with identity and lifecycle)
- Domain Services (stateless business logic)
- Validators (business rule validation)
- Domain Errors (business-specific error types)

**Rules**:

- ❌ NO imports from infrastructure packages
- ❌ NO database, API, or file I/O operations
- ❌ NO external library dependencies (except stdlib)
- ✅ Pure functions preferred
- ✅ Testable without mocks

**Structure**:

```
internal/domain/
├── account/      # Account types, validators, business rules
├── transaction/  # Transaction types, state machine, validators
├── wallet/       # Wallet type definitions
├── key/          # Key value objects and validators
├── multisig/     # Multisig validators and business rules
└── coin/         # Cryptocurrency type definitions
```

### Application Layer (`internal/application/`)

**Purpose**: Orchestrate business operations through use cases

**Contains**:

- **Use Cases** (`usecase/`): Single business operations with clear input/output
- **Ports** (`ports/`): Interface definitions for infrastructure dependencies
- **DTOs** (`dto/`): Data Transfer Objects for port interfaces

**Rules**:

- Use cases depend on port interfaces, not concrete implementations
- Ports define contracts that infrastructure must fulfill
- DTOs use domain types, not infrastructure types

**Structure**:

```
internal/application/
├── usecase/
│   ├── keygen/     # Key generation use cases
│   │   ├── interfaces.go
│   │   ├── btc/
│   │   ├── eth/
│   │   ├── xrp/
│   │   └── shared/
│   ├── sign/       # Signing use cases
│   └── watch/      # Watch wallet use cases
├── ports/
│   ├── bitcoin/    # Bitcoin API interface
│   ├── persistence/# Repository interfaces
│   └── storage/    # File storage interfaces
└── dto/
    └── btc/        # Bitcoin-specific DTOs
```

### Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: Implement interfaces defined by application ports

**Contains**:

- API Clients (Bitcoin Core RPC, Ethereum JSON-RPC, Ripple gRPC)
- Database Repositories (MySQL via SQLC)
- File Storage (address files, transaction files)
- Network Communication (WebSocket clients)

**Rules**:

- ❌ NO interface definitions (only implementations)
- ✅ Implements port interfaces from `application/ports/`
- ✅ Maps infrastructure types to application DTOs

**Structure**:

```
internal/infrastructure/
├── api/
│   ├── bitcoin/    # Bitcoin/BCH RPC clients
│   ├── ethereum/   # Ethereum JSON-RPC clients
│   └── ripple/     # Ripple gRPC clients
├── database/
│   ├── mysql/      # Connection management
│   └── sqlc/       # Generated database code (DO NOT EDIT)
├── repository/
│   ├── cold/       # Cold wallet repository
│   └── watch/      # Watch wallet repository
├── storage/
│   └── file/       # File-based storage
└── wallet/key/     # Key generation infrastructure
```

### Interface Adapters Layer (`internal/interface-adapters/`)

**Purpose**: Adapt between external interfaces and application layer

**Contains**:

- CLI Commands (Cobra-based command implementations)
- Wallet Adapters (wallet-specific implementations)
- HTTP Handlers (if applicable)

**Rules**:

- Commands depend on use cases, NOT services directly
- Convert between external formats and application DTOs

**Structure**:

```
internal/interface-adapters/
├── cli/
│   ├── keygen/     # Keygen wallet commands
│   ├── sign/       # Sign wallet commands
│   └── watch/      # Watch wallet commands
└── wallet/
    ├── interfaces.go  # Wallet interfaces
    ├── btc/           # Bitcoin implementations
    ├── eth/           # Ethereum implementations
    └── xrp/           # XRP implementations
```
