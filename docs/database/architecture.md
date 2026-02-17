# Database Architecture

This document describes the database architecture and operations for the go-crypto-wallet project.

## Table of Contents

- [Overview](#overview)
- [Supported Databases](#supported-databases)
- [Architecture](#architecture)
- [Schema Design](#schema-design)
- [Setup and Configuration](#setup-and-configuration)
- [Common Operations](#common-operations)
- [Database Management](#database-management)
- [Schema Migrations with Atlas](#schema-migrations-with-atlas)
- [SQLite for E2E Testing](#sqlite-for-e2e-testing)
- [Troubleshooting](#troubleshooting)
- [Migration Guide](#migration-guide)

## Overview

The project supports **two database backends**:

| Database | Use Case | Features |
|----------|----------|----------|
| **MySQL** | Production, full testing | Docker container, schema separation |
| **SQLite** | E2E testing, CI/CD | Local file, fast startup, no Docker DB required |

## Supported Databases

### MySQL (Production)

The project uses a **single MySQL 8.4 container** with **three separate schemas** to manage wallet data:

- **`watch` schema**: Online wallet data (addresses, transactions, payment requests)
- **`keygen` schema**: Key generation data (seeds, account keys, full public keys)
- **`sign` schema**: Signing wallet data (auth account keys, seeds)

This consolidated approach provides:
- ✅ Reduced resource usage (single MySQL instance)
- ✅ Simplified deployment and maintenance
- ✅ Data isolation through schema separation
- ✅ Easier backup and restore operations
- ✅ Single point of configuration

## Architecture

### Container Setup

```yaml
services:
  wallet-mysql:
    image: mysql:8.4
    container_name: wallet-mysql
    ports:
      - "${MYSQL_PORT:-3306}:3306"
    volumes:
      - wallet-mysql:/var/lib/mysql
      - "./docker/mysql/sqls:/sqls"
      - "./docker/mysql/conf.d:/etc/mysql/conf.d"
      - "./docker/mysql/init.d:/docker-entrypoint-initdb.d"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_USER: hiromaily
      MYSQL_PASSWORD: hiromaily
```

### Directory Structure

```
docker/mysql/
├── archive/                     # Archived SQL schema files (reference only)
│   ├── definition_watch.sql     # Original watch schema (archived)
│   ├── definition_keygen.sql    # Original keygen schema (archived)
│   ├── definition_sign.sql     # Original sign schema (archived)
│   └── payment_request.sql     # Original payment request table (archived)
├── conf.d/
│   └── custom.cnf              # Server-level configuration
├── init.d/
│   └── 01_init_all_schemas.sql # Schema initialization (creates empty schemas)
├── insert/
│   └── ganache.example.sql     # Test data for Ganache
└── scripts/
    └── (utility scripts)
```

**Note:** Schema definitions are now managed by Atlas migrations in `tools/atlas/migrations/`. The archived SQL files are kept for reference only.

### Initialization Process

When the container starts for the first time:

1. **User Creation**: MySQL creates users via environment variables
   - `root@'%'` with password `root`
   - `hiromaily@'%'` with password `hiromaily`

2. **Schema Creation**: Executes `01_init_all_schemas.sql` to create empty schemas
   ```sql
   -- Create watch schema (empty)
   CREATE DATABASE `watch` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

   -- Create keygen schema (empty)
   CREATE DATABASE `keygen` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

   -- Create sign schema (empty)
   CREATE DATABASE `sign` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

3. **Configuration**: Applies server settings from `custom.cnf`
   ```ini
   [mysqld]
   character-set-server=utf8mb4
   collation-server=utf8mb4_unicode_ci
   ```

4. **Schema Migration**: After the container starts, apply Atlas migrations
   ```bash
   make atlas-migrate-docker
   ```
   
   This will create all tables and apply schema definitions using Atlas migrations.

## Schema Design

### Watch Schema (`watch`)

**Purpose**: Manages online wallet operations including address tracking, transaction monitoring, and payment requests.

**Tables**:
- `address` - Wallet addresses for all account types
- `btc_tx` - Bitcoin/BCH transaction records
- `btc_tx_input` - Bitcoin transaction inputs
- `btc_tx_output` - Bitcoin transaction outputs
- `eth_detail_tx` - Ethereum transaction details
- `xrp_detail_tx` - XRP transaction details
- `tx` - Generic transaction records
- `payment_request` - Payment request queue

**Access Pattern**: High read/write - monitors blockchain, creates transactions

### Keygen Schema (`keygen`)

**Purpose**: Stores key generation data for offline key generation wallet.

**Tables**:
- `account_key` - Generated account keys (HD wallet)
- `auth_fullpubkey` - Full public keys for multisig authentication
- `xrp_account_key` - XRP-specific account keys
- `seed` - Encrypted seed phrases

**Access Pattern**: Write-heavy during key generation, read-only during export

**Security**: This schema contains sensitive key material - should be in offline/cold storage in production

### Sign Schema (`sign`)

**Purpose**: Stores signing wallet data for offline transaction signing.

**Tables**:
- `auth_account_key` - Authentication account keys for signing
- `seed` - Encrypted seed phrases for signing wallet

**Access Pattern**: Read-only during signing operations

**Security**: This schema contains sensitive signing keys - should be in offline/cold storage in production

## Setup and Configuration

### Initial Setup

1. **Start the database**:
   ```bash
   docker compose up -d wallet-mysql
   ```

2. **Wait for database to be ready** (about 30 seconds):
   ```bash
   docker compose exec wallet-mysql mysqladmin ping -uroot -proot --silent
   ```

3. **Apply Atlas migrations**:
   ```bash
   make atlas-migrate-docker
   ```
   
   This will create all tables and schema definitions.

4. **Verify schemas and tables created**:
   ```bash
   docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"
   docker compose exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"
   ```

   Expected output:
   ```
   Database
   keygen
   sign
   watch
   (plus system databases)
   ```

3. **Verify server configuration**:
   ```bash
   docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'character_set_server';"
   docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'collation_server';"
   ```

   Expected: `utf8mb4` and `utf8mb4_unicode_ci`

### Application Configuration

Each wallet type (watch, keygen, sign) connects to the same database host but specifies different schema names:

**Watch Wallet** (`config/wallet/*_watch.toml`):
```toml
[mysql]
host = "127.0.0.1:3306"
dbname = "watch"
user = "hiromaily"
pass = "hiromaily"
```

**Keygen Wallet** (`config/wallet/*_keygen.toml`):
```toml
[mysql]
host = "127.0.0.1:3306"
dbname = "keygen"
user = "hiromaily"
pass = "hiromaily"
```

**Sign Wallet** (`config/wallet/*_sign.toml`):
```toml
[mysql]
host = "127.0.0.1:3306"
dbname = "sign"
user = "hiromaily"
pass = "hiromaily"
```

## Common Operations

### Database Access

**Using Docker Exec**:
```bash
# Access watch schema
docker compose exec wallet-mysql mysql -uroot -proot watch

# Access keygen schema
docker compose exec wallet-mysql mysql -uroot -proot keygen

# Access sign schema
docker compose exec wallet-mysql mysql -uroot -proot sign
```

**From Host Machine**:
```bash
# Access watch schema
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 watch

# Access keygen schema
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 keygen

# Access sign schema
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 sign
```

### Schema Export (Backup)

Export schema structure without data:

```bash
# Export watch schema
make dump-schema-watch

# Export keygen schema
make dump-schema-keygen

# Export sign schema
make dump-schema-sign

# Export all schemas
make dump-schema-all
```

Output location: `data/dump/sql/dump_*.sql`

### Data Export (Full Backup)

Export schema with data:

```bash
# Backup watch schema
docker compose exec wallet-mysql mysqldump -uroot -proot watch > backups/watch_$(date +%Y%m%d).sql

# Backup keygen schema
docker compose exec wallet-mysql mysqldump -uroot -proot keygen > backups/keygen_$(date +%Y%m%d).sql

# Backup sign schema
docker compose exec wallet-mysql mysqldump -uroot -proot sign > backups/sign_$(date +%Y%m%d).sql

# Backup all schemas in one file
docker compose exec wallet-mysql mysqldump -uroot -proot --databases watch keygen sign > backups/all_schemas_$(date +%Y%m%d).sql
```

### Data Restore

Restore from backup:

```bash
# Restore watch schema
docker compose exec -T wallet-mysql mysql -uroot -proot watch < backups/watch_20241227.sql

# Restore keygen schema
docker compose exec -T wallet-mysql mysql -uroot -proot keygen < backups/keygen_20241227.sql

# Restore sign schema
docker compose exec -T wallet-mysql mysql -uroot -proot sign < backups/sign_20241227.sql

# Restore all schemas
docker compose exec -T wallet-mysql mysql -uroot -proot < backups/all_schemas_20241227.sql
```

### Reset Database

Complete database reset (WARNING: deletes all data):

```bash
# Stop and remove container with volumes
docker compose down -v

# Restart - will reinitialize schemas
docker compose up -d wallet-mysql
```

### Reset Individual Schema

Reset specific schema while keeping others:

```bash
# Reset watch schema
docker compose exec wallet-mysql mysql -uroot -proot -e "DROP DATABASE watch; CREATE DATABASE watch CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
make atlas-migrate-docker  # This will apply migrations for all schemas

# Reset keygen schema
docker compose exec wallet-mysql mysql -uroot -proot -e "DROP DATABASE keygen; CREATE DATABASE keygen CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
make atlas-migrate-docker

# Reset sign schema
docker compose exec wallet-mysql mysql -uroot -proot -e "DROP DATABASE sign; CREATE DATABASE sign CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
make atlas-migrate-docker
```

**Note:** After dropping and recreating a schema, run `make atlas-migrate-docker` to apply migrations. This will only apply migrations to schemas that need them.

### Reset Payment Request Table

```bash
# Drop and recreate the table using Atlas migration
# Or manually:
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DROP TABLE IF EXISTS payment_request;"
make atlas-migrate-docker
```

**Note:** The payment request table is now managed by Atlas migrations. See `tools/atlas/migrations/watch/` for the current schema definition.

## Database Management

### View Schema Information

```bash
# List all tables in watch schema
docker compose exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"

# List all tables in keygen schema
docker compose exec wallet-mysql mysql -uroot -proot keygen -e "SHOW TABLES;"

# List all tables in sign schema
docker compose exec wallet-mysql mysql -uroot -proot sign -e "SHOW TABLES;"

# Describe a specific table
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE address;"
```

### Check Database Size

```bash
# Size of each schema
docker compose exec wallet-mysql mysql -uroot -proot -e "
SELECT
  table_schema AS 'Schema',
  ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables
WHERE table_schema IN ('watch', 'keygen', 'sign')
GROUP BY table_schema;"
```

### Monitor Active Connections

```bash
# Show active connections
docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW PROCESSLIST;"

# Show connections per schema
docker compose exec wallet-mysql mysql -uroot -proot -e "
SELECT db, COUNT(*) as connections
FROM information_schema.processlist
WHERE db IN ('watch', 'keygen', 'sign')
GROUP BY db;"
```

### View Logs

```bash
# View database container logs
docker compose logs wallet-mysql

# Follow logs
docker compose logs -f wallet-mysql

# View last 100 lines
docker compose logs --tail=100 wallet-mysql
```

## Schema Migrations with Atlas

The project uses [Atlas](https://atlasgo.io/) for managing database schema migrations. Atlas provides version-controlled migrations, migration history tracking, and rollback capabilities.

### Installation

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

### Migration Structure

Atlas supports both HCL schema definitions and SQL migrations:

```
tools/atlas/
├── atlas.hcl              # Atlas configuration
├── schemas/               # HCL schema definitions (declarative)
│   ├── watch.hcl          # Watch schema definition
│   ├── keygen.hcl         # Keygen schema definition
│   └── sign.hcl           # Sign schema definition
├── migrations/            # SQL migration files (versioned)
│   ├── watch/             # Watch schema migrations
│   ├── keygen/            # Keygen schema migrations
│   └── sign/              # Sign schema migrations
└── README.md              # Detailed Atlas documentation
```

### HCL Schema Management

The project uses HCL (HashiCorp Configuration Language) files for declarative schema management. HCL files define the desired state of the database schema.

#### Apply HCL Schema

Apply HCL schema definitions directly to the database:

```bash
make atlas-schema-apply
```

#### Show Schema Diff

Compare the current database state with HCL schema definition:

```bash
make atlas-schema-diff SCHEMA=watch
```

#### Generate Migration from HCL Diff

Generate a migration file based on the difference between database and HCL schema:

```bash
make atlas-schema-diff-migration SCHEMA=watch
```

### Common Operations

#### Apply Migrations

Apply all pending migrations for all schemas:

```bash
# Local environment (requires Atlas CLI installed)
make atlas-migrate

# Docker environment (uses wallet-mysql-migrate service)
make atlas-migrate-docker
```

The Docker environment uses a dedicated migration service (`wallet-mysql-migrate`) that runs Atlas in a container. This service:
- Automatically waits for the database to be ready (health check)
- Runs in the same network as the database
- Has access to the `tools/atlas` directory via volume mount
- Can be run manually: `docker compose run --rm wallet-mysql-migrate migrate apply --dir file://migrations/watch --url "mysql://root:root@wallet-mysql:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"`

#### Check Migration Status

View migration status for all schemas:

```bash
# Local environment
make atlas-status

# Docker environment
make atlas-status-docker
```

This shows:
- Applied migrations
- Pending migrations
- Migration history

#### Rollback Migrations

Rollback the last migration for a specific schema:

```bash
# Local environment
make atlas-rollback SCHEMA=watch
make atlas-rollback SCHEMA=keygen
make atlas-rollback SCHEMA=sign

# Docker environment
make atlas-rollback-docker SCHEMA=watch
make atlas-rollback-docker SCHEMA=keygen
make atlas-rollback-docker SCHEMA=sign
```

#### Validate Migrations

Validate all migration files before applying:

```bash
# Local environment
make atlas-validate

# Docker environment
make atlas-validate-docker
```

#### Create New Migration

Create a new migration file:

```bash
# Local environment
make atlas-new SCHEMA=watch NAME=add_new_table
make atlas-new SCHEMA=keygen NAME=update_account_key
make atlas-new SCHEMA=sign NAME=add_index

# Docker environment
make atlas-new-docker SCHEMA=watch NAME=add_new_table
make atlas-new-docker SCHEMA=keygen NAME=update_account_key
make atlas-new-docker SCHEMA=sign NAME=add_index
```

### Migration History

Atlas automatically creates a migration history table (`atlas_schema_migrations`) in each schema to track applied migrations. This table should not be modified manually.

### Integration with Existing Setup

Atlas migrations work alongside the existing SQL files:

- **Legacy SQL files** (`docker/mysql/sqls/`): Preserved for reference and backward compatibility
- **Atlas migrations** (`tools/atlas/migrations/`): Used for version-controlled schema changes

For new schema changes:
1. Create an Atlas migration instead of modifying SQL files directly
2. Apply migrations using `make atlas-migrate`
3. Update sqlc schema files if needed for code generation

### Best Practices

1. **Always validate** migrations before applying:
   ```bash
   make atlas-validate
   ```

2. **Check status** before applying:
   ```bash
   make atlas-status
   ```

3. **Test migrations** on development database first

4. **Never modify existing migration files** - create new migrations instead

5. **Keep migrations small and focused** - one logical change per migration

6. **Document complex migrations** with comments in SQL files

### Troubleshooting Atlas

#### Migration Fails

1. Check error message for details
2. Verify database connection
3. Ensure schema exists
4. Review migration file syntax

#### Connection Issues

1. Verify MySQL is running: `docker compose ps wallet-mysql`
2. Check connection string in Makefile targets
3. Verify credentials are correct

#### Rollback Issues

1. Check migration history: `make atlas-status`
2. Verify migration file exists
3. Check database connection

For more detailed information, see [Atlas README](../../tools/atlas/README.md).

## SQLite for E2E Testing

SQLite provides a lightweight alternative for E2E testing without requiring Docker MySQL.

### Benefits

- **Faster startup**: No Docker MySQL container needed
- **Parallel testing**: Each test can use separate database files
- **Lighter CI/CD**: Reduced infrastructure requirements
- **Simpler debugging**: Direct file access

### Configuration

#### Config Files

All wallet config files support SQLite:

```yaml
database:
  type: "sqlite"  # mysql or sqlite
  mysql:
    host: "127.0.0.1:3306"
    dbname: "watch"
  sqlite:
    path: "./data/sqlite/btc/e2e.db"
    debug: true
```

#### Environment Variables

Override database type via environment variables:

```bash
export WALLET_DATABASE_TYPE=sqlite
export WALLET_DATABASE_SQLITE_PATH=./data/sqlite/btc/e2e.db
```

### E2E Script Usage

E2E scripts support the `DB_TYPE` environment variable:

```bash
# SQLite (default) - faster, no Docker MySQL
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

These schemas are converted from MySQL with the following changes:

| MySQL | SQLite |
|-------|--------|
| `ENUM('a','b')` | `TEXT CHECK(col IN ('a','b'))` |
| `AUTO_INCREMENT` | `AUTOINCREMENT` |
| `TINYINT(1)` | `INTEGER` |
| `DATETIME DEFAULT CURRENT_TIMESTAMP` | `TEXT DEFAULT CURRENT_TIMESTAMP` |

### SQLite Data Location

SQLite database files are stored in:

```
data/sqlite/
└── btc/
    └── e2e.db  # Combined E2E test database
```

**Note**: Database files are gitignored (see `.gitignore`: `data/sqlite/**/*.db`)

### SQLC Code Generation

Generate SQLite-specific SQLC code:

```bash
make sqlc-sqlite
```

Generated files: `internal/infrastructure/database/sqlite/sqlcgen/`

## Troubleshooting

### Container Won't Start

**Check logs**:
```bash
docker compose logs wallet-mysql
```

**Common issues**:
1. Port already in use:
   ```bash
   # Check what's using port 3306
   lsof -i :3306

   # Use different port
   MYSQL_PORT=3307 docker compose up -d wallet-mysql
   ```

2. Volume permission issues:
   ```bash
   # Remove and recreate volume
   docker compose down -v
   docker compose up -d wallet-mysql
   ```

### Cannot Connect to Database

**Verify container is running**:
```bash
docker compose ps wallet-mysql
```

**Check container health**:
```bash
docker compose exec wallet-mysql mysqladmin ping -uroot -proot
```

**Verify users exist**:
```bash
docker compose exec wallet-mysql mysql -uroot -proot -e "SELECT User, Host FROM mysql.user WHERE User IN ('root', 'hiromaily');"
```

**Test connection from host**:
```bash
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 -e "SELECT 1;"
```

### Schema Not Found

**List existing schemas**:
```bash
docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"
```

**Reinitialize schemas**:
```bash
docker compose exec wallet-mysql mysql -uroot -proot < docker/mysql/init.d/01_init_all_schemas.sql
```

### Character Set Issues

**Check current settings**:
```bash
docker compose exec wallet-mysql mysql -uroot -proot -e "
SHOW VARIABLES LIKE 'character_set%';
SHOW VARIABLES LIKE 'collation%';"
```

**Expected values**:
- `character_set_server`: `utf8mb4`
- `collation_server`: `utf8mb4_unicode_ci`

**Fix**: Ensure `docker/mysql/conf.d/custom.cnf` is properly mounted and restart container.

### Slow Queries

**Enable slow query log**:
```bash
docker compose exec wallet-mysql mysql -uroot -proot -e "
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
SHOW VARIABLES LIKE 'slow_query%';"
```

**View slow query log**:
```bash
docker compose exec wallet-mysql cat /var/lib/mysql/slow-query.log
```

## Migration Guide

### From Old Three-Container Setup

If migrating from the previous three-container setup (`watch-db`, `keygen-db`, `sign-db`):

#### 1. Backup Existing Data

```bash
# Backup from old containers
docker compose exec watch-db mysqldump -uroot -proot watch > migration/watch_backup.sql
docker compose exec keygen-db mysqldump -uroot -proot keygen > migration/keygen_backup.sql
docker compose exec sign-db mysqldump -uroot -proot sign > migration/sign_backup.sql
```

#### 2. Update Configuration

All configuration files have been updated in the repository. If you have custom configs, update them:

```toml
# Change from:
host = "127.0.0.1:3307"  # or 3308, 3309

# To:
host = "127.0.0.1:3306"

# Keep dbname unchanged:
dbname = "watch"  # or "keygen", "sign"
```

#### 3. Stop Old Containers

```bash
docker compose stop watch-db keygen-db sign-db
docker compose rm -f watch-db keygen-db sign-db
```

#### 4. Start New Container

```bash
docker compose up -d wallet-mysql
```

#### 5. Restore Data (Optional)

If you need to restore your backed-up data:

```bash
# Wait for container to initialize
sleep 30

# Restore each schema
docker compose exec -T wallet-mysql mysql -uroot -proot watch < migration/watch_backup.sql
docker compose exec -T wallet-mysql mysql -uroot -proot keygen < migration/keygen_backup.sql
docker compose exec -T wallet-mysql mysql -uroot -proot sign < migration/sign_backup.sql
```

#### 6. Verify Migration

```bash
# Check schemas exist
docker compose exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"

# Check tables in each schema
docker compose exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"
docker compose exec wallet-mysql mysql -uroot -proot keygen -e "SHOW TABLES;"
docker compose exec wallet-mysql mysql -uroot -proot sign -e "SHOW TABLES;"

# Verify data (example)
docker compose exec wallet-mysql mysql -uroot -proot watch -e "SELECT COUNT(*) FROM address;"
```

#### 7. Cleanup Old Volumes (Optional)

After verifying everything works:

```bash
docker volume rm go-crypto-wallet_watch-db
docker volume rm go-crypto-wallet_keygen-db
docker volume rm go-crypto-wallet_sign-db
```

## Best Practices

### Security

1. **Production Deployment**:
   - Change default passwords immediately
   - Use strong passwords for `root` and `hiromaily` users
   - Limit remote access (use `localhost` instead of `%` for Host)
   - Enable SSL/TLS for connections
   - Store `keygen` and `sign` schemas in offline/cold storage

2. **Secrets Management**:
   - Never commit passwords to version control
   - Use environment variables or secrets management tools
   - Rotate passwords regularly

### Backup Strategy

1. **Automated Backups**:
   ```bash
   # Daily backup script example
   #!/bin/bash
   BACKUP_DIR="/path/to/backups"
   DATE=$(date +%Y%m%d_%H%M%S)

   docker compose exec wallet-mysql mysqldump -uroot -proot \
     --single-transaction \
     --databases watch keygen sign \
     > "$BACKUP_DIR/wallet_backup_$DATE.sql"

   # Keep only last 30 days
   find "$BACKUP_DIR" -name "wallet_backup_*.sql" -mtime +30 -delete
   ```

2. **Backup Frequency**:
   - **watch schema**: Daily or more frequent (active transaction data)
   - **keygen schema**: After key generation operations
   - **sign schema**: After key import operations

3. **Off-site Backups**:
   - Store backups in multiple locations
   - Encrypt backups containing sensitive data (keygen, sign)

### Performance Optimization

1. **Connection Pooling**: Applications should use connection pooling

2. **Indexes**: Verify indexes exist for frequently queried columns

3. **Query Optimization**: Use `EXPLAIN` to analyze slow queries

4. **Resource Limits**: Adjust MySQL configuration for your workload
   ```ini
   # Example additional settings in custom.cnf
   [mysqld]
   max_connections = 100
   innodb_buffer_pool_size = 256M
   ```

### Monitoring

1. **Health Checks**: Container includes health check via `mysqladmin ping`

2. **Metrics**: Consider integrating with monitoring tools:
   - Prometheus + MySQL Exporter
   - Grafana dashboards
   - CloudWatch (AWS)

3. **Alerts**: Set up alerts for:
   - Database connection failures
   - Disk space usage
   - Slow queries
   - Replication lag (if using replication)

## References

- [MySQL 8.4 Documentation](https://dev.mysql.com/doc/refman/8.4/en/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Project Installation Guide](../Installation.md)
- [Issue #87: Database Consolidation](../issues/database_consolidation.md)
