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
	*WSClient    // WebSocket operations (public + admin)
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRP creates XRP object
func NewXRP(
	wsPublic *websocket.WS,
	wsAdmin *websocket.WS,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *config.Ripple,
) (*XRP, error) {
	xrp := &XRP{
		WSClient:     newWSClient(wsPublic, wsAdmin),
		coinTypeCode: coinTypeCode,
	}

	if conf.NetworkType != NetworkTypeXRPMainNet.String() {
		xrp.chainConf = &chaincfg.TestNet3Params
	} else {
		xrp.chainConf = &chaincfg.MainNetParams
	}

	return xrp, nil
}

// Close disconnects all WebSocket connections.
// This overrides the promoted Close from *WSClient.
func (r *XRP) Close() error {
	if r.WSClient != nil {
		_ = r.WSClient.Close()
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
