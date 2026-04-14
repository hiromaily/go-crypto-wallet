### Supported Databases

#### PostgreSQL (Production Default)

The project uses a **single PostgreSQL 18.2 container** with **four separate databases** (watch, keygen, sign, sign2). PostgreSQL uses named `enum` types, `identity` columns, `timestamptz`, and `bytea` for binary data.

#### MySQL (Production Alternative)

The project uses a **single MySQL 8.4 container** with **four separate schemas** (watch, keygen, sign, sign2). MySQL uses inline `enum()`, `auto_increment`, `datetime`, and `blob` for binary data.

**Note**: MySQL 8.4 uses `caching_sha2_password` authentication by default, which requires SSL for local connections. Use `?tls=true` in connection strings.

#### SQLite (E2E Testing)

SQLite provides a lightweight alternative for E2E testing without Docker. Uses `CHECK` constraints for enums, `INTEGER PRIMARY KEY` for auto-increment, and `TEXT` for timestamps.

---
