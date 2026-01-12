# Agent Guidelines for go-crypto-wallet

This document defines the **behavior and values** for AI agents working on this project.
For detailed documentation, see [llms.txt](llms.txt) and [ARCHITECTURE.md](ARCHITECTURE.md).

## Project Identity

- **Type**: Multi-signature cryptocurrency wallet (BTC, BCH, ETH, XRP, ERC-20)
- **Security Model**: Offline cold wallets (keygen, sign) + Online watch wallet
- **Architecture**: Clean Architecture with strict layer separation
- **Status**: Under refactoring based on Clean Code principles

## Core Values (Priority Order)

1. **Security First** - Private key protection is non-negotiable
2. **Clean Architecture** - Domain layer has ZERO infrastructure dependencies
3. **Incremental Changes** - No breaking changes without rollback plan
4. **Code Quality** - Follow language-specific linting and testing standards

## Expected Behavior

### Always Do

- Read relevant documentation before making changes
- Run verification commands after code changes
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Consider impact on offline wallet operations

### Never Do

- ❌ Log private keys or sensitive information
- ❌ Edit files marked `DO NOT EDIT` (auto-generated)
- ❌ Push directly to `main` branch
- ❌ Run `git merge` or `gh pr merge`

### Ask Before

- Making security-related changes
- Breaking changes to public APIs
- Changes affecting multiple layers

## Documentation Map

| Need | Document |
|------|----------|
| Project overview | [llms.txt](llms.txt) |
| Architecture design | [ARCHITECTURE.md](ARCHITECTURE.md) |
| **Standards (SSOT)** | [docs/standards/](docs/standards/) |
| Coding conventions | [docs/standards/coding-conventions.md](docs/standards/coding-conventions.md) |
| Security | [docs/standards/security.md](docs/standards/security.md) |
| Testing | [docs/standards/testing.md](docs/standards/testing.md) |
| Workflow | [docs/standards/workflow.md](docs/standards/workflow.md) |
| Database changes | [docs/guidelines/database.md](docs/guidelines/database.md) |
| Auto-generated files | [docs/guidelines/code-generation.md](docs/guidelines/code-generation.md) |
| Internal packages | [internal/AGENTS.md](internal/AGENTS.md) |
| Public packages | [pkg/AGENTS.md](pkg/AGENTS.md) |

## Quick Reference

### Verification Commands

| Language | Lint | Build | Test |
|----------|------|-------|------|
| Go | `make go-lint` | `make check-build` | `make gotest` |
| TypeScript | `npm run lint` | `npm run build` | `npm test` |
| Shell | `make shfmt` | - | - |

### Git Operations

```bash
# Allowed
git add, git commit, git push

# NOT Allowed
git merge, gh pr merge, push to main
```

## See Also

- [llms.txt](llms.txt) - AI-friendly project sitemap
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [docs/standards/](docs/standards/) - Project standards (SSOT)
- [docs/guidelines/](docs/guidelines/) - Detailed guidelines
