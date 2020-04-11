package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/tracer"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/opentracing/opentracing-go"
)

// Registry is for registry interface
type Registry interface {
	NewWalleter() wallet.Walleter
}

type registry struct {
	conf        *config.Config
	mysqlClient *sqlx.DB
	walletType  wallet.WalletType
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config, walletType wallet.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewWalleter is to register for walleter interface
func (r *registry) NewWalleter() wallet.Walleter {
	//TODO: should be interface
	r.setFilePath()

	return wallet.NewWallet(
		r.newBTC(),
		r.newLogger(),
		r.newTracer(),
		r.newStorager(),
		r.walletType,
	)
}

func (r *registry) newBTC() api.Bitcoiner {
	// Connection to Bitcoin core
	bit, err := api.Connection(&r.conf.Bitcoin, enum.CoinType(r.conf.CoinType))
	if err != nil {
		panic(fmt.Sprintf("btc.Connection error: %s", err))
	}
	return bit
}

func (r *registry) newLogger() *zap.Logger {
	return logger.NewZapLogger(&r.conf.Logger)
}

func (r *registry) newTracer() opentracing.Tracer {
	return tracer.NewTracer(r.conf.Tracer)
}

//TODO: change return as interface
func (r *registry) newStorager() *model.DB {
	// if there are multiple options, set proper one
	// storager interface as MySQL
	return db.NewGenreItemRepository(r.NewMySQLSlaveClient())
	//return model.NewDB(r.newMySQLClient())
}

func (r *registry) newMySQLClient() *sqlx.DB {
	if r.mysqlClient == nil {
		dbConn, err := rdb.NewMySQL(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	return r.mysqlClient
}

//TODO: move to somewhere
func (r *registry) setFilePath() {
	// TxFile
	if r.conf.TxFile.BasePath != "" {
		txfile.SetFilePath(r.conf.TxFile.BasePath)
	}

	// PubkeyCSV
	if r.conf.PubkeyFile.BasePath != "" {
		key.SetFilePath(r.conf.PubkeyFile.BasePath)
	}
}
