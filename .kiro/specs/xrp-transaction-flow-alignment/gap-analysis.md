# Gap Analysis: XRP Transaction Flow Alignment

## Executive Summary

This gap analysis evaluates the implementation challenges for aligning XRP transaction flow with Bitcoin patterns through JSON-based file format and native Go signing. The analysis reveals a **mature file handling infrastructure** that can be extended with minimal changes, **existing xrpl-go integration** that simplifies submission, but **significant refactoring needed** for interface segregation and native signing implementation.

**Key Findings**:

- **Infrastructure Reuse**: 85% of transaction file handling infrastructure is reusable
- **xrpl-go Already Integrated**: Transaction submission library is present and functional
- **Interface Segregation Gap**: Current use cases depend on monolithic XRPer interface
- **Native Signing Missing**: No Go-native signing implementation (currently gRPC-based)
- **Multi-Signature Metadata Absent**: File format lacks signature count tracking

**Recommended Approach**: **Option C - Hybrid** (extend file repository + new interface definitions + refactor use cases)

---

## 1. Current State Investigation

### 1.1 Existing Assets

#### Transaction File Infrastructure

**Location**: `internal/infrastructure/storage/file/transaction/transaction.go`

**Current Capabilities**:

- ✅ File naming convention: `{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.{ext}`
- ✅ Directory creation and path validation
- ✅ PSBT file support (`.psbt` extension, base64-encoded)
- ✅ Hex file support (`.hex` extension for BCH)
- ✅ Text-based slice operations (`ReadFileSlice`, `WriteFileSlice`)
- ✅ Metadata extraction from file paths (`GetFileNameType`, `ValidateFilePath`)

**Interface Definition** (`internal/application/ports/file/interface.go`):

```go
type TransactionFileRepositorier interface {
    CreateFilePath(actionType, txType, txID, signedCount)
    ReadFile(path) / WriteFile(path, hexTx)
    ReadFileSlice(path) / WriteFileSlice(path, data)
    ReadPSBTFile(path) / WritePSBTFile(path, psbtBase64)
    WriteHexFile(path, hexTx)
    ValidateFilePath(filePath, expectedTxType)
    GetFileNameType(filePath)
}
```

**Gap**: No `ReadJSONFile` / `WriteJSONFile` methods for structured JSON transaction files.

#### XRP Use Cases

**Create Transaction** (`internal/application/usecase/watch/xrp/create_transaction.go`):

- Current format: CSV-like text with sender account type header
- Uses `WriteFileSlice` for multi-line output
- Depends on full `xrpCreateTxClient` interface (embedsBalanceChecker + TransactionPreparer)
- Output format: `senderAccount\nUUID,txJSON\n...`

**Sign Transaction** (`internal/application/usecase/sign/xrp/sign_transaction.go`):

- Current format: Parses CSV text format
- Uses gRPC-based signing via `xrpSignClient`  (TransactionSigner interface)
- No signature count validation or multi-signature metadata
- Output: CSV format with UUID, signed transaction ID, and blob

**Send Transaction** (`internal/application/usecase/watch/xrp/send_transaction.go`):

- Uses `ReadFileSlice` to parse signed transactions
- Submission already uses **xrpl-go library** for network calls
- No validation of signature completion status

#### XRP API Interfaces

**Location**: `internal/application/ports/api/xrp/interface.go`

**Current Structure**:

- `XRPer` (monolithic): Embeds XRPAdminer, XRPPublicer, XRPAPIer + additional methods
- `XRPAPIer` (comprehensive): 30+ methods including account, transaction, escrow, NFT operations
- No segregated `AccountInfoProvider` interface
- No segregated `TransactionSigner` interface for offline signing

**Current Dependencies**:

- Create use case: Depends on `xrpCreateTxClient` (custom interface embedding BalanceChecker + TransactionPreparer)
- Sign use case: Depends on `xrpSignClient` (custom interface for TransactionSigner)
- Send use case: Depends on full `XRPAPIer`

#### XRP Implementation

**xrpl-go Integration** (`internal/infrastructure/api/xrp/xrplgo/`):

- ✅ Already integrated: `github.com/xrpscan/xrpl-go v0.2.11`
- ✅ Transaction submission implemented (`SubmitTransaction`)
- ✅ Account info retrieval (`GetAccountInfo`)
- ✅ Ledger validation waiting (`WaitValidation`)
- ✅ Transaction status checking (`GetTransaction`)

**gRPC-based Signing** (`internal/infrastructure/api/xrp/xrp/`):

- Currently uses gRPC service for signing operations
- Protocol Buffers auto-generated: `*.pb.go` (DO NOT EDIT)
- **Gap**: No native Go signing implementation using xrpl-go's signing capabilities

### 1.2 Architecture Patterns

**File Naming Convention**:

```
Pattern: {baseDir}{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.{extension}

Examples:
- BTC:  ./data/tx/btc/deposit_8_unsigned_0_1634744535097796209.psbt
- BCH:  ./data/tx/bch/payment_42_signed_1_1634744535097796209.hex
- XRP:  ./data/tx/xrp/transfer_5_unsigned_0_1634744535097796209.json (target)
```

**Multi-Signature File Progression**:

```
transfer_5_unsigned_0_*.json      # Created by watch wallet
transfer_5_signed_1_*.json        # After keygen wallet signature
transfer_5_signed_2_*.json        # After sign wallet signature (ready if 2-of-3)
transfer_5_sent_2_*.json          # Broadcast to network
transfer_5_done_2_*.json          # Confirmed on ledger
```

**Interface Segregation Pattern** (from Bitcoin):

```go
// BTC Create Transaction Use Case
type createTxBTCClient interface {
    apibtc.ChainConfigProvider
    apibtc.AmountConverter
    apibtc.UTXOProvider
    apibtc.RawTransactionCreator
    apibtc.AddressOperator
    apibtc.BalanceChecker
    apibtc.PSBTCreator
}
```

XRP needs similar segregation for `AccountInfoProvider` and `TransactionSigner`.

### 1.3 Integration Surfaces

**Database Schema**:

- `xrp_tx` table: Tracks transaction metadata (txID, status, etc.)
- `xrp_tx_detail` table: Stores XRP-specific details
- `xrp_account_key` table: Stores master seeds for signing

**DTO Structures** (`internal/application/dto/xrp/transaction.go`):

```go
type TxInput struct {
    TransactionType    string
    Account            string
    Amount             string
    Destination        string
    Fee                string
    Flags              uint64
    LastLedgerSequence uint64
    Sequence           uint64
    SigningPubKey      string
    TxnSignature       string
    Hash               string
}

type Instructions struct {
    Fee                    string
    MaxFee                 string
    MaxLedgerVersion       uint64
    MaxLedgerVersionOffset uint64
    Sequence               uint64
    SignersCount           uint64  // Available for multisig
}
```

---

## 2. Requirements Feasibility Analysis

### 2.1 Technical Needs from Requirements

#### Requirement 1: JSON Transaction File Format

**Data Models**:

- JSON structure matching `dtoxrp.TxInput` plus metadata:
  - Version (e.g., "1.0.0")
  - Chain ("XRP")
  - Network ("mainnet" | "testnet")
  - CreatedAt (ISO 8601 timestamp)
  - Transactions array with:
    - UUID (transaction identifier)
    - UnsignedData (dtoxrp.TxInput serialized)
    - SenderAccount (account address)
    - SenderAccountType (domain account type)
    - SignatureCount (current signatures)
    - RequiredSignatures (threshold)
    - SignedBlob (optional, when signed)

**Services/Components**:

- Extend `TransactionFileRepositorier` interface with JSON methods
- JSON serialization/deserialization for transaction files
- File extension handler (`.json` instead of `.psbt` or `.hex`)

**Validation**:

- JSON schema validation for version compatibility
- UUID format validation
- Transaction data completeness validation

#### Requirement 2: Create Transaction Use Case Refactoring

**Data Models**:

- New `AccountInfoProvider` interface extracting account-related methods from XRPer
- Create transaction output DTO with JSON file path

**Services/Components**:

- Refactor `createTransactionUseCase` to use JSON file writer
- Define segregated `AccountInfoProvider` interface
- Update file generation logic to include comprehensive metadata

**Business Rules**:

- Preserve all transaction details for offline signing
- Initialize signature count to 0
- Set required signatures based on account configuration

#### Requirement 3: Sign Transaction Use Case Refactoring

**Data Models**:

- New `TransactionSigner` interface for offline signing methods
- Native Go signing using xrpl-go library

**Services/Components**:

- Implement native Go signing (no gRPC dependency)
- JSON transaction file parser
- Signature count incrementer
- Multi-signature completion detector

**Business Rules**:

- Skip re-signing if transaction already complete
- Increment signature count after successful signing
- Mark transaction as complete when `signatureCount >= requiredSignatures`

#### Requirement 4: Send Transaction Use Case Refactoring

**Services/Components**:

- JSON file parser for signed transactions
- xrpl-go submission (already implemented in `xrplgo` package)
- Signature completion validator

**Validation**:

- Verify `isComplete` flag before submission
- Validate signed blob format
- Check signature count meets threshold

#### Requirement 5-6: Offline Signing + Multi-Signature

**Security**:

- No network calls in sign wallet operations
- All required data in JSON file (no database lookups except secrets)
- Sequential file naming for audit trail

**Workflow**:

- File-based data transfer between online/offline systems
- Signature count tracking in file metadata
- Completion status flag

#### Requirement 7: Backward Compatibility

**Migration Strategy**:

- Detect file format: Try JSON parse, fallback to CSV
- Deprecation warnings for text format
- Dual-format support during migration period

#### Requirement 9: Interface Segregation

**Interface Definitions**:

```go
// New interfaces to define in application/ports/api/xrp/

type AccountInfoProvider interface {
    GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error)
    GetBalance(ctx context.Context, addr string) (float64, error)
}

type TransactionSigner interface {
    // Native Go signing (no gRPC)
    SignTransactionNative(ctx context.Context, txInput *dtoxrp.TxInput, secret string) (signedTxID, txBlob string, err error)
}

type TransactionSubmitter interface {
    SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error)
    WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error)
    GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)
}
```

### 2.2 Gaps and Constraints

#### Missing Capabilities

| Capability | Current State | Gap | Constraint |
|------------|---------------|-----|------------|
| JSON transaction files | Not implemented | Medium | Need JSON schema design + serialization |
| AccountInfoProvider interface | Not defined | Low | Extract from existing XRPer |
| TransactionSigner interface | gRPC-based | Medium | Need native Go signing implementation |
| Native Go signing | Not implemented | High | Requires xrpl-go signing research |
| Multi-signature metadata | Not tracked in files | Low | Extend JSON structure |
| Signature completion detection | Not implemented | Low | Add validation logic |
| JSON file read/write methods | Not implemented | Low | Extend file repository |
| Backward compatibility parser | Not implemented | Medium | Dual-format support complexity |

#### Research Needed

1. **xrpl-go Native Signing**:
   - Question: How to use xrpl-go for offline transaction signing?
   - Investigate: xrpl-go wallet/crypto signing APIs
   - Expected: Similar to `xrpl-js` signing patterns

2. **Multi-Signature Signing Flow**:
   - Question: How does xrpl-go handle multi-signature transaction signing?
   - Investigate: SignerListSet transaction signing and combination
   - Expected: Sequential signing with blob updates

3. **JSON Schema Versioning**:
   - Question: What versioning strategy for JSON transaction file format?
   - Options: Semantic versioning (1.0.0), date-based, or integer versions
   - Recommendation: Semantic versioning for clarity

#### Constraints

- **Clean Architecture**: Must maintain strict layer separation (no infrastructure in use cases)
- **Existing File Naming**: Must preserve current file path pattern for consistency
- **Database Schema**: XRP tables already exist; minimal changes preferred
- **Testing**: Must maintain backward compatibility during migration

### 2.3 Complexity Signals

**Complexity Indicators**:

- ✅ Simple CRUD: File read/write operations
- 🟡 Algorithmic Logic: JSON serialization/deserialization, signature count tracking
- 🟡 Workflows: Multi-phase transaction signing (unsigned → signed-1 → signed-N → sent)
- 🔴 External Integrations: Native Go signing with xrpl-go (requires research)

**Risk Areas**:

- Native Go signing implementation (unfamiliar API)
- Multi-signature transaction blob manipulation
- Backward compatibility with in-flight text-format transactions

---

## 3. Implementation Approach Options

### Option A: Extend Existing Components

**Which files/modules to extend**:

1. **`internal/application/ports/file/interface.go`**:
   - Add `ReadJSONTransactionFile(path) (TransactionFile, error)`
   - Add `WriteJSONTransactionFile(path, data TransactionFile) (string, error)`

2. **`internal/infrastructure/storage/file/transaction/transaction.go`**:
   - Implement JSON read/write methods
   - Add JSON schema validation

3. **`internal/application/usecase/watch/xrp/create_transaction.go`**:
   - Modify to use `WriteJSONTransactionFile` instead of `WriteFileSlice`
   - Change output format from CSV to JSON

4. **`internal/application/usecase/sign/xrp/sign_transaction.go`**:
   - Modify to parse JSON instead of CSV
   - Implement native Go signing (replace gRPC call)
   - Add signature count tracking

5. **`internal/application/usecase/watch/xrp/send_transaction.go`**:
   - Modify to parse JSON signed files
   - Add signature completion validation

**Compatibility Assessment**:

- ✅ No breaking changes to interface consumers outside XRP
- ✅ Existing PSBT/Hex methods remain unchanged
- ⚠️ XRP use cases require updates (expected as part of refactor)
- ✅ File repository backward compatible (new methods added, existing preserved)

**Complexity and Maintainability**:

- 📊 TransactionFileRepository file size: ~280 lines → ~400 lines (manageable)
- 📊 Create/Sign/Send use cases: Moderate changes to parsing logic
- ✅ Single Responsibility Principle: File repository handles multiple formats (acceptable for infrastructure)
- ✅ Cognitive load: Medium (JSON methods clearly separated from PSBT/Hex)

**Trade-offs**:

- ✅ **Pros**:
  - Minimal new files (faster implementation)
  - Leverages existing file path and naming infrastructure
  - Consistent with PSBT/Hex pattern already in repository
  - Easy to test (extend existing test suite)

- ❌ **Cons**:
  - TransactionFileRepository grows to handle 3 formats (text, PSBT, JSON)
  - XRP-specific logic mixed with general file operations
  - Backward compatibility adds conditional logic complexity

### Option B: Create New Components

**Rationale for new creation**:

- Distinct JSON transaction file handling vs. existing text/PSBT/Hex
- XRP-specific transaction lifecycle (metadata, multi-sig) differs from BTC/BCH
- Clean separation for xrpl-go signing vs. gRPC-based signing

**New Components**:

1. **`internal/application/ports/file/xrp_transaction.go`**:

   ```go
   type XRPTransactionFileHandler interface {
       ReadJSONFile(path string) (*XRPTransactionFile, error)
       WriteJSONFile(path string, data *XRPTransactionFile) (string, error)
       ValidateJSONSchema(data *XRPTransactionFile) error
   }
   ```

2. **`internal/infrastructure/storage/file/transaction/xrp_json.go`**:
   - JSON transaction file handler implementation
   - Schema validation
   - Version compatibility checking

3. **`internal/application/ports/api/xrp/segregated_interfaces.go`**:

   ```go
   type AccountInfoProvider interface { ... }
   type TransactionSigner interface { ... }
   type TransactionSubmitter interface { ... }
   ```

4. **`internal/infrastructure/api/xrp/xrplgo/native_signer.go`**:
   - Native Go signing using xrpl-go
   - Implements `TransactionSigner` interface

**Integration Points**:

- XRP use cases depend on `XRPTransactionFileHandler` instead of `TransactionFileRepositorier`
- Create use case uses `AccountInfoProvider` instead of full `XRPer`
- Sign use case uses `TransactionSigner` (native Go implementation)
- Send use case uses `TransactionSubmitter` (xrpl-go client)

**Responsibility Boundaries**:

- `XRPTransactionFileHandler`: JSON file I/O, schema validation, versioning
- `TransactionFileRepositorier`: Generic file operations (text, PSBT, Hex)
- `AccountInfoProvider`: Account queries only (balance, info)
- `TransactionSigner`: Offline signing only (no network calls)
- `TransactionSubmitter`: Ledger submission only (no signing)

**Trade-offs**:

- ✅ **Pros**:
  - Clean separation of XRP-specific logic
  - Easier to test in isolation (no PSBT/Hex coupling)
  - Interface Segregation Principle fully applied
  - Simpler file repository (stays focused on generic operations)
  - Future XRP enhancements isolated from BTC/BCH/ETH

- ❌ **Cons**:
  - More files to navigate (6-8 new files)
  - Duplicates file path logic (CreateFilePath, ValidateFilePath)
  - Requires careful interface design upfront
  - Integration in DI layer becomes more complex

### Option C: Hybrid Approach ⭐ **Recommended**

**Combination Strategy**:

**Part 1: Extend File Repository** (minimal changes):

- Add JSON file methods to `TransactionFileRepositorier` interface
- Implement in `TransactionFileRepository` (reuse existing path logic)
- **Rationale**: File path and naming convention is reusable infrastructure

**Part 2: Create Segregated Interfaces** (new files):

- Define `AccountInfoProvider`, `TransactionSigner`, `TransactionSubmitter` interfaces
- Implement native Go signing in `xrplgo/native_signer.go`
- **Rationale**: Interface segregation is architectural improvement

**Part 3: Refactor Use Cases** (update existing):

- Update create/sign/send use cases to use JSON methods and segregated interfaces
- Add backward compatibility parser as private helper in use cases
- **Rationale**: Use cases are the integration point; changes expected here

**Phased Implementation**:

**Phase 1**: Infrastructure (Low Risk)

- [ ] Add `ReadJSONTransactionFile` / `WriteJSONTransactionFile` to port interface
- [ ] Implement JSON methods in `TransactionFileRepository`
- [ ] Define transaction file JSON schema structure
- [ ] Write unit tests for JSON serialization/deserialization

**Phase 2**: Interface Segregation (Medium Risk)

- [ ] Define `AccountInfoProvider` interface in `application/ports/api/xrp/`
- [ ] Define `TransactionSigner` interface
- [ ] Implement `AccountInfoProvider` in `xrplgo` client (delegate existing methods)
- [ ] Research xrpl-go native signing API

**Phase 3**: Native Signing (High Risk - Research Heavy)

- [ ] Implement `TransactionSigner` with native Go signing
- [ ] Test signing output matches gRPC-based signing
- [ ] Validate signed transaction blob format

**Phase 4**: Use Case Refactoring (Medium Risk)

- [ ] Refactor create transaction to use JSON format and `AccountInfoProvider`
- [ ] Refactor sign transaction to use JSON parsing and native signing
- [ ] Refactor send transaction to validate JSON and submit via xrpl-go
- [ ] Add backward compatibility for text format

**Phase 5**: Testing & Migration (Low Risk)

- [ ] Integration tests for full transaction flow (create → sign → send)
- [ ] Multi-signature scenario tests (2-of-3, 3-of-5)
- [ ] Backward compatibility tests with legacy text files
- [ ] Documentation and migration guide

**Risk Mitigation**:

- **Phase 3 Blocker**: If native signing research fails, implement signing wrapper around gRPC (defer full native implementation)
- **Backward Compatibility**: Feature flag to enable/disable JSON format during migration
- **Incremental Rollout**: Deploy to testnet first, validate multi-sig flows, then mainnet

**Trade-offs**:

- ✅ **Pros**:
  - Balanced approach (reuse + new components)
  - Phased implementation reduces risk
  - Infrastructure changes minimal
  - Interface segregation achieved
  - Backward compatibility maintained

- ❌ **Cons**:
  - More complex planning required
  - Phases 2-4 interdependent (careful sequencing needed)
  - File repository handles 4 formats (text, PSBT, Hex, JSON) - acceptable complexity
  - Migration period requires dual-format support

---

## 4. Implementation Complexity & Risk

### 4.1 Effort Estimation

**Overall Effort**: **M (Medium - 5-7 days)**

**Breakdown by Phase**:

- Phase 1 (Infrastructure): **S** (1-2 days) - Straightforward JSON I/O
- Phase 2 (Interface Segregation): **S** (1 day) - Interface extraction
- Phase 3 (Native Signing): **M** (2-3 days) - Research + implementation
- Phase 4 (Use Case Refactoring): **S** (1-2 days) - Update existing logic
- Phase 5 (Testing): **S** (1 day) - Extend existing test patterns

**Justification**:

- Extends established patterns (file handling, interface segregation)
- Manageable complexity (JSON serialization is standard library)
- xrpl-go research mitigated by existing `xrplgo` integration knowledge
- Clear scope (10 requirements, 4 use cases, 3 interfaces)

### 4.2 Risk Assessment

**Overall Risk**: **Medium**

**High-Risk Areas**:

1. **Native Go Signing with xrpl-go**:
   - **Risk**: xrpl-go signing API unfamiliar; may not support all transaction types
   - **Mitigation**: Research xrpl-go documentation and examples first; fallback to gRPC wrapper if needed
   - **Impact**: High (blocks Requirement 3)

2. **Multi-Signature Transaction Blob Handling**:
   - **Risk**: Incorrect signature count or blob manipulation breaks multisig
   - **Mitigation**: Comprehensive tests with 2-of-3 and 3-of-5 scenarios; validate against XRP testnet
   - **Impact**: High (security and correctness critical)

**Medium-Risk Areas**:

1. **Backward Compatibility Parsing**:
   - **Risk**: Edge cases in text format parsing cause migration failures
   - **Mitigation**: Thorough testing with existing text files; dual-format support with feature flag
   - **Impact**: Medium (migration disruption)

2. **JSON Schema Evolution**:
   - **Risk**: Future schema changes break backward compatibility
   - **Mitigation**: Semantic versioning in JSON; schema validation with clear error messages
   - **Impact**: Medium (maintenance burden)

**Low-Risk Areas**:

1. **File Repository Extension**:
   - **Risk**: Low (similar to existing PSBT/Hex methods)
   - **Mitigation**: Unit tests cover new methods
   - **Impact**: Low

2. **Interface Segregation**:
   - **Risk**: Low (extracting methods from existing interface)
   - **Mitigation**: Existing implementations satisfy new interfaces
   - **Impact**: Low

**Overall Justification**:

- Known tech (Go stdlib JSON, xrpl-go library)
- Moderate integrations (file I/O, interface refactoring)
- Clear performance path (file-based, no heavy computation)
- Security well-defined (offline signing, no network in sign wallet)

---

## 5. Recommendations for Design Phase

### 5.1 Preferred Approach

**Option C - Hybrid Approach** for the following reasons:

1. **Pragmatic Balance**:
   - Reuses file path infrastructure (don't duplicate working code)
   - Creates new interfaces where segregation adds value (ISP compliance)
   - Minimal impact on existing BTC/BCH/ETH code (isolated changes)

2. **Risk Management**:
   - Phased implementation allows early validation
   - Infrastructure changes tested before use case refactoring
   - Native signing research can proceed in parallel with Phase 1-2

3. **Long-Term Maintainability**:
   - Interface segregation improves testability and clarity
   - JSON format provides structured metadata for future enhancements
   - Backward compatibility ensures smooth migration

### 5.2 Key Design Decisions

#### JSON Transaction File Schema

**Recommendation**: Define strict schema with semantic versioning

```json
{
  "version": "1.0.0",
  "chain": "XRP",
  "network": "mainnet",
  "createdAt": "2026-02-14T12:34:56Z",
  "transactions": [
    {
      "uuid": "01933e82-b6c9-7890-a123-456789abcdef",
      "unsignedData": {
        "transactionType": "Payment",
        "account": "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
        "destination": "rLHzPsX6oXkzU9LV8WMNNBx6D3dVyEfPu6",
        "amount": "1000000",
        "fee": "12",
        "sequence": 123,
        "lastLedgerSequence": 8820051,
        "flags": 2147483648
      },
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

**Versioning Strategy**:

- Major version: Breaking schema changes (e.g., field removal)
- Minor version: Backward-compatible additions (e.g., new optional fields)
- Patch version: Documentation or validation fixes

#### Interface Segregation Design

**Recommendation**: Three focused interfaces

```go
// Minimal account queries (for create transaction)
type AccountInfoProvider interface {
    GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error)
    GetBalance(ctx context.Context, addr string) (float64, error)
}

// Offline signing only (for sign transaction)
type TransactionSigner interface {
    SignTransactionNative(ctx context.Context, txInput *dtoxrp.TxInput, secret string) (signedTxID, txBlob string, err error)
}

// Network submission only (for send transaction)
type TransactionSubmitter interface {
    SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error)
    WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error)
    GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)
}
```

**Benefits**:

- Clear dependency boundaries (create, sign, send)
- Testability (mock minimal interfaces)
- Offline capability (TransactionSigner has no network methods)

#### Backward Compatibility Strategy

**Recommendation**: Detect-and-Parse pattern with deprecation warnings

```go
// In sign transaction use case
func (u *signTransactionUseCase) parseTransactionFile(filePath string) (*TransactionData, error) {
    // Try JSON first
    if strings.HasSuffix(filePath, ".json") {
        return u.parseJSONFile(filePath)
    }

    // Fallback to text format with warning
    logger.Warn("Using deprecated text format transaction file", "path", filePath)
    return u.parseTextFile(filePath)
}
```

**Migration Timeline**:

- Week 1-2: Deploy JSON format (new transactions use JSON)
- Week 3-4: Monitor for text format usage (log warnings)
- Month 2+: Deprecate text format (documentation update)
- Month 6: Remove text format support (breaking change release)

### 5.3 Research Items for Design Phase

1. **xrpl-go Native Signing API**:
   - **Investigate**: `xrpl-go` transaction signing methods
   - **Questions**:
     - How to sign `dtoxrp.TxInput` using xrpl-go wallet/crypto APIs?
     - Does xrpl-go support multi-signature transaction signing?
     - What format does xrpl-go produce for signed transaction blobs?
   - **Expected Outcome**: Code examples and API signatures for native signing implementation
   - **Fallback**: Wrapper around existing gRPC signing if xrpl-go signing insufficient

2. **Multi-Signature Blob Manipulation**:
   - **Investigate**: How XRP multi-signature transactions combine signatures
   - **Questions**:
     - Does each signer sign the same unsigned transaction, or does the blob evolve?
     - How to validate signature count matches SignerQuorum?
     - Can xrpl-go combine multiple signed blobs?
   - **Expected Outcome**: Multi-signature signing workflow documentation
   - **Resources**: XRP Ledger multi-signing documentation, xrpl-go examples

3. **JSON Schema Validation Library**:
   - **Investigate**: Go libraries for JSON schema validation (e.g., `gojsonschema`, `jsonschema`)
   - **Questions**:
     - Should schema validation be strict or lenient for forward compatibility?
     - How to handle version migration (1.0.0 → 1.1.0)?
   - **Expected Outcome**: Chosen library and validation strategy
   - **Alternative**: Manual validation with struct tags

4. **File Extension Precedence**:
   - **Investigate**: How to determine transaction file format when file path doesn't have extension
   - **Questions**:
     - Should `.json` extension be mandatory?
     - How to handle existing files without extensions?
   - **Expected Outcome**: File detection strategy (extension-based vs. content-based)
   - **Recommendation**: Require `.json` extension for new files; content-based fallback for legacy

---

## 6. Requirement-to-Asset Map

| Requirement | Current Asset | Gap | Approach |
|-------------|---------------|-----|----------|
| **Req 1**: JSON Transaction File Format | `TransactionFileRepositorier` (text/PSBT/Hex) | Missing JSON read/write | **Extend**: Add JSON methods to interface + implementation |
| **Req 2**: Create Transaction Refactoring | `createTransactionUseCase` (text format, full XRPer) | Missing `AccountInfoProvider`, JSON output | **Extend**: Refactor use case + **New**: Define interface |
| **Req 3**: Sign Transaction Refactoring | `signTransactionUseCase` (text format, gRPC signing) | Missing native Go signing, JSON parsing | **Extend**: Refactor use case + **New**: Implement native signer |
| **Req 4**: Send Transaction Refactoring | `sendTransactionUseCase` (text format) | Missing JSON parsing, validation | **Extend**: Refactor use case to parse JSON |
| **Req 5**: Offline Signing Support | Existing offline architecture | Missing offline-compatible interfaces | **New**: `TransactionSigner` with no network methods |
| **Req 6**: Multi-Signature Flow Support | File naming (`signedCount`) | Missing signature metadata in files | **Extend**: JSON schema with signature count fields |
| **Req 7**: Backward Compatibility | N/A | Missing text format fallback | **New**: Detection + parsing logic in use cases |
| **Req 8**: Testing and Validation | Existing test patterns | Missing JSON-specific tests | **Extend**: Add JSON test cases to existing suites |
| **Req 9**: Port Interface Segregation | Monolithic `XRPer` interface | Missing segregated interfaces | **New**: Define `AccountInfoProvider`, `TransactionSigner`, `TransactionSubmitter` |
| **Req 10**: Error Handling and Logging | Existing error wrapping patterns | Missing JSON parsing errors | **Extend**: Add JSON validation error messages |

---

## 7. Summary

### 7.1 Gap Summary

- **Infrastructure**: 85% reusable (file path, naming, repository pattern)
- **Interfaces**: 0% compliant (need full segregation for ISP)
- **Signing**: 0% native Go (currently gRPC-dependent)
- **File Format**: 0% JSON (currently text-based)

### 7.2 Recommended Path Forward

**Option C - Hybrid Approach** with **5-phase implementation**:

1. **Phase 1** (Infrastructure): Add JSON file methods → Low risk, high reuse
2. **Phase 2** (Interfaces): Define segregated interfaces → Low risk, clear value
3. **Phase 3** (Native Signing): Research and implement xrpl-go signing → High risk, mitigated by fallback
4. **Phase 4** (Use Cases): Refactor create/sign/send → Medium risk, phased testing
5. **Phase 5** (Testing): Integration tests and migration → Low risk, validation-focused

**Effort**: M (5-7 days)
**Risk**: Medium (native signing research + multi-sig complexity)

### 7.3 Success Criteria for Design Phase

- [ ] JSON transaction file schema defined with versioning strategy
- [ ] Three segregated interfaces documented (`AccountInfoProvider`, `TransactionSigner`, `TransactionSubmitter`)
- [ ] xrpl-go native signing approach researched and validated (or fallback identified)
- [ ] Multi-signature workflow documented (signature count tracking, blob handling)
- [ ] Backward compatibility strategy finalized (detection, parsing, deprecation timeline)
- [ ] File extension and detection logic specified
- [ ] Test scenarios defined for each requirement (unit, integration, migration)

---

**Next Command**: `/kiro:spec-design xrp-transaction-flow-alignment -y`
