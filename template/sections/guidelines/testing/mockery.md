### Mock Generation with Mockery

> **Code Generation**: For mock generation commands and configuration, see [Code Generation Guide](../../../../docs/guidelines/code-generation.md#mock-code-mockery).

This project uses [mockery v3](https://github.com/vektra/mockery) to generate mock implementations from Go interfaces.

#### Configuration

Mock generation is configured in `.mockery.yaml` at the project root.

**Key Settings:**

- `all: false` - Only generate mocks for explicitly listed interfaces
- `template: testify` - Generate testify-compatible mocks with `EXPECT()` support
- Mocks are placed in `mocks/` subdirectories alongside implementations

#### Mock Directory Structure

```text
internal/infrastructure/
├── api/bitcoin/
│   ├── btc/bitcoin.go              # Implementation
│   └── mocks/
│       └── mock_bitcoiner.go       # Generated mock
├── repository/
│   ├── watch/repository.go         # Implementation
│   └── mocks/
│       └── mock_*.go               # Generated mocks for persistence interfaces
└── storage/file/
    └── transaction/
        ├── transaction.go          # Implementation
        └── mocks/
            └── mock_transaction_file_repositorier.go
```

#### Commands

```bash
# Generate all mocks
make mockery

# Clean all generated mocks
make clean-mocks

# Regenerate mocks (clean + generate)
make clean-mocks && make mockery
```

#### Adding New Mock Interfaces

To add a new interface for mock generation:

1. Edit `.mockery.yaml`
2. Add the interface under the appropriate package:

```yaml
packages:
  github.com/hiromaily/go-crypto-wallet/internal/your/package:
    config:
      dir: "internal/your/package/mocks"
      pkgname: "mocks"
    interfaces:
      YourInterface:
```

1. Run `make mockery`

#### Moving Mocks Directories

⚠️ **IMPORTANT**: When refactoring code and moving implementation files that have associated mocks, you **MUST** update `.mockery.yaml` to reflect the new directory structure.

**Why this matters**:

- Mockery generates mocks based on the `dir` path specified in `.mockery.yaml`
- If you move implementation code but don't update the configuration, mocks will be generated in the wrong location
- The `make mockery` target automatically cleans all mocks before generating, so old mocks will be removed, but new ones won't be created in the correct location if the config is wrong

**Steps when moving mocks directories**:

1. Move the implementation code to the new location
2. **Update `.mockery.yaml`** - Change the `dir` path for the affected interface(s)
3. Update any import paths in code that reference the old mock location
4. Run `make mockery` to regenerate mocks in the new location

**Example**: Moving transaction file repository:

```yaml
# .mockery.yaml - Before
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/mocks"
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:

# .mockery.yaml - After (code moved to transaction/ subdirectory)
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/transaction/mocks"  # Updated!
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:
```

**Note**: The `make mockery` target has `clean-mocks` as a dependency, so it will automatically remove all existing mocks before generating new ones. This ensures no stale mocks remain when paths change.

#### Using Generated Mocks

```go
import (
    "testing"

    bitcoinmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/mocks"
    repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/mocks"
    storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

func TestWithMocks(t *testing.T) {
    // Create mock (automatically registers cleanup with t.Cleanup)
    mockClient := bitcoinmocks.NewMockBitcoiner(t)

    // Set expectations with EXPECT()
    mockClient.EXPECT().
        GetBlockCount().
        Return(int64(100), nil)

    // Use mock in test
    result, err := mockClient.GetBlockCount()
    // Expectations are automatically verified at test end
}
```

#### Mock Best Practices

**Do:**

- Pass `t *testing.T` to mock constructors for automatic cleanup
- Use `EXPECT()` for type-safe expectation setting
- Set expectations before calling the code under test
- Keep mock setups minimal and focused

**Don't:**

- Don't manually verify expectations (automatic with `t`)
- Don't create mocks without passing `t`
- Don't over-mock (mock only direct dependencies)
