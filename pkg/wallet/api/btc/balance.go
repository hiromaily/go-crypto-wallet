package btc

import (
	"encoding/json"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/pkg/errors"
)

// GetBalance gets balance
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

// GetBalanceByAccount gets balance by account
func (b *Bitcoin) GetBalanceByAccount(accountType account.AccountType) (btcutil.Amount, error) {
	unspentList, _, err := b.ListUnspentByAccount(accountType)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btc.ListUnspentByAccount(%s)", accountType.String())
	}
	var totalAmout float64
	for _, tx := range unspentList {
		totalAmout += tx.Amount
	}
	return b.FloatToAmount(totalAmout)
}

