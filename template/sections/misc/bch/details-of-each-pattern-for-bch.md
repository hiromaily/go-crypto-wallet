## Details of Each Pattern for BCH

### BCH Pattern 3: BCH CashAddr P2SH 3-of-3 Multisig

**✅ Fully implemented and verified in `scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh`**

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 3-of-3 (Keygen + Sign1 + Sign2)
Address Format: bchreg:p... (P2SH multisig in regtest)
Key Derivation: m/44'/1'/account'/change/index (testnet/regtest)
```

**Implementation Status:**

- ✅ Infrastructure and wallet setup
- ✅ HD key generation and fullpubkey export/import
- ✅ 3-of-3 multisig address creation
- ✅ Payment request creation and database storage
- ✅ **UTXO retrieval** (Fixed in PR #426 - CashAddr format normalization)
- ✅ Transaction creation and signing workflow
- ⚠️ PSBT generation (Known issue with BCH label format)

**Workflow:**

1. Generate Seed in Keygen/Sign1/Sign2
2. Generate HD Key in Keygen (10 keys per account)
3. Generate HD Key in Sign1/Sign2
4. Export fullpubkey from Sign1/Sign2
5. Import fullpubkey to Keygen
6. Create 3-of-3 Multisig addresses in Keygen
7. Export addresses from Keygen (CashAddr format)
8. Import addresses to Watch wallet
9. Generate Test UTXO (regtest)
10. Create unsigned transaction → Sign 3 times (ECDSA) → Broadcast

**Signing Flow:**

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st ECDSA signature)
    ↓
Sign1 Wallet (2nd ECDSA signature)
    ↓
Sign2 Wallet (3rd ECDSA signature)
    ↓
Watch Wallet (broadcast)
```

**Note:** Unlike BTC's SegWit patterns, BCH transactions don't benefit from witness data separation, resulting in larger transaction sizes. BCH compensates with very low fees (~1 sat/byte).

For more BCH patterns, see [BCH Technical Reference](../../../../docs/chains/bch/README.md#e2e-transaction-patterns)

---
