## Prerequisites

### Technical Requirements

#### Bitcoin Core Version

```bash
# Check Bitcoin Core version
bitcoin-cli -version

# Required: v22.0 or higher for Taproot support
# Recommended: v25.0+ for best Taproot support
```

If using older version:

1. Backup wallet data
2. Stop Bitcoin Core
3. Upgrade to v25.0 or higher
4. Restart and verify sync

#### Go Version

```bash
# Check Go version
go version

# Required: Go 1.21 or higher
```

#### Database Schema

MuSig2 uses existing database tables with additional columns:

```sql
-- Check if taproot_address column exists
DESCRIBE account_key;

-- Should see:
-- taproot_address VARCHAR(255) NULL
-- multisig_address VARCHAR(255) NULL (for compatibility)
```

If columns are missing, run migration:

```bash
# Apply database migrations
make db-migrate
```

### Infrastructure Requirements

#### Test Environment

**CRITICAL**: Always test migration in testnet first.

Required test setup:

1. **Testnet Bitcoin Core node**
   - Testnet fully synced
   - RPC credentials configured
2. **Test Wallets**
   - Keygen wallet (offline test system)
   - Sign wallet(s) (offline test systems)
   - Watch wallet (online test system)
3. **Test Funds**
   - Small amount of testnet BTC for testing
   - Source: Bitcoin testnet faucets

#### Production Environment

Before migrating production:

1. **Backup Strategy**
   - Full database backup
   - Keystore backups
   - Configuration file backups
2. **Monitoring**
   - Transaction monitoring
   - Error alerting
   - Nonce tracking
3. **Rollback Plan**
   - Documented rollback procedures
   - Tested rollback process
   - Traditional multisig capability preserved

### Team Readiness

#### Required Knowledge

Your team should understand:

1. **MuSig2 Basics**
   - Two-round protocol
   - Nonce security requirements
   - Key aggregation
   - Read: `docs/chains/btc/musig2/user-guide.md`

2. **Security Implications**
   - Nonce reuse consequences
   - Signature verification
   - Key management
   - Read: `docs/security/musig2_security.md`

3. **Operational Procedures**
   - File management
   - Error recovery
   - Monitoring
   - Read: `docs/chains/btc/musig2/user-guide.md` (Best Practices)

#### Training Checklist

- [ ] All operators reviewed MuSig2 documentation
- [ ] Team completed testnet practice runs
- [ ] Security procedures documented and reviewed
- [ ] Error scenarios practiced
- [ ] Rollback procedures tested
- [ ] Incident response plan updated

### Risk Assessment

Before proceeding, assess these risks:

#### Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Nonce reuse** | Medium | Critical | Database constraints + monitoring |
| **Key loss** | Low | Critical | Robust backup procedures |
| **Software bugs** | Low | High | Thorough testing + gradual rollout |
| **Network issues** | Medium | Medium | Offline wallet architecture |
| **Human error** | High | Medium | Clear procedures + training |

#### Operational Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Team learning curve** | High | Low | Extended testing period |
| **Coordination overhead** | Medium | Medium | Clear file management process |
| **Monitoring gaps** | Medium | Medium | Enhanced monitoring setup |
| **Documentation gaps** | Low | Medium | This migration guide |

#### Business Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Downtime** | Low | High | Gradual migration |
| **Fund loss** | Very Low | Critical | Testnet validation + small batches |
| **Regulatory issues** | Low | High | Legal review before migration |
| **Customer confusion** | Medium | Low | Clear communication plan |

---
