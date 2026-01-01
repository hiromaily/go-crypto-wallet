# Agent Guidelines for go-crypto-wallet

This document provides guidelines for AI agents working on this project.

## Project Overview

- This project is a cryptocurrency wallet implementation in Go supporting BTC, BCH, ETH, XRP, and ERC-20 tokens
- The project is currently under refactoring based on Clean Architecture and Clean Code principles
- Security is of utmost importance (private key management, offline wallets)
- The project follows the `pkg` layout pattern

## Quick Navigation

Use this navigation guide to find relevant documentation for specific tasks:

- **[Core Principles](agents/core.md)** - Security, error handling, panic usage, context management, logging, and core patterns
- **[Architecture Guidelines](agents/architecture.md)** - Clean Architecture principles, layer separation, directory structure, and dependency direction
- **[Coding Standards](agents/coding-standards.md)** - Linting, formatting, naming conventions, import organization, and code style
- **[Database Management](agents/database.md)** - Database schema changes, Atlas migrations, and SQLC code generation
- **[Code Generation](agents/code-generation.md)** - Auto-generated files, code generation tools (Atlas, SQLC, protobuf, ABI)
- **[Workflow Guidelines](agents/workflow.md)** - Git operations, dependency management, refactoring status, and verification commands
- **[Testing Guidelines](agents/testing.md)** - Testing strategy, test organization, and coverage goals by layer
- **[Multi-Chain Support](agents/multi-chain.md)** - Cryptocurrency support (BTC, ETH, XRP), wallet types, and blockchain communication

### Directory-Specific Guidelines

- **[`internal/` Directory Guidelines](internal/AGENTS.md)** - Clean Architecture layers, dependency rules, and internal package guidelines
- **[`pkg/` Directory Guidelines](pkg/AGENTS.md)** - Public package guidelines and critical rule about `internal/` dependencies

### Custom Commands

- **[Custom Slash Commands](.claude/commands/README.md)** - Custom commands for Claude Desktop (fix-issue, fix-pr-review, fix-linter, create-github-issue)

## Core Priorities

1. **Security First**: This is a financial-related project handling private keys and cryptocurrency transactions. Security is non-negotiable.
2. **Clean Architecture**: Maintain clear layer separation (domain, application, infrastructure, interface-adapters).
3. **Code Quality**: Follow linting standards, write tests, and maintain high code quality.
4. **Incremental Changes**: Make changes incrementally without breaking existing functionality.

## Essential Rules

### Security

- **NEVER** log private keys or sensitive information
- Always conduct security review for changes involving sensitive information
- Consider the impact on offline wallet operations (keygen, sign)

### Architecture

- Follow Clean Architecture principles with clear layer separation
- Domain layer has **ZERO infrastructure dependencies**
- Use dependency injection and abstract with interfaces

### Code Quality

- Run verification commands after making changes: `make lint-fix`, `make tidy`, `make check-build`, `make gotest`
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)
- Wrap errors with `fmt.Errorf` + `%w` for proper error chains

### Git Operations

- **Allowed**: `git add`, `git commit`, `git push` to GitHub
- **NOT allowed**: `git merge` operations, `gh pr merge`, commits/pushes to `main`/`master` branches

## Common Workflows

### Making Code Changes

1. Create a feature branch
2. Make your changes following Clean Architecture principles
3. Run verification commands:
   - `make lint-fix` - Fix linting issues
   - `make tidy` - Organize dependencies
   - `make check-build` - Verify build
   - `make gotest` - Run tests
4. Commit changes with descriptive message
5. Push and create pull request

See [Workflow Guidelines](agents/workflow.md) for detailed workflow documentation.

### Changing Database Schema

1. Modify HCL schema files in `tools/atlas/schemas/`
2. Format and lint: `make atlas-fmt` and `make atlas-lint`
3. Regenerate migrations: `make atlas-dev-reset`
4. Verify migration: `docker compose down -v && docker compose up`
5. Regenerate SQLC code: `make sqlc`
6. Verify build: `make check-build`

See [Database Management](agents/database.md) for detailed database workflow.

### Adding New Use Case

1. Define use case interface in `internal/application/usecase/{wallet-type}/interfaces.go`
2. Implement use case in `internal/application/usecase/{wallet-type}/{coin}/`
3. Create constructor tests
4. Update CLI commands in `internal/interface-adapters/cli/` to use new use case
5. Wire up dependencies in `internal/di/`

See [Architecture Guidelines](agents/architecture.md) for use case patterns and guidelines.

## When to Use Each Document

| Task | Document to Read |
|------|------------------|
| Understanding security requirements | [Core Principles](agents/core.md) |
| Adding new business logic | [Architecture Guidelines](agents/architecture.md) |
| Fixing linting issues | [Coding Standards](agents/coding-standards.md) |
| Changing database schema | [Database Management](agents/database.md) |
| Working with auto-generated files | [Code Generation](agents/code-generation.md) |
| Creating pull requests | [Workflow Guidelines](agents/workflow.md) |
| Writing tests | [Testing Guidelines](agents/testing.md) |
| Adding cryptocurrency support | [Multi-Chain Support](agents/multi-chain.md) |
| Working in `internal/` packages | [Internal Guidelines](internal/AGENTS.md) |
| Working in `pkg/` packages | [Pkg Guidelines](pkg/AGENTS.md) |

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)

## Getting Help

If you need more detailed information about a specific topic, refer to the relevant document in the Quick Navigation section above. Each document provides focused, comprehensive guidelines for its domain.
