package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

// GetTotalBalance returns total amount and addresses
func (e *Ethereum) GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainETH.UserAmount) {
	total := new(big.Int)
	userAmounts := make([]domainETH.UserAmount, 0, len(addrs))
	for _, addr := range addrs {
		balance, err := e.GetBalance(ctx, addr, domainETH.QuantityTagPending)
		if err != nil {
			continue
		}
		if balance.Uint64() != 0 {
			total = total.Add(total, balance)
			userAmounts = append(userAmounts, domainETH.UserAmount{Address: addr, Amount: balance.Uint64()})
		}
	}
	return total, userAmounts
}

// FloatToBigInt converts Ether(float64) to Wei(*big.Int).
// Kept as a method to satisfy the ERC20er interface, which ERC20 implements
// using decimal-aware conversion. For Ethereum, this always assumes 18 decimals.
func (*Ethereum) FloatToBigInt(v float64) *big.Int {
	return pkgeth.FromFloatEther(v)
}

// BalanceAt returns balance of address
func (e *Ethereum) BalanceAt(ctx context.Context, hexAddr string) (*big.Int, error) {
	account := common.HexToAddress(hexAddr)
	balance, err := e.ethClient.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to call ethClient.BalanceAt(): %w", err)
	}
	return balance, nil
}
