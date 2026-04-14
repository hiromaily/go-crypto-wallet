### Migration Guide

#### From Old Three-Container Setup

If migrating from the previous three-container setup (`watch-db`, `keygen-db`, `sign-db`):

##### 1. Backup Existing Data

```bash
docker compose exec watch-db mysqldump -uroot -proot watch > migration/watch_backup.sql
docker compose exec keygen-db mysqldump -uroot -proot keygen > migration/keygen_backup.sql
docker compose exec sign-db mysqldump -uroot -proot sign > migration/sign_backup.sql
```

##### 2. Update Configuration

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

##### 3. Stop Old Containers

```bash
docker compose stop watch-db keygen-db sign-db
docker compose rm -f watch-db keygen-db sign-db
```

##### 4. Start New Container

```bash
# PostgreSQL (recommended)
docker compose --profile postgres up -d

# Or MySQL
docker compose --profile mysql up -d
```

##### 5. Restore Data (Optional)

```bash
# Wait for container and migrations
docker compose wait wallet-postgres-migrate-watch wallet-postgres-migrate-keygen wallet-postgres-migrate-sign
```

##### 6. Cleanup Old Volumes (Optional)

```bash
docker volume rm go-crypto-wallet_watch-db
docker volume rm go-crypto-wallet_keygen-db
docker volume rm go-crypto-wallet_sign-db
```

---
