<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/architecture/wallet-architecture.tpl.md · Run `make docs` to regenerate.
-->

# Wallet Architecture

## Three Wallet Types

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Watch Wallet  │     │  Keygen Wallet  │     │   Sign Wallet   │
│    (Online)     │     │   (Offline)     │     │   (Offline)     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ • Monitor txs   │     │ • Generate keys │     │ • Auth signing  │
│ • Create unsig  │     │ • Create multis │     │ • Second+ sign  │
│ • Send signed   │     │ • First sign    │     │ • Export pubkey │
│ • Import pubkey │     │ • Export pubkey │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │    CSV/File Export    │    CSV/File Export    │
        └───────────────────────┴───────────────────────┘
```

## Security Model

1. **Keygen Wallet** (Offline): Generates HD wallet seeds and keys. Never connects to network.
2. **Sign Wallet** (Offline): Provides authorization signatures. Each operator has own instance.
3. **Watch Wallet** (Online): Only stores public keys. Cannot sign transactions.
