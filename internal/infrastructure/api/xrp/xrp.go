package xrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// XRP includes client to call JSON-RPC
// This type implements the interfaces defined in internal/application/ports/api/xrp
type XRP struct {
	wsPublic     *websocket.WS
	wsAdmin      *websocket.WS
	API          *XRPAPI
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRP creates XRP object
func NewXRP(
	wsPublic *websocket.WS,
	wsAdmin *websocket.WS,
	api *XRPAPI,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *config.Ripple,
) (*XRP, error) {
	xrp := &XRP{
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
func (r *XRP) Close() error {
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
func (r *XRP) CoinTypeCode() domainCoin.CoinTypeCode {
	return r.coinTypeCode
}

// GetChainConf returns chain conf
func (r *XRP) GetChainConf() *chaincfg.Params {
	return r.chainConf
}
