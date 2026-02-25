---
paths: ["**/*.go"]
---

# Go Type Alias Rules

## Overview

Rules for using type aliases and type definitions in Go files.

## Definitions

Go has two distinct constructs:

| Syntax | Name | Description |
|--------|------|-------------|
| `type X = Y` | **Type alias** | X and Y are exactly the same type (interchangeable) |
| `type X Y` | **Type definition** | X is a new distinct type based on Y |

Both are subject to the rules below.

## Rules

### 1. Aliasing Self-Defined Types Is Prohibited

Creating a type alias for a type that is already defined within this project is prohibited.
Use the original type directly.

```go
// ❌ BAD: aliasing a self-defined type adds no value and creates confusion
type MyBitcoiner = apibtc.Bitcoiner
type MyAccountType = domainAccount.AccountType
type MyAddressList = []domainAddress.Address

// ✅ GOOD: use the original type directly
var client apibtc.Bitcoiner
var acct domainAccount.AccountType
var addrs []domainAddress.Address
```

**Exception**: Auto-generated code (e.g., protobuf, SQLC) may use type aliases. Do NOT edit those files.

### 2. Type Aliases and Definitions on Primitives Are Allowed

Both `type X = Primitive` and `type X Primitive` are **allowed and encouraged** when the underlying type is a primitive (`string`, `int`, `uint32`, `int64`, `float64`, `bool`, `byte`) and it adds domain semantics.

```go
// ✅ GOOD: adds type safety and domain meaning
type CoinTypeCode string
type AccountType = string
type AddrStatus int
type WalletType int
type AuthType int
```

These prevent accidental misuse (e.g., passing a raw `string` where a `CoinTypeCode` is expected).

### 3. Import Path Aliases Are Discouraged

Package import aliases (`import foo "github.com/.../bar"`) should be avoided unless:

- The package name conflicts with a local variable or another import
- The package name is ambiguous (e.g., two packages named `btc`)
- The auto-generated package name is excessively long

```go
// ❌ BAD: alias is unnecessary, package name is already clear
import btcapi "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"

// ✅ GOOD: use the declared package name directly
import "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
// then: btc.Bitcoiner, btc.TransactionSender, etc.

// ✅ ACCEPTABLE: alias is needed to resolve a name collision
import (
    apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
    infra  "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc"
)
```

## Summary

| Case | Allowed |
|------|---------|
| `type X = ProjectType` (self-defined) | ❌ No |
| `type X ProjectType` (self-defined) | ❌ No |
| `type X = string` | ✅ Yes (with domain meaning) |
| `type X = int` | ✅ Yes (with domain meaning) |
| `type X string` | ✅ Yes (with domain meaning) |
| `type X int` | ✅ Yes (with domain meaning) |
| Import alias (no collision) | ❌ Avoid |
| Import alias (resolves collision) | ✅ Acceptable |

## Related Rules

- `.claude/rules/go/conventions.md` - General Go conventions
- `.claude/rules/go/internal-interface.md` - Interface definition rules
