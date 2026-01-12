# Fix PR Review Command

Refer to [docs/commands/fix-pr-review.md](../../docs/commands/fix-pr-review.md) for full documentation.

## Quick Reference

```
/fix-pr-review #123
```

## Parameters

- `{pr_number}`: Pull request number

## Process Summary

1. Fetch PR and review comments
2. Categorize comments (security > functionality > quality)
3. Address each comment
4. Commit and push

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```
