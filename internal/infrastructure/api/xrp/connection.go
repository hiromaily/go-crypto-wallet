package xrp

import (
	"context"
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// NewRippleFromCoinType creates Ripple instance according to coinType
func NewRippleFromCoinType(
	wsPublic *websocket.WS, wsAdmin *websocket.WS, api *RippleAPI, conf *config.Ripple,
	coinTypeCode domainCoin.CoinTypeCode,
) (apixrp.Rippler, error) {
	//nolint:exhaustive
	switch coinTypeCode {
	case domainCoin.XRP:
		ripple, err := NewRipple(context.Background(), wsPublic, wsAdmin, api, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call xrp.NewRipple(): %w", err)
		}
		return ripple, err
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
