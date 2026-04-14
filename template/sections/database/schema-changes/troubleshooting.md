### Troubleshooting

#### Issue: Atlas Migration Fails

**Symptoms**:

```
Error: atlas migrate apply failed
```

**Solutions**:

1. **Check HCL syntax**:

   ```bash
   make atlas-lint
   ```

2. **Verify database is running**:

   ```bash
   docker compose ps wallet-mysql
   ```

3. **Check migration history**:

   ```bash
   make atlas-status-docker
   ```

4. **Reset and retry**:

   ```bash
   docker compose down -v
   docker compose up -d wallet-mysql
   make atlas-dev-reset
   make atlas-migrate-docker
   ```

#### Issue: SQLC Generation Fails

**Symptoms**:

```
Error: sqlc generate failed
```

**Solutions**:

1. **Check schema file syntax**:

   ```bash
   # Validate MySQL syntax
   docker compose exec wallet-mysql mysql -uroot -proot watch < tools/sqlc/schemas/mysql/01_watch.sql
   ```

2. **Check query file syntax**:

   ```bash
   # Run sqlc with verbose output
   cd tools/sqlc && sqlc generate --experimental
   ```

3. **Verify engine configuration**:

   ```yaml
   # tools/sqlc/sqlc_mysql.yml for MySQL
   version: "2"
   sql:
     - engine: "mysql"  # Ensure correct engine
       queries: "queries"
       schema: "schemas"
   ```

#### Issue: Schema Mismatch Between Databases

**Symptoms**:

- Tests pass with MySQL but fail with SQLite
- Repository code works differently across databases

**Solutions**:

1. **Compare schema files**:

   ```bash
   # Compare table structures
   diff -u tools/sqlc/schemas/mysql/01_watch.sql tools/sqlc/schemas/sqlite/01_watch.sql
   ```

2. **Verify data type mappings**:
   - Review [Data Type Mapping Strategy](#data-type-mapping-strategy)
   - Ensure ENUM → CHECK constraint conversion is correct

3. **Check sqlc generated models**:

   ```bash
   # Compare generated models
   diff -u internal/infrastructure/database/mysql/sqlcgen/models.go \
           internal/infrastructure/database/sqlite/sqlcgen/models.go
   ```

#### Issue: Migration Conflict

**Symptoms**:

```
Error: migration checksum mismatch
```

**Solutions**:

1. **Regenerate from scratch**:

   ```bash
   make atlas-dev-reset
   ```

2. **Clear migration history**:

   ```bash
   docker compose exec wallet-mysql mysql -uroot -proot watch \
     -e "DROP TABLE IF EXISTS atlas_schema_revisions;"
   make atlas-migrate-docker
   ```
