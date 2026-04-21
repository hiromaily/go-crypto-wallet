# BCH: Interface Separation Requirements

## Overview

This document defines the interface separation requirements between BTC and BCH.
The current `Bitcoiner` interface mixes BTC-only methods, which causes architectural issues.

## Current Problem

The `Bitcoiner` interface at `internal/application/ports/api/btc/interface.go` contains methods that are:

1. **Common to BTC and BCH** (should remain in shared interface)
2. **BTC-Only** (should be in a separate BTC-specific interface)

### Current Interface Structure (Problematic)

```go
// internal/application/ports/api/btc/interface.go
type Bitcoiner interface {
    // ✅ Common methods (BTC + BCH)
    GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)
    CreateRawTransaction(...) (*wire.MsgTx, error)
    SignRawTransactionWithKey(...) (*wire.MsgTx, bool, error)
    // ... other common methods

    // ❌ BTC-Only methods (should NOT be in shared interface)
    CreatePSBT(...) (string, error)           // PSBT - BTC only
    SignPSBTWithKey(...) (string, bool, error) // PSBT - BTC only
    ImportDescriptors(...) (...)              // Descriptor - BTC only
    GetDescriptorInfo(...) (...)              // Descriptor - BTC only
    // ... other BTC-only methods
}
```

## Proposed Interface Separation

### Option 1: Separate Interfaces (Recommended)

```go
// internal/application/ports/api/btc/interface.go

// BitcoinCompatible - Methods common to BTC and BCH
type BitcoinCompatible interface {
    // Address operations
    GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)
    ValidateAddress(addr string) (*dtobtc.ValidateAddressResult, error)
    DecodeAddress(addr string) (btcutil.Address, error)

    // Transaction operations (Raw TX - works for both)
    CreateRawTransaction(...) (*wire.MsgTx, error)
    SignRawTransaction(...) (*wire.MsgTx, bool, error)
    SignRawTransactionWithKey(...) (*wire.MsgTx, bool, error)
    SendTransactionByHex(hex string) (*chainhash.Hash, error)
    DecodeRawTransaction(hexTx string) (*dtobtc.RawTransaction, error)

    // Balance and UTXO
    GetBalance() (btcutil.Amount, error)
    ListUnspent(confirmationNum uint64) ([]dtobtc.UnspentOutput, error)

    // Import (traditional methods)
    ImportPrivKey(privKeyWIF *btcutil.WIF) error
    ImportAddress(pubkey string) error
    ImportAddressWithLabel(addr, label string, rescan bool) error

    // Multisig (P2SH - works for both)
    AddMultisigAddress(...) (*dtobtc.MultisigAddress, error)

    // ... other common methods
}

// BTCOnly - Methods specific to BTC (PSBT, Descriptor, MuSig2)
type BTCOnly interface {
    // PSBT operations (BIP174)
    CreatePSBT(...) (string, error)
    ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error)
    ValidatePSBT(psbtBase64 string) error
    SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
    WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error)
    FinalizePSBT(psbtBase64 string) (string, error)
    ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
    IsPSBTComplete(psbtBase64 string) (bool, error)
    GetPSBTFee(psbtBase64 string) (int64, error)

    // Descriptor operations
    ImportDescriptors(requests []dtobtc.ImportDescriptorsRequest) ([]dtobtc.ImportDescriptorsResponse, error)
    GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error)
    ListDescriptors(privateDescriptors bool) (*dtobtc.ListDescriptorsResult, error)
}

// Bitcoiner - Full BTC interface (embeds both)
type Bitcoiner interface {
    BitcoinCompatible
    BTCOnly
}

// BCHer - BCH-specific interface (only compatible methods)
type BCHer interface {
    BitcoinCompatible
    // BCH-specific methods can be added here
}
```

### Option 2: Runtime Checks (Current Workaround)

Until interface separation is implemented, add runtime checks in DI layer:

```go
// internal/di/container.go

func (c *container) NewWatchImportDescriptorUseCase() watchusecase.ImportDescriptorUseCase {
    // Block BCH from using BTC-only features
    if c.conf.CoinTypeCode == domainCoin.BCH {
        panic("Descriptor import is not supported for BCH. " +
              "BCH should use ImportAddress methods instead.")
    }
    // ... rest of implementation
}
```

## Method Classification

### Common Methods (BTC + BCH)

| Category | Methods |
|----------|---------|
| Address | `GetAddressInfo`, `ValidateAddress`, `DecodeAddress`, `GetAddressesByLabel` |
| Amount | `AmountString`, `FloatToAmount`, `StrToAmount`, etc. |
| Balance | `GetBalance`, `GetBalanceByListUnspent`, `GetBalanceByAccount` |
| Block | `GetBlockCount` |
| Config | `Close`, `GetChainConf`, `CoinTypeCode`, etc. |
| Fee | `EstimateSmartFee`, `GetTransactionFee`, `GetFee` |
| Import | `ImportPrivKey`, `ImportAddress`, `ImportAddressWithLabel` |
| Label | `SetLabel` |
| Multisig | `AddMultisigAddress` (P2SH only for BCH) |
| Network | `GetNetworkInfo`, `GetBlockchainInfo` |
| Transaction | `CreateRawTransaction`, `SignRawTransaction`, `SignRawTransactionWithKey`, `SendTransactionByHex`, `DecodeRawTransaction`, `ToHex`, `ToMsgTx` |
| UTXO | `ListUnspent`, `ListUnspentByAccount`, `LockUnspent`, `UnlockUnspent` |
| Wallet | `BackupWallet`, `LoadWallet`, etc. |

### BTC-Only Methods

| Category | Methods | BCH Alternative |
|----------|---------|-----------------|
| PSBT | `CreatePSBT`, `ParsePSBT`, `ValidatePSBT`, `SignPSBTWithKey`, `WalletProcessPsbt`, `FinalizePSBT`, `ExtractTransaction`, `IsPSBTComplete`, `GetPSBTFee` | Use Raw TX methods |
| Descriptor | `ImportDescriptors`, `GetDescriptorInfo`, `ListDescriptors` | Use `ImportAddress` |
| Import | `ImportMulti` (uses descriptors internally) | Use `ImportAddress` |

## Implementation Priority

1. **Immediate**: Add runtime checks in DI layer (Issue #435)
2. **Short-term**: Add `IsBTCOnly()` helper function
3. **Long-term**: Refactor to separate interfaces

## Helper Functions to Add

```go
// internal/domain/coin/types.go

// IsBTCOnly returns true only for BTC (excludes BCH and other coins)
func IsBTCOnly(val CoinTypeCode) bool {
    return val == BTC
}

// IsBTCCompatible returns true for coins that share BTC base functionality
// but may have different advanced features
func IsBTCCompatible(val CoinTypeCode) bool {
    return val == BTC || val == BCH
}

// Existing function - rename for clarity
// IsBTCGroup is DEPRECATED, use IsBTCCompatible instead
func IsBTCGroup(val CoinTypeCode) bool {
    return IsBTCCompatible(val)
}
```

## Related Issues

- #435 - BCH: Use case layer incorrectly uses BTC-only APIs
- #434 - BCH: Implement SIGHASH_FORKID support

## Related Documentation

- `.claude/rules/bch/btc-only-files.md` - List of BTC-only files
- `docs/chains/bch/README.md` - BCH technical context
