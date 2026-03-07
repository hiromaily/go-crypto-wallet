# Research & Design Decisions

---

## Summary

- **Feature**: `xrp-multisig-transaction`
- **Discovery Scope**: Extension (existing XRP infrastructure, extending multisig support)
- **Key Findings**:
  - The serial file-passing signing path is already fully implemented; the gap is entirely in setup (SignerListSet) and wire-up (DI, CLI, repositories)
  - XRPL `SignerListSet` follows the same transaction structure as `Payment`, enabling direct reuse of `PrepareTransaction` infrastructure patterns
  - DB schema and SQLC-generated code for all 5 multisig tables already exist; only concrete repository adapter implementations are missing

---

## Research Log

### XRPL SignerListSet Transaction Format

- **Context**: Need to implement `PrepareSignerListSetTransaction` in infrastructure layer
- **Sources Consulted**: Existing `PrepareTransaction` in `infrastructure/api/xrp/public/transaction.go`; `dto/xrp/transaction.go` which already defines `SignerListSetTxInput`; XRPL docs reference in code comments
- **Findings**:
  - `SignerListSetTxInput` DTO already exists with all required fields (`Account`, `SignerQuorum`, `SignerEntries`, `Fee`, `Sequence`, `LastLedgerSequence`)
  - The infrastructure wire type for `SignerListSet` must include JSON-tagged `SignerEntry` objects as nested arrays
  - XRPL protocol: `SignerEntries` is an array of `{"SignerEntry": {"Account": "r...", "SignerWeight": N}}`
  - Fee calculation follows same pattern as `Payment`: use minimum fee "12" drops unless overridden
  - Sequence and LastLedgerSequence come from `AccountInfo` RPC, same as `Payment`
- **Implications**: `PrepareSignerListSetTransaction` can reuse `AccountInfo` lookup and `toDTOXxx` pattern verbatim; only the JSON wire type differs

### Existing `create_transaction.go` Multisig Mode Gap

- **Context**: FR-2 requires multisig TX creation with `required_signatures > 1`
- **Sources Consulted**: `usecase/watch/xrp/create_transaction.go`; `dto/xrp/transaction_file.go`
- **Findings**:
  - `generateHexFile()` writes files using `WriteFileSlice()` with serialized text format (not the JSON format)
  - The JSON file format (`XRPTransactionFile`) is only used by `ReadXRPJSONFile` / `WriteXRPJSONFile` in sign and send paths
  - `create_transaction.go` uses `WriteFileSlice()` with format `"senderAccount\nuuid,txJSON"` — this is a plain-text file, NOT the JSON file
  - The sign use case (`sign_transaction.go`) reads `ReadXRPJSONFile` (JSON format with `RequiredSignatures` field)
  - **Critical gap**: The current Watch TX creation writes plain text; the Sign wallet reads JSON format — these two formats are INCOMPATIBLE for the existing single-sig flow already working (P1)
  - Investigating further: `send_transaction.go` reads `ReadFileSlice()` (text format), but the sign wallet writes JSON format... this implies two different file paths for single-sig vs multisig
- **Re-investigation**: P1 single-sig flow uses the text file format throughout (Watch writes text → Keygen reads text and signs → Watch reads text to submit). The JSON format (`ReadXRPJSONFile`) is a separate newer implementation prepared for multisig but not yet connected to the Watch TX creation.
- **Implications**: FR-2 requires `create_transaction.go` to write JSON-format files (with `required_signatures`) when multisig is needed. This is a new code path. Single-sig can remain text-based or migrate to JSON — design decision needed.

### Repository Implementation Pattern

- **Context**: FR-6 requires concrete repository implementations
- **Sources Consulted**: SQLC generated files in `infrastructure/database/postgres/sqlcgen/xrp_signer_list.sql.go`; existing repository implementations in `infrastructure/repository/`
- **Findings**:
  - SQLC generates typed query functions (`GetXRPSignerListByAccountID`, `InsertXRPSignerList`, etc.) for all 5 tables
  - The repository adapter pattern is: struct wraps `*sqlcgen.Queries`, methods convert SQLC types → domain entities
  - `XRPSignerList` and `XRPSignerEntry` belong to the keygen DB (cold wallet)
  - `XRPPendingMultisig` and `XRPMultisigSignature` belong to the watch DB
  - `XRPRegularKey` belongs to the keygen DB
- **Implications**: Four repository adapter structs needed (one per interface), split across cold/ and watch/ directories. Must implement all methods in the port interface.

### DI Panic Locations

- **Context**: FR-7 requires fixing DI panics
- **Sources Consulted**: `internal/di/container.go` lines 1559–1584
- **Findings**:
  - `newXRPWatchSetRegularKeyUseCase()` → panic ("gRPC removed")
  - `newXRPWatchSetSignerListUseCase()` → panic ("gRPC removed")
  - `newXRPWatchAddMultisigSignatureUseCase()` → panic ("gRPC removed")
  - `newXRPWatchCreateMultisigTxUseCase()` → already wired (calls `watchusecasexrp.NewCreateMultisigTxUseCase`)
  - `newXRPWatchSubmitMultisigTxUseCase()` → already wired
  - The `AddMultisigSignature` use case requires `XRPPendingMultisigRepositorier` + `XRPMultisigSignatureRepositorier` + `TransactionCombiner`
  - `TransactionCombiner` (`CombineTransaction`) is NOT implemented anywhere
- **Implications**: For serial file-based flow, `AddMultisigSignature` DI can remain stubbed (it's the DB-centric approach not needed for file-based flow). `SetSignerList` DI must be fully wired.

### File Format Alignment

- **Context**: Understanding which file format to use for the new multisig TX creation path
- **Sources Consulted**: `dto/xrp/transaction_file.go`; `sign_transaction.go`; `send_transaction.go`
- **Findings**:
  - JSON format (`XRPTransactionFile`) was designed for multisig from the start (has `RequiredSignatures`, `SignatureCount`, `IsComplete`, `SignedBlob`)
  - `sign_transaction.go` ALREADY handles JSON format with full multisig logic
  - The EXISTING single-sig P1 test uses text format through a DIFFERENT code path (`CreateRawTransaction → WriteFileSlice → ReadFileSlice → SubmitTransaction`)
  - For multisig, we need a PARALLEL path using JSON format: `CreateRawTransaction → WriteXRPJSONFile → ReadXRPJSONFile (sign) → ReadXRPJSONFile or ReadFileSlice (send)`
  - `send_transaction.go` reads text format (ReadFileSlice), but multisig blob is in JSON format `signed_blob` field — the send path must be extended to also support JSON format OR the JSON file must produce a compatible text file after signing
- **Decision**: The `create_transaction.go` should support a multisig mode that outputs JSON format. The `send_transaction.go` must be extended to detect and read JSON format files. See Design Decision below.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| Serial file-passing | Each signer extends the blob; file transferred between wallets | Matches BTC PSBT flow; fully offline; already implemented in sign.go | Requires strict ordering of signers | **Selected** |
| Parallel + combine | Each signer independently signs raw JSON; Watch combines blobs | More flexible signer ordering | Requires `CombineTransaction` implementation; Watch must be online during signing coordination | Out of scope |
| DB-centric aggregation | Watch holds pending multisig state; signers submit blobs back to Watch | Central coordination; retry possible | Requires network access for signer coordination; complex state machine | Existing stubs; not primary for this spec |

---

## Design Decisions

### Decision: JSON File Format for Multisig TX Creation

- **Context**: `create_transaction.go` currently writes text-format files; `sign_transaction.go` reads JSON-format files; these are incompatible for multisig
- **Alternatives Considered**:
  1. Migrate all TX creation to JSON format (single-sig + multisig)
  2. Add a parallel multisig code path that writes JSON format while keeping text path for single-sig
  3. Post-process the text file into JSON after creation
- **Selected Approach**: Option 2 — add a multisig-specific creation path in `create_transaction.go` that calls `WriteXRPJSONFile`. The single-sig text path remains unchanged to avoid regression. The `send_transaction.go` is extended to detect JSON files (by file extension or content) and extract the signed blob accordingly.
- **Rationale**: Zero regression risk on P1 (single-sig); clear separation of code paths; minimal change surface
- **Trade-offs**: Two code paths for TX creation; some duplication in file format handling; `send_transaction.go` must handle both formats
- **Follow-up**: Future cleanup could unify to JSON format for all XRP TX files

### Decision: SignerList Stored in Keygen DB

- **Context**: `XRPSignerList` and `XRPSignerEntry` domain entities are defined; the use case (`set_signer_list.go`) injects `repocold.XRPSignerListRepositorier` (cold wallet DB)
- **Alternatives Considered**:
  1. Store signer list in Watch DB (watch can read it when creating multisig TX)
  2. Store in Keygen DB (cold) only
  3. Duplicate across both
- **Selected Approach**: Per existing use case code, the signer list is stored in the keygen DB (cold). The Watch wallet reads signer list data from its own copy after import, OR the `create_transaction.go` reads from a Watch-side copy.
- **Rationale**: The existing code uses `repocold.XRPSignerListRepositorier` in `set_signer_list.go`. The Watch TX creation use case also needs access to the signer quorum. Watch-side requires Watch to have access to the same repository. Since both Watch and Keygen wallets run independently, the Watch wallet needs its own local record of the signer quorum. The `set_signer_list` use case (Watch wallet command) can store the configured signer list in the Watch DB as well.
- **Trade-offs**: The Watch wallet must import/store the signer list configuration locally. This is acceptable as the operator configures it via CLI before creating multisig transactions.
- **Follow-up**: Determine whether to add a `XRPSignerListRepositorier` to the Watch DB schema, or pass quorum as CLI flag on `create` command.

### Decision: `AddMultisigSignature` DI Stub

- **Context**: `AddMultisigSignature` use case requires `CombineTransaction` which is not implemented
- **Alternatives Considered**:
  1. Implement `CombineTransaction` (WebSocket-based)
  2. Leave DI stub but remove panic (return error instead)
  3. Delete the use case
- **Selected Approach**: Change panic to a proper error return in the DI factory. Keep the use case code. `CombineTransaction` implementation is deferred.
- **Rationale**: The primary flow (serial file-based) doesn't need `CombineTransaction`. Keeping the use case allows future extension without architecture changes.
- **Trade-offs**: `AddMultisigSignature` remains unusable at runtime (returns error) but compiles cleanly.

---

## Risks & Mitigations

- **File format split risk**: Two code paths for TX creation may diverge over time — mitigate by documenting clearly and adding integration tests for both paths
- **Watch DB signer quorum access**: The `create_transaction.go` needs the signer quorum from DB; if the signer list is only in the keygen DB, the Watch wallet cannot access it — mitigate by storing signer list in Watch DB too, or accepting quorum as a CLI flag
- **XRPL fee calculation for SignerListSet**: SignerListSet has a higher minimum fee when many signers are involved (N+2 drops per signer) — mitigate by using the `fee_cushion` or computing from signer count
- **Serial signing order**: If a signer loses the file, the entire multisig chain must restart — mitigate with E2E test and clear operator documentation

---

## References

- XRPL SignerListSet documentation: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
- XRPL Multi-signing guide: https://xrpl.org/docs/concepts/accounts/multi-signing
- Existing Payment transaction implementation: `internal/infrastructure/api/xrp/public/transaction.go`
- Existing multisig signing: `internal/application/usecase/sign/xrp/sign_transaction.go`
- SQLC generated signer list queries: `internal/infrastructure/database/postgres/sqlcgen/xrp_signer_list.sql.go`
