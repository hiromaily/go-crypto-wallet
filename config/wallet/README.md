# Wallet Configuration Files

This directory contains configuration files for each wallet type and cryptocurrency.

## Directory Structure

```
config/wallet/
├── btc/               # Bitcoin wallet configurations
│   ├── watch.yaml     # BTC Watch-only wallet (online)
│   ├── keygen.yaml    # BTC Keygen wallet (offline recommended)
│   ├── sign.yaml      # BTC Sign wallet (offline recommended)
│   ├── sign1.yaml     # Multisig signer 1
│   └── sign2.yaml     # Multisig signer 2
├── bch/               # Bitcoin Cash wallet configurations
│   ├── watch.yaml
│   ├── keygen.yaml
│   ├── sign.yaml
│   ├── sign1.yaml
│   └── sign2.yaml
├── eth/               # Ethereum wallet configurations
│   ├── watch.yaml
│   ├── keygen.yaml
│   └── sign.yaml
├── xrp/               # XRP (Ripple) wallet configurations
│   ├── watch.yaml
│   └── keygen.yaml
├── account/           # Account type configurations
│   ├── account.yaml        # Single-sig account
│   ├── account_2of3.yaml   # 2-of-3 multisig account
│   └── account_3of3.yaml   # 3-of-3 multisig account
└── archive/           # Legacy TOML configurations (deprecated)
    ├── btc/
    ├── bch/
    ├── eth/
    ├── xrp/
    └── account/
```

## Configuration File Types

### Wallet Config (`{chain}/{wallet_type}.yaml`)

| File | Purpose |
|------|---------|
| `btc/watch.yaml` | BTC Watch-only wallet (online) |
| `btc/keygen.yaml` | BTC Keygen wallet (offline recommended) |
| `btc/sign.yaml` | BTC Sign wallet (offline recommended) |
| `btc/sign1.yaml` / `btc/sign2.yaml` | Multisig signers 1/2 |

### Account Config (`account/account*.yaml`)

| File | Purpose | Use Case |
|------|---------|----------|
| `account/account.yaml` | Single-sig account configuration | Pattern 1, 3, 5, 9 |
| `account/account_2of3.yaml` | 2-of-3 multisig account configuration | Pattern 2, 4 |
| `account/account_3of3.yaml` | 3-of-3 multisig account configuration | Pattern 8 |

## ⚠️ Important: Configuration File Policy

**These configuration files are NOT intended to be edited or overwritten per script execution.**

### Override via Environment Variables

Use **environment variables** to run CLI with different settings:

```bash
# Override config without editing files
export WALLET_ADDRESS_TYPE=legacy
export WALLET_BITCOIN_HOST=127.0.0.1:18443
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh
```

#### Environment Variable Naming

```
WALLET_{KEY}              # Top-level settings
WALLET_{SECTION}_{KEY}    # Nested settings
```

#### Common Environment Variables

| Config Key | Environment Variable |
|------------|---------------------|
| `address_type` | `WALLET_ADDRESS_TYPE` |
| `bitcoin.host` | `WALLET_BITCOIN_HOST` |
| `bitcoin.user` | `WALLET_BITCOIN_USER` |
| `bitcoin.pass` | `WALLET_BITCOIN_PASS` |
| `mysql.host` | `WALLET_MYSQL_HOST` |
| `mysql.dbname` | `WALLET_MYSQL_DBNAME` |
| `logger.level` | `WALLET_LOGGER_LEVEL` |

#### Priority Order

1. **Environment Variables** (highest priority)
2. **Config File**
3. **Default Values** (lowest priority)

See [pkg/config/README.md](../../pkg/config/README.md) for details.

## Bitcoin RPC Host Configuration

The `bitcoin.host` setting supports two formats: **domain only** and **with path**.

### Format 1: Domain Only

```yaml
bitcoin:
  host: "127.0.0.1:20332"
```

### Format 2: With Path (`/wallet/<name>`)

```yaml
bitcoin:
  host: "127.0.0.1:18332/wallet/watch"
```

### Which Format to Use?

| Wallet Type | Recommended | Reason |
|-------------|-------------|--------|
| **Watch** | With path | Imports Descriptors for address/UTXO management in bitcoind |
| **Keygen** | With path | Imports Descriptors for address/UTXO management in bitcoind |
| **Sign** | Domain only | Does not use Descriptors. Manages private keys within the app |

### When Path Format is Required

When using bitcoind's [Descriptor Wallet](https://bitcoincore.org/en/doc/24.0.0/rpc/wallet/importdescriptors/) feature, you must create the wallet in bitcoind first:

```bash
# Create wallets in bitcoind
bitcoin-cli createwallet "watch"
bitcoin-cli createwallet "keygen"

# Execute RPC with specific wallet
bitcoin-cli -rpcwallet=watch getwalletinfo
```

The path format (`host: "127.0.0.1:18332/wallet/watch"`) is equivalent to `bitcoin-cli -rpcwallet=watch`.

#### Creation Script

Script to create wallets in Docker container:

```bash
./scripts/operation/btc/create-bitcoind-wallet.sh
```

### Why Sign Wallet Doesn't Use Path

Sign wallet does not use bitcoind's wallet feature because:

1. **Offline Environment**: Sign wallet is designed for offline operation for security
2. **Private Key Management**: Manages private keys within the application, not in bitcoind
3. **Transaction Signing**: Receives unsigned PSBTs and performs signing within the app

## Related Documentation

- [pkg/config/README.md](../../pkg/config/README.md) - Config package details
- [docs/crypto/btc/descriptor/user-guide.md](../../docs/crypto/btc/descriptor/user-guide.md) - Descriptor Wallet guide
- [scripts/operation/btc/](../../scripts/operation/btc/) - BTC operation scripts
