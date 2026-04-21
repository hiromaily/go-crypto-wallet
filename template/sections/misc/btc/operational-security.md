## Operational Security

### Common Pitfalls

#### Pitfall #1: Nonce Reuse

**Description**: Using the same nonce to sign two different transactions.

**Causes**:

- Database failure (nonce not recorded properly)
- Operator error (re-running Round 1 for same transaction)
- Software bug (nonce not generated uniquely)
- File management error (using old PSBT with same nonce)

**Consequences**:

```
CRITICAL: Private key leakage
→ Immediate fund loss risk
→ All funds accessible with that key are at risk
```

**Prevention**:

1. Database unique constraints (technical control)
2. Application validation (technical control)
3. Operator training (human control)
4. File naming conventions (process control)
5. Monitoring (detective control)

**Detection**:

```bash
# Check for nonce reuse
mysql> SELECT nonce, COUNT(*) FROM musig2_nonces GROUP BY nonce HAVING COUNT(*) > 1;
```

**Response**: See [Incident Response - Nonce Reuse](#nonce-reuse-incident)

#### Pitfall #2: Wrong Signer Order

**Description**: Signers use different ordering of public keys when aggregating.

**Causes**:

- Configuration mismatch between wallets
- Manual key entry (typos, wrong order)
- Software version mismatch

**Consequences**:

```
Result: Different aggregated public keys
→ Signature verification fails
→ Funds remain safe, but transaction unusable
→ Need to restart signing process
```

**Prevention**:

1. Use consistent key sorting (sort alphabetically by public key)
2. Configuration management (all wallets use same config)
3. Automated key import (no manual entry)

**Detection**:

```go
// Verify all signers compute same aggregated key

aggregatedKey1 := keygen.GetAggregatedKey()
aggregatedKey2 := sign1.GetAggregatedKey()

if !bytes.Equal(aggregatedKey1, aggregatedKey2) {
    return errors.New("aggregated key mismatch - check key order")
}
```

**Response**:

- Fix key ordering configuration
- Restart signing process with correct order

#### Pitfall #3: Missing Partial Signature

**Description**: Not all signers provide partial signatures.

**Causes**:

- Operator error (forgot to sign)
- File transfer failure (partial signature lost)
- Hardware failure (signing device offline)

**Consequences**:

```
Result: Incomplete signature set
→ Cannot aggregate final signature
→ Transaction cannot be broadcast
→ Funds remain safe, transaction pending
```

**Prevention**:

1. Tracking system (checklist for each transaction)
2. Automated validation (check signature count before aggregation)
3. Redundant file backups

**Detection**:

```go
// Before aggregation, verify all signatures present

func validatePartialSignatures(sigs []PartialSignature, expectedCount int) error {
    if len(sigs) < expectedCount {
        return fmt.Errorf("missing partial signatures: got %d, expected %d",
            len(sigs), expectedCount)
    }
    return nil
}
```

**Response**:

- Identify which signer is missing
- Have that signer complete their partial signature
- Continue aggregation

#### Pitfall #4: File Management Errors

**Description**: Using wrong PSBT file, overwriting files, losing files.

**Causes**:

- Confusing file names
- Manual file management (human error)
- Lack of version control for files

**Consequences**:

```
Result varies:
- Using old PSBT → may include wrong transaction data
- Overwriting file → loss of signatures, need to restart
- Losing file → may need to recreate transaction
```

**Prevention**:

1. Strict file naming convention:

   ```
   payment_{request_id}_{stage}_{step}_{timestamp}.psbt

   Example:
   payment_15_unsigned_0_{timestamp}.psbt   (initial PSBT)
   payment_15_nonce_1_{timestamp}.psbt      (after keygen nonce)
   payment_15_nonce_2_{timestamp}.psbt      (after sign1 nonce)
   payment_15_nonce_3_{timestamp}.psbt      (all nonces collected)
   payment_15_signed_1_{timestamp}.psbt     (after keygen signature)
   payment_15_signed_2_{timestamp}.psbt     (after sign1 signature)
   payment_15_signed_3_{timestamp}.psbt     (all partial signatures)
   payment_15_final_{timestamp}.psbt        (final, ready to broadcast)
   ```

2. File checksums (detect corruption)

   ```bash
   sha256sum payment_15_unsigned_0_1704067200.psbt > payment_15_unsigned_0_1704067200.psbt.sha256
   ```

3. Automated file management (scripts handle naming)

4. File backup (keep copies of all intermediate files)

**Detection**:

- Checksum verification before use
- PSBT analysis (check expected contents)

**Response**:

- If wrong file used: Restart from last correct checkpoint
- If file corrupted: Restore from backup or recreate

#### Pitfall #5: Network Synchronization Issues

**Description**: Offline wallets have stale data.

**Causes**:

- Offline wallets not updated with latest public keys
- Configuration out of sync between wallets
- Key rotation not communicated

**Consequences**:

```
Result: Signing failures
→ Different aggregated keys computed
→ Signature verification fails
→ Need to re-sync and restart
```

**Prevention**:

1. Regular sync schedule (weekly)
2. Manifest files tracking current keys/config
3. Version numbers in configuration files
4. Checksum verification after sync

**Sync Procedure**:

```bash
# 1. Export current public keys from watch wallet
watch export-keys --output keys_v5.json

# 2. Transfer to offline systems (encrypted USB)
# Physically transfer USB to offline systems

# 3. Import on keygen wallet
keygen import-keys --input keys_v5.json --verify

# 4. Import on sign wallets
sign import-keys --input keys_v5.json --verify

# 5. Verify all systems have same keys
# Check checksums of key files on all systems
```

**Detection**:

- Pre-flight checks before signing
- Aggregated key verification

**Response**:

- Re-sync offline wallets
- Verify sync with checksums
- Restart signing with synced config

### Best Practices

#### Operator Training

**Minimum Training Requirements**:

1. **MuSig2 Basics** (2 hours)
   - Read: `docs/chains/btc/musig2/user-guide.md`
   - Understand two-round protocol
   - Understand nonce security

2. **Security Requirements** (1 hour)
   - Read this document
   - Understand threat model
   - Understand nonce reuse consequences

3. **Operational Procedures** (3 hours)
   - Documented SOPs for each role
   - File management protocols
   - Error recovery procedures

4. **Hands-On Practice** (8 hours)
   - Testnet practice (10+ transactions)
   - Error scenario practice
   - Emergency procedures drill

5. **Assessment** (1 hour)
   - Written test on security concepts
   - Practical test on testnet
   - Must score 90%+ to qualify

**Ongoing Training**:

- Quarterly refresher (1 hour)
- Review of any incidents
- Updates for procedure changes

#### File Management

**Standard File Naming Convention**:

```
{tx_type}_{request_id}_{stage}_{step}_{timestamp}.psbt

Fields:
- tx_type: "payment", "deposit", "sweep", etc.
- request_id: Numeric ID from database
- stage: "unsigned", "nonce", "signed", "final"
- step: Numeric counter (0, 1, 2, 3)
- timestamp: Unix timestamp (for uniqueness)

Examples:
payment_42_unsigned_0_1704067200.psbt    # Initial unsigned PSBT
payment_42_nonce_1_1704067201.psbt       # After keygen adds nonce
payment_42_nonce_2_1704067202.psbt       # After sign1 adds nonce
payment_42_nonce_3_1704067203.psbt       # All nonces collected
payment_42_signed_1_1704067204.psbt      # After keygen signature
payment_42_signed_2_1704067205.psbt      # After sign1 signature
payment_42_signed_3_1704067206.psbt      # All partial signatures
payment_42_final_1704067230.psbt         # Final signed PSBT
```

**File Organization**:

```
data/tx/btc/
├── pending/              # Active transactions being signed
│   ├── payment_42_unsigned_0_1704067200.psbt
│   ├── payment_42_nonce_1_1704067201.psbt
│   ├── payment_42_nonce_2_1704067202.psbt
│   ├── payment_42_nonce_3_1704067203.psbt
│   ├── payment_42_signed_1_1704067204.psbt
│   ├── payment_42_signed_2_1704067205.psbt
│   └── payment_42_signed_3_1704067206.psbt
├── completed/            # Successfully broadcast transactions
│   └── payment_42_final_1704067230.psbt
└── failed/               # Failed transactions (for investigation)
    └── payment_15_failed_reason.txt
```

**File Lifecycle**:

```
1. Creation (watch wallet)
   └─> data/tx/btc/pending/payment_42_unsigned_0.psbt

2. Round 1: Nonce Generation
   ├─> Keygen adds nonce → payment_42_unsigned_0_...1.psbt
   ├─> Sign1 adds nonce  → payment_42_unsigned_0_...2.psbt
   └─> Sign2 adds nonce  → payment_42_nonce_0_...3.psbt

3. Round 2: Signing
   ├─> Keygen signs → payment_42_unsigned_0_...1.psbt
   ├─> Sign1 signs  → payment_42_unsigned_1_...2.psbt
   └─> Sign2 signs  → payment_42_unsigned_2_...3.psbt

4. Aggregation (watch wallet)
   └─> Watch aggregates → payment_42_signed_3.psbt

5. Broadcast (watch wallet)
   └─> Transaction broadcast, move to completed/
       └─> data/tx/btc/completed/payment_42_signed_3.psbt

6. Archive (monthly)
   └─> Move old completed transactions to archive/
       └─> data/tx/btc/archive/2024-01/payment_42_signed_3.psbt
```

#### Monitoring

**Key Metrics to Monitor**:

1. **Nonce Uniqueness**

   ```sql
   -- Run every 5 minutes
   SELECT COUNT(*) as duplicate_nonces
   FROM (
       SELECT nonce, COUNT(*) as cnt
       FROM musig2_nonces
       GROUP BY nonce
       HAVING COUNT(*) > 1
   ) duplicates;

   -- Alert if result > 0
   ```

2. **Signing Success Rate**

   ```sql
   -- Run daily
   SELECT
       DATE(created_at) as date,
       COUNT(*) as total_transactions,
       SUM(CASE WHEN status = 'broadcast' THEN 1 ELSE 0 END) as successful,
       SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
       (SUM(CASE WHEN status = 'broadcast' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) as success_rate
   FROM transactions
   WHERE coin = 'btc' AND tx_type = 'musig2'
   GROUP BY DATE(created_at)
   ORDER BY date DESC
   LIMIT 30;

   -- Alert if success_rate < 95%
   ```

3. **Pending Transaction Age**

   ```sql
   -- Run hourly
   SELECT
       id,
       request_id,
       TIMESTAMPDIFF(HOUR, created_at, NOW()) as age_hours
   FROM transactions
   WHERE status = 'pending' AND tx_type = 'musig2'
   HAVING age_hours > 24;

   -- Alert if any transactions older than 24 hours
   ```

4. **File System Health**

   ```bash
   # Check for stale PSBT files
   find data/tx/btc/pending/ -name "*.psbt" -mtime +7

   # Alert if any files older than 7 days in pending/
   ```

**Alerting Channels**:

- **Email**: All alerts (operators + management)
- **SMS**: Critical alerts only (nonce reuse, key compromise)
- **Slack**: All alerts + daily summaries
- **Dashboard**: Real-time metrics display

---
