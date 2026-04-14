### Testing Schema Changes

#### 1. Local MySQL Test

```bash
# Reset database and apply migrations
docker compose down -v
docker compose up -d wallet-mysql
make atlas-migrate-docker

# Verify schema
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE new_table;"

# Run integration tests
make integration-test
```

#### 2. SQLite E2E Test

```bash
# Test with SQLite (fast, no Docker)
make btc-e2e-reset P=1 DB=sqlite

# Verify database file
sqlite3 ./data/sqlite/btc/e2e.db "PRAGMA table_info(new_table);"
```

#### 3. PostgreSQL Test (Coming Soon)

```bash
# Reset and apply migrations
docker compose down -v
docker compose up -d wallet-db-postgres
make atlas-migrate-docker-postgres

# Verify schema
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\d new_table"

# Run integration tests
make integration-test-postgres
```

#### 4. Cross-Database Compatibility Test

Create a test that verifies schema consistency across all databases:

```go
// internal/infrastructure/repository/watch/compatibility_test.go
func TestSchemaParity(t *testing.T) {
    // Test that MySQL, SQLite, and PostgreSQL schemas are equivalent
    // - Same tables exist
    // - Same columns exist (ignoring type differences)
    // - Same primary keys
    // - Same foreign keys
}
```
