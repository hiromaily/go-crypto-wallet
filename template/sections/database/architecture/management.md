### Database Management

#### View Schema Information

**PostgreSQL**:

```bash
# List all tables in watch database
docker exec wallet-postgres psql -U postgres -d watch -c "\dt"

# Describe a specific table
docker exec wallet-postgres psql -U postgres -d watch -c "\d address"
```

**MySQL**:

```bash
# List all tables in watch schema
docker exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"

# Describe a specific table
docker exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE address;"
```

#### View Logs

```bash
# PostgreSQL logs
docker compose logs wallet-postgres
docker compose logs -f wallet-postgres

# MySQL logs
docker compose logs wallet-mysql
docker compose logs -f wallet-mysql
```

---
