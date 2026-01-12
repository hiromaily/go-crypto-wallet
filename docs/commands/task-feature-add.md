# Feature Addition Task Command

## Purpose

Execute a feature addition workflow with proper context loading.

## Parameters

- `{description}`: Feature description or issue number
- `{chain}` (optional): Target cryptocurrency (BTC, BCH, ETH, XRP)

## Process

### 1. Load Context

Required documents:
- `docs/ai-agents/task-contexts/feature-add.md`
- `docs/ai-agents/guidelines/workflow.md`
- `docs/ai-agents/guidelines/architecture.md`

Chain-specific (if `{chain}` specified):
- `docs/ai-agents/task-contexts/chains/{chain}.md`

### 2. Plan Feature

1. Understand requirements
2. Design following Clean Architecture
3. Identify affected layers
4. Plan test cases

### 3. Implement

1. Domain layer changes (if needed)
2. Application layer (use cases, ports)
3. Infrastructure layer (implementations)
4. Interface adapters (CLI, handlers)

### 4. Verify and Commit

```bash
make go-lint && make tidy && make check-build && make gotest

git add <files>
git commit -m "feat: {description}

Closes #{issue_number}"

gh pr create --title "Feature: {description}"
```

## Examples

```
/task-feature-add Issue #123 Chain: ETH
/task-feature-add Add new transaction type support Chain: BTC
/task-feature-add Implement batch processing
```

## Related Documents

- [Feature Add Context](../ai-agents/task-contexts/feature-add.md)
- [Architecture Guidelines](../ai-agents/guidelines/architecture.md)
