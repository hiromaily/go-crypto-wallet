---
name: db-migration
description: Database schema and migration workflow. Use when modifying database schemas in tools/atlas/ or SQLC queries in tools/sqlc/.
---

# Database Migration Workflow

Workflow for database schema and migration changes.

## Prerequisites

**Use `git-workflow` Skill** for branch, commit, and PR workflow.

**Use `go-development` Skill** for Go code verification after SQLC generation.

## Applicable Files

| Path | Description |
|------|-------------|
| `tools/atlas/schemas/` | HCL schema definitions |
| `tools/sqlc/` | SQLC query definitions |
| `docker/*.sql` | Initial data/schemas |

## Workflow

### 1. Modify Schema

Edit HCL files in `tools/atlas/schemas/`:

```hcl
table "new_table" {
  schema = schema.wallet
  column "id" {
    type = bigint
    auto_increment = true
  }
  primary_key {
    columns = [column.id]
  }
}
```

### 2. Format and Lint

```bash
make atlas-fmt
make atlas-lint
```

### 3. Generate Migrations

```bash
make atlas-dev-reset
```

### 4. Test Migration

```bash
# Reset database and apply
docker compose down -v
docker compose up -d

# Verify tables
docker compose exec mysql mysql -u root -p -e "SHOW TABLES;"
```

### 5. Regenerate SQLC

```bash
make sqlc
```

### 6. Verify Go Code

```bash
make check-build
make gotest
```

## Guidelines

### Schema Changes

- [ ] Backward compatible (if possible)
- [ ] Indexes for frequently queried columns
- [ ] Foreign keys with proper cascades
- [ ] Appropriate column types

### SQLC Queries

- [ ] Named parameters
- [ ] Proper return types
- [ ] Efficient queries (avoid N+1)

## Verification Checklist

- [ ] `make atlas-fmt` passes
- [ ] `make atlas-lint` passes
- [ ] Migration applies cleanly
- [ ] `make sqlc` generates correctly
- [ ] `make check-build` passes
- [ ] `make gotest` passes

## Auto-Generated Files

**DO NOT EDIT** directly:
- `internal/infrastructure/database/sqlc/*.go`

Regenerate with `make sqlc` instead.

## Commit Format

```
feat(db): {brief description}

- {schema change 1}
- {schema change 2}

Closes #{issue_number}
```

## Related Skills

- `git-workflow` - Branch, commit, PR workflow
- `go-development` - Go verification after SQLC generation
