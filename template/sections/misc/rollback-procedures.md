## Rollback Procedures

### When to Rollback

**Rollback Immediately If**:

1. **Critical Security Issue**
   - Nonce reuse detected
   - Private key leakage suspected
   - Signature verification failures

2. **Operational Failures**
   - Unable to complete signing workflow
   - Multiple consecutive transaction failures
   - File management system breaks down

3. **Software Bugs**
   - Critical bugs in MuSig2 implementation discovered
   - Database integrity issues
   - Wallet crashes during MuSig2 operations

**Do NOT Rollback For**:

- Minor coordination issues (solvable with training)
- Individual transaction failures (investigate first)
- Operator errors (training issue, not system issue)

### Rollback Checklist

Before initiating rollback:

- [ ] **Pause Operations**: Stop creating new MuSig2 addresses
- [ ] **Assess Impact**: Determine scope of rollback needed
- [ ] **Secure Funds**: Verify all funds are safe and accessible
- [ ] **Document Reason**: Record detailed reason for rollback
- [ ] **Notify Team**: Alert all operators of rollback decision
- [ ] **Backup Data**: Backup current state before rollback

### Rollback Process

#### Step 1: Stop MuSig2 Operations

```bash
# 1. Disable MuSig2 in configuration
vi config/wallet/watch_btc.toml

[multisig]
use_musig2 = false  # Set to false

# 2. Restart watch wallet
# (If needed - configuration changes may not require restart)

# 3. Verify MuSig2 disabled
keygen create multisig-address --account deposit
# Should create P2WSH address, not Taproot
```

#### Step 2: Assess Current State

```bash
# 1. Count MuSig2 addresses created
mysql> SELECT
    account,
    COUNT(*) as musig2_count,
    SUM(CASE WHEN taproot_address IS NOT NULL THEN 1 ELSE 0 END) as with_funds
FROM account_key
WHERE coin = 'btc' AND taproot_address IS NOT NULL
GROUP BY account;

# 2. Calculate total BTC in MuSig2 addresses
watch btc api listunspent | \
    jq '[.[] | select(.address | startswith("bc1p")) | .amount] | add'

# 3. Check pending MuSig2 transactions
# Any transactions in signing process that need completion
```

#### Step 3: Complete In-Flight Transactions

```bash
# IMPORTANT: Complete any partially signed MuSig2 transactions
# Don't leave funds in limbo

# 1. List all pending PSBTs
ls data/tx/btc/payment_*_unsigned*.psbt

# 2. Complete signing for each pending transaction
# Follow normal MuSig2 workflow to completion

# 3. Broadcast all pending transactions
# Ensure all funds are properly settled
```

#### Step 4: Sweep MuSig2 Funds to P2WSH

**CRITICAL**: This step moves funds from MuSig2 back to traditional multisig.

```bash
# 1. Create target P2WSH addresses
keygen create multisig-address --account deposit --count 10

# 2. Create sweep transaction (MuSig2 → P2WSH)
watch create sweep --from-type musig2 --to-type p2wsh --account deposit

# 3. Sign using MuSig2 workflow (these are MuSig2 inputs)
# Round 1: Nonce generation
keygen musig2 nonce --file sweep_musig2_to_p2wsh_unsigned_0.psbt
sign musig2 nonce --file sweep_musig2_to_p2wsh_unsigned_0_...1.psbt
sign musig2 nonce --file sweep_musig2_to_p2wsh_unsigned_0_...2.psbt

# Round 2: Partial signatures
keygen musig2 sign --file sweep_musig2_to_p2wsh_nonce_0.psbt
sign musig2 sign --file sweep_musig2_to_p2wsh_unsigned_0_...1.psbt
sign musig2 sign --file sweep_musig2_to_p2wsh_unsigned_1_...2.psbt

# Aggregation
watch musig2 aggregate --file sweep_musig2_to_p2wsh_unsigned_2_...3.psbt

# 4. Broadcast
watch send --file sweep_musig2_to_p2wsh_signed_3.psbt

# 5. Wait for confirmation
# Verify funds received at P2WSH addresses
```

#### Step 5: Validate Rollback

```bash
# 1. Verify all funds back in P2WSH addresses
watch btc api listunspent | \
    jq '[.[] | select(.address | startswith("bc1q")) | .amount] | add'
# Should equal: Original amount - fees

# 2. Verify no significant funds remain in MuSig2 addresses
watch btc api listunspent | \
    jq '[.[] | select(.address | startswith("bc1p")) | .amount] | add'
# Should be: 0 or dust amounts

# 3. Verify P2WSH signing works
# Create and sign test transaction using traditional workflow
```

#### Step 6: Document Rollback

```markdown
# MuSig2 Rollback Report

## Rollback Details
- Date: YYYY-MM-DD
- Reason: ___
- Funds Affected: ___ BTC
- Addresses Rolled Back: ___

## Root Cause Analysis
- Issue: ___
- Impact: ___
- Why it happened: ___
- Prevention: ___

## Actions Taken
1. ...
2. ...
3. ...

## Lessons Learned
- ...
- ...

## Future Plans
- [ ] Address root cause
- [ ] Additional testing required
- [ ] Team training gaps identified
- [ ] Consider re-migration timeline: ___
```

### Partial Rollback

In some cases, you may want partial rollback (keep some MuSig2, revert others):

```bash
# Example: Keep deposit account as MuSig2, revert payment account

# 1. Disable MuSig2 for specific account only
# (Configuration doesn't support per-account, so handle manually)

# 2. Sweep only payment account MuSig2 → P2WSH
watch create sweep --from-type musig2 --to-type p2wsh --account payment

# 3. Continue using MuSig2 for deposit account
# No changes needed

# 4. Update procedures to reflect hybrid setup
```

---
