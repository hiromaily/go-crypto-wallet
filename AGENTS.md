# Agent Guidelines for go-crypto-wallet

This document provides guidelines for AI agents working on this project.

## Project Overview

- This project is a cryptocurrency wallet implementation supporting BTC, BCH, ETH, XRP, and ERC-20 tokens
- **Primary language**: Go (backend, wallet operations)
- **Additional languages**: TypeScript/JavaScript (`apps/`), Solidity (smart contracts), SQL, HCL (Atlas schemas), Protobuf
- The project is currently under refactoring based on Clean Architecture and Clean Code principles
- Security is of utmost importance (private key management, offline wallets)
- The project follows the `pkg` layout pattern

## Quick Navigation

Use this navigation guide to find relevant documentation for specific tasks:

### Guidelines (Basic Guidelines)

- **[Core Principles](docs/ai-agents/guidelines/core.md)** - Security, error handling, panic usage, context management, logging, and core patterns
- **[Architecture Guidelines](docs/ai-agents/guidelines/architecture.md)** - Clean Architecture principles, layer separation, directory structure, and dependency direction
- **[Coding Standards](docs/ai-agents/guidelines/coding-standards.md)** - Linting, formatting, naming conventions, import organization, and code style
- **[Database Management](docs/ai-agents/guidelines/database.md)** - Database schema changes, Atlas migrations, and SQLC code generation
- **[Code Generation](docs/ai-agents/guidelines/code-generation.md)** - Auto-generated files, code generation tools (Atlas, SQLC, protobuf, ABI, mocks). **Important**: When moving code with mocks, update `.mockery.yaml` configuration.
- **[Workflow Guidelines](docs/ai-agents/guidelines/workflow.md)** - Git operations, dependency management, refactoring status, and verification commands
- **[Required Tools and Versions](docs/ai-agents/guidelines/requirements.md)** - Tool requirements and version information
- **[Testing Guidelines](docs/ai-agents/guidelines/testing.md)** - Testing strategy, test organization, and coverage goals by layer
- **[Multi-Chain Support](docs/ai-agents/guidelines/multi-chain.md)** - Cryptocurrency support (BTC, ETH, XRP), wallet types, and blockchain communication

### Task Contexts (Task-specific Context)

- **[Task-Oriented Context](docs/ai-agents/task-oriented-context.md)** - Task-oriented context management
- **[Task Contexts](docs/ai-agents/task-contexts/README.md)** - Task type-specific context list
- **[Verification Matrix](docs/ai-agents/task-contexts/verification.md)** - Verification commands by file type
- **[Task Analysis](docs/ai-agents/task-analysis.md)** - Issue/Commit pattern analysis

### Directory-Specific Guidelines

- **[`internal/` Directory Guidelines](internal/AGENTS.md)** - Clean Architecture layers, dependency rules, and internal package guidelines
- **[`pkg/` Directory Guidelines](pkg/AGENTS.md)** - Public package guidelines and critical rule about `internal/` dependencies

### Custom Commands and Agent Skills

- **[Agent Skills](docs/ai-agents/agent-skills.md)** - Modern Agent Skills format for AI assistants (github-issue-creation). Includes installation, usage guide, and examples.
- **[Custom Slash Commands](.claude/commands/README.md)** - Legacy custom commands for Claude Desktop (fix-issue, fix-pr-review, fix-linter). Note: create-github-issue has been migrated to Agent Skills.

## Core Priorities

1. **Security First**: This is a financial-related project handling private keys and cryptocurrency transactions. Security is non-negotiable.
2. **Clean Architecture**: Maintain clear layer separation (domain, application, infrastructure, interface-adapters) for Go code.
3. **Code Quality**: Follow language-specific linting standards, write tests, and maintain high code quality.
4. **Incremental Changes**: Make changes incrementally without breaking existing functionality.
5. **Language-Appropriate Patterns**: Use idiomatic patterns and best practices for each language (Go, TypeScript, Solidity, etc.).

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

- Run appropriate verification commands after making changes (see language-specific commands below)
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)
- Follow language-specific error handling patterns (e.g., Go: `fmt.Errorf` + `%w`)

#### Language-Specific Verification Commands

| Language | Format | Lint | Build | Test |
|----------|--------|------|-------|------|
| Go | `make go-fmt` | `make go-lint` | `make check-build` | `make gotest` |
| TypeScript/JS | `npm run format` | `npm run lint` | `npm run build` | `npm test` |
| Shell | `make shfmt` | - | - | - |
| Makefile | - | `make mk-lint` | - | - |

**Note**: For TypeScript/JavaScript apps in `apps/` directory, navigate to the specific app directory (e.g., `apps/ripple-lib-server/`) and run `npm install` before running npm commands.

### Git Operations

- **Allowed**: `git add`, `git commit`, `git push` to GitHub
- **NOT allowed**: `git merge` operations, `gh pr merge`, commits/pushes to `main` branch

## Common Workflows

### Making Code Changes

1. Create a feature branch
2. Make your changes following Clean Architecture principles
3. Run verification commands appropriate for the modified files:
   - **Go files**: `make go-lint`, `make tidy`, `make check-build`, `make gotest`
   - **TypeScript/JS files**: `npm run lint`, `npm run build`, `npm test` (in the relevant `apps/` subdirectory)
   - **Shell scripts**: `make shfmt`
   - **Makefiles**: `make mk-lint`
4. Commit changes with descriptive message
5. Push and create pull request

See [Workflow Guidelines](docs/ai-agents/guidelines/workflow.md) for detailed workflow documentation.

### Changing Database Schema

1. Modify HCL schema files in `tools/atlas/schemas/`
2. Format and lint: `make atlas-fmt` and `make atlas-lint`
3. Regenerate migrations: `make atlas-dev-reset`
4. Verify migration: `docker compose down -v && docker compose up`
5. Regenerate SQLC code: `make sqlc`
6. Verify build: `make check-build`

See [Database Management](docs/ai-agents/guidelines/database.md) for detailed database workflow.

### Adding New Use Case

1. Define use case interface in `internal/application/usecase/{wallet-type}/interfaces.go`
2. Implement use case in `internal/application/usecase/{wallet-type}/{coin}/`
3. Create constructor tests
4. Update CLI commands in `internal/interface-adapters/cli/` to use new use case
5. Wire up dependencies in `internal/di/`

See [Architecture Guidelines](docs/ai-agents/guidelines/architecture.md) for use case patterns and guidelines.

## When to Use Each Document

| Task | Document to Read |
|------|------------------|
| Understanding security requirements | [Core Principles](docs/ai-agents/guidelines/core.md) |
| Adding new business logic | [Architecture Guidelines](docs/ai-agents/guidelines/architecture.md) |
| Fixing linting issues (Go) | [Coding Standards](docs/ai-agents/guidelines/coding-standards.md) |
| Fixing linting issues (TS/JS) | Check `package.json` in relevant `apps/` directory |
| Changing database schema | [Database Management](docs/ai-agents/guidelines/database.md) |
| Working with auto-generated files | [Code Generation](docs/ai-agents/guidelines/code-generation.md) |
| Creating pull requests | [Workflow Guidelines](docs/ai-agents/guidelines/workflow.md) |
| Writing tests | [Testing Guidelines](docs/ai-agents/guidelines/testing.md) |
| Adding cryptocurrency support | [Multi-Chain Support](docs/ai-agents/guidelines/multi-chain.md) |
| Working in `internal/` packages | [Internal Guidelines](internal/AGENTS.md) |
| Working in `pkg/` packages | [Pkg Guidelines](pkg/AGENTS.md) |
| Working in `apps/` directory | Check `package.json` and `README.md` in relevant app |

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- When modifying files, use the appropriate linting/testing tools for that language

## Getting Help

If you need more detailed information about a specific topic, refer to the relevant document in the Quick Navigation section above. Each document provides focused, comprehensive guidelines for its domain.

## Multi-Language Project Structure

This project contains multiple languages with different tooling:

| Directory | Language | Tools |
|-----------|----------|-------|
| `internal/`, `pkg/`, `cmd/` | Go | golangci-lint, go test |
| `apps/ripple-lib-server/` | TypeScript | npm, eslint |
| `apps/erc20-token/` | JavaScript, Solidity | npm, truffle |
| `proto/` | Protobuf | protoc |
| `tools/atlas/schemas/` | HCL | atlas fmt, atlas lint |
| `tools/sqlc/` | SQL | sqlc |
| `scripts/` | Shell | shfmt |

When modifying files, use the appropriate verification tools for that language.

## AI Agent Documentation Structure

All AI-agent related documentation is consolidated under `docs/ai-agents/`:

```
docs/ai-agents/
├── guidelines/           # Basic Guidelines (formerly agents/)
│   ├── architecture.md
│   ├── coding-standards.md
│   ├── core.md
│   ├── database.md
│   ├── testing.md
│   ├── workflow.md
│   ├── code-generation.md
│   ├── multi-chain.md
│   └── requirements.md
├── task-contexts/        # Task-specific Context
│   ├── bug-fix.md
│   ├── feature-add.md
│   ├── refactoring.md
│   ├── db-change.md
│   ├── documentation.md
│   ├── verification.md
│   └── chains/           # Chain-specific
│       ├── btc.md
│       ├── bch.md
│       ├── eth.md
│       └── xrp.md
├── agent-skills.md       # Agent Skills Guide
├── task-oriented-context.md
└── task-analysis.md      # Pattern Analysis
```
