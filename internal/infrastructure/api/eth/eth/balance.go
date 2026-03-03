package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

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

// invalidAddressBalanceSentinel is the specific balance value returned by some ETH nodes
// when querying the balance of an invalidly formatted address.
const invalidAddressBalanceSentinel uint64 = 416778046407207737

// BalanceAt returns balance of address
// if wrong address is given, response of balance would be strange like `416778046407207737`
func (e *Ethereum) BalanceAt(ctx context.Context, hexAddr string) (*big.Int, error) {
	account := common.HexToAddress(hexAddr)
	balance, err := e.ethClient.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to call ethClient.BalanceAt(): %w", err)
	}
	if balance.Uint64() == invalidAddressBalanceSentinel {
		return nil, fmt.Errorf("received invalid address sentinel balance %d: possibly invalid address",
			invalidAddressBalanceSentinel)
	}
	return balance, nil
}
