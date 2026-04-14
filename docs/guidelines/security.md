# Security Standards

Security requirements for go-crypto-wallet. **Security is non-negotiable** in this financial project.

## Critical Rules

| Rule | Description |
|------|-------------|
| Never log private keys | No sensitive data in logs, errors, or commits |
| Never hardcode secrets | Use secure input methods, not CLI arguments |
| Zero-clear memory | Clear sensitive data from memory when done |
| Security review | Required for changes involving sensitive data |

## Security-Critical Areas

These areas require extra caution:

- `internal/infrastructure/wallet/key/` - Key generation
- `internal/domain/key/` - Key value objects
- Any code handling private keys, seeds, or passwords

## Offline Wallet Considerations

This project uses a security model with offline wallets:

| Wallet | Environment | Security Level |
|--------|-------------|----------------|
| Watch | Online | Public keys only |
| Keygen | **Offline** | Generates private keys |
| Sign | **Offline** | Signs transactions |

Always consider the impact of changes on offline wallet operations.

## Security Scans

```bash
make go-check-vuln  # Run vulnerability scan
```

Run for:

- Security-related changes
- Dependency updates
- Encryption/decryption logic changes

## When to Ask for Review

- Any changes to key management code
- Changes to encryption/decryption
- Authentication/authorization changes
- New dependencies that handle sensitive data

## Detailed Guidelines

See [core.md](./core.md) for full security guidelines including error handling, panic usage, and core patterns.
