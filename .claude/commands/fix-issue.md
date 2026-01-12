# Fix Issue #{issue_number}

Fix a GitHub issue following the standard workflow.

## Skill Reference

**Use the `go-development` Skill** for branch management, verification, and commit workflow.

## Process

1. **Fetch issue**: `gh issue view {issue_number}`
2. **Create branch**: Follow `go-development` Skill
3. **Implement fix**: Follow Clean Architecture
4. **Verify**: Run verification commands from Skill
5. **Self-review**: Complete checklist from Skill
6. **Commit & PR**: Follow Skill workflow

## Parameters

- `{issue_number}`: GitHub issue number (e.g., `#123` or `123`)

## Example

```
/fix-issue #123
```
