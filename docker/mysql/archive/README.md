# Archived Database Initialization Files

This directory contains the original database initialization files that were used before migrating to Atlas for database schema management.

## Directory Structure

```
archive/
├── init.d/
│   └── 01_init_all_schemas.sql    # Old initialization script (SQL-based)
├── scripts/
│   └── init.sh                     # Old initialization script (bash-based)
└── sqls/
    ├── definition_watch.sql        # Watch schema table definitions
    ├── definition_keygen.sql      # Keygen schema table definitions
    ├── definition_sign.sql         # Sign schema table definitions
    └── payment_request.sql         # Payment request table (watch schema)
```

## Migration to Atlas

All these files have been replaced by Atlas migrations located in `tools/atlas/migrations/`.

**Current schema management:**

- Schema initialization: `docker/mysql/init.d/01_init_all_schemas.sql` (creates schemas only)
- Schema definitions: `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL format, for reference)
- Migrations: `tools/atlas/migrations/{db_dialect}/*/` (SQL migrations, version-controlled)

**Current initialization process:**

1. MySQL container starts (`docker compose up wallet-mysql`)
2. `init.d/01_init_all_schemas.sql` creates empty schemas (watch, keygen, sign)
3. Atlas migration services automatically apply migrations when database is ready
4. All tables are created via Atlas migrations

**To manually apply migrations (if needed):**

```bash
# Migrations are automatically applied on docker compose --profile mysql up
# For manual re-application:
make atlas-migrate-docker
```

## Why Archived?

These files are kept for reference purposes:

- **Historical reference**: Understanding the original schema structure
- **Migration verification**: Comparing old SQL definitions with new Atlas migrations
- **Rollback scenarios**: Reference for emergency rollback (if needed)
- **Documentation**: Understanding the evolution of database initialization

## Important Notes

- **DO NOT** use these files for new database initialization
- **DO NOT** modify these files - they are archived for reference only
- Always use Atlas migrations (`tools/atlas/migrations/`) for schema management
- The current initialization script (`docker/mysql/init.d/01_init_all_schemas.sql`) only creates empty schemas and grants permissions
- Table creation is handled entirely by Atlas migrations
