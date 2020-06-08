package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var (
	txRepo         *watchrepo.BTCTxRepository
	accountKeyRepo *coldrepo.AccountKeyRepository
)

// NewTxRepository returns TxRepository for test
func NewTxRepository() watchrepo.BTCTxRepositorier {
	if txRepo != nil {
		return txRepo
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc_watch.toml", projPath)
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
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

	txRepo = watchrepo.NewBTCTxRepository(db, coin.BTC, logger)
	return txRepo
}

// NewAccountKeyRepository returns AccountKeyRepository for test
func NewAccountKeyRepository() coldrepo.AccountKeyRepositorier {
	if accountKeyRepo != nil {
		return accountKeyRepo
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc_watch.toml", projPath)
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
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
