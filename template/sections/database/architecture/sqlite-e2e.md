### SQLite for E2E Testing

SQLite provides a lightweight alternative for E2E testing without requiring Docker.

#### Benefits

- **Faster startup**: No Docker container needed
- **Parallel testing**: Each test can use separate database files
- **Lighter CI/CD**: Reduced infrastructure requirements
- **Simpler debugging**: Direct file access

#### Configuration

All wallet config files support SQLite via the `database.type` field:

```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/sqlite/btc/watch.db"
    max_open_conns: 2  # Minimum 2 to prevent deadlock
    debug: true
```

#### E2E Script Usage

```bash
# SQLite (default) - faster, no Docker
make btc-e2e-reset P=1

# MySQL - traditional Docker-based testing
make btc-e2e-reset P=1 DB=mysql
```

#### SQLite Schema Files

SQLite-compatible schemas are located in:

```
tools/sqlc/schemas/sqlite/
├── 01_watch.sql   # Watch wallet schema
├── 02_keygen.sql  # Keygen wallet schema
└── 03_sign.sql    # Sign wallet schema
```

These schemas are converted from MySQL/PostgreSQL with the following changes:

| MySQL/PostgreSQL | SQLite |
|-----------------|--------|
| `ENUM('a','b')` / named enum | `TEXT CHECK(col IN ('a','b'))` |
| `AUTO_INCREMENT` / `identity` | `AUTOINCREMENT` |
| `TINYINT(1)` / `boolean` | `INTEGER` |
| `DATETIME` / `timestamptz` | `TEXT DEFAULT CURRENT_TIMESTAMP` |
| `BLOB` / `bytea` | `BLOB` |

#### SQLite Data Location

```
data/sqlite/
└── btc/
    └── watch.db  # E2E test database
```

**Note**: Database files are gitignored (`data/sqlite/**/*.db`)

---
