# Fix Issue

Fix a GitHub issue following the standard workflow.

## Parameters

- `{issue_number}`: GitHub issue number

## Workflow

See `.codex/rules/general.md` for detailed workflow.

```bash
gh issue view {issue_number}
# Create branch, implement, verify, commit, PR
```
