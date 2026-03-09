# Database Management

This document provides a quick reference for database schema management and SQLC code generation in the go-crypto-wallet project.

**📘 For detailed schema change workflows, see [Database Schema Changes Guide](schema-changes.md)**

## Supported Databases

The project supports three database backends:

| Database | Status | Use Case | Configuration |
|----------|--------|----------|---------------|
| **MySQL 8.4** | ✅ Production | Production, full integration testing | Docker container |
| **SQLite** | ✅ Production | E2E testing, CI/CD, lightweight testing | Local file |
| **PostgreSQL 18.2** | 🚧 In Development | Production alternative, advanced features | Docker container |

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

**🚀 Quick Workflow**:

1. Edit HCL schema (`tools/atlas/schemas/{db_dialect}/*.hcl`)
2. Format and lint (`make atlas-fmt && make atlas-lint`)
3. Regenerate migrations (`make atlas-dev-reset`)
4. Apply to database (`docker compose down -v && docker compose --profile mysql up`)
5. Update SQLC schemas for all databases (MySQL, SQLite, PostgreSQL)
6. Regenerate SQLC code (`make sqlc-all`)
7. Verify build (`make check-build`)

**📘 For complete step-by-step workflow with examples, see [Database Schema Changes Guide](schema-changes.md)**

### Schema Files (Source of Truth)

There are 3 HCL schema files corresponding to each wallet type:

- `tools/atlas/schemas/{db_dialect}/watch.hcl` - Watch wallet schema (online wallet)
- `tools/atlas/schemas/{db_dialect}/keygen.hcl` - Keygen wallet schema (offline, key generation)
- `tools/atlas/schemas/{db_dialect}/sign.hcl` - Sign wallet schema (offline, signing)

**CRITICAL**: These HCL files are the **single source of truth**. Never edit migration SQL files or generated code directly.

### Multi-Database Workflow

When making schema changes, you must update schemas for **all supported databases**:

1. **MySQL**: Automatically updated via Atlas migrations
2. **SQLite**: Manually convert MySQL schema with type mappings
3. **PostgreSQL** *(coming soon)*: Manually convert MySQL schema with type mappings

**Data Type Mapping Quick Reference**:

| MySQL | SQLite | PostgreSQL |
|-------|--------|------------|
| `AUTO_INCREMENT` | `AUTOINCREMENT` | `BIGSERIAL` |
| `TINYINT(1)` | `INTEGER` | `BOOLEAN` |
| `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| `DATETIME` | `TEXT` | `TIMESTAMP` |
| `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |

See [Database Schema Changes Guide](schema-changes.md) for complete mapping table.

### Important Principles

- **HCL as source of truth** - Always modify HCL files first
- **Schema parity** - All databases must have identical table/column names
- **Test all databases** - Verify changes work with MySQL, SQLite, and PostgreSQL
- **Atomic commits** - Commit HCL changes, migrations, and generated code together

## Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)
**Source**: `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL schema files)
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/{database}/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/{database}/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/{database}/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: These files are auto-generated and should **NEVER** be edited manually.

## SQLC Schema Files (from Database Dumps)

**Tool**: Custom shell script (`scripts/db/extract-sqlc-schema.sh`)
**Source**: MySQL database dumps (`data/dump/sql/dump_*.sql`)
**Command**: `make extract-sqlc-schema-all` (or individual: `make extract-sqlc-schema-watch`, `make extract-sqlc-schema-keygen`, `make extract-sqlc-schema-sign`)

**Generated Files**:

- `tools/sqlc/schemas/{db_dialect}/01_watch.sql` - Watch schema for SQLC
- `tools/sqlc/schemas/{db_dialect}/02_keygen.sql` - Keygen schema for SQLC
- `tools/sqlc/schemas/{db_dialect}/03_sign.sql` - Sign schema for SQLC

**Note**: These schema files are extracted from MySQL database dumps. The source of truth is the Atlas HCL files (`tools/atlas/schemas/{db_dialect}/*.hcl`). To update schemas, modify the HCL files and run the database migration flow.

## Database Code (SQLC)

**Tool**: [sqlc](https://sqlc.dev/)
**Source**: `tools/sqlc/schemas/{db_dialect}/*.sql` (auto-generated) and `tools/sqlc/queries/{db_dialect}/*.sql` (manually edited)
**Command**: `make sqlc` (or `cd tools/sqlc && sqlc generate`)

### MySQL SQLC

**Generated Files**:

- `internal/infrastructure/database/mysql/sqlcgen/*.go` (15 files)
  - `models.go` - Database models
  - `db.go` - Database connection code
  - `*.sql.go` - Query functions

### SQLite SQLC

**Config**: `tools/sqlc/sqlc_sqlite.yml`
**Schema**: `tools/sqlc/schemas/sqlite/*.sql`
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

**Manual Editing**: The SQL query files in `tools/sqlc/queries/mysql/*.sql` are **manually edited** and should be modified when adding new database queries.

**Location**: `tools/sqlc/queries/mysql/`

**Workflow:**

1. Write SQL queries in `tools/sqlc/queries/mysql/*.sql`
2. Run `make sqlc` to generate Go code
3. Use the generated code in your repositories

## See Also

- [Code Generation Guidelines](/guidelines/code-generation) - Complete overview of all code generation tools
- [Architecture Guidelines](architecture.md) - Infrastructure layer guidelines for repositories
- [Workflow Guidelines](/guidelines/workflow) - Dependency management and verification commands
