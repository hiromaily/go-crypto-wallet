# Documentation References - PostgreSQL Integration

**Date**: 2024-02-15
**Status**: Updated to reference new database management documentation

## Overview

This document tracks all references to database management documentation within the PostgreSQL integration specification. These references provide implementation guidance and detailed workflows for developers.

## Referenced Documentation

### Primary Guides

1. **[Database Schema Changes Guide](../../../docs/guidelines/database-schema-changes.md)**
   - Complete step-by-step workflow for schema modifications
   - Multi-database considerations (MySQL, SQLite, PostgreSQL)
   - Three detailed scenarios with examples
   - Testing procedures for all databases
   - Troubleshooting guide
   - Best practices

2. **[Database Quick Reference](../../../docs/guidelines/database-quick-reference.md)**
   - Command cheat sheet
   - Data type mapping tables
   - Common workflows
   - File location guide
   - Testing commands
   - Troubleshooting quick fixes

3. **[Database Management](../../../docs/guidelines/database.md)**
   - Overview of database tools and architecture
   - Quick workflow summary
   - Multi-database workflow

4. **[Database Architecture](../../../docs/development/database.md)**
   - Detailed database setup and operations
   - Container configuration
   - Common operations

## References in Specification Files

### design.md

**Line ~13**: Overview section
```markdown
**📚 Implementation References**:
- [Database Schema Changes Guide] - Complete workflow for multi-database schema modifications
- [Database Quick Reference] - Command cheat sheet and data type mappings
- [Database Management] - Overview of database tools and architecture
```

**Line ~481**: Data Type Mappings section
```markdown
**📘 Complete reference**: See [Database Quick Reference - Data Type Mapping]
for comprehensive mapping across MySQL, SQLite, and PostgreSQL.
```

**Line ~1027**: Migration Strategy section
```markdown
**📘 Complete workflow guide**: See [Database Schema Changes Guide]
for detailed step-by-step schema modification workflow across all databases.
```

### research.md

**Line ~289**: New "Project Documentation" section
```markdown
## Project Documentation

**Database Management Guides** (Created 2024-02-15):
- [Database Schema Changes Guide] - Complete workflow for multi-database schema modifications
- [Database Quick Reference] - Command cheat sheet and data type mapping tables
- [Database Management] - Overview of database tools and architecture
- [Database Architecture] - Detailed database setup and operations
```

### tasks.md

**Line ~11**: Overview section
```markdown
**📘 Implementation References**:
- [Database Schema Changes Guide] - Complete workflow for schema modifications
- [Database Quick Reference] - Command cheat sheet and data type mappings
- See design.md and research.md for architectural decisions and technology research
```

**Line ~37**: Phase 2 header
```markdown
**📘 Reference**: See [Database Quick Reference - Data Type Mapping]
for MySQL→PostgreSQL conversion table.
```

**Line ~103**: Phase 3 header
```markdown
**📘 Reference**: Follow [Database Schema Changes Guide - Scenario 1]
for complete schema extraction and sqlc workflow.
```

## Key Topics Covered in Documentation

### Schema Change Workflow

From [Database Schema Changes Guide](../../../docs/guidelines/database-schema-changes.md):

1. **Step-by-Step Workflows**:
   - Adding a new column
   - Adding a new table
   - Modifying existing columns

2. **Multi-Database Considerations**:
   - Schema parity requirements
   - Data type mapping strategy
   - Repository pattern for each database
   - DI container switching logic

3. **Testing Procedures**:
   - Local MySQL test
   - SQLite E2E test
   - PostgreSQL test (coming soon)
   - Cross-database compatibility tests

### Data Type Mappings

From [Database Quick Reference](../../../docs/guidelines/database-quick-reference.md):

| Concept | MySQL | SQLite | PostgreSQL |
|---------|-------|--------|------------|
| Auto ID | `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| Boolean | `TINYINT(1)` | `INTEGER (0/1)` | `BOOLEAN` |
| Enum | `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| Decimal | `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |
| Timestamp | `DATETIME` | `TEXT (ISO8601)` | `TIMESTAMP` |
| Text (sized) | `VARCHAR(255)` | `TEXT` | `VARCHAR(255)` |
| Text (large) | `TEXT` | `TEXT` | `TEXT` |

### Common Commands

From [Database Quick Reference](../../../docs/guidelines/database-quick-reference.md):

```bash
# Schema change workflow
make atlas-fmt && make atlas-lint
make atlas-dev-reset
docker compose down -v && docker compose up -d
make atlas-migrate-docker
make sqlc && make sqlc-sqlite && make sqlc-postgresql

# Testing
make go-lint && make check-build && make gotest
```

## Implementation Guidance

### For Task 2.x (Atlas Migration & Docker Compose)

**Reference**: [Database Quick Reference - Data Type Mapping](../../../docs/guidelines/database-quick-reference.md#-data-type-mapping)

**Use For**:
- Converting MySQL types to PostgreSQL
- Understanding ENUM → TEXT CHECK conversion
- Verifying AUTO_INCREMENT → BIGSERIAL mapping

### For Task 4.x (Schema Extraction & Code Generation)

**Reference**: [Database Schema Changes Guide - Scenario 1](../../../docs/guidelines/database-schema-changes.md#scenario-1-adding-a-new-column)

**Use For**:
- Understanding the complete workflow: HCL → migrations → database → dump → extract → sqlc
- Following the established pattern for schema extraction
- Implementing pg_dump wrapper scripts

### For Task 5.x-7.x (Repository Implementations)

**Reference**: [Database Schema Changes Guide - Repository Pattern](../../../docs/guidelines/database-schema-changes.md#repository-pattern)

**Use For**:
- Understanding the repository pattern structure
- DI container integration examples
- Type-safe sqlc-generated code usage

### For Task 11.x-14.x (Testing)

**Reference**: [Database Schema Changes Guide - Testing Schema Changes](../../../docs/guidelines/database-schema-changes.md#testing-schema-changes)

**Use For**:
- Testing with all three databases
- Cross-database compatibility testing
- Integration test patterns

### For Task 15.x (Documentation)

**Reference**: All database management guides

**Use For**:
- Migration documentation structure
- Configuration examples
- Best practices and common pitfalls

## Benefits of Documentation References

1. **Single Source of Truth**:
   - Database workflow documented once, referenced everywhere
   - Updates to workflow automatically apply to all implementations

2. **Implementation Guidance**:
   - Step-by-step instructions for complex tasks
   - Examples and common patterns
   - Troubleshooting guidance

3. **Consistency**:
   - All three databases (MySQL, SQLite, PostgreSQL) follow same patterns
   - Unified command interface
   - Predictable workflows

4. **Maintainability**:
   - Documentation updates don't require spec changes
   - New database additions can reference existing guides
   - Easy to update best practices

## Next Steps

### During Implementation

1. **Before Starting Tasks**:
   - Read referenced documentation sections
   - Understand the complete workflow
   - Review data type mappings

2. **During Development**:
   - Use Quick Reference for commands
   - Follow Schema Changes Guide for workflows
   - Reference examples in documentation

3. **Testing**:
   - Follow testing procedures from guides
   - Verify multi-database compatibility
   - Run cross-database tests

4. **Documentation Updates**:
   - Update migration guides as needed
   - Add PostgreSQL-specific examples
   - Document lessons learned

## Maintenance

**When to Update References**:
- New documentation added to `docs/guidelines/`
- Workflow changes or improvements
- New database features added
- Breaking changes in tools (Atlas, sqlc)

**How to Update**:
1. Update relevant documentation in `docs/guidelines/`
2. Verify references in specs still point to correct sections
3. Test workflows to ensure accuracy
4. Update this tracking document

---

**Last Updated**: 2024-02-15
**Maintained By**: go-crypto-wallet team
