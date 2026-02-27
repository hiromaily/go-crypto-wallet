---
paths: ["**/*test.go"]
---

# Go Test Rules

## Overview

Rules for writing and modifying test files (`*_test.go`) in go-crypto-wallet.

**Full reference (SSOT)**: @docs/guidelines/testing.md

---

## Test File Organization

| Test type | Filename | Package | Build tag |
|-----------|----------|---------|-----------|
| Unit test | `*_test.go` | same package (e.g. `package btc`) | none |
| Integration test (DB / external) | `*_integration_test.go` | external (e.g. `package btc_test`) | `//go:build integration` |
| E2E / complex multi-service | `*_test.go` in `internal/integration_test/` | `package integration_test` | `//go:build integration` |

```go
// unit_test.go — same package
package btc

// repository_integration_test.go — external package + build tag
//go:build integration

package postgres_test
```

---

## Test Framework — Mandatory

Always use [testify](https://github.com/stretchr/testify). **Never** use `t.Errorf`, `t.Fatalf`, or `reflect.DeepEqual` directly.

| Package | When to use |
|---------|-------------|
| `require` | Fatal — test stops immediately on failure |
| `assert` | Non-fatal — test continues after failure |

**Rule**: Use `require` when subsequent test code depends on the assertion passing; use `assert` otherwise.

---

## Table-Driven Tests — Mandatory

Always use `t.Run` with a slice of test cases:

```go
func TestFoo(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "ok", false},
        {"empty", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Foo(tt.input)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

---

## Mock Usage — Mandatory

**Always use mockery-generated mocks. Never write mock structs by hand.**

See @.claude/rules/internal/mockery.md for placement convention and exceptions.

**Import alias convention**:

| Mocks package | Import alias |
|---------------|--------------|
| `infrastructure/api/btc/mocks` | `btcapiamocks` |
| `infrastructure/api/eth/mocks` | `ethapiamocks` |
| `infrastructure/api/xrp/mocks` | `xrpapiamocks` |
| `infrastructure/repository/cold/mocks` | `coldmocks` |
| `infrastructure/repository/watch/mocks` | `repomocks` |
| `infrastructure/storage/file/transaction/mocks` | `storagemocks` |

---

## Helper Functions

Always call `t.Helper()` in helper functions so failures point to the call site.

---

## Test Utilities (`testutil/`)

**Rule**: A global `pkg/testutil` package is **prohibited**. Co-locate test utilities with the package they support.

```
pkg/db/testutil/       ← helpers for pkg/db (e.g. OpenTestDB)
internal/foo/testutil/ ← helpers for internal/foo
```

Create a `testutil/` subdirectory inside the package the utilities belong to. Never create a standalone `pkg/testutil/` or `internal/testutil/` package.

---

## What NOT to Do

```go
// ❌ Standard library assertions — use testify
t.Errorf("got %v, want %v", got, expected)

// ❌ Hand-written mock structs — use mockery
type stubRepo struct{}

// ❌ Manual AssertExpectations — automatic via NewMock*(t)
mock.AssertExpectations(t)

// ❌ reflect.DeepEqual — use assert.Equal
if !reflect.DeepEqual(expected, got) { t.Error(...) }

// ❌ Global testutil package
import "github.com/hiromaily/go-crypto-wallet/pkg/testutil"
```

---

## Parallelism

Add `t.Parallel()` to tests with no shared mutable state. Do **not** add it to tests sharing DB connections.

---

## Verification

```bash
make go-lint   # Required: zero lint errors
make go-test   # Required: all unit tests pass
```

---

## Related Rules

- @.claude/rules/internal/mockery.md — mock generation placement convention
- @.claude/rules/go/conventions.md — Go coding conventions
- @docs/guidelines/testing.md — full testing guidelines (SSOT)
