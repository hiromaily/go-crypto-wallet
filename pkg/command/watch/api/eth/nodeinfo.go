package eth

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

func runNodeInfo(eth ethgrp.Ethereumer) error {
	peerInfo, err := eth.NodeInfo(context.Background())
	if err != nil {
		return fmt.Errorf("fail to call eth.NodeInfo() %w", err)
	}

	fmt.Printf("nodeinfo: %v\n", peerInfo)

	return nil
}
