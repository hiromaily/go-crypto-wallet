## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Use air-gapped systems for key generation and signing
- Implement proper entropy for key generation

### Address Security

- **Always use CashAddr** to prevent BTC/BCH cross-sends
- Validate address format before sending
- Double-check network prefix (`bitcoincash:` vs `bchtest:`)

### Transaction Security

- Verify all transaction details before signing
- Implement multi-signature for high-value accounts
- Validate change addresses

### Replay Protection

- BCH transactions include SIGHASH_FORKID
- Transactions are not valid on BTC chain
- Always verify fork ID in signatures

### 51% Attack Considerations

BCH has lower hashrate than BTC, making it theoretically more vulnerable to 51% attacks:

- Wait for more confirmations for large transactions
- Consider 10+ confirmations for high-value transfers

---
