package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// GetAccount returns account name of address
// `getaccount` should be called because getaccount RPC is gone from version 0.18
func (b *Bitcoin) GetAccount(addr string) (string, error) {
	// actually `getaddressinfo` is called
	res, err := b.GetAddressInfo(addr)
	if err != nil {
		return "", errors.Wrap(err, "fail to call btc.GetAddressInfo()")
	}
	if len(res.Labels) == 0 {
		return "", nil
	}
	return res.Labels[0], nil
}

// ListAccounts list of balance in accounts
func (b *Bitcoin) ListAccounts() (map[string]btcutil.Amount, error) {
	listAmts, err := b.client.ListAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.ListAccounts()")
	}

	return listAmts, nil
}

// GetBalanceByAccount get balance for account
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
