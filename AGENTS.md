# Agent Guidelines for go-crypto-wallet

This document provides guidelines for AI agents working on this project.

## Project Overview

- This project is a cryptocurrency wallet implementation in Go supporting BTC, BCH, ETH, XRP, and ERC-20 tokens
- The project is currently under refactoring based on Clean Architecture and Clean Code principles
- Security is of utmost importance (private key management, offline wallets)
- The project follows the `pkg` layout pattern

## Architecture Principles

- Follow Clean Architecture principles
- Maintain clear layer separation (domain, application, infrastructure)
- Use dependency injection and abstract with interfaces
- Follow the `pkg` layout pattern

### Domain Layer Guidelines

The `internal/domain/` package contains pure business logic with **ZERO infrastructure dependencies**.

**Key Principles:**

- Domain layer has NO dependencies on infrastructure (no database, no API clients, no file I/O)
- Domain defines interfaces; infrastructure implements them (Dependency Inversion Principle)
- All domain logic must be testable without mocks (pure functions preferred)
- Domain is the most stable layer - changes here affect all other layers

**Domain Layer Structure:**

- **Types & Value Objects**: Immutable objects defined by values (AccountType, TxType, CoinTypeCode)
- **Entities**: Objects with unique identity and lifecycle (not yet fully implemented)
- **Validators**: Business rule validation functions
- **Domain Services**: Stateless services with business logic

**Important:**

- When adding new business logic, first consider if it belongs in the domain layer
- Use domain validators for input validation before infrastructure operations
- Business rules should be in domain, not scattered across services

### Application Layer (Use Case) Guidelines

The `internal/application/usecase/` package implements the use case layer following Clean Architecture principles.

**Key Principles:**

- Use cases orchestrate business logic by coordinating domain objects and infrastructure services
- Each use case represents a single business operation with clear input and output
- Use cases act as thin wrappers that transform DTOs, delegate to services, and wrap errors with context
- Use cases depend on domain layer and infrastructure layer through interfaces (Dependency Inversion)
- Organized by wallet type (watch, keygen, sign) and cryptocurrency (btc, eth, xrp, shared)

**Use Case Structure:**

```go
// Use case interface definition
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

// Input/Output DTOs
type XxxInput struct {
    Param1 string
    Param2 int
}

type XxxOutput struct {
    Result string
}

// Implementation
type xxxUseCase struct {
    service ServiceInterface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    result, err := u.service.SomeMethod(input.Param1, input.Param2)
    if err != nil {
        return nil, fmt.Errorf("failed to execute xxx: %w", err)
    }
    return &XxxOutput{Result: result}, nil
}
```

**DTO Conventions:**

- **Input DTOs**: Contain all parameters needed for the use case operation
- **Output DTOs**: Contain all results returned by the use case
- DTOs use domain types (not primitive types when domain types exist)
- DTOs are passed by value for inputs, returned as pointers for outputs

**Error Handling:**

- Wrap service errors with context using `fmt.Errorf` with `%w`
- Error messages should describe the use case operation that failed
- Return domain errors when business rule violations occur
- Let infrastructure errors propagate with added context

**Organization Structure:**

```text
internal/application/usecase/
├── keygen/
│   ├── interfaces.go              # Use case interfaces
│   ├── btc/                       # Bitcoin-specific use cases
│   ├── eth/                       # Ethereum-specific use cases
│   ├── xrp/                       # XRP-specific use cases
│   └── shared/                    # Shared use cases (all coins)
├── sign/
│   ├── interfaces.go
│   ├── btc/
│   ├── eth/
│   ├── xrp/
│   └── shared/
└── watch/
    ├── interfaces.go
    ├── btc/
    ├── eth/
    ├── xrp/
    └── shared/
```

**Testing Approach:**

Use cases currently have constructor tests that verify:

- Use case can be instantiated with dependencies
- Correct interface implementation

For comprehensive testing strategy, see `docs/TESTING_STRATEGY.md`.

**When to Create a New Use Case:**

- New command functionality is added (commands should use use cases, not services directly)
- Existing service logic needs to be exposed to commands with different DTO structure
- Business logic needs to coordinate multiple services
- Transaction boundaries need to be defined

**Important:**

- Commands in `internal/interface-adapters/cli/` should ONLY depend on use cases, NOT services directly
- Use cases should be small and focused on a single operation
- Avoid business logic in use cases; delegate to domain or services
- Use cases are the entry point to application logic from command layer

## Coding Standards

- Follow `golangci-lint` configuration (`.golangci.yml`)
- Format code with `make format` (uses `gofumpt` and `goimports` via golangci-lint)
  - Import order: standard → third-party → local
- Use `make lint-fix` to run linting and formatting together (executes lint checks and format fixes)
- Maintain consistent naming conventions (lowercase package names, exported functions start with uppercase)

## Common Commands

After making code changes, use these commands to verify code correctness:

- `make lint-fix`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

**Command Constraints**:

- **DO NOT** use `go build -v` directly; use `make check-build` instead
- **DO NOT** use `go tool golangci-lint` directly; use `make lint-fix` instead

## Error Handling

- Wrap errors with `fmt.Errorf` + `%w`
- Use `errors.Is`/`errors.As` for error checking
- Include context information in error messages

## Panic Usage

Following the principle of separation of concerns, the project clearly separates the instance construction phase
from the instance usage phase.
Therefore, `panic` is only allowed during instance construction.
Specifically, `panic` is acceptable in:

- `main.go` files (application entry points)
- `internal/di` package (dependency injection container)
- `pkg/di` package (legacy dependency injection container - for backward compatibility)

**Important:**

- `panic` should **NOT** be used in business logic, service layers, or infrastructure layers
- Use proper error handling with error returns in all other packages
- The separation allows for fail-fast behavior during initialization while maintaining robust error handling during runtime

## Context Management

- Add `context.Context` to all API calls
- Implement timeouts and cancellation
- Implement graceful shutdown

## Security

- **NEVER** log private keys or sensitive information
- Encrypt or zero-clear private keys in memory when possible
- Do not pass passwords via command-line arguments; use secure input methods
- Conduct security review when making changes involving sensitive information

## Wallet Types Understanding

- **Watch Wallet**: Online, public keys only, creates and sends transactions
- **Keygen Wallet**: Offline, generates keys, first signature for multisig
- **Sign Wallet**: Offline, second and subsequent signatures for multisig

## Directory Structure

- `cmd/`: Application entry points (keygen, sign, watch)
- `internal/`: Internal packages (application-specific, not for external use)
  - `domain/`: **Domain layer** - Pure business logic (ZERO infrastructure dependencies)
    - `account/`: Account types, validators, and business rules
    - `transaction/`: Transaction types, state machine, validators
    - `wallet/`: Wallet types and definitions
    - `key/`: Key value objects and validators
    - `multisig/`: Multisig validators and business rules
    - `coin/`: Cryptocurrency type definitions
  - `application/`: **Application layer** - Use case layer (Clean Architecture)
    - `usecase/`: Use case implementations organized by wallet type
      - `keygen/`: Key generation use cases (btc, eth, xrp, shared)
      - `sign/`: Signing use cases (btc, eth, xrp, shared)
      - `watch/`: Watch wallet use cases (btc, eth, xrp, shared)
  - `infrastructure/`: **Infrastructure layer** - External dependencies and implementations
    - `api/`: External API clients
      - `bitcoin/`: Bitcoin/BCH Core RPC API clients (btc, bch)
      - `ethereum/`: Ethereum JSON-RPC API clients (eth, erc20)
      - `ripple/`: Ripple gRPC API clients (xrp)
    - `contract/`: Smart contract utilities (ERC-20 token ABI generated code)
    - `database/`: Database connections and generated code
      - `mysql/`: MySQL connection management
      - `sqlc/`: SQLC generated database code
    - `repository/`: Data persistence implementations
      - `cold/`: Cold wallet repository (keygen, sign)
      - `watch/`: Watch wallet repository
    - `storage/`: File storage implementations
      - `file/`: File-based storage (address, transaction)
    - `network/`: Network communication
      - `websocket/`: WebSocket client implementations
    - `wallet/key/`: Key generation logic - Infrastructure layer
  - `interface-adapters/`: **Interface Adapters layer** - Adapters between use cases and external interfaces
    - `cli/`: CLI command adapters (keygen, sign, watch)
      - `keygen/`: Keygen command implementations (api, create, export, imports, sign)
      - `sign/`: Sign command implementations (create, export, imports, sign)
      - `watch/`: Watch command implementations (api, create, imports, monitor, send)
    - `wallet/`: Wallet adapter interfaces and implementations
      - `interfaces.go`: Wallet interfaces (Keygener, Signer, Watcher)
      - `btc/`: Bitcoin wallet implementations
      - `eth/`: Ethereum wallet implementations
      - `xrp/`: XRP wallet implementations
  - `wallet/service/`: **Application layer** - Business logic orchestration (legacy/transitional)
    - `keygen/`: Key generation services (btc, eth, xrp, shared)
    - `sign/`: Signing services (btc, eth, xrp, shared)
    - `watch/`: Watch wallet services (btc, eth, xrp, shared)
  - `di/`: Dependency injection container
- `pkg/`: Shared packages (reusable, for external use)
  - `config/`: Configuration management utilities
    - `testutil/`: Test utilities for configuration
  - `logger/`: Logging utilities (structured logging, noop logger, slog support)
  - `converter/`: Data conversion utilities
  - `debug/`: Debug utilities
  - `serial/`: Serialization utilities
  - `testutil/`: Test utilities for various components (btc, eth, xrp, repository, suite)
  - `uuid/`: UUID generation utilities
  - `db/mysql/`: MySQL database connection utilities
  - `decimal/`: Decimal number utilities
  - `grpc/`: gRPC client utilities
  - `websocket/`: WebSocket client utilities
  - `di/`: Legacy dependency injection container (for backward compatibility)
  
  **Important**: See `pkg/AGENTS.md` for detailed guidelines on working with `pkg/` packages.
  **Critical Rule**: Packages in `pkg/` MUST NOT import or depend on any packages in `internal/` directory.
- `data/`: Generated files, configuration files
  - `address/`: Address data files (bch, btc, eth, xrp)
  - `config/`: Configuration files (account, wallet configs, node configs)
  - `contract/`: Contract ABI files
  - `keystore/`: Keystore files
  - `proto/`: Protocol buffer definitions (rippleapi)
  - `tx/`: Transaction data files (bch, btc, eth, xrp)
- `scripts/`: Operation scripts
  - `operation/`: Wallet operation scripts
  - `setup/`: Setup scripts for blockchain nodes

**Architecture Dependency Direction:**

```text
Interface Adapters (interface-adapters/*) → Application Layer (application/usecase, wallet/service) → Domain Layer (domain/*) ← Infrastructure Layer (infrastructure/*)
```

## Refactoring Status

- Refer to `REFACTORING_CHECKLIST.md` for current refactoring tasks
- Make changes incrementally without breaking existing functionality
- Run tests before and after refactoring

## Testing

- Use `//go:build integration` tag for integration tests
- Separate unit tests and integration tests
- Measure and improve test coverage

## Dependency Management

- Use `go mod tidy` to organize dependencies
- Follow procedures in `REFACTORING_CHECKLIST.md` when updating dependencies
- Run security scans (`govulncheck`)

## Logging

- Use structured logging
- Set appropriate log levels
- **NEVER** log sensitive information (private keys, passwords, etc.)

## Patterns to Avoid

- Using `log.Fatal` (except in `main`)
- Using `panic` outside of instance construction (i.e., outside of `main.go` and `pkg/di` package)
- Leaving commented-out code
- Unused imports, variables, or functions
- Ignoring errors (detected by `errcheck`)
- Using `go build -v` directly (use `make check-build` instead)
- Using `go tool golangci-lint` directly (use `make lint-fix` instead)
- For `unused-receiver` lint errors: **Remove the receiver entirely** instead of renaming it to `_`.
  Renaming to `_` will only cause the same error to appear for other receivers,
  so convert the method to a function from the start.

## Recommended Patterns

- Abstraction through interfaces
- Dependency injection
- Proper error wrapping with context
- Use of `context.Context`

## Documentation

- Add godoc comments to exported functions and methods
- Add package-level comments
- Include usage examples for complex logic

## Multi-Chain Support

- **BTC/BCH**: Bitcoin Core RPC API
- **ETH**: Ethereum JSON-RPC API, ERC-20 token support
- **XRP**: Communication via gRPC with ripple-lib-server

## Database Schema Changes

This project uses [Atlas](https://atlasgo.io/) for database schema management with HCL (HashiCorp Configuration Language) as the source of truth.

### Schema Files

There are 3 schema files corresponding to each wallet type:

- `tools/atlas/schemas/watch.hcl` - Watch wallet schema (online wallet)
- `tools/atlas/schemas/keygen.hcl` - Keygen wallet schema (offline, key generation)
- `tools/atlas/schemas/sign.hcl` - Sign wallet schema (offline, signing)

### How to Change Database Schema

**Step 1: Modify the HCL schema file**

Edit the appropriate `.hcl` file in `tools/atlas/schemas/` directory. These files are the single source of truth for database schema.

**Step 2: Format and lint the schema files**

Run the following commands to format and validate the HCL schema files:

```bash
make atlas-fmt
make atlas-lint
```

This will:

- Format all HCL schema files for consistency
- Validate the schema syntax and structure
- Ensure no errors exist before generating migrations

**Step 3: Regenerate migration files**

Run the following command to regenerate migration files from scratch:

```bash
make atlas-dev-reset
```

This command will:

- Delete all existing migration files
- Generate new migration files from the HCL schemas
- Prompt for confirmation before proceeding

**Step 4: Verify the migration**

Test the migration by recreating the database:

```bash
docker compose down -v
docker compose up
```

This will:

- Remove existing database volumes (`-v` flag)
- Start fresh containers and apply migrations
- Verify that no errors occur during migration

**Step 5: Regenerate SQLC code (if needed)**

If the schema changes affect queries, regenerate SQLC code:

```bash
make sqlc
```

**Step 6: Verify the build**

```bash
make check-build
```

### Important Notes

- **Always modify HCL files first** - Never edit migration SQL files directly
- **HCL files are the source of truth** - Migration files are auto-generated from HCL
- **Test locally before committing** - Always run the full `docker compose down -v && docker compose up` cycle

## Auto-Generated Files

This project uses several code generation tools.
**All auto-generated files contain `DO NOT EDIT` comments and must never be manually modified.**

### Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)  
**Source**: `tools/atlas/schemas/*.hcl` (HCL schema files)  
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: See [Database Schema Changes](#database-schema-changes) section for detailed workflow.

### SQLC Schema Files (from Database Dumps)

**Tool**: Custom shell script (`scripts/db/extract-sqlc-schema.sh`)  
**Source**: MySQL database dumps (`data/dump/sql/dump_*.sql`)  
**Command**: `make extract-sqlc-schema-all` (or individual: `make extract-sqlc-schema-watch`, `make extract-sqlc-schema-keygen`, `make extract-sqlc-schema-sign`)

**Generated Files**:

- `tools/sqlc/schemas/01_watch.sql` - Watch schema for SQLC
- `tools/sqlc/schemas/02_keygen.sql` - Keygen schema for SQLC
- `tools/sqlc/schemas/03_sign.sql` - Sign schema for SQLC

**Note**: These schema files are extracted from MySQL database dumps. The source of truth is the Atlas HCL files (`tools/atlas/schemas/*.hcl`). To update schemas, modify the HCL files and run the database migration flow.

### Database Code (SQLC)

**Tool**: [sqlc](https://sqlc.dev/)  
**Source**: `tools/sqlc/schemas/*.sql` (auto-generated) and `tools/sqlc/queries/*.sql` (manually edited)  
**Command**: `make sqlc` (or `cd tools/sqlc && sqlc generate`)

**Generated Files**:

- `internal/infrastructure/database/sqlc/*.go` (15 files)
  - `models.go` - Database models
  - `db.go` - Database connection code
  - `*.sql.go` - Query functions (account_key, address, auth_account_key,
    auth_fullpubkey, btc_tx, btc_tx_input, btc_tx_output, eth_detail_tx,
    payment_request, seed, tx, xrp_account_key, xrp_detail_tx)

**Note**: The legacy location `pkg/db/rdb/sqlcgen/*.go` is no longer generated and can be safely deleted.

**Note**: SQLC generates type-safe Go code from SQL queries and schemas.

### Protocol Buffer Code (Go)

**Tool**: [buf](https://buf.build/) with protoc-gen-go and protoc-gen-go-grpc  
**Source**: `data/proto/rippleapi/*.proto`  
**Command**: `make protoc-go` (or `buf generate`)

**Generated Files**:

- `internal/infrastructure/api/ripple/xrp/*.pb.go` (6 files)
  - `account.pb.go` - Account message types
  - `account_grpc.pb.go` - Account gRPC service code
  - `address.pb.go` - Address message types
  - `address_grpc.pb.go` - Address gRPC service code
  - `transaction.pb.go` - Transaction message types
  - `transaction_grpc.pb.go` - Transaction gRPC service code

**Note**: Protocol buffers are used for XRP (Ripple) gRPC communication.

### Smart Contract ABI Code

**Tool**: [abigen](https://geth.ethereum.org/docs/tools/abigen) (from go-ethereum)  
**Source**: `data/contract/token.abi`  
**Command**: `make generate-abi` (or `abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go`)

**Generated Files**:

- `internal/infrastructure/contract/token-abi.go` - ERC-20 token contract bindings

**Note**: ABI code is generated from Ethereum smart contract ABI JSON files.

### Protocol Buffer Code (JavaScript/TypeScript)

**Tool**: protoc with JavaScript/TypeScript plugins  
**Source**: `data/proto/rippleapi/*.proto`  
**Command**: `web/ripple-lib-server/scripts/protoc-ts.sh`

**Generated Files**:

- `web/ripple-lib-server/src/pb/*.js` - JavaScript/TypeScript protocol buffer code
  - `account_pb.js`, `account_grpc_pb.js`
  - `address_pb.js`, `address_grpc_pb.js`
  - `transaction_pb.js`, `transaction_grpc_pb.js`
  - `gogo/protobuf/gogoproto/gogo_pb.js`

**Note**: Used by the ripple-lib-server web project.

### Web Project Build Artifacts

**Tool**: Various build tools (Truffle, webpack, etc.)  
**Generated Files**:

- `web/erc20-token/build/` - Compiled smart contracts and frontend assets

**Note**: These are build outputs from the ERC-20 token web project.

### Dependency Lock Files

**Tool**: Go modules, npm/yarn  
**Generated Files**:

- `go.sum` - Go module checksums
- `web/*/yarn.lock` - Yarn package lock files
- `web/*/package-lock.json` - npm package lock files

**Note**: These files track exact dependency versions and should be committed to version control.

### Important Rules

1. **Never manually edit auto-generated files** - Changes will be overwritten on next generation
2. **Edit source files instead**:
   - Atlas: Edit `tools/atlas/schemas/*.hcl` (HCL schema files)
   - SQLC Schemas: **DO NOT EDIT** `tools/sqlc/schemas/*.sql` - these are auto-generated from database dumps. Edit `tools/atlas/schemas/*.hcl` instead.
   - SQLC Queries: Edit `tools/sqlc/queries/*.sql` (manually edited)
   - Protocol Buffers: Edit `data/proto/rippleapi/*.proto`
   - ABI: Edit `data/contract/token.abi` (or regenerate from Solidity source)
3. **Regenerate after source changes** - Run the appropriate make command after modifying source files
4. **Verify generation** - Run `make check-build` after regenerating to ensure code compiles

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments
  (typically auto-generated files from tools like sqlc, protoc, or go generate)
- See [Auto-Generated Files](#auto-generated-files) section for complete list
- **Git Operations**:
  - Allowed: `git add`, `git commit`, and `git push` to GitHub
  - **NOT allowed**: `git merge` operations
  - **NOT allowed**: `gh` command merge operations (e.g., `gh pr merge`)
  - **NOT allowed**: `git commit` and `git push` to `main` or `master` branches
