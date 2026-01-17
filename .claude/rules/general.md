# Claude Rules - General

## Overview

General rules for Claude Code when working on go-crypto-wallet.

## Behavior Guidelines

- Follow @AGENTS.md for behavior guidelines
- Refer to @docs/standards/ for detailed conventions
- Read @ARCHITECTURE.md for system design

## Code Quality

Refer to @docs/standards/coding-conventions.md

Key rules:
- Follow Clean Architecture layer separation
- Domain layer has ZERO infrastructure dependencies
- Follow language-specific conventions (see Skills below)
- Run file-type-appropriate verification commands before committing (see table below; `*.md` requires none)

## Testing

Refer to @docs/standards/testing.md

## Workflow

Refer to @docs/standards/workflow.md

Key rules:
- **Check current branch FIRST** before starting any task (never work on `main`)
- Create feature branches for changes
- Follow conventional commit messages
- Run file-type-appropriate verification commands before committing (see table below; `*.md` requires none)

## Language-Specific Skills

Use appropriate skills based on the file type:

| File Type | Skill | Verification |
|-----------|-------|--------------|
| Go (`*.go`) | `go-development` | `make go-lint`, `make check-build` |
| TypeScript (`*.ts`) | `typescript-development` | `npm run lint`, `npm run build` |
| Shell (`*.sh`) | `shell-scripts` | `make shfmt` |
| Solidity (`*.sol`) | `solidity-development` | (see skill) |
| SQL/HCL | `db-migration` | `make atlas-fmt`, `make atlas-lint` |
| Makefile | `makefile-update` | `make mk-lint` |
| Markdown (`*.md`) | `docs-update` | (none required) |

## Auto-Generated Files

**DO NOT EDIT** files containing `DO NOT EDIT` comments.
See language-specific skills for details on auto-generated file locations.
