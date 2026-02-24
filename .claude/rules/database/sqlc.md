---
paths: ["tools/sqlc/**/*.sql", "tools/sqlc/*.yml", "internal/infrastructure/database/*/sqlcgen/*.go"]
---

# sqlc File Rules

## Overview

Rules for modifying SQLC-related files in go-crypto-wallet.
The project supports **three database backends**: PostgreSQL (default), MySQL, and SQLite.

## File Categories

| Path | Description | Editable |
| --- | --- | --- |
| `tools/sqlc/queries/postgres/*.sql` | SQLC query definitions (PostgreSQL) | Yes |
| `tools/sqlc/queries/mysql/*.sql` | SQLC query definitions (MySQL) | Yes |
| `tools/sqlc/schemas/postgres/*.sql` | SQLC schema files (PostgreSQL) | Auto-generated |
| `tools/sqlc/schemas/mysql/*.sql` | SQLC schema files (MySQL) | Auto-generated |
| `tools/sqlc/schemas/sqlite/*.sql` | SQLC schema files (SQLite) | Auto-generated |
| `tools/sqlc/sqlc_postgres.yml` | SQLC config (PostgreSQL) | Yes |
| `tools/sqlc/sqlc_mysql.yml` | SQLC config (MySQL) | Yes |
| `tools/sqlc/sqlc_sqlite.yml` | SQLC config (SQLite) | Yes |
| `internal/infrastructure/database/postgres/sqlcgen/*.go` | Generated Go code (PostgreSQL) | Auto-generated |
| `internal/infrastructure/database/mysql/sqlcgen/*.go` | Generated Go code (MySQL) | Auto-generated |
| `internal/infrastructure/database/sqlite/sqlcgen/*.go` | Generated Go code (SQLite) | Auto-generated |
| `tools/atlas/migrations/**/*.sql` | Atlas migration files | Auto-generated |
| `docker/*.sql` | Initial data/seed scripts | Yes |

## Verification Commands

### SQLC Validation

| Command | Purpose | Required |
| --- | --- | --- |
| `make sqlc-compile` | Check syntax and type errors (postgres + mysql) | Yes |
| `make sqlc-vet` | Check for potential issues (postgres + mysql) | Yes |
| `make sqlc-validate` | Compile + vet combined | Recommended |
| `make sqlc-lint` | Alias for sqlc-vet | |

### After Changes

```bash
# 1. Validate queries
make sqlc-validate

# 2. Generate Go code (default: postgres)
make sqlc

# 3. Verify Go build
make check-build
```

## SQLC Query Files

### Locations

- PostgreSQL: `tools/sqlc/queries/postgres/*.sql` (default)
- MySQL: `tools/sqlc/queries/mysql/*.sql`

### Query Syntax

**PostgreSQL** (uses `$1`, `$2`, ... for parameters):

```sql
-- name: GetAccountByID :one
SELECT * FROM account_key
WHERE id = $1;

-- name: ListAccountsByAccountType :many
SELECT * FROM account_key
WHERE account_type = $1
ORDER BY id;

-- name: CreateAccount :execresult
INSERT INTO account_key (
    coin_type, account_type, account, full_pubkey_idx
) VALUES ($1, $2, $3, $4);
```

**MySQL** (uses `?` for parameters):

```sql
-- name: GetAccountByID :one
SELECT * FROM account_key
WHERE id = ?;
```

### Query Annotations

| Annotation | Description |
| --- | --- |
| `:one` | Returns single row |
| `:many` | Returns multiple rows |
| `:exec` | Executes without return |
| `:execresult` | Returns result with affected rows |
| `:execrows` | Returns number of affected rows |

## Auto-Generated Files

**DO NOT EDIT** the following files:

### SQLC Schema Files

- `tools/sqlc/schemas/postgres/*.sql`
- `tools/sqlc/schemas/mysql/*.sql`
- `tools/sqlc/schemas/sqlite/*.sql`

These are extracted from database dumps. Regenerate with:

```bash
# Extract schemas from running database (default: postgres)
make extract-sqlc-schema-all

# For MySQL
make extract-sqlc-schema-all DB_DIALECT=mysql
```

### SQLC Generated Go Code

- `internal/infrastructure/database/postgres/sqlcgen/*.go`
- `internal/infrastructure/database/mysql/sqlcgen/*.go`
- `internal/infrastructure/database/sqlite/sqlcgen/*.go`

Regenerate with:

```bash
# PostgreSQL (default)
make sqlc

# MySQL
make sqlc-mysql

# SQLite
make sqlc-sqlite

# All backends
make sqlc-all
```

### Atlas Migration Files

- `tools/atlas/migrations/**/*.sql`
- `tools/atlas/migrations/*/atlas.sum`

These are generated from HCL schemas. Use HCL files as the source of truth.

## Workflow for Query Changes

```bash
# 1. Edit query files
# tools/sqlc/queries/postgres/*.sql (and/or mysql)

# 2. Validate
make sqlc-validate

# 3. Generate Go code
make sqlc          # postgres (default)
make sqlc-mysql    # mysql
make sqlc-all      # all backends

# 4. Verify build
make check-build

# 5. Run tests
make gotest
```

## Full Regeneration from Running Database

```bash
# Regenerate schemas + Go code from current database
make regenerate-sqlc-from-current-db              # postgres (default)
make regenerate-sqlc-from-current-db DB_DIALECT=mysql  # mysql
```

## Best Practices

### Query Naming

- Use descriptive names: `GetAccountByID`, `ListPendingTransactions`
- Prefix with action: `Get`, `List`, `Create`, `Update`, `Delete`

### Efficiency

- Use appropriate indexes
- Avoid `SELECT *` in production queries when possible
- Use pagination for large result sets

### Security

- Never include sensitive data in query comments
- Use parameterized queries (SQLC handles this)

### Multi-Dialect Consistency

- Keep query names consistent across postgres and mysql query files
- When adding a new query, add it to both dialects

## Quick Checklist

- [ ] `make sqlc-compile` passes
- [ ] `make sqlc-vet` passes
- [ ] Query names are descriptive
- [ ] Appropriate return type annotation (`:one`, `:many`, etc.)
- [ ] `make sqlc` generates without errors
- [ ] `make check-build` passes

## Related Documentation

- @docs/database/db-management.md - Database management guide
- @tools/sqlc/ - SQLC configuration and files

## Related Skills

- `db-migration` - Database schema workflow
- `go-development` - Go verification after SQLC generation
