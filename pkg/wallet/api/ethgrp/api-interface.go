package ethgrp

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// Ethereumer Ethereum Interface
type Ethereumer interface {
	// admin
	NodeInfo() (*p2p.NodeInfo, error)
	// client
	BalanceAt(hexAddr string) (*big.Int, error)
	SendRawTx(tx *types.Transaction) error
	// eth
	Syncing() (*eth.ResponseSyncing, bool, error)
	// ethereum
	Close()
	GetChainConf() *chaincfg.Params
	// key
	ToECDSA(privKey string) (*ecdsa.PrivateKey, error)
	GetKeyDir(accountType account.AccountType) string
	// net
	NetVersion() (uint16, error)
	// personal
	ImportRawKey(hexKey, passPhrase string) (string, error)
	// web3
	ClientVersion() (string, error)
}
