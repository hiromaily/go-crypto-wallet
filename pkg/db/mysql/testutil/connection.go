package testutil

import (
	"database/sql"
	"log"
	"os"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
)

// shared database connection
var dbConn *sql.DB

// GetDB returns shared database connection for tests
func GetDB() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	dbConn, err = mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	return dbConn
}
