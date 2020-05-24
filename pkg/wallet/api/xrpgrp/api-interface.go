package xrpgrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Rippler Ripple Interface
type Rippler interface {
	// method_account
	AccountChannels(sender, receiver string) (*xrp.ResponseAccountChannels, error)

	// ripple
	Close()
	CoinTypeCode() coin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}
