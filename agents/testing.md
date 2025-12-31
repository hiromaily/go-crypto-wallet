# Testing Guidelines

This document describes the testing strategy and requirements for the go-crypto-wallet project.

## Testing Principles

- Use `//go:build integration` tag for integration tests
- Separate unit tests and integration tests
- Measure and improve test coverage
- Write tests for all exported functions and methods
- Keep tests maintainable and readable

## Testing by Layer

### Domain Layer Testing

**Approach**: Pure unit tests without mocks

**Characteristics:**

- Test business logic in isolation
- No infrastructure dependencies required
- Fast, deterministic tests
- No mocks needed (pure functions)

**What to Test:**

- Value object validation
- Business rule enforcement
- Domain validators
- State transitions
- Entity lifecycle

**Example:**

```go
func TestAccountType_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   AccountType
        wantErr bool
    }{
        {"valid client", AccountTypeClient, false},
        {"valid receipt", AccountTypeReceipt, false},
        {"invalid empty", AccountType(""), true},
        {"invalid unknown", AccountType("unknown"), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Application Layer Testing (Use Cases)

**Approach**: Test with mocked infrastructure

**Current Status**: Use cases have constructor tests that verify:

- Use case can be instantiated with dependencies
- Correct interface implementation

**Future Testing Strategy** (see `docs/TESTING_STRATEGY.md`):

- Test use case orchestration logic
- Mock infrastructure services
- Verify DTO transformation
- Test error handling and wrapping

**What to Test:**

- Use case input validation
- Service coordination
- Error wrapping with context
- DTO transformation
- Business flow orchestration

### Infrastructure Layer Testing

**Approach**: Unit tests with mocked external dependencies + integration tests

**Unit Tests:**

- Mock external systems (database, API clients)
- Test error handling
- Test retry logic
- Test data transformation

**Integration Tests:**

- Use `//go:build integration` tag
- Test with real external systems (when possible)
- Use test databases/containers
- Verify end-to-end functionality

**What to Test:**

- Repository CRUD operations
- API client request/response handling
- Database connection management
- File I/O operations
- Network communication

### Interface Adapters Layer Testing

**Approach**: Test with mocked use cases

**What to Test:**

- Command argument parsing
- Output formatting
- Error message formatting
- CLI flag handling
- Use case integration (with mocked use cases)

## Test Organization

**File Naming:**

- Test files: `*_test.go` (same package)
- Integration tests: `*_integration_test.go` with `//go:build integration` tag

**Package Organization:**

```text
internal/domain/account/
├── account.go           # Domain code
└── account_test.go      # Unit tests

internal/infrastructure/repository/watch/
├── repository.go                      # Implementation
├── repository_test.go                 # Unit tests (mocked database)
└── repository_integration_test.go     # Integration tests (real database)
```

**Integration Test Tags:**

```go
//go:build integration

package repository_test

import "testing"

func TestRepository_Integration(t *testing.T) {
    // Integration test with real database
}
```

## Running Tests

**Unit Tests:**

```bash
make gotest
```

**Integration Tests:**

```bash
make gotest-integration
```

**Test Coverage:**

```bash
go test -cover ./...
```

**Verbose Output:**

```bash
go test -v ./...
```

## Test Utilities

The project provides test utilities in `pkg/testutil/`:

- `btc.go`: Bitcoin test utilities
- `eth.go`: Ethereum test utilities
- `xrp.go`: XRP test utilities
- `repository.go`: Repository test utilities
- `suite.go`: Test suite utilities

**Using Test Utilities:**

```go
import "github.com/hiromaily/go-crypto-wallet/pkg/testutil"

func TestSomething(t *testing.T) {
    // Use test utilities
    btcAddr := testutil.GenerateBTCTestAddress()
    // ...
}
```

## Table-Driven Tests

Use table-driven tests for multiple test cases:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "error case",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Test Best Practices

**Do:**

- Write tests for all exported functions
- Use table-driven tests for multiple cases
- Test both success and error paths
- Use descriptive test names
- Keep tests simple and focused
- Mock external dependencies in unit tests
- Use integration tests for end-to-end verification

**Don't:**

- Don't test implementation details
- Don't write flaky tests
- Don't skip error handling in tests
- Don't use sleeps for timing (use channels or mocks)
- Don't test private functions directly (test through public API)
- Don't write tests that depend on external state

## Test Coverage Goals

- **Domain Layer**: 80%+ coverage (pure business logic)
- **Application Layer**: 70%+ coverage (orchestration)
- **Infrastructure Layer**: 60%+ coverage (external dependencies)
- **Interface Adapters**: 70%+ coverage (user-facing logic)

**Note**: Coverage is a guideline, not a strict requirement. Focus on testing critical paths and business logic.

## See Also

- [Architecture Guidelines](architecture.md) - Layer structure and responsibilities
- [Coding Standards](coding-standards.md) - Code quality and verification commands
- [Workflow Guidelines](workflow.md) - Running tests in CI/CD workflow
- `docs/TESTING_STRATEGY.md` - Comprehensive testing strategy document
