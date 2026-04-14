# Database Schema Change Workflow

This document provides a comprehensive guide for modifying or adding database columns and schemas in the go-crypto-wallet project, which supports multiple database backends.

## Table of Contents

- [Overview](#overview)
- [Supported Databases](#supported-databases)
- [Quick Reference](#quick-reference)
- [Step-by-Step Workflow](#step-by-step-workflow)
- [Multi-Database Considerations](#multi-database-considerations)
- [Testing Schema Changes](#testing-schema-changes)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Overview

The project uses a **declarative schema management** approach with [Atlas](https://atlasgo.io/) and type-safe code generation with [sqlc](https://sqlc.dev/).

**Key Principle**: HCL schema files are the **single source of truth**. Never edit migration SQL files or generated code directly.

### Architecture

```
HCL Schemas (Source of Truth)
  ↓ Atlas generates
Database Migration Files (.sql)
  ↓ Applied to
Running Databases (MySQL/SQLite/PostgreSQL)
  ↓ Dumped to
SQLC Schema Files (.sql)
  ↓ sqlc generates
Type-Safe Go Code (sqlcgen/*.go)
  ↓ Used by
Repository Implementations
```

## Supported Databases

See [Database Architecture](./architecture.md) for supported backends (MySQL 8.4, SQLite, PostgreSQL). Database type is selected via `database.type` in wallet config files.

## Quick Reference

### Common Commands

| Task | Command |
|------|---------|
| Format HCL schemas | `make atlas-fmt` |
| Lint HCL schemas | `make atlas-lint` |
| Regenerate all migrations | `make atlas-dev-reset` |
| Apply migrations (local) | `make atlas-migrate` |
| Apply migrations (Docker) | `make atlas-migrate-docker` |
| Generate MySQL sqlc code | `make sqlc` |
| Generate SQLite sqlc code | `make sqlc-sqlite` |
| Generate PostgreSQL sqlc code | `make sqlc-postgres` |
| Generate all sqlc code | `make sqlc-all` |
| Verify build | `make check-build` |

### File Locations

```
tools/atlas/
├── schemas/                    # Source of truth (HCL)
│   ├── watch.hcl              # Watch wallet schema
│   ├── keygen.hcl             # Keygen wallet schema
│   └── sign.hcl               # Sign wallet schema
├── migrations/                 # Generated migrations (DO NOT EDIT)
│   ├── watch/*.sql            # Watch migrations
│   ├── keygen/*.sql           # Keygen migrations
│   └── sign/*.sql             # Sign migrations
└── atlas.hcl                  # Atlas configuration

tools/sqlc/
├── queries/
│   ├── mysql/                  # SQL queries - MySQL (? placeholders)
│   │   ├── address.sql
│   │   ├── btc_tx.sql
│   │   └── ...
│   └── postgres/             # SQL queries - PostgreSQL ($1,$2 placeholders)
│       └── ...
├── schemas/
│   ├── mysql/                  # MySQL schema files (extracted from DB)
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   ├── postgres/             # PostgreSQL schema files (extracted from DB)
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   └── sqlite/                 # SQLite schema files (manually converted)
│       ├── 01_watch.sql
│       ├── 02_keygen.sql
│       └── 03_sign.sql
├── sqlc.yml                   # MySQL sqlc config
├── sqlc_sqlite.yml            # SQLite sqlc config
└── sqlc_postgres.yml        # PostgreSQL sqlc config
```

## Step-by-Step Workflow

### Scenario 1: Adding a New Column

**Example**: Add `email` column to `address` table in watch schema.

#### Step 1: Modify HCL Schema

Edit `tools/atlas/schemas/{db_dialect}/watch.hcl`:

```hcl
table "address" {
  schema = schema.watch
  column "id" {
    type = bigint
    auto_increment = true
  }
  // ... existing columns ...

  // NEW: Add email column
  column "email" {
    type = varchar(255)
    null = true
  }

  primary_key {
    columns = [column.id]
  }
}
```

#### Step 2: Format and Lint

```bash
make atlas-fmt
make atlas-lint
```

**Output**:

```
✓ Format complete
✓ No linting errors
```

#### Step 3: Regenerate Migrations

```bash
make atlas-dev-reset
```

**What happens**:

1. Prompts for confirmation (deletes existing migrations)
2. Generates new migration files from HCL
3. Creates migrations for all schemas (watch, keygen, sign)

**Output**:

```
Are you sure you want to delete all migration files? [y/N]: y
✓ Deleted tools/atlas/migrations/watch/*
✓ Generated tools/atlas/migrations/watch/20240215120000.sql
✓ Checksum updated
```

#### Step 4: Apply Migrations to Database

```bash
# Stop and recreate database containers
docker compose down -v
docker compose up -d wallet-mysql

# Apply migrations
make atlas-migrate-docker
```

**Output**:

```
✓ Migrating to version 20240215120000 (1 migration)
  └─ watch: OK
  └─ keygen: OK
  └─ sign: OK
```

#### Step 5: Update SQLC Schema Files

```bash
# Extract schema from running MySQL database
make dump-schema-all
make extract-sqlc-schema-all

# Convert to SQLite format
# (This is manual - see tools/sqlc/schemas/sqlite/)
# Copy and modify MySQL schema with SQLite-specific changes

# PostgreSQL
# Convert to PostgreSQL format following data type mappings
```

**Data Type Conversion Reference**:

| MySQL | SQLite | PostgreSQL |
|-------|--------|------------|
| `VARCHAR(255)` | `TEXT` | `VARCHAR(255)` |
| `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| `TINYINT(1)` | `INTEGER` | `BOOLEAN` |
| `ENUM('a','b')` | `TEXT CHECK(col IN ('a','b'))` | `TEXT CHECK(col IN ('a','b'))` |
| `DATETIME` | `TEXT` | `TIMESTAMP` |
| `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |

#### Step 6: Regenerate SQLC Code

```bash
# Generate for MySQL
make sqlc

# Generate for SQLite
make sqlc-sqlite

# Generate for PostgreSQL
make sqlc-postgres
```

**Generated files**:

- `internal/infrastructure/database/mysql/sqlcgen/*.go`
- `internal/infrastructure/database/sqlite/sqlcgen/*.go`
- (Future) `internal/infrastructure/database/postgres/sqlcgen/*.go`

#### Step 7: Update Queries (if needed)

If you need to query the new column, edit `tools/sqlc/queries/mysql/address.sql`:

```sql
-- name: GetAddressWithEmail :one
SELECT id, wallet_address, email, created_at
FROM address
WHERE id = ?;
```

Then regenerate:

```bash
make sqlc
make sqlc-sqlite
```

#### Step 8: Update Repository Code

Update repository implementations to use the new field:

```go
// Example: internal/infrastructure/repository/watch/mysql/address_sqlc.go
func (r *AddressRepositorySqlc) GetByIDWithEmail(ctx context.Context, id int64) (*domain.Address, error) {
    addr, err := r.queries.GetAddressWithEmail(ctx, id)
    if err != nil {
        return nil, err
    }
    return convertToAddressWithEmail(addr), nil
}
```

#### Step 9: Verify Build

```bash
make go-lint
make check-build
```

#### Step 10: Test

```bash
# Unit tests
make go-test

# Integration tests (MySQL)
make integration-test

# E2E tests (SQLite)
make btc-e2e-reset P=1
```

---

### Scenario 2: Adding a New Table

**Example**: Add `audit_log` table to watch schema.

#### Step 1: Add Table Definition to HCL

Edit `tools/atlas/schemas/{db_dialect}/watch.hcl`:

```hcl
table "audit_log" {
  schema = schema.watch

  column "id" {
    type = bigint
    auto_increment = true
  }

  column "entity_type" {
    type = varchar(50)
    null = false
  }

  column "entity_id" {
    type = bigint
    null = false
  }

  column "action" {
    type = enum("create", "update", "delete")
    null = false
  }

  column "user_id" {
    type = varchar(100)
    null = true
  }

  column "changes" {
    type = text
    null = true
  }

  column "created_at" {
    type = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_entity" {
    columns = [column.entity_type, column.entity_id]
  }

  index "idx_created_at" {
    columns = [column.created_at]
  }
}
```

#### Step 2: Follow Same Workflow

Follow Steps 2-10 from Scenario 1.

**Additional Considerations**:

- Create corresponding queries in `tools/sqlc/queries/mysql/audit_log.sql`
- Implement repository interface and implementations for all databases
- Add integration tests for the new table

---

### Scenario 3: Modifying Existing Column

**Example**: Change `wallet_address` from `VARCHAR(255)` to `VARCHAR(500)`.

#### Step 1: Modify HCL Schema

Edit `tools/atlas/schemas/{db_dialect}/watch.hcl`:

```hcl
table "address" {
  // ...
  column "wallet_address" {
    type = varchar(500)  // Changed from 255
    null = false
  }
  // ...
}
```

#### Step 2-10: Follow Same Workflow

Same as Scenario 1, Steps 2-10.

**Important**: Atlas will generate an `ALTER TABLE` migration.

---

## Multi-Database Considerations

### Schema Parity Requirements

**CRITICAL**: All three databases (MySQL, SQLite, PostgreSQL) MUST maintain identical:

- Table names
- Column names
- Column order
- Primary keys
- Foreign keys (where supported)
- Index names (for consistency)

**Different**: Data types are mapped to database-specific equivalents.

### Data Type Mapping Strategy

When adding a new column, choose types carefully:

| Concept | MySQL | SQLite | PostgreSQL |
|---------|-------|--------|------------|
| Auto-incrementing ID | `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| Boolean flag | `TINYINT(1)` | `INTEGER (0/1)` | `BOOLEAN` |
| Enum values | `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| Decimal precision | `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |
| Timestamps | `DATETIME` | `TEXT (ISO8601)` | `TIMESTAMP` |
| Variable text | `VARCHAR(n)` | `TEXT` | `VARCHAR(n)` |
| Large text | `TEXT` | `TEXT` | `TEXT` |

### Enum Handling

**MySQL**:

```sql
coin ENUM('btc','bch','eth','xrp') NOT NULL
```

**SQLite & PostgreSQL**:

```sql
coin TEXT NOT NULL CHECK (coin IN ('btc','bch','eth','xrp'))
```

**Rationale**: TEXT with CHECK constraints provides schema evolution flexibility (no ACCESS EXCLUSIVE locks in PostgreSQL).

### Workflow for All Databases

When making schema changes:

1. **Update HCL schema** (single source of truth)
2. **Regenerate Atlas migrations** for MySQL
3. **Apply migrations** to MySQL database
4. **Extract MySQL schema** via dump
5. **Convert to SQLite schema** (manual data type mapping)
6. **Convert to PostgreSQL schema** (manual data type mapping)
7. **Regenerate sqlc code** for all three databases
8. **Update repositories** for all three databases (if needed)
9. **Test** with all three databases

### Repository Pattern

Each database has separate repository implementations:

```
internal/infrastructure/repository/
├── watch/
│   ├── mysql/
│   │   ├── address_sqlc.go
│   │   ├── btc_tx_sqlc.go
│   │   └── ...
│   ├── sqlite/
│   │   ├── address_sqlc.go
│   │   ├── btc_tx_sqlc.go
│   │   └── ...
│   └── postgres/           # Coming soon
│       ├── address_sqlc.go
│       ├── btc_tx_sqlc.go
│       └── ...
└── cold/
    ├── mysql/
    ├── sqlite/
    └── postgres/           # Coming soon
```

**DI Container** switches implementation based on `database.type`:

```go
// internal/di/container.go
func (c *Container) CreateAddressRepository() repository.AddressRepository {
    switch c.config.Database.Type {
    case "mysql":
        return mysql.NewAddressRepositorySqlc(c.mysqlDB, c.coinTypeCode)
    case "sqlite":
        return sqlite.NewAddressRepositorySqlc(c.sqliteDB, c.coinTypeCode)
    case "postgres":
        return postgres.NewAddressRepositorySqlc(c.postgresDB, c.coinTypeCode)
    default:
        panic("unsupported database type")
    }
}
```

## Testing Schema Changes

### 1. Local MySQL Test

```bash
# Reset database and apply migrations
docker compose down -v
docker compose up -d wallet-mysql
make atlas-migrate-docker

# Verify schema
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE new_table;"

# Run integration tests
make integration-test
```

### 2. SQLite E2E Test

```bash
# Test with SQLite (fast, no Docker)
make btc-e2e-reset P=1 DB=sqlite

# Verify database file
sqlite3 ./data/sqlite/btc/e2e.db "PRAGMA table_info(new_table);"
```

### 3. PostgreSQL Test (Coming Soon)

```bash
# Reset and apply migrations
docker compose down -v
docker compose up -d wallet-db-postgres
make atlas-migrate-docker-postgres

# Verify schema
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\d new_table"

# Run integration tests
make integration-test-postgres
```

### 4. Cross-Database Compatibility Test

Create a test that verifies schema consistency across all databases:

```go
// internal/infrastructure/repository/watch/compatibility_test.go
func TestSchemaParity(t *testing.T) {
    // Test that MySQL, SQLite, and PostgreSQL schemas are equivalent
    // - Same tables exist
    // - Same columns exist (ignoring type differences)
    // - Same primary keys
    // - Same foreign keys
}
```

## Common Patterns

### Pattern 1: Adding Created/Updated Timestamps

**MySQL**:

```hcl
column "created_at" {
  type = datetime
  default = sql("CURRENT_TIMESTAMP")
}

column "updated_at" {
  type = datetime
  null = true
}
```

**SQLite** (in `.sql` file):

```sql
created_at TEXT DEFAULT CURRENT_TIMESTAMP,
updated_at TEXT
```

**PostgreSQL**:

```sql
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP
```

### Pattern 2: Nullable vs NOT NULL

**Prefer NULL for optional fields**:

```hcl
column "email" {
  type = varchar(255)
  null = true  // Optional field
}
```

**Use NOT NULL for required fields**:

```hcl
column "wallet_address" {
  type = varchar(500)
  null = false  // Required field
}
```

### Pattern 3: Foreign Keys

**MySQL/PostgreSQL** (HCL):

```hcl
table "btc_tx_input" {
  // ...
  column "tx_id" {
    type = bigint
    null = false
  }

  foreign_key "fk_btc_tx_input_tx" {
    columns     = [column.tx_id]
    ref_columns = [table.btc_tx.column.id]
    on_delete   = CASCADE
    on_update   = CASCADE
  }
}
```

**SQLite** (limited FK support):

```sql
-- SQLite supports FK but they're not enforced by default
-- PRAGMA foreign_keys = ON; required
tx_id INTEGER NOT NULL,
FOREIGN KEY (tx_id) REFERENCES btc_tx(id) ON DELETE CASCADE
```

### Pattern 4: Indexes for Performance

**HCL**:

```hcl
table "address" {
  // ...
  index "idx_coin_account" {
    columns = [column.coin, column.account]
  }

  index "idx_is_allocated" {
    columns = [column.is_allocated]
  }
}
```

### Pattern 5: Unique Constraints

**HCL**:

```hcl
table "address" {
  // ...
  column "wallet_address" {
    type = varchar(500)
    null = false
  }

  index "idx_wallet_address_unique" {
    unique = true
    columns = [column.wallet_address]
  }
}
```

## Troubleshooting

### Issue: Atlas Migration Fails

**Symptoms**:

```
Error: atlas migrate apply failed
```

**Solutions**:

1. **Check HCL syntax**:

   ```bash
   make atlas-lint
   ```

2. **Verify database is running**:

   ```bash
   docker compose ps wallet-mysql
   ```

3. **Check migration history**:

   ```bash
   make atlas-status-docker
   ```

4. **Reset and retry**:

   ```bash
   docker compose down -v
   docker compose up -d wallet-mysql
   make atlas-dev-reset
   make atlas-migrate-docker
   ```

### Issue: SQLC Generation Fails

**Symptoms**:

```
Error: sqlc generate failed
```

**Solutions**:

1. **Check schema file syntax**:

   ```bash
   # Validate MySQL syntax
   docker compose exec wallet-mysql mysql -uroot -proot watch < tools/sqlc/schemas/mysql/01_watch.sql
   ```

2. **Check query file syntax**:

   ```bash
   # Run sqlc with verbose output
   cd tools/sqlc && sqlc generate --experimental
   ```

3. **Verify engine configuration**:

   ```yaml
   # tools/sqlc/sqlc_mysql.yml for MySQL
   version: "2"
   sql:
     - engine: "mysql"  # Ensure correct engine
       queries: "queries"
       schema: "schemas"
   ```

### Issue: Schema Mismatch Between Databases

**Symptoms**:

- Tests pass with MySQL but fail with SQLite
- Repository code works differently across databases

**Solutions**:

1. **Compare schema files**:

   ```bash
   # Compare table structures
   diff -u tools/sqlc/schemas/mysql/01_watch.sql tools/sqlc/schemas/sqlite/01_watch.sql
   ```

2. **Verify data type mappings**:
   - Review [Data Type Mapping Strategy](#data-type-mapping-strategy)
   - Ensure ENUM → CHECK constraint conversion is correct

3. **Check sqlc generated models**:

   ```bash
   # Compare generated models
   diff -u internal/infrastructure/database/mysql/sqlcgen/models.go \
           internal/infrastructure/database/sqlite/sqlcgen/models.go
   ```

### Issue: Migration Conflict

**Symptoms**:

```
Error: migration checksum mismatch
```

**Solutions**:

1. **Regenerate from scratch**:

   ```bash
   make atlas-dev-reset
   ```

2. **Clear migration history**:

   ```bash
   docker compose exec wallet-mysql mysql -uroot -proot watch \
     -e "DROP TABLE IF EXISTS atlas_schema_revisions;"
   make atlas-migrate-docker
   ```

## Best Practices

### 1. Always Use HCL as Source of Truth

✅ **DO**:

```bash
# Edit HCL file
vim tools/atlas/schemas/watch.hcl

# Regenerate migrations
make atlas-dev-reset
```

❌ **DON'T**:

```bash
# Never edit migration files directly
vim tools/atlas/migrations/watch/20240215120000.sql  # WRONG!
```

### 2. Test Locally Before Committing

```bash
# Complete local test cycle
docker compose down -v
docker compose up -d wallet-mysql
make atlas-migrate-docker
make sqlc-all
make go-lint
make check-build
make go-test
```

### 3. Keep Migrations Small and Focused

✅ **Good**: One logical change per commit

- Add `email` column to `address` table
- Create `audit_log` table
- Add index on `created_at`

❌ **Bad**: Multiple unrelated changes

- Add `email` column, create `audit_log` table, modify `btc_tx` structure

### 4. Document Schema Changes

Add comments to HCL files:

```hcl
// Added 2024-02-15: Email notifications feature (#566)
column "email" {
  type = varchar(255)
  null = true
}
```

### 5. Backward Compatibility

When modifying schemas:

- **Adding columns**: Always make them nullable or provide default values
- **Removing columns**: Deprecated first, remove in next major version
- **Changing types**: Ensure data can be migrated without loss

### 6. Security for Sensitive Columns

For sensitive data (seeds, private keys):

- Use appropriate encryption
- Never log sensitive column values
- Consider separate schemas for hot/cold storage

```hcl
// Keygen schema (OFFLINE - sensitive data)
table "seed" {
  schema = schema.keygen

  column "encrypted_seed" {
    type = text
    null = false
  }

  // Never store unencrypted
}
```

### 7. Index Strategy

Add indexes for:

- Foreign key columns
- Frequently queried columns
- Columns used in WHERE clauses
- Columns used in ORDER BY

```hcl
index "idx_coin_account" {
  columns = [column.coin, column.account]
}
```

### 8. Version Control

Commit in this order:

1. HCL schema changes
2. Generated migration files
3. Updated SQLC schema files
4. Generated SQLC code
5. Repository implementations
6. Tests

### 9. CI/CD Integration

Ensure CI tests schema changes:

```yaml
# .github/workflows/test.yml
- name: Test schema migrations
  run: |
    docker compose up -d wallet-mysql
    make atlas-migrate-docker
    make sqlc-all
    make check-build
    make go-test
```

### 10. Documentation

Update these files when schema changes:

- `docs/database/architecture.md` - Database architecture
- This file - If new patterns emerge
- `tools/atlas/README.md` - Atlas-specific details
- Schema-specific docs (if applicable)

## Additional Resources

- [Atlas Documentation](https://atlasgo.io/)
- [SQLC Documentation](https://sqlc.dev/)
- [Database Architecture](./architecture.md)
- [Code Generation Guidelines](../guidelines/code-generation.md)
- PostgreSQL Integration Spec (`.kiro/specs/postgres-integration/`)

---

**Last Updated**: 2024-02-15
**Maintained By**: go-crypto-wallet team
