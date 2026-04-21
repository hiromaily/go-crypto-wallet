## Transaction Types

### Legacy Transactions (Type 0)

Uses a single `gasPrice` field. Still supported on all networks.

```go
types.LegacyTx{
    Nonce:    nonce,
    GasPrice: gasPrice,
    Gas:      gasLimit,
    To:       &toAddr,
    Value:    amount,
    Data:     data,
}
```

### EIP-1559 Transactions (Type 2)

Introduced in the London hard fork (2021). Separates gas price into base fee (burned) + priority fee (to validator).

```go
types.DynamicFeeTx{
    ChainID:   chainID,
    Nonce:     nonce,
    GasTipCap: maxPriorityFeePerGas,  // tip to validator
    GasFeeCap: maxFeePerGas,           // max total fee willing to pay
    Gas:       gasLimit,
    To:        &toAddr,
    Value:     amount,
    Data:      data,
}
```

**Fee Calculation (EIP-1559):**

```go
maxPriorityFeePerGas = 2 Gwei  // configurable via config.Ethereum.MaxPriorityFeePerGas
maxFeePerGas = (baseFeePerGas * 2) + maxPriorityFeePerGas
```

**EIP-1559 Detection:**

```go
// Check if network supports EIP-1559 by presence of baseFeePerGas in latest block
supported = eth.SupportsEIP1559(ctx)
```

**Current Implementation Note:** The Watch Wallet currently creates Legacy transactions (`types.LegacyTx`). The EIP-1559 path (`CreateRawTransactionEIP1559`) exists in the infrastructure layer but is not yet used by the Watch use case.

---
