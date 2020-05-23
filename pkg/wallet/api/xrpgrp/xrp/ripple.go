package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Ripple includes client to call JSON-RPC
type Ripple struct {
	chainConf    *chaincfg.Params
	coinTypeCode coin.CoinTypeCode //eth
	logger       *zap.Logger
	ctx          context.Context
}

// NewRipple creates Ripple object
func NewRipple(
	ctx context.Context,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Ripple,
	logger *zap.Logger) (*Ripple, error) {

	xrp := &Ripple{
		coinTypeCode: coinTypeCode,
		logger:       logger,
		ctx:          ctx,
	}

	xrp.chainConf = &chaincfg.TestNet3Params

	return xrp, nil
}

// Close disconnect to server
func (r *Ripple) Close() {
	//if e.rpcClient != nil {
	//	e.rpcClient.Close()
	//}
}

// CoinTypeCode returns coinTypeCode
func (r *Ripple) CoinTypeCode() coin.CoinTypeCode {
	return r.coinTypeCode
}

// GetChainConf returns chain conf
func (r *Ripple) GetChainConf() *chaincfg.Params {
	return r.chainConf
}
