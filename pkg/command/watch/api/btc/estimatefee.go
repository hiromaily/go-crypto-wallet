package btc

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
)

func runEstimateFee(btc bitcoin.Bitcoiner) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee()
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
