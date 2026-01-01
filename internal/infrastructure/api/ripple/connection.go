package ripple

import (
	"context"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// NewRipple creates Ripple instance according to coinType
func NewRipple(
	wsPublic *websocket.WS, wsAdmin *websocket.WS, api *xrp.RippleAPI, conf *config.Ripple,
	coinTypeCode domainCoin.CoinTypeCode,
) (Rippler, error) {
	//nolint:exhaustive
	switch coinTypeCode {
	case domainCoin.XRP:
		ripple, err := xrp.NewRipple(context.Background(), wsPublic, wsAdmin, api, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call xrp.NewRipple(): %w", err)
		}
		return ripple, err
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
