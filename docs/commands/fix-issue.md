# Fix Issue Command

Fix a GitHub issue following the standard workflow.

## Parameters

- `{issue_number}`: GitHub issue number (e.g., `#123`)

## Quick Reference

```bash
# 1. Fetch issue
gh issue view {issue_number}

# 2. Create branch
git fetch origin && git checkout main && git reset --hard origin/main
git checkout -b fix/issue-{number}-{description}

# 3. Implement fix (follow Clean Architecture)

# 4. Verify
make go-lint && make tidy && make check-build && make gotest

# 5. Commit
git add <files>
git commit -m "fix: {description}

Closes #{issue_number}"

# 6. Create PR
git push -u origin {branch-name}
gh pr create --title "Fix: {description}" --body "Closes #{issue_number}"
```
