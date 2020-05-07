package ethgrp

import (
	"github.com/ethereum/go-ethereum/p2p"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// Ethereumer Ethereum Interface
type Ethereumer interface {
	// admin
	NodeInfo() (*p2p.PeerInfo, error)
	// eth
	Syncing() (*eth.ResponseSyncing, bool, error)
	// ethereum
	Close()
	// net
	NetVersion() (uint16, error)
	// web3
	ClientVersion() (string, error)
}
