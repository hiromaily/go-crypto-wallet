## Transaction Architecture

### Account Model

Ethereum uses an account-based model (not UTXO):

```
Account State = {
    nonce:    uint64    // Number of transactions sent
    balance:  *big.Int  // Wei balance
    codeHash: []byte    // Empty for EOA
    storage:  map       // Empty for EOA
}
```

**Key Concepts:**

- Each transaction from an address must use the next sequential nonce
- Transactions are identified by hash (not UTXO references)
- Balance changes are atomic (no partial spending)

### Nonce Management

Critical for transaction ordering and double-spend prevention:

```go
// Get pending nonce (includes in-flight transactions)
nonce, err = eth.GetTransactionCount(ctx, fromAddr, QuantityTagPending)

// For batch operations, increment additionalNonce per transaction
effectiveNonce = nonce + additionalNonce
```

**Rules:**

- Always use `pending` state to account for unconfirmed transactions
- Nonces must be sequential — gaps cause transactions to be stuck
- Increment `additionalNonce` when creating multiple transactions in the same batch

### Gas

Gas is the unit of computation cost on Ethereum:

```
Transaction Fee = Gas Used × Effective Gas Price

// Legacy (pre-EIP-1559):
Fee = gasLimit × gasPrice

// EIP-1559 (post-London):
Fee = gasLimit × (baseFee + priorityFee)
```

**Gas Estimation:**

```go
estimatedGas, err = eth.EstimateGas(ctx, msg)
```

---
