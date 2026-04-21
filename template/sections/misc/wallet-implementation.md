## Wallet Implementation

### Wallet Types in This System

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, broadcast, monitor | Online |
| **Keygen** | Generate keys, first signature | Offline (air-gapped) |
| **Sign** | Additional signatures (multisig) | Offline (air-gapped) |

### Account Types

| Account | Purpose | Multisig |
|---------|---------|----------|
| **client** | Customer deposit addresses | No |
| **deposit** | Aggregate client funds | No |
| **payment** | Outgoing payments | Yes (2-of-3 or 3-of-3) |
| **stored** | Cold storage | Yes |

For the common 3-wallet transaction flow (chain-agnostic), see [docs/transaction-flow.md](../../transaction-flow.md).
For BTC-specific procedures and Mermaid diagrams, see [operations/wallet-flow.md](operations/wallet-flow.md).

---
