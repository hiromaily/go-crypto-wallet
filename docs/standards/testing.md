# Testing Standards

Testing strategy and requirements for go-crypto-wallet.

## Test Commands

```bash
make gotest              # Run unit tests
make gotest-integration  # Run integration tests
```

## Test Organization

| Type | Location | Build Tag |
|------|----------|-----------|
| Unit tests | Same package (`*_test.go`) | None |
| Integration tests | Same package | `//go:build integration` |

## Testing by Layer

| Layer | Approach | Dependencies |
|-------|----------|--------------|
| Domain | Pure unit tests | None (no mocks needed) |
| Application | Unit tests with mocks | Mocked interfaces |
| Infrastructure | Integration tests | Real external systems |
| Interface Adapters | Integration tests | Full stack |

## Test Guidelines

### Unit Tests

- Test all exported functions
- Use table-driven tests for multiple cases
- Keep tests fast and deterministic
- No external dependencies

### Integration Tests

- Use build tag: `//go:build integration`
- Test with real databases/APIs when possible
- Clean up test data after tests

### Example: Table-Driven Test

```go
func TestValidateAddress(t *testing.T) {
    tests := []struct {
        name    string
        address string
        wantErr bool
    }{
        {"valid", "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", false},
        {"invalid", "invalid-address", true},
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

## Detailed Guidelines

See [docs/guidelines/testing.md](../guidelines/testing.md) for full testing strategy.
