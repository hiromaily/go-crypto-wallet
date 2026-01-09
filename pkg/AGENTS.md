# `pkg` Directory Guidelines

This document provides guidelines for AI agents working on packages in the `pkg/` directory.

## Overview

The `pkg/` directory contains **public packages** that can be imported by external code.
These packages provide shared utilities, configuration management, logging, and other common functionality.

## Directory Structure

The `pkg/` directory is organized as follows:

- `config/`: Configuration management utilities
  - `testutil/`: Test utilities for configuration
- `converter/`: Data conversion utilities
- `db/mysql/`: MySQL database connection utilities
- `debug/`: Debug utilities
- `decimal/`: Decimal number utilities
- `di/`: Legacy dependency injection container (for backward compatibility)
- `grpc/`: gRPC client utilities
- `logger/`: Logging utilities (structured logging, noop logger, slog support)
- `serial/`: Serialization utilities
- `testutil/`: Test utilities for various components
  - `btc.go`: Bitcoin test utilities
  - `eth.go`: Ethereum test utilities
  - `xrp.go`: XRP test utilities
  - `repository.go`: Repository test utilities
  - `suite.go`: Test suite utilities
- `uuid/`: UUID generation utilities
- `websocket/`: WebSocket client utilities

## Critical Rule: No Dependencies on `internal/`

**⚠️ THE MOST IMPORTANT RULE:**

**Packages in `pkg/` MUST NOT import or depend on any packages in `../internal/` directory.**

This is the fundamental architectural principle for the `pkg/` directory.

### Why This Rule Exists

- `pkg/` packages are **public APIs** that can be used by external code
- `internal/` contains **internal implementation details** following Clean Architecture
- Mixing these would break encapsulation and create circular dependencies
- External code should not need to know about internal architecture

## Common Commands

After making code changes, use these commands to verify code correctness:

- `make go-lint`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make gotest`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

## Patterns to Avoid

- ❌ Importing from `internal/` directory
- ❌ Using `panic` (except in `pkg/di` legacy code)
- ❌ Leaving commented-out code
- ❌ Unused imports, variables, or functions
- ❌ Ignoring errors (detected by `errcheck`)
- ❌ Logging sensitive information (private keys, passwords, etc.)

## Recommended Patterns

- ✅ Self-contained utility functions
- ✅ Clear, focused package responsibilities
- ✅ Proper error wrapping with context
- ✅ Use of `context.Context` for cancellation
- ✅ Comprehensive documentation
- ✅ Unit tests for all exported functions

## Important Notes

- This is a financial-related project; make changes carefully
- `pkg/` packages are public APIs; breaking changes affect external code
- Always verify that changes don't break existing functionality
- Consider backward compatibility when making changes
- **DO NOT** edit files that contain `DO NOT EDIT` comments
  (typically auto-generated files from tools like sqlc, protoc, or go generate)

## References

### Root Documentation

- **[Root AGENTS.md](../AGENTS.md)**: Overall project guidelines, navigation, and quick reference
- **[Core Principles](../docs/ai-agents/guidelines/core.md)**: Security, error handling, panic usage, and core patterns
- **[Coding Standards](../docs/ai-agents/guidelines/coding-standards.md)**: Linting, formatting, and code style
- **[Testing Guidelines](../docs/ai-agents/guidelines/testing.md)**: Testing strategy and requirements
- **[Workflow Guidelines](../docs/ai-agents/guidelines/workflow.md)**: Git operations and dependency management

### Directory Documentation

- **[Internal Guidelines](../internal/AGENTS.md)**: Guidelines for `internal/` directory
