# Research & Design Decisions

---
**Feature**: `refactoring-chain-rpc`
**Discovery Scope**: Extension / Complex Refactoring (brownfield, multi-layer)
**Key Findings**:
- `pkg/chains/*/rpc/` does not exist; only three factory functions are present — all RPC logic lives in `internal/`.
- The BTC mapper (`mapper.go`, 768 lines) and XRP converter (`converter.go`, 635 lines) hold most conversion logic; these are pure field renames, not domain transformations.
- ETH already returns domain types (`*big.Int`, `*domainETH.*`) for most methods — its Phase 2 scope is minimal.

---

## Research Log

### RPC Client Abstraction in `pkg/`

- **Context**: `pkg/chains/*/rpc/` functions must accept a client parameter. The concrete types (`*rpcclient.Client`, `*ethrpc.Client`, `*websocket.WS`) all live in third-party packages. We need to decide whether to accept the concrete type or define a minimal interface.
- **Findings**:
  - `*rpcclient.Client` (btcsuite/btcd) has a `RawRequest(method string, params []json.RawMessage) (json.RawMessage, error)` method — this is the only method used by all BTC RPC call functions.
  - `*ethrpc.Client` (go-ethereum) has `CallContext(ctx context.Context, result any, method string, args ...any) error` — this is the only method used by ETH `rpc_*.go` files. `*ethclient.Client` uses `SuggestGasTipCap` and `SendTransaction` as higher-level wrappers.
  - `*websocket.WS` (internal pkg) has `Call(ctx context.Context, req, resp any) error`.
- **Implications**: Define a minimal interface in each `pkg/chains/*/rpc/` package. Concrete types satisfy these interfaces automatically (Go's implicit interface satisfaction). This enables unit testing without a live node.

### DTO Field Naming — Idiomatic vs Protocol

- **Context**: `GetAddressInfoResult` (current infrastructure type) uses JSON-protocol field names (`Ismine`, `Iscompressed`). `dtobtc.AddressInfo` (current DTO) uses Go-idiomatic names (`IsMine`, `IsCompressed`). The mapper converts between them.
- **Findings**: The two types carry identical data — the only difference is field naming convention. The mapper is essentially a renaming exercise.
- **Implications**: `pkg/chains/btc/rpc/` response types should use **Go-idiomatic field names** with JSON struct tags matching the protocol. This eliminates the need for most mapper functions in Phase 2 since the `pkg/` type IS the DTO.

### ETH Scope — Is Phase 2 Needed?

- **Context**: ETH `rpc_*.go` files already return domain types (`*big.Int`, `*domainETH.ResponseSyncing`, `*domainETH.BlockInfo`) directly. Unlike BTC, there is no `mapper.go` (only a small `converters.go` at 130 lines).
- **Findings**:
  - ETH Phase 1 (move `rpc_*.go` to `pkg/chains/eth/rpc/`) is straightforward.
  - ETH has `ResponseSyncing` defined as an infrastructure type in `rpc_eth.go`; it gets mapped to `domainETH.ResponseSyncing`. This is the main Phase 2 item.
  - `TxCreateParams` is defined in the port interface file itself (not a separate DTO file) — this can stay in the application layer.
- **Implications**: ETH Phase 2 is minimal — primarily moving `ResponseSyncing` to `pkg/` and updating the one port method that uses it.

### XRP — Library Types vs. pkg/ Types

- **Context**: XRP's `xrpapi_tx.go` (1,518 lines) wraps `github.com/peersyst/xrpl-go`. The peersyst types are used internally; the port interface returns `dtoxrp.*` types.
- **Findings**: XRP's RPC is WebSocket-based. The "RPC call functions" are the `public_account.go`, `admin_keygen.go`, and `public_transaction.go` files that call `wsPublic.Call()` / `wsAdmin.Call()`. The `xrpapi_tx.go` file is a higher-level orchestration layer (not pure RPC).
- **Implications**: Phase 1 for XRP should focus on the WebSocket call functions (`public_*.go`, `admin_keygen.go`). The `xrpapi_tx.go` orchestration layer stays in infrastructure initially. XRP peersyst types should not be exposed in `pkg/` — the `pkg/chains/xrp/rpc/` types should be self-defined structs matching the XRP WebSocket protocol.

### BCH Reuse of BTC RPC Package

- **Context**: BCH (`BitcoinCash`) embeds `Bitcoin` via struct embedding. BCH uses the same Bitcoin RPC protocol. The `FlexibleLabels` custom unmarshaler already handles BTC/BCH label format differences.
- **Findings**: BCH can call the same `pkg/chains/btc/rpc/` functions as BTC — no duplication needed. The `FlexibleLabels` type in `pkg/` handles the wire format difference transparently.
- **Implications**: BCH's `BitcoinCash` thin adapter calls `btcrpc.GetAddressInfo(b.Client, addr)` exactly like BTC. Override in `btc/bch/` only when BCH needs different domain mapping (not RPC call differences).

### Domain-Enriched vs. Pure-Wire BTC DTOs

- **Context**: BTC has 40+ DTO types. We need to classify which move to `pkg/` and which stay in the application/domain layer.
- **Findings** (by analyzing dto.go and mapper.go):

| DTO Type | Classification | Decision |
|----------|---------------|----------|
| `AddressInfo` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `ValidateAddressResult` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `NetworkInfo`, `BlockchainInfo` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `TransactionResult`, `RawTransaction` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `FundRawTransactionResult` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `DescriptorInfo`, `ListDescriptorsResult` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `ImportDescriptorsRequest/Response` | Pure wire (request payload) | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `MultisigAddress` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `LoggingResult` | Pure wire rename | Move to `pkg/chains/btc/rpc/`; DTO eliminated |
| `UnspentOutput` | Wire + domain `AccountType` | **Keep** in application/dto or domain layer |
| `ParsedPSBT`, `BIP32Derivation`, all `ParsedPSBT*` | Domain-enriched (PSBT parsing) | **Keep** in application/dto or domain layer |
| `PreviousTx` | PSBT input, domain context | **Keep** |

- **Implications**: Most BTC DTOs can be eliminated. Only PSBT-related and UTXO-with-domain-fields types stay.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| One-Phase Full Migration | Move all RPC code and update all port interfaces simultaneously | Clean end state, no transient duplication | ~300+ file changes, cannot ship incrementally, breaks all mocks at once | Too risky for a large codebase |
| Phase 1 only (RPC extraction) | Move RPC code to `pkg/`; infrastructure still converts to DTOs | Low risk, testable, shippable independently | Does not achieve DTO elimination yet | Good for validating `pkg/` API design |
| Two-Phase Hybrid (recommended) | Phase 1: Move RPC to `pkg/`; Phase 2: Eliminate DTOs + update port interfaces | Balanced risk, independently shippable phases, ETH done in Phase 1 | Two rounds of review | Best fit for project constraints |

---

## Design Decisions

### Decision: `RPCCaller` Interface Per Chain in `pkg/`

- **Context**: `pkg/chains/*/rpc/` functions need a client parameter but cannot import concrete third-party types without making those transitive dependencies of `pkg/`.
- **Alternatives Considered**:
  1. Accept concrete type (`*rpcclient.Client`) — tight coupling, cannot mock
  2. Define minimal interface in `pkg/chains/*/rpc/` — implicit satisfaction, mockable
- **Selected Approach**: Define a minimal interface per chain (e.g., `RPCCaller` for BTC/ETH, `WSCaller` for XRP). The concrete client types satisfy these interfaces without any changes.
- **Rationale**: Aligns with Go best practices ("accept interfaces, return structs"). Enables unit testing without a live node.
- **Trade-offs**: Adds one small interface type per chain package. Minor.

### Decision: Go-Idiomatic Field Names in `pkg/` Response Types

- **Context**: Current infrastructure RPC response types use JSON-protocol names (`Ismine`). Current DTOs use idiomatic names (`IsMine`). The mapper just renames fields.
- **Alternatives Considered**:
  1. Use protocol names in `pkg/` (e.g., `Ismine`) — faithful to wire format but unidiomatic Go
  2. Use idiomatic names with JSON tags in `pkg/` — clean Go API; JSON tags handle protocol faithfully
- **Selected Approach**: Use Go-idiomatic field names in `pkg/chains/*/rpc/` types with JSON struct tags matching the protocol.
- **Rationale**: Eliminates the need for mapper functions in Phase 2; `pkg/` types replace DTOs directly.
- **Trade-offs**: Diverges from raw protocol field names in struct definitions, but JSON serialization remains correct via tags.

### Decision: Two-Phase Migration

- **Context**: Moving RPC code AND reorganizing DTOs simultaneously creates a huge change surface.
- **Alternatives Considered**:
  1. All at once — high risk, unpredictable test failures
  2. Two phases — Phase 1 (RPC extraction), Phase 2 (DTO elimination + port interface updates)
- **Selected Approach**: Two phases. Phase 1 is pure code movement (no behavioral change). Phase 2 handles type system changes.
- **Rationale**: Phase 1 can be reviewed and shipped independently. Phase 2 can be done chain by chain. ETH starts first (simplest).
- **Trade-offs**: Phase 1 infrastructure adapters temporarily call `pkg/` functions AND convert to DTOs (transient duplication). This is acceptable since Phase 2 eliminates it.

### Decision: XRP `xrpapi_tx.go` Stays in Infrastructure for Phase 1

- **Context**: `xrpapi_tx.go` (1,518 lines) is a high-level orchestration layer wrapping the peersyst library, not a pure RPC call file.
- **Selected Approach**: Phase 1 moves only the direct WebSocket call files (`public_*.go`, `admin_keygen.go`) to `pkg/chains/xrp/rpc/`. The orchestration layer stays in infrastructure.
- **Rationale**: The orchestration layer involves complex business logic (transaction preparation, signing, combining) that is not "pure RPC call." Forcing it into `pkg/` would expose XRP domain logic in a public package.

---

## Risks & Mitigations

- **Port interface signature changes cascade to 170+ import sites** — Mitigation: Do Phase 2 chain by chain (ETH first, then BTC, then XRP). Regenerate mocks per chain.
- **Test displacement** — Some test files in `internal/infrastructure/api/` test RPC call functions that will move to `pkg/`. Mitigation: Co-locate tests with moved functions in `pkg/chains/*/rpc/`; update test imports.
- **BCH break risk** — BCH embeds BTC; if BTC's `GetAddressInfo` signature changes (Phase 2), BCH's method resolution changes too. Mitigation: BCH must override any method whose signature changes, placing BCH-specific conversion in `btc/bch/`.
- **`make mockery` must be re-run after Phase 2** — Mitigation: Re-run `make mockery` at end of each Phase 2 chain block; CI catches mock drift.

---

## References

- `internal/infrastructure/api/btc/btc/address.go` — Example: `GetAddressInfo()` pattern to move
- `internal/infrastructure/api/btc/btc/mapper.go` — 768-line mapper to eliminate in Phase 2
- `internal/infrastructure/api/eth/eth/rpc_eth_gas.go` — Minimal ETH RPC pattern
- `internal/application/ports/api/btc/interface.go` — Port interface method signatures (Phase 2 target)
- `pkg/CLAUDE.md` — Hard rule: `pkg/` cannot import from `internal/`
- `.claude/rules/internal/btc-shared.md` — BCH override pattern (BCH must override, not modify BTC)
- `.claude/rules/internal/clean-architecture.md` — Layer dependency rules
