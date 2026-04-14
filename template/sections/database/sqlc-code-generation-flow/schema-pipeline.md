### Schema Extraction Pipeline

The extract scripts transform raw database dumps into clean SQL that SQLC can parse:

```mermaid
flowchart LR
  A[Running DB Container] --> B[pg_dump / mysqldump\nSchema-only export]
  B --> C[Raw dump file\ndata/dump/sql/dump_*.dialect.sql]
  C --> D[Extract script\nRemove dialect-specific noise\nFilter out atlas_schema_revisions]
  D --> E[Clean SQL schema\ntools/sqlc/schemas/dialect/*.sql]
  E --> F[sqlc generate\nType-safe Go code]
```

**What the extract scripts do:**

1. Read raw database dump file
2. Extract only `CREATE TABLE` statements
3. Exclude internal tables (`atlas_schema_revisions`; for `sign` schema, also excludes `seed` and `musig2_nonces`)
4. Remove dialect-specific noise (PostgreSQL: `SET`, `SELECT`, `ALTER TABLE OWNER`; MySQL: `DROP TABLE`, conditional comments)
5. Remove schema prefixes (`public.` for PostgreSQL, backticks for MySQL)
6. Add `DO NOT EDIT` comment header
7. Output formatted SQL for SQLC consumption
