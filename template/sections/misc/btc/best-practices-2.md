## Best Practices

### Security

> For general offline wallet security (air-gapped signing, file transfer security), see [MuSig2 Best Practices](../../../../docs/chains/btc/musig2/user-guide.md#best-practices).

1. **Private Key Protection**
   - ✅ Store seeds and private keys in secure offline storage
   - ✅ Use hardware security modules (HSMs) for production
   - ✅ Implement proper access controls

2. **Transaction Verification**
   - ✅ Always verify transaction amounts before signing
   - ✅ Check recipient addresses carefully
   - ✅ Verify fees are reasonable

### Operations

1. **File Naming**
   - ✅ Use auto-generated filenames (don't rename PSBT files)
   - ✅ Keep files organized by transaction type
   - ✅ Archive old PSBTs regularly

2. **Workflow Documentation**
   - ✅ Document your signing workflow
   - ✅ Create checklists for each transaction type
   - ✅ Train operators on PSBT procedures

3. **Testing**
   - ✅ Test on testnet before mainnet
   - ✅ Test with small amounts first
   - ✅ Verify full workflow end-to-end

4. **Monitoring**
   - ✅ Monitor transaction confirmations
   - ✅ Track fee rates and adjust as needed
   - ✅ Set up alerts for transaction failures

### Performance

1. **Fee Management**
   - ✅ Monitor mempool for optimal fee rates
   - ✅ Use higher fees for urgent transactions
   - ✅ Consider consolidating UTXOs during low fee periods

2. **UTXO Management**
   - ✅ Avoid creating dust outputs
   - ✅ Consolidate UTXOs when fees are low
   - ✅ Monitor UTXO set size

3. **Transaction Batching**
   - ✅ Batch multiple payments into single transaction
   - ✅ Reduces overall fees
   - ✅ Increases efficiency

### Backup and Recovery

1. **Regular Backups**
   - ✅ Backup seeds and private keys securely
   - ✅ Backup wallet databases regularly
   - ✅ Test recovery procedures

2. **Disaster Recovery**
   - ✅ Document recovery procedures
   - ✅ Store backups in multiple secure locations
   - ✅ Test recovery periodically

3. **Key Management**
   - ✅ Use BIP39 mnemonics for seed backup
   - ✅ Store seed backups in secure, offline locations
   - ✅ Consider multi-signature setup for critical keys

---
