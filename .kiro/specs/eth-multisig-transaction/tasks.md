# Implementation Plan

## Task 1. Define the multisig transaction file format, port interfaces, and file I/O contract

- [x] 1.1 Define the `ETHMultisigTransactionFile` data transfer object
  - Create the struct with all Safe execution parameters: safe address, recipient, value (Wei decimal string), call data (hex), operation flag, gas fields (all default to zero), Safe nonce, chain ID, and the precomputed `safeTxHash`
  - Include signature accumulation fields: threshold count and a list of `{signer_address, signature_hex}` entries; `TxType` transitions from `"unsigned"` to `"signed"` when the signature count reaches threshold
  - Include a `UUID` field (generated at proposal time) that serves as the unique proposal identifier in file names
  - Implement `Validate()` enforcing non-empty required fields, valid chain ID, EIP-55 checksum addresses, no duplicate signer entries, and consistency between `TxType` and signature count
  - Define sentinel error variables (`ErrSafeTxHashMismatch`, `ErrDuplicateSigner`, `ErrNotFullySigned`) alongside the struct
  - _Requirements: 2_

- [x] 1.2 Define Safe operation port interfaces and the execution parameter struct
  - Define `SafeExecParams` as a plain struct in the ports package carrying all fields needed to call `execTransaction`: Safe address, recipient, value, data, operation, gas fields, nonce, and the final concatenated signatures bytes
  - Define `SafeInfo` struct (owners list, threshold, nonce, balance in Wei) as a return type for the info interface
  - Define `SafeTxHashComputer` interface with a method to call `getTransactionHash` on the deployed Safe contract and return the hash as a 0x-prefixed hex string
  - Define `SafeNonceReader` interface with a method to fetch the Safe's current nonce from the contract
  - Define `SafeExecuter` interface with a method accepting `SafeExecParams` and returning the submitted Ethereum transaction hash
  - Define `SafeInfoReader` interface with a method returning `SafeInfo` for a given Safe address
  - Define `SafeClientDeps` as a composed interface of all four (for DI use only); add a comment making explicit that use cases depend on the narrow interfaces, not the combined one
  - _Requirements: 3_

- [x] 1.3 Define the multisig file I/O port interface
  - Define `MultisigFileRepositorier` as a separate interface from the existing `TransactionFileRepositorier` with methods to write and read `ETHMultisigTransactionFile` as JSON, and a helper to generate the file path using action type, UUID, and signed count
  - Ensure the existing `TransactionFileRepositorier` interface and all its implementations and mocks are untouched
  - _Requirements: 2, 3_

## Task 2. Generate Safe ABI bindings and implement the Safe contract client

- [x] 2.1 Generate Safe v1.4.1 ABI bindings
  - Obtain the official Safe v1.4.1 ABI JSON from the `safe-global/safe-smart-account` repository
  - Run `abigen` to generate a typed Go client and place the output in the Safe contract bindings package under `internal/infrastructure/contract/safe/`
  - The generated file must carry the `// Code generated - DO NOT EDIT.` header
  - Add a `make safe-abi` Makefile target with a comment documenting the ABI source URL
  - _Requirements: 1_

- [x] 2.2 Implement the `SafeClient` infrastructure component
  - Create a `SafeClient` struct in `internal/infrastructure/api/eth/safe/` wrapping the generated ABI bindings and an Ethereum JSON-RPC client
  - Fetch and cache the chain ID at construction time via `client.ChainID(ctx)` — never hardcode it
  - Implement `GetSafeTxHash`: call `getTransactionHash` on the Safe contract for the given parameters and return the result as a 0x-prefixed hex string
  - Implement `GetSafeNonce`: call `Nonce()` on the Safe contract binding
  - Implement `ExecuteSafeTransaction`: build an EIP-1559 transaction calling `execTransaction`, estimate gas, submit it, and return the transaction hash
  - Implement `GetSafeInfo`: call `GetOwners()`, `GetThreshold()`, `Nonce()` on the Safe binding and `BalanceAt` on the client; return a `SafeInfo` struct
  - The struct satisfies `SafeClientDeps` (all four narrow interfaces)
  - _Requirements: 1, 3_

- [x] 2.3 Implement the `MultisigFileRepositorier` in the file infrastructure
  - Add `WriteETHMultisigJSONFile`, `ReadETHMultisigJSONFile`, and `CreateMultisigFilePath` to the existing concrete file implementation in `internal/infrastructure/file/`
  - `CreateMultisigFilePath` produces paths of the form `{actionType}_multisig_{uuid}_{signedCount}.json`
  - The concrete implementation now satisfies both `TransactionFileRepositorier` and `MultisigFileRepositorier` without any interface change to the former
  - _Requirements: 2_

## Task 3. Implement Watch wallet use cases for multisig transaction lifecycle

- [x] 3.1 (P) Implement the multisig transaction creation use case
  - Define `CreateMultisigTransactionUseCase` interface and its input/output structs in the Watch use case interfaces file; input carries the Safe address, recipient, Ether amount, threshold, and action type
  - Generate a UUID for the proposal, convert Ether to Wei, call `SafeNonceReader` to fetch the current Safe nonce, then call `SafeTxHashComputer` with all-zero gas parameters and empty call data for a simple ETH transfer
  - Populate and write an `ETHMultisigTransactionFile` with `TxType: "unsigned"` and an empty signatures list via `MultisigFileRepositorier`
  - Output the generated file path; no database record is created
  - _Requirements: 4_

- [x] 3.2 (P) Implement the multisig transaction submission use case
  - Define `SendMultisigTransactionUseCase` interface and its input/output structs; input carries the file path, output carries the submitted transaction hash
  - Read and validate the `ETHMultisigTransactionFile`; return `ErrNotFullySigned` if `TxType` is still `"unsigned"`
  - Sort the `Signatures` slice by signer address in ascending order (hex string comparison, case-insensitive) before concatenation — this ordering is required by the Safe contract
  - Decode each `SignatureHex` from hex to bytes and concatenate into a single byte slice; map all file fields to `SafeExecParams`
  - Call `SafeExecuter.ExecuteSafeTransaction` and poll for the transaction receipt using the existing receipt-polling pattern; log the transaction hash on success
  - _Requirements: 6_

- [x] 3.3 (P) Implement the Safe info use case
  - Define `SafeInfoUseCase` interface and its input/output structs; input carries the Safe address, output carries owners list, threshold, nonce, and balance as a Wei decimal string
  - Call `SafeInfoReader.GetSafeInfo` and map the result to the output struct
  - _Requirements: 7_

## Task 4. Implement the offline EIP-712 signing use case

- [ ] 4.1 (P) Implement the EIP-712 hash recomputation helper
  - Implement a pure-Go function (no network calls) that recomputes the `safeTxHash` from the fields stored in `ETHMultisigTransactionFile`
  - Follow the exact hash chain: domain separator typehash → domain separator → Safe TX typehash → struct hash → final EIP-191 prefix hash using `crypto.Keccak256` and ABI-encoding via `go-ethereum/accounts/abi`
  - The function returns the computed hash as a byte slice; the caller compares it to the `SafeTxHash` field in the file and aborts with `ErrSafeTxHashMismatch` if they differ
  - Cover the helper with unit tests using known Safe transaction vectors to confirm the output matches the on-chain value
  - _Requirements: 5_

- [ ] 4.2 (P) Implement the offline multisig signing use case
  - Define `SignMultisigTransactionUseCase` interface and its input/output structs in the keygen use case interfaces file; input carries the file path and signer address (supplied by the CLI); output carries the new file path, completion flag, and counts
  - Read the `ETHMultisigTransactionFile`, invoke the EIP-712 recomputation helper, and abort if the hash does not match
  - Look up the signer's private key by the provided signer address in the account key repository; derive the child private key via the existing BIP-44 derivation helper
  - Sign the hash with `crypto.Sign`, then increment the last byte of the 65-byte result by 27 to produce a Safe-compatible v value; encode the result as a 0x-prefixed hex string
  - Check that the signer address is not already present in `Signatures`; return `ErrDuplicateSigner` if it is
  - Append `{SignerAddress, SignatureHex}` to `Signatures`; set `TxType` to `"signed"` and `IsComplete` to true when `len(Signatures) >= Threshold`
  - Write the updated file with an incremented counter in the path via `MultisigFileRepositorier`
  - _Requirements: 5_

## Task 5. Extend wallet adapters and add CLI commands

- [ ] 5.1 Extend the ETHWatch wallet adapter and add Watch CLI commands
  - Add `createMultisigTxUseCase`, `sendMultisigTxUseCase`, and `safeInfoUseCase` fields to the `ETHWatch` struct; update its constructor accordingly
  - Add adapter methods `CreateMultisigTx`, `SendMultisigTx`, and `GetSafeInfo` that delegate to the respective use cases
  - Add a `watch create multisig` Cobra command with `--safe`, `--to`, `--amount`, `--threshold`, and `--action-type` flags; validate non-empty address fields at the CLI layer before calling the adapter
  - Add a `watch send multisig` Cobra command with a `--file` flag
  - Add a `watch safe info` Cobra command with a `--safe` flag; place it in a new `safe/` subdirectory under the Watch CLI
  - _Requirements: 7, 8_

- [ ] 5.2 (P) Extend the ETHKeygen adapter to route multisig signing
  - Add a `signMultisigTxUseCase` field to `ETHKeygen`; update its constructor
  - In the existing `SignTx` adapter method, detect the file format by reading the JSON and checking for the `safe_address` field: if present, route to `signMultisigTxUseCase.Sign()` with the signer address derived from CLI context; otherwise keep the existing single-sig path
  - Pass `--signer-address` from the CLI flag through to the use case input
  - _Requirements: 5, 8_

- [ ] 5.3 (P) Activate the ETHSign wallet adapter for multisig signing
  - Add a `signMultisigTxUseCase` field to `ETHSign` (currently the struct has no signing capability) and update its constructor
  - Wire the `SignTx` adapter method (currently a no-op returning empty values) to `signMultisigTxUseCase.Sign()` following the same file-type detection and signer address routing as the Keygen adapter
  - The Sign wallet handles multisig only; single-sig ETH signing remains Keygen-only and is unaffected
  - _Requirements: 5, 8_

## Task 6. Wire all new components in the DI container

- [ ] 6.1 Add DI factory functions for SafeClient and all new use cases
  - Add `newSafeClient()` factory that constructs a `SafeClient` using the existing Ethereum client, fetches the chain ID once, and returns the cached instance; the factory is called by each use case factory that needs a Safe port
  - Add factories for `CreateMultisigTransactionUseCase`, `SendMultisigTransactionUseCase`, and `SafeInfoUseCase` for the Watch wallet, injecting the appropriate narrow Safe interface from `SafeClient`
  - Add factories for `SignMultisigTransactionUseCase` for both the Keygen and Sign wallet containers, injecting the account key repository and the multisig file repository
  - _Requirements: 9_

- [ ] 6.2 Update wallet adapter constructors in the DI container
  - Pass the three new Watch use cases into `NewETHWatch` and add the Sign multisig use case to `NewETHKeygen` and `NewETHSign`; verify that no panics occur and that all ETH wallet types compile cleanly
  - Confirm the existing single-sig factories and their injected dependencies are unchanged
  - _Requirements: 9_

## Task 7. Implement the E2E P2 test for 2-of-2 Safe multisig payment

- [ ] 7.1 Write the Foundry deployment script for Safe
  - Create `DeploySafe.s.sol` in `apps/eth-contracts/script/` that accepts owner addresses and a threshold as constructor arguments and deploys a Safe v1.4.1 proxy pointing to the canonical singleton
  - The script prints the deployed Safe address to stdout in a parseable format for use by the shell script
  - The deployment targets Anvil (local) and requires no mainnet credentials
  - _Requirements: 10_

- [ ] 7.2 Write the E2E P2 shell script
  - Create `scripts/operation/eth/e2e/e2e-p2.sh` following the structure of the existing `e2e-p1.sh`
  - Phase 0 — Provisioning: initialise a primary keygen DB (signer 1) and a separate keygen DB that acts as the Sign wallet (signer 2); generate HD keys in each; export signer addresses
  - Phase 1 — Setup: call the Foundry `DeploySafe.s.sol` script with both signer addresses and threshold=2; fund the resulting Safe contract address via Anvil
  - Phase 2 — Create: invoke `watch create multisig` to produce the unsigned multisig file
  - Phase 3 — Sign 1: invoke the Keygen wallet `sign tx` with `--signer-address signer1_addr`; verify signature count is 1
  - Phase 4 — Sign 2: invoke the second keygen DB (Sign wallet role) `sign tx` with `--signer-address signer2_addr`; verify `TxType` is `"signed"`
  - Phase 5 — Send: invoke `watch send multisig`; capture the transaction hash
  - Phase 6 — Verify: assert the recipient ETH balance increased and the Safe ETH balance decreased by the expected amount
  - Support `--cleanup`, `--reset`, `--non-interactive`, and `--verbose` flags consistent with the existing E2E scripts
  - _Requirements: 10_

- [ ] 7.3 Add Makefile targets for the P2 E2E test
  - Add `eth-e2e-p2` target that runs `e2e-p2.sh` interactively
  - Add `eth-e2e-p2-ci` target that runs with `--non-interactive` for CI use
  - Place the targets in `make/wallet/eth_e2e.mk` alongside the existing P1 targets
  - _Requirements: 10_
