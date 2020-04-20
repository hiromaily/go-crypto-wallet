package btc

import (
	"github.com/pkg/errors"
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
//  => replaced by `listlabels`
//func (b *Bitcoin) ListAccounts() (map[string]btcutil.Amount, error) {
//	listAmts, err := b.client.ListAccounts()
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call client.ListAccounts()")
//	}
//
//	return listAmts, nil
//}
