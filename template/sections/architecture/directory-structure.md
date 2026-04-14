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
- `proto/` ... Protocol Buffers definitions (**Edition 2024**, see [docs/proto.md](../../../docs/reference/proto.md))
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
  - Communicates with `apps/xrpl-grpc-server/` [Deprecated]
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
