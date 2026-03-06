---
description: BTC Package Placement Rules
globs: 
  - internal/infrastructure/api/btc/btc/*.go
  - pkg/chains/btc/*.go
alwaysApply: false
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
- **Facade methods** that delegate to `pkgrpc` but hide the `pkgrpc` layer from callers

```go
// Good: uses Bitcoin struct properties for business logic
func (b *Bitcoin) AdjustFee(fee float64) float64 {
    return fee * b.feeRange.Max
}

// Good: facade method — callers use this instead of b.pkgrpc.EstimateSmartFee()
//       pkgrpc is NOT exposed through the Bitcoiner interface
func (b *Bitcoin) EstimateSmartFee(confirmationBlock int) (float64, error) {
    return b.pkgrpc.EstimateSmartFee(confirmationBlock)
}
```

### Anti-Pattern: Redundant Delegation Wrappers

Do NOT create methods in `internal/` that simply delegate to `pkg/` when the method is
not part of the `Bitcoiner` interface and provides no value as a facade.

```go
// Bad: pure pass-through with no callers
func (b *Bitcoin) SomeRPCMethod() error {
    return b.pkgrpc.SomeRPCMethod()
}
```

### Facade Pattern: The Only Entrypoint to `pkgrpc`

The `Bitcoin` struct exposes all `pkgrpc` operations as its own facade methods.
`pkgrpc` is a private field — it is NOT accessible outside the `btc` package.
`GetPkgRPC()` does NOT exist. BCH overrides use promoted facade methods from `Bitcoin`.

Callers outside the infrastructure layer use methods directly on the `Bitcoiner` interface:

```go
// Good: caller uses direct method — no knowledge of pkgrpc needed
func runEstimateFee(btc apibtc.FeeEstimator) error {
    fee, err := btc.EstimateSmartFee(confirmations)
    ...
}

// Good: caller uses WalletSecurityManager interface
func runEncryptWallet(btc apibtc.WalletSecurityManager, passphrase string) error {
    return btc.EncryptWallet(passphrase)
}
```

### How Internal Infrastructure Accesses `pkgrpc`

Within `internal/infrastructure/api/btc/btc/` (same package): use `b.pkgrpc` directly.

BCH overrides in `internal/infrastructure/api/btc/bch/` call promoted facade methods:

```go
// bch/address.go: calls Bitcoin facade method (not pkgrpc directly)
result, err := b.Bitcoin.GetAddressInfo(addr)
labels, err := b.GetAddressesByLabelMap(labelName)
```

## Decision Checklist

| Question | Yes → | No → |
|----------|-------|------|
| Does it use `RawRequest` for a single RPC call? | `pkg/chains/btc/rpc/` | ↓ |
| Is it a standalone function with no `internal/` deps? | `pkg/chains/btc/` | ↓ |
| Does it depend on `Bitcoin` struct or `internal/` packages? | `internal/infrastructure/api/btc/btc/` | ↓ |
| Is it a facade hiding pkgrpc from callers (no bypass via interface)? | **OK** — add to `Bitcoin` struct and `Bitcoiner` interface | ↓ |
| Is it a pure wrapper with no callers and no interface entry? | **Remove it** — use facade pattern instead | — |
