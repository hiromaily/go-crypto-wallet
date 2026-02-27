---
task_type: test
description: Context for test addition and modification tasks
version: 1.0.0
---

# Task: Test (Test Addition/Modification)

## When to Use

- Adding new test cases
- Fixing or improving existing tests
- Improving test coverage
- Improving test reliability
- Fixing flaky tests

## Required Context Documents

### Must Read

| Document | Path | Purpose |
|----------|------|---------|
| Testing Standards | `docs/guidelines/testing.md` | Test strategy, coverage requirements |
| Coding Standards | `docs/guidelines/coding-conventions.md` | Coding conventions |
| Workflow Guidelines | `docs/guidelines/workflow.md` | Verification steps |

### Conditional Read

| Condition | Document | Path |
|-----------|----------|------|
| Integration tests | Architecture | `ARCHITECTURE.md` |
| E2E tests | E2E Script Rules | `.claude/rules/btc/e2e-script.md` |
| Mock creation | Internal Guidelines | `internal/AGENTS.md` |
| DB-related tests | Database Guidelines | `docs/database/db-management.md` |

## Test Types

### Unit Tests

- Test single packages/functions
- Replace external dependencies with mocks
- Fast execution

```bash
# Run all
make go-test

# Specific package
go test ./internal/domain/...
```

### Integration Tests

- Test multiple component interactions
- May use real DB/external services

```bash
# Run
make go-test-integration
```

### E2E Tests

- Test end-to-end workflows
- Test in near-production environment

```bash
# BTC E2E
make btc-e2e-descriptor
```

## Task-Specific Rules

### 1. Test File Placement

```
Unit Tests:
  {package}_test.go  # Same package

Integration Tests:
  test/integration/  # Dedicated directory

E2E Tests:
  scripts/operation/{chain}/e2e/  # E2E scripts
```

### 2. Test Naming Conventions

```go
// ✅ Good: Clear what is being tested
func TestCreateTransaction_WithValidInput_ReturnsTransaction(t *testing.T)
func TestCreateTransaction_WithInvalidAmount_ReturnsError(t *testing.T)

// ❌ Bad: Unclear
func TestCreateTransaction(t *testing.T)
func Test1(t *testing.T)
```

### 3. Table-Driven Tests

```go
// ✅ Good: Table-driven tests
func TestValidateAddress(t *testing.T) {
    tests := []struct {
        name    string
        address string
        wantErr bool
    }{
        {"valid address", "bc1q...", false},
        {"empty address", "", true},
        {"invalid format", "invalid", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAddress(tt.address)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAddress() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 4. Mock Usage

```go
// Mock generation
//go:generate mockery --name=Repository --output=./mocks

// Usage in tests
func TestUseCase(t *testing.T) {
    mockRepo := mocks.NewRepository(t)
    mockRepo.On("Get", mock.Anything).Return(entity, nil)

    uc := NewUseCase(mockRepo)
    // ...
}
```

### 5. Security Considerations

- **NEVER**: Include real private keys in tests
- **ALWAYS**: Use test keys/addresses
- **CHECK**: Ensure sensitive information is not logged

## Pre-Task Checklist

- [ ] Understood the functionality being tested
- [ ] Reviewed existing test patterns
- [ ] Determined test type (unit/integration/e2e)
- [ ] Identified required mocks
- [ ] Listed edge cases

## Verification Commands

```bash
# Unit Tests
make go-test              # Run all tests
go test -v ./path/to/... # Specific package
go test -run TestName    # Specific test

# Integration Tests
make go-test-integration

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# E2E Tests (BTC)
make btc-e2e-descriptor
```

## Test Workflow

### Step 1: Analyze Test Target

1. Read the code being tested
2. Identify inputs and outputs
3. List edge cases

### Step 2: Design Tests

1. Enumerate test cases
2. Cover both success and error paths
3. Consider boundary values

### Step 3: Implement Tests

1. Create/update test files
2. Use table-driven tests
3. Use appropriate assertions

### Step 4: Verify

```bash
# Confirm tests pass
make go-test

# Check coverage (optional)
go test -coverprofile=coverage.out ./...
```

### Step 5: Commit

```bash
git add <test-files>
git commit -m "test: add tests for {feature}

- Add unit tests for {component}
- Cover edge cases for {scenario}
- Improve coverage for {package}"
```

## Testing Requirements

### Required

- [ ] Success path tests
- [ ] Major error path tests
- [ ] Tests can run independently

### Recommended

- [ ] Edge case tests
- [ ] Boundary value tests
- [ ] Performance tests (for critical processing)

## Examples

### Example 1: Adding Unit Tests

```
User: "Add unit tests for Transaction creation"

Agent Actions:
1. Read docs/guidelines/testing.md
2. Review internal/domain/transaction/ implementation
3. Review existing test patterns
4. Design test cases
5. Implement with table-driven tests
6. Verify with make go-test
```

### Example 2: Adding Integration Tests

```
User: "Add integration tests for BTC address generation"

Agent Actions:
1. Read docs/guidelines/testing.md
2. Review existing patterns in test/integration/
3. Identify required external dependencies
4. Implement tests
5. Verify with make go-test-integration
```

### Example 3: Fixing Flaky Tests

```
User: "Fix the intermittently failing test"

Agent Actions:
1. Analyze test failure patterns
2. Identify race conditions or timing issues
3. Add appropriate synchronization/retries
4. Run multiple times to confirm stability
```

## Related Documents

- [Testing Standards](../guidelines/testing.md) - Detailed test strategy
- [Go Development Skill](../../.claude/skills/go-development/SKILL.md) - Go development workflow
- [Task Classification](../guidelines/task-classification.md) - Task classification SSOT
