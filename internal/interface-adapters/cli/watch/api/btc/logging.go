package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runLogging(btc apibtc.Bitcoiner) error {
	// logging
	logData, err := btc.Logging()
	if err != nil {
		return fmt.Errorf("fail to call BTC.Logging() %w", err)
	}
	grok.Value(logData)

	return nil
}
