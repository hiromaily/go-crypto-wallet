package bch

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// GetBalanceByAccount overrides Bitcoin's GetBalanceByAccount to ensure BCH-specific
// ListUnspentByAccount is called, which uses BCH address decoding
func (b *BitcoinCash) GetBalanceByAccount(
	accountType domainAccount.AccountType, confirmationNum uint64,
) (btcutil.Amount, error) {
	// Call BCH's ListUnspentByAccount which handles CashAddr format
	unspentList, err := b.ListUnspentByAccount(accountType, confirmationNum)
	if err != nil {
		return 0, fmt.Errorf("fail to call bch.ListUnspentByAccount(%s): %w", accountType.String(), err)
	}

	// Calculate total from unspent list
	var total btcutil.Amount
	for _, unspent := range unspentList {
		total += unspent.Amount
	}

	return total, nil
}
