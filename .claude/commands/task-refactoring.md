# Refactoring Task Command

Refer to @docs/commands/task-refactoring.md for full documentation.

## Quick Reference

```
/task-refactoring Issue #123
/task-refactoring Extract common validation logic
/task-refactoring Simplify error handling
```

## Parameters

- `{description}`: Refactoring description or issue number

## Process Summary

1. Load context documents
2. Analyze current implementation
3. Make incremental changes (no behavior change)
4. Verify and create PR

## Required Context

- @docs/ai-agents/task-contexts/refactoring.md
- @docs/ai-agents/guidelines/architecture.md

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```
