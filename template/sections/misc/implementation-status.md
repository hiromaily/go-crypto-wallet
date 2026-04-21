## Implementation Status

### Key Type Implementation Status

| Key Type | BTC | BCH |
|----------|-----|-----|
| P2PKH (Legacy) | ✅ Implemented | N/A |
| P2SH-P2WPKH (BIP49) | ✅ Implemented | N/A |
| P2WPKH (BIP84) | ✅ Implemented | N/A |
| P2TR (BIP86) | ✅ Implemented | N/A |
| CashAddr | N/A | ✅ Implemented |

### Signature Pattern Implementation Status

| Signature Pattern | BTC | BCH |
|-------------------|-----|-----|
| Single-sig | ✅ Implemented | ✅ Implemented |
| 2-of-3 Multisig | ✅ Implemented (Fixed in #357) | ⚠️ Partial |
| 3-of-3 Multisig | ✅ Implemented | ✅ Implemented |
| MuSig2 | 🔜 In development | N/A |

### Descriptor Support

| Feature | BTC | BCH |
|---------|-----|-----|
| Descriptor Export | ✅ Implemented | ❌ Not supported |
| Descriptor Import | ✅ Implemented | ❌ Not supported |
| Bitcoin Core Integration | ✅ Implemented | N/A |

---
