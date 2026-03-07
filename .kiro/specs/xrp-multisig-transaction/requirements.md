# Requirements Document

## Project Description (Input)

Implement XRPL Multi-signature Transaction support using a file-based offline signing workflow
that aligns with the existing BTC multisig flow described in `docs/transaction-flow.md`.
The implementation must not deviate significantly from the current BTC/ETH approach and must use
the XRPL native multisig protocol as of 2026.

---

## Background

### XRPL Native Multisig

XRP Ledger supports native multi-signing via **SignerList**. Key characteristics:

- **SignerList**: An on-ledger object attached to an account that lists authorized signers and their weights
- **SignerQuorum**: The minimum total weight required to authorize a transaction
- **Multisign**: Each signer independently signs the raw transaction JSON using `wallet.Multisign()`
- **Serial accumulation**: Signatures are accumulated by passing the previous signer's blob to the next signer; the `Signers` array grows with each step
- **Combination**: All individually-signed blobs can also be combined in parallel via `xrpl_submit` multi-sig combine
- **Zero master key**: Optionally, the account master key can be disabled after configuring SignerList (not required for this feature)

### Signing Approaches

| Approach | Description | Alignment with BTC flow |
|----------|-------------|------------------------|
| **Serial file-passing** | Each signer signs the previous blob; file is passed between wallets | ✅ Matches BTC PSBT flow |
| **Parallel + combine** | Each signer independently signs the raw JSON; Watch combines blobs | ❌ Requires Watch to aggregate |

**This specification uses the serial file-passing approach** for alignment with the BTC workflow.

### Current Implementation State

The following components are **already implemented** and must not be re-implemented:

| Component | Location | Status |
|-----------|----------|--------|
| `SignTransactionNative` with `existingSignedBlob` | `infrastructure/api/xrp/public/sign.go` | ✅ Done |
| `signTransactionUseCase` with serial accumulation | `usecase/sign/xrp/sign_transaction.go` | ✅ Done |
| `XRPTransactionFile` with `RequiredSignatures`, `SignatureCount`, `IsComplete` | `dto/xrp/transaction_file.go` | ✅ Done |
| `sendTransactionUseCase` (submits blob when `is_complete`) | `usecase/watch/xrp/send_transaction.go` | ✅ Done |
| Domain entities: `XRPSignerList`, `XRPSignerEntry`, `XRPPendingMultisig` | `domain/chains/xrp/` | ✅ Done |
| Repository port interfaces for signer list, pending multisig | `ports/repository/cold/` and `watch/` | ✅ Done |
| DB schema + SQLC generated code (all 5 tables) | `tools/atlas/`, `infrastructure/database/*/sqlcgen/` | ✅ Done |
| Use case interfaces in `watch/interfaces.go` | `usecase/watch/interfaces.go` | ✅ Done |
| Use case implementations (application logic) | `usecase/watch/xrp/set_signer_list.go`, etc. | ✅ Done |

The following components are **missing** and must be implemented:

| Component | Gap | Required Action |
|-----------|-----|-----------------|
| `PrepareSignerListSetTransaction` | Port interface defined, not implemented | Implement in infrastructure |
| `PrepareSetRegularKeyTransaction` | Port interface defined, not implemented | Implement in infrastructure (optional) |
| `CombineTransaction` | Port interface defined, not implemented | Not needed for serial approach |
| Repository implementations for XRP multisig tables | No concrete repos | Implement for keygen + watch DB |
| DI wiring for `SetSignerList` | Currently panics with "gRPC removed" message | Fix DI wiring |
| DI wiring for `AddMultisigSignature` | Currently panics | Fix or remove (not needed for serial) |
| CLI commands for multisig setup (`set-signer-list`) | Missing entirely | Add to Watch wallet CLI |
| `XRPWatch` wallet adapter for multisig use cases | Not included | Update adapter |
| Multisig-mode transaction creation (Watch) | `create_transaction.go` always uses `required_signatures: 1` | Support `required_signatures: N` |
| Keygen-side signer account provisioning | No CLI/flow for setting up signer keys | Add keygen flow |
| E2E test for multisig flow (P2) | Only P1 (single-sig) exists | Add P2 E2E test |

---

## Requirements

### FR-1: SignerList Configuration (Watch Wallet — One-Time Setup)

**Goal**: The Watch wallet must be able to create and broadcast a `SignerListSet` transaction that
configures an XRP account for multi-signature authorization.

**Acceptance Criteria**:

1. The Watch wallet CLI exposes a `set-signer-list` command that accepts:
   - `--account` — the XRP account address to configure (required)
   - `--quorum` — minimum signature weight required (required, must be ≥ 1)
   - `--signer` flag (repeatable) — signer entries in `address:weight` format (required, 1–8 entries)
2. The command produces an unsigned `SignerListSet` transaction file compatible with the existing file format
3. The file can be signed by the keygen wallet using the existing `sign` command (single-sig, master key)
4. After the signed file is sent via the Watch wallet `send` command, the signer list is active on-ledger
5. The Watch wallet stores the configured signer list and entries in its local database (`xrp_signer_list`, `xrp_signer_entry` tables in the keygen DB)
6. Validation: `SignerQuorum` must not exceed the sum of all signer weights
7. Validation: Maximum 8 signer entries (XRPL protocol limit)
8. The `SetSignerListUseCase` DI wiring must be fully functional (currently panics; must be fixed)

**Implementation Notes**:
- Uses `SignerListPreparer.PrepareSignerListSetTransaction()` port interface
- Infrastructure must implement `PrepareSignerListSetTransaction` using the WebSocket `xrpl_submit` RPC
- The signed TX is a standard single-sig transaction (signed by the account's master key, not by signers)

---

### FR-2: Multisig Transaction Creation (Watch Wallet)

**Goal**: The Watch wallet must be able to create unsigned payment transactions that require multiple
signatures, using the same file-based mechanism as BTC.

**Acceptance Criteria**:

1. The Watch wallet supports a multisig mode for `create deposit`, `create payment`, and `create transfer` commands
2. When multisig mode is active, the unsigned transaction file has `required_signatures` set to the `SignerQuorum` count (N > 1) rather than 1
3. The Watch wallet reads the configured `SignerQuorum` from the local signer list database
4. The unsigned transaction JSON in the file is the raw `Payment` transaction JSON (without `Signers` array), identical to the single-sig case
5. The `XRPTransactionFile` format's `required_signatures` field is already defined; no changes to the DTO are needed
6. `signature_count` starts at 0; `is_complete` starts as `false`
7. The account address used as the sender must have a configured signer list in the database

**Implementation Notes**:
- The `create_transaction.go` use case must detect if the sender account has a signer list configured and set `required_signatures` accordingly
- Alternatively, a new flag or dedicated command (`create multisig-payment`) may be used — design decision for Phase 2

---

### FR-3: Offline Multisig Signing (Keygen / Sign Wallet)

**Goal**: Each authorized signer signs the transaction using the existing offline signing mechanism.
The file is passed between wallets until enough signatures are collected.

**Acceptance Criteria**:

1. The existing `sign --file <path>` command handles multisig transactions without modification
2. When `required_signatures > 1`, the signer calls `SignTransactionNative` with `isMultiSig: true`
3. The first signer (Keygen wallet) receives the unsigned file and produces a file with `signature_count: 1`
4. Each subsequent signer (Sign wallet) receives the partially-signed file and appends their signature via serial blob accumulation
5. When `signature_count >= required_signatures`, the file's `is_complete` field is set to `true`
6. The file naming follows the existing pattern: `{action}_{type}_{txID}_{signedCount}.json`
7. The returned `IsComplete` flag guides the operator to either pass the file to the next signer or send it via Watch

**Note**: This requirement is **already implemented** via `sign/xrp/sign_transaction.go`. It is
listed here as a confirmed acceptance criterion — no new implementation is required for FR-3.

---

### FR-4: Multisig Transaction Submission (Watch Wallet)

**Goal**: The Watch wallet submits the fully-signed transaction blob when all required signatures
have been collected.

**Acceptance Criteria**:

1. The existing `send --file <path>` command handles fully-signed multisig transactions without modification
2. When `is_complete: true`, the transaction blob (with accumulated `Signers` array) is submitted to the XRPL node
3. The Watch wallet waits for ledger validation using the existing `waitValidation` polling loop
4. Transaction confirmation is verified via `GetTransaction`

**Note**: This requirement is **already implemented** via `watch/xrp/send_transaction.go`. It is
listed here as a confirmed acceptance criterion — no new implementation is required for FR-4.

---

### FR-5: Signer Key Provisioning (Keygen Wallet)

**Goal**: The Keygen wallet must be able to generate and manage keys for individual signer accounts,
and export their addresses for Watch wallet configuration.

**Acceptance Criteria**:

1. The existing `create seed` and `create hdkey` commands are reused for signer account setup
2. The keygen wallet can generate keys for additional account types (e.g., `authorization` accounts used as signers)
3. Signer addresses can be exported via the existing `export address` command
4. The Watch wallet imports signer addresses to build the `set-signer-list` command inputs
5. The signer's private key remains offline in the Keygen or Sign wallet's local DB at all times

**Note**: Key generation is **already implemented** via existing keygen commands. This requirement
confirms that no new keygen infrastructure is needed for FR-5 — only the workflow documentation
and E2E test coverage for the multisig setup path.

---

### FR-6: Repository Implementations for Multisig Tables

**Goal**: Concrete repository implementations must exist for the XRP multisig database tables
so that the use cases can persist signer lists and operate end-to-end.

**Acceptance Criteria**:

1. Keygen DB repository implementations:
   - `XRPSignerListRepositorier` — implemented for all configured DB drivers (postgres, mysql, sqlite)
   - `XRPSignerEntryRepositorier` — implemented for all configured DB drivers
   - `XRPRegularKeyRepositorier` — implemented for all configured DB drivers (used by `SetRegularKey`)
2. Watch DB repository implementations:
   - `XRPPendingMultisigRepositorier` — implemented (used by `CreateMultisigTxUseCase`, `SubmitMultisigTxUseCase`)
   - `XRPMultisigSignatureRepositorier` — implemented (used by `AddMultisigSignatureUseCase`)
3. All implementations satisfy the port interfaces defined in `internal/application/ports/repository/`
4. SQLC generated code (already exists) is used as the data access layer

**Notes on scope**:
- `XRPPendingMultisigRepositorier` and `XRPMultisigSignatureRepositorier` are used by the DB-centric
  use cases (`create_multisig_tx`, `add_multisig_signature`). These use cases are NOT the primary flow
  for this spec (serial file-based signing is primary). However, the repository implementations are still
  required because these use cases are already defined and the DI container already partially wires them.
- For the **file-based serial flow**, only the `XRPSignerListRepositorier` + `XRPSignerEntryRepositorier`
  are strictly required (to read the quorum when creating a multisig TX).

---

### FR-7: DI Container Fixes

**Goal**: Fix the DI container so that multisig-related use cases are properly wired instead of panicking.

**Acceptance Criteria**:

1. `newXRPWatchSetSignerListUseCase()` is fully implemented and wired (no panic)
2. `newXRPWatchSetRegularKeyUseCase()` is fully implemented and wired (optional; may remain stubbed if FR-8 is out of scope)
3. `newXRPWatchAddMultisigSignatureUseCase()` is wired if repository implementations exist; otherwise removed or clearly stubbed with a non-panic error
4. `newXRPWatchCreateMultisigTxUseCase()` and `newXRPWatchSubmitMultisigTxUseCase()` are already wired; verify they compile and run

---

### FR-8: Watch Wallet CLI Commands (Multisig Operations)

**Goal**: Expose CLI commands for the multisig lifecycle on the Watch wallet.

**Acceptance Criteria**:

1. `watch api xrp set-signer-list` command is added (from FR-1)
2. All commands follow the existing Cobra CLI pattern used in `internal/interface-adapters/cli/watch/`
3. `XRPWatch` wallet adapter is updated to include the new use cases

---

### FR-9: E2E Test — Multisig Payment (P2)

**Goal**: A complete end-to-end test validates the XRP multisig flow against a local `rippled`
node in standalone mode.

**Acceptance Criteria**:

1. A new E2E test pattern `P2` covers the 2-of-2 multisig payment flow:
   - Setup: Two signer accounts created; `SignerListSet` configured on sender account
   - Create: Unsigned multisig payment transaction file created by Watch wallet
   - Sign 1: Keygen wallet signs (signer 1), producing `signed_1` file
   - Sign 2: Sign wallet signs (signer 2), producing `signed_2` file with `is_complete: true`
   - Send: Watch wallet broadcasts the fully-signed transaction
   - Verify: Transaction confirmed on ledger; balances updated correctly
2. `make xrp-e2e-p2` and `make xrp-e2e-p2-ci` targets are added to the Makefile
3. The test runs against a `rippled` standalone node (Docker Compose)
4. The test passes reliably in CI

---

## Out of Scope

The following are explicitly **out of scope** for this specification:

| Item | Rationale |
|------|-----------|
| Parallel signature combination (`CombineTransaction`) | Not needed for serial file-based flow |
| `SetRegularKey` full implementation | Optional; XRP multisig works without it |
| Disabling the master key (`AccountSet: asfDisableMaster`) | Advanced use case; not required for basic multisig |
| Web UI or API for multisig coordination | No HTTP/REST layer in this project |
| Destination tag support | Separate feature gap; does not affect multisig |
| EscrowCreate, NFToken, payment channels | Separate transaction types; out of scope |
| Multisig for `SignerListSet` transaction itself (meta-multisig) | Only single-sig setup of the SignerList is required |

---

## Constraints

1. **Architecture**: Must follow Clean Architecture. Domain layer has zero infrastructure dependencies.
2. **Offline**: Keygen and Sign wallets must never make network calls during signing (FR-3, FR-5).
3. **File transfer**: Transaction files move between wallets via physical transfer (USB). No networked key sharing.
4. **Library**: Use `github.com/Peersyst/xrpl-go` for offline signing (already in use).
5. **Protocol**: Target XRPL mainnet protocol. The serial accumulation approach must produce a valid multisig transaction blob.
6. **DB drivers**: Repository implementations must support postgres, mysql, and sqlite (matching existing codebase pattern).
7. **No gRPC**: The former gRPC-based `CombineTransaction` must not be reintroduced. WebSocket-native implementation only.
