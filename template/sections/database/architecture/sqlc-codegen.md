### SQLC Code Generation

The project uses [SQLC](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

#### Configuration

SQLC has separate configuration files per dialect:

| Config File | Engine | Output |
|-------------|--------|--------|
| `tools/sqlc/sqlc_postgres.yml` | PostgreSQL | `internal/infrastructure/database/postgres/sqlcgen/` |
| `tools/sqlc/sqlc_mysql.yml` | MySQL | `internal/infrastructure/database/mysql/sqlcgen/` |
| `tools/sqlc/sqlc_sqlite.yml` | SQLite | `internal/infrastructure/database/sqlite/sqlcgen/` |

#### Structure

```
tools/sqlc/
├── sqlc_postgres.yml             # PostgreSQL SQLC config
├── sqlc_mysql.yml                # MySQL SQLC config
├── sqlc_sqlite.yml               # SQLite SQLC config
├── schemas/                      # Schema files (extracted from running DB)
│   ├── postgres/*.sql
│   ├── mysql/*.sql
│   └── sqlite/*.sql
└── queries/                      # SQL query files (20 per dialect)
    ├── postgres/*.sql
    └── mysql/*.sql
```

**Query files** (20 per dialect): `address`, `auth_account_key`, `auth_fullpubkey`, `btc_account_key`, `btc_tx`, `btc_tx_input`, `btc_tx_output`, `eth_account_key`, `eth_detail_tx`, `musig2_nonces`, `payment_request`, `seed`, `tx`, `xrp_account_key`, `xrp_detail_tx`, `xrp_multisig_signature`, `xrp_pending_multisig`, `xrp_regular_key`, `xrp_signer_entry`, `xrp_signer_list`

#### Common Operations

```bash
# Generate Go code
make sqlc                          # PostgreSQL (default)
make sqlc-mysql                    # MySQL
make sqlc-sqlite                   # SQLite
make sqlc-all                      # All dialects

# Validate SQL queries
make sqlc-compile                  # Check syntax (all dialects)
make sqlc-vet                      # Vet for potential issues
make sqlc-validate                 # Both compile and vet

# Extract schema from running database
make extract-sqlc-schema-all                     # PostgreSQL (default)
make extract-sqlc-schema-all DB_DIALECT=mysql     # MySQL

# Full regeneration (extract + generate)
make regenerate-sqlc-from-current-db              # PostgreSQL (default)
make regenerate-sqlc-from-current-db DB_DIALECT=mysql  # MySQL
```

#### SQL Formatting

```bash
make sqlfluff-format    # Format all SQL query files
make sqlfluff-lint      # Lint SQL query files
make sqlfluff-fix       # Format and auto-fix SQL
```

---
