# Atlas Database Migration Tool

This directory contains Atlas configuration and migration files for managing database schemas in the go-crypto-wallet project.

## Overview

Atlas is a modern database schema migration tool written in Go. It provides:
- Version-controlled migrations
- Migration history tracking
- Rollback capabilities
- Schema validation

## Installation

Install Atlas CLI using Homebrew (macOS):

```bash
brew install arigaio/tap/atlas
```

Alternatively, you can install using Go:

```bash
go install ariga.io/atlas/cmd/atlas@latest
```

Verify installation:

```bash
atlas version
```

## Project Structure

```
tools/atlas/
├── atlas.hcl              # Atlas configuration file
├── schemas/                # HCL schema definitions (declarative)
│   ├── watch.hcl          # Watch schema definition
│   ├── keygen.hcl         # Keygen schema definition
│   └── sign.hcl           # Sign schema definition
├── migrations/            # SQL migration files
│   ├── watch/             # Watch schema migrations
│   ├── keygen/            # Keygen schema migrations
│   └── sign/              # Sign schema migrations
└── README.md              # This file
```

## Schemas

The project uses three separate MySQL schemas:

- **watch**: Online wallet data (addresses, transactions, payment requests)
- **keygen**: Key generation data (seeds, account keys, full public keys)
- **sign**: Signing wallet data (auth account keys, seeds)

## Usage

### HCL Schema Management (Declarative)

The project uses HCL (HashiCorp Configuration Language) files for declarative schema management. HCL files define the desired state of the database schema.

#### Apply HCL Schema

Apply HCL schema definitions directly to the database:

```bash
make atlas-schema-apply
```

This will apply the HCL schema files (`schemas/*.hcl`) to their respective databases.

#### Show Schema Diff

Compare the current database state with HCL schema definition:

```bash
make atlas-schema-diff SCHEMA=watch
make atlas-schema-diff SCHEMA=keygen
make atlas-schema-diff SCHEMA=sign
```

#### Generate Migration from HCL Diff

Generate a migration file based on the difference between database and HCL schema:

```bash
make atlas-schema-diff-migration SCHEMA=watch
```

This creates a new migration file that will bring the database in line with the HCL schema.

### Apply Migrations

Apply all pending migrations for all schemas:

```bash
make atlas-migrate
```

Apply migrations in Docker environment:

```bash
make atlas-migrate-docker
```

### Check Migration Status

View migration status for all schemas:

```bash
make atlas-status
```

### Rollback Migrations

Rollback the last migration for a specific schema:

```bash
make atlas-rollback SCHEMA=watch
make atlas-rollback SCHEMA=keygen
make atlas-rollback SCHEMA=sign
```

### Validate Migrations

Validate all migration files:

```bash
make atlas-validate
```

### Create New Migration

Create a new migration file:

```bash
make atlas-new SCHEMA=watch NAME=add_new_table
make atlas-new SCHEMA=keygen NAME=update_account_key
make atlas-new SCHEMA=sign NAME=add_index
```

## Migration File Naming

Atlas uses timestamp-based naming for migration files:

- Format: `YYYYMMDDHHMMSS_description.sql`
- Example: `20240101000000_initial_watch_schema.sql`

## Manual Atlas Commands

If you need to run Atlas commands directly:

### Watch Schema

```bash
# Apply migrations
atlas migrate apply \
  --dir file://tools/atlas/migrations/watch \
  --url "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"

# Check status
atlas migrate status \
  --dir file://tools/atlas/migrations/watch \
  --url "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
```

### Keygen Schema

```bash
atlas migrate apply \
  --dir file://tools/atlas/migrations/keygen \
  --url "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
```

### Sign Schema

```bash
atlas migrate apply \
  --dir file://tools/atlas/migrations/sign \
  --url "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
```

## Migration History

Atlas automatically creates a migration history table (`atlas_schema_migrations`) in each schema to track applied migrations. This table should not be modified manually.

## Best Practices

1. **Always validate migrations** before applying:
   ```bash
   make atlas-validate
   ```

2. **Check migration status** before applying:
   ```bash
   make atlas-status
   ```

3. **Test migrations** on a development database before applying to production

4. **Never modify existing migration files** - create new migrations instead

5. **Keep migrations small and focused** - one logical change per migration

6. **Document complex migrations** with comments in the SQL file

## HCL Schema vs SQL Migrations

The project supports both approaches:

### HCL Schema (Declarative)
- **Location**: `tools/atlas/schemas/*.hcl`
- **Purpose**: Define the desired state of the database schema
- **Usage**: Use `atlas schema apply` to apply directly, or generate migrations from diffs
- **Benefits**: 
  - Single source of truth for schema definition
  - Easy to see the complete schema structure
  - Can generate migrations automatically from diffs

### SQL Migrations (Versioned)
- **Location**: `tools/atlas/migrations/*/`
- **Purpose**: Version-controlled, incremental schema changes
- **Usage**: Use `atlas migrate apply` to apply migrations in order
- **Benefits**:
  - Full migration history
  - Can rollback changes
  - Better for production deployments

### Workflow

**Recommended workflow for schema changes:**

1. **Update HCL schema file** (`schemas/*.hcl`) with desired changes
2. **Generate migration** from diff:
   ```bash
   make atlas-schema-diff-migration SCHEMA=watch
   ```
3. **Review the generated migration** file
4. **Apply the migration**:
   ```bash
   make atlas-migrate
   ```
5. **Update sqlc schema files** if needed for code generation
6. **Run sqlc generate** to update generated code

## Integration with sqlc

Atlas migrations and HCL schemas work alongside sqlc schema files:

- **Atlas HCL schemas**: Declarative schema definitions (`tools/atlas/schemas/`)
- **Atlas migrations**: Version-controlled schema changes (`tools/atlas/migrations/`)
- **sqlc schemas**: Used for code generation (`tools/sqlc/schemas/`)

When creating new tables or modifying existing ones:
1. Update the HCL schema file (`schemas/*.hcl`)
2. Generate a migration from the diff (or apply directly)
3. Apply the migration
4. Update sqlc schema files if needed for code generation
5. Run `sqlc generate` to update generated code

## Troubleshooting

### Migration Fails

If a migration fails:
1. Check the error message
2. Verify the database connection
3. Check if the schema exists
4. Review migration file syntax

### Rollback Issues

If rollback fails:
1. Check migration history: `make atlas-status`
2. Verify the migration file exists
3. Check database connection

### Connection Issues

If you can't connect to the database:
1. Verify MySQL is running: `docker compose ps wallet-db`
2. Check connection string in Makefile targets
3. Verify credentials are correct

## Related Documentation

- [Database Architecture Documentation](../../docs/development/database.md)
- [Atlas Official Documentation](https://atlasgo.io/)
- [Atlas MySQL Guide](https://atlasgo.io/guides/mysql)

