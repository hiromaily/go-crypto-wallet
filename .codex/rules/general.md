# Codex Rules for go-crypto-wallet

## Overview

This file defines general rules for OpenAI Codex CLI when working on this repository.

## General Behavior

- Follow [AGENTS.md](../../AGENTS.md) for behavior guidelines
- Refer to [docs/standards/](../../docs/standards/) for detailed conventions
- Read [ARCHITECTURE.md](../../ARCHITECTURE.md) for system design

## Security

**Security is paramount** - this project handles private keys and cryptocurrency transactions.

Refer to [docs/standards/security.md](../../docs/standards/security.md) for full details.

Key rules:
- **NEVER** log private keys or sensitive information
- Always validate inputs at boundaries
- Consider offline wallet implications

## Code Quality

Refer to [docs/standards/coding-conventions.md](../../docs/standards/coding-conventions.md)

Key rules:
- Follow Clean Architecture layer separation
- Domain layer has ZERO infrastructure dependencies
- Use error wrapping: `fmt.Errorf("context: %w", err)`

## Testing

Refer to [docs/standards/testing.md](../../docs/standards/testing.md)

## Workflow

Refer to [docs/standards/workflow.md](../../docs/standards/workflow.md)

Key rules:
- Create feature branches for changes
- Follow conventional commit messages
- Run verification commands before committing

## Auto-Generated Files

**DO NOT EDIT** files containing `DO NOT EDIT` comments:
- `internal/infrastructure/database/sqlc/*.go`
- `internal/infrastructure/api/xrp/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`
