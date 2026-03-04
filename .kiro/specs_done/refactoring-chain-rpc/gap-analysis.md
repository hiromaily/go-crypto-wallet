# Gap Analysis: refactoring-chain-rpc

## 1. Current State Investigation

### 1.1 Directory Layout & Asset Inventory

#### `pkg/chains/` — Current State (Tiny)

`pkg/chains/` currently contains **only RPC client factory functions**:

| File | Lines | Content |
|------|-------|---------|
| `pkg/chains/btc.go` | ~30 | `NewBitcoinRPCClient()` — creates `*rpcclient.Client` |
| `pkg/chains/eth.go` | ~27 | `NewEthereumRPCClient()` — creates `*ethrpc.Client` |
| `pkg/chains/xrp.go` | ~17 | `NewWebSocketClient()` — creates `*websocket.WS` |
| `pkg/chains/btc/` | ~5 files | BIP32 derivation, seed utilities |
| `pkg/chains/eth/eth.go` | ~1 file | ETH chain value objects |
| `pkg/chains/xrp/` | ~10 files | XRP address, keygen, sign, serialize |

**Gap:** The target `pkg/chains/{btc,eth,xrp}/rpc/` directories do not exist. There is no RPC call logic in `pkg/` at all.

#### `internal/infrastructure/api/` — Current State (Large)

| Chain | Files | Lines | Key Files |
|-------|-------|-------|-----------|
| BTC (`btc/btc/`) | 55 | ~11,000 | `address.go`, `transaction.go`, `unspent.go`, `psbt*.go`, `descriptor*.go`, `mapper.go` (768 lines) |
| ETH (`eth/eth/`) | 25 | ~3,500 | `rpc_eth.go`, `rpc_eth_tx.go`, `rpc_eth_gas.go`, `rpc_admin.go`, `rpc_personal.go`, `rpc_miner.go`, `rpc_net.go`, `rpc_web3.go`, `converters.go` (130 lines) |
| XRP (`xrp/`) | 39 | ~4,500 | `public_account.go`, `admin_keygen.go`, `xrpapi_tx.go` (1,518 lines), `converter.go` (635 lines) |

#### Port Interfaces — Current State

| Chain | Interface | Methods |
|-------|-----------|---------|
| BTC | `Bitcoiner` (monolithic, DI-only) + 25 focused | 173 methods in Bitcoiner |
| ETH | `Ethereumer` (monolithic, DI-only) + 15 focused | 50+ methods |
| XRP | `XRPer` (monolithic, DI-only) + 8 focused | 80+ methods |

**All port interface methods that return data currently return application DTOs** (e.g., `GetAddressInfo() (*dtobtc.AddressInfo, error)`).

#### DTO Layer — Current State

| Chain | DTO Package | Types |
|-------|-------------|-------|
| BTC | `internal/application/dto/btc/` | 40+ types (AddressInfo, ParsedPSBT, UnspentOutput, etc.) |
| ETH | Embedded in port (`TxCreateParams`) + `internal/domain/chains/eth/` | Mostly domain types |
| XRP | `internal/application/dto/xrp/` | 20+ types (TxInput, ResponseWalletPropose, etc.) |

**DTO usage in use cases:** ~170 import references across `internal/application/usecase/`.

---

### 1.2 Conventions & Patterns

**Critical architectural constraint (from `pkg/CLAUDE.md`):**
> Packages in `pkg/` MUST NOT import from `internal/` directory.

**Clean Architecture dependency rules (`internal/`):**
```
Infrastructure → application/ports + application/dto + domain
pkg/ → zero imports from internal/
```

**RPC call patterns by chain:**

| Chain | Client Type | Call Mechanism | Response Handling |
|-------|------------|----------------|-------------------|
| BTC | `*rpcclient.Client` | `Client.RawRequest(method, params)` | `json.Unmarshal → Infrastructure struct → mapper → DTO` |
| ETH | `*ethrpc.Client` + `*ethclient.Client` | `rpcClient.CallContext(ctx, &result, method)` | Hex decode + `json.Unmarshal → domain type` (mostly direct) |
| XRP | `*websocket.WS` (public + admin) | `ws.Call(ctx, &req, &res)` | Struct mapping → converter → DTO |

---

### 1.3 Integration Surfaces

**DI wiring (`internal/di/container.go`):**
```go
// 1. Create RPC clients using pkg/chains factories
rpcClient, _ := cryptocurrency.NewBitcoinRPCClient(conf)

// 2. Create infrastructure implementations
btc := apibtcimpl.NewBitcoin(rpcClient, conf...)

// 3. Store as monolithic port interface
type container struct { btc apibtc.Bitcoiner }
```

**Mock generation (`.mockery.yaml`):**
- Mocks generated from `internal/application/ports/api/` interfaces
- Placed in `internal/infrastructure/api/{chain}/mocks/`
- BTC: 1 monolithic mock; ETH: 15 focused mocks; XRP: 7 focused mocks

---

## 2. Requirements Feasibility Analysis

### 2.1 Requirement-to-Asset Map

| Req | Summary | Existing Assets | Gap | Classification |
|-----|---------|-----------------|-----|----------------|
| 1 | RPC Package Structure | `pkg/chains/` exists (tiny) | `pkg/chains/{btc,eth,xrp}/rpc/` missing | **Missing** |
| 2 | RPC Types in pkg | Types scattered in infra files | No types in `pkg/chains/*/rpc/` | **Missing** |
| 3 | Remove infra-layer DTO conversion | `mapper.go` (768L), `converter.go` (635L) | Must be moved/removed | **Constraint** |
| 4 | App-layer type conversion | DTOs in `internal/application/dto/` | Converters must move to app layer | **Missing** |
| 5 | Port Interface Preservation | All interfaces defined in `application/ports/` | Method signatures reference DTOs being changed | **Constraint** |
| 6 | Infrastructure as thin adapter | Infrastructure structs are full implementations | Must be thinned out | **Missing** |
| 7 | Multi-chain coverage | BTC (55 files), ETH (25 files), XRP (39 files) | All chains affected | **Constraint** |
| 8 | Behavioral equivalence | 50+ test files across chains | Tests co-located with infra; must stay passing | **Constraint** |

### 2.2 Critical Gaps

#### Gap 1: `pkg/chains/*/rpc/` Does Not Exist (Missing)
The entire `pkg/chains/{btc,eth,xrp}/rpc/` directory tree must be created from scratch. No scaffolding exists.

#### Gap 2: RPC Response Types Are Embedded in Infrastructure Files (Missing)
BTC RPC response structs (e.g., `GetAddressInfoResult`, `GetNetworkInfoResult`) are defined inline in their respective `*.go` files, not in a dedicated types file. They must be extracted and moved to `pkg/chains/btc/rpc/`.

#### Gap 3: Port Interface Signatures Reference Application DTOs (Constraint)
**All data-returning methods on port interfaces currently return application DTOs.** For example:
```go
// Current:
GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)

// Target (after refactoring):
GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error)
// OR: stays returning a domain type, with conversion moved to app layer
```

Changing these signatures requires updating every use case and mock that references them — potentially 170+ import sites.

#### Gap 4: Mapper/Converter Ownership Is Ambiguous (Unknown)
- `mapper.go` (768 lines) converts infrastructure types → application DTOs
- `converter.go` (635 lines, XRP) converts peersyst/xrpapi types → XRP DTOs
- **Research Needed:** Which conversions belong in the application layer (use-case-driven) vs. which are genuinely infrastructure concerns (e.g., raw wire format → Go struct)?

#### Gap 5: No Application-Layer Converters Exist (Missing)
The application layer currently has no converter functions. They must be created in `internal/application/` to perform `pkg/chains/*/rpc/` types → domain/use-case types.

#### Gap 6: ETH Already Uses Domain Types Directly (Constraint)
Ethereum is partially already following the target pattern — ETH port methods mostly return domain types (`*domainETH.ResponseSyncing`, `*big.Int`) rather than application DTOs. **Research Needed:** Clarify whether ETH still needs `pkg/chains/eth/rpc/` extraction or if the existing `rpc_*.go` files are already "thin enough."

#### Gap 7: XRP Uses Third-Party Library Types (Unknown)
XRP's `xrpapi_tx.go` (1,518 lines) heavily uses the `peersyst/xrpl-go` library's types internally. **Research Needed:** Whether these library types should be exposed in `pkg/chains/xrp/rpc/` or wrapped.

---

### 2.3 Complexity Signals

- **BTC:** Complex — 40+ DTO types, 768-line mapper, PSBT operations, BCH reuse concern
- **ETH:** Moderate — already somewhat domain-oriented, but 30+ RPC methods across 8 files
- **XRP:** Complex — 1,518-line transaction file, peersyst library wrapping, DTO-heavy port interface

---

## 3. Implementation Approach Options

### Option A: Full Migration — Move Everything at Once

Extract all RPC call functions and types to `pkg/chains/*/rpc/` in one pass, simultaneously update port interfaces, eliminate mappers, and add application-layer converters.

**Which files change:**
- Create: `pkg/chains/{btc,eth,xrp}/rpc/*.go` (new RPC packages)
- Modify: All 119 infrastructure files → thin adapters
- Modify: All 3 port interface files (method signatures)
- Modify: 170+ use case DTO import sites
- Modify: All 23 mock files (regenerate after interface changes)
- Delete: `mapper.go`, `converter.go`, portions of infrastructure type files

**Trade-offs:**
- ✅ Clean result; single consistent refactoring
- ✅ All circular dependency risks resolved at once
- ❌ Massive change surface (~300+ files); high risk of merge conflicts
- ❌ Cannot ship partial progress; all-or-nothing
- ❌ Port interface changes break all existing mocks until regenerated

---

### Option B: Chain-by-Chain Migration with Port Interface Stability

Migrate one chain at a time (ETH first as it's already partly there, then BTC, then XRP). Keep port interface signatures **stable** by having the thin infrastructure adapter still perform the final conversion — but delegating the RPC call itself to `pkg/chains/*/rpc/`.

**Strategy:**
1. Create `pkg/chains/*/rpc/` with standalone RPC call functions returning raw wire types
2. Infrastructure adapter calls the `pkg` function, then converts inline (simple wrapper)
3. Port interface signatures unchanged in this phase
4. In a later phase, push conversion up to the application layer

**Trade-offs:**
- ✅ Minimal disruption to use cases and mocks in phase 1
- ✅ Shippable incrementally per chain
- ✅ Enables validation of `pkg/chains/*/rpc/` API before committing to full migration
- ❌ Temporarily double-converts (RPC type → infra type → DTO) until phase 2
- ❌ Does not fully satisfy Req 3 (remove infra-layer DTO conversion) until phase 2

---

### Option C: Hybrid — Separate RPC Extraction from DTO Elimination (Recommended)

Split into two clearly sequenced sub-refactorings:

**Phase 1 — RPC Extraction (Req 1, 2, 6):**
- Create `pkg/chains/{btc,eth,xrp}/rpc/` packages with RPC call functions and wire types
- Infrastructure becomes thin adapters calling `pkg` functions
- Port interface signatures unchanged; adapters still convert to existing DTOs

**Phase 2 — DTO Layer Reorganization (Req 3, 4, 5):**
- Move mapper/converter logic from infrastructure to application layer
- Update port interface method signatures to return `pkg/chains/*/rpc/` types or domain types
- Eliminate application DTOs that are structurally identical to RPC response types
- Regenerate mocks; update use cases

**Why recommended:**
- Phase 1 is a pure move of RPC code with zero behavioral change — low risk, independently reviewable
- Phase 2 is the high-risk type-system change — isolated, can be done chain by chain
- Aligns with the project's existing chain-by-chain boundary structure
- ETH should be Phase 1 pilot (simplest; mostly domain types already)
- BTC has most complexity (mapper.go) — do last

**Trade-offs:**
- ✅ Lowest risk per phase
- ✅ Each phase independently reviewable and deployable
- ✅ Phase 1 satisfies Reqs 1, 2, 6; Phase 2 satisfies Reqs 3, 4, 5
- ❌ Requires two rounds of review and coordination
- ❌ Phase 1 infra adapters temporarily have both a `pkg` call and a DTO conversion (transient complexity)

---

## 4. Implementation Complexity & Risk

### Effort

**Total estimate: XL (2+ weeks)**

| Phase | Effort | Justification |
|-------|--------|---------------|
| Phase 1 — ETH RPC extraction | S (1–2 days) | 8 `rpc_*.go` files, mostly direct; ETH is already domain-oriented |
| Phase 1 — BTC RPC extraction | M (3–5 days) | 55 files, complex PSBT/descriptor logic, BCH reuse |
| Phase 1 — XRP RPC extraction | M (3–5 days) | 39 files, 1,518-line tx file, peersyst wrapping |
| Phase 2 — DTO layer (all chains) | L (1–2 weeks) | 170 import sites, interface signature changes, mock regeneration |

### Risk

**Overall risk: Medium-High**

| Risk Area | Level | Justification |
|-----------|-------|---------------|
| `pkg/` circular dependency | Low | Hard rule enforced; `pkg/` already has no `internal/` imports |
| Port interface stability | High | Changing 173 method signatures cascades to mocks, use cases, DI |
| BCH/BTC RPC reuse | Low | BCH already embeds `Bitcoin` struct; `pkg/chains/btc/rpc/` is naturally shared |
| ETH domain types already correct | Low | ETH minimal DTO usage means smaller surface for Phase 2 |
| Test displacement | Medium | 50+ test files co-located with infra; tests must follow moved logic |
| XRP peersyst library wrapping | Medium | Library types may not be suitable for direct exposure in `pkg/` |
| Mock regeneration | Medium | All mocks must be regenerated after Phase 2 interface changes |

---

## 5. Research Items for Design Phase

1. **ETH scope clarification:** Do `rpc_*.go` files in ETH need to move to `pkg/chains/eth/rpc/`, or are they already sufficiently separated? ETH port methods mostly return `*big.Int` and domain types — does Phase 2 even apply?

2. **XRP peersyst library types:** `xrpapi_tx.go` wraps `github.com/peersyst/xrpl-go`. Should `pkg/chains/xrp/rpc/` expose peersyst types directly, or define a thin wrapper layer?

3. **BTC RPC client interface in pkg:** The `pkg/chains/btc/rpc/` functions need to accept the BTC RPC client. Should this be `*rpcclient.Client` (concrete) or a minimal interface? An interface enables unit testing in `pkg/` without a live node.

4. **Test ownership after migration:** Where do unit tests for RPC functions live after moving to `pkg/`? They would need to move to `pkg/chains/{btc,eth,xrp}/rpc/` — confirm this is acceptable given mock setup overhead.

5. **DTO types that are genuinely domain-specific:** Some BTC DTOs (e.g., `ParsedPSBT`, `BIP32Derivation`) encode domain knowledge, not just RPC wire format. These should remain as domain types, not be eliminated. Need inventory of which BTC DTOs are "pure RPC wire" vs. "domain-enriched."

6. **`application/dto/` retention policy:** After Phase 2, some DTO packages may become empty or near-empty. Should the DTO packages be repurposed (renamed to `types`) or deprecated?

---

## 6. Recommendations for Design Phase

**Preferred approach:** Option C (Hybrid, two-phase)

**Key decisions for design phase:**
1. Define the public API surface of `pkg/chains/{btc,eth,xrp}/rpc/` — what functions, what types, what client parameter type
2. Specify which port interface method signatures change in Phase 2 (start with the simplest, e.g., `GetNetworkInfo`)
3. Define the application-layer converter pattern (function location, naming, test strategy)
4. Determine ETH scope — likely Phase 1 only with minimal Phase 2 impact
5. Design BCH reuse of `pkg/chains/btc/rpc/` without code duplication
6. Plan mock regeneration cadence — regenerate per chain or all at once at end of Phase 2

**Starting point (pilot):** Begin with ETH (simplest, least DTO coupling) to validate the `pkg/chains/*/rpc/` API design before applying to BTC and XRP.
