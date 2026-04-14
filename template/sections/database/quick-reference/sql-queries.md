### 🔍 Useful SQL Queries

#### List All Tables

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db ".tables"

# PostgreSQL
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\dt"
```

#### Describe Table Structure

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE address;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db "PRAGMA table_info(address);"

# PostgreSQL
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\d address"
```

#### Check Migration Status

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch \
  -e "SELECT * FROM atlas_schema_revisions ORDER BY version DESC LIMIT 5;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db \
  "SELECT * FROM atlas_schema_revisions ORDER BY version DESC LIMIT 5;"
```
