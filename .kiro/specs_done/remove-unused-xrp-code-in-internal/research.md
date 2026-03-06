# Research & Design Decisions

---

## Summary

- **Feature**: `remove-unused-xrp-code-in-internal`
- **Discovery Scope**: Simple Addition (pure deletion/cleanup, no new components)
- **Key Findings**:
  - All dead code is concentrated in `internal/infrastructure/api/xrp/` and two support files (`di/container.go`, `application/ports/api/xrp/interface.go`)
  - `converter.go` is almost entirely deletable — only 4 of its 15+ functions survive (those that convert types still used by WebSocket-based methods)
  - The `XRP` struct itself survives after stripping the `API` field; the overall `XRP`/`WSClient` architecture is preserved unchanged

## Research Log

### Dead Code Classification

- **Context**: Identify which files/methods to delete vs. preserve.
- **Findings**:
  - **100% dead files** (all methods call `r.API.*`): `xrpapi_address.go`, `xrpapi_tx_escrow.go`, `xrpapi_tx_payment_channel.go`, `xrpapi_tx_nftoken.go`, `xrpapi_tx_account.go`
  - **Partially dead**: `xrpapi_tx.go` — 3 gRPC methods (`signTransactionJSON`, `CombineTransaction`, `SignTransactionNative`) are dead; 4 WebSocket methods survive
  - **converter.go** — 2 functions (`ToInfraInstructions`, `ToDTOInstructions`) use `protogen` and are deleted; 10+ DTO-to-local-struct converters for deleted struct types (EscrowCreate, PaymentChannel, NFToken, SignerListSet, TrustSet, ResponseGenerateAddress) are deleted because their source struct types no longer exist; 4 converters survive (`ToInfraTxInput`, `ToDTOTxInput`, `ToInfraXRPKeyType`, `ToDTOXRPKeyType`)
- **Implications**: `converter.go` after deletion will only import `dtoxrp` (no `protogen` import needed)

### Surviving Types in xrpapi_tx.go

- **Context**: `xrpapi_tx.go` defines local struct types used by both WebSocket and gRPC methods; need to determine which survive.
- **Findings**:
  - `TxInput` — used by `PrepareTransaction` (builds it) and `SubmitTransaction` (maps from response) → **KEEP**
  - `SentTx` — used by `SubmitTransaction` → **KEEP**
  - `TxInfo`, `TxSpecification`, `TxSpecSource`, `TxAmount`, `TxTotalPrice`, `TxSpecDestination`, `TxOutcome`, `TxOrderbookChange` — only referenced by `GetTransaction` return mapping → **KEEP** (but simplify if only used as intermediate in one method)
  - `unquoteJSON()` — only called by gRPC `Prepare*` methods → **DELETE**
- **Implications**: `xrpapi_tx.go` shrinks by ~3 methods but retains its struct definitions

### XRPAPIProvider Interface Surface

- **Context**: `XRPAPIProvider` embeds into `XRPer`. Need to verify no use case depends on `XRPAPIProvider` directly.
- **Findings**:
  - Grep shows `XRPAPIProvider` appears in: `interface.go` (definition), `doc.go` (comment), `mock_xrpapi_provider.go` (auto-generated mock). No use case files reference `XRPAPIProvider` directly — they all use focused interfaces.
  - `XRPer` embeds `XRPAPIProvider` — after removal, `XRPer` will no longer embed it; `XRPer` itself may also become removable in the subsequent `refactor-xrp-struct-infra` spec.
- **Implications**: Safe to remove `XRPAPIProvider` from `interface.go` without breaking use cases

### DI Container gRPC wiring

- **Context**: `container.go` has `newXRPAPI()` which calls `pkgContainer.NewGRPCClient()`.
- **Findings**:
  - `newXRPAPI()` is only called inside `newXRP()`. After removing it, `newXRP()` no longer needs a gRPC client.
  - `NewGRPCClient()` for XRP is called only in this one place; after deletion, no XRP gRPC connection is established at all.
  - `newXRP()` passes `c.newXRPAPI()` as third argument to `NewXRPFromCoinType` → that parameter disappears.

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations |
|--------|-------------|-----------|---------------------|
| Delete files + strip methods | Remove entire dead files; edit partial files to remove gRPC methods only | Minimal surface area change; preserves working code | Risk: missing a cross-file reference (handled by build verification) |
| Replace with stubs | Keep file structure, stub out gRPC methods with `return errors.New("not implemented")` | No cascading deletions | Leaves dead code; contradicts the goal; keeps `XRPAPIProvider` interface in place |

**Selected**: Delete files + strip methods. Stubs are already present in some cases and still considered dead code because the port interface they satisfy (`XRPAPIProvider`) is itself being deleted.

## Design Decisions

### Decision: Delete `converter.go` Functions Whose Source Types Disappear

- **Context**: Many converter functions in `converter.go` convert from local infra struct types (e.g., `EscrowCreateTxInput`, `SignerListSetTxInput`) that are defined in files being deleted.
- **Selected Approach**: Delete all converter functions whose source or destination struct type is defined in a deleted file. After deletion, `converter.go` retains only the 4 converters that reference surviving types.
- **Rationale**: Go compiler will catch any missed reference, so deleting aggressively and fixing build errors is safe.
- **Trade-offs**: Converter functions that might be useful later (e.g., for future re-implementation) are deleted; they can be re-added when needed.

### Decision: Keep `XRP` Struct (Strip Field, Not Delete Struct)

- **Context**: The `XRP` struct could be deleted entirely, but it still holds `WSClient` (WebSocket) and implements `SignTransaction`.
- **Selected Approach**: Remove only the `API *xrplclient.XRPLClient` field. The struct itself and all WebSocket-delegating methods survive.
- **Rationale**: This spec's scope is only dead gRPC code. The `XRP` struct elimination is handled by `refactor-xrp-struct-infra`.

## Risks & Mitigations

- **Missing reference in converter.go** — a converter function for a type in a deleted file might be referenced elsewhere. Mitigation: `make check-build` will catch at compile time.
- **mockery.yaml out of sync** — if `XRPAPIProvider` is still listed in `.mockery.yaml` after deletion, `make mockery` would fail. Mitigation: remove it from the config file as part of Requirement 5.
- **Integration tests referencing gRPC** — `xrpapi_address_test.go` and `xrpapi_account_test.go` are `//go:build integration` tagged; they test deleted methods. Mitigation: delete the test files entirely.
