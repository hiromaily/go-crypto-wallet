# go-crypto-wallet

<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/xrp-img.jpg?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/ethereum-img.png?raw=true">
<img align="right" width="159px" src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/bitcoin-img.svg?sanitize=true">

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![codebeat badge](https://codebeat.co/badges/792a7c07-2352-4b7e-8083-0a323368b26f)](https://codebeat.co/projects/github-com-hiromaily-go-crypto-wallet-master)
[![GitHub release](https://img.shields.io/badge/release-v5.0.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/LICENSE)

Wallet functionalities to create raw transaction, to sign on unsigned transaction,
to send signed transaction for BTC, BCH, ETH, XRP and so on.  

## What kind of coin can be used?

- Bitcoin
- Bitcoin Cash
- Ethereum
- ERC-20 Token
- Ripple

## Current development

- This project is under refactoring based on `Clean Code`, `Clean Architecture`, [`Refactoring`](https://martinfowler.com/articles/refactoring-2nd-ed.html)
  - ✅ Domain layer separated (`internal/domain/`) - Pure business logic with zero infrastructure dependencies
  - ✅ Application layer (`internal/application/usecase/`) - Use case implementations following Clean Architecture
  - ✅ Infrastructure layer (`internal/infrastructure/`) - External dependencies (API clients, database, repositories)
  - ✅ Interface adapters layer (`internal/interface-adapters/`) - CLI commands and wallet adapters
  - 🔄 Migration from legacy `pkg/wallet/service/` to new `internal/application/usecase/` in progress
  - ✅ Integration tests separated using build tags (`//go:build integration`)
  - ✅ Go 1.25.5 with updated major dependencies (btcsuite/btcd v0.25.0, ethereum/go-ethereum v1.16.7)
  - ✅ golangci-lint v2.7.2 for code quality checks
- ✅ **Taproot (BIP341/BIP86) Support** - Full support for P2TR addresses with Schnorr signatures (requires Bitcoin Core v22.0+)
  - 30-50% transaction size/fee reduction compared to legacy multisig
  - Enhanced privacy with indistinguishable spend patterns
  - See [Taproot User Guide](./docs/TAPROOT_GUIDE.md) for setup and usage

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

![generate keys](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/0_key%20generation%20diagram.png?raw=true)

#### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address

![create tx](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

#### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address

![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/master/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)

## Requirements

### Core Dependencies

- **Go**: 1.25.5
- **golangci-lint**: v2.7.2+ (for development)
- **direnv**: For environment variable management
- **Docker**: For running blockchain nodes and databases

### Blockchain Nodes

- **BTC**: [Bitcoin Core](https://bitcoin.org/en/bitcoin-core/) 0.18+ (Bitcoin node)
- **BCH**: [Bitcoin ABC](https://www.bitcoinabc.org/) 0.21+ (Bitcoin Cash node)
- **ETH**:
  - [go-ethereum](https://github.com/ethereum/go-ethereum) (Geth client)
  - [Ganache](https://www.trufflesuite.com/ganache) (for local development)
  - [erc20-token](https://github.com/hiromaily/go-crypto-wallet/tree/master/web/erc20-token) (ERC-20 token contract)
- **XRP**:
  - [rippled](https://xrpl.org/manage-the-rippled-server.html) (Ripple node)
  - [ripple-lib-server](https://github.com/hiromaily/go-crypto-wallet/tree/master/web/ripple-lib-server) (gRPC server)

### Database

- **MySQL**: 5.7+ (for wallet data persistence)

### Major Go Dependencies

- **btcsuite/btcd**: v0.25.0 (Bitcoin library)
- **ethereum/go-ethereum**: v1.16.7 (Ethereum library)
- **spf13/cobra**: v1.10.2 (CLI framework)
- **spf13/viper**: v1.21.0 (Configuration management)
- **google.golang.org/grpc**: v1.78.0 (gRPC for XRP communication)
- **golang.org/x/crypto**: v0.46.0 (Cryptographic functions)

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
- `data/` ... Generated files and configuration
  - `address/` ... Generated address files (bch, btc, eth, xrp)
  - `config/` ... Configuration TOML files
  - `contract/` ... Contract ABI files
  - `keystore/` ... Keystore files for Ethereum
  - `proto/` ... Protocol buffer definitions (rippleapi)
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

## Installation

[Installation](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/Installation.md)

## Operation example

- [For Bitcoin](https://github.com/hiromaily/go-crypto-wallet/blob/master/docs/btc/OperationExample.md)
- [operation scripts](https://github.com/hiromaily/go-crypto-wallet/tree/master/scripts/operation)

## Command example

- [Makefile](https://github.com/hiromaily/go-crypto-wallet/blob/master/Makefile) - Main Makefile with modular includes
- Makefile modules (in `make/` directory):
  - [watch_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/master/make/watch_op.mk) - Watch wallet operations
  - [keygen_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/master/make/keygen_op.mk) - Keygen wallet operations
  - [sign_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/master/make/sign_op.mk) - Sign wallet operations
  - And other specialized modules for builds, tests, Docker, etc.

## TODO

### Basics

- [ ] Remove [github.com/cpacia/bchutil](https://github.com/cpacia/bchutil) due to outdated code. Try to replace to [github.com/gcash/bchd](https://github.com/gcash/bchd)
- [x] Separate dependent test as Integration Test using tag (`//go:build integration`)
- [ ] Complete migration from `pkg/wallet/service/` to `pkg/application/usecase/`
- [ ] Add ATOM tokens on [Cosmos Hub](https://hub.cosmos.network/main/hub-overview/overview.html)
- [ ] Add [Polkadot](https://polkadot.network/technology/)
- [ ] Various monitoring patterns to detect suspicious operations.
- [ ] Add Github Action as CI
- [ ] Generate mnemonic instead of seed. [bip-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)

### For BTC/BCH

- ✅ **Taproot (BIP341/BIP86) Support** - P2TR addresses with Schnorr signatures (See [Taproot Guide](./docs/TAPROOT_GUIDE.md))
- ✅ Native SegWit-Bech32 (P2WPKH) addresses supported
- ✅ P2SH-SegWit addresses supported
- ✅ Legacy P2PKH addresses supported
- [ ] Setup [Signet](https://en.bitcoin.it/wiki/Signet) environment for development use
- [ ] Fix `overpaying fee issue` on Signet. It says 725% overpaying.
- [ ] Multisig-address is used only once because of security reason, so after tx is sent,
  related receiver addresses should be updated by is_allocated=true.
- [ ] Sent tx is not proceeded in bitcoin network if fee is not enough comparatively.
  So re-sending tx functionality is required adding more fee.

### For ERC20 token

- [ ] Add any useful APIs using contract equivalent to ETH APIs
- [ ] Monitoring for ERC20 token

### For ETH

- [ ] Make sure that `quantity-tag` is used properly. e.g. when getting balance,
  which quantity-tag should be used, latest or pending.
- [ ] Handling secret of private key properly. Password could be passed from command line argument.

### For XRP

- [ ] Handling secret of private key properly. Password could be passed from command line argument.

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

## Project layout patterns

- The `pkg` layout pattern, refer to the
  [linked](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2) URLs for details.
