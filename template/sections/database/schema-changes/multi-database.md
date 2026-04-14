### Multi-Database Considerations

#### Schema Parity Requirements

**CRITICAL**: All three databases (MySQL, SQLite, PostgreSQL) MUST maintain identical:

- Table names
- Column names
- Column order
- Primary keys
- Foreign keys (where supported)
- Index names (for consistency)

**Different**: Data types are mapped to database-specific equivalents.

#### Data Type Mapping Strategy

When adding a new column, choose types carefully:

| Concept | MySQL | SQLite | PostgreSQL |
|---------|-------|--------|------------|
| Auto-incrementing ID | `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| Boolean flag | `TINYINT(1)` | `INTEGER (0/1)` | `BOOLEAN` |
| Enum values | `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| Decimal precision | `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |
| Timestamps | `DATETIME` | `TEXT (ISO8601)` | `TIMESTAMP` |
| Variable text | `VARCHAR(n)` | `TEXT` | `VARCHAR(n)` |
| Large text | `TEXT` | `TEXT` | `TEXT` |

#### Enum Handling

**MySQL**:

```sql
coin ENUM('btc','bch','eth','xrp') NOT NULL
```

**SQLite & PostgreSQL**:

```sql
coin TEXT NOT NULL CHECK (coin IN ('btc','bch','eth','xrp'))
```

**Rationale**: TEXT with CHECK constraints provides schema evolution flexibility (no ACCESS EXCLUSIVE locks in PostgreSQL).

#### Workflow for All Databases

When making schema changes:

1. **Update HCL schema** (single source of truth)
2. **Regenerate Atlas migrations** for MySQL
3. **Apply migrations** to MySQL database
4. **Extract MySQL schema** via dump
5. **Convert to SQLite schema** (manual data type mapping)
6. **Convert to PostgreSQL schema** (manual data type mapping)
7. **Regenerate sqlc code** for all three databases
8. **Update repositories** for all three databases (if needed)
9. **Test** with all three databases

#### Repository Pattern

Each database has separate repository implementations:

```
internal/infrastructure/repository/
├── watch/
│   ├── mysql/
│   │   ├── address_sqlc.go
│   │   ├── btc_tx_sqlc.go
│   │   └── ...
│   ├── sqlite/
│   │   ├── address_sqlc.go
│   │   ├── btc_tx_sqlc.go
│   │   └── ...
│   └── postgres/           # Coming soon
│       ├── address_sqlc.go
│       ├── btc_tx_sqlc.go
│       └── ...
└── cold/
    ├── mysql/
    ├── sqlite/
    └── postgres/           # Coming soon
```

**DI Container** switches implementation based on `database.type`:

```go
// internal/di/container.go
func (c *Container) CreateAddressRepository() repository.AddressRepository {
    switch c.config.Database.Type {
    case "mysql":
        return mysql.NewAddressRepositorySqlc(c.mysqlDB, c.coinTypeCode)
    case "sqlite":
        return sqlite.NewAddressRepositorySqlc(c.sqliteDB, c.coinTypeCode)
    case "postgres":
        return postgres.NewAddressRepositorySqlc(c.postgresDB, c.coinTypeCode)
    default:
        panic("unsupported database type")
    }
}
```
