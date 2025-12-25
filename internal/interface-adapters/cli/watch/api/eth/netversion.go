package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
)

func runNetVersion(eth ethereum.Ethereumer) error {
	version, err := eth.NetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.NetVersion() %w", err)
	}

	fmt.Printf("net version: %d\n", version)

	return nil
}
