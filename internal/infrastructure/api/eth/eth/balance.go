package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
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

// BalanceAt returns balance of address
func (e *Ethereum) BalanceAt(ctx context.Context, hexAddr string) (*big.Int, error) {
	account := common.HexToAddress(hexAddr)
	balance, err := e.ethClient.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to call ethClient.BalanceAt(): %w", err)
	}
	return balance, nil
}
