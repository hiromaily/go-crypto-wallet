# Descriptor Development Guide

Guidance for contributors extending descriptor functionality.

## Prerequisites
- Go toolchain (see `agents/requirements.md`).
- Bitcoin Core available for compatibility testing when running `integration`-tag tests or smoke scripts.
- Avoid editing auto-generated files; descriptors live in application/infrastructure layers.

## Local Verification

1. **Unit tests**
   ```bash
   go test ./...
   ```

2. **Integration workflow test (no external deps)**
   ```bash
   go test ./test/integration -run TestDescriptorWorkflow_EndToEnd -v
   ```

3. **Bitcoin Core compatibility (requires node)**
   ```bash
   BTC_CORE_COMPAT_DESCRIPTOR_FILE=/path/to/descriptors.json \
   BITCOIN_CLI_ARGS="-regtest -rpcuser=test -rpcpassword=test" \
   go test -tags=integration ./test/integration -run TestDescriptorCompatibility_ImportAndDerive -v
   ```

4. **Performance benchmarks**
   ```bash
   go test ./internal/infrastructure/api/bitcoin/btc -bench=BenchmarkGenerateDescriptors -benchmem
   ```

## Manual Smoke Test

```bash
scripts/test/test_descriptor_compatibility.sh ./descriptors.json
```

Useful when a full Bitcoin Core instance is available; reports import/derive results quickly.

## Security Checklist

- Never log or write extended private keys.
- Validate descriptor checksums and derivation paths before import.
- Keep descriptor backups encrypted and offline.
- Multisig imports remain intentionally disabled in the watch wallet until hardened.

## Extending Functionality

- Add new descriptor types in the infrastructure layer (`DescriptorService`, `DescriptorParser`) and domain enums.
- Wire through use cases and CLI commands via dependency injection (see `internal/di/container.go`).
- Update docs: user guide, API reference, and migration notes when behavior changes.
- Add unit tests for parsing/formatting and an integration test demonstrating the end-to-end flow.
