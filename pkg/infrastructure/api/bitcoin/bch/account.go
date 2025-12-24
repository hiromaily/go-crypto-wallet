package bch

import (
	"fmt"
)

// GetAccount returns account name of address
// `getaccount` should be called because getaccount RPC is gone from version 0.18
func (b *BitcoinCash) GetAccount(addr string) (string, error) {
	// actually `getaddressinfo` is called
	res, err := b.GetAddressInfo(addr)
	if err != nil {
		return "", fmt.Errorf("fail to call btc.GetAddressInfo()in bch: %w", err)
	}
	if len(res.Labels) == 0 {
		return "", nil
	}
	return res.Labels[0], nil
}
