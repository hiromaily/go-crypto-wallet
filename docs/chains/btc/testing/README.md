# Testing

This directory contains test procedures and verification documentation.

## Contents

| Document | Description |
|----------|-------------|
| [pattern3-verification.md](pattern3-verification.md) | Pattern 3 transaction verification procedures |

## Related Testing Documentation

- [../taproot/testing.md](../taproot/testing.md) - Taproot-specific tests
- [../operations/e2e-transaction-patterns.md](../operations/e2e-transaction-patterns.md) - E2E test patterns

## Test Scripts

E2E test scripts are located in:

- `scripts/operation/btc/e2e/` - Bitcoin E2E test scripts

## Running Tests

```bash
# Run Go unit tests
make gotest

# Run BTC E2E tests
make btc-e2e-setup
make btc-e2e-pattern1  # Single-sig
make btc-e2e-pattern2  # Traditional multisig
# ... etc
```

## Related Documentation

- [../README.md](../README.md) - Main BTC documentation index
- [../../../guidelines/testing.md](../../../guidelines/testing.md) - Project testing standards
