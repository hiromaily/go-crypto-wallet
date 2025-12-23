package wallet

import (
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet instead.
// This package provides backward compatibility aliases.

// WalletType wallet type
// Deprecated: Use domain/wallet.WalletType
type WalletType = domainWallet.WalletType

// wallet_type
// Deprecated: Use constants from domain/wallet package
const (
	WalletTypeWatchOnly = domainWallet.WalletTypeWatchOnly
	WalletTypeKeyGen    = domainWallet.WalletTypeKeyGen
	WalletTypeSign      = domainWallet.WalletTypeSign
)

// WalletTypeValue value
// Deprecated: Use domain/wallet.WalletTypeValue
var WalletTypeValue = domainWallet.WalletTypeValue
