package btc

import (
	"fmt"

	"github.com/bookerzzz/grok"
)

func runLogging(btc btcWatchAPICmds) error {
	// logging
	logData, err := btc.Logging()
	if err != nil {
		return fmt.Errorf("fail to call BTC.Logging() %w", err)
	}
	grok.Value(logData)

	return nil
}
