## Network & Confirmation Guidelines

### Supported Networks

| Network | Chain ID | Purpose | RPC Port |
|---------|----------|---------|----------|
| **Mainnet** | 1 | Production | 8545 |
| **Sepolia** | 11155111 | Public testnet | 8545 |
| **Holesky** | 17000 | Staking testnet | 8545 |
| **Anvil (local)** | 31337 | Local development | 8545 |

**Chain ID Mapping in this system:**

```go
// netID 1 → chainID 1 (mainnet)
// other → chainID 4 (treated as testnet; historical Rinkeby ID)
```

### Confirmation Guidelines

| Confirmations | Risk Level | Typical Use Case |
|---------------|------------|------------------|
| 0 (pending) | High | Development only |
| 1 | Medium | Low-value transfers |
| 6 | Low | Standard commerce |
| 12 | Very Low | Large transactions |
| 32+ | Negligible | High-value custody |

---
