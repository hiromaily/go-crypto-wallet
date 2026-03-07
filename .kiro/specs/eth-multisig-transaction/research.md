# Research & Design Decisions

---

## Summary

- **Feature**: `eth-multisig-transaction`
- **Discovery Scope**: Complex Integration — new Safe contract infrastructure, new use cases across 3 wallets, new file format
- **Key Findings**:
  - Safe v1.4.1 contracts are well-documented; `execTransaction` and `getTransactionHash` are stable API surfaces used since v1.0. Go bindings generated via `abigen` from the official ABI are straightforward.
  - EIP-712 `safeTxHash` recomputation in pure `go-ethereum/crypto` is feasible without any Safe SDK. The domain separator and struct hash computation maps cleanly to `crypto.Keccak256`.
  - Safe signature validation requires signatures sorted by signer address (ascending) and v-byte adjusted to 27/28 (`go-ethereum` returns 0/1). This is a critical correctness constraint.
  - ETH Sign wallet (`ETHSign.SignTx`) is currently a no-op. It can be wired to `SignMultisigTransactionUseCase` without breaking any existing flow.
  - No DB tables are needed — the multisig file is the sole state carrier (parallels BTC PSBT flow).

---

## Research Log

### Safe v1.4.1 Contract API

- **Context**: Need to interact with Safe's `getTransactionHash` and `execTransaction` on-chain.
- **Findings**:
  - `getTransactionHash` takes all transaction parameters plus the Safe's current `nonce` and returns a `bytes32` hash conforming to EIP-712.
  - `execTransaction` validates that the aggregated `signatures` parameter contains at least `threshold` valid owner signatures sorted in ascending address order.
  - Safe contract supports three signature types: `v=0` (contract sig), `v=1` (approved hash), `v=27/28` (EOA eth_sign / eth_signTypedData). Our flow uses EOA off-chain signing, so v should be 27 or 28 after adjustment.
  - The Safe contract's domain separator uses `chainId` and `verifyingContract` (the Safe address). This matches EIP-712 standard.
- **Implications**: The port interface `GetSafeTxHash` can delegate to the on-chain call; offline recomputation in the signing use case must reproduce the identical hash.

### EIP-712 Hash Computation

- **Context**: Signers must recompute `safeTxHash` offline to verify the file before signing.
- **Findings**:
  - Domain separator: `keccak256(abi.encode(keccak256("EIP712Domain(uint256 chainId,address verifyingContract)"), chainId, safeAddress))`
  - Safe TX type hash (constant): `keccak256("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address payable refundReceiver,uint256 nonce)")`
  - Struct hash: `keccak256(abi.encode(SAFE_TX_TYPEHASH, to, value, keccak256(data), operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, nonce))`
  - Final hash: `keccak256("\x19\x01", domainSeparator, structHash)` — standard EIP-712 encoding
  - All `abi.encode` packing uses ABI encoding (left-padded to 32 bytes for all types).
  - `go-ethereum/accounts/abi` provides `abi.Arguments.Pack()` for this purpose.
- **Implications**: Signing use case needs pure Go implementation of the above hash. `crypto.Sign` in `go-ethereum/crypto` signs the raw hash; returned v=0/1 must be incremented to 27/28 for Safe compatibility.

### Existing Contract Binding Pattern

- **Context**: No existing Safe contract bindings in the codebase. One existing contract: ERC-20 token ABI at `internal/infrastructure/contract/token-abi.go`.
- **Findings**:
  - The existing token-abi.go is hand-written (struct + ABI string). It does NOT use `abigen`.
  - The codebase uses `go-ethereum/accounts/abi/bind` for binding calls.
  - For Safe, using `abigen` produces a full typed Go client automatically, eliminating manual ABI string management.
- **Implications**: Place generated Safe bindings in `internal/infrastructure/contract/safe/` with `DO NOT EDIT` header. Add a `make safe-abi` target to `Makefile`.

### File-Based State Carrier (No DB)

- **Context**: XRP multisig uses DB tables for signer list state; BTC uses PSBT files. ETH multisig needs to decide.
- **Findings**:
  - For ETH multisig, the Safe's `nonce` is the only on-chain state needed at proposal time (fetched via `getTransactionHash`).
  - All signature accumulation state lives in the JSON file itself — no cross-wallet DB sync needed.
  - UUID in the file name provides unique proposal identity without a DB-assigned ID.
- **Implications**: No DB schema changes, no SQLC changes. The `ETHMultisigTransactionFile` is the sole state artifact.

### Sign Wallet Activation

- **Context**: `ETHSign.SignTx()` currently logs "no functionality" and returns empty values.
- **Findings**:
  - `ETHSign` struct does not currently have a `signTxUseCase` field.
  - The Sign wallet CLI command `sign tx` exists and routes to `ETHSign.SignTx()`.
  - Adding `signMultisigTxUseCase keygenusecase.SignMultisigTransactionUseCase` to `ETHSign` and wiring it introduces no breaking change to the single-sig flow.
- **Implications**: `ETHSign` struct needs a new field + updated constructor. The Sign wallet DI container gets a new factory function `newETHSignMultisigSignTransactionUseCase()`.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| File-only (no DB) | All state in JSON file; UUID as identity | No schema changes, matches BTC PSBT pattern, simple | No persistent audit trail, can't list pending proposals via CLI | **Selected** — sufficient for 2-of-2 E2E scope |
| DB-backed | Store proposal in watch DB; file contains hash reference | Queryable history, re-submission possible | Requires schema changes, cross-wallet DB coupling | Ruled out — adds complexity for no user-visible benefit in scope |
| SDK-based (safe-core-sdk-go) | Use a Go Safe SDK if one exists | Less manual work | No mature Go SDK as of 2026; adds external dependency | Ruled out — pure go-ethereum is sufficient |

---

## Design Decisions

### Decision: `SafeExecParams` as port boundary type

- **Context**: FR-3 requires `SafeExecuter` port interface. Initial design passed `*ETHMultisigTransactionFile` directly — violates port boundary rules.
- **Alternatives Considered**:
  1. Pass full DTO to `SafeExecuter` — simple but couples port to file format
  2. Define `SafeExecParams` struct in ports package — decouples execution from file format
- **Selected Approach**: `SafeExecParams` struct defined in `internal/application/ports/api/eth/` alongside the interface. Use case converts from DTO to `SafeExecParams` before calling the port.
- **Rationale**: Consistent with project's DTO rules (port interfaces use only primitive types or types defined within the ports package).
- **Trade-offs**: Slight duplication of field definitions between DTO and params struct; justified by clean layer separation.

### Decision: Shared `SignMultisigTransactionUseCase` for Keygen and Sign wallets

- **Context**: Both Keygen (signer 1) and Sign (signer N) perform identical signing logic.
- **Alternatives Considered**:
  1. Duplicate implementations per wallet — clear but redundant
  2. Shared implementation in `usecase/keygen/eth/sign_multisig_transaction.go` — reused by both wallets
- **Selected Approach**: Single implementation in `usecase/keygen/eth/` package; both Keygen and Sign DI containers instantiate it via their respective factory functions.
- **Rationale**: Logic is identical; wallet type does not affect signing logic. DI wiring differences (key repo access) are handled at factory level.

### Decision: v-byte adjustment (+27) at signing time

- **Context**: `crypto.Sign` returns v=0 or v=1. Safe contract expects v=27 or v=28 for EOA signatures.
- **Selected Approach**: `SignMultisigTransactionUseCase` adjusts the last byte of the 65-byte signature: `sig[64] += 27` before encoding to hex.
- **Rationale**: This matches `eth_signTypedData` behavior. No `SigToEIP155` needed since we use a fixed chain ID from the file.

### Decision: Foundry script for Safe deployment (not Watch wallet CLI)

- **Context**: FR-7 originally included `watch safe deploy`. Safe deployment requires a deployer private key — the Watch wallet holds no private keys.
- **Selected Approach**: Safe deployment handled exclusively by `apps/eth-contracts/script/DeploySafe.s.sol` Foundry script, called from the E2E shell script.
- **Rationale**: Correct security boundary. Watch wallet = online, key-less. Deployment is an ops-time activity, not a runtime wallet operation.

---

## Risks & Mitigations

- **Hash mismatch between Watch (on-chain) and Signer (offline)** — Mitigated by: signer aborts with error if recomputed hash ≠ file's `safe_tx_hash`. On-chain computation via `getTransactionHash` is authoritative.
- **Wrong v-byte (Safe rejects signature)** — Mitigated by: unit tests that sign a known hash and verify using Safe's signature validation logic.
- **Safe nonce race** — Mitigated by: proposal file embeds the nonce at creation time; if another transaction is executed against the Safe before submission, the nonce in the file is stale and `execTransaction` will revert. Operator must recreate the proposal.
- **`abigen` Safe ABI version mismatch** — Mitigated by: pin Safe v1.4.1 ABI source in Makefile target; document ABI source URL.

---

## References

- Safe Smart Account — official contracts: `github.com/safe-global/safe-smart-account`
- EIP-712: Typed structured data hashing and signing — `eips.ethereum.org/EIPS/eip-712`
- go-ethereum crypto package: `pkg.go.dev/github.com/ethereum/go-ethereum/crypto`
- Safe signature encoding: `docs.safe.global/advanced/smart-account-signatures`
