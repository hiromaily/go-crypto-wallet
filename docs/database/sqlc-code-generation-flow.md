# sqlc code generation flowchart

```mermaid
flowchart TD
  A[Start] --> B[前提: migration が各DBに完全適用済み\nwallet, keygen, sign]

  B --> C{対象のDB種別はどれ?}

  C --> PG[PostgreSQL]
  C --> MY[MySQL]
  C --> SQ[SQLite]

  %% PostgreSQL path
  PG --> PG1[Schema dump\nコマンド例\nmake dump-schema-postgres-all\npg_dump via docker exec wallet-postgres]

  %% MySQL path
  MY --> MY1[Schema dump\nコマンド例\nmake dump-schema-all\nmysqldump via docker exec wallet-mysql]

  %% SQLite path
  SQ --> SQ1[Schema extract by Atlas\nコマンド例\natlas schema inspect -u sqlite://...\n--format {{ sql . }}]

  %% Normalize per dialect
  PG1 --> NPG[スキーマ正規化\nscripts/db/extract-sqlc-schema-postgres.sh\n出力: tools/sqlc/schemas/postgres/*.sql]
  MY1 --> NMY[スキーマ正規化\nscripts/db/extract-sqlc-schema.sh\n出力: tools/sqlc/schemas/mysql/*.sql]
  SQ1 --> NSQ[スキーマは手動管理\ntools/sqlc/schemas/sqlite/*.sql]

  %% sqlc generation per dialect
  NPG --> GPG[sqlc generate -f sqlc_postgres.yml\nクエリ: queries/postgres/*.sql\n出力: database/postgres/sqlcgen/]
  NMY --> GMY[sqlc generate\nクエリ: queries/mysql/*.sql\n出力: database/mysql/sqlcgen/]
  NSQ --> GSQ[sqlc generate -f sqlc_sqlite.yml\nクエリ: queries/mysql/*.sql\n出力: database/sqlite/sqlcgen/]

  GPG --> Z[Done]
  GMY --> Z
  GSQ --> Z
```

## Actual Make commands

### MySQL (default)

```bash
# Full regeneration after HCL schema change
make regenerate-all-from-atlas

# Or step by step:
make extract-sqlc-schema-all    # dump + normalize
make sqlc                       # generate Go code
```

### PostgreSQL

```bash
# Full regeneration after HCL schema change
make regenerate-all-from-atlas-postgres

# Or step by step:
make extract-sqlc-schema-postgres-all  # dump + normalize
make sqlc-postgres                     # generate Go code
```

### SQLite

```bash
# SQLite schemas are manually maintained
make sqlc-sqlite                # generate Go code
```

### All dialects

```bash
make sqlc-all                   # generate for MySQL + SQLite + PostgreSQL
make sqlc-validate              # validate all configs
```
