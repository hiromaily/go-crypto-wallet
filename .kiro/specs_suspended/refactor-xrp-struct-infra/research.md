# Research Log — refactor-xrp-struct-infra

## Summary

**Discovery type**: Light (extension of existing system)

**Scope**: Full analysis of XRP infrastructure layer in `internal/infrastructure/api/xrp/`, port interfaces in `internal/application/ports/api/xrp/`, use case wiring in `internal/application/usecase/*/xrp/`, DI container in `internal/di/container.go`, and interface-adapters in `internal/interface-adapters/wallet/xrp/`.

**Key findings**:
1. Use cases already apply ISP locally — each defines a narrow local interface (e.g., `xrpMonitorClient`) and the DI container passes the monolithic `XRP` through `apixrp.XRPer`. Decomposition at the infrastructure level will eliminate this over-provision.
2. WebSocket methods are already split between public and admin connections inside `WSClient`; the separation is a matter of creating two distinct exported structs.
3. `PrepareTransaction`, `SignTransaction`, and `CombineTransaction` currently use gRPC (`r.API.TxClient`). `PrepareTransaction` has an existing WebSocket-based implementation on `WSClient` that serves as the non-gRPC path.

---

## Research Log

### Topic 1 — Current infrastructure method distribution

| Method | Receiver | Transport |
|--------|----------|-----------|
| `GetAccountInfo` | `WSClient` | Public WebSocket (xrprpc.AccountInfo) |
| `GetBalance`, `GetTotalBalance` | `WSClient` | Public WebSocket (via GetAccountInfo) |
| `AccountChannels`, `AccountInfo`, `ServerInfo` | `WSClient` | Public WebSocket (rpc package) |
| `SubmitTransaction` | `WSClient` | Public WebSocket (xrprpc.Submit) |
| `WaitValidation` | `WSClient` | Public WebSocket + Admin WebSocket (ledger_accept, optional) |
| `GetTransaction` | `WSClient` | Public WebSocket (xrprpc.Tx) |
| `PrepareTransaction` | `WSClient` | Public WebSocket (account_info for sequence) |
| `ValidationCreate`, `WalletProposeWithKey`, `WalletPropose` | `WSClient` | Admin WebSocket |
| `GenerateAddress`, `GenerateXAddress`, `IsValidAddress` | `XRP` | gRPC → `API.AddressClient` |
| `CombineTransaction` | `XRP` | gRPC → `API.TxClient` |
| `SignTransaction` | `XRP` | Offline (PeersystSigner) |
| `SignTransactionNative` | `XRP` | Stub (not implemented) |
| `CreateRawTransaction` | `XRP` | Delegates to WSClient.GetAccountInfo + WSClient.PrepareTransaction |
| `PrepareSetRegularKeyTransaction`, `PrepareSignerListSetTransaction`, etc. | `XRP` | gRPC → `API.TxClient` |

**Implication**: `WSClient` already has all public/admin separation. The `XRP` struct adds gRPC delegation on top.

---

### Topic 2 — Port interface usage in use cases

| Use case file | Port interface used |
|--------------|---------------------|
| `watch/xrp/create_transaction.go` | `AccountInfoProvider`, `TransactionPreparer` (separate params) |
| `watch/xrp/monitor_transaction.go` | `BalanceChecker` (via local interface) |
| `watch/xrp/send_transaction.go` | `TransactionSubmitter` |
| `watch/xrp/set_regular_key.go` | `RegularKeyPreparer` |
| `watch/xrp/set_signer_list.go` | `SignerListPreparer` |
| `watch/xrp/create_multisig_tx.go` | `TransactionPreparer` |
| `watch/xrp/add_multisig_signature.go` | `TransactionCombiner` |
| `watch/xrp/submit_multisig_tx.go` | `TransactionSubmitter` |
| `keygen/xrp/generate_key.go` | `KeyGenerator` |
| `keygen/xrp/sign_transaction.go` | `TransactionSigner` (via local interface) |

**All use cases already use narrow interfaces** — no use case takes `XRPer` directly as its dependency type. The DI container passes a concrete `*XRP` that satisfies the narrow interface implicitly.

---

### Topic 3 — DI container current wiring

`internal/di/container.go`:
- Holds `xrp apixrp.XRPer` as a single field (cached instance)
- `newXRP()` creates `*apixrpimpl.XRP` via `NewXRPFromCoinType(wsPublic, wsAdmin, xrpAPI, conf, coinType)`
- Same instance passed to all XRP use case factories (`c.newXRP()` is called per use case factory method)
- `wsXrpPublic`, `wsXrpAdmin`, `xrpAPI` cached as separate fields

**After decomposition**:
- `xrp apixrp.XRPer` field and `newXRP()` are removed
- 5 separate cached fields are added, one per struct
- Each use case factory receives only the relevant interface(s)

---

### Topic 4 — Interface-adapters holding XRPer

`internal/interface-adapters/wallet/xrp/keygen.go`:
- `XRPKeygen.XRP apixrp.XRPer` — used only for `XRP.Close()` in `Done()`

`internal/interface-adapters/wallet/xrp/watch.go` (not yet read but from earlier sessions):
- Holds `XRP apixrp.XRPer` for `Close()`

**Impact**: These structs need `XRP apixrp.XRPer` replaced with `apixrp.Closer` or individual port interfaces. Alternatively `Done()` can be updated to close each connection separately.

---

### Topic 5 — gRPC backward compatibility path

`pkg/chains/xrp/xrplclient/client.go` is not modified per requirements.

`xrplclient.XRPLClient` exposes:
- `TxClient protogen.XRPTransactionAPIClient`
- `AccountClient protogen.XRPAccountAPIClient`
- `AddressClient protogen.XRPAddressAPIClient`

These three fields can be passed directly to the corresponding new struct constructors:
- `NewTxClient(xrplClient.TxClient, ...)`
- `NewAccountClient(xrplClient.AccountClient)`
- `NewAddressClient(xrplClient.AddressClient)`

Switching gRPC on/off requires only DI wiring changes.

---

### Topic 6 — Non-gRPC local implementations

**txClient (transaction preparation, signing, combining)**:
- `PrepareTransaction` → existing `WSClient.PrepareTransaction` logic (queries `xrprpc.AccountInfo`)
- `SignTransaction` → existing `PeersystSigner.SignTransactionNative` (offline, `Peersyst/xrpl-go`)
- `CombineTransaction` → local multisig combine using `binary-codec` from `Peersyst/xrpl-go`
- `Prepare*` (SetRegularKey, SignerListSet, etc.) → local tx building with `Peersyst/xrpl-go`

**accountClient (account queries)**:
- `GetAccountInfo` → existing `WSClient.GetAccountInfo` logic (queries `xrprpc.AccountInfo`)
- `GetBalance`, `GetTotalBalance` → existing `WSClient.GetBalance` logic

**addressClient (address generation/validation)**:
- `GenerateAddress`, `GenerateXAddress` → existing `pkg/chains/xrp` keygen utilities
- `IsValidAddress` → existing `pkg/chains/xrp/address.go` validation

**Library compliance** (`.claude/rules/chains/xrp/3rd-library.md`):
- Offline ops → `github.com/Peersyst/xrpl-go`
- WebSocket node queries → `pkg/chains/xrp/xrplgo` or `xrprpc`
- `github.com/XRPLF/xrpl-go` MUST NOT be added

---

### Topic 7 — `WaitValidation` admin dependency

Current `WSClient.WaitValidation` optionally uses the admin connection to call `ledger_accept` in standalone mode. After decomposition:
- `publicClient` handles `WaitValidation` using only the public connection
- The admin ledger_accept call is removed (non-standalone environments not affected)
- Standalone mode requires manual ledger advancement outside the wallet binary

This is an acceptable trade-off since the admin part is already guarded by `if w.admin != nil` and errors are silently ignored.

---

### Topic 8 — Port interface obsolescence

After decomposition:
- `XRPer` (monolithic) → **eliminated**
- `XRPAPIProvider` (sub-monolithic) → **eliminated** (replaced by per-struct narrow interfaces)
- `XRPPublicer` → **retained** (implemented by `publicClient`)
- `XRPAdminer` → **retained** (implemented by `adminClient`)
- `AccountInfoProvider`, `BalanceChecker`, `TransactionSubmitter`, `TransactionPreparer`, `TransactionCombiner`, `RegularKeyPreparer`, `SignerListPreparer`, `KeyGenerator`, `Closer` → **retained** (implemented by appropriate new structs)
- New composed interfaces may be added to `interface.go` if multiple use cases need the same combination

---

## Architecture Decision Records

### ADR-1: Adapter pattern for txClient, accountClient, addressClient
**Decision**: Each struct wraps a `protogen.*APIClient` field and exposes port interface methods through type conversion.
**Rationale**: Allows DI container to inject either gRPC client or a local implementation without changing the struct.
**Alternative rejected**: Embedding the protogen interface directly would expose gRPC methods to use cases, violating ISP.

### ADR-2: Local non-gRPC implementations as separate structs
**Decision**: `localTxImpl`, `localAccountImpl`, `localAddressImpl` implement `protogen.*APIClient` locally.
**Rationale**: Keeps the adapter structs thin; swapping between gRPC and local requires only DI wiring change.

### ADR-3: publicClient holds only public WebSocket
**Decision**: `publicClient` wraps only `public *websocket.WS`; admin ledger_accept in `WaitValidation` is removed.
**Rationale**: Strict single-responsibility; the ledger_accept is a non-essential standalone optimization.

### ADR-4: txClient depends on AccountInfoProvider for CreateRawTransaction
**Decision**: `txClient` constructor accepts `apixrp.AccountInfoProvider` to enable the balance-check logic in `CreateRawTransaction`.
**Rationale**: `CreateRawTransaction` must validate sender balance before preparing the transaction; the balance data comes from the public node (publicClient).
