## Incident Response

### Incident Classification

#### Severity 1 (CRITICAL)

**Nonce Reuse or Private Key Compromise Suspected**

- **Impact**: Immediate risk of fund loss
- **Response Time**: Immediate (within 15 minutes)
- **Escalation**: CEO, CTO, Security Officer

**Actions**:

1. STOP all operations
2. Assess scope of compromise
3. Assume key is compromised
4. Plan emergency fund sweep

#### Severity 2 (HIGH)

**Repeated Transaction Failures, System Unavailability**

- **Impact**: Operations disrupted, funds safe
- **Response Time**: Within 1 hour
- **Escalation**: Technical Lead, Operations Manager

**Actions**:

1. Investigate root cause
2. Activate backup procedures
3. Communicate ETA to stakeholders

#### Severity 3 (MEDIUM)

**Single Transaction Failure, Minor Issues**

- **Impact**: Minimal, isolated issue
- **Response Time**: Within 4 hours
- **Escalation**: On-duty operator

**Actions**:

1. Debug and resolve
2. Document in incident log
3. Monitor for recurrence

### Nonce Reuse Incident

**Trigger**: Nonce reuse detected (database alert, manual discovery).

#### Immediate Actions (0-15 minutes)

```
1. STOP ALL OPERATIONS
   - No new nonce generation
   - No signing operations
   - No transaction broadcasts

2. ALERT TEAM
   - Page on-call security officer
   - Alert all MuSig2 operators
   - Notify management

3. PRESERVE EVIDENCE
   - Do NOT modify database
   - Snapshot current state
   - Save all log files

4. ASSESS SCOPE
   - Which nonce was reused?
   - Which signer is affected?
   - Which transactions?
   - Was nonce used to sign different messages?
```

#### Investigation (15 minutes - 2 hours)

```
1. DETERMINE IF KEY IS COMPROMISED
   Query:
   - Were two different messages signed with same nonce?
   - Have those signatures been broadcast?
   - Are they visible on the blockchain?

   If YES to all:
   └─> ASSUME PRIVATE KEY IS COMPROMISED

   If NO (nonce reuse caught before signing):
   └─> KEY LIKELY SAFE, but verify carefully

2. IDENTIFY ROOT CAUSE
   - Software bug?
   - Database failure?
   - Operator error?
   - File reuse?

3. CALCULATE EXPOSURE
   - How many BTC controlled by compromised key?
   - Which addresses?
   - Which UTXOs?
```

#### Mitigation (2-24 hours)

**If Key is Compromised**:

```
EMERGENCY FUND SWEEP

1. PREPARE NEW ADDRESSES
   - Generate NEW keys (not from same seed)
   - Use traditional P2WSH temporarily (faster setup)
   - Verify new keys are properly backed up

2. CREATE SWEEP TRANSACTION
   - Collect all UTXOs controlled by compromised key
   - Send to new safe addresses
   - Use HIGH fee (priority confirmation)

3. SIGN WITH REMAINING KEYS
   - 2-of-3 multisig: Use the 2 uncompromised keys
   - Sign immediately (race against attacker)

4. BROADCAST IMMEDIATELY
   - Submit to mempool with high priority
   - Submit to multiple nodes
   - Monitor for confirmation

5. MONITOR MEMPOOL
   - Watch for competing transactions (attacker trying to steal)
   - If attacker transaction seen, consider RBF (Replace-By-Fee)

6. CONFIRMATION
   - Wait for 1 confirmation (10 min average)
   - Funds are safe once confirmed
```

**If Key is Safe** (nonce reuse caught before signing):

```
1. FIX ROOT CAUSE
   - Patch software bug
   - Repair database
   - Update procedures

2. VERIFY FIX
   - Test in testnet
   - Verify nonce uniqueness enforced

3. GENERATE NEW NONCES
   - Discard old nonces
   - Regenerate fresh nonces
   - Continue with transaction

4. RESUME OPERATIONS
   - Gradual restart
   - Enhanced monitoring
```

#### Post-Incident (24-48 hours)

```
1. ROOT CAUSE ANALYSIS
   - Detailed investigation report
   - Timeline of events
   - Contributing factors
   - Why controls failed

2. CORRECTIVE ACTIONS
   - Fix identified weaknesses
   - Enhance controls
   - Update procedures

3. PREVENTIVE MEASURES
   - Additional safeguards
   - Enhanced monitoring
   - Operator retraining

4. LESSONS LEARNED
   - Team meeting
   - Update documentation
   - Share knowledge

5. FOLLOW-UP
   - Verify corrective actions effective
   - Monitor for recurrence
   - Schedule review in 30 days
```

### Communication Plan

#### Internal Communication

**Severity 1 (Critical)**:

- Immediately alert: CEO, CTO, Security Officer, All Operators
- Method: SMS + Phone Call (ensure acknowledgment)
- Updates: Every 15 minutes until resolved

**Severity 2 (High)**:

- Alert: Technical Lead, Operations Manager, On-call team
- Method: Slack + Email
- Updates: Hourly

**Severity 3 (Medium)**:

- Inform: Operations team
- Method: Slack
- Updates: As resolved

#### External Communication (if applicable)

**Customers/Partners**:

- Inform IF funds are at risk or operations disrupted
- Method: Email + status page update
- Frequency: As situation develops
- Transparency: Provide facts, no speculation

**Regulators** (if applicable):

- Report IF required by regulations
- Consult legal team first
- Follow compliance procedures

---
