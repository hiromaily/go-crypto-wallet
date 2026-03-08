# Research & Design Decisions

---

## Summary

- **Feature**: `eth-multisig-transaction-mpc-tss`
- **Discovery Scope**: Complex Integration (new cryptographic subsystem + new networking layer)
- **Key Findings**:
  - `github.com/bnb-chain/tss-lib/v2` (v2.0.1) is the industry-standard Go TSS/ECDSA library and is not yet in go.mod.
  - MPC-TSS nodes **cannot be air-gapped**: the TSS protocol is an interactive multi-round message exchange that requires real-time network connectivity between all T participants — fundamentally different from Safe's offline file-transfer model.
  - The existing project's Clean Architecture is a good fit: port interfaces cleanly isolate the TSS library and the inter-node transport so that neither leaks into use cases.

---

## Research Log

### TSS Library Selection

- **Context**: Requirements specify using a Go TSS library. We need to verify the best option for 2025-2026.
- **Sources Consulted**: GitHub bnb-chain/tss-lib, Coinbase kryptology, ZenGo awesome-tss list.
- **Findings**:
  - `github.com/bnb-chain/tss-lib/v2` v2.0.1: battle-tested (Binance, THORChain, ZetaChain), secp256k1 ECDSA, full audit (Kudelski 2019), v2 fixed the Fireblocks/Tsshock vulnerabilities (2023). Protocol: GG18, ~9 signing rounds.
  - `coinbase/kryptology`: broader scope, uses GG20 (declared obsolete by its authors), or DKLS18 for 2-of-2. Less suitable for T-of-N.
  - No TSS library is currently in the project's `go.mod`.
- **Implications**: Use `github.com/bnb-chain/tss-lib/v2` as the single TSS dependency. Confine it to `internal/infrastructure/`.

### tss-lib v2 API Surface

- **Context**: Need to understand the library's keygen and signing API to design port interfaces correctly.
- **Findings**:
  - **Pre-params** (expensive, one-time per node): `keygen.GeneratePreParams(timeout)` produces Paillier key pair, NTilde, H1, H2. Must be persisted with the key shard.
  - **DKG**: each party creates a `keygen.LocalParty`, drives rounds via `out chan tss.Message` and `end chan keygen.LocalPartySaveData`. Requires out-of-band message routing between all N parties.
  - **Signing**: each party creates a `signing.LocalParty` (takes `LocalPartySaveData` from DKG), drives rounds similarly. Output is `common.SignatureData` containing R, S, Signature (concatenated) bytes.
  - Session identity: parties are identified by `tss.PartyID`; all parties must share the same sorted `tss.SortedPartyIDs` at ceremony start.
  - The library is **message-based and stateless between rounds** — the transport is the caller's responsibility.
- **Implications**: The infrastructure layer wraps the tss-lib message channels. The port interface (`MPCTransactionSigner`) abstracts the multi-round exchange into a single synchronous `Sign(ctx, hash) ([]byte, error)` call.

### Inter-node Transport

- **Context**: TSS requires interactive multi-round message exchange between all T participants simultaneously.
- **Findings**:
  - MPC-TSS is inherently online and networked — air-gap separation (USB file transfer) is **not possible** for the signing round.
  - gRPC is the simplest reliable transport for Go services; existing codebase already has a proto setup (XRP gRPC). Using gRPC for MPC messaging is consistent.
  - The transport must be pluggable behind a port to allow future replacement with libp2p.
- **Implications**: Design `MPCNodeTransport` port interface in `application/ports/`. Infrastructure provides a gRPC implementation. The requirements state "No gRPC in Safe flow" — this constraint applies to Safe only; MPC-TSS has a network transport by design.

### Wallet Model Change for MPC

- **Context**: The existing 3-wallet model (Watch online / Keygen+Sign air-gapped) must be evaluated for MPC.
- **Findings**:
  - Safe P3: Keygen and Sign wallets are fully air-gapped; signing is done offline, results in a file.
  - MPC P4: All T signing nodes must exchange protocol messages in real time. Air-gap is not possible.
  - MPC nodes (previously Keygen/Sign roles) must expose a **listening service** (gRPC server) during signing sessions.
  - The DKG ceremony is also interactive but can be run in a controlled, secure network environment.
- **Implications**: For P4, the Keygen and Sign wallets run as gRPC server daemons (`keygen serve mpc`) when acting as MPC nodes. The Watch wallet is the client that initiates signing sessions. This is documented explicitly as a departure from the air-gap model.

### Key Shard Storage

- **Context**: Each MPC node must persist its `LocalPartySaveData` and pre-params securely.
- **Findings**:
  - `LocalPartySaveData` is a Go struct serializable to JSON. It contains the private share, public keys, commitments.
  - Pre-params are also JSON-serializable.
  - Existing project uses go-ethereum `keystore` (scrypt) for ECDSA key encryption. Similar approach (AES-GCM + passphrase-derived key) is appropriate for shard files.
- **Implications**: New `MPCShardStore` in `internal/infrastructure/storage/file/mpc/` persists an encrypted JSON bundle: `{pre_params, save_data, party_metadata}`.

### Ethereum Signature Format

- **Context**: Need tss-lib output to be compatible with `types.Transaction.WithSignature`.
- **Findings**:
  - `common.SignatureData.Signature` is `r[32] || s[32] || v[1]` where `v` is 0 or 1 (recovery bit).
  - `go-ethereum`'s `types.Transaction.WithSignature(signer, sig)` expects a 65-byte signature where the `v` byte is chain-ID-encoded (for EIP-155) or 27/28 (for legacy).
  - The infrastructure layer must adjust `v` after TSS output: `v_adjusted = v_recovery_bit + 27` for legacy, or use `signer.SignatureValues` to encode EIP-155.
  - After applying, verify the recovered sender matches the expected distributed EOA address.
- **Implications**: Signature post-processing (v adjustment + sender verification) lives in the infrastructure TSS client, not in the use case.

### Transaction File Format

- **Context**: Need a new DTO distinct from `ETHTransactionFile` (single-sig) and `ETHMultisigTransactionFile` (Safe).
- **Findings**:
  - `ETHMPCTransactionFile` needs: `from`, `to`, `value`, `nonce`, `gas`, fee fields, `chain_id`, `tx_hash` (the pre-image to sign), `raw_tx_hex` (unsigned, for verification), `signed_tx_hex` (filled after TSS), `tx_type`, `uuid`, `action_type`.
  - The `tx_hash` field is the 32-byte Keccak256 hash that TSS nodes will sign. Nodes verify `hash(raw_tx_hex) == tx_hash` before participating.
  - Unlike Safe, there is only one signature (65 bytes) — no `signatures` array.
  - File naming: `{action_type}_mpc_{uuid}.json` (unsigned) → `{action_type}_mpc_{uuid}_signed.json`.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| **A: Watch-coordinated gRPC** (Selected) | Watch wallet is gRPC client. MPC nodes run gRPC servers. Watch sends hash, nodes exchange TSS messages, Watch collects signature. | Consistent with existing Watch-centric model; simple client/server split | Nodes must be network-reachable from Watch | Aligns with existing architecture |
| B: P2P libp2p | All nodes are peers; any node can initiate; Watch receives final sig out-of-band | Fully decentralized, resilient | High complexity; adds libp2p dependency; harder to debug | Useful if Watch fails, but over-engineering for this scope |
| C: File-based relay | TSS messages written to files, parties exchange via USB | Maintains air-gap | Not feasible — TSS requires real-time exchange (multiple rounds in seconds) | Rejected |

---

## Design Decisions

### Decision: Use `tss-lib/v2` as the sole TSS dependency

- **Context**: Need a production-ready Go library for T-of-N ECDSA on secp256k1.
- **Alternatives Considered**:
  1. `coinbase/kryptology` — GG20 (obsolete per authors) or DKLS18 (2-of-2 only)
  2. Custom implementation — insecure, not feasible
- **Selected Approach**: `github.com/bnb-chain/tss-lib/v2` confined to `internal/infrastructure/api/eth/mpc/`
- **Rationale**: Battle-tested in production blockchains; audited; active maintenance; v2 patches critical security issues.
- **Trade-offs**: GG18 protocol is 9-round (~seconds of latency); heavier dependency than custom ECDSA.
- **Follow-up**: Pin exact version (v2.0.1) and add a security audit note to CHANGELOG.

### Decision: MPC nodes are networked (not air-gapped)

- **Context**: TSS requires real-time inter-node message exchange; file-based transfer is not feasible.
- **Alternatives Considered**:
  1. Store-and-forward via shared filesystem — rounds would take hours; not practical
  2. Air-gapped with USB — impossible; multiple parties must respond to each other in lockstep
- **Selected Approach**: MPC nodes run a gRPC server (`keygen serve mpc`); Watch wallet is the orchestrator/client.
- **Rationale**: The TSS protocol is fundamentally online; this is documented clearly as P4's difference from P3.
- **Trade-offs**: Lower security isolation than Safe's air-gap model; mitigated by network ACLs and TLS.
- **Follow-up**: Document network security requirements (mTLS, firewall rules) in the E2E test script and README.

### Decision: Pre-params generated once, stored with key shard

- **Context**: tss-lib requires expensive Paillier key pre-computation per party. Can be done separately.
- **Alternatives Considered**:
  1. Generate inside DKG round 1 — slower DKG ceremony, but simpler
  2. Pre-generate separately — recommended by library; separates concerns
- **Selected Approach**: `keygen pre-params` is a standalone command; pre-params file is referenced by `keygen dkg`.
- **Rationale**: Pre-generation can be done weeks in advance, on more powerful hardware; DKG ceremony runs faster.
- **Trade-offs**: Extra command step in E2E; pre-params file must be secured.
- **Follow-up**: E2E script generates pre-params inline for simplicity; production deployments pre-generate.

### Decision: `MPCTransactionSigner` port hides all multi-round complexity

- **Context**: The use case layer should not know about TSS rounds, message routing, or library types.
- **Alternatives Considered**:
  1. Expose round-by-round API in port — more flexible but leaks TSS model into use cases
  2. Single `Sign(ctx, hash) ([]byte, error)` — simple; all complexity in infrastructure
- **Selected Approach**: `MPCTransactionSigner.SignTransaction(ctx, MPCSigningRequest) (*MPCSigningResult, error)` — one call, returns 65-byte signature.
- **Rationale**: Consistent with existing `TxSigner` and `SafeExecuter` ports; use cases stay clean.
- **Trade-offs**: No progress visibility during the multi-round exchange (use logging in infra layer).

---

## Risks & Mitigations

- **tss-lib dependency size** — The library pulls in Paillier and other heavy crypto; adds ~10MB to binary. Mitigation: acceptable for wallet software.
- **TSS round latency** — 9 rounds of message exchange; network latency multiplies. Mitigation: use localhost/LAN for E2E tests; add context timeout.
- **Key shard loss** — If a node's shard file is lost, the distributed key is irrecoverable. Mitigation: key resharing (out of scope here); operators must backup shard files.
- **gRPC security** — MPC nodes exposed as gRPC servers. Mitigation: require mTLS; note in E2E test and production documentation.
- **Signature `v` adjustment bug** — Incorrect v-byte leads to a signature that doesn't recover the right address. Mitigation: always verify `types.Sender(signer, signedTx) == from` after signing; unit-test this path.
- **Parallel DKG/signing sessions** — tss-lib parties must share exact sorted party IDs; session ID mismatches cause protocol failure. Mitigation: Watch wallet generates a UUID session ID, included in every gRPC message.

---

## References

- [bnb-chain/tss-lib GitHub](https://github.com/bnb-chain/tss-lib) — canonical implementation
- [tss-lib v2 release notes](https://github.com/bnb-chain/tss-lib/releases/tag/v2.0.0) — Tsshock fix
- [Fireblocks Tsshock disclosure (2023)](https://www.fireblocks.com/blog/gg18-and-gg20-paillier-key-vulnerability-technical-report/) — CVE addressed in v2
- [CertiK: GG18 9-round protocol](https://www.certik.com/resources/blog/threshold-cryptography-iii-binance-tss-libs-9-round-threshold-ecdsa) — round-count analysis
- [ZenGo awesome-tss](https://github.com/ZenGo-X/awesome-tss) — comprehensive TSS survey
- [go-ethereum types.Transaction.WithSignature](https://pkg.go.dev/github.com/ethereum/go-ethereum/core/types#Transaction.WithSignature) — signature application API
- [EIP-155](https://eips.ethereum.org/EIPS/eip-155) — replay protection (v encoding)
