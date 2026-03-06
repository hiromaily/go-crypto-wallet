package xrp

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	xrpadmin "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/admin"
	xrppublic "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/public"
	pkgxrp "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// chainConf resolves the chain configuration from the Ripple config.
func chainConf(conf *config.Ripple) *chaincfg.Params {
	if conf.NetworkType != pkgxrp.NetworkTypeXRPMainNet.String() {
		return &chaincfg.TestNet3Params
	}
	return &chaincfg.MainNetParams
}

// NewPublicXRP creates a PublicXRP instance for the given WebSocket and coin type.
func NewPublicXRP(
	ws *websocket.WS, conf *config.Ripple, coinTypeCode domainCoin.CoinTypeCode,
) (*xrppublic.PublicXRP, error) {
	if ws == nil {
		return nil, errors.New("public WebSocket connection is required")
	}
	return xrppublic.NewPublicXRP(ws, coinTypeCode, chainConf(conf)), nil
}

// NewAdminXRP creates an AdminXRP instance for the given WebSocket and coin type.
// Returns nil, nil when wsAdmin is nil (admin not configured).
func NewAdminXRP(
	ws *websocket.WS, conf *config.Ripple, coinTypeCode domainCoin.CoinTypeCode,
) (*xrpadmin.AdminXRP, error) {
	if ws == nil {
		return nil, nil
	}
	return xrpadmin.NewAdminXRP(ws, coinTypeCode, chainConf(conf)), nil
}

// NewPublicXRPFromCoinType creates a PublicXRP instance for the given coin type.
// Returns an error for unsupported coin types.
func NewPublicXRP(
	wsPublic *websocket.WS, conf *config.Ripple, coinTypeCode domainCoin.CoinTypeCode,
) (apixrp.XRPPublicClient, error) {
	switch coinTypeCode {
	case domainCoin.XRP:
		return NewPublicXRP(wsPublic, conf, coinTypeCode)
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}

// NewAdminXRPFromCoinType creates an AdminXRP instance for the given coin type.
// Returns nil, nil when wsAdmin is nil (admin not configured).
func NewAdminXRPFromCoinType(
	wsAdmin *websocket.WS, conf *config.Ripple, coinTypeCode domainCoin.CoinTypeCode,
) (apixrp.XRPAdminClient, error) {
	switch coinTypeCode {
	case domainCoin.XRP:
		client, err := NewAdminXRP(wsAdmin, conf, coinTypeCode)
		if err != nil || client == nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
