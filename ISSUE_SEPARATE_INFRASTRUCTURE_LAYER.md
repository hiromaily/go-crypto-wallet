# Separate Infrastructure Layer from Application Layer

## Overview

This issue implements the "Separate infrastructure layer" task from Phase 3.2 of the refactoring plan. The goal is to clearly separate infrastructure concerns (database, external APIs, file I/O, network connections) from the application layer and organize them according to Clean Architecture principles.

## Background

Currently, infrastructure components are scattered across the codebase and not clearly organized. While some infrastructure exists (API clients, repositories, database connections), the boundaries between infrastructure and application layers are not well-defined, making it difficult to test, maintain, and replace infrastructure components.

### Current Issues

1. **Unclear Boundaries**: Infrastructure components are mixed with application logic in some areas
2. **Inconsistent Organization**: Infrastructure is spread across multiple packages (`pkg/wallet/api/`, `pkg/repository/`, `pkg/db/`, `pkg/tx/file.go`, etc.)
3. **Tight Coupling**: Some infrastructure components have direct dependencies on application-specific types
4. **Hard to Replace**: Infrastructure components are not easily swappable for testing or different implementations
5. **Interface Mismatch**: Repository interfaces in infrastructure packages may not align with domain interfaces

### Example of Current Structure

```go
// pkg/repository/watchrepo/interfaces.go
type AddressRepositorier interface {
    GetAll(accountType account.AccountType) ([]*models.Address, error)
    // Uses infrastructure models, not domain entities
}

// pkg/wallet/api/btcgrp/api-interface.go
type Bitcoiner interface {
    // Infrastructure-specific interface, not aligned with domain
}
```

## Objectives

1. Organize infrastructure components by concern (database, API clients, file I/O, etc.)
2. Ensure infrastructure implements domain interfaces (Dependency Inversion Principle)
3. Create clear boundaries between infrastructure and application layers
4. Make infrastructure easily replaceable and mockable
5. Remove application logic from infrastructure components
6. Maintain backward compatibility during the refactoring process

## Proposed Structure

Create a clear infrastructure layer that implements domain interfaces and is organized by technical concern:

```text
pkg/
├── domain/                    # Domain layer (already separated)
│   ├── wallet/
│   ├── account/
│   ├── transaction/
│   └── ...
├── infrastructure/            # Infrastructure layer (NEW)
│   ├── database/              # Database infrastructure
│   │   ├── mysql/
│   │   │   ├── connection.go  # Database connection management
│   │   │   └── migrations/    # Database migrations (if needed)
│   │   └── sqlc/              # SQL code generation results
│   ├── repository/            # Repository implementations
│   │   ├── watch/             # Watch wallet repositories
│   │   │   ├── address.go     # Implements domain.AddressRepository
│   │   │   ├── transaction.go # Implements domain.TransactionRepository
│   │   │   └── ...
│   │   └── cold/              # Cold wallet repositories
│   │       ├── account_key.go # Implements domain.AccountKeyRepository
│   │       └── ...
│   ├── api/                   # External API clients
│   │   ├── bitcoin/           # Bitcoin/BCH RPC client
│   │   │   ├── client.go      # RPC connection
│   │   │   └── adapter.go     # Adapts to domain interfaces
│   │   ├── ethereum/          # Ethereum JSON-RPC client
│   │   │   ├── client.go
│   │   │   └── adapter.go
│   │   └── ripple/            # XRP gRPC client
│   │       ├── client.go
│   │       └── adapter.go
│   ├── storage/               # File I/O and storage
│   │   ├── file/              # File-based storage
│   │   │   ├── transaction.go # Transaction file storage
│   │   │   └── address.go     # Address file storage
│   │   └── keystore/          # Key storage (if applicable)
│   └── network/               # Network connections
│       ├── websocket/         # WebSocket connections
│       └── grpc/              # gRPC connections
├── uuid/                      # Shared utility (stays in pkg/)
├── logger/                    # Shared utility (stays in pkg/)
├── config/                    # Shared utility (stays in pkg/)
├── converter/                 # Shared utility (stays in pkg/)
├── serial/                    # Shared utility (stays in pkg/)
├── debug/                     # Shared utility (stays in pkg/)
├── di/                        # Dependency injection (stays in pkg/)
├── testutil/                  # Test utilities (stays in pkg/)
├── wallet/
│   ├── key/                   # Cryptographic key generation (shared utility)
│   │   ├── hd_wallet.go        # HD wallet implementation
│   │   └── seed.go             # Seed generation
│   └── service/               # Application layer (depends on domain and infrastructure)
│       ├── btc/
│       ├── eth/
│       └── xrp/
└── contract/                  # Smart contract interactions (shared utility)
```

**Key Requirements:**

- Infrastructure layer (`pkg/infrastructure/`) implements domain interfaces
- Dependency direction: `wallet/service/` → `domain/` ← `infrastructure/`
- Infrastructure is organized by technical concern, not by business domain
- Infrastructure components are easily replaceable (interface-based)
- No business logic in infrastructure layer (only technical implementation)
- **Shared utilities** (uuid, logger, config, etc.) remain in `pkg/` as they are used across all layers

## Implementation Steps

### Step 1: Create Infrastructure Package Structure

- [ ] Create `pkg/infrastructure/` directory as the root for all infrastructure
- [ ] Create subdirectories: `database/`, `repository/`, `api/`, `storage/`, `network/`
- [ ] **Do NOT move shared utilities** (uuid, logger, config, converter, serial, debug, di, testutil, wallet/key, contract) - these stay in `pkg/`
- [ ] Add package-level documentation explaining the infrastructure layer
- [ ] Document that infrastructure implements domain interfaces (Dependency Inversion)
- [ ] Clarify distinction between infrastructure (external systems) and shared utilities (cross-layer helpers)

### Step 2: Organize Database Infrastructure

- [ ] Move `pkg/db/rdb/` → `pkg/infrastructure/database/mysql/`
- [ ] Move SQL code generation results to `pkg/infrastructure/database/sqlc/`
- [ ] Create `pkg/infrastructure/database/mysql/connection.go` for connection management
- [ ] Ensure database layer has no business logic, only connection and query execution
- [ ] Update all imports across the codebase

### Step 3: Reorganize Repository Implementations

- [ ] Move `pkg/repository/watchrepo/` → `pkg/infrastructure/repository/watch/`
- [ ] Move `pkg/repository/coldrepo/` → `pkg/infrastructure/repository/cold/`
- [ ] Update repository implementations to use domain entities (not infrastructure models)
- [ ] Create adapters to convert between domain entities and database models
- [ ] Ensure repositories implement domain repository interfaces (from `pkg/domain/`)
- [ ] Update all imports

### Step 4: Reorganize API Clients

- [ ] Move `pkg/wallet/api/btcgrp/` → `pkg/infrastructure/api/bitcoin/`
- [ ] Move `pkg/wallet/api/ethgrp/` → `pkg/infrastructure/api/ethereum/`
- [ ] Move `pkg/wallet/api/xrpgrp/` → `pkg/infrastructure/api/ripple/`
- [ ] Create adapter layers that implement domain interfaces (if domain defines API interfaces)
- [ ] Separate connection logic from API logic
- [ ] Ensure API clients have no business logic, only API communication
- [ ] Update all imports

### Step 5: Organize File Storage

- [ ] Move `pkg/tx/file.go` → `pkg/infrastructure/storage/file/transaction.go`
- [ ] Move `pkg/address/file.go` → `pkg/infrastructure/storage/file/address.go`
- [ ] Ensure file storage implements domain storage interfaces (if defined)
- [ ] Remove any business logic from file storage (only I/O operations)
- [ ] Update all imports

### Step 6: Organize Network Infrastructure

- [ ] Move `pkg/ws/` → `pkg/infrastructure/network/websocket/`
- [ ] Organize gRPC connection code in `pkg/infrastructure/network/grpc/`
- [ ] Ensure network layer has no business logic
- [ ] Update all imports

### Step 7: Review Shared Utilities (No Changes Needed)

- [ ] **Keep `pkg/uuid/` in place** - General utility used across all layers
- [ ] **Keep `pkg/logger/` in place** - General logging utility used across all layers
- [ ] **Keep `pkg/config/` in place** - Configuration management used by all layers
- [ ] **Keep `pkg/converter/` in place** - General conversion utilities
- [ ] **Keep `pkg/serial/` in place** - General serialization utilities
- [ ] **Keep `pkg/debug/` in place** - Debug utilities
- [ ] **Keep `pkg/di/` in place** - Dependency injection container (application concern)
- [ ] **Keep `pkg/testutil/` in place** - Test utilities
- [ ] **Keep `pkg/wallet/key/` in place** - Cryptographic key generation utility used across all layers
- [ ] **Keep `pkg/contract/` in place** - Smart contract interaction utility used across all layers
- [ ] Document that these are shared utilities, not infrastructure components

### Step 8: Create Adapters for Domain Interfaces

If domain layer defines repository or service interfaces, create adapters:

- [ ] Create adapters in `pkg/infrastructure/repository/` that implement domain interfaces
- [ ] Convert between domain entities and infrastructure models in adapters
- [ ] Ensure adapters are thin (only conversion logic, no business rules)
- [ ] Update application layer to use domain interfaces, not infrastructure interfaces

### Step 9: Update Application Layer

- [ ] Update `pkg/wallet/service/` to depend on domain interfaces, not infrastructure directly
- [ ] Use dependency injection to inject infrastructure implementations
- [ ] Remove direct infrastructure dependencies from application services
- [ ] Ensure application layer only knows about domain interfaces

### Step 10: Update Dependency Injection

- [ ] Update `pkg/di/container.go` to wire infrastructure implementations
- [ ] Ensure container returns domain interfaces, not infrastructure types
- [ ] Update factory functions to create infrastructure and return domain interfaces
- [ ] Test dependency injection works correctly

### Step 11: Update Tests

- [ ] Update test files to use new infrastructure paths
- [ ] Create mock implementations of infrastructure for testing
- [ ] Ensure integration tests still work with reorganized infrastructure
- [ ] Update test utilities in `pkg/testutil/` if needed

### Step 12: Update Documentation

- [ ] Add package documentation to `pkg/infrastructure/`
- [ ] Document infrastructure organization and responsibilities
- [ ] Update `AGENTS.md` with infrastructure layer guidelines
- [ ] Update architecture diagrams if they exist
- [ ] Document how to add new infrastructure components

## Infrastructure Component Guidelines

### Database Infrastructure

- **Responsibility**: Database connections, query execution, transaction management
- **What it does**: Connects to database, executes queries, manages connections
- **What it doesn't do**: Business logic, validation, domain rules
- **Interfaces**: Implements domain repository interfaces

### API Clients

- **Responsibility**: Communication with external blockchain APIs
- **What it does**: Makes HTTP/gRPC calls, handles network errors, serializes/deserializes data
- **What it doesn't do**: Business logic, transaction validation, domain rules
- **Interfaces**: May implement domain service interfaces or be wrapped by adapters

### File Storage

- **Responsibility**: File I/O operations
- **What it does**: Reads/writes files, manages file paths, handles file errors
- **What it doesn't do**: Business logic, file content validation (beyond format)
- **Interfaces**: Implements domain storage interfaces (if defined)

### Network Infrastructure

- **Responsibility**: Network connections (WebSocket, gRPC, HTTP)
- **What it does**: Establishes connections, handles network errors, manages connection lifecycle
- **What it doesn't do**: Business logic, message routing decisions
- **Interfaces**: Provides connection objects to API clients

## Shared Utilities (Stay in `pkg/`)

The following packages are **shared utilities** used across all layers (domain, application, infrastructure) and should **remain in `pkg/`**:

- **`uuid/`**: UUID generation utility
- **`logger/`**: Logging utility used by all layers
- **`config/`**: Configuration management
- **`converter/`**: General conversion utilities
- **`serial/`**: Serialization utilities
- **`debug/`**: Debug utilities
- **`di/`**: Dependency injection container (application concern)
- **`testutil/`**: Test utilities
- **`wallet/key/`**: Cryptographic key generation utility (HD wallet, seed generation)
- **`contract/`**: Smart contract interaction utility

These are not infrastructure components because they:

- Are used by domain, application, and infrastructure layers
- Don't represent external systems or technical implementations
- Are general-purpose utilities that support all layers
- Cryptographic operations (`wallet/key/`, `contract/`) are application-wide utilities, not external system integrations

## Migration Strategy

To maintain backward compatibility:

1. **Gradual Migration**: Move components incrementally, starting with least-used components
2. **Dual Support**: Keep old packages temporarily with deprecation notices
3. **Update Imports**: Update imports package by package
4. **Remove Old Code**: After all imports are updated, remove deprecated packages

## Acceptance Criteria

- [ ] All infrastructure is organized under `pkg/infrastructure/` by technical concern
- [ ] Shared utilities (uuid, logger, config, wallet/key, contract, etc.) remain in `pkg/` (not moved to infrastructure)
- [ ] Infrastructure components implement domain interfaces (Dependency Inversion)
- [ ] No business logic exists in infrastructure layer
- [ ] Application layer depends on domain interfaces, not infrastructure directly
- [ ] Infrastructure is easily replaceable (interface-based)
- [ ] Clear distinction between infrastructure (external systems) and shared utilities (cross-layer helpers)
- [ ] All existing tests pass
- [ ] No breaking changes to public APIs (maintain backward compatibility)
- [ ] Code builds successfully (`make check-build`)
- [ ] Linter passes (`make lint-fix`)
- [ ] Documentation is updated

## Testing Strategy

1. **Unit Tests for Infrastructure**:
   - Test infrastructure components in isolation
   - Mock external dependencies (database, network)
   - Test error handling and edge cases

2. **Integration Tests**:
   - Test infrastructure with real dependencies (database, APIs)
   - Ensure infrastructure integrates correctly with domain and application layers
   - Test adapter conversions between domain entities and infrastructure models

3. **Regression Tests**:
   - Run all existing tests to ensure no functionality is broken
   - Test wallet operations (keygen, sign, watch) end-to-end
   - Verify infrastructure replacement works (e.g., swapping database implementations)

## Dependencies

- This task should be completed after:
  - Phase 3.2: Separate domain logic (already completed)
- This task should be completed before or in parallel with:
  - Phase 3.2: Organize application layer

## Estimated Effort

- **Complexity**: High
- **Time Estimate**: 2-3 weeks
- **Files Affected**: ~80-100 files
- **Risk Level**: Medium (requires careful refactoring to avoid breaking changes)

## Notes

- Follow Clean Architecture principles strictly
- Infrastructure layer should implement domain interfaces (Dependency Inversion)
- Keep infrastructure components focused on technical concerns only
- Make infrastructure easily testable and replaceable
- Refer to `AGENTS.md` for coding standards
- Run `make check-build` and `make lint-fix` after each step
- Consider creating a migration guide for developers

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
- Project's `REFACTORING_PLAN.md` Phase 3.2
- Project's `AGENTS.md` for coding standards
- Project's `ISSUE_SEPARATE_DOMAIN_LOGIC.md` for domain layer structure
