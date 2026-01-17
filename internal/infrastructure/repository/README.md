# Repository Infrastructure

This directory contains database-specific repository implementations following the Repository pattern.

## Directory Structure

```
repository/
├── cold/                    # Cold wallet repositories (keygen/sign)
│   ├── mysql/              # MySQL implementations
│   ├── sqlite/             # SQLite implementations
│   └── mocks/              # Test mocks
├── watch/                   # Watch wallet repositories
│   ├── mysql/              # MySQL implementations
│   ├── sqlite/             # SQLite implementations
│   └── mocks/              # Test mocks
├── doc.go                   # Package documentation
└── README.md                # This file
```

## Why Separate mysql/ and sqlite/ Directories?

### TL;DR

Go generics **cannot** consolidate these implementations because sqlc generates fundamentally different types for each database engine.

### Technical Explanation

sqlc generates database-specific types that differ in:

| Aspect | MySQL | SQLite |
|--------|-------|--------|
| ENUM columns | Custom types (e.g., `BtcAccountKeyCoin`) | `string` |
| Integer columns | `int8`, `int16`, etc. | `int64` |
| Timestamps | `sql.NullTime` | `sql.NullString` (TEXT) |
| Param structs | Database-specific field types | Database-specific field types |

#### Example: Generated Param Types

```go
// MySQL (internal/infrastructure/database/mysql/sqlcgen/)
type GetMaxBtcAccountKeyIndexParams struct {
    Coin    BtcAccountKeyCoin      // Custom ENUM type
    Account BtcAccountKeyAccount   // Custom ENUM type
}

// SQLite (internal/infrastructure/database/sqlite/sqlcgen/)
type GetMaxBtcAccountKeyIndexParams struct {
    Coin    string
    Account string
}
```

#### Example: Timestamp Handling

```go
// MySQL - direct time.Time access
if sqlcKey.UpdatedAt.Valid {
    key.UpdatedAt = &sqlcKey.UpdatedAt.Time
}

// SQLite - requires parsing TEXT format
if sqlcKey.UpdatedAt.Valid {
    t, err := time.Parse("2006-01-02 15:04:05", sqlcKey.UpdatedAt.String)
    key.UpdatedAt = &t
}
```

### Why Generics Don't Help

Go generics work for algorithm abstraction over types, but cannot:

1. **Parameterize struct field types** - Each sqlcgen model has different field types
2. **Unify method signatures** - `*sqlcgen.Queries` methods accept different param types
3. **Abstract type conversions** - Conversion logic depends on concrete types

### Design Decision

We intentionally maintain separate implementations because:

1. **Type Safety**: sqlc's type-specific generation catches errors at compile time
2. **Database Semantics**: MySQL ENUMs and SQLite TEXT have different behaviors
3. **Optimization**: Each implementation can use database-specific features
4. **Clarity**: Explicit separation makes debugging easier

### Reducing Duplication

To minimize code duplication while keeping type safety:

1. **Shared validation logic** can be extracted to helper packages
2. **Domain conversion logic** that doesn't touch sqlcgen types can be shared
3. **Test utilities** can be shared via the `mocks/` directory

## Repository Implementation Guidelines

See [Repository Pattern Rules](../../../.claude/rules/go/repository.md) for implementation details.

### Key Points

- Repository interfaces are defined in `internal/application/ports/repository/`
- Implementations convert between domain entities and sqlcgen types
- Use `context.Context` for all database operations
- Wrap errors with operation context using `fmt.Errorf("operation: %w", err)`

## Related Documentation

- [Database Guidelines](../../../docs/guidelines/database.md)
- [Repository Pattern Rules](../../../.claude/rules/go/repository.md)
- [sqlc Configuration](../../../tools/sqlc/)
- [Issue #385](https://github.com/hiromaily/go-crypto-wallet/issues/385) - SQLite implementation tracking
