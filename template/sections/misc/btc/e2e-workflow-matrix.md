## E2E Workflow Matrix

### BTC Pattern Matrix

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Support |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **P2PKH (BIP44)** | **Single-sig** | **`1...`** | **✅ e2e/e2e-p1-p2pkh-singlesig.sh** |
| **2** | **P2PKH (BIP44)** | **2-of-3 Multisig** | **`3...` (P2SH wrapped)** | **✅ e2e/e2e-p2-p2pkh-2of3.sh** |
| **3** | **P2SH-P2WPKH (BIP49)** | **Single-sig** | **`3...`** | **✅ e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh** |
| **4** | **P2SH-P2WSH (BIP49)** | **2-of-3 Multisig** | **`3...`/`2...`** | **✅ e2e/e2e-p4-p2sh-p2wsh-2of3.sh** |
| **5** | **P2WPKH (BIP84)** | **Single-sig** | **`bc1q...`** | **✅ e2e/e2e-p5-p2wpkh-singlesig.sh** |
| **6** | **P2WSH (BIP84)** | **2-of-3 Multisig** | **`bc1q...`** | **✅ e2e/e2e-p6-p2wsh-2of3.sh** |
| **7** | **P2WSH (BIP84)** | **3-of-3 Multisig** | **`bc1q...`** | **✅ e2e/e2e-p7-p2wsh-3of3.sh** |
| **8** | **P2SH-P2WSH** | **3-of-3 Multisig** | **`3...`** | **✅ e2e/e2e-p8-p2sh-p2wsh-3of3.sh** |
| **9** | **P2TR (BIP86)** | **Single-sig** | **`bc1p...`** | **✅ e2e/e2e-p9-p2tr-singlesig.sh** |
| **10** | **P2TR (BIP86)** | **MuSig2 (N-of-N)** | **`bc1p...`** | **✅ e2e/e2e-p10-p2tr-musig2.sh** |
| **11** | **P2TR (BIP86)** | **Tapscript (M-of-N)** | **`bc1p...`** | **✅ e2e/e2e-p11-p2tr-tapscript.sh** |

### BCH Pattern Matrix

BCH supports **fewer patterns** than BTC due to lack of SegWit, Taproot, and Schnorr signatures.

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Support |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **CashAddr P2PKH** | **Single-sig** | **`bitcoincash:q...`** | **✅ e2e/e2e-p1-p2pkh-singlesig.sh** |
| **2** | **CashAddr P2SH** | **2-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e/e2e-p2-p2sh-2of3.sh** |
| **3** | **CashAddr P2SH** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e/e2e-p3-p2sh-3of3.sh** |

**BCH Limitations:**

- ❌ No SegWit (no P2WPKH, P2WSH, P2SH-P2WPKH patterns)
- ❌ No Taproot (no P2TR patterns)
- ❌ No Schnorr signatures (ECDSA only)
- ❌ No MuSig2 (no signature aggregation)
- ❌ No Descriptor support (use address export/import instead)
- ❌ No PSBT format (use raw transaction hex)
- ❌ No Bech32/Bech32m encoding (use CashAddr)
- ❌ No BIP49/84/86 derivation paths (BIP44 only)
- ⚠️ Fee unit: sat/Byte (not sat/vByte - no witness discount)

For detailed BCH patterns, see [BCH Technical Reference](../../../../docs/chains/bch/README.md#e2e-transaction-patterns).

---
