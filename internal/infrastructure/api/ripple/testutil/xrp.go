package testutil

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency"
	"github.com/hiromaily/go-crypto-wallet/pkg/grpc"
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
	wsPublicClient, err := cryptocurrency.NewWebSocketClient(conf.Ripple.WebsocketPublicURL)
	if err != nil {
		return nil, fmt.Errorf("fail to create xrp public websocket client: %w", err)
	}
	wsAdminClient, err := cryptocurrency.NewWebSocketClient(conf.Ripple.WebsocketAdminURL)
	if err != nil {
		return nil, fmt.Errorf("fail to create xrp admin websocket client: %w", err)
	}
	// client
	conn, err := grpc.NewClient(conf.Ripple.API.URL)
	if err != nil {
		return nil, fmt.Errorf("fail to create api instance: %w", err)
	}
	grpcAPI := xrp.NewRippleAPI(conn)

	xr, err = ripple.NewRipple(wsPublicClient, wsAdminClient, grpcAPI, &conf.Ripple, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create xrp instance: %w", err)
	}
	return xr, nil
}

// XRPTestSuite is a test suite for XRP
type XRPTestSuite struct {
	suite.Suite
	XRP ripple.Rippler
}

func (xts *XRPTestSuite) SetupTest() {
	xrp, err := GetXRP()
	xts.NoError(err)
	xts.XRP = xrp
}

func (xts *XRPTestSuite) TearDownTest() {
	_ = xts.XRP.Close() // Best effort cleanup
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
