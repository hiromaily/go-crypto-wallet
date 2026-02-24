# Database Architecture

This document describes the database architecture and operations for the go-crypto-wallet project.

---

## Table of Contents

- [Overview](#overview)
- [Supported Databases](#supported-databases)
- [Architecture](#architecture)
- [Schema Design](#schema-design)
- [Setup and Configuration](#setup-and-configuration)
- [Common Operations](#common-operations)
- [Database Management](#database-management)
- [Schema Migrations with Atlas](#schema-migrations-with-atlas)
- [SQLC Code Generation](#sqlc-code-generation)
- [SQLite for E2E Testing](#sqlite-for-e2e-testing)
- [Troubleshooting](#troubleshooting)
- [Migration Guide](#migration-guide)

---

## Overview

The project supports **three database backends**:

| Database | Version | Use Case | Features |
|----------|---------|----------|----------|
| **PostgreSQL** | 18.2 | Production (default) | Docker container, schema separation, `identity` columns |
| **MySQL** | 8.4 | Production (alternative) | Docker container, schema separation, `auto_increment` |
| **SQLite** | - | E2E testing, CI/CD | Local file, fast startup, no Docker required |

All backends use **four separate databases/schemas** to manage wallet data:

- **`watch`**: Online wallet data (addresses, transactions, payment requests)
- **`keygen`**: Key generation data (seeds, account keys, full public keys)
- **`sign`**: Signing wallet data (auth account keys, seeds)
- **`sign2`**: Second signing wallet (same schema as `sign`, separate database)

This approach provides:

- Reduced resource usage (single DB instance per dialect)
- Simplified deployment and maintenance
- Data isolation through database/schema separation
- Easier backup and restore operations
- Single point of configuration

---

## Supported Databases

### PostgreSQL (Production Default)

The project uses a **single PostgreSQL 18.2 container** with **four separate databases** (watch, keygen, sign, sign2). PostgreSQL uses named `enum` types, `identity` columns, `timestamptz`, and `bytea` for binary data.

### MySQL (Production Alternative)

The project uses a **single MySQL 8.4 container** with **four separate schemas** (watch, keygen, sign, sign2). MySQL uses inline `enum()`, `auto_increment`, `datetime`, and `blob` for binary data.

**Note**: MySQL 8.4 uses `caching_sha2_password` authentication by default, which requires SSL for local connections. Use `?tls=true` in connection strings.

### SQLite (E2E Testing)

SQLite provides a lightweight alternative for E2E testing without Docker. Uses `CHECK` constraints for enums, `INTEGER PRIMARY KEY` for auto-increment, and `TEXT` for timestamps.

---

## Architecture

### Container Setup

```yaml
services:
  # PostgreSQL (default)
  wallet-postgres:
    image: postgres:18.2
    profiles: ["postgres"]
    ports:
      - "${POSTGRESQL_PORT:-5432}:5432"
    volumes:
      - wallet-postgres:/var/lib/postgres
      - "./docker/postgres/init.d:/docker-entrypoint-initdb.d"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres

  # MySQL (alternative)
  wallet-mysql:
    image: mysql:8.4
    profiles: ["mysql"]
    ports:
      - "${MYSQL_PORT:-3306}:3306"
    volumes:
      - wallet-mysql:/var/lib/mysql
      - "./docker/mysql/conf.d:/etc/mysql/conf.d"
      - "./docker/mysql/init.d:/docker-entrypoint-initdb.d"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_USER: hiromaily
      MYSQL_PASSWORD: hiromaily
```

Both services include health checks and are activated via Docker Compose profiles (`--profile postgres` or `--profile mysql`).

### Migration Services

Atlas migration services run automatically when databases start. Each schema has a dedicated migration service:

```yaml
# Example: PostgreSQL watch migration
wallet-postgres-migrate-watch:
  image: arigaio/atlas:1.1.0
  profiles: ["postgres"]
  command:
    - migrate
    - apply
    - --dir
    - "file://migrations/postgres/watch"
    - --url
    - "postgres://postgres:postgres@wallet-postgres:5432/watch?sslmode=disable"
  depends_on:
    wallet-postgres:
      condition: service_healthy
  restart: "no"
```

Migration services exist for: `watch`, `keygen`, `sign`, `sign2` (both MySQL and PostgreSQL).

### Directory Structure

```
docker/
├── mysql/
│   ├── archive/                       # Archived SQL schema files (reference only)
│   ├── conf.d/
│   │   └── custom.cnf                 # Server-level configuration (utf8mb4)
│   ├── init.d/
│   │   └── 01_init_all_schemas_2.sql  # Creates watch, keygen, sign, sign2 databases
│   └── insert/
│       └── ganache.example.sql        # Test data for Ganache
├── postgres/
│   └── init.d/
│       └── 01_create_databases.sh     # Creates watch, keygen, sign, sign2 databases
```

**Note:** Schema definitions are managed by Atlas (HCL files in `tools/atlas/schemas/`). The init scripts only create empty databases.

### Initialization Process

When the container starts for the first time:

1. **Database Creation**: Init scripts create four empty databases

   **PostgreSQL** (`01_create_databases.sh`):

   ```sql
   CREATE DATABASE watch;
   CREATE DATABASE keygen;
   CREATE DATABASE sign;
   CREATE DATABASE sign2;
   ```

   **MySQL** (`01_init_all_schemas_2.sql`):

   ```sql
   CREATE DATABASE `watch` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `keygen` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `sign` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `sign2` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. **User/Permission Setup**: Grants privileges to application users

3. **Schema Migration**: Atlas migration services automatically apply migrations after the database is healthy

---

## Schema Design

### Watch Schema (`watch`)

**Purpose**: Manages online wallet operations including address tracking, transaction monitoring, and payment requests.

**Tables**:

| Table | Description |
|-------|-------------|
| `address` | Wallet addresses for all coin/account types (btc, bch, eth, xrp, hyt) |
| `btc_tx` | Bitcoin/BCH transaction records |
| `btc_tx_input` | Bitcoin transaction inputs (UTXOs) |
| `btc_tx_output` | Bitcoin transaction outputs |
| `tx` | Generic transaction records (ETH, XRP, HYT) |
| `eth_detail_tx` | Ethereum transaction details |
| `xrp_detail_tx` | XRP transaction details |
| `xrp_pending_multisig` | Pending XRP multi-signature transactions awaiting signatures |
| `xrp_multisig_signature` | Collected signatures for XRP multi-signature transactions |
| `payment_request` | Payment request queue |

**Access Pattern**: High read/write - monitors blockchain, creates transactions

### Keygen Schema (`keygen`)

**Purpose**: Stores key generation data for offline key generation wallet.

**Tables**:

| Table | Description |
|-------|-------------|
| `seed` | Encrypted seed phrases (all coins) |
| `btc_account_key` | BTC/BCH keys with multiple address formats (P2PKH, P2SH-SegWit, Bech32, Taproot) |
| `eth_account_key` | Ethereum keys (address, public key, private key) |
| `xrp_account_key` | XRP-specific account keys |
| `auth_fullpubkey` | Full public keys for multisig authentication (with BIP32 extended key support) |
| `xrp_signer_list` | XRP SignerList configuration for multi-signature accounts |
| `xrp_signer_entry` | Individual signer entries within an XRP signer list |
| `xrp_regular_key` | XRP regular key assignments for enhanced security |
| `musig2_nonces` | MuSig2 nonce commitments (shared with sign schema) |

**Access Pattern**: Write-heavy during key generation, read-only during export

**Security**: This schema contains sensitive key material - should be in offline/cold storage in production

### Sign Schema (`sign` / `sign2`)

**Purpose**: Stores signing wallet data for offline transaction signing. `sign2` uses the same schema for a second signing wallet instance.

**Tables**:

| Table | Description |
|-------|-------------|
| `seed` | Encrypted seed phrases for signing wallet (BTC/BCH only) |
| `auth_account_key` | Authentication account keys for signing (with multiple address formats) |
| `musig2_nonces` | MuSig2 nonce commitments (shared with keygen schema) |

**Access Pattern**: Read-only during signing operations

**Security**: This schema contains sensitive signing keys - should be in offline/cold storage in production

### Key Type Support

The `btc_account_key` and `auth_account_key` tables support multiple key types:

| Key Type | BIP Standard | Address Format |
|----------|-------------|----------------|
| `bip44` | BIP-44 | P2PKH (legacy) |
| `bip49` | BIP-49 | P2SH-SegWit |
| `bip84` | BIP-84 | Bech32 (native SegWit) |
| `bip86` | BIP-86 | Taproot |
| `musig2` | MuSig2 | Multi-signature |

### Type Differences Across Dialects

| Feature | PostgreSQL | MySQL | SQLite |
|---------|-----------|-------|--------|
| Auto-increment | `identity` | `auto_increment` | `INTEGER PRIMARY KEY` |
| Numeric | `numeric(26,10)` | `decimal(26,10)` | `TEXT` |
| Enums | Named types | Inline `enum()` | `CHECK` constraints |
| Timestamps | `timestamptz` | `datetime` | `TEXT` |
| Boolean | `boolean` | `bool` (tinyint) | `INTEGER` |
| Binary | `bytea` | `blob` | `BLOB` |

---

## Setup and Configuration

### Initial Setup

1. **Start the database** (choose one dialect):

   ```bash
   # PostgreSQL (default)
   docker compose --profile postgres up -d

   # MySQL
   docker compose --profile mysql up -d
   ```

2. **Migrations run automatically** via Atlas migration services. Wait for completion:

   ```bash
   # PostgreSQL
   docker compose wait wallet-postgres-migrate-watch wallet-postgres-migrate-keygen wallet-postgres-migrate-sign

   # MySQL
   docker compose wait wallet-mysql-migrate-watch wallet-mysql-migrate-keygen wallet-mysql-migrate-sign
   ```

3. **Verify databases and tables**:

   ```bash
   # PostgreSQL
   docker exec wallet-postgres psql -U postgres -c "\l"
   docker exec wallet-postgres psql -U postgres -d watch -c "\dt"

   # MySQL
   docker exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"
   docker exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"
   ```

### Application Configuration

Each wallet type connects to the same database host but specifies different database names. Configuration files support all three backends:

**Example** (`config/wallet/btc/watch.yaml`):

```yaml
database:
  type: "sqlite"  # mysql, sqlite, postgres
  mysql:
    host: "127.0.0.1:3306"
    dbname: "watch"
    user: "hiromaily"
    pass: "hiromaily"
    debug: true
  postgres:
    host: "127.0.0.1"
    port: 5432
    dbname: "watch"
    user: "postgres"
    pass: "postgres"
    sslmode: "disable"
    debug: true
  sqlite:
    path: "./data/sqlite/btc/watch.db"
    max_open_conns: 2  # Minimum 2 to prevent deadlock
    debug: true
```

**Connection details**:

| Dialect | Host | Port | User | Password |
|---------|------|------|------|----------|
| PostgreSQL | `127.0.0.1` | `5432` | `postgres` | `postgres` |
| MySQL | `127.0.0.1` | `3306` | `hiromaily` | `hiromaily` |
| SQLite | N/A | N/A | N/A | N/A |

---

## Common Operations

### Database Access

**PostgreSQL**:

```bash
# Access watch database
docker exec -it wallet-postgres psql -U postgres -d watch

# Access keygen database
docker exec -it wallet-postgres psql -U postgres -d keygen

# Access sign database
docker exec -it wallet-postgres psql -U postgres -d sign
```

**MySQL**:

```bash
# Access watch schema
docker exec -it wallet-mysql mysql -uroot -proot watch

# Access keygen schema
docker exec -it wallet-mysql mysql -uroot -proot keygen

# Access sign schema
docker exec -it wallet-mysql mysql -uroot -proot sign
```

### Schema Export (Backup)

Export schema structure without data:

```bash
# Export all schemas (uses DB_DIALECT, default: postgres)
make dump-schema-all

# Export for specific dialect
make dump-schema-all DB_DIALECT=mysql
make dump-schema-all DB_DIALECT=postgres

# Export individual schemas
make dump-schema-watch
make dump-schema-keygen
make dump-schema-sign
```

Output location: `data/dump/sql/dump_*.{dialect}.sql`

### Data Export (Full Backup)

**PostgreSQL**:

```bash
docker exec wallet-postgres pg_dump -U postgres watch > backups/watch_$(date +%Y%m%d).sql
docker exec wallet-postgres pg_dump -U postgres keygen > backups/keygen_$(date +%Y%m%d).sql
docker exec wallet-postgres pg_dump -U postgres sign > backups/sign_$(date +%Y%m%d).sql
```

**MySQL**:

```bash
docker exec wallet-mysql mysqldump -uroot -proot watch > backups/watch_$(date +%Y%m%d).sql
docker exec wallet-mysql mysqldump -uroot -proot keygen > backups/keygen_$(date +%Y%m%d).sql
docker exec wallet-mysql mysqldump -uroot -proot sign > backups/sign_$(date +%Y%m%d).sql
```

### Reset Database

Complete database reset (WARNING: deletes all data):

```bash
# Reset specific dialect
make reset-docker                    # PostgreSQL (default)
make reset-docker DB_DIALECT=mysql   # MySQL

# Remove all database containers
make remove-all-dbs
```

---

## Database Management

### View Schema Information

**PostgreSQL**:

```bash
# List all tables in watch database
docker exec wallet-postgres psql -U postgres -d watch -c "\dt"

# Describe a specific table
docker exec wallet-postgres psql -U postgres -d watch -c "\d address"
```

**MySQL**:

```bash
# List all tables in watch schema
docker exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"

# Describe a specific table
docker exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE address;"
```

### View Logs

```bash
# PostgreSQL logs
docker compose logs wallet-postgres
docker compose logs -f wallet-postgres

# MySQL logs
docker compose logs wallet-mysql
docker compose logs -f wallet-mysql
```

---

## Schema Migrations with Atlas

The project uses [Atlas](https://atlasgo.io/) for managing database schema migrations. Atlas provides version-controlled migrations, migration history tracking, and rollback capabilities.

### Installation

```bash
# Homebrew (macOS)
brew install arigaio/tap/atlas

# Or via Go
go install ariga.io/atlas/cmd/atlas@latest

# Verify
atlas version
```

### Migration Structure

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

### Atlas Configuration

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

### Common Operations

#### Format and Lint

```bash
make atlas-fmt          # Format all HCL schema files
make atlas-fmt-check    # Check formatting (CI mode)
make atlas-lint         # Lint all schemas (both dialects)
make atlas-validate     # Validate Atlas configuration
```

#### Apply Schema Changes

```bash
# Apply HCL schema directly to database (all schemas)
make atlas-schema-apply-all                     # PostgreSQL (default)
make atlas-schema-apply-all DB_DIALECT=mysql     # MySQL

# Apply specific schema
make atlas-schema-apply SCHEMA=watch
```

#### Migration Management

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

#### Regenerate Migrations from Scratch

```bash
# Regenerate migrations from HCL schemas
make atlas-dev-reset                             # PostgreSQL (default)
make atlas-dev-reset DB_DIALECT=mysql             # MySQL
```

#### Clean and Recreate Databases

```bash
# Drop all databases and recreate from HCL (WARNING: destructive)
make atlas-dev-clean                             # PostgreSQL (default)
make atlas-dev-clean DB_DIALECT=mysql             # MySQL
```

### Full Regeneration Workflow

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

### Best Practices

1. **Always format and lint** before committing: `make atlas-fmt && make atlas-lint`
2. **Check status** before applying: `make atlas-migrate-status`
3. **Never modify existing migration files** - create new migrations instead
4. **Keep migrations small and focused** - one logical change per migration
5. **Test on development database first**
6. **Update both MySQL and PostgreSQL HCL schemas** when making changes

For more detailed information, see [Atlas README](../../tools/atlas/README.md).

---

## SQLC Code Generation

The project uses [SQLC](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

### Configuration

SQLC has separate configuration files per dialect:

| Config File | Engine | Output |
|-------------|--------|--------|
| `tools/sqlc/sqlc_postgres.yml` | PostgreSQL | `internal/infrastructure/database/postgres/sqlcgen/` |
| `tools/sqlc/sqlc_mysql.yml` | MySQL | `internal/infrastructure/database/mysql/sqlcgen/` |
| `tools/sqlc/sqlc_sqlite.yml` | SQLite | `internal/infrastructure/database/sqlite/sqlcgen/` |

### Structure

```
tools/sqlc/
├── sqlc_postgres.yml             # PostgreSQL SQLC config
├── sqlc_mysql.yml                # MySQL SQLC config
├── sqlc_sqlite.yml               # SQLite SQLC config
├── schemas/                      # Schema files (extracted from running DB)
│   ├── postgres/*.sql
│   ├── mysql/*.sql
│   └── sqlite/*.sql
└── queries/                      # SQL query files (20 per dialect)
    ├── postgres/*.sql
    └── mysql/*.sql
```

**Query files** (20 per dialect): `address`, `auth_account_key`, `auth_fullpubkey`, `btc_account_key`, `btc_tx`, `btc_tx_input`, `btc_tx_output`, `eth_account_key`, `eth_detail_tx`, `musig2_nonces`, `payment_request`, `seed`, `tx`, `xrp_account_key`, `xrp_detail_tx`, `xrp_multisig_signature`, `xrp_pending_multisig`, `xrp_regular_key`, `xrp_signer_entry`, `xrp_signer_list`

### Common Operations

```bash
# Generate Go code
make sqlc                          # PostgreSQL (default)
make sqlc-mysql                    # MySQL
make sqlc-sqlite                   # SQLite
make sqlc-all                      # All dialects

# Validate SQL queries
make sqlc-compile                  # Check syntax (all dialects)
make sqlc-vet                      # Vet for potential issues
make sqlc-validate                 # Both compile and vet

# Extract schema from running database
make extract-sqlc-schema-all                     # PostgreSQL (default)
make extract-sqlc-schema-all DB_DIALECT=mysql     # MySQL

# Full regeneration (extract + generate)
make regenerate-sqlc-from-current-db              # PostgreSQL (default)
make regenerate-sqlc-from-current-db DB_DIALECT=mysql  # MySQL
```

### SQL Formatting

```bash
make sqlfluff-format    # Format all SQL query files
make sqlfluff-lint      # Lint SQL query files
make sqlfluff-fix       # Format and auto-fix SQL
```

---

## SQLite for E2E Testing

SQLite provides a lightweight alternative for E2E testing without requiring Docker.

### Benefits

- **Faster startup**: No Docker container needed
- **Parallel testing**: Each test can use separate database files
- **Lighter CI/CD**: Reduced infrastructure requirements
- **Simpler debugging**: Direct file access

### Configuration

All wallet config files support SQLite via the `database.type` field:

```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/sqlite/btc/watch.db"
    max_open_conns: 2  # Minimum 2 to prevent deadlock
    debug: true
```

### E2E Script Usage

```bash
# SQLite (default) - faster, no Docker
make btc-e2e-reset P=1

# MySQL - traditional Docker-based testing
make btc-e2e-reset P=1 DB=mysql
```

### SQLite Schema Files

SQLite-compatible schemas are located in:

```
tools/sqlc/schemas/sqlite/
├── 01_watch.sql   # Watch wallet schema
├── 02_keygen.sql  # Keygen wallet schema
└── 03_sign.sql    # Sign wallet schema
```

These schemas are converted from MySQL/PostgreSQL with the following changes:

| MySQL/PostgreSQL | SQLite |
|-----------------|--------|
| `ENUM('a','b')` / named enum | `TEXT CHECK(col IN ('a','b'))` |
| `AUTO_INCREMENT` / `identity` | `AUTOINCREMENT` |
| `TINYINT(1)` / `boolean` | `INTEGER` |
| `DATETIME` / `timestamptz` | `TEXT DEFAULT CURRENT_TIMESTAMP` |
| `BLOB` / `bytea` | `BLOB` |

### SQLite Data Location

```
data/sqlite/
└── btc/
    └── watch.db  # E2E test database
```

**Note**: Database files are gitignored (`data/sqlite/**/*.db`)

---

## Troubleshooting

### Container Won't Start

**Check logs**:

```bash
docker compose logs wallet-postgres  # or wallet-mysql
```

**Common issues**:

1. Port already in use:

   ```bash
   lsof -i :5432  # or :3306

   # Use different port
   POSTGRESQL_PORT=5433 docker compose --profile postgres up -d
   MYSQL_PORT=3307 docker compose --profile mysql up -d
   ```

2. Volume permission issues:

   ```bash
   docker compose --profile postgres down -v
   docker compose --profile postgres up -d
   ```

### Cannot Connect to Database

**Verify container is running**:

```bash
docker compose ps
```

**PostgreSQL**:

```bash
docker exec wallet-postgres pg_isready -U postgres
psql -h 127.0.0.1 -U postgres -d watch -c "SELECT 1;"
```

**MySQL**:

```bash
docker exec wallet-mysql mysqladmin ping -uroot -proot
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 watch -e "SELECT 1;"
```

### Schema Not Found

```bash
# PostgreSQL
docker exec wallet-postgres psql -U postgres -c "\l"

# MySQL
docker exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"
```

### Migration Fails

1. Check error message in migration service logs
2. Verify database connection
3. Ensure database exists
4. Review migration file syntax
5. Check migration status: `make atlas-migrate-status`

### Character Set Issues (MySQL)

```bash
docker exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'character_set_server';"
docker exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'collation_server';"
```

Expected: `utf8mb4` and `utf8mb4_unicode_ci`

---

## Migration Guide

### From Old Three-Container Setup

If migrating from the previous three-container setup (`watch-db`, `keygen-db`, `sign-db`):

#### 1. Backup Existing Data

```bash
docker compose exec watch-db mysqldump -uroot -proot watch > migration/watch_backup.sql
docker compose exec keygen-db mysqldump -uroot -proot keygen > migration/keygen_backup.sql
docker compose exec sign-db mysqldump -uroot -proot sign > migration/sign_backup.sql
```

#### 2. Update Configuration

```toml
# Change from:
host = "127.0.0.1:3307"  # or 3308, 3309

# To (MySQL):
host = "127.0.0.1:3306"
```

Or switch to PostgreSQL:

```yaml
database:
  type: "postgres"
  postgres:
    host: "127.0.0.1"
    port: 5432
    dbname: "watch"
    user: "postgres"
    pass: "postgres"
    sslmode: "disable"
```

#### 3. Stop Old Containers

```bash
docker compose stop watch-db keygen-db sign-db
docker compose rm -f watch-db keygen-db sign-db
```

#### 4. Start New Container

```bash
# PostgreSQL (recommended)
docker compose --profile postgres up -d

# Or MySQL
docker compose --profile mysql up -d
```

#### 5. Restore Data (Optional)

```bash
# Wait for container and migrations
docker compose wait wallet-postgres-migrate-watch wallet-postgres-migrate-keygen wallet-postgres-migrate-sign
```

#### 6. Cleanup Old Volumes (Optional)

```bash
docker volume rm go-crypto-wallet_watch-db
docker volume rm go-crypto-wallet_keygen-db
docker volume rm go-crypto-wallet_sign-db
```

---

## References

- [PostgreSQL 18 Documentation](https://www.postgresql.org/docs/18/)
- [MySQL 8.4 Documentation](https://dev.mysql.com/doc/refman/8.4/en/)
- [Atlas Documentation](https://atlasgo.io/docs)
- [SQLC Documentation](https://docs.sqlc.dev/)
- [Atlas README](../../tools/atlas/README.md)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
