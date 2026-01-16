# PSBT Signature Mismatch - Root Cause Analysis

**Issue**: #320
**Date**: 2026-01-12
**Status**: Root cause identified ✅, implementation pending

## Issue Summary

**Problem**: PSBT finalization fails with "witness script has 3 pubkeys, PSBT has 2 partial sigs, but none matched"

**Prerequisite**: This issue only surfaces after PR #321 (descriptor export fix) is merged. The descriptor export must include the keygen key first.

---

## Root Cause

**The Sign wallet stores WIF private keys at a STATIC address index (typically 0), but the descriptor-based workflow derives public keys DYNAMICALLY based on the actual address index being used.**

### Detailed Explanation

#### Descriptor Export (Keygen Wallet)

The descriptor contains an **extended public key** at the account level:

```
[09464f4b/49'/1'/1]tpubDDQ8...0/*
                   ^^^^^^^^^^ account-level xpub
                             ^^^ derives children at /0/0, /0/1, /0/2, ...
```

- **Derivation level**: m/49'/1'/1 (account level)
- **Child derivation**: /0/* (change=0, address_index varies)
- **Result**: Pubkeys at m/49'/1'/1/0/0, m/49'/1'/1/0/1, m/49'/1'/1/0/2, ...

#### Auth Key Generation (Sign Wallet)

File: `internal/infrastructure/wallet/key/hd_wallet.go`

The `CreateKey` method (lines 109-121):
1. Creates private key at **ACCOUNT level**: m/49'/1'/1
2. Derives child keys using `createKeysWithIndex()`

The `createKeysWithIndex` method (lines 218-279):
1. Line 224: Derives **change level**: m/49'/1'/1/0
2. Line 234: Derives **address index**: m/49'/1'/1/0/idxFrom
3. Line 241: Extracts **private key** from child
4. Line 249: Creates **WIF** from private key
5. Stores WIF in `auth_account_key` table

**Result**: Each row in `auth_account_key` stores a WIF at a **FIXED address index** (typically 0):
- `idxFrom=0`: WIF at m/49'/1'/1/0/0
- No extended private key stored
- No way to derive keys at other indices

#### PSBT Signing

File: `internal/application/usecase/sign/btc/sign_transaction.go`

When signing (lines 150-178):
1. Retrieves auth key from database (line 158)
2. Gets WIF: `authKey.WalletImportFormat` (line 172)
3. Passes WIF to `SignPSBTWithKey()`

File: `internal/infrastructure/api/btc/btc/psbt.go`

The signing logic (lines 694-798):
1. Decodes WIF to private key (line 269-274)
2. Derives public key: `privKey.PrivKey.PubKey().SerializeCompressed()` (line 779)
3. Signs PSBT with this public key

**Problem**: The WIF is ALWAYS at index 0 (m/49'/1'/1/0/0), regardless of which address index is being signed!

#### PSBT Creation (Watch Wallet)

When creating a PSBT for sending funds:
1. Watch wallet selects UTXOs from various addresses (e.g., address index 5: m/49'/1'/1/0/5)
2. Creates PSBT with witness script derived from descriptor
3. Witness script contains pubkey at m/49'/1'/1/0/5

#### Finalization (Watch Wallet)

File: `internal/infrastructure/api/btc/btc/psbt.go`

The finalization logic (lines 425-558):
1. Extracts pubkeys from witness script (line 445-462)
   - Contains keygen pubkey at m/49'/1'/1/0/5
   - Contains auth1 pubkey at m/49'/1'/1/0/5
   - Contains auth2 pubkey at m/49'/1'/1/0/5
2. Attempts to match PSBT partial signatures (lines 479-500)
   - Keygen signature uses pubkey at m/49'/1'/1/0/5 ✅
   - Auth1 signature uses pubkey at m/49'/1'/1/0/0 ❌ (wrong index!)
   - Auth2 signature uses pubkey at m/49'/1'/1/0/0 ❌ (wrong index!)
3. Matching fails because auth signatures use index 0 keys, not index 5 keys

---

## Why This Happens

### Current System Architecture

```
┌─────────────────────────────────────────────────────┐
│ Descriptor Export (Keygen)                          │
│ - Extended public key at account level: m/49'/1'/1  │
│ - Derives children dynamically: /0/0, /0/1, /0/2... │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│ Descriptor Import (Watch)                           │
│ - Imports descriptor with xpub at m/49'/1'/1        │
│ - Generates addresses: index 0, 1, 2, ..., 999      │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│ PSBT Creation (Watch) - Address Index 5 Example     │
│ - Selects UTXO at address index 5                   │
│ - Witness script has pubkey at m/49'/1'/1/0/5       │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│ PSBT Signing (Sign Wallets) ❌ MISMATCH              │
│ - Auth1: Uses WIF at m/49'/1'/1/0/0 (STATIC)        │
│ - Auth2: Uses WIF at m/49'/1'/1/0/0 (STATIC)        │
│ - Generates pubkey at index 0, not index 5!         │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│ PSBT Finalization (Watch) ❌ FAILURE                 │
│ - Witness script expects pubkeys at index 5         │
│ - Signatures have pubkeys at index 0                │
│ - No match found!                                    │
└─────────────────────────────────────────────────────┘
```

### Key Insight

**Descriptor-based workflows require deriving keys at the CORRECT address index during signing.**

The current system:
- ✅ Keygen wallet: Stores extended private key as seed, derives children correctly
- ❌ Sign wallets: Store WIF at fixed index 0, cannot derive other indices

---

## Solution Options

### Option 1: Store Extended Private Keys in Sign Wallets (Recommended)

**Change**: Store account-level extended private keys instead of WIF in `auth_account_key` table.

**Schema Change**:
```sql
ALTER TABLE auth_account_key
  ADD COLUMN account_extended_privkey VARCHAR(255) AFTER wallet_import_format;
```

**Signing Flow**:
1. Sign wallet retrieves extended private key from database
2. Parses PSBT to determine address index from derivation path in BIP32 derivation field
3. Derives child private key at correct index: m/49'/1'/1/0/{index}
4. Signs with derived private key
5. Pubkey matches witness script ✅

**Pros**:
- Supports any address index
- Matches descriptor-based workflow
- No need to store thousands of WIFs

**Cons**:
- Requires database migration
- Extended private keys are more sensitive (store account-level key, not individual address keys)
- Need to update key generation and signing logic

### Option 2: Pass Address Index to Signing Function

**Change**: Modify signing interface to accept address index parameter.

**Signing Flow**:
1. Watch wallet determines address index when creating PSBT
2. Passes address index to sign wallets (via PSBT file metadata or separate parameter)
3. Sign wallet derives WIF at specified index
4. Signs with correct WIF

**Pros**:
- No database migration required
- Can keep current WIF storage

**Cons**:
- Need to generate and store WIFs for all possible indices (0-999) in advance
- Large storage overhead (3 wallets × 3 accounts × 1000 indices = 9000 rows)
- Doesn't scale well

### Option 3: Extract Derivation Path from PSBT (Recommended)

**Change**: Parse PSBT BIP32 derivation fields to determine address index during signing.

**PSBT Format**:
PSBTs contain BIP32 derivation info for each public key:
```
PSBT_IN_BIP32_DERIVATION = {pubkey: (fingerprint, derivation_path)}
```

**Signing Flow**:
1. Parse PSBT to extract BIP32 derivation paths
2. Identify address index from path: m/49'/1'/1/0/INDEX
3. Derive private key at correct index from account-level extended key
4. Sign with derived key

**Pros**:
- Self-contained (PSBT has all necessary info)
- No need to pass index separately
- Matches BIP174 standard

**Cons**:
- Requires storing extended private keys (same as Option 1)
- More complex PSBT parsing logic

---

## Recommended Fix

**Implement Option 1 + Option 3 Hybrid:**

1. **Store extended private keys** in `auth_account_key` table:
   ```sql
   ALTER TABLE auth_account_key
     ADD COLUMN account_extended_privkey VARCHAR(255);
   ```

2. **Update key generation** to store both WIF (for backward compatibility) and extended private key:
   ```go
   // In generate_auth_key.go, store account-level xpriv
   // Retrieve account-level extended private key from HD key generator
   accountXpriv := accountPrivKey.String()
   ```

3. **Update signing logic** to:
   - Parse PSBT BIP32 derivation fields
   - Extract address index from derivation path
   - Derive child private key at correct index
   - Sign with derived key

4. **Migration path**:
   - Existing WIF can remain for non-descriptor workflows
   - New descriptor workflow uses extended private key + derivation

---

## Files Requiring Changes

### 1. Database Schema
- `tools/sqlc/schemas/03_sign.sql` - Add `account_extended_privkey` column
- `tools/sqlc/queries/auth_account_key.sql` - Update queries
- `tools/atlas/migrations/sign/` - Create migration

### 2. Key Generation
- `internal/infrastructure/wallet/key/hd_wallet.go` - Return account-level xpriv from CreateKey
- `internal/application/usecase/sign/shared/generate_auth_key.go` - Store xpriv to repository
- `internal/infrastructure/repository/cold/auth_account_key_sqlc.go` - Handle xpriv storage

### 3. PSBT Parsing
- `internal/infrastructure/api/btc/btc/psbt.go` - Add BIP32 derivation path parsing
- New method: `extractAddressIndexFromPSBT(psbt)` - Parse BIP32_DERIVATION fields

### 4. Signing Logic
- `internal/application/usecase/sign/btc/sign_transaction.go` - Retrieve xpriv, derive child key
- `internal/infrastructure/api/btc/btc/psbt.go` - Update SignPSBTWithKey to support extended keys
- New method: `deriveKeyAtIndex(xpriv, index)` - Derive child key from extended key

### 5. Testing
- Update all signing tests to use extended keys
- Add tests for multi-index signing
- Test PSBT BIP32 derivation parsing

---

## Verification Steps

After implementing the fix:

1. **Generate auth keys** with extended private key storage
2. **Export descriptor** from keygen wallet (already working after PR #321)
3. **Import descriptor** into watch wallet (already working)
4. **Create PSBT** for address at index 5
5. **Sign PSBT** with keygen wallet (should derive key at index 5)
6. **Sign PSBT** with sign1 wallet (should derive key at index 5)
7. **Finalize PSBT** - should succeed with matching signatures ✅

---

## Implementation Checklist

- [ ] Database schema changes
  - [ ] Add `account_extended_privkey` column to `auth_account_key`
  - [ ] Create Atlas migration
  - [ ] Update SQLC queries
  - [ ] Regenerate sqlcgen code
- [ ] Key generation updates
  - [ ] Modify `hd_wallet.go` to return account xpriv
  - [ ] Update `generate_auth_key.go` to store xpriv
  - [ ] Update repository to handle xpriv
- [ ] PSBT parsing
  - [ ] Add BIP32 derivation path extraction
  - [ ] Parse address index from derivation path
  - [ ] Add unit tests for parsing logic
- [ ] Signing logic
  - [ ] Update use case to retrieve xpriv
  - [ ] Implement child key derivation at index
  - [ ] Update SignPSBTWithKey interface
  - [ ] Add unit tests for multi-index signing
- [ ] Integration testing
  - [ ] Update E2E tests
  - [ ] Test with multiple address indices
  - [ ] Verify PSBT finalization succeeds

---

## Dependencies

This fix **depends on PR #321** (descriptor export fix) being merged first:
- PR #321 ensures descriptors contain keygen + auth keys (2-of-3)
- This PR fixes the signing to use correct address indices

**Merge order**:
1. Merge PR #321 (descriptor export fix)
2. Implement and merge this PR (PSBT signing fix)

---

**Analyzed by**: Claude Code
**Related PR**: #321 (descriptor export fix - prerequisite)
