package xrpgrp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// NewRipple creates Ripple instance according to coinType
func NewRipple(conf *config.Ripple, logger *zap.Logger, coinTypeCode coin.CoinTypeCode) (Rippler, error) {
	switch coinTypeCode {
	case coin.XRP:
		eth, err := xrp.NewRipple(context.Background(), coinTypeCode, conf, logger)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call xrp.NewRipple()")
		}
		return eth, err
	}
	return nil, errors.Errorf("coinType %s is not defined", coinTypeCode.String())
}
