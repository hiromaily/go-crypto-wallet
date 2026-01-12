# Fix Issue Command

## Purpose

Resolve GitHub issues by implementing fixes, creating commits, and submitting pull requests.

## Parameters

- `{issue_number}`: GitHub issue number (e.g., `#123` or `123`)
  - Multiple issues: `#123,#124` or `123 124` (comma/space-separated)
- `{base_branch}` (optional): Branch to base work on
  - `new` (default): Create from latest `origin/main`
  - `current`: Work on current branch
  - `<branch-name>`: Create from specified branch

## Process

### 1. Pre-Flight

1. Parse issue number(s) from input
2. Fetch issue details: `gh issue view {issue_number}`
3. Verify issue exists and is open
4. Prepare working branch based on mode

### 2. Implementation

For each issue:

1. **Analyze**: Understand problem, identify affected files
2. **Plan**: Break down solution, identify test cases
3. **Implement**: Follow Clean Architecture, coding standards
4. **Test**: Run `make gotest`, add new tests
5. **Verify**: Run verification commands
6. **Commit**: Create commit with descriptive message

### 3. Pull Request

1. Push branch: `git push origin {branch-name}`
2. Create PR: `gh pr create --title "Fix: {title} (Closes #{issue})" --body "..."`

## Verification

```bash
make go-lint
make tidy
make check-build
make gotest
```

## Commit Message Format

```
fix: resolve issue #{issue_number} - {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}
```

## Related Documents

- [Workflow Guidelines](../ai-agents/guidelines/workflow.md)
- [Coding Standards](../standards/coding-conventions.md)
- [Testing Standards](../standards/testing.md)
