package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

func runNetVersion(eth ethgrp.Ethereumer) error {
	version, err := eth.NetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.NetVersion() %w", err)
	}

	fmt.Printf("net version: %d\n", version)

	return nil
}
