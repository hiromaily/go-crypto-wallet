---
description: Claude Rules - Use Case Pattern
globs: 
  - "internal/application/usecase/**/*.go"
alwaysApply: false
---


# Claude Rules - Use Case Pattern

## Overview

Rules for implementing Use Cases in go-crypto-wallet.
Use Cases represent single business operations and are the core of the Application layer.

## Applicable Files

- `internal/application/usecase/**/*.go`

## Use Case Structure

### Interface Definition

Each use case is defined as an interface in the appropriate `interfaces.go` file:

```go
// ✅ GOOD: Interface with clear single responsibility
type CreateTransactionUseCase interface {
    Execute(ctx context.Context, input CreateTransactionInput) (CreateTransactionOutput, error)
}

// ✅ GOOD: Multiple related operations in one interface
type GenerateSeedUseCase interface {
    Generate(ctx context.Context) (GenerateSeedOutput, error)
    Store(ctx context.Context, input StoreSeedInput) (StoreSeedOutput, error)
}
```

### Implementation Structure

```go
// Private struct (unexported)
type createTransactionUseCase struct {
    btcClient  apibtc.Bitcoiner           // Ports only
    addrRepo   repowatch.AddressRepositorier
    txRepo     repowatch.BTCTxRepositorier
    // ... other dependencies
}

// Public constructor
func NewCreateTransactionUseCase(
    btcClient apibtc.Bitcoiner,
    addrRepo repowatch.AddressRepositorier,
    txRepo repowatch.BTCTxRepositorier,
    // ... pass all dependencies
) watchusecase.CreateTransactionUseCase {
    return &createTransactionUseCase{
        btcClient: btcClient,
        addrRepo:  addrRepo,
        txRepo:    txRepo,
    }
}

// Implementation
func (u *createTransactionUseCase) Execute(
    ctx context.Context,
    input watchusecase.CreateTransactionInput,
) (watchusecase.CreateTransactionOutput, error) {
    // Business logic here
}
```

## Method Naming Conventions

| Operation Type | Method Name | Example |
|----------------|-------------|---------|
| Create/Generate | `Execute`, `Generate`, `Create` | `CreateTransactionUseCase.Execute` |
| Read/Get | `Get`, `Find`, `List` | `GetAccountUseCase.Get` |
| Update | `Update`, `Modify` | `UpdateStatusUseCase.Update` |
| Delete | `Delete`, `Remove` | `DeleteAccountUseCase.Delete` |
| Import/Export | `Import`, `Export` | `ImportAddressUseCase.Execute` |
| Sign | `Sign` | `SignTransactionUseCase.Sign` |
| Monitor | `UpdateTxStatus`, `MonitorBalance` | `MonitorTransactionUseCase.UpdateTxStatus` |

## Input/Output DTO Patterns

### Input DTOs

```go
// ✅ GOOD: Clear, focused input
type CreateTransactionInput struct {
    ActionType        string                      // Use domain types when possible
    SenderAccount     domainAccount.AccountType
    ReceiverAccount   domainAccount.AccountType
    Amount            float64
    PaymentRequestIDs []int64
}

// ❌ BAD: Raw types without context
type CreateTransactionInput struct {
    Type      string
    From      string
    To        string
    Amount    string  // Should be float64 or decimal
}
```

### Output DTOs

```go
// ✅ GOOD: Contains result and metadata
type CreateTransactionOutput struct {
    TransactionHex string
    FileName       string
}

type SignTransactionOutput struct {
    FilePath      string
    IsDone        bool
    SignedCount   int
    UnsignedCount int
}
```

## Mandatory Rules

### 1. Context as First Parameter

All use case methods MUST accept `context.Context` as the first parameter:

```go
// ✅ GOOD
func (u *usecase) Execute(ctx context.Context, input Input) (Output, error)

// ❌ BAD
func (u *usecase) Execute(input Input) (Output, error)
```

### 2. Depend on Ports, Not Infrastructure

Use cases depend on interfaces from `internal/application/ports/`:

```go
// ✅ GOOD: Depends on port interfaces
type createTransactionUseCase struct {
    btcClient apibtc.Bitcoiner              // from ports/api/btc
    addrRepo  repowatch.AddressRepositorier // from ports/repository/watch
}

// ❌ BAD: Depends on infrastructure directly
type createTransactionUseCase struct {
    btcClient *apibtcimpl.Bitcoin           // infrastructure implementation
    addrRepo  *watchmysql.AddressRepository // infrastructure implementation
}
```

### 3. Single Responsibility

Each use case handles one business operation:

```go
// ✅ GOOD: Single responsibility
type CreateTransactionUseCase interface {
    Execute(ctx context.Context, input CreateTransactionInput) (CreateTransactionOutput, error)
}

// ❌ BAD: Multiple unrelated responsibilities
type TransactionUseCase interface {
    Create(...)
    Send(...)
    Monitor(...)
    Delete(...)
}
```

### 4. No Business Logic in Interface Adapters

Business logic belongs in use cases, not CLI commands:

```go
// ✅ GOOD: CLI just calls use case
func (c *CreateCommand) Run() error {
    output, err := c.createTxUseCase.Execute(ctx, input)
    if err != nil {
        return err
    }
    fmt.Printf("Created: %s\n", output.FileName)
    return nil
}

// ❌ BAD: Business logic in CLI
func (c *CreateCommand) Run() error {
    // DON'T put validation, calculation, DB access here
    addresses := c.repo.GetAddresses(...)
    balance := calculateBalance(addresses)
    if balance < amount {
        return errors.New("insufficient balance")
    }
    // ...
}
```

## Directory Structure

```
internal/application/usecase/
├── keygen/
│   ├── interfaces.go          # Interface definitions
│   ├── btc/
│   │   ├── create_multisig.go # BTC-specific implementation
│   │   ├── export_address.go
│   │   └── ...
│   ├── eth/
│   │   └── ...
│   ├── xrp/
│   │   └── ...
│   └── shared/
│       └── generate_seed.go   # Shared across all coins
├── sign/
│   ├── interfaces.go
│   └── ...
└── watch/
    ├── interfaces.go
    └── ...
```

## Coin-Specific vs Shared Use Cases

### Shared Use Cases (`shared/` directory)

Use cases that work identically across all coins:

```go
// internal/application/usecase/keygen/shared/generate_seed.go
func NewGenerateSeedUseCase(seedRepo repocold.SeedRepositorier) keygenusecase.GenerateSeedUseCase {
    return &generateSeedUseCase{seedRepo: seedRepo}
}
```

### Coin-Specific Use Cases (`btc/`, `eth/`, `xrp/` directories)

Use cases with coin-specific logic:

```go
// internal/application/usecase/keygen/btc/create_multisig.go
func NewCreateMultisigAddressUseCase(
    btc apibtc.Bitcoiner,
    authFullPubKeyRepo repocold.AuthFullPubkeyRepositorier,
    // BTC-specific dependencies
) keygenusecase.CreateMultisigAddressUseCase {
    // ...
}
```

## Error Handling

Wrap errors with context and operation name:

```go
func (u *createTransactionUseCase) Execute(ctx context.Context, input Input) (Output, error) {
    addresses, err := u.addrRepo.GetAll(ctx)
    if err != nil {
        return Output{}, fmt.Errorf("failed to get addresses: %w", err)
    }

    tx, err := u.btcClient.CreateRawTransaction(addresses, outputs)
    if err != nil {
        return Output{}, fmt.Errorf("failed to create raw transaction: %w", err)
    }

    return Output{TransactionHex: tx.Hex}, nil
}
```

## Testing

Use cases should be tested with mocked ports:

```go
func TestCreateTransaction(t *testing.T) {
    // Arrange
    mockBTC := &mocks.MockBitcoiner{}
    mockAddrRepo := &mocks.MockAddressRepo{}

    usecase := NewCreateTransactionUseCase(mockBTC, mockAddrRepo, ...)

    // Setup expectations
    mockAddrRepo.On("GetAll", mock.Anything).Return(testAddresses, nil)
    mockBTC.On("CreateRawTransaction", mock.Anything, mock.Anything).Return(testTx, nil)

    // Act
    output, err := usecase.Execute(context.Background(), input)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, expectedHex, output.TransactionHex)
}
```

## Related Rules (in this directory)

- @.claude/rules/go/di.md - Dependency injection rules
- @.claude/rules/go/repository.md - Repository patterns
- @.claude/rules/go/conventions.md - Go coding conventions

## Related Documentation

- @ARCHITECTURE.md - Layer responsibilities
- @internal/application/ports/ - Port interface definitions
