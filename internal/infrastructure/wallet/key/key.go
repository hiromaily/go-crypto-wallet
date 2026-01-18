// Package key provides implementations of wallet key generation.
//
// DEPRECATED: The interfaces (Generator, GeneratorFactory) have been moved to
// internal/application/ports/wallet following Clean Architecture principles.
// Import from there instead:
//
//	import portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
//
// This package now contains only the implementations.
package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key/strategy"
)

// CreateCoinKeyStrategy creates the appropriate CoinKeyStrategy based on coin type.
// This is a re-export of strategy.CreateCoinKeyStrategy for convenience.
func CreateCoinKeyStrategy(
	coinTypeCode domainCoin.CoinTypeCode,
	conf *chaincfg.Params,
) (portsWallet.CoinKeyStrategy, error) {
	return strategy.CreateCoinKeyStrategy(coinTypeCode, conf)
}
