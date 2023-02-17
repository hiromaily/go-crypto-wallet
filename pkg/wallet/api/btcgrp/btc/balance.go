package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
)

// GetBalance gets balance
// - It would include dirty outputs already spent tx, so it maybe useless
//   - wallet does not have the "avoid reuse" feature enabled
//   - `bitcoin-cli getbalance "*" 6 true true`
func (b *Bitcoin) GetBalance() (btcutil.Amount, error) {
	input1, err := json.Marshal("*")
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call json.Marchal(dummy)")
	}
	input2, err := json.Marshal(b.confirmationBlock)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call json.Marchal(%d)", b.confirmationBlock)
	}

	rawResult, err := b.Client.RawRequest("getbalance", []json.RawMessage{input1, input2})
	if err != nil {
		return 0, errors.Wrap(err, "fail to call json.RawRequest(getbalance)")
	}

	var amount float64
	err = json.Unmarshal(rawResult, &amount)
	if err != nil {
		return 0, errors.Wrap(err, "fail to json.Unmarshal(rawResult)")
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
		return 0, errors.Wrapf(err, "fail to call btc.ListUnspentByAccount(%s)", accountType.String())
	}
	var totalAmout float64
	for _, tx := range unspentList {
		totalAmout += tx.Amount
	}
	return b.FloatToAmount(totalAmout)
}
