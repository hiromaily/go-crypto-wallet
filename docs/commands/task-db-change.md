# Database Change Task Command

## Purpose

Execute a database schema change workflow with proper context loading.

## Parameters

- `{description}`: Change description or issue number

## Process

### 1. Load Context

Required documents:
- `docs/ai-agents/task-contexts/db-change.md`
- `docs/ai-agents/guidelines/database.md`
- `docs/ai-agents/guidelines/workflow.md`

### 2. Plan Schema Change

1. Understand requirements
2. Design schema modifications
3. Plan migration strategy
4. Consider backward compatibility

### 3. Implement

1. Modify HCL schema: `tools/atlas/schemas/`
2. Format and lint: `make atlas-fmt && make atlas-lint`
3. Generate migrations: `make atlas-dev-reset`
4. Verify migration: `docker compose down -v && docker compose up`
5. Regenerate SQLC: `make sqlc`

### 4. Verify and Commit

```bash
make check-build && make gotest

git add <files>
git commit -m "feat(db): {description}

Closes #{issue_number}"

gh pr create --title "DB: {description}"
```

## Guidelines

- **Backward compatibility**: Consider impact on existing data
- **Migration safety**: Test migrations thoroughly
- **Auto-generated files**: SQLC generates code - regenerate, don't edit

## Examples

```
/task-db-change Issue #123
/task-db-change Add new transaction status column
/task-db-change Create audit log table
```

## Related Documents

- [Database Change Context](../ai-agents/task-contexts/db-change.md)
- [Database Guidelines](../ai-agents/guidelines/database.md)
