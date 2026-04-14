### SQLC Configuration Summary

| Dialect | Engine | Queries | Schema Source | Output |
|---|---|---|---|---|
| PostgreSQL | `postgresql` | `queries/postgres/*.sql` | `schemas/postgres/*.sql` | `database/postgres/sqlcgen/` |
| MySQL | `mysql` | `queries/mysql/*.sql` | `schemas/mysql/*.sql` | `database/mysql/sqlcgen/` |
| SQLite | `sqlite` | `queries/mysql/*.sql` | `schemas/sqlite/*.sql` | `database/sqlite/sqlcgen/` |

> **Note**: SQLite reuses MySQL query files (`queries/mysql/*.sql`) because their SQL syntax is compatible.
