# Database Management

This document describes the database schema management workflow and SQLC code generation for the go-crypto-wallet project.

## Supported Databases

The project supports two database backends:

| Database | Use Case | Configuration |
|----------|----------|---------------|
| **MySQL** | Production, full integration testing | Docker container |
| **SQLite** | E2E testing, CI/CD, lightweight testing | Local file |

### Database Type Configuration

Database type is configured in wallet YAML files:

```yaml
database:
  type: "sqlite"  # or "mysql"
  mysql:
    host: "127.0.0.1:3306"
    dbname: "watch"
    user: "hiromaily"
    pass: "hiromaily"
  sqlite:
    path: "./data/sqlite/btc/e2e.db"
    debug: true
```

Or override via environment variables:

```bash
export WALLET_DATABASE_TYPE=sqlite
export WALLET_DATABASE_SQLITE_PATH=./data/sqlite/btc/e2e.db
```

## Database Schema Changes

This project uses [Atlas](https://atlasgo.io/) for database schema management with HCL (HashiCorp Configuration Language) as the source of truth.

### Schema Files

There are 3 schema files corresponding to each wallet type:

- `tools/atlas/schemas/watch.hcl` - Watch wallet schema (online wallet)
- `tools/atlas/schemas/keygen.hcl` - Keygen wallet schema (offline, key generation)
- `tools/atlas/schemas/sign.hcl` - Sign wallet schema (offline, signing)

### How to Change Database Schema

**Step 1: Modify the HCL schema file**

Edit the appropriate `.hcl` file in `tools/atlas/schemas/` directory. These files are the single source of truth for database schema.

**Step 2: Format and lint the schema files**

Run the following commands to format and validate the HCL schema files:

```bash
make atlas-fmt
make atlas-lint
```

This will:

- Format all HCL schema files for consistency
- Validate the schema syntax and structure
- Ensure no errors exist before generating migrations

**Step 3: Regenerate migration files**

Run the following command to regenerate migration files from scratch:

```bash
make atlas-dev-reset
```

This command will:

- Delete all existing migration files
- Generate new migration files from the HCL schemas
- Prompt for confirmation before proceeding

**Step 4: Verify the migration**

Test the migration by recreating the database:

```bash
docker compose down -v
docker compose up
```

This will:

- Remove existing database volumes (`-v` flag)
- Start fresh containers and apply migrations
- Verify that no errors occur during migration

**Step 5: Regenerate SQLC code (if needed)**

If the schema changes affect queries, regenerate SQLC code:

```bash
make sqlc
```

**Step 6: Verify the build**

```bash
make check-build
```

### Important Notes

- **Always modify HCL files first** - Never edit migration SQL files directly
- **HCL files are the source of truth** - Migration files are auto-generated from HCL
- **Test locally before committing** - Always run the full `docker compose down -v && docker compose up` cycle

## Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)
**Source**: `tools/atlas/schemas/*.hcl` (HCL schema files)
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: These files are auto-generated and should **NEVER** be edited manually.

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

### MySQL SQLC

**Generated Files**:

- `internal/infrastructure/database/mysql/sqlcgen/*.go` (15 files)
  - `models.go` - Database models
  - `db.go` - Database connection code
  - `*.sql.go` - Query functions

### SQLite SQLC

**Config**: `tools/sqlc/sqlc_sqlite.yml`
**Schema**: `tools/sqlc/schemas_sqlite/*.sql`
**Command**: `make sqlc-sqlite`

**Generated Files**:

- `internal/infrastructure/database/sqlite/sqlcgen/*.go`

**Schema Conversion Notes** (MySQL → SQLite):

| MySQL | SQLite |
|-------|--------|
| `ENUM('a','b','c')` | `TEXT CHECK(column IN ('a','b','c'))` |
| `AUTO_INCREMENT` | `AUTOINCREMENT` |
| `TINYINT(1)` | `INTEGER` |
| `DATETIME DEFAULT CURRENT_TIMESTAMP` | `TEXT DEFAULT CURRENT_TIMESTAMP` |

**Note**: SQLC generates type-safe Go code from SQL queries and schemas.

## SQLC Query Files

**Manual Editing**: The SQL query files in `tools/sqlc/queries/*.sql` are **manually edited** and should be modified when adding new database queries.

**Location**: `tools/sqlc/queries/`

**Workflow:**

1. Write SQL queries in `tools/sqlc/queries/*.sql`
2. Run `make sqlc` to generate Go code
3. Use the generated code in your repositories

## See Also

- [Code Generation Guidelines](code-generation.md) - Complete overview of all code generation tools
- [Architecture Guidelines](architecture.md) - Infrastructure layer guidelines for repositories
- [Workflow Guidelines](workflow.md) - Dependency management and verification commands
