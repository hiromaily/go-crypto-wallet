package eth

import (
	"math/big"
)

// UserAmount user address and amount
type UserAmount struct {
	Address string
	Amount  uint64
}

// GetTotalBalance returns total amount and addresses
func (e *Ethereum) GetTotalBalance(addrs []string) (*big.Int, []UserAmount) {
	total := new(big.Int)
	userAmounts := make([]UserAmount, 0, len(addrs))
	for _, addr := range addrs {
		balance, err := e.GetBalance(addr, QuantityTagPending)
		if err != nil {
			continue
		}
		if balance.Uint64() != 0 {
			total = total.Add(total, balance)
			userAmounts = append(userAmounts, UserAmount{Address: addr, Amount: balance.Uint64()})
		}
	}
	return total, userAmounts
}
