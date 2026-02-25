package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"
)

func runGetNetworkInfo(btc btcWatchAPICmds) error {
	// call getnetworkinfo
	infoData, err := btc.GetNetworkInfo()
	if err != nil {
		return fmt.Errorf("fail to call BTC.GetNetworkInfo() %w", err)
	}
	grok.Value(infoData)

	return nil
}
