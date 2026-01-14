# BTC E2E Scripts

This directory contains Bitcoin E2E (End-to-End) test scripts.
Each script automates the complete workflow from infrastructure setup to transaction execution.

## Documentation

For detailed transaction pattern explanations, technical references, and implementation status, see:

- **[E2E Transaction Patterns Guide](../../../../docs/crypto/btc/e2e_transaction_patterns.md)** - Key types, signature patterns, and workflow details

## Script List

| Script | Pattern | Signature Requirement | Address Format |
|--------|---------|----------------------|----------------|
| `e2e-p2pkh-singlesig.sh` | P2PKH Single-sig (Pattern 1) | Single-sig | `1...` / `m...` |
| `e2e-p2pkh-2of3.sh` | P2PKH 2-of-3 Multisig (Pattern 2) | 2-of-3 | `3...` / `2...` |
| `e2e-p2sh-p2wsh-3of3.sh` | P2SH-P2WSH 3-of-3 Multisig (Pattern 8) | 3-of-3 | `3...` / `2...` |

## Usage

### Basic Execution

```bash
# Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh

# 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh

# 3-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh
```

### Make Targets

```bash
# Single-sig
make btc-e2e-p2pkh-singlesig

# 2-of-3 Multisig
make btc-e2e-p2pkh-2of3

# 3-of-3 Multisig
make btc-e2e-p2sh-p2wsh-3of3
```

### Common Options

| Option | Description |
|--------|-------------|
| `--cleanup` | Stop containers and cleanup state |
| `--reset` | Full reset and run from scratch |
| `--verbose` | Enable verbose output |
| `--non-interactive` | Run without prompts (for CI/CD) |
| `-h, --help` | Display help message |

## Required Configuration

Each script requires matching `address_type` in the corresponding config files:

### Pattern 1: P2PKH Single-sig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml
address_type: "legacy"
```

### Pattern 2: P2PKH 2-of-3 Multisig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml, btc_sign1.yaml, btc_sign2.yaml
address_type: "legacy"
```

### Pattern 8: P2SH-P2WSH 3-of-3 Multisig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml, btc_sign1.yaml, btc_sign2.yaml
address_type: "p2sh-segwit"
```

## Environment Variables

```bash
# RPC credentials (defaults are for regtest/development only)
RPC_USER=xyz
RPC_PASSWORD=xyz
```

## Related Documentation

- [E2E Transaction Patterns Guide](../../../../docs/crypto/btc/e2e_transaction_patterns.md) - Pattern details
- [BTC Technical Reference](../../../../docs/crypto/btc/README.md) - Bitcoin technical reference
- [Descriptor Examples](../../../../docs/crypto/btc/descriptor_examples.md) - Descriptor examples
- [PSBT Developer Guide](../../../../docs/crypto/btc/psbt_developer_guide.md) - PSBT developer guide
