package btc

import (
	"fmt"
)

func runEstimateFee(btc btcWatchAPICmds) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee()
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
