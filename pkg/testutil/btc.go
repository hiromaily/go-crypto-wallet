package testutil

import (
	"fmt"
	"log"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

var bc api.Bitcoiner

//FIXME: hard coded
func GetBTC() api.Bitcoiner {
	if bc != nil {
		return bc
	}
	// create bitcoin instance
	// config
	//projPath := os.Getenv("PROJECT_PATH")
	//TODO: delete it
	projPath := "/Users/hy/work/go/src/github.com/hiromaily/go-bitcoin"

	if projPath == "" {
		log.Fatalf("$PROJECT_PATH should be defined as environment variable")
	}

	confPath := fmt.Sprintf("%s/data/config/btc/wallet.toml", projPath)
	conf, err := config.New(confPath, types.WalletTypeWatchOnly)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overridden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := api.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		log.Fatalf("fail to create bitcoin core client: %v", err)
	}
	bc, err = api.NewBitcoin(client, &conf.Bitcoin, logger, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create btc instance: %v", err)
	}
	return bc
}
