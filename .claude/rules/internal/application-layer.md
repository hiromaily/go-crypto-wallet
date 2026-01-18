---
paths: ["internal/application/**/*.go"]
---

# Application Layer Rules

## Overview

Rules for working with the Application layer (`internal/application/`).
This layer orchestrates use cases and defines port interfaces.

## Structure

```
internal/application/
├── dto/           # Data Transfer Objects
│   ├── btc/       # Bitcoin DTOs
│   └── ripple/    # XRP DTOs
├── ports/         # Interface definitions (abstractions)
│   ├── api/       # External API interfaces (btc, eth, xrp)
│   ├── file/      # File storage interfaces
│   ├── repository/# Repository interfaces (cold, watch)
│   └── wallet/    # Wallet key generator interfaces
└── usecase/       # Use case implementations
    ├── keygen/    # Key generation (btc, bch, eth, xrp, shared)
    ├── sign/      # Signing (btc, bch, eth, xrp, shared)
    ├── watch/     # Watch wallet (btc, bch, eth, xrp, shared)
    └── shared/    # Shared helpers
```

## Port Interfaces (`application/ports/`)

**CRITICAL**: Interfaces are defined where they are USED, not where they are implemented.

### Rules

- Define interfaces in `application/ports/` for infrastructure abstractions
- Use application-layer DTOs in method signatures (NOT infrastructure types)
- Infrastructure packages implement these interfaces

### Example

```go
// application/ports/api/btc/interface.go
package btc

import btcdto "internal/application/dto/btc"

type Bitcoiner interface {
    GetBalance() (btcutil.Amount, error)
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)  // Uses DTO
}
```

## DTOs (`application/dto/`)

### Rules

- DTOs define data structures for port interface method signatures
- Keep DTOs simple - no business logic
- Map infrastructure types to DTOs in infrastructure layer

### Example

```go
// application/dto/btc/dto.go
package btc

type AddressInfo struct {
    Address      string
    ScriptPubKey string
    IsWitness    bool
    Labels       []string
}
```

## Use Cases (`application/usecase/`)

### Rules

- Each use case represents a single business operation
- Depend on port interfaces, not concrete implementations
- Delegate business logic to domain layer
- Wrap errors with context

### Pattern

```go
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

type xxxUseCase struct {
    repo     repository.SomeRepositorier  // Port interface
    apiClient api.SomeAPIer               // Port interface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    // 1. Validate input (using domain validators)
    // 2. Call infrastructure via port interfaces
    // 3. Transform and return result
}
```

## Allowed/Forbidden

| Action | Allowed |
|--------|---------|
| Import from `domain/` | Yes |
| Import from `application/ports/` | Yes |
| Import from `application/dto/` | Yes |
| Import from `infrastructure/` | No |
| Import from `interface-adapters/` | No |
| Define interfaces | Yes (in `ports/`) |
| Contain business logic | No (delegate to domain) |

## Related Rules

- @.claude/rules/internal/clean-architecture.md - Layer dependencies
- @.claude/rules/go/usecase.md - Use case patterns
- @.claude/rules/go/repository.md - Repository interfaces
