### Application Layer (Use Case) Guidelines

The `internal/application/usecase/` package implements the use case layer following Clean Architecture principles.

**Key Principles:**

- Use cases orchestrate business logic by coordinating domain objects and infrastructure services
- Each use case represents a single business operation with clear input and output
- Use cases act as thin wrappers that transform DTOs, delegate to services, and wrap errors with context
- **Use cases depend on interfaces defined in `application/ports/`** (Dependency Inversion)
- Organized by wallet type (watch, keygen, sign) and cryptocurrency (btc, eth, xrp, shared)

**Use Case Structure:**

```go
// Use case interface definition
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

// Input/Output DTOs
type XxxInput struct {
    Param1 string
    Param2 int
}

type XxxOutput struct {
    Result string
}

// Implementation
type xxxUseCase struct {
    service ServiceInterface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    result, err := u.service.SomeMethod(input.Param1, input.Param2)
    if err != nil {
        return nil, fmt.Errorf("failed to execute xxx: %w", err)
    }
    return &XxxOutput{Result: result}, nil
}
```

**DTO Conventions:**

- **Input DTOs**: Contain all parameters needed for the use case operation
- **Output DTOs**: Contain all results returned by the use case
- DTOs use domain types (not primitive types when domain types exist)
- DTOs are passed by value for inputs, returned as pointers for outputs

**Error Handling:**

- Wrap service errors with context using `fmt.Errorf` with `%w`
- Error messages should describe the use case operation that failed
- Return domain errors when business rule violations occur
- Let infrastructure errors propagate with added context

**Organization Structure:**

```text
internal/application/usecase/
├── keygen/
│   ├── interfaces.go              # Use case interfaces
│   ├── btc/                       # Bitcoin-specific use cases
│   ├── eth/                       # Ethereum-specific use cases
│   ├── xrp/                       # XRP-specific use cases
│   └── shared/                    # Shared use cases (all coins)
├── sign/
│   ├── interfaces.go
│   ├── btc/
│   ├── eth/
│   ├── xrp/
│   └── shared/
└── watch/
    ├── interfaces.go
    ├── btc/
    ├── eth/
    ├── xrp/
    └── shared/
```

**Testing Approach:**

Use cases currently have constructor tests that verify:

- Use case can be instantiated with dependencies
- Correct interface implementation

For comprehensive testing strategy, see [Testing Guidelines](../../../docs/guidelines/testing.md).

**When to Create a New Use Case:**

- New command functionality is added (commands should use use cases, not services directly)
- Existing service logic needs to be exposed to commands with different DTO structure
- Business logic needs to coordinate multiple services
- Transaction boundaries need to be defined

**Important:**

- Commands in `internal/interface-adapters/cli/` should ONLY depend on use cases, NOT services directly
- Use cases should be small and focused on a single operation
- Avoid business logic in use cases; delegate to domain or services
- Use cases are the entry point to application logic from command layer
