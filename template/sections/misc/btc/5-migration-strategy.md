## 5. Migration Strategy

### 5.1 Recommended Approach: Clean Break

**Decision**: Replace CSV with PSBT immediately (no backward compatibility)

**Rationale**:

1. **Simplicity**: Single code path, no format detection logic
2. **Security**: PSBT is standardized and well-tested
3. **Compatibility**: PSBT works with other Bitcoin tools
4. **Future-Proof**: Foundation for advanced features (MuSig2, hardware wallets)

### 5.2 Migration Steps

#### Phase 1: Implementation (Issues #93-#98)

1. Implement PSBT infrastructure (#93)
2. Update file repository (#94)
3. Update Watch wallet (#95)
4. Update Keygen wallet (#96)
5. Update Sign wallet (#97)
6. Update finalization (#98)

#### Phase 2: Testing (#99)

1. End-to-end integration tests
2. Compatibility testing with Bitcoin Core
3. Performance benchmarking

#### Phase 3: Deployment

1. Complete all pending CSV transactions
2. Deploy PSBT-enabled binaries
3. Archive old CSV files (keep for audit)
4. Update operational procedures

### 5.3 Handling Existing CSV Files

**Options**:

**Option A: Complete Before Migration** (Recommended)

- Finish all pending CSV transactions
- Deploy PSBT after queue is clear
- Archive CSV files post-migration

**Option B: Conversion Tool** (If needed)

- Create CSV-to-PSBT conversion utility
- Convert pending transactions
- Validate converted PSBTs

**Option C: Dual Support** (Not recommended)

- Maintain both CSV and PSBT code paths
- Adds complexity
- Only if gradual migration required

**Selected**: Option A (Complete Before Migration)

---
