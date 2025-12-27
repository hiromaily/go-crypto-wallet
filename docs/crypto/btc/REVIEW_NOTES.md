# Self-Review Notes for PSBT Implementation Design

**Reviewer**: Claude Sonnet 4.5 (Self-Review)
**Date**: 2025-12-27
**Documents Reviewed**:
- `psbt_implementation.md` (Technical Design)
- `psbt_poc_example.md` (POC Examples)

---

## Critical Issues Found 🔴

### 1. **Base64 Decoding Error in Code Examples**

**Location**: `psbt_implementation.md` lines 245, 297
**Severity**: HIGH
**Issue**: The code examples show passing base64 string directly to `psbt.NewFromRawBytes()`, but this function expects raw bytes, not base64-encoded string.

**Current (Incorrect)**:
```go
packet, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(psbtBase64)), true)
```

**Should Be**:
```go
// Decode base64 first
psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
if err != nil {
    return "", false, fmt.Errorf("failed to decode base64: %w", err)
}

// Then parse PSBT
packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
if err != nil {
    return "", false, fmt.Errorf("invalid PSBT: %w", err)
}
```

**Impact**: Code won't work as written. Will cause runtime errors.
**Fix**: Update both Keygen and Sign wallet code examples with correct base64 decoding.

---

### 2. **Missing Imports in POC Example**

**Location**: `psbt_poc_example.md` line 26, 151
**Severity**: MEDIUM
**Issue**: POC example uses `chainhash` and `ecdsa` without importing them.

**Missing Imports**:
```go
"github.com/btcsuite/btcd/chaincfg/chainhash"
"github.com/btcsuite/btcd/btcec/v2/ecdsa"
```

**Impact**: POC won't compile. This is documented code, should be correct.
**Fix**: Already fixed by renaming to `.md` format, but should add correct imports in example.

---

### 3. **Watch Wallet PSBT Creation Ambiguity**

**Location**: `psbt_implementation.md` Section 2.3
**Severity**: MEDIUM
**Issue**: Document recommends using `walletcreatefundedpsbt` RPC for Watch wallet, but current codebase (`create_transaction.go`) creates raw transactions manually with `CreateRawTransaction`.

**Questions**:
1. Will Watch wallet use RPC to create PSBT, or convert existing raw tx to PSBT?
2. If using RPC, how do we add previous transaction metadata (PrevTxs) that offline wallets need?
3. If converting raw tx, should document use `converttopsbt` RPC or btcd package?

**Recommendation**:
- **Option A** (Recommended): Create raw tx as currently done, then convert to PSBT using btcd package, adding all metadata
- **Option B**: Use `walletcreatefundedpsbt` but may lose control over input selection

**Impact**: Implementation approach needs clarification.
**Fix**: Add detailed section on Watch wallet PSBT creation workflow.

---

## Moderate Issues Found 🟡

### 4. **Inconsistent Parameter Types in Interface**

**Location**: `psbt_implementation.md` Section 3.1 (lines 360-383)
**Severity**: MEDIUM
**Issue**: `SignPSBTWithKey` method signature includes `prevTxs []PrevTx` parameter, but:
- For **offline wallets**: Previous tx metadata should already be in PSBT
- For **Watch wallet**: Uses RPC, doesn't need prevTxs parameter

**Current**:
```go
SignPSBTWithKey(psbtBase64 string, wifs []string, prevTxs []PrevTx) (string, bool, error)
```

**Should Be** (for offline signing):
```go
// Offline signing - all metadata in PSBT
SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
```

**Rationale**: BIP174 PSBT design includes all metadata needed for signing within the PSBT itself. If we need to pass `prevTxs` separately, it defeats the purpose of PSBT.

**Impact**: Interface design doesn't align with BIP174 principles.
**Fix**: Remove `prevTxs` parameter, ensure Watch wallet adds all metadata to PSBT during creation.

---

### 5. **Missing Error Handling for Partial Signatures**

**Location**: `psbt_implementation.md` lines 295-343
**Severity**: MEDIUM
**Issue**: Sign wallet example doesn't check if PSBT has partial signatures before signing.

**Should Add**:
```go
// Verify PSBT has at least one partial signature
hasPartialSig := false
for _, input := range packet.Inputs {
    if len(input.PartialSigs) > 0 {
        hasPartialSig = true
        break
    }
}
if !hasPartialSig {
    return "", false, errors.New("PSBT must have at least one signature from Keygen wallet")
}
```

**Impact**: Sign wallet might try to sign unsigned PSBT (should come from Keygen first).
**Fix**: Add validation in Sign wallet code example.

---

### 6. **File Naming Convention Edge Case**

**Location**: `psbt_implementation.md` Section 4.2 (lines 534-543)
**Severity**: LOW
**Issue**: File naming shows `.psbt` extension, but `CreateFilePath` method in current codebase doesn't add extensions.

**Current Codebase**:
```go
// transaction.go:63
return fmt.Sprintf("%s%s_%d_%s_%d_", baseDir, actionType.String(), txID, txType, signedCount)
```

**Proposed**:
```
deposit_8_unsigned_0_1534744535097796209.psbt
```

**Issue**: Need to update `CreateFilePath` to append `.psbt` or update `WritePSBTFile` to add it.

**Impact**: File naming implementation needs clarification.
**Fix**: Specify in #94 (File Repository) issue how extension is added.

---

## Minor Issues Found 🟢

### 7. **Typo in Section 2.3 Line 245**

**Location**: `psbt_implementation.md` line 245
**Issue**: Comment says "Parse PSBT using btcd package" but code shows wrong implementation.
**Fix**: Already covered in Critical Issue #1.

---

### 8. **Incomplete `hasPartialSignatures` Function**

**Location**: `psbt_implementation.md` line 303
**Severity**: LOW
**Issue**: Code references `hasPartialSignatures(packet)` function but doesn't provide implementation.

**Should Add**:
```go
func hasPartialSignatures(packet *psbt.Packet) bool {
    for _, input := range packet.Inputs {
        if len(input.PartialSigs) > 0 {
            return true
        }
    }
    return false
}
```

**Impact**: Helper function not defined.
**Fix**: Add to code examples or note as pseudo-code.

---

### 9. **Missing Signature Verification**

**Location**: `psbt_implementation.md` Section 2.3
**Severity**: LOW
**Issue**: Code examples don't show signature verification before finalization.

**Should Add** (in Watch wallet finalization):
```go
// Verify all signatures before finalization
for i, input := range packet.Inputs {
    if !input.IsSigned() {
        return fmt.Errorf("input %d is not fully signed", i)
    }
}
```

**Impact**: Edge case handling missing.
**Fix**: Add validation step in finalization examples.

---

### 10. **Documentation References BIP370 (PayJoin)**

**Location**: `psbt_implementation.md` Section 790 (References)
**Severity**: LOW
**Issue**: BIP 370 (PayJoin) is mentioned but not relevant to current implementation.

**Fix**: Remove BIP 370 reference or add note "future consideration".

---

## Structural Improvements Needed 📋

### 11. **Add Troubleshooting Section**

**Recommendation**: Add section on common errors and solutions:
- "PSBT is not fully signed" - Check signature count
- "Invalid PSBT format" - Check base64 encoding
- "Missing witness UTXO" - Ensure Watch wallet adds metadata
- "Signature verification failed" - Check key correspondence

### 12. **Add Version Compatibility Matrix**

**Recommendation**: Add clear table:

| Component | Minimum Version | Recommended | Notes |
|-----------|----------------|-------------|-------|
| Bitcoin Core | v0.17.0 | v22.0+ | v22.0+ for Taproot |
| btcd | v0.24.0 | v0.25.0 | PSBT support |
| Go | 1.19 | 1.21+ | Project requirement |

### 13. **Clarify RPC vs btcd Usage**

**Recommendation**: Add decision tree:

```
Watch Wallet PSBT Creation:
├─ Option A: Use CreateRawTransaction + btcd PSBT package ✅ (Recommended)
│  └─ Pros: Full control, consistent with current code, can add all metadata
│
└─ Option B: Use walletcreatefundedpsbt RPC
   └─ Pros: Simpler, let Bitcoin Core handle input selection
   └─ Cons: Less control, may not include all metadata needed offline
```

### 14. **Add Sequence Diagram**

**Recommendation**: Add more detailed sequence diagram showing:
1. Watch creates unsigned tx with `CreateRawTransaction`
2. Watch converts to PSBT with metadata
3. Watch writes `.psbt` file
4. Keygen reads, parses, signs
5. Keygen writes partially signed `.psbt`
6. Sign reads, parses, adds signature
7. Sign writes fully signed `.psbt`
8. Watch reads, finalizes, extracts, broadcasts

---

## Recommendations for Fixes

### Priority 1 (MUST FIX before implementation)

1. ✅ Fix base64 decoding in Keygen/Sign wallet examples
2. ✅ Clarify Watch wallet PSBT creation approach (RPC vs manual)
3. ✅ Remove `prevTxs` parameter from `SignPSBTWithKey` interface
4. ✅ Add validation for partial signatures in Sign wallet

### Priority 2 (SHOULD FIX before implementation)

5. ✅ Add missing imports to POC example
6. ✅ Define helper functions used in examples
7. ✅ Clarify file extension handling
8. ✅ Add troubleshooting section

### Priority 3 (NICE TO HAVE)

9. Add version compatibility matrix
10. Add detailed sequence diagram
11. Add decision tree for RPC vs btcd usage
12. Clean up irrelevant BIP references

---

## Overall Assessment

### Strengths ✅

1. **Comprehensive Coverage**: Document covers all aspects (library validation, architecture, migration, risks)
2. **Clear Rationale**: Hybrid approach well-justified
3. **Good Structure**: Logical flow from research → design → implementation
4. **Risk Assessment**: Thorough risk analysis with mitigation strategies
5. **Offline Focus**: Correctly prioritizes offline wallet security

### Weaknesses ⚠️

1. **Code Examples Have Errors**: Base64 decoding issue is critical
2. **Implementation Details Missing**: Watch wallet PSBT creation needs clarification
3. **Interface Design Flaw**: `prevTxs` parameter contradicts BIP174 design
4. **Edge Cases**: Some error handling scenarios not covered

### Overall Grade: **B+ (85/100)**

**Verdict**: **APPROVED WITH REVISIONS**

The document is fundamentally sound and provides a solid foundation for implementation. However, **critical code errors must be fixed** before proceeding to implementation phase. The technical approach (hybrid) is correct and well-reasoned.

---

## Action Items

- [ ] Fix base64 decoding in code examples (Critical)
- [ ] Clarify Watch wallet PSBT creation workflow (Critical)
- [ ] Update interface design to remove `prevTxs` parameter (High)
- [ ] Add validation examples for partial signatures (Medium)
- [ ] Add troubleshooting section (Medium)
- [ ] Review and approve fixed version (Required)

---

**Recommendation**: Fix Critical and High priority issues, then **PROCEED** with implementation.

**Reviewer**: Claude Sonnet 4.5
**Review Date**: 2025-12-27
**Status**: ✅ Review Complete - Revisions Required
