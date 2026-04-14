### 🆘 Troubleshooting

#### Migration Fails

```bash
# Check lint errors
make atlas-lint

# Reset and retry
docker compose down -v
docker compose up -d wallet-mysql
make atlas-dev-reset
make atlas-migrate-docker
```

#### SQLC Generation Fails

```bash
# Check schema syntax
docker compose exec wallet-mysql mysql -uroot -proot watch < tools/sqlc/schemas/mysql/01_watch.sql

# Run sqlc with verbose output
cd tools/sqlc && sqlc generate --experimental
```

#### Build Fails After Schema Change

```bash
# Ensure all code is regenerated
make sqlc
make sqlc-sqlite

# Verify imports and types
make go-lint
make check-build
```

#### Schema Mismatch Between Databases

```bash
# Compare schemas
diff -u tools/sqlc/schemas/mysql/01_watch.sql tools/sqlc/schemas/sqlite/01_watch.sql

# Verify data type conversions match the mapping table above
```

---

**For detailed explanations and complete workflows, see [Database Schema Changes Guide](../../../../docs/database/schema-changes.md)**
