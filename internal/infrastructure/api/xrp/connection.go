package xrp

import (
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// NewXRPFromCoinType creates XRP instance according to coinType
// FIXME: coinType doesn't care
func NewXRPFromCoinType(
	wsPublic *websocket.WS, wsAdmin *websocket.WS, conf *config.Ripple,
	coinTypeCode domainCoin.CoinTypeCode,
) (apixrp.XRPer, error) {
	switch coinTypeCode {
	case domainCoin.XRP:
		xrp, err := NewXRP(wsPublic, wsAdmin, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call xrp.NewXRP(): %w", err)
		}
		return xrp, err
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
