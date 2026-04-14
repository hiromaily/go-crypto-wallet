### Troubleshooting

#### Container Won't Start

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

#### Cannot Connect to Database

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

#### Schema Not Found

```bash
# PostgreSQL
docker exec wallet-postgres psql -U postgres -c "\l"

# MySQL
docker exec wallet-mysql mysql -uroot -proot -e "SHOW DATABASES;"
```

#### Migration Fails

1. Check error message in migration service logs
2. Verify database connection
3. Ensure database exists
4. Review migration file syntax
5. Check migration status: `make atlas-migrate-status`

#### Character Set Issues (MySQL)

```bash
docker exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'character_set_server';"
docker exec wallet-mysql mysql -uroot -proot -e "SHOW VARIABLES LIKE 'collation_server';"
```

Expected: `utf8mb4` and `utf8mb4_unicode_ci`

---
