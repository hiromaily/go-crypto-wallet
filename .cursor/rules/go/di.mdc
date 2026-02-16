---
description: Dependency Injection
globs: 
  - "internal/di/container.go"
  - "pkg/di/container.go"
alwaysApply: false
---


# Dependency Injection

## Overview

Rules for the Dependency Injection (DI) layer in go-crypto-wallet.
The DI layer is responsible for wiring all dependencies and should follow Clean Architecture principles.

## Applicable Files

- `internal/di/container.go` - Application DI container
- `pkg/di/container.go` - Infrastructure/pkg DI container

## Core Principles

### 1. Pure Constructor Pattern (Mandatory)

Constructors must be **pure assembly functions**. They only:

- Receive dependencies as parameters
- Assign parameters to struct fields
- Return the constructed instance

**DO NOT** create internal objects within constructors:

```go
// ✅ GOOD: Pure constructor - just assembles the struct
func NewCreateTransactionUseCase(
    btcClient apibtc.Bitcoiner,
    dbConn *sql.DB,
    addrRepo repowatch.AddressRepositorier,
    txRepo repowatch.BTCTxRepositorier,
) watchusecase.CreateTransactionUseCase {
    return &createTransactionUseCase{
        btcClient: btcClient,
        dbConn:    dbConn,
        addrRepo:  addrRepo,
        txRepo:    txRepo,
    }
}

// ❌ BAD: Constructor creates internal dependencies
func NewCreateTransactionUseCase(
    conf *config.Config,
) watchusecase.CreateTransactionUseCase {
    // DON'T create dependencies inside constructor
    btcClient := bitcoin.NewBitcoin(conf)
    dbConn := mysql.NewConnection(conf.Database)
    return &createTransactionUseCase{
        btcClient: btcClient,
        dbConn:    dbConn,
    }
}
```

### 2. No Factory Pattern for Domain/Application Objects

Factory patterns add unnecessary indirection. Use direct constructor calls instead:

```go
// ✅ GOOD: Direct constructor call
useCase := keygenusecasebtc.NewImportPrivateKeyUseCase(
    btcClient,
    accountKeyRepo,
)

// ❌ BAD: Factory pattern
factory := usecaseFactory.New()
useCase := factory.CreateImportPrivateKeyUseCase()
```

**Exception**: Factory pattern is acceptable only for infrastructure components that require runtime type selection (e.g., key generators based on address type).

### 3. Interface-Based Dependencies

Always depend on interfaces, not concrete types:

```go
// ✅ GOOD: Depends on interface
type createTransactionUseCase struct {
    btcClient apibtc.Bitcoiner           // interface
    addrRepo  repowatch.AddressRepositorier // interface
}

// ❌ BAD: Depends on concrete type
type createTransactionUseCase struct {
    btcClient *apibtcimpl.Bitcoin       // concrete type
    addrRepo  *watchmysql.AddressRepository // concrete type
}
```

### 4. Lazy Initialization in Container

The DI container may use lazy initialization for expensive resources:

```go
// ✅ GOOD: Lazy initialization with caching in container
func (c *container) newBTC() apibtc.Bitcoiner {
    if c.btc == nil {
        var err error
        c.btc, err = bitcoin.NewBitcoin(
            c.newRPCClient(),
            &c.conf.Bitcoin,
            c.conf.CoinTypeCode,
        )
        if err != nil {
            panic(err)
        }
    }
    return c.btc
}
```

This pattern is **only allowed in the DI container**, not in domain/application code.

### 5. Single Responsibility per Constructor

Each constructor creates exactly one type of object:

```go
// ✅ GOOD: One constructor, one responsibility
func NewBTCWatch(...) *BTCWatch { ... }
func NewETHWatch(...) *ETHWatch { ... }

// ❌ BAD: One constructor with branching logic
func NewWatch(coinType string, ...) Watch {
    switch coinType {
    case "BTC":
        return &BTCWatch{...}
    case "ETH":
        return &ETHWatch{...}
    }
}
```

Branching logic belongs in the DI container, not in constructors.

## DI Container Responsibilities

The DI container (`internal/di/container.go`) should:

1. **Hold configuration** - Store config and shared state
2. **Wire dependencies** - Create and connect all objects
3. **Manage lifecycles** - Handle singleton-like instances
4. **Handle coin-type switching** - Branch logic based on coin type

```go
// ✅ GOOD: Coin-type switching in container
func (c *container) NewKeygener() wallets.Keygener {
    switch {
    case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
        return c.newBTCKeygener()
    case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
        return c.newETHKeygener()
    default:
        panic(fmt.Sprintf("coinType[%s] not implemented", c.conf.CoinTypeCode))
    }
}
```

## Naming Conventions

| Pattern | Convention | Example |
|---------|------------|---------|
| Public constructor | `New` + TypeName | `NewBTCWatch` |
| Private container method | `new` + TypeName | `newBTCWalleter` |
| Use case constructor | `New` + UseCaseName | `NewCreateTransactionUseCase` |
| Repository constructor | `New` + RepoName | `NewAddressRepositorySqlc` |

## Dependency Flow

```
cmd/
  └── calls internal/di/container.go
          └── creates Use Cases (internal/application/usecase/)
                └── with Repositories (internal/application/ports/repository/)
                      └── implemented by Infrastructure (internal/infrastructure/repository/)
```

All dependencies flow **inward** toward the domain layer.

## Anti-Patterns to Avoid

| Anti-Pattern | Why It's Bad | Alternative |
|--------------|--------------|-------------|
| Service Locator | Hidden dependencies, hard to test | Constructor injection |
| God Container | Single container does everything | Split by layer/feature |
| Circular Dependencies | Coupling, initialization issues | Refactor interfaces |
| Config in Constructor | Hard to test, coupling | Pass config values directly |

## Testing Implications

Pure constructors enable easy testing:

```go
// Test can inject mocks directly
func TestCreateTransaction(t *testing.T) {
    mockBTC := &mockBitcoiner{}
    mockRepo := &mockAddressRepo{}

    useCase := NewCreateTransactionUseCase(
        mockBTC,
        nil, // dbConn not needed for this test
        mockRepo,
        // ...
    )

    // Test with full control over dependencies
}
```

## Related Rules (in this directory)

- @.claude/rules/go/conventions.md - Go coding conventions
- @.claude/rules/go/usecase.md - Use case patterns
- @.claude/rules/go/repository.md - Repository patterns

## Related Documentation

- @ARCHITECTURE.md - System architecture
- @docs/guidelines/coding-conventions.md - Coding standards
- @internal/AGENTS.md - Internal package guidelines
