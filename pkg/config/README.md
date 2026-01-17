# Configuration Package

This package handles loading and validation of wallet configuration files.

## Supported Formats

- TOML (`.toml`)
- YAML (`.yaml`, `.yml`)

## Environment Variable Override

Configuration values can be overridden using environment variables. This is useful for:

- CI/CD pipelines
- Docker containers
- E2E testing with different configurations

### Priority Order

1. **Environment Variables** (highest priority)
2. **Config File**
3. **Default Values** (lowest priority)

### Naming Convention

Environment variables follow this pattern:

```
WALLET_{KEY}
WALLET_{SECTION}_{KEY}
```

- **Prefix**: `WALLET_`
- **Separator**: `_` (underscore)
- **Case**: UPPERCASE

### Examples

| Config File Key | Environment Variable |
|-----------------|----------------------|
| `address_type` | `WALLET_ADDRESS_TYPE` |
| `bitcoin.host` | `WALLET_BITCOIN_HOST` |
| `bitcoin.user` | `WALLET_BITCOIN_USER` |
| `mysql.host` | `WALLET_MYSQL_HOST` |
| `mysql.dbname` | `WALLET_MYSQL_DBNAME` |
| `logger.level` | `WALLET_LOGGER_LEVEL` |

### Usage Examples

#### Override address type for E2E testing

```bash
# Run E2E test with legacy address type
export WALLET_ADDRESS_TYPE=legacy
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh
```

#### Override database connection

```bash
export WALLET_MYSQL_HOST=localhost:3307
export WALLET_MYSQL_USER=testuser
export WALLET_MYSQL_PASS=testpass
keygen -c config/wallet/btc/keygen.yaml create seed
```

#### Override Bitcoin node connection

```bash
export WALLET_BITCOIN_HOST=127.0.0.1:18443
export WALLET_BITCOIN_USER=rpcuser
export WALLET_BITCOIN_PASS=rpcpassword
watch -c config/wallet/btc/watch.yaml balance
```

### Docker / CI Usage

```yaml
# docker-compose.yaml
services:
  keygen:
    environment:
      - WALLET_ADDRESS_TYPE=legacy
      - WALLET_BITCOIN_HOST=btc-node:18332
      - WALLET_MYSQL_HOST=db:3306
```

```yaml
# GitHub Actions
jobs:
  e2e-test:
    env:
      WALLET_ADDRESS_TYPE: legacy
      WALLET_BITCOIN_HOST: 127.0.0.1:18332
```

## Configuration Structure

See `wallet.go` for the complete `WalletRoot` structure with all available configuration options.

### Key Configuration Fields

| Field | Description | Environment Variable |
|-------|-------------|----------------------|
| `address_type` | Address type (legacy, p2sh-segwit, bech32, taproot) | `WALLET_ADDRESS_TYPE` |
| `bitcoin.host` | Bitcoin node RPC host | `WALLET_BITCOIN_HOST` |
| `bitcoin.network_type` | Network (mainnet, testnet3, regtest) | `WALLET_BITCOIN_NETWORK_TYPE` |
| `mysql.host` | MySQL host | `WALLET_MYSQL_HOST` |
| `mysql.dbname` | MySQL database name | `WALLET_MYSQL_DBNAME` |
| `logger.level` | Log level (debug, info, warn, error) | `WALLET_LOGGER_LEVEL` |

### Automatic Key Type Derivation

`key_type` is automatically derived from `address_type` to ensure consistency:

| address_type | Derived key_type |
|--------------|------------------|
| `legacy` | `bip44` |
| `p2sh-segwit` | `bip49` |
| `bech32` | `bip84` |
| `taproot` | `bip86` |

See `internal/domain/address/types.go` `AddrType.ToKeyType()` for implementation.
