package btc

import (
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runEstimateFee(btc apibtc.WatchAPIClient) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee(int(btc.ConfirmationBlock()))
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
