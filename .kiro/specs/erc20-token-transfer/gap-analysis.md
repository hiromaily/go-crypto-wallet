# Gap Analysis: erc20-token-transfer

## Summary

- **Scope**: Go-only changes — no Solidity, no new CLI commands, no new use cases required
- **Key insight**: The DI routing is already fully automatic for any `IsETHGroup` coin, so adding HYC to the domain type system automatically activates all watch wallet flows
- **Primary gap**: `ERC20.SupportsEIP1559()` is hardcoded `false`; the `Ethereum` struct already has the full EIP-1559 implementation but `ERC20` does not embed it
- **Configuration gap**: `keygen.yaml` and `sign.yaml` have no `erc20s` section; `watch.yaml` has only `hyt`, missing `hyc`
- **Recommended approach**: Option A (extend existing), all changes are additive — no new files needed except tests

---

## 1. Current State Investigation

### Key Files Identified

| File | Role | Status for HYC |
|------|------|----------------|
| `internal/domain/coin/types.go` | Token type constants and validation | Missing HYC in all 3 registries |
| `internal/infrastructure/api/eth/erc20/erc20.go` | ERC-20 transaction infrastructure | EIP-1559 hardcoded `false` |
| `internal/infrastructure/api/eth/eth/transaction.go` | Native ETH with full EIP-1559 | Complete; can be reused |
| `internal/di/container.go` | DI routing for all use cases | Already routes by `IsETHGroup`/`IsERC20Token` |
| `config/wallet/eth/watch.yaml` | Watch wallet config | Has `hyt`, missing `hyc` and `erc20_token` |
| `config/wallet/eth/keygen.yaml` | Keygen wallet config | No `erc20s` section at all |
| `config/wallet/eth/sign.yaml` | Sign wallet config | No `erc20s` section at all |
| `pkg/config/wallet.go` | Config struct `Ethereum.ERC20s` | Complete, generic map |

### DI Routing Analysis (Critical Finding)

The DI container already handles routing automatically based on domain type checks:

```go
// container.go: watch create-tx routing
case domainCoin.IsETHGroup(c.conf.CoinTypeCode):     // ← HYC will match once registered
    return c.newETHWatchCreateTransactionUseCase()

// container.go: inside newETHWatchCreateTransactionUseCase
if domainCoin.IsERC20Token(c.conf.CoinTypeCode.String()) {  // ← HYC picks ERC20 API
    targetEthAPI = c.newERC20()
} else {
    targetEthAPI = c.newETH()
}
```

**Implication**: Adding HYC to `ERC20Map` and `CoinTypeCodeValue` is sufficient to activate the watch wallet's create/send/monitor flows. No changes to use cases, CLI adapters, or DI routing logic are required.

### EIP-1559 Implementation Gap

`Ethereum.SupportsEIP1559()` (in `eth/transaction.go`) is fully implemented:
- Checks if `clientType == ClientVersionAnvil` → returns `true`
- Falls back to checking `baseFeePerGas` in the latest block header

`ERC20.SupportsEIP1559()` (in `erc20/erc20.go`):
```go
func (*ERC20) SupportsEIP1559(_ context.Context) bool {
    return false  // ← hardcoded; no embedding of Ethereum struct
}
```

The `ERC20` struct holds an `*ethclient.Client` directly — it does NOT embed `*eth.Ethereum`. The TODO comment on line 28-29 of `erc20.go` explicitly notes this as a known gap:
> `// TODO: Ethereum struct in internal/infrastructure/api/eth/eth/ethereum.go must be embedded`

### Config Gap Detail

| Config File | `erc20s` key | `erc20_token` key |
|-------------|-------------|-------------------|
| `watch.yaml` | Has `hyt` entry, missing `hyc` | Not set (defaults to empty) |
| `keygen.yaml` | Absent entirely | Absent |
| `sign.yaml` | Absent entirely | Absent |

The config struct `Ethereum.ERC20s map[ERC20Token]ERC20` is fully generic — adding entries requires no code changes.

---

## 2. Requirements Feasibility Analysis

### Requirement 1: HYC Domain Registration

| Need | Gap | Complexity |
|------|-----|------------|
| `CoinTypeERC20HYC CoinType = 9002` | Missing from `types.go` | Trivial — 1-line addition |
| `HYC CoinTypeCode = "hyc"` + `CoinTypeCodeValue` entry | Missing | Trivial — 2-line addition |
| `TokenHYC ERC20Token = "hyc"` + `ERC20Map` entry | Missing | Trivial — 2-line addition |

No downstream code changes needed; DI and use cases already call `IsERC20Token`/`IsETHGroup`.

### Requirement 2: Config YAML

| Need | Gap | Complexity |
|------|-----|------------|
| `watch.yaml`: add `hyc` under `erc20s` | Present for `hyt`, add `hyc` | Trivial — YAML addition |
| `keygen.yaml`: add `erc20s.hyc` block | Entire `erc20s` key absent | Trivial — YAML addition |
| `sign.yaml`: add `erc20s.hyc` block | Entire `erc20s` key absent | Trivial — YAML addition |
| Contract address for `hyc` | **Unknown** — depends on anvil deployment | Research Needed |

> **Research Needed**: The `contract_address` for HYC in config depends on the deployed contract address from the `erc20-token` spec. This value must be obtained from the Foundry deployment output.

### Requirement 3: EIP-1559 for ERC-20

| Need | Gap | Complexity |
|------|-----|------------|
| `SupportsEIP1559()` real impl | Currently `return false` | Small — embed or copy detection logic |
| `CreateRawTransactionEIP1559()` real impl | Currently delegates to legacy | Medium — build DynamicFeeTx with contract `data` field |
| Reuse `Ethereum` struct methods | Not embedded (TODO comment) | Medium — structural refactor needed |

**Option A** (embed `*eth.Ethereum` in `ERC20`): Fixes both EIP-1559 gap and the `getNonce` / `SupportsEIP1559` FIXME duplication simultaneously. Requires changing `ERC20` struct fields and constructor.

**Option B** (duplicate minimal logic): Copy only `SupportsEIP1559` check and `DynamicFeeTx` construction without embedding. Faster but perpetuates the code duplication TODO.

### Requirements 4–7: Watch Wallet Flows

**Finding: No code changes needed.** Once domain registration is complete:

- `create-tx` with `erc20_token: "hyc"` → `IsETHGroup` → `newETHWatchCreateTransactionUseCase` → `IsERC20Token` → `newERC20()` → existing `CreateTransactionUseCase`
- `send-tx` → `IsETHGroup` → `newETHWatchSendTransactionUseCase` → existing `SendTransactionUseCase`
- `monitor-tx` → `IsETHGroup` → `newETHWatchMonitorTransactionUseCase` → existing `MonitorTransactionUseCase`

The CLI adapter (`newETHWalleter`) and all use cases already support any `ERC20er` implementation. No new CLI commands, use cases, or adapter changes required.

### Requirement 8: Tests

| Need | Gap | Complexity |
|------|-----|------------|
| Tests for `IsCoinTypeCode("hyc")`, `IsERC20Token("hyc")`, `IsETHGroup("hyc")` | No `types_test.go` for ERC20 tokens | Small — pure unit tests |
| Tests for `SupportsEIP1559()` true/false path | No tests in `erc20/` package | Small — mock `ethclient` needed |
| Tests for `CreateRawTransactionEIP1559()` Type 2 fields | No EIP-1559 tests in `erc20/` | Medium — mock chain calls |

---

## 3. Implementation Approach Options

### Option A: Extend Existing Components (Recommended)

**Strategy**: Additive changes only. Embed `*eth.Ethereum` into `ERC20` struct.

**Files modified** (no new files):

1. `internal/domain/coin/types.go` — add 5 lines (HYC constants + map entries)
2. `internal/infrastructure/api/eth/erc20/erc20.go` — embed `*eth.Ethereum`, implement real `SupportsEIP1559` and `CreateRawTransactionEIP1559`
3. `config/wallet/eth/watch.yaml` — add `hyc` ERC20 entry
4. `config/wallet/eth/keygen.yaml` — add `erc20s.hyc` section
5. `config/wallet/eth/sign.yaml` — add `erc20s.hyc` section

**New test files**:

6. `internal/domain/coin/types_erc20_test.go` — domain registration tests
7. `internal/infrastructure/api/eth/erc20/erc20_test.go` — EIP-1559 infrastructure tests

**Trade-offs**:

- ✅ Resolves both FIXME/TODO comments in `erc20.go` simultaneously
- ✅ Minimal surface area — only 5 source files changed
- ✅ No new interfaces needed — `ERC20er` already has all required methods
- ❌ `newERC20()` in DI must be updated to pass `*eth.Ethereum` reference to embedded struct

### Option B: Minimal Change (Partial Fix)

**Strategy**: Add only domain registration and configs. For EIP-1559, keep the delegation to `CreateRawTransaction` (legacy fallback) but add real `SupportsEIP1559` using a standalone `ethclient` call.

**Trade-offs**:

- ✅ Smallest possible change
- ❌ Perpetuates code duplication (FIXME comments remain)
- ❌ Requirement 3 is only partially satisfied (EIP-1559 detection works but still falls back to legacy tx)

### Option C: Not Applicable

No new components are warranted — the Clean Architecture layers are already correctly structured.

---

## 4. Effort and Risk

| Requirement | Approach | Effort | Risk | Notes |
|-------------|----------|--------|------|-------|
| Req 1: Domain registration | Option A | S | Low | Pure constant additions, zero logic |
| Req 2: Config YAML | Option A | S | Low | Contract address is unknown (Research Needed) |
| Req 3: EIP-1559 for ERC-20 | Option A | M | Medium | Embedding `*eth.Ethereum`; DI constructor change |
| Reqs 4–7: Watch flows | N/A | S | Low | No code changes; automatic via DI routing |
| Req 8: Tests | Option A | S | Low | Standard unit test patterns |
| **Total** | Option A | **M** | **Low** | Config blocked on contract address |

---

## 5. Recommendations for Design Phase

### Preferred Approach

**Option A** with embedding `*eth.Ethereum` into `ERC20`. This is the direction the existing TODO comment already points to, and it resolves the EIP-1559 gap cleanly without code duplication.

### Key Design Decisions

1. **Embedding vs composition**: Should `ERC20` embed `*eth.Ethereum` (Go embedding) or hold it as a named field? Embedding is idiomatic Go but increases coupling. Named field (`eth *eth.Ethereum`) is more explicit.

2. **DI constructor update**: `newERC20()` in `container.go` currently creates its own `ethclient.Client`. With embedding, it should reuse `c.newETH()` (the existing `*eth.Ethereum` instance) to avoid two separate connections.

3. **EIP-1559 fee strategy for ERC-20**: The `Ethereum.CreateRawTransactionEIP1559` calculates `maxFeePerGas = (baseFee * 2) + tip`. Should the ERC-20 variant use the same formula, or read `MaxFeePerGas` / `MaxPriorityFeePerGas` from config (as stated in Req 3.3)? Config-driven is safer for tokens with unpredictable gas needs.

### Research Items

- **Contract address for HYC**: Must be obtained from the `erc20-token` spec's Foundry deployment output (`forge script` stdout or `broadcast/` directory). This is the only external dependency for Req 2.

- **DI circular dependency**: Confirm that `newERC20()` can safely accept `c.newETH()` without creating duplicate `ethclient.Client` connections or circular initialization in the DI container.
