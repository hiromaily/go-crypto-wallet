package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/bitcoin"
)

func runLogging(btc portsBitcoin.Bitcoiner) error {
	// logging
	logData, err := btc.Logging()
	if err != nil {
		return fmt.Errorf("fail to call BTC.Logging() %w", err)
	}
	grok.Value(logData)

	return nil
}
