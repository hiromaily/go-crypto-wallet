package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetBalance gets balance
// - It would include dirty outputs already spent tx, so it maybe useless
//   - wallet does not have the "avoid reuse" feature enabled
//   - `bitcoin-cli getbalance "*" 6 true true`
func (b *Bitcoin) GetBalance() (btcutil.Amount, error) {
	amount, err := btcrpc.GetBalance(b.Client, int(b.confirmationBlock))
	if err != nil {
		return 0, fmt.Errorf("fail to call btcrpc.GetBalance(): %w", err)
	}

	return b.FloatToAmount(amount)
}

// GetBalanceByListUnspent gets balance by rpc `listunspent`
func (b *Bitcoin) GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error) {
	listunspentResult, err := b.ListUnspent(confirmationNum)
	if err != nil {
		return 0, err
	}

	// Sum up all amounts (already in btcutil.Amount)
	var sum btcutil.Amount
	for _, unspent := range listunspentResult {
		sum += unspent.Amount
	}
	return sum, nil
}

// GetBalanceByAccount gets balance by account
func (b *Bitcoin) GetBalanceByAccount(
	accountType domainAccount.AccountType, confirmationNum uint64,
) (btcutil.Amount, error) {
	unspentList, err := b.ListUnspentByAccount(accountType, confirmationNum)
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.ListUnspentByAccount(%s): %w", accountType.String(), err)
	}

	// Sum up all amounts (already in btcutil.Amount)
	var totalAmount btcutil.Amount
	for _, tx := range unspentList {
		totalAmount += tx.Amount
	}
	return totalAmount, nil
}
