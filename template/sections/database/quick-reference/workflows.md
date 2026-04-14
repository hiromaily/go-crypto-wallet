### 🎯 Common Workflows

#### Schema Change Workflow

```bash
# 1. Edit HCL schema
vim tools/atlas/schemas/{db_dialect}/watch.hcl

# 2. Format and validate
make atlas-fmt && make atlas-lint

# 3. Regenerate migrations
make atlas-dev-reset

# 4. Apply to database
docker compose down -v && docker compose up -d wallet-mysql
make atlas-migrate-docker

# 5. Extract and convert schemas
make dump-schema-all
make extract-sqlc-schema-all
# Manually convert to SQLite/PostgreSQL

# 6. Regenerate code
make sqlc        # MySQL
make sqlc-sqlite # SQLite
make sqlc-postgres # PostgreSQL

# 7. Verify
make go-lint && make check-build && make go-test
```

#### Add New Query

```bash
# 1. Add query to SQL file
vim tools/sqlc/queries/mysql/address.sql

# 2. Regenerate code
make sqlc && make sqlc-sqlite

# 3. Use in repository
vim internal/infrastructure/repository/watch/mysql/address_sqlc.go
```

#### Database Reset

```bash
# Complete reset (all data lost)
docker compose down -v
docker compose up -d wallet-mysql
make atlas-migrate-docker

# Reset specific schema
docker compose exec wallet-mysql mysql -uroot -proot \
  -e "DROP DATABASE watch; CREATE DATABASE watch;"
make atlas-migrate-docker
```
