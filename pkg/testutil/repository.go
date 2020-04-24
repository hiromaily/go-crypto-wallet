package testutil

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/repository"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

var txRepo repository.TxRepository

func NewTxRepository() repository.TxRepository {
	if txRepo != nil {
		return txRepo
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-bitcoin", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc/wallet.toml", projPath)
	conf, err := config.New(confPath, types.WalletTypeWatchOnly)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overridden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)

	// db
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	txRepo = repository.NewTxRepository(db, coin.BTC, logger)
	return txRepo
}
