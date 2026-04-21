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
