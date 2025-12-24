# Organize Application Layer

## Overview

This issue implements the "Organize application layer" task from Phase 3.2 of the refactoring plan. The goal is to reorganize the application layer (`pkg/wallet/service/`) to follow Clean Architecture principles with clear organization by use cases, consistent structure, and proper separation of concerns.

## Background

Currently, the application layer has inconsistent organization patterns and mixed concerns. Services are organized partially by coin type (btc/, eth/, xrp/) and partially by wallet type (coldsrv/, watchsrv/), making it difficult to understand the architecture and maintain the codebase.

### Current Issues

1. **Inconsistent Organization**: Services are organized by different criteria:
   - Coin-specific services: `btc/`, `eth/`, `xrp/`
   - Wallet-type services: `coldsrv/`, `watchsrv/`
   - Interfaces at root level: `cold_interface.go`, `watch_interface.go`

2. **Mixed Concerns**: Services combine multiple responsibilities:
   - Transaction creation, monitoring, and sending in the same package
   - Key generation, signing, and import/export mixed together
   - No clear separation between use cases

3. **Unclear Boundaries**: It's difficult to understand:
   - Which services belong to which wallet type (watch, keygen, sign)
   - Which services are coin-specific vs. shared
   - How services relate to each other

4. **Interface-Implementation Separation**: Interfaces are at the root level while implementations are scattered across subdirectories, making it hard to find related code.

### Example of Current Structure

```text
pkg/wallet/service/
├── cold_interface.go          # Interfaces at root
├── watch_interface.go         # Interfaces at root
├── coldsrv/                   # Shared cold wallet services
│   ├── hd_walleter.go
│   ├── seeder.go
│   └── address_exporter.go
├── watchsrv/                  # Shared watch wallet services
│   ├── address_importer.go
│   └── payment_request_creator.go
├── btc/
│   ├── coldsrv/
│   │   ├── signer.go
│   │   ├── keygensrv/
│   │   │   ├── privkey_importer.go
│   │   │   ├── multisigaddress.go
│   │   │   └── fullpubkey_importer.go
│   │   └── signsrv/
│   │       ├── privkey_importer.go
│   │       └── fullpubkey_exporter.go
│   └── watchsrv/
│       ├── tx_creator.go
│       ├── tx_creator_deposit.go
│       ├── tx_creator_payment.go
│       ├── tx_creator_transfer.go
│       ├── tx_monitor.go
│       ├── tx_sender.go
│       └── address_importer.go
├── eth/
│   ├── keygensrv/
│   │   ├── privkey_importer.go
│   │   └── signer.go
│   └── watchsrv/
│       ├── tx_creator.go
│       ├── tx_creator_deposit.go
│       ├── tx_creator_payment.go
│       ├── tx_creator_transfer.go
│       ├── tx_monitor.go
│       └── tx_sender.go
└── xrp/
    ├── keygensrv/
    │   ├── key_generator.go
    │   └── signer.go
    └── watchsrv/
        ├── tx_creator.go
        ├── tx_creator_deposit.go
        ├── tx_creator_payment.go
        ├── tx_creator_transfer.go
        ├── tx_monitor.go
        └── tx_sender.go
```

## Objectives

1. Organize application services by wallet type (watch, keygen, sign) as the primary structure
2. Group coin-specific implementations within each wallet type
3. Co-locate interfaces with their implementations
4. Separate use cases clearly (transaction creation, monitoring, signing, etc.)
5. Make the structure consistent and easy to navigate
6. Maintain backward compatibility during the refactoring process

## Proposed Structure

Reorganize the application layer to follow a consistent pattern: **wallet type → coin → use case**

```text
pkg/wallet/service/
├── watch/                      # Watch wallet services (online, public keys only)
│   ├── interfaces.go           # Watch wallet interfaces (TxCreator, TxMonitorer, etc.)
│   ├── shared/                 # Shared watch wallet services (coin-agnostic)
│   │   ├── address_importer.go
│   │   └── payment_request_creator.go
│   ├── btc/                    # BTC/BCH watch wallet services
│   │   ├── tx_creator.go       # Transaction creation (deposit, payment, transfer)
│   │   ├── tx_monitor.go       # Transaction monitoring
│   │   ├── tx_sender.go        # Transaction sending
│   │   └── address_importer.go # BTC-specific address import
│   ├── eth/                    # ETH watch wallet services
│   │   ├── tx_creator.go
│   │   ├── tx_monitor.go
│   │   └── tx_sender.go
│   └── xrp/                    # XRP watch wallet services
│       ├── tx_creator.go
│       ├── tx_monitor.go
│       └── tx_sender.go
├── keygen/                     # Keygen wallet services (offline, first signature)
│   ├── interfaces.go           # Keygen wallet interfaces (HDWalleter, PrivKeyer, etc.)
│   ├── shared/                 # Shared keygen services (coin-agnostic)
│   │   ├── hd_walleter.go      # HD wallet key generation
│   │   ├── seeder.go           # Seed generation
│   │   └── address_exporter.go # Address export
│   ├── btc/                    # BTC/BCH keygen services
│   │   ├── privkey_importer.go
│   │   ├── multisigaddress.go
│   │   └── fullpubkey_importer.go
│   ├── eth/                    # ETH keygen services
│   │   └── privkey_importer.go
│   └── xrp/                    # XRP keygen services
│       └── key_generator.go
└── sign/                       # Sign wallet services (offline, subsequent signatures)
    ├── interfaces.go           # Sign wallet interfaces (Signer, FullPubkeyExporter, etc.)
    ├── shared/                 # Shared sign services (coin-agnostic)
    │   └── (none currently, but structure ready)
    ├── btc/                    # BTC/BCH sign services
    │   ├── signer.go
    │   ├── privkey_importer.go
    │   └── fullpubkey_exporter.go
    ├── eth/                    # ETH sign services
    │   └── signer.go
    └── xrp/                    # XRP sign services
        └── signer.go
```

**Key Requirements:**

- Primary organization by wallet type (watch, keygen, sign)
- Secondary organization by coin (btc, eth, xrp) within each wallet type
- Shared services in `shared/` subdirectory when coin-agnostic
- Interfaces co-located with implementations in each wallet type
- Clear separation of use cases (tx_creator, tx_monitor, tx_sender, etc.)
- Application layer depends on domain layer and uses infrastructure through domain interfaces

## Implementation Steps

### Step 1: Create New Directory Structure

- [ ] Create `pkg/wallet/service/watch/` directory
- [ ] Create `pkg/wallet/service/keygen/` directory
- [ ] Create `pkg/wallet/service/sign/` directory
- [ ] Create `shared/` subdirectories under each wallet type
- [ ] Create coin-specific subdirectories (`btc/`, `eth/`, `xrp/`) under each wallet type
- [ ] Add package-level documentation explaining the organization

### Step 2: Reorganize Watch Wallet Services

- [ ] Move `pkg/wallet/service/watch_interface.go` → `pkg/wallet/service/watch/interfaces.go`
- [ ] Move `pkg/wallet/service/watchsrv/address_importer.go` → `pkg/wallet/service/watch/shared/address_importer.go`
- [ ] Move `pkg/wallet/service/watchsrv/payment_request_creator.go` → `pkg/wallet/service/watch/shared/payment_request_creator.go`
- [ ] Move `pkg/wallet/service/btc/watchsrv/*` → `pkg/wallet/service/watch/btc/`
- [ ] Move `pkg/wallet/service/eth/watchsrv/*` → `pkg/wallet/service/watch/eth/`
- [ ] Move `pkg/wallet/service/xrp/watchsrv/*` → `pkg/wallet/service/watch/xrp/`
- [ ] Update package declarations in moved files
- [ ] Update all imports across the codebase

### Step 3: Reorganize Keygen Wallet Services

- [ ] Move `pkg/wallet/service/cold_interface.go` → `pkg/wallet/service/keygen/interfaces.go`
  - Extract only keygen-related interfaces (HDWalleter, Seeder, AddressExporter, PrivKeyer, FullPubKeyImporter, Multisiger)
- [ ] Move `pkg/wallet/service/coldsrv/hd_walleter.go` → `pkg/wallet/service/keygen/shared/hd_walleter.go`
- [ ] Move `pkg/wallet/service/coldsrv/seeder.go` → `pkg/wallet/service/keygen/shared/seeder.go`
- [ ] Move `pkg/wallet/service/coldsrv/address_exporter.go` → `pkg/wallet/service/keygen/shared/address_exporter.go`
- [ ] Move `pkg/wallet/service/btc/coldsrv/keygensrv/*` → `pkg/wallet/service/keygen/btc/`
- [ ] Move `pkg/wallet/service/eth/keygensrv/*` → `pkg/wallet/service/keygen/eth/`
- [ ] Move `pkg/wallet/service/xrp/keygensrv/*` → `pkg/wallet/service/keygen/xrp/`
- [ ] Update package declarations in moved files
- [ ] Update all imports across the codebase

### Step 4: Reorganize Sign Wallet Services

- [ ] Create `pkg/wallet/service/sign/interfaces.go`
  - Extract sign-related interfaces from `cold_interface.go` (Signer, FullPubkeyExporter)
- [ ] Move `pkg/wallet/service/btc/coldsrv/signer.go` → `pkg/wallet/service/sign/btc/signer.go`
- [ ] Move `pkg/wallet/service/btc/coldsrv/signsrv/*` → `pkg/wallet/service/sign/btc/`
- [ ] Move `pkg/wallet/service/eth/keygensrv/signer.go` → `pkg/wallet/service/sign/eth/signer.go`
  - Note: ETH signer is currently in keygensrv, needs to be moved to sign
- [ ] Move `pkg/wallet/service/xrp/keygensrv/signer.go` → `pkg/wallet/service/sign/xrp/signer.go`
  - Note: XRP signer is currently in keygensrv, needs to be moved to sign
- [ ] Update package declarations in moved files
- [ ] Update all imports across the codebase

### Step 5: Update Interface Definitions

- [ ] Review and update `watch/interfaces.go`:
  - Ensure all watch wallet interfaces are properly defined
  - Add godoc comments for each interface
  - Group related interfaces logically
- [ ] Review and update `keygen/interfaces.go`:
  - Ensure all keygen wallet interfaces are properly defined
  - Add godoc comments for each interface
  - Group related interfaces logically
- [ ] Review and update `sign/interfaces.go`:
  - Ensure all sign wallet interfaces are properly defined
  - Add godoc comments for each interface
  - Group related interfaces logically

### Step 6: Update Service Implementations

- [ ] Ensure all service implementations in `watch/` implement interfaces from `watch/interfaces.go`
- [ ] Ensure all service implementations in `keygen/` implement interfaces from `keygen/interfaces.go`
- [ ] Ensure all service implementations in `sign/` implement interfaces from `sign/interfaces.go`
- [ ] Update service constructors to use domain types and infrastructure interfaces
- [ ] Remove any direct infrastructure dependencies (use domain interfaces instead)

### Step 7: Update Dependency Injection

- [ ] Update `pkg/di/container.go` to use new service paths
- [ ] Update factory functions to create services from new locations
- [ ] Ensure DI container returns interfaces, not concrete types
- [ ] Test dependency injection works correctly with new structure

### Step 8: Update Command Layer

- [ ] Update `pkg/command/watch/` to use new service paths
- [ ] Update `pkg/command/keygen/` to use new service paths
- [ ] Update `pkg/command/sign/` to use new service paths
- [ ] Ensure commands use service interfaces, not concrete implementations
- [ ] Update all imports

### Step 9: Update Tests

- [ ] Move test files to match new structure
- [ ] Update test imports to use new paths
- [ ] Update test setup to use new service constructors
- [ ] Ensure all tests pass with new structure
- [ ] Update integration tests if needed

### Step 10: Clean Up Old Structure

- [ ] Remove empty `pkg/wallet/service/btc/` directory
- [ ] Remove empty `pkg/wallet/service/eth/` directory
- [ ] Remove empty `pkg/wallet/service/xrp/` directory
- [ ] Remove empty `pkg/wallet/service/coldsrv/` directory
- [ ] Remove empty `pkg/wallet/service/watchsrv/` directory
- [ ] Remove old interface files if all interfaces have been moved
- [ ] Verify no remaining references to old paths

### Step 11: Update Documentation

- [ ] Add package documentation to `pkg/wallet/service/` explaining the organization
- [ ] Add package documentation to each wallet type directory (watch/, keygen/, sign/)
- [ ] Document the organization pattern (wallet type → coin → use case)
- [ ] Update `AGENTS.md` with application layer organization guidelines
- [ ] Update architecture diagrams if they exist
- [ ] Document how to add new services following the pattern

## Application Layer Organization Guidelines

### Wallet Types

The application layer is organized by three wallet types:

1. **Watch Wallet** (`watch/`): Online wallet that monitors addresses and creates/sends transactions
   - Uses public keys only
   - Can create and send transactions
   - Monitors transaction status
   - Imports addresses from external sources

2. **Keygen Wallet** (`keygen/`): Offline wallet that generates keys and creates first signature
   - Generates HD wallet keys
   - Imports private keys
   - Creates multisig addresses
   - Exports addresses and public keys

3. **Sign Wallet** (`sign/`): Offline wallet that provides subsequent signatures for multisig
   - Signs transactions
   - Imports private keys
   - Exports public keys

### Organization Pattern

```
wallet/service/
├── {wallet_type}/          # watch, keygen, or sign
│   ├── interfaces.go       # Interfaces for this wallet type
│   ├── shared/             # Coin-agnostic services
│   └── {coin}/             # btc, eth, or xrp
│       └── {use_case}.go  # Specific use case implementation
```

### Use Cases

Common use cases across wallet types:

- **Transaction Creation** (`tx_creator.go`): Creates deposit, payment, or transfer transactions
- **Transaction Monitoring** (`tx_monitor.go`): Monitors transaction status and balances
- **Transaction Sending** (`tx_sender.go`): Sends signed transactions to the network
- **Key Generation** (`hd_walleter.go`, `key_generator.go`): Generates cryptographic keys
- **Key Import/Export** (`privkey_importer.go`, `fullpubkey_importer.go`, etc.): Imports/exports keys
- **Address Management** (`address_importer.go`, `address_exporter.go`): Manages addresses
- **Multisig Operations** (`multisigaddress.go`): Creates and manages multisig addresses
- **Signing** (`signer.go`): Signs transactions

### Interface Organization

- Interfaces are co-located with implementations in each wallet type directory
- Each wallet type has its own `interfaces.go` file
- Interfaces should be grouped logically (e.g., transaction-related, key-related)
- Interfaces should use domain types, not infrastructure types

### Shared Services

- Services that work across multiple coins go in `shared/` subdirectory
- Examples: HD wallet generation (works for BTC/BCH), seed generation
- Shared services should not have coin-specific logic

## Migration Strategy

To maintain backward compatibility:

1. **Gradual Migration**: Move services incrementally, starting with least-used services
2. **Dual Support**: Keep old packages temporarily with deprecation notices (if needed)
3. **Update Imports**: Update imports package by package
4. **Remove Old Code**: After all imports are updated, remove deprecated packages

## Acceptance Criteria

- [ ] All services are organized by wallet type (watch, keygen, sign)
- [ ] Coin-specific services are grouped under their respective wallet types
- [ ] Shared services are in `shared/` subdirectories
- [ ] Interfaces are co-located with implementations
- [ ] Clear separation of use cases (tx_creator, tx_monitor, etc.)
- [ ] Application layer depends on domain layer and uses infrastructure through domain interfaces
- [ ] All existing tests pass
- [ ] No breaking changes to public APIs (maintain backward compatibility)
- [ ] Code builds successfully (`make check-build`)
- [ ] Linter passes (`make lint-fix`)
- [ ] Documentation is updated
- [ ] Dependency injection works correctly
- [ ] Commands work correctly with new structure

## Testing Strategy

1. **Unit Tests**:
   - Test each service in isolation
   - Mock domain services and infrastructure
   - Test service interfaces

2. **Integration Tests**:
   - Test services with real domain and infrastructure
   - Test wallet operations end-to-end (keygen, sign, watch)
   - Ensure services work correctly with new structure

3. **Regression Tests**:
   - Run all existing tests to ensure no functionality is broken
   - Test wallet operations (keygen, sign, watch) end-to-end
   - Verify commands work correctly

## Dependencies

- This task should be completed after:
  - Phase 3.2: Separate domain logic (already completed)
  - Phase 3.2: Separate infrastructure layer (already completed)

## Estimated Effort

- **Complexity**: Medium-High
- **Time Estimate**: 1-2 weeks
- **Files Affected**: ~50-70 files
- **Risk Level**: Medium (requires careful refactoring to avoid breaking changes)

## Notes

- Follow Clean Architecture principles strictly
- Application layer should orchestrate domain services and infrastructure
- Keep services focused on single use cases
- Use dependency injection for all dependencies
- Refer to `AGENTS.md` for coding standards
- Run `make check-build` and `make lint-fix` after each step
- Consider creating a migration guide for developers

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Application Services Pattern](https://martinfowler.com/eaaCatalog/applicationService.html)
- Project's `REFACTORING_PLAN.md` Phase 3.2
- Project's `AGENTS.md` for coding standards
- Project's `ISSUE_SEPARATE_DOMAIN_LOGIC.md` for domain layer structure
- Project's `ISSUE_SEPARATE_INFRASTRUCTURE_LAYER.md` for infrastructure layer structure

