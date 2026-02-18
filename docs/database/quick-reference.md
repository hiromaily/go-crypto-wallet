# Database Quick Reference Card

Quick reference for common database operations in go-crypto-wallet project.

## 🎯 Common Workflows

### Schema Change Workflow

```bash
# 1. Edit HCL schema
vim tools/atlas/schemas/watch.hcl

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
make sqlc-postgresql # PostgreSQL

# 7. Verify
make go-lint && make check-build && make gotest
```

### Add New Query

```bash
# 1. Add query to SQL file
vim tools/sqlc/queries/mysql/address.sql

# 2. Regenerate code
make sqlc && make sqlc-sqlite

# 3. Use in repository
vim internal/infrastructure/repository/watch/mysql/address_sqlc.go
```

### Database Reset

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

## 📁 File Locations

```
tools/atlas/
├── schemas/              # ✏️  EDIT HERE - Source of truth
│   ├── watch.hcl
│   ├── keygen.hcl
│   └── sign.hcl
└── migrations/           # 🔒 AUTO-GENERATED - Do not edit
    ├── watch/*.sql
    ├── keygen/*.sql
    └── sign/*.sql

tools/sqlc/
├── queries/
│   ├── mysql/            # ✏️  EDIT HERE - MySQL queries (? placeholders)
│   │   ├── address.sql
│   │   ├── btc_tx.sql
│   │   └── *.sql
│   └── postgresql/       # ✏️  EDIT HERE - PostgreSQL queries ($1,$2 placeholders)
│       └── *.sql
├── schemas/
│   ├── mysql/            # 🔄 EXTRACTED - From MySQL dump
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   ├── postgresql/       # 🔄 EXTRACTED - From PostgreSQL dump
│   │   └── *.sql
│   └── sqlite/           # ✏️  CONVERTED - Manual type mapping
│       └── *.sql

internal/infrastructure/database/
├── mysql/sqlcgen/        # 🔒 AUTO-GENERATED
├── sqlite/sqlcgen/       # 🔒 AUTO-GENERATED
└── postgresql/sqlcgen/   # 🔒 AUTO-GENERATED (coming soon)
```

**Legend**:

- ✏️  Manual editing allowed/required
- 🔒 Auto-generated - Do not edit
- 🔄 Extracted from database

## ⚡ Make Commands

### Atlas (Schema Migrations)

| Command | Description |
|---------|-------------|
| `make atlas-fmt` | Format HCL schema files |
| `make atlas-lint` | Validate HCL schema files |
| `make atlas-dev-reset` | Regenerate all migrations from scratch |
| `make atlas-migrate` | Apply migrations (local) |
| `make atlas-migrate-docker` | Apply migrations (Docker) |
| `make atlas-status` | Show migration status (local) |
| `make atlas-status-docker` | Show migration status (Docker) |
| `make atlas-validate` | Validate migration files |

### SQLC (Code Generation)

| Command | Description |
|---------|-------------|
| `make sqlc` | Generate MySQL SQLC code |
| `make sqlc-sqlite` | Generate SQLite SQLC code |
| `make sqlc-postgresql` | Generate PostgreSQL SQLC code |
| `make sqlc-all` | Generate code for all databases |

### Schema Extraction

| Command | Description |
|---------|-------------|
| `make dump-schema-watch` | Dump watch schema from MySQL |
| `make dump-schema-keygen` | Dump keygen schema from MySQL |
| `make dump-schema-sign` | Dump sign schema from MySQL |
| `make dump-schema-all` | Dump all schemas from MySQL |
| `make extract-sqlc-schema-all` | Extract SQLC-compatible schema files |

### Database Operations

| Command | Description |
|---------|-------------|
| `docker compose up -d wallet-mysql` | Start MySQL database |
| `docker compose down -v` | Stop and remove database (data lost) |
| `docker compose exec wallet-mysql mysql -uroot -proot watch` | Access watch schema |
| `docker compose logs wallet-mysql` | View database logs |

## 🗄️ Database Configuration

### MySQL (Production)

```toml
[database]
type = "mysql"

[database.mysql]
host = "127.0.0.1:3306"
dbname = "watch"  # or "keygen", "sign"
user = "hiromaily"
pass = "hiromaily"
```

### SQLite (E2E Testing)

```toml
[database]
type = "sqlite"

[database.sqlite]
path = "./data/sqlite/btc/e2e.db"
debug = true
```

### PostgreSQL (Coming Soon)

```toml
[database]
type = "postgresql"

[database.postgresql]
host = "127.0.0.1"
port = 5432
dbname = "watch"  # or "keygen", "sign"
user = "hiromaily"
pass = "hiromaily"
sslmode = "prefer"
```

## 🔄 Data Type Mapping

| Concept | MySQL | SQLite | PostgreSQL |
|---------|-------|--------|------------|
| **Auto ID** | `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| **Boolean** | `TINYINT(1)` | `INTEGER (0/1)` | `BOOLEAN` |
| **Enum** | `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| **Decimal** | `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |
| **Timestamp** | `DATETIME` | `TEXT (ISO8601)` | `TIMESTAMP` |
| **Text (sized)** | `VARCHAR(255)` | `TEXT` | `VARCHAR(255)` |
| **Text (large)** | `TEXT` | `TEXT` | `TEXT` |

## 🧪 Testing Commands

```bash
# Unit tests
make gotest

# Integration tests (MySQL)
make integration-test

# E2E tests (SQLite)
make btc-e2e-reset P=1

# E2E tests (MySQL)
make btc-e2e-reset P=1 DB=mysql

# Verify build
make go-lint
make check-build
```

## 🔍 Useful SQL Queries

### List All Tables

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch -e "SHOW TABLES;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db ".tables"

# PostgreSQL
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\dt"
```

### Describe Table Structure

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch -e "DESCRIBE address;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db "PRAGMA table_info(address);"

# PostgreSQL
docker compose exec wallet-db-postgres psql -U hiromaily -d watch -c "\d address"
```

### Check Migration Status

```bash
# MySQL
docker compose exec wallet-mysql mysql -uroot -proot watch \
  -e "SELECT * FROM atlas_schema_revisions ORDER BY version DESC LIMIT 5;"

# SQLite
sqlite3 ./data/sqlite/btc/e2e.db \
  "SELECT * FROM atlas_schema_revisions ORDER BY version DESC LIMIT 5;"
```

## ❌ Common Mistakes to Avoid

| ❌ Don't Do This | ✅ Do This Instead |
|------------------|-------------------|
| Edit migration SQL files | Edit HCL schemas, regenerate migrations |
| Edit generated SQLC code | Modify queries or schemas, regenerate code |
| Create MySQL-only schemas | Ensure SQLite/PostgreSQL equivalents exist |
| Skip `atlas-fmt` and `atlas-lint` | Always format and validate before regenerating |
| Commit without testing | Run full test cycle before commit |
| Use different column names | Maintain identical names across all databases |
| Modify only one database | Update all three databases (MySQL, SQLite, PostgreSQL) |

## 📚 Documentation Links

- **Complete Workflow**: [Database Schema Changes Guide](schema-changes.md)
- **Database Architecture**: [Development Database Docs](../development/database.md)
- **Atlas Details**: `tools/atlas/README.md`
- **Code Generation**: [Code Generation Guidelines](code-generation.md)
- **PostgreSQL Integration**: `.kiro/specs/postgresql-integration/`

## 🆘 Troubleshooting

### Migration Fails

```bash
# Check lint errors
make atlas-lint

# Reset and retry
docker compose down -v
docker compose up -d wallet-mysql
make atlas-dev-reset
make atlas-migrate-docker
```

### SQLC Generation Fails

```bash
# Check schema syntax
docker compose exec wallet-mysql mysql -uroot -proot watch < tools/sqlc/schemas/mysql/01_watch.sql

# Run sqlc with verbose output
cd tools/sqlc && sqlc generate --experimental
```

### Build Fails After Schema Change

```bash
# Ensure all code is regenerated
make sqlc
make sqlc-sqlite

# Verify imports and types
make go-lint
make check-build
```

### Schema Mismatch Between Databases

```bash
# Compare schemas
diff -u tools/sqlc/schemas/mysql/01_watch.sql tools/sqlc/schemas/sqlite/01_watch.sql

# Verify data type conversions match the mapping table above
```

---

**📘 For detailed explanations and complete workflows, see [Database Schema Changes Guide](schema-changes.md)**
