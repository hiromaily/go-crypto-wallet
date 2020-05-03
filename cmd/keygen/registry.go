package main

import (
	"database/sql"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

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
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service/coldsrv/keygensrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// Registry is for registry interface
type Registry interface {
	NewKeygener() wallets.Keygener
}

type registry struct {
	conf        *config.Config
	walletType  wallet.WalletType
	logger      *zap.Logger
	btc         api.Bitcoiner
	rpcClient   *rpcclient.Client
	mysqlClient *sql.DB
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType wallet.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewKeygener is to register for keygener interface
// Which is better ?
// - create each interface getter to difine as interface
// - return struct itself
//func (r *registry) NewKeygener() *keygen.Keygen {
func (r *registry) NewKeygener() wallets.Keygener {
	return wallets.NewKeygen(
		r.newBTC(),
		r.newMySQLClient(),
		r.newSeeder(),
		r.newHdWallter(),
		r.newPrivKeyer(),
		r.newFullPubKeyImporter(),
		r.newMultisiger(),
		r.newAddressExporter(),
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
	return coldsrv.NewAccountHDWalletRepo(
		r.newAccountKeyRepo(),
	)
}

func (r *registry) newPrivKeyer() keygensrv.PrivKeyer {
	return keygensrv.NewPrivKey(
		r.newBTC(),
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.walletType,
	)
}

func (r *registry) newFullPubKeyImporter() keygensrv.FullPubKeyImporter {
	return keygensrv.NewFullPubkeyImport(
		r.newBTC(),
		r.newLogger(),
		r.newAuthFullPubKeyRepo(),
		r.newPubkeyFileStorager(),
		r.walletType,
	)
}

func (r *registry) newMultisiger() keygensrv.Multisiger {
	return keygensrv.NewMultisig(
		r.newBTC(),
		r.newLogger(),
		r.newAuthFullPubKeyRepo(),
		r.newAccountKeyRepo(),
		r.walletType,
	)
}

func (r *registry) newAddressExporter() keygensrv.AddressExporter {
	return keygensrv.NewAddressExport(
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.newAddressFileStorager(),
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

func (r *registry) newAuthFullPubKeyRepo() coldrepo.AuthFullPubkeyRepositorier {
	return coldrepo.NewAuthFullPubkeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAuthKeyRepo() coldrepo.AuthAccountKeyRepositorier {
	return coldrepo.NewAuthAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		"",
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

func (r *registry) newAddressFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.AddressFile.BasePath,
		r.newLogger(),
	)
}

func (r *registry) newPubkeyFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.PubkeyFile.BasePath,
		r.newLogger(),
	)
}

func (r *registry) newTxFileStorager() tx.FileRepositorier {
	return tx.NewFileRepository(
		r.conf.TxFile.BasePath,
		r.newLogger(),
	)
}
