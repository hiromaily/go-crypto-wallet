package btc

import (
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/bitcoin"
)

func runGetAddressInfo(btc portsBitcoin.Bitcoiner, addr string) error {
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
