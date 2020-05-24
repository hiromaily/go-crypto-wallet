package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var xr xrpgrp.Rippler

// GetXRP returns xrp instance
//FIXME: hard coded
func GetXRP() xrpgrp.Rippler {
	if xr != nil {
		return xr
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/xrp_watch.toml", projPath)
	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly, coin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overridden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := xrpgrp.NewWSClient(&conf.Ripple)
	if err != nil {
		log.Fatalf("fail to create ethereum rpc client: %v", err)
	}
	xr, err = xrpgrp.NewRipple(client, nil, &conf.Ripple, logger, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create xrp instance: %v", err)
	}
	return xr
}
