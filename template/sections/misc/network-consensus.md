## Network & Consensus

### Networks

| Network | Purpose | Port | RPC Port | Magic Bytes |
|---------|---------|------|----------|-------------|
| **Mainnet** | Production | 8333 | 8332 | 0xF9BEB4D9 |
| **Testnet3** | Public testing | 18333 | 18332 | 0x0B110907 |
| **Signet** | Controlled testing | 38333 | 38332 | 0x0A03CF40 |
| **Regtest** | Local development | 18444 | 18443 | 0xFABFB5DA |

### Confirmation Guidelines

| Confirmations | Risk Level | Typical Use Case |
|---------------|------------|------------------|
| 0 (unconfirmed) | High | Very small amounts, trusted parties |
| 1 | Medium | Small retail transactions |
| 3 | Low | Most commerce |
| 6 | Very Low | Large transactions |
| 100+ | None | Coinbase maturity |

---
