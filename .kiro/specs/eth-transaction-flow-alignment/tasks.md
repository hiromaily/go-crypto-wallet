# Implementation Plan

- [x] 1. Modernize ETH configuration with network types, node selection, and fee parameters
- [x] 1.1 (P) Add network type constants and replace deprecated network references
  - Define `EthNetworkType` constants: mainnet, sepolia, holesky, local (removing Goerli, Rinkeby, Ropsten)
  - Add `ChainID` field with auto-population from `NetworkType` when not explicitly set
  - Add EIP-1559 fee config fields: `MaxFeePerGasCap` (absolute fee ceiling in Gwei)
  - Add `KeystorePassword` field to replace the hardcoded `Password = "password"` constant
  - Validate on load that the network name is supported; return an error listing valid values for unknown names
  - Add PostgreSQL connection support alongside existing MySQL and SQLite
  - _Requirements: 3.1, 3.3, 7.1, 7.2, 7.3, 7.4, 7.5_

- [x] 1.2 (P) Add node-type selection to distinguish Anvil from Geth
  - Define `EthNodeType` constants: `anvil` and `geth`
  - Add `NodeType` field to the ETH config struct, defaulting to `anvil`
  - Route any node-specific behavior differences (e.g., key import method) through `NodeType` at the infrastructure layer so all use cases remain node-agnostic
  - _Requirements: 9.1, 9.2_

- [x] 2. Define ISP-compliant port interfaces for ETH blockchain operations
- [x] 2.1a (DONE PR #575) Add lifecycle, key access, transaction sign/send, raw key import, and node API interfaces
  - Added `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHTransactionSender`, `ETHRawKeyImporter`, `ETHNodeAPIClient`
  - Added composed `ETHKeygenSignClient`, `ETHWatchClient`
  - Added to `internal/application/ports/api/eth/interface.go` (not a separate `interfaces_small.go`)
  - `Ethereumer` preserved with DI-layer-only usage restriction comment
  - _Requirements: 5.2 (partial)_

- [x] 2.1b (P) Add EIP-1559 and monitoring port interfaces (required for Tasks 7, 9)
  - Define `BalanceChecker`, `TxCreator`, `GasEstimator`, `TxSigner`, `TxSender`, `TxMonitor`, `AddressValidator`, and `ChainConfigProvider` as separate interfaces in `interface.go`
  - Define composed interfaces `WatchTxCreationDeps` and `KeygenSignTxDeps` for use case convenience
  - Added `TransactionReceipt` domain type to `internal/domain/ethereum/types.go` for `TxMonitor`
  - `ETHTransactionSender` converted to type alias for `TxSender` to eliminate duplication
  - _Requirements: 5.2_

- [x] 2.2 (P) Add port-level DTOs to prevent infrastructure type leakage
  - `TxCreateParams` struct already present in ETH ports package (added in PR #579)
  - All new interfaces in 2.1b use only domain types — no infrastructure types in port signatures
  - _Requirements: 5.5_

- [x] 3. (P) Implement JSON transaction file format for air-gapped exchange
  - Define `ETHTransactionFile` struct with fields: version, tx_type, eth_tx_type, chain_id, nonce, from, to, value, gas, gas_price (legacy), max_fee_per_gas, max_priority_fee_per_gas, data, raw_tx_hex, signed_tx_hex
  - Use `tx.MarshalBinary()` (EIP-2718 envelope, not RLP) to correctly encode both legacy and EIP-1559 transactions
  - Implement write and read functions with field validation; return descriptive errors identifying missing or invalid fields
  - Follow BTC naming convention for file names: `{action}_{txID}_{type}_{signedCount}_{timestamp}.json`
  - Include `eth_tx_type` (0=legacy, 2=EIP-1559) so the offline signer can distinguish formats without parsing the raw hex
  - _Requirements: 2.5, 6.1, 6.2, 6.4, 6.5_

- [x] 4. Update the Ethereum client for EIP-1559 fee estimation and modern transaction encoding
- [x] 4.1 Add dynamic fee estimation methods to the Ethereum client
  - Add `SuggestGasTipCap` wrapper method calling `ethClient.SuggestGasTipCap(ctx)` with fallback to the config default when the RPC call fails
  - Add `SupportsEIP1559` detection method that inspects the connected node's block header for base fee support
  - Compute `maxFeePerGas` from the latest block base fee using `baseFee × 2 + tip`
  - _Requirements: 2.3, 3.5_

- [x] 4.2 Fix transaction encoding to use EIP-2718 binary format
  - Replace `rlp.EncodeToBytes`/`rlp.Decode` with `tx.MarshalBinary()`/`tx.UnmarshalBinary()` in the ETH transaction encoding layer
  - Verify encoding handles both Type 0 (legacy) and Type 2 (EIP-1559) transactions correctly
  - _Requirements: 2.1, 2.2_

- [x] 4.3 Add offline-capable signing and update signer selection
  - Add `SignOnRawTransaction` variant that accepts `*ecdsa.PrivateKey` directly (no keystore password required)
  - Replace `types.NewLondonSigner(chainID)` with `types.LatestSignerForChainID(chainID)` for forward compatibility with future transaction types
  - Depends on 4.2 (encoding fix must be in place)
  - _Requirements: 2.4_

- [x] 5. Verify and strengthen HD key generation in the Keygen wallet
- [x] 5.1 (P) Ensure BIP-39/BIP-44 key generation is robust and persists accountXpriv
  - Verify BIP-44 derivation path `m/44'/60'/0'/0/x` is used consistently for all ETH key generation
  - Verify Ethereum addresses are derived from secp256k1 public keys via Keccak-256 (last 20 bytes)
  - Store the account-level extended private key (`accountXpriv`) in the database for later offline signing
  - Return an error without creating any partial key material if mnemonic generation fails due to insufficient entropy
  - Ensure no raw private key material is stored or logged outside the offline keygen environment
  - _Requirements: 3.2, 4.1, 4.2, 4.4, 4.5, 4.6_

- [x] 5.2 (P) Implement accountXpub export for the Watch wallet
  - Export the account-level extended public key (`accountXpub`) to a file so the Watch wallet can derive and verify child addresses without ever holding private keys
  - Follow the BTC `ExportFullPubkey` use case pattern
  - _Requirements: 4.3_

- [ ] 6. Implement Keygen offline transaction signing use case
- [x] 6.1 Implement offline signing logic for ETH transactions
  - Read unsigned transaction JSON file using the transaction file repo (Task 3)
  - Retrieve `accountXpriv` from the database and derive the child private key at the BIP-44 index matching the transaction's target account
  - Return a clear error identifying the missing key index when the required key is not available
  - Sign using `types.SignTx` with `LatestSignerForChainID(chainID)`; detect transaction type from `tx.Type()` to select the correct signing approach
  - Write the signed transaction JSON file; leave the original unsigned file unchanged
  - Make no network or RPC calls during signing (fully offline operation)
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.4_

- [x] 6.2 Wire DI container for ETH Keygen signing
  - Replace the current "not implemented yet" stub in the DI container with a fully instantiated ETH Keygen signing use case
  - Inject all required dependencies: account key repository, transaction file repository, chain config provider
  - _Requirements: 1.5_

- [x] 7. (P) Update Watch wallet transaction creation for EIP-1559
  - Check EIP-1559 support via `SupportsEIP1559`; fall back to constructing a legacy `LegacyTx` when the node does not support it
  - Estimate `maxPriorityFeePerGas` using `SuggestGasTipCap` with config fallback; compute `maxFeePerGas = baseFee × 2 + tip`
  - Construct an EIP-1559 `DynamicFeeTx` with both fee fields when supported
  - Serialize the unsigned transaction to a JSON file using the transaction file repo
  - Save a transaction record in the database with status `TxTypeUnsigned`
  - Depends on: Tasks 1, 2, 3, 4
  - _Requirements: 2.1, 2.2, 2.3, 6.1_

- [ ] 8. (P) Update Watch wallet transaction broadcast to read from signed files
  - Read the signed transaction JSON file using the transaction file repo
  - Deserialize and validate the signature before broadcasting
  - Submit via `eth_sendRawTransaction` and persist the returned txHash to the database as `TxTypeSent`
  - Depends on: Tasks 2, 3, 4
  - _Requirements: 6.3_

- [ ] 9. Fix Watch wallet transaction monitoring
- [ ] 9.1 (P) Implement confirmation counting and success status transition
  - Query `eth_getTransactionReceipt` and `eth_blockNumber` to compute the current confirmation count
  - Transition transaction status from `TxTypeSent` to `TxTypeDone` once confirmations reach the configured threshold
  - Update `is_allocated` in `account_pubkey_table` after a confirmed successful send to prevent double-spending
  - _Requirements: 8.1, 8.2, 8.4_

- [ ] 9.2 (P) Detect failed transactions and handle node connectivity failures
  - Detect transactions that fail or revert on-chain; update status to a failure state with the revert reason if available
  - Implement exponential backoff retry when the Ethereum node becomes unreachable; log the connectivity issue at WARN level including the retry count
  - _Requirements: 8.3, 8.5_

- [ ] 10. Migrate all ETH use cases to Clean Architecture compliance
- [~] 10.1 Remove direct infrastructure imports from ETH use cases (Partially done — PR #575)
  - [x] Removed `ethereum` infrastructure package imports from all use cases (PR #575)
  - [ ] Remove remaining `apiethimpl` infrastructure imports (3 files still import for `Password` constant):
    - `internal/application/usecase/keygen/eth/sign_transaction.go`
    - `internal/application/usecase/sign/eth/sign_transaction.go`
    - `internal/application/usecase/keygen/eth/import_private_key.go`
  - Completion blocked by Task 10.3 (configurable password injection)
  - _Requirements: 5.1_

- [~] 10.2 (P) Restrict/deprecate monolithic Ethereumer interface (Partially done — PR #575)
  - [x] `Ethereumer` in `internal/application/ports/api/eth/interface.go` now carries a DI-layer-only usage restriction comment (PR #575)
  - [x] The deprecated `internal/infrastructure/api/eth/api-interface.go` has been deleted (PR #575)
  - [ ] Final removal of `Ethereumer` from `interface.go` — evaluate after Tasks 6.2, 7, 8, 9 complete DI wiring migration
  - Depends on 10.1
  - _Requirements: 5.3_

- [ ] 10.3 (P) Replace hardcoded password constant with configurable injection
  - Remove the hardcoded `Password = "password"` constant from the ETH infrastructure package
  - Inject the keystore password via the configuration system using the `KeystorePassword` field added in Task 1.1
  - Depends on 10.1
  - _Requirements: 5.4_

- [ ] 11. Docker Compose node profiles and E2E verification scripts
- [ ] 11.1 (P) Add Docker Compose service profiles for Anvil and Geth with multi-DB support
  - Add `anvil` profile service using `ghcr.io/foundry-rs/foundry:latest` configured with chain ID, pre-funded accounts, and block time
  - Add `geth` profile service using `ethereum/client-go:stable` as an alternative testing node
  - Support multiple database backends (PostgreSQL, MySQL, SQLite) consistent with BTC/BCH compose structure
  - Ensure both node profiles share the same Watch/Keygen wallet service definitions
  - _Requirements: 3.4, 9.3, 9.4_

- [ ] 11.2 Create E2E shell scripts for the full ETH transaction flow
  - Mirror the BTC/BCH E2E script structure at `scripts/operation/btc/e2e/`
  - Cover the complete single-sig flow: keygen seed/key → accountXpub export → watch address import → create unsigned tx → keygen sign → watch send → monitor confirmation
  - Accept `NODE_TYPE=anvil|geth` as an environment variable for node selection
  - Provide individual scripts per operation pattern (consistent with BTC patterns 1–11)
  - Depends on: Tasks 5, 6, 7, 8, 9
  - _Requirements: 9.5, 9.6_

- [ ] 12. Unit and integration test coverage
- [ ]* 12.1 Unit tests for Keygen offline signing use case
  - Test that signing an unsigned tx file produces a correctly signed output with a valid signature
  - Test that requesting a missing key index returns a clear, descriptive error
  - _Requirements: 1.1, 1.2, 1.3_

- [ ]* 12.2 Unit tests for JSON transaction file format
  - Test write/read roundtrip for both legacy (Type 0) and EIP-1559 (Type 2) transaction files
  - Test that malformed files or missing required fields produce descriptive parse errors
  - _Requirements: 6.4, 6.5_

- [ ]* 12.3 Integration tests for Watch → Keygen → Watch file exchange flow
  - Create unsigned transaction file → sign offline → broadcast signed file → verify `TxTypeDone` status transition
  - Test EIP-1559 transaction creation against an Anvil node
  - Test legacy fallback path when `SupportsEIP1559` returns false
  - _Requirements: 1.1, 2.1, 2.2, 6.1, 6.2, 6.3_
