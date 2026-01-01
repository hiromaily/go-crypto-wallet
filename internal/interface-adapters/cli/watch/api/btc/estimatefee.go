package btc

import (
	"fmt"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
)

func runEstimateFee(btc portsBtc.Bitcoiner) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee()
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
