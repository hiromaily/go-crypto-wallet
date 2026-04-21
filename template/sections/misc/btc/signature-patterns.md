## Signature Patterns

### Single-Sig

Pattern where a single private key is used for signing.

```
┌─────────────────────────────────────────────────────────┐
│                  SINGLE-SIG FLOW                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign with single key                │
│          ↓                                              │
│  3. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Characteristics:**

- Simple and fast
- Completed with a single signature
- Risk concentrated since there's only one private key

### Multi-Sig

Pattern where multiple private keys are used for signing. M-of-N (M signatures required out of N keys).

#### 3-of-3 Multisig

```
┌─────────────────────────────────────────────────────────┐
│                  3-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Sign2 Wallet: Sign (3rd signature)                 │
│          ↓                                              │
│  5. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 2-of-3 Multisig

```
┌─────────────────────────────────────────────────────────┐
│                  2-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Watch Wallet: Broadcast transaction                 │
│                                                         │
│  (Sign2 Wallet not required - completed with 2 sigs)   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### MuSig2 (Signature Aggregation)

Aggregate signature protocol based on Schnorr signatures. N-of-N multisig becomes the same size as single-sig.

```
┌─────────────────────────────────────────────────────────┐
│                    MUSIG2 FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Round 1: Nonce Generation (can be parallelized)       │
│  ├─ Keygen Wallet: Generate nonce                       │
│  ├─ Sign1 Wallet: Generate nonce                        │
│  └─ Sign2 Wallet: Generate nonce                        │
│          ↓                                              │
│  Round 2: Signing (sequential)                         │
│  ├─ Keygen Wallet: Create partial signature             │
│  ├─ Sign1 Wallet: Create partial signature              │
│  └─ Sign2 Wallet: Create partial signature              │
│          ↓                                              │
│  Aggregation:                                           │
│  └─ Watch Wallet: Aggregate & broadcast                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**MuSig2 Benefits:**

- Transaction size reduced by 30-50%
- Improved privacy (indistinguishable from single-sig)
- Reduced fees

---
