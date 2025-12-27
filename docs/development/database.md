# Database Architecture

This document describes the database architecture and operations for the go-crypto-wallet project.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Schema Design](#schema-design)
- [Setup and Configuration](#setup-and-configuration)
- [Common Operations](#common-operations)
- [Database Management](#database-management)
- [Troubleshooting](#troubleshooting)
- [Migration Guide](#migration-guide)

## Overview

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
  wallet-db:
    image: mysql:8.4
    container_name: wallet-db
    ports:
      - "${MYSQL_PORT:-3306}:3306"
    volumes:
      - wallet-db:/var/lib/mysql
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
├── conf.d/
│   └── custom.cnf              # Server-level configuration
├── init.d/
│   └── 01_init_all_schemas.sql # Schema initialization
├── sqls/
│   ├── definition_watch.sql    # Watch schema tables
│   ├── definition_keygen.sql   # Keygen schema tables
│   ├── definition_sign.sql     # Sign schema tables
│   └── payment_request.sql     # Payment request table (watch)
├── insert/
│   └── ganache.example.sql     # Test data for Ganache
└── scripts/
    └── (utility scripts)
```

### Initialization Process

When the container starts for the first time:

1. **User Creation**: MySQL creates users via environment variables
   - `root@'%'` with password `root`
   - `hiromaily@'%'` with password `hiromaily`

2. **Schema Initialization**: Executes `01_init_all_schemas.sql`
   ```sql
   -- Create watch schema
   CREATE DATABASE `watch`;
   USE `watch`;
   source /sqls/definition_watch.sql;
   source /sqls/payment_request.sql;

   -- Create keygen schema
   CREATE DATABASE `keygen`;
   USE `keygen`;
   source /sqls/definition_keygen.sql;

   -- Create sign schema
   CREATE DATABASE `sign`;
   USE `sign`;
   source /sqls/definition_sign.sql;
   ```

3. **Configuration**: Applies server settings from `custom.cnf`
   ```ini
   [mysqld]
   character-set-server=utf8mb4
   collation-server=utf8mb4_unicode_ci
   ```

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
   docker compose up -d wallet-db
   ```

2. **Verify schemas created**:
   ```bash
   docker compose exec wallet-db mysql -uroot -proot -e "SHOW DATABASES;"
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
   docker compose exec wallet-db mysql -uroot -proot -e "SHOW VARIABLES LIKE 'character_set_server';"
   docker compose exec wallet-db mysql -uroot -proot -e "SHOW VARIABLES LIKE 'collation_server';"
   ```

   Expected: `utf8mb4` and `utf8mb4_unicode_ci`

### Application Configuration

Each wallet type (watch, keygen, sign) connects to the same database host but specifies different schema names:

**Watch Wallet** (`data/config/*_watch.toml`):
```toml
[mysql]
host = "127.0.0.1:3306"
dbname = "watch"
user = "hiromaily"
pass = "hiromaily"
```

**Keygen Wallet** (`data/config/*_keygen.toml`):
```toml
[mysql]
host = "127.0.0.1:3306"
dbname = "keygen"
user = "hiromaily"
pass = "hiromaily"
```

**Sign Wallet** (`data/config/*_sign.toml`):
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
docker compose exec wallet-db mysql -uroot -proot watch

# Access keygen schema
docker compose exec wallet-db mysql -uroot -proot keygen

# Access sign schema
docker compose exec wallet-db mysql -uroot -proot sign
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
docker compose exec wallet-db mysqldump -uroot -proot watch > backups/watch_$(date +%Y%m%d).sql

# Backup keygen schema
docker compose exec wallet-db mysqldump -uroot -proot keygen > backups/keygen_$(date +%Y%m%d).sql

# Backup sign schema
docker compose exec wallet-db mysqldump -uroot -proot sign > backups/sign_$(date +%Y%m%d).sql

# Backup all schemas in one file
docker compose exec wallet-db mysqldump -uroot -proot --databases watch keygen sign > backups/all_schemas_$(date +%Y%m%d).sql
```

### Data Restore

Restore from backup:

```bash
# Restore watch schema
docker compose exec -T wallet-db mysql -uroot -proot watch < backups/watch_20241227.sql

# Restore keygen schema
docker compose exec -T wallet-db mysql -uroot -proot keygen < backups/keygen_20241227.sql

# Restore sign schema
docker compose exec -T wallet-db mysql -uroot -proot sign < backups/sign_20241227.sql

# Restore all schemas
docker compose exec -T wallet-db mysql -uroot -proot < backups/all_schemas_20241227.sql
```

### Reset Database

Complete database reset (WARNING: deletes all data):

```bash
# Stop and remove container with volumes
docker compose down -v

# Restart - will reinitialize schemas
docker compose up -d wallet-db
```

### Reset Individual Schema

Reset specific schema while keeping others:

```bash
# Reset watch schema
docker compose exec wallet-db mysql -uroot -proot -e "DROP DATABASE watch; CREATE DATABASE watch CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
docker compose exec wallet-db mysql -uroot -proot watch < docker/mysql/sqls/definition_watch.sql
docker compose exec wallet-db mysql -uroot -proot watch < docker/mysql/sqls/payment_request.sql

# Reset keygen schema
docker compose exec wallet-db mysql -uroot -proot -e "DROP DATABASE keygen; CREATE DATABASE keygen CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
docker compose exec wallet-db mysql -uroot -proot keygen < docker/mysql/sqls/definition_keygen.sql

# Reset sign schema
docker compose exec wallet-db mysql -uroot -proot -e "DROP DATABASE sign; CREATE DATABASE sign CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
docker compose exec wallet-db mysql -uroot -proot sign < docker/mysql/sqls/definition_sign.sql
```

### Reset Payment Request Table

```bash
# Using make
make reset-payment-request-docker

# Direct command
docker compose exec wallet-db mysql -uroot -proot watch -e "$(cat ./docker/mysql/sqls/payment_request.sql)"
```

## Database Management

### View Schema Information

```bash
# List all tables in watch schema
docker compose exec wallet-db mysql -uroot -proot watch -e "SHOW TABLES;"

# List all tables in keygen schema
docker compose exec wallet-db mysql -uroot -proot keygen -e "SHOW TABLES;"

# List all tables in sign schema
docker compose exec wallet-db mysql -uroot -proot sign -e "SHOW TABLES;"

# Describe a specific table
docker compose exec wallet-db mysql -uroot -proot watch -e "DESCRIBE address;"
```

### Check Database Size

```bash
# Size of each schema
docker compose exec wallet-db mysql -uroot -proot -e "
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
docker compose exec wallet-db mysql -uroot -proot -e "SHOW PROCESSLIST;"

# Show connections per schema
docker compose exec wallet-db mysql -uroot -proot -e "
SELECT db, COUNT(*) as connections
FROM information_schema.processlist
WHERE db IN ('watch', 'keygen', 'sign')
GROUP BY db;"
```

### View Logs

```bash
# View database container logs
docker compose logs wallet-db

# Follow logs
docker compose logs -f wallet-db

# View last 100 lines
docker compose logs --tail=100 wallet-db
```

## Troubleshooting

### Container Won't Start

**Check logs**:
```bash
docker compose logs wallet-db
```

**Common issues**:
1. Port already in use:
   ```bash
   # Check what's using port 3306
   lsof -i :3306

   # Use different port
   MYSQL_PORT=3307 docker compose up -d wallet-db
   ```

2. Volume permission issues:
   ```bash
   # Remove and recreate volume
   docker compose down -v
   docker compose up -d wallet-db
   ```

### Cannot Connect to Database

**Verify container is running**:
```bash
docker compose ps wallet-db
```

**Check container health**:
```bash
docker compose exec wallet-db mysqladmin ping -uroot -proot
```

**Verify users exist**:
```bash
docker compose exec wallet-db mysql -uroot -proot -e "SELECT User, Host FROM mysql.user WHERE User IN ('root', 'hiromaily');"
```

**Test connection from host**:
```bash
mysql -h 127.0.0.1 -u hiromaily -phiromaily -P 3306 -e "SELECT 1;"
```

### Schema Not Found

**List existing schemas**:
```bash
docker compose exec wallet-db mysql -uroot -proot -e "SHOW DATABASES;"
```

**Reinitialize schemas**:
```bash
docker compose exec wallet-db mysql -uroot -proot < docker/mysql/init.d/01_init_all_schemas.sql
```

### Character Set Issues

**Check current settings**:
```bash
docker compose exec wallet-db mysql -uroot -proot -e "
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
docker compose exec wallet-db mysql -uroot -proot -e "
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
SHOW VARIABLES LIKE 'slow_query%';"
```

**View slow query log**:
```bash
docker compose exec wallet-db cat /var/lib/mysql/slow-query.log
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
docker compose up -d wallet-db
```

#### 5. Restore Data (Optional)

If you need to restore your backed-up data:

```bash
# Wait for container to initialize
sleep 30

# Restore each schema
docker compose exec -T wallet-db mysql -uroot -proot watch < migration/watch_backup.sql
docker compose exec -T wallet-db mysql -uroot -proot keygen < migration/keygen_backup.sql
docker compose exec -T wallet-db mysql -uroot -proot sign < migration/sign_backup.sql
```

#### 6. Verify Migration

```bash
# Check schemas exist
docker compose exec wallet-db mysql -uroot -proot -e "SHOW DATABASES;"

# Check tables in each schema
docker compose exec wallet-db mysql -uroot -proot watch -e "SHOW TABLES;"
docker compose exec wallet-db mysql -uroot -proot keygen -e "SHOW TABLES;"
docker compose exec wallet-db mysql -uroot -proot sign -e "SHOW TABLES;"

# Verify data (example)
docker compose exec wallet-db mysql -uroot -proot watch -e "SELECT COUNT(*) FROM address;"
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

   docker compose exec wallet-db mysqldump -uroot -proot \
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
