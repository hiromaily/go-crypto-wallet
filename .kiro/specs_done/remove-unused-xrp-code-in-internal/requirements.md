# Requirements Document

## Project Description (Input)

`apps/xrpl-grpc-server` is a gRPC server providing XRP Ledger functionality
(`XRPTransactionAPIClient`, `XRPAccountAPIClient`, `XRPAddressAPIClient`). This server is not yet
operational, so all code under `internal/` that calls it via `*xrplclient.XRPLClient` is dead code.
The goal is to **delete only that dead gRPC-dependent logic from `internal/`**, leaving the working
WebSocket-based implementations untouched. `apps/xrpl-grpc-server` itself and `pkg/` are out of scope.

### Dead Code Identified

The following code under `internal/` is dead because it makes calls through `r.API.*` (the
`*xrplclient.XRPLClient` gRPC client) which requires the non-operational server:

**Entire files (all methods gRPC-only):**

- `internal/infrastructure/api/xrp/xrpapi_address.go` — `GenerateAddress`, `GenerateXAddress`, `IsValidAddress` (all via `r.API.AddressClient`)
- `internal/infrastructure/api/xrp/xrpapi_tx_escrow.go` — `PrepareEscrowCreate/Finish/CancelTransaction` (all via `r.API.TxClient`)
- `internal/infrastructure/api/xrp/xrpapi_tx_payment_channel.go` — all payment channel prepare methods (via `r.API.TxClient`)
- `internal/infrastructure/api/xrp/xrpapi_tx_nftoken.go` — all NFToken prepare methods (via `r.API.TxClient`)

**Partial files (gRPC methods within otherwise WebSocket files):**

- `internal/infrastructure/api/xrp/xrpapi_tx.go` — `signTransactionJSON()`, `CombineTransaction()`, `SignTransactionNative()` (stub) use `r.API.TxClient`; WebSocket methods `PrepareTransaction`, `SubmitTransaction`, `WaitValidation`, `GetTransaction` are NOT dead
- `internal/infrastructure/api/xrp/xrpapi_tx_account.go` — all `Prepare*` and `Sign*` methods use `r.API.TxClient` or `signTransactionJSON()`; entire file is dead
- `internal/infrastructure/api/xrp/converter.go` — `ToInfraInstructions`, `ToDTOInstructions` convert to/from `protogen.Instructions` (only needed by deleted methods); other converters may be kept if still referenced

**Supporting infrastructure:**

- `internal/infrastructure/api/xrp/xrp.go` — `API *xrplclient.XRPLClient` field and its `r.API.Close()` call in `Close()`
- `internal/infrastructure/api/xrp/connection.go` — `api *xrplclient.XRPLClient` parameter in `NewXRPFromCoinType` and passed to `NewXRP`
- `internal/infrastructure/api/xrp/testutil/xrp.go` — creates gRPC connection and passes to `NewXRPFromCoinType`
- `internal/di/container.go` — `xrpAPI *xrplclient.XRPLClient` field and `newXRPAPI()` factory method
- `internal/application/ports/api/xrp/interface.go` — `XRPAPIProvider` interface (composites all the gRPC-backed port methods)

**Mocks to remove:**

- `internal/infrastructure/api/xrp/mocks/mock_xrpapi_provider.go` — mock for `XRPAPIProvider`

**Tests to remove:**

- `internal/infrastructure/api/xrp/xrpapi_address_test.go` — integration tests for deleted address methods
- `internal/infrastructure/api/xrp/xrpapi_account_test.go` — integration tests referencing deleted methods

### Preserved Code

The following WebSocket-based code is NOT dead and must be preserved:

- `WSClient.PrepareTransaction` — uses `xrprpc.AccountInfo` WebSocket
- `WSClient.SubmitTransaction` — uses `xrprpc.Submit` WebSocket
- `WSClient.WaitValidation` — uses `xrprpc.LedgerCurrent` / `LedgerAccept` WebSocket
- `WSClient.GetTransaction` — uses `xrprpc.GetTx` WebSocket
- `XRP.SignTransaction` — uses `PeersystSigner` (offline, no gRPC)
- All `public_*.go` and `balance.go` methods (WebSocket-based)
- `signer/` package
- Focused port interfaces: `AccountInfoProvider`, `BalanceChecker`, `TransactionSubmitter`, `TransactionPreparer`, `TransactionCombiner`, `RegularKeyPreparer`, `SignerListPreparer`, `KeyGenerator`, `Closer`, `XRPPublicer`, `XRPAdminer`

---

## Requirements

### Requirement 1: Delete gRPC-Only Infrastructure Files

**Objective**: Remove all files whose content is entirely dead gRPC-calling code.

#### Acceptance Criteria

1. The following files shall be deleted from `internal/infrastructure/api/xrp/`:
   - `xrpapi_address.go`
   - `xrpapi_tx_escrow.go`
   - `xrpapi_tx_payment_channel.go`
   - `xrpapi_tx_nftoken.go`
   - `xrpapi_tx_account.go`
2. The following test files that test only the deleted methods shall be deleted:
   - `xrpapi_address_test.go`
   - `xrpapi_account_test.go`
3. The mock file for `XRPAPIProvider` shall be deleted:
   - `mocks/mock_xrpapi_provider.go`

---

### Requirement 2: Remove gRPC Methods from Partial Files

**Objective**: Strip the gRPC-calling methods from files that also contain valid WebSocket code.

#### Acceptance Criteria

1. From `internal/infrastructure/api/xrp/xrpapi_tx.go`, remove:
   - The `signTransactionJSON()` helper (used only by gRPC sign methods)
   - `(*XRP).CombineTransaction()` (calls `r.API.TxClient.CombineTransaction`)
   - `(*XRP).SignTransactionNative()` (stub returning error; only needed by `XRPAPIProvider`)
   - Local types `TxInput`, `SentTx`, `TxInfo`, `TxSpecification`, `TxSpecSource`, `TxAmount`, `TxTotalPrice`, `TxSpecDestination`, `TxOutcome`, `TxOrderbookChange` — delete if no longer referenced after deletions above; keep if still used by remaining WebSocket methods
   - The `unquoteJSON()` helper — delete if only used by deleted gRPC methods
2. From `internal/infrastructure/api/xrp/converter.go`, remove:
   - `ToInfraInstructions()` (converts to `protogen.Instructions`)
   - `ToDTOInstructions()` (converts from `protogen.Instructions`)
   - Any other converter functions that are only referenced by deleted files
   - Keep converter functions still referenced by surviving code (e.g., `ToDTOTxInput`, `toXRPClientSentTx` if needed)
3. The `import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"` shall be removed from all files in `internal/` after the above deletions.

---

### Requirement 3: Remove gRPC Client from XRP Struct and Constructors

**Objective**: Strip the `API *xrplclient.XRPLClient` field from `XRP` so the struct no longer holds a gRPC client.

#### Acceptance Criteria

1. In `internal/infrastructure/api/xrp/xrp.go`:
   - Remove the `API *xrplclient.XRPLClient` field from the `XRP` struct
   - Remove the `api *xrplclient.XRPLClient` parameter from `NewXRP()`
   - Remove `r.API.Close()` from the `Close()` method
   - Remove the `xrplclient` import
2. In `internal/infrastructure/api/xrp/connection.go`:
   - Remove the `api *xrplclient.XRPLClient` parameter from `NewXRPFromCoinType()`
   - Update the internal `NewXRP()` call accordingly
   - Remove the `xrplclient` import
3. In `internal/infrastructure/api/xrp/testutil/xrp.go`:
   - Remove gRPC connection setup (`grpc.NewClient`, `xrplclient.NewXRPLClient`)
   - Update `GetXRPPublicClient()` call to `NewXRPFromCoinType` without gRPC parameter
   - Remove `xrplclient` and `pkg/grpc` imports

---

### Requirement 4: Remove gRPC Client from DI Container

**Objective**: Remove the gRPC client lifecycle from the dependency injection container.

#### Acceptance Criteria

1. In `internal/di/container.go`:
   - Remove the `xrpAPI *xrplclient.XRPLClient` field
   - Remove the `newXRPAPI()` factory method
   - Update `newXRP()` to not pass a gRPC client to `NewXRPFromCoinType`
   - Remove the `xrplclient` import
2. The `pkgContainer.NewGRPCClient()` call for XRP shall be removed; if no other callers remain, the gRPC client construction for XRP is eliminated entirely.

---

### Requirement 5: Remove `XRPAPIProvider` Interface

**Objective**: Delete the monolithic `XRPAPIProvider` port interface which only aggregated the now-deleted gRPC-backed methods.

#### Acceptance Criteria

1. In `internal/application/ports/api/xrp/interface.go`:
   - Remove the `XRPAPIProvider` interface definition
   - Remove the `XRPAPIProvider` embedding from the `XRPer` composite interface
2. Update `.mockery.yaml` to remove `XRPAPIProvider` from the mockery configuration if it is still listed there.
3. The focused interfaces that are still valid and implemented by WebSocket code (`AccountInfoProvider`, `BalanceChecker`, `TransactionSubmitter`, `TransactionPreparer`, `TransactionCombiner`, `RegularKeyPreparer`, `SignerListPreparer`, `KeyGenerator`, `Closer`, `XRPPublicer`, `XRPAdminer`, `CoinTypeProvider`) shall remain unchanged.

---

### Requirement 6: Build and Lint Pass After Deletion

**Objective**: Ensure the codebase compiles and lints cleanly after all deletions.

#### Acceptance Criteria

1. `make check-build` shall pass with no compilation errors after all deletions.
2. `make go-lint` shall produce no new lint errors after all deletions.
3. `make go-test` shall pass for all unit tests under `internal/` (integration tests requiring the gRPC server are excluded).
4. No `import` of `github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplclient` or `github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen` shall remain in any file under `internal/` after the deletions (those packages themselves are NOT deleted — only their consumers in `internal/`).
