package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

func runSyncing(eth apieth.ETHNodeAPIClient) error {
	syncResult, isSyncing, err := eth.Syncing(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.Syncing() %w", err)
	}

	fmt.Printf("is syncing? : %t, %v\n", isSyncing, syncResult)

	return nil
}
