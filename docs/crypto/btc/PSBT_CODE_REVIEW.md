# Code Review: PSBT Infrastructure Implementation

**Reviewer**: Claude Sonnet 4.5 (Self-Review)
**Date**: 2025-12-27
**Files Reviewed**:
- `internal/infrastructure/api/bitcoin/btc/psbt.go` (432 lines)
- `internal/infrastructure/api/bitcoin/btc/psbt_test.go` (431 lines)
- `internal/infrastructure/api/bitcoin/api-interface.go` (modified)

---

## Critical Issues Found 🔴

### 1. **Taproot Signature Support is BROKEN**

**Location**: `psbt.go` lines 235-250
**Severity**: CRITICAL
**Issue**: The signing implementation only supports ECDSA signatures, not Schnorr signatures required for Taproot (P2TR).

**Current Code**:
```go
// Line 235
sigHashes := txscript.NewTxSigHashes(parsed.Packet.UnsignedTx, nil)
hash, err := txscript.CalcWitnessSigHash(
    witnessUtxo.PkScript,
    sigHashes,
    txscript.SigHashAll,
    parsed.Packet.UnsignedTx,
    i,
    witnessUtxo.Value,
)

// Line 250
signature := ecdsa.Sign(privKey.PrivKey, hash)
```

**Problems**:
1. `CalcWitnessSigHash` only works for SegWit v0 (P2WPKH, P2WSH), NOT Taproot
2. `ecdsa.Sign` produces ECDSA signature, but Taproot requires Schnorr (BIP340)
3. `NewTxSigHashes(..., nil)` - second parameter should be PrevOutputFetcher for Taproot

**Impact**: **Taproot signing will completely fail**. This contradicts the claim of "P2TR support".

**Fix Required**:
```go
// Detect script type
scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(witnessUtxo.PkScript, b.GetChainConf())
if err != nil {
    return "", false, fmt.Errorf("failed to extract script type: %w", err)
}

var hash []byte
if scriptClass == txscript.WitnessV1TaprootTy {
    // Taproot (SegWit v1) - use Schnorr signature
    prevOutputFetcher := NewPrevOutputFetcher(parsed.Packet) // Need to implement
    sigHashes := txscript.NewTxSigHashes(parsed.Packet.UnsignedTx, prevOutputFetcher)

    hash, err = txscript.CalcTaprootSignatureHash(
        sigHashes,
        txscript.SigHashDefault, // Taproot uses SigHashDefault
        parsed.Packet.UnsignedTx,
        i,
        prevOutputFetcher,
    )
    if err != nil {
        return "", false, fmt.Errorf("failed to calculate taproot sig hash: %w", err)
    }

    // Sign with Schnorr
    signature, err := schnorr.Sign(privKey.PrivKey, hash)
    if err != nil {
        return "", false, fmt.Errorf("failed to create schnorr signature: %w", err)
    }
    sigBytes = signature.Serialize() // Schnorr sigs don't have sighash type appended

} else {
    // SegWit v0 (P2WPKH, P2WSH) - use ECDSA signature
    sigHashes := txscript.NewTxSigHashes(parsed.Packet.UnsignedTx, nil)
    hash, err = txscript.CalcWitnessSigHash(
        witnessUtxo.PkScript,
        sigHashes,
        txscript.SigHashAll,
        parsed.Packet.UnsignedTx,
        i,
        witnessUtxo.Value,
    )
    if err != nil {
        return "", false, fmt.Errorf("failed to calculate witness sig hash: %w", err)
    }

    // Sign with ECDSA
    signature := ecdsa.Sign(privKey.PrivKey, hash)
    sigBytes = append(signature.Serialize(), byte(txscript.SigHashAll))
}
```

---

### 2. **Input-PrevTx Mapping Validation Missing**

**Location**: `psbt.go` lines 56-96
**Severity**: HIGH
**Issue**: `CreatePSBT` doesn't validate that `prevTxs` corresponds to actual transaction inputs.

**Current Code**:
```go
for i, prevTx := range prevTxs {
    if i >= len(packet.UnsignedTx.TxIn) {
        return "", fmt.Errorf("prevTxs index %d exceeds number of inputs %d", i, len(packet.UnsignedTx.TxIn))
    }
    // ... adds metadata
}
```

**Problems**:
1. Doesn't validate `prevTxs[i].Txid` matches `msgTx.TxIn[i].PreviousOutPoint.Hash`
2. Doesn't validate `prevTxs[i].Vout` matches `msgTx.TxIn[i].PreviousOutPoint.Index`
3. If `len(prevTxs) < len(msgTx.TxIn)`, some inputs get no metadata (silent failure)
4. The check `i >= len(packet.UnsignedTx.TxIn)` happens AFTER processing starts

**Impact**: Could add wrong metadata to inputs, causing signing failures or security issues.

**Fix Required**:
```go
// Validate prevTxs count matches inputs BEFORE processing
if len(prevTxs) != len(msgTx.TxIn) {
    return "", fmt.Errorf("prevTxs count (%d) must match transaction inputs (%d)",
        len(prevTxs), len(msgTx.TxIn))
}

for i, prevTx := range prevTxs {
    // Validate prevTx corresponds to this input
    prevHash, err := chainhash.NewHashFromStr(prevTx.Txid)
    if err != nil {
        return "", fmt.Errorf("invalid prevTx txid for input %d: %w", i, err)
    }

    if !prevHash.IsEqual(&msgTx.TxIn[i].PreviousOutPoint.Hash) {
        return "", fmt.Errorf("prevTx[%d] txid mismatch: expected %s, got %s",
            i, msgTx.TxIn[i].PreviousOutPoint.Hash, prevHash)
    }

    if prevTx.Vout != msgTx.TxIn[i].PreviousOutPoint.Index {
        return "", fmt.Errorf("prevTx[%d] vout mismatch: expected %d, got %d",
            i, msgTx.TxIn[i].PreviousOutPoint.Index, prevTx.Vout)
    }

    // ... rest of metadata addition
}
```

---

### 3. **Hex Decoding Implementation is UNSAFE**

**Location**: `psbt.go` lines 414-431
**Severity**: MEDIUM-HIGH
**Issue**: `decodeHexScript` uses `fmt.Sscanf` incorrectly, causing potential buffer overflows.

**Current Code**:
```go
script := make([]byte, len(hexScript)/2)
_, err := fmt.Sscanf(hexScript, "%x", &script)
```

**Problems**:
1. `fmt.Sscanf` with `%x` and slice doesn't respect buffer size
2. If `hexScript` has odd length, `len(hexScript)/2` rounds down, causing issues
3. Doesn't validate that all characters are valid hex

**Impact**: Could cause panics or incorrect script parsing.

**Fix Required**:
```go
import "encoding/hex"

func (*Bitcoin) decodeHexScript(hexScript string) ([]byte, error) {
    if hexScript == "" {
        return nil, errors.New("empty hex script")
    }

    // Remove "0x" prefix if present
    if len(hexScript) >= 2 && hexScript[:2] == "0x" {
        hexScript = hexScript[2:]
    }

    // Use proper hex decoding
    script, err := hex.DecodeString(hexScript)
    if err != nil {
        return nil, fmt.Errorf("failed to decode hex script: %w", err)
    }

    return script, nil
}
```

---

## High Priority Issues 🟡

### 4. **GetPSBTFee Ignores Legacy Transactions**

**Location**: `psbt.go` lines 370-376
**Severity**: MEDIUM
**Issue**: Fee calculation only considers `WitnessUtxo`, ignoring `NonWitnessUtxo` used by legacy P2PKH.

**Current Code**:
```go
for _, input := range parsed.Packet.Inputs {
    if input.WitnessUtxo != nil {
        totalInput += input.WitnessUtxo.Value
    }
}
```

**Impact**: Fee calculation returns 0 or incorrect values for legacy transactions.

**Fix Required**:
```go
for i, input := range parsed.Packet.Inputs {
    if input.WitnessUtxo != nil {
        // SegWit/Taproot input
        totalInput += input.WitnessUtxo.Value
    } else if input.NonWitnessUtxo != nil {
        // Legacy input - need to find the correct output
        outIndex := parsed.Packet.UnsignedTx.TxIn[i].PreviousOutPoint.Index
        if int(outIndex) < len(input.NonWitnessUtxo.TxOut) {
            totalInput += input.NonWitnessUtxo.TxOut[outIndex].Value
        }
    }
}
```

---

### 5. **ValidatePSBT Rejects Legacy Transactions**

**Location**: `psbt.go` lines 177-182
**Severity**: MEDIUM
**Issue**: Validation requires `WitnessUtxo` for all inputs, rejecting valid legacy P2PKH transactions.

**Current Code**:
```go
for i, input := range parsed.Packet.Inputs {
    if input.WitnessUtxo == nil {
        return fmt.Errorf("input %d missing witness UTXO (required for SegWit/Taproot)", i)
    }
}
```

**Impact**: Cannot validate PSBTs with legacy inputs.

**Fix Required**:
```go
for i, input := range parsed.Packet.Inputs {
    if input.WitnessUtxo == nil && input.NonWitnessUtxo == nil {
        return fmt.Errorf("input %d missing UTXO information (need WitnessUtxo or NonWitnessUtxo)", i)
    }

    // For SegWit, WitnessUtxo is required
    // For Legacy, NonWitnessUtxo is required
    // Both is also valid (belt and suspenders approach)
}
```

---

### 6. **No Support for Witness Script (P2WSH)**

**Location**: `psbt.go` lines 81-90
**Severity**: MEDIUM
**Issue**: Only adds redeem script (P2SH), not witness script (P2WSH).

**Current Implementation**: Only handles `RedeemScript` for P2SH.

**Impact**: P2WSH multisig (SegWit native multisig) cannot be signed properly.

**Fix Required**:
```go
// Add redeem script for P2SH if provided
if prevTx.RedeemScript != "" {
    redeemScript, err := b.decodeHexScript(prevTx.RedeemScript)
    if err != nil {
        return "", fmt.Errorf("failed to decode redeemScript for input %d: %w", i, err)
    }
    if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
        return "", fmt.Errorf("failed to add redeem script for input %d: %w", i, err)
    }
}

// Add witness script for P2WSH if provided (need to add WitnessScript field to PrevTx)
if prevTx.WitnessScript != "" {
    witnessScript, err := b.decodeHexScript(prevTx.WitnessScript)
    if err != nil {
        return "", fmt.Errorf("failed to decode witnessScript for input %d: %w", i, err)
    }
    if err := updater.AddInWitnessScript(witnessScript, i); err != nil {
        return "", fmt.Errorf("failed to add witness script for input %d: %w", i, err)
    }
}
```

**Note**: This requires adding `WitnessScript string` field to the `PrevTx` struct in `transaction.go`.

---

## Medium Priority Issues 🟠

### 7. **Unused Type Definitions**

**Location**: `psbt.go` lines 18-28
**Severity**: LOW
**Issue**: `PSBTInput` and `PSBTOutput` types are defined but never used.

**Fix**: Remove unused types or document their intended future use.

---

### 8. **Inefficient Key Matching in Signing**

**Location**: `psbt.go` lines 232-271
**Severity**: LOW
**Issue**: Tries every key on every input, logging "Signature not applicable" many times.

**Impact**: Performance degradation and verbose logs for multisig.

**Optimization**: Could extract public keys from PSBT inputs and match before signing.

---

### 9. **No PSBT Version Check**

**Location**: `psbt.go` (missing)
**Severity**: LOW
**Issue**: Doesn't check or handle PSBT version (v0 vs v2).

**Impact**: May not be compatible with PSBT v2 (BIP370) in the future.

**Recommendation**: Add version check in `ParsePSBT`.

---

### 10. **No BIP32 Derivation Path Support**

**Location**: `psbt.go` (missing)
**Severity**: LOW
**Issue**: Doesn't handle BIP32 derivation paths in PSBT.

**Impact**: Cannot work optimally with hardware wallets and HD wallet tooling.

**Recommendation**: Add in future iteration when HD wallet support is added.

---

## Test Coverage Issues 🧪

### 11. **Tests Require Bitcoin Core Connection**

**Location**: `psbt_test.go` line 1
**Severity**: MEDIUM
**Issue**: All tests have `//go:build integration` tag, requiring Bitcoin Core.

**Impact**: Cannot run unit tests in CI without Bitcoin Core setup.

**Recommendation**: Add pure unit tests with mocked Bitcoin interface for basic validation.

---

### 12. **No Taproot-Specific Tests**

**Location**: `psbt_test.go` (missing)
**Severity**: HIGH (given Critical Issue #1)
**Issue**: No tests for Taproot signing workflow.

**Impact**: Taproot support is untested and broken (see Critical Issue #1).

**Fix Required**: Add integration tests for:
- P2TR address creation
- P2TR PSBT signing with Schnorr
- P2TR key path spending
- P2TR script path spending (if supported)

---

### 13. **No Legacy P2PKH Tests**

**Location**: `psbt_test.go` (missing)
**Severity**: MEDIUM
**Issue**: No tests for legacy P2PKH transactions.

**Impact**: Cannot verify legacy transaction support works correctly.

---

## Documentation Issues 📝

### 14. **Misleading Comments About Taproot**

**Location**: `psbt.go` lines 39-41, 191-193
**Severity**: MEDIUM
**Issue**: Comments claim "supports P2TR" but implementation is broken.

**Fix**: Either fix Taproot support or update comments to clarify limitations.

---

## Security Considerations 🔒

### 15. **No Private Key Zeroing**

**Location**: `psbt.go` lines 206-214
**Severity**: LOW
**Issue**: WIF private keys remain in memory after use.

**Recommendation**: Zero private key memory after use (though Go GC makes this difficult).

---

### 16. **No Rate Limiting on Signing Attempts**

**Location**: `psbt.go` lines 232-271
**Severity**: LOW
**Issue**: Could sign with many keys in a loop without limit.

**Impact**: Potential DoS if attacker provides many invalid keys.

**Recommendation**: Add max keys limit (e.g., 100 keys).

---

## Recommendations for Fixes

### Priority 1 (MUST FIX before merging)

1. ✅ **Fix Taproot signing** (Critical Issue #1) - This is a showstopper
2. ✅ **Add input validation** (Critical Issue #2) - Security issue
3. ✅ **Fix hex decoding** (Critical Issue #3) - Safety issue
4. ✅ **Fix fee calculation for legacy** (Issue #4)
5. ✅ **Fix validation for legacy** (Issue #5)

### Priority 2 (SHOULD FIX before production use)

6. ✅ Add P2WSH witness script support (Issue #6)
7. ✅ Add Taproot integration tests (Issue #12)
8. ✅ Add legacy P2PKH tests (Issue #13)
9. ✅ Update misleading documentation (Issue #14)

### Priority 3 (NICE TO HAVE)

10. Remove unused types (Issue #7)
11. Add PSBT version check (Issue #9)
12. Add unit tests without Bitcoin Core (Issue #11)

---

## Overall Assessment

### Strengths ✅

1. **Good structure**: Clean separation of concerns, well-organized functions
2. **Comprehensive API**: Covers all major PSBT operations
3. **Error handling**: Proper error wrapping with context
4. **Logging**: Good use of debug logging for troubleshooting
5. **Documentation**: Well-commented functions with clear descriptions

### Critical Weaknesses ⚠️

1. **Taproot support is broken**: Completely non-functional (Critical Issue #1)
2. **Input validation missing**: Security vulnerability (Critical Issue #2)
3. **Legacy support incomplete**: Fee calculation and validation issues
4. **No Taproot tests**: Claims are unverified

### Overall Grade: **C (70/100)**

**Verdict**: **REQUIRES MAJOR REVISIONS**

The infrastructure has a solid foundation, but **critical Taproot support is broken** and there are several security/correctness issues that must be fixed before this can be merged.

The biggest concern is that the code **claims to support Taproot** but the implementation will fail completely for P2TR inputs. This needs to be fixed immediately or the Taproot claims should be removed.

---

## Action Items

### Before Merging:
- [ ] Fix Taproot signing (CRITICAL - Issue #1)
- [ ] Add input-prevTx validation (HIGH - Issue #2)
- [ ] Fix hex decoding (MEDIUM-HIGH - Issue #3)
- [ ] Fix fee calculation for legacy (MEDIUM - Issue #4)
- [ ] Fix validation for legacy (MEDIUM - Issue #5)
- [ ] Add P2WSH support or document limitation (MEDIUM - Issue #6)

### Before Production:
- [ ] Add Taproot integration tests (Issue #12)
- [ ] Add legacy P2PKH tests (Issue #13)
- [ ] Update documentation to reflect actual capabilities (Issue #14)

### Technical Debt:
- [ ] Add unit tests without Bitcoin Core (Issue #11)
- [ ] Add PSBT version checking (Issue #9)
- [ ] Remove unused types (Issue #7)
- [ ] Optimize key matching (Issue #8)

---

**Reviewer**: Claude Sonnet 4.5
**Review Date**: 2025-12-27
**Status**: ⚠️ Major Revisions Required - DO NOT MERGE

**Recommendation**: Fix Critical Issues #1-3 before merging. Consider removing Taproot claims until Issue #1 is fully resolved and tested.
