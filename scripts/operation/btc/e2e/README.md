# BTC E2E Scripts

This directory contains Bitcoin E2E (End-to-End) test scripts.
Each script automates the complete workflow from infrastructure setup to transaction execution.

## Documentation

For detailed transaction pattern explanations, technical references, and implementation status, see:

- **[E2E Transaction Patterns Guide](../../../../docs/crypto/btc/e2e_transaction_patterns.md)** - Key types, signature patterns, and workflow details

## Script List

| Script | Pattern | Signature Requirement | Address Format |
|--------|---------|----------------------|----------------|
| `e2e-p1-p2pkh-singlesig.sh` | P2PKH Single-sig (Pattern 1) | Single-sig | `1...` / `m...` |
| `e2e-p2-p2pkh-2of3.sh` | P2PKH 2-of-3 Multisig (Pattern 2) | 2-of-3 | `3...` / `2...` |
| `e2e-p3-p2sh-p2wpkh-singlesig.sh` | P2SH-P2WPKH Single-sig (Pattern 3) | Single-sig | `3...` / `2...` |
| `e2e-p4-p2sh-p2wsh-2of3.sh` | P2SH-P2WSH 2-of-3 Multisig (Pattern 4) | 2-of-3 | `3...` / `2...` |
| `e2e-p5-p2wpkh-singlesig.sh` | P2WPKH Native SegWit Single-sig (Pattern 5) | Single-sig | `bc1q...` / `bcrt1q...` |
| `e2e-p6-p2wsh-2of3.sh` | P2WSH Native SegWit 2-of-3 Multisig (Pattern 6) | 2-of-3 | `bc1q...` / `bcrt1q...` |
| `e2e-p7-p2wsh-3of3.sh` | P2WSH Native SegWit 3-of-3 Multisig (Pattern 7) | 3-of-3 | `bc1q...` / `bcrt1q...` |
| `e2e-p8-p2sh-p2wsh-3of3.sh` | P2SH-P2WSH 3-of-3 Multisig (Pattern 8) | 3-of-3 | `3...` / `2...` |

## Usage

### Basic Execution

```bash
# Pattern 1: Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh

# Pattern 2: 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh

# Pattern 3: P2SH-P2WPKH Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh

# Pattern 4: P2SH-P2WSH 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh

# Pattern 5: P2WPKH Native SegWit Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh

# Pattern 6: P2WSH Native SegWit 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh

# Pattern 7: P2WSH Native SegWit 3-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh

# Pattern 8: P2SH-P2WSH 3-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh
```

### Make Targets

```bash
# Pattern 1: Single-sig
make btc-e2e-p1

# Pattern 2: 2-of-3 Multisig
make btc-e2e-p2

# Pattern 3: P2SH-P2WPKH Single-sig
make btc-e2e-p3

# Pattern 4: P2SH-P2WSH 2-of-3 Multisig
make btc-e2e-p4

# Pattern 5: P2WPKH Native SegWit Single-sig
make btc-e2e-p5

# Pattern 6: P2WSH Native SegWit 2-of-3 Multisig
make btc-e2e-p6

# Pattern 7: P2WSH Native SegWit 3-of-3 Multisig
make btc-e2e-p7

# Pattern 8: P2SH-P2WSH 3-of-3 Multisig
make btc-e2e-p8
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

### Pattern 3: P2SH-P2WPKH Single-sig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml
address_type: "p2sh-segwit"
```

### Pattern 4: P2SH-P2WSH 2-of-3 Multisig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml, btc_sign1.yaml, btc_sign2.yaml
address_type: "p2sh-segwit"
```

### Pattern 5: P2WPKH Native SegWit Single-sig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml
address_type: "bech32"
```

### Pattern 6: P2WSH Native SegWit 2-of-3 Multisig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml, btc_sign1.yaml, btc_sign2.yaml
address_type: "bech32"
```

### Pattern 7: P2WSH Native SegWit 3-of-3 Multisig

```yaml
# config/wallet/btc_watch.yaml, btc_keygen.yaml, btc_sign1.yaml, btc_sign2.yaml
address_type: "bech32"
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

# MySQL credentials (defaults are for regtest/development only)
MYSQL_ROOT_PASSWORD=root
```

## Related Documentation

- [E2E Transaction Patterns Guide](../../../../docs/crypto/btc/e2e_transaction_patterns.md) - Pattern details
- [BTC Technical Reference](../../../../docs/crypto/btc/README.md) - Bitcoin technical reference
- [Descriptor Examples](../../../../docs/crypto/btc/descriptor_examples.md) - Descriptor examples
- [PSBT Developer Guide](../../../../docs/crypto/btc/psbt_developer_guide.md) - PSBT developer guide
