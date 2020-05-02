package main

import (
	"database/sql"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/opentracing/opentracing-go"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/repository/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tracer"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/wallet"
)

// Registry is for registry interface
type Registry interface {
	NewWalleter() wallets.Walleter
}

type registry struct {
	conf        *config.Config
	walletType  wtype.WalletType
	logger      *zap.Logger
	rpcClient   *rpcclient.Client
	mysqlClient *sql.DB
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType wtype.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewWalleter is to register for walleter interface
func (r *registry) NewWalleter() wallets.Walleter {
	return wallet.NewWallet(
		r.newBTC(),
		r.newLogger(),
		r.newTracer(),
		r.newRepository(),
		r.newAddressFileStorager(),
		r.newTxFileStorager(),
		r.walletType,
	)
}

func (r *registry) newRPCClient() *rpcclient.Client {
	var err error
	if r.rpcClient == nil {
		r.rpcClient, err = api.NewRPCClient(&r.conf.Bitcoin)
	}
	if err != nil {
		panic(err)
	}
	return r.rpcClient
}

func (r *registry) newBTC() api.Bitcoiner {
	bit, err := api.NewBitcoin(r.newRPCClient(), &r.conf.Bitcoin, r.newLogger(), r.conf.CoinTypeCode)
	if err != nil {
		panic(err)
	}
	return bit
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(&r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newTracer() opentracing.Tracer {
	return tracer.NewTracer(r.conf.Tracer)
}

func (r *registry) newRepository() walletrepo.WalletRepository {
	// if there are multiple options, set proper one
	// storager interface as MySQL
	return walletrepo.NewWalletRepository(
		r.newMySQLClient(),
		r.newLogger(),
		r.newTxRepo(),
		r.newTxInputRepo(),
		r.newTxOutputRepo(),
		r.newPaymentRequestRepo(),
		r.newAddressRepo(),
	)
}

func (r *registry) newTxRepo() walletrepo.TxRepository {
	return walletrepo.NewTxRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxInputRepo() walletrepo.TxInputRepository {
	return walletrepo.NewTxInputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxOutputRepo() walletrepo.TxOutputRepository {
	return walletrepo.NewTxOutputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newPaymentRequestRepo() walletrepo.PaymentRequestRepository {
	return walletrepo.NewPaymentRequestRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAddressRepo() walletrepo.AddressRepository {
	return walletrepo.NewAddressRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	if r.conf.MySQL.Debug {
		boil.DebugMode = true
	}
	return r.mysqlClient
}

func (r *registry) newAddressFileStorager() address.FileStorager {
	return address.NewFileRepository(
		r.conf.PubkeyFile.BasePath,
		r.newLogger(),
	)
}

func (r *registry) newTxFileStorager() tx.FileStorager {
	return tx.NewFileRepository(
		r.conf.TxFile.BasePath,
		r.newLogger(),
	)
}
