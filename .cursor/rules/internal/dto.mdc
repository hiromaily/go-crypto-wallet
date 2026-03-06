---
description: DTO and Data Definition Placement Rules
globs: ["internal/**/*.go"]
alwaysApply: false
---


# DTO and Data Definition Placement Rules

## Overview

Rules for where to define data types (DTOs, wire types, presentation types) and where to place
conversion logic in the Clean Architecture layers of `internal/`.

## Layer Summary

```
Interface Adapters  → defines: presentation input/output types (CLI args, HTTP body, gRPC message)
                      converts: presentation → application DTO before calling use case

Application Layer   → defines: canonical cross-boundary DTOs (internal/application/dto/)
                      port interfaces use ONLY these types in their signatures

Infrastructure      → defines: internal wire/DB types (JSON-tagged structs, ORM models)
                      converts: wire types ↔ application DTOs at the boundary

Domain Layer        → defines: domain entities, value objects (with business logic — NOT DTOs)
                      no conversion, no DTOs
```

## Rules

### 1. Application DTOs Are the Canonical Cross-Boundary Types

All types used in `application/ports/` interface signatures MUST be defined in `internal/application/dto/`.
Infrastructure types MUST NOT appear in port interface signatures.

```go
// ✅ GOOD: port interface uses application DTO
// internal/application/ports/api/xrp/interface.go
type TransactionPreparer interface {
    CreateRawTransaction(
        ctx context.Context,
        senderAccount, receiverAccount string,
        amount float64,
        instructions *dtoxrp.Instructions,  // application DTO
    ) (*dtoxrp.TxInput, string, error)       // application DTO
}

// ❌ BAD: port interface returns an infrastructure package type
type KeyGenerator interface {
    WalletPropose(ctx context.Context, passphrase string) (*xrprpcadmin.ResponseWalletPropose, error)
    //                                                      ^^^^^^^^^ infrastructure type in port
}
```

### 2. Infrastructure Defines Wire Types Locally, Converts at the Boundary

Infrastructure packages may define their own internal types optimised for wire format (JSON tags,
protocol buffer fields, etc.). These must be converted to/from application DTOs before crossing
the boundary.

Conversion functions must be **private** unless explicitly needed by a sibling package.

```go
// internal/infrastructure/api/xrp/public/transaction.go

// Internal wire type — JSON tags, infra-specific detail
type TxInput struct {
    TransactionType string `json:"TransactionType"`
    Account         string `json:"Account"`
    // ...
}

// Private converters — live in the infrastructure package, not in application/
func toInfraTxInput(dto *dtoxrp.TxInput) *TxInput { ... }
func toDTOTxInput(infra *TxInput) *dtoxrp.TxInput  { ... }

// Public method satisfies the port interface — returns application DTO
func (p *PublicXRP) CreateRawTransaction(...) (*dtoxrp.TxInput, string, error) {
    infra := buildTxInput(...)
    return toDTOTxInput(infra), jsonStr, nil
}
```

### 3. Interface Adapters Convert Presentation Input Before Calling Use Cases

CLI flags, HTTP request bodies, gRPC messages are presentation details. Handlers must convert
them to application DTOs and never pass raw presentation types into use cases.

```go
// internal/interface-adapters/cli/watch/api/xrp/sendcoin.go

func (h *SendCoinHandler) Handle(args CLIArgs) error {
    instructions := &dtoxrp.Instructions{  // convert CLI args → application DTO
        Fee: args.Fee,
    }
    return h.useCase.Execute(ctx, args.Sender, args.Receiver, args.Amount, instructions)
}
```

### 4. Domain Types Are NOT DTOs

Domain entities and value objects contain business logic (validation, invariants). They are
not plain data carriers. Do not use them as DTOs across layer boundaries.

| Domain type | Used as DTO? |
|-------------|-------------|
| `domain/address.Address` | No — pass `string` or an application DTO |
| `domain/coin.CoinTypeCode` | Acceptable as a value type if it has no infra deps |
| `domain/transaction.Transaction` | No — map to an application DTO |

### 5. Where Conversion Code Lives

The **outer** layer is always responsible for conversion. This keeps inner layers clean.

| Conversion direction | Where the conversion code lives |
|---|---|
| Infra wire type → Application DTO | Infrastructure package (`toDTOXxx` functions) |
| Application DTO → Infra wire type | Infrastructure package (`toInfraXxx` functions) |
| Presentation input → Application DTO | Interface-adapters handler/controller |
| Application DTO → Presentation output | Interface-adapters handler/controller |

### 6. DTO Placement by Sub-Directory

```
internal/application/dto/
├── address.go          # shared address DTOs (all chains)
├── address_format.go   # address format helpers
├── btc/
│   ├── dto.go          # BTC-specific DTOs (AddressInfo, etc.)
│   └── fullpubkey.go   # BTC public key DTOs
├── eth/
│   └── transaction_file.go  # ETH file-based transaction DTOs
└── xrp/
    ├── transaction.go       # XRP transaction DTOs (TxInput, Instructions, ...)
    ├── transaction_file.go  # XRP file-based transaction DTOs
    └── response.go          # XRP response DTOs
```

Add new DTOs to the appropriate chain sub-package. Chain-agnostic DTOs go in the root `dto/` package.

## Quick Checklist

- [ ] Port interface method signatures use only `internal/application/dto/` types (or primitives/domain value types)
- [ ] Infrastructure wire types are private to their package
- [ ] Conversion functions (`toDTOXxx`, `toInfraXxx`) live in the infrastructure package
- [ ] Handlers in interface-adapters convert presentation input to DTOs before calling use cases
- [ ] No infrastructure package types appear as return types in `application/ports/` interfaces

## Related Rules

- @.claude/rules/internal/clean-architecture.md - Layer dependency rules
- @.claude/rules/internal/application-layer.md - Application layer specifics
- @.claude/rules/internal/infrastructure-layer.md - Infrastructure layer specifics
- @.claude/rules/internal/interface-adapters-layer.md - Interface adapters layer specifics
- @.claude/rules/internal/internal-interface.md - Interface definition placement rules
