### ⚡ Make Commands

#### Atlas (Schema Migrations)

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

#### SQLC (Code Generation)

| Command | Description |
|---------|-------------|
| `make sqlc` | Generate MySQL SQLC code |
| `make sqlc-sqlite` | Generate SQLite SQLC code |
| `make sqlc-postgres` | Generate PostgreSQL SQLC code |
| `make sqlc-all` | Generate code for all databases |

#### Schema Extraction

| Command | Description |
|---------|-------------|
| `make dump-schema-watch` | Dump watch schema from MySQL |
| `make dump-schema-keygen` | Dump keygen schema from MySQL |
| `make dump-schema-sign` | Dump sign schema from MySQL |
| `make dump-schema-all` | Dump all schemas from MySQL |
| `make extract-sqlc-schema-all` | Extract SQLC-compatible schema files |

#### Database Operations

| Command | Description |
|---------|-------------|
| `docker compose up -d wallet-mysql` | Start MySQL database |
| `docker compose down -v` | Stop and remove database (data lost) |
| `docker compose exec wallet-mysql mysql -uroot -proot watch` | Access watch schema |
| `docker compose logs wallet-mysql` | View database logs |
