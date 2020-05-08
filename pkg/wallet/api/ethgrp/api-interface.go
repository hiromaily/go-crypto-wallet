package ethgrp

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/p2p"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// Ethereumer Ethereum Interface
type Ethereumer interface {
	// admin
	NodeInfo() (*p2p.NodeInfo, error)
	// eth
	Syncing() (*eth.ResponseSyncing, bool, error)
	// ethereum
	Close()
	GetChainConf() *chaincfg.Params
	// net
	NetVersion() (uint16, error)
	// personal
	ImportRawKey(hexKey, passPhrase string) (string, error)
	// web3
	ClientVersion() (string, error)
}
