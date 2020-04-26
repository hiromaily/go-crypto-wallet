package main

import (
	"database/sql"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tracer"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwallet"
)

// Registry is for registry interface
type Registry interface {
	NewSigner() wallets.Signer
}

type registry struct {
	conf        *config.Config
	mysqlClient *sql.DB
	logger      *zap.Logger
	rpcClient   *rpcclient.Client
	btc         api.Bitcoiner
	walletType  types.WalletType
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType types.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewSigner is to register for Signer interface
func (r *registry) NewSigner() wallets.Signer {
	return coldwallet.NewColdWalet(
		r.newBTC(),
		r.newLogger(),
		r.newTracer(),
		r.newRepository(),
		r.newKeyGenerator(),
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
	var err error
	if r.btc == nil {
		r.btc, err = api.NewBitcoin(r.newRPCClient(), &r.conf.Bitcoin, r.newLogger(), r.conf.CoinTypeCode)
		if err != nil {
			panic(err)
		}
	}
	return r.btc
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

//func (r *registry) newStorager() rdb.SignatureStorager {
func (r *registry) newRepository() coldrepo.ColdRepository {
	// if there are multiple options, set proper one
	// storager interface as MySQL
	return coldrepo.NewColdWalletRepository(
		r.newMySQLClient(),
		r.newLogger(),
		r.newSeedRepo(),
		r.newAccountKeyRepo(),
		r.newMultisigRepo(),
	)
}

func (r *registry) newSeedRepo() coldrepo.SeedRepository {
	return coldrepo.NewSeedRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAccountKeyRepo() coldrepo.AccountKeyRepository {
	return coldrepo.NewAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newMultisigRepo() coldrepo.MultisigHistoryRepository {
	return coldrepo.NewMultisigHistoryRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newKeyGenerator() key.Generator {
	return key.NewHDKey(
		key.PurposeTypeBIP44,
		r.conf.CoinTypeCode,
		r.newBTC().GetChainConf(),
		r.newLogger())
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	return r.mysqlClient
}

func (r *registry) newAddressFileStorager() address.Storager {
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
