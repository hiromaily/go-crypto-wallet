## PSBT (Partially Signed Bitcoin Transactions)

PSBT (BIP174) is the standard format for offline/multi-party signing workflows.

### PSBT Workflow

```
1. Creator (Watch Wallet - Online)
   └── Create unsigned PSBT with UTXO data

2. Updater (Optional)
   └── Add metadata (derivation paths, etc.)

3. Signer(s) (Offline Wallets)
   └── Add partial signatures

4. Combiner (Optional)
   └── Combine multiple PSBTs

5. Finalizer (Watch Wallet)
   └── Create final scriptSig/witness

6. Extractor
   └── Extract broadcastable transaction
```

See [psbt/](../../../../docs/chains/btc/psbt/README.md) for detailed documentation.

---
