<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/guidelines/multi-chain.tpl.md · Run `make docs` to regenerate.
-->

# Multi-Chain Support

This document describes the multi-chain cryptocurrency support in the go-crypto-wallet project.

## Supported Cryptocurrencies

The project supports the following cryptocurrencies:

- **BTC** (Bitcoin)
- **BCH** (Bitcoin Cash)
- **ETH** (Ethereum)
- **XRP** (Ripple)
- **ERC-20 tokens** (Ethereum-based tokens)

## Wallet Types Understanding

The project implements three types of wallets, each with specific roles:

### Watch Wallet

- **Status**: Online
- **Keys**: Public keys only (no private keys)
- **Purpose**: Creates and sends transactions
- **Security**: Lower risk (no private keys exposed to online systems)
- **Use Cases**:
  - Monitor addresses and balances
  - Create unsigned transactions
  - Broadcast signed transactions
  - Query blockchain state

### Keygen Wallet

- **Status**: Offline
- **Keys**: Generates and stores private keys
- **Purpose**: Generates keys, provides first signature for multisig
- **Security**: Highest security (offline, air-gapped if possible)
- **Use Cases**:
  - Generate new key pairs
  - Create multisig addresses
  - Provide first signature in multisig transactions
  - Export public keys to watch wallet

### Sign Wallet

- **Status**: Offline
- **Keys**: Stores private keys for signing
- **Purpose**: Provides second and subsequent signatures for multisig
- **Security**: High security (offline)
- **Use Cases**:
  - Import unsigned transactions
  - Sign transactions with stored keys
  - Complete multisig signing process
  - Export signed transactions

## Blockchain Communication

### BTC/BCH (Bitcoin Core RPC API)

**Implementation**: Bitcoin Core RPC API

**Location**: `internal/infrastructure/api/btc/`

**Features**:

- Transaction creation and broadcasting
- UTXO management
- Address generation
- Block and transaction queries

**API Clients**:

- `btc/`: Bitcoin Core RPC client
- `bch/`: Bitcoin Cash Core RPC client

### ETH (Ethereum JSON-RPC API)

**Implementation**: Ethereum JSON-RPC API

**Location**: `internal/infrastructure/api/eth/`

**Features**:

- Transaction creation and broadcasting
- Smart contract interaction
- ERC-20 token support
- Gas price estimation
- Nonce management

**API Clients**:

- `eth/`: Ethereum JSON-RPC client
- `erc20/`: ERC-20 token client

### XRP

TODO: updating

**Implementation**: Communication via

**Location**: `internal/infrastructure/api/xrp/`

**Features**:

- Transaction creation and submission
- Account management
- Balance queries
- Payment tracking

**API Client**:

- `xrp/`: Ripple gRPC client [Deprecated]

**Note**: XRP uses protocol buffers for gRPC communication. See [Code Generation Guidelines](./code-generation.md) for protobuf code generation.

## Multi-Chain Architecture

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

## Adding New Cryptocurrency Support

When adding support for a new cryptocurrency, follow these steps:

### 1. Domain Layer

Add cryptocurrency-specific types and validators:

- Add coin type to `internal/domain/coin/`
- Add validators for addresses, amounts, transactions
- Define domain interfaces for the new cryptocurrency

### 2. Infrastructure Layer

Implement API clients and repositories:

- Create API client in `internal/infrastructure/api/{coin}/`
- Implement repository interfaces in `internal/infrastructure/repository/`
- Add storage implementations if needed

### 3. Application Layer

Create use cases for each wallet type:

- Key generation use cases in `internal/application/usecase/keygen/{coin}/`
- Signing use cases in `internal/application/usecase/sign/{coin}/`
- Watch wallet use cases in `internal/application/usecase/watch/{coin}/`

### 4. Interface Adapters Layer

Implement CLI commands and wallet adapters:

- Add CLI commands in `internal/interface-adapters/cli/{wallet-type}/`
- Implement wallet adapter in `internal/interface-adapters/wallet/{coin}/`

### 5. Dependency Injection

Wire up dependencies:

- Add DI container setup in `internal/di/`
- Register new cryptocurrency services

### 6. Configuration

Add configuration support:

- Update configuration files in `config/wallet/`
- Add coin-specific settings

### 7. Testing

Add tests for all layers:

- Domain layer: Pure unit tests
- Application layer: Use case tests with mocked infrastructure
- Infrastructure layer: Unit tests + integration tests
- Interface adapters: Command tests with mocked use cases

## Cryptocurrency-Specific Considerations

### Bitcoin (BTC) / Bitcoin Cash (BCH)

- UTXO-based transaction model
- Multisig support (P2SH, P2WSH)
- Fee estimation and management
- Address types (P2PKH, P2SH, P2WPKH, P2WSH)

### Ethereum (ETH)

- Account-based transaction model
- Nonce management
- Gas price and gas limit
- Smart contract interaction
- ERC-20 token support

### Ripple (XRP)

- Account-based transaction model
- Sequence number management
- Destination tag support
- Account reserve requirements
- Memo field support

## See Also

- [Architecture Guidelines](./architecture.md) - Layer structure and organization
- [Code Generation Guidelines](./code-generation.md) - Protocol buffer generation for XRP
- [Testing Guidelines](./testing.md) - Testing multi-chain implementations
- `pkg/testutil/` - Test utilities for BTC, ETH, XRP
