### Makefile Targets

#### SQLC Code Generation

| Target | Description |
|---|---|
| `make sqlc` | Generate PostgreSQL Go code (default) |
| `make sqlc-postgres` | Alias for `make sqlc` |
| `make sqlc-mysql` | Generate MySQL Go code |
| `make sqlc-sqlite` | Generate SQLite Go code |
| `make sqlc-all` | Generate all dialects (PostgreSQL + MySQL + SQLite) |

#### Schema Dumping

| Target | Description |
|---|---|
| `make dump-schema-watch [DB_DIALECT=postgres\|mysql]` | Export watch schema from running DB |
| `make dump-schema-keygen [DB_DIALECT=postgres\|mysql]` | Export keygen schema from running DB |
| `make dump-schema-sign [DB_DIALECT=postgres\|mysql]` | Export sign schema from running DB |
| `make dump-schema-all [DB_DIALECT=postgres\|mysql]` | Export all three schemas |
| `make dump-schema-all-mysql` | Convenience alias for MySQL |

#### Schema Extraction (Dump + Normalize)

| Target | Description |
|---|---|
| `make extract-sqlc-schema-all [DB_DIALECT=postgres\|mysql]` | Dump + normalize all schemas |
| `make extract-sqlc-schema-all-mysql` | Convenience alias for MySQL |
| `make clean-sqlc-schemas [DB_DIALECT=postgres\|mysql]` | Remove old schema SQL files |

#### Validation

| Target | Description |
|---|---|
| `make sqlc-compile` | Compile SQL queries (PostgreSQL + MySQL) |
| `make sqlc-vet` | Check SQL queries for issues |
| `make sqlc-validate` | Combined compile + vet |
| `make sqlc-lint` | Alias for `sqlc-vet` |

#### Full Regeneration

| Target | Description |
|---|---|
| `make regenerate-sqlc-from-current-db [DB_DIALECT=postgres\|mysql]` | Extract schemas + generate Go code |
| `make regenerate-all-from-atlas [DB_DIALECT=postgres\|mysql]` | Full pipeline: Atlas reset + Docker reset + extract + generate |
| `make regenerate-all-from-atlas-mysql` | Convenience alias for MySQL |

> **Default dialect is PostgreSQL.** Most commands default to `DB_DIALECT=postgres` if not specified.
