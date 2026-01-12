# Bug Fix Task Command

## Purpose

Execute a bug fix workflow with proper context loading.

## Parameters

- `{description}`: Bug description or issue number
- `{chain}` (optional): Target cryptocurrency (BTC, BCH, ETH, XRP)

## Process

### 1. Load Context

Required documents:
- `docs/ai-agents/task-contexts/bug-fix.md`
- `docs/ai-agents/guidelines/workflow.md`
- `docs/ai-agents/guidelines/core.md`

Chain-specific (if `{chain}` specified):
- `docs/ai-agents/task-contexts/chains/{chain}.md`

### 2. Analyze Problem

1. If issue number: `gh issue view {issue_number}`
2. Identify reproduction steps
3. Determine affected scope
4. Find root cause

### 3. Implement Fix

Follow bug fix workflow from task context document.

### 4. Verify and Commit

```bash
make go-lint && make tidy && make check-build && make gotest

git add <files>
git commit -m "fix: {description}

Closes #{issue_number}"

gh pr create --title "Fix: {description}"
```

## Examples

```
/task-bug-fix Issue #123 Chain: BTC
/task-bug-fix GetAddressInfo returns wrong result Chain: BCH
/task-bug-fix Database connection timeout
```

## Related Documents

- [Bug Fix Context](../ai-agents/task-contexts/bug-fix.md)
- [Testing Standards](../standards/testing.md)
