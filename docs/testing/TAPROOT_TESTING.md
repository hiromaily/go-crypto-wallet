# Taproot Testing Guide

This document describes the comprehensive test suite for Taproot (BIP341/BIP86) support in the go-crypto-wallet project.

## Overview

Taproot support has been implemented across all three wallet types (Watch, Keygen, Sign) with comprehensive test coverage at multiple layers of the architecture.

## Test Organization

Tests are organized following Clean Architecture principles:

```
tests/
├── Domain Layer Tests       - Pure business logic validation
├── Infrastructure Tests     - Key generation, Bitcoin Core RPC integration
└── Application Layer Tests  - Use case constructor tests
```

## Running Tests

### Run All Unit Tests

```bash
make gotest
```

### Run Specific Taproot Tests

```bash
# Domain layer validator tests
go test -v ./internal/domain/key/ -run "Taproot"

# BIP86 key generation tests
go test -v ./internal/infrastructure/wallet/key/ -run "BIP86"

# Integration tests (requires Bitcoin Core RPC)
go test -v -tags=integration ./internal/infrastructure/api/bitcoin/btc/
```

### Test Coverage

```bash
# Generate coverage report
go test -cover ./...

# Generate detailed HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Categories

### 1. Domain Layer Tests

**Location:** `internal/domain/key/validator_test.go`

**Purpose:** Validate Taproot address validation logic without external dependencies.

**Key Tests:**
- `TestValidateWalletKey_TaprootSupport` - Verifies Taproot addresses are properly validated
- `TestValidateWalletKey_BackwardCompatibility` - Ensures legacy address types still work

**Coverage:**
- Taproot-only wallet keys
- Mixed address format wallet keys (legacy + Taproot)
- Mainnet and testnet Taproot addresses
- Address format validation (bc1p/tb1p prefix, length, character set)

**Run:**
```bash
go test -v ./internal/domain/key/ -run "TestValidateWalletKey"
```

### 2. Infrastructure Layer Tests - Key Generation

**Location:** `internal/infrastructure/wallet/key/`

**Test Files:**
- `bip86_generator_test.go` - BIP86 key generator unit tests
- `bip86_integration_test.go` - BIP86 integration tests

**Key Tests:**

#### BIP86 Generator Tests
- `TestBIP86Generator` - Basic BIP86 key generation
- `TestBIP86GeneratorConsistency` - Deterministic key generation
- `TestBIP86vsHDKeyEquivalence` - Compatibility with existing HDKey implementation

#### BIP86 Integration Tests
- `TestBIP86IntegrationRealWalletScenario` - Real-world wallet operation scenarios
- `TestBIP86IntegrationKeyConsistency` - Wallet recovery consistency
- `TestBIP86IntegrationMultipleAccounts` - Multi-account wallet support
- `TestBIP86IntegrationAddressValidation` - Network-specific address validation

**Coverage:**
- BIP86 derivation paths (`m/86'/0'/account'/0/index`)
- Taproot address generation for all account types (Client, Deposit, Payment, Stored)
- Mainnet, Testnet3, and Signet network support
- Key uniqueness across accounts
- Wallet recovery consistency
- Schnorr public key derivation

**Run:**
```bash
go test -v ./internal/infrastructure/wallet/key/ -run "BIP86"
```

### 3. Application Layer Tests - Use Cases

**Location:** `internal/application/usecase/*/`

**Test Files:**
- `internal/application/usecase/watch/shared/import_address_test.go`
- `internal/application/usecase/keygen/shared/generate_seed_test.go`
- `internal/application/usecase/sign/shared/store_seed_test.go`

**Key Tests:**
- Constructor tests for use cases
- Interface compliance verification

**Note:** Full integration tests for use cases require:
- Database setup (MySQL)
- Test fixtures for address files
- Bitcoin Core RPC connection

**Run:**
```bash
go test -v ./internal/application/usecase/...
```

## Test Scenarios Covered

### Taproot Address Generation
- ✅ Generate Taproot addresses from BIP86 derivation paths
- ✅ Support all network types (Mainnet, Testnet, Signet)
- ✅ Generate unique addresses across multiple accounts
- ✅ Deterministic key generation for wallet recovery

### Taproot Address Validation
- ✅ Validate Taproot address format (bc1p/tb1p prefix)
- ✅ Validate address length (62 characters)
- ✅ Validate bech32m character set
- ✅ Support Taproot-only and mixed address format wallet keys

### Taproot Transaction Support
- ✅ Transaction creation with Taproot outputs (via Bitcoin Core RPC)
- ✅ Schnorr signature generation (via Bitcoin Core RPC)
- ✅ Multisig Taproot transaction signing

### Backward Compatibility
- ✅ Legacy address types (P2PKH, P2SH-SegWit, Bech32) continue to work
- ✅ Mixed wallet keys with multiple address formats

## Integration Testing

### Prerequisites

Integration tests require:
1. **Bitcoin Core v22.0+** (for Taproot/Schnorr support)
2. **MySQL Database** (for address/transaction storage)
3. **Configuration Files** with Taproot settings

### Setup Bitcoin Core for Testing

```bash
# Start Bitcoin Core in regtest mode
bitcoind -regtest -daemon

# Create test wallet
bitcoin-cli -regtest createwallet "test_taproot"

# Generate blocks for testing
bitcoin-cli -regtest -rpcwallet=test_taproot -generate 101
```

### Configuration

Set `address_type = "taproot"` in your test configuration:

```toml
# data/config/btc_keygen_bip86_test.toml
key_type = "bip86"        # Required for Taproot
address_type = "taproot"  # legacy, p2sh-segwit, bech32, taproot
```

### Run Integration Tests

```bash
# Run with integration tag
go test -v -tags=integration ./internal/infrastructure/api/bitcoin/btc/
```

## End-to-End Testing Workflow

### 1. Keygen Wallet - Generate Taproot Keys

```bash
# Configure for Taproot
export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen_bip86_test.toml

# Generate seed
./keygen --coin btc create seed

# Generate HD wallet keys
./keygen --coin btc create hdkey --account client --count 10

# Export addresses
./keygen --coin btc export address --account client
```

### 2. Watch Wallet - Import Taproot Addresses

```bash
# Configure for Taproot
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml

# Import Taproot addresses
./watch --coin btc import address --account client \
  --filepath ./data/address/btc/bip86_test/client_1234567890.csv

# Create unsigned transaction
./watch --coin btc create transaction --account deposit
```

### 3. Keygen Wallet - Sign Taproot Transaction

```bash
# Sign transaction (Schnorr signature)
./keygen --coin btc sign --file ./data/tx/btc/deposit_1_unsigned_0_1234567890.tx
```

### 4. Sign Wallet - Multisig Taproot Signing

```bash
# Configure for Taproot
export BTC_SIGN_WALLET_CONF=./data/config/btc_sign.toml

# Sign multisig transaction (Schnorr signature)
./sign --coin btc sign --file ./data/tx/btc/payment_5_unsigned_1_1234567890.tx
```

### 5. Watch Wallet - Send Taproot Transaction

```bash
# Send signed transaction
./watch --coin btc send transaction --account deposit \
  --hex 020000000001...

# Monitor transaction
./watch --coin btc monitor transaction
```

## Test Coverage Summary

| Layer | Component | Test Type | Coverage |
|-------|-----------|-----------|----------|
| Domain | Key Validators | Unit | ✅ Comprehensive |
| Infrastructure | BIP86 Generator | Unit | ✅ Comprehensive |
| Infrastructure | BIP86 Integration | Integration | ✅ Comprehensive |
| Application | Use Cases | Constructor | ✅ Basic |
| End-to-End | Full Workflow | Manual | ⚠️  Requires setup |

## Known Limitations

### Integration Tests
- **Database Required:** Full use case tests require MySQL database setup
- **Bitcoin Core Required:** Transaction signing tests require Bitcoin Core v22.0+
- **Network Access:** Some tests require testnet/mainnet connectivity

### Manual Testing
- Testnet/Mainnet transaction broadcast requires real network connectivity
- Multisig scenarios require multiple wallet instances

## Troubleshooting

### Test Failures

**Issue:** `fail to call NewAddressTaproot()`
**Solution:** Ensure Bitcoin Core v22.0+ is running and configured correctly

**Issue:** `encodedPrevsAddrs must be set in csv file`
**Solution:** Transaction files must include previous transaction data (required since Bitcoin Core v17)

**Issue:** `key index exceeds maximum`
**Solution:** Verify key index is within valid range (0 to 2^31-1)

### Bitcoin Core Version

Verify Bitcoin Core version supports Taproot:
```bash
bitcoin-cli --version
# Should be v22.0 or higher
```

### Address Validation

Verify Taproot address format:
```bash
# Mainnet: bc1p... (62 characters)
# Testnet/Signet: tb1p... (62 characters)
```

## Future Test Enhancements

### Phase 1h (Documentation)
- User guide for Taproot usage
- API reference documentation
- Migration guide from legacy addresses

### Potential Additions
- **Performance Tests:** Benchmark Schnorr signature generation
- **Stress Tests:** Large transaction batch signing
- **Fuzz Tests:** Input validation edge cases
- **Network Tests:** Testnet/Mainnet transaction broadcast
- **Multisig Tests:** Complex multisig scenarios with Taproot

## References

- [BIP341 - Taproot: SegWit version 1 spending rules](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP86 - Key Derivation for Single Key P2TR Outputs](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)
- [BIP340 - Schnorr Signatures for secp256k1](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [Bitcoin Core v22.0 Release Notes](https://github.com/bitcoin/bitcoin/blob/master/doc/release-notes/release-notes-22.0.md)

## Contributing

When adding new Taproot functionality:
1. Add domain layer tests for business logic validation
2. Add infrastructure tests for implementation details
3. Update integration tests for end-to-end scenarios
4. Document test scenarios in this file
5. Ensure backward compatibility with existing address types

---

**Last Updated:** Phase 1g - Comprehensive Taproot testing implementation
**Related Issue:** #89 - Implement Taproot Support (Phase 1)
