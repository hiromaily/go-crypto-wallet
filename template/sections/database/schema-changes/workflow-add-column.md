### Step-by-Step Workflow

#### Scenario 1: Adding a New Column

**Example**: Add `email` column to `address` table in watch schema.

##### Step 1: Modify HCL Schema

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

##### Step 2: Format and Lint

```bash
make atlas-fmt
make atlas-lint
```

**Output**:

```
✓ Format complete
✓ No linting errors
```

##### Step 3: Regenerate Migrations

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

##### Step 4: Apply Migrations to Database

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

##### Step 5: Update SQLC Schema Files

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

##### Step 6: Regenerate SQLC Code

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

##### Step 7: Update Queries (if needed)

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

##### Step 8: Update Repository Code

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

##### Step 9: Verify Build

```bash
make go-lint
make check-build
```

##### Step 10: Test

```bash
# Unit tests
make go-test

# Integration tests (MySQL)
make integration-test

# E2E tests (SQLite)
make btc-e2e-reset P=1
```

---
