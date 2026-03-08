# Requirements Document

## Project Description (Input)

Implement Ethereum Multi-signature Transaction support using MPC (Multi-Party Computation) /
Threshold Signature Scheme (TSS) as an **additional multisig pattern** alongside the Safe (Gnosis Safe)
smart-contract-based multisig already implemented.
The Safe implementation (E2E Pattern 3) is **preserved unchanged** — MPC-TSS adds a new E2E Pattern 4.
Unlike Safe, which aggregates signatures on-chain via `execTransaction`, TSS produces a single
standard ECDSA signature `(v, r, s)` that is indistinguishable from a regular EOA transaction to
the Ethereum network, and is broadcast via `eth_sendRawTransaction`.

---

## Background

### Ethereum Transaction Patterns (Overview)

This project supports multiple Ethereum transaction patterns. This spec adds Pattern 4 only;
Patterns 1–3 are complete and must not be modified.

| Pattern | Name | Signing Mechanism | Broadcast Method | Status |
|---------|------|-------------------|-----------------|--------|
| **P1** | Single-sig EOA | Single private key (Keygen wallet) | `eth_sendRawTransaction` | ✅ Done |
| **P2** | _(reserved)_ | — | — | — |
| **P3** | Safe 2-of-2 Multisig | EIP-712 per-owner signatures + Safe `execTransaction` | On-chain via Safe contract | ✅ Done |
| **P4** | MPC-TSS Multisig | Threshold ECDSA (T-of-N, no full key) | `eth_sendRawTransaction` | 🚧 This spec |

### Ethereum Multisig via MPC/TSS

The existing Safe implementation requires an on-chain smart contract and a gas-heavy `execTransaction`
call. MPC/TSS provides an alternative where:

- **No smart contract is required**: the distributed key produces a standard EOA address.
- **No on-chain aggregation**: signing happens off-chain; the result is a single 65-byte ECDSA signature.
- **Lower gas cost**: identical to a standard ETH transfer from an EOA.
- **Stronger key isolation**: the full private key never exists in one place at any point in time.

### TSS Workflow (T-of-N)

| Phase | Description |
|-------|-------------|
| **DKG** | N nodes run Distributed Key Generation; each receives a private key shard. The corresponding Ethereum address is derived from the joint public key. No node ever holds the full private key. |
| **Transaction Preparation** | A coordinator node (Watch wallet) builds a standard `types.Transaction`. |
| **Signing Round** | T-of-N nodes use their shards to collectively sign `tx.Hash()` via interactive MPC protocol. |
| **Signature Reconstruction** | The protocol outputs a standard 65-byte ECDSA signature `(r, s, v)`. |
| **Broadcast** | The signed transaction is submitted via `eth_sendRawTransaction`. |

### Relationship to Existing Safe Implementation

| Feature | Safe (Smart Contract) | MPC (TSS) |
|---|---|---|
| Transaction Type | Contract Call (`execTransaction`) | Legacy / DynamicFee Tx |
| Signature Type | EIP-712 (aggregated, per-owner) | ECDSA `(r, s, v)` — single signature |
| On-chain Identity | Safe Proxy address | Distributed EOA address |
| Gas Cost | High (Safe contract logic) | Same as standard EOA |
| Key Storage | Individual EOA keys, per owner | Distributed key shards, no single-point key |
| Library | `go-ethereum` only | TSS library + `go-ethereum` |

### Current Implementation State

The following components are **already implemented** and must not be re-implemented:

| Component | Location | Status |
|-----------|----------|--------|
| ETH single-sig EOA create / sign / send | `usecase/watch/eth/`, `usecase/keygen/eth/` | ✅ Done |
| ETH Safe 2-of-2 multisig (P3) | `usecase/watch/eth/`, `usecase/keygen/eth/`, `usecase/sign/eth/` | ✅ Done |
| `ETHTransactionFile` DTO | `application/dto/eth/transaction_file.go` | ✅ Done |
| `ETHMultisigTransactionFile` DTO | `application/dto/eth/multisig_transaction_file.go` | ✅ Done |
| `TxCreator`, `TxSender`, `TxSigner` small ports | `application/ports/api/eth/interfaces_small.go` | ✅ Done |
| Safe ports (`SafeNonceReader`, `SafeTxHashComputer`, `SafeExecutor`, `SafeInfoReader`) | `application/ports/api/eth/interfaces_safe.go` | ✅ Done |
| ETH key derivation (BIP-44, `deriveChildPrivKey`) | `usecase/keygen/eth/sign_transaction.go` | ✅ Done |
| ETH wallet adapters (keygen, sign, watch) | `interface-adapters/wallet/eth/` | ✅ Done |
| `go-ethereum` library | `go.mod` | ✅ In use |
| E2E P1 (single-sig) and P3 (Safe multisig) | `scripts/operation/eth/e2e/` | ✅ Done |

The following are **missing** and must be implemented:

| Component | Gap |
|-----------|-----|
| TSS library dependency | No MPC/TSS library in `go.mod` |
| DKG use case (Keygen wallet) | No distributed key generation flow |
| TSS key shard storage | No storage for per-node key shards |
| `MPCSigner` port interface | No TSS-based signer port |
| `MPCSigningSessionManager` port | No session/round management port |
| TSS infrastructure implementation | No TSS signing client in `infrastructure/` |
| `CreateMPCTransaction` use case (Watch) | Transaction proposal for TSS flow |
| `SignMPCTransaction` use case (MPC nodes) | TSS signing round participation |
| `SendMPCTransaction` use case (Watch) | Broadcast signed TSS transaction |
| MPC node communication layer | No inter-node transport (libp2p or gRPC) |
| CLI commands for MPC flow | No CLI for DKG and TSS signing |
| E2E P4 test script (MPC-TSS multisig) | Only P1/P3 exist |

---

## Requirements

### Requirement 1: TSS Library Integration

**Objective:** As a wallet developer, I want a Go TSS library integrated into the project,
so that MPC-based key generation and threshold signing can be implemented without custom cryptography.

#### Acceptance Criteria

1. The ETH MPC Signing System shall integrate a production-grade Go TSS/ECDSA library
   (e.g., `github.com/bnb-chain/tss-lib/v2` or an equivalent 2026-standard library) as a
   declared `go.mod` dependency.
2. The ETH MPC Signing System shall expose TSS capabilities exclusively through port interfaces
   in `internal/application/ports/` — no TSS library types shall appear in the use case or
   domain layers.
3. The ETH MPC Signing System shall use the `secp256k1` curve, matching Ethereum's cryptographic
   primitive.
4. If the chosen TSS library is updated or replaced, the ETH MPC Signing System shall require
   only changes in the infrastructure layer, leaving use case and domain layers unmodified.
5. The ETH MPC Signing System shall keep all TSS library imports confined to
   `internal/infrastructure/` — no direct TSS library calls in `application/` or `domain/`.

---

### Requirement 2: Distributed Key Generation (DKG)

**Objective:** As a wallet operator, I want to run a DKG ceremony across N nodes so that
a distributed Ethereum address is established without any single node possessing the full private key.

#### Acceptance Criteria

1. When a DKG ceremony is initiated, the ETH MPC Signing System shall coordinate N participant
   nodes to execute a T-of-N DKG protocol, where T (threshold) and N (total participants) are
   configurable inputs.
2. When DKG completes successfully, the ETH MPC Signing System shall derive a single Ethereum
   address from the joint public key, using the standard derivation (`Keccak256(pubkey)[12:]`).
3. When DKG completes, each participating node's ETH MPC Signing System shall store only its own
   private key shard locally — the full private key shall never be reconstructed or stored anywhere.
4. When DKG completes, the ETH MPC Signing System shall export the joint Ethereum address and
   joint public key to a file consumable by the Watch wallet's address import flow.
5. If any participant node fails during DKG, the ETH MPC Signing System shall abort the ceremony
   and report an error — partial shards from a failed ceremony shall not be persisted.
6. The ETH MPC Signing System shall support at minimum a 2-of-3 DKG configuration for E2E testing.
7. Where pre-computation parameters are required by the TSS library, the ETH MPC Signing System
   shall generate and persist them separately from the DKG ceremony itself, so that the DKG phase
   can proceed without recomputing pre-parameters each time.

---

### Requirement 3: TSS Key Shard Storage

**Objective:** As a wallet operator, I want each MPC node's key shard to be stored securely,
so that it can be loaded during a TSS signing session without exposing raw key material.

#### Acceptance Criteria

1. The ETH MPC Signing System shall store each node's TSS key shard in an encrypted file on the
   local filesystem.
2. The ETH MPC Signing System shall encrypt key shard files using a passphrase-derived key
   (minimum equivalent to AES-256), never storing the passphrase in plaintext.
3. When loading a key shard for signing, the ETH MPC Signing System shall decrypt it in memory
   and zero the plaintext after use.
4. If a key shard file is missing or corrupted, the ETH MPC Signing System shall return an error
   and refuse to proceed with signing.
5. The ETH MPC Signing System shall never log, print, or expose raw key shard bytes in any
   output, error message, or diagnostic trace.
6. Where pre-computation parameters accompany the key shard (e.g., Paillier key pairs), the ETH
   MPC Signing System shall store them in the same encrypted file as the shard.

---

### Requirement 4: MPC Signer Port Interface

**Objective:** As an architect, I want a dedicated port interface for TSS-based signing
so that the use case layer remains decoupled from the TSS library and inter-node networking.

#### Acceptance Criteria

1. The ETH MPC Signing System shall define an `MPCSigner` port interface in
   `internal/application/ports/api/eth/` with a method signature compatible with
   `Sign(ctx context.Context, hash []byte) ([]byte, error)`, returning a 65-byte ECDSA signature.
2. The ETH MPC Signing System shall define an `MPCSessionManager` port interface in
   `internal/application/ports/api/eth/` with methods to:
   - Initiate a new signing session for a given hash and set of participant node IDs
   - Wait for session completion and retrieve the resulting signature
3. The `MPCSigner` port interface shall accept only primitive types and standard library types —
   no TSS library-specific types shall appear in the port signature.
4. The ETH MPC Signing System shall ensure that the infrastructure implementation of `MPCSigner`
   satisfies Go's implicit interface compliance at compile time.
5. The ETH MPC Signing System shall define a separate `MPCKeyShardLoader` port interface in
   `internal/application/ports/` to abstract key shard loading from the signing logic.

---

### Requirement 5: TSS Signing Infrastructure

**Objective:** As a developer, I want a TSS infrastructure component that implements the MPC
port interfaces so that signing rounds can be executed against real TSS library logic.

#### Acceptance Criteria

1. The ETH MPC Signing System shall implement `MPCSigner` in
   `internal/infrastructure/api/eth/mpc/` using the chosen TSS library.
2. When `Sign(ctx, hash)` is called, the ETH MPC Signing System shall initiate a T-of-N signing
   round and coordinate the required message exchange between participant nodes.
3. When the signing round completes, the ETH MPC Signing System shall assemble a 65-byte
   Ethereum-compatible signature `(r || s || v)` compatible with
   `types.Transaction.WithSignature()`.
4. The ETH MPC Signing System shall enforce that `v` is set to the Ethereum-correct value
   (either 27/28 for legacy or chain-ID-derived for EIP-155) matching the transaction's signer type.
5. The ETH MPC Signing System shall implement the node-to-node communication channel as a
   pluggable transport (gRPC or libp2p), isolated behind the `MPCSessionManager` port so that
   the transport can be swapped without modifying use cases.
6. If fewer than T participants are available during a signing round, the ETH MPC Signing System
   shall return an error and not produce a partial or invalid signature.
7. The ETH MPC Signing System shall not hold or reconstruct the full private key at any point
   during the signing round.

---

### Requirement 6: MPC Transaction Creation (Watch Wallet)

**Objective:** As a Watch wallet operator, I want to create a standard Ethereum transaction
for MPC signing so that the transaction parameters are shared with all TSS signing nodes.

#### Acceptance Criteria

1. When the `create mpc` command is invoked, the ETH MPC Signing System shall build a standard
   `types.Transaction` (Legacy or EIP-1559) using the existing `TxCreator` port, reusing the
   single-sig transaction construction logic.
2. The ETH MPC Signing System shall compute the transaction hash (`tx.Hash()`) and write it,
   along with the unsigned raw transaction bytes, to a new `ETHMPCTransactionFile` JSON file.
3. The `ETHMPCTransactionFile` shall include: `from`, `to`, `value`, `nonce`, `gas`, fee fields,
   `chain_id`, `tx_hash` (the hash to be signed by TSS), `raw_tx_hex`, and `tx_type: "unsigned"`.
4. The ETH MPC Signing System shall not call `execTransaction` or interact with any smart
   contract — the MPC flow always uses `eth_sendRawTransaction`.
5. The ETH MPC Signing System shall print the generated file path to stdout on success.
6. The `ETHMPCTransactionFile` format shall be distinct from `ETHTransactionFile` (single-sig)
   and `ETHMultisigTransactionFile` (Safe) — no format mixing is permitted.

---

### Requirement 7: MPC Signing Session (MPC Node Participants)

**Objective:** As an MPC node operator, I want to participate in a TSS signing round so that
my key shard contributes to producing the final ECDSA signature for a pending transaction.

#### Acceptance Criteria

1. When a signing session is initiated, the ETH MPC Signing System shall load the node's key
   shard from encrypted local storage before starting the round.
2. The ETH MPC Signing System shall verify that the `tx_hash` in the `ETHMPCTransactionFile`
   matches the hash derived from `raw_tx_hex` before participating in the signing round.
3. If `tx_hash` verification fails, the ETH MPC Signing System shall abort and return an error —
   no signing shall occur.
4. While the signing round is in progress, the ETH MPC Signing System shall exchange the required
   MPC messages with the other T-1 participant nodes via the configured transport.
5. When the signing round completes, the ETH MPC Signing System shall write the resulting 65-byte
   signature into the `ETHMPCTransactionFile` (populating `signed_tx_hex` by applying the
   signature to the raw transaction via `types.Transaction.WithSignature()`).
6. The ETH MPC Signing System shall update `tx_type` to `"signed"` in the output file when the
   signature is successfully applied.
7. The ETH MPC Signing System shall derive and verify the sender address from the signed
   transaction, confirming it matches the `from` field in the file.
8. If fewer than T nodes confirm participation before a configurable timeout, the ETH MPC Signing
   System shall cancel the session and return an error.

---

### Requirement 8: MPC Transaction Broadcast (Watch Wallet)

**Objective:** As a Watch wallet operator, I want to broadcast a TSS-signed transaction to the
Ethereum network so that it is treated as a standard EOA transaction by the network.

#### Acceptance Criteria

1. When the `send mpc` command is invoked with a `--file` flag pointing to a signed
   `ETHMPCTransactionFile`, the ETH MPC Signing System shall validate that `tx_type == "signed"`.
2. If `tx_type != "signed"`, the ETH MPC Signing System shall reject the file with a descriptive
   error message.
3. The ETH MPC Signing System shall broadcast the signed transaction via `eth_sendRawTransaction`,
   reusing the existing `TxSender` port — no Safe-specific or contract-call path is used.
4. The ETH MPC Signing System shall print the transaction hash to stdout on successful broadcast.
5. The ETH MPC Signing System shall poll for transaction receipt using the existing monitoring
   pattern, confirming inclusion in a block before reporting success.

---

### Requirement 9: CLI Commands for MPC Flow

**Objective:** As a wallet operator, I want CLI commands for the complete MPC transaction
lifecycle so that the MPC flow can be driven from the command line consistently with existing flows.

#### Acceptance Criteria

1. The ETH MPC Signing System shall expose a `watch create mpc` command with flags:
   `--from` (sender address), `--to`, `--amount`, and `--action-type`.
2. The ETH MPC Signing System shall expose a `watch send mpc` command with a `--file` flag.
3. The ETH MPC Signing System shall expose a `keygen dkg` command to initiate a DKG ceremony
   with flags: `--participants` (N) and `--threshold` (T).
4. The ETH MPC Signing System shall expose a `keygen sign mpc` command with a `--file` flag,
   triggering this node's participation in the TSS signing round.
5. All commands shall validate inputs at the CLI layer before invoking use cases, returning
   descriptive errors for invalid inputs.
6. All CLI commands shall follow the existing Cobra command structure in
   `internal/interface-adapters/cli/`.

---

### Requirement 10: DI Container Wiring

**Objective:** As a developer, I want all MPC use cases and infrastructure components wired
in the DI container so that the MPC flow composes correctly at startup.

#### Acceptance Criteria

1. The ETH MPC Signing System shall instantiate the TSS infrastructure implementation
   (`MPCSignerImpl`) in `internal/di/container.go` and inject it into all use cases that
   require `MPCSigner`.
2. The ETH MPC Signing System shall wire the `CreateMPCTransactionUseCase`,
   `SendMPCTransactionUseCase` (Watch), and `SignMPCTransactionUseCase` (Keygen nodes) in the
   DI container.
3. All factory functions shall follow the existing `newXxx` naming convention.
4. If the MPC transport (gRPC / libp2p) is not configured, the ETH MPC Signing System shall
   return a descriptive startup error rather than panicking.
5. All existing single-sig EOA and Safe multisig DI wiring shall remain unchanged — no regressions.

---

### Requirement 11: E2E Test — MPC-TSS Multisig Payment (P4)

**Objective:** As a QA engineer, I want a complete end-to-end test for the MPC-TSS flow against
a local Ethereum node so that the full T-of-N signing lifecycle is validated automatically.

#### Acceptance Criteria

1. A new E2E test script `scripts/operation/eth/e2e/e2e-p4.sh` shall cover the 2-of-3 MPC-TSS
   flow with the following steps:
   - DKG: 3 nodes run DKG ceremony; joint Ethereum address derived; key shards stored locally
   - Fund: ETH sent to the joint address via Anvil pre-funded account
   - Create: Watch wallet creates `ETHMPCTransactionFile` via `create mpc`
   - Sign: 2 of 3 nodes participate in TSS signing round via `sign mpc`; `tx_type` set to `"signed"`
   - Send: Watch wallet broadcasts via `send mpc`
   - Verify: Transaction confirmed; recipient balance updated correctly
2. `make eth-e2e-p4` and `make eth-e2e-p4-ci` targets shall be added to the Makefile.
3. The test shall run against Anvil (local) using deterministic accounts and pass reliably in CI.
4. The E2E parallel runner (`e2e-parallel-runner.sh`) shall support P4 alongside P1 and P3.
5. All existing P1 (single-sig) and P3 (Safe multisig) E2E tests shall continue to pass
   without modification after P4 is added.

---

## Out of Scope

| Item | Rationale |
|------|-----------|
| ERC-20 token MPC transfers | Initial scope is native ETH only; extension in a future spec |
| Safe + MPC hybrid (Safe owners replaced by MPC nodes) | Separate design concern; not required for baseline MPC |
| MPC key resharing / threshold change | Advanced TSS feature; not required for initial implementation |
| Web UI or HTTP coordination server for MPC sessions | No HTTP/REST layer in this project |
| BTC or XRP MPC signing | Separate chain specs required |
| Hardware Security Module (HSM) integration | Future extension; not in scope |
| Production key management / backup ceremony | Operational concern; E2E test uses local ephemeral shards |

---

## Constraints

1. **Architecture**: Must follow Clean Architecture. Domain layer has zero infrastructure dependencies.
2. **No full private key**: The full TSS private key must never be reconstructed or stored anywhere
   in any flow — DKG or signing.
3. **Offline signing principle**: Key shard operations must be executable without network access
   to the Ethereum node; only the Watch wallet communicates with the node.
4. **Transport decoupling**: The inter-node MPC communication transport (gRPC / libp2p) must be
   isolated behind a port interface and not leak into use case or domain layers.
5. **Library**: Use `go-ethereum` for all Ethereum transaction operations. The TSS library is used
   only for key generation and signing round logic.
6. **Additive-only — no regressions**: MPC-TSS (P4) is strictly additive. All existing
   single-sig EOA (P1) and Safe multisig (P3) flows — including CLI commands, use cases, ports,
   infrastructure, DI wiring, and E2E scripts — must remain unchanged and fully functional.
   The Safe implementation is not deprecated, replaced, or feature-frozen by this spec.
7. **No Safe SDK**: Do not add a Safe SDK dependency. The MPC flow has no Safe contract interaction.
8. **Standard transaction**: The MPC signing output must produce a transaction compatible with
   `eth_sendRawTransaction` — not `execTransaction`.
9. **File-based state**: The `ETHMPCTransactionFile` is the sole state carrier between wallet hops,
   consistent with the existing Safe multisig file-based pattern.
