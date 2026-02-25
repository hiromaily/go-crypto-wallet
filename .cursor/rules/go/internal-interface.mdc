---
description: Interface Definition Rules for `internal/`
globs: ["internal/**"]
alwaysApply: false
---


# Interface Definition Rules for `internal/`

## Overview

Rules for defining and using interfaces inside the `internal/` directory.
These rules enforce the Interface Segregation Principle (ISP) throughout the codebase.

## Core Principles

**Interfaces are defined where they are used, not where they are implemented.**
The `application/ports` layer owns all interface definitions.
The `infrastructure` layer contains only concrete implementations — never interface definitions.

**Depend only on the methods you actually call.**
Use cases and adapters must declare only the interface surface they need — never a broader one.

## Rules

### 0. Interface Ownership: `application/ports` Defines, `infrastructure` Implements

All interfaces for external dependencies (blockchain APIs, repositories, file storage) MUST be defined in `internal/application/ports/`.
The `infrastructure/` layer MUST NOT define any interfaces — only concrete structs that implement them.

```
internal/application/ports/api/btc/interface.go  ← interface definitions live here
internal/infrastructure/api/btc/bitcoin.go        ← concrete implementation only
```

```go
// ❌ BAD: interface defined in infrastructure layer
// internal/infrastructure/api/btc/bitcoin.go
type BitcoinAPI interface {
    GetBalance() (float64, error)
}

// ✅ GOOD: interface defined in application/ports, infrastructure only implements
// internal/application/ports/api/btc/interface.go
type BalanceChecker interface {
    GetBalance() (btcutil.Amount, error)
}

// internal/infrastructure/api/btc/bitcoin.go
type Bitcoin struct { ... }
func (b *Bitcoin) GetBalance() (btcutil.Amount, error) { ... } // implements BalanceChecker
```

### 1. Never Use `Bitcoiner` or `BCHer` Directly

`apibtc.Bitcoiner` and `apibtc.BCHer` are legacy monolithic interfaces.
They must NEVER appear as parameters or field types in use cases or adapters.

```go
// ❌ BAD: depends on the full monolithic interface
type myUseCase struct {
    btc apibtc.Bitcoiner
}

// ✅ GOOD: depends on only what is needed
type myUseCase struct {
    btc apibtc.MultisigManager
}
```

### 2. No Single-Interface Wrappers

A local type that simply re-wraps one existing interface adds no value. Use the named interface directly.

```go
// ❌ BAD: redundant wrapper
type sendTxBCHClient interface {
    apibtc.TransactionSender
}

// ✅ GOOD: use the named interface directly
type sendTransactionUseCase struct {
    bchClient apibtc.TransactionSender
}
```

### 3. No Duplicate Local Interfaces Across Packages

If the same interface combination is needed in two or more files, promote it to a named interface in `internal/application/ports/api/btc/interface.go`.

```go
// ❌ BAD: identical definition in keygen/btc/ AND sign/btc/
type importPrivKeyBTCClient interface {
    apibtc.ChainConfigProvider
    apibtc.PrivateKeyImporter
    apibtc.AddressOperator
}

// ✅ GOOD: promote to interface.go, use everywhere
// In interface.go:
type PrivateKeyImportClient interface {
    ChainConfigProvider
    PrivateKeyImporter
    AddressOperator
}
// In use case files:
func NewImportPrivateKeyUseCase(btc apibtc.PrivateKeyImportClient, ...) { ... }
```

### 4. When to Define a Local Interface

Define a local interface **only** when:

- It combines 2+ small interfaces that do **not** yet form a named composite in `apibtc`
- **AND** it is used in exactly one file (not duplicated)

Once two files share the same combination, promote it to `interface.go`.

```go
// ✅ ACCEPTABLE: unique combination used only in this file
type createTxBCHClient interface {
    apibtc.TransactionCreationDeps
    apibtc.SomeOtherInterface
}
```

### 5. Use Named Composed Interfaces from `interface.go`

When a composed interface already exists in `apibtc`, use it — do not redefine locally.

| Named interface | Composed from |
|---|---|
| `BTCTransactionSigner` | `ChainConfigProvider` + `PSBTSigner` + `PSBTHandler` |
| `BCHTransactionSigner` | `RawTransactionConverter` + `RawTransactionSigner` |
| `TransactionCreationDeps` | `ChainConfigProvider` + `AmountConverter` + `UTXOProvider` + `RawTransactionCreator` + `AddressOperator` + `BalanceChecker` |
| `PrivateKeyImportClient` | `ChainConfigProvider` + `PrivateKeyImporter` + `AddressOperator` |
| `AddressImportClient` | `ChainConfigProvider` + `LegacyAddressImporter` + `AddressOperator` |

### 6. DI Layer and Infrastructure Use `Bitcoiner`

The DI container and infrastructure factories use `apibtc.Bitcoiner` (the full monolithic interface).
`Bitcoiner` is the ONLY place where the full interface is permitted — all other layers use small focused interfaces.
Use cases never receive `Bitcoiner` — the DI layer narrows it to the required small interface.

```go
// DI container (internal/di/container.go) — Bitcoiner is fine here
btc apibtc.Bitcoiner

// Use case constructor — only what is needed
func NewSignTransactionUseCase(btc apibtc.BTCTransactionSigner, ...) { ... }
```

### 7. Where to Add New Composed Interfaces

New composed interfaces belong in:

```
internal/application/ports/api/btc/interface.go
```

Place them in the **"Composed Interfaces for Common Use Cases"** section with a descriptive comment explaining which use cases they serve.

## Decision Flowchart

```
Need an interface for a use case or adapter?
│
├─ Is there already a named interface in apibtc that covers exactly the needed methods?
│   └─ YES → Use it directly. Done.
│
├─ Is this combination used in 2+ files?
│   └─ YES → Add a named composed interface to interface.go. Use it everywhere.
│
├─ Is it a wrapper around exactly 1 existing interface?
│   └─ YES → Use that interface directly. No wrapper needed.
│
└─ Only 1 file, combines 2+ interfaces not yet composited → Define local interface in that file.
    When a second file needs the same combo → promote to interface.go.
```

## Related Files

- `internal/application/ports/api/btc/interface.go` — all small + composed BTC/BCH interfaces
- `.claude/rules/internal/application-layer.md` — application layer rules
- `.claude/rules/internal/clean-architecture.md` — layer dependency rules
