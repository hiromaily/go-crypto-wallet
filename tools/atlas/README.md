# Atlas Database Migration Tool

This directory contains Atlas configuration and migration files for managing database schemas in the go-crypto-wallet project.

## Overview

Atlas is a modern database schema migration tool written in Go. It provides:

- Version-controlled migrations
- Migration history tracking
- Rollback capabilities
- Schema validation

**Supported databases**: MySQL and PostgreSQL only.

> **Note**: SQLite is NOT managed by Atlas. SQLite schemas are manually maintained in `tools/sqlc/schemas/sqlite/` and used only for SQLC code generation. See [Database Management](../../docs/database/db-management.md) for details.

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
├── atlas.hcl              # Atlas configuration file (with lint rules)
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

## Atlas Configuration

The `atlas.hcl` file contains:

- **Three environments**: `local_mysql_watch`, `local_mysql_keygen`, `local_mysql_sign`
- **Diff configuration**: Protects against destructive changes (drop schema/table/column)
- **Lint configuration**: Validates schema for:
  - Destructive operations (errors on dangerous changes)
  - Naming conventions (lowercase with underscores)

## Schemas

The project uses three separate MySQL schemas:

- **watch**: Online wallet data (addresses, transactions, payment requests)
- **keygen**: Key generation data (seeds, account keys, full public keys)
- **sign**: Signing wallet data (auth account keys, seeds)

## Development Workflows

Atlas supports two distinct workflows depending on your development phase:

### Development Workflow (Schema-First Approach)

**When to use**: Active development phase with frequent schema changes.

**Characteristics**:

- Modify HCL schema files (`schemas/*.hcl`) directly
- Regenerate migrations from scratch using `make atlas-dev-reset`
- Single migration file represents the current complete schema state
- Fast iteration on schema design
- `docker compose --profile mysql up` automatically applies the latest migration

**Workflow**:

1. Modify HCL schema files in `schemas/` directory
2. Format and lint schemas:

   ```bash
   make atlas-fmt
   make atlas-lint
   ```

3. Regenerate migrations from scratch (WARNING: deletes existing migrations):

   ```bash
   make atlas-dev-reset
   ```

4. Apply migrations:

   ```bash
   docker compose --profile mysql up
   # or
   make atlas-migrate-apply
   ```

**Alternative: Direct Schema Apply** (bypasses migrations):

For rapid prototyping, you can apply HCL schemas directly to the database:

```bash
# Apply all schemas
make atlas-schema-apply

# Apply specific schema
make atlas-schema-apply-one SCHEMA=watch
```

**Note**: To see what changes will be applied, use `atlas-migrate-diff` to generate a migration file which shows the SQL changes.

### Production Workflow (Migration History Approach)

**When to use**: Production or stable development phase where migration history matters.

**Characteristics**:

- Incremental migrations with version control
- Each migration represents a specific change
- Migration history is preserved
- Supports rollback capabilities
- Safe for production deployments

**Workflow**:

1. Initialize production migration history (one-time):

   ```bash
   make atlas-prod-init
   ```

2. Modify HCL schema files in `schemas/` directory
3. Format and lint schemas:

   ```bash
   make atlas-fmt
   make atlas-lint
   ```

4. Generate incremental migration:

   ```bash
   make atlas-migrate-diff SCHEMA=watch NAME=add_new_column
   ```

5. Review generated migration in `migrations/<schema>/`
6. Apply migration:

   ```bash
   make atlas-migrate-apply
   ```

7. Check migration status:

   ```bash
   make atlas-migrate-status
   ```

### Switching Between Workflows

**Development → Production**:

```bash
# Initialize migration history from current state
make atlas-prod-init
```

**Production → Development** (for prototyping):

```bash
# WARNING: This deletes all migrations and databases!
# Only use in non-production environments
make atlas-dev-clean
```

## Available Make Targets

Run `make atlas-help` to see all available targets:

```bash
make atlas-help
```

**Quick Reference**:

- **Format & Lint**: `make atlas-fmt`, `make atlas-lint` (requires Docker)
- **Schema Management**: `make atlas-schema-apply`, `make atlas-schema-apply-one`
- **Migration Management**: `make atlas-migrate-status`, `make atlas-migrate-apply`, `make atlas-migrate-diff`
- **Development**: `make atlas-dev-reset`, `make atlas-dev-clean`
- **Production**: `make atlas-prod-init`
- **Utilities**: `make atlas-validate`, `make atlas-help`

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

**Docker environment (Automatic)**:

Migrations are **automatically applied** when you run `docker compose --profile mysql up`. The `compose.yaml` file defines three migration services:

- `wallet-mysql-migrate-watch`: Applies migrations for the watch schema
- `wallet-mysql-migrate-keygen`: Applies migrations for the keygen schema
- `wallet-mysql-migrate-sign`: Applies migrations for the sign schema

These services:

- Automatically wait for the database to be ready (using `depends_on` with `service_healthy` condition)
- Run in the same network as the database
- Have access to the `tools/atlas` directory via volume mount
- Use the official `arigaio/atlas:1.0.0` Docker image
- Exit after completion (restart: "no")

**Manual execution** (if needed):

```bash
make atlas-migrate-docker
```

**Local environment** (requires Atlas CLI installed):

```bash
make atlas-migrate
```

### Check Migration Status

View migration status for all schemas:

**Local environment**:

```bash
make atlas-status
```

**Docker environment**:

```bash
make atlas-status-docker
```

### Rollback Migrations

Rollback the last migration for a specific schema:

**Local environment**:

```bash
make atlas-rollback SCHEMA=watch
make atlas-rollback SCHEMA=keygen
make atlas-rollback SCHEMA=sign
```

**Docker environment**:

```bash
make atlas-rollback-docker SCHEMA=watch
make atlas-rollback-docker SCHEMA=keygen
make atlas-rollback-docker SCHEMA=sign
```

### Validate Migrations

Validate all migration files:

**Local environment**:

```bash
make atlas-validate
```

**Docker environment**:

```bash
make atlas-validate-docker
```

### Create New Migration

Create a new migration file:

**Local environment**:

```bash
make atlas-new SCHEMA=watch NAME=add_new_table
make atlas-new SCHEMA=keygen NAME=update_account_key
make atlas-new SCHEMA=sign NAME=add_index
```

**Docker environment**:

```bash
make atlas-new-docker SCHEMA=watch NAME=add_new_table
make atlas-new-docker SCHEMA=keygen NAME=update_account_key
make atlas-new-docker SCHEMA=sign NAME=add_index
```

## Migration File Naming

Atlas uses timestamp-based naming for migration files:

- Format: `YYYYMMDDHHMMSS_description.sql`
- Example: `20240101000000_initial_watch_schema.sql`

## Manual Atlas Commands

If you need to run Atlas commands directly (without Makefile):

### Local Environment

**Watch Schema**:

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

**Keygen Schema**:

```bash
atlas migrate apply \
  --dir file://tools/atlas/migrations/keygen \
  --url "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
```

**Sign Schema**:

```bash
atlas migrate apply \
  --dir file://tools/atlas/migrations/sign \
  --url "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
```

### Docker Environment

**Note**: Migrations are automatically applied when you run `docker compose --profile mysql up`. The following commands are for manual execution if needed.

Using the migration services:

**Watch Schema**:

```bash
# Apply migrations manually
docker compose run --rm wallet-mysql-migrate-watch migrate apply \
  --dir file://migrations/watch \
  --url "mysql://root:root@wallet-mysql:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"

# Check status
docker compose run --rm wallet-mysql-migrate-watch migrate status \
  --dir file://migrations/watch \
  --url "mysql://root:root@wallet-mysql:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
```

**Keygen Schema**:

```bash
# Apply migrations manually
docker compose run --rm wallet-mysql-migrate-keygen migrate apply \
  --dir file://migrations/keygen \
  --url "mysql://root:root@wallet-mysql:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"

# Check status
docker compose run --rm wallet-mysql-migrate-keygen migrate status \
  --dir file://migrations/keygen \
  --url "mysql://root:root@wallet-mysql:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
```

**Sign Schema**:

```bash
# Apply migrations manually
docker compose run --rm wallet-mysql-migrate-sign migrate apply \
  --dir file://migrations/sign \
  --url "mysql://root:root@wallet-mysql:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"

# Check status
docker compose run --rm wallet-mysql-migrate-sign migrate status \
  --dir file://migrations/sign \
  --url "mysql://root:root@wallet-mysql:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
```

**Note**: In Docker environment, the working directory is `/app/atlas` (mounted from `./tools/atlas`), so paths are relative to that directory (e.g., `file://migrations/watch` instead of `file://tools/atlas/migrations/watch`).

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

- **Location**: `tools/atlas/schemas/{db_dialect}/*.hcl`
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
   - **Docker environment**: Run `docker compose --profile mysql up` (migrations are applied automatically)
   - **Local environment**: Run `make atlas-migrate`
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
3. Apply the migration:
   - **Docker environment**: Run `docker compose --profile mysql up` (migrations are applied automatically)
   - **Local environment**: Run `make atlas-migrate`
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

1. Verify MySQL is running: `docker compose ps wallet-mysql`
2. Check connection string in Makefile targets
3. Verify credentials are correct

## Related Documentation

- [Database Architecture Documentation](../../docs/database/architecture.md)
- [Atlas Official Documentation](https://atlasgo.io/)
- [Atlas MySQL Guide](https://atlasgo.io/guides/mysql)
