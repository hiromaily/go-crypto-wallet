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
3. **Single Source of Truth (SSOT)** - One authoritative location for each piece of information
4. **Incremental Changes** - No breaking changes without rollback plan
5. **Code Quality** - Follow language-specific linting and testing standards

## Expected Behavior

### Always Do

- **Check current branch before starting any task** (see `git-workflow` skill)
- Read relevant documentation before making changes
- Run verification commands after code changes
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Consider impact on offline wallet operations

### Never Do

- ❌ Log private keys or sensitive information
- ❌ Edit files marked `DO NOT EDIT` (auto-generated)
- ❌ Push directly to `main` branch
- ❌ Run `git merge` or `gh pr merge`
- ❌ Run `protoc` or `buf` commands directly (always use Makefile targets like `make proto`, `make proto-ts`)

### Ask Before

- Making security-related changes
- Breaking changes to public APIs
- Changes affecting multiple layers

## SSOT Structure

**When modifying rules, skills, or documentation, always edit the SSOT location.**

### AI Agent Configuration

| Category | SSOT Location | Other Locations |
|----------|---------------|-----------------|
| Rules | `.claude/rules/*.md` | `.cursor/rules/*.mdc` (auto-generated) |
| Skills | `.claude/skills/*/SKILL.md` | `.cursor/skills/` (symlink) |
| Commands | `.claude/commands/` | `.cursor/commands/` (reference only) |

**Sync Process:**

- `.cursor/rules/` → Auto-generated via `make sync-cursor-rules`
- `.cursor/skills/` → Symlink to `.claude/skills/`

### Project Documentation

| Category | SSOT Location | Notes |
|----------|---------------|-------|
| Standards | `docs/standards/` | Coding, testing, security, workflow |
| Guidelines | `docs/guidelines/` | Database, code-generation |
| Architecture | `ARCHITECTURE.md` | System design |
| Agent behavior | `AGENTS.md` (this file) | Entry point for all agents |

### Key Principle

> **Don't Repeat Yourself (DRY)**: Define once, reference everywhere.
> When information exists in multiple places, update the SSOT and reference it from others.

## Documentation Map

| Need | Document |
|------|----------|
| Project overview | [llms.txt](llms.txt) |
| Architecture design | [ARCHITECTURE.md](ARCHITECTURE.md) |
| **AI Agent instruction design** | [docs/design/ai-agents-instruction.md](docs/design/ai-agents-instruction.md) |
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
