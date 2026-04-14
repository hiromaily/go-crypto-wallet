### 3. Makefile Targets

#### Development Workflow (Schema-First)

| Target | Description |
|---|---|
| `make atlas-fmt` | Format all HCL schema files |
| `make atlas-lint` | Lint HCL schemas (requires Docker) |
| `make atlas-validate` | Validate Atlas configuration |
| `make atlas-dev-reset [DB_DIALECT=postgres\|mysql]` | Regenerate migrations from HCL (drops and recreates) |
| `make atlas-dev-clean [DB_DIALECT=postgres\|mysql]` | Clean databases and reapply from HCL |

#### Production Workflow (Migration-History)

| Target | Description |
|---|---|
| `make atlas-migrate-diff SCHEMA=<schema> NAME=<name> [DB_DIALECT=...]` | Generate a new incremental migration |
| `make atlas-migrate-apply-all [DB_DIALECT=postgres\|mysql]` | Apply all pending migrations |
| `make atlas-migrate-status [DB_DIALECT=postgres\|mysql]` | Show migration status |
| `make atlas-migrate-hash-all [DB_DIALECT=postgres\|mysql]` | Rehash migration directory |

#### Schema Direct Apply

| Target | Description |
|---|---|
| `make atlas-schema-apply-all [DB_DIALECT=postgres\|mysql]` | Apply HCL schemas directly (bypasses migrations) |
| `make atlas-schema-apply SCHEMA=<schema> [DB_DIALECT=...]` | Apply a single schema directly |

#### Full Regeneration

| Target | Description |
|---|---|
| `make regenerate-all-from-atlas [DB_DIALECT=postgres\|mysql]` | Full workflow: reset Docker, regenerate migrations, extract SQLC schemas, generate Go code |

---
