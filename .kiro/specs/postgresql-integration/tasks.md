# Implementation Plan

## Overview

This implementation plan adds PostgreSQL 18.2 as a third database backend alongside MySQL and SQLite, following the existing repository pattern. Tasks are organized to build incrementally, starting with configuration and schema setup, then code generation, connection management, repository implementations, build integration, and finally testing and documentation.

**Total Estimated Effort**: 5-7 days (based on gap analysis)

---

## Tasks

### Phase 1: Foundation & Configuration

- [ ] 1. Add PostgreSQL configuration support
- [ ] 1.1 (P) Extend database configuration struct for PostgreSQL
  - Add PostgreSQL struct to pkg/config/wallet.go with fields: Host, Port, DB, User, Pass, SSLMode, Debug
  - Update Database struct to include PostgreSQL field
  - Update validation tag from `oneof=mysql sqlite` to `oneof=mysql sqlite postgresql`
  - _Requirements: 1.1, 1.6_

- [ ] 1.2 (P) Implement PostgreSQL configuration validation
  - Add PostgreSQL case to validateDatabase() switch statement
  - Validate required fields: Host, DB, User, Pass (return descriptive errors if empty)
  - Support optional fields: Port (default 5432), SSLMode (default "prefer"), Debug
  - _Requirements: 1.2, 1.3, 1.4, 1.5, 1.7_

- [ ] 2. Create PostgreSQL schema files
- [ ] 2.1 Convert watch schema from MySQL to PostgreSQL
  - Convert tools/sqlc/schemas/01_watch.sql to PostgreSQL syntax
  - Map AUTO_INCREMENT to BIGSERIAL for primary keys
  - Convert tinyint(1) to BOOLEAN for is_allocated fields
  - Convert ENUM types to TEXT with CHECK constraints (coin, account types)
  - Convert datetime to TIMESTAMP, decimal(26,10) to NUMERIC(26,10)
  - Add PostgreSQL COMMENT syntax for table and column descriptions
  - Save as tools/sqlc/schemas_postgresql/01_watch.sql
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [ ] 2.2 Convert keygen schema from MySQL to PostgreSQL
  - Convert tools/sqlc/schemas/02_keygen.sql to PostgreSQL syntax
  - Apply same data type mappings as watch schema
  - Preserve all foreign key constraints and indexes
  - Save as tools/sqlc/schemas_postgresql/02_keygen.sql
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [ ] 2.3 Convert sign schema from MySQL to PostgreSQL
  - Convert tools/sqlc/schemas/03_sign.sql to PostgreSQL syntax
  - Apply same data type mappings as watch and keygen schemas
  - Validate CHECK constraints match MySQL ENUM values exactly
  - Save as tools/sqlc/schemas_postgresql/03_sign.sql
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [ ] 2.4 Validate PostgreSQL schema syntax
  - Create sqlc_postgresql.yml configuration file with engine: postgresql
  - Point schema to ./schemas_postgresql/*.sql
  - Run sqlc compile -f sqlc_postgresql.yml to verify syntax
  - Fix any PostgreSQL-specific syntax errors reported by sqlc
  - _Requirements: 2.5, 3.2_

### Phase 2: Code Generation & Connection

- [ ] 3. Configure sqlc for PostgreSQL code generation
- [ ] 3.1 (P) Create sqlc PostgreSQL configuration
  - Create tools/sqlc/sqlc_postgresql.yml with version 2
  - Configure engine: postgresql, queries: ./queries/*.sql
  - Set schema: ./schemas_postgresql/*.sql
  - Configure Go output to internal/infrastructure/database/postgresql/sqlcgen
  - _Requirements: 3.1, 3.2, 3.4_

- [ ] 3.2 (P) Generate PostgreSQL sqlcgen package
  - Run sqlc generate -f sqlc_postgresql.yml
  - Verify generated files: db.go, models.go, *.sql.go
  - Confirm PostgreSQL-specific types (BOOLEAN, NUMERIC, TEXT with CHECK)
  - Validate query interfaces match MySQL/SQLite implementations
  - _Requirements: 3.1, 3.3, 3.5_

- [ ] 4. Implement PostgreSQL database connection
- [ ] 4.1 (P) Create PostgreSQL connection factory
  - Create pkg/db/postgresql/connection.go package
  - Implement NewPostgreSQL(conf *config.PostgreSQL) (*sql.DB, error)
  - Use pgx driver: import _ "github.com/jackc/pgx/v5/stdlib"
  - Build connection string: postgres://user:pass@host:port/dbname?sslmode=mode
  - Configure connection pool: SetMaxOpenConns(25), SetMaxIdleConns(5)
  - Set connection lifetime: SetConnMaxLifetime(5min), SetConnMaxIdleTime(1min)
  - Verify connection with Ping() before returning
  - Wrap errors with descriptive context (host, port, database)
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 4.2 (P) Add PostgreSQL client to DI container
  - Add NewPostgreSQLClient() method to pkg DI container
  - Initialize connection using postgresql.NewPostgreSQL()
  - Cache connection in container (singleton pattern)
  - Panic with descriptive error if connection fails during initialization
  - _Requirements: 4.1, 4.2_

### Phase 3: Repository Implementations

- [ ] 5. Implement PostgreSQL repositories for watch wallet
- [ ] 5.1 (P) Create address repository for PostgreSQL
  - Create internal/infrastructure/repository/watch/postgresql/address_sqlc.go
  - Implement AddressRepositorySqlc struct with queries field
  - Add NewAddressRepositorySqlc(dbConn, coinTypeCode) constructor
  - Implement converter functions: convertToAddress, convertFromAddress
  - Implement repository methods: GetAll, GetByID, Insert, Update (match interface)
  - Handle NULL values in UpdatedAt field using sql.NullTime
  - _Requirements: 3.3, 4.1_

- [ ] 5.2 (P) Create BTC transaction repositories for PostgreSQL
  - Create btc_tx_sqlc.go, btc_tx_input_sqlc.go, btc_tx_output_sqlc.go
  - Implement BTCTxRepositorySqlc, TxInputRepositorySqlc, TxOutputRepositorySqlc
  - Add type conversion functions for BTC-specific entities
  - Implement all repository interface methods defined in ports
  - Handle decimal precision for cryptocurrency amounts (NUMERIC 26,10)
  - _Requirements: 3.3, 4.1_

- [ ] 5.3 (P) Create ETH and XRP transaction repositories for PostgreSQL
  - Create eth_detail_tx_sqlc.go, xrp_detail_tx_sqlc.go, xrp_pending_multisig_sqlc.go
  - Create xrp_multisig_signature_sqlc.go for XRP multisig support
  - Implement repository interfaces for ETH and XRP transactions
  - Add conversion functions for chain-specific entity types
  - _Requirements: 3.3, 4.1_

- [ ] 5.4 (P) Create payment and transaction repositories for PostgreSQL
  - Create payment_request_sqlc.go, tx_sqlc.go
  - Implement PaymentRequestRepositorySqlc and TxRepositorySqlc
  - Handle transaction status enums with CHECK constraints
  - _Requirements: 3.3, 4.1_

- [ ] 6. Implement PostgreSQL repositories for keygen/sign wallets
- [ ] 6.1 (P) Create account key repositories for PostgreSQL
  - Create internal/infrastructure/repository/cold/postgresql/ directory
  - Create btc_account_key_sqlc.go, eth_account_key_sqlc.go, xrp_account_key_sqlc.go
  - Create auth_account_key_sqlc.go, auth_fullpubkey_sqlc.go
  - Implement all account key repository interfaces
  - Handle private key storage types (WIF for BTC, keystore for ETH)
  - _Requirements: 3.3, 4.1_

- [ ] 6.2 (P) Create seed and nonce repositories for PostgreSQL
  - Create seed_sqlc.go for mnemonic seed storage
  - Create nonce_repository_sqlc.go for MuSig2 nonce management
  - Implement repository interfaces for seed and nonce operations
  - _Requirements: 3.3, 4.1_

- [ ] 6.3 (P) Create XRP signer list repositories for PostgreSQL
  - Create xrp_signer_list_sqlc.go, xrp_signer_entry_sqlc.go
  - Create xrp_regular_key_sqlc.go for XRP key management
  - Implement XRP multisig signer management repositories
  - _Requirements: 3.3, 4.1_

- [ ] 7. Integrate PostgreSQL repositories into DI container
- [ ] 7.1 Add PostgreSQL cases to watch repository factories
  - Update internal/di/container.go repository factory methods
  - Add "postgresql" case to each repository switch statement (~12 watch repos)
  - Return watchpostgresql.NewXxxRepositorySqlc(pgClient, coinType)
  - Verify all switch statements handle "mysql", "sqlite", "postgresql" consistently
  - _Requirements: 4.1, 4.2_

- [ ] 7.2 Add PostgreSQL cases to cold repository factories
  - Add "postgresql" case to keygen/sign repository factories (~8 cold repos)
  - Return coldpostgresql.NewXxxRepositorySqlc(pgClient, coinType)
  - Ensure default case panics with "unsupported database type" error
  - _Requirements: 4.1, 4.2_

### Phase 4: Build System & Atlas Integration

- [ ] 8. Extend Makefile targets for PostgreSQL
- [ ] 8.1 (P) Add PostgreSQL to sqlc build targets
  - Update make/db_sqlc.mk: sqlc-compile target to include PostgreSQL
  - Add sqlc compile -f sqlc_postgresql.yml with success message
  - Update sqlc-vet target to validate PostgreSQL queries
  - Update sqlc (generate) target to create PostgreSQL Go code
  - Ensure all three databases processed: MySQL, SQLite, PostgreSQL
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ] 8.2 (P) Update Atlas version to 1.1.0
  - Update Docker Compose image: arigaio/atlas:1.0.0 → arigaio/atlas:1.1.0
  - Update all migration services (watch, keygen, sign) to use Atlas 1.1
  - Update make/db_atlas.mk comments referencing Atlas version
  - Test existing MySQL migrations with Atlas 1.1 (verify compatibility)
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.6, 10.7_

- [ ] 9. Configure Atlas for PostgreSQL migrations
- [ ] 9.1 (P) Add PostgreSQL environments to Atlas configuration
  - Edit tools/atlas/atlas.hcl to add local_postgresql_watch environment
  - Set URL: postgres://postgres:postgres@localhost:5432/watch?sslmode=disable
  - Configure src: file://schemas/watch.hcl (reuse existing HCL)
  - Set migration dir: file://migrations/postgresql_watch
  - Configure dev database: docker://postgres/18/watch
  - Add similar environments for keygen and sign databases
  - _Requirements: 7.1, 7.2, 7.3, 7.5, 7.7_

- [ ] 9.2 (P) Add PostgreSQL admin environments to Atlas configuration
  - Add admin_postgresql_watch environment without database in URL
  - Configure for schema-level operations (clean, drop)
  - Add admin environments for keygen and sign
  - _Requirements: 7.1, 7.7_

- [ ] 9.3 (P) Generate initial PostgreSQL migrations from HCL
  - Run atlas migrate diff initial_schema --env local_postgresql_watch
  - Repeat for keygen and sign environments
  - Verify migration SQL contains PostgreSQL syntax (BIGSERIAL, CHECK, etc.)
  - Validate migrations don't reference MySQL-specific features
  - _Requirements: 7.3, 7.4, 7.6_

- [ ] 9.4 (P) Extend Atlas Makefile targets for PostgreSQL
  - Update make/db_atlas.mk: Add ATLAS_POSTGRESQL_SCHEMAS variable
  - Create atlas-lint-postgresql target for schema validation
  - Update combined atlas-lint to include PostgreSQL schemas
  - Add PostgreSQL to atlas-validate and atlas-status targets
  - _Requirements: 5.3, 7.6, 10.6_

### Phase 5: Docker Compose Integration

- [ ] 10. Add PostgreSQL to Docker Compose environment
- [ ] 10.1 Configure PostgreSQL database service
  - Add wallet-db-postgresql service to compose.yaml
  - Use image: postgres:18.2
  - Configure volume: wallet-db-postgresql:/var/lib/postgresql/data
  - Mount init scripts: ./docker/postgresql/init.d:/docker-entrypoint-initdb.d
  - Set environment: POSTGRES_USER=postgres, POSTGRES_PASSWORD=postgres
  - Expose port: ${POSTGRESQL_PORT:-5432}:5432
  - Add healthcheck: pg_isready -U postgres (interval 10s, retries 5)
  - Add to db network for connectivity with migration services
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.6, 8.8_

- [ ] 10.2 (P) Create PostgreSQL initialization script
  - Create docker/postgresql/init.d/01_create_databases.sh
  - Use psql with ON_ERROR_STOP=1 for safety
  - Create three databases: watch, keygen, sign
  - Grant all privileges to postgres user for each database
  - Make script executable (chmod +x)
  - _Requirements: 8.4, 8.5_

- [ ] 10.3 Add PostgreSQL migration services to Docker Compose
  - Add wallet-db-migrate-postgresql-watch service
  - Use image: arigaio/atlas:1.1.0 (matching Atlas upgrade)
  - Configure depends_on: wallet-db-postgresql with service_healthy condition
  - Set command: migrate apply --dir file://migrations/postgresql_watch
  - Set connection URL: postgres://postgres:postgres@wallet-db-postgresql:5432/watch
  - Add similar services for keygen and sign migrations
  - Set restart: "no" for one-time migration execution
  - _Requirements: 8.7, 8.8, 10.2_

- [ ] 10.4 (P) Add PostgreSQL volume definition
  - Add wallet-db-postgresql volume to compose.yaml volumes section
  - Ensure PostgreSQL and MySQL volumes coexist for migration testing
  - _Requirements: 8.2, 8.8_

### Phase 6: Testing & Validation

- [ ] 11. Implement PostgreSQL connection tests
- [ ] 11.1 (P) Create connection factory unit tests
  - Create pkg/db/postgresql/connection_test.go
  - Test successful connection with valid configuration
  - Test connection failure with invalid host (verify error message)
  - Test connection failure with invalid credentials
  - Test connection pool configuration (max open, max idle, lifetime)
  - Test SSL mode variations (disable, require, verify-ca, verify-full)
  - _Requirements: 9.1_

- [ ] 12. Implement PostgreSQL repository integration tests
- [ ] 12.1 Create address repository integration tests
  - Create test file for PostgreSQL address repository
  - Test CRUD operations: Insert, GetByID, GetAll, Update
  - Test type conversion: domain.Address ↔ sqlcgen.Address
  - Test NULL handling for optional fields (UpdatedAt)
  - Test CHECK constraint violations (invalid coin/account values)
  - Use PostgreSQL test database separate from production
  - _Requirements: 9.2, 9.4_

- [ ] 12.2 Create transaction repository integration tests
  - Test BTC transaction repositories (tx, input, output)
  - Test ETH and XRP transaction repositories
  - Verify decimal precision for cryptocurrency amounts (NUMERIC 26,10)
  - Test concurrent access (multiple connections)
  - Validate query filtering (by coin type, status, dates)
  - _Requirements: 9.2, 9.4_

- [ ] 12.3 Create account key repository integration tests
  - Test account key repositories for BTC, ETH, XRP, auth
  - Test seed and nonce repository operations
  - Test XRP signer list repositories
  - Verify private key storage and retrieval
  - _Requirements: 9.2, 9.4_

- [ ] 13. Implement cross-database compatibility tests
- [ ] 13.1 Create cross-database query validation tests
  - Execute same queries against MySQL, SQLite, PostgreSQL
  - Compare result sets for consistency (row counts, data values)
  - Test data type precision across databases (especially NUMERIC)
  - Verify enum-like field validation (CHECK vs ENUM behavior)
  - _Requirements: 6.1, 6.5, 9.2, 9.5_

- [ ] 13.2 (P) Test schema migration compatibility
  - Apply Atlas migrations to PostgreSQL test database
  - Run atlas migrate apply for watch, keygen, sign schemas
  - Verify schema structure matches HCL definitions
  - Test migration idempotency (apply twice, no errors)
  - Validate atlas schema diff shows no differences after migration
  - _Requirements: 7.4, 9.3_

- [ ] 14. Implement end-to-end integration tests
- [ ] 14.1 Test wallet initialization with PostgreSQL
  - Load wallet config with database.type = "postgresql"
  - Initialize DI container and verify repository creation
  - Execute basic wallet operations (address generation, transaction creation)
  - Verify connection pool behavior under load
  - Test graceful shutdown (all connections closed)
  - _Requirements: 4.1, 4.5, 9.2_

- [ ] 14.2 Test configuration validation errors
  - Test missing required fields (host, dbname, user, pass)
  - Verify descriptive error messages for each validation failure
  - Test invalid database.type value handling
  - _Requirements: 1.2, 1.3, 1.4, 1.5_

### Phase 7: Migration Documentation

- [ ] 15. Create PostgreSQL migration documentation
- [ ] 15.1 Document MySQL to PostgreSQL migration
  - Create docs/migration/mysql-to-postgresql.md
  - Document data type mappings (MySQL → PostgreSQL)
  - Provide pgloader configuration example
  - List migration steps: backup, schema creation, data transfer, validation
  - Document rollback procedure
  - Include data integrity validation queries
  - _Requirements: 6.2, 6.3_

- [ ] 15.2 Document SQLite to PostgreSQL migration
  - Create docs/migration/sqlite-to-postgresql.md
  - Document SQLite-specific considerations (type affinity)
  - Provide migration tool recommendations
  - Document validation steps for data consistency
  - _Requirements: 6.2, 6.4_

- [ ] 15.3 Create PostgreSQL configuration examples
  - Add sample config files: config/wallet/postgresql/watch.toml
  - Add keygen.toml and sign.toml examples
  - Document all configuration parameters (host, port, sslmode, etc.)
  - Add troubleshooting section for common connection issues
  - _Requirements: 1.6, 6.3, 6.4_

- [ ] 15.4 Update main README with PostgreSQL support
  - Add PostgreSQL to supported databases section
  - Update database selection instructions
  - Link to migration documentation
  - Add PostgreSQL version requirement (18.2)
  - Document Docker Compose setup for PostgreSQL
  - _Requirements: 10.7_

---

## Requirements Coverage

All 10 requirements mapped to implementation tasks:

| Requirement | Tasks Covering |
|-------------|----------------|
| 1. Configuration Support | 1.1, 1.2, 14.2 |
| 2. Schema Generation | 2.1, 2.2, 2.3, 2.4 |
| 3. SQLC Code Generation | 3.1, 3.2, 5.1, 5.2, 5.3, 5.4, 6.1, 6.2, 6.3 |
| 4. Connection Management | 4.1, 4.2, 7.1, 7.2, 14.1 |
| 5. Build Integration | 8.1, 9.4 |
| 6. Migration Path | 13.1, 15.1, 15.2, 15.3 |
| 7. Atlas Migration | 9.1, 9.2, 9.3, 9.4, 13.2 |
| 8. Docker Compose | 10.1, 10.2, 10.3, 10.4 |
| 9. Testing | 11.1, 12.1, 12.2, 12.3, 13.1, 13.2, 14.1, 14.2 |
| 10. Atlas Upgrade | 8.2, 10.3 |

---

## Task Dependencies

**Sequential Dependencies:**
- 1.1 → 1.2 (config struct before validation)
- 2.1, 2.2, 2.3 → 2.4 (schemas before validation)
- 3.1 → 3.2 (config before generation)
- 4.1 → 4.2 (connection factory before DI integration)
- 3.2 → 5.x, 6.x (sqlcgen before repositories)
- 5.x, 6.x → 7.1, 7.2 (repositories before DI wiring)
- 9.1, 9.2 → 9.3 (Atlas config before migration generation)
- 10.1, 10.2 → 10.3 (database and init before migrations)
- All implementation → Testing (11, 12, 13, 14)
- All testing → Documentation (15)

**Parallel Opportunities (marked with P):**
- Configuration tasks (1.1, 1.2) - independent
- Schema conversion tasks (2.1, 2.2, 2.3) - no file conflicts
- Watch repositories (5.1-5.4) - separate files
- Cold repositories (6.1-6.3) - separate files
- Build/Atlas tasks (8.1, 8.2, 9.1, 9.2, 9.3, 9.4) - different tools/files
- Docker tasks (10.2, 10.4) - separate files
- Test files (11.1, 12.1-12.3, 13.2) - isolated test suites
- Documentation (15.1-15.4) - separate markdown files
