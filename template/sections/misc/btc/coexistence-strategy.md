## Coexistence Strategy

### Running Both Address Types

During migration (and optionally long-term), you'll run both P2WSH and MuSig2 addresses simultaneously.

#### Configuration

```toml
# config/wallet/watch_btc.toml

[multisig]
require_num = 2
pubkey_num = 3
use_musig2 = true          # Enable MuSig2 for new addresses
allow_legacy = true        # Allow traditional P2WSH (for backward compat)

[taproot]
enabled = true
```

#### Address Creation Strategy

```bash
# Option 1: Automatic (recommended)
# All new addresses use MuSig2 by default
keygen create musig2-address --account deposit --count 10

# Legacy addresses created only if explicitly requested
keygen create multisig-address --account deposit --count 5 --type p2wsh

# Option 2: Manual selection per transaction
# Specify address type for each payment
# (Requires custom tooling)
```

#### Transaction Handling

Different workflows for different address types:

**P2WSH Inputs (Traditional)**:

```bash
# 1. Create unsigned PSBT
watch create payment --from-address bc1q... --to bc1q... --amount 0.01

# 2. Sign (traditional workflow - single round)
keygen sign --file payment_15_unsigned.psbt
sign sign --file payment_15_signed_keygen.psbt

# 3. Combine and broadcast
watch send --file payment_15_signed_final.psbt
```

**MuSig2 Inputs (Taproot)**:

```bash
# 1. Create unsigned PSBT
watch create payment --from-address bc1p... --to bc1p... --amount 0.01

# 2. Round 1: Nonce generation
keygen musig2 nonce --file payment_15_unsigned_0.psbt
sign musig2 nonce --file payment_15_unsigned_0_...1.psbt
sign musig2 nonce --file payment_15_unsigned_0_...2.psbt

# 3. Round 2: Partial signatures
keygen musig2 sign --file payment_15_nonce_0.psbt
sign musig2 sign --file payment_15_unsigned_0_...1.psbt
sign musig2 sign --file payment_15_unsigned_1_...2.psbt

# 4. Aggregate and broadcast
watch musig2 aggregate --file payment_15_unsigned_2_...3.psbt
watch send --file payment_15_signed_3.psbt
```

**Mixed Inputs (Both Types)**:

```bash
# Transaction with both P2WSH and MuSig2 inputs
# 1. Create unsigned PSBT (watch wallet handles mixing)
watch create payment --amount 0.05 --to bc1p...

# 2. Sign P2WSH inputs (traditional workflow)
keygen sign --file payment_15_unsigned.psbt  # Signs only P2WSH inputs

# 3. Sign MuSig2 inputs (two-round workflow)
# Round 1: Nonces for MuSig2 inputs only
keygen musig2 nonce --file payment_15_partial_signed.psbt
sign musig2 nonce --file payment_15_partial_signed_...1.psbt
sign musig2 nonce --file payment_15_partial_signed_...2.psbt

# Round 2: Partial signatures for MuSig2 inputs
keygen musig2 sign --file payment_15_nonce_0.psbt
sign musig2 sign --file payment_15_signed_0_...1.psbt
sign musig2 sign --file payment_15_signed_1_...2.psbt

# 4. Aggregate MuSig2 signatures and broadcast
watch musig2 aggregate --file payment_15_signed_2_...3.psbt
watch send --file payment_15_signed_3.psbt
```

#### Database Queries for Coexistence

```sql
-- Monitor address type distribution
SELECT
    account,
    SUM(CASE WHEN p2wsh_address IS NOT NULL AND taproot_address IS NULL THEN 1 ELSE 0 END) as p2wsh_only,
    SUM(CASE WHEN taproot_address IS NOT NULL THEN 1 ELSE 0 END) as musig2,
    COUNT(*) as total
FROM account_key
WHERE coin = 'btc'
GROUP BY account;

-- Find UTXOs by address type
-- P2WSH UTXOs:
SELECT * FROM utxo WHERE address LIKE 'bc1q%';

-- MuSig2 UTXOs:
SELECT * FROM utxo WHERE address LIKE 'bc1p%';
```

### Long-Term Coexistence Considerations

#### Pros of Long-Term Coexistence

- Flexibility to choose address type per use case
- Backward compatibility maintained
- Gradual learning curve for team
- Risk mitigation (eggs not all in one basket)

#### Cons of Long-Term Coexistence

- Increased operational complexity
- Two workflows to maintain
- Potential for confusion/errors
- Higher maintenance burden

#### Recommendation

**Short-term** (during migration): Coexistence is necessary and beneficial.

**Long-term** (after 6-12 months): Consider full migration to MuSig2 for simplicity, unless:

- Regulatory requirements mandate traditional multisig
- External integrations require P2WSH
- Risk tolerance demands redundancy

---
