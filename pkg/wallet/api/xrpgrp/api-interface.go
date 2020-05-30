package xrpgrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Rippler Ripple Interface
type Rippler interface {
	// admin_keygen
	ValidationCreate(secret string) (*xrp.ResponseValidationCreate, error)
	WalletProposeWithKey(seed string, keyType xrp.XRPKeyType) (*xrp.ResponseWalletPropose, error)
	WalletPropose(passphrase string) (*xrp.ResponseWalletPropose, error)

	// public_account
	AccountChannels(sender, receiver string) (*xrp.ResponseAccountChannels, error)
	// public_server_info
	ServerInfo() (*xrp.ResponseServerInfo, error)
	// ripple
	Close()
	CoinTypeCode() coin.CoinTypeCode
	GetChainConf() *chaincfg.Params

	// RippleAPI
	PrepareTransaction(senderAccount, receiverAccount string, amount float64) error
}
