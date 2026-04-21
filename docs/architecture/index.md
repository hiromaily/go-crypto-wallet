<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/architecture/index.tpl.md · Run `make docs` to regenerate.
-->

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

# Architecture

This project follows **Clean Architecture** principles with clear layer separation:

```text
Interface Adapters Layer (internal/interface-adapters/*)
    ↓ depends on
Application Layer (internal/application/usecase/*)
    ↓ depends on
Domain Layer (internal/domain/*)
    ↑ implements interfaces defined by
Infrastructure Layer (internal/infrastructure/*)
```

## Key Principles

- **Domain Layer** (`internal/domain/`): Pure business logic with **zero infrastructure dependencies**
  - Defines interfaces that infrastructure must implement
  - Contains business rules, validators, and value objects
  - Most stable layer - changes affect all other layers

- **Application Layer** (`internal/application/usecase/`): Use cases orchestrate business logic
  - Coordinates domain objects and infrastructure services
  - Organized by wallet type (keygen, sign, watch) and cryptocurrency (btc, eth, xrp)
  - Each use case represents a single business operation

- **Infrastructure Layer** (`internal/infrastructure/`): External dependencies and implementations
  - Implements interfaces defined by domain layer (Dependency Inversion Principle)
  - Contains API clients, database repositories, file storage, and network communication
  - Easily replaceable and mockable for testing

- **Interface Adapters Layer** (`internal/interface-adapters/`): Adapters between use cases and external interfaces
  - CLI commands, HTTP handlers, and wallet adapters
  - Converts between external formats and application DTOs

- **Dependency Direction**: Outer layers depend on inner layers, never the reverse

## Architecture Dependency Flow

```
┌─────────────────────────────────────────┐
│   Interface Adapters (CLI, HTTP)        │
└──────────────────┬──────────────────────┘
                   │ depends on
                   ↓
┌─────────────────────────────────────────┐
│   Application Layer (Use Cases)          │
└───────────┬───────────────────┬─────────┘
            │ depends on        │ depends on
            ↓                   ↓
┌───────────────────┐  ┌──────────────────────┐
│   Domain Layer    │  │ Infrastructure Layer │
│ (Business Logic)  │  │ (External Systems)   │
└───────────────────┘  └──────────┬───────────┘
                                 │ implements
                                 ↓
                        ┌────────────────────┐
                        │ Domain Interfaces  │
                        └────────────────────┘
```

For detailed architecture guidelines, see `AGENTS.md`.

# Components inside repository

- xrpl-grpc-server [Deprecated]
  - ./apps/xrpl-grpc-server
- eth-contracts
  - ./apps/eth-contracts
