package eth

import (
	"context"
	"math/big"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
)

// UserAmount user address and amount
type UserAmount struct {
	Address string
	Amount  uint64
}

// GetTotalBalance returns total amount and addresses
func (e *Ethereum) GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainETH.UserAmount) {
	total := new(big.Int)
	userAmounts := make([]UserAmount, 0, len(addrs))
	for _, addr := range addrs {
		balance, err := e.GetBalance(ctx, addr, domainETH.QuantityTagPending)
		if err != nil {
			continue
		}
		if balance.Uint64() != 0 {
			total = total.Add(total, balance)
			userAmounts = append(userAmounts, UserAmount{Address: addr, Amount: balance.Uint64()})
		}
	}
	// Convert to domain types
	return total, ToDomainUserAmounts(userAmounts)
}
