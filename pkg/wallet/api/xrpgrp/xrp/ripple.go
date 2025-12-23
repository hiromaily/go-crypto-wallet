package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/ws"
)

// Ripple includes client to call JSON-RPC
type Ripple struct {
	wsPublic     *ws.WS
	wsAdmin      *ws.WS
	API          *RippleAPI
	chainConf    *chaincfg.Params
	coinTypeCode coin.CoinTypeCode // eth
}

// NewRipple creates Ripple object
func NewRipple(
	ctx context.Context,
	wsPublic *ws.WS,
	wsAdmin *ws.WS,
	api *RippleAPI,
	coinTypeCode coin.CoinTypeCode,
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
		r.wsPublic.Close()
	}
	if r.wsAdmin != nil {
		r.wsAdmin.Close()
	}
	if r.API != nil {
		r.API.Close()
	}
	return nil
}

// CoinTypeCode returns coinTypeCode
func (r *Ripple) CoinTypeCode() coin.CoinTypeCode {
	return r.coinTypeCode
}

// GetChainConf returns chain conf
func (r *Ripple) GetChainConf() *chaincfg.Params {
	return r.chainConf
}
