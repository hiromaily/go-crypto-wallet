package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runGetNetworkInfo(btc apibtc.Bitcoiner) error {
	// call getnetworkinfo
	infoData, err := btc.GetNetworkInfo()
	if err != nil {
		return fmt.Errorf("fail to call BTC.GetNetworkInfo() %w", err)
	}
	grok.Value(infoData)

	return nil
}
