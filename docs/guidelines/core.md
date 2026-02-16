# Core Principles

This document contains immutable core principles that apply across all domains of the go-crypto-wallet project.

## Security

Security is of utmost importance in this financial-related project.

**Critical Rules:**

- **NEVER** log private keys or sensitive information
- Encrypt or zero-clear private keys in memory when possible
- Do not pass passwords via command-line arguments; use secure input methods
- Conduct security review when making changes involving sensitive information
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)

**When to Run Security Scans:**

- For security-related changes, run `make go-check-vuln`
- Review encryption/decryption logic carefully
- Extra caution for private key management and wallet operations

## Error Handling

Proper error handling is essential for robustness and debugging.

**Rules:**

- Wrap errors with `fmt.Errorf` + `%w` to preserve error chains
- Use `errors.Is`/`errors.As` for error checking
- Include context information in error messages
- Never ignore errors (detected by `errcheck` linter)

**Example:**

```go
result, err := service.DoSomething()
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}
```

## Panic Usage

Following the principle of separation of concerns, the project clearly separates the instance construction phase from the instance usage phase. Therefore, `panic` is only allowed during instance construction.

**Where `panic` is Acceptable:**

- `main.go` files (application entry points)
- `internal/di` package (dependency injection container)
- `pkg/di` package (legacy dependency injection container - for backward compatibility)

**Important:**

- `panic` should **NOT** be used in business logic, service layers, or infrastructure layers
- Use proper error handling with error returns in all other packages
- The separation allows for fail-fast behavior during initialization while maintaining robust error handling during runtime

## Context Management

Context management is crucial for proper request handling, timeouts, and cancellation.

**Rules:**

- Add `context.Context` to all API calls
- Implement timeouts and cancellation
- Implement graceful shutdown
- Pass context through all layers (use cases, services, repositories)

## Logging

Structured logging helps with debugging and monitoring.

**Rules:**

- Use structured logging (via `pkg/logger/`)
- Set appropriate log levels
- **NEVER** log sensitive information (private keys, passwords, etc.)
- Avoid excessive logging in hot paths

## Patterns to Avoid

These patterns should be avoided across the entire codebase:

- Using `log.Fatal` (except in `main`)
- Using `panic` outside of instance construction (i.e., outside of `main.go` and `pkg/di` package)
- Leaving commented-out code
- Unused imports, variables, or functions
- Ignoring errors (detected by `errcheck`)
- Using `go build -v` directly (use `make check-build` instead)
- Using `go tool golangci-lint` directly (use `make go-lint` instead)
- For `unused-receiver` lint errors: **Remove the receiver entirely** instead of renaming it to `_`. Renaming to `_` will only cause the same error to appear for other receivers, so convert the method to a function from the start.

## Recommended Patterns

These patterns should be used throughout the codebase:

- Abstraction through interfaces
- Dependency injection
- Proper error wrapping with context
- Use of `context.Context`
- Clear, focused package responsibilities
- Self-contained utility functions
- Comprehensive documentation
- Unit tests for all exported functions

## Documentation

Good documentation makes the codebase more maintainable.

**Rules:**

- Add godoc comments to exported functions and methods
- Add package-level comments
- Include usage examples for complex logic
- Keep documentation up-to-date with code changes

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)

## See Also

- [Architecture Guidelines](architecture.md) - Clean Architecture principles and layer guidelines
- [Coding Standards](coding-conventions.md) - Linting, formatting, and naming conventions
- [Workflow Guidelines](workflow.md) - Git operations and dependency management
