## Migration Process

### Overview

The migration follows a **gradual, phased approach** to minimize risk:

```
Phase 1: Planning (Week 1)
├─ Assess current setup
├─ Document infrastructure
├─ Create migration plan
└─ Set up rollback procedures

Phase 2: Infrastructure Setup (Week 2)
├─ Upgrade Bitcoin Core
├─ Configure wallets for MuSig2
├─ Update database schema
└─ Set up monitoring

Phase 3: Testing (Week 2)
├─ Create testnet MuSig2 addresses
├─ Test nonce generation
├─ Test signing workflow
└─ Verify transaction broadcast

Phase 4: Gradual Migration (Week 3-4)
├─ Create new MuSig2 addresses
├─ Start using MuSig2 for new transactions
├─ Monitor both address types
└─ Gradually increase MuSig2 usage

Phase 5: Fund Sweeping (Week 5)
├─ Sweep funds from P2WSH to P2TR addresses
├─ Retire old P2WSH addresses
└─ Monitor final transactions

Phase 6: Validation (Week 5)
├─ Verify all funds migrated
├─ Validate transaction history
├─ Document lessons learned
└─ Update operational procedures
```

### Phase 1: Planning and Preparation

**Duration**: 3-5 days

**Objective**: Understand current setup and create detailed migration plan.

#### Step 1.1: Document Current Setup

Create inventory of your multisig infrastructure:

```bash
# 1. Count traditional multisig addresses
mysql> SELECT
    account,
    COUNT(*) as address_count,
    SUM(CASE WHEN p2wsh_address IS NOT NULL THEN 1 ELSE 0 END) as p2wsh_count
FROM account_key
WHERE coin = 'btc'
GROUP BY account;

# Example output:
# +----------+---------------+-------------+
# | account  | address_count | p2wsh_count |
# +----------+---------------+-------------+
# | deposit  |           100 |         100 |
# | payment  |           200 |         200 |
# | stored   |            50 |          50 |
# +----------+---------------+-------------+

# 2. Check UTXOs per address type
watch btc api gettxoutsetinfo

# 3. Document wallet configuration
ls -la config/wallet/
# Review: watch_btc.toml, keygen_btc.toml, sign_btc.toml
```

Document findings:

- Total P2WSH addresses: ___
- Active addresses (with UTXOs): ___
- Total BTC in P2WSH addresses: ___
- Number of transactions per month: ___

#### Step 1.2: Calculate Expected Savings

```bash
# Average transaction size
P2WSH_SIZE=385  # bytes (2-of-3 multisig)
MUSIG2_SIZE=225 # bytes (2-of-3 MuSig2)

# Calculate savings
MONTHLY_TXS=100  # Your monthly transaction count
FEE_RATE=50      # sat/vB (adjust to current network conditions)

P2WSH_MONTHLY_FEES=$((P2WSH_SIZE * FEE_RATE * MONTHLY_TXS))
MUSIG2_MONTHLY_FEES=$((MUSIG2_SIZE * FEE_RATE * MONTHLY_TXS))
SAVINGS=$((P2WSH_MONTHLY_FEES - MUSIG2_MONTHLY_FEES))

echo "P2WSH monthly fees: $P2WSH_MONTHLY_FEES sats"
echo "MuSig2 monthly fees: $MUSIG2_MONTHLY_FEES sats"
echo "Monthly savings: $SAVINGS sats"

# Example output:
# P2WSH monthly fees: 1,925,000 sats (~0.019 BTC)
# MuSig2 monthly fees: 1,125,000 sats (~0.011 BTC)
# Monthly savings: 800,000 sats (~0.008 BTC)
```

#### Step 1.3: Identify Risks

Create risk register:

| Risk ID | Description | Likelihood | Impact | Mitigation |
|---------|-------------|-----------|---------|------------|
| R1 | Nonce reuse during signing | Medium | Critical | Database constraints + testing |
| R2 | Team unfamiliar with MuSig2 | High | Medium | Training + practice runs |
| R3 | Software bugs in MuSig2 code | Low | High | Extensive testnet testing |
| R4 | Bitcoin Core compatibility | Low | Medium | Verify version ≥ 22.0 |
| R5 | Coordination failures | Medium | Medium | Clear procedures + monitoring |

#### Step 1.4: Create Migration Timeline

```markdown
## Migration Timeline

### Week 1: Planning
- Day 1-2: Document current setup
- Day 3: Team training on MuSig2
- Day 4: Create detailed migration plan
- Day 5: Review and approval

### Week 2: Setup and Testing
- Day 1: Upgrade Bitcoin Core to v25.0+
- Day 2: Configure wallets for MuSig2
- Day 3-4: Testnet testing (address creation, signing)
- Day 5: Review test results

### Week 3: Initial Production Migration
- Day 1: Create first 10 MuSig2 addresses (deposit account)
- Day 2-3: Test production signing workflow
- Day 4-5: Monitor and evaluate

### Week 4: Gradual Expansion
- Day 1-2: Create MuSig2 addresses for payment account
- Day 3-4: Increase MuSig2 usage to 50% of new transactions
- Day 5: Monitor and evaluate

### Week 5: Fund Sweeping and Validation
- Day 1-3: Sweep funds from P2WSH to P2TR
- Day 4: Validate all funds migrated correctly
- Day 5: Document lessons learned, update procedures
```

#### Step 1.5: Prepare Rollback Plan

Document rollback triggers:

**Rollback Immediately If**:

- Nonce reuse detected
- Private key leakage
- Unable to sign transactions
- Critical software bugs discovered

**Rollback Procedure** (see [Rollback Procedures](#rollback-procedures) section)

### Phase 2: Infrastructure Setup

**Duration**: 2-3 days

**Objective**: Prepare technical infrastructure for MuSig2.

#### Step 2.1: Upgrade Bitcoin Core

```bash
# 1. Check current version
bitcoin-cli -version
# Must be ≥ v22.0 for Taproot support
# Recommended: v25.0+ for best support

# 2. Stop Bitcoin Core (if upgrading)
systemctl stop bitcoind

# 3. Backup wallet data
cp -r ~/.bitcoin/wallets ~/.bitcoin/wallets_backup_$(date +%Y%m%d)

# 4. Install Bitcoin Core v25.0+
# Download from: https://bitcoincore.org/en/download/
wget https://bitcoincore.org/bin/bitcoin-core-25.0/bitcoin-25.0-x86_64-linux-gnu.tar.gz
tar -xzf bitcoin-25.0-x86_64-linux-gnu.tar.gz
sudo install -m 0755 -o root -g root -t /usr/local/bin bitcoin-25.0/bin/*

# 5. Start Bitcoin Core
systemctl start bitcoind

# 6. Verify version
bitcoin-cli -version
# Should show: Bitcoin Core version v25.0.0

# 7. Wait for sync (if needed)
bitcoin-cli getblockchaininfo | grep verification
```

#### Step 2.2: Configure Wallets for MuSig2

```bash
# 1. Update watch wallet configuration
vi config/wallet/watch_btc.toml

# Add or verify these settings:
[bitcoin]
host = "127.0.0.1:8332"
user = "your-rpc-user"
pass = "your-rpc-password"

[multisig]
require_num = 2        # 2-of-3 multisig
pubkey_num = 3         # Total signers
use_musig2 = true      # Enable MuSig2

[taproot]
enabled = true         # Enable Taproot support

# 2. Update keygen wallet configuration
vi config/wallet/keygen_btc.toml

# Verify multisig and taproot settings match

# 3. Update sign wallet configuration
vi config/wallet/sign_btc.toml

# Verify multisig and taproot settings match
```

#### Step 2.3: Verify Database Schema

```bash
# 1. Check if schema is up to date
mysql -u wallet_user -p wallet_db

mysql> SHOW COLUMNS FROM account_key LIKE '%taproot%';
# Should show:
# +-----------------+--------------+------+-----+---------+-------+
# | Field           | Type         | Null | Key | Default | Extra |
# +-----------------+--------------+------+-----+---------+-------+
# | taproot_address | varchar(255) | YES  |     | NULL    |       |
# +-----------------+--------------+------+-----+---------+-------+

# 2. If column is missing, apply migration
# (This should already be done, but verify)
make db-migrate
```

#### Step 2.4: Set Up Monitoring

```bash
# 1. Create monitoring script for nonce tracking
vi scripts/monitoring/check_nonce_usage.sh

#!/bin/bash
# Check for nonce reuse (if nonce table exists)
mysql -u wallet_user -p wallet_db -e "
SELECT
    signer_id,
    COUNT(*) as usage_count
FROM musig2_nonces
WHERE used = true
GROUP BY signer_id, nonce
HAVING COUNT(*) > 1;
"

# If any results, nonce was reused - ALERT!

# 2. Set up transaction monitoring
vi scripts/monitoring/check_musig2_transactions.sh

#!/bin/bash
# Monitor MuSig2 transaction success rate
bitcoin-cli listtransactions "*" 100 | \
    jq '.[] | select(.address | startswith("bc1p"))' | \
    jq -s 'length'

echo "MuSig2 transactions in last 100: $count"

# 3. Create alerting (example using email)
# Integrate with your monitoring system (Prometheus, Grafana, etc.)
```

### Phase 3: Testing

**Duration**: 2-3 days

**Objective**: Validate MuSig2 functionality in testnet environment.

#### Step 3.1: Create Test MuSig2 Addresses

```bash
# 1. Start with testnet
# Ensure watch_btc.toml points to testnet Bitcoin Core

# 2. Create MuSig2 address for deposit account
keygen create musig2-address --account deposit

# Example output:
# Created 1 MuSig2 Taproot addresses for account 'deposit'
# Address: tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr

# 3. Verify address in database
mysql> SELECT id, account, taproot_address, addr_status
FROM account_key
WHERE coin = 'btc' AND account = 'deposit' AND taproot_address IS NOT NULL
LIMIT 1;

# Should show:
# +----+---------+------------------------------------------------------------------+-------------+
# | id | account | taproot_address                                                  | addr_status |
# +----+---------+------------------------------------------------------------------+-------------+
# |  1 | deposit | tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr |           3 |
# +----+---------+------------------------------------------------------------------+-------------+

# 4. Export address to watch wallet
keygen export musig2-address --account deposit

# 5. Import address in watch wallet
watch imports musig2-address --account deposit
```

#### Step 3.2: Test Nonce Generation Workflow

```bash
# 1. Watch wallet creates unsigned PSBT (using testnet funds)
# First, send some testnet BTC to the MuSig2 address
# Get testnet BTC from: https://testnet-faucet.com/btc-testnet/

# 2. Wait for confirmation, then create payment transaction
watch create payment --account payment --address tb1q... --amount 0.001

# Output: payment_15_unsigned_0.psbt created

# 3. Keygen generates nonce
keygen musig2 nonce --file payment_15_unsigned_0.psbt

# Output: payment_15_unsigned_0_...1.psbt created (with keygen nonce)

# 4. Sign1 generates nonce
sign musig2 nonce --file payment_15_unsigned_0_...1.psbt

# Output: payment_15_unsigned_0_...2.psbt created (with sign1 nonce)

# 5. Sign2 generates nonce
sign musig2 nonce --file payment_15_unsigned_0_...2.psbt

# Output: payment_15_nonce_0_...3.psbt created (with all nonces)

# 6. Verify all nonces are present
# PSBT should now contain nonces from all 3 signers
```

#### Step 3.3: Test Signing Workflow

```bash
# 1. Keygen creates partial signature
keygen musig2 sign --file payment_15_nonce_0.psbt

# Output: payment_15_unsigned_0_...1.psbt (with keygen partial signature)

# 2. Sign1 creates partial signature
sign musig2 sign --file payment_15_unsigned_0_...1.psbt

# Output: payment_15_unsigned_1_...2.psbt (with sign1 partial signature)

# 3. Sign2 creates partial signature
sign musig2 sign --file payment_15_unsigned_1_...2.psbt

# Output: payment_15_unsigned_2_...3.psbt (with all partial signatures)

# 4. Watch wallet aggregates signatures
watch musig2 aggregate --file payment_15_unsigned_3.psbt

# Output: payment_15_signed_3.psbt (with final aggregated signature)

# 5. Broadcast transaction
watch send --file payment_15_signed_3.psbt

# Output: Transaction ID: abcdef123456...
```

#### Step 3.4: Verify Transaction Broadcast

```bash
# 1. Check transaction in mempool
bitcoin-cli getmempoolentry <txid>

# 2. Check transaction size
bitcoin-cli getrawtransaction <txid> true | jq '.vsize'
# Should be ~200-250 bytes for MuSig2

# Compare with P2WSH size (~370-400 bytes) - should see 40%+ reduction

# 3. Wait for confirmation
watch -n 10 'bitcoin-cli gettransaction <txid> | jq .confirmations'

# 4. Verify on block explorer
# Testnet: https://blockstream.info/testnet/tx/<txid>
# Should show single signature (looks like single-sig transaction)
```

#### Step 3.5: Practice Error Scenarios

**Important**: Practice handling errors to build confidence.

```bash
# Scenario 1: Missing nonce
# Try to sign without generating nonce first
keygen musig2 sign --file payment_15_unsigned_0.psbt
# Expected: Error - nonces not present

# Scenario 2: Wrong file order
# Try to use Sign1 before Keygen
sign musig2 nonce --file payment_15_unsigned_0.psbt
# Expected: May work, but follow documented order

# Scenario 3: File corruption
# Corrupt PSBT file and try to use it
# Expected: Error - invalid PSBT format

# Practice recovery for each scenario
```

### Phase 4: Gradual Migration

**Duration**: 1-2 weeks

**Objective**: Gradually transition from P2WSH to MuSig2 in production.

#### Step 4.1: Create Initial Production MuSig2 Addresses

```bash
# Start with small batch (10 addresses)
# Choose low-risk account first (e.g., deposit)

# 1. Create 10 MuSig2 addresses
keygen create musig2-address --account deposit --count 10

# 2. Export to watch wallet
keygen export musig2-address --account deposit

# 3. Import in watch wallet
watch imports musig2-address --account deposit

# 4. Verify addresses
mysql> SELECT
    id, account, taproot_address, addr_status
FROM account_key
WHERE coin = 'btc' AND account = 'deposit' AND taproot_address IS NOT NULL
ORDER BY id DESC LIMIT 10;

# Should show 10 new addresses with status = 3 (MultisigAddressGenerated)
```

#### Step 4.2: Monitor Initial Transactions

```bash
# 1. Start monitoring script
scripts/monitoring/check_musig2_transactions.sh

# 2. Process first MuSig2 transaction
# Use new MuSig2 address for next incoming deposit

# 3. Complete signing workflow (Round 1 + Round 2)
# Follow documented procedure exactly

# 4. Verify transaction broadcast
# Check block explorer, transaction size, confirmation

# 5. Document any issues encountered
# Note: Resolution, lessons learned
```

#### Step 4.3: Gradual Increase

Week 1:

```bash
# Create 10 more addresses, use for 10% of new transactions
keygen create musig2-address --account deposit --count 10
```

Week 2:

```bash
# Create 50 more addresses, use for 30% of new transactions
keygen create musig2-address --account deposit --count 50
keygen create musig2-address --account payment --count 20
```

Week 3:

```bash
# Create 100 more addresses, use for 50% of new transactions
keygen create musig2-address --account deposit --count 100
keygen create musig2-address --account payment --count 50
```

Week 4:

```bash
# Use MuSig2 for 100% of new addresses
# All new addresses use MuSig2 by default
```

#### Step 4.4: Parallel Operation

During migration, both address types coexist:

```sql
-- Monitor address usage
SELECT
    account,
    SUM(CASE WHEN p2wsh_address IS NOT NULL AND taproot_address IS NULL THEN 1 ELSE 0 END) as p2wsh_count,
    SUM(CASE WHEN taproot_address IS NOT NULL THEN 1 ELSE 0 END) as musig2_count,
    COUNT(*) as total
FROM account_key
WHERE coin = 'btc'
GROUP BY account;

-- Example output:
-- +----------+-------------+--------------+-------+
-- | account  | p2wsh_count | musig2_count | total |
-- +----------+-------------+--------------+-------+
-- | deposit  |         100 |           70 |   170 |
-- | payment  |         200 |           50 |   250 |
-- +----------+-------------+--------------+-------+
```

### Phase 5: Fund Sweeping

**Duration**: 3-5 days

**Objective**: Move funds from old P2WSH addresses to new MuSig2 addresses.

#### Step 5.1: Plan Fund Sweeping

```bash
# 1. Identify P2WSH addresses with funds
watch btc api listunspent | jq '.[] | select(.address | startswith("bc1q"))'

# 2. Calculate total BTC in P2WSH addresses
watch btc api listunspent | \
    jq '[.[] | select(.address | startswith("bc1q")) | .amount] | add'

# Output: Total BTC in P2WSH addresses

# 3. Estimate sweep transaction fees
# Number of inputs: ___
# Estimated size: inputs * 370 bytes + outputs * 34 bytes
# Fee at 50 sat/vB: ___ sats

# 4. Choose target MuSig2 address(es)
# Option 1: Consolidate to single MuSig2 address
# Option 2: Distribute to multiple MuSig2 addresses (better privacy)
```

#### Step 5.2: Execute Fund Sweeping in Batches

**CRITICAL**: Sweep in small batches to minimize risk.

```bash
# Batch 1: Sweep 5 P2WSH addresses (smallest UTXOs first)
# 1. Create sweep transaction
watch create sweep --from-type p2wsh --to-type musig2 --count 5 --account deposit

# 2. Sign using traditional P2WSH process (NOT MuSig2)
# These are P2WSH inputs, use traditional signing
keygen sign --file sweep_1_unsigned.psbt
sign sign --file sweep_1_signed_keygen.psbt

# 3. Combine and broadcast
watch send --file sweep_1_signed_final.psbt

# 4. Wait for confirmation
# Verify funds received at MuSig2 address

# 5. Repeat for next batch
# Batch 2: Next 10 addresses
# Batch 3: Next 20 addresses
# Continue until all funds swept
```

#### Step 5.3: Validate Fund Migration

```bash
# 1. Verify all P2WSH addresses are empty
watch btc api listunspent | jq '.[] | select(.address | startswith("bc1q"))' | jq length
# Should return: 0 (no UTXOs in P2WSH addresses)

# 2. Verify total BTC in MuSig2 addresses matches expected
watch btc api listunspent | \
    jq '[.[] | select(.address | startswith("bc1p")) | .amount] | add'
# Should match: Original total - fees

# 3. Document sweep transactions
# Transaction IDs: ___
# Total fees paid: ___ sats
# Final BTC in MuSig2: ___ BTC
```

### Phase 6: Validation and Cleanup

**Duration**: 2-3 days

**Objective**: Finalize migration, validate success, update procedures.

#### Step 6.1: Final Validation

```bash
# 1. Verify all addresses migrated
mysql> SELECT
    account,
    SUM(CASE WHEN taproot_address IS NOT NULL THEN 1 ELSE 0 END) as musig2_addresses,
    SUM(CASE WHEN p2wsh_address IS NOT NULL AND taproot_address IS NULL THEN 1 ELSE 0 END) as legacy_only
FROM account_key
WHERE coin = 'btc'
GROUP BY account;

# legacy_only should be 0 for active accounts

# 2. Verify no funds remain in P2WSH addresses
# (Already checked in Phase 5.3)

# 3. Test MuSig2 transaction end-to-end
# Create, sign, broadcast one more transaction to confirm everything works
```

#### Step 6.2: Update Operational Procedures

```bash
# 1. Update documentation
# - Standard operating procedures
# - Signing workflows
# - Error recovery procedures
# - Monitoring runbooks

# 2. Archive old procedures
mv docs/procedures/p2wsh_signing.md docs/procedures/archive/

# 3. Create new MuSig2 procedures
cp docs/chains/btc/musig2/user-guide.md docs/procedures/musig2_signing.md
# Customize for your operational environment
```

#### Step 6.3: Document Lessons Learned

Create migration report:

```markdown
# MuSig2 Migration Report

## Migration Summary
- Start Date: YYYY-MM-DD
- Completion Date: YYYY-MM-DD
- Duration: ___ weeks
- Addresses Migrated: ___
- Funds Migrated: ___ BTC
- Sweep Transactions: ___
- Total Fees Paid: ___ sats

## Success Metrics
- Transaction Size Reduction: ___%
- Fee Savings: ___%
- Zero incidents
- Zero fund loss

## Issues Encountered
1. Issue: ...
   Resolution: ...
   Prevention: ...

2. Issue: ...
   Resolution: ...
   Prevention: ...

## Lessons Learned
- ...
- ...

## Recommendations for Future Migrations
- ...
- ...
```

#### Step 6.4: Schedule Post-Migration Review

```bash
# 1. Schedule review meeting (1 week after completion)
# Attendees: All operators, technical lead, stakeholders

# 2. Review topics:
# - Migration success metrics
# - Issues and resolutions
# - Process improvements
# - Training needs
# - Monitoring adjustments

# 3. Update procedures based on feedback
```

---
