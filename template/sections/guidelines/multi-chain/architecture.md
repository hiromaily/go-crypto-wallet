### Multi-Chain Architecture

The project organizes multi-chain support by wallet type and cryptocurrency:

```text
internal/application/usecase/
├── keygen/
│   ├── btc/     # Bitcoin-specific key generation
│   ├── eth/     # Ethereum-specific key generation
│   ├── xrp/     # XRP-specific key generation
│   └── shared/  # Shared key generation logic
├── sign/
│   ├── btc/     # Bitcoin-specific signing
│   ├── eth/     # Ethereum-specific signing
│   ├── xrp/     # XRP-specific signing
│   └── shared/  # Shared signing logic
└── watch/
    ├── btc/     # Bitcoin-specific watching
    ├── eth/     # Ethereum-specific watching
    ├── xrp/     # XRP-specific watching
    └── shared/  # Shared watching logic
```
