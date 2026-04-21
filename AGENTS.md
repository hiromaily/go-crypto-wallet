<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/AGENTS.tpl.md · Run `make docs` to regenerate.
-->

# Agent Guidelines for go-crypto-wallet

This document defines the **behavior and values** for AI agents working on this project.
For detailed documentation, see [llms.txt](./llms.txt) and [ARCHITECTURE.md](./ARCHITECTURE.md).

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
| Guidelines | `docs/guidelines/` | Coding, testing, security, workflow, database, code-generation |
| Architecture | `ARCHITECTURE.md` | System design |
| Agent behavior | `AGENTS.md` (this file) | Entry point for all agents |

### Key Principle

> **Don't Repeat Yourself (DRY)**: Define once, reference everywhere.
> When information exists in multiple places, update the SSOT and reference it from others.

## Documentation Map

| Need | Document |
|------|----------|
| Project overview | [llms.txt](./llms.txt) |
| Architecture design | [ARCHITECTURE.md](./ARCHITECTURE.md) |
| **AI Agent instruction design** | [docs/ai/design.md](./docs/ai/design.md) |
| **Guidelines** | [docs/guidelines/](./docs/guidelines) |
| Coding conventions | [docs/guidelines/coding-conventions.md](./docs/guidelines/coding-conventions.md) |
| Security | [docs/guidelines/security.md](./docs/guidelines/security.md) |
| Testing | [docs/guidelines/testing.md](./docs/guidelines/testing.md) |
| Workflow | [docs/guidelines/workflow.md](./docs/guidelines/workflow.md) |
| Database changes | [docs/database/db-management.md](./docs/database/db-management.md) |
| Auto-generated files | [docs/guidelines/code-generation.md](./docs/guidelines/code-generation.md) |
| CLI commands | [internal/interface-adapters/cli/README.md](./internal/interface-adapters/cli/README.md) | Command × Chain × UseCase matrix (SSOT) |
| Internal packages | [internal/AGENTS.md](./internal/AGENTS.md) |
| Public packages | [pkg/AGENTS.md](./pkg/AGENTS.md) |

## Quick Reference

### Identifying BTC Address Types

| Prefix | Type | BIP | SegWit |
|--------|------|-----|--------|
| `1...` | P2PKH | BIP44 | ❌ |
| `3...` | P2SH or P2SH-P2WPKH | BIP16/BIP49 | △ |
| `bc1q...` | P2WPKH or P2WSH | BIP84 | ✅ |
| `bc1p...` | P2TR (Taproot) | BIP86 | ✅ |

### Identifying BCH Address Types

| Prefix | Type | Multisig |
|--------|------|----------|
| `bitcoincash:q...` | P2PKH | ❌ |
| `bitcoincash:p...` | P2SH | ✅ |

### Transaction Size Comparison

| Pattern | Weight | vBytes | Notes |
|---------|--------|--------|-------|
| P2PKH Single-sig (1-in, 2-out) | ~680 | ~170 | Legacy |
| P2WPKH Single-sig (1-in, 2-out) | ~440 | ~110 | Native SegWit |
| P2TR Single-sig (1-in, 2-out) | ~396 | ~99 | Taproot |
| 2-of-3 P2WSH Multisig | ~1,100 | ~275 | Traditional Multisig |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | Signature Aggregation |

---

## See Also

- [llms.txt](./llms.txt) - AI-friendly project sitemap
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [docs/guidelines/](./docs/guidelines) - Project guidelines and standards
