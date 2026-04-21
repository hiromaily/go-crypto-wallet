## Fee Management

### Fee Components (EIP-1559)

| Component | Description | Set By |
|-----------|-------------|--------|
| **Base Fee** | Burned by protocol, adjusts per block | Protocol |
| **Priority Fee** (tip) | Paid to validator | Sender (default: 2 Gwei) |
| **Max Fee** | Maximum total gas price willing to pay | Sender |

### Fee Sources

| Source | Method |
|--------|--------|
| **Gas Price (legacy)** | `eth_gasPrice` RPC |
| **Base Fee** | Latest block header (`baseFeePerGas`) |
| **Priority Fee** | Configurable (default 2 Gwei) |
| **Gas Estimate** | `eth_estimateGas` RPC |

### Fee Calculation per Transaction Type

**Deposit / Transfer (full balance sweep):**

```
transferAmount = balance - (gasPrice × estimatedGas)
```

**Payment (fixed amount):**

```
senderCost = amount + (gasPrice × estimatedGas)
recipientReceives = amount
```

---
