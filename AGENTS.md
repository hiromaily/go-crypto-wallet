# Agent Guidelines for go-crypto-wallet

This document provides guidelines for AI agents working on this project.

## Project Overview

- This project is a cryptocurrency wallet implementation in Go supporting BTC, BCH, ETH, XRP, and ERC-20 tokens
- The project is currently under refactoring based on Clean Architecture and Clean Code principles
- Security is of utmost importance (private key management, offline wallets)
- The project follows the `pkg` layout pattern

## Architecture Principles

- Follow Clean Architecture principles
- Maintain clear layer separation (domain, application, infrastructure)
- Use dependency injection and abstract with interfaces
- Follow the `pkg` layout pattern

## Coding Standards

- Follow `golangci-lint` configuration (`.golangci.yml`)
- Use `goimports` to maintain import order (standard → third-party → local)
- Format code with `gofmt`/`goimports`
- Maintain consistent naming conventions (lowercase package names, exported functions start with uppercase)

## Common Commands

After making code changes, use these commands to verify code correctness:

- `make lint-fix`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: Run these commands after code changes to ensure code quality and correctness.

## Error Handling

- **DO NOT** use `log.Fatal` outside of `main` function (see `REFACTORING_CHECKLIST.md`)
- Wrap errors with `fmt.Errorf` + `%w`
- Use `errors.Is`/`errors.As` for error checking
- Include context information in error messages

## Context Management

- Add `context.Context` to all API calls
- Implement timeouts and cancellation
- Implement graceful shutdown

## Security

- **NEVER** log private keys or sensitive information
- Encrypt or zero-clear private keys in memory when possible
- Do not pass passwords via command-line arguments; use secure input methods
- Conduct security review when making changes involving sensitive information

## Wallet Types Understanding

- **Watch Wallet**: Online, public keys only, creates and sends transactions
- **Keygen Wallet**: Offline, generates keys, first signature for multisig
- **Sign Wallet**: Offline, second and subsequent signatures for multisig

## Directory Structure

- `cmd/`: Application entry points (keygen, sign, watch)
- `pkg/`: Package code
  - `wallet/api/`: External API clients (btcgrp, ethgrp, xrpgrp)
  - `wallet/service/`: Business logic (btc, eth, xrp, coldsrv, watchsrv)
  - `wallet/key/`: Key generation logic
- `data/`: Generated files, configuration files
- `scripts/`: Operation scripts

## Refactoring Status

- Refer to `REFACTORING_CHECKLIST.md` for current refactoring tasks
- Make changes incrementally without breaking existing functionality
- Run tests before and after refactoring

## Testing

- Use `//go:build integration` tag for integration tests
- Separate unit tests and integration tests
- Measure and improve test coverage

## Dependency Management

- Use `go mod tidy` to organize dependencies
- Follow procedures in `REFACTORING_CHECKLIST.md` when updating dependencies
- Run security scans (`govulncheck`)

## Logging

- Use structured logging
- Set appropriate log levels
- **NEVER** log sensitive information (private keys, passwords, etc.)

## Patterns to Avoid

- Using `log.Fatal` (except in `main`)
- Leaving commented-out code
- Unused imports, variables, or functions
- Ignoring errors (detected by `errcheck`)

## Recommended Patterns

- Abstraction through interfaces
- Dependency injection
- Proper error wrapping with context
- Use of `context.Context`

## Documentation

- Add godoc comments to exported functions and methods
- Add package-level comments
- Include usage examples for complex logic

## Multi-Chain Support

- **BTC/BCH**: Bitcoin Core RPC API
- **ETH**: Ethereum JSON-RPC API, ERC-20 token support
- **XRP**: Communication via gRPC with ripple-lib-server

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
