package ethgrp

import (
	"github.com/ethereum/go-ethereum/p2p"
)

// Ethereumer Ethereum Interface
type Ethereumer interface {
	// admin
	NodeInfo() (*p2p.PeerInfo, error)
	// ethereum
	Close()
	// net
	NetVersion() (uint16, error)
	// web3
	ClientVersion() (string, error)
}
