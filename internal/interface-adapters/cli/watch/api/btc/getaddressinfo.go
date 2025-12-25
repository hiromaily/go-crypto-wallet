package btc

import (
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

func runGetAddressInfo(btc bitcoin.Bitcoiner, addr string) error {
	// validator
	if addr == "" {
		return errors.New("address option [-address] is required")
	}

	// call getaddressinfo
	addrData, err := btc.GetAddressInfo(addr)
	if err != nil {
		return fmt.Errorf("fail to call BTC.GetAddressInfo() %w", err)
	}
	grok.Value(addrData)

	return nil
}
