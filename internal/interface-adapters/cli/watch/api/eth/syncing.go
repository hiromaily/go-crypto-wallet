package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
)

func runSyncing(eth ethereum.Ethereumer) error {
	syncResult, isSyncing, err := eth.Syncing(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.Syncing() %w", err)
	}

	fmt.Printf("is syncing? : %t, %v\n", isSyncing, syncResult)

	return nil
}
