## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Private keys only exist on offline Keygen/Sign wallets
- Keys stored encrypted via scrypt in local keystore files
- Watch Wallet holds no private keys — only public addresses

### Keystore Security

```
// Encryption: scrypt with standard parameters
N = 1 << 18 (StandardScryptN)
P = 1       (StandardScryptP)
r = 8
```

> **Warning:** The keystore password is currently hardcoded in the codebase. This must be externalized to a secure configuration source before production use.

### Transaction Security

- Always validate destination address with `common.IsHexAddress(addr)`
- Verify nonce before signing to prevent replay
- Check sender address after signing (compare recovered address with expected)
- EIP-155 replay protection is enforced via `LondonSigner`

### Anvil / Local Development

- Anvil is detected by version string containing "Anvil"
- Key operations use local filesystem (not RPC) for Anvil compatibility
- Chain ID is forced to testnet value for non-mainnet environments

---
