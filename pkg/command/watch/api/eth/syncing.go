package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

func runSyncing(eth ethgrp.Ethereumer) error {
	syncResult, isSyncing, err := eth.Syncing(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.Syncing() %w", err)
	}

	fmt.Printf("is syncing? : %t, %v\n", isSyncing, syncResult)

	return nil
}
