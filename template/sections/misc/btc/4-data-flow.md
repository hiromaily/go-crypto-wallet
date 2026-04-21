## 4. Data Flow

### 4.1 Transaction Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    PSBT Transaction Flow                     │
└─────────────────────────────────────────────────────────────┘

1. CREATE (Watch Wallet - Online)
   ┌──────────────┐
   │ Watch Wallet │
   │   (Online)   │
   └──────┬───────┘
          │ RPC: walletcreatefundedpsbt
          │ - Auto-select inputs
          │ - Calculate fees
          │ - Add change output
          ↓
   ┌──────────────┐
   │ Unsigned     │
   │ PSBT File    │
   │ (Base64)     │
   └──────┬───────┘
          │ deposit_8_unsigned_0_1234.psbt
          │ Signatures: 0/2
          ↓

2. SIGN #1 (Keygen Wallet - Offline)
   ┌──────────────┐
   │    Keygen    │
   │   Wallet     │
   │  (Offline)   │
   └──────┬───────┘
          │ btcd: psbt.NewFromRawBytes
          │ btcd: updater.Sign (first key)
          │ btcd: packet.Serialize
          ↓
   ┌──────────────┐
   │ Partially    │
   │ Signed PSBT  │
   └──────┬───────┘
          │ deposit_8_unsigned_1_1235.psbt
          │ Signatures: 1/2
          ↓

3. SIGN #2 (Sign Wallet - Offline)
   ┌──────────────┐
   │     Sign     │
   │    Wallet    │
   │  (Offline)   │
   └──────┬───────┘
          │ btcd: psbt.NewFromRawBytes
          │ btcd: updater.Sign (second key)
          │ btcd: packet.Serialize
          ↓
   ┌──────────────┐
   │ Fully Signed │
   │ PSBT File    │
   └──────┬───────┘
          │ deposit_8_signed_2_1236.psbt
          │ Signatures: 2/2 ✓
          ↓

4. FINALIZE & BROADCAST (Watch Wallet - Online)
   ┌──────────────┐
   │ Watch Wallet │
   │   (Online)   │
   └──────┬───────┘
          │ RPC: finalizepsbt
          │ - Combine signatures
          │ - Create final scriptSig/witness
          │ - Extract transaction
          ↓
   ┌──────────────┐
   │ Final TX Hex │
   └──────┬───────┘
          │ RPC: sendrawtransaction
          ↓
   ┌──────────────┐
   │  Blockchain  │
   │   (Mined)    │
   └──────────────┘
```

### 4.2 File Flow

**Filename Convention**: `{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt`

**Examples**:

```
deposit_8_unsigned_0_1534744535097796209.psbt   # Created by Watch
deposit_8_unsigned_1_1534744536000000000.psbt   # Signed by Keygen (1/2)
deposit_8_signed_2_1534744537000000000.psbt     # Signed by Sign (2/2, complete)
```

### 4.3 PSBT Metadata Requirements

**For all inputs**, PSBT must include:

- Previous output amount (required for SegWit/Taproot)
- Previous output scriptPubKey
- Transaction ID and output index

**For P2SH/P2SH-SegWit**:

- Redeem script

**For P2WSH**:

- Witness script

**For P2TR (Taproot)**:

- Taproot internal key
- Taproot merkle root (if script path)

**Optional (for hardware wallets)**:

- BIP32 derivation paths
- Public keys

---
