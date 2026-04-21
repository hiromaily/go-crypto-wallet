## Appendix: Attack Scenarios

### Scenario A1: Nonce Reuse Attack

**Attacker Goal**: Extract private key by observing two signatures with same nonce.

**Attack Steps**:

1. Attacker observes transaction 1 on blockchain:
   - Message: m1 (transaction hash)
   - Signature: s1
   - Public nonce: R

2. Attacker observes transaction 2 on blockchain:
   - Message: m2 (different transaction hash)
   - Signature: s2
   - Public nonce: R (same as transaction 1!)

3. Attacker computes private key:

   ```
   s1 = k + Hash(R || P || m1) * x
   s2 = k + Hash(R || P || m2) * x

   Subtract: s1 - s2 = (Hash(R || P || m1) - Hash(R || P || m2)) * x

   Solve for x:
   x = (s1 - s2) / (Hash(R || P || m1) - Hash(R || P || m2))
   ```

4. Attacker now has private key, can steal all funds

**Defense**:

- Nonce uniqueness enforced at multiple layers
- Monitoring detects nonce reuse before broadcast
- Emergency fund sweep if compromise detected

### Scenario A2: Rogue Key Attack (Prevented)

**Attacker Goal**: Gain full control over multisig by manipulating their public key.

**Attack Steps** (Attempted):

1. Honest parties have keys: P1, P2
2. Attacker claims key: P3' = P3 - P1 - P2
3. Naive aggregation: P_agg = P1 + P2 + P3' = P3
4. Attacker controls P_agg with only their private key

**Why This Fails in MuSig2**:

- MuSig2 uses deterministic coefficients
- Coefficients depend on all public keys
- Attacker cannot manipulate aggregated key
- Attack is mathematically prevented

### Scenario A3: File Substitution Attack

**Attacker Goal**: Substitute malicious PSBT file to redirect funds.

**Attack Steps**:

1. Attacker gains access to file transport (USB drive)
2. Attacker replaces PSBT file with malicious version
   - Changes output address to attacker's address
3. Signers unknowingly sign malicious transaction
4. Funds sent to attacker

**Defense**:

1. File integrity checks (checksums)
2. Review transaction details before signing
3. Physical security for file transport
4. Encrypted file containers

**Detection**:

- Checksum verification fails
- Manual review notices wrong address

### Scenario A4: Timing Attack on Nonce Generation

**Attacker Goal**: Predict nonces by observing timing.

**Attack Steps**:

1. Attacker measures time taken for nonce generation
2. Attempts to correlate timing with random values
3. Predicts nonce values
4. Uses predicted nonces to forge signatures

**Why This Fails**:

- `crypto/rand` provides cryptographically secure randomness
- Timing does not leak information about random values
- 2^256 keyspace makes prediction infeasible

### Scenario A5: Social Engineering

**Attacker Goal**: Trick operator into signing malicious transaction.

**Attack Steps**:

1. Attacker impersonates manager
2. Requests urgent transaction signing
3. Operator skips verification steps
4. Signs transaction sending funds to attacker

**Defense**:

1. Verification procedures mandatory (no exceptions)
2. Multi-person approval for transactions
3. Out-of-band confirmation for unusual requests
4. Security awareness training

---
