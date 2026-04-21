<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/directory-structure.tpl.md · Run `make docs` to regenerate.
-->

# Directory Structure

This document describes the current directory structure and dependency relationships of the go-crypto-wallet project.

## Overview

This project follows **Clean Architecture** principles with clear layer separation.
The project uses both `pkg/` and `internal/` directories:

- **`internal/`**: New architecture following Clean Architecture (domain, application, infrastructure, interface-adapters)
- **`pkg/`**: Public packages containing shared utilities (configuration, logging, database, etc.)
  - **Critical Rule**: Packages in `pkg/` MUST NOT import or depend on `internal/` directory
  - See `pkg/AGENTS.md` for detailed guidelines

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
│   ├── ports/              # Port interfaces
│   │   ├── api/            # API port interfaces
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
│   │   │   ├── btc/        # Bitcoin-specific API clients
│   │   │   └── bch/        # Bitcoin Cash-specific API clients
│   │   ├── ethereum/       # Ethereum JSON-RPC API clients
│   │   │   ├── eth/        # Ethereum-specific API clients
│   │   │   ├── erc20/      # ERC-20 token API clients
│   │   │   └── ethtx/      # Ethereum transaction utilities
│   │   └── ripple/         # Ripple gRPC API clients
│   │       └── xrp/        # XRP-specific API clients
│   ├── config/             # Configuration management
│   │   └── account/        # Account configuration (multisig, TOML)
│   ├── contract/           # Smart contract utilities
│   │   └── token-abi.go    # ERC-20 token ABI generated code
│   ├── database/           # Database connections and generated code
│   │   ├── models/         # Database models
│   │   │   └── rdb/        # RDB models
│   │   └── sqlc/           # SQLC generated database code
│   ├── repository/         # Data persistence implementations
│   │   ├── cold/           # Cold wallet repository (keygen, sign)
│   │   └── watch/          # Watch wallet repository
│   ├── storage/            # File storage implementations
│   │   └── file/           # File-based storage
│   │       ├── address/    # Address file storage (bch, xrp subdirs)
│   │       ├── fullpubkey/ # Full public key file storage
│   │       └── transaction.go # Transaction file storage
│   └── wallet/             # Wallet infrastructure
│       └── key/            # Key generation logic (HD wallet, seeds, BIP generators)
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
└── di/                     # Dependency injection container
    └── container.go        # DI container implementation
```

## Package Directory Structure (Shared Utilities)

The `pkg/` directory contains **public packages** that can be imported by external code.
These packages provide shared utilities, configuration management, logging, and other common functionality.

**⚠️ Critical Rule**: Packages in `pkg/` MUST NOT import or depend on any packages in `internal/` directory.

For detailed guidelines, see `pkg/AGENTS.md`.

```text
pkg/
├── AGENTS.md               # Guidelines for working with pkg/ packages
├── config/                 # Configuration management utilities
│   └── testutil/           # Test utilities for configuration
├── converter/              # Data conversion utilities
├── db/                     # Database utilities
│   └── mysql/              # MySQL database connection utilities
├── debug/                  # Debug utilities
├── decimal/                # Decimal number utilities
├── di/                     # Legacy dependency injection container
│                           # (for backward compatibility)
├── grpc/                   # gRPC client utilities
├── logger/                 # Logging utilities
│                           # (structured logging, noop logger, slog support)
├── serial/                 # Serialization utilities
├── testutil/               # Test utilities for various components
│                           # (btc, eth, xrp, repository, suite)
├── uuid/                   # UUID generation utilities
└── websocket/              # WebSocket client utilities
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

## Configuration and Data Directories

```text
config/                      # Git管理: アプリケーション設定ファイル
├── wallet/                 # ウォレット設定ファイル
│   ├── account.toml
│   ├── *_keygen.toml       # Keygen wallet configs
│   ├── *_sign.toml         # Sign wallet configs
│   └── *_watch.toml        # Watch wallet configs
└── blockchain/             # ブロックチェーンノード設定
    ├── bitcoind/
    ├── openeth/
    └── rippled/

proto/                       # Git管理: Protocol Buffers定義（コード生成のソース）
└── rippleapi/              # Ripple gRPC proto files

contracts/                    # Git管理: スマートコントラクトABI（コード生成のソース）
└── token.abi

data/                        # 生成されるファイル（.gitignore対象）
├── address/                # Generated address files (bch, btc, eth, xrp)
├── certs/                  # Certificates for Docker volumes
├── dump/                   # Database dumps
├── fullpubkey/             # Generated full public key files
├── keystore/               # Keystore files
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

The project follows Clean Architecture principles:

- **New Architecture**: `internal/` directory with Clean Architecture (domain, application, infrastructure, interface-adapters)
- **Shared Utilities**: `pkg/` directory contains public packages for external use (configuration, logging, utilities)
  - **Critical Rule**: Packages in `pkg/` MUST NOT import or depend on `internal/` directory
  - See `pkg/AGENTS.md` for detailed guidelines

For detailed refactoring status, see `docs/issues/REFACTORING_CHECKLIST.md`.

## References

- `AGENTS.md` - Detailed architecture guidelines
- `pkg/AGENTS.md` - Guidelines for working with `pkg/` packages
- [README.md](https://github.com/hiromaily/go-crypto-wallet/blob/main/README.md) - Project overview
- [docs/issues/REFACTORING_PLAN.md](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/issues/REFACTORING_PLAN.md) - Refactoring plan
