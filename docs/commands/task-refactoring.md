# Refactoring Task Command

## Purpose

Execute a refactoring workflow with proper context loading.

## Parameters

- `{description}`: Refactoring description or issue number

## Process

### 1. Load Context

Required documents:
- `docs/ai-agents/task-contexts/refactoring.md`
- `docs/ai-agents/guidelines/workflow.md`
- `docs/ai-agents/guidelines/architecture.md`

### 2. Analyze Current State

1. Understand current implementation
2. Identify code smells or issues
3. Plan refactoring approach
4. Ensure no behavior changes

### 3. Implement

1. Make incremental changes
2. Verify tests pass after each step
3. Update documentation if needed

### 4. Verify and Commit

```bash
make go-lint && make tidy && make check-build && make gotest

git add <files>
git commit -m "refactor: {description}

Closes #{issue_number}"

gh pr create --title "Refactor: {description}"
```

## Guidelines

- **No behavior changes**: Refactoring must preserve functionality
- **Incremental**: Make small, verifiable changes
- **Test coverage**: Ensure adequate tests exist before refactoring

## Examples

```
/task-refactoring Issue #123
/task-refactoring Extract common validation logic
/task-refactoring Simplify error handling in repository layer
```

## Related Documents

- [Refactoring Context](../ai-agents/task-contexts/refactoring.md)
- [Architecture Guidelines](../ai-agents/guidelines/architecture.md)
