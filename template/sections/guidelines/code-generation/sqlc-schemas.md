### SQLC Schema Files (from Database Dumps)

**Tool**: Custom shell script (`scripts/db/extract-sqlc-schema.sh`)
**Source**: MySQL database dumps (`data/dump/sql/dump_*.sql`)
**Command**: `make extract-sqlc-schema-all` (or individual: `make extract-sqlc-schema-watch`, `make extract-sqlc-schema-keygen`, `make extract-sqlc-schema-sign`)

**Generated Files**:

- `tools/sqlc/schemas/mysql/01_watch.sql` - Watch schema for SQLC
- `tools/sqlc/schemas/mysql/02_keygen.sql` - Keygen schema for SQLC
- `tools/sqlc/schemas/mysql/03_sign.sql` - Sign schema for SQLC

**Note**: These schema files are extracted from MySQL database dumps. The source of truth is the Atlas HCL files (`tools/atlas/schemas/{db_dialect}/*.hcl`). To update schemas, modify the HCL files and run the database migration flow.
