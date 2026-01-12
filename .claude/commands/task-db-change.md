# Database Change Task Command

Refer to @docs/commands/task-db-change.md for full documentation.

## Quick Reference

```
/task-db-change Issue #123
/task-db-change Add transaction status column
/task-db-change Create audit log table
```

## Parameters

- `{description}`: Change description or issue number

## Process Summary

1. Load context documents
2. Modify HCL schema in `tools/atlas/schemas/`
3. Run `make atlas-fmt && make atlas-lint`
4. Generate migrations: `make atlas-dev-reset`
5. Regenerate SQLC: `make sqlc`
6. Verify and create PR

## Required Context

- @docs/ai-agents/task-contexts/db-change.md
- @docs/ai-agents/guidelines/database.md

## Verification

```bash
make check-build && make gotest
```
