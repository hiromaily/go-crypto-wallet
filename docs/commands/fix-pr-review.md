# Fix PR Review Command

Address review comments on a pull request.

## Parameters

- `{pr_number}`: Pull request number (e.g., `#123`)

## Quick Reference

```bash
# 1. Fetch PR and comments
gh pr view {pr_number}
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments

# 2. Checkout PR branch
gh pr checkout {pr_number}

# 3. Fix comments (priority: security > functionality > quality)

# 4. Verify
make go-lint && make tidy && make check-build && make gotest

# 5. Push (updates existing PR)
git add <files>
git commit -m "fix(pr): address review comments"
git push
```
