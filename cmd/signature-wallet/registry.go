package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/signaturerepo"
	"github.com/hiromaily/go-bitcoin/pkg/tracer"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
)

// Registry is for registry interface
type Registry interface {
	NewSigner() wallets.Signer
}

type registry struct {
	conf        *config.Config
	mysqlClient *sqlx.DB
	walletType  wallets.WalletType
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType wallets.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewSigner is to register for Signer interface
func (r *registry) NewSigner() wallets.Signer {
	//TODO: should be interface
	r.setFilePath()

	//FIXME: wallet.NewWallet doesn't have rdb.SignatureStorager
	// How should it be fixed?? NewSignature should be defined based on NewWallet
	return wallets.NewSignature(
		r.newBTC(),
		r.newLogger(),
		r.newTracer(),
		r.newStorager(),
		r.walletType,
	)
}

func (r *registry) newBTC() api.Bitcoiner {
	// Connection to Bitcoin core
	// TODO: coinType should be judged here
	// TODO: name should be NewBitcoin or NewBitcoinCash
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

func (r *registry) newStorager() rdb.SignatureStorager {
	// if there are multiple options, set proper one
	// storager interface as MySQL
	return signaturerepo.NewSignatureRepository(
		r.newMySQLClient(),
		r.newLogger(),
	)
}

func (r *registry) newMySQLClient() *sqlx.DB {
	if r.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&r.conf.MySQL)
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
