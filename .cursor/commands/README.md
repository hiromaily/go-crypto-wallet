# Cursor Commands

Minimal commands referencing shared workflow.

## Available Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Fix a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |

## Workflow

All commands follow the workflow in `.cursor/rules/general.mdc`:

1. Branch from `main`
2. Implement following Clean Architecture
3. Verify: `make go-lint && make tidy && make check-build && make gotest`
4. Self-review
5. Commit & PR
