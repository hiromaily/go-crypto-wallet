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
