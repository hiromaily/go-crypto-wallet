package xrppublic

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// PublicXRP implements public XRP node operations.
// It embeds *xrprpcpublic.PublicRPC to promote all low-level public RPC methods,
// and adds higher-level business logic methods.
type PublicXRP struct {
	*xrprpcpublic.PublicRPC
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
}

// NewPublicXRP creates a new PublicXRP with the given WebSocket connection, coin type,
// and pre-resolved chain configuration.
func NewPublicXRP(ws *websocket.WS, coinTypeCode domainCoin.CoinTypeCode, chainConf *chaincfg.Params) *PublicXRP {
	return &PublicXRP{
		PublicRPC:    xrprpcpublic.NewPublicRPC(ws),
		coinTypeCode: coinTypeCode,
		chainConf:    chainConf,
	}
}

// CoinTypeCode returns the coin type code.
func (p *PublicXRP) CoinTypeCode() domainCoin.CoinTypeCode {
	return p.coinTypeCode
}

// GetChainConf returns the chain configuration.
func (p *PublicXRP) GetChainConf() *chaincfg.Params {
	return p.chainConf
}
