---
paths: ["**/*test.go"]
---

# Go Test Rules

## Overview

Rules for writing and modifying test files (`*_test.go`) in go-crypto-wallet.

**Full reference**: @docs/guidelines/testing.md

## Package Naming

| Test type | Package declaration | When |
|---|---|---|
| Unit test | `package btc` | Same package — direct access to unexported symbols |
| Integration test | `package btc_test` | External package — tests that depend on a database or other external infrastructure |

```go
// unit_test.go — same package
package btc

// integration_test.go — external package + build tag
//go:build integration

package btc_test
```

## Test Framework — Mandatory

Always use [testify](https://github.com/stretchr/testify). **Never** use `t.Errorf`, `t.Fatalf`, or `reflect.DeepEqual` directly.

| Package | When to use |
|---------|-------------|
| `require` | Fatal assertion — test stops immediately on failure |
| `assert` | Non-fatal assertion — test continues after failure |

```go
// ✅ GOOD
result, err := useCase.Execute(ctx, input)
require.NoError(t, err)          // stops if error (subsequent code depends on result)
assert.Equal(t, expected, result.TxID)
assert.Contains(t, result.Message, "success")

// ❌ BAD
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

**Rule**: Use `require` when subsequent test code depends on the assertion passing; use `assert` otherwise.

## Mock Usage — Mandatory

**Always use mockery-generated mocks. Never write mock structs by hand.**

See @.claude/rules/internal/mockery.md for the full placement convention and exceptions.

```go
// ✅ GOOD: generated mock
import coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"

func TestFoo(t *testing.T) {
    repo := coldmocks.NewMockBTCAccountKeyRepositorier(t) // auto-registers cleanup
    repo.EXPECT().
        GetOneMaxID(domainAccount.AccountTypeDeposit).
        Return(key, nil)
    // AssertExpectations called automatically — do NOT call it manually
}

// ❌ BAD: hand-written mock struct
type stubRepo struct{}
func (s *stubRepo) GetOneMaxID(accountType AccountType) (*Key, error) { ... }
```

**Import alias convention**:

| Mocks package | Import alias |
|---|---|
| `infrastructure/api/btc/mocks` | `btcapiamocks` |
| `infrastructure/api/eth/mocks` | `ethapiamocks` |
| `infrastructure/api/xrp/mocks` | `xrpapiamocks` |
| `infrastructure/repository/cold/mocks` | `coldmocks` |
| `infrastructure/repository/watch/mocks` | `repomocks` |
| `infrastructure/storage/file/transaction/mocks` | `storagemocks` |

## Helper Functions

Always call `t.Helper()` in helper functions so failures point to the call site:

```go
func newTestDependencies(t *testing.T) *testDependencies {
    t.Helper()  // ← required
    return &testDependencies{
        repo: coldmocks.NewMockBTCAccountKeyRepositorier(t),
    }
}
```

## Subtests and Table-Driven Tests

Use `t.Run` for subtests and table-driven tests:

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

## Parallelism

Add `t.Parallel()` to tests that have no shared mutable state:

```go
func TestPureDomainLogic(t *testing.T) {
    t.Parallel()
    // pure computation, no shared state
}
```

Do NOT add `t.Parallel()` to tests that share global state or use the same database connections.

## Skipping Tests

Use `t.Skip` with a clear reason when a test requires external infrastructure:

```go
func TestWithDatabase(t *testing.T) {
    t.Skip("Skipping until DB transaction is mocked — requires comprehensive DB mock setup")
}
```

For integration tests, use the build tag instead of `t.Skip`:

```go
//go:build integration

package repository_test
```

## Temporary Files and Directories

Use `t.TempDir()` — it cleans up automatically after the test:

```go
tmpDir := t.TempDir()    // auto-removed after test
tmpFile, err := os.CreateTemp(t.TempDir(), "prefix-*.ext")
require.NoError(t, err)
```

## What NOT to Do

```go
// ❌ Manual AssertExpectations — automatic via NewMock*(t)
mock.AssertExpectations(t)

// ❌ Hand-written mock structs for ports interfaces
type MockFoo struct { mock.Mock }
func (m *MockFoo) Method() error { ... }

// ❌ reflect.DeepEqual instead of assert.Equal
if !reflect.DeepEqual(expected, got) { t.Error(...) }

// ❌ Standard library assertions
t.Errorf("got %v, want %v", got, expected)

// ❌ Sleeping for timing
time.Sleep(100 * time.Millisecond)
```

## Verification

Before committing test-only changes:

```bash
make go-lint        # Required: zero lint errors
make gotest         # Required: all tests pass
```

## Related Rules

- @.claude/rules/internal/mockery.md — mock generation placement convention
- @.claude/rules/go/conventions.md — Go coding conventions
- @docs/guidelines/testing.md — full testing guidelines (SSOT)
