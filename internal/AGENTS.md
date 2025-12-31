# `internal/` Directory Guidelines

This document provides guidelines for AI agents working on packages in the `internal/` directory.

## Overview

The `internal/` directory contains **internal application code** that follows Clean Architecture principles.
These packages are application-specific and are **NOT** intended for external use (unlike `pkg/` packages).

**Key Characteristics:**

- Application-specific code (not reusable libraries)
- Follows Clean Architecture with clear layer separation
- Organized by architectural layers (domain, application, infrastructure, interface-adapters)
- Implements dependency inversion (domain defines interfaces, infrastructure implements them)

## Architecture Layers

The `internal/` directory is organized into four main layers following Clean Architecture:

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Pure business logic with **ZERO infrastructure dependencies**

**Key Principles:**

- **NO dependencies** on infrastructure (no database, no API clients, no file I/O)
- Defines interfaces that infrastructure must implement (Dependency Inversion Principle)
- Contains pure business logic that can be tested without mocks
- Most stable layer - changes here affect all other layers

**Structure:**

- `account/`: Account types, validators, and business rules
- `transaction/`: Transaction types, state machine, validators
- `wallet/`: Wallet types and definitions
- `key/`: Key value objects and validators
- `multisig/`: Multisig validators and business rules
- `coin/`: Cryptocurrency type definitions

**When Working in Domain Layer:**

- ✅ Add business rules and validators
- ✅ Define value objects and entities
- ✅ Create domain interfaces for infrastructure to implement
- ✅ Write pure functions (no side effects)
- ❌ **NEVER** import from `infrastructure/`, `application/`, or `interface-adapters/`
- ❌ **NEVER** use database, API clients, or file I/O
- ❌ **NEVER** add logging (use `pkg/logger/` if absolutely necessary, but prefer no logging)

**Testing:**

- Use pure unit tests without mocks
- Test business logic in isolation
- No infrastructure dependencies required

### 2. Application Layer (`internal/application/`)

**Purpose**: Use case orchestration and business logic coordination

**Key Principles:**

- Orchestrates business logic by coordinating domain objects and infrastructure services
- Each use case represents a single business operation
- Depends on domain layer and infrastructure layer through interfaces
- Organized by wallet type (watch, keygen, sign) and cryptocurrency (btc, eth, xrp, shared)

**Structure:**

- `usecase/`: Use case implementations
  - `keygen/`: Key generation use cases
  - `sign/`: Signing use cases
  - `watch/`: Watch wallet use cases

**When Working in Application Layer:**

- ✅ Create use cases that orchestrate domain and infrastructure
- ✅ Transform DTOs between layers
- ✅ Wrap errors with context
- ✅ Depend on domain layer and infrastructure interfaces
- ❌ **NEVER** import from `interface-adapters/` (dependency direction violation)
- ❌ **NEVER** contain business logic (delegate to domain)
- ❌ **NEVER** directly access infrastructure implementations (use interfaces)

**Use Case Pattern:**

```go
type XxxUseCase interface {
    Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
}

type xxxUseCase struct {
    service ServiceInterface  // Domain or infrastructure interface
}

func (u *xxxUseCase) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    // Validate input using domain validators
    if err := domain.ValidateXxx(input); err != nil {
        return nil, fmt.Errorf("invalid input: %w", err)
    }
    
    // Delegate to service
    result, err := u.service.SomeMethod(input.Param1)
    if err != nil {
        return nil, fmt.Errorf("failed to execute xxx: %w", err)
    }
    
    return &XxxOutput{Result: result}, nil
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: External dependencies and technical implementations

**Key Principles:**

- Implements interfaces defined by domain layer (Dependency Inversion Principle)
- Contains NO business logic (only technical implementation)
- Easily replaceable and mockable for testing
- Organized by technical concern (database, api, repository, storage, network)

**Structure:**

- `api/`: External API clients (Bitcoin, Ethereum, Ripple)
- `database/`: Database connections and generated code (MySQL, sqlc)
- `repository/`: Data persistence implementations
- `storage/`: File storage implementations
- `network/`: Network communication (WebSocket, gRPC)
- `contract/`: Smart contract utilities
- `wallet/key/`: Key generation logic (infrastructure concern)
- `config/`: Configuration management

**When Working in Infrastructure Layer:**

- ✅ Implement domain interfaces
- ✅ Handle external system communication
- ✅ Convert between domain entities and external formats
- ✅ Manage technical concerns (database, network, file I/O)
- ❌ **NEVER** contain business logic (delegate to domain)
- ❌ **NEVER** validate business rules (use domain validators)
- ❌ **NEVER** import from `application/` or `interface-adapters/`
- ❌ **NEVER** make domain decisions

**Infrastructure Component Guidelines:**

- **Database**: Connection management, query execution, transaction management
- **Repository**: CRUD operations, convert between domain entities and database models
- **API Clients**: Communicate with external blockchain APIs, handle network errors
- **Storage**: File I/O for transaction and address data
- **Network**: Connection management (WebSocket, gRPC)

### 4. Interface Adapters Layer (`internal/interface-adapters/`)

**Purpose**: Adapters between use cases and external interfaces

**Key Principles:**

- Converts between external formats and application DTOs
- Depends ONLY on application layer (use cases)
- Provides adapters for CLI, HTTP, and wallet interfaces

**Structure:**

- `cli/`: CLI command adapters (keygen, sign, watch)
- `wallet/`: Wallet adapter interfaces and implementations
- `http/`: HTTP handlers (if applicable)

**When Working in Interface Adapters Layer:**

- ✅ Call use cases from commands
- ✅ Convert CLI arguments to use case DTOs
- ✅ Convert use case outputs to CLI output
- ✅ Handle user-facing errors and formatting
- ❌ **NEVER** import from `infrastructure/` or `domain/` directly
- ❌ **NEVER** contain business logic (delegate to use cases)
- ❌ **NEVER** call services directly (use use cases)

**Command Pattern:**

```go
type CreateCommand struct {
    useCase CreateUseCase  // Application layer use case
}

func (c *CreateCommand) Run(args []string) error {
    // Convert CLI args to use case input
    input := CreateInput{
        Param1: args[0],
        Param2: args[1],
    }
    
    // Call use case
    output, err := c.useCase.Execute(context.Background(), input)
    if err != nil {
        return fmt.Errorf("failed to create: %w", err)
    }
    
    // Format output for CLI
    fmt.Printf("Created: %s\n", output.Result)
    return nil
}
```

## Dependency Direction Rules

**CRITICAL**: Dependencies must flow in ONE direction only:

```text
Interface Adapters → Application Layer → Domain Layer ← Infrastructure Layer
```

**Rules:**

1. **Domain Layer**: NO dependencies on other layers
2. **Application Layer**: Can depend on Domain Layer only
3. **Infrastructure Layer**: Can depend on Domain Layer only (implements domain interfaces)
4. **Interface Adapters**: Can depend on Application Layer only

**Violations to Avoid:**

- ❌ Domain importing from infrastructure
- ❌ Application importing from infrastructure directly (use interfaces)
- ❌ Application importing from interface-adapters
- ❌ Infrastructure importing from application
- ❌ Interface-adapters importing from infrastructure or domain directly

## `internal/` vs `pkg/` Directory

**`internal/` Directory:**

- Application-specific code
- Follows Clean Architecture layers
- NOT intended for external use
- Can import from `pkg/` (shared utilities)

**`pkg/` Directory:**

- Shared, reusable packages
- Can be imported by external code
- **MUST NOT** import from `internal/` (critical rule)
- Contains utilities (logger, config, converter, etc.)

**When to Use Each:**

- Use `internal/` for application-specific business logic
- Use `pkg/` for shared utilities that could be used by external code
- If unsure, prefer `internal/` (can be refactored to `pkg/` later if needed)

## Legacy Code (`internal/wallet/service/`)

**Status**: Legacy/transitional code being refactored to use case layer

**Guidelines:**

- New code should use `internal/application/usecase/` instead
- Legacy services are being gradually replaced by use cases
- When refactoring, move logic to appropriate layer:
  - Business logic → `domain/`
  - Orchestration → `application/usecase/`
  - Technical implementation → `infrastructure/`

## Dependency Injection (`internal/di/`)

**Purpose**: Wire up dependencies between layers

**Guidelines:**

- `panic` is allowed here (instance construction phase)
- Wire use cases, services, repositories, and API clients
- Follow dependency direction rules
- Keep wiring logic separate from business logic

## Common Patterns

### Adding New Business Logic

1. **Check if it belongs in Domain Layer** (business rules, validators)
2. **If orchestration needed**, create a use case in Application Layer
3. **If external system interaction**, implement in Infrastructure Layer
4. **If user interface**, create adapter in Interface Adapters Layer

### Adding New Cryptocurrency Support

1. **Domain**: Add coin type, validators, business rules
2. **Application**: Add use cases for new coin
3. **Infrastructure**: Add API client, repository implementations
4. **Interface Adapters**: Add CLI commands, wallet adapters

### Refactoring Legacy Code

1. Identify the layer the code belongs to
2. Move business logic to domain layer
3. Create use cases in application layer
4. Move technical implementation to infrastructure layer
5. Update interface adapters to use new use cases
6. Update dependency injection

## Testing Strategy

**Domain Layer:**

- Pure unit tests without mocks
- Test business logic in isolation
- Fast, deterministic tests

**Application Layer:**

- Test use cases with mocked infrastructure
- Verify orchestration logic
- Test error handling and DTO transformation

**Infrastructure Layer:**

- Unit tests with mocked external dependencies
- Integration tests with real external systems
- Test error handling and retries

**Interface Adapters Layer:**

- Test command parsing and output formatting
- Test use case integration (with mocked use cases)
- Test error handling and user-facing messages

## Common Commands

After making code changes in `internal/`, use these commands:

- `make lint-fix`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

## Important Notes

- This is a financial-related project; make changes carefully
- Follow Clean Architecture principles strictly
- Maintain dependency direction rules
- Security is critical (private key management, offline wallets)
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)

## References

### Root Documentation

- **[Root AGENTS.md](../AGENTS.md)**: Overall project guidelines, navigation, and quick reference
- **[Core Principles](../agents/core.md)**: Security, error handling, panic usage, and core patterns
- **[Architecture Guidelines](../agents/architecture.md)**: Clean Architecture principles and layer guidelines
- **[Coding Standards](../agents/coding-standards.md)**: Linting, formatting, and code style
- **[Database Management](../agents/database.md)**: Database schema changes and SQLC workflow
- **[Code Generation](../agents/code-generation.md)**: Auto-generated files and code generation tools
- **[Workflow Guidelines](../agents/workflow.md)**: Git operations and dependency management
- **[Testing Guidelines](../agents/testing.md)**: Testing strategy and requirements

### Directory Documentation

- **[Pkg Guidelines](../pkg/AGENTS.md)**: Guidelines for `pkg/` directory
- `internal/domain/doc.go`: Domain layer documentation
- `internal/infrastructure/doc.go`: Infrastructure layer documentation
- `internal/application/doc.go`: Application layer documentation
