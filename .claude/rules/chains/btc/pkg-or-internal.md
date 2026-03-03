---
paths:
  - internal/infrastructure/api/btc/btc/*.go
  - pkg/chains/btc/*.go
---

# BTC Package Placement Rules

## `pkg/chains/btc/*.go` — Standalone Functions

- Functions that can be defined independently without infrastructure or struct dependencies
- Pure logic that depends only on stdlib or external Bitcoin libraries
- Example: descriptor parsing, address utilities, multisig helpers

### `pkg/chains/btc/rpc/*.go` — Raw RPC Operations

- Functions that perform Bitcoin Core RPC calls via `RawRequest`
- Each function wraps a single RPC command with request/response marshalling
- Must NOT depend on `internal/` packages

```go
// Good: direct RPC call via RawRequest
rawResult, err := c.client.RawRequest("getaddressinfo", []json.RawMessage{input})
```

## `internal/infrastructure/api/btc/btc/` — Infrastructure-Dependent Logic

- Code that depends on `internal/` packages (domain, application, other infrastructure)
- Code that requires `Bitcoin` struct properties (`chainConf`, `coinTypeCode`, `version`, `feeRange`, etc.)
- Business logic that orchestrates multiple RPC calls or combines RPC results with domain logic

```go
// Good: uses Bitcoin struct properties for business logic
func (b *Bitcoin) AdjustFee(fee float64) float64 {
    return fee * b.feeRange.Max
}
```

### Anti-Pattern: Pure Delegation Wrappers

Do NOT create methods in `internal/` that simply delegate to `pkg/` without adding any logic. These add indirection with no value.

```go
// Bad: pure wrapper that adds nothing
func (b *Bitcoin) GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error) {
    return b.pkgrpc.GetAddressInfo(addr)
}
```

If a caller needs an RPC result, it should call `pkgrpc` directly instead of going through a wrapper.

### How Callers Access `pkgrpc`

Callers outside `internal/infrastructure/` access RPC via the `PKGRPCProvider` interface (`internal/application/ports/api/btc`):

```go
// Interface with GetPkgRPC() accessor
type PKGRPCProvider interface {
    GetPkgRPC() btcrpc.BTCRPC
}

// Caller uses GetPkgRPC() directly
func runEncryptWallet(btc apibtc.PKGRPCProvider, passphrase string) error {
    return btc.GetPkgRPC().EncryptWallet(passphrase)
}
```

## Decision Checklist

| Question | Yes → | No → |
|----------|-------|------|
| Does it use `RawRequest` for a single RPC call? | `pkg/chains/btc/rpc/` | ↓ |
| Is it a standalone function with no `internal/` deps? | `pkg/chains/btc/` | ↓ |
| Does it depend on `Bitcoin` struct or `internal/` packages? | `internal/infrastructure/api/btc/btc/` | ↓ |
| Is it just a delegation wrapper to `pkgrpc`? | **Remove it** — call `pkgrpc` directly | — |
