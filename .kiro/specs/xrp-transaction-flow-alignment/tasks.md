# Implementation Plan

## Overview

This implementation plan translates the XRP Transaction Flow Alignment design into executable tasks. Tasks follow a logical progression from infrastructure setup through use case refactoring to comprehensive testing. Parallel-capable tasks are marked with `(P)` to enable concurrent development.

## Task Breakdown

### 1. Infrastructure and Dependencies Setup

- [x] 1.1 Add Peersyst xrpl-go dependency
  - Add `github.com/Peersyst/xrpl-go` to go.mod with appropriate version
  - Verify compatibility with existing xrpscan/xrpl-go dependency (no conflicts)
  - Run `go mod tidy` and verify build passes
  - _Requirements: 3.2, 5.1_

- [x] 1.2 Create segregated port interfaces
  - Define `AccountInfoProvider` interface in `application/ports/api/xrp/account_info.go`
  - Define `TransactionSigner` interface in `application/ports/api/xrp/transaction_signer.go`
  - Define `TransactionSubmitter` interface in `application/ports/api/xrp/transaction_submitter.go`
  - Document interface contracts with method signatures and preconditions
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

### 2. JSON Transaction File Format Implementation

- [x] 2.1 (P) Create XRP transaction file DTOs with validation
  - Define `XRPTransactionFile` struct in `application/dto/xrp/transaction_file.go`
  - Define `XRPTransactionEntry` struct with all required fields (UUID, unsigned data, signature count, etc.)
  - Add JSON struct tags following Google JSON Style Guide (lowercase field names)
  - Implement `Validate()` method with semantic versioning pattern check (Major.Minor.Patch)
  - Validate chain equals "XRP" and network is "mainnet" or "testnet"
  - Enforce signature count invariants (0 <= signatureCount <= requiredSignatures)
  - Verify signedBlob is null when signatureCount is 0
  - _Requirements: 1.2, 1.4, 1.5, 2.4_

- [x] 2.2 Extend transaction file repository with JSON methods
  - Implement `ReadJSONTransactionFile()` in TransactionFileRepository
  - Implement `WriteJSONTransactionFile()` in TransactionFileRepository
  - Reuse existing file path naming convention (`{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.json`)
  - Add JSON validation (version field, chain/network consistency, non-empty transactions)
  - Enforce `.json` file extension
  - _Requirements: 1.1, 1.3, 2.1, 3.1, 4.1_

### 3. Native Go Transaction Signing Implementation

- [x] 3.1 Implement PeersystSigner with offline signing capability
  - Create `PeersystSigner` struct in `infrastructure/api/xrp/peersyst_signer.go`
  - Implement `SignTransactionNative()` method using Peersyst wallet.Sign()
  - Handle wallet derivation from seed using Peersyst wallet.FromSeed()
  - Convert `dtoxrp.TxInput` to Peersyst transaction format
  - Return hex-encoded transaction hash and signed blob
  - Ensure no network calls during signing operations (offline capability)
  - Implement deterministic signing (same input + secret = same signature)
  - _Requirements: 3.2, 5.1, 5.2, 5.3_

- [x] 3.2 Add multi-signature signing support
  - Detect multi-sig transactions (empty SigningPubKey or Signers array present)
  - Implement multi-sig flow using Peersyst wallet.Multisign()
  - Combine multiple signatures into single signed blob
  - Validate multi-sig fee calculation ((N+1) × base fee)
  - _Requirements: 6.1, 6.2, 6.3_

### 4. Create Transaction Use Case Refactoring

- [ ] 4.1 Refactor CreateTransactionUseCase for JSON output with error handling
  - Replace dependency from full `XRPer` interface to `AccountInfoProvider`
  - Generate JSON transaction file using `WriteJSONTransactionFile()`
  - Initialize transactions with signatureCount = 0 and signedBlob = null
  - Set requiredSignatures based on account configuration (multi-sig threshold)
  - Include comprehensive metadata (version "1.0.0", chain "XRP", network, timestamp)
  - Wrap errors with operation context using `fmt.Errorf("context: %w", err)`
  - Validate account has sufficient balance before transaction creation
  - Validate destination address format
  - Return descriptive errors for account info query failures
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 10.1, 10.5_

### 5. Sign Transaction Use Case Refactoring

- [ ] 5.1 Refactor SignTransactionUseCase for JSON parsing and native signing
  - Replace CSV file parsing with `ReadJSONTransactionFile()`
  - Replace gRPC signing dependency with `TransactionSigner` interface
  - Retrieve signing secrets from XRPAccountKeyRepositorier
  - Track signature count and increment after each signature
  - Wrap signing errors with transaction UUID context
  - Never log secret values (security requirement)
  - Return descriptive errors for secret retrieval failures
  - Log signature count and completion status at DEBUG level
  - _Requirements: 3.1, 3.2, 3.6, 10.1, 10.2, 10.4, 10.5, 10.6_

- [ ] 5.2 Implement signature completion logic
  - Calculate completion status (signatureCount >= requiredSignatures)
  - Set `isComplete` flag when quorum is met
  - Skip re-signing if transaction already complete
  - Generate output file with updated signature count in filename
  - _Requirements: 3.3, 3.4, 3.7, 6.3, 6.4_

- [ ] 5.3 Add multi-signature workflow support
  - Support sequential signing by multiple operators
  - Preserve existing signatures when adding new signature
  - Update signed blob with combined signatures
  - Generate sequential file names (unsigned → signed-1 → signed-2 → signed-final)
  - _Requirements: 6.2, 6.5_

### 6. Send Transaction Use Case Refactoring

- [ ] 6.1 Refactor SendTransactionUseCase for JSON validation with error handling
  - Replace CSV parsing with `ReadJSONTransactionFile()`
  - Validate all transactions have `isComplete == true` before submission
  - Verify `signedBlob` is non-null and hex-encoded
  - Extract signed blobs for ledger submission
  - Return descriptive errors for incomplete transactions (insufficient signatures)
  - Propagate XRP Ledger error codes and messages (tefPAST_SEQ, tecUNFUNDED_PAYMENT, etc.)
  - Log file path and error details at ERROR level for parsing failures
  - Handle batch transaction submission (collect successes and errors)
  - _Requirements: 4.1, 4.3, 4.5, 4.6, 10.1, 10.3_

- [ ] 6.2 Implement transaction submission via TransactionSubmitter
  - Replace direct xrpl-go calls with `TransactionSubmitter` interface dependency
  - Submit signed transaction blobs to XRP Ledger
  - Return transaction hash and ledger version on success
  - Update database transaction status (sent → done)
  - _Requirements: 4.2, 4.4_

### 7. Backward Compatibility Implementation

- [ ] 7.1 Implement legacy text format detection
  - Create `parseTransactionFile` helper that detects format (JSON vs text) based on extension
  - Fall back to text parsing when JSON parsing fails
  - Log deprecation warnings when processing text format files
  - Ensure JSON format always used for new transactions
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

### 8. Unit Testing

- [ ] 8.1 (P) Test JSON transaction file serialization
  - Verify JSON marshaling/unmarshaling correctness for XRPTransactionFile
  - Test schema validation (version, chain, network fields)
  - Validate signature count invariants enforcement
  - Test file path naming convention
  - Test error handling for invalid JSON and missing files
  - _Requirements: 8.1_

- [ ] 8.2 (P) Test PeersystSigner implementation
  - Verify wallet derivation from seed produces correct addresses
  - Test transaction signing output format (hex-encoded blob and hash)
  - Validate multi-signature signing flow produces correct combined signatures
  - Test error handling for invalid seed and malformed transactions
  - Verify offline signing capability (no network calls during signing)
  - Test deterministic signing (same input produces same signature)
  - _Requirements: 8.2, 5.1, 5.2, 5.3_

- [ ] 8.3 (P) Test CreateTransactionUseCase
  - Verify JSON file generation with correct metadata structure
  - Test signature count initialized to 0
  - Verify required signatures set based on account configuration
  - Test error wrapping for account info query failures
  - Test balance validation logic
  - Test destination address format validation
  - _Requirements: 8.5_

- [ ] 8.4 (P) Test SignTransactionUseCase
  - Verify signature count tracking and incrementation
  - Test completion status detection (signatureCount >= requiredSignatures)
  - Verify skip re-signing logic for completed transactions
  - Test error handling for missing secrets
  - Verify secret values never logged
  - Test error context wrapping with transaction UUID
  - _Requirements: 8.5_

- [ ] 8.5 (P) Test SendTransactionUseCase
  - Verify completion validation before submission
  - Test error handling for incomplete transactions
  - Verify signed blob extraction and submission
  - Test batch transaction handling (successes and errors)
  - Test XRP Ledger error code propagation
  - _Requirements: 8.5_

### 9. Integration Testing

- [ ] 9.1 Test full transaction flow (create → sign → send)
  - Create unsigned transaction → verify JSON file structure
  - Sign transaction → verify signature count incremented and blob populated
  - Send transaction → verify ledger submission (testnet)
  - Validate database updates at each phase
  - _Requirements: 8.3_

- [ ] 9.2 Test multi-signature scenarios
  - Create transaction requiring 2-of-3 signatures
  - Sign with Keygen wallet (first signature)
  - Sign with Sign wallet (second signature)
  - Verify completion status set correctly
  - Submit to testnet and confirm transaction
  - Test 3-of-5 multisig scenario for comprehensive validation
  - _Requirements: 8.4, 6.1, 6.3_

- [ ] 9.3 Test Peersyst + xrpscan compatibility
  - Sign transaction with Peersyst wallet.Sign()
  - Submit signed blob using xrpscan/xrpl-go SubmitTransaction()
  - Verify XRP Ledger accepts the signed blob
  - Confirm transaction appears on ledger explorer
  - _Requirements: 3.2_

- [ ] 9.4 (P) Test backward compatibility
  - Parse legacy text format transaction files
  - Verify deprecation warnings logged
  - Ensure text format parsing still functional
  - Confirm new transactions always use JSON format
  - _Requirements: 7.1, 7.2, 7.4, 8.6_

### 10. System Integration and Validation

- [ ] 10.1 Wire segregated interfaces to implementations
  - Connect `AccountInfoProvider` interface to xrpscan/xrpl-go client implementation
  - Connect `TransactionSigner` interface to PeersystSigner implementation
  - Connect `TransactionSubmitter` interface to xrpscan/xrpl-go client implementation
  - Update dependency injection configuration
  - _Requirements: 9.5_

- [ ] 10.2 Update use case initialization in DI container
  - Inject `AccountInfoProvider` into CreateTransactionUseCase (replace full XRPer)
  - Inject `TransactionSigner` into SignTransactionUseCase
  - Inject `TransactionSubmitter` into SendTransactionUseCase
  - Verify all interface contracts satisfied
  - _Requirements: 2.2, 3.2, 4.2_

- [ ] 10.3 Run end-to-end validation
  - Execute deposit flow (collect client funds) on testnet
  - Execute payment flow (2-of-2 multi-sig) on testnet
  - Execute transfer flow (internal account transfer) on testnet
  - Verify all transactions confirmed on XRP Ledger
  - _Requirements: 8.3_

## Requirements Coverage Matrix

| Requirement | Tasks |
|-------------|-------|
| 1.1 | 2.2 |
| 1.2 | 2.1 |
| 1.3 | 2.2 |
| 1.4 | 2.1 |
| 1.5 | 2.1 |
| 2.1 | 2.2, 4.1 |
| 2.2 | 4.1, 10.2 |
| 2.3 | 4.1 |
| 2.4 | 2.1, 4.1 |
| 2.5 | 4.1 |
| 3.1 | 2.2, 5.1 |
| 3.2 | 1.1, 3.1, 5.1, 9.3, 10.2 |
| 3.3 | 5.2 |
| 3.4 | 5.2 |
| 3.5 | 2.2 |
| 3.6 | 5.1 |
| 3.7 | 5.2 |
| 4.1 | 2.2, 6.1 |
| 4.2 | 6.2, 10.2 |
| 4.3 | 6.1 |
| 4.4 | 6.2 |
| 4.5 | 6.1 |
| 4.6 | 6.1 |
| 5.1 | 3.1, 8.2 |
| 5.2 | 3.1, 8.2 |
| 5.3 | 3.1, 8.2 |
| 6.1 | 9.2 |
| 6.2 | 3.2, 5.3 |
| 6.3 | 3.2, 5.2, 9.2 |
| 6.4 | 5.2 |
| 6.5 | 5.3 |
| 7.1 | 7.1, 9.4 |
| 7.2 | 7.1, 9.4 |
| 7.3 | 7.1 |
| 7.4 | 7.1, 9.4 |
| 7.5 | 7.1 |
| 8.1 | 8.1 |
| 8.2 | 8.2 |
| 8.3 | 9.1, 10.3 |
| 8.4 | 9.2 |
| 8.5 | 8.3, 8.4, 8.5 |
| 8.6 | 9.4 |
| 9.1 | 1.2 |
| 9.2 | 1.2 |
| 9.3 | 1.2 |
| 9.4 | 1.2 |
| 9.5 | 10.1 |
| 10.1 | 5.1, 6.1 |
| 10.2 | 5.1 |
| 10.3 | 6.1 |
| 10.4 | 5.1 |
| 10.5 | 4.1, 5.1 |
| 10.6 | 5.1 |

## Task Dependencies

### Critical Path
1. Infrastructure (1.1, 1.2) → All subsequent tasks
2. JSON Format (2.1) → File Repository (2.2) → Use Cases (4.1, 5.1, 6.1)
3. Native Signing (3.1) → Multi-sig (3.2) → Use Case Integration (5.1, 5.2)
4. Use Case Refactoring (4.1, 5.1, 6.1) → Integration (10.1, 10.2) → Validation (10.3)

### Parallel Execution Groups
- **Group A** (after 1.1, 1.2): 2.1
- **Group B** (after 2.2, 3.2): All use case refactoring (4.1, 5.1, 6.1)
- **Group C** (after 3.1, 5.1, 6.1): All unit tests (8.1-8.5)
- **Group D** (after 9.1): 9.4

## Key Improvements from Review

### Integration of Related Work
1. **Merged validation into DTO creation** (2.3 → 2.1): Validation is intrinsic to DTO implementation
2. **Removed separate offline testing task** (3.3): Offline capability verified in unit tests (8.2)
3. **Integrated error handling** (4.2, 5.4, 6.3): Error handling is part of correct implementation, not optional

### Better Requirements Mapping
1. **Offline signing requirements (5.1-5.3)**: Now correctly map to implementation (3.1) and testing (8.2)
2. **Error handling requirements (10.1-10.6)**: Integrated into respective use case tasks

### More Cohesive Tasks
- Each task now represents a complete, integrated piece of functionality
- Error handling and validation are built into implementation tasks
- Testing tasks focus on verification, not implementation

## Estimated Effort

- **Total Tasks**: 10 major tasks, 28 sub-tasks (reduced from 33)
- **Average Sub-task**: 1.5-2.5 hours
- **Total Effort**: ~45-70 hours
- **Critical Path**: ~22-30 hours (with parallelization)
- **Testing Phase**: ~18 hours (can overlap with implementation)

## Next Steps

1. Review refined task breakdown and confirm scope alignment
2. Approve tasks in `.kiro/specs/xrp-transaction-flow-alignment/spec.json`
3. Execute tasks using `/kiro:spec-impl xrp-transaction-flow-alignment <task-id>`
4. Run verification commands after each task (see go-development skill)
