# Technical Design: PostgreSQL Integration

## Overview

This feature adds PostgreSQL 18.2 as a third database backend option to go-crypto-wallet, which currently supports MySQL and SQLite. Users will be able to select PostgreSQL by setting `database.type = "postgresql"` in the wallet configuration file, with the application transparently handling PostgreSQL-specific connection management, schema generation, and query execution.

**Purpose**: Enable wallet operators to leverage PostgreSQL's advanced features, reliability, and performance characteristics for their wallet infrastructure, while maintaining complete feature parity with existing MySQL and SQLite implementations.

**Users**: Wallet operators deploying production watch/keygen/sign wallets will benefit from PostgreSQL's ACID compliance, advanced indexing, and operational tooling. Developers will use PostgreSQL in local Docker environments for testing and development.

**Impact**: This change extends the existing database abstraction layer without modifying the core wallet logic. All three databases (MySQL, SQLite, PostgreSQL) will coexist with identical query interfaces and schemas, allowing seamless database migration and selection based on deployment requirements.

### Goals

- Add PostgreSQL 18.2 as a fully supported database backend with feature parity to MySQL/SQLite
- Upgrade Atlas migration tool from version 1.0 to 1.1 for improved PostgreSQL support
- Maintain backward compatibility with existing MySQL and SQLite deployments
- Provide comprehensive migration documentation and tooling for database transitions
- Ensure type-safe query execution through sqlc code generation
- Enable local development with PostgreSQL via Docker Compose

### Non-Goals

- Database abstraction layer refactoring (maintain existing repository pattern)
- PostgreSQL-specific optimization features beyond standard functionality
- Automatic data migration between databases (provide manual tooling and documentation)
- Performance benchmarking across database engines (focus on functional equivalence)
- Multi-database support in a single wallet instance (one database type per deployment)

## Architecture

### Existing Architecture Analysis

The go-crypto-wallet project follows **Clean Architecture** with strict layer separation:

**Layer Structure**:
- **Domain** (`internal/domain/`): Pure business logic with ZERO infrastructure dependencies
- **Application** (`internal/application/`): Use case orchestration and port interfaces
- **Infrastructure** (`internal/infrastructure/`): Database, API clients, and external service implementations
- **Interface Adapters** (`internal/interface-adapters/`): CLI commands and wallet adapters

**Current Database Support Pattern**:
- Database selection via `config.Database.Type` field (currently `mysql` or `sqlite`)
- Separate repository implementations per database engine (repository pattern)
- sqlc-generated type-safe query code in database-specific packages:
  - `internal/infrastructure/database/mysql/sqlcgen/`
  - `internal/infrastructure/database/sqlite/sqlcgen/`
- Shared query definitions in `tools/sqlc/queries/*.sql`
- DI container (`internal/di/container.go`) switches implementation based on config
- Three database schemas: watch (online), keygen (offline), sign (offline)

**Integration Points**:
- Configuration validation in `pkg/config/wallet.go`
- Connection factories in `pkg/db/{mysql,sqlite}/`
- Repository interfaces in `internal/application/ports/repository/`
- Repository implementations in `internal/infrastructure/repository/{cold,watch}/{mysql,sqlite}/`
- Build system targets in `make/db_sqlc.mk` and `make/db_atlas.mk`

**Constraints to Preserve**:
- Clean Architecture dependency direction (Infrastructure→Application→Domain)
- Repository pattern (one implementation per database per entity)
- Type safety through sqlc code generation
- Independent wallet schemas (watch, keygen, sign)
- Offline wallet isolation (keygen/sign have no network dependencies)

### Architecture Pattern & Boundary Map

**Selected Pattern**: Repository Pattern with Database-Specific Implementations

**Rationale**: Extends the proven MySQL/SQLite dual-database pattern to include PostgreSQL. This approach provides clear type safety, explicit error handling, and independent database evolution while maintaining architectural consistency.

```mermaid
graph TB
    subgraph "Interface Adapters Layer"
        CLI[CLI Commands]
    end

    subgraph "Application Layer"
        UseCase[Use Cases]
        Ports[Port Interfaces]
    end

    subgraph "Infrastructure Layer - PostgreSQL Addition"
        PGRepo[PostgreSQL Repositories]
        PGSqlcGen[PostgreSQL sqlcgen]
        PGConnection[PostgreSQL Connection Factory]

        MySQLRepo[MySQL Repositories]
        MySQLSqlcGen[MySQL sqlcgen]

        SQLiteRepo[SQLite Repositories]
        SQLiteSqlcGen[SQLite sqlcgen]
    end

    subgraph "Configuration & DI"
        Config[Database Config]
        DI[DI Container]
    end

    subgraph "External Systems"
        PostgreSQLDB[(PostgreSQL 18.2)]
        MySQLDB[(MySQL 8.4)]
        SQLiteDB[(SQLite)]
    end

    CLI --> UseCase
    UseCase --> Ports
    DI --> Config
    DI --> PGRepo
    DI --> MySQLRepo
    DI --> SQLiteRepo

    PGRepo --> PGSqlcGen
    MySQLRepo --> MySQLSqlcGen
    SQLiteRepo --> SQLiteSqlcGen

    PGRepo -.implements.-> Ports
    MySQLRepo -.implements.-> Ports
    SQLiteRepo -.implements.-> Ports

    PGConnection --> PostgreSQLDB
    PGSqlcGen --> PostgreSQLDB
```

**Architecture Integration**:
- Selected pattern: Repository Pattern (maintains existing MySQL/SQLite approach)
- Domain/feature boundaries: Database implementation isolated to Infrastructure layer; Application layer remains database-agnostic
- Existing patterns preserved: Clean Architecture layers, Repository pattern, DI container selection
- New components rationale:
  - PostgreSQL connection factory: Manages pgx driver lifecycle and connection pooling
  - PostgreSQL sqlcgen package: Type-safe query execution (auto-generated by sqlc)
  - PostgreSQL repository implementations: Conform to existing repository interfaces
  - PostgreSQL Atlas environments: Schema migration tracking
- Steering compliance: Maintains Clean Architecture dependency direction, Zero domain dependencies on infrastructure

### Technology Stack

| Layer | Choice / Version | Role in Feature | Notes |
|-------|------------------|-----------------|-------|
| Data / Storage | PostgreSQL 18.2 | Primary database backend | Official Docker image `postgres:18.2` |
| Data / Storage | pgx v5 | PostgreSQL Go driver | `github.com/jackc/pgx/v5` - 50-100% faster than lib/pq, actively maintained |
| Data / Storage | sqlc 1.x | Type-safe SQL code generator | Generates PostgreSQL-specific Go code from queries |
| Infrastructure / Runtime | Atlas 1.1.0 | Database schema migration tool | Upgraded from 1.0 for improved PostgreSQL support |
| Infrastructure / Runtime | Docker Compose | Local development environment | PostgreSQL 18.2 container + migration services |

**Key Technology Decisions** (see `research.md` for detailed investigation):
- **pgx over lib/pq**: 50-100% performance improvement, active maintenance, built-in connection pooling ([source](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/))
- **TEXT+CHECK over ENUM**: Schema evolution flexibility, SHARE UPDATE EXCLUSIVE lock (allows concurrent updates) ([source](https://making.close.com/posts/native-enums-or-check-constraints-in-postgresql/))
- **Atlas 1.1**: Enhanced PostgreSQL support for partitions, functions, procedures ([Atlas releases](https://github.com/ariga/atlas/releases))

## System Flows

### Database Selection Flow

```mermaid
sequenceDiagram
    participant App as Application Start
    participant Config as Config Loader
    participant DI as DI Container
    participant PGConn as PostgreSQL Connection Factory
    participant PGRepo as PostgreSQL Repository
    participant PG as PostgreSQL Database

    App->>Config: Load wallet config
    Config->>Config: Parse database.type = "postgresql"
    Config->>Config: Validate PostgreSQL config section
    App->>DI: Initialize container
    DI->>Config: Read database.type
    DI->>PGConn: NewPostgreSQL(conf.Database.PostgreSQL)
    PGConn->>PG: Connect with pgx driver
    PG-->>PGConn: Connection established
    PGConn-->>DI: Return *sql.DB
    DI->>PGRepo: NewRepositorySqlc(dbConn, coinType)
    PGRepo-->>DI: Repository instance
    DI-->>App: Container ready
```

**Key Decisions**:
- Config validation occurs before connection attempt (fail-fast principle)
- DI container owns database selection logic via switch on database.type
- Connection pooling configured at factory level (pgx native pool support)
- Repository creation deferred until needed (lazy initialization)

### Schema Migration Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant HCL as Atlas HCL Schema
    participant Atlas as Atlas CLI
    participant PG as PostgreSQL Database

    Dev->>HCL: Modify schema (watch.hcl)
    Dev->>Atlas: atlas migrate diff --env local_postgresql_watch
    Atlas->>HCL: Read schema definition
    Atlas->>PG: Compare with current DB state
    PG-->>Atlas: Return current schema
    Atlas->>Atlas: Generate migration SQL
    Atlas-->>Dev: Create migration file
    Dev->>Atlas: atlas migrate apply --env local_postgresql_watch
    Atlas->>PG: Execute migration SQL
    PG-->>Atlas: Migration applied
    Atlas-->>Dev: Success confirmation
```

**Key Decisions**:
- HCL schemas remain database-agnostic (reuse existing watch.hcl, keygen.hcl, sign.hcl)
- Migrations are database-specific (separate directories per engine)
- Dev database (`docker://postgres/18`) validates migrations before application
- Migration failures prevent deployment (safety first)

## Requirements Traceability

| Requirement | Summary | Components | Interfaces | Flows |
|-------------|---------|------------|------------|-------|
| 1.1-1.7 | PostgreSQL configuration support and validation | Database Config struct, validateDatabase() | config.Database.PostgreSQL | Config validation |
| 2.1-2.5 | PostgreSQL schema generation and type conversion | Schema files in schemas_postgresql/ | N/A (file artifacts) | N/A |
| 3.1-3.5 | sqlc code generation for PostgreSQL | sqlc_postgresql.yml, PostgreSQL sqlcgen package | sqlcgen.Queries interface | N/A |
| 4.1-4.5 | PostgreSQL connection management | PostgreSQL Connection Factory (pkg/db/postgresql/) | NewPostgreSQL() function | Database selection flow |
| 5.1-5.5 | Build system integration | Makefile targets (db_sqlc.mk, db_atlas.mk) | make sqlc, make atlas-* | N/A |
| 6.1-6.5 | Migration path and data compatibility | Migration documentation, data type mappings | N/A (documentation) | Data migration process |
| 7.1-7.7 | Atlas PostgreSQL migration support | Atlas configuration (atlas.hcl), PostgreSQL environments | atlas migrate diff/apply | Schema migration flow |
| 8.1-8.8 | Docker Compose PostgreSQL integration | compose.yaml PostgreSQL service, migration services | docker-compose up | Container startup |
| 9.1-9.5 | Testing and validation | Integration test suite, PostgreSQL test database | Test helper functions | Test execution |
| 10.1-10.7 | Atlas version upgrade to 1.1 | Docker image update, Makefile updates | atlas CLI commands | Migration compatibility |

## Components and Interfaces

### Component Summary

| Component | Domain/Layer | Intent | Req Coverage | Key Dependencies (P0/P1) | Contracts |
|-----------|--------------|--------|--------------|--------------------------|-----------|
| PostgreSQLConfig | Config/Infrastructure | PostgreSQL connection configuration | 1.1-1.7 | None | State |
| PostgreSQLConnectionFactory | Infrastructure/DB | PostgreSQL connection lifecycle management | 4.1-4.5 | pgx driver (P0) | Service |
| PostgreSQLRepositories | Infrastructure/Repository | Database operations for PostgreSQL | 3.1-3.5, 4.1-4.5 | PostgreSQL sqlcgen (P0) | Service |
| PostgreSQLSqlcGen | Infrastructure/DB | Type-safe query execution (auto-generated) | 3.1-3.5 | pgx driver (P0) | Service |
| AtlasPostgreSQLEnv | Infrastructure/Migration | Schema migration for PostgreSQL | 7.1-7.7 | Atlas 1.1 (P0), PostgreSQL 18.2 (P0) | Service |
| PostgreSQLDockerService | Infrastructure/Container | Local PostgreSQL environment | 8.1-8.8 | Docker (P0) | Service |
| BuildSystemTargets | Build/Infrastructure | Code generation and validation | 5.1-5.5, 10.1-10.7 | sqlc (P0), Atlas (P0) | Batch |

### Configuration Layer

#### PostgreSQLConfig

| Field | Detail |
|-------|--------|
| Intent | Define PostgreSQL connection parameters and validation rules |
| Requirements | 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7 |

**Responsibilities & Constraints**:
- Store PostgreSQL-specific connection parameters (host, port, database, user, password, SSL mode)
- Enforce required field validation (host, dbname, user, pass)
- Support optional connection pool configuration (max connections, idle connections, connection lifetime)
- Domain boundary: Configuration layer only, no business logic

**Dependencies**:
- Inbound: Config loader from TOML/YAML files
- Outbound: PostgreSQL Connection Factory consumes this configuration
- External: None

**Contracts**: State [x]

##### State Management

**State Model**:
```go
// PostgreSQL connection configuration
type PostgreSQL struct {
    Host     string `toml:"host" yaml:"host" mapstructure:"host"`
    Port     int    `toml:"port" yaml:"port" mapstructure:"port"`
    DB       string `toml:"dbname" yaml:"dbname" mapstructure:"dbname"`
    User     string `toml:"user" yaml:"user" mapstructure:"user"`
    Pass     string `toml:"pass" yaml:"pass" mapstructure:"pass"`
    SSLMode  string `toml:"sslmode" yaml:"sslmode" mapstructure:"sslmode"`
    Debug    bool   `toml:"debug" yaml:"debug" mapstructure:"debug"`
}

// Database selection configuration (updated)
type Database struct {
    Type       string     `toml:"type" yaml:"type" mapstructure:"type" validate:"required,oneof=mysql sqlite postgresql"`
    MySQL      MySQL      `toml:"mysql" yaml:"mysql" mapstructure:"mysql"`
    SQLite     SQLite     `toml:"sqlite" yaml:"sqlite" mapstructure:"sqlite"`
    PostgreSQL PostgreSQL `toml:"postgresql" yaml:"postgresql" mapstructure:"postgresql"`
}
```

**Validation Rules** (in validateDatabase()):
```go
case "postgresql":
    if c.Database.PostgreSQL.Host == "" {
        return errors.New("database.postgresql.host is required when database.type is postgresql")
    }
    if c.Database.PostgreSQL.DB == "" {
        return errors.New("database.postgresql.dbname is required when database.type is postgresql")
    }
    if c.Database.PostgreSQL.User == "" {
        return errors.New("database.postgresql.user is required when database.type is postgresql")
    }
    if c.Database.PostgreSQL.Pass == "" {
        return errors.New("database.postgresql.pass is required when database.type is postgresql")
    }
```

**Persistence & Consistency**: Loaded from TOML configuration file; immutable after parsing

**Implementation Notes**:
- Integration: Extend existing Database struct in `pkg/config/wallet.go`
- Validation: Add PostgreSQL case to validateDatabase() switch statement
- Risks: None - follows existing MySQL/SQLite pattern exactly

### Infrastructure Layer - Database

#### PostgreSQLConnectionFactory

| Field | Detail |
|-------|--------|
| Intent | Establish and manage PostgreSQL database connections using pgx driver |
| Requirements | 4.1, 4.2, 4.3, 4.4, 4.5 |

**Responsibilities & Constraints**:
- Create PostgreSQL connections using pgx driver with proper SSL configuration
- Configure connection pool (max open connections, max idle connections, connection lifetime)
- Verify connection health via Ping before returning
- Return database/sql compatible *sql.DB instance for repository usage
- Transaction boundary: Per-connection (repository manages transactions)

**Dependencies**:
- Inbound: DI container requests connection during initialization
- Outbound: pgx driver (github.com/jackc/pgx/v5/stdlib) (P0 - Critical)
- External: PostgreSQL database server (P0 - Critical)

**Contracts**: Service [x]

##### Service Interface

```go
package postgresql

import (
    "database/sql"
    "fmt"

    _ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver

    "github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewPostgreSQL connects to PostgreSQL server and returns a database connection
// using the pgx driver for optimal performance
func NewPostgreSQL(conf *config.PostgreSQL) (*sql.DB, error)
```

**Preconditions**:
- conf must not be nil
- conf.Host, conf.DB, conf.User, conf.Pass must be non-empty (validated by config layer)
- PostgreSQL server must be accessible at conf.Host:conf.Port

**Postconditions**:
- Returns *sql.DB with established connection OR error
- Connection verified via Ping
- Connection pool configured appropriately for PostgreSQL

**Invariants**:
- Connection pool size > 0 (default: 10 max open connections)
- SSL mode defaults to "prefer" if not specified
- Connection lifetime defaults to unlimited if not specified

**Implementation Notes**:
- Integration: New package `pkg/db/postgresql/connection.go`
- Connection string format: `postgres://user:pass@host:port/dbname?sslmode=<mode>`
- Connection pool configuration follows PostgreSQL best practices (see research.md)
- Error handling: Wrap errors with descriptive context using fmt.Errorf
- Risks: Connection pool configuration may need tuning based on production load

#### PostgreSQLRepositories (Example: AddressRepository)

| Field | Detail |
|-------|--------|
| Intent | Implement repository interfaces using PostgreSQL sqlcgen-generated code |
| Requirements | 3.1, 3.2, 3.3, 3.4, 3.5, 4.1, 4.2 |

**Responsibilities & Constraints**:
- Implement repository interfaces defined in `internal/application/ports/repository/`
- Convert between domain entities and sqlcgen-generated types
- Execute type-safe queries via sqlcgen.Queries
- Handle database errors and map to domain errors
- Transaction boundary: Single repository operation (no multi-entity transactions)
- Data ownership: Repository owns database schema access for its entity

**Dependencies**:
- Inbound: Use cases invoke repository methods (P0 - Critical)
- Outbound: PostgreSQL sqlcgen.Queries interface (P0 - Critical), domain entity types (P0 - Critical)
- External: PostgreSQL database via *sql.DB (P0 - Critical)

**Contracts**: Service [x]

##### Service Interface

```go
package postgresql

import (
    "context"
    "database/sql"

    domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
    domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
    "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgresql/sqlcgen"
)

// AddressRepositorySqlc implements AddressRepositorier using PostgreSQL sqlcgen
type AddressRepositorySqlc struct {
    queries      *sqlcgen.Queries
    coinTypeCode domainCoin.CoinTypeCode
}

// NewAddressRepositorySqlc creates a new PostgreSQL address repository
func NewAddressRepositorySqlc(dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode) *AddressRepositorySqlc

// Repository interface methods (implement all from ports/repository/watch/address.go)
func (r *AddressRepositorySqlc) GetAll(ctx context.Context) ([]*domainAddress.Address, error)
func (r *AddressRepositorySqlc) GetByID(ctx context.Context, id int64) (*domainAddress.Address, error)
func (r *AddressRepositorySqlc) Insert(ctx context.Context, addr *domainAddress.Address) error
func (r *AddressRepositorySqlc) Update(ctx context.Context, addr *domainAddress.Address) error
```

**Type Conversion Functions** (private helpers):
```go
// convertToAddress converts sqlcgen.Address to domain.Address entity
func convertToAddress(sqlcAddr *sqlcgen.Address) (*domainAddress.Address, error)

// convertFromAddress converts domain.Address entity to sqlcgen.Address
func convertFromAddress(addr *domainAddress.Address) *sqlcgen.Address
```

**Preconditions**:
- dbConn must be non-nil and connected
- coinTypeCode must be valid (btc, bch, eth, xrp)
- Context must not be nil for all methods

**Postconditions**:
- Returns domain entities OR error (never partial results)
- Database state consistent after operations
- Errors include context for debugging

**Invariants**:
- Repository only operates on its assigned coin type
- All queries scoped by coin type
- Domain entities remain independent of database schema

**Implementation Notes**:
- Integration: Create ~40 repository files in `internal/infrastructure/repository/{cold,watch}/postgresql/`
- Pattern: Copy-paste from MySQL implementations, update import paths to PostgreSQL sqlcgen
- Validation: Repository validates domain entity invariants before database operations
- Risks: Type conversion errors if sqlcgen schema differs from domain model

### Infrastructure Layer - Schema Management

#### PostgreSQL Schema Files

| Field | Detail |
|-------|--------|
| Intent | Define PostgreSQL-compatible database schemas for watch, keygen, and sign databases |
| Requirements | 2.1, 2.2, 2.3, 2.4, 2.5 |

**Responsibilities & Constraints**:
- Store PostgreSQL DDL for all wallet tables
- Convert MySQL-specific types to PostgreSQL equivalents
- Preserve referential integrity constraints
- Domain boundary: Schema definition only, no data
- Three separate schema files: 01_watch.sql, 02_keygen.sql, 03_sign.sql

**Data Type Mappings**:

| MySQL Type | PostgreSQL Type | Conversion Notes |
|------------|-----------------|------------------|
| AUTO_INCREMENT | SERIAL / BIGSERIAL | Use BIGSERIAL for BIGINT columns |
| tinyint(1) | BOOLEAN | PostgreSQL native boolean |
| enum('a','b') | TEXT CHECK(col IN ('a','b')) | CHECK constraint for validation |
| decimal(26,10) | NUMERIC(26,10) | Identical precision |
| datetime | TIMESTAMP | Default without time zone |
| varchar(n) | VARCHAR(n) | Direct mapping |
| text | TEXT | Direct mapping |
| bigint | BIGINT | Direct mapping |

**Enum Conversion Example**:
```sql
-- MySQL
coin enum('btc','bch','eth','xrp','hyt') NOT NULL COMMENT 'coin type code'

-- PostgreSQL
coin TEXT NOT NULL CHECK (coin IN ('btc','bch','eth','xrp','hyt')) 
-- COMMENT ON COLUMN table.coin IS 'coin type code';
```

**Implementation Notes**:
- Location: `tools/sqlc/schemas_postgresql/01_watch.sql`, `02_keygen.sql`, `03_sign.sql`
- Conversion method: Manual conversion with validation (see research.md for rationale)
- Comments: Use PostgreSQL COMMENT syntax for column/table descriptions
- Risks: Type precision differences require thorough testing (especially NUMERIC for crypto amounts)

#### sqlc PostgreSQL Configuration

| Field | Detail |
|-------|--------|
| Intent | Configure sqlc to generate PostgreSQL-specific Go code |
| Requirements | 3.1, 3.2, 3.3, 3.4, 3.5 |

**Configuration File** (`tools/sqlc/sqlc_postgresql.yml`):
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./queries/*.sql"
    schema: "./schemas_postgresql/*.sql"
    gen:
      go:
        package: "sqlcgen"
        out: "../../internal/infrastructure/database/postgresql/sqlcgen"
```

**Generated Code Structure**:
- `internal/infrastructure/database/postgresql/sqlcgen/`
  - `db.go` - Database interface wrapper
  - `models.go` - Go structs matching database schema
  - `*.sql.go` - Query implementations (one file per query file)

**Implementation Notes**:
- Shared queries: Uses same `tools/sqlc/queries/*.sql` files as MySQL/SQLite
- Query validation: sqlc validates PostgreSQL syntax and types at generation time
- Risks: If queries contain database-specific syntax, sqlc will fail compilation (intentional safety)

#### Atlas PostgreSQL Environments

| Field | Detail |
|-------|--------|
| Intent | Define Atlas migration environments for PostgreSQL schemas |
| Requirements | 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 10.1, 10.2 |

**Atlas Configuration Extension** (`tools/atlas/atlas.hcl`):
```hcl
# PostgreSQL Local Environments
env "local_postgresql_watch" {
  url     = "postgres://postgres:postgres@localhost:5432/watch?sslmode=disable"
  src     = "file://schemas/watch.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgresql_watch"
  }
  dev = "docker://postgres/18/watch"
}

env "local_postgresql_keygen" {
  url     = "postgres://postgres:postgres@localhost:5432/keygen?sslmode=disable"
  src     = "file://schemas/keygen.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgresql_keygen"
  }
  dev = "docker://postgres/18/keygen"
}

env "local_postgresql_sign" {
  url     = "postgres://postgres:postgres@localhost:5432/sign?sslmode=disable"
  src     = "file://schemas/sign.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgresql_sign"
  }
  dev = "docker://postgres/18/sign"
}

# PostgreSQL Admin Environments (for schema operations)
env "admin_postgresql_watch" {
  url     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
  src     = "file://schemas/watch.hcl"
  schemas = ["public"]
}
# ... similar for keygen and sign
```

**Implementation Notes**:
- HCL schemas: Reuse existing `schemas/watch.hcl`, `keygen.hcl`, `sign.hcl` (database-agnostic)
- Migrations: Separate directories per database (`migrations/postgresql_watch/`, etc.)
- Dev database: Uses PostgreSQL 18 Docker container for migration validation
- Admin environments: Required for schema-level operations (clean, drop)
- Risks: None - follows MySQL pattern exactly

### Infrastructure Layer - Docker & Build

#### PostgreSQL Docker Service

| Field | Detail |
|-------|--------|
| Intent | Provide PostgreSQL 18.2 local development environment via Docker Compose |
| Requirements | 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7, 8.8 |

**Docker Compose Service Definition** (`compose.yaml`):
```yaml
services:
  # PostgreSQL Database Service
  wallet-db-postgresql:
    image: postgres:18.2
    container_name: wallet-db-postgresql
    volumes:
      - wallet-db-postgresql:/var/lib/postgresql/data
      - "./docker/postgresql/init.d:/docker-entrypoint-initdb.d"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      # Databases created by init script
    ports:
      - "${POSTGRESQL_PORT:-5432}:5432"
    networks:
      - db
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  # PostgreSQL Migration Services (Atlas 1.1.0)
  wallet-db-migrate-postgresql-watch:
    image: arigaio/atlas:1.1.0
    container_name: wallet-db-migrate-postgresql-watch
    volumes:
      - "./tools/atlas:/app/atlas"
    working_dir: /app/atlas
    depends_on:
      wallet-db-postgresql:
        condition: service_healthy
    networks:
      - db
    command:
      - migrate
      - apply
      - --dir
      - file://migrations/postgresql_watch
      - --url
      - "postgres://postgres:postgres@wallet-db-postgresql:5432/watch?sslmode=disable"
    restart: "no"

  # Similar services for keygen and sign...
```

**Database Initialization Script** (`docker/postgresql/init.d/01_create_databases.sh`):
```bash
#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE watch;
    CREATE DATABASE keygen;
    CREATE DATABASE sign;
    GRANT ALL PRIVILEGES ON DATABASE watch TO postgres;
    GRANT ALL PRIVILEGES ON DATABASE keygen TO postgres;
    GRANT ALL PRIVILEGES ON DATABASE sign TO postgres;
EOSQL
```

**Volumes**:
```yaml
volumes:
  wallet-db-postgresql: {}
```

**Implementation Notes**:
- PostgreSQL version: 18.2 (official postgres image)
- Coexistence: PostgreSQL and MySQL containers run simultaneously for migration testing
- Healthcheck: Uses pg_isready for readiness verification
- Atlas upgrade: All migration services use arigaio/atlas:1.1.0 (including MySQL/SQLite)
- Risks: Port conflict if 5432 already in use (use POSTGRESQL_PORT env var)

#### Build System Integration

| Field | Detail |
|-------|--------|
| Intent | Extend Makefile targets to include PostgreSQL code generation and validation |
| Requirements | 5.1, 5.2, 5.3, 5.4, 5.5, 10.1, 10.4, 10.6 |

**Makefile Extensions** (`make/db_sqlc.mk`):
```makefile
# SQLC compilation for PostgreSQL
.PHONY: sqlc-compile
sqlc-compile:
	@echo "Compiling MySQL SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile
	@echo "✓ MySQL SQL compilation successful"
	@echo "Compiling SQLite SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL compilation successful"
	@echo "Compiling PostgreSQL SQL queries and schemas..."
	@cd tools/sqlc && sqlc compile -f sqlc_postgresql.yml
	@echo "✓ PostgreSQL SQL compilation successful"

# SQLC vet for PostgreSQL
.PHONY: sqlc-vet
sqlc-vet:
	@echo "Vetting MySQL SQL queries..."
	@cd tools/sqlc && sqlc vet
	@echo "✓ MySQL SQL queries passed vetting"
	@echo "Vetting SQLite SQL queries..."
	@cd tools/sqlc && sqlc vet -f sqlc_sqlite.yml
	@echo "✓ SQLite SQL queries passed vetting"
	@echo "Vetting PostgreSQL SQL queries..."
	@cd tools/sqlc && sqlc vet -f sqlc_postgresql.yml
	@echo "✓ PostgreSQL SQL queries passed vetting"

# Generate SQLC code for all databases
.PHONY: sqlc
sqlc:
	@echo "Generating MySQL sqlc code..."
	@cd tools/sqlc && sqlc generate
	@echo "✓ MySQL sqlc code generated"
	@echo "Generating SQLite sqlc code..."
	@cd tools/sqlc && sqlc generate -f sqlc_sqlite.yml
	@echo "✓ SQLite sqlc code generated"
	@echo "Generating PostgreSQL sqlc code..."
	@cd tools/sqlc && sqlc generate -f sqlc_postgresql.yml
	@echo "✓ PostgreSQL sqlc code generated"
```

**Makefile Extensions** (`make/db_atlas.mk`):
```makefile
# Atlas PostgreSQL environments
ATLAS_POSTGRESQL_SCHEMAS := postgresql_watch postgresql_keygen postgresql_sign

# Lint PostgreSQL schemas
.PHONY: atlas-lint-postgresql
atlas-lint-postgresql:
	@echo "Linting PostgreSQL Atlas HCL schema files..."
	@for schema in $(ATLAS_POSTGRESQL_SCHEMAS); do \
		echo "=== Linting $$schema schema ==="; \
		(cd tools/atlas && atlas schema lint --config file://atlas.hcl --env local_$$schema) || exit 1; \
	done
	@echo "✓ All PostgreSQL schemas passed linting"

# Combined lint (all databases)
.PHONY: atlas-lint
atlas-lint: atlas-lint-mysql atlas-lint-postgresql
	@echo "✓ All database schemas passed linting"
```

**Contracts**: Batch [x]

**Batch Contract**:
- Trigger: Manual execution via `make sqlc` or `make atlas-lint`
- Input: Schema files (`schemas_postgresql/*.sql`) and HCL files (`schemas/*.hcl`)
- Output: Generated Go code or validation results
- Idempotency: Code generation is deterministic (same input → same output)
- Recovery: Regenerate code on failure; validation errors prevent commits

**Implementation Notes**:
- Atlas version update: Change all `arigaio/atlas:1.0.0` references to `1.1.0` in compose.yaml and Makefile comments
- Integration: Extend existing targets rather than create new top-level targets
- Validation: All three databases must pass before allowing commit
- Risks: Build time increases linearly with number of databases (acceptable trade-off)

### Application Layer - Dependency Injection

#### DI Container Updates

| Field | Detail |
|-------|--------|
| Intent | Add PostgreSQL repository selection to DI container |
| Requirements | 4.1, 4.2 |

**DI Container Modifications** (`internal/di/container.go`):
```go
// Example: Address Repository factory
func (c *container) newAddressRepo() repowatch.AddressRepositorier {
    switch c.conf.Database.Type {
    case "mysql":
        return watchmysql.NewAddressRepositorySqlc(
            c.pkgContainer.NewDatabaseClient(),
            c.conf.CoinTypeCode,
        )
    case "sqlite":
        return watchsqlite.NewAddressRepositorySqlc(
            c.pkgContainer.NewSQLiteClient(),
            c.conf.CoinTypeCode,
        )
    case "postgresql":
        return watchpostgresql.NewAddressRepositorySqlc(
            c.pkgContainer.NewPostgreSQLClient(),
            c.conf.CoinTypeCode,
        )
    default:
        panic("unsupported database type: " + c.conf.Database.Type)
    }
}
```

**pkg Container Extension** (`pkg/di/container.go` or similar):
```go
// NewPostgreSQLClient creates PostgreSQL database connection
func (c *pkgContainer) NewPostgreSQLClient() *sql.DB {
    if c.postgreSQLClient == nil {
        var err error
        c.postgreSQLClient, err = postgresql.NewPostgreSQL(&c.conf.Database.PostgreSQL)
        if err != nil {
            panic(fmt.Sprintf("failed to connect to PostgreSQL: %v", err))
        }
    }
    return c.postgreSQLClient
}
```

**Implementation Notes**:
- Modification count: ~20 repository factory methods require PostgreSQL case addition
- Pattern: Identical structure to MySQL/SQLite cases
- Risks: Easy to miss a factory method during implementation (mitigate with grep search)

## Data Models

### Physical Data Model (PostgreSQL Schemas)

**Structure**: Three separate databases (watch, keygen, sign) with PostgreSQL-specific schema syntax.

**Watch Schema** (Example Table):
```sql
-- Table: address (watch database)
CREATE TABLE address (
    id BIGSERIAL PRIMARY KEY,
    coin TEXT NOT NULL CHECK (coin IN ('btc','bch','eth','xrp','hyt')),
    account TEXT NOT NULL CHECK (account IN ('client','deposit','payment','stored')),
    wallet_address VARCHAR(255) NOT NULL UNIQUE,
    is_allocated BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_address_coin ON address(coin);
CREATE INDEX idx_address_account ON address(account);

COMMENT ON TABLE address IS 'table for account pubkey';
COMMENT ON COLUMN address.coin IS 'coin type code';
COMMENT ON COLUMN address.account IS 'account type';
COMMENT ON COLUMN address.is_allocated IS 'true: address is allocated(used)';
```

**Keygen Schema** (Example Table):
```sql
-- Table: btc_account_key (keygen database)
CREATE TABLE btc_account_key (
    id BIGSERIAL PRIMARY KEY,
    coin TEXT NOT NULL CHECK (coin IN ('btc','bch')),
    account TEXT NOT NULL CHECK (account IN ('client','deposit','payment','stored')),
    idx INTEGER NOT NULL,
    address VARCHAR(255) NOT NULL,
    wif VARCHAR(255) NOT NULL,
    full_pubkey_idx INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_btc_account_key_coin ON btc_account_key(coin);
CREATE INDEX idx_btc_account_key_account ON btc_account_key(account);
CREATE UNIQUE INDEX idx_btc_account_key_coin_account_idx ON btc_account_key(coin, account, idx);

COMMENT ON TABLE btc_account_key IS 'table for btc account pubkey';
```

**Sign Schema** (Similar structure to keygen for offline signing)

**Key Constraints**:
- Primary keys: BIGSERIAL (auto-incrementing 64-bit integers)
- Foreign keys: Preserved from MySQL schema
- Unique constraints: Preserved from MySQL schema
- Check constraints: Replace MySQL ENUMs
- Indexes: Preserved from MySQL schema for query performance

### Data Contracts & Integration

**Type Conversion Interface** (sqlcgen to domain):
```go
// PostgreSQL sqlcgen types (auto-generated)
type Address struct {
    ID            int64
    Coin          string // TEXT column with CHECK constraint
    Account       string
    WalletAddress string
    IsAllocated   bool
    UpdatedAt     sql.NullTime
}

// Domain entity (clean, infrastructure-independent)
type Address struct {
    ID            int64
    CoinTypeCode  CoinTypeCode // domain enum
    AccountType   AccountType  // domain enum
    WalletAddress string
    IsAllocated   bool
    UpdatedAt     *time.Time
}

// Conversion function
func convertToAddress(sqlcAddr *sqlcgen.Address) (*domainAddress.Address, error) {
    addr := &domainAddress.Address{
        ID:            sqlcAddr.ID,
        CoinTypeCode:  domainCoin.CoinTypeCode(sqlcAddr.Coin),
        AccountType:   domainAccount.AccountType(sqlcAddr.Account),
        WalletAddress: sqlcAddr.WalletAddress,
        IsAllocated:   sqlcAddr.IsAllocated,
    }
    if sqlcAddr.UpdatedAt.Valid {
        addr.UpdatedAt = &sqlcAddr.UpdatedAt.Time
    }
    return addr, nil
}
```

**Cross-Database Data Management**:
- No distributed transactions (each wallet database is independent)
- Data synchronization: File-based export/import for offline wallets (keygen/sign)
- Eventual consistency: Not applicable (no distributed state)
- Migration strategy: Database-specific dump/restore tools (pg_dump, mysqldump)

## Error Handling

### Error Strategy

PostgreSQL-specific error handling follows the existing repository error pattern:

**Error Categories**:
1. **Connection Errors** (Infrastructure): Network failures, authentication failures, SSL handshake errors
2. **Query Errors** (Application): Syntax errors, type mismatches, constraint violations
3. **Transaction Errors** (Application): Deadlocks, serialization failures, timeout errors
4. **Data Errors** (Domain): Validation failures, business rule violations

**Error Responses**:
- **Connection Errors**: Fail fast with descriptive error; log connection parameters (excluding password); suggest network/configuration checks
- **Query Errors**: Wrap with context (query name, parameters); sqlc provides compile-time safety for most syntax errors
- **Transaction Errors**: Retry logic for transient errors (deadlocks); abort for fatal errors
- **Data Errors**: Return domain-specific errors; do not expose database details to callers

### Error Categories and Responses

**PostgreSQL-Specific Errors**:
```go
// Connection establishment errors
if err := db.Ping(); err != nil {
    return nil, fmt.Errorf("failed to ping PostgreSQL database at %s:%d: %w", 
        conf.Host, conf.Port, err)
}

// Query execution errors (repository layer)
rows, err := r.queries.GetAllAddresses(ctx, sqlcgen.GetAllAddressesParams{...})
if err != nil {
    return nil, fmt.Errorf("failed to execute GetAllAddresses query: %w", err)
}

// Constraint violation errors (CHECK constraint on enum-like fields)
if err := r.queries.InsertAddress(ctx, params); err != nil {
    if strings.Contains(err.Error(), "violates check constraint") {
        return fmt.Errorf("invalid coin type value: %w", err)
    }
    return fmt.Errorf("failed to insert address: %w", err)
}
```

**Monitoring**:
- Log all connection failures with host/port/database information (exclude credentials)
- Log query execution failures with query name and duration
- Metrics: Connection pool utilization, query latency percentiles, error rates by category
- Health check: Periodic Ping verification in monitoring service

## Testing Strategy

### Unit Tests

**Connection Factory Tests** (`pkg/db/postgresql/connection_test.go`):
- Test successful connection with valid configuration
- Test connection failure with invalid host
- Test connection failure with invalid credentials
- Test connection pool configuration (max connections, idle connections)
- Test SSL mode variations (disable, require, verify-ca, verify-full)

**Repository Tests** (example: `internal/infrastructure/repository/watch/postgresql/address_sqlc_test.go`):
- Test type conversion (sqlcgen.Address ↔ domain.Address)
- Test CRUD operations (Insert, GetByID, GetAll, Update)
- Test query filtering (by coin type, account type)
- Test NULL handling (optional fields like UpdatedAt)
- Test constraint violations (invalid enum values in CHECK constraints)

### Integration Tests

**Database Operations** (`internal/integration_test/postgresql/`):
- Test end-to-end flow: config load → connection → repository → query execution
- Test transaction handling (commit, rollback)
- Test concurrent access (multiple connections)
- Test all CRUD operations across all repository types
- Test data consistency: Insert via PostgreSQL repo, read via domain query

**Schema Migration Tests**:
- Test Atlas migration application (`atlas migrate apply`)
- Test migration rollback scenarios
- Test schema validation (`atlas schema diff` shows no differences after migration)
- Test idempotency (apply migrations multiple times)

**Cross-Database Compatibility Tests**:
- Test identical query results across MySQL, SQLite, PostgreSQL
- Test data type precision (especially NUMERIC for cryptocurrency amounts)
- Test enum-like field validation (CHECK constraints work like MySQL ENUMs)

### Performance Tests

**Connection Pool Tests**:
- Test connection acquisition under load (100 concurrent requests)
- Test connection pool exhaustion handling
- Test connection lifetime and idle timeout behavior

**Query Performance Tests** (baseline, not optimization):
- Measure query latency for common operations (GetAll, GetByID, Insert)
- Compare performance across MySQL, SQLite, PostgreSQL (informational only)
- Establish baseline metrics for regression detection

## Migration Strategy

### Phase 1: Schema Creation (Manual Conversion)

```mermaid
flowchart TB
    Start[Start: MySQL Schemas] --> ReadMySQL[Read MySQL Schema Files]
    ReadMySQL --> ConvertTypes[Convert Data Types]
    ConvertTypes --> ConvertEnums[Convert ENUMs to CHECK]
    ConvertEnums --> ConvertAutoInc[Convert AUTO_INCREMENT to SERIAL]
    ConvertAutoInc --> AddComments[Add PostgreSQL COMMENT syntax]
    AddComments --> Validate[Validate Syntax]
    Validate --> WritePostgreSQL[Write PostgreSQL Schema Files]
    WritePostgreSQL --> Review[Manual Review]
    Review --> |Issues Found| ConvertTypes
    Review --> |Valid| Complete[Complete: PostgreSQL Schemas]
```

**Conversion Checklist**:
- [ ] Convert all AUTO_INCREMENT columns to SERIAL/BIGSERIAL
- [ ] Convert all tinyint(1) columns to BOOLEAN
- [ ] Convert all ENUM columns to TEXT with CHECK constraints
- [ ] Convert all datetime columns to TIMESTAMP
- [ ] Verify DECIMAL precision matches (26,10)
- [ ] Convert MySQL comments to PostgreSQL COMMENT syntax
- [ ] Remove ENGINE=InnoDB and MySQL-specific syntax
- [ ] Validate schema with sqlc compile
- [ ] Test CREATE TABLE statements in PostgreSQL 18.2

**Validation**: `cd tools/sqlc && sqlc compile -f sqlc_postgresql.yml`

### Phase 2: Data Migration (MySQL → PostgreSQL)

**Migration Tools**:
- **pgloader** (recommended): Automated migration tool with type conversion
- **Manual dump/restore**: For custom data transformations

**pgloader Configuration** (example for watch database):
```lisp
LOAD DATABASE
    FROM mysql://hiromaily:hiromaily@localhost:3306/watch
    INTO postgres://postgres:postgres@localhost:5432/watch

WITH include drop, create tables, create indexes, reset sequences,
     workers = 8, concurrency = 1

SET MySQL PARAMETERS
    net_read_timeout = '120',
    net_write_timeout = '120'

SET PostgreSQL PARAMETERS
    maintenance_work_mem to '128MB',
    work_mem to '12MB'

CAST type tinyint to boolean
     drop typemod,
CAST type decimal(26,10) to numeric(26,10)
     keep typemod;
```

**Migration Process**:
1. Backup source database (mysqldump or pg_dump)
2. Create PostgreSQL schemas (Atlas migration)
3. Run pgloader for data transfer
4. Validate data integrity (row counts, sum checks for amounts)
5. Test wallet operations with migrated data
6. Rollback triggers: Data loss, precision errors, constraint violations

**Rollback Strategy**:
- Keep original database untouched during migration
- Test migration on copy first
- Document rollback procedure (restore from backup)

## Security Considerations

**PostgreSQL Connection Security**:
- SSL mode configuration: Support disable, require, verify-ca, verify-full
- Password security: Never log passwords; use environment variables or secrets management
- Connection string sanitization: Redact passwords in logs and errors

**Authentication**:
- Use role-based access control in PostgreSQL (separate users for watch/keygen/sign)
- Minimum privilege principle: Grant only necessary permissions to application user
- Credential rotation: Document procedure for updating database passwords

**Data Protection**:
- Encryption at rest: PostgreSQL supports tablespace encryption (deployment decision)
- Encryption in transit: Enforce SSL for production deployments
- Sensitive data: Never log private keys, seeds, or passwords (existing rule maintained)

**Audit Logging**:
- PostgreSQL audit logs for DDL operations (schema changes)
- Application logs for data access patterns
- Compliance: Maintain audit trail for regulatory requirements

## Performance & Scalability

**Target Metrics** (informational baseline, not optimization targets):
- Connection pool utilization: < 80% under normal load
- Query latency: P95 < 100ms for single-record queries, P95 < 500ms for batch operations
- Connection acquisition: < 10ms under normal load
- Transaction throughput: > 100 TPS for wallet operations

**Connection Pool Configuration**:
```go
// PostgreSQL connection pool (pkg/db/postgresql/connection.go)
db.SetMaxOpenConns(25)           // Maximum concurrent connections
db.SetMaxIdleConns(5)            // Idle connections for quick reuse
db.SetConnMaxLifetime(5 * time.Minute) // Prevent stale connections
db.SetConnMaxIdleTime(1 * time.Minute) // Close unused connections
```

**Indexing Strategy**:
- Preserve all MySQL indexes for PostgreSQL schemas
- Add coin type and account type indexes for filtering queries
- Consider composite indexes for multi-column filters (coin + account)
- Monitor query plans with EXPLAIN ANALYZE during integration testing

**Scaling Approaches**:
- Vertical scaling: Increase PostgreSQL server resources (CPU, memory, storage)
- Read replicas: PostgreSQL replication for read-heavy workloads (future consideration)
- Connection pooling: External pgbouncer for high-concurrency scenarios (if needed)
- Partitioning: Table partitioning for large transaction history tables (future consideration)

**Caching Strategy**:
- Application-level caching: Not in scope for this feature
- PostgreSQL query cache: Automatic (PostgreSQL handles internally)
- Connection reuse: Connection pool provides implicit caching

