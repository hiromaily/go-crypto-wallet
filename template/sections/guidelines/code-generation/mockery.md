### Mock Code (Mockery)

**Tool**: [mockery v3](https://github.com/vektra/mockery)
**Source**: Interface definitions in Go files
**Configuration**: `.mockery.yaml`
**Command**: `make mockery` (or `go tool github.com/vektra/mockery/v3`)

**Generated Files**:

- `internal/infrastructure/api/btc/mocks/mock_bitcoiner.go` - Bitcoin API mock
- `internal/infrastructure/repository/mocks/mock_*.go` - Repository interface mocks
- `internal/infrastructure/storage/file/transaction/mocks/mock_transaction_file_repositorier.go` - File storage mock

**Mock Directory Structure**:

Mocks are placed in `mocks/` subdirectories alongside their implementations:

```text
internal/infrastructure/
├── api/bitcoin/
│   ├── btc/bitcoin.go          # Implementation
│   └── mocks/mock_bitcoiner.go # Generated mock
├── repository/
│   └── mocks/mock_*.go         # Persistence interface mocks
└── storage/file/
    └── transaction/
        ├── transaction.go      # Implementation
        └── mocks/
            └── mock_transaction_file_repositorier.go  # Storage interface mocks
```

**Adding New Mocks**:

1. Edit `.mockery.yaml`
2. Add the interface under the appropriate package section
3. Run `make mockery`

**Moving Mocks Directories**:

⚠️ **IMPORTANT**: When moving implementation code that has associated mocks, you **MUST** also update `.mockery.yaml` to reflect the new directory structure.

**Steps when moving mocks directories**:

1. Move the implementation code to the new location
2. **Update `.mockery.yaml`** - Change the `dir` path in the configuration for the affected interface(s)
3. Run `make mockery` - This will automatically:
   - Clean all existing mocks (via `clean-mocks` dependency)
   - Regenerate mocks in the new location based on updated `.mockery.yaml`

**Example**: If moving `internal/infrastructure/storage/file/transaction.go` to `internal/infrastructure/storage/file/transaction/transaction.go`:

```yaml
# Before
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/mocks"
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:

# After
github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage:
  config:
    dir: "internal/infrastructure/storage/file/transaction/mocks"  # Updated path
    pkgname: "mocks"
  interfaces:
    TransactionFileRepositorier:
```

**Note**: The `make mockery` target automatically runs `clean-mocks` first, so old mocks will be removed before generating new ones. This ensures no stale mocks remain when paths change.

**Note**: See [Testing Guidelines](../../../../docs/guidelines/testing.md) for mock usage examples and best practices.
