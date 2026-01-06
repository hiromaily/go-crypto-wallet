package testutil

import (
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
)

var (
	// sqlc repositories
	btcTxRepoSqlc          *watch.BTCTxRepositorySqlc
	txRepoSqlc             *watch.TxRepositorySqlc
	addressRepoSqlc        *watch.AddressRepositorySqlc
	paymentRequestRepoSqlc *watch.PaymentRequestRepositorySqlc
	btcTxInputRepoSqlc     *watch.TxInputRepositorySqlc
	btcTxOutputRepoSqlc    *watch.TxOutputRepositorySqlc
	ethDetailTXRepoSqlc    *watch.ETHDetailTXInputRepositorySqlc
	xrpDetailTXRepoSqlc    *watch.XRPDetailTxInputRepositorySqlc
)

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc for test
func NewBTCTxRepositorySqlc() persistence.BTCTxRepositorier {
	if btcTxRepoSqlc != nil {
		return btcTxRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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
func NewTxRepositorySqlc() persistence.TxRepositorier {
	if txRepoSqlc != nil {
		return txRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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
func NewAddressRepositorySqlc() persistence.AddressRepositorier {
	if addressRepoSqlc != nil {
		return addressRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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
func NewPaymentRequestRepositorySqlc() persistence.PaymentRequestRepositorier {
	if paymentRequestRepoSqlc != nil {
		return paymentRequestRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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
func NewBTCTxInputRepositorySqlc() persistence.TxInputRepositorier {
	if btcTxInputRepoSqlc != nil {
		return btcTxInputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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
func NewBTCTxOutputRepositorySqlc() persistence.TxOutputRepositorier {
	if btcTxOutputRepoSqlc != nil {
		return btcTxOutputRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc_watch.toml"
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

// NewETHDetailTXRepositorySqlc returns ETHDetailTXInputRepositorySqlc for test
func NewETHDetailTXRepositorySqlc() persistence.ETHDetailTXRepositorier {
	if ethDetailTXRepoSqlc != nil {
		return ethDetailTXRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTXRepoSqlc = watch.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)
	return ethDetailTXRepoSqlc
}

// NewXrpDetailTxRepositorySqlc returns XRPDetailTxInputRepositorySqlc for test
func NewXrpDetailTxRepositorySqlc() persistence.XRPDetailTXRepositorier {
	if xrpDetailTXRepoSqlc != nil {
		return xrpDetailTXRepoSqlc
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/xrp_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTXRepoSqlc = watch.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)
	return xrpDetailTXRepoSqlc
}
