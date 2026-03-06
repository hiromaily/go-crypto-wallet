package xrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplclient"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// XRP includes client to call JSON-RPC
// This type implements the interfaces defined in internal/application/ports/api/xrp
type XRP struct {
	API          *xrplclient.XRPLClient // gRPC operations (legacy, being phased out)
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRP creates XRP object
func NewXRP(
	api *xrplclient.XRPLClient,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *config.Ripple,
) (*XRP, error) {
	xrp := &XRP{
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

// Close disconnects all connections (WebSocket and gRPC).
// This overrides the promoted Close from *WSClient.
func (r *XRP) Close() error {
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
