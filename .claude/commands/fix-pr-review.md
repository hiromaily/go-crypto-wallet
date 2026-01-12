# Fix PR Review #{pr_number}

Address review comments on a pull request.

## Skill Reference

**Use the `go-development` Skill** for verification and commit workflow.

## Process

1. **Fetch PR**: `gh pr view {pr_number}`
2. **Get comments**: `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments`
3. **Prioritize**: security > functionality > quality > docs
4. **Fix**: Address each comment
5. **Verify**: Run verification commands from Skill
6. **Push**: `git push` (updates existing PR)

## Parameters

- `{pr_number}`: Pull request number (e.g., `#123` or `123`)

## Example

```
/fix-pr-review #123
```
