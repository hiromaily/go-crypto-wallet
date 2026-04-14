### Actual Commands

#### PostgreSQL (default)

```bash
# Full regeneration after HCL schema change
make regenerate-all-from-atlas

# Or step by step:
make extract-sqlc-schema-all           # dump + normalize
make sqlc                              # generate Go code
```

#### MySQL

```bash
# Full regeneration after HCL schema change
make regenerate-all-from-atlas DB_DIALECT=mysql

# Or step by step:
make extract-sqlc-schema-all DB_DIALECT=mysql   # dump + normalize
make sqlc-mysql                                  # generate Go code
```

#### SQLite

```bash
# SQLite schemas are manually maintained, no extraction needed
make sqlc-sqlite                       # generate Go code
```

#### All Dialects

```bash
make sqlc-all                          # generate for PostgreSQL + MySQL + SQLite
make sqlc-validate                     # validate all configs
```
