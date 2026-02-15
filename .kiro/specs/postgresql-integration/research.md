# Research & Design Decisions

## Summary
- **Feature**: `postgresql-integration`
- **Discovery Scope**: Extension (adding PostgreSQL to existing MySQL/SQLite support)
- **Key Findings**:
  - pgx driver provides 50-100% better performance than lib/pq and is actively maintained
  - TEXT with CHECK constraints preferred over native ENUMs for easier schema evolution
  - PostgreSQL 18.2 official Docker image available with multi-architecture support
  - Atlas 1.1 provides enhanced PostgreSQL support with improved migration capabilities

## Research Log

### PostgreSQL Driver Selection (pgx vs lib/pq)

**Context**: Go ecosystem offers two main PostgreSQL drivers - need to select the optimal one for performance and maintainability.

**Sources Consulted**:
- [pgx vs lib/pq performance comparison](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/)
- [Go database driver benchmarks](https://github.com/jackc/go_db_bench)
- [GitLab's recommendation for pgx](https://gitlab.com/gitlab-org/gitlab/-/merge_requests/49135)
- [sqlc + pgx article](https://brandur.org/sqlc)

**Findings**:
- **Performance**: pgx is 10-20% faster than lib/pq when using database/sql interface, and 50-100% faster with native pgx interface
- **Binary protocol**: pgx uses PostgreSQL binary format, reducing parsing overhead compared to lib/pq's text-based communication
- **Connection pooling**: pgx includes built-in pgxpool for connection management, lib/pq requires external solutions
- **Maintenance**: lib/pq is no longer actively maintained; pgx is actively developed
- **Compatibility**: pgx provides database/sql compatibility layer, allowing drop-in replacement
- **Statement preparation**: pgx automatically prepares and caches statements by default (3x performance improvement in some workloads)

**Implications**:
- **Decision**: Use `github.com/jackc/pgx/v5` as PostgreSQL driver
- Connection factory in `pkg/db/postgresql/` will use pgx driver
- Maintains compatibility with existing database/sql patterns used for MySQL/SQLite

### Enum Handling Strategy (PostgreSQL ENUM vs TEXT with CHECK)

**Context**: MySQL schemas use ENUM types extensively; need strategy for PostgreSQL conversion.

**Sources Consulted**:
- [Native enums vs CHECK constraints](https://making.close.com/posts/native-enums-or-check-constraints-in-postgresql/)
- [Enums vs Check Constraints (Crunchy Data)](https://www.crunchydata.com/blog/enums-vs-check-constraints-in-postgres)
- [PostgreSQL ENUM documentation](https://www.postgresql.org/docs/current/datatype-enum.html)
- [sqlc datatype mappings](https://docs.sqlc.dev/en/stable/reference/datatypes.html)

**Findings**:
- **PostgreSQL ENUM advantages**: Type safety, ordering, space efficiency (stored as catalog references)
- **PostgreSQL ENUM drawbacks**:
  - ALTER TYPE requires ACCESS EXCLUSIVE lock (blocks all table access)
  - Difficult to modify (add/remove values requires full table scan)
  - Cannot compare values from different enum types
- **TEXT + CHECK advantages**:
  - Easy modification (constraints use SHARE UPDATE EXCLUSIVE lock, allows concurrent updates)
  - Simple migration path from MySQL
  - No full table scan on constraint changes
  - Better flexibility for evolving schemas
- **TEXT + CHECK drawbacks**: Slightly larger storage (full text vs catalog reference), less type safety at DB level
- **sqlc support**: Both approaches work with sqlc; ENUMs map to Go type aliases, TEXT maps to string

**Implications**:
- **Decision**: Use TEXT columns with CHECK constraints for enum-like fields
- Rationale: Schema evolution flexibility outweighs storage efficiency in this use case
- Migration safety: SHARE UPDATE EXCLUSIVE lock allows concurrent operations during schema changes
- Example conversion: `coin ENUM('btc','bch','eth','xrp')` → `coin TEXT CHECK (coin IN ('btc','bch','eth','xrp'))`

### PostgreSQL Docker Image Version

**Context**: Need to verify PostgreSQL 18.2 availability and compatibility for local development.

**Sources Consulted**:
- [PostgreSQL official Docker Hub](https://hub.docker.com/_/postgres/)
- [PostgreSQL Docker tags](https://hub.docker.com/_/postgres/tags)
- [PostgreSQL Docker documentation](https://www.docker.com/blog/how-to-use-the-postgres-docker-official-image/)

**Findings**:
- PostgreSQL 18.2 available as `postgres:18.2` tag on Docker Hub
- Multi-architecture support: amd64, arm64v8, arm32v7, and others
- PGDATA environment variable changed in PostgreSQL 18: now `/var/lib/postgresql/18/docker` (version-specific)
- Official postgres image actively maintained by Docker Library
- Image includes standard tools: psql, pg_dump, pg_restore, etc.

**Implications**:
- Use `postgres:18.2` for Docker Compose service
- Update PGDATA path in initialization scripts if needed
- Healthcheck command: `pg_isready -U postgres`
- Init scripts location: `/docker-entrypoint-initdb.d/`

### Atlas Migration Tool Version Upgrade

**Context**: Current version is 1.0.0, upgrading to 1.1.0 for improved PostgreSQL support.

**Sources Consulted**:
- [Atlas GitHub releases](https://github.com/ariga/atlas/releases)
- [Atlas PostgreSQL guide](https://atlasgo.io/guides/postgres/automatic-migrations)
- [Atlas documentation](https://atlasgo.io/versioned/apply)

**Findings**:
- Atlas supports declarative schema management for PostgreSQL
- PostgreSQL partitions, functions, procedures, and domains now supported
- Docker image: `arigaio/atlas:1.1.0` available
- Dev database URL format: `docker://postgres/18` for PostgreSQL 18
- No breaking changes identified from 1.0 to 1.1 (need to verify on official releases page)
- HCL schemas are database-agnostic; migrations are engine-specific

**Implications**:
- Update Docker Compose to use `arigaio/atlas:1.1.0`
- Existing MySQL/SQLite migrations remain unchanged
- Create PostgreSQL-specific environments in atlas.hcl
- Reuse existing HCL schemas (watch.hcl, keygen.hcl, sign.hcl) for PostgreSQL

### Data Type Mappings (MySQL to PostgreSQL)

**Context**: Need precise mappings to ensure data integrity and compatibility.

**Findings**:
| MySQL Type | PostgreSQL Type | Notes |
|------------|-----------------|-------|
| AUTO_INCREMENT | SERIAL / BIGSERIAL | Use BIGSERIAL for BIGINT columns |
| tinyint(1) | BOOLEAN | PostgreSQL native boolean type |
| enum('a','b') | TEXT CHECK(...) | Use CHECK constraints per decision above |
| decimal(26,10) | NUMERIC(26,10) | Direct mapping, identical precision |
| datetime | TIMESTAMP | PostgreSQL TIMESTAMP without time zone |
| varchar(n) | VARCHAR(n) | Direct mapping |
| text | TEXT | Direct mapping |
| bigint | BIGINT | Direct mapping |
| int | INTEGER | Direct mapping |

**Implications**:
- Schema conversion scripts need to handle AUTO_INCREMENT → SERIAL transformation
- Enum values must be extracted and converted to CHECK constraints
- All other types have straightforward mappings

## Architecture Pattern Evaluation

### Pattern: Repository Pattern with Database-Specific Implementations

**Description**: Follow existing MySQL/SQLite pattern - one repository implementation per database engine.

**Strengths**:
- Proven pattern (working well for MySQL/SQLite)
- Clear separation of concerns
- Type-safe sqlc-generated code per database
- Independent evolution of database-specific optimizations
- No abstraction leakage

**Risks / Limitations**:
- Code duplication (~40 repository files)
- Must maintain three parallel implementations
- Schema drift risk if not carefully managed

**Notes**:
- Aligns with existing Clean Architecture in codebase
- DI container already handles database selection via switch statements
- Acceptable trade-off: code duplication for type safety and clarity

### Alternative Considered: Generic Repository Layer

**Description**: Create database-agnostic repository using reflection or code generation.

**Strengths**:
- Single implementation for all databases
- Easier to add future databases

**Risks / Limitations**:
- Significant refactoring of existing codebase
- Loss of type safety
- Performance overhead (reflection)
- Not aligned with current architecture

**Decision**: NOT RECOMMENDED - maintain repository pattern

## Design Decisions

### Decision: Use pgx Driver for PostgreSQL

**Context**: Need to select PostgreSQL driver for Go that balances performance, maintainability, and compatibility.

**Alternatives Considered**:
1. **lib/pq** - Traditional driver, database/sql compatible, no longer maintained
2. **pgx** - Modern driver, superior performance, actively maintained

**Selected Approach**: pgx (github.com/jackc/pgx/v5)

**Rationale**:
- 50-100% better performance than lib/pq
- Active maintenance and community support
- Built-in connection pooling (pgxpool)
- Automatic statement preparation and caching
- Full database/sql compatibility layer
- Better PostgreSQL-specific feature support

**Trade-offs**:
- ✅ Benefits: Performance, features, maintenance
- ❌ Compromises: None significant - pgx is drop-in compatible

**Follow-up**: Verify connection pool configuration matches existing MySQL/SQLite patterns

### Decision: TEXT with CHECK Constraints for Enums

**Context**: MySQL schemas use ENUM types extensively; PostgreSQL offers native ENUMs or TEXT+CHECK alternatives.

**Alternatives Considered**:
1. **PostgreSQL native ENUM** - Type-safe, space-efficient, difficult to modify
2. **TEXT with CHECK constraints** - Flexible, easy migration, slightly larger storage
3. **Lookup tables** - Referential integrity, normalized, adds join overhead

**Selected Approach**: TEXT with CHECK constraints

**Rationale**:
- Schema evolution flexibility (SHARE UPDATE EXCLUSIVE vs ACCESS EXCLUSIVE lock)
- Simpler migration from MySQL ENUMs
- No blocking operations during constraint modifications
- Acceptable storage overhead for this use case
- Maintains similar validation semantics to MySQL

**Trade-offs**:
- ✅ Benefits: Migration safety, flexibility, concurrent updates
- ❌ Compromises: ~4 bytes per value vs catalog reference (negligible for wallet data)

**Follow-up**: Document enum migration pattern in migration guide

### Decision: Reuse Existing HCL Schemas for PostgreSQL

**Context**: Atlas supports database-agnostic HCL schemas; need to decide if separate PostgreSQL HCL files are needed.

**Alternatives Considered**:
1. **Reuse existing HCL** - Single source of truth, database-agnostic
2. **Separate PostgreSQL HCL** - Database-specific optimizations, independent evolution

**Selected Approach**: Reuse existing HCL schemas (watch.hcl, keygen.hcl, sign.hcl)

**Rationale**:
- HCL schemas are already database-agnostic
- Single source of truth reduces maintenance burden
- Atlas generates database-specific migrations from HCL
- No PostgreSQL-specific schema features needed at this time

**Trade-offs**:
- ✅ Benefits: Single source of truth, reduced maintenance
- ❌ Compromises: Cannot use PostgreSQL-specific features easily

**Follow-up**: If PostgreSQL-specific features needed later, can create separate HCL files

### Decision: Manual Schema Conversion (MySQL → PostgreSQL)

**Context**: Need to convert 3 MySQL schema files to PostgreSQL syntax.

**Alternatives Considered**:
1. **Automated conversion tool** - Fast, consistent, requires tooling development
2. **Manual conversion** - Controlled, educational, time-consuming
3. **Atlas-generated from HCL** - Automated, requires HCL as source

**Selected Approach**: Manual conversion with validation

**Rationale**:
- Only 3 schema files (watch, keygen, sign)
- Opportunity to review and validate each table
- Ensures understanding of schema nuances
- Atlas can generate from HCL for future changes

**Trade-offs**:
- ✅ Benefits: Quality control, learning opportunity
- ❌ Compromises: Initial time investment (1-2 days)

**Follow-up**: Create conversion checklist and validation tests

## Risks & Mitigations

- **Risk**: Data type precision loss during migration
  - **Mitigation**: Comprehensive integration tests comparing decimal handling across databases; explicit test cases for min/max cryptocurrency amounts

- **Risk**: Query syntax incompatibility between databases
  - **Mitigation**: sqlc validates syntax per database engine; integration tests run against all three databases

- **Risk**: Connection pool configuration differences
  - **Mitigation**: Document PostgreSQL connection pool best practices; use conservative defaults initially; monitor connection usage in production

- **Risk**: Enum migration complexity
  - **Mitigation**: Document conversion pattern; create helper function for CHECK constraint generation; validate all enum values preserved

- **Risk**: Atlas 1.0 → 1.1 breaking changes
  - **Mitigation**: Test migration on sample schemas first; review release notes thoroughly; maintain fallback to 1.0 if needed

- **Risk**: PostgreSQL 18-specific behaviors
  - **Mitigation**: Use Docker for consistent local environment; document PostgreSQL version requirements; test against multiple PostgreSQL versions if possible

## References

**PostgreSQL Drivers**:
- [pgx vs lib/pq comparison](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/)
- [pgx GitHub repository](https://github.com/jackc/pgx)
- [sqlc + pgx integration](https://brandur.org/sqlc)

**Enum Handling**:
- [Native enums vs CHECK constraints](https://making.close.com/posts/native-enums-or-check-constraints-in-postgresql/)
- [Enums vs Check Constraints - Crunchy Data](https://www.crunchydata.com/blog/enums-vs-check-constraints-in-postgres)
- [PostgreSQL ENUM documentation](https://www.postgresql.org/docs/current/datatype-enum.html)

**Docker & Infrastructure**:
- [PostgreSQL official Docker image](https://hub.docker.com/_/postgres/)
- [PostgreSQL Docker tags](https://hub.docker.com/_/postgres/tags)

**Atlas Migration Tool**:
- [Atlas GitHub releases](https://github.com/ariga/atlas/releases)
- [Atlas PostgreSQL guide](https://atlasgo.io/guides/postgres/automatic-migrations)
- [Atlas documentation](https://atlasgo.io/)

**sqlc**:
- [sqlc datatype documentation](https://docs.sqlc.dev/en/stable/reference/datatypes.html)
- [sqlc PostgreSQL support](https://docs.sqlc.dev/en/stable/reference/language-support.html)
