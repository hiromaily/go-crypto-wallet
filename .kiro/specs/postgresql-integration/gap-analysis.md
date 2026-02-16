# Gap Analysis: PostgreSQL Integration

## Executive Summary

This analysis examines the implementation gap for adding PostgreSQL support to go-crypto-wallet, which currently supports MySQL and SQLite. The project follows Clean Architecture with repository pattern, using sqlc for type-safe SQL generation and Atlas for schema migrations.

**Scope**: Add PostgreSQL 18.2 as a third database option alongside MySQL and SQLite, and upgrade Atlas from 1.0 to 1.1.

**Key Findings**:
- **Existing Pattern**: Well-established dual-database support provides clear blueprint for PostgreSQL addition
- **Implementation Approach**: Extend existing components following the MySQL/SQLite pattern (Option A + minimal new components)
- **Effort Estimate**: M (Medium, 3-7 days) - Primarily mechanical replication of existing patterns
- **Risk Level**: Low - Clear precedent exists, familiar technologies, minimal architectural changes

## Current State Analysis

### 1. Configuration Layer (`pkg/config/`)

**Existing Assets**:
- `pkg/config/wallet.go`: Database configuration with type selection
  - `Database.Type` validates `oneof=mysql sqlite` (line 145)
  - Separate `MySQL` and `SQLite` structs with connection parameters
  - Validation logic in `validateDatabase()` (lines 264-291)

**Pattern Identified**:
```go
type Database struct {
    Type   string `validate:"required,oneof=mysql sqlite"`
    MySQL  MySQL
    SQLite SQLite
}
```

**Gap**: No PostgreSQL struct or validation support

---

### 2. Database Connection Layer (`pkg/db/`)

**Existing Assets**:
- `pkg/db/mysql/connection.go`: MySQL connection factory
  - Uses `sql.Open("mysql", dsn)`
  - Connection pool configuration
  - Ping verification
- `pkg/db/sqlite/connection.go`: SQLite connection factory
  - Uses `sql.Open("sqlite", path)`
  - SQLite-specific pragmas (foreign keys, WAL mode, busy timeout)

**Pattern Identified**:
- Package per database type
- Factory function `New{Database}(conf *config.{Database}) (*sql.DB, error)`
- Database-specific optimizations

**Gap**: No `pkg/db/postgresql/` package

---

### 3. SQLC Configuration (`tools/sqlc/`)

**Existing Assets**:
- `sqlc.yml`: MySQL configuration
  - Engine: `mysql`
  - Schema: `./schemas/*.sql`
  - Output: `internal/infrastructure/database/mysql/sqlcgen`
- `sqlc_sqlite.yml`: SQLite configuration
  - Engine: `sqlite`
  - Schema: `./schemas_sqlite/*.sql`
  - Output: `internal/infrastructure/database/sqlite/sqlcgen`
- Shared queries: `./queries/*.sql` (used by both engines)

**Pattern Identified**:
- Separate config file per database
- Separate schema directory per database
- Same query files reused across all databases
- Generated code in engine-specific directories

**Gap**: No `sqlc_postgresql.yml` or `schemas_postgresql/` directory

---

### 4. Schema Definitions

**Existing Assets**:
- MySQL schemas: `tools/sqlc/schemas/*.sql`
  - Auto-generated from Atlas migrations via `make extract-sqlc-schema-all`
  - Contains MySQL-specific syntax (AUTO_INCREMENT, ENGINE=InnoDB, enum types)
- SQLite schemas: `tools/sqlc/schemas_sqlite/*.sql`
  - Manually created SQLite-compatible versions
  - Different data types (INTEGER PRIMARY KEY instead of AUTO_INCREMENT)

**Pattern Identified**:
- MySQL schemas are source of truth (generated from HCL via Atlas)
- SQLite schemas are adapted versions
- Three schema files: 01_watch.sql, 02_keygen.sql, 03_sign.sql

**Gap**: No PostgreSQL schema files with PostgreSQL-specific syntax

**Data Type Mappings Needed**:
| MySQL | PostgreSQL |
|-------|------------|
| AUTO_INCREMENT | SERIAL / BIGSERIAL |
| tinyint(1) | BOOLEAN |
| enum('a','b') | TEXT CHECK(...) or custom ENUM type |
| decimal(26,10) | NUMERIC(26,10) |
| datetime | TIMESTAMP |
| varchar(n) | VARCHAR(n) |
| text | TEXT |
| bigint | BIGINT |

---

### 5. Atlas Configuration (`tools/atlas/`)

**Existing Assets**:
- `atlas.hcl`: Atlas 1.0.0 configuration
  - MySQL-only environments (local_watch, local_keygen, local_sign)
  - Dev database: `docker://mysql/8/`
  - HCL schemas: `schemas/watch.hcl`, `schemas/keygen.hcl`, `schemas/sign.hcl`
- Migrations: `migrations/watch/`, `migrations/keygen/`, `migrations/sign/`

**Pattern Identified**:
```hcl
env "local_watch" {
  url = "mysql://root:root@127.0.0.1:3306/watch?..."
  src = "file://schemas/watch.hcl"
  schemas = ["watch"]
  migration {
    dir = "file://migrations/watch"
  }
  dev = "docker://mysql/8/watch"
}
```

**Gap**:
- No PostgreSQL environments
- Atlas 1.0.0 (needs upgrade to 1.1.0)
- No PostgreSQL dev database configuration

---

### 6. Repository Layer (`internal/infrastructure/repository/`)

**Existing Assets**:
- Interface definitions: `internal/application/ports/repository/`
  - `watch/`: Address, BTC/ETH/XRP transactions, payments
  - `cold/`: Account keys, seed, auth, XRP signer lists
- MySQL implementations: `repository/{cold,watch}/mysql/*_sqlc.go`
  - Pattern: `New{Entity}RepositorySqlc(dbConn *sql.DB, coinTypeCode) *{Entity}RepositorySqlc`
  - Uses `sqlcgen.New(dbConn)` to create queries
  - Converter functions: `convertTo{Entity}`, `convertFrom{Entity}`
- SQLite implementations: `repository/{cold,watch}/sqlite/*_sqlc.go`
  - Identical structure to MySQL
  - Same interface, different sqlcgen package

**Pattern Identified**:
- One repository file per entity per database
- Private converter functions (domain ↔ sqlcgen types)
- Repository struct holds `queries *sqlcgen.Queries`

**Count**: ~20 repository interfaces × 2 databases = ~40 implementation files

**Gap**: No PostgreSQL repository implementations (~40 files needed)

---

### 7. Dependency Injection (`internal/di/container.go`)

**Existing Assets**:
- Database selection via `switch c.conf.Database.Type` (appears ~20 times)
- Pattern:
```go
func (c *container) newBTCTxRepo() repowatch.BTCTxRepositorier {
    switch c.conf.Database.Type {
    case "mysql":
        return watchmysql.NewBTCTxRepositorySqlc(
            c.pkgContainer.NewDatabaseClient(),
            c.conf.CoinTypeCode,
        )
    case "sqlite":
        return watchsqlite.NewBTCTxRepositorySqlc(
            c.pkgContainer.NewSQLiteClient(),
            c.conf.CoinTypeCode,
        )
    default:
        panic("unsupported database type: " + c.conf.Database.Type)
    }
}
```

**Gap**: Each switch statement needs PostgreSQL case added

---

### 8. Build System (`make/`)

**Existing Assets**:
- `make/db_sqlc.mk`:
  - `sqlc-compile`, `sqlc-vet`: Validate both MySQL and SQLite
  - `extract-sqlc-schema-all`: Extract MySQL schemas from database
  - Pattern: Run commands for each database variant
- `make/db_atlas.mk`:
  - `atlas-fmt`, `atlas-lint`: Format and validate HCL schemas
  - `atlas-dev-reset`: Regenerate migrations from HCL
  - Only MySQL environments configured

**Gap**:
- sqlc targets need PostgreSQL variant
- Atlas targets need PostgreSQL environments
- Schema extraction needs PostgreSQL support

---

### 9. Docker Compose (`compose.yaml`)

**Existing Assets**:
- `wallet-mysql` service: MySQL 8.4 container
- Migration services: `wallet-mysql-migrate-{watch,keygen,sign}` using Atlas 1.0.0
- Pattern:
```yaml
x-migration-base: &migration-base
  image: arigaio/atlas:1.0.0
  depends_on:
    wallet-mysql:
      condition: service_healthy
```

**Gap**:
- No PostgreSQL 18.2 container
- No PostgreSQL migration services
- Atlas image version 1.0.0 (needs upgrade to 1.1.0)

---

## Requirements Feasibility Analysis

### Technical Needs by Requirement

#### R1: Configuration Support
- **Needs**: PostgreSQL struct, validation logic, SSLMode parameter
- **Complexity**: Simple - extend existing pattern
- **Unknowns**: None

#### R2: PostgreSQL Schema Generation
- **Needs**: Data type conversion script/tool, PostgreSQL schemas in `schemas_postgresql/`
- **Complexity**: Moderate - requires careful type mapping
- **Unknowns**:
  - Research Needed: Enum handling strategy (TEXT CHECK vs PostgreSQL ENUM)
  - Research Needed: Decimal precision compatibility

#### R3: SQLC Code Generation
- **Needs**: `sqlc_postgresql.yml`, PostgreSQL sqlcgen package
- **Complexity**: Simple - mechanical replication
- **Unknowns**: None (sqlc supports PostgreSQL engine)

#### R4: Database Connection Management
- **Needs**: PostgreSQL driver, connection factory, SSL configuration
- **Complexity**: Simple - use `github.com/lib/pq` or `pgx`
- **Unknowns**:
  - Research Needed: Driver selection (lib/pq vs pgx - recommend pgx for better performance and features)

#### R5: Build Integration
- **Needs**: Makefile targets for PostgreSQL sqlc/atlas operations
- **Complexity**: Simple - extend existing targets
- **Unknowns**: None

#### R6: Migration Path
- **Needs**: Data export/import scripts, type compatibility testing
- **Complexity**: Moderate - requires data validation
- **Unknowns**:
  - Research Needed: Best practices for MySQL→PostgreSQL data migration
  - Research Needed: Enum value migration strategy

#### R7: Atlas Migration Support
- **Needs**: PostgreSQL environments in atlas.hcl, HCL schemas (optional), migrations
- **Complexity**: Simple - replicate MySQL pattern
- **Unknowns**: None (Atlas supports PostgreSQL)

#### R8: Docker Compose Integration
- **Needs**: PostgreSQL 18.2 container, healthcheck, init scripts, migration services
- **Complexity**: Simple - replicate MySQL pattern
- **Unknowns**: None

#### R9: Testing
- **Needs**: Integration tests, test database setup
- **Complexity**: Moderate - adapt existing tests
- **Unknowns**: None (can reuse existing test patterns)

#### R10: Atlas Upgrade
- **Needs**: Update version references, verify compatibility
- **Complexity**: Simple - version bump
- **Unknowns**:
  - Research Needed: Atlas 1.0 → 1.1 breaking changes and migration path

---

## Implementation Approach Analysis

### Option A: Extend Existing Components (RECOMMENDED)

**Strategy**: Follow the exact pattern used for MySQL/SQLite dual support

**Components to Extend**:

1. **Configuration** (`pkg/config/wallet.go`):
   - Add `PostgreSQL` struct (similar to MySQL)
   - Update `Database.Type` validation to `oneof=mysql sqlite postgresql`
   - Add PostgreSQL validation in `validateDatabase()`

2. **Connection Factory** (new package `pkg/db/postgresql/`):
   - Create `connection.go` with `NewPostgreSQL(conf *config.PostgreSQL) (*sql.DB, error)`
   - Use `pgx` driver for better PostgreSQL support
   - Configure SSL mode, connection pooling

3. **SQLC Configuration** (`tools/sqlc/`):
   - Create `sqlc_postgresql.yml` with engine: postgresql
   - Create `schemas_postgresql/` directory with converted schemas
   - Generate to `internal/infrastructure/database/postgresql/sqlcgen`

4. **Repository Implementations** (new packages):
   - `internal/infrastructure/repository/cold/postgresql/*_sqlc.go` (~15 files)
   - `internal/infrastructure/repository/watch/postgresql/*_sqlc.go` (~12 files)
   - Copy-paste from MySQL implementations, update import paths

5. **DI Container** (`internal/di/container.go`):
   - Add `case "postgresql":` to each database type switch (~20 locations)

6. **Atlas Configuration** (`tools/atlas/atlas.hcl`):
   - Add PostgreSQL environments (local_postgresql_watch, etc.)
   - Update dev database to `docker://postgres/18`
   - Optionally create PostgreSQL-specific HCL schemas or reuse existing

7. **Docker Compose** (`compose.yaml`):
   - Add `wallet-db-postgresql` service (PostgreSQL 18.2)
   - Add migration services for PostgreSQL
   - Update Atlas image to 1.1.0

8. **Makefile** (`make/db_sqlc.mk`, `make/db_atlas.mk`):
   - Extend sqlc targets to include PostgreSQL
   - Add PostgreSQL schema extraction (if needed)
   - Extend atlas targets for PostgreSQL environments

**Trade-offs**:
- ✅ Minimal risk - proven pattern
- ✅ Consistent with existing architecture
- ✅ Fast development - copy-paste with modifications
- ✅ Easy review - clear diff from MySQL code
- ❌ Code duplication (but acceptable given repository pattern)
- ❌ ~40 new repository files (but well-isolated)

**Compatibility Assessment**:
- No breaking changes to existing MySQL/SQLite users
- New database type is opt-in via configuration
- Existing interfaces unchanged
- Test coverage remains independent per database

---

### Option B: Create Abstraction Layer

**Strategy**: Create database-agnostic repository layer to reduce duplication

**Components**:
- Generic repository implementations using reflection or code generation
- Database-specific adapters for sqlcgen types
- Unified sqlcgen wrapper

**Trade-offs**:
- ✅ Reduces code duplication long-term
- ✅ Easier to add future databases
- ❌ Significant refactoring of existing code
- ❌ Adds complexity and indirection
- ❌ Higher risk of bugs
- ❌ Not aligned with current architecture
- ❌ Much larger effort (2-3 weeks)

**Recommendation**: NOT RECOMMENDED for this iteration. Current pattern is working well.

---

### Option C: Hybrid - Extend + Shared Utilities

**Strategy**: Extend existing components (Option A) + create shared utilities for common operations

**Additional Components**:
- Shared converter utilities for common patterns
- Schema conversion tool (MySQL → PostgreSQL)
- Test helpers for multi-database testing

**Trade-offs**:
- ✅ Reduces some duplication
- ✅ Provides tooling for maintenance
- ✅ Maintains consistency with existing pattern
- ❌ Slightly more complex than pure Option A
- ⚠️ May be premature optimization

**Recommendation**: Consider for schema conversion tool only, keep repositories as-is

---

## Effort and Risk Assessment

### Effort Estimate: M (Medium, 3-7 days)

**Breakdown by Requirement**:
| Requirement | Effort | Justification |
|-------------|--------|---------------|
| R1: Configuration | S (0.5 day) | Simple struct addition and validation |
| R2: Schema Generation | M (1-2 days) | Manual conversion + validation of 3 schema files |
| R3: SQLC Generation | S (0.5 day) | Config file + run generation |
| R4: Connection Management | S (0.5 day) | Simple factory pattern |
| R5: Build Integration | S (0.5 day) | Extend existing Makefile targets |
| R6: Migration Path | M (1 day) | Documentation + migration scripts |
| R7: Atlas Migration | S (0.5 day) | Add environments to atlas.hcl |
| R8: Docker Compose | S (0.5 day) | Add PostgreSQL container + migration services |
| R9: Testing | M (1-2 days) | Integration tests + validation |
| R10: Atlas Upgrade | S (0.5 day) | Version update + verification |
| **Repository Implementations** | M (2 days) | ~40 files, but copy-paste with modifications |

**Total**: 3-7 days for experienced developer familiar with the codebase

---

### Risk Assessment: Low

**Risk Factors**:

| Risk | Level | Mitigation |
|------|-------|------------|
| Data type incompatibility | Medium | Thorough schema validation, integration tests |
| Enum handling differences | Medium | Research PostgreSQL enum vs CHECK constraints |
| Query syntax differences | Low | sqlc validates syntax per engine |
| Atlas 1.0→1.1 breaking changes | Low | Check Atlas release notes, test migrations |
| PostgreSQL driver issues | Low | Use well-maintained pgx driver |
| Connection pool configuration | Low | Follow PostgreSQL best practices |
| Migration data loss | Medium | Test migration scripts on copies, provide rollback |
| Repository implementation bugs | Low | Copy proven MySQL code, unit test each |

**Overall**: Low risk due to:
- Clear precedent with MySQL/SQLite
- Familiar technologies (PostgreSQL, sqlc, Atlas)
- No architectural changes required
- Incremental approach possible (add PostgreSQL without affecting existing DBs)

---

## Research Items for Design Phase

1. **PostgreSQL Driver Selection**:
   - Compare `lib/pq` vs `pgx` (recommend pgx for better features)
   - SSL/TLS configuration requirements
   - Connection string format

2. **Enum Handling Strategy**:
   - PostgreSQL custom ENUM types vs TEXT with CHECK constraints
   - Migration path from MySQL enums
   - sqlc support for each approach

3. **Decimal Precision**:
   - Verify `NUMERIC(26,10)` handles same range as MySQL `DECIMAL(26,10)`
   - Test edge cases for cryptocurrency amounts

4. **Atlas 1.0 → 1.1 Changes**:
   - Review Atlas 1.1 release notes
   - Identify breaking changes affecting existing workflows
   - Test migration compatibility

5. **Data Migration Tools**:
   - Evaluate pgloader for MySQL→PostgreSQL data migration
   - Test with sample wallet data
   - Document edge cases and manual steps

6. **PostgreSQL 18.2 Features**:
   - Check if any new features benefit this use case
   - Verify Docker image availability and stability

---

## Recommendations for Design Phase

### Preferred Approach

**Option A (Extend Existing Components)** is strongly recommended because:
- Proven pattern exists (MySQL/SQLite)
- Low risk, fast implementation
- Maintains architectural consistency
- No impact on existing database users

### Key Design Decisions Needed

1. **PostgreSQL Driver**: Select pgx over lib/pq for better performance and native features
2. **Enum Strategy**: Use TEXT with CHECK constraints (simpler migration, better compatibility)
3. **Schema Source**: Manually create PostgreSQL schemas (don't auto-generate from MySQL)
4. **Atlas HCL**: Reuse existing HCL schemas (engine-agnostic), create PostgreSQL-specific migrations
5. **Migration Path**: Provide pgloader-based scripts with manual verification steps

### Implementation Order

**Phase 1: Foundation** (1-2 days)
1. Atlas upgrade to 1.1.0
2. PostgreSQL Docker container + healthcheck
3. Configuration structs and validation
4. Connection factory

**Phase 2: Schema & Code Generation** (2-3 days)
5. Convert schemas to PostgreSQL syntax
6. Create sqlc_postgresql.yml
7. Generate PostgreSQL sqlcgen code
8. Create Atlas PostgreSQL environments

**Phase 3: Repository Layer** (2 days)
9. Implement PostgreSQL repositories (copy-paste + modify)
10. Update DI container switches
11. Update Makefile targets

**Phase 4: Testing & Documentation** (1-2 days)
12. Integration tests with PostgreSQL
13. Data migration scripts and documentation
14. Update README and configuration examples

### Success Criteria

- [ ] All 10 requirements fully implemented
- [ ] Integration tests pass with PostgreSQL backend
- [ ] Data migration from MySQL to PostgreSQL verified
- [ ] Atlas 1.1.0 works with all databases (MySQL, SQLite, PostgreSQL)
- [ ] Docker Compose supports all three databases simultaneously
- [ ] Documentation complete for PostgreSQL setup and migration
- [ ] No regressions in MySQL or SQLite functionality

---

## Constraint Summary

**Architectural Constraints**:
- Must follow Clean Architecture (domain ↔ application ↔ infrastructure)
- Must use repository pattern (one impl per entity per database)
- Must maintain interface compatibility (no breaking changes)

**Technical Constraints**:
- Must use sqlc for type-safe SQL generation
- Must use Atlas for schema migrations
- Must support offline wallets (keygen/sign) and online wallet (watch)
- Must handle three separate database schemas (watch, keygen, sign)

**Compatibility Constraints**:
- Must not affect existing MySQL/SQLite users
- Must maintain data consistency across databases
- Must preserve decimal precision for cryptocurrency amounts
- Must support concurrent operation of multiple database types (for migration)
