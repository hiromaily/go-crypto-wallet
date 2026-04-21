## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Use air-gapped systems for key generation and signing
- Implement proper entropy for key generation
- Use hardware security modules (HSMs) for production

### Transaction Security

- Verify all transaction details before signing
- Implement multi-signature for high-value accounts
- Use PSBT for offline signing workflows
- Validate change addresses

### Nonce Security (MuSig2)

- **CRITICAL:** Never reuse nonces in MuSig2
- Generate cryptographically secure random nonces
- Delete nonces immediately after signing

See [musig2/security.md](musig2/security.md) for details.

---
