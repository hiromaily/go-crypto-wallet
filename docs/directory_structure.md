# Directory Structure

This document describes the current directory structure and dependency relationships of the go-crypto-wallet project.

## Overview

This project follows **Clean Architecture** principles with clear layer separation.
The project uses both `pkg/` and `internal/` directories:

- **`internal/`**: New architecture following Clean Architecture (domain, application, infrastructure, interface-adapters)
- **`pkg/`**: Shared utilities that can be imported by external code (public APIs)
  - These packages must NOT depend on `internal/` directory
  - Contains configuration, logging, test utilities, and other common functionality

## Root Directory Structure

```text
.
├── cmd/                    # Application entry points
├── internal/               # New architecture (Clean Architecture)
├── pkg/                    # Shared utilities and legacy code
├── data/                   # Generated files and configuration
├── docker/                 # Docker resources
├── docs/                   # Documentation
├── scripts/                # Operation scripts
├── tools/                  # Development tools
└── web/                    # Web-related projects
```

## Internal Directory Structure (New Architecture)

The `internal/` directory contains the new architecture following Clean Architecture principles:

```text
internal/
├── domain/                 # Domain Layer - Pure business logic
│   ├── account/            # Account types, validators, business rules
│   ├── coin/               # Cryptocurrency type definitions
│   ├── key/                # Key value objects and validators
│   ├── multisig/           # Multisig validators and business rules
│   ├── transaction/        # Transaction types, state machine, validators
│   └── wallet/             # Wallet types and definitions
│
├── application/            # Application Layer - Use case layer
│   ├── ports/              # Port interfaces (API, persistence)
│   │   ├── api/            # API port interfaces (for future use)
│   │   └── persistence/    # Persistence port interfaces
│   └── usecase/            # Use case implementations
│       ├── keygen/         # Key generation use cases
│       │   ├── btc/        # Bitcoin-specific use cases
│       │   ├── eth/        # Ethereum-specific use cases
│       │   ├── xrp/        # XRP-specific use cases
│       │   └── shared/     # Shared use cases (all coins)
│       ├── sign/           # Signing use cases
│       │   ├── btc/
│       │   ├── eth/
│       │   ├── xrp/
│       │   └── shared/
│       └── watch/          # Watch wallet use cases
│           ├── btc/
│           ├── eth/
│           ├── xrp/
│           └── shared/
│
├── infrastructure/         # Infrastructure Layer - External dependencies
│   ├── api/                # External API clients
│   │   ├── bitcoin/        # Bitcoin/BCH Core RPC API clients
│   │   │   ├── bch/        # Bitcoin Cash API clients
│   │   │   └── btc/        # Bitcoin API clients
│   │   ├── ethereum/       # Ethereum JSON-RPC API clients
│   │   │   ├── erc20/      # ERC-20 token API clients
│   │   │   ├── eth/        # Ethereum API clients
│   │   │   └── ethtx/      # Ethereum transaction utilities
│   │   └── ripple/         # Ripple gRPC API clients
│   │       └── xrp/        # XRP API clients
│   ├── config/             # Configuration implementations
│   │   └── account/        # Account configuration
│   ├── contract/           # Smart contract utilities
│   │   └── token-abi.go    # ERC-20 token ABI
│   ├── database/           # Database connections and generated code
│   │   ├── models/         # Database models
│   │   │   └── rdb/        # RDB models
│   │   ├── mysql/          # MySQL connection management
│   │   └── sqlc/           # SQLC generated database code
│   ├── repository/         # Data persistence implementations
│   │   ├── cold/           # Cold wallet repository (keygen, sign)
│   │   └── watch/          # Watch wallet repository
│   ├── storage/            # File storage implementations
│   │   └── file/           # File-based storage
│   │       ├── address/    # Address file storage
│   │       │   ├── bch/    # Bitcoin Cash address utilities
│   │       │   └── xrp/    # XRP address utilities
│   │       ├── fullpubkey/ # Full public key file storage
│   │       └── transaction.go  # Transaction file storage
│   ├── wallet/             # Wallet infrastructure implementations
│   │   └── key/            # Key generation logic (HD wallet, seeds)
│   └── network/            # Network communication
│       └── websocket/      # WebSocket client implementations
│
├── interface-adapters/     # Interface Adapters Layer
│   ├── cli/                # CLI command implementations
│   │   ├── keygen/         # Keygen commands (api, create, export, imports, sign)
│   │   ├── sign/           # Sign commands (create, export, imports, sign)
│   │   └── watch/          # Watch commands (api, create, imports, monitor, send)
│   ├── http/               # HTTP handlers and middleware
│   └── wallet/             # Wallet adapter interfaces and implementations
│       ├── interfaces.go   # Wallet interfaces (Keygener, Signer, Watcher)
│       ├── btc/            # Bitcoin wallet implementations
│       ├── eth/            # Ethereum wallet implementations
│       └── xrp/            # XRP wallet implementations
│
├── wallet/                 # Deprecated - Backward compatibility aliases
│   ├── keygener.go         # Type alias to interface-adapters/wallet.Keygener
│   ├── signer.go           # Type alias to interface-adapters/wallet.Signer
│   └── watcher.go          # Type alias to interface-adapters/wallet.Watcher
│
└── di/                     # Dependency injection container
    └── container.go        # DI container implementation
```

## Package Directory Structure (Shared Utilities)

The `pkg/` directory contains shared utilities that can be imported by external code.
These packages are public APIs and must NOT depend on `internal/` directory:

```text
pkg/
├── AGENTS.md               # Guidelines for pkg/ directory
├── config/                 # Configuration management
│   └── testutil/           # Test utilities for config
├── converter/              # Data conversion utilities
├── debug/                  # Debug utilities
├── di/                     # Dependency injection utilities
│   └── infrastructure.go   # Infrastructure DI setup
├── logger/                 # Logging utilities
│   ├── global.go           # Global logger
│   ├── logger.go           # Logger implementation
│   ├── noop.go             # No-op logger
│   └── slog.go             # Structured logging
├── serial/                 # Serialization utilities
├── testutil/               # Test utilities
│   ├── btc.go              # Bitcoin test utilities
│   ├── eth.go              # Ethereum test utilities
│   ├── repository.go       # Repository test utilities
│   ├── suite.go            # Test suite utilities
│   └── xrp.go              # XRP test utilities
└── uuid/                   # UUID generation utilities
    ├── google.go           # Google UUID implementation
    └── types.go            # UUID types
```

## Command Entry Points

```text
cmd/
├── keygen/                 # Keygen wallet entry point
│   └── main.go
├── sign/                   # Sign wallet entry point
│   └── main.go
├── watch/                  # Watch wallet entry point
│   └── main.go
└── tools/                  # Development tools
    ├── eth-key/            # Ethereum key management tool
    └── get-eth-key/        # Ethereum key retrieval tool
```

## Data Directory

```text
data/
├── address/                # Generated address files (bch, btc, eth, xrp)
├── certs/                  # Certificates for Docker volumes
├── config/                 # Configuration TOML files
│   ├── account.toml
│   ├── *_keygen.toml       # Keygen wallet configs
│   ├── *_sign.toml         # Sign wallet configs
│   ├── *_watch.toml        # Watch wallet configs
│   └── [blockchain]/       # Blockchain node configs
├── contract/               # Contract ABI files
├── keystore/               # Keystore files
├── proto/                  # Protocol buffer definitions
│   └── rippleapi/          # Ripple gRPC proto files
└── tx/                     # Transaction data files (bch, btc, eth, xrp)
```

## Architecture Dependency Relationships

The project follows Clean Architecture with clear dependency direction:

### Dependency Flow

```text
┌─────────────────────────────────────────────────────────────┐
│                    Interface Adapters Layer                  │
│  (internal/interface-adapters/cli, internal/interface-      │
│   adapters/http)                                             │
└───────────────────────────┬─────────────────────────────────┘
                            │ depends on
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  (internal/application/usecase)                             │
│  - Orchestrates business logic                              │
│  - Coordinates domain objects and infrastructure            │
└───────────────┬───────────────────────┬─────────────────────┘
                │ depends on            │ depends on
                ↓                       ↓
┌───────────────────────────────┐  ┌───────────────────────────────┐
│      Domain Layer             │  │   Infrastructure Layer        │
│  (internal/domain/*)          │  │  (internal/infrastructure/*) │
│  - Pure business logic        │  │  - External API clients      │
│  - ZERO infrastructure deps   │  │  - Database repositories     │
│  - Defines interfaces         │  │  - File storage              │
│  - Business rules & validators│  │  - Network communication     │
└───────────────────────────────┘  └───────────────┬───────────────┘
                                                   │ implements
                                                   ↓
                                        ┌──────────────────────────┐
                                        │   Domain Interfaces      │
                                        │  (defined in domain/)    │
                                        └──────────────────────────┘
```

### Key Principles

1. **Domain Layer** (`internal/domain/`)
   - Has **ZERO** dependencies on infrastructure
   - Contains pure business logic
   - Defines interfaces that infrastructure must implement
   - Most stable layer - changes affect all layers

2. **Application Layer** (`internal/application/`)
   - Depends on domain layer
   - Orchestrates business logic using domain objects
   - Uses infrastructure through domain interfaces
   - Organized by wallet type (keygen, sign, watch) and coin (btc, eth, xrp)

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements interfaces defined by domain layer
   - Contains NO business logic (only technical implementation)
   - Easily replaceable and mockable for testing
   - Handles external system communication

4. **Interface Adapters Layer** (`internal/interface-adapters/`)
   - Depends on application layer (use cases)
   - Converts between external formats and application DTOs
   - CLI commands and HTTP handlers

### Dependency Rules

- ✅ **Allowed**: Application → Domain
- ✅ **Allowed**: Infrastructure → Domain (implements domain interfaces)
- ✅ **Allowed**: Interface Adapters → Application
- ❌ **Forbidden**: Domain → Application
- ❌ **Forbidden**: Domain → Infrastructure
- ❌ **Forbidden**: Application → Infrastructure (directly, must go through domain interfaces)

## Wallet Types

The project supports three wallet types:

1. **Watch Wallet** (`watch/`)
   - Online wallet
   - Public keys only
   - Creates and sends transactions
   - Monitors blockchain

2. **Keygen Wallet** (`keygen/`)
   - Offline wallet
   - Generates keys
   - First signature for multisig transactions

3. **Sign Wallet** (`sign/`)
   - Offline wallet
   - Second and subsequent signatures for multisig transactions

## Supported Cryptocurrencies

- **BTC**: Bitcoin
- **BCH**: Bitcoin Cash
- **ETH**: Ethereum
- **ERC-20**: ERC-20 tokens
- **XRP**: Ripple

## Migration Status

The project follows Clean Architecture with clear separation:

- **New Architecture**: `internal/` directory with Clean Architecture (domain, application, infrastructure, interface-adapters)
- **Shared Utilities**: `pkg/` directory contains shared utilities that can be imported by external code
  - These packages are public APIs and must NOT depend on `internal/` directory
  - Contains configuration, logging, test utilities, and other common functionality
- **Backward Compatibility**: Type aliases in `internal/wallet/` for compatibility (temporary, will be removed)

For detailed refactoring status, see `docs/issues/REFACTORING_CHECKLIST.md`.

## References

- [AGENTS.md](../AGENTS.md) - Detailed architecture guidelines
- [README.md](../README.md) - Project overview
- [docs/issues/REFACTORING_PLAN.md](issues/REFACTORING_PLAN.md) - Refactoring plan
