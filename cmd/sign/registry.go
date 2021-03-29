package main

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv/signsrv"
	commonsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/coldsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
)

// Registry is for registry interface
type Registry interface {
	NewSigner() wallets.Signer
}

type registry struct {
	conf        *config.WalletRoot
	accountConf *account.AccountRoot
	walletType  wallet.WalletType
	authType    account.AuthType
	logger      *zap.Logger
	btc         btcgrp.Bitcoiner
	rpcClient   *rpcclient.Client
	mysqlClient *sql.DB
	multisig    account.MultisigAccounter
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.WalletRoot, accountConf *account.AccountRoot, walletType wallet.WalletType, authName string) Registry {
	// validate
	if !account.ValidateAuthType(authName) {
		panic(fmt.Sprintf("authName is invalid. this should be embedded when building: %s", authName))
	}

	return &registry{
		conf:        conf,
		accountConf: accountConf,
		walletType:  walletType,
		authType:    account.AuthTypeMap[authName],
	}
}

// NewSigner is to register for Signer interface
// NewKeygener is to register for keygener interface
func (r *registry) NewSigner() wallets.Signer {
	switch r.conf.CoinTypeCode {
	case coin.BTC, coin.BCH:
		return btcwallet.NewBTCSign(
			r.newBTC(),
			r.newMySQLClient(),
			r.authType,
			r.conf.AddressType,
			r.newSeeder(),
			r.newHdWallter(),
			r.newPrivKeyer(),
			r.newFullPubkeyExporter(),
			r.newSigner(),
			r.walletType,
		)
	case coin.ETH:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	}
}

func (r *registry) newSeeder() service.Seeder {
	return commonsrv.NewSeed(
		r.newLogger(),
		r.newSeedRepo(),
		r.walletType,
	)
}

func (r *registry) newHdWallter() service.HDWalleter {
	return commonsrv.NewHDWallet(
		r.newLogger(),
		r.newHdWalletRepo(),
		r.newKeyGenerator(),
		r.conf.CoinTypeCode,
		r.walletType,
	)
}

func (r *registry) newHdWalletRepo() commonsrv.HDWalletRepo {
	return commonsrv.NewAuthHDWalletRepo(
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

func (r *registry) newFullPubkeyExporter() service.FullPubkeyExporter {
	return signsrv.NewFullPubkeyExport(
		r.newLogger(),
		r.newAuthKeyRepo(),
		r.newPubkeyFileStorager(),
		r.conf.CoinTypeCode,
		r.authType,
		r.walletType,
	)
}

func (r *registry) newSigner() service.Signer {
	return coldsrv.NewSign(
		r.newBTC(),
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.newAuthKeyRepo(),
		r.newTxFileStorager(),
		r.newMultiAccount(),
		r.walletType,
	)
}

func (r *registry) newRPCClient() *rpcclient.Client {
	var err error
	if r.rpcClient == nil {
		r.rpcClient, err = btcgrp.NewRPCClient(&r.conf.Bitcoin)
	}
	if err != nil {
		panic(err)
	}
	return r.rpcClient
}

func (r *registry) newBTC() btcgrp.Bitcoiner {
	if r.btc == nil {
		var err error
		r.btc, err = btcgrp.NewBitcoin(
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

func (r *registry) newMultiAccount() account.MultisigAccounter {
	if r.multisig == nil {
		if r.accountConf.Multisigs == nil {
			panic("account config is required to call newMultiAccount()")
		}
		r.multisig = account.NewMultisigAccounts(r.accountConf.Multisigs)
	}
	return r.multisig
}
