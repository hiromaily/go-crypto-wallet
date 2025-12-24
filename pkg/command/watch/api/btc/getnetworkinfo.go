package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
)

func runGetNetworkInfo(btc bitcoin.Bitcoiner) error {
	// call getnetworkinfo
	infoData, err := btc.GetNetworkInfo()
	if err != nil {
		return fmt.Errorf("fail to call BTC.GetNetworkInfo() %w", err)
	}
	grok.Value(infoData)

	return nil
}
