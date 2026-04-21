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
