package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
)

func runClientVersion(eth ethereum.Ethereumer) error {
	version, err := eth.ClientVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.ClientVersion() %w", err)
	}

	fmt.Println("client version: " + version)

	return nil
}
