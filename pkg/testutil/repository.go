package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/repository/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

var (
	txRepo         walletrepo.TxRepository
	accountKeyRepo *coldrepo.AccountKeyRepository
)

// NewTxRepository returns TxRepository for test
func NewTxRepository() walletrepo.TxRepository {
	if txRepo != nil {
		return txRepo
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-bitcoin", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc/watch.toml", projPath)
	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly)
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

	txRepo = walletrepo.NewTxRepository(db, coin.BTC, logger)
	return txRepo
}

// NewAccountKeyRepository returns AccountKeyRepository for test
func NewAccountKeyRepository() coldrepo.AccountKeyRepositorier {
	if accountKeyRepo != nil {
		return accountKeyRepo
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-bitcoin", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc/watch.toml", projPath)
	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly)
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

	accountKeyRepo = coldrepo.NewAccountKeyRepository(db, coin.BTC, logger)
	return accountKeyRepo
}
