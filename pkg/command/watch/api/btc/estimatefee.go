package btc

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runEstimateFee(btc btcgrp.Bitcoiner) error {
	// estimate fee
	feePerKb, err := btc.EstimateSmartFee()
	if err != nil {
		return fmt.Errorf("fail to call BTC.EstimateSmartFee() %w", err)
	}
	fmt.Printf("EstimateSmartFee: %f\n", feePerKb)

	return nil
}
