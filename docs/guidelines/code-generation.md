# Code Generation

This document describes all code generation tools used in the go-crypto-wallet project.

## Overview

This project uses several code generation tools. **All auto-generated files contain `DO NOT EDIT` comments and must never be manually modified.**
## Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)
**Source**: `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL schema files)
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: See [Database Management Guidelines](../database/architecture.md) for detailed workflow.
## SQLC Schema Files (from Database Dumps)

**Tool**: Custom shell script (`scripts/db/extract-sqlc-schema.sh`)
**Source**: MySQL database dumps (`data/dump/sql/dump_*.sql`)
**Command**: `make extract-sqlc-schema-all` (or individual: `make extract-sqlc-schema-watch`, `make extract-sqlc-schema-keygen`, `make extract-sqlc-schema-sign`)

**Generated Files**:

- `tools/sqlc/schemas/mysql/01_watch.sql` - Watch schema for SQLC
- `tools/sqlc/schemas/mysql/02_keygen.sql` - Keygen schema for SQLC
- `tools/sqlc/schemas/mysql/03_sign.sql` - Sign schema for SQLC

**Note**: These schema files are extracted from MySQL database dumps. The source of truth is the Atlas HCL files (`tools/atlas/schemas/{db_dialect}/*.hcl`). To update schemas, modify the HCL files and run the database migration flow.
## Database Code (SQLC)

**Tool**: [sqlc](https://sqlc.dev/)
**Source**: `tools/sqlc/schemas/mysql/*.sql` (auto-generated) and `tools/sqlc/queries/mysql/*.sql` (manually edited)
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
## Protocol Buffer Code (Go) [DEPRECATED for XRP]

> **⚠️ DEPRECATED**: XRP protocol buffers (`proto/rippleapi/`) are no longer used.
>
> **Status**: `apps/xrpl-grpc-server/` and `proto/xrpapi/` are deprecated. XRP functionality now uses native Go implementation with xrpl-go libraries.
>
> **Migration**: See `.kiro/specs/xrp-transaction-flow-alignment/` for current architecture.

**Tool**: protoc (or buf when Edition 2024 is supported)
**Source**: `proto/rippleapi/*.proto` [DEPRECATED]
**Command**: `make proto` [NO LONGER NEEDED FOR XRP]

**Generated Files** [DEPRECATED]:

- `internal/infrastructure/api/xrp/xrp/*.pb.go` (6 files) [No longer generated]
  - `account.pb.go` - Account message types [Deprecated]
  - `account_grpc.pb.go` - Account gRPC service code [Deprecated]
  - `address.pb.go` - Address message types [Deprecated]
  - `address_grpc.pb.go` - Address gRPC service code [Deprecated]
  - `transaction.pb.go` - Transaction message types [Deprecated]
  - `transaction_grpc.pb.go` - Transaction gRPC service code [Deprecated]

**Note**: Protocol buffers were previously used for XRP (Ripple) gRPC communication. This is no longer the case.
## Smart Contract ABI Code

**Tool**: [abigen](https://geth.ethereum.org/docs/tools/abigen) (from go-ethereum)
**Source**: `contracts/token.abi`
**Command**: `make generate-abi` (or `abigen --abi ./contracts/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go`)

**Generated Files**:

- `internal/infrastructure/contract/token-abi.go` - ERC-20 token contract bindings

**Note**: ABI code is generated from Ethereum smart contract ABI JSON files.
## Mock Code (Mockery)

**Tool**: [mockery v3](https://github.com/vektra/mockery)
**Source**: Interface definitions in Go files
**Configuration**: `.mockery.yaml`
**Command**: `make mockery` (or `go tool github.com/vektra/mockery/v3`)

**Generated Files**:

- `internal/infrastructure/api/btc/mocks/mock_bitcoiner.go` - Bitcoin API mock
- `internal/infrastructure/repository/mocks/mock_*.go` - Repository interface mocks
- `internal/infrastructure/storage/file/transaction/mocks/mock_transaction_file_repositorier.go` - File storage mock

**Mock Directory Structure**:

Mocks are placed in `mocks/` subdirectories alongside their implementations:

```text
internal/infrastructure/
├── api/bitcoin/
│   ├── btc/bitcoin.go          # Implementation
│   └── mocks/mock_bitcoiner.go # Generated mock
├── repository/
│   └── mocks/mock_*.go         # Persistence interface mocks
└── storage/file/
    └── transaction/
        ├── transaction.go      # Implementation
        └── mocks/
            └── mock_transaction_file_repositorier.go  # Storage interface mocks
```

**Adding New Mocks**:

1. Edit `.mockery.yaml`
2. Add the interface under the appropriate package section
3. Run `make mockery`

**Moving Mocks Directories**:

⚠️ **IMPORTANT**: When moving implementation code that has associated mocks, you **MUST** also update `.mockery.yaml` to reflect the new directory structure.

**Steps when moving mocks directories**:

1. Move the implementation code to the new location
2. **Update `.mockery.yaml`** - Change the `dir` path in the configuration for the affected interface(s)
3. Run `make mockery` - This will automatically:
   - Clean all existing mocks (via `clean-mocks` dependency)
   - Regenerate mocks in the new location based on updated `.mockery.yaml`

**Example**: If moving `internal/infrastructure/storage/file/transaction.go` to `internal/infrastructure/storage/file/transaction/transaction.go`:

```yaml
# Before
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/mocks"
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:

# After
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/transaction/mocks"  # Updated path
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:
```

**Note**: The `make mockery` target automatically runs `clean-mocks` first, so old mocks will be removed before generating new ones. This ensures no stale mocks remain when paths change.

**Note**: See [Testing Guidelines](./testing.md) for mock usage examples and best practices.
## Protocol Buffer Code (TypeScript) [DEPRECATED]

> **⚠️ DEPRECATED**: XRP gRPC server (`apps/xrpl-grpc-server/`) is no longer used.
>
> **Status**: All XRP functionality migrated to native Go implementation. TypeScript protobuf code generation is no longer needed.

**Tool**: protoc with protoc-gen-es and protoc-gen-connect-es [NO LONGER USED]
**Source**: `proto/rippleapi/*.proto` [DEPRECATED]
**Command**: `make proto-ts` [NO LONGER NEEDED]

**Generated Files** [DEPRECATED]:

- `apps/xrpl-grpc-server/src/protogen/*.ts` - TypeScript protocol buffer code [No longer generated]
  - `account_pb.ts` - Account message types [Deprecated]
  - `account_connect.ts` - Account service client [Deprecated]
  - `address_pb.ts` - Address message types [Deprecated]
  - `address_connect.ts` - Address service client [Deprecated]
  - `transaction_pb.ts` - Transaction message types [Deprecated]
  - `transaction_connect.ts` - Transaction service client [Deprecated]
## Web Project Build Artifacts

**Tool**: Various build tools (Truffle, webpack, etc.)

**Generated Files**:

- `web/erc20-token/build/` - Compiled smart contracts and frontend assets

**Note**: These are build outputs from the ERC-20 token web project.
## Dependency Lock Files

**Tool**: Go modules, npm/yarn

**Generated Files**:

- `go.sum` - Go module checksums
- `web/*/yarn.lock` - Yarn package lock files
- `web/*/package-lock.json` - npm package lock files

**Note**: These files track exact dependency versions and should be committed to version control.
## Important Rules

1. **Never manually edit auto-generated files** - Changes will be overwritten on next generation
2. **Edit source files instead**:
   - Atlas: Edit `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL schema files)
   - SQLC Schemas: **DO NOT EDIT** `tools/sqlc/schemas/{db_dialect}/*.sql` - these are auto-generated from database dumps. Edit `tools/atlas/schemas/{db_dialect}/*.hcl` instead.
   - SQLC Queries: Edit `tools/sqlc/queries/{db_dialect}/*.sql` (manually edited)
   - Mockery: Edit `.mockery.yaml` to add new interfaces, then run `make mockery`
   - Protocol Buffers: Edit `proto/rippleapi/*.proto`
   - ABI: Edit `contracts/token.abi` (or regenerate from Solidity source)
3. **Regenerate after source changes** - Run the appropriate make command after modifying source files
4. **Verify generation** - Run `make check-build` after regenerating to ensure code compiles
## Quick Reference

| Tool | Source Files | Command | Generated Files |
|------|--------------|---------|-----------------|
| Atlas | `tools/atlas/schemas/{db_dialect}/*.hcl` | `make atlas-dev-reset` | `tools/atlas/migrations/*/*.sql` |
| SQLC Schema Extract | `data/dump/sql/dump_*.sql` | `make extract-sqlc-schema-all` | `tools/sqlc/schemas/mysql/*.sql` |
| SQLC | `tools/sqlc/schemas/mysql/*.sql` + `tools/sqlc/queries/mysql/*.sql` | `make sqlc` | `internal/infrastructure/database/mysql/sqlcgen/*.go` |
| Mockery | `.mockery.yaml` + interface definitions | `make mockery` | `internal/infrastructure/*/mocks/*.go` |
| ~~Protocol Buffers (Go)~~ [DEPRECATED] | ~~`proto/rippleapi/*.proto`~~ | ~~`make proto`~~ | ~~(XRP protobufs no longer used)~~ |
| Smart Contract ABI | `contracts/token.abi` | `make generate-abi` | `internal/infrastructure/contract/token-abi.go` |
| ~~Protocol Buffers (TS)~~ [DEPRECATED] | ~~`proto/rippleapi/*.proto`~~ | ~~`make proto-ts`~~ | ~~(XRP gRPC server no longer used)~~ |
## See Also

- [Database Management Guidelines](../database/architecture.md) - Detailed database schema workflow
- [Testing Guidelines](./testing.md) - Mock usage and unit testing patterns
- [Coding Standards](./coding-conventions.md) - Verification commands
- [Core Principles](./core.md) - Rules about editing auto-generated files
