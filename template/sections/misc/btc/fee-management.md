## Fee Management

### Fee Estimation

Bitcoin Core provides fee estimation via RPC:

```bash
# Estimate fee for confirmation in N blocks
bitcoin-cli estimatesmartfee <conf_target> [estimate_mode]

# Modes: UNSET, ECONOMICAL, CONSERVATIVE
```

### Fee Rate Sources

| Source | Endpoint/Method |
|--------|-----------------|
| **Bitcoin Core** | `estimatesmartfee` RPC |
| **Mempool.space** | `https://mempool.space/api/v1/fees/recommended` |
| **Blockstream** | `https://blockstream.info/api/fee-estimates` |

### Fee Optimization Strategies

1. **SegWit/Taproot** - Use native SegWit or Taproot for smaller transactions
2. **UTXO Consolidation** - Consolidate UTXOs during low-fee periods
3. **Batching** - Combine multiple payments in single transaction
4. **RBF** - Use Replace-by-Fee for fee bumping if needed

---
