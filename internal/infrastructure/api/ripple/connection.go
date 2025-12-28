package ripple

import (
	"context"
	"errors"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// NewWSClient try to connect Ripple Server by web socket
func NewWSClient(conf *config.Ripple) (*websocket.WS, *websocket.WS, error) {
	publicURL := conf.WebsocketPublicURL
	if publicURL == "" {
		if publicURL = xrp.GetPublicWSServer(conf.NetworkType).String(); publicURL == "" {
			return nil, nil, errors.New("websocket URL is not found")
		}
	}
	public, err := websocket.New(context.Background(), publicURL)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call websocket.New() for public API: %s: %w", publicURL, err)
	}

	// acceptable without adminClient
	adminURL := conf.WebsocketAdminURL
	if adminURL == "" {
		return public, nil, nil
	}
	admin, err := websocket.New(context.Background(), adminURL)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call websocket.New() for admin API: %s: %w", adminURL, err)
	}

	return public, admin, nil
}

// NewRipple creates Ripple instance according to coinType
func NewRipple(
	wsPublic *websocket.WS, wsAdmin *websocket.WS, api *xrp.RippleAPI, conf *config.Ripple,
	coinTypeCode domainCoin.CoinTypeCode,
) (Rippler, error) {
	switch coinTypeCode {
	case domainCoin.XRP:
		ripple, err := xrp.NewRipple(context.Background(), wsPublic, wsAdmin, api, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call xrp.NewRipple(): %w", err)
		}
		return ripple, err
	case domainCoin.BTC, domainCoin.BCH, domainCoin.LTC, domainCoin.ETH, domainCoin.ERC20, domainCoin.HYT:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
