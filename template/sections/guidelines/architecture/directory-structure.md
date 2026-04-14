### Directory Structure

- `cmd/`: Application entry points (keygen, sign, watch)
- `internal/`: Internal packages (application-specific, not for external use)
  - `domain/`: **Domain layer** - Pure business logic (ZERO infrastructure dependencies)
    - `account/`: Account types, validators, and business rules
    - `transaction/`: Transaction types, state machine, validators
    - `wallet/`: Wallet types and definitions
    - `key/`: Key value objects and validators
    - `multisig/`: Multisig validators and business rules
    - `coin/`: Cryptocurrency type definitions
  - `application/`: **Application layer** - Use case layer (Clean Architecture)
    - `dto/`: **Data Transfer Objects** - Application-layer DTOs for port interfaces
      - `btc/`: Bitcoin DTOs (AddressInfo, UnspentOutput, TransactionResult, etc.)
      - Other coin DTOs (eth, xrp) as needed
    - `ports/`: **Interface definitions (abstractions)** - Contracts for infrastructure implementations
      - `btc/`: Bitcoiner interface (Bitcoin/BCH API abstraction)
      - `persistence/`: Repository interfaces (database abstractions)
      - `storage/`: File storage interfaces (TransactionFileRepositorier)
    - `usecase/`: Use case implementations organized by wallet type
      - `keygen/`: Key generation use cases (btc, eth, xrp, shared)
      - `sign/`: Signing use cases (btc, eth, xrp, shared)
      - `watch/`: Watch wallet use cases (btc, eth, xrp, shared)
  - `infrastructure/`: **Infrastructure layer** - **Implementation only** (NO interface definitions)
    - `api/`: External API clients
      - `bitcoin/`: Bitcoin/BCH Core RPC API clients (btc, bch)
      - `ethereum/`: Ethereum JSON-RPC API clients (eth, erc20)
      - `ripple/`: Ripple gRPC API clients (xrp)
    - `contract/`: Smart contract utilities (ERC-20 token ABI generated code)
    - `database/`: Database connections and generated code
      - `mysql/`: MySQL connection management
      - `sqlc/`: SQLC generated database code
    - `repository/`: Data persistence implementations
      - `cold/`: Cold wallet repository (keygen, sign)
      - `watch/`: Watch wallet repository
    - `storage/`: File storage implementations
      - `file/`: File-based storage (address, transaction)
    - `network/`: Network communication
      - `websocket/`: WebSocket client implementations
    - `wallet/key/`: Key generation logic - Infrastructure layer
  - `interface-adapters/`: **Interface Adapters layer** - Adapters between use cases and external interfaces
    - `cli/`: CLI command adapters (keygen, sign, watch)
      - `keygen/`: Keygen command implementations (api, create, export, imports, sign)
      - `sign/`: Sign command implementations (create, export, imports, sign)
      - `watch/`: Watch command implementations (api, create, imports, monitor, send)
    - `wallet/`: Wallet adapter interfaces and implementations
      - `interfaces.go`: Wallet interfaces (Keygener, Signer, Watcher)
      - `btc/`: Bitcoin wallet implementations
      - `eth/`: Ethereum wallet implementations
      - `xrp/`: XRP wallet implementations
  - `wallet/service/`: **Application layer** - Business logic orchestration (legacy/transitional)
    - `keygen/`: Key generation services (btc, eth, xrp, shared)
    - `sign/`: Signing services (btc, eth, xrp, shared)
    - `watch/`: Watch wallet services (btc, eth, xrp, shared)
  - `di/`: Dependency injection container
- `pkg/`: Shared packages (reusable, for external use)
  - `config/`: Configuration management utilities
    - `testutil/`: Test utilities for configuration
  - `logger/`: Logging utilities (structured logging, noop logger, slog support)
  - `converter/`: Data conversion utilities
  - `debug/`: Debug utilities
  - `serial/`: Serialization utilities
  - `testutil/`: Test utilities for various components (btc, eth, xrp, repository, suite)
  - `uuid/`: UUID generation utilities
  - `db/mysql/`: MySQL database connection utilities
  - `decimal/`: Decimal number utilities
  - `grpc/`: gRPC client utilities
  - `websocket/`: WebSocket client utilities
  - `di/`: Legacy dependency injection container (for backward compatibility)

  **Important**: See `pkg/AGENTS.md` for detailed guidelines on working with `pkg/` packages.
  **Critical Rule**: Packages in `pkg/` MUST NOT import or depend on any packages in `internal/` directory.
- `data/`: Generated files, configuration files
  - `address/`: Address data files (bch, btc, eth, xrp)
  - `config/`: Configuration files (account, wallet configs, node configs)
  - `contract/`: Contract ABI files
  - `keystore/`: Keystore files
  - `proto/`: Protocol buffer definitions (rippleapi)
  - `tx/`: Transaction data files (bch, btc, eth, xrp)
- `scripts/`: Operation scripts
  - `operation/`: Wallet operation scripts
  - `setup/`: Setup scripts for blockchain nodes

**Architecture Dependency Direction:**

```text
┌─────────────────────────────────────────────────────────────────────┐
│                    Clean Architecture Layers                         │
│                                                                      │
│  Interface Adapters (interface-adapters/*)                          │
│           ↓                                                          │
│  Application Layer (application/usecase/) ← depends on              │
│           ↓                                  ↓                       │
│  Application Ports (application/ports/) ← implemented by            │
│           ↓                                  ↓                       │
│  Domain Layer (domain/*)              Infrastructure Layer          │
│                                        (infrastructure/*)            │
│                                        - Implementations ONLY        │
│                                        - NO interface definitions    │
└─────────────────────────────────────────────────────────────────────┘

Key Principles:
1. Infrastructure implements interfaces defined in application/ports/
2. Application layer depends on application/ports/ (abstractions), not infrastructure (concrete)
3. This follows Dependency Inversion Principle (DIP)
4. Dependency flow: Interface Adapters → Application (Use Cases + Ports) ← Infrastructure
```
