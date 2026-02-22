# Design Document: Revise DB Tooling (Atlas Migration + SQLC Code Generation Flow)

## Overview

Align the Atlas migration and SQLC code generation tooling
with the documented multi-dialect flow described in
`docs/database/atlas-migration-flow.md` and
`docs/database/sqlc-code-generation-flow.md`.
The documentation describes support for MySQL, PostgreSQL,
and SQLite across 3 databases (watch, keygen, sign),
but the tooling code only fully supports MySQL.

## Background

### Current State

The project recently restructured migrations to
`migrations/{dialect}/{db}/` and added PostgreSQL Atlas
environments to `atlas.hcl`.
However, the surrounding tooling was not updated:

**Bugs:**

1. `db_atlas.mk` uses `--env local_$$schema`
   (expands to `local_watch`) but `atlas.hcl` env names
   are `local_mysql_watch` - all atlas-lint/validate/status
   targets are broken
2. Docker Compose MySQL migration services reference
   `file://migrations/watch` but actual path is
   `migrations/mysql/watch`

**Missing PostgreSQL tooling:**

1. No PostgreSQL Docker migration services in `compose.yaml`
2. No `sqlc_postgresql.yml` configuration
3. No PostgreSQL schema extraction script or Make targets
4. No PostgreSQL query files
   (PostgreSQL uses `$1, $2` placeholders, not `?`)
5. `atlas-dev-reset.sh` is hardcoded to MySQL
6. PostgreSQL migrations lack `atlas.sum` files
7. No `regenerate-all-from-atlas-postgresql` orchestration
   target

### Goals

1. Fix broken Makefile targets and Docker Compose paths
2. Add complete PostgreSQL SQLC pipeline
   (config, schema extraction, query files, code generation)
3. Make Atlas dev scripts dialect-aware (MySQL/PostgreSQL)
4. Ensure all tooling targets support `DB_DIALECT` parameter

## Design

### Task 1: Fix Makefile env naming in `db_atlas.mk`

**File**: `make/db_atlas.mk`

Add `DB_DIALECT` variable (default: `mysql`) and update
all env references:

```makefile
# Before (broken)
ATLAS_ENV_WATCH := local_mysql_watch
# ...
(cd tools/atlas && atlas schema lint --env local_$$schema)

# After (dialect-aware)
DB_DIALECT ?= mysql
# ...
(cd tools/atlas && atlas schema lint \
  --env local_$(DB_DIALECT)_$$schema)
```

Affected targets:

- `atlas-lint`, `atlas-validate`, `atlas-migrate-status`
- `atlas-migrate-apply-all`, `atlas-migrate-diff`,
  `atlas-migrate-hash-all`
- `atlas-schema-apply-all`, `atlas-schema-apply`,
  `atlas-dev-clean`

Add convenience aliases:
`atlas-lint-mysql`, `atlas-lint-postgresql`

### Task 2: Fix Docker Compose MySQL migration paths

**File**: `compose.yaml`

```yaml
# Before (broken)
- file://migrations/watch

# After (correct)
- file://migrations/mysql/watch
```

Fix all 4 MySQL migration services
(watch, keygen, sign, sign2).

### Task 3: Add PostgreSQL Docker migration services

**File**: `compose.yaml`

Add YAML anchor and 4 services:

```yaml
x-postgresql-migration-base: &postgresql-migration-base
  image: arigaio/atlas:1.1.0
  profiles: ["postgres"]
  volumes:
    - "./tools/atlas:/app/atlas"
  working_dir: /app/atlas
  depends_on:
    wallet-postgres:
      condition: service_healthy
  networks:
    - db
  environment:
    - ATLAS_NO_UPDATE_NOTIFIER=true
  restart: "no"

wallet-postgres-migrate-watch:
  <<: *postgresql-migration-base
  command:
    - migrate
    - apply
    - --dir
    - "file://migrations/postgresql/watch"
    - --url
    - "postgres://postgres:postgres@wallet-postgres:5432/watch?sslmode=disable"
# ... keygen, sign, sign2
```

### Task 4: Generate PostgreSQL atlas.sum files

Run `atlas migrate hash` for each PostgreSQL migration
directory:

```bash
cd tools/atlas
atlas migrate hash \
  --config file://atlas.hcl --env local_postgresql_watch
atlas migrate hash \
  --config file://atlas.hcl --env local_postgresql_keygen
atlas migrate hash \
  --config file://atlas.hcl --env local_postgresql_sign
```

### Task 5: Make `atlas-dev-reset.sh` dialect-aware

**File**: `scripts/db/atlas-dev-reset.sh`

- Accept `DB_DIALECT` env var (default: `mysql`)
- Derive `MIGRATIONS_DIR` as
  `${ATLAS_DIR}/migrations/${DB_DIALECT}`
- Derive env names as `local_${DB_DIALECT}_watch`, etc.
- Update success message to suggest
  `make reset-docker-${DB_DIALECT}`

### Task 6: Update `db.mk` orchestration targets

**File**: `make/db.mk`

Add `regenerate-all-from-atlas-postgresql`:

1. `atlas-dev-reset` with `DB_DIALECT=postgresql`
2. `docker compose down -v`
   `&& docker compose --profile postgres up -d`
3. `docker compose wait wallet-postgres-migrate-watch ...`
4. `extract-sqlc-schema-postgresql-all`
5. `sqlc-postgresql`

### Task 7: Add PostgreSQL SQLC pipeline

#### 7a: Create `sqlc_postgresql.yml`

**File**: `tools/sqlc/sqlc_postgresql.yml` (new)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./queries_postgresql/*.sql"
    schema: "./schemas_postgresql/*.sql"
    gen:
      go:
        package: "sqlcgen"
        out: "../../internal/infrastructure/database/postgresql/sqlcgen"
```

#### 7b: Create PostgreSQL schema extraction script

**File**: `scripts/db/extract-sqlc-schema-postgresql.sh`
(new)

- Use `pg_dump --schema-only --no-owner --no-privileges`
  via `docker exec wallet-postgres`
- Extract CREATE TABLE statements
- Exclude `atlas_schema_revisions` (all schemas),
  `seed` and `musig2_nonces` (sign only)
- Remove PostgreSQL noise
  (SET statements, comments, extensions)
- Output to `tools/sqlc/schemas_postgresql/`

#### 7c: Add PostgreSQL targets to `db_sqlc.mk`

**File**: `make/db_sqlc.mk`

New targets:

| Target | Description |
| ------ | ----------- |
| `dump-schema-postgresql-{watch,keygen,sign}` | `pg_dump` via Docker |
| `dump-schema-postgresql-all` | All three dumps |
| `clean-sqlc-schemas-postgresql` | Clean `schemas_postgresql/` |
| `extract-sqlc-schema-postgresql-{watch,keygen,sign}` | Extract + format |
| `extract-sqlc-schema-postgresql-all` | All three extractions |
| `regenerate-sqlc-postgresql-from-current-db` | Extract + generate |

#### 7d: Add PostgreSQL to `codegen.mk`

**File**: `make/codegen.mk`

```makefile
.PHONY: sqlc-postgresql
sqlc-postgresql:
    cd tools/sqlc && sqlc generate -f sqlc_postgresql.yml

.PHONY: sqlc-all
sqlc-all: sqlc sqlc-sqlite sqlc-postgresql
```

#### 7e: Add PostgreSQL to sqlc validation targets

**File**: `make/db_sqlc.mk`

Add PostgreSQL compile/vet/validate alongside MySQL
and SQLite.

#### 7f: Create PostgreSQL query files

**Directory**: `tools/sqlc/queries_postgresql/`
(new, 20 files)

MySQL queries use `?` placeholders;
PostgreSQL requires `$1, $2, ...`.
Queries cannot be shared.
This matches the existing pattern
(`schemas/` vs `schemas_sqlite/` vs `schemas_postgresql/`).

Conversion rules:

- Replace `?` with `$1, $2, $3, ...`
  (positional, incrementing per query)
- For `:execresult` inserts, add `RETURNING id` and change
  to `:one`
  (PostgreSQL doesn't support `LastInsertId()`)
- Keep same `-- name:` annotations and query structure

20 files to create
(mechanical conversion from MySQL originals):

- `address.sql`, `auth_account_key.sql`,
  `auth_fullpubkey.sql`
- `btc_account_key.sql`, `btc_tx.sql`,
  `btc_tx_input.sql`, `btc_tx_output.sql`
- `eth_account_key.sql`, `eth_detail_tx.sql`
- `musig2_nonces.sql`, `payment_request.sql`,
  `seed.sql`, `tx.sql`
- `xrp_account_key.sql`, `xrp_detail_tx.sql`,
  `xrp_multisig_signature.sql`
- `xrp_pending_multisig.sql`, `xrp_regular_key.sql`,
  `xrp_signer_entry.sql`, `xrp_signer_list.sql`

### Task 8: Update docs flow diagram

**File**: `docs/database/sqlc-code-generation-flow.md`

The doc references `strip_schema.sh` which doesn't exist.
Update to reference `extract-sqlc-schema.sh` and
`extract-sqlc-schema-postgresql.sh`.

## Task 9: Refactor `tools/sqlc/` directory structure

### Motivation

The current `tools/sqlc/` directory uses inconsistent
flat naming with `_dialect` suffixes
(`queries_postgresql/`, `schemas_sqlite/`,
`schemas_postgresql/`). This doesn't match the
`tools/atlas/migrations/{dialect}/{db}/` pattern used
elsewhere in the project.

### Current Structure

```
tools/sqlc/
├── queries/                    # MySQL queries (shared with SQLite via ?)
├── queries_postgresql/         # PostgreSQL queries ($1, $2, ...)
├── schemas/                    # MySQL schemas (auto-generated)
├── schemas_sqlite/             # SQLite schemas (manually maintained)
├── schemas_postgresql/         # PostgreSQL schemas (auto-generated)
├── sqlc.yml                    # MySQL config
├── sqlc_sqlite.yml             # SQLite config
└── sqlc_postgresql.yml         # PostgreSQL config
```

### Target Structure

```
tools/sqlc/
├── queries/
│   ├── mysql/                  # ← from queries/ (MySQL ? placeholders)
│   └── postgresql/             # ← from queries_postgresql/ ($1,$2 placeholders)
├── schemas/
│   ├── mysql/                  # ← from schemas/ (auto-generated from mysqldump)
│   ├── postgresql/             # ← from schemas_postgresql/ (auto-generated from pg_dump)
│   └── sqlite/                 # ← from schemas_sqlite/ (manually maintained)
├── sqlc.yml                    # Update paths
├── sqlc_sqlite.yml             # Update paths
└── sqlc_postgresql.yml         # Update paths
```

**Note on `queries/sqlite/`**: Not needed.
SQLite uses `?` placeholders (same as MySQL),
so SQLite reuses `queries/mysql/*.sql` via its sqlc config.

**Note on `schemas/sqlite/`**: IS needed.
SQLite schemas differ from MySQL
(TEXT vs VARCHAR, CHECK constraints vs ENUM, etc.)
and are manually maintained.

### Files to Update

#### 9a: Move directories

```bash
# Queries
mv tools/sqlc/queries tools/sqlc/queries_tmp
mkdir -p tools/sqlc/queries/mysql
mv tools/sqlc/queries_tmp/* tools/sqlc/queries/mysql/
rmdir tools/sqlc/queries_tmp
mv tools/sqlc/queries_postgresql tools/sqlc/queries/postgresql

# Schemas
mv tools/sqlc/schemas tools/sqlc/schemas_tmp
mkdir -p tools/sqlc/schemas/mysql
mv tools/sqlc/schemas_tmp/* tools/sqlc/schemas/mysql/
rmdir tools/sqlc/schemas_tmp
mv tools/sqlc/schemas_sqlite tools/sqlc/schemas/sqlite
mv tools/sqlc/schemas_postgresql tools/sqlc/schemas/postgresql
```

#### 9b: Update sqlc config files

**`tools/sqlc/sqlc_mysql.yml`** (MySQL):

```yaml
queries: "./queries/mysql/*.sql"
schema: "./schemas/mysql/*.sql"
```

**`tools/sqlc/sqlc_sqlite.yml`**:

```yaml
queries: "./queries/mysql/*.sql"    # shared with MySQL
schema: "./schemas/sqlite/*.sql"
```

**`tools/sqlc/sqlc_postgresql.yml`**:

```yaml
queries: "./queries/postgresql/*.sql"
schema: "./schemas/postgresql/*.sql"
```

#### 9c: Update Makefile targets

**`make/db_sqlc.mk`** - Update all schema output paths:

| Before | After |
| ------ | ----- |
| `tools/sqlc/schemas/*.sql` | `tools/sqlc/schemas/mysql/*.sql` |
| `tools/sqlc/schemas_postgresql/*.sql` | `tools/sqlc/schemas/postgresql/*.sql` |
| `clean-sqlc-schemas`: `rm -f tools/sqlc/schemas/*.sql` | `rm -f tools/sqlc/schemas/mysql/*.sql` |
| `clean-sqlc-schemas-postgresql`: `rm -f tools/sqlc/schemas_postgresql/*.sql` | `rm -f tools/sqlc/schemas/postgresql/*.sql` |

**`make/codegen.mk`** - Update comments referencing paths

**`make/db.mk`** - Update sqlfluff paths:

```makefile
# Before
@sqlfluff format tools/sqlc/queries/*.sql
# After
@sqlfluff format tools/sqlc/queries/mysql/*.sql
```

#### 9d: Update extraction scripts

**`scripts/db/extract-sqlc-schema.sh`** -
Update example path in usage comment

**`scripts/db/extract-sqlc-schema-postgresql.sh`** -
Update example path in usage comment

#### 9e: Update E2E operation scripts

**`scripts/operation/btc/btc_common.sh`** and
**`scripts/operation/bch/bch_common.sh`**:

```bash
# Before
SQLITE_WATCH_SCHEMA="${SQLITE_WATCH_SCHEMA:-tools/sqlc/schemas_sqlite/01_watch.sql}"
# After
SQLITE_WATCH_SCHEMA="${SQLITE_WATCH_SCHEMA:-tools/sqlc/schemas/sqlite/01_watch.sql}"
```

#### 9f: Update documentation

Files referencing old paths (docs only, no code impact):

- `docs/database/schema-changes.md`
- `docs/database/quick-reference.md`
- `docs/database/db-management.md`
- `docs/database/architecture.md`
- `docs/database/sqlc-code-generation-flow.md`
- `docs/guidelines/code-generation.md`
- `.claude/rules/sql.md`
- `.claude/skills/db-migration/SKILL.md`
- `.kiro/specs/postgresql-integration/` (multiple files)

### Verification

1. `make sqlc-all` - generates for all 3 dialects
   with new paths
2. `make sqlc-validate` - validates all configs
3. `make extract-sqlc-schema-all` - MySQL extraction
   outputs to `schemas/mysql/`
4. `make extract-sqlc-schema-postgresql-all` - PostgreSQL
   extraction outputs to `schemas/postgresql/`
5. E2E scripts reference correct SQLite schema paths

## Scope Exclusions

- **Go repository code**
  (PostgreSQL sqlcgen usage, DI wiring) - separate task
- **SQLite Atlas migrations** - not used;
  SQLite schemas are maintained manually for sqlc only

## Files Modified Summary

| File | Action |
| ---- | ------ |
| `make/db_atlas.mk` | Fix env naming, add DB_DIALECT |
| `make/db_sqlc.mk` | Add PostgreSQL dump/extract/validate targets |
| `make/db.mk` | Add `regenerate-all-from-atlas-postgresql` |
| `make/codegen.mk` | Add `sqlc-postgresql`, update `sqlc-all` |
| `compose.yaml` | Fix MySQL paths, add PostgreSQL migration services |
| `scripts/db/atlas-dev-reset.sh` | Make dialect-aware |
| `scripts/db/extract-sqlc-schema-postgresql.sh` | New: PostgreSQL schema extraction |
| `tools/sqlc/sqlc_postgresql.yml` | New: PostgreSQL SQLC config |
| `tools/sqlc/queries_postgresql/*.sql` | New: 20 PostgreSQL query files |
| `tools/atlas/migrations/postgresql/*/atlas.sum` | New: generated hash files |
| `docs/database/sqlc-code-generation-flow.md` | Fix script reference |

## Verification

1. `make atlas-lint` - passes without env name errors
   (MySQL default)
2. `make atlas-lint DB_DIALECT=postgresql` - passes for
   PostgreSQL
3. `make reset-docker-mysql` - migration services find
   correct paths
4. `make reset-docker-postgres` - PostgreSQL migrations
   apply successfully
5. `make regenerate-all-from-atlas` - full MySQL pipeline
   end-to-end
6. `make regenerate-all-from-atlas-postgresql` - full
   PostgreSQL pipeline end-to-end
7. `make sqlc-validate` - validates MySQL, SQLite,
   and PostgreSQL configs
