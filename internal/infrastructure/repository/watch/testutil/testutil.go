package testutil

import (
	"log"
	"os"

	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	watchmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
)

var (
	// sqlc repositories
	btcTxRepoSqlc          *watchmysql.BTCTxRepositorySqlc
	txRepoSqlc             *watchmysql.TxRepositorySqlc
	addressRepoSqlc        *watchmysql.AddressRepositorySqlc
	paymentRequestRepoSqlc *watchmysql.PaymentRequestRepositorySqlc
	btcTxInputRepoSqlc     *watchmysql.TxInputRepositorySqlc
	btcTxOutputRepoSqlc    *watchmysql.TxOutputRepositorySqlc
	ethDetailTXRepoSqlc    *watchmysql.ETHDetailTXInputRepositorySqlc
	xrpDetailTXRepoSqlc    *watchmysql.XRPDetailTxInputRepositorySqlc
)

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc for test
func NewBTCTxRepositorySqlc() repowatch.BTCTxRepositorier {
	if btcTxRepoSqlc != nil {
		return btcTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxRepoSqlc = watchmysql.NewBTCTxRepositorySqlc(db, domainCoin.BTC)
	return btcTxRepoSqlc
}

// NewTxRepositorySqlc returns TxRepositorySqlc for test
func NewTxRepositorySqlc() repowatch.TxRepositorier {
	if txRepoSqlc != nil {
		return txRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	txRepoSqlc = watchmysql.NewTxRepositorySqlc(db, domainCoin.BTC)
	return txRepoSqlc
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc for test
func NewAddressRepositorySqlc() repowatch.AddressRepositorier {
	if addressRepoSqlc != nil {
		return addressRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	addressRepoSqlc = watchmysql.NewAddressRepositorySqlc(db, domainCoin.BTC)
	return addressRepoSqlc
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc for test
func NewPaymentRequestRepositorySqlc() repowatch.PaymentRequestRepositorier {
	if paymentRequestRepoSqlc != nil {
		return paymentRequestRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	paymentRequestRepoSqlc = watchmysql.NewPaymentRequestRepositorySqlc(db, domainCoin.BTC)
	return paymentRequestRepoSqlc
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc for test
func NewBTCTxInputRepositorySqlc() repowatch.TxInputRepositorier {
	if btcTxInputRepoSqlc != nil {
		return btcTxInputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxInputRepoSqlc = watchmysql.NewBTCTxInputRepositorySqlc(db, domainCoin.BTC)
	return btcTxInputRepoSqlc
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc for test
func NewBTCTxOutputRepositorySqlc() repowatch.TxOutputRepositorier {
	if btcTxOutputRepoSqlc != nil {
		return btcTxOutputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	btcTxOutputRepoSqlc = watchmysql.NewBTCTxOutputRepositorySqlc(db, domainCoin.BTC)
	return btcTxOutputRepoSqlc
}

// NewETHDetailTXRepositorySqlc returns ETHDetailTXInputRepositorySqlc for test
func NewETHDetailTXRepositorySqlc() repowatch.ETHDetailTXRepositorier {
	if ethDetailTXRepoSqlc != nil {
		return ethDetailTXRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/eth/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTXRepoSqlc = watchmysql.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)
	return ethDetailTXRepoSqlc
}

// NewXrpDetailTxRepositorySqlc returns XRPDetailTxInputRepositorySqlc for test
func NewXrpDetailTxRepositorySqlc() repowatch.XRPDetailTXRepositorier {
	if xrpDetailTXRepoSqlc != nil {
		return xrpDetailTXRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/xrp/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTXRepoSqlc = watchmysql.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)
	return xrpDetailTXRepoSqlc
}
