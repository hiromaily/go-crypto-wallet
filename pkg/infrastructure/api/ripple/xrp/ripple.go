package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/network/websocket"
)

// Ripple includes client to call JSON-RPC
type Ripple struct {
	wsPublic     *websocket.WS
	wsAdmin      *websocket.WS
	API          *RippleAPI
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode // eth
}

// NewRipple creates Ripple object
func NewRipple(
	ctx context.Context,
	wsPublic *websocket.WS,
	wsAdmin *websocket.WS,
	api *RippleAPI,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *config.Ripple,
) (*Ripple, error) {
	xrp := &Ripple{
		wsPublic:     wsPublic,
		wsAdmin:      wsAdmin,
		API:          api,
		coinTypeCode: coinTypeCode,
	}

	if conf.NetworkType != NetworkTypeXRPMainNet.String() {
		xrp.chainConf = &chaincfg.TestNet3Params
	} else {
		xrp.chainConf = &chaincfg.MainNetParams
	}

	return xrp, nil
}

// Close disconnect to server
func (r *Ripple) Close() error {
	if r.wsPublic != nil {
		_ = r.wsPublic.Close() // Best effort cleanup
	}
	if r.wsAdmin != nil {
		_ = r.wsAdmin.Close() // Best effort cleanup
	}
	if r.API != nil {
		r.API.Close()
	}
	return nil
}

// CoinTypeCode returns coinTypeCode
func (r *Ripple) CoinTypeCode() domainCoin.CoinTypeCode {
	return r.coinTypeCode
}

// GetChainConf returns chain conf
func (r *Ripple) GetChainConf() *chaincfg.Params {
	return r.chainConf
}
