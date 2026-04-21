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
