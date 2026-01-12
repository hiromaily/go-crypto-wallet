# Codex Prompts

Minimal prompts referencing shared workflow.

## Available Prompts

| Prompt | Description |
|--------|-------------|
| `fix-issue #123` | Fix a GitHub issue |
| `fix-linter` | Fix linter errors |
| `fix-pr-review #123` | Address PR review comments |

## Workflow

All prompts follow the workflow in `.codex/rules/general.md`:

1. Branch from `main`
2. Implement following Clean Architecture
3. Verify: `make go-lint && make tidy && make check-build && make gotest`
4. Self-review
5. Commit & PR
