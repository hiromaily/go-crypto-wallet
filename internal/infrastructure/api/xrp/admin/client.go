package xrpadmin

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	xrprpcadmin "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// AdminXRP implements admin XRP node operations.
// It embeds *xrprpcadmin.AdminRPC to promote all low-level admin RPC methods,
// and overrides methods that require type conversion between DTO and RPC types.
type AdminXRP struct {
	*xrprpcadmin.AdminRPC
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAdminXRP creates a new AdminXRP with the given WebSocket connection, coin type,
// and pre-resolved chain configuration.
func NewAdminXRP(ws *websocket.WS, coinTypeCode domainCoin.CoinTypeCode, chainConf *chaincfg.Params) *AdminXRP {
	return &AdminXRP{
		AdminRPC:     xrprpcadmin.NewAdminRPC(ws),
		coinTypeCode: coinTypeCode,
		chainConf:    chainConf,
	}
}

// CoinTypeCode returns the coin type code.
func (a *AdminXRP) CoinTypeCode() domainCoin.CoinTypeCode {
	return a.coinTypeCode
}

// GetChainConf returns the chain configuration.
func (a *AdminXRP) GetChainConf() *chaincfg.Params {
	return a.chainConf
}
