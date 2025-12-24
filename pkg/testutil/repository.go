package testutil

import (
	"database/sql"
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/database/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
)

var (
	// shared database connection
	dbConn *sql.DB
	// sqlc repositories
	btcTxRepoSqlc          *watch.BTCTxRepositorySqlc
	txRepoSqlc             *watch.TxRepositorySqlc
	addressRepoSqlc        *watch.AddressRepositorySqlc
	paymentRequestRepoSqlc *watch.PaymentRequestRepositorySqlc
	btcTxInputRepoSqlc     *watch.TxInputRepositorySqlc
	btcTxOutputRepoSqlc    *watch.TxOutputRepositorySqlc
	ethDetailTxRepoSqlc    *watch.EthDetailTxInputRepositorySqlc
	xrpDetailTxRepoSqlc    *watch.XrpDetailTxInputRepositorySqlc
)

// GetDB returns shared database connection for tests
func GetDB() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
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

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc for test
func NewBTCTxRepositorySqlc() watch.BTCTxRepositorier {
	if btcTxRepoSqlc != nil {
		return btcTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxRepoSqlc = watch.NewBTCTxRepositorySqlc(db, domainCoin.BTC)
	return btcTxRepoSqlc
}

// NewTxRepositorySqlc returns TxRepositorySqlc for test
func NewTxRepositorySqlc() watch.TxRepositorier {
	if txRepoSqlc != nil {
		return txRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	txRepoSqlc = watch.NewTxRepositorySqlc(db, domainCoin.BTC)
	return txRepoSqlc
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc for test
func NewAddressRepositorySqlc() watch.AddressRepositorier {
	if addressRepoSqlc != nil {
		return addressRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	addressRepoSqlc = watch.NewAddressRepositorySqlc(db, domainCoin.BTC)
	return addressRepoSqlc
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc for test
func NewPaymentRequestRepositorySqlc() watch.PaymentRequestRepositorier {
	if paymentRequestRepoSqlc != nil {
		return paymentRequestRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	paymentRequestRepoSqlc = watch.NewPaymentRequestRepositorySqlc(db, domainCoin.BTC)
	return paymentRequestRepoSqlc
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc for test
func NewBTCTxInputRepositorySqlc() watch.TxInputRepositorier {
	if btcTxInputRepoSqlc != nil {
		return btcTxInputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxInputRepoSqlc = watch.NewBTCTxInputRepositorySqlc(db, domainCoin.BTC)
	return btcTxInputRepoSqlc
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc for test
func NewBTCTxOutputRepositorySqlc() watch.TxOutputRepositorier {
	if btcTxOutputRepoSqlc != nil {
		return btcTxOutputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxOutputRepoSqlc = watch.NewBTCTxOutputRepositorySqlc(db, domainCoin.BTC)
	return btcTxOutputRepoSqlc
}

// NewEthDetailTxRepositorySqlc returns EthDetailTxInputRepositorySqlc for test
func NewEthDetailTxRepositorySqlc() watch.EthDetailTxRepositorier {
	if ethDetailTxRepoSqlc != nil {
		return ethDetailTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTxRepoSqlc = watch.NewEthDetailTxInputRepositorySqlc(db, domainCoin.ETH)
	return ethDetailTxRepoSqlc
}

// NewXrpDetailTxRepositorySqlc returns XrpDetailTxInputRepositorySqlc for test
func NewXrpDetailTxRepositorySqlc() watch.XrpDetailTxRepositorier {
	if xrpDetailTxRepoSqlc != nil {
		return xrpDetailTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/xrp_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTxRepoSqlc = watch.NewXrpDetailTxInputRepositorySqlc(db, domainCoin.XRP)
	return xrpDetailTxRepoSqlc
}
