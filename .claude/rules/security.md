# Claude Rules - Security

## Overview

Security rules for Claude Code when working on go-crypto-wallet.
**Security is paramount** - this project handles private keys and cryptocurrency transactions.

## Detailed Rules

Refer to @docs/guidelines/security.md for full security requirements.

## Quick Reference

### Critical Rules

1. **NEVER** log private keys or sensitive information
2. **NEVER** commit secrets or credentials
3. Always validate inputs at system boundaries
4. Consider offline wallet implications for keygen/sign operations
5. Security-related changes must be reviewed

### Security-Critical Areas

- `internal/infrastructure/wallet/key/` - Key generation
- `internal/domain/key/` - Key value objects
- Any code handling private keys, seeds, or passwords

### Wallet Architecture

| Wallet | Environment | Security Level |
|--------|-------------|----------------|
| Watch | Online | Public keys only |
| Keygen | **Offline** | Generates private keys |
| Sign | **Offline** | Signs transactions |

## Verification

Before committing security-related changes:

```bash
make go-check-vuln
make go-lint
make gotest
```
