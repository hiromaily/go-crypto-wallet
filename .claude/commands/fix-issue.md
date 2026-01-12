# Fix Issue Command

Refer to @docs/commands/fix-issue.md for full documentation.

## Quick Reference

```
/fix-issue #123
/fix-issue #123,#124
/fix-issue #123 current
```

## Parameters

- `{issue_number}`: GitHub issue number(s)
- `{base_branch}`: `new` (default), `current`, or branch name

## Process Summary

1. Fetch issue details: `gh issue view {issue_number}`
2. Create feature branch (unless `current` mode)
3. Implement fix following Clean Architecture
4. Run verification commands
5. Create commit and PR

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```

## Related

- @docs/standards/workflow.md
- @docs/standards/coding-conventions.md
