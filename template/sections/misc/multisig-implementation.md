## Multisig Implementation

### Traditional P2SH Multisig

BCH supports only traditional P2SH multisig (no P2WSH or MuSig2):

**Redeem Script Structure:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

**Example 2-of-3:**

```
2 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

### Multisig Characteristics

| Property | Value |
|----------|-------|
| **Maximum Keys** | 15 (OP_CHECKMULTISIG limit) |
| **Address Type** | P2SH (CashAddr `p` prefix) |
| **Transaction Size** | ~370-400 bytes (2-of-3) |
| **Privacy** | Low (multisig visible on-chain) |
| **Signature Type** | ECDSA |

### Multisig Workflow

```
1. Create Redeem Script
   └── <M> <PubKey1> ... <PubKeyN> <N> OP_CHECKMULTISIG

2. Create P2SH Address
   └── Hash160(redeemScript) → CashAddr P2SH

3. Create Transaction
   └── Reference UTXO with P2SH scriptPubKey

4. Sign Transaction (Sequential)
   └── Signer 1 → Signer 2 → ... → Signer M

5. Assemble scriptSig
   └── OP_0 <Sig1> <Sig2> ... <SigM> <redeemScript>

6. Broadcast
   └── Send to network
```

### Limitations (vs Bitcoin)

| Feature | BCH | BTC |
|---------|-----|-----|
| P2SH Multisig | ✅ | ✅ |
| P2WSH Multisig | ❌ | ✅ |
| Taproot Multisig | ❌ | ✅ |
| MuSig2 | ❌ | ✅ |
| Privacy | Low | High (with MuSig2) |
| Fee Efficiency | Standard | 30-50% lower with MuSig2 |

---
