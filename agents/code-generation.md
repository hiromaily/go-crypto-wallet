# Code Generation

This document describes all code generation tools used in the go-crypto-wallet project.

## Overview

This project uses several code generation tools. **All auto-generated files contain `DO NOT EDIT` comments and must never be manually modified.**

## Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)
**Source**: `tools/atlas/schemas/*.hcl` (HCL schema files)
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: See [Database Management Guidelines](database.md) for detailed workflow.

## SQLC Schema Files (from Database Dumps)

**Tool**: Custom shell script (`scripts/db/extract-sqlc-schema.sh`)
**Source**: MySQL database dumps (`data/dump/sql/dump_*.sql`)
**Command**: `make extract-sqlc-schema-all` (or individual: `make extract-sqlc-schema-watch`, `make extract-sqlc-schema-keygen`, `make extract-sqlc-schema-sign`)

**Generated Files**:

- `tools/sqlc/schemas/01_watch.sql` - Watch schema for SQLC
- `tools/sqlc/schemas/02_keygen.sql` - Keygen schema for SQLC
- `tools/sqlc/schemas/03_sign.sql` - Sign schema for SQLC

**Note**: These schema files are extracted from MySQL database dumps. The source of truth is the Atlas HCL files (`tools/atlas/schemas/*.hcl`). To update schemas, modify the HCL files and run the database migration flow.

## Database Code (SQLC)

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

## Protocol Buffer Code (Go)

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

## Smart Contract ABI Code

**Tool**: [abigen](https://geth.ethereum.org/docs/tools/abigen) (from go-ethereum)
**Source**: `data/contract/token.abi`
**Command**: `make generate-abi` (or `abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go`)

**Generated Files**:

- `internal/infrastructure/contract/token-abi.go` - ERC-20 token contract bindings

**Note**: ABI code is generated from Ethereum smart contract ABI JSON files.

## Protocol Buffer Code (JavaScript/TypeScript)

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
   - Atlas: Edit `tools/atlas/schemas/*.hcl` (HCL schema files)
   - SQLC Schemas: **DO NOT EDIT** `tools/sqlc/schemas/*.sql` - these are auto-generated from database dumps. Edit `tools/atlas/schemas/*.hcl` instead.
   - SQLC Queries: Edit `tools/sqlc/queries/*.sql` (manually edited)
   - Protocol Buffers: Edit `data/proto/rippleapi/*.proto`
   - ABI: Edit `data/contract/token.abi` (or regenerate from Solidity source)
3. **Regenerate after source changes** - Run the appropriate make command after modifying source files
4. **Verify generation** - Run `make check-build` after regenerating to ensure code compiles

## Quick Reference

| Tool | Source Files | Command | Generated Files |
|------|--------------|---------|-----------------|
| Atlas | `tools/atlas/schemas/*.hcl` | `make atlas-dev-reset` | `tools/atlas/migrations/*/*.sql` |
| SQLC Schema Extract | `data/dump/sql/dump_*.sql` | `make extract-sqlc-schema-all` | `tools/sqlc/schemas/*.sql` |
| SQLC | `tools/sqlc/schemas/*.sql` + `tools/sqlc/queries/*.sql` | `make sqlc` | `internal/infrastructure/database/sqlc/*.go` |
| Protocol Buffers (Go) | `data/proto/rippleapi/*.proto` | `make protoc-go` | `internal/infrastructure/api/ripple/xrp/*.pb.go` |
| Smart Contract ABI | `data/contract/token.abi` | `make generate-abi` | `internal/infrastructure/contract/token-abi.go` |
| Protocol Buffers (JS/TS) | `data/proto/rippleapi/*.proto` | `web/ripple-lib-server/scripts/protoc-ts.sh` | `web/ripple-lib-server/src/pb/*.js` |

## See Also

- [Database Management Guidelines](database.md) - Detailed database schema workflow
- [Coding Standards](coding-standards.md) - Verification commands
- [Core Principles](core.md) - Rules about editing auto-generated files
