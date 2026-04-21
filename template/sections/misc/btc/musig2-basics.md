## MuSig2 Basics

### MuSig2 Address Format

MuSig2 uses **Taproot (P2TR)** addresses:

```
bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr
└─┬┘└──────────────────────────────────────────────────────────┘
  │                         Taproot address
bc1p = Bech32m Taproot address prefix
```

**Key Properties:**

- Starts with `bc1p` on mainnet, `tb1p` on testnet
- 62 characters long (Bech32m encoding)
- Contains aggregated public key (looks like single-sig)
- Supports Schnorr signatures only

### PSBT Files with MuSig2 Data

MuSig2 transactions use PSBT format with extension fields:

#### File Naming Convention

```
{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
```

**Examples:**

```
payment_15_unsigned_0_1735680000000000000.psbt       # Unsigned (no nonces)
payment_15_nonce_0_1735680000000000001.psbt          # After Round 1 (nonces)
payment_15_unsigned_1_1735680000000000002.psbt       # After partial sig 1
payment_15_unsigned_2_1735680000000000003.psbt       # After partial sig 2
payment_15_signed_3_1735680000000000004.psbt         # Fully signed (aggregated)
```

#### PSBT Extension Fields for MuSig2

MuSig2 data is stored in PSBT proprietary fields:

- **Public Nonces**: Stored per-input (66 bytes per signer)
- **Partial Signatures**: Stored per-input (32 bytes per signer)
- **Aggregated Signature**: Final signature after aggregation (64 bytes)

**Note**: These fields are automatically managed by wallet commands.

### MuSig2 States

A MuSig2 transaction progresses through these states:

```
1. Unsigned (no nonces) ──────> Watch Wallet creates PSBT
                    │
2. Nonces Generated ──────────> Round 1: All wallets generate nonces
                    │
3. Partially Signed (1 sig) ──> Keygen Wallet signs (Round 2)
                    │
4. Partially Signed (2 sig) ──> Sign Wallet 1 signs (Round 2)
                    │
5. Partially Signed (3 sig) ──> Sign Wallet 2 signs (Round 2)
                    │
6. Aggregated & Finalized ────> Watch Wallet aggregates & broadcasts
```

---
