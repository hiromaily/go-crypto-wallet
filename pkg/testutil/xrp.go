package testutil

import (
	"fmt"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var (
	xr xrpgrp.Rippler
)

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
	// ws client
	wsClient, wsAdmin, err := xrpgrp.NewWSClient(&conf.Ripple)
	if err != nil {
		log.Fatalf("fail to create ethereum rpc client: %v", err)
	}
	// client
	conn, err := xrpgrp.NewGRPCClient(&conf.Ripple.API)
	if err != nil {
		log.Fatalf("fail to create api instance: %v", err)
	}
	grpcAPI := xrp.NewRippleAPI(conn, logger)

	xr, err = xrpgrp.NewRipple(wsClient, wsAdmin, grpcAPI, nil, &conf.Ripple, logger, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create xrp instance: %v", err)
	}
	return xr
}

// GetRippleAPI returns RippleAPIer
//func GetRippleAPI() xrpgrp.RippleAPIer {
//	if api != nil {
//		return api
//	}
//
//	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
//	confPath := fmt.Sprintf("%s/data/config/xrp_watch.toml", projPath)
//	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly, coin.XRP)
//	if err != nil {
//		log.Fatalf("fail to create config: %v", err)
//	}
//	//TODO: if config should be overridden, here
//
//	// client
//	conn, err := xrpgrp.NewGRPCClient(&conf.Ripple.API)
//	if err != nil {
//		log.Fatalf("fail to create api instance: %v", err)
//	}
//	if conn == nil {
//		log.Fatal("connection is nil")
//	}
//	logger := logger.NewZapLogger(&conf.Logger)
//	api = xrp.NewRippleAPI(conn, logger)
//
//	return api
//}
