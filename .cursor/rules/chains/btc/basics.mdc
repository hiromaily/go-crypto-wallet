---
description: Claude Rules - Bitcoin (BTC) Development
globs: 
  - internal/application/usecase/keygen/btc/*.go
  - internal/application/usecase/sign/btc/*.go
  - internal/application/usecase/watch/btc/*.go
  - internal/application/dto/btc/*.go
  - internal/application/ports/api/btc/*.go
  - internal/application/ports/repository/watch/btc_transaction.go
  - internal/infrastructure/api/btc/connection.go
  - internal/infrastructure/api/btc/btc/*.go
  - internal/interface-adapters/cli/keygen/api/btc/*.go
  - internal/interface-adapters/cli/watch/api/btc/*.go
  - internal/interface-adapters/wallet/btc/*.go
alwaysApply: false
---


# Claude Rules - Bitcoin (BTC) Development

## Overview

Rules for Claude Code when working on BTC-related code in go-crypto-wallet.
Always consult the documentation under `docs/chains/btc/` for context before making changes.

## Documentation Map

### Start Here

| Task                            | Primary Doc                                                                                                               |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Understand wallet architecture  | [docs/chains/btc/architecture.md](../../../docs/chains/btc/architecture.md)                                               |
| Bitcoin protocol fundamentals   | [docs/chains/btc/README.md](../../../docs/chains/btc/README.md)                                                           |
| Transaction flow (BTC-specific) | [docs/chains/btc/operations/wallet-flow.md](../../../docs/chains/btc/operations/wallet-flow.md)                           |
| E2E transaction patterns        | [docs/chains/btc/operations/e2e-transaction-patterns.md](../../../docs/chains/btc/operations/e2e-transaction-patterns.md) |

### By Topic

| Topic                                | Documents                                                                                                                                                                                                                           |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Address types & derivation paths** | [overview/address-types.md](../../../docs/chains/btc/overview/address-types.md), [overview/technical-reference.md](../../../docs/chains/btc/overview/technical-reference.md)                                                        |
| **Key generation**                   | [keygen/improvements-2025.md](../../../docs/chains/btc/keygen/improvements-2025.md), [keygen/interface-design.md](../../../docs/chains/btc/keygen/interface-design.md)                                                              |
| **PSBT (offline signing)**           | [psbt/user-guide.md](../../../docs/chains/btc/psbt/user-guide.md), [psbt/developer-guide.md](../../../docs/chains/btc/psbt/developer-guide.md), [psbt/implementation.md](../../../docs/chains/btc/psbt/implementation.md)           |
| **Output Descriptors**               | [descriptor/architecture.md](../../../docs/chains/btc/descriptor/architecture.md), [descriptor/api.md](../../../docs/chains/btc/descriptor/api.md), [descriptor/development.md](../../../docs/chains/btc/descriptor/development.md) |
| **Taproot (BIP86/P2TR)**             | [taproot/user-guide.md](../../../docs/chains/btc/taproot/user-guide.md), [taproot/testing.md](../../../docs/chains/btc/taproot/testing.md)                                                                                          |
| **MuSig2 multisig**                  | [musig2/architecture.md](../../../docs/chains/btc/musig2/architecture.md), [musig2/user-guide.md](../../../docs/chains/btc/musig2/user-guide.md), [musig2/security.md](../../../docs/chains/btc/musig2/security.md)                 |
| **Testing & verification**           | [testing/README.md](../../../docs/chains/btc/testing/README.md), [testing/pattern3-verification.md](../../../docs/chains/btc/testing/pattern3-verification.md)                                                                      |

## Key Code Paths

### Implementation

```
internal/application/usecase/keygen/btc/   # BTC keygen use cases
internal/application/usecase/sign/btc/    # BTC signing use cases
internal/application/usecase/watch/btc/   # BTC watch use cases
internal/infrastructure/api/btc/          # Bitcoin Core RPC client
internal/infrastructure/api/btc/bch/      # BCH overrides (DO NOT modify BTC for BCH)
internal/infrastructure/wallet/key/       # Key generation (SECURITY-CRITICAL)
```

### CLI

```
internal/interface-adapters/cli/keygen/api/btc/   # Keygen CLI commands
internal/interface-adapters/cli/watch/api/btc/    # Watch CLI commands
```

### E2E Scripts

```
scripts/operation/btc/e2e/   # E2E test shell scripts
scripts/operation/btc/       # BTC operation scripts
```

## Critical Rules

### Address Types

- **Recommended:** P2TR (Taproot, `bc1p` prefix) via BIP86 derivation `m/86'/0'/account'/change/index`
- Use **bech32m** encoding for Taproot addresses (NOT bech32)
- See `btc-terminology` skill to avoid confusion between bech32m (encoding) and taproot (address type)

### PSBT Workflow

- Watch wallet (online) creates unsigned PSBT
- Keygen/Sign wallets (offline/air-gapped) add signatures
- Watch wallet finalizes and broadcasts
- Always use [psbt/developer-guide.md](../../../docs/chains/btc/psbt/developer-guide.md) for implementation

### MuSig2 Security

- **NEVER reuse nonces** — doing so leaks the private key
- Generate fresh nonces for every signing session
- Delete nonces immediately after signing
- See [musig2/security.md](../../../docs/chains/btc/musig2/security.md) for full requirements

### BCH Separation

- **DO NOT** modify BTC code for BCH issues
- BCH overrides live in `internal/infrastructure/api/btc/bch/`
- Use the `bch-development` skill for BCH tasks

## Applicable Skills

| Situation                        | Skill             |
| -------------------------------- | ----------------- |
| Generating/updating mocks        | `mockery`         |
| Writing Go code                  | `go-development`  |
| BTC address/encoding terminology | `btc-terminology` |
| BCH-specific code                | `bch-development` |
| Running wallet CLI commands      | `wallet-cli`      |
