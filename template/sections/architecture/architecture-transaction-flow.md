## Architecture & Transaction Flow

| Document | Description |
|----------|-------------|
| [architecture.md](./architecture.md) | **[SSOT] ETH wallet architecture** — wallet roles, use case assignments, Clean Architecture boundary maps (EOA + Safe), port interfaces, offline signing detail |
| [multisig.md](./multisig.md) | **ETH Safe multisig** — Safe v1.4.1 implementation, EIP-712 signing flow, file format, CLI commands, E2E Pattern 3 |
| [docs/transaction-flow.md](../../transaction-flow.md) | Chain-agnostic 3-wallet setup, signing, and monitoring flows |

For Ethereum-specific concerns on top of the common flow, see [ETH-Specific Flow Details](#eth-specific-flow-details) below.

> **Key point:** ETH supports both single-sig EOA and Safe multisig flows.
> For single-sig, only Watch and Keygen wallets are required.
> For Safe multisig (E2E Pattern 3), all three wallets are used: Watch proposes and submits, Keygen and Sign wallets each sign offline.
> See [multisig.md](./multisig.md) for the Safe multisig implementation details.

---
