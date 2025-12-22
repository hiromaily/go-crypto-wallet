package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
)

// GetBalance gets balance
// - It would include dirty outputs already spent tx, so it maybe useless
//   - wallet does not have the "avoid reuse" feature enabled
//   - `bitcoin-cli getbalance "*" 6 true true`
func (b *Bitcoin) GetBalance() (btcutil.Amount, error) {
	input1, err := json.Marshal("*")
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marchal(dummy): %w", err)
	}
	input2, err := json.Marshal(b.confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marchal(%d): %w", b.confirmationBlock, err)
	}

	rawResult, err := b.Client.RawRequest("getbalance", []json.RawMessage{input1, input2})
	if err != nil {
		return 0, fmt.Errorf("fail to call json.RawRequest(getbalance): %w", err)
	}

	var amount float64
	err = json.Unmarshal(rawResult, &amount)
	if err != nil {
		return 0, fmt.Errorf("fail to json.Unmarshal(rawResult): %w", err)
	}

	return b.FloatToAmount(amount)
}

// GetBalanceByListUnspent gets balance by rpc `listunspent`
func (b *Bitcoin) GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error) {
	listunspentResult, err := b.ListUnspent(confirmationNum)
	if err != nil {
		return 0, err
	}
	sum := b.getUnspentListAmount(listunspentResult)
	return b.FloatToAmount(sum)
}

// GetBalanceByAccount gets balance by account
func (b *Bitcoin) GetBalanceByAccount(accountType account.AccountType, confirmationNum uint64) (btcutil.Amount, error) {
	unspentList, err := b.ListUnspentByAccount(accountType, confirmationNum)
	if err != nil {
		return 0, fmt.Errorf("fail to call btc.ListUnspentByAccount(%s): %w", accountType.String(), err)
	}
	var totalAmout float64
	for _, tx := range unspentList {
		totalAmout += tx.Amount
	}
	return b.FloatToAmount(totalAmout)
}
