# Pattern 3 (P2SH-P2WPKH) Verification Report

## Overview

This document verifies the successful implementation of Pattern 3: P2SH-P2WPKH (BIP49 Nested SegWit) single-signature transactions.

## Issue Reference

- **Issue**: #362 - P2SH-P2WPKH transaction validation failure
- **Fix**: PR #363 - Implement P2SH-P2WPKH (BIP49) support with BIP143 compliance
- **Verification Date**: 2026-01-15

## Test Environment

| Component | Value |
|-----------|-------|
| Network | Bitcoin Regtest (local) |
| Test Script | `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh` |
| Command | `make btc-e2e-reset P=3` |

## Test Results

### ✅ All Phases Completed Successfully

1. **Infrastructure Setup**
   - Database container started and healthy
   - Bitcoin nodes (watch, keygen) started and healthy

2. **Wallet Configuration**
   - Bitcoin Core wallets created successfully
   - RPC endpoints configured correctly

3. **Key Generation**
   - HD keys generated for keygen wallet (client, deposit, payment, stored)
   - BIP49 derivation path used: `m/49'/1'/account'/change/index`

4. **Descriptor Export/Import**
   - Descriptors exported from keygen wallet
   - Format: `sh(wpkh([fingerprint/49'/1'/1']tpub.../0/*))`
   - Descriptors imported into watch wallet
   - 2000 addresses generated (1000 per descriptor: receive + change)

5. **P2SH-P2WPKH Address Creation**
   - Address format verified: `2...` (P2SH format in regtest)
   - Sample address: `2N3r5aBFJWBGavJpX7HNpaFRoZWWdyr1jTd`

6. **Test UTXO Generation**
   - 101 blocks generated to payment address
   - Coinbase maturity achieved (100 confirmations)
   - Balance verified: 50.00000000 BTC

7. **Payment Request Creation**
   - 3 payment requests created successfully
   - Sender address: `2N3r5aBFJWBGavJpX7HNpaFRoZWWdyr1jTd` (P2SH format verified)
   - Receiver addresses generated (P2SH-SegWit format for testing)

8. **Transaction Creation and Signing**
   - ✅ Unsigned transaction created: `payment_1_unsigned_0_1768471249785717000.psbt`
   - ✅ Transaction signed with keygen wallet (single signature)
   - ✅ Signed transaction: `payment_1_signed_0_1768471249856265000.psbt`

9. **Transaction Broadcast**
   - ✅ **Transaction sent successfully**
   - ✅ **Transaction ID**: `5b17761f9ea9f5cee3f5a613b368d4a69128cf8177df817ee4d842fc2580bcd7`
   - ✅ Transaction accepted into mempool

## Technical Verification

### BIP143 Compliance

The implementation correctly follows BIP143 (Segregated Witness Transaction Signature Verification):

1. **ScriptCode Construction**: P2PKH format for signature hash calculation

   ```
   OP_DUP OP_HASH160 <20-byte-pubkey-hash> OP_EQUALVERIFY OP_CHECKSIG
   ```

2. **Signature Hash Calculation**: Uses `CalcWitnessSigHash` with correct scriptCode

3. **PSBT Finalization**: Custom `finalizeP2SHP2WPKHInput()` function
   - ScriptSig: Contains redeemScript (P2WPKH script)
   - Witness: `[<signature>, <pubkey>]`

### Key Implementation Files

| File | Changes |
|------|---------|
| `internal/infrastructure/api/btc/btc/descriptor.go` | Added `deriveP2WPKHRedeemScript()` for `sh(wpkh(...))` descriptors |
| `internal/infrastructure/api/btc/btc/psbt.go` | BIP143-compliant signing and `finalizeP2SHP2WPKHInput()` finalization |

### Descriptor Format Verified

```
Descriptor: sh(wpkh([5ff97435/49'/1'/1']tpubDDhYwcV7NpBCsKUG6mU2L5qL8GBHwucnpYig9ziE7yWEzLzYQtJDgy3JatAqfRyxedR2A7p1CijnkV1LasdY3XAx9HGBcnv9DDgDw54TtNn/0/*))

Components:
- Type: sh(wpkh(...)) - P2SH-wrapped P2WPKH
- Derivation: [5ff97435/49'/1'/1'] - BIP49 path for testnet account 1
- Extended Pubkey: tpubDD...
- Range: /0/* - Receive addresses (0-999)
```

### Address Format Verified

```
Expected: 2... (P2SH format in regtest)
Actual:   2N3r5aBFJWBGavJpX7HNpaFRoZWWdyr1jTd ✅

Address Type Breakdown:
- Prefix: 2 (P2SH in regtest/testnet)
- Contains: SegWit P2WPKH script (OP_0 <20-byte-hash>)
- Backward Compatible: Can receive from any wallet
```

## Comparison with Other Patterns

| Aspect | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) | Pattern 5 (P2WPKH) |
|--------|-------------------|-------------------------|---------------------|
| BIP | BIP44 | **BIP49** | BIP84 |
| Address | `m...`/`n...` | **`2...`** | `tb1q...` |
| SegWit | No | **Yes (wrapped)** | Yes (native) |
| Legacy Compatible | Yes | **Yes** | No |
| Transaction Size | Largest | Medium | Smallest |
| Status | ✅ Working | ✅ **Working** | ✅ Working |

## Root Cause Analysis (Issue #362)

### Original Problem

Transaction validation was failing with:

```
mandatory-script-verify-flag-failed (Script evaluated without error but
finished with a false/empty top stack element)
```

### Root Cause

The original implementation was missing proper PSBT finalization for P2SH-P2WPKH:

1. ScriptSig construction was incorrect
2. Witness serialization was not properly implemented
3. Custom finalization function was needed to work around btcd limitations

### Solution (PR #363)

Implemented `finalizeP2SHP2WPKHInput()` function that:

1. ✅ Properly constructs scriptSig with redeemScript
2. ✅ Correctly serializes witness as `[<signature>, <pubkey>]`
3. ✅ Uses ScriptBuilder for proper length prefixing
4. ✅ Clears partial signatures after finalization

## Conclusion

✅ **Pattern 3 (P2SH-P2WPKH Single-sig) is fully implemented and verified**

- BIP143 compliance confirmed
- BIP49 derivation path working correctly
- Transactions create, sign, and broadcast successfully
- Address format matches specification (`2...` for regtest)
- E2E test passes completely without errors

## References

- BIP16: Pay to Script Hash (P2SH)
- BIP49: Derivation scheme for P2WPKH-nested-in-P2SH based accounts
- BIP141: Segregated Witness (Consensus layer)
- BIP143: Transaction Signature Verification for Version 0 Witness Program
- Issue #362: <https://github.com/hiromaily/go-crypto-wallet/issues/362>
- PR #363: <https://github.com/hiromaily/go-crypto-wallet/pull/363>
