package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
)

func runGetNetworkInfo(btc portsBitcoin.Bitcoiner) error {
	// call getnetworkinfo
	infoData, err := btc.GetNetworkInfo()
	if err != nil {
		return fmt.Errorf("fail to call BTC.GetNetworkInfo() %w", err)
	}
	grok.Value(infoData)

	return nil
}
