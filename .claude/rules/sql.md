---
paths: ["**/*.sql"]
---

# SQL File Rules

## Overview

Rules for modifying SQL files (`*.sql`) in go-crypto-wallet.

## File Categories

| Path                              | Description               | Editable          |
| --------------------------------- | ------------------------- | ----------------- |
| `tools/sqlc/queries/*.sql`        | SQLC query definitions    | ✅ Yes            |
| `tools/sqlc/schemas/*.sql`        | SQLC schema files         | ❌ Auto-generated |
| `tools/atlas/migrations/**/*.sql` | Atlas migration files     | ❌ Auto-generated |
| `docker/*.sql`                    | Initial data/seed scripts | ✅ Yes            |

## Verification Commands

### SQLC (Query Files)

| Command              | Purpose                      | Required    |
| -------------------- | ---------------------------- | ----------- |
| `make sqlc-compile`  | Check syntax and type errors | ✅ Yes      |
| `make sqlc-vet`      | Check for potential issues   | ✅ Yes      |
| `make sqlc-validate` | Compile + vet combined       | Recommended |

### After Changes

```bash
# 1. Validate queries
make sqlc-validate

# 2. Generate Go code
make sqlc

# 3. Verify Go build
make check-build
```

## SQLC Query Files

### Location

`tools/sqlc/queries/*.sql`

### Query Syntax

```sql
-- name: GetAccountByID :one
SELECT * FROM account_key
WHERE id = ?;

-- name: ListAccountsByAccountType :many
SELECT * FROM account_key
WHERE account_type = ?
ORDER BY id;

-- name: CreateAccount :execresult
INSERT INTO account_key (
    coin_type, account_type, account, full_pubkey_idx
) VALUES (?, ?, ?, ?);
```

### Query Annotations

| Annotation    | Description                       |
| ------------- | --------------------------------- |
| `:one`        | Returns single row                |
| `:many`       | Returns multiple rows             |
| `:exec`       | Executes without return           |
| `:execresult` | Returns result with affected rows |
| `:execrows`   | Returns number of affected rows   |

## Auto-Generated Files

**DO NOT EDIT** the following files:

### SQLC Schema Files

- `tools/sqlc/schemas/*.sql`

These are extracted from database dumps. Regenerate with:

```bash
make extract-sqlc-schema-all
```

### SQLC Generated Go Code

- `internal/infrastructure/database/sqlc/*.go`

Regenerate with:

```bash
make sqlc
```

### Atlas Migration Files

- `tools/atlas/migrations/**/*.sql`
- `tools/atlas/migrations/*/atlas.sum`

These are generated from HCL schemas. Use HCL files as the source of truth.

## Workflow for Query Changes

```bash
# 1. Edit query files
# tools/sqlc/queries/*.sql

# 2. Validate
make sqlc-validate

# 3. Generate Go code
make sqlc

# 4. Verify build
make check-build

# 5. Run tests
make gotest
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
