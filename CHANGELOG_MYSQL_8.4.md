# Changelog - MySQL 8.4 Upgrade

## [Unreleased] - 2025-12-22

### Changed

#### MySQL Version Upgrade (Issue #27)
- Upgraded MySQL Docker images from 5.7 to 8.4 for all database services
  - `watch-db`: mysql:5.7 → mysql:8.4
  - `keygen-db`: mysql:5.7 → mysql:8.4
  - `sign-db`: mysql:5.7 → mysql:8.4

#### Character Set Migration
- Migrated from `utf8` to `utf8mb4` for full Unicode support
- Updated all custom.cnf files to use `utf8mb4` with `utf8mb4_unicode_ci` collation
- Updated database creation statements to use `utf8mb4`
- This aligns configuration with Go application connection strings (already using `utf8mb4`)

#### Authentication Updates
- Updated user creation syntax to MySQL 8.0+ compatible format
- Configured `mysql_native_password` as default authentication plugin for compatibility
- Replaced deprecated `GRANT ... IDENTIFIED BY` syntax with separate `CREATE USER` and `GRANT` statements
- Added `FLUSH PRIVILEGES` to ensure changes are applied

### Files Modified

1. **docker-compose.yml**
   - Updated MySQL image tags from 5.7 to 8.4 for watch-db, keygen-db, and sign-db

2. **docker/mysql/watch/conf.d/custom.cnf**
   - Changed character-set-server from utf8 to utf8mb4
   - Added collation-server=utf8mb4_unicode_ci
   - Added default-authentication-plugin=mysql_native_password

3. **docker/mysql/keygen/conf.d/custom.cnf**
   - Changed character-set-server from utf8 to utf8mb4
   - Added collation-server=utf8mb4_unicode_ci
   - Added default-authentication-plugin=mysql_native_password

4. **docker/mysql/sign/conf.d/custom.cnf**
   - Changed character-set-server from utf8 to utf8mb4
   - Added collation-server=utf8mb4_unicode_ci
   - Added default-authentication-plugin=mysql_native_password

5. **docker/mysql/watch/init.d/user.sql**
   - Updated to use CREATE USER IF NOT EXISTS with explicit authentication plugin
   - Separated user creation from privilege grants
   - Added FLUSH PRIVILEGES

6. **docker/mysql/keygen/init.d/user.sql**
   - Updated to use CREATE USER IF NOT EXISTS with explicit authentication plugin
   - Separated user creation from privilege grants
   - Added FLUSH PRIVILEGES

7. **docker/mysql/sign/init.d/user.sql**
   - Updated to use CREATE USER IF NOT EXISTS with explicit authentication plugin
   - Separated user creation from privilege grants
   - Added FLUSH PRIVILEGES

8. **docker/mysql/watch/init.d/watch.sql**
   - Updated CREATE DATABASE to use utf8mb4 with utf8mb4_unicode_ci collation

9. **docker/mysql/keygen/init.d/keygen.sql**
   - Updated CREATE DATABASE to use utf8mb4 with utf8mb4_unicode_ci collation

10. **docker/mysql/sign/init.d/sign.sql**
    - Updated CREATE DATABASE to use utf8mb4 with utf8mb4_unicode_ci collation

### Added

- **MYSQL_8.4_MIGRATION_GUIDE.md**: Comprehensive migration guide including:
  - Overview of changes
  - Step-by-step migration instructions
  - Fresh installation guide
  - Existing data migration procedure
  - Verification steps
  - Rollback procedure
  - Troubleshooting section
  - Security considerations
  - Performance notes
  - Testing checklist

### Technical Details

#### Compatibility Considerations

**Authentication Plugin:**
- Using `mysql_native_password` instead of MySQL 8.4's default `caching_sha2_password`
- Ensures compatibility with existing Go MySQL driver configurations
- No application code changes required

**Character Set:**
- `utf8mb4` is a superset of `utf8` (utf8mb3)
- Provides full Unicode support including 4-byte characters (emojis, special symbols)
- Existing data is compatible and will work without issues
- Go application connection string already uses `charset=utf8mb4`

**SQL Mode:**
- MySQL 8.4 has stricter SQL mode defaults
- Existing SQL scripts are compatible
- No schema changes required

**Reserved Keywords:**
- MySQL 8.4 has additional reserved keywords
- Current schema doesn't use any new reserved keywords
- No schema changes required

#### Breaking Changes

**Data Volume Incompatibility:**
- MySQL 8.4 cannot directly read MySQL 5.7 data files
- Data migration requires:
  - Backup via mysqldump
  - Volume recreation
  - Data restoration

**User Creation Syntax:**
- Old: `GRANT ALL PRIVILEGES ON *.* TO user@host IDENTIFIED BY 'password'`
- New: `CREATE USER IF NOT EXISTS 'user'@'host' IDENTIFIED WITH mysql_native_password BY 'password'`
- Prevents errors during container initialization

### Migration Impact

**Development:**
- Requires volume recreation and data backup/restore
- Approximately 5-10 minutes downtime per database

**Testing:**
- All wallet types need connectivity verification
- Functional testing required for key operations
- See MYSQL_8.4_MIGRATION_GUIDE.md for complete checklist

**Production:**
- Requires careful planning and scheduled maintenance window
- Full backup and rollback plan required
- Follow MYSQL_8.4_MIGRATION_GUIDE.md procedures

### Benefits

1. **Security:**
   - Active MySQL 5.7 support ended October 2023
   - MySQL 8.4 receives regular security updates
   - Modern authentication mechanisms available

2. **Performance:**
   - Improved InnoDB buffer pool management
   - Enhanced query optimizer
   - Better JSON data type performance
   - Optimized window functions

3. **Features:**
   - Full Unicode support (utf8mb4)
   - Modern SQL features
   - Better compliance with SQL standards
   - Improved error handling

4. **Compatibility:**
   - Character set now matches Go application configuration
   - Consistent configuration across all services
   - Future-proof for new MySQL features

### Testing

Recommended testing after upgrade:

1. **Container Verification:**
   - All containers start successfully
   - MySQL version is 8.4.x
   - Character set is utf8mb4

2. **Connectivity:**
   - Watch wallet connects to watch-db
   - Keygen wallet connects to keygen-db
   - Sign wallet connects to sign-db

3. **Functional Testing:**
   - Key generation (BTC, BCH, ETH, XRP)
   - Transaction creation
   - Transaction signing
   - Payment request processing

See MYSQL_8.4_MIGRATION_GUIDE.md for complete testing checklist.

### Rollback

If issues are encountered, rollback is possible by:
1. Reverting all configuration files
2. Recreating volumes
3. Restoring from MySQL 5.7 backups

See MYSQL_8.4_MIGRATION_GUIDE.md for detailed rollback procedure.

### References

- Issue: #27
- MySQL 8.4 Release Notes: https://dev.mysql.com/doc/relnotes/mysql/8.4/en/
- MySQL 8.0 Upgrade Guide: https://dev.mysql.com/doc/refman/8.0/en/upgrading.html
- Migration Guide: MYSQL_8.4_MIGRATION_GUIDE.md

### Contributors

- AI Assistant (Claude Sonnet 4.5)

---

**Note:** This upgrade addresses Issue #27 and brings the project up to date with current MySQL versions while maintaining backward compatibility through careful configuration choices.
