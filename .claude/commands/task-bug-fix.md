# Bug Fix Task Command

Refer to @docs/commands/task-bug-fix.md for full documentation.

## Quick Reference

```
/task-bug-fix Issue #123 Chain: BTC
/task-bug-fix GetAddressInfo error Chain: BCH
/task-bug-fix Database connection timeout
```

## Parameters

- `{description}`: Bug description or issue number
- `{chain}` (optional): BTC, BCH, ETH, XRP

## Process Summary

1. Load context documents
2. Analyze problem and root cause
3. Implement fix
4. Verify and create PR

## Required Context

- @docs/ai-agents/task-contexts/bug-fix.md
- @docs/ai-agents/guidelines/workflow.md
- Chain-specific: @docs/ai-agents/task-contexts/chains/{chain}.md

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```
