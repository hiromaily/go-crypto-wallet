## FAQ

### General Questions

#### Q: How long does migration take?

**A**: Typical timeline: 4-6 weeks for full migration (including planning, testing, gradual rollout, and fund sweeping). If you're only creating new MuSig2 addresses without sweeping old funds, you can complete in 2-3 weeks.

#### Q: Can I migrate incrementally?

**A**: Yes! This is the recommended approach. Start with 10 addresses, test thoroughly, then gradually increase. Both P2WSH and MuSig2 can coexist indefinitely.

#### Q: What if I encounter issues mid-migration?

**A**: You can pause at any time. Complete any in-flight transactions, then stop creating new MuSig2 addresses. See [Rollback Procedures](#rollback-procedures) for details.

#### Q: Do I need to sweep old P2WSH funds immediately?

**A**: No. You can keep old P2WSH addresses active and use them alongside MuSig2. Sweep when convenient (low fee environment, scheduled maintenance window, etc.).

### Technical Questions

#### Q: Are MuSig2 transactions more expensive to create?

**A**: MuSig2 transactions are **cheaper on-chain** (30-50% lower fees) but require **more coordination** (two rounds vs one round). The cost savings far outweigh the coordination overhead for most use cases.

#### Q: Can I mix P2WSH and MuSig2 inputs in one transaction?

**A**: Yes! A single transaction can spend from both P2WSH and MuSig2 UTXOs. Each input type uses its respective signing workflow.

#### Q: What happens if a nonce is reused?

**A**: **CRITICAL SECURITY ISSUE**. Nonce reuse leaks the private key. The system has multiple protections:

1. Database unique constraints prevent storage of duplicate nonces
2. Application-level validation before signing
3. Error handling stops signing if duplicate detected

Never override these protections. See [Security Documentation](/chains/btc/musig2/security) for details.

#### Q: How do I verify a MuSig2 transaction before broadcast?

**A**:

```bash
# 1. Decode PSBT
watch btc api decodepsbt $(cat payment_15_signed_3.psbt)

# 2. Verify signature present
watch btc api analyzepsbt $(cat payment_15_signed_3.psbt)

# 3. Extract and verify transaction
watch btc api finalizepsbt $(cat payment_15_signed_3.psbt) true

# 4. Verify transaction validity (doesn't broadcast)
watch btc api testmempoolaccept '["<tx_hex>"]'
```

#### Q: Can I recover if I lose a PSBT file during signing?

**A**:

- **During Round 1 (nonce generation)**: Regenerate nonces. The Watch wallet recreates the PSBT from the transaction.
- **During Round 2 (signing)**: More complex. You may need to restart the entire signing process.
- **Best Practice**: Always backup PSBT files after each step.

### Operational Questions

#### Q: How do I train my team on MuSig2?

**A**: Recommended approach:

1. **Read Documentation**: All operators read `docs/chains/btc/musig2/user-guide.md` and `docs/security/musig2_security.md`
2. **Testnet Practice**: Each operator practices full workflow on testnet (5-10 transactions)
3. **Shadow Production**: Observe experienced operator in production
4. **Supervised Production**: Perform under supervision
5. **Independent**: Solo operations after 10 successful supervised transactions

#### Q: What monitoring should I set up?

**A**:

1. **Nonce Tracking**: Alert on any nonce reuse attempts
2. **Transaction Success Rate**: Alert if success rate drops below 95%
3. **PSBT File Management**: Alert on missing or stale files
4. **Signature Verification**: Alert on signature verification failures
5. **Fee Estimation**: Monitor actual fee savings vs expected

#### Q: How do I handle emergency situations?

**A**: See [Rollback Procedures](#rollback-procedures). Key points:

- Always complete in-flight transactions first
- Document the issue thoroughly
- Don't panic-rollback for operator errors (training issue)
- Do rollback for security issues or critical bugs

#### Q: Can I automate the signing workflow?

**A**: **Partial automation possible**, but:

- ✅ Nonce generation can be automated (Round 1)
- ✅ Partial signature creation can be automated (Round 2)
- ❌ **Never automate**: Private key access, nonce reuse checks, final broadcast approval
- **Recommendation**: Keep human oversight at critical points

### Migration-Specific Questions

#### Q: Should I migrate all accounts at once?

**A**: **No**. Start with low-risk accounts (e.g., deposit), then gradually expand. This allows you to:

- Build team confidence
- Identify issues in low-risk environment
- Adjust procedures based on lessons learned

#### Q: What if I discover an issue after sweeping funds?

**A**:

1. **Stop**: Don't create more MuSig2 transactions
2. **Assess**: Determine if funds are at risk
3. **Secure**: If funds are safe, plan fix; if at risk, execute emergency rollback
4. **Fix**: Address root cause
5. **Validate**: Test fix thoroughly
6. **Resume**: Continue migration or complete rollback

#### Q: How do I calculate break-even point for migration?

**A**:

```bash
# Migration costs:
SETUP_TIME=40        # hours (team time)
HOURLY_RATE=100      # USD (team hourly cost)
SWEEP_FEES=50000     # sats (estimated fees to sweep old UTXOs)
BTC_PRICE_USD=40000  # Example BTC price in USD

TOTAL_COST=$(echo "scale=2; $SETUP_TIME * $HOURLY_RATE + ($SWEEP_FEES / 100000000) * $BTC_PRICE_USD" | bc -l)

# Monthly savings:
TX_PER_MONTH=100
FEE_RATE=50          # sat/vB
P2WSH_SIZE=385       # bytes
MUSIG2_SIZE=225      # bytes
SAVINGS_PER_TX=$(echo "($P2WSH_SIZE - $MUSIG2_SIZE) * $FEE_RATE" | bc)
MONTHLY_SAVINGS=$(echo "scale=2; ($SAVINGS_PER_TX * $TX_PER_MONTH / 100000000) * $BTC_PRICE_USD" | bc -l)

# Break-even time:
if [ "$(echo "$MONTHLY_SAVINGS > 0" | bc -l)" -eq 1 ]; then
    BREAK_EVEN_MONTHS=$(echo "scale=2; $TOTAL_COST / $MONTHLY_SAVINGS" | bc -l)
else
    BREAK_EVEN_MONTHS="N/A"
fi

echo "Break-even in $BREAK_EVEN_MONTHS months"

# For 100 tx/month at 50 sat/vB, BTC at $40,000:
# Monthly savings: ~$128 USD
# Break-even: 3-5 months
```

#### Q: Can I migrate without downtime?

**A**: **Yes**. The gradual migration approach allows zero-downtime migration:

1. Create MuSig2 addresses alongside P2WSH
2. Gradually shift new transactions to MuSig2
3. P2WSH remains fully operational during migration
4. Sweep funds during scheduled maintenance or low-activity periods

#### Q: What if Bitcoin Core introduces MuSig2 wallet support?

**A**: Future Bitcoin Core versions may add native MuSig2 wallet support. When that happens:

- Your addresses remain valid (standard P2TR)
- Your keys remain valid
- You may be able to import keys into Bitcoin Core wallet
- **Recommendation**: Monitor Bitcoin Core release notes for MuSig2 wallet features

---
