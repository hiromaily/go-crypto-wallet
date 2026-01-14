# go-crypto-wallet

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/xrp-img.jpg?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/ethereum-img.png?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![Test](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml/badge.svg)](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml)
[![GitHub release](https://img.shields.io/badge/release-v5.0.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Wallet functionalities to create raw transaction, to sign on unsigned transaction,
to send signed transaction for BTC, BCH, ETH, XRP and so on.

## What kind of coin can be used?

- Bitcoin
- Bitcoin Cash
- Ethereum
- ERC-20 Token
- Ripple

## Requirements

### Core Dependencies

| Tool | Version | Description |
|------|---------|-------------|
| Go | 1.25.5 | Programming language |
| MySQL | 8.4+ | Database (via Docker) |
| Atlas | 1.0.0 | Database schema migration |
| sqlc | 1.30.0 | SQL code generator |
| Docker | latest | Container runtime |
| Docker Compose | latest | Container orchestration |
| [golangci-lint](https://github.com/golangci/golangci-lint) | v2.7.2+ | Linter (for development) |
| [buf](https://buf.build/) | latest | Protocol buffer management |

### Blockchain Nodes

| Chain | Node | Version | Notes |
|-------|------|---------|-------|
| BTC | [Bitcoin Core](https://bitcoin.org/en/bitcoin-core/) | 28.0+ | supports v28-v30, [Docker image](https://hub.docker.com/r/bitcoin/bitcoin) |
| BCH | [Bitcoin ABC](https://www.bitcoinabc.org/) | 0.21+ | Bitcoin Cash node |
| ETH | [go-ethereum](https://github.com/ethereum/go-ethereum) | latest | Geth client |
| ETH | [Anvil](https://getfoundry.sh/anvil/overview/) | latest | For local development (Foundry) |
| XRP | [rippled](https://xrpl.org/manage-the-rippled-server.html) | latest | Ripple node |

### Major Go Dependencies

| Package | Version | Description |
|---------|---------|-------------|
| btcsuite/btcd | v0.25.0 | Bitcoin library |
| ethereum/go-ethereum | v1.16.7 | Ethereum library |
| spf13/cobra | v1.10.2 | CLI framework |
| spf13/viper | v1.21.0 | Configuration management |
| google.golang.org/grpc | v1.78.0 | gRPC for XRP communication |
| golang.org/x/crypto | v0.46.0 | Cryptographic functions |

## Current development

- This project is under refactoring based on `Clean Code`, `Clean Architecture`, [`Refactoring`](https://martinfowler.com/articles/refactoring-2nd-ed.html)
  - ✅ Domain layer separated (`internal/domain/`) - Pure business logic with zero infrastructure dependencies
  - ✅ Application layer (`internal/application/usecase/`) - Use case implementations following Clean Architecture
  - ✅ Infrastructure layer (`internal/infrastructure/`) - External dependencies (API clients, database, repositories)
  - ✅ Interface adapters layer (`internal/interface-adapters/`) - CLI commands and wallet adapters
  - ✅ Integration tests separated using build tags (`//go:build integration`)
  - ✅ Go 1.25.5 with updated major dependencies (btcsuite/btcd v0.25.0, ethereum/go-ethereum v1.16.7)
  - ✅ golangci-lint v2.7.2 for code quality checks
- ✅ **Taproot (BIP341/BIP86) Support** - Full support for P2TR addresses with Schnorr signatures
  - 30-50% transaction size/fee reduction compared to legacy multisig
  - Enhanced privacy with indistinguishable spend patterns
  - See [Taproot User Guide](./docs/TAPROOT_GUIDE.md) for setup and usage
- ✅ **MuSig2 (BIP327) Support** - Simple Two-Round Schnorr multisignatures for efficient multisig transactions
  - Single aggregated signature on-chain (looks like single-sig)
  - 30-50% smaller transactions compared to traditional P2WSH multisig
  - Parallel nonce generation (Round 1) for faster workflow
  - Maximum privacy - indistinguishable from single-signature transactions
  - See [MuSig2 User Guide](./docs/crypto/btc/musig2_guide.md) for setup and usage
- ✅ **Descriptor Wallets (BIP380)** - Generate, export, and import descriptors for BTC accounts (Bitcoin Core compatible)
  - Export receive/change descriptors in Bitcoin Core JSON or text formats
  - Watch wallet import/validate commands for descriptor onboarding (single-key; multisig import intentionally disabled)
  - See [Descriptor User Guide](./docs/descriptor_user_guide.md) and [Descriptor API](./docs/api/descriptor_api.md)

## Expected use cases

### 1.Deposit functionality

- Pubkey addresses are given to our users first.
- Users would want to deposit coins on our system.
- After users sent coins to their given addresses, these all amount of coins are sent
  to our safe addresses managed offline by cold wallet

### 2.Payment functionality

- Users would want to withdraw their coins to specific addresses.
- Transaction is created and sent after payment is requested by users.

### 3.Transfer functionality

- Internal use. Each accounts can transfer coins among internal accounts.

## Wallet Type

This is explained for BTC/BCH for now.
There are mainly 3 wallets separately and these wallets are expected to be installed in each different devices.

### 1.Watch only wallet

- Only this wallet run online to access to BTC/BCH Nodes.
- Only pubkey address is stored. Private key is NOT stored for security reason. That's why this is called `watch only wallet`.
- Major functionalities are
  - creating unsigned transaction
  - sending signed transaction
  - monitoring transaction status.

### 2.Keygen wallet as cold wallet

- Key management functionalities for accounts.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts
  - generating keys based on `HD Wallet`
  - generating multisig addressed according to account setting
  - exporting pubkey addresses as csv file which is imported from `Watch only wallet`
  - signing on unsigned transaction as first sign. However, multisig addresses could not be completed by only this wallet.

### 3.Sign wallet as cold wallet (Auth wallet)

- The internal authorization operators would use this wallet to sign on unsigned transaction for multisig addresses.
- Each of operators would be given own authorization account and Sing wallet apps.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts for own auth account
  - generating keys based on `HD Wallet` for own auth account
  - exporting full-pubkey addresses as csv file which is imported from `Keygen wallet` to generate multisig address
  - signing on unsigned transaction as second or more signs for multisig addresses.

## Workflow diagram

### BTC

#### 1. Generate keys

![generate keys](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/0_key%20generation%20diagram.png?raw=true)

#### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address

![create tx](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

#### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address

![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)

## Directory Structure

The project follows **Clean Architecture** principles with clear layer separation.
The codebase is organized into `internal/` (new architecture) and `pkg/` (shared utilities and legacy code).

### Root Directory

- `cmd/` ... Application entry points
  - `keygen/` ... Keygen wallet entry point
  - `sign/` ... Sign wallet entry point
  - `watch/` ... Watch wallet entry point
  - `tools/` ... Development tools
- `internal/` ... **New architecture** following Clean Architecture
- `pkg/` ... Shared utilities and legacy/transitional code
- `config/` ... Application configuration files (Git managed)
  - `wallet/` ... Wallet configuration TOML files
  - `blockchain/` ... Blockchain node configuration files
- `proto/` ... Protocol Buffers definitions (Git managed, code generation source)
- `contracts/` ... Smart contract ABI files (Git managed, code generation source)
- `data/` ... Generated files (ignored by Git)
  - `address/` ... Generated address files (bch, btc, eth, xrp)
  - `keystore/` ... Keystore files for Ethereum
  - `tx/` ... Transaction data files
- `docker/` ... Docker resources for blockchain nodes and databases
- `docs/` ... Documentation
- `scripts/` ... Operation and setup scripts
- `tools/` ... Development tools (sqlc configuration)
- `web/` ... Web-related projects
  - `erc20-token/` ... ERC-20 token contract
  - `ripple-lib-server/` ... Ripple gRPC server

### `internal/` Directory Structure (New Architecture)

The `internal/` directory contains the new architecture following Clean Architecture:

#### Domain Layer (`internal/domain/`)

Pure business logic with **zero infrastructure dependencies**:

- `domain/account/` ... Account types, validators, and business rules
- `domain/transaction/` ... Transaction types, state machine, validators
- `domain/wallet/` ... Wallet types and definitions
- `domain/key/` ... Key value objects and validators
- `domain/multisig/` ... Multisig validators and business rules
- `domain/coin/` ... Cryptocurrency type definitions

#### Application Layer (`internal/application/`)

Use case layer following Clean Architecture:

- `application/usecase/keygen/` ... Key generation use cases (btc, eth, xrp, shared)
- `application/usecase/sign/` ... Signing use cases (btc, eth, xrp, shared)
- `application/usecase/watch/` ... Watch wallet use cases (btc, eth, xrp, shared)

#### Infrastructure Layer (`internal/infrastructure/`)

External dependencies and implementations:

- `infrastructure/api/bitcoin/` ... Bitcoin/BCH Core RPC API clients
  - [API References](https://developer.bitcoin.org/reference/rpc/index.html)
- `infrastructure/api/ethereum/` ... Ethereum JSON-RPC API clients
  - [API References](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- `infrastructure/api/ripple/` ... Ripple gRPC API clients
  - Communicates with [ripple-lib-server](./web/ripple-lib-server/)
- `infrastructure/database/` ... Database connections and generated code
  - `mysql/` ... MySQL connection management
  - `sqlc/` ... SQLC generated database code
- `infrastructure/repository/` ... Data persistence implementations
  - `cold/` ... Cold wallet repository (keygen, sign)
  - `watch/` ... Watch wallet repository
- `infrastructure/storage/` ... File storage implementations
  - `file/` ... File-based storage (address, transaction)
- `infrastructure/network/` ... Network communication
  - `websocket/` ... WebSocket client implementations
- `infrastructure/wallet/key/` ... Key generation logic (HD wallet, seeds)

#### Interface Adapters Layer (`internal/interface-adapters/`)

Adapters between use cases and external interfaces:

- `interface-adapters/cli/` ... CLI command implementations
  - `keygen/` ... Keygen commands (api, create, export, imports, sign)
  - `sign/` ... Sign commands (create, export, imports, sign)
  - `watch/` ... Watch commands (api, create, imports, monitor, send)
- `interface-adapters/http/` ... HTTP handlers and middleware
- `interface-adapters/wallet/` ... Wallet adapter interfaces and implementations
  - `interfaces.go` ... Wallet interfaces (Keygener, Signer, Watcher)
  - `btc/` ... Bitcoin wallet implementations
  - `eth/` ... Ethereum wallet implementations
  - `xrp/` ... XRP wallet implementations

#### Dependency Injection

- `internal/di/` ... Dependency injection container

### `pkg/` Directory Structure (Shared Utilities)

The `pkg/` directory contains shared utilities and legacy/transitional code:

- `config/` ... Configuration management
- `logger/` ... Logging utilities
- `address/` ... Address formatting and utilities (bch, xrp)
- `contract/` ... Smart contract utilities (ERC-20 token ABI)
- `converter/` ... Data conversion utilities
- `debug/` ... Debug utilities
- `fullpubkey/` ... Full public key formatting utilities
- `serial/` ... Serialization utilities
- `testutil/` ... Test utilities (btc, eth, xrp, repository, suite)
- `uuid/` ... UUID generation utilities

**Note**: Some legacy packages in `pkg/` are being migrated to `internal/` as part of the refactoring effort.

## Components inside repository

- ripple-lib-server
  - ./web/ripple-lib-server
- erc20-token
  - ./web/erc20-token

## Development Environment

### DevContainer (Recommended for AI-Assisted Development)

This project provides an **optional** DevContainer configuration for a standardized, isolated development environment. DevContainer is particularly useful when working with AI coding assistants like Claude Code, GitHub Copilot, or Cursor.

**Key Benefits:**

- ✅ **Safe AI Development**: Isolated environment protects your host system from accidental AI-generated changes
- ✅ **Consistent Setup**: Pre-configured with Go 1.25.5, golangci-lint v2.7.2, Atlas v1.0.0, and GitHub CLI
- ✅ **Quick Start**: New developers can start coding in minutes
- ✅ **Zero Impact**: Local development workflow remains completely unchanged

**Quick Start:**

```bash
# 1. Open project in VS Code or Cursor
code .

# 2. Click "Reopen in Container" when prompted
# (or press F1 → "Dev Containers: Reopen in Container")

# 3. Start developing!
# All tools are pre-installed and ready
```

**Documentation:**

- 📖 [Complete DevContainer Guide](./docs/development/devcontainer.md) - Setup, usage, and troubleshooting
- 🤖 [AI-Assisted Development with DevContainer](./docs/development/devcontainer.md#using-with-ai-tools) - Claude Code, Copilot integration

**Note:** DevContainer is completely optional. Continue with local development if you prefer.

### Local Development

For traditional local development setup, follow the installation guide below.

## Installation

[Installation](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/Installation.md)

## Operation example

- [For Bitcoin](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/btc/OperationExample.md)
- [operation scripts](https://github.com/hiromaily/go-crypto-wallet/tree/main/scripts/operation)

## Command example

- [Makefile](https://github.com/hiromaily/go-crypto-wallet/blob/main/Makefile) - Main Makefile with modular includes
- Makefile modules (in `make/` directory):
  - [watch_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/watch_op.mk) - Watch wallet operations
  - [keygen_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/keygen_op.mk) - Keygen wallet operations
  - [sign_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/sign_op.mk) - Sign wallet operations
  - And other specialized modules for builds, tests, Docker, etc.

## Architecture

This project follows **Clean Architecture** principles with clear layer separation:

```text
Interface Adapters Layer (internal/interface-adapters/*)
    ↓ depends on
Application Layer (internal/application/usecase/*)
    ↓ depends on
Domain Layer (internal/domain/*)
    ↑ implements interfaces defined by
Infrastructure Layer (internal/infrastructure/*)
```

### Key Principles

- **Domain Layer** (`internal/domain/`): Pure business logic with **zero infrastructure dependencies**
  - Defines interfaces that infrastructure must implement
  - Contains business rules, validators, and value objects
  - Most stable layer - changes affect all other layers

- **Application Layer** (`internal/application/usecase/`): Use cases orchestrate business logic
  - Coordinates domain objects and infrastructure services
  - Organized by wallet type (keygen, sign, watch) and cryptocurrency (btc, eth, xrp)
  - Each use case represents a single business operation

- **Infrastructure Layer** (`internal/infrastructure/`): External dependencies and implementations
  - Implements interfaces defined by domain layer (Dependency Inversion Principle)
  - Contains API clients, database repositories, file storage, and network communication
  - Easily replaceable and mockable for testing

- **Interface Adapters Layer** (`internal/interface-adapters/`): Adapters between use cases and external interfaces
  - CLI commands, HTTP handlers, and wallet adapters
  - Converts between external formats and application DTOs

- **Dependency Direction**: Outer layers depend on inner layers, never the reverse

### Architecture Dependency Flow

```
┌─────────────────────────────────────────┐
│   Interface Adapters (CLI, HTTP)        │
└──────────────────┬──────────────────────┘
                   │ depends on
                   ↓
┌─────────────────────────────────────────┐
│   Application Layer (Use Cases)          │
└───────────┬───────────────────┬─────────┘
            │ depends on        │ depends on
            ↓                   ↓
┌───────────────────┐  ┌──────────────────────┐
│   Domain Layer    │  │ Infrastructure Layer │
│ (Business Logic)  │  │ (External Systems)   │
└───────────────────┘  └──────────┬───────────┘
                                 │ implements
                                 ↓
                        ┌────────────────────┐
                        │ Domain Interfaces  │
                        └────────────────────┘
```

For detailed architecture guidelines, see [AGENTS.md](./AGENTS.md).
