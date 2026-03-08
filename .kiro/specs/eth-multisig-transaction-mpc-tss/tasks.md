# Implementation Plan

## Feature: ETH Multisig — MPC/TSS (Pattern 4)

---

- [x] 1. Add TSS library dependency and define all MPC port interfaces
- [x] 1.1 Add `github.com/bnb-chain/tss-lib/v2` v2.0.1 and `google.golang.org/grpc` to `go.mod`
  - Run `go get github.com/bnb-chain/tss-lib/v2@v2.0.1` and verify the module graph resolves cleanly
  - Add `google.golang.org/grpc` (any compatible version) required by the gRPC transport implementation
  - Confirm `make check-build` passes after the dependency additions
  - _Requirements: 1.1, 1.3, 1.5_

- [x] 1.2 Define all MPC port interfaces in `internal/application/ports/api/eth/interfaces_mpc.go`
  - Define request/result structs `MPCSigningRequest`, `MPCSigningResult`, `DKGParams`, `DKGResult` using only Go primitives and `context.Context`
  - Define `MPCTransactionSigner` with `SignTransaction(ctx, req) (*MPCSigningResult, error)`
  - Define `MPCKeyGeneratorPort` with `RunDKG(ctx, params) (*DKGResult, error)` and `GeneratePreParams(ctx, outputPath) error`
  - Define `MPCKeyShardStorage` with `LoadShard` and `SaveShard` methods
  - Define `MPCOutboundTransport` (Watch-side: `Send`, `Receive`, `Close`) and `MPCInboundTransport` (Node-side: `Listen`, `Receive`, `Close`)
  - Define composed interfaces `MPCCoordinatorDeps` (embeds `MPCTransactionSigner`) and `MPCNodeDeps` (embeds `MPCKeyGeneratorPort`) for DI wiring
  - Verify no TSS library types, gRPC types, or protobuf types appear anywhere in this file
  - _Requirements: 1.2, 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 1.3 Define `ETHMPCTransactionFile` DTO and `MPCFileRepositorier` port
  - Add `ETHMPCTransactionFile` struct to `internal/application/dto/eth/mpc_transaction_file.go` with all fields: version, tx_type, uuid, action_type, from/to/value/nonce/gas_limit/chain_id, EIP-1559 fee fields, tx_hash, raw_tx_hex, signed_tx_hex (omitempty), threshold, party_ids
  - Implement `Validate()` method enforcing version ≥ 1, valid tx_type enum, chain_id > 0, valid EIP-55 addresses, non-empty tx_hash, `len(PartyIDs) >= Threshold`, and when signed: non-empty SignedTxHex
  - Define sentinel errors in a `var` block following the `ETHMultisigTransactionFile` pattern
  - Define `MPCFileRepositorier` interface in `internal/application/ports/file/mpc_file.go` with `ReadETHMPCJSONFile`, `WriteETHMPCJSONFile`, and `CreateMPCFilePath`
  - _Requirements: 6.2, 6.3, 6.6, 8.1_

---

- [x] 2. Implement encrypted key shard storage and MPC transaction file I/O
- [x] 2.1 (P) Implement `MPCShardStore` — encrypted shard file read/write
  - Implement `internal/infrastructure/storage/file/mpc/shard_store.go` satisfying `MPCKeyShardStorage`
  - Derive a 32-byte AES-256-GCM key from the passphrase using scrypt (N=1<<15, r=8, p=1)
  - Encrypt the JSON bundle `{version, party_id, all_party_ids, threshold, eth_address, pre_params, save_data}` as a single AES-256-GCM ciphertext
  - On `LoadShard`: decrypt, unmarshal, verify `party_id` matches expected node identity, return opaque JSON bytes, zero plaintext after return
  - On `SaveShard`: marshal to JSON, encrypt, write atomically (write temp file then rename) to prevent corruption
  - Return descriptive error (never raw key bytes) if file is missing or passphrase is wrong
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

- [x] 2.2 (P) Implement the `MPCFileRepositorier` for `ETHMPCTransactionFile` JSON persistence
  - Implement `internal/infrastructure/storage/file/mpc/mpc_file_repository.go` satisfying `MPCFileRepositorier`
  - `WriteETHMPCJSONFile`: marshal to indented JSON, write to `{action_type}_mpc_{uuid}.json` (unsigned) or `{action_type}_mpc_{uuid}_signed.json` (signed)
  - `ReadETHMPCJSONFile`: unmarshal and call `Validate()` before returning; return `ErrNotUnsigned` or `ErrNotSigned` sentinel appropriately
  - `CreateMPCFilePath`: compose file path deterministically from action type, UUID, and signed flag
  - _Requirements: 6.2, 7.5, 8.1_

- [x] 2.3* Implement unit tests for `MPCShardStore` and `MPCFileRepositorier`
  - `shard_store_test.go`: save+load round-trip with correct passphrase; verify decryption failure with wrong passphrase; verify missing file error
  - `mpc_file_repository_test.go`: write unsigned file, read it back; attempt to read a signed file as unsigned and verify sentinel error
  - _Requirements: 3.1, 3.4, 6.2_

---

- [ ] 3. Define gRPC proto schema and implement MPC transport layer
- [x] 3.1 Define proto schema for MPC message relay and generate gRPC stubs
  - Add `mpc.proto` to the project's existing proto directory, following the XRP gRPC proto conventions
  - Define RPC methods for: submitting a signing request (session ID + hash + party config), relaying a TSS round message (session ID + from_party + payload bytes), and collecting results
  - Generate Go stubs via `make proto` (or equivalent) and verify the generated files compile
  - _Requirements: 5.5_

- [x] 3.2 (P) Implement `GRPCOutboundTransport` — Watch-side gRPC client
  - Implement in `internal/infrastructure/api/eth/mpc/grpc_outbound.go` satisfying `MPCOutboundTransport`
  - `Send(ctx, peerAddr, msg)`: open (or reuse) a gRPC connection to `peerAddr` and stream the TSS message bytes
  - `Receive(ctx)`: return a channel that delivers inbound messages routed back from node responses
  - `Close()`: drain the channel and close all open gRPC connections
  - Include exponential backoff on connection failures; surface the error after context deadline
  - _Requirements: 5.5_

- [x] 3.3 (P) Implement `GRPCInboundTransport` — Node-side gRPC server
  - Implement in `internal/infrastructure/api/eth/mpc/grpc_inbound.go` satisfying `MPCInboundTransport`
  - `Listen(ctx, listenAddr)`: start a gRPC server accepting TSS relay messages from the coordinator
  - `Receive(ctx)`: expose a channel of inbound message bytes delivered by the gRPC server handler
  - `Close()`: gracefully stop the gRPC server and close the channel
  - Validate session ID on each inbound message; discard messages with unknown or mismatched session IDs
  - _Requirements: 5.5, 7.1_

---

- [ ] 4. Implement TSS coordinator and node server using `tss-lib/v2`
- [x] 4.1 Implement `MPCCoordinator` — Watch-side TSS signing coordinator
  - Implement in `internal/infrastructure/api/eth/mpc/coordinator.go` satisfying `MPCCoordinatorDeps` (compile-time check: `var _ apieth.MPCCoordinatorDeps = (*MPCCoordinator)(nil)`)
  - `SignTransaction(ctx, req)`: validate `len(req.Hash) == 32` and `len(req.PartyIDs) >= req.Threshold`; open connections to all T nodes via `MPCOutboundTransport`; fan out signing request
  - Act as message bus: receive TSS round messages from each node via `MPCOutboundTransport.Receive()` and forward to the appropriate T-1 other nodes until all 9 GG18 rounds complete
  - Assemble the 65-byte signature from `common.SignatureData`: `r[32] || s[32] || v[1]`, adjusting `v` to `v_recovery + 27` for legacy signers or chain-ID-encoded for EIP-1559 signers
  - Never hold or reconstruct the full private key; the coordinator is a relay-only participant
  - Return wrapped error (with session ID) if any node drops or context times out
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_

- [x] 4.2 Implement `MPCNodeServer` — Node-side DKG and signing participant
  - Implement in `internal/infrastructure/api/eth/mpc/node_server.go` satisfying `MPCNodeDeps` (compile-time check: `var _ apieth.MPCNodeDeps = (*MPCNodeServer)(nil)`)
  - `GeneratePreParams(ctx, outputPath)`: call `keygen.GeneratePreParams(timeout)` from `tss-lib/v2`; marshal to JSON; write to `outputPath`
  - `RunDKG(ctx, params)`: load pre-params from `params.PreParamsPath`; create `keygen.LocalParty` with the sorted party ID set; drive DKG rounds by reading from `MPCInboundTransport.Receive()` and routing outbound messages back to the coordinator; on completion, marshal `LocalPartySaveData` to JSON and pass to `MPCKeyShardStorage.SaveShard`; return joint Ethereum address and public key
  - Abort DKG if any round returns an error; do not persist a partial shard
  - Handle signing participation: on receiving a session request from coordinator, load shard, create `signing.LocalParty`, drive signing rounds via the inbound message channel
  - Include session ID in every outbound message to prevent cross-session routing
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 7.1, 7.2, 7.3, 7.4_

---

- [x] 5. Implement application use cases for the MPC flow
- [x] 5.1 (P) Implement `RunDKGUseCase` on the Keygen wallet
  - Add `RunDKGUseCase` interface, `RunDKGInput`, and `RunDKGOutput` to `internal/application/usecase/keygen/interfaces.go`
  - Implement in `internal/application/usecase/keygen/eth/run_dkg.go`
  - Validate `len(input.AllPartyIDs) >= input.Threshold` before calling the port; return a descriptive error if invalid
  - Delegate DKG execution to `MPCKeyGeneratorPort.RunDKG`; if the port returns an error, propagate it without saving any shard
  - Print the joint Ethereum address to stdout on success
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

- [x] 5.2 (P) Implement `CreateMPCTransactionUseCase` on the Watch wallet
  - Add `CreateMPCTransactionUseCase` interface and DTOs to `internal/application/usecase/watch/interfaces.go`
  - Implement in `internal/application/usecase/watch/eth/create_mpc_transaction.go`
  - Validate EIP-55 addresses before building the transaction; return error on invalid input
  - Use `TxCreator` port to build a standard EIP-1559 `types.Transaction`; compute `tx.Hash()` as the TSS pre-image
  - Populate all `ETHMPCTransactionFile` fields (including `tx_hash`, `raw_tx_hex`, `tx_type: "unsigned"`, UUID, threshold, party_ids) and write the file via `MPCFileRepositorier`; do not interact with any smart contract
  - Print the generated file path to stdout on success
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [x] 5.3 (P) Implement `SendMPCTransactionUseCase` on the Watch wallet
  - Add `SendMPCTransactionUseCase` interface and DTOs to `internal/application/usecase/watch/interfaces.go`
  - Implement in `internal/application/usecase/watch/eth/send_mpc_transaction.go`
  - Read `ETHMPCTransactionFile` via `MPCFileRepositorier`; return `ErrNotUnsigned` if `TxType != "unsigned"`
  - Build `MPCSigningRequest` from file fields (session UUID, tx_hash, party_ids, threshold, peer_addrs from input) and call `MPCTransactionSigner.SignTransaction`
  - Apply the 65-byte signature to the raw transaction via `types.Transaction.WithSignature`; verify that `types.Sender(signer, signedTx) == file.From` before proceeding; return `ErrSenderMismatch` if not
  - Broadcast the signed transaction via `TxSender.SendSignedRawTransaction`; write the signed file via `MPCFileRepositorier`; print the transaction hash to stdout
  - If the TSS session times out, preserve the unsigned file and return a wrapped error with session ID
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

- [x] 5.4 (P) Implement `ServeMPCUseCase` on the Keygen wallet (node daemon)
  - Add `ServeMPCUseCase` interface and `ServeMPCInput` to `internal/application/usecase/keygen/interfaces.go`
  - Implement in `internal/application/usecase/keygen/eth/serve_mpc.go`
  - On startup, load the key shard from `MPCKeyShardStorage.LoadShard` using the passphrase from input
  - Call `MPCInboundTransport.Listen(ctx, listenAddr)` to start accepting coordinator connections
  - Drive the `tss-lib/v2` signing state machine by reading from `MPCInboundTransport.Receive()`; for each inbound signing request, verify `hash(raw_tx_hex) == tx_hash` before participating — abort and return `ErrTxHashMismatch` if they do not match
  - Block until context is cancelled (graceful shutdown); write the signed partial result via the transport channel
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8_

- [x] 5.5* Implement unit tests for each use case
  - `run_dkg_test.go`: mock `MPCKeyGeneratorPort` and `MPCKeyShardStorage`; verify DKG abort on invalid threshold; verify no shard saved on port error
  - `create_mpc_transaction_test.go`: mock `TxCreator`, `GasEstimator`, `MPCFileRepositorier`; verify all file fields are populated; verify invalid address is rejected
  - `send_mpc_transaction_test.go`: mock `MPCTransactionSigner`, `TxSender`, `MPCFileRepositorier`; verify `ErrNotUnsigned` guard; verify `ErrSenderMismatch` is returned when recovered address does not match `from`
  - `serve_mpc_test.go`: mock `MPCKeyShardStorage`, `MPCInboundTransport`; verify `ErrTxHashMismatch` is returned when hash verification fails
  - _Requirements: 7.2, 7.3, 8.1, 8.2_

---

- [x] 6. Update wallet adapters and add CLI commands for the MPC flow
- [x] 6.1 (P) Extend `ETHWatch` wallet adapter with MPC use case fields and methods
  - Add `createMPCTxUseCase watchusecase.CreateMPCTransactionUseCase` and `sendMPCTxUseCase watchusecase.SendMPCTransactionUseCase` fields to `internal/interface-adapters/wallet/eth/watch.go`
  - Add public method `CreateMPCTx(ctx, from, to string, amount float64, threshold int, partyIDs []string, actionType string) (filePath string, txHash string, err error)` delegating to the use case
  - Add public method `SendMPCTx(ctx, filePath string, peerAddrs []string) (txHash string, err error)` delegating to the use case
  - _Requirements: 9.1, 9.2, 9.5, 9.6_

- [x] 6.2 (P) Extend `ETHKeygen` wallet adapter with DKG and serve use case fields and methods
  - Add `runDKGUseCase keygenusecase.RunDKGUseCase` and `serveMPCUseCase keygenusecase.ServeMPCUseCase` fields to `internal/interface-adapters/wallet/eth/keygen.go`
  - Add public method `RunDKG(ctx context.Context, input keygenusecase.RunDKGInput) (keygenusecase.RunDKGOutput, error)` delegating to the use case
  - Add public method `ServeMPC(ctx context.Context, input keygenusecase.ServeMPCInput) error` delegating to the use case
  - Add public method `GeneratePreParams(ctx context.Context, outputPath string) error` calling `MPCKeyGeneratorPort.GeneratePreParams`
  - _Requirements: 9.3, 9.4, 9.5, 9.6_

- [x] 6.3 (P) Add `watch create mpc` and `watch send mpc` Cobra CLI commands
  - Add `internal/interface-adapters/cli/watch/create/mpc.go`: validate `--from`, `--to`, `--amount`, `--action-type`, `--threshold`, `--party-ids` flags; call `ETHWatch.CreateMPCTx`; print the file path on success
  - Add `internal/interface-adapters/cli/watch/send/mpc.go`: validate `--file` and `--peer-addrs` flags; call `ETHWatch.SendMPCTx`; print the transaction hash on success
  - Return a descriptive error for any missing or invalid CLI input without invoking use cases
  - _Requirements: 9.1, 9.2, 9.5, 9.6_

- [x] 6.4 (P) Add `keygen dkg`, `keygen pre-params`, and `keygen serve mpc` Cobra CLI commands
  - Add `internal/interface-adapters/cli/keygen/dkg/dkg.go` with `dkg` subcommand (flags: `--threshold`, `--parties`, `--party-id`, `--peers`, `--pre-params-path`, `--shard-output`, `--passphrase`); call `ETHKeygen.RunDKG`
  - Add `pre-params` subcommand (flag: `--output`); call `ETHKeygen.GeneratePreParams`; print the output path on success
  - Add `internal/interface-adapters/cli/keygen/serve/mpc.go` with `serve mpc` subcommand (flags: `--listen-addr`, `--shard-path`, `--passphrase`, `--party-id`, `--all-party-ids`); call `ETHKeygen.ServeMPC`; run until SIGINT
  - Validate all required flags at the CLI layer and return descriptive errors before calling into use cases
  - _Requirements: 9.3, 9.4, 9.5, 9.6_

---

- [ ] 7. Wire all MPC components in the DI container
- [ ] 7.1 Add MPC infrastructure factory functions and wire use cases in `internal/di/container.go`
  - Add factory functions `newMPCShardStore`, `newGRPCOutboundTransport`, `newGRPCInboundTransport`, `newMPCCoordinator`, `newMPCNodeServer` following the existing `newXxx` naming convention
  - Wire `RunDKGUseCase`, `ServeMPCUseCase`, `CreateMPCTransactionUseCase`, and `SendMPCTransactionUseCase` constructors, injecting the corresponding port implementations
  - If gRPC listen address is not configured at startup, return a descriptive error rather than panicking; add a configuration check in the startup path
  - Inject new use cases into `ETHWatch` and `ETHKeygen` adapter constructors
  - _Requirements: 10.1, 10.2, 10.3, 10.4_

- [ ] 7.2 Verify all existing DI wiring remains unchanged and the build passes
  - Confirm `make check-build` and `make go-lint` both pass with no regressions in P1 (single-sig) and P3 (Safe multisig) wiring
  - Confirm all pre-existing unit tests pass
  - _Requirements: 10.5_

---

- [ ] 8. Add E2E test for MPC-TSS Pattern 4 (P4) and Makefile targets
- [ ] 8.1 Write `scripts/operation/eth/e2e/e2e-p4.sh` covering the full 2-of-3 MPC-TSS flow
  - Generate pre-params for 3 nodes (`keygen pre-params`)
  - Run DKG ceremony for all 3 nodes concurrently; capture the joint Ethereum address output
  - Fund the joint address from an Anvil pre-funded account
  - Create an `ETHMPCTransactionFile` via `watch create mpc`; verify the file exists and contains `tx_type: unsigned`
  - Start 2 of the 3 nodes as gRPC server daemons (`keygen serve mpc`) in the background
  - Run `watch send mpc` with the unsigned file and the 2 node addresses; verify exit code 0
  - Check the recipient balance on Anvil confirms the expected transfer amount
  - Verify the signed `ETHMPCTransactionFile` exists with `tx_type: signed` and a non-empty `signed_tx_hex`
  - _Requirements: 11.1, 11.3_

- [ ] 8.2 (P) Add `make eth-e2e-p4` and `make eth-e2e-p4-ci` Makefile targets
  - Add targets to the Makefile that invoke `e2e-p4.sh`, following the existing `eth-e2e-p1` and `eth-e2e-p3` target conventions
  - _Requirements: 11.2_

- [ ] 8.3 (P) Update the E2E parallel runner to include P4
  - Update `scripts/operation/eth/e2e/e2e-parallel-runner.sh` to add P4 as a parallel pattern alongside P1 and P3
  - Assign P4 an isolated set of Anvil accounts to prevent nonce conflicts with P1 and P3 runs
  - Confirm P1 and P3 test runs remain unmodified and pass in the updated parallel runner
  - _Requirements: 11.4, 11.5_
