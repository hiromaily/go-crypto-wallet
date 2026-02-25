package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

func runClientVersion(eth apieth.ETHNodeAPIClient) error {
	version, err := eth.ClientVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.ClientVersion() %w", err)
	}

	fmt.Println("client version: " + version)

	return nil
}
