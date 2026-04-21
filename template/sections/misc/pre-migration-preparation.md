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
cp -r config/wallet/ backups/pre-psbt-migration-$(date +%Y%m%d)/config/

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
cp config/wallet/btc/watch.yaml config/wallet/btc/watch_testnet.yaml

# Update RPC endpoint to testnet node
vim config/wallet/btc/watch_testnet.yaml
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
