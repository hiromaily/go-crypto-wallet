---
description: Interface Adapters Layer Rules
globs: ["internal/interface-adapters/**/*.go"]
alwaysApply: false
---


# Interface Adapters Layer Rules

## Overview

Rules for working with the Interface Adapters layer (`internal/interface-adapters/`).
This layer adapts external interfaces (CLI, HTTP) to application use cases.

## Structure

```
internal/interface-adapters/
├── cli/            # CLI command adapters
│   ├── app/        # App initialization, flags, validation
│   ├── keygen/     # Key generation commands
│   │   ├── api/    # API commands (btc, eth)
│   │   ├── create/ # Creation commands
│   │   ├── export/ # Export commands
│   │   └── imports/# Import commands
│   ├── sign/       # Sign wallet commands
│   └── watch/      # Watch wallet commands
│       ├── api/    # API commands (btc, eth, xrp)
│       ├── create/ # Creation commands
│       ├── imports/# Import commands
│       ├── monitor/# Monitor commands
│       └── send/   # Send commands
├── http/           # HTTP handlers and middleware
└── wallet/         # Wallet adapters
    ├── btc/        # Bitcoin (keygen, sign, watch)
    ├── eth/        # Ethereum (keygen, sign, watch)
    └── xrp/        # XRP (keygen, sign, watch)
```

## Key Principle

**Interface Adapters depend ONLY on the Application layer.**

- Call use cases from `application/usecase/`
- Convert external input (CLI args, HTTP requests) to use case input
- Convert use case output to external output (CLI output, HTTP responses)

## Command Pattern

```go
type CreateCommand struct {
    useCase usecase.CreateUseCase  // Application layer use case
}

func (c *CreateCommand) Run(args []string) error {
    // 1. Convert CLI args to use case input
    input := usecase.CreateInput{
        Param1: args[0],
    }

    // 2. Call use case
    output, err := c.useCase.Execute(context.Background(), input)
    if err != nil {
        return fmt.Errorf("failed to create: %w", err)
    }

    // 3. Format output for CLI
    fmt.Printf("Created: %s\n", output.Result)
    return nil
}
```

## Allowed/Forbidden

| Action | Allowed |
|--------|---------|
| Import from `application/usecase/` | Yes |
| Import from `application/dto/` | Yes |
| Import from `domain/` | No |
| Import from `infrastructure/` | No |
| Contain business logic | No (delegate to use cases) |
| Call infrastructure directly | No (use use cases) |

## Guidelines

### CLI Commands

- Parse and validate CLI arguments
- Convert to use case DTOs
- Call use cases
- Format output for terminal

### HTTP Handlers

- Parse HTTP requests
- Call use cases
- Return appropriate HTTP responses
- Handle errors with proper status codes

### Wallet Adapters

- Provide wallet-specific interfaces
- Delegate to appropriate use cases
- Handle wallet type differences

## Related Rules

- @.claude/rules/internal/clean-architecture.md - Layer dependencies
- @.claude/rules/internal/application-layer.md - Use case interfaces
- @.claude/rules/go/conventions.md - Go conventions
