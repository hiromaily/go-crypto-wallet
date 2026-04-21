## Architecture Overview

### Clean Architecture Layers

The PSBT implementation follows Clean Architecture principles:

```
┌─────────────────────────────────────────────────┐
│         Interface Adapters Layer                 │
│  (CLI, Wallet Adapters)                          │
│  - internal/interface-adapters/cli/              │
│  - internal/interface-adapters/wallet/           │
└───────────────┬─────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Application Layer (Use Cases)            │
│  - internal/application/usecase/watch/btc/       │
│  - internal/application/usecase/keygen/btc/      │
│  - internal/application/usecase/sign/btc/        │
└───────────────┬─────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Domain Layer (Business Logic)            │
│  - internal/domain/transaction/                  │
│  - internal/domain/account/                      │
│  - internal/domain/key/                          │
└─────────────────────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Infrastructure Layer                     │
│  - internal/infrastructure/api/btc/btc/      │
│  - internal/infrastructure/storage/file/         │
│  - internal/infrastructure/repository/           │
└─────────────────────────────────────────────────┘
```

### PSBT Flow Through Layers

```
User Command (CLI)
    │
    ▼
Interface Adapter (e.g., watch/btc.BTCWatch)
    │
    ▼
Use Case (e.g., CreateTransactionUseCase)
    │
    ├──> Infrastructure: Bitcoin API (CreatePSBT)
    ├──> Infrastructure: File Storage (WritePSBTFile)
    └──> Infrastructure: Database (InsertTransaction)
```

---

## Design Philosophy

### Why Clean Architecture?

This project handles sensitive financial operations including private key management and cryptocurrency transactions. Clean Architecture provides:

1. **Testability**: Core business logic can be tested without external dependencies
2. **Maintainability**: Clear boundaries make changes predictable and safe
3. **Security**: Separation of concerns helps identify and isolate security-critical code
4. **Flexibility**: Infrastructure can be replaced without affecting business logic

### Core Principles

1. **Dependency Rule**: Dependencies point inward. Outer layers depend on inner layers, never the reverse.
2. **Dependency Inversion**: High-level modules define interfaces; low-level modules implement them.
3. **Single Responsibility**: Each layer has one reason to change.
4. **Interface Segregation**: Interfaces are defined by the layer that uses them, not the layer that implements them.

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Interface Adapters Layer                      │
│                 (CLI Commands, HTTP Handlers)                    │
│                  internal/interface-adapters/                    │
└────────────────────────────┬────────────────────────────────────┘
                             │ depends on
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Application Layer                            │
│            (Use Cases, Ports, DTOs)                              │
│                  internal/application/                           │
└──────────────┬─────────────────────────────────┬────────────────┘
               │ depends on                       │ depends on
               ▼                                  ▼
┌──────────────────────────────┐    ┌─────────────────────────────┐
│        Domain Layer          │    │    Infrastructure Layer     │
│    (Pure Business Logic)     │    │   (External Dependencies)   │
│      internal/domain/        │    │   internal/infrastructure/  │
└──────────────────────────────┘    └──────────────┬──────────────┘
                                                   │ implements
                                                   ▼
                                    ┌─────────────────────────────┐
                                    │   Application Ports         │
                                    │ (Interface Definitions)     │
                                    │  internal/application/ports/│
                                    └─────────────────────────────┘
```

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

## Dependency Injection

The `internal/di/` package wires together all dependencies:

```go
// Example: Watch wallet dependency injection
func NewWatchWallet(cfg *config.WalletConfig) (*WatchWallet, error) {
    // Infrastructure
    bitcoinClient := bitcoin.NewClient(cfg.Bitcoin)
    repository := watchrepo.NewRepository(db)

    // Application (Use Cases)
    createTxUseCase := watch.NewCreateTransactionUseCase(bitcoinClient, repository)

    // Interface Adapters (Commands)
    return &WatchWallet{
        createTxCmd: cli.NewCreateTxCommand(createTxUseCase),
    }, nil
}
```

## Wallet Architecture

### Three Wallet Types

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Watch Wallet  │     │  Keygen Wallet  │     │   Sign Wallet   │
│    (Online)     │     │   (Offline)     │     │   (Offline)     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ • Monitor txs   │     │ • Generate keys │     │ • Auth signing  │
│ • Create unsig  │     │ • Create multis │     │ • Second+ sign  │
│ • Send signed   │     │ • First sign    │     │ • Export pubkey │
│ • Import pubkey │     │ • Export pubkey │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │    CSV/File Export    │    CSV/File Export    │
        └───────────────────────┴───────────────────────┘
```

### Security Model

1. **Keygen Wallet** (Offline): Generates HD wallet seeds and keys. Never connects to network.
2. **Sign Wallet** (Offline): Provides authorization signatures. Each operator has own instance.
3. **Watch Wallet** (Online): Only stores public keys. Cannot sign transactions.

## Data Flow Examples

### Creating and Signing a Transaction

```
1. Watch Wallet (Online)
   └── Create unsigned transaction
   └── Export to file

2. Keygen Wallet (Offline)
   └── Import unsigned transaction
   └── Sign (first signature for multisig)
   └── Export partially signed transaction

3. Sign Wallet (Offline)
   └── Import partially signed transaction
   └── Sign (additional signatures)
   └── Export fully signed transaction

4. Watch Wallet (Online)
   └── Import signed transaction
   └── Broadcast to network
```

## Shared Packages (`pkg/`)

The `pkg/` directory contains reusable utilities that can be used across the application:

```
pkg/
├── config/       # Configuration loading utilities
├── logger/       # Structured logging (slog-based)
├── db/mysql/     # Database connection utilities
├── grpc/         # gRPC client utilities
├── websocket/    # WebSocket client utilities
├── uuid/         # UUID generation
├── decimal/      # Decimal number handling
└── retry/        # Retry utilities
```

**Critical Rule**: Packages in `pkg/` MUST NOT import from `internal/`.

## Testing Strategy

### Unit Tests

Location: `internal/application/usecase/*/btc/*_test.go`

**Current Approach:**

- Constructor tests verify use case instantiation
- Interface compliance tests verify correct interface implementation

**Example:**

```go
func TestNewSignTransactionUseCase(t *testing.T) {
    t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
        useCase := btc.NewSignTransactionUseCase(
            nil, // btc
            nil, // accountKeyRepo
            nil, // txFileRepo
            nil, // multisigAccount
            domainWallet.WalletTypeKeygen,
            "auth1",
        )
        assert.NotNil(t, useCase)
    })

    t.Run("returns correct interface type", func(t *testing.T) {
        useCase := btc.NewSignTransactionUseCase(...)
        assert.Implements(t, (*keygusecase.SignTransactionUseCase)(nil), useCase)
    })
}
```

### Integration Tests

**Requirements for Full Integration Tests:**

1. **Mock Bitcoin Client**
   - CreatePSBT
   - SignPSBTWithKey
   - FinalizePSBT
   - ExtractTransaction
   - IsPSBTComplete

2. **Mock Repositories**
   - TransactionFileRepository (read/write PSBT)
   - AccountKeyRepository (get keys)
   - AuthKeyRepository (get auth keys)
   - BTCTxRepository (database operations)

3. **Test Fixtures**
   - Sample PSBTs (unsigned, partially signed, fully signed)
   - Sample private keys (WIF format)
   - Sample transaction data

**Example Integration Test:**

```go
func TestSignTransactionUseCase_Integration(t *testing.T) {
    // Setup mocks
    mockBTC := &mockBitcoinClient{}
    mockKeyRepo := &mockAccountKeyRepository{}
    mockFileRepo := &mockTransactionFileRepository{}

    // Create use case
    useCase := btc.NewSignTransactionUseCase(
        mockBTC,
        mockKeyRepo,
        mockFileRepo,
        nil,
        domainWallet.WalletTypeKeygen,
        "auth1",
    )

    // Setup test data
    unsignedPSBT := loadTestPSBT("testdata/unsigned.psbt")
    mockFileRepo.On("ReadPSBTFile", mock.Anything).Return(unsignedPSBT, nil)
    mockKeyRepo.On("GetAll", mock.Anything, mock.Anything).Return(testKeys, nil)
    mockBTC.On("SignPSBTWithKey", mock.Anything, mock.Anything).Return(signedPSBT, true, nil)

    // Execute
    output, err := useCase.Sign(context.Background(), input)

    // Assert
    assert.NoError(t, err)
    assert.True(t, output.IsComplete)
    assert.NotEmpty(t, output.SignedData)
}
```

### End-to-End Tests

**Manual E2E Test on Testnet:**

```bash
# 1. Create unsigned PSBT
./watch create deposit --fee 0.00001

# 2. Sign with Keygen
./keygen sign --file deposit_*_unsigned_0_*.psbt

# 3. Broadcast
./watch send --file deposit_*_signed_1_*.psbt

# 4. Verify on blockchain
bitcoin-cli -testnet getrawtransaction <txid> 1
```

**Automated E2E Tests:**

See `docs/TESTING_STRATEGY.md` for comprehensive testing approach.

---

## See Also

- [AGENTS.md](./AGENTS.md) - Agent behavior guidelines
- [llms.txt](./llms.txt) - AI-friendly project overview
- [docs/guidelines/architecture.md](./docs/guidelines/architecture.md) - Detailed architecture guidelines
- [internal/AGENTS.md](./internal/AGENTS.md) - Internal package guidelines
- [pkg/AGENTS.md](./pkg/AGENTS.md) - Public package guidelines
