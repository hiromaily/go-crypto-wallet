## Best Practices

### Security

1. **Nonce Management** (Critical)
   - ✅ Generate fresh nonces for every transaction
   - ✅ Store nonces securely during Round 1-2
   - ✅ Delete nonces immediately after signing
   - ❌ **NEVER reuse nonces** (will leak private key)
   - ❌ Never share secret nonces with anyone

2. **Air-Gapped Signing**
   - ✅ Keep Keygen and Sign wallets offline at all times
   - ✅ Use dedicated, isolated computers for offline wallets
   - ✅ Never connect offline wallets to networks

3. **File Transfer Security**
   - ✅ Use dedicated USB drives for PSBT transfer
   - ✅ Virus scan USB drives before use
   - ✅ Verify file checksums after transfer

4. **Transaction Verification**
   - ✅ Verify transaction amounts before signing
   - ✅ Check recipient addresses carefully
   - ✅ Verify aggregated signature before broadcasting

### Operations

1. **Workflow Management**
   - ✅ Track PSBT state at each step
   - ✅ Use file checksums for integrity verification
   - ✅ Document workflow for each transaction type
   - ✅ Create checklists for Round 1 and Round 2

2. **Parallel vs Sequential Operations**
   - ✅ **Round 1** (nonce generation): Can be done **in parallel**
   - ✅ **Round 2** (signing): Must be done **sequentially** after nonces collected
   - ✅ Collect all nonces before proceeding to Round 2

3. **Testing**
   - ✅ Test on testnet/signet before mainnet
   - ✅ Test with small amounts first
   - ✅ Verify full workflow end-to-end
   - ✅ Test nonce uniqueness enforcement

4. **Monitoring**
   - ✅ Monitor transaction confirmations
   - ✅ Track fee savings compared to traditional multisig
   - ✅ Set up alerts for transaction failures

### Performance

1. **Fee Optimization**
   - ✅ MuSig2 transactions are 30-50% smaller
   - ✅ Fees are proportionally lower
   - ✅ Monitor mempool for optimal fee rates
   - ✅ Consider consolidating UTXOs during low fee periods

2. **Transaction Batching**
   - ✅ Batch multiple payments into single transaction
   - ✅ Further reduces fees per payment
   - ✅ Increases privacy (harder to link payments)

### Backup and Recovery

1. **Regular Backups**
   - ✅ Backup seeds and private keys securely
   - ✅ Backup wallet databases regularly
   - ✅ Store backups in multiple secure locations
   - ✅ Test recovery procedures

2. **Key Management**
   - ✅ Use BIP39 mnemonics for seed backup
   - ✅ Store seed backups offline in secure locations
   - ✅ Consider hardware security modules (HSMs) for production

---
