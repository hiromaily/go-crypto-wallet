### Schema Migrations with Atlas

The project uses [Atlas](https://atlasgo.io/) for managing database schema migrations. Atlas provides version-controlled migrations, migration history tracking, and rollback capabilities.

#### Installation

```bash
# Homebrew (macOS)
brew install arigaio/tap/atlas

# Or via Go
go install ariga.io/atlas/cmd/atlas@latest

# Verify
atlas version
```

#### Migration Structure

Atlas uses HCL schema definitions as source of truth and generates SQL migration files:

```
tools/atlas/
├── atlas.hcl                          # Atlas configuration (environments for each dialect/schema)
├── schemas/                           # HCL schema definitions (source of truth)
│   ├── mysql/
│   │   ├── watch.hcl
│   │   ├── keygen.hcl
│   │   └── sign.hcl
│   └── postgres/
│       ├── watch.hcl
│       ├── keygen.hcl
│       └── sign.hcl
└── migrations/                        # Generated SQL migration files
    ├── mysql/
    │   ├── watch/                     # Watch schema migrations + atlas.sum
    │   ├── keygen/                    # Keygen schema migrations + atlas.sum
    │   └── sign/                      # Sign schema migrations + atlas.sum
    └── postgres/
        ├── watch/                     # Watch schema migrations + atlas.sum
        ├── keygen/                    # Keygen schema migrations + atlas.sum
        └── sign/                      # Sign schema migrations + atlas.sum
```

#### Atlas Configuration

The `atlas.hcl` defines environments for each dialect/schema combination:

```hcl
# Example: PostgreSQL watch environment
env "local_postgres_watch" {
  url     = "postgres://postgres:postgres@127.0.0.1:5432/watch?sslmode=disable"
  src     = "file://schemas/postgres/watch.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgres/watch"
  }
  dev = "docker://postgres/18/watch"
}
```

**Available environments**: `local_{mysql|postgres}_{watch|keygen|sign}`, `admin_{mysql|postgres}_{watch|keygen|sign}`

**Features**:

- Destructive change protection (drop_schema, drop_table, drop_column)
- Lint rules for naming conventions: `^[a-z][a-z0-9_]*$`
- Dev databases for migration generation

#### Common Operations

##### Format and Lint

```bash
make atlas-fmt          # Format all HCL schema files
make atlas-fmt-check    # Check formatting (CI mode)
make atlas-lint         # Lint all schemas (both dialects)
make atlas-validate     # Validate Atlas configuration
```

##### Apply Schema Changes

```bash
# Apply HCL schema directly to database (all schemas)
make atlas-schema-apply-all                     # PostgreSQL (default)
make atlas-schema-apply-all DB_DIALECT=mysql     # MySQL

# Apply specific schema
make atlas-schema-apply SCHEMA=watch
```

##### Migration Management

```bash
# Check migration status
make atlas-migrate-status                        # PostgreSQL (default)
make atlas-migrate-status DB_DIALECT=mysql        # MySQL

# Apply all pending migrations
make atlas-migrate-apply-all

# Generate new migration from HCL diff
make atlas-migrate-diff SCHEMA=watch NAME=add_new_table

# Hash migration directory
make atlas-migrate-hash-all
```

##### Regenerate Migrations from Scratch

```bash
# Regenerate migrations from HCL schemas
make atlas-dev-reset                             # PostgreSQL (default)
make atlas-dev-reset DB_DIALECT=mysql             # MySQL
```

##### Clean and Recreate Databases

```bash
# Drop all databases and recreate from HCL (WARNING: destructive)
make atlas-dev-clean                             # PostgreSQL (default)
make atlas-dev-clean DB_DIALECT=mysql             # MySQL
```

#### Full Regeneration Workflow

When HCL schemas change, regenerate everything:

```bash
make regenerate-all-from-atlas                   # PostgreSQL (default)
make regenerate-all-from-atlas DB_DIALECT=mysql   # MySQL
```

This runs 5 steps:

1. Regenerate Atlas migrations from HCL
2. Restart Docker DB
3. Wait for DB and migrations to complete
4. Extract SQLC schemas from DB
5. Generate SQLC Go code

#### Best Practices

1. **Always format and lint** before committing: `make atlas-fmt && make atlas-lint`
2. **Check status** before applying: `make atlas-migrate-status`
3. **Never modify existing migration files** - create new migrations instead
4. **Keep migrations small and focused** - one logical change per migration
5. **Test on development database first**
6. **Update both MySQL and PostgreSQL HCL schemas** when making changes

For more detailed information, see [Atlas README](https://github.com/hiromaily/go-crypto-wallet/blob/main/tools/atlas/README.md).

---
