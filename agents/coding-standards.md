# Coding Standards

This document describes the coding standards and conventions for the go-crypto-wallet project.

## Linting and Formatting

- Follow `golangci-lint` configuration (`.golangci.yml`)
- Format code with `make format` (uses `gofumpt` and `goimports` via golangci-lint)
  - Import order: standard → third-party → local
- Use `make lint-fix` to run linting and formatting together (executes lint checks and format fixes)
- Maintain consistent naming conventions (lowercase package names, exported functions start with uppercase)

## Common Commands

After making code changes, use these commands to verify code correctness:

- `make lint-fix`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

**Command Constraints**:

- **DO NOT** use `go build -v` directly; use `make check-build` instead
- **DO NOT** use `go tool golangci-lint` directly; use `make lint-fix` instead

## Naming Conventions

**Packages:**

- Use lowercase package names
- Avoid underscores in package names
- Keep package names short and descriptive

**Exported Symbols:**

- Start with uppercase letter
- Use camelCase (e.g., `GetAccountKey`, `CreateTransaction`)

**Unexported Symbols:**

- Start with lowercase letter
- Use camelCase (e.g., `calculateFee`, `validateAddress`)

**Constants:**

- Use camelCase for exported constants (e.g., `MaxRetries`)
- Use camelCase for unexported constants (e.g., `defaultTimeout`)

**Interfaces:**

- Name interfaces after the behavior they represent (e.g., `Validator`, `Repository`)
- For single-method interfaces, use the method name + "er" suffix (e.g., `Reader`, `Writer`)

## Import Organization

Organize imports in the following order:

1. Standard library packages
2. Third-party packages
3. Local packages (this project)

Example:

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third-party packages
    "github.com/btcsuite/btcd/btcutil"
    "github.com/pkg/errors"

    // Local packages
    "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
    "github.com/hiromaily/go-crypto-wallet/pkg/logger"
)
```

The `goimports` tool (via `make format`) will automatically organize imports in this order.

## Code Style

**Function Length:**

- Keep functions short and focused
- Prefer small, composable functions over large, complex ones
- If a function exceeds 50 lines, consider refactoring

**Error Handling:**

- Always check errors
- Wrap errors with context using `fmt.Errorf` with `%w`
- Return errors early (early return pattern)

**Comments:**

- Add godoc comments to all exported functions, methods, types, and constants
- Keep comments up-to-date with code changes
- Use complete sentences in comments
- Explain "why" rather than "what" in implementation comments

**Variable Naming:**

- Use short names for short-lived variables (e.g., `i` for loop index, `err` for errors)
- Use descriptive names for long-lived variables
- Avoid single-letter names except for:
  - Loop indices (`i`, `j`, `k`)
  - Short-lived variables in small scopes
  - Receivers (use consistent short names like `r` for `Repository`)

## Linter-Specific Guidelines

**unused-receiver:**

- For `unused-receiver` lint errors: **Remove the receiver entirely** instead of renaming it to `_`
- Renaming to `_` will only cause the same error to appear for other receivers
- Convert the method to a function from the start

**errcheck:**

- Never ignore errors
- If an error must be ignored, add a comment explaining why
- Use `_ = err` with a comment for intentionally ignored errors

**gofmt / gofumpt:**

- Always run `make format` before committing
- Consistent formatting helps with code reviews

**goimports:**

- Import organization is handled by `goimports`
- Run `make format` to organize imports automatically

## Testing Standards

See [Testing Guidelines](testing.md) for comprehensive testing standards.

**Basic Guidelines:**

- Write tests for all exported functions
- Use table-driven tests where appropriate
- Use `//go:build integration` tag for integration tests
- Keep unit tests fast and deterministic

## See Also

- [Core Principles](core.md) - Error handling, panic usage, and core patterns
- [Testing Guidelines](testing.md) - Comprehensive testing strategy
- [Workflow Guidelines](workflow.md) - Dependency management and verification commands
