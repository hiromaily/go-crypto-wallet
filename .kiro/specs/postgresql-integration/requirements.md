# Requirements Document

## Project Description (Input)

The go-crypto-wallet project currently supports two database backends:

- **MySQL** - Primary database for production deployments
- **SQLite** - Lightweight option for testing and development

The project uses [sqlc](https://github.com/sqlc-dev/sqlc) as a type-safe SQL code generator (ORM-alternative) and [Atlas](https://github.com/ariga/atlas) version 1.0 for schema migrations. Users select their preferred database engine through the `database.type` configuration field in TOML config files.

**Goals**:

1. **Add PostgreSQL Support**: Add PostgreSQL as a third database option, maintaining feature parity with MySQL and SQLite. Users should be able to select PostgreSQL by setting `database.type = "postgres"` in the configuration file, with the application handling PostgreSQL-specific connection management, schema generation, and query execution transparently.

2. **Upgrade Atlas**: Upgrade Atlas from version 1.0 to the current version 1.1, taking advantage of improved PostgreSQL support and new features while ensuring backward compatibility with existing migrations.

## Introduction

The go-crypto-wallet project currently supports MySQL and SQLite as database backends, using sqlc for type-safe SQL query generation and Atlas for schema migrations. This specification adds PostgreSQL as a third database option, allowing users to select their preferred database engine via configuration.

The implementation must maintain the existing architecture where database selection is determined by the `database.type` configuration field, with separate sqlc-generated code for each database engine. PostgreSQL 18.2 will be used as the target version, integrated into the Docker Compose development environment alongside MySQL.

## Requirements

### Requirement 1: Configuration Support for PostgreSQL

**Objective:** As a wallet operator, I want to configure PostgreSQL as the database backend, so that I can use PostgreSQL's advanced features and reliability for my wallet infrastructure.

#### Acceptance Criteria

1. When the configuration file specifies `database.type = "postgres"`, the Wallet System shall accept this as a valid database type
2. If `database.type` is "postgres" and `database.postgres.host` is empty, then the Wallet System shall return a validation error
3. If `database.type` is "postgres" and `database.postgres.dbname` is empty, then the Wallet System shall return a validation error
4. If `database.type` is "postgres" and `database.postgres.user` is empty, then the Wallet System shall return a validation error
5. If `database.type` is "postgres" and `database.postgres.pass` is empty, then the Wallet System shall return a validation error
6. The Wallet System shall support PostgreSQL connection parameters including host, port, database name, username, password, and SSL mode
7. Where database.type is "postgres", the Wallet System shall use the PostgreSQL configuration section and ignore MySQL and SQLite configurations

### Requirement 2: PostgreSQL Schema Generation

**Objective:** As a developer, I want PostgreSQL-compatible schema files generated from existing schemas, so that the database structure matches MySQL and SQLite implementations.

#### Acceptance Criteria

1. When PostgreSQL schema generation is executed, the Schema Generation Tool shall create PostgreSQL-compatible schema files for all wallet databases (watch, keygen, sign)
2. The Schema Generation Tool shall convert MySQL-specific data types to PostgreSQL equivalents (e.g., `AUTO_INCREMENT` to `SERIAL`, `tinyint(1)` to `BOOLEAN`, `enum` to appropriate types)
3. The Schema Generation Tool shall preserve all table structures, relationships, and constraints from MySQL schemas
4. The Schema Generation Tool shall store PostgreSQL schemas in a dedicated directory separate from MySQL and SQLite schemas
5. If schema conversion encounters incompatible features, then the Schema Generation Tool shall generate PostgreSQL-idiomatic alternatives

### Requirement 3: SQLC Code Generation for PostgreSQL

**Objective:** As a developer, I want sqlc to generate type-safe Go code for PostgreSQL queries, so that the application can interact with PostgreSQL using the same query interface as MySQL and SQLite.

#### Acceptance Criteria

1. When sqlc code generation is executed for PostgreSQL, the Build System shall generate Go code in a dedicated PostgreSQL package directory
2. The Build System shall configure sqlc to use PostgreSQL engine for syntax validation and type checking
3. The Build System shall generate identical query interfaces for PostgreSQL as exist for MySQL and SQLite implementations
4. The Build System shall use the same query definitions from `tools/sqlc/queries/*.sql` for all database engines
5. If query syntax is incompatible with PostgreSQL, then the Build System shall fail with a clear error message indicating the incompatibility

### Requirement 4: Database Connection Management

**Objective:** As a wallet operator, I want the application to establish and manage PostgreSQL connections based on configuration, so that I can run the wallet with PostgreSQL as the backend.

#### Acceptance Criteria

1. When the application starts with `database.type = "postgres"`, the Database Connection Manager shall establish a connection to PostgreSQL using the provided credentials
2. If PostgreSQL connection fails, then the Database Connection Manager shall return a descriptive error including the connection failure reason
3. While the application is running with PostgreSQL, the Database Connection Manager shall maintain connection pool settings appropriate for PostgreSQL
4. The Database Connection Manager shall support PostgreSQL-specific connection options including SSL mode and connection timeout
5. When the application shuts down, the Database Connection Manager shall properly close all PostgreSQL connections

### Requirement 5: Build and Verification Integration

**Objective:** As a developer, I want build commands to include PostgreSQL schema and code generation, so that PostgreSQL support is maintained automatically during development.

#### Acceptance Criteria

1. When `make sqlc` is executed, the Build System shall generate code for all three database engines (MySQL, SQLite, PostgreSQL)
2. When PostgreSQL schema files are modified, the Build System shall regenerate sqlc code for PostgreSQL
3. The Build System shall provide a verification command to check PostgreSQL schema validity
4. The Build System shall provide a command to validate PostgreSQL queries independently
5. If any PostgreSQL-related build step fails, then the Build System shall report the failure with clear context about which step failed

### Requirement 6: Migration Path and Compatibility

**Objective:** As a wallet operator, I want to migrate existing wallet data to PostgreSQL, so that I can transition from MySQL or SQLite to PostgreSQL without data loss.

#### Acceptance Criteria

1. The Database Schema shall maintain identical table structures across MySQL, SQLite, and PostgreSQL implementations
2. The Application shall use consistent data type mappings that preserve value ranges and precision across all database engines
3. The Application shall provide documentation for migrating data from MySQL to PostgreSQL
4. The Application shall provide documentation for migrating data from SQLite to PostgreSQL
5. When reading data from any database engine, the Application shall return data in a consistent format regardless of the underlying database

### Requirement 7: Atlas Migration Support for PostgreSQL

**Objective:** As a developer, I want Atlas to manage PostgreSQL schema migrations, so that database schema changes are tracked and versioned consistently across all database engines.

#### Acceptance Criteria

1. When Atlas configuration is updated for PostgreSQL, the Atlas Configuration shall define environments for all three PostgreSQL schemas (watch, keygen, sign)
2. The Atlas Configuration shall support PostgreSQL connection URLs with appropriate parameters (host, port, user, password, database, SSL mode)
3. When `atlas migrate diff` is executed for PostgreSQL, the Atlas CLI shall generate migration files based on HCL schema differences
4. When `atlas migrate apply` is executed for PostgreSQL, the Atlas CLI shall apply pending migrations to the PostgreSQL database
5. The Atlas Configuration shall use PostgreSQL-compatible dev database for migration validation (e.g., `docker://postgres/18`)
6. When migration validation fails, then the Atlas CLI shall report specific schema incompatibilities with PostgreSQL
7. The Atlas Configuration shall support both local and admin environments for PostgreSQL similar to MySQL configuration

### Requirement 8: Docker Compose Integration for PostgreSQL

**Objective:** As a developer, I want PostgreSQL 18.2 to run via Docker Compose, so that I can easily set up a local PostgreSQL development environment consistent with the team.

#### Acceptance Criteria

1. When Docker Compose is started, the Container Orchestration shall create a PostgreSQL 18.2 container with proper configuration
2. The Container Orchestration shall configure PostgreSQL container with dedicated volume for data persistence
3. The Container Orchestration shall expose PostgreSQL on a configurable port (default: 5432)
4. The Container Orchestration shall configure PostgreSQL with environment variables for database names (watch, keygen, sign), user, and password
5. When PostgreSQL container starts, the Container Orchestration shall execute initialization scripts to create the three wallet databases
6. The Container Orchestration shall include healthcheck configuration to verify PostgreSQL readiness before migration services start
7. The Container Orchestration shall configure Atlas migration services for PostgreSQL that automatically apply migrations on startup
8. The Container Orchestration shall ensure PostgreSQL and MySQL can coexist in the same Docker Compose environment for migration testing

### Requirement 9: Testing and Validation

**Objective:** As a developer, I want comprehensive tests for PostgreSQL integration, so that I can verify PostgreSQL functionality works correctly.

#### Acceptance Criteria

1. The Test Suite shall include unit tests for PostgreSQL connection establishment
2. The Test Suite shall include integration tests that verify all CRUD operations work with PostgreSQL
3. The Test Suite shall include tests that verify data consistency across database migrations
4. When integration tests are executed, the Test Suite shall use a PostgreSQL test database separate from production
5. The Test Suite shall verify that all existing queries work correctly with PostgreSQL backend

### Requirement 10: Atlas Version Upgrade

**Objective:** As a developer, I want to upgrade Atlas from version 1.0 to version 1.1, so that I can leverage improved PostgreSQL support, bug fixes, and new features in the latest stable release.

#### Acceptance Criteria

1. When Atlas is upgraded, the Build System shall update all Atlas references from version 1.0 to version 1.1
2. The Docker Compose configuration shall use the Atlas 1.1 Docker image (`arigaio/atlas:1.1.0`)
3. When Atlas 1.1 is used, the Migration System shall verify that all existing migrations remain compatible and functional
4. The Build System shall validate that all Atlas CLI commands (`migrate diff`, `migrate apply`, `schema apply`, etc.) work correctly with version 1.1
5. If Atlas 1.1 introduces breaking changes affecting existing workflows, then the Migration System shall provide migration guidance or compatibility adaptations
6. The Build System shall update Makefile targets to use Atlas 1.1 syntax and features where applicable
7. When Atlas upgrade is complete, the Documentation shall reflect the new version number in all relevant files (compose.yaml, Makefile, documentation)
