---
paths: ["internal/infrastructure/**/*.go"]
---

# Infrastructure Layer Rules

## Overview

Rules for working with the Infrastructure layer (`internal/infrastructure/`).
This layer contains **implementations only** - no interface definitions.

## Structure

```
internal/infrastructure/
├── api/            # External API clients
│   ├── btc/        # Bitcoin (btc/, bch/, mocks/, testutil/)
│   ├── eth/        # Ethereum (eth/, erc20/, ethtx/, testutil/)
│   └── xrp/        # Ripple (xrp/, mocks/, testutil/)
├── contract/       # Smart contract utilities
├── database/       # SQLC-generated code (DO NOT EDIT)
│   ├── mysql/sqlcgen/
│   └── sqlite/sqlcgen/
├── repository/     # Repository implementations
│   ├── cold/       # Cold wallet (mysql/, sqlite/, mocks/)
│   └── watch/      # Watch wallet (mysql/, sqlite/, mocks/)
├── storage/        # File storage implementations
│   └── file/       # address/, descriptor/, fullpubkey/, transaction/
└── wallet/key/     # Key generation (BIP32, BIP44, BIP84, BIP86, HD wallet)
```

## Critical Principle

**Infrastructure contains NO interface definitions.**

All interfaces are defined in `application/ports/`. Infrastructure packages only implement them.

## Auto-Generated Files (DO NOT EDIT)

```
internal/infrastructure/database/mysql/sqlcgen/*.go
internal/infrastructure/database/sqlite/sqlcgen/*.go
internal/infrastructure/api/xrp/xrp/*.pb.go
```

## Type Conversion Pattern

Infrastructure must convert between infrastructure types and domain entities:

```go
// Private conversion functions in repository
func convertToAddress(sqlcAddr *sqlcgen.Address) (*domainAddress.Address, error) {
    return &domainAddress.Address{
        ID:            sqlcAddr.ID,
        CoinTypeCode:  domainCoin.CoinTypeCode(sqlcAddr.Coin),
        AccountType:   domainAccount.AccountType(sqlcAddr.Account),
        WalletAddress: sqlcAddr.WalletAddress,
    }, nil
}

// Public method returns domain entity
func (r *AddressRepositorySqlc) GetAll(...) ([]*domainAddress.Address, error) {
    rows, err := r.queries.GetAllAddresses(...)
    if err != nil {
        return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
    }
    // Convert each row to domain entity
    return convertAll(rows, convertToAddress)
}
```

## Allowed/Forbidden

| Action | Allowed |
|--------|---------|
| Import from `domain/` | Yes |
| Import from `application/ports/` | Yes (to implement interfaces) |
| Import from `application/dto/` | Yes (for API type mapping) |
| Import from `application/usecase/` | No |
| Import from `interface-adapters/` | No |
| Define interfaces | No (define in `application/ports/`) |
| Contain business logic | No (delegate to domain) |

## Implementation Guidelines

### Repository Implementations

- Implement interfaces from `application/ports/repository/`
- Use private `convertToXxx` and `convertFromXxx` functions
- Return domain entities, not sqlcgen types

### API Client Implementations

- Implement interfaces from `application/ports/api/`
- Map API responses to application DTOs
- Handle network errors and retries

### Storage Implementations

- Implement interfaces from `application/ports/file/`
- Handle file I/O errors properly

## Testing

- Use mocks in `mocks/` directories for unit tests — see @.claude/rules/internal/mockery.md
- Use `testutil/` for test helpers
- Integration tests with real systems when needed

### `testutil/` Files During Interface Refactoring

**Rule**: When renaming or splitting an infrastructure client interface (e.g. monolithic `XRPer` →
`XRPPublicClient`), always check `testutil/` subdirectories in the same package tree.
These files often hold a package-level `var` and factory function referencing the old interface
type and constructor, and they are **not** caught at compilation time until `make go-lint` runs.

Checklist when refactoring an infrastructure API client:
1. Search for `testutil/` under `internal/infrastructure/api/<chain>/`
2. Update the package-level `var` type (e.g. `var xr apixrp.XRPer` → `var xr apixrp.XRPPublicClient`)
3. Update the factory function return type and constructor call
4. Remove any now-unused arguments (e.g. the admin WebSocket client if testutil only needs public)

## Mock Generation Rule

Every interface defined under `internal/application/ports/` MUST have a mockery-generated mock
in the corresponding `mocks/` subdirectory. See @.claude/rules/internal/mockery.md for:

- Placement convention (ports package → infrastructure mocks directory)
- How to add a new interface to `.mockery.yaml`
- Test usage pattern (`NewMock*(t)` + `.EXPECT()`)
- Exceptions (`Ethereumer`, type aliases, use case interfaces)

## Related Rules

- @.claude/rules/internal/clean-architecture.md - Layer dependencies
- @.claude/rules/internal/mockery.md - Mock generation rules (mockery)
- @.claude/rules/go/repository.md - Repository pattern details
- @.claude/rules/go/di.md - Dependency injection
