# Implementation Plan

> **Already Implemented (no tasks required)**
>
> - **FR-3** (Offline multisig signing): `signTransactionUseCase` with `existingSignedBlob` serial accumulation is complete; `SignTransactionNative(isMultiSig=true)` produces the correct XRPL `Signers` array.
> - **FR-5** (Signer key provisioning): Existing `create seed` / `create hdkey` keygen commands are reused for signer account setup; no new keygen infrastructure is needed.

---

- [x] 1. (P) Add the ability to prepare an unsigned `SignerListSet` transaction for a given XRP account
  - Fetch the account's current sequence number and ledger index via the existing account_info WebSocket RPC (same pattern as the existing payment transaction preparation)
  - Construct the SignerListSet transaction body with the signer entries array, quorum, fee, sequence, and last ledger sequence
  - Fee calculation: use `(N+1) * 12` drops where N is the number of signer entries, matching the XRPL protocol minimum for SignerListSet
  - Return both the typed DTO and the JSON-serialized transaction string; no signing is performed here
  - Satisfies the `SignerListPreparer` port interface already defined in ports/api/xrp/
  - _Requirements: 1_

- [x] 2. (P) Implement the keygen DB repository adapters for signer list data

- [x] 2.1 (P) Build the signer list aggregate repository backed by keygen DB SQLC queries
  - Wrap all SQLC-generated keygen queries for `xrp_signer_list`; implement every method of the `XRPSignerListRepositorier` port interface
  - Convert SQLC row types to `XRPSignerList` domain entities using private conversion functions; never expose SQLC types outside the repository
  - Implement for all three DB drivers: postgres, mysql, and sqlite (three separate adapter files following the existing cold repository layout)
  - The deactivate-then-insert sequence is non-atomic; document this as an accepted CLI limitation in code comments
  - _Requirements: 1, 6_

- [x] 2.2 (P) Build the signer entry repository backed by keygen DB SQLC queries
  - Wrap all SQLC-generated keygen queries for `xrp_signer_entry`; implement every method of the `XRPSignerEntryRepositorier` port interface
  - Convert rows to `XRPSignerEntry` domain entities; list ID is passed explicitly (no cross-table join needed)
  - Implement for all three DB drivers, collocated with the signer list adapters
  - _Requirements: 1, 6_

- [x] 3. (P) Implement the watch DB repository adapters for pending multisig state

- [x] 3.1 (P) Build the pending multisig transaction repository backed by watch DB SQLC queries
  - Wrap all SQLC-generated watch queries for `xrp_pending_multisig`; implement every method of the `XRPPendingMultisigRepositorier` port interface
  - Convert rows to `XRPPendingMultisig` domain entities; implement for all three DB drivers under the watch repository layout
  - These adapters enable `CreateMultisigTxUseCase` and `SubmitMultisigTxUseCase` (already wired in DI) to compile and run against real DB
  - _Requirements: 6_

- [x] 3.2 (P) Build the multisig signature repository backed by watch DB SQLC queries
  - Wrap all SQLC-generated watch queries for `xrp_multisig_signature`; implement every method of the `XRPMultisigSignatureRepositorier` port interface
  - Convert rows to `XRPMultisigSignature` domain entities; implement for all three DB drivers, collocated with the pending multisig adapters
  - Required for `AddMultisigSignatureUseCase` to compile (the use case is currently stubbed but the repository must still satisfy the interface)
  - _Requirements: 6_

- [ ] 4. (P) Extend the Watch wallet transaction creation flow to support multisig JSON-format files
  - Add a `MultisigQuorum` field to the create transaction input (zero value keeps the existing single-sig text-format path fully unchanged)
  - When `MultisigQuorum >= 2`, write an `XRPTransactionFile` JSON file instead of the legacy text file; set `required_signatures` to the quorum value, `signature_count` to 0, and `is_complete` to false
  - The unsigned payment transaction JSON (without `Signers` array) is identical between single-sig and multisig paths
  - The existing `generateHexFile()` text-format path is untouched; the new JSON path is a parallel branch invoked only when the quorum field is set
  - _Requirements: 2_

- [ ] 5. (P) Extend the Watch wallet send flow to detect and process JSON-format signed transaction files
  - Use content-based format detection: attempt to parse the file as `XRPTransactionFile` JSON first; if parsing fails or the transactions slice is empty, fall back to the existing text-format `ReadFileSlice` path
  - Do not use file extension for detection — both the legacy text format and the new JSON format share the `.json` extension
  - For a successfully parsed JSON file: find entries where `is_complete == true` and extract the `signed_blob` for submission; return a clear error if `is_complete` is false
  - If both parsers fail, return an error identifying the unsupported file format
  - The existing submission path (`SubmitTransaction`) is unchanged
  - _Requirements: 4_

- [ ] 6. Fix the DI container so all XRP multisig use cases wire correctly without panicking at startup

- [ ] 6.1 Wire `SetSignerListUseCase` in its DI factory function
  - Replace the current "gRPC removed" startup panic with proper construction of the use case
  - Inject: `PublicXRP` (satisfying `SignerListPreparer`), the UUID handler, `XRPSignerListRepository` (keygen DB, from Task 2.1), `XRPSignerEntryRepository` (keygen DB, from Task 2.2), and the transaction file repository
  - Requires Tasks 1, 2.1, and 2.2 to be complete
  - _Requirements: 7_

- [ ] 6.2 Replace the `AddMultisigSignature` DI panic with a no-op use case struct
  - Implement a `notImplementedAddMultisigSignatureUseCase` struct whose `Execute()` returns `fmt.Errorf("AddMultisigSignature is not yet implemented")`
  - The struct must fully satisfy the `AddMultisigSignatureUseCase` interface; a nil interface return is forbidden (causes a silent nil pointer panic at the first call site, which is worse than the current named startup panic)
  - The use case application code remains intact for future extension; only the DI factory is changed
  - _Requirements: 7_

- [ ] 6.3 Verify that `CreateMultisigTxUseCase` and `SubmitMultisigTxUseCase` DI factories compile and run after injecting the new watch DB repositories
  - These factories are already wired in the container; update them to inject `XRPPendingMultisigRepository` and `XRPMultisigSignatureRepository` from Tasks 3.1 and 3.2
  - Confirm the full DI chain builds without errors: run `make check-build` after all repository adapters are in place
  - _Requirements: 7_

- [ ] 7. Add the `set-signer-list` Watch wallet CLI command and update the wallet adapter

- [ ] 7.1 (P) Implement the `watch send multisig set-signer-list` Cobra command
  - Accept `--account` (required), `--quorum` (required, must be ≥ 1), and `--signers` (required, comma-separated `address:weight,...` string) flags
  - Parse the `--signers` string into structured signer entries via `parseSignerEntries()`; validate quorum ≥ 1 and signer count between 1 and 8 in the CLI layer before calling the use case
  - Delegate to `SetSignerListUseCase.Execute()`; on success, print the unsigned transaction file path to stdout
  - Follow the existing Cobra command registration pattern used in the XRP Watch CLI command group
  - _Requirements: 1, 8_

- [ ] 7.2 (P) Update the XRPWatch wallet adapter and the create commands to pass multisig quorum
  - Add `setSignerListUseCase` as a field in the `XRPWatch` adapter struct; update the constructor and add a `SetSignerList()` method delegating to the use case
  - Register the new `set-signer-list` command in the XRP Watch CLI subcommand tree
  - Add an optional `--quorum` flag to the `create deposit`, `create payment`, and `create transfer` commands; when supplied (value ≥ 2), set `MultisigQuorum` in `CreateTransactionInput` to activate the JSON file path from Task 4
  - _Requirements: 2, 8_

- [ ] 8. Add the P2 end-to-end test for the complete 2-of-2 multisig payment flow

- [ ] 8.1 Implement the P2 E2E test script covering the full multisig lifecycle against a local rippled node
  - **Setup**: fund two signer accounts from the genesis wallet; call `watch send multisig set-signer-list` to configure a 2-of-2 signer list on the sender account; sign and broadcast the `SignerListSet` TX via the existing watch send command
  - **Create**: Watch wallet creates an unsigned multisig payment JSON file with `required_signatures=2`, `signature_count=0`, `is_complete=false`
  - **Sign 1**: Keygen wallet signs the unsigned file; assert that the output file has `signature_count=1` and `is_complete=false`
  - **Sign 2**: Sign wallet receives the partially-signed file from Keygen and signs it; assert `signature_count=2` and `is_complete=true`
  - **Send & Verify**: Watch wallet reads the complete file and broadcasts the fully-signed blob; call `ledger_accept` on the rippled standalone node; assert the transaction is validated and balances updated correctly
  - Runs against the existing Docker Compose rippled standalone setup (same infrastructure as P1)
  - _Requirements: 9_

- [ ] 8.2 Add Makefile targets for the P2 multisig E2E test
  - Add `xrp-e2e-p2` (interactive run), `xrp-e2e-p2-ci` (CI-friendly non-interactive run), and `xrp-e2e-p2-reset` (ledger reset / cleanup) targets
  - Follow the structure and variable conventions of the existing `xrp-e2e-p1` / `xrp-e2e-p1-ci` targets
  - _Requirements: 9_
