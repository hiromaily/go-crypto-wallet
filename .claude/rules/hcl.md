---
paths: ["**/*.hcl"]
---

# HCL File Rules

## Overview

Rules for modifying HCL (HashiCorp Configuration Language) files (`*.hcl`) in go-crypto-wallet.

## Applicable Directories

| Path                    | Description                                   |
| ----------------------- | --------------------------------------------- |
| `tools/atlas/schemas/`  | Database schema definitions (source of truth) |
| `tools/atlas/atlas.hcl` | Atlas configuration                           |

## Schema Files

| File                             | Description                    |
| -------------------------------- | ------------------------------ |
| `tools/atlas/schemas/watch.hcl`  | Watch wallet schema (online)   |
| `tools/atlas/schemas/keygen.hcl` | Keygen wallet schema (offline) |
| `tools/atlas/schemas/sign.hcl`   | Sign wallet schema (offline)   |

## Verification Commands

| Command               | Purpose                      | Required    |
| --------------------- | ---------------------------- | ----------- |
| `make atlas-fmt`      | Format HCL files             | ✅ Yes      |
| `make atlas-lint`     | Lint and validate schemas    | ✅ Yes      |
| `make atlas-validate` | Validate Atlas configuration | Recommended |

## Complete Workflow

After modifying HCL schema files:

```bash
# 1. Format
make atlas-fmt

# 2. Lint
make atlas-lint

# 3. Regenerate migrations
make atlas-dev-reset

# 4. Reset mysql database and apply
docker compose down -v
docker compose --profile mysql up -d

# 5. Regenerate SQLC (if schema affects queries)
make extract-sqlc-schema-all
make sqlc

# 6. Verify
make check-build
make gotest
```

## HCL Schema Syntax

### Table Definition

```hcl
table "account_key" {
  schema = schema.watch

  column "id" {
    type           = bigint
    auto_increment = true
  }

  column "coin_type" {
    type = varchar(10)
  }

  column "account" {
    type = varchar(100)
  }

  column "created_at" {
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_account_key_coin_type" {
    columns = [column.coin_type]
  }
}
```

### Foreign Key

```hcl
foreign_key "fk_address_account" {
  columns     = [column.account_key_id]
  ref_columns = [table.account_key.column.id]
  on_delete   = CASCADE
  on_update   = CASCADE
}
```

### Unique Constraint

```hcl
index "idx_unique_account" {
  columns = [column.coin_type, column.account]
  unique  = true
}
```

## Auto-Generated Files

**DO NOT EDIT** files generated from HCL:

| Path                                 | Description         |
| ------------------------------------ | ------------------- |
| `tools/atlas/migrations/**/*.sql`    | Migration SQL files |
| `tools/atlas/migrations/*/atlas.sum` | Migration checksums |

Regenerate with `make atlas-dev-reset`.

## Best Practices

### Schema Changes

- Make backward compatible changes when possible
- Add indexes for frequently queried columns
- Use appropriate column types
- Define foreign keys with proper cascades

### Column Types

| MySQL Type   | HCL Type       |
| ------------ | -------------- |
| BIGINT       | `bigint`       |
| INT          | `int`          |
| VARCHAR(n)   | `varchar(n)`   |
| TEXT         | `text`         |
| TIMESTAMP    | `timestamp`    |
| BOOLEAN      | `bool`         |
| DECIMAL(p,s) | `decimal(p,s)` |

### Naming Conventions

- Table names: lowercase, snake_case
- Column names: lowercase, snake_case
- Index names: `idx_{table}_{column(s)}`
- Foreign key names: `fk_{table}_{reference}`

## Quick Checklist

- [ ] `make atlas-fmt` passes
- [ ] `make atlas-lint` passes
- [ ] Migration applies cleanly (`docker compose down -v && docker compose --profile mysql up -d`)
- [ ] `make sqlc` generates correctly (if schema changed)
- [ ] `make check-build` passes
- [ ] `make gotest` passes

## Related Documentation

- @docs/database/database.md - Database management guide
- @tools/atlas/atlas.hcl - Atlas configuration

## Related Skills

- `db-migration` - Full database migration workflow
- `go-development` - Go verification after SQLC generation
