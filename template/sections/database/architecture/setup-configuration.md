### Setup and Configuration

#### Initial Setup

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

#### Application Configuration

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
