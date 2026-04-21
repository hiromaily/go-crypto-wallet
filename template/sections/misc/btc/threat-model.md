## Threat Model

### Assets to Protect

1. **Private Keys**
   - **Keygen private key** (offline wallet)
   - **Sign wallet private keys** (offline wallets)
   - **Impact if compromised**: Complete loss of funds

2. **Nonces**
   - **Secret nonces** (generated in Round 1, never shared)
   - **Impact if reused**: Private key leakage → loss of funds

3. **Transaction Data**
   - **Unsigned PSBTs** (before signing)
   - **Partial signatures** (during signing)
   - **Impact if compromised**: Transaction details leaked, no fund loss

### Threat Actors

#### External Attackers

**Capabilities**:

- Observe blockchain transactions
- Observe network traffic (if not encrypted)
- Attempt to exploit software vulnerabilities
- Social engineering attacks

**Goals**:

- Steal funds by exploiting vulnerabilities
- Learn about transaction patterns (privacy attack)
- Disrupt operations (DoS)

**Mitigation**:

- Offline wallets (air-gapped) for key operations
- Encrypted file transport
- Regular security updates
- Employee security training

#### Insider Threats (Malicious)

**Capabilities**:

- Access to operational systems
- Knowledge of procedures
- Ability to manipulate files or processes

**Goals**:

- Steal funds
- Sabotage operations

**Mitigation**:

- Multi-signature requirement (no single person can steal funds)
- Audit logging
- Access controls
- Background checks

#### Insider Threats (Accidental)

**Capabilities**:

- Authorized access to systems
- Perform operations

**Goals**:

- None (accidental errors, not malicious)

**Mitigation**:

- Training and procedures
- Monitoring and alerts
- Error detection (multiple validation layers)

### Attack Surfaces

#### 1. Nonce Generation

**Vulnerability**: Weak random number generation
**Impact**: Predictable nonces → private key recovery
**Mitigation**:

- Use cryptographically secure RNG (`crypto/rand` in Go)
- btcd library handles nonce generation securely
- Never implement custom nonce generation

#### 2. Nonce Storage

**Vulnerability**: Nonce reuse due to storage issues
**Impact**: Private key leakage
**Mitigation**:

- Database unique constraints
- Application-level validation
- Monitoring for reuse attempts

#### 3. File Transport

**Vulnerability**: PSBT file interception/modification
**Impact**: Transaction manipulation, DoS
**Mitigation**:

- Encrypted transport (if over network)
- File integrity checks (checksums)
- Physical security for USB transport

#### 4. Signature Aggregation

**Vulnerability**: Accepting invalid partial signatures
**Impact**: Transaction broadcast failure, fund stuck
**Mitigation**:

- Verify each partial signature
- Verify final aggregated signature
- Test transaction validity before broadcast

#### 5. Software Bugs

**Vulnerability**: Implementation errors in MuSig2 code
**Impact**: Various (depends on bug)
**Mitigation**:

- Use well-tested library (`btcd/btcec/v2/schnorr/musig2`)
- Extensive testing
- Code review
- Security audits

---
