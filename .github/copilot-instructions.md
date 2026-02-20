# GitHub Copilot Instructions for go-crypto-wallet

## Project Overview

This is a cryptocurrency wallet implementation supporting BTC, BCH, ETH, XRP, and ERC-20 tokens.
**Security is paramount** - this project handles private keys and cryptocurrency transactions.

## Quick Reference

- **Architecture**: Clean Architecture with clear layer separation
- **Primary Language**: Go (backend, wallet operations)
- **Additional Languages**: TypeScript/JavaScript (apps/), Solidity, SQL, HCL, Protobuf

## General Guidelines

- Follow Clean Architecture principles (see [ARCHITECTURE.md](../ARCHITECTURE.md))
- Refer to [docs/guidelines/](../docs/guidelines/) for detailed coding conventions
- Follow [AGENTS.md](../AGENTS.md) for behavior guidelines

## Code Suggestions

- Use Go 1.25+ features where appropriate
- Follow error wrapping pattern: `fmt.Errorf("context: %w", err)`
- Prefer interfaces for dependencies (Dependency Inversion)
- Domain layer has **ZERO infrastructure dependencies**

## Security Requirements

See [docs/guidelines/security.md](../docs/guidelines/security.md) for full details.

Key rules:

- **NEVER** log private keys or sensitive information
- Always validate inputs at boundaries
- Consider offline wallet implications (keygen, sign)
- Security-related changes must be reviewed

## Testing

See [docs/guidelines/testing.md](../docs/guidelines/testing.md) for testing strategy.

## Coding Conventions

See [docs/guidelines/coding-conventions.md](../docs/guidelines/coding-conventions.md) for standards.

## Workflow

See [docs/guidelines/workflow.md](../docs/guidelines/workflow.md) for Git operations and PR guidelines.

## Auto-Generated Files

**DO NOT EDIT** files containing `DO NOT EDIT` comments:

- `internal/infrastructure/database/sqlc/*.go`
- `internal/infrastructure/api/xrp/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`
