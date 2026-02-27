# Requirements Document

## Introduction

This specification defines the requirements for modernizing the Ethereum (ETH) transaction flow in the go-crypto-wallet project.
The current ETH implementation has significant gaps compared to the mature BTC implementation:
the Sign wallet is entirely non-functional (all methods are stubs),
multiple use cases violate Clean Architecture by importing directly from the infrastructure layer,
and the default transaction type is legacy (pre-EIP-1559).
This feature aligns the ETH flow with the BTC reference implementation by introducing offline signing capability,
Foundry/Anvil as the development node, robust HD key generation, and EIP-1559 transaction support.

**Scope Note:** Multisig support (Safe/Gnosis Safe smart contract wallets) is explicitly excluded from this specification and will be designed in a separate spec.

**Coin Type Scope:** This specification targets the `eth` coin type only. ERC-20 token support is out of scope and will be handled separately.

### Current State Summary

| Capability | BTC | ETH (Current) |
|---|---|---|
| Sign wallet | Fully implemented (PSBT multisig) | All stubs — non-functional |
| DI Sign wiring | Fully wired | Returns "not implemented yet" |
| Offline signing | PSBT file-based | None (relies on online Geth keystore) |
| Transaction type | N/A | Legacy only (EIP-1559 code exists but unused) |
| Port compliance | ISP-compliant small interfaces | ISP interfaces added (PR #575) — `Ethereumer` retained for DI layer only |
| Architecture | Clean (ports used correctly) | Partial: `ethereum` infra import removed; `apiethimpl.Password` import remains in 3 use cases |
| Dev node | Bitcoin Core (regtest) | Geth (config references deprecated Goerli) |

## Requirements

### Requirement 1: Offline Transaction Signing (Keygen Wallet)

**Objective:** As a wallet operator, I want the ETH Keygen wallet to function as a fully operational offline signing component, so that private keys never need to be exposed on an online machine.

> **Note:** ETH uses single-sig EOA only. The Sign Wallet is not required. The Keygen Wallet performs all transaction signing, consistent with the BTC keygen signing pattern. See [docs/chains/eth/architecture.md](../../../../docs/chains/eth/architecture.md).

#### Acceptance Criteria

1. When the Keygen wallet receives an unsigned ETH transaction file, the ETH Keygen Wallet shall deserialize the transaction, sign it with the appropriate private key, and produce a signed transaction file.
2. The ETH Keygen Wallet shall derive child private keys from the stored `accountXpriv` using BIP-44 derivation path (`m/44'/60'/0'/0/x`) matching the transaction's target account index.
3. If a private key for the requested account index is not available, the ETH Keygen Wallet shall return an error identifying the missing key index.
4. The ETH Keygen Wallet shall operate without any network connectivity or RPC connection to an Ethereum node.
5. When the DI container constructs a Keygen wallet for ETH coin type, the DI Container shall instantiate a fully functional ETH keygen signer instead of returning "not implemented yet".

### Requirement 2: EIP-1559 Transaction Support

**Objective:** As a wallet operator, I want all new ETH transactions to use EIP-1559 (Type 2) dynamic fee transactions by default, so that fee estimation is more predictable and gas costs are optimized.

#### Acceptance Criteria

1. When creating a new ETH transaction, the Watch Wallet shall construct an EIP-1559 `DynamicFeeTx` with `maxFeePerGas` and `maxPriorityFeePerGas` fields.
2. When the connected Ethereum node does not support EIP-1559, the Watch Wallet shall fall back to creating a legacy `LegacyTx` with `gasPrice`.
3. The ETH Transaction Creator shall estimate `maxPriorityFeePerGas` via `eth_maxPriorityFeePerGas` RPC and compute `maxFeePerGas` from the latest base fee.
4. When signing an ETH transaction, the ETH Keygen Signer shall use `types.LatestSignerForChainID(chainID)` to produce the correct signature, supporting both legacy and EIP-1559 transaction types.
5. The ETH Transaction File Format shall include the transaction type indicator so the Sign wallet can distinguish between legacy and EIP-1559 transactions during offline signing.

### Requirement 3: Foundry/Anvil Development Environment

**Objective:** As a developer, I want to use Foundry's Anvil as the local Ethereum development node, so that I have a fast, modern, and actively maintained testing environment replacing deprecated alternatives.

#### Acceptance Criteria

1. The ETH Configuration shall support an Anvil node endpoint (default: `http://localhost:8545`) as the JSON-RPC target.
2. When importing private keys into a local development node, the Keygen Wallet shall use Anvil-compatible key import methods (keystore file import, not `personal_importRawKey` RPC).
3. The ETH Configuration shall replace deprecated network references (Goerli) with current networks (Sepolia for testnet, local Anvil for development).
4. The ETH Development Setup shall provide Anvil startup configuration (chain ID, pre-funded accounts, block time) in project documentation or scripts.
5. While running against an Anvil node, the ETH Client shall function identically to running against a production Geth/Erigon node via standard `eth_*` JSON-RPC.

### Requirement 4: Robust Key Generation

**Objective:** As a wallet operator, I want ETH key generation to be robust, secure, and aligned with BTC's HD wallet approach, so that key management is consistent across chains.

#### Acceptance Criteria

1. The ETH Keygen Wallet shall generate HD wallet keys using BIP-39 mnemonic and BIP-44 derivation path (`m/44'/60'/0'/0/x`).
2. The ETH Keygen Wallet shall store the account-level extended private key (`accountXpriv`) in the database for later child key derivation by the Sign wallet.
3. The ETH Keygen Wallet shall export the account-level extended public key (`accountXpub`) to a file so the Watch Wallet can derive and verify child addresses independently without holding private keys.
4. The ETH Key Strategy shall derive Ethereum addresses from secp256k1 public keys via Keccak-256 hashing (last 20 bytes).
5. The ETH Keygen Wallet shall never store or transmit raw private keys in plaintext outside of the offline keygen environment.
6. If the HD seed or mnemonic generation fails due to insufficient entropy, the ETH Keygen Wallet shall return an error without creating partial key material.

### Requirement 5: Clean Architecture Compliance

**Objective:** As a developer, I want all ETH use cases to follow Clean Architecture by depending only on port interfaces (not infrastructure implementations), so that the codebase is maintainable and testable.

#### Acceptance Criteria

1. The ETH Use Cases shall import only from `internal/application/ports/api/eth` and never directly from `internal/infrastructure/api/eth` or `internal/infrastructure/api/eth/eth`. (Partial — PR #575 removed `ethereum` infrastructure imports; `apiethimpl.Password` constant import remains in 3 use cases, tracked by criterion 5.4.)
2. The ETH Port Interface shall be decomposed into small, ISP-compliant interfaces (e.g., `BalanceChecker`, `TxCreator`, `TxSigner`, `TxSender`, `GasEstimator`) matching the BTC port pattern. (Partial — PR #575 added `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHTransactionSender`, `ETHRawKeyImporter`, `ETHNodeAPIClient`, `ETHKeygenSignClient`, `ETHWatchClient`; EIP-1559 and monitoring interfaces still needed for Tasks 7, 9.)
3. ~~The deprecated interface definitions in `internal/infrastructure/api/eth/api-interface.go` shall be removed after all use cases are migrated to port interfaces.~~ **Done (PR #575)**: `internal/infrastructure/api/eth/api-interface.go` has been deleted.
4. The hardcoded password constant (`Password = "password"`) shall be replaced with a configurable secret injected via the configuration system.
5. When a use case requires an infrastructure-specific type, the ETH Ports Layer shall define a corresponding DTO or port-level type to avoid infrastructure leakage.

### Requirement 6: Transaction File Format and Flow Alignment

**Objective:** As a wallet operator, I want the ETH transaction flow (Watch → Keygen → Watch) to follow the same file-based exchange pattern as BTC single-sig, so that the air-gapped signing workflow is consistent across chains.

#### Acceptance Criteria

1. When creating an unsigned transaction, the ETH Watch Wallet shall serialize the transaction to a structured file format that includes all data needed for offline signing (nonce, chain ID, gas parameters, to, value, data, transaction type).
2. When the Sign Wallet signs a transaction file, the ETH Sign Wallet shall produce a signed file ready for broadcast.
3. When broadcasting a signed transaction, the ETH Watch Wallet shall deserialize the signed file and submit via `eth_sendRawTransaction` RPC.
4. The ETH Transaction File Format shall use JSON encoding with clearly labeled fields for human readability and debugging.
5. If a transaction file is malformed or missing required fields, the ETH Wallet Component shall return a descriptive parse error identifying the missing or invalid fields.

### Requirement 7: Configuration and Network Modernization

**Objective:** As a wallet operator, I want the ETH configuration to support modern networks and PostgreSQL, so that the deployment environment is current and consistent with the BTC configuration.

#### Acceptance Criteria

1. The ETH Configuration shall support PostgreSQL as a database backend in addition to MySQL and SQLite.
2. The ETH Configuration shall define network types for mainnet, Sepolia (testnet), Holesky (testnet), and local (Anvil).
3. The ETH Configuration shall include EIP-1559 fee parameters (max fee cap, priority fee floor) with sensible defaults.
4. The ETH Configuration shall support chain ID specification per network for EIP-155 replay protection.
5. When an unsupported or deprecated network name is specified, the ETH Configuration Loader shall return an error listing the supported networks.

### Requirement 8: Transaction Monitoring and Status Tracking

**Objective:** As a wallet operator, I want ETH transaction monitoring to be fully functional (not partially commented out), so that I can track transaction confirmations and completion status.

#### Acceptance Criteria

1. When a signed transaction is broadcast, the ETH Monitor shall track its confirmation count against the configured threshold.
2. When a transaction reaches the required confirmation count, the ETH Monitor shall update the transaction status to `TxTypeDone` in the database.
3. When a transaction is detected as failed or reverted on-chain, the ETH Monitor shall update the status to a failure state with the revert reason if available.
4. The ETH Monitor shall update `is_allocated` in `account_pubkey_table` after a successful send to prevent double-spending from the same account.
5. If the Ethereum node becomes unreachable during monitoring, the ETH Monitor shall retry with exponential backoff and log the connectivity issue.

### Requirement 9: Ethereum Node Flexibility, Docker Compose, and E2E Verification

**Objective:** As a developer, I want to be able to switch between Foundry Anvil and go-ethereum (Geth) via a single configuration value, and to verify the full ETH transaction flow using an E2E script, so that the implementation is validated against both node implementations.

#### Acceptance Criteria

1. The ETH Configuration shall support both Foundry Anvil (`https://www.getfoundry.sh/anvil`) and go-ethereum Geth (`https://github.com/ethereum/go-ethereum`) as Ethereum node backends, selectable via a `node_type` config field (e.g., `anvil` or `geth`).
2. Any node-specific behavior differences (e.g., key import method, supported RPC methods) shall be handled internally so that use cases remain node-agnostic.
3. The Docker Compose configuration shall define separate service profiles for Anvil and Geth so that either can be launched independently (e.g., `docker compose --profile anvil up` / `docker compose --profile geth up`).
4. The Docker Compose configuration shall support multiple database backends (PostgreSQL, MySQL, SQLite) consistent with the BTC/BCH compose structure.
5. An E2E shell script shall be provided at `scripts/operation/eth/e2e/` covering the full transaction flow: key generation → address export/import → transaction creation → offline signing → broadcast → confirmation monitoring.
6. The E2E script shall be runnable against both Anvil and Geth by passing the node type as a parameter or reading from the environment.
