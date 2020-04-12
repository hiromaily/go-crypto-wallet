package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

var bc api.Bitcoiner

func GetBTC() api.Bitcoiner {
	if bc != nil {
		return bc
	}
	// create bitcoin instance
	// config
	projPath := os.Getenv("PROJECT_PATH")
	//TODO: delete it
	projPath = "/Users/hy/work/go/src/github.com/hiromaily/go-bitcoin"

	if projPath == "" {
		log.Fatalf("$PROJECT_PATH should be defined as environment variable")
	}

	confPath := fmt.Sprintf("%s/data/config/btc/wallet.toml", projPath)
	conf, err := config.New(confPath, types.WalletTypeWatchOnly)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overriden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := api.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		log.Fatalf("fail to create bitcoin core client: %v", err)
	}
	bc, err = api.NewBitcoin(client, &conf.Bitcoin, logger, ctype.CoinType(conf.CoinType))
	if err != nil {
		log.Fatalf("fail to create btc instance: %v", err)
	}
	return bc
}
