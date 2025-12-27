# PSBT Migration Guide

This guide provides step-by-step instructions for migrating from the legacy CSV transaction format to the PSBT (Partially Signed Bitcoin Transaction) format in go-crypto-wallet.

## Table of Contents

1. [Overview](#overview)
2. [Migration Timeline](#migration-timeline)
3. [Pre-Migration Preparation](#pre-migration-preparation)
4. [Migration Process](#migration-process)
5. [Post-Migration Verification](#post-migration-verification)
6. [Rollback Procedure](#rollback-procedure)
7. [FAQ](#faq)

---

## Overview

### What Changed?

The go-crypto-wallet system has migrated from a custom CSV-based transaction format to the industry-standard PSBT format (BIP 174).

**Key Changes:**

| Aspect | Before (CSV) | After (PSBT) |
|--------|-------------|--------------|
| **File Format** | Custom CSV | Base64 PSBT |
| **Extension** | Various | `.psbt` |
| **Encoding** | Plain text | Base64 |
| **Standard** | Custom | BIP 174 |
| **Compatibility** | go-crypto-wallet only | Bitcoin Core, hardware wallets |
| **Validation** | Basic | Comprehensive |
| **Metadata** | Limited | Rich (UTXO info, scripts) |

### Why Migrate?

**Benefits of PSBT:**

1. **Standardization**
   - Industry-standard format (BIP 174)
   - Compatible with Bitcoin Core, Electrum, hardware wallets
   - Better interoperability

2. **Security**
   - Structured data format prevents parsing errors
   - Built-in validation
   - Better error detection

3. **Functionality**
   - Richer transaction metadata
   - Support for complex scripts
   - Better debugging capabilities

4. **Future-Proofing**
   - Foundation for hardware wallet integration
   - Enables Taproot multisig features
   - Supports MuSig2 (future)

### Migration Scope

**Affected Components:**

- ✅ Watch Wallet - Transaction creation and broadcasting
- ✅ Keygen Wallet - Transaction signing (first signature)
- ✅ Sign Wallet - Transaction signing (second signature)
- ✅ Transaction file storage
- ✅ Transaction file naming convention

**Not Affected:**

- ❌ Database schema (no changes)
- ❌ Key generation process
- ❌ Address generation
- ❌ Configuration files
- ❌ Command-line interface (same commands)

---

## Migration Timeline

### Recommended Timeline

```
Week 1-2: Preparation
├── Review documentation
├── Test on testnet
├── Backup production data
└── Schedule maintenance window

Week 3: Migration
├── Deploy PSBT-enabled binaries
├── Complete pending CSV transactions
├── Test PSBT workflow
└── Resume operations with PSBT

Week 4+: Post-Migration
├── Monitor operations
├── Archive CSV files
├── Update procedures
└── Close migration
```

### Critical Milestones

| Milestone | Timeline | Status |
|-----------|----------|--------|
| **Preparation Complete** | T-2 weeks | [ ] |
| **Testnet Validation** | T-1 week | [ ] |
| **Maintenance Window** | T-0 (planned date) | [ ] |
| **Migration Complete** | T+1 day | [ ] |
| **Production Validation** | T+1 week | [ ] |
| **CSV Archive** | T+1 month | [ ] |

---

## Pre-Migration Preparation

### Step 1: Review Current State

#### 1.1 Check for Pending Transactions

```bash
# On Watch wallet
cd /path/to/go-crypto-wallet

# Check for unsigned CSV files
ls -la data/tx/btc/*unsigned*

# Check for partially signed CSV files (non-PSBT files)
find data/tx/btc/ -type f ! -name "*.psbt"

# Query database for pending transactions
sqlite3 data/db/btc_watch.db "SELECT id, tx_type, action FROM btc_tx WHERE tx_type != 'sent';"
```

**Action Required:**
- Complete all pending CSV transactions before migration
- Or convert to PSBT format (see conversion section)

#### 1.2 Verify Binary Versions

```bash
# Check current wallet versions
./watch version
./keygen version
./sign version

# Expected output after migration:
# Version: vX.X.X (with PSBT support)
# PSBT Support: Enabled
```

#### 1.3 Backup Current State

```bash
# Create backup directory
mkdir -p backups/pre-psbt-migration-$(date +%Y%m%d)

# Backup databases
cp -r data/db/ backups/pre-psbt-migration-$(date +%Y%m%d)/db/

# Backup transaction files
cp -r data/tx/ backups/pre-psbt-migration-$(date +%Y%m%d)/tx/

# Backup configuration
cp -r data/config/ backups/pre-psbt-migration-$(date +%Y%m%d)/config/

# Create backup archive
tar -czf psbt-migration-backup-$(date +%Y%m%d).tar.gz \
    backups/pre-psbt-migration-$(date +%Y%m%d)/

# Verify backup
tar -tzf psbt-migration-backup-$(date +%Y%m%d).tar.gz | head -20
```

**Critical:** Store backup in secure, offline location.

### Step 2: Test on Testnet

#### 2.1 Setup Testnet Environment

```bash
# Clone production config for testnet
cp data/config/btc_watch.toml data/config/btc_watch_testnet.toml

# Update RPC endpoint to testnet node
vim data/config/btc_watch_testnet.toml
# Change:
# host = "mainnet-node:8332"
# To:
# host = "testnet-node:18332"
```

#### 2.2 Test Complete PSBT Workflow

**Test Case 1: Deposit Transaction (Single-Sig)**

```bash
# 1. Create unsigned PSBT
./watch create deposit --fee 0.00001 --config testnet

# Expected: deposit_X_unsigned_0_*.psbt file created

# 2. Sign with Keygen
./keygen sign --file data/tx/btc/deposit_*_unsigned_0_*.psbt --config testnet

# Expected: deposit_X_signed_1_*.psbt file created

# 3. Broadcast
./watch send --file data/tx/btc/deposit_*_signed_1_*.psbt --config testnet

# Expected: Transaction hash returned
```

**Test Case 2: Payment Transaction (Multisig 2-of-2)**

```bash
# 1. Create unsigned PSBT
./watch create payment --fee 0.00001 --config testnet

# 2. First signature (Keygen)
./keygen sign --file data/tx/btc/payment_*_unsigned_0_*.psbt --config testnet

# Expected: payment_X_unsigned_1_*.psbt (partially signed)

# 3. Second signature (Sign)
./sign sign --file data/tx/btc/payment_*_unsigned_1_*.psbt --config testnet

# Expected: payment_X_signed_2_*.psbt (fully signed)

# 4. Broadcast
./watch send --file data/tx/btc/payment_*_signed_2_*.psbt --config testnet

# Expected: Transaction hash returned
```

#### 2.3 Verify Test Results

**Checklist:**
- [ ] Unsigned PSBT created successfully
- [ ] First signature completed (Keygen)
- [ ] Second signature completed (Sign, for multisig)
- [ ] PSBT finalized successfully
- [ ] Transaction broadcast successfully
- [ ] Transaction confirmed on testnet blockchain
- [ ] Database updated correctly

### Step 3: Prepare Production Environment

#### 3.1 Schedule Maintenance Window

**Recommended Window:**
- **Duration**: 2-4 hours
- **Timing**: Off-peak hours
- **Communication**: Notify stakeholders 1 week in advance

**Maintenance Window Checklist:**
- [ ] Schedule announced to stakeholders
- [ ] Backup procedures verified
- [ ] Rollback plan documented
- [ ] Support team on standby
- [ ] Monitoring enabled

#### 3.2 Prepare Deployment Package

```bash
# Build PSBT-enabled binaries
make build

# Verify PSBT support
./build/watch version | grep "PSBT Support"
./build/keygen version | grep "PSBT Support"
./build/sign version | grep "PSBT Support"

# Create deployment package
mkdir -p deploy/psbt-migration
cp build/watch deploy/psbt-migration/
cp build/keygen deploy/psbt-migration/
cp build/sign deploy/psbt-migration/

# Create checksums
cd deploy/psbt-migration
sha256sum * > SHA256SUMS

# Package for distribution
cd ..
tar -czf psbt-migration-binaries-$(date +%Y%m%d).tar.gz psbt-migration/
```

#### 3.3 Update Documentation

- [ ] Update operational procedures
- [ ] Update runbooks
- [ ] Update training materials
- [ ] Update monitoring dashboards

---

## Migration Process

### Phase 1: Complete Pending Transactions

#### Step 1: Process All CSV Transactions

```bash
# List all pending transactions
sqlite3 data/db/btc_watch.db \
  "SELECT id, action, tx_type, created_at FROM btc_tx WHERE tx_type != 'sent' ORDER BY id;"

# Complete each pending transaction using CSV format
# (Do not create new transactions during this phase)
```

**Options:**

**Option A: Complete CSV Transactions**
- Finish all pending CSV transactions normally
- Quickest, safest option
- Recommended if pending count is low (<10)

**Option B: Convert to PSBT** (Advanced)
- Convert pending CSV to PSBT format
- Requires custom conversion script
- Only if many pending transactions

#### Step 2: Verify No Pending Transactions

```bash
# Check database
sqlite3 data/db/btc_watch.db \
  "SELECT COUNT(*) as pending FROM btc_tx WHERE tx_type != 'sent';"

# Expected output: 0
```

### Phase 2: Deploy PSBT Binaries

#### Step 1: Stop Wallet Operations

```bash
# Stop any running wallet processes
pkill -f "./watch"
pkill -f "./keygen"
pkill -f "./sign"

# Verify all processes stopped
ps aux | grep -E "(watch|keygen|sign)" | grep -v grep
```

#### Step 2: Backup Current Binaries

```bash
# Backup existing binaries
mkdir -p backups/binaries-csv-$(date +%Y%m%d)
cp watch backups/binaries-csv-$(date +%Y%m%d)/
cp keygen backups/binaries-csv-$(date +%Y%m%d)/
cp sign backups/binaries-csv-$(date +%Y%m%d)/
```

#### Step 3: Deploy PSBT Binaries

```bash
# Extract deployment package
tar -xzf psbt-migration-binaries-$(date +%Y%m%d).tar.gz

# Verify checksums
cd psbt-migration
sha256sum -c SHA256SUMS

# Deploy binaries
cp watch /path/to/production/
cp keygen /path/to/production/
cp sign /path/to/production/

# Verify permissions
chmod +x /path/to/production/watch
chmod +x /path/to/production/keygen
chmod +x /path/to/production/sign
```

#### Step 4: Verify Deployment

```bash
# Verify PSBT support enabled
/path/to/production/watch version
# Expected: "PSBT Support: Enabled"

/path/to/production/keygen version
# Expected: "PSBT Support: Enabled"

/path/to/production/sign version
# Expected: "PSBT Support: Enabled"
```

### Phase 3: Test PSBT Workflow

#### Step 1: Create Test Transaction (Small Amount)

```bash
# Create small deposit transaction
./watch create deposit --amount 0.001 --fee 0.00001

# Verify PSBT file created
ls -la data/tx/btc/deposit_*_unsigned_0_*.psbt
```

#### Step 2: Complete Test Transaction

```bash
# Sign with Keygen
./keygen sign --file data/tx/btc/deposit_*_unsigned_0_*.psbt

# Verify signed PSBT created
ls -la data/tx/btc/deposit_*_signed_1_*.psbt

# Broadcast (WARNING: This will send real Bitcoin)
./watch send --file data/tx/btc/deposit_*_signed_1_*.psbt

# Verify transaction broadcast
# Check Bitcoin block explorer for transaction hash
```

#### Step 3: Verify Database Updates

```bash
# Check transaction status in database
sqlite3 data/db/btc_watch.db \
  "SELECT id, action, tx_type, sent_tx_hash FROM btc_tx ORDER BY id DESC LIMIT 1;"

# Expected: tx_type = 'sent', sent_tx_hash populated
```

### Phase 4: Resume Operations

#### Step 1: Enable Production Operations

```bash
# Resume normal transaction operations
# Monitor first few transactions closely
```

#### Step 2: Monitor Metrics

Monitor the following during first 24-48 hours:

- **Transaction Success Rate**
  ```bash
  # Check successful broadcasts
  sqlite3 data/db/btc_watch.db \
    "SELECT COUNT(*) as success FROM btc_tx WHERE tx_type = 'sent' AND created_at > datetime('now', '-24 hours');"
  ```

- **Transaction Failures**
  ```bash
  # Check for failures (should be 0)
  grep -i "error" logs/watch.log | grep -i "psbt" | tail -20
  ```

- **File Format**
  ```bash
  # Verify all new files are PSBT
  ls -la data/tx/btc/*.psbt | tail -10
  ```

- **Signing Success**
  ```bash
  # Check signing operations
  grep -i "signing completed" logs/keygen.log | tail -10
  grep -i "signing completed" logs/sign.log | tail -10
  ```

---

## Post-Migration Verification

### Verification Checklist

#### Day 1: Immediate Verification

- [ ] Test transaction created successfully (PSBT format)
- [ ] Test transaction signed successfully (Keygen)
- [ ] Test transaction signed successfully (Sign, if multisig)
- [ ] Test transaction broadcast successfully
- [ ] Transaction confirmed on blockchain
- [ ] Database updated correctly
- [ ] No errors in logs
- [ ] Monitoring shows normal metrics

#### Week 1: Short-Term Monitoring

- [ ] All transaction types tested (deposit, payment, transfer)
- [ ] All address types working (P2PKH, P2WPKH, P2TR, etc.)
- [ ] Multisig transactions completing successfully
- [ ] No PSBT-related errors
- [ ] Performance metrics stable
- [ ] Operators comfortable with PSBT workflow

#### Month 1: Long-Term Validation

- [ ] CSV files archived
- [ ] Documentation updated
- [ ] Training completed
- [ ] No rollback needed
- [ ] Migration officially closed

### Archive Legacy CSV Files

After successful migration (recommend waiting 1 month):

```bash
# Create archive directory
mkdir -p archive/csv-legacy-$(date +%Y%m%d)

# Move old CSV files (non-PSBT)
find data/tx/btc/ -type f ! -name "*.psbt" -exec mv {} archive/csv-legacy-$(date +%Y%m%d)/ \;

# Compress archive
tar -czf csv-legacy-$(date +%Y%m%d).tar.gz archive/csv-legacy-$(date +%Y%m%d)/

# Move to long-term storage
mv csv-legacy-$(date +%Y%m%d).tar.gz /secure/archive/location/

# Optional: Remove local archive after verification
# rm -rf archive/csv-legacy-$(date +%Y%m%d)/
```

---

## Rollback Procedure

### When to Rollback

Rollback if you encounter:

- ❌ Critical errors preventing transaction creation
- ❌ Signing failures
- ❌ Broadcasting failures
- ❌ Data corruption
- ❌ Unexpected behavior

### Rollback Steps

#### Step 1: Stop Operations Immediately

```bash
# Stop all wallet processes
pkill -f "./watch"
pkill -f "./keygen"
pkill -f "./sign"

# Prevent new transactions
# (Communicate to operators)
```

#### Step 2: Restore CSV Binaries

```bash
# Restore previous binaries
cp backups/binaries-csv-$(date +%Y%m%d)/watch ./
cp backups/binaries-csv-$(date +%Y%m%d)/keygen ./
cp backups/binaries-csv-$(date +%Y%m%d)/sign ./

# Verify restoration
./watch version
# Expected: CSV-based version
```

#### Step 3: Verify Database Integrity

```bash
# Check database consistency
sqlite3 data/db/btc_watch.db "PRAGMA integrity_check;"

# If corrupted, restore from backup
cp backups/pre-psbt-migration-$(date +%Y%m%d)/db/btc_watch.db data/db/
```

#### Step 4: Resume CSV Operations

```bash
# Test CSV transaction creation
./watch create deposit --fee 0.0001

# Verify CSV file created (not PSBT)
ls -la data/tx/btc/deposit_*

# Complete test transaction
# Verify end-to-end flow working
```

#### Step 5: Investigate Issues

```bash
# Collect logs
cp logs/watch.log investigation/
cp logs/keygen.log investigation/
cp logs/sign.log investigation/

# Collect failed PSBT files
cp data/tx/btc/*.psbt investigation/

# Review error messages
grep -i "error" logs/*.log > investigation/errors.txt
```

#### Step 6: Plan Re-Migration

After fixing issues:

1. Identify root cause
2. Apply fixes
3. Test on testnet again
4. Schedule new migration window
5. Retry migration with corrected procedures

---

## FAQ

### General Questions

**Q: Will the migration cause downtime?**

A: Yes, a brief maintenance window (2-4 hours) is required to complete pending transactions and deploy new binaries. However, the system will be offline only during binary deployment (15-30 minutes).

**Q: Do I need to regenerate keys or addresses?**

A: No, keys and addresses are not affected by PSBT migration. Only the transaction file format changes.

**Q: Can I use PSBT and CSV simultaneously?**

A: No, the PSBT-enabled version removes CSV support. You must complete all CSV transactions before migration.

**Q: What happens to existing transactions after migration?**

A: Completed transactions (already broadcast) are not affected. Only pending transactions need to be completed or converted.

### Technical Questions

**Q: Are PSBT files larger than CSV files?**

A: Yes, PSBT files contain more metadata (UTXO info, scripts, derivation paths). However, files are base64-encoded and compressed, so the size difference is minimal (typically 20-30% larger).

**Q: Can I inspect PSBT contents?**

A: Yes, use Bitcoin Core:
```bash
bitcoin-cli decodepsbt "$(cat transaction.psbt)"
```

**Q: Are PSBTs compatible with Bitcoin Core?**

A: Yes, PSBT is a Bitcoin standard (BIP 174) fully supported by Bitcoin Core v0.17+.

**Q: Can I convert CSV files to PSBT?**

A: Technically possible but complex. Recommended approach is to complete CSV transactions before migration.

**Q: Does PSBT support all address types?**

A: Yes, PSBT supports all Bitcoin address types: P2PKH, P2SH, P2WPKH, P2WSH, and P2TR (Taproot).

### Operational Questions

**Q: How do I transfer PSBT files between wallets?**

A: Same as CSV files - use USB drives, QR codes, or secure file transfer. PSBT files have `.psbt` extension.

**Q: Do commands change with PSBT?**

A: No, commands remain the same. The format change is transparent to operators:
```bash
# Same commands work with PSBT
./watch create deposit --fee 0.0001
./keygen sign --file transaction.psbt
./watch send --file transaction.psbt
```

**Q: How do I know if a PSBT is fully signed?**

A: Check the filename:
- `_unsigned_0_*.psbt` - No signatures
- `_unsigned_1_*.psbt` - Partially signed (1 signature)
- `_signed_2_*.psbt` - Fully signed (2 signatures)

**Q: What if I lose a PSBT file?**

A: PSBT files can be recreated from database transaction records. Contact support for recovery procedures.

### Troubleshooting

**Q: Error: "PSBT is not fully signed"**

A: The PSBT needs more signatures. For multisig:
1. Check current signature count in filename
2. Send to next signer (Keygen → Sign)
3. Verify all required signatures collected

**Q: Error: "Invalid PSBT format"**

A: PSBT file may be corrupted:
1. Verify file integrity (checksum)
2. Don't edit PSBT files manually
3. Re-create transaction if needed

**Q: Error: "Failed to read PSBT file"**

A: Check file extension:
```bash
# Must have .psbt extension
mv transaction transaction.psbt
```

**Q: Can I edit PSBT files?**

A: No, never edit PSBT files manually. They are base64-encoded binary data. Use wallet commands to create and sign PSBTs.

---

## Migration Support

### Pre-Migration Support

**Documentation:**
- [PSBT User Guide](psbt_user_guide.md)
- [PSBT Developer Guide](psbt_developer_guide.md)
- [PSBT Implementation](psbt_implementation.md)

**Testing:**
- Test on Bitcoin testnet before mainnet
- Test with small amounts first
- Test complete workflow end-to-end

### Migration Assistance

**Resources:**
- GitHub Issues: [https://github.com/hiromaily/go-crypto-wallet/issues](https://github.com/hiromaily/go-crypto-wallet/issues)
- Documentation: `docs/crypto/btc/`
- Example Scripts: `scripts/examples/`

**Emergency Contact:**
- Critical issues during migration
- Rollback assistance
- Data recovery

---

## Appendices

### Appendix A: Migration Checklist

```
Pre-Migration:
[ ] Backup databases
[ ] Backup transaction files
[ ] Backup configuration files
[ ] Complete pending CSV transactions
[ ] Test on testnet
[ ] Schedule maintenance window
[ ] Notify stakeholders
[ ] Prepare rollback plan

Migration:
[ ] Stop wallet operations
[ ] Verify no pending transactions
[ ] Backup current binaries
[ ] Deploy PSBT binaries
[ ] Verify PSBT support enabled
[ ] Create test PSBT transaction
[ ] Complete test transaction
[ ] Verify database updates
[ ] Resume operations

Post-Migration:
[ ] Monitor for 24 hours
[ ] Test all transaction types
[ ] Verify all address types
[ ] Check logs for errors
[ ] Update documentation
[ ] Train operators
[ ] Archive CSV files (after 1 month)
[ ] Close migration
```

### Appendix B: Conversion Script (CSV to PSBT)

Note: This is an advanced procedure. Recommended only if you have many pending transactions and completing them manually is impractical.

```bash
#!/bin/bash
# csv_to_psbt_converter.sh
# WARNING: Test thoroughly before using in production

for csv_file in data/tx/btc/*_unsigned_*.csv; do
    echo "Converting: $csv_file"

    # Extract transaction data from CSV
    tx_data=$(cat "$csv_file")

    # Create new transaction in PSBT format
    # (This requires custom implementation based on CSV structure)
    ./watch create-from-csv --csv "$csv_file" --output psbt

    echo "Created: ${csv_file%.csv}.psbt"
done
```

**Note:** The actual conversion logic depends on your CSV structure and is not provided in this template.

### Appendix C: Monitoring Queries

```bash
# Transaction success rate (last 24 hours)
sqlite3 data/db/btc_watch.db <<EOF
SELECT
    COUNT(*) as total_transactions,
    SUM(CASE WHEN tx_type = 'sent' THEN 1 ELSE 0 END) as successful,
    ROUND(100.0 * SUM(CASE WHEN tx_type = 'sent' THEN 1 ELSE 0 END) / COUNT(*), 2) as success_rate
FROM btc_tx
WHERE created_at > datetime('now', '-24 hours');
EOF

# PSBT file count by status
find data/tx/btc -name "*.psbt" | \
    awk -F'_' '{print $(NF-2)"_"$(NF-1)}' | \
    sort | uniq -c

# Recent errors in logs
tail -100 logs/watch.log | grep -i "error" | grep -i "psbt"
```

---

**Last Updated**: 2025-01-27
**Version**: 1.0 (PSBT Phase 2 Complete)
