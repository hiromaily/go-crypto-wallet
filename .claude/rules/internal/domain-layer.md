---
paths: ["internal/domain/**/*.go"]
---

# Domain Layer Rules

## Overview

Rules for working in the domain layer (`internal/domain/`).
The domain layer contains **pure business logic with ZERO infrastructure dependencies**.

## Key Principles

1. **NO dependencies** on infrastructure (no database, no API clients, no file I/O)
2. Defines interfaces that infrastructure must implement (Dependency Inversion)
3. Contains pure business logic testable without mocks
4. Most stable layer - changes here affect all other layers

## Structure

| Package | Purpose |
|---------|---------|
| `account/` | Account types, validators, multisig configurations |
| `address/` | Address entity and types |
| `auth/` | Authentication account keys and full public keys |
| `bitcoin/` | Bitcoin-specific types (address, transactions, inputs/outputs) |
| `coin/` | Cryptocurrency type definitions |
| `ethereum/` | Ethereum-specific types (account keys, transactions) |
| `key/` | Key value objects, validators, fingerprint, seed |
| `multisig/` | Multisig validators, MuSig2 implementation |
| `payment/` | Payment request types |
| `transaction/` | Transaction entity, types, validators |
| `wallet/` | Wallet types, descriptor builder |
| `xrp/` | XRP-specific types (constants, account keys) |

## Allowed

- Add business rules and validators
- Define value objects and entities
- Create domain interfaces for infrastructure to implement
- Write pure functions (no side effects)

## Forbidden

- **NEVER** import from `infrastructure/`, `application/`, or `interface-adapters/`
- **NEVER** use database, API clients, or file I/O
- **NEVER** add logging (use `pkg/logger/` only if absolutely necessary)

## Testing

- Use pure unit tests **without mocks**
- Test business logic in isolation
- No infrastructure dependencies required
- Tests should be fast and deterministic

## Example: Domain Interface

```go
// internal/domain/multisig/nonce_repository.go
// Domain defines the interface, infrastructure implements it
type NonceRepositorier interface {
    Store(ctx context.Context, nonce *NonceCommitment) error
    GetBySignerAndTx(ctx context.Context, signerID, txID string) (*NonceCommitment, error)
}
```

## Related Rules

- @.claude/rules/internal/clean-architecture.md - Layer dependency rules
- @.claude/rules/go/conventions.md - Go conventions
