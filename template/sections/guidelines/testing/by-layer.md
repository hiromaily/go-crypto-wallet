### Testing by Layer

#### Domain Layer Testing

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
import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

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
            if tt.wantErr {
                require.Error(t, err, "Validate() should return error for %v", tt.input)
            } else {
                assert.NoError(t, err, "Validate() should not return error for %v", tt.input)
            }
        })
    }
}
```

#### Application Layer Testing (Use Cases)

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
    bitcoinmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/mocks"
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

#### Infrastructure Layer Testing

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

#### Interface Adapters Layer Testing

**Approach**: Test with mocked use cases

**What to Test:**

- Command argument parsing
- Output formatting
- Error message formatting
- CLI flag handling
- Use case integration (with mocked use cases)
