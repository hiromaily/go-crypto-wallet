# Migration Guide: Traditional Multisig to MuSig2

This guide provides a comprehensive strategy for migrating from traditional P2WSH multisig to MuSig2 Taproot addresses in go-crypto-wallet.

## Table of Contents

1. [Introduction](#introduction)
2. [Should You Migrate?](#should-you-migrate)
3. [Prerequisites](#prerequisites)
4. [Understanding the Differences](#understanding-the-differences)
5. [Migration Process](#migration-process)
6. [Rollback Procedures](#rollback-procedures)
7. [Coexistence Strategy](#coexistence-strategy)
8. [Compatibility Considerations](#compatibility-considerations)
9. [FAQ](#faq)

---

## Introduction

### What This Guide Covers

This guide walks you through migrating from traditional Bitcoin multisig (P2WSH) to MuSig2 Taproot (P2TR) addresses. It covers:

- **Decision making**: When migration makes sense for your use case
- **Risk assessment**: Understanding what can go wrong
- **Step-by-step process**: Detailed migration phases with commands
- **Rollback procedures**: How to safely reverse the migration if needed
- **Coexistence**: Running both address types simultaneously

### Migration Timeline

Typical migration timeline for a production system:

```
Week 1: Planning and Assessment (Phase 1)
Week 2: Setup and Testing (Phase 2-3)
Week 3-4: Gradual Migration (Phase 4)
Week 5: Fund Sweeping and Validation (Phase 5-6)
```

**Important**: This is a gradual, low-risk migration. You can run both traditional and MuSig2 addresses simultaneously.

---

## Should You Migrate?

### Benefits of Migration

#### Cost Savings

MuSig2 transactions are 30-50% smaller than traditional P2WSH multisig:

| Transaction Type | P2WSH Multisig (2-of-3) | MuSig2 (2-of-3) | Savings |
|-----------------|-------------------------|-----------------|---------|
| **Transaction Size** | ~370-400 bytes | ~200-250 bytes | **40-45%** |
| **Fee (10 sat/vB)** | 3,700-4,000 sats | 2,000-2,500 sats | **40-45%** |
| **Fee (50 sat/vB)** | 18,500-20,000 sats | 10,000-12,500 sats | **40-45%** |

**Annual Savings Example**: If you make 1,000 transactions per year at 50 sat/vB:
- Traditional P2WSH: ~19,250,000 sats (~0.19 BTC)
- MuSig2: ~11,250,000 sats (~0.11 BTC)
- **Savings**: ~8,000,000 sats (~0.08 BTC)

#### Privacy Improvements

```
Traditional P2WSH Multisig (bc1q...):
├─ Visible on-chain: Multiple signatures
├─ Reveals: This is a multisig transaction
└─ Privacy: Low (blockchain analysis can identify multisig)

MuSig2 Taproot (bc1p...):
├─ Visible on-chain: Single aggregated signature
├─ Reveals: Looks like normal single-sig
└─ Privacy: High (indistinguishable from single-signature)
```

#### Modern Bitcoin Standard

- Uses Schnorr signatures (BIP340) - more efficient than ECDSA
- Leverages Taproot (BIP341) - future-proof for script enhancements
- Supported by Bitcoin Core 22.0+ (November 2021 onwards)

### Trade-offs to Consider

#### Two-Round Signing Process

Traditional P2WSH:
```
1. Create unsigned transaction
2. Each signer signs independently (can be parallel)
3. Combine signatures
4. Broadcast
```

MuSig2:
```
1. Create unsigned transaction
2. Round 1: Each signer generates nonce (parallel)
3. Collect all nonces
4. Round 2: Each signer creates partial signature (sequential)
5. Aggregate signatures
6. Broadcast
```

**Impact**: MuSig2 requires more coordination steps, but offers significant benefits.

#### Operational Complexity

| Aspect | Traditional P2WSH | MuSig2 |
|--------|------------------|--------|
| **Signing Steps** | 1 round | 2 rounds |
| **File Exchanges** | 1 PSBT file | 2-3 PSBT files (nonce + signatures) |
| **Error Modes** | Simpler | More complex (nonce reuse risk) |
| **Team Training** | Familiar | Requires training |
| **Debugging** | Well-documented | Newer, fewer resources |

#### Technology Maturity

- **Traditional P2WSH**: Battle-tested since 2017 (SegWit activation)
- **MuSig2**: Newer standard (BIP327 finalized 2023), less battle-tested
- **Library Support**: `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2` v2.3.6

### Decision Matrix

Use this matrix to decide whether migration makes sense for your use case:

| Your Situation | Recommendation |
|----------------|----------------|
| **High transaction volume** (>100 tx/month) | ✅ **Migrate** - Cost savings will be significant |
| **Privacy is critical** | ✅ **Migrate** - MuSig2 provides better privacy |
| **Low transaction volume** (<10 tx/month) | ⚠️ **Consider** - Savings may not justify complexity |
| **Risk-averse operations** | ⚠️ **Wait** - Let technology mature further |
| **Team unfamiliar with MuSig2** | ⚠️ **Test First** - Extended testing period recommended |
| **Legacy systems** | ❌ **Stay** - Migration overhead may be too high |
| **Simple operations preferred** | ❌ **Stay** - Traditional multisig is simpler |

### When NOT to Migrate

- **Critical production systems without extensive testing**: Test thoroughly first
- **Team lacks technical expertise**: Ensure adequate training
- **Regulatory uncertainty**: Check compliance requirements
- **Legacy integrations**: External systems may not support Taproot
- **Immediate need**: Migration requires planning and testing

---

## Prerequisites

### Technical Requirements

#### Bitcoin Core Version

```bash
# Check Bitcoin Core version
bitcoin-cli -version

# Required: v22.0 or higher for Taproot support
# Recommended: v25.0+ for best Taproot support
```

If using older version:
1. Backup wallet data
2. Stop Bitcoin Core
3. Upgrade to v25.0 or higher
4. Restart and verify sync

#### Go Version

```bash
# Check Go version
go version

# Required: Go 1.21 or higher
```

#### Database Schema

MuSig2 uses existing database tables with additional columns:

```sql
-- Check if taproot_address column exists
DESCRIBE account_key;

-- Should see:
-- taproot_address VARCHAR(255) NULL
-- multisig_address VARCHAR(255) NULL (for compatibility)
```

If columns are missing, run migration:
```bash
# Apply database migrations
make db-migrate
```

### Infrastructure Requirements

#### Test Environment

**CRITICAL**: Always test migration in testnet first.

Required test setup:
1. **Testnet Bitcoin Core node**
   - Testnet fully synced
   - RPC credentials configured
2. **Test Wallets**
   - Keygen wallet (offline test system)
   - Sign wallet(s) (offline test systems)
   - Watch wallet (online test system)
3. **Test Funds**
   - Small amount of testnet BTC for testing
   - Source: Bitcoin testnet faucets

#### Production Environment

Before migrating production:
1. **Backup Strategy**
   - Full database backup
   - Keystore backups
   - Configuration file backups
2. **Monitoring**
   - Transaction monitoring
   - Error alerting
   - Nonce tracking
3. **Rollback Plan**
   - Documented rollback procedures
   - Tested rollback process
   - Traditional multisig capability preserved

### Team Readiness

#### Required Knowledge

Your team should understand:

1. **MuSig2 Basics**
   - Two-round protocol
   - Nonce security requirements
   - Key aggregation
   - Read: `docs/crypto/btc/musig2_guide.md`

2. **Security Implications**
   - Nonce reuse consequences
   - Signature verification
   - Key management
   - Read: `docs/security/musig2_security.md`

3. **Operational Procedures**
   - File management
   - Error recovery
   - Monitoring
   - Read: `docs/crypto/btc/musig2_guide.md` (Best Practices)

#### Training Checklist

- [ ] All operators reviewed MuSig2 documentation
- [ ] Team completed testnet practice runs
- [ ] Security procedures documented and reviewed
- [ ] Error scenarios practiced
- [ ] Rollback procedures tested
- [ ] Incident response plan updated

### Risk Assessment

Before proceeding, assess these risks:

#### Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Nonce reuse** | Medium | Critical | Database constraints + monitoring |
| **Key loss** | Low | Critical | Robust backup procedures |
| **Software bugs** | Low | High | Thorough testing + gradual rollout |
| **Network issues** | Medium | Medium | Offline wallet architecture |
| **Human error** | High | Medium | Clear procedures + training |

#### Operational Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Team learning curve** | High | Low | Extended testing period |
| **Coordination overhead** | Medium | Medium | Clear file management process |
| **Monitoring gaps** | Medium | Medium | Enhanced monitoring setup |
| **Documentation gaps** | Low | Medium | This migration guide |

#### Business Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Downtime** | Low | High | Gradual migration |
| **Fund loss** | Very Low | Critical | Testnet validation + small batches |
| **Regulatory issues** | Low | High | Legal review before migration |
| **Customer confusion** | Medium | Low | Clear communication plan |

---

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
cp docs/crypto/btc/musig2_guide.md docs/procedures/musig2_signing.md
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

## Compatibility Considerations

### Database Schema Compatibility

#### Backward Compatibility

The schema maintains backward compatibility:

```sql
-- Old queries still work
SELECT * FROM account_key WHERE p2wsh_address IS NOT NULL;

-- New queries for MuSig2
SELECT * FROM account_key WHERE taproot_address IS NOT NULL;

-- Both types
SELECT * FROM account_key WHERE
    p2wsh_address IS NOT NULL OR taproot_address IS NOT NULL;
```

#### Schema Evolution

| Version | P2WSH Support | Taproot Support | MuSig2 Support |
|---------|---------------|-----------------|----------------|
| v1.0    | ✅ p2wsh_address | ❌ | ❌ |
| v2.0    | ✅ p2wsh_address | ✅ taproot_address | ❌ |
| v3.0    | ✅ p2wsh_address | ✅ taproot_address | ✅ (via taproot_address + musig2_nonces table) |

### Wallet Configuration Compatibility

#### Keygen Wallet

```toml
# v2.x configuration (legacy)
[multisig]
require_num = 2
pubkey_num = 3

# v3.x configuration (MuSig2)
[multisig]
require_num = 2
pubkey_num = 3
use_musig2 = false  # Default: false (backward compat)

[taproot]
enabled = true      # Required for MuSig2
```

#### Sign Wallet

Same configuration structure as Keygen wallet.

#### Watch Wallet

```toml
# Additional Bitcoin Core settings for Taproot
[bitcoin]
host = "127.0.0.1:8332"
# Minimum version: 22.0 for Taproot support
# Recommended: 25.0+
```

### External System Compatibility

#### Block Explorers

- **P2WSH**: Supported by all explorers
- **MuSig2 (P2TR)**: Supported by explorers with Taproot support
  - ✅ blockstream.info
  - ✅ mempool.space
  - ✅ blockchain.com (since late 2021)
  - ⚠️ Older explorers may not recognize bc1p... addresses

#### Wallets and Exchanges

- **Sending to MuSig2 addresses**: Most modern wallets/exchanges support sending to bc1p... addresses (as of 2023)
- **Receiving from MuSig2 addresses**: All wallets/exchanges accept from bc1p... addresses (looks like normal payment)

**Compatibility Matrix**:

| Service Type | P2WSH (bc1q...) | MuSig2 (bc1p...) | Notes |
|--------------|-----------------|------------------|-------|
| **Bitcoin Core 22.0+** | ✅ | ✅ | Full support |
| **Bitcoin Core 0.21.x** | ✅ | ❌ | No Taproot |
| **Hardware Wallets** | ✅ | ✅ (most) | Ledger: Firmware 2.1.0+, Trezor: Firmware 2.4.2+ |
| **Major Exchanges** | ✅ | ✅ (most) | Verify with specific exchange |
| **Mobile Wallets** | ✅ | ✅ (varies) | Check wallet documentation |

#### API Compatibility

If you expose APIs for wallet integration:

```javascript
// Old API (P2WSH)
POST /api/v1/address/create
{
  "account": "deposit",
  "type": "p2wsh"  // explicit type
}

// New API (MuSig2)
POST /api/v1/address/create
{
  "account": "deposit",
  "type": "musig2"  // or "taproot"
}

// Backward-compatible default
POST /api/v1/address/create
{
  "account": "deposit"
  // type: "p2wsh" by default for backward compat
}
```

---

## FAQ

### General Questions

#### Q: How long does migration take?

**A**: Typical timeline: 4-6 weeks for full migration (including planning, testing, gradual rollout, and fund sweeping). If you're only creating new MuSig2 addresses without sweeping old funds, you can complete in 2-3 weeks.

#### Q: Can I migrate incrementally?

**A**: Yes! This is the recommended approach. Start with 10 addresses, test thoroughly, then gradually increase. Both P2WSH and MuSig2 can coexist indefinitely.

#### Q: What if I encounter issues mid-migration?

**A**: You can pause at any time. Complete any in-flight transactions, then stop creating new MuSig2 addresses. See [Rollback Procedures](#rollback-procedures) for details.

#### Q: Do I need to sweep old P2WSH funds immediately?

**A**: No. You can keep old P2WSH addresses active and use them alongside MuSig2. Sweep when convenient (low fee environment, scheduled maintenance window, etc.).

### Technical Questions

#### Q: Are MuSig2 transactions more expensive to create?

**A**: MuSig2 transactions are **cheaper on-chain** (30-50% lower fees) but require **more coordination** (two rounds vs one round). The cost savings far outweigh the coordination overhead for most use cases.

#### Q: Can I mix P2WSH and MuSig2 inputs in one transaction?

**A**: Yes! A single transaction can spend from both P2WSH and MuSig2 UTXOs. Each input type uses its respective signing workflow.

#### Q: What happens if a nonce is reused?

**A**: **CRITICAL SECURITY ISSUE**. Nonce reuse leaks the private key. The system has multiple protections:
1. Database unique constraints prevent storage of duplicate nonces
2. Application-level validation before signing
3. Error handling stops signing if duplicate detected

Never override these protections. See [Security Documentation](../security/musig2_security.md) for details.

#### Q: How do I verify a MuSig2 transaction before broadcast?

**A**:
```bash
# 1. Decode PSBT
watch btc api decodepsbt $(cat payment_15_signed_3.psbt)

# 2. Verify signature present
watch btc api analyzepsbt $(cat payment_15_signed_3.psbt)

# 3. Extract and verify transaction
watch btc api finalizepsbt $(cat payment_15_signed_3.psbt) true

# 4. Verify transaction validity (doesn't broadcast)
watch btc api testmempoolaccept '["<tx_hex>"]'
```

#### Q: Can I recover if I lose a PSBT file during signing?

**A**:
- **During Round 1 (nonce generation)**: Regenerate nonces. The Watch wallet recreates the PSBT from the transaction.
- **During Round 2 (signing)**: More complex. You may need to restart the entire signing process.
- **Best Practice**: Always backup PSBT files after each step.

### Operational Questions

#### Q: How do I train my team on MuSig2?

**A**: Recommended approach:
1. **Read Documentation**: All operators read `docs/crypto/btc/musig2_guide.md` and `docs/security/musig2_security.md`
2. **Testnet Practice**: Each operator practices full workflow on testnet (5-10 transactions)
3. **Shadow Production**: Observe experienced operator in production
4. **Supervised Production**: Perform under supervision
5. **Independent**: Solo operations after 10 successful supervised transactions

#### Q: What monitoring should I set up?

**A**:
1. **Nonce Tracking**: Alert on any nonce reuse attempts
2. **Transaction Success Rate**: Alert if success rate drops below 95%
3. **PSBT File Management**: Alert on missing or stale files
4. **Signature Verification**: Alert on signature verification failures
5. **Fee Estimation**: Monitor actual fee savings vs expected

#### Q: How do I handle emergency situations?

**A**: See [Rollback Procedures](#rollback-procedures). Key points:
- Always complete in-flight transactions first
- Document the issue thoroughly
- Don't panic-rollback for operator errors (training issue)
- Do rollback for security issues or critical bugs

#### Q: Can I automate the signing workflow?

**A**: **Partial automation possible**, but:
- ✅ Nonce generation can be automated (Round 1)
- ✅ Partial signature creation can be automated (Round 2)
- ❌ **Never automate**: Private key access, nonce reuse checks, final broadcast approval
- **Recommendation**: Keep human oversight at critical points

### Migration-Specific Questions

#### Q: Should I migrate all accounts at once?

**A**: **No**. Start with low-risk accounts (e.g., deposit), then gradually expand. This allows you to:
- Build team confidence
- Identify issues in low-risk environment
- Adjust procedures based on lessons learned

#### Q: What if I discover an issue after sweeping funds?

**A**:
1. **Stop**: Don't create more MuSig2 transactions
2. **Assess**: Determine if funds are at risk
3. **Secure**: If funds are safe, plan fix; if at risk, execute emergency rollback
4. **Fix**: Address root cause
5. **Validate**: Test fix thoroughly
6. **Resume**: Continue migration or complete rollback

#### Q: How do I calculate break-even point for migration?

**A**:
```bash
# Migration costs:
SETUP_TIME=40        # hours (team time)
HOURLY_RATE=100      # USD (team hourly cost)
SWEEP_FEES=50000     # sats (estimated fees to sweep old UTXOs)
BTC_PRICE_USD=40000  # Example BTC price in USD

TOTAL_COST=$(echo "scale=2; $SETUP_TIME * $HOURLY_RATE + ($SWEEP_FEES / 100000000) * $BTC_PRICE_USD" | bc -l)

# Monthly savings:
TX_PER_MONTH=100
FEE_RATE=50          # sat/vB
P2WSH_SIZE=385       # bytes
MUSIG2_SIZE=225      # bytes
SAVINGS_PER_TX=$(echo "($P2WSH_SIZE - $MUSIG2_SIZE) * $FEE_RATE" | bc)
MONTHLY_SAVINGS=$(echo "scale=2; ($SAVINGS_PER_TX * $TX_PER_MONTH / 100000000) * $BTC_PRICE_USD" | bc -l)

# Break-even time:
if [ "$(echo "$MONTHLY_SAVINGS > 0" | bc -l)" -eq 1 ]; then
    BREAK_EVEN_MONTHS=$(echo "scale=2; $TOTAL_COST / $MONTHLY_SAVINGS" | bc -l)
else
    BREAK_EVEN_MONTHS="N/A"
fi

echo "Break-even in $BREAK_EVEN_MONTHS months"

# For 100 tx/month at 50 sat/vB, BTC at $40,000:
# Monthly savings: ~$128 USD
# Break-even: 3-5 months
```

#### Q: Can I migrate without downtime?

**A**: **Yes**. The gradual migration approach allows zero-downtime migration:
1. Create MuSig2 addresses alongside P2WSH
2. Gradually shift new transactions to MuSig2
3. P2WSH remains fully operational during migration
4. Sweep funds during scheduled maintenance or low-activity periods

#### Q: What if Bitcoin Core introduces MuSig2 wallet support?

**A**: Future Bitcoin Core versions may add native MuSig2 wallet support. When that happens:
- Your addresses remain valid (standard P2TR)
- Your keys remain valid
- You may be able to import keys into Bitcoin Core wallet
- **Recommendation**: Monitor Bitcoin Core release notes for MuSig2 wallet features

---

## Conclusion

This migration guide provides a comprehensive strategy for transitioning from traditional P2WSH multisig to MuSig2 Taproot addresses. Key takeaways:

### Success Factors

1. **Thorough Testing**: Test extensively on testnet before production
2. **Gradual Rollout**: Start small, build confidence, expand gradually
3. **Team Training**: Ensure all operators understand MuSig2 security requirements
4. **Monitoring**: Set up comprehensive monitoring before migration
5. **Rollback Plan**: Always have a tested rollback procedure ready

### Expected Outcomes

After successful migration, you'll achieve:
- **40-50% fee savings** on multisig transactions
- **Improved privacy** (transactions indistinguishable from single-sig)
- **Modern infrastructure** using latest Bitcoin standards
- **Future-proof** wallet architecture

### Next Steps

1. Review [Security Documentation](../security/musig2_security.md)
2. Read [MuSig2 User Guide](../crypto/btc/musig2_guide.md)
3. Create your migration plan using this guide as template
4. Schedule testnet testing
5. Execute migration following documented procedures

### Support and Resources

- **Documentation**: `docs/crypto/btc/musig2_guide.md`
- **Security**: `docs/security/musig2_security.md`
- **Architecture**: `docs/architecture/musig2_architecture.md`
- **Issues**: GitHub Issues - Report bugs or ask questions
- **Bitcoin Core**: https://bitcoincore.org/ (v22.0+ required)

---

**Document Version**: 1.0
**Last Updated**: 2025-01-30
**Author**: go-crypto-wallet Team
**Related Issues**: #141, #171
