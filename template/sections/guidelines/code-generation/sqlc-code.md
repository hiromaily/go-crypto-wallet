### Database Code (SQLC)

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
