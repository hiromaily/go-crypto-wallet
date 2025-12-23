# Separate Domain Logic from Application and Infrastructure Layers

## Overview

This issue implements the "Separate domain logic" task from Phase 3.2 of the refactoring plan. The goal is to extract pure domain logic (business rules, domain models, and domain services) from the current service layer and organize it according to Clean Architecture principles.

## Background

Currently, domain logic is mixed with application logic and infrastructure concerns in `pkg/wallet/service/`. This violates Clean Architecture principles and makes the codebase harder to maintain, test, and evolve.

### Current Issues

1. **Mixed Concerns**: Business rules, validation logic, and domain models are scattered across service files that also contain infrastructure dependencies (database connections, API clients, repositories)
2. **Tight Coupling**: Domain logic directly depends on infrastructure (e.g., `*sql.DB`, repository interfaces, API clients)
3. **Hard to Test**: Domain logic cannot be tested in isolation without mocking infrastructure
4. **Unclear Boundaries**: It's difficult to distinguish between domain rules and application orchestration

### Example of Current Structure

```go
// pkg/wallet/service/btc/watchsrv/tx_creator.go
type TxCreate struct {
    btc             btcgrp.Bitcoiner  // Infrastructure
    dbConn          *sql.DB           // Infrastructure
    addrRepo        watchrepo.AddressRepositorier  // Infrastructure
    // ... business logic mixed with infrastructure
}
```

## Objectives

1. Create a dedicated `pkg/domain/` package for pure domain logic
2. Extract domain entities, value objects, and business rules
3. Remove infrastructure dependencies from domain logic
4. Define clear domain interfaces that infrastructure will implement
5. Maintain backward compatibility during the refactoring process

## Proposed Structure

Create the domain layer as `pkg/domain/` following Clean Architecture principles. The domain layer must be independent with no infrastructure dependencies. The application layer (`pkg/wallet/service/`) will depend on the domain layer.

```text
pkg/
├── domain/                    # Pure domain logic (no infrastructure dependencies)
│   ├── wallet/
│   │   ├── entity.go          # Wallet entity
│   │   ├── types.go           # WalletType, etc.
│   │   └── service.go         # Domain services (pure business logic)
│   ├── account/
│   │   ├── entity.go          # Account entity
│   │   ├── types.go           # AccountType, AuthType (move from pkg/account/)
│   │   └── validator.go       # Account validation rules
│   ├── transaction/
│   │   ├── entity.go          # Transaction entity
│   │   ├── types.go           # TxType (move from pkg/tx/)
│   │   ├── validator.go       # Transaction validation rules
│   │   └── service.go         # Transaction domain services
│   ├── key/
│   │   ├── entity.go          # Key entity
│   │   ├── valueobject.go     # WalletKey, Address, etc.
│   │   └── service.go         # Key generation domain logic
│   ├── multisig/
│   │   ├── entity.go          # MultisigAddress entity
│   │   ├── validator.go       # Multisig validation rules
│   │   └── service.go         # Multisig domain services
│   └── coin/
│       ├── entity.go          # Coin entity
│       └── types.go           # CoinTypeCode (move from pkg/wallet/coin/)
├── wallet/
│   └── service/               # Application layer (depends on domain)
│       ├── btc/
│       ├── eth/
│       └── xrp/
└── ...
```

**Key Requirements:**

- Domain layer (`pkg/domain/`) must be independent with zero infrastructure dependencies
- Dependency direction: `wallet/service/` → `domain/` (Application depends on Domain)
- Domain logic spans multiple features (wallet, account, transaction, key, multisig, coin)
- Domain layer must be testable without any infrastructure (pure functions where possible)

## Implementation Steps

### Step 1: Create Domain Package Structure

- [ ] Create `pkg/domain/` directory (independent domain layer, not nested in service)
- [ ] Create subdirectories for each domain: `wallet/`, `account/`, `transaction/`, `key/`, `multisig/`, `coin/`
- [ ] Add package-level documentation explaining the domain layer and its independence from infrastructure
- [ ] Document the dependency direction: Application layer (`pkg/wallet/service/`) depends on Domain layer (`pkg/domain/`)

### Step 2: Extract Domain Types and Value Objects

- [ ] Move `pkg/wallet/types.go` → `pkg/domain/wallet/types.go`
- [ ] Move `pkg/account/types.go` → `pkg/domain/account/types.go` (keep utility functions)
- [ ] Move `pkg/tx/types.go` → `pkg/domain/transaction/types.go`
- [ ] Move `pkg/wallet/coin/types.go` → `pkg/domain/coin/types.go`
- [ ] Extract value objects from `pkg/wallet/key/key.go` → `pkg/domain/key/valueobject.go`
- [ ] Update all imports across the codebase

### Step 3: Create Domain Entities

- [ ] Create `pkg/domain/transaction/entity.go` with `Transaction` entity (business rules only, no DB fields)
- [ ] Create `pkg/domain/account/entity.go` with `Account` entity
- [ ] Create `pkg/domain/key/entity.go` with `Key` entity
- [ ] Create `pkg/domain/multisig/entity.go` with `MultisigAddress` entity
- [ ] Ensure entities contain only business logic and validation, no infrastructure concerns

### Step 4: Extract Domain Business Rules

Identify and extract pure business logic from service files:

#### Transaction Domain Rules

- [ ] Extract transaction validation rules from:
  - `pkg/wallet/service/btc/watchsrv/tx_creator.go`
  - `pkg/wallet/service/eth/watchsrv/tx_creator.go`
  - `pkg/wallet/service/xrp/watchsrv/tx_creator.go`
- [ ] Create `pkg/domain/transaction/validator.go` with validation functions:
  - `ValidateAmount(amount, balance) error`
  - `ValidateSenderReceiver(sender, receiver AccountType) error`
  - `ValidateTransactionType(txType TxType) error`
- [ ] Create `pkg/domain/transaction/service.go` for transaction domain services:
  - `CalculateFee(amount, feeRate) Amount`
  - `CalculateChange(inputTotal, requiredAmount, fee) Amount`

#### Account Domain Rules

- [ ] Extract account validation from service files
- [ ] Create `pkg/domain/account/validator.go`:
  - `ValidateAccountType(accountType AccountType) error`
  - `CanTransferFrom(accountType AccountType) bool`
  - `CanReceiveTo(accountType AccountType) bool`

#### Multisig Domain Rules

- [ ] Extract multisig validation from `pkg/wallet/service/btc/coldsrv/keygensrv/multisigaddress.go`
- [ ] Create `pkg/domain/multisig/validator.go`:
  - `ValidateMultisigConfig(requiredSigs, totalSigs int) error`
  - `ValidateRedeemScript(redeemScript string) error`
- [ ] Create `pkg/domain/multisig/service.go`:
  - `GenerateMultisigAddress(publicKeys []string, requiredSigs int) (address, redeemScript string, error)`

#### Key Domain Rules

- [ ] Extract key generation business rules from `pkg/wallet/service/coldsrv/hd_walleter.go`
- [ ] Create `pkg/domain/key/service.go`:
  - `DeriveKeyPath(accountType AccountType, index uint32) string`
  - `ValidateKeyIndex(index uint32) error`

### Step 5: Define Domain Interfaces

Create repository interfaces in the domain layer (not in infrastructure):

- [ ] Create `pkg/domain/transaction/repository.go` with `TransactionRepository` interface
- [ ] Create `pkg/domain/account/repository.go` with `AccountRepository` interface
- [ ] Create `pkg/domain/key/repository.go` with `KeyRepository` interface
- [ ] These interfaces should use domain entities, not infrastructure models

### Step 6: Refactor Service Layer

Update service layer to use domain logic:

- [ ] Update `pkg/wallet/service/btc/watchsrv/tx_creator.go`:
  - Remove business rule logic
  - Call domain services and validators
  - Keep only orchestration and infrastructure coordination
- [ ] Update `pkg/wallet/service/eth/watchsrv/tx_creator.go` similarly
- [ ] Update `pkg/wallet/service/xrp/watchsrv/tx_creator.go` similarly
- [ ] Update `pkg/wallet/service/coldsrv/hd_walleter.go`:
  - Use domain key service for business logic
  - Keep only repository coordination
- [ ] Update `pkg/wallet/service/btc/coldsrv/keygensrv/multisigaddress.go`:
  - Use domain multisig service
  - Keep only infrastructure calls

### Step 7: Update Infrastructure Layer

- [ ] Ensure repository implementations in `pkg/repository/` implement domain interfaces
- [ ] Add adapters if needed to convert between domain entities and infrastructure models
- [ ] Update API clients in `pkg/wallet/api/` to work with domain entities where appropriate

### Step 8: Update Tests

- [ ] Create unit tests for domain logic (no infrastructure mocks needed)
- [ ] Update existing tests to use domain entities
- [ ] Ensure domain tests are fast and don't require database/network

### Step 9: Update Documentation

- [ ] Add package documentation to `pkg/domain/`
- [ ] Document domain entities and their business rules
- [ ] Update `AGENTS.md` with domain layer guidelines
- [ ] Update architecture diagrams if they exist

## Domain Logic Identification Guide

When extracting domain logic, look for:

1. **Business Rules**:
   - "Sender balance must be sufficient"
   - "Multisig requires at least 2 signatures"
   - "Account type X cannot receive from account type Y"

2. **Validation Logic**:
   - Input validation
   - State validation
   - Business rule enforcement

3. **Domain Calculations**:
   - Fee calculations
   - Change calculations
   - Balance validations

4. **Domain Models**:
   - Entities with identity and lifecycle
   - Value objects (immutable)
   - Domain events (if applicable)

**What NOT to include in domain layer:**

- Database operations
- API calls to external services
- File I/O
- Logging (domain should be silent, application layer logs)
- Configuration access
- Repository implementations

## Acceptance Criteria

- [ ] All domain logic is in `pkg/domain/` with no infrastructure dependencies
- [ ] Domain packages can be imported without importing infrastructure packages
- [ ] Domain logic is testable without mocks (pure functions where possible)
- [ ] Service layer uses domain services and validators
- [ ] All existing tests pass
- [ ] No breaking changes to public APIs (maintain backward compatibility)
- [ ] Code builds successfully (`make check-build`)
- [ ] Linter passes (`make lint-fix`)
- [ ] Documentation is updated

## Testing Strategy

1. **Unit Tests for Domain**:
   - Test domain entities and value objects
   - Test domain services (pure functions)
   - Test validators
   - No infrastructure dependencies

2. **Integration Tests**:
   - Test service layer with domain logic
   - Ensure domain logic integrates correctly with infrastructure

3. **Regression Tests**:
   - Run all existing tests to ensure no functionality is broken
   - Test wallet operations (keygen, sign, watch) end-to-end

## Migration Strategy

To maintain backward compatibility:

1. **Gradual Migration**: Move code incrementally, starting with types and value objects
2. **Dual Support**: Keep old packages temporarily, mark as deprecated
3. **Update Imports**: Update imports package by package
4. **Remove Old Code**: After all imports are updated, remove deprecated code

## Dependencies

- This task should be completed after or in parallel with:
  - Phase 2.1: Error Handling Improvements (for proper error types in domain)
  - Phase 3.1: Dependency Injection Improvements (for clean interfaces)

## Estimated Effort

- **Complexity**: High
- **Time Estimate**: 2-3 weeks
- **Files Affected**: ~50-70 files
- **Risk Level**: Medium (requires careful refactoring to avoid breaking changes)

## Notes

- Follow Clean Architecture principles strictly
- Domain layer should have zero dependencies on infrastructure
- Use dependency inversion: domain defines interfaces, infrastructure implements them
- Keep domain logic pure and testable
- Refer to `AGENTS.md` for coding standards
- Run `make check-build` and `make lint-fix` after each step

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- Project's `REFACTORING_PLAN.md` Phase 3.2
- Project's `AGENTS.md` for coding standards
