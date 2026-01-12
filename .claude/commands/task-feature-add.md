# Feature Addition Task Command

Refer to @docs/commands/task-feature-add.md for full documentation.

## Quick Reference

```
/task-feature-add Issue #123 Chain: ETH
/task-feature-add Add new transaction type Chain: BTC
/task-feature-add Implement batch processing
```

## Parameters

- `{description}`: Feature description or issue number
- `{chain}` (optional): BTC, BCH, ETH, XRP

## Process Summary

1. Load context documents
2. Plan feature following Clean Architecture
3. Implement (domain → application → infrastructure → interface)
4. Verify and create PR

## Required Context

- @docs/ai-agents/task-contexts/feature-add.md
- @docs/ai-agents/guidelines/architecture.md
- Chain-specific: @docs/ai-agents/task-contexts/chains/{chain}.md

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```
