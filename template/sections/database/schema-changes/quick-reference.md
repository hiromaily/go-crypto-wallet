### Quick Reference

#### Common Commands

| Task | Command |
|------|---------|
| Format HCL schemas | `make atlas-fmt` |
| Lint HCL schemas | `make atlas-lint` |
| Regenerate all migrations | `make atlas-dev-reset` |
| Apply migrations (local) | `make atlas-migrate` |
| Apply migrations (Docker) | `make atlas-migrate-docker` |
| Generate MySQL sqlc code | `make sqlc` |
| Generate SQLite sqlc code | `make sqlc-sqlite` |
| Generate PostgreSQL sqlc code | `make sqlc-postgres` |
| Generate all sqlc code | `make sqlc-all` |
| Verify build | `make check-build` |

#### File Locations

```
tools/atlas/
├── schemas/                    # Source of truth (HCL)
│   ├── watch.hcl              # Watch wallet schema
│   ├── keygen.hcl             # Keygen wallet schema
│   └── sign.hcl               # Sign wallet schema
├── migrations/                 # Generated migrations (DO NOT EDIT)
│   ├── watch/*.sql            # Watch migrations
│   ├── keygen/*.sql           # Keygen migrations
│   └── sign/*.sql             # Sign migrations
└── atlas.hcl                  # Atlas configuration

tools/sqlc/
├── queries/
│   ├── mysql/                  # SQL queries - MySQL (? placeholders)
│   │   ├── address.sql
│   │   ├── btc_tx.sql
│   │   └── ...
│   └── postgres/             # SQL queries - PostgreSQL ($1,$2 placeholders)
│       └── ...
├── schemas/
│   ├── mysql/                  # MySQL schema files (extracted from DB)
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   ├── postgres/             # PostgreSQL schema files (extracted from DB)
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   └── sqlite/                 # SQLite schema files (manually converted)
│       ├── 01_watch.sql
│       ├── 02_keygen.sql
│       └── 03_sign.sql
├── sqlc.yml                   # MySQL sqlc config
├── sqlc_sqlite.yml            # SQLite sqlc config
└── sqlc_postgres.yml        # PostgreSQL sqlc config
```
