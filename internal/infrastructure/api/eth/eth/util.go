package eth

import (
	"math/big"

	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

// FloatToBigInt converts Ether(float64) to Wei(*big.Int).
// Kept as a method to satisfy the ERC20er interface, which ERC20 implements
// using decimal-aware conversion. For Ethereum, this always assumes 18 decimals.
func (*Ethereum) FloatToBigInt(v float64) *big.Int {
	return pkgeth.FromFloatEther(v)
}
