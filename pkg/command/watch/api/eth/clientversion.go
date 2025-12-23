package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

func runClientVersion(eth ethgrp.Ethereumer) error {
	version, err := eth.ClientVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.ClientVersion() %w", err)
	}

	fmt.Println("client version: " + version)

	return nil
}
