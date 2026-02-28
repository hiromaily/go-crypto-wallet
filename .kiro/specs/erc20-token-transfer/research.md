# Research & Design Decisions

---

## Summary

- **Feature**: `erc20-token-transfer`
- **Discovery Scope**: Extension — brownfield Go codebase modification
- **Key Findings**:
  - The DI routing is already fully automated via `IsETHGroup`/`IsERC20Token` checks; no new CLI commands, use cases, or adapters required
  - `ERC20.SupportsEIP1559()` is hardcoded `false`; the `Ethereum` struct holds the correct implementation and the existing TODO comment explicitly calls for embedding `*eth.Ethereum` in `ERC20`
  - Embedding `*eth.Ethereum` as a named field (composition, not Go promotion) is safer to avoid accidental method promotion that could break the `ERC20er` interface contract

---

## Research Log

### DI Routing Automation

- **Context**: Determine whether new CLI commands or DI wiring are needed for HYC
- **Sources Consulted**: `internal/di/container.go` (direct read)
- **Findings**:
  - `NewWatchCreateTransactionUseCase` dispatches via `domainCoin.IsETHGroup(c.conf.CoinTypeCode)` → `newETHWatchCreateTransactionUseCase()`
  - Inside `newETHWatchCreateTransactionUseCase`, `domainCoin.IsERC20Token(...)` selects `newERC20()` vs `newETH()`
  - `NewWatchSendTransactionUseCase` and `NewWatchMonitorTransactionUseCase` both dispatch via `IsETHGroup`
  - The same `ETHWalleter` instance handles all ETH-group coins including ERC-20 tokens
- **Implications**: Adding HYC to `ERC20Map` and `CoinTypeCodeValue` is sufficient to activate all watch wallet flows. Zero changes to use cases, CLI adapters, or routing logic.

### EIP-1559 Implementation in `Ethereum` vs `ERC20`

- **Context**: `ERC20.SupportsEIP1559()` returns hardcoded `false`; Req 3 requires real detection
- **Sources Consulted**: `internal/infrastructure/api/eth/eth/transaction.go`, `internal/infrastructure/api/eth/erc20/erc20.go`
- **Findings**:
  - `Ethereum.SupportsEIP1559()`: checks `clientType == ClientVersionAnvil` OR queries `blockInfo.BaseFeePerGas != nil`
  - `Ethereum.CreateRawTransactionEIP1559()`: builds `types.DynamicFeeTx` with `GasTipCap` and `GasFeeCap`, same encoding pipeline (`ethtx.EncodeTx`)
  - `ERC20.getNonce()` duplicates `Ethereum.getNonce()` (acknowledged FIXME on line 279)
  - `ERC20` holds `*ethclient.Client` directly; `Ethereum` holds `*ethclient.Client` + `*ethrpc.Client` + full config
- **Implications**: Embedding `*eth.Ethereum` provides `SupportsEIP1559`, `BlockNumber`, `GetBlockByNumber`, `SuggestGasTipCap`, and `getNonce` for free, eliminating all duplication.

### Embedding Strategy: Named Field vs Go Promotion

- **Context**: How to embed `*eth.Ethereum` in `ERC20` without leaking unrelated methods
- **Sources Consulted**: Go spec (interface satisfaction), existing `ERC20er` interface
- **Findings**:
  - Go type promotion (`type ERC20 struct { *eth.Ethereum }`) would promote all ~40 Ethereum methods onto ERC20, making `ERC20` accidentally satisfy `Ethereumer` — undesirable coupling
  - Named field (`eth *eth.Ethereum`) keeps methods private to the ERC20 implementation, calling them explicitly (`e.eth.SupportsEIP1559(ctx)`)
  - The compile-time check `var _ apieth.ERC20er = (*ERC20)(nil)` remains clean either way
- **Implications**: Use named field `eth *eth.Ethereum` in `ERC20` struct. The DI container passes `c.newETH()` (cast to a minimal helper interface if needed) to the `NewERC20` constructor.

### DI Constructor: Avoiding Duplicate ethclient Connections

- **Context**: `newERC20()` currently creates its own `ethclient.Client`; `newETH()` creates another
- **Sources Consulted**: `internal/di/container.go` lines 503-527
- **Findings**:
  - `newETH()` is already lazy-initialized and cached: `if c.eth == nil { ... }`
  - Calling `c.newETH()` inside `newERC20()` is safe — it returns the cached instance
  - `*eth.Ethereum` exposes `SupportsEIP1559`, `BlockNumber`, `GetBlockByNumber`, `SuggestGasTipCap` — all needed by ERC-20 EIP-1559
  - No circular initialization: `newERC20` depends on `newETH`; `newETH` does not depend on `newERC20`
- **Implications**: `newERC20()` should call `c.newETH()` and pass the result to `NewERC20`. A minimal `eip1559Detector` interface can be defined in the `erc20` package to keep the dependency explicit.

### EIP-1559 Fee Strategy for ERC-20

- **Context**: Req 3.3 states fees should come from wallet configuration; `Ethereum.CreateRawTransactionEIP1559` derives fees dynamically from block header
- **Sources Consulted**: `pkg/config/wallet.go`, `Ethereum.CreateRawTransactionEIP1559`
- **Findings**:
  - `config.Ethereum` already has `MaxPriorityFeePerGas uint64` and `MaxFeePerGas` cap fields
  - The native ETH implementation ignores these config fields and always queries the node dynamically
  - For ERC-20 tokens interacting with contracts, gas estimation includes the `data` field (calldata), making dynamic fee calculation more important
- **Implications**: Follow the same dynamic fee strategy as `Ethereum.CreateRawTransactionEIP1559` (query `SuggestGasTipCap` + `baseFee * 2`). Config `MaxFeePerGas` can serve as an upper-bound safety cap. Return error if computed `maxFeePerGas` exceeds the configured cap (when cap > 0).

### Contract Address for Config

- **Context**: `config/wallet/eth/*.yaml` requires a `contract_address` for HYC
- **Sources Consulted**: `erc20-token` spec deployment output pattern
- **Implications**: The contract address is set to a placeholder (`0x...`). Operators must replace it with the address from `forge script` deployment output. This is documented as an operator step, not a code gap.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Decision |
|--------|-------------|-----------|---------------------|----------|
| Named field composition | `ERC20` holds `eth *eth.Ethereum` as private field | No accidental method promotion; explicit calls | Slightly more verbose | **Selected** |
| Go promotion (embedding) | `ERC20` embeds `*eth.Ethereum` | Less code | Leaks all ~40 Ethereum methods; breaks interface isolation | Rejected |
| Duplicate minimal logic | Copy `SupportsEIP1559` check without embedding | No structural change | Perpetuates FIXME duplication; violates DRY | Rejected |

---

## Design Decisions

### Decision: Named Field Composition for `eth.Ethereum` in `ERC20`

- **Context**: `ERC20` needs `SupportsEIP1559`, `BlockNumber`, `GetBlockByNumber`, `SuggestGasTipCap` from `Ethereum` without exposing the entire `Ethereumer` interface
- **Alternatives Considered**:
  1. Go promotion (`*eth.Ethereum` embedded anonymously) — promotes all methods
  2. Named field (`eth *eth.Ethereum`) — explicit delegation
  3. Duplicate only needed methods — code duplication
- **Selected Approach**: Named field `eth *eth.Ethereum`; `ERC20` calls `e.eth.SupportsEIP1559(ctx)`, `e.eth.SuggestGasTipCap(ctx)`, etc.
- **Rationale**: Maintains interface isolation (`ERC20er` contract stays clean), resolves the existing TODO, eliminates all duplication
- **Trade-offs**: `NewERC20` constructor gains one new parameter; DI container must be updated
- **Follow-up**: Verify that `c.newETH()` is safe to call during `newERC20()` initialization (no circular dep)

### Decision: Dynamic EIP-1559 Fee Calculation for ERC-20

- **Context**: Req 3.3 says fees from config; but dynamic fees ensure correctness under varying base fees
- **Alternatives Considered**:
  1. Config-static fees — simple but can cause failed txs if base fee spikes
  2. Dynamic fees (same as `Ethereum.CreateRawTransactionEIP1559`) — correct but ignores config
  3. Dynamic with config cap — dynamic base + config max cap as safety guard
- **Selected Approach**: Dynamic fees (option 3): `maxFeePerGas = (baseFee * 2) + tip`; if config `MaxFeePerGas > 0`, return error when computed value exceeds cap
- **Rationale**: Consistent with native ETH implementation; safety cap prevents accidental overpayment
- **Trade-offs**: Slightly more complex; requires config validation at startup
- **Follow-up**: Config validation in `pkg/config/wallet.go` `HasERC20Config` can add cap check

### Decision: DI Passes `c.newETH()` to `NewERC20`

- **Context**: Avoid duplicate `ethclient.Client` connections; share the cached Ethereum instance
- **Selected Approach**: `newERC20()` calls `c.newETH()` and passes the result to `NewERC20` as the new `eth *eth.Ethereum` parameter
- **Rationale**: `newETH()` is lazy-cached; no circular initialization; single connection to node
- **Trade-offs**: Tighter coupling between ERC20 and Ethereum constructors in DI
- **Follow-up**: Confirm container lazy-init order: `newETH()` has no dependency on `newERC20()`

---

## Risks & Mitigations

- **HYC contract address placeholder** — Config files use `0x...` placeholder. Mitigation: document in README that operator must replace after `forge deploy`. Validation via `HasERC20Config` will still pass (non-empty string).
- **`newETH()` circular init** — If `Ethereum` construction ever calls `newERC20()`, circular init occurs. Mitigation: confirm DI wiring direction; add a comment in `container.go`.
- **`ERC20` no longer independent** — After embedding, `ERC20` requires a valid `*eth.Ethereum`. Mitigation: test constructor with nil guard; panic in DI is acceptable per coding conventions.
- **EIP-1559 gas estimation with contract data** — `estimateGas` in `ERC20` already includes the `transfer` calldata, which is correct for token transfers. The fee cap check is additive, not blocking.

---

## References

- [go-ethereum types.DynamicFeeTx](https://pkg.go.dev/github.com/ethereum/go-ethereum/core/types#DynamicFeeTx) — EIP-1559 transaction struct fields
- [EIP-1559 specification](https://eips.ethereum.org/EIPS/eip-1559) — base fee mechanism
- [go-ethereum ethclient](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) — `SuggestGasTipCap`, `BlockByNumber`
