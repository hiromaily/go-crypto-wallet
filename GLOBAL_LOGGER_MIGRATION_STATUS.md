# Global Logger Migration Status

## Issue #45: Refactor - Migrate to Global Logger Pattern

### Overview
Migrating from dependency-injected logger fields to global logger pattern across the entire go-crypto-wallet codebase.

### Progress: ~60% Complete

## Completed Work

### 1. Container Layer (✅ DONE)
- **File**: `pkg/di/container.go`
- Removed `logger logger.Logger` field from `container` struct
- Removed `newLogger()` method
- Updated all 100+ `c.newLogger()` calls to inline logger initialization or removed parameters
- Updated `SetGlobal()` calls in `NewKeygener()`, `NewWalleter()`, `NewSigner()` to use inline logger creation

### 2. Wallet Constructors (✅ DONE - 5 files)
- `pkg/wallet/wallets/ethwallet/keygen.go` - Removed logger parameter and field
- `pkg/wallet/wallets/ethwallet/watch.go` - Removed logger parameter and field
- `pkg/wallet/wallets/xrpwallet/keygen.go` - Removed logger parameter and field
- `pkg/wallet/wallets/xrpwallet/watch.go` - Removed logger parameter and field
- `pkg/wallet/wallets/btcwallet/watch.go` - Removed logger parameter and field

### 3. BTC Watch Services (✅ DONE - 6 files)
- `pkg/wallet/service/btc/watchsrv/address_importer.go` - Updated constructor and all logger calls
- `pkg/wallet/service/btc/watchsrv/tx_creator.go` - Updated constructor and all logger calls
- `pkg/wallet/service/btc/watchsrv/tx_creator_deposit.go` - Added logger import, updated calls
- `pkg/wallet/service/btc/watchsrv/tx_creator_payment.go` - Added logger import, updated calls
- `pkg/wallet/service/btc/watchsrv/tx_sender.go` - Updated constructor and all logger calls
- `pkg/wallet/service/btc/watchsrv/tx_monitor.go` - Updated constructor and all logger calls

### 4. ETH Watch Services (✅ DONE - 6 files)
- `pkg/wallet/service/eth/watchsrv/tx_creator.go` - Updated constructor and all logger calls
- `pkg/wallet/service/eth/watchsrv/tx_creator_deposit.go` - Added logger import, updated calls
- `pkg/wallet/service/eth/watchsrv/tx_creator_payment.go` - Added logger import, updated calls
- `pkg/wallet/service/eth/watchsrv/tx_creator_transfer.go` - Added logger import, updated calls
- `pkg/wallet/service/eth/watchsrv/tx_sender.go` - Updated constructor and all logger calls
- `pkg/wallet/service/eth/watchsrv/tx_monitor.go` - Updated constructor and all logger calls

### 5. Common Watch Services (✅ DONE - 2 files)
- `pkg/wallet/service/watchsrv/address_importer.go` - Removed logger parameter and field
- `pkg/wallet/service/watchsrv/payment_request_creator.go` - Removed logger parameter and field

## Remaining Work (Based on Build Errors)

### Priority 1: XRP Watch Services (⚠️ TODO - 3 files)
- [ ] `pkg/wallet/service/xrp/watchsrv/tx_creator.go` - Remove logger parameter from `NewTxCreate`
- [ ] `pkg/wallet/service/xrp/watchsrv/tx_sender.go` - Remove logger parameter from `NewTxSend`
- [ ] `pkg/wallet/service/xrp/watchsrv/tx_monitor.go` - Remove logger parameter from `NewTxMonitor`

### Priority 2: API Constructors (⚠️ TODO - 5 constructor functions)
- [ ] `pkg/wallet/api/btcgrp/bitcoin.go` - Remove logger parameter from `NewBitcoin`, update all logger calls
- [ ] `pkg/wallet/api/ethgrp/ethereum.go` - Remove logger parameter from `NewEthereum`, update all logger calls
- [ ] `pkg/wallet/api/ethgrp/erc20/erc20.go` - Remove logger parameter from `NewERC20`, update all logger calls
- [ ] `pkg/wallet/api/xrpgrp/ripple.go` - Remove logger parameter from `NewRipple`, update all logger calls
- [ ] `pkg/wallet/api/xrpgrp/xrp/grpc_api.go` - Remove logger parameter from `NewRippleAPI`, update all logger calls

### Priority 3: Watch Repository Constructors (⚠️ TODO - 8 files)
- [ ] `pkg/repository/watchrepo/btc_tx_sqlc.go` - Remove logger from `NewBTCTxRepositorySqlc`
- [ ] `pkg/repository/watchrepo/btc_tx_input_sqlc.go` - Remove logger from `NewBTCTxInputRepositorySqlc`
- [ ] `pkg/repository/watchrepo/btc_tx_output_sqlc.go` - Remove logger from `NewBTCTxOutputRepositorySqlc`
- [ ] `pkg/repository/watchrepo/tx_sqlc.go` - Remove logger from `NewTxRepositorySqlc`
- [ ] `pkg/repository/watchrepo/eth_detail_tx_sqlc.go` - Remove logger from `NewEthDetailTxInputRepositorySqlc`
- [ ] `pkg/repository/watchrepo/xrp_detail_tx_sqlc.go` - Remove logger from `NewXrpDetailTxInputRepositorySqlc`
- [ ] `pkg/repository/watchrepo/payment_request_sqlc.go` - Remove logger from `NewPaymentRequestRepositorySqlc`
- [ ] `pkg/repository/watchrepo/address_sqlc.go` - Remove logger from `NewAddressRepositorySqlc`

### Priority 4: Cold Repository Constructors (⚠️ TODO - 5 files)
- [ ] `pkg/repository/coldrepo/seed_sqlc.go` - Remove logger from `NewSeedRepositorySqlc`
- [ ] `pkg/repository/coldrepo/account_key_sqlc.go` - Remove logger from `NewAccountKeyRepositorySqlc`
- [ ] `pkg/repository/coldrepo/xrp_account_key_sqlc.go` - Remove logger from `NewXRPAccountKeyRepositorySqlc`
- [ ] `pkg/repository/coldrepo/auth_fullpubkey_sqlc.go` - Remove logger from `NewAuthFullPubkeyRepositorySqlc`
- [ ] `pkg/repository/coldrepo/auth_account_key_sqlc.go` - Remove logger from `NewAuthAccountKeyRepositorySqlc`

### Priority 5: File Repository Constructors (⚠️ TODO - 2 files)
- [ ] `pkg/address/file.go` - Remove logger parameter from `NewFileRepository`
- [ ] `pkg/tx/file.go` - Remove logger parameter from `NewFileRepository`

### Priority 6: Key Package (⚠️ TODO - 1 file)
- [ ] `pkg/wallet/key/hdkey.go` - Remove logger parameter from `NewHDKey`

### Priority 7: Cold Services (⚠️ TODO - estimated 15+ files)
- [ ] All files in `pkg/wallet/service/coldsrv/` - Seed, HDWallet, AddressExport, Sign constructors
- [ ] All files in `pkg/wallet/service/btc/coldsrv/keygensrv/` - PrivKey, FullPubkeyImport, Multisig constructors
- [ ] All files in `pkg/wallet/service/btc/coldsrv/signsrv/` - PrivKey, FullPubkeyExport constructors
- [ ] All files in `pkg/wallet/service/eth/keygensrv/` - PrivKey, Sign constructors
- [ ] All files in `pkg/wallet/service/xrp/keygensrv/` - XRPKeyGenerate, Sign constructors

## Migration Pattern

For each file, the following changes are applied:

### 1. Constructor Function Signature
```go
// Before
func NewExample(logger logger.Logger, other params) *Example {
    return &Example{
        logger: logger,
        // ...
    }
}

// After
func NewExample(other params) *Example {
    return &Example{
        // logger field removed
        // ...
    }
}
```

### 2. Struct Definition
```go
// Before
type Example struct {
    logger logger.Logger
    // ...
}

// After
type Example struct {
    // logger field removed
    // ...
}
```

### 3. Logger Calls
```go
// Before
e.logger.Debug("message", args...)
e.logger.Info("message", args...)
e.logger.Warn("message", args...)
e.logger.Error("message", args...)

// After
logger.Debug("message", args...)
logger.Info("message", args...)
logger.Warn("message", args...)
logger.Error("message", args...)
```

### 4. Import Management
- Add `"github.com/hiromaily/go-crypto-wallet/pkg/logger"` if logger calls exist
- Remove unused logger import if no logger calls remain in the file

## Testing Strategy

After completing all migrations, run:

```bash
# Fix linting issues
make lint-fix

# Verify build
make check-build

# Run tests
make gotest
```

## Next Steps

1. Complete Priority 1-3 to get the build passing
2. Run tests to verify functionality
3. Complete Priority 4-7 for full migration
4. Run full test suite
5. Commit and create pull request

## Notes

- The global logger is initialized in three container methods: `NewKeygener()`, `NewWalleter()`, and `NewSigner()`
- Global logger functions are thread-safe and available via: `logger.Debug()`, `logger.Info()`, `logger.Warn()`, `logger.Error()`
- All logger calls use the same signature: `logger.Level(msg string, args ...any)`
- This migration removes 100+ logger fields and simplifies the dependency injection structure

## References

- Issue: https://github.com/hiromaily/go-crypto-wallet/issues/45
- Global Logger Implementation: `pkg/logger/global.go`
