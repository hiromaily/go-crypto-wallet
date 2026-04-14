### Quick Reference

| Tool | Source Files | Command | Generated Files |
|------|--------------|---------|-----------------|
| Atlas | `tools/atlas/schemas/{db_dialect}/*.hcl` | `make atlas-dev-reset` | `tools/atlas/migrations/*/*.sql` |
| SQLC Schema Extract | `data/dump/sql/dump_*.sql` | `make extract-sqlc-schema-all` | `tools/sqlc/schemas/mysql/*.sql` |
| SQLC | `tools/sqlc/schemas/mysql/*.sql` + `tools/sqlc/queries/mysql/*.sql` | `make sqlc` | `internal/infrastructure/database/mysql/sqlcgen/*.go` |
| Mockery | `.mockery.yaml` + interface definitions | `make mockery` | `internal/infrastructure/*/mocks/*.go` |
| ~~Protocol Buffers (Go)~~ [DEPRECATED] | ~~`proto/rippleapi/*.proto`~~ | ~~`make proto`~~ | ~~(XRP protobufs no longer used)~~ |
| Smart Contract ABI | `contracts/token.abi` | `make generate-abi` | `internal/infrastructure/contract/token-abi.go` |
| ~~Protocol Buffers (TS)~~ [DEPRECATED] | ~~`proto/rippleapi/*.proto`~~ | ~~`make proto-ts`~~ | ~~(XRP gRPC server no longer used)~~ |
