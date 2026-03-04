# Research & Design Decisions

---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design for XRP transaction flow alignment.

**Usage**:

- Log research activities and outcomes during the discovery phase.
- Document design decision trade-offs that are too detailed for `design.md`.
- Provide references and evidence for future audits or reuse.

---

## Summary

- **Feature**: `xrp-transaction-flow-alignment`
- **Discovery Scope**: Complex Integration (extending existing system with new transaction file format and native signing)
- **Key Findings**:
  - xrpl-go library (Peersyst implementation) provides native Go signing with `wallet.Sign()` and `wallet.Multisign()` functions
  - XRP multi-signature transactions require SignerListSet setup, empty SigningPubKey field, and (N+1) fee multiplier
  - JSON versioning best practices recommend semantic versioning with apiVersion field for backward compatibility
  - Existing file repository infrastructure (85% reusable) can be extended with JSON methods
  - Interface Segregation Principle requires three focused interfaces: AccountInfoProvider, TransactionSigner, TransactionSubmitter

## Research Log

### xrpl-go Native Signing Capabilities

- **Context**: Need to replace gRPC-based signing with native Go implementation for offline signing capability
- **Sources Consulted**:
  - [XRPLF/xrpl-go](https://github.com/XRPLF/xrpl-go) - Official XRP Ledger Foundation implementation
  - [Peersyst/xrpl-go](https://pkg.go.dev/github.com/Peersyst/xrpl-go) - Comprehensive wallet and signing functionality
  - [xrpscan/xrpl-go](https://github.com/xrpscan/xrpl-go) - Currently integrated in project (v0.2.11)
  - [Wallet package documentation](https://pkg.go.dev/github.com/Peersyst/xrpl-go/xrpl/wallet)

- **Findings**:
  - **Peersyst implementation** provides the most comprehensive wallet functionality:
    - `Sign(transaction)` - Signs a transaction offline, returns transaction blob and hash
    - `Multisign(transaction)` - Signs multisigned transactions offline
    - `Wallet` utility for deriving keypairs from seed, mnemonic, or entropy
  - **XRPLF implementation** supports serialization to signing transactions
  - **xrpscan implementation** (currently used) focuses on websocket API interactions
  - All libraries support native Go signing without external dependencies
  - Transaction signing returns 64-character hexadecimal transaction ID (hash)

- **Implications**:
  - **Recommended**: Integrate Peersyst/xrpl-go for signing functionality alongside existing xrpscan/xrpl-go
  - Offline signing fully supported without gRPC dependency
  - Multi-signature signing workflow natively available
  - May need to manage two xrpl-go implementations (xrpscan for submission, Peersyst for signing)

### XRP Multi-Signature Transaction Workflow

- **Context**: Need to understand XRP multi-sig requirements for JSON transaction file metadata design
- **Sources Consulted**:
  - [XRP Ledger Multi-Signing Concepts](https://xrpl.org/docs/concepts/accounts/multi-signing)
  - [Set Up Multi-Signing Tutorial](https://xrpl.org/docs/tutorials/how-tos/manage-account-settings/set-up-multi-signing)
  - [Send a Multi-Signed Transaction](https://xrpl.org/send-a-multi-signed-transaction.html)
  - [Secure Signing](https://xrpl.org/docs/concepts/transactions/secure-signing)

- **Findings**:
  - **SignerList Setup** (prerequisite):
    - Must submit a SignerListSet transaction first (single-signature)
    - Supports 1-32 addresses in signer list
    - Each signer has a weight; transactions require total weight ≥ quorum
  - **Multi-Signature Transaction Requirements**:
    - `SigningPubKey` field must be empty string
    - `Signers` field must contain array of signatures
    - Total signer weight must meet or exceed quorum
    - Fee must be at least (N+1) times normal transaction cost, where N = number of signatures
    - All transaction fields must be defined before collecting signatures
  - **Signing Process**:
    - Create transaction JSON with all fields populated
    - Each signer signs independently using `sign_for` method (not `sign`)
    - Collect signatures from all signers
    - Append signatures to `Signers` array
    - Submit when quorum met

- **Implications**:
  - JSON transaction file must track:
    - Current signature count vs required signatures (quorum)
    - Signer list configuration (if available)
    - SigningPubKey must be empty for multi-sig
  - File format must support sequential signing (unsigned → partial → complete)
  - Fee calculation must account for signature count
  - Transaction fields immutable after first signature (no modifications during signing rounds)

### JSON Transaction File Format Versioning

- **Context**: Need to design JSON file format that supports future evolution and backward compatibility
- **Sources Consulted**:
  - [safe-json: Automatic JSON format versioning](https://hackage.haskell.org/package/safe-json)
  - [JSON schema version attribute example](https://gist.github.com/mattyod/3608613)
  - [JSON-API Versioning Strategy](https://github.com/json-api/json-api/issues/406)
  - [REST API Standards & Best Practices 2026](https://www.boltic.io/blog/rest-api-standards)
  - [Google JSON Style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)

- **Findings**:
  - **Versioning Approaches**:
    - Semantic versioning (Major.Minor.Patch) most common for data formats
    - `apiVersion` or `version` field in root JSON object
    - Migration functions for converting between versions
  - **Best Practices**:
    - Include version in every file (not just filename)
    - Major version for breaking changes (field removal)
    - Minor version for backward-compatible additions
    - Patch version for documentation or validation fixes
  - **Backward Compatibility**:
    - Version field allows detection of format
    - Parser can apply migration logic based on version
    - Deprecation warnings for old versions

- **Implications**:
  - JSON transaction file should include `version` field (e.g., "1.0.0")
  - Design schema to minimize future breaking changes
  - Plan for version detection and migration in parsers
  - Consider optional fields for forward compatibility

### Existing File Repository Architecture Analysis

- **Context**: Determine if transaction file repository can be extended or needs replacement
- **Sources**: Internal codebase analysis (gap-analysis.md findings)
- **Findings**:
  - **Current Capabilities**:
    - File naming convention: `{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.{ext}`
    - PSBT support (`.psbt`, base64-encoded for Bitcoin)
    - Hex support (`.hex` for Bitcoin Cash)
    - Text slice operations (`ReadFileSlice`, `WriteFileSlice`)
    - Metadata extraction from file paths
  - **Reusability Assessment**:
    - 85% of infrastructure reusable (path naming, directory creation, validation)
    - Pattern established for multiple formats (PSBT, Hex, Text)
    - Interface already supports extension (new methods can be added)
  - **Missing**:
    - No JSON file methods (`ReadJSONTransactionFile`, `WriteJSONTransactionFile`)
    - No structured transaction metadata support

- **Implications**:
  - **Recommended Approach**: Extend existing `TransactionFileRepositorier` interface
  - Add JSON-specific methods following PSBT/Hex pattern
  - Reuse file path naming and validation logic
  - Minimal disruption to existing BTC/BCH/ETH implementations

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| **Clean Architecture (Current)** | Strict layer separation: Domain → Application (Use Cases + Ports) → Infrastructure → Interface Adapters | Clear boundaries, testable, dependency inversion | Requires discipline to maintain layer separation | **Selected** - Already established in project |
| **Interface Segregation** | Split monolithic XRPer interface into focused interfaces (AccountInfoProvider, TransactionSigner, TransactionSubmitter) | Minimal dependencies, offline-capable signing, clear contracts | Requires interface refactoring and DI updates | **Selected** - Aligns with ISP and offline requirements |
| **JSON-First File Format** | Replace CSV text format with structured JSON | Metadata support, validation, extensibility | Migration complexity, file size increase | **Selected** - Essential for multi-sig tracking |
| **Native Go Signing** | Use xrpl-go library directly instead of gRPC | Offline signing, no network dependency, simpler architecture | Need to integrate new library (Peersyst/xrpl-go) | **Selected** - Requirement for offline operations |

## Design Decisions

### Decision: Extend TransactionFileRepositorier vs Create XRPTransactionFileHandler

- **Context**: Need JSON transaction file support; gap analysis presented three options (extend, create new, hybrid)
- **Alternatives Considered**:
  1. **Option A**: Extend existing `TransactionFileRepositorier` with JSON methods
  2. **Option B**: Create new `XRPTransactionFileHandler` interface
  3. **Option C**: Hybrid (extend file repo + new interfaces for API segregation)

- **Selected Approach**: **Option C - Hybrid**
  - Extend `TransactionFileRepositorier` with `ReadJSONTransactionFile` and `WriteJSONTransactionFile` methods
  - Create new segregated API interfaces (`AccountInfoProvider`, `TransactionSigner`, `TransactionSubmitter`)
  - Refactor use cases to use JSON methods and segregated interfaces

- **Rationale**:
  - File repository extension follows established PSBT/Hex pattern (consistency)
  - 85% infrastructure reuse (file naming, path validation, directory creation)
  - Interface segregation is orthogonal concern (separate decision)
  - Minimal disruption to existing BTC/BCH/ETH code

- **Trade-offs**:
  - ✅ Reuses mature infrastructure
  - ✅ Consistent with existing multi-format support
  - ✅ Simple testing (extend existing test suite)
  - ⚠️ File repository handles 4 formats (text, PSBT, Hex, JSON) - acceptable complexity

- **Follow-up**:
  - Monitor file repository size; refactor if exceeds ~500 lines
  - Ensure JSON methods follow same error handling patterns as PSBT/Hex

### Decision: xrpl-go Library Selection for Native Signing

- **Context**: Need to replace gRPC-based signing with native Go implementation
- **Alternatives Considered**:
  1. **Peersyst/xrpl-go**: Comprehensive wallet with Sign/Multisign functions
  2. **XRPLF/xrpl-go**: Official implementation with serialization support
  3. **xrpscan/xrpl-go**: Currently used, websocket-focused
  4. **Keep gRPC**: Maintain existing implementation

- **Selected Approach**: **Peersyst/xrpl-go for signing + xrpscan/xrpl-go for submission**
  - Add `github.com/Peersyst/xrpl-go` dependency for wallet/signing functionality
  - Keep `github.com/xrpscan/xrpl-go` for existing submission logic (already integrated)
  - Implement `TransactionSigner` interface using Peersyst wallet functions

- **Rationale**:
  - Peersyst provides most mature wallet/signing API (`wallet.Sign`, `wallet.Multisign`)
  - Offline signing explicitly supported (no network dependency)
  - Multi-signature signing natively available
  - xrpscan client already working for submission (no need to replace)

- **Trade-offs**:
  - ✅ Native Go signing (no gRPC dependency)
  - ✅ Multi-signature support out-of-box
  - ✅ Offline operation guaranteed
  - ⚠️ Two xrpl-go libraries in go.mod (acceptable for specialized use)
  - ⚠️ Need to test Peersyst signing output compatibility with xrpscan submission

- **Follow-up**:
  - Verify Peersyst signed transaction blobs work with xrpscan submission
  - Investigate if XRPLF implementation could replace both (future optimization)

### Decision: JSON Transaction File Schema Design

- **Context**: Define JSON structure that supports current workflow and future evolution
- **Alternatives Considered**:
  1. **Flat Structure**: Single-level JSON with all fields
  2. **Nested Metadata + Transactions**: Separate envelope metadata from transaction array
  3. **PSBT-like Binary Format**: Binary encoding similar to Bitcoin PSBT

- **Selected Approach**: **Nested Metadata + Transactions Array**

```json
{
  "version": "1.0.0",
  "chain": "XRP",
  "network": "mainnet",
  "createdAt": "2026-02-14T12:34:56Z",
  "transactions": [
    {
      "uuid": "01933e82-b6c9-7890-a123-456789abcdef",
      "unsignedData": { /* dtoxrp.TxInput fields */ },
      "senderAccount": "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
      "senderAccountType": "client",
      "signatureCount": 0,
      "requiredSignatures": 2,
      "signedBlob": null,
      "isComplete": false
    }
  ]
}
```

- **Rationale**:
  - Clear separation of file metadata vs transaction data
  - Version field enables format evolution
  - Transaction array supports batch operations (matches existing pattern)
  - Signature tracking fields enable multi-sig workflow validation

- **Trade-offs**:
  - ✅ Self-documenting structure
  - ✅ JSON schema validation possible
  - ✅ Extensible (new fields can be added to envelope without breaking)
  - ✅ Human-readable for debugging
  - ⚠️ Larger file size than text format (acceptable for metadata value)

- **Follow-up**:
  - Define JSON schema for validation
  - Test file size impact with realistic transaction counts

### Decision: Interface Segregation for XRP API

- **Context**: Current use cases depend on monolithic `XRPer` interface; need minimal, focused interfaces for ISP compliance
- **Alternatives Considered**:
  1. **Keep Monolithic XRPer**: No changes to existing interface
  2. **Full Decomposition**: Break XRPer into 10+ small interfaces
  3. **Three Focused Interfaces**: AccountInfoProvider, TransactionSigner, TransactionSubmitter

- **Selected Approach**: **Three Focused Interfaces**

```go
// For create transaction use case (online operations)
type AccountInfoProvider interface {
    GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error)
    GetBalance(ctx context.Context, addr string) (float64, error)
}

// For sign transaction use case (offline operations)
type TransactionSigner interface {
    SignTransactionNative(ctx context.Context, txInput *dtoxrp.TxInput, secret string) (signedTxID, txBlob string, err error)
}

// For send transaction use case (online operations)
type TransactionSubmitter interface {
    SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error)
    WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error)
    GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)
}
```

- **Rationale**:
  - Aligns with use case boundaries (create/sign/send)
  - `TransactionSigner` has zero network methods (offline-capable)
  - Minimal dependencies (only methods actually used)
  - Clear separation of concerns (account queries, signing, submission)

- **Trade-offs**:
  - ✅ Interface Segregation Principle fully applied
  - ✅ Offline signing guaranteed (no network methods in signer)
  - ✅ Clear contracts for each use case
  - ✅ Easier to test (mock minimal interfaces)
  - ⚠️ Requires use case refactoring (expected as part of alignment)

- **Follow-up**:
  - Update DI layer to wire segregated interfaces
  - Ensure existing XRPer implementations satisfy all three interfaces

## Risks & Mitigations

- **Risk 1**: Peersyst/xrpl-go signed transaction blobs incompatible with xrpscan/xrpl-go submission
  - **Mitigation**: Integration tests validating signing output → submission workflow; fallback to XRPLF if needed

- **Risk 2**: JSON file size impact on file transfer for offline signing
  - **Mitigation**: Monitor file sizes in testing; compress if exceeds reasonable limits (~10KB per transaction)

- **Risk 3**: Multi-signature flow complexity (sequential signing, signature count validation)
  - **Mitigation**: Comprehensive test scenarios (2-of-3, 3-of-5); validate against XRP testnet before mainnet

- **Risk 4**: Backward compatibility with in-flight text format transactions during migration
  - **Mitigation**: Dual-format support with format detection; deprecation warnings; staged rollout

- **Risk 5**: xrpl-go library maintenance and updates
  - **Mitigation**: Monitor library releases; abstract signing behind TransactionSigner interface for easy replacement

## References

### Official Documentation

- [XRP Ledger Multi-Signing](https://xrpl.org/docs/concepts/accounts/multi-signing) - Multi-signature transaction requirements
- [Set Up Multi-Signing Tutorial](https://xrpl.org/docs/tutorials/how-tos/manage-account-settings/set-up-multi-signing) - SignerList configuration
- [Secure Signing](https://xrpl.org/docs/concepts/transactions/secure-signing) - Best practices for transaction signing

### xrpl-go Libraries

- [XRPLF/xrpl-go](https://github.com/XRPLF/xrpl-go) - Official XRP Ledger Foundation implementation
- [Peersyst/xrpl-go](https://pkg.go.dev/github.com/Peersyst/xrpl-go) - Comprehensive wallet and signing functionality (selected for signing)
- [xrpscan/xrpl-go](https://github.com/xrpscan/xrpl-go) - Currently integrated (v0.2.11, used for submission)
- [Wallet Package Documentation](https://pkg.go.dev/github.com/Peersyst/xrpl-go/xrpl/wallet) - Sign and Multisign function reference

### JSON Standards

- [Google JSON Style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml) - API versioning and structure best practices
- [REST API Standards & Best Practices 2026](https://www.boltic.io/blog/rest-api-standards) - Modern API design patterns
- [JSON-API Versioning Strategy](https://github.com/json-api/json-api/issues/406) - Semantic versioning approaches

### Internal References

- [Gap Analysis](./gap-analysis.md) - Implementation gap evaluation and option analysis
- [Requirements](./requirements.md) - EARS-formatted acceptance criteria
