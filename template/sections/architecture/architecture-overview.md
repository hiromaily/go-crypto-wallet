## Architecture Overview

### Clean Architecture Layers

The PSBT implementation follows Clean Architecture principles:

```
┌─────────────────────────────────────────────────┐
│         Interface Adapters Layer                 │
│  (CLI, Wallet Adapters)                          │
│  - internal/interface-adapters/cli/              │
│  - internal/interface-adapters/wallet/           │
└───────────────┬─────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Application Layer (Use Cases)            │
│  - internal/application/usecase/watch/btc/       │
│  - internal/application/usecase/keygen/btc/      │
│  - internal/application/usecase/sign/btc/        │
└───────────────┬─────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Domain Layer (Business Logic)            │
│  - internal/domain/transaction/                  │
│  - internal/domain/account/                      │
│  - internal/domain/key/                          │
└─────────────────────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│         Infrastructure Layer                     │
│  - internal/infrastructure/api/btc/btc/      │
│  - internal/infrastructure/storage/file/         │
│  - internal/infrastructure/repository/           │
└─────────────────────────────────────────────────┘
```

### PSBT Flow Through Layers

```
User Command (CLI)
    │
    ▼
Interface Adapter (e.g., watch/btc.BTCWatch)
    │
    ▼
Use Case (e.g., CreateTransactionUseCase)
    │
    ├──> Infrastructure: Bitcoin API (CreatePSBT)
    ├──> Infrastructure: File Storage (WritePSBTFile)
    └──> Infrastructure: Database (InsertTransaction)
```

---
