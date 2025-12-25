package wallets

import (
	wallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// Signer is a backward compatibility alias
//
// Deprecated: Use github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet.Signer instead
type Signer = wallet.Signer
