---
name: db-migration
description: Database schema and migration workflow. Use when modifying database schemas in tools/atlas/ or SQLC queries in tools/sqlc/.
---

# Database Migration Workflow

Workflow for database schema and migration changes.

## Prerequisites

- **Use `git-workflow` Skill** for branch, commit, and PR workflow.
- **Refer to `.claude/rules/hcl.md`** for HCL schema rules (SSOT).
- **Refer to `.claude/rules/sql.md`** for SQL query rules (SSOT).

## Applicable Files

| Path                                     | Description                              |
| ---------------------------------------- | ---------------------------------------- |
| `tools/atlas/schemas/{db_dialect}/*.hcl` | HCL schema definitions (source of truth) |
| `tools/sqlc/queries/{db_dialect}/*.sql`  | SQLC query definitions                   |

## Workflow

### 1. Modify Schema (HCL)

Edit HCL files in `tools/atlas/schemas/{db_dialect}/`.

Dialect directories:

- `tools/atlas/schemas/postgres/*.hcl`
- `tools/atlas/schemas/mysql/*.hcl`

### 2. Verify HCL (from rules/hcl.md)

```bash
make atlas-fmt && make atlas-lint
```

### 3. Regenerate All Artifacts

Run **all three** targets after modifying HCL schemas:

```bash
# Regenerate migrations, extract SQLC schemas, and generate Go code for PostgreSQL
make regenerate-all-from-atlas

# Regenerate for MySQL
make regenerate-all-from-atlas-mysql

# Convert PostgreSQL schemas to SQLite and generate Go code for SQLite
make regenerate-all-from-atlas-sqlite
```

Each target handles: Atlas migrations → Docker DB reset → SQLC schema extraction → SQLC code generation.

### 4. Verify Migration Files

After running the regenerate targets, confirm each DB dialect has exactly one migration SQL file:

- `tools/atlas/migrations/postgres/{db-name}/*.sql` ← one file, updated
- `tools/atlas/migrations/mysql/{db-name}/*.sql` ← one file, updated

DB names: `keygen`, `sign`, `watch`

### 5. Verify Go Code

```bash
make check-build && make go-test
```

## Self-Review Checklist

- [ ] HCL format/lint passes
- [ ] `make regenerate-all-from-atlas` succeeds (postgres)
- [ ] `make regenerate-all-from-atlas-mysql` succeeds (mysql)
- [ ] `make regenerate-all-from-atlas-sqlite` succeeds (sqlite)
- [ ] Each `tools/atlas/migrations/{postgres,mysql}/{db-name}/` has exactly one `.sql` file and it is updated
- [ ] SQLC generates correctly
- [ ] Go build passes

## Related

- `.claude/rules/hcl.md` - HCL rules (SSOT)
- `.claude/rules/sql.md` - SQL rules (SSOT)
- `go-development` - Go verification after SQLC generation
- `git-workflow` - Branch, commit, PR workflow
