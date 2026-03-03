---
paths:
  - internal/infrastructure/api/eth/eth/*.go
  - pkg/chains/eth/*.go
---

# ETH Package Placement Rules

## `pkg/chains/eth/*.go` — Standalone Functions

- Functions that can be defined independently without infrastructure or struct dependencies
- Pure logic that depends only on stdlib or external Ethereum libraries
- Example: transaction encoding (`EncodeSignedTx`), offline signing (`SignTxOffline`)

### `pkg/chains/eth/rpc/*.go` — Raw RPC Operations

- Methods on `rpcClient` that perform Ethereum JSON-RPC calls via `RPCCaller.CallContext`
- Each method wraps a single RPC command with request/response marshalling
- Must NOT depend on `internal/` packages

```go
// Good: direct RPC call via RPCCaller.CallContext
func (c *rpcClient) BlockNumber(ctx context.Context) (*big.Int, error) {
    var result string
    err := c.caller.CallContext(ctx, &result, "eth_blockNumber")
    ...
}
```

#### Two Client Abstractions

The pkg layer separates two client types:

| Interface    | Backed by              | Purpose                                   |
| ------------ | ---------------------- | ----------------------------------------- |
| `RPCCaller`  | `*ethrpc.Client`       | JSON-RPC calls via `CallContext` + `Close` |
| `ETHCaller`  | `*ethclient.Client`    | Typed calls: `SendTransaction`, `TransactionReceipt` |

Both are defined in `pkg/chains/eth/rpc/rpc.go` and must NOT import `internal/`.

## `internal/infrastructure/api/eth/eth/` — Infrastructure-Dependent Logic

- Code that depends on `internal/` packages (domain, application, other infrastructure)
- Code that requires `Ethereum` struct properties (`ethClient`, `isParity`, `conf`, `coinTypeCode`, etc.)
- Business logic that adds parity guards, config fallbacks, or type conversions at the domain boundary

```go
// Good: uses Ethereum struct state for business logic
func (e *Ethereum) AdminDataDir(ctx context.Context) (string, error) {
    if e.isParity {
        return "", nil  // parity guard using struct field
    }
    return e.pkgrpc.AdminDataDir(ctx)
}

// Good: type conversion at domain boundary
func (e *Ethereum) GetBalance(
    ctx context.Context, hexAddr string, quantityTag domainETH.QuantityTag,
) (*big.Int, error) {
    return e.pkgrpc.GetBalance(ctx, hexAddr, ethrpc.QuantityTag(quantityTag))
}
```

### Anti-Pattern: Pure Delegation Wrappers

Do NOT create methods in `internal/` that simply delegate to `pkgrpc` without adding any logic. These add indirection with no value.

```go
// Bad: pure wrapper that adds nothing
func (e *Ethereum) BlockNumber(ctx context.Context) (*big.Int, error) {
    return e.pkgrpc.BlockNumber(ctx)
}
```

If a caller needs an RPC result, it should call `pkgrpc` directly via `GetPkgRPC()` instead of going through a wrapper.

### Method Classification in `internal/infrastructure/api/eth/eth/`

| Method | File | Classification | Reason |
|--------|------|----------------|--------|
| `SuggestGasTipCap` | `rpc_eth_gas.go` | Keep | uses `e.ethClient`, config fallback (`e.conf`) |
| `GetBalance` | `rpc_eth.go` | Keep | domain type conversion (`domainETH.QuantityTag` → `ethrpc.QuantityTag`) |
| `GetTransactionCount` | `rpc_eth.go` | Keep | domain type conversion (same) |
| `GetTxReceipt` | `rpc_eth_tx.go` | Keep | domain type conversion (`pkgrpc.TransactionReceipt` → `domainETH.TransactionReceipt`) |
| All others | `rpc_*.go` | Anti-pattern | pure `e.pkgrpc.Xxx(...)` delegation |

### How Callers Access `pkgrpc`

Callers outside `internal/infrastructure/` access RPC via the `PKGRPCProvider` interface (`internal/application/ports/api/eth`):

```go
// Interface with GetPkgRPC() accessor
type PKGRPCProvider interface {
    GetPkgRPC() ethrpc.ETHRPC
}

// Caller uses GetPkgRPC() directly instead of a delegation wrapper
func runBlockNumber(eth apiETH.PKGRPCProvider) (*big.Int, error) {
    return eth.GetPkgRPC().BlockNumber(ctx)
}
```

This pattern eliminates pure delegation wrappers from `Ethereum` while still allowing use cases to call raw RPC methods when needed.

## Decision Checklist

| Question | Yes → | No → |
|----------|-------|------|
| Does it call `CallContext` for a single JSON-RPC method? | `pkg/chains/eth/rpc/` | ↓ |
| Is it a standalone function with no `internal/` deps? | `pkg/chains/eth/` | ↓ |
| Does it use `e.isParity`, `e.conf`, or `e.ethClient`? | `internal/infrastructure/api/eth/eth/` | ↓ |
| Does it convert between `domainETH.*` and `ethrpc.*` types? | `internal/infrastructure/api/eth/eth/` | ↓ |
| Is it just `return e.pkgrpc.Xxx(...)` with no added logic? | **Remove it** — call `pkgrpc` via `GetPkgRPC()` | — |
