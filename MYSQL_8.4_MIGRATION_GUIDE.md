# MySQL 8.4 Migration Guide

This document provides instructions for migrating from MySQL 5.7 to MySQL 8.4 for the go-crypto-wallet project.

## Overview

The project has been updated to use MySQL 8.4 for all three database services:
- `watch-db` (Watch Wallet)
- `keygen-db` (Keygen Wallet)
- `sign-db` (Sign Wallet)

## What Changed

### Docker Images
- Updated from `mysql:5.7` to `mysql:8.4` in `docker-compose.yml`

### Character Set
- Changed from `utf8` to `utf8mb4` with `utf8mb4_unicode_ci` collation
- This provides full Unicode support including 4-byte characters (emojis, etc.)
- Go application connection string already uses `utf8mb4`

### Authentication
- User creation syntax updated to MySQL 8.4 compatible format
- Using `caching_sha2_password` (MySQL 8.4 default, mysql_native_password was removed)
- Changed from deprecated `GRANT ... IDENTIFIED BY` syntax
- Go MySQL driver supports caching_sha2_password natively

### Configuration Files Updated

1. **docker-compose.yml**
   - All three database services now use `mysql:8.4`

2. **custom.cnf files** (watch/keygen/sign)
   - `character-set-server=utf8mb4`
   - `collation-server=utf8mb4_unicode_ci`
   - Removed `default-authentication-plugin` (deprecated, removed in MySQL 8.4)

3. **user.sql files** (watch/keygen/sign)
   - Updated to use `CREATE USER IF NOT EXISTS` with default authentication
   - Uses caching_sha2_password by default (mysql_native_password removed in 8.4)
   - Separated user creation from privilege grants

4. **Database initialization SQL files**
   - Updated `CREATE DATABASE` statements to use `utf8mb4`

## Migration Steps

### For Fresh Installation

If you're starting with a fresh installation, simply run:

```bash
docker-compose up -d watch-db keygen-db sign-db
```

The new MySQL 8.4 containers will be created with the correct configuration.

### For Existing Data Migration

If you have existing data in MySQL 5.7 volumes, follow these steps:

#### Step 1: Backup Existing Data

```bash
# Backup watch-db
docker-compose exec watch-db mysqldump -u hiromaily -phiromaily watch > backup_watch.sql

# Backup keygen-db
docker-compose exec keygen-db mysqldump -u hiromaily -phiromaily keygen > backup_keygen.sql

# Backup sign-db
docker-compose exec sign-db mysqldump -u hiromaily -phiromaily sign > backup_sign.sql
```

#### Step 2: Stop and Remove Old Containers

```bash
docker-compose stop watch-db keygen-db sign-db
docker-compose rm -f watch-db keygen-db sign-db
```

#### Step 3: Remove Old Volumes

**WARNING: This will delete all existing database data. Ensure you have backups.**

```bash
docker volume rm go-crypto-wallet_watch-db
docker volume rm go-crypto-wallet_keygen-db
docker volume rm go-crypto-wallet_sign-db
```

#### Step 4: Start New MySQL 8.4 Containers

```bash
docker-compose up -d watch-db keygen-db sign-db
```

Wait for containers to be healthy:

```bash
docker-compose ps
```

#### Step 5: Restore Data

```bash
# Restore watch-db
docker exec -i watch-db mysql -u hiromaily -phiromaily watch < backup_watch.sql

# Restore keygen-db
docker exec -i keygen-db mysql -u hiromaily -phiromaily keygen < backup_keygen.sql

# Restore sign-db
docker exec -i sign-db mysql -u hiromaily -phiromaily sign < backup_sign.sql
```

## Verification Steps

### 1. Check Container Status

```bash
docker-compose ps
```

All database containers should be in "Up" state.

### 2. Verify MySQL Version

```bash
docker-compose exec watch-db mysql --version
docker-compose exec keygen-db mysql --version
docker-compose exec sign-db mysql --version
```

Should output: `mysql  Ver 8.4.x`

### 3. Check Character Set Configuration

```bash
# Watch DB
docker-compose exec watch-db mysql -u hiromaily -phiromaily -e "SHOW VARIABLES LIKE 'character_set%';"

# Keygen DB
docker-compose exec keygen-db mysql -u hiromaily -phiromaily -e "SHOW VARIABLES LIKE 'character_set%';"

# Sign DB
docker-compose exec sign-db mysql -u hiromaily -phiromaily -e "SHOW VARIABLES LIKE 'character_set%';"
```

Should show `utf8mb4` for relevant variables.

### 4. Check Database Character Set

```bash
# Watch DB
docker-compose exec watch-db mysql -u hiromaily -phiromaily -e "SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'watch';"

# Keygen DB
docker-compose exec keygen-db mysql -u hiromaily -phiromaily -e "SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'keygen';"

# Sign DB
docker-compose exec sign-db mysql -u hiromaily -phiromaily -e "SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'sign';"
```

Should output:
```
DEFAULT_CHARACTER_SET_NAME: utf8mb4
DEFAULT_COLLATION_NAME: utf8mb4_unicode_ci
```

### 5. Test Application Connectivity

Run the wallet applications to verify database connectivity:

```bash
# Test watch wallet connectivity
./cmd/watch/watch --help

# Test keygen wallet connectivity
./cmd/keygen/keygen --help

# Test sign wallet connectivity
./cmd/sign/sign --help
```

### 6. Functional Testing

Perform functional tests for each wallet type:

#### Watch Wallet
- Test address import
- Test transaction monitoring
- Test payment request creation

#### Keygen Wallet
- Test key generation for supported coins (BTC, BCH, ETH, XRP)
- Test seed management
- Test address export

#### Sign Wallet
- Test transaction signing
- Test multisig operations
- Test unsigned transaction import

## Rollback Procedure

If you encounter issues and need to rollback to MySQL 5.7:

### Step 1: Stop MySQL 8.4 Containers

```bash
docker-compose stop watch-db keygen-db sign-db
docker-compose rm -f watch-db keygen-db sign-db
```

### Step 2: Revert Configuration Files

```bash
git checkout docker-compose.yml
git checkout docker/mysql/watch/conf.d/custom.cnf
git checkout docker/mysql/keygen/conf.d/custom.cnf
git checkout docker/mysql/sign/conf.d/custom.cnf
git checkout docker/mysql/watch/init.d/user.sql
git checkout docker/mysql/keygen/init.d/user.sql
git checkout docker/mysql/sign/init.d/user.sql
git checkout docker/mysql/watch/init.d/watch.sql
git checkout docker/mysql/keygen/init.d/keygen.sql
git checkout docker/mysql/sign/init.d/sign.sql
```

### Step 3: Remove MySQL 8.4 Volumes

```bash
docker volume rm go-crypto-wallet_watch-db
docker volume rm go-crypto-wallet_keygen-db
docker volume rm go-crypto-wallet_sign-db
```

### Step 4: Start MySQL 5.7 Containers

```bash
docker-compose up -d watch-db keygen-db sign-db
```

### Step 5: Restore from Backup

```bash
docker exec -i watch-db mysql -u hiromaily -phiromaily watch < backup_watch.sql
docker exec -i keygen-db mysql -u hiromaily -phiromaily keygen < backup_keygen.sql
docker exec -i sign-db mysql -u hiromaily -phiromaily sign < backup_sign.sql
```

## Known Differences from MySQL 5.7

### Authentication Plugin
- MySQL 8.4 uses `caching_sha2_password` by default
- `mysql_native_password` plugin was **removed** in MySQL 8.4
- This upgrade uses `caching_sha2_password` for all users
- Go MySQL driver (v1.5.0+) supports caching_sha2_password natively
- Provides enhanced security over the deprecated mysql_native_password

### SQL Mode
- MySQL 8.4 has stricter SQL mode by default
- `NO_ZERO_DATE` and `NO_ZERO_IN_DATE` are enabled by default
- If you encounter SQL errors, check the SQL mode: `SELECT @@sql_mode;`

### Reserved Keywords
- Some words that were not reserved in MySQL 5.7 are now reserved in MySQL 8.4
- Most common: `ADMIN`, `ARRAY`, `MEMBER`, `PERSIST`, `ROLE`
- The project's schema doesn't use these keywords

### Character Set
- `utf8` in MySQL 5.7 was actually `utf8mb3` (3-byte UTF-8)
- `utf8mb4` provides full Unicode support
- Existing `utf8` data will be compatible with `utf8mb4`

## Troubleshooting

### Issue: Container fails to start

**Symptom:**
```
ERROR: Container watch-db exited with code 1
```

**Solution:**
1. Check container logs: `docker-compose logs watch-db`
2. Verify volume permissions
3. Ensure no port conflicts: `netstat -an | grep 3307`

### Issue: Authentication error

**Symptom:**
```
Error 1045: Access denied for user 'hiromaily'@'%'
```

**Solution:**
1. Verify user creation in init scripts
2. Check authentication plugin configuration
3. Recreate containers and volumes

### Issue: Character set mismatch

**Symptom:**
```
Error: Incorrect string value
```

**Solution:**
1. Verify character set in custom.cnf
2. Check database and table character sets
3. Ensure Go connection string uses `charset=utf8mb4`

### Issue: Connection timeout

**Symptom:**
```
Error: dial tcp 127.0.0.1:3307: connect: connection refused
```

**Solution:**
1. Wait for MySQL initialization to complete (30-60 seconds)
2. Check container is running: `docker-compose ps`
3. Verify port mapping in docker-compose.yml

## Security Considerations

### Authentication Plugin

This upgrade uses `caching_sha2_password` (MySQL 8.4 default):

- **Security Benefits:**
  - More secure password hashing algorithm
  - Better protection against brute-force attacks
  - Recommended by MySQL for production use

- **Compatibility:**
  - Go MySQL driver (github.com/go-sql-driver/mysql) v1.5.0+ supports it natively
  - No code changes required if using a recent driver version
  - `mysql_native_password` was removed in MySQL 8.4 (no fallback option)

- **Connection String:**
  - Existing connection string works without modification
  - Driver automatically negotiates caching_sha2_password

### Network Exposure

The current configuration exposes database ports to localhost:
- Watch DB: 3307
- Keygen DB: 3308
- Sign DB: 3309

For production:
- Remove port mappings if not needed
- Use Docker networks for container-to-container communication
- Implement firewall rules

### Password Management

Current configuration uses plaintext passwords in:
- docker-compose.yml environment variables
- init.d/user.sql files

For production:
- Use Docker secrets or external secret management
- Implement password rotation policies
- Use strong, randomly generated passwords

## Performance Considerations

MySQL 8.4 includes several performance improvements:

- **InnoDB**: Enhanced buffer pool management
- **Query Optimizer**: Improved cost model
- **JSON**: Better JSON data type performance
- **Window Functions**: Optimized execution

No application code changes are needed to benefit from these improvements.

## Testing Checklist

- [ ] All three database containers start successfully
- [ ] MySQL version is 8.4.x
- [ ] Character set is utf8mb4
- [ ] User authentication works
- [ ] Watch wallet can connect to watch-db
- [ ] Keygen wallet can connect to keygen-db
- [ ] Sign wallet can connect to sign-db
- [ ] BTC key generation works
- [ ] BCH key generation works
- [ ] ETH key generation works
- [ ] XRP key generation works
- [ ] Transaction creation works
- [ ] Transaction signing works
- [ ] Transaction sending works
- [ ] Payment request processing works

## References

- [MySQL 8.4 Release Notes](https://dev.mysql.com/doc/relnotes/mysql/8.4/en/)
- [MySQL 8.0 Upgrade Guide](https://dev.mysql.com/doc/refman/8.0/en/upgrading.html)
- [MySQL Character Sets](https://dev.mysql.com/doc/refman/8.4/en/charset.html)
- [MySQL Authentication Plugins](https://dev.mysql.com/doc/refman/8.4/en/authentication-plugins.html)

## Support

If you encounter issues during migration:

1. Check container logs: `docker-compose logs <service-name>`
2. Review this guide's troubleshooting section
3. Open an issue on the project repository with:
   - Error messages
   - Container logs
   - Steps to reproduce
   - MySQL version: `docker-compose exec <service> mysql --version`
   - Go version: `go version`
