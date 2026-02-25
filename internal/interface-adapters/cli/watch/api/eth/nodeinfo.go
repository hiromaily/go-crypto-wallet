package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

func runNodeInfo(eth apieth.ETHNodeAPIClient) error {
	peerInfo, err := eth.NodeInfo(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.NodeInfo() %w", err)
	}

	fmt.Printf("nodeinfo: %v\n", peerInfo)

	return nil
}
