## Overview

### What is MuSig2?

MuSig2 is a cryptographic protocol that enables multiple parties to create a **single aggregated Schnorr signature** that looks identical to a single-signature transaction on the blockchain. Unlike traditional multisig (P2SH, P2WSH) where multiple signatures are stored on-chain, MuSig2 aggregates multiple signatures into one, providing significant benefits.

### Benefits Over Traditional Multisig

| Feature | Traditional P2WSH Multisig | MuSig2 |
|---------|---------------------------|--------|
| **On-Chain Appearance** | Multiple signatures visible | Single signature (looks like single-sig) |
| **Transaction Size** | ~370-400 bytes (2-of-3) | ~200-250 bytes (30-50% smaller) |
| **Privacy** | Multisig is visible | Indistinguishable from single-sig |
| **Fees** | Higher (proportional to size) | 30-50% lower |
| **Signature Algorithm** | ECDSA | Schnorr (BIP340) |
| **Address Type** | P2WSH (bc1q...) | P2TR Taproot (bc1p...) |
| **Compatibility** | Older standard | Modern (Bitcoin Core 22.0+) |

### When to Use MuSig2

- ✅ **New multisig setups** - Best privacy and efficiency
- ✅ **High-volume operations** - Significant fee savings over time
- ✅ **Privacy-focused applications** - Transactions look like single-sig
- ✅ **Modern infrastructure** - Requires Bitcoin Core 22.0+ and Taproot support
- ⚠️ **Legacy multisig** - Traditional P2WSH still supported for backward compatibility

### How MuSig2 Works (Two-Round Protocol)

```
Round 1: Nonce Generation (Parallel)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Generate Nonce 1                   │
│  Sign Wallet 1 → Generate Nonce 2  (can run in     │
│  Sign Wallet 2 → Generate Nonce 3   parallel)      │
└─────────────────────────────────────────────────────┘
                        ↓
            Exchange nonces via PSBT files
                        ↓
Round 2: Signing (Sequential)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Create Partial Signature 1         │
│  Sign Wallet 1 → Create Partial Signature 2         │
│  Sign Wallet 2 → Create Partial Signature 3         │
└─────────────────────────────────────────────────────┘
                        ↓
            Collect partial signatures
                        ↓
Aggregation (Watch Wallet)
┌─────────────────────────────────────────────────────┐
│  Watch Wallet → Aggregate Partial Signatures        │
│              → Verify Final Signature               │
│              → Broadcast Transaction                │
└─────────────────────────────────────────────────────┘
```

**Key Security Feature:**

- Each wallet generates a **nonce** (random value) in Round 1
- Nonces must be **unique per transaction** and **never reused**
- Reusing nonces can leak private keys - this is critical!

---
