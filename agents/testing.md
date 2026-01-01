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

**Approach**: Test with mocked infrastructure using [mockery](https://github.com/vektra/mockery)

**What to Test:**

- Use case input validation
- Service coordination and orchestration
- Error wrapping with context
- DTO transformation
- Business flow orchestration

**Example with Mocks:**

```go
package btc_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
    bitcoinmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/mocks"
    repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/mocks"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
    // Create mocks
    mockBtcClient := bitcoinmocks.NewMockBitcoiner(t)
    mockAddrRepo := repomocks.NewMockAddressRepositorier(t)

    // Set up expectations
    mockBtcClient.EXPECT().
        ListUnspentByAccount("deposit").
        Return(nil, nil)

    // Create use case with mocks
    useCase := btc.NewCreateTransactionUseCase(
        mockBtcClient,
        mockAddrRepo,
        // ... other dependencies
    )

    // Execute and verify
    result, err := useCase.Execute(ctx, params)
    require.NoError(t, err)
    // ... assertions
}
```

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

## Mock Generation with Mockery

This project uses [mockery v3](https://github.com/vektra/mockery) to generate mock implementations from Go interfaces.

### Configuration

Mock generation is configured in `.mockery.yaml` at the project root.

**Key Settings:**

- `all: false` - Only generate mocks for explicitly listed interfaces
- `template: testify` - Generate testify-compatible mocks with `EXPECT()` support
- Mocks are placed in `mocks/` subdirectories alongside implementations

### Mock Directory Structure

```text
internal/infrastructure/
├── api/bitcoin/
│   ├── btc/bitcoin.go              # Implementation
│   └── mocks/
│       └── mock_bitcoiner.go       # Generated mock
├── repository/
│   ├── watch/repository.go         # Implementation
│   └── mocks/
│       └── mock_*.go               # Generated mocks for persistence interfaces
└── storage/file/
    ├── transaction.go              # Implementation
    └── mocks/
        └── mock_transaction_file_repositorier.go
```

### Commands

```bash
# Generate all mocks
make mockery

# Clean all generated mocks
make clean-mocks

# Regenerate mocks (clean + generate)
make clean-mocks && make mockery
```

### Adding New Mock Interfaces

To add a new interface for mock generation:

1. Edit `.mockery.yaml`
2. Add the interface under the appropriate package:

```yaml
packages:
  github.com/hiromaily/go-crypto-wallet/internal/your/package:
    config:
      dir: "internal/your/package/mocks"
      pkgname: "mocks"
    interfaces:
      YourInterface:
```

3. Run `make mockery`

### Using Generated Mocks

```go
import (
    "testing"

    bitcoinmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/mocks"
    repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/mocks"
    storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/mocks"
)

func TestWithMocks(t *testing.T) {
    // Create mock (automatically registers cleanup with t.Cleanup)
    mockClient := bitcoinmocks.NewMockBitcoiner(t)

    // Set expectations with EXPECT()
    mockClient.EXPECT().
        GetBlockCount().
        Return(int64(100), nil)

    // Use mock in test
    result, err := mockClient.GetBlockCount()
    // Expectations are automatically verified at test end
}
```

### Mock Best Practices

**Do:**

- Pass `t *testing.T` to mock constructors for automatic cleanup
- Use `EXPECT()` for type-safe expectation setting
- Set expectations before calling the code under test
- Keep mock setups minimal and focused

**Don't:**

- Don't manually verify expectations (automatic with `t`)
- Don't create mocks without passing `t`
- Don't over-mock (mock only direct dependencies)

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
- Use mockery-generated mocks for infrastructure dependencies
- Use `EXPECT()` for type-safe mock expectations
- Pass `t *testing.T` to mock constructors
- Use integration tests for end-to-end verification

**Don't:**

- Don't test implementation details
- Don't write flaky tests
- Don't skip error handling in tests
- Don't use sleeps for timing (use channels or mocks)
- Don't test private functions directly (test through public API)
- Don't write tests that depend on external state
- Don't manually verify mock expectations (automatic with testify)
- Don't over-mock (only mock direct dependencies of the code under test)

## Test Coverage Goals

- **Domain Layer**: 80%+ coverage (pure business logic)
- **Application Layer**: 70%+ coverage (orchestration)
- **Infrastructure Layer**: 60%+ coverage (external dependencies)
- **Interface Adapters**: 70%+ coverage (user-facing logic)

**Note**: Coverage is a guideline, not a strict requirement. Focus on testing critical paths and business logic.

## See Also

- [Architecture Guidelines](architecture.md) - Layer structure and responsibilities
- [Coding Standards](coding-standards.md) - Code quality and verification commands
- [Code Generation](code-generation.md) - Mock generation and other code generation tools
- [Workflow Guidelines](workflow.md) - Running tests in CI/CD workflow
