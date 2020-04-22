package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// GetBalance gets balance
// - It would include dirty outputs already spent tx, so it maybe useless
//  - wallet does not have the "avoid reuse" feature enabled
//  - `bitcoin-cli getbalance "*" 6 true true`
func (b *Bitcoin) GetBalance() (btcutil.Amount, error) {
	input1, err := json.Marshal("*")
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call json.Marchal(dummy)")
	}
	input2, err := json.Marshal(uint64(b.confirmationBlock))
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call json.Marchal(%d)", b.confirmationBlock)
	}

	rawResult, err := b.client.RawRequest("getbalance", []json.RawMessage{input1, input2})
	if err != nil {
		return 0, errors.Wrap(err, "fail to call json.RawRequest(getbalance)")
	}

	var amount float64
	err = json.Unmarshal([]byte(rawResult), &amount)
	if err != nil {
		return 0, errors.Wrap(err, "fail to json.Unmarshal(rawResult)")
	}

	return b.FloatToAmount(amount)
}

func (b *Bitcoin) GetBalanceByListUnspent() (btcutil.Amount, error) {
	listunspentResult, err := b.ListUnspent()
	if err != nil {
		return 0, err
	}
	sum := b.getUnspentListAmount(listunspentResult)
	return b.FloatToAmount(sum)
}

// GetBalanceByAccount gets balance by account
func (b *Bitcoin) GetBalanceByAccount(accountType account.AccountType) (btcutil.Amount, error) {
	unspentList, err := b.ListUnspentByAccount(accountType)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btc.ListUnspentByAccount(%s)", accountType.String())
	}
	var totalAmout float64
	for _, tx := range unspentList {
		totalAmout += tx.Amount
	}
	return b.FloatToAmount(totalAmout)
}
