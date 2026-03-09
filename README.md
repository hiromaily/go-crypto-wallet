# go-crypto-wallet

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/xrp-img.jpg?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/ethereum-img.png?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![Test](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml/badge.svg)](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml)
[![GitHub release](https://img.shields.io/badge/release-v6.2.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Wallet functionalities to create raw transaction, to sign on unsigned transaction,
to send signed transaction for BTC, BCH, ETH, XRP and so on.

## What kind of coin can be used?

- Bitcoin
- Bitcoin Cash
- Ethereum
- ERC-20 Token
- XRP Ledger (Ripple)

## Requirements

### Core Dependencies

| Tool | Version | Description |
|------|---------|-------------|
| Go | 1.25.6 | Programming language |
| Atlas | 1.1.0 | Database schema migration |
| sqlc | 1.30.0 | SQL code generator |
| Docker | latest | Container runtime |
| Docker Compose | latest | Container orchestration |
| [golangci-lint](https://github.com/golangci/golangci-lint) | v2.10.0+ | Linter (for development) |
| [protoc](https://grpc.io/docs/protoc-installation/) | 33.0+ | Protocol buffer compiler (**Edition 2024**) |
| [buf](https://buf.build/) | latest | Protocol buffer management (lint, format) |

### Database

Supported databases (choose one)

| Tool | Version | Description |
|------|---------|-------------|
| PostgreSQL | 18.2+ | Database (via Docker) |
| MySQL | 8.4+ | Database (via Docker) |
| SQLite | 3.0+ | Database |

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
- ✅ **Taproot (BIP341/BIP86) Support** - Full support for P2TR addresses with Schnorr signatures
  - 30-50% transaction size/fee reduction compared to legacy multisig
  - Enhanced privacy with indistinguishable spend patterns
  - See [Taproot User Guide](./docs/TAPROOT_GUIDE.md) for setup and usage
- ✅ **MuSig2 (BIP327) Support** - Simple Two-Round Schnorr multisignatures for efficient multisig transactions
  - Single aggregated signature on-chain (looks like single-sig)
  - 30-50% smaller transactions compared to traditional P2WSH multisig
  - Parallel nonce generation (Round 1) for faster workflow
  - Maximum privacy - indistinguishable from single-signature transactions
  - See [MuSig2 User Guide](./docs/chains/btc/musig2/user-guide.md) for setup and usage
- ✅ **Descriptor Wallets (BIP380)** - Generate, export, and import descriptors for BTC accounts (Bitcoin Core compatible)
  - Export receive/change descriptors in Bitcoin Core JSON or text formats
  - Watch wallet import/validate commands for descriptor onboarding (single-key; multisig import intentionally disabled)
  - See [Descriptor User Guide](./docs/descriptor_user_guide.md) and [Descriptor API](./docs/api/descriptor_api.md)

## ✨ Comprehensive E2E Testing by Chain

This project includes **fully automated E2E tests** covering transaction patterns from legacy to cutting-edge. Each pattern is implemented and verified through real transactions on regtest.

### Bitcoin (BTC) Transaction Patterns

BTC supports 11 transaction patterns from legacy P2PKH to cutting-edge Taproot with MuSig2.

| Pattern | Type | Address Format | Signature | Status |
|---------|------|----------------|-----------|--------|
| **P1** | P2PKH Single-sig | `m.../n...` | 1-of-1 | ✅ Verified |
| **P2** | P2PKH 2-of-3 Multisig | `2...` | 2-of-3 | ✅ Verified |
| **P3** | P2SH-P2WPKH Single-sig | `2...` | 1-of-1 | ✅ Verified |
| **P4** | P2SH-P2WSH 2-of-3 | `2...` | 2-of-3 | ✅ Verified |
| **P5** | P2WPKH Native SegWit | `bcrt1q...` | 1-of-1 | ✅ Verified |
| **P6** | P2WSH 2-of-3 | `bcrt1q...` | 2-of-3 | ✅ Verified |
| **P7** | P2WSH 3-of-3 | `bcrt1q...` | 3-of-3 | ✅ Verified |
| **P8** | P2SH-P2WSH 3-of-3 | `2...` | 3-of-3 | ✅ Verified |
| **P9** | P2TR Taproot Single-sig | `bcrt1p...` | Schnorr | ✅ Verified |
| **P10** | P2TR MuSig2 N-of-N | `bcrt1p...` | Aggregated | 🔧 Framework |
| **P11** | P2TR Tapscript M-of-N | `bcrt1p...` | Script Path | 🔧 Framework |

### Bitcoin Cash (BCH) Transaction Patterns

BCH supports P2PKH and P2SH patterns with CashAddr format. SegWit and Taproot are **not supported** on BCH.

| Pattern | Type | Address Format | Signature | Status |
|---------|------|----------------|-----------|--------|
| **P1** | P2PKH Single-sig | `bchreg:q...` | 1-of-1 | ✅ Verified |
| **P2** | P2SH 2-of-3 Multisig | `bchreg:p...` | 2-of-3 | ✅ Verified |
| **P3** | P2SH 3-of-3 Multisig | `bchreg:p...` | 3-of-3 | ✅ Verified |

> **Note**: BCH uses CashAddr format (`bitcoincash:q...` for P2PKH, `bitcoincash:p...` for P2SH on mainnet).
> BCH signing requires SIGHASH_FORKID (0x41) for replay protection.

### Ethereum (ETH) Transaction Patterns

ETH supports EIP-1559 (Type 2) transactions with EOA single-sig, ERC-20 token transfers, Safe multisig payments, and MPC-TSS threshold signing.

| Pattern | Type | Address Format | Transaction Type | Status |
|---------|------|----------------|------------------|--------|
| **P1** | Single-sig (EOA) | `0x...` | EIP-1559 (Type 2) | ✅ Verified |
| **P2** | ERC-20 HYT Token Transfer | `0x...` | EIP-1559 (Type 2) | ✅ Verified |
| **P3** | Safe 2-of-2 Multisig Payment | `0x...` | Safe `execTransaction` | ✅ Verified |
| **P4** | MPC-TSS 2-of-3 Threshold Signing | `0x...` | EIP-1559 (Type 2) | ✅ Verified |

> **Note**: Supports both Anvil and Geth nodes. Database backend is configurable (SQLite/MySQL/PostgreSQL).
> P2 deploys the HYT ERC-20 contract via Foundry (forge) onto the local Anvil node before running the transfer workflow.
> P3 deploys a [Safe v1.4.1](https://safe.global/) proxy contract via Foundry, then exercises the full 2-of-2 multisig
> workflow: propose → offline EIP-712 sign (signer 1) → offline EIP-712 sign (signer 2) → broadcast `execTransaction`.
> P4 implements threshold ECDSA signing via MPC-TSS (tss-lib): 3 nodes run distributed key generation (DKG) offline,
> then 2-of-3 nodes collaborate online to produce a single ECDSA signature without any party holding the full private key.

### XRP Ledger (XRP) Transaction Patterns

XRP supports single-sig payment transfer with classic addresses (ed25519 or secp256k1).

| Pattern | Type | Address Format | Signing | Status |
|---------|------|----------------|---------|--------|
| **P1** | Single-sig Payment Transfer | `r...` (base58) | Keygen wallet (offline) | ✅ Verified |
| **P2** | 2-of-2 Multisig Payment Transfer | `r...` (base58) | Keygen wallet (offline, sequential) | ✅ Verified |

> **Note**: XRP uses classic addresses (`r...` prefix). The Sign wallet is not required — the
> Keygen wallet handles offline signing directly. E2E runs against a local **rippled standalone mode**
> node (Docker), equivalent to Bitcoin regtest. Ledgers are advanced manually via `ledger_accept`.
> P2 exercises the full XRP multisig workflow: set-signer-list → create-multisig-tx →
> add-multisig-signature (×2) → submit-multisig-tx.

### Quick Start

```bash
# BTC: Run any pattern with a single command
make btc-e2e P=1    # P2PKH Single-sig
make btc-e2e P=9    # P2TR Taproot

# BCH: Run supported patterns
make bch-e2e P=1    # P2PKH Single-sig
make bch-e2e P=2    # P2SH 2-of-3 Multisig
make bch-e2e P=3    # P2SH 3-of-3 Multisig

# ETH: Run individual patterns
make eth-e2e-p1     # Single-sig EIP-1559
make eth-e2e-p2     # ERC-20 HYT Token Transfer
make eth-e2e-p3     # Safe 2-of-2 Multisig Payment
make eth-e2e-p4     # MPC-TSS 2-of-3 Threshold Signing

# XRP: Run patterns
make xrp-e2e-p1     # Single-sig Payment Transfer
make xrp-e2e-p2     # 2-of-2 Multisig Payment Transfer

# ETH: Run all patterns in parallel
make eth-e2e-parallel          # Run P1, P2, P3, and P4 in parallel
make eth-e2e-ci-all            # CI mode (non-interactive)

# XRP: Run all patterns in parallel
make xrp-e2e-parallel          # Run P1 and P2 in parallel
make xrp-e2e-ci-all            # CI mode (non-interactive)

# Fresh start with full reset
make btc-e2e-reset P=5
make bch-e2e-reset P=3
make eth-e2e-p1-reset
make eth-e2e-p2-reset
make eth-e2e-p3-reset
make eth-e2e-p4-reset
make xrp-e2e-p1-reset
make xrp-e2e-p2-reset

# CI/CD mode (non-interactive)
make btc-e2e-ci P=3
make bch-e2e-ci P=3
make eth-e2e-p1-ci
make eth-e2e-p2-ci
make eth-e2e-p3-ci
make eth-e2e-p4-ci
make xrp-e2e-p1-ci
make xrp-e2e-p2-ci
```

### Why This Matters

- 🔒 **Production-Ready**: Every transaction pattern is tested end-to-end
- 🔄 **Regression Testing**: Catch breaking changes before they reach production
- 📚 **Learning Resource**: Real working examples of all Bitcoin script types
- 🚀 **CI/CD Ready**: Automated testing for continuous integration
- 🧩 **Modular Design**: Shared utilities in `btc_common.sh` reduce code duplication by 80%

See [E2E Transaction Patterns Guide](./docs/chains/btc/operations/e2e-transaction-patterns.md) for detailed documentation.

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
- `proto/` ... Protocol Buffers definitions (**Edition 2024**, see [docs/proto.md](./docs/proto.md))
- `contracts/` ... Smart contract ABI files (Git managed, code generation source)
- `data/` ... Generated files (ignored by Git)
  - `address/` ... Generated address files (bch, btc, eth, xrp)
  - `keystore/` ... Keystore files for Ethereum
  - `tx/` ... Transaction data files
- `docker/` ... Docker resources for blockchain nodes and databases
- `docs/` ... Documentation
- `scripts/` ... Operation and setup scripts
- `tools/` ... Development tools (sqlc configuration)
- `apps/` ... Web-related projects
  - `eth-contracts/` ... Ethereum smart contracts
  - `xrpl-grpc-server/` ... XRP Ledger (Ripple) gRPC server [Deprecated]

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

- `infrastructure/api/btc/` ... Bitcoin/BCH Core RPC API clients
  - [API References](https://developer.bitcoin.org/reference/rpc/index.html)
- `infrastructure/api/eth/` ... Ethereum JSON-RPC API clients
  - [API References](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- `infrastructure/api/xrp/` ... Ripple gRPC API clients
  - Communicates with [xrpl-grpc-server](./apps/xrpl-grpc-server/) [Deprecated]
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

- xrpl-grpc-server [Deprecated]
  - ./apps/xrpl-grpc-server
- eth-contracts
  - ./apps/eth-contracts

## Development Environment

### DevContainer (Recommended for AI-Assisted Development)

This project provides an **optional** DevContainer configuration for a standardized, isolated development environment. DevContainer is particularly useful when working with AI coding assistants like Claude Code, GitHub Copilot, or Cursor.

**Key Benefits:**

- ✅ **Safe AI Development**: Isolated environment protects your host system from accidental AI-generated changes
- ✅ **Consistent Setup**: Pre-configured with Go 1.25.6, golangci-lint v2.8.0, Atlas v1.0.0, and GitHub CLI
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

- 📖 [Complete DevContainer Guide](./docs/devcontainer.md) - Setup, usage, and troubleshooting
- 🤖 [AI-Assisted Development with DevContainer](./docs/devcontainer.md#using-with-ai-tools) - Claude Code, Copilot integration

**Note:** DevContainer is completely optional. Continue with local development if you prefer.

### Local Development

For traditional local development setup, follow the installation guide below.

## Installation

[Installation](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/Installation.md)

## Operation example

- [For Bitcoin](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/btc/OperationExample.md)
- [operation scripts](https://github.com/hiromaily/go-crypto-wallet/tree/main/scripts/operation)

## Command example

- [CLI Command Reference](./internal/interface-adapters/cli/README.md) - Command × Chain × UseCase matrix (SSOT): which commands each chain supports, corresponding use case interfaces, and role of each command
- [docs/commands.md](./docs/commands.md) - Detailed flags, options, and usage examples per command
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
