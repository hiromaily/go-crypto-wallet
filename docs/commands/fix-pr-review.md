# Fix PR Review Command

## Purpose

Address review comments on a pull request.

## Parameters

- `{pr_number}`: Pull request number (e.g., `#123` or `123`)

## Process

### 1. Fetch PR Information

1. Get PR details: `gh pr view {pr_number}`
2. Get review comments: `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments`
3. Checkout PR branch

### 2. Analyze Comments

Categorize by type (priority order):
1. **Security**: Security concerns, sensitive data handling
2. **Functionality**: Bugs, logic errors, edge cases
3. **Architecture**: Design patterns, Clean Architecture violations
4. **Code quality**: Style, naming, structure
5. **Documentation**: Missing comments, unclear code

### 3. Implement Fixes

For each comment:
1. Address the specific feedback
2. Follow project coding standards
3. Test the fix
4. Document changes

### 4. Commit and Push

```bash
git add <files>
git commit -m "fix(pr): address review comments for PR #{pr_number}

- {fix 1}
- {fix 2}

Addresses feedback on PR #{pr_number}"

git push origin {branch-name}
```

## Verification

```bash
make go-lint
make tidy
make check-build
make gotest
```

## Related Documents

- [Workflow Guidelines](../ai-agents/guidelines/workflow.md)
- [Coding Standards](../standards/coding-conventions.md)
