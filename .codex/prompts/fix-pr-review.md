# Fix PR Review Prompt

Refer to [docs/commands/fix-pr-review.md](../../docs/commands/fix-pr-review.md) for full documentation.

## Quick Reference

```
fix-pr-review #123
```

## Parameters

- `{pr_number}`: Pull request number

## Process Summary

1. Fetch PR: `gh pr view {pr_number}`
2. Get comments: `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments`
3. Categorize comments (security > functionality > quality)
4. Address each comment
5. Commit and push

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```

## Related

- [Workflow Guidelines](../../docs/standards/workflow.md)
- [Coding Conventions](../../docs/standards/coding-conventions.md)
