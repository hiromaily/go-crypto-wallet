## Quick Reference

### Identifying BTC Address Types

| Prefix | Type | BIP | SegWit |
|--------|------|-----|--------|
| `1...` | P2PKH | BIP44 | ❌ |
| `3...` | P2SH or P2SH-P2WPKH | BIP16/BIP49 | △ |
| `bc1q...` | P2WPKH or P2WSH | BIP84 | ✅ |
| `bc1p...` | P2TR (Taproot) | BIP86 | ✅ |

### Identifying BCH Address Types

| Prefix | Type | Multisig |
|--------|------|----------|
| `bitcoincash:q...` | P2PKH | ❌ |
| `bitcoincash:p...` | P2SH | ✅ |

### Transaction Size Comparison

| Pattern | Weight | vBytes | Notes |
|---------|--------|--------|-------|
| P2PKH Single-sig (1-in, 2-out) | ~680 | ~170 | Legacy |
| P2WPKH Single-sig (1-in, 2-out) | ~440 | ~110 | Native SegWit |
| P2TR Single-sig (1-in, 2-out) | ~396 | ~99 | Taproot |
| 2-of-3 P2WSH Multisig | ~1,100 | ~275 | Traditional Multisig |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | Signature Aggregation |

---
