## Understanding the Differences

### Address Format Comparison

#### Traditional P2WSH Multisig

```
Address Format: bc1q... (mainnet) or tb1q... (testnet)
Address Type: P2WSH (Pay-to-Witness-Script-Hash)
Witness Program: 32 bytes (SHA256 hash of script)
Script: Multisig script visible when spent

Example Address:
bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3

On-Chain Appearance When Spent:
- Multiple signatures visible (2-of-3 shows 2 signatures)
- Script revealed (shows it's multisig)
- Higher witness size (~200+ bytes for signatures)
```

#### MuSig2 Taproot

```
Address Format: bc1p... (mainnet) or tb1p... (testnet)
Address Type: P2TR (Pay-to-Taproot)
Witness Program: 32 bytes (Taproot output key)
Script: Aggregated key - no script visible

Example Address:
bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr

On-Chain Appearance When Spent:
- Single aggregated signature (64 bytes)
- Looks identical to single-signature transaction
- No multisig indication
- Smaller witness size (~100 bytes total)
```

### Database Schema Comparison

#### Traditional P2WSH (account_key table)

```sql
account_key
├─ id BIGINT PRIMARY KEY
├─ coin VARCHAR(10)
├─ account VARCHAR(20)
├─ idx INT
├─ wallet_address VARCHAR(255)        -- P2PKH/P2SH-SegWit
├─ p2wpkh_address VARCHAR(255)        -- Native SegWit
├─ p2wsh_address VARCHAR(255)         -- P2WSH multisig <-- OLD
├─ multisig_address VARCHAR(255) NULL -- For compatibility
└─ taproot_address VARCHAR(255) NULL  -- Not used yet
```

#### MuSig2 Taproot (account_key table)

```sql
account_key
├─ id BIGINT PRIMARY KEY
├─ coin VARCHAR(10)
├─ account VARCHAR(20)
├─ idx INT
├─ wallet_address VARCHAR(255)        -- P2PKH/P2SH-SegWit
├─ p2wpkh_address VARCHAR(255)        -- Native SegWit
├─ p2wsh_address VARCHAR(255)         -- Traditional multisig (backward compat)
├─ multisig_address VARCHAR(255) NULL -- Stores taproot_address for MuSig2
└─ taproot_address VARCHAR(255) NULL  -- P2TR MuSig2 <-- NEW
```

**Key Difference**: Both `taproot_address` and `multisig_address` store the Taproot address for compatibility.

### Transaction Size Comparison

#### Traditional P2WSH 2-of-3 Multisig Transaction

```
Transaction Structure:
├─ Version (4 bytes)
├─ Input Count (1 byte)
├─ Inputs
│  └─ Previous Output (36 bytes)
│  └─ Script Sig (empty for SegWit, ~1 byte)
│  └─ Sequence (4 bytes)
├─ Output Count (1 byte)
├─ Outputs
│  └─ Amount (8 bytes)
│  └─ Script PubKey (~34 bytes for P2WSH)
├─ Witness Data (in witness structure)
│  └─ Witness Stack Count (1 byte)
│  └─ Empty placeholder (1 byte)
│  └─ Signature 1 (~72 bytes)
│  └─ Signature 2 (~72 bytes)
│  └─ Redeem Script (~105 bytes for 2-of-3)
└─ Locktime (4 bytes)

Total: ~370-400 bytes (including witness data)
```

#### MuSig2 2-of-3 Taproot Transaction

```
Transaction Structure:
├─ Version (4 bytes)
├─ Input Count (1 byte)
├─ Inputs
│  └─ Previous Output (36 bytes)
│  └─ Script Sig (empty, ~1 byte)
│  └─ Sequence (4 bytes)
├─ Output Count (1 byte)
├─ Outputs
│  └─ Amount (8 bytes)
│  └─ Script PubKey (~34 bytes for P2TR)
├─ Witness Data (in witness structure)
│  └─ Witness Stack Count (1 byte)
│  └─ Aggregated Signature (64 bytes)  <-- Single signature!
└─ Locktime (4 bytes)

Total: ~200-250 bytes (including witness data)

Size Reduction: 40-45% smaller
```

### Signing Process Comparison

#### Traditional P2WSH Signing

```
1. Watch Wallet creates unsigned PSBT
   payment_15_unsigned.psbt

2. Each signer signs independently (parallel):
   Keygen: payment_15_unsigned.psbt → payment_15_signed_keygen.psbt
   Sign1:  payment_15_unsigned.psbt → payment_15_signed_sign1.psbt
   Sign2:  payment_15_unsigned.psbt → payment_15_signed_sign2.psbt

3. Watch Wallet combines signatures:
   Combine all signatures into one PSBT
   payment_15_signed_final.psbt

4. Watch Wallet broadcasts transaction
```

**Advantages**:

- Simple: One round of signing
- Parallel: All signers can sign simultaneously
- Independent: Signers don't need to coordinate

**Disadvantages**:

- Larger transactions (multiple signatures on-chain)
- Higher fees
- Privacy leak (multisig visible)

#### MuSig2 Signing

```
1. Watch Wallet creates unsigned PSBT
   payment_15_unsigned_0.psbt

2. Round 1: Nonce Generation (parallel):
   Keygen: payment_15_unsigned_0.psbt → payment_15_unsigned_0_...1.psbt (add nonce)
   Sign1:  payment_15_unsigned_0_...1.psbt → payment_15_unsigned_0_...2.psbt (add nonce)
   Sign2:  payment_15_unsigned_0_...2.psbt → payment_15_nonce_0_...3.psbt (add nonce)

3. Round 2: Partial Signatures (sequential):
   Keygen: payment_15_nonce_0.psbt → payment_15_unsigned_0_...1.psbt (partial sig)
   Sign1:  payment_15_unsigned_0_...1.psbt → payment_15_unsigned_1_...2.psbt (partial sig)
   Sign2:  payment_15_unsigned_1_...2.psbt → payment_15_unsigned_2_...3.psbt (partial sig)

4. Watch Wallet aggregates signatures:
   payment_15_unsigned_3.psbt → payment_15_signed_3.psbt (final signature)

5. Watch Wallet broadcasts transaction
```

**Advantages**:

- Smaller transactions (single aggregated signature)
- Lower fees (30-50% reduction)
- Better privacy (looks like single-sig)

**Disadvantages**:

- More complex: Two rounds of signing
- Coordination: Must collect all nonces before Round 2
- Security: Nonce reuse can leak private keys

---
