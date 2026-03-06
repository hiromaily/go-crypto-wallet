# Implementation Plan

- [x] 1. Delete all gRPC-only infrastructure files
- [x] 1.1 (P) Delete the five gRPC-only infrastructure implementation files
  - Remove `xrpapi_address.go` — all 3 methods depend on `r.API.AddressClient`
  - Remove `xrpapi_tx_escrow.go` — all escrow prepare methods depend on `r.API.TxClient`
  - Remove `xrpapi_tx_payment_channel.go` — all payment channel prepare methods depend on `r.API.TxClient`
  - Remove `xrpapi_tx_nftoken.go` — all NFToken prepare methods depend on `r.API.TxClient`
  - Remove `xrpapi_tx_account.go` — all account transaction methods depend on `r.API.TxClient` or `signTransactionJSON`
  - Use `git rm` so the deletions are tracked
  - _Requirements: 1.1_

- [x] 1.2 (P) Delete the integration test files and the XRPAPIProvider mock
  - Remove `xrpapi_address_test.go` — integration tests for deleted address methods
  - Remove `xrpapi_account_test.go` — integration tests referencing deleted methods
  - Remove `mocks/mock_xrpapi_provider.go` — auto-generated mock for the interface being removed
  - Use `git rm` for all three files
  - _Requirements: 1.2, 1.3_

- [x] 2. Strip gRPC code from partial implementation files
- [x] 2.1 (P) Remove the dead converter functions from the converter file
  - Delete `ToInfraInstructions()` and `ToDTOInstructions()` — convert to/from `protogen.Instructions`
  - Delete `ToDTOResponseGenerateAddress()` and `ToDTOResponseGenerateXAddress()` — only used by deleted `xrpapi_address.go`
  - Delete `ToDTOSignerListSetTxInput()` and `ToDTOTrustSetTxInput()` — source types defined in deleted `xrpapi_tx_account.go`
  - Delete `ToDTOEscrowCreateTxInput()`, `ToDTOEscrowFinishTxInput()`, `ToDTOEscrowCancelTxInput()` — source types defined in deleted `xrpapi_tx_escrow.go`
  - Delete `ToDTOPaymentChannelCreateTxInput()`, `ToDTOPaymentChannelFundTxInput()`, `ToDTOPaymentChannelClaimTxInput()` — source types defined in deleted `xrpapi_tx_payment_channel.go`
  - Delete all NFToken converter functions — source types defined in deleted `xrpapi_tx_nftoken.go`
  - Retain the 4 surviving functions: `ToInfraTxInput`, `ToDTOTxInput`, `ToInfraXRPKeyType`, `ToDTOXRPKeyType`
  - Remove the `protogen` import after deletions; only `dtoxrp` import should remain
  - _Requirements: 2.2, 2.3_

- [x] 2.2 (P) Remove the four gRPC-calling items from the transaction file
  - Delete the `signTransactionJSON()` private helper — only called by deleted gRPC sign methods
  - Delete `(*XRP).CombineTransaction()` — calls `r.API.TxClient.CombineTransaction()`
  - Delete `(*XRP).SignTransactionNative()` — stub returning error; only satisfies the deleted `XRPAPIProvider` interface
  - Delete the `unquoteJSON()` private helper — only called by gRPC Prepare* methods in deleted files
  - Retain all struct type definitions (`TxInput`, `SentTx`, `TxInfo`, and related types) — used by surviving WebSocket methods
  - Retain all WebSocket-based methods: `PrepareTransaction`, `SignTransaction`, `SubmitTransaction`, `WaitValidation`, `GetTransaction`, `toXRPClientSentTx`
  - _Requirements: 2.1, 2.3_

- [x] 3. Remove the gRPC client from the XRP struct and its constructor chain
- [x] 3.1 (P) Strip the API field and its cleanup call from the main XRP struct
  - Remove the `API *xrplclient.XRPLClient` field from the `XRP` struct definition
  - Remove the `api *xrplclient.XRPLClient` parameter from the `NewXRP()` constructor
  - Remove the `r.API.Close()` call from the `Close()` method
  - Remove the `xrplclient` import from the file
  - _Requirements: 3.1_

- [x] 3.2 (P) Remove the gRPC client parameter from the connection factory
  - Remove the `api *xrplclient.XRPLClient` parameter from `NewXRPFromCoinType()`
  - Update the internal `NewXRP()` call to match the new signature without the gRPC argument
  - Remove the `xrplclient` import from the file
  - _Requirements: 3.2_

- [x] 3.3 (P) Remove gRPC connection setup from the test utility
  - Remove `grpc.NewClient` call and related gRPC connection setup
  - Remove `xrplclient.NewXRPLClient` call
  - Update the `NewXRPFromCoinType` call to no longer pass a gRPC client argument
  - Remove `xrplclient` and `pkg/grpc` imports from the file
  - _Requirements: 3.3_

- [x] 4. Remove gRPC wiring from the DI container
  - Remove the `xrpAPI *xrplclient.XRPLClient` field from the container struct
  - Remove the `newXRPAPI()` factory method entirely
  - Update the `newXRP()` call to no longer pass a gRPC client to `NewXRPFromCoinType`
  - Remove the `xrplclient` import; verify no other XRP gRPC calls remain in container
  - Depends on Task 3.2 completing the updated constructor signature
  - _Requirements: 4.1, 4.2_

- [x] 5. Remove XRPAPIProvider from the port interface layer
- [x] 5.1 (P) Remove the XRPAPIProvider interface definition and its composite embedding
  - Delete the entire `XRPAPIProvider` interface definition from the ports interface file
  - Remove the `XRPAPIProvider` embedding from the `XRPer` composite interface
  - Check `doc.go` for any stale comment references to `XRPAPIProvider` and remove them
  - Retain all focused interfaces: `AccountInfoProvider`, `BalanceChecker`, `TransactionSubmitter`, `TransactionPreparer`, `TransactionCombiner`, `RegularKeyPreparer`, `SignerListPreparer`, `KeyGenerator`, `Closer`, `XRPPublicer`, `XRPAdminer`, `CoinTypeProvider`
  - _Requirements: 5.1, 5.3_

- [x] 5.2 (P) Remove XRPAPIProvider from the mockery configuration
  - Remove the `XRPAPIProvider` entry from the XRP section of `.mockery.yaml`
  - Verify no other mockery targets reference the removed interface
  - _Requirements: 5.2_

- [x] 6. Verify build, lint, and tests pass
  - Run `make check-build` and confirm zero compilation errors
  - Run `make go-lint` and confirm no new lint errors introduced
  - Run `make go-test` and confirm all unit tests under `internal/` pass
  - Run `grep -r "xrplclient\|protogen" internal/` and confirm zero results
  - _Requirements: 6.1, 6.2, 6.3, 6.4_
