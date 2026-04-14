### Environments Defined in `atlas.hcl`

| Environment | Purpose |
|---|---|
| `local_postgres_watch`, `local_postgres_keygen`, `local_postgres_sign` | Local PostgreSQL development |
| `local_mysql_watch`, `local_mysql_keygen`, `local_mysql_sign` | Local MySQL development |
| `admin_postgres_watch`, `admin_postgres_keygen`, `admin_postgres_sign` | Admin (allows `ModifySchema` for drop/create) |
| `admin_mysql_watch`, `admin_mysql_keygen`, `admin_mysql_sign` | Admin (allows `ModifySchema` for drop/create) |
