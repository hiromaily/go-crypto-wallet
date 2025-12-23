package btc

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runValidateAddress(btc btcgrp.Bitcoiner, address string) error {
	// validate args
	if address == "" {
		return errors.New("address option [-address] is required")
	}

	// validate address
	_, err := btc.ValidateAddress(address)
	if err != nil {
		return fmt.Errorf("fail to call BTC.ValidateAddress() %w", err)
	}

	return nil
}
