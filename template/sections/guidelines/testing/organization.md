### Test Organization

**File Naming:**

- Test files: `*_test.go` (same package)
- Integration tests: `*_integration_test.go` with `//go:build integration` tag

**Package Organization:**

```text
internal/domain/account/
├── account.go           # Domain code
└── account_test.go      # Unit tests

internal/infrastructure/repository/watch/
├── repository.go                      # Implementation
├── repository_test.go                 # Unit tests (mocked database)
└── repository_integration_test.go     # Integration tests (real database)
```

**Integration Test Tags:**

```go
//go:build integration

package repository_test

import "testing"

func TestRepository_Integration(t *testing.T) {
    // Integration test with real database
}
```
