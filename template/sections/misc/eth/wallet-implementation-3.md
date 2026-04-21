## Wallet Implementation

> **Architecture SSOT:** Wallet roles, use case assignments, Clean Architecture boundary map, and signing flows are documented in [architecture.md](../../../../docs/chains/eth/architecture.md). This section provides a quick-reference summary only.

### Wallet Roles (Summary)

| Wallet | Single-sig EOA | Safe Multisig | Network |
|--------|---------------|---------------|---------|
| **Watch** | Create transactions, broadcast, monitor | Propose multisig tx, submit `execTransaction`, monitor | Online |
| **Keygen** | Generate keys, sign transactions | Generate keys, sign as Safe owner 1 | Offline (air-gapped) |
| **Sign** | Not used | Sign as Safe owner 2…n | Offline (air-gapped) |

See [architecture.md](../../../../docs/chains/eth/architecture.md) for the complete use case assignment table, architecture boundary map, and offline signing detail.

### Database Schema (ETH-specific tables)

| Table | Description |
|-------|-------------|
| `account_key` | Key generation tracking with status |
| `address` | Generated Ethereum addresses |
| `full_pubkey` | Extended public keys |
| `eth_tx` | Unsigned/signed transaction records |
| `eth_detail_tx` | Per-transaction detail (nonce, gas, amounts) |
| `payment_request` | Outbound payment requests |

### Address Status Lifecycle

```
AddrStatusHDKeyGenerated
    → (import privkey)
AddrStatusPrivKeyImported
```

### Transaction Status Lifecycle

```
TxTypeSent ──(confirmations >= threshold)──> TxTypeDone ──(notified)──> TxTypeNotified
```

---
