package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

func runNetVersion(eth apieth.Ethereumer) error {
	version, err := eth.NetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.NetVersion() %w", err)
	}

	fmt.Printf("net version: %d\n", version)

	return nil
}
