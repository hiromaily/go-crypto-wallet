package testutil

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
)

var xr ripple.Rippler

// GetXRP returns xrp instance
// FIXME: hard coded
func GetXRP() (ripple.Rippler, error) {
	if xr != nil {
		return xr, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/xrp_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = domainCoin.XRP

	// ws client
	wsClient, wsAdmin, err := ripple.NewWSClient(&conf.Ripple)
	if err != nil {
		return nil, fmt.Errorf("fail to create ethereum rpc client: %w", err)
	}
	// client
	conn, err := ripple.NewGRPCClient(&conf.Ripple.API)
	if err != nil {
		return nil, fmt.Errorf("fail to create api instance: %w", err)
	}
	grpcAPI := xrp.NewRippleAPI(conn)

	xr, err = ripple.NewRipple(wsClient, wsAdmin, grpcAPI, &conf.Ripple, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create xrp instance: %w", err)
	}
	return xr, nil
}

// GetRippleAPI returns RippleAPIer
// func GetRippleAPI() ripple.RippleAPIer {
//	if api != nil {
//		return api
//	}
//
//	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
//	confPath := fmt.Sprintf("%s/data/config/xrp_watch.toml", projPath)
//	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly, domainCoin.XRP)
//	if err != nil {
//		log.Fatalf("fail to create config: %v", err)
//	}
//	//TODO: if config should be overridden, here
//
//	// client
//	conn, err := ripple.NewGRPCClient(&conf.Ripple.API)
//	if err != nil {
//		log.Fatalf("fail to create api instance: %v", err)
//	}
//	if conn == nil {
//		log.Fatal("connection is nil")
//	}
//	logger := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
//	api = xrp.NewRippleAPI(conn, logger)
//
//	return api
//}
