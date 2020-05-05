package main

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv/signsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// Registry is for registry interface
type Registry interface {
	NewSigner() wallets.Signer
}

type registry struct {
	conf        *config.Config
	walletType  wallet.WalletType
	authType    account.AuthType
	logger      *zap.Logger
	btc         api.Bitcoiner
	rpcClient   *rpcclient.Client
	mysqlClient *sql.DB
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType wallet.WalletType, authName string) Registry {
	// validate
	if !account.ValidateAuthType(authName) {
		panic(fmt.Sprintf("authName is invalid. this should be embedded when building: %s", authName))
	}

	return &registry{
		conf:       conf,
		walletType: walletType,
		authType:   account.AuthTypeMap[authName],
	}
}

// NewSigner is to register for Signer interface
// NewKeygener is to register for keygener interface
func (r *registry) NewSigner() wallets.Signer {
	return wallets.NewSign(
		r.newBTC(),
		r.newMySQLClient(),
		r.authType,
		r.newSeeder(),
		r.newHdWallter(),
		r.newPrivKeyer(),
		r.newFullPubkeyExporter(),
		r.newSigner(),
		r.walletType,
	)
}

func (r *registry) newSeeder() coldsrv.Seeder {
	return coldsrv.NewSeed(
		r.newLogger(),
		r.newSeedRepo(),
		r.walletType,
	)
}

func (r *registry) newHdWallter() coldsrv.HDWalleter {
	return coldsrv.NewHDWallet(
		r.newLogger(),
		r.newHdWalletRepo(),
		r.newKeyGenerator(),
		r.conf.CoinTypeCode,
		r.walletType,
	)
}

func (r *registry) newHdWalletRepo() coldsrv.HDWalletRepo {
	return coldsrv.NewAuthHDWalletRepo(
		r.newAuthKeyRepo(),
		r.authType,
	)
}

func (r *registry) newPrivKeyer() signsrv.PrivKeyer {
	return signsrv.NewPrivKey(
		r.newBTC(),
		r.newLogger(),
		r.newAuthKeyRepo(),
		r.authType,
		r.walletType,
	)
}

func (r *registry) newFullPubkeyExporter() signsrv.FullPubkeyExporter {
	return signsrv.NewFullPubkeyExport(
		r.newLogger(),
		r.newAuthKeyRepo(),
		r.newPubkeyFileStorager(),
		r.conf.CoinTypeCode,
		r.authType,
		r.walletType,
	)
}
func (r *registry) newSigner() coldsrv.Signer {
	return coldsrv.NewSign(
		r.newBTC(),
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.newAuthKeyRepo(),
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
	if r.btc == nil {
		var err error
		r.btc, err = api.NewBitcoin(
			r.newRPCClient(),
			&r.conf.Bitcoin,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
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

func (r *registry) newSeedRepo() coldrepo.SeedRepositorier {
	return coldrepo.NewSeedRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAccountKeyRepo() coldrepo.AccountKeyRepositorier {
	return coldrepo.NewAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAuthKeyRepo() coldrepo.AuthAccountKeyRepositorier {
	return coldrepo.NewAuthAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.authType,
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
	if r.conf.MySQL.Debug {
		boil.DebugMode = true
	}
	return r.mysqlClient
}

func (r *registry) newPubkeyFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.FilePath.FullPubKey,
		r.newLogger(),
	)
}

func (r *registry) newTxFileStorager() tx.FileRepositorier {
	return tx.NewFileRepository(
		r.conf.FilePath.Tx,
		r.newLogger(),
	)
}
