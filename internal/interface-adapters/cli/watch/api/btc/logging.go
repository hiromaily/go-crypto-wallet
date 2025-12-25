package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

func runLogging(btc bitcoin.Bitcoiner) error {
	// logging
	logData, err := btc.Logging()
	if err != nil {
		return fmt.Errorf("fail to call BTC.Logging() %w", err)
	}
	grok.Value(logData)

	return nil
}
