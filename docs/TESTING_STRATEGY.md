# Testing Strategy for Use Case Layer

## Current Status (Phase 8)

As part of the Clean Architecture refactoring (#69), we've introduced a use case layer that wraps service implementations. This document outlines the testing strategy for use cases.

## Testing Approach

### Use Case Layer Characteristics

The use case layer currently acts as thin wrappers around service implementations with these responsibilities:
- Input/output DTO transformation
- Error wrapping with context
- Consistent interface across wallet types

Most use cases follow this pattern:
```go
func (u *someUseCase) Execute(ctx context.Context, input Input) (Output, error) {
    result, err := u.service.SomeMethod(input.Param)
    if err != nil {
        return Output{}, fmt.Errorf("failed to ...: %w", err)
    }
    return Output{Result: result}, nil
}
```

### Current Test Coverage

**Constructor Tests** ✅
- Verify use case can be instantiated
- Verify correct interface implementation
- Located in `*_test.go` files alongside implementations

Example test files:
- `pkg/application/usecase/watch/shared/import_address_test.go`
- `pkg/application/usecase/keygen/shared/generate_seed_test.go`
- `pkg/application/usecase/sign/shared/store_seed_test.go`

### Testing Limitations

**Service Dependencies**
Most services are concrete types, not interfaces:
```go
// Current (concrete type)
type useCaseImpl struct {
    service *concreteService.Service
}

// Ideal for testing (interface)
type useCaseImpl struct {
    service ServiceInterface
}
```

This makes mocking difficult without significant refactoring.

**Current Test Scope**
- ✅ Constructor validation
- ✅ Interface compliance
- ❌ Method behavior (requires mocks or integration tests)
- ❌ Error handling paths (requires mocks or integration tests)

## Future Improvements

### Phase 1: Service Interface Extraction

Refactor services to use interfaces:

```go
// In pkg/wallet/service/interfaces.go
type Seeder interface {
    Generate() ([]byte, error)
    Store(seed string) ([]byte, error)
}

// Use cases depend on interface
type generateSeedUseCase struct {
    seeder service.Seeder  // interface instead of *concrete.Service
}
```

**Benefits:**
- Enable proper unit testing with mocks
- Improve testability without changing behavior
- Better dependency inversion

**Effort:** Medium (affects ~20 services)

### Phase 2: Mock Implementations

Create mock implementations for testing:

```go
type mockSeeder struct {
    generateFunc func() ([]byte, error)
    storeFunc    func(seed string) ([]byte, error)
}
```

**Tools to consider:**
- `github.com/stretchr/testify/mock`
- `github.com/golang/mock`
- Manual mocks (current approach)

### Phase 3: Integration Tests

Add integration tests that test full flows:
- Require test database setup
- Test file fixtures
- Proper cleanup

Location: `pkg/application/usecase/integration_test.go` or similar

## Testing Guidelines

### When to Write Tests

**Always:**
- New use cases should have constructor tests
- Changes to use case logic require test updates

**Consider:**
- Integration tests for critical paths
- Mock-based tests when services become interfaces

**Skip:**
- Extensive unit tests for thin wrappers (low ROI)
- Tests that duplicate service layer tests

### Test Organization

```
pkg/application/usecase/
├── watch/
│   ├── shared/
│   │   ├── import_address.go
│   │   └── import_address_test.go      # Constructor tests
│   ├── btc/
│   │   ├── create_transaction.go
│   │   └── create_transaction_test.go  # Constructor tests
│   └── integration_test.go             # Future: integration tests
```

### Running Tests

```bash
# Run all use case tests
go test ./pkg/application/usecase/...

# Run specific package
go test ./pkg/application/usecase/watch/shared/...

# With coverage
go test -cover ./pkg/application/usecase/...
```

## Migration Path

1. ✅ **Phase 8 (Current)**: Establish testing pattern with constructor tests
2. **Future**: Extract service interfaces
3. **Future**: Add comprehensive unit tests with mocks
4. **Future**: Add integration test suite

## Related Issues

- #69 - Use Case Layer Implementation
- Future: Service Interface Extraction (TBD)

## Questions?

For questions about testing strategy, refer to:
- This document
- `AGENTS.md` for architecture guidelines
- `CLAUDE.md` for coding standards
