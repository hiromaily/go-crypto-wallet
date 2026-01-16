package btc

import (
	"fmt"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
)

func runEstimateFee(btc portsBitcoin.Bitcoiner) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee()
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
