# Archived SQL Schema Files

This directory contains the original SQL schema definition files that were used before migrating to Atlas for database schema management.

## Files

- `definition_watch.sql` - Watch schema table definitions
- `definition_keygen.sql` - Keygen schema table definitions
- `definition_sign.sql` - Sign schema table definitions
- `payment_request.sql` - Payment request table (watch schema)

## Migration to Atlas

These files have been replaced by Atlas migrations located in `tools/atlas/migrations/`.

**Current schema management:**

- Schema definitions: `tools/atlas/schemas/*.hcl` (HCL format)
- Migrations: `tools/atlas/migrations/*/` (SQL migrations)

**To apply migrations:**

```bash
make atlas-migrate-docker
```

## Why Archived?

These files are kept for reference purposes:

- Historical reference
- Migration verification
- Rollback scenarios (if needed)

**Note:** These files should not be used for new database initialization. Always use Atlas migrations for schema management.
