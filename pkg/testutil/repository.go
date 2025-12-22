package testutil

import (
	"database/sql"
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
	// shared database connection
	dbConn         *sql.DB
	txRepo         *watchrepo.BTCTxRepository
	accountKeyRepo *coldrepo.AccountKeyRepository
	// sqlc repositories
	btcTxRepoSqlc          *watchrepo.BTCTxRepositorySqlc
	txRepoSqlc             *watchrepo.TxRepositorySqlc
	addressRepoSqlc        *watchrepo.AddressRepositorySqlc
	paymentRequestRepoSqlc *watchrepo.PaymentRequestRepositorySqlc
	btcTxInputRepoSqlc     *watchrepo.TxInputRepositorySqlc
	btcTxOutputRepoSqlc    *watchrepo.TxOutputRepositorySqlc
	ethDetailTxRepoSqlc    *watchrepo.EthDetailTxInputRepositorySqlc
	xrpDetailTxRepoSqlc    *watchrepo.XrpDetailTxInputRepositorySqlc
)

// GetDB returns shared database connection for tests
func GetDB() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	dbConn, err = mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	return dbConn
}

// NewTxRepository returns TxRepository for test
func NewTxRepository() watchrepo.BTCTxRepositorier {
	if txRepo != nil {
		return txRepo
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	// TODO: if config should be overridden, here

	// logger
	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)

	// db
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	txRepo = watchrepo.NewBTCTxRepository(db, coin.BTC, zapLog)
	return txRepo
}

// NewAccountKeyRepository returns AccountKeyRepository for test
func NewAccountKeyRepository() coldrepo.AccountKeyRepositorier {
	if accountKeyRepo != nil {
		return accountKeyRepo
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	// TODO: if config should be overridden, here

	// logger
	zapLogger := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)

	// db
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	accountKeyRepo = coldrepo.NewAccountKeyRepository(db, coin.BTC, zapLogger)
	return accountKeyRepo
}

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc for test
func NewBTCTxRepositorySqlc() watchrepo.BTCTxRepositorier {
	if btcTxRepoSqlc != nil {
		return btcTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxRepoSqlc = watchrepo.NewBTCTxRepositorySqlc(db, coin.BTC, zapLog)
	return btcTxRepoSqlc
}

// NewTxRepositorySqlc returns TxRepositorySqlc for test
func NewTxRepositorySqlc() watchrepo.TxRepositorier {
	if txRepoSqlc != nil {
		return txRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	txRepoSqlc = watchrepo.NewTxRepositorySqlc(db, coin.BTC, zapLog)
	return txRepoSqlc
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc for test
func NewAddressRepositorySqlc() watchrepo.AddressRepositorier {
	if addressRepoSqlc != nil {
		return addressRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	addressRepoSqlc = watchrepo.NewAddressRepositorySqlc(db, coin.BTC, zapLog)
	return addressRepoSqlc
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc for test
func NewPaymentRequestRepositorySqlc() watchrepo.PaymentRequestRepositorier {
	if paymentRequestRepoSqlc != nil {
		return paymentRequestRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	paymentRequestRepoSqlc = watchrepo.NewPaymentRequestRepositorySqlc(db, coin.BTC, zapLog)
	return paymentRequestRepoSqlc
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc for test
func NewBTCTxInputRepositorySqlc() watchrepo.TxInputRepositorier {
	if btcTxInputRepoSqlc != nil {
		return btcTxInputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxInputRepoSqlc = watchrepo.NewBTCTxInputRepositorySqlc(db, coin.BTC, zapLog)
	return btcTxInputRepoSqlc
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc for test
func NewBTCTxOutputRepositorySqlc() watchrepo.TxOutputRepositorier {
	if btcTxOutputRepoSqlc != nil {
		return btcTxOutputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxOutputRepoSqlc = watchrepo.NewBTCTxOutputRepositorySqlc(db, coin.BTC, zapLog)
	return btcTxOutputRepoSqlc
}

// NewEthDetailTxRepositorySqlc returns EthDetailTxInputRepositorySqlc for test
func NewEthDetailTxRepositorySqlc() watchrepo.EthDetailTxRepositorier {
	if ethDetailTxRepoSqlc != nil {
		return ethDetailTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTxRepoSqlc = watchrepo.NewEthDetailTxInputRepositorySqlc(db, coin.ETH, zapLog)
	return ethDetailTxRepoSqlc
}

// NewXrpDetailTxRepositorySqlc returns XrpDetailTxInputRepositorySqlc for test
func NewXrpDetailTxRepositorySqlc() watchrepo.XrpDetailTxRepositorier {
	if xrpDetailTxRepoSqlc != nil {
		return xrpDetailTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/xrp_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTxRepoSqlc = watchrepo.NewXrpDetailTxInputRepositorySqlc(db, coin.XRP, zapLog)
	return xrpDetailTxRepoSqlc
}
