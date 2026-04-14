# Database Management

This document provides a quick reference for database schema management and SQLC code generation in the go-crypto-wallet project.

**For detailed schema change workflows, see [Database Schema Changes Guide](./schema-changes.md)**

## Supported Databases

See [Database Architecture](./architecture.md) for the full list of supported backends, schema design, and setup instructions.

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

**Quick Workflow**:

1. Edit HCL schema (`tools/atlas/schemas/{db_dialect}/*.hcl`)
2. Format and lint (`make atlas-fmt && make atlas-lint`)
3. Regenerate migrations (`make atlas-dev-reset`)
4. Apply to database (`docker compose down -v && docker compose --profile mysql up`)
5. Update SQLC schemas for all databases (MySQL, SQLite, PostgreSQL)
6. Regenerate SQLC code (`make sqlc-all`)
7. Verify build (`make check-build`)

**For complete step-by-step workflow with examples, see [Database Schema Changes Guide](./schema-changes.md)**

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

See [Database Schema Changes Guide](./schema-changes.md) for complete mapping table.

### Important Principles

- **HCL as source of truth** - Always modify HCL files first
- **Schema parity** - All databases must have identical table/column names
- **Test all databases** - Verify changes work with MySQL, SQLite, and PostgreSQL
- **Atomic commits** - Commit HCL changes, migrations, and generated code together

## Database Migrations (Atlas)

See [Code Generation Guide](../guidelines/code-generation.md) for full Atlas workflow and commands.

## SQLC Schema Files (from Database Dumps)

See [Code Generation Guide](../guidelines/code-generation.md) for the full schema extraction process.

## Database Code (SQLC)

See [Code Generation Guide](../guidelines/code-generation.md) for SQLC configuration and generation commands.

## SQLC Query Files

**Manual Editing**: The SQL query files in `tools/sqlc/queries/mysql/*.sql` are **manually edited** and should be modified when adding new database queries.

**Location**: `tools/sqlc/queries/mysql/`

**Workflow:**

1. Write SQL queries in `tools/sqlc/queries/mysql/*.sql`
2. Run `make sqlc` to generate Go code
3. Use the generated code in your repositories

## See Also

- [Code Generation Guidelines](../guidelines/code-generation.md) - Complete overview of all code generation tools
- [Architecture Guidelines](./architecture.md) - Infrastructure layer guidelines for repositories
- [Workflow Guidelines](../guidelines/workflow.md) - Dependency management and verification commands
