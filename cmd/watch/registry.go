package main

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/opentracing/opentracing-go"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tracer"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/watchsrv"
	ethsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/watchsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
)

// Registry is for registry interface
type Registry interface {
	NewWalleter() wallets.Watcher
}

type registry struct {
	conf         *config.Config
	walletType   wtype.WalletType
	logger       *zap.Logger
	btc          btcgrp.Bitcoiner
	eth          ethgrp.Ethereumer
	xrp          xrpgrp.Rippler
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	mysqlClient  *sql.DB
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType wtype.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewWalleter is to register for walleter interface
func (r *registry) NewWalleter() wallets.Watcher {
	switch r.conf.CoinTypeCode {
	case coin.BTC, coin.BCH:
		return r.newBTCWalleter()
	case coin.ETH:
		return r.newETHWalleter()
	case coin.XRP:
		return r.newXRPWalleter()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	}
}

func (r *registry) newBTCWalleter() wallets.Watcher {
	return btcwallet.NewBTCWatch(
		r.newBTC(),
		r.newMySQLClient(),
		r.newLogger(),
		r.newTracer(),
		r.conf.AddressType,
		r.newBTCAddressImporter(),
		r.newTxCreator(),
		r.newTxSender(),
		r.newTxMonitorer(),
		r.newPaymentRequestCreator(),
		r.walletType,
	)
}

func (r *registry) newETHWalleter() wallets.Watcher {
	return ethwallet.NewETHWatch(
		r.newETH(),
		r.newMySQLClient(),
		r.newLogger(),
		r.newETHAddressImporter(),
		r.newETHTxCreator(),
		r.newETHTxSender(),
		r.newETHTxMonitorer(),
		r.newETHPaymentRequestCreator(),
		r.walletType,
	)
}

func (r *registry) newXRPWalleter() wallets.Watcher {
	return xrpwallet.NewXRPWatch(
		r.newXRP(),
		r.newMySQLClient(),
		r.newLogger(),
		r.walletType,
	)
}

func (r *registry) newBTCAddressImporter() service.AddressImporter {
	return watchsrv.NewAddressImport(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newAddressFileRepo(),
		r.conf.CoinTypeCode,
		r.conf.AddressType,
		r.walletType,
	)
}

func (r *registry) newETHAddressImporter() ethsrv.AddressImporter {
	return ethsrv.NewAddressImport(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newAddressFileRepo(),
		r.conf.CoinTypeCode,
		r.conf.AddressType,
		r.walletType,
	)
}

func (r *registry) newTxCreator() service.TxCreator {
	return watchsrv.NewTxCreate(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newTxRepo(),
		r.newTxInputRepo(),
		r.newTxOutputRepo(),
		r.newPaymentRequestRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newETHTxCreator() ethsrv.TxCreator {
	return ethsrv.NewTxCreate(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newETHTxRepo(),
		r.newETHTxDetailRepo(),
		r.newPaymentRequestRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newTxSender() service.TxSender {
	return watchsrv.NewTxSend(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newTxRepo(),
		r.newTxOutputRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newETHTxSender() service.TxSender {
	return ethsrv.NewTxSend(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newETHTxRepo(),
		r.newETHTxDetailRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newTxMonitorer() service.TxMonitorer {
	return watchsrv.NewTxMonitor(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newTxRepo(),
		r.newTxInputRepo(),
		r.newPaymentRequestRepo(),
		r.walletType,
	)
}

func (r *registry) newETHTxMonitorer() service.TxMonitorer {
	if r.conf.Ethereum.ConfirmationNum == 0 {
		panic("confirmation_num of ethereum in config is required")
	}

	return ethsrv.NewTxMonitor(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newETHTxDetailRepo(),
		r.conf.Ethereum.ConfirmationNum,
		r.walletType,
	)
}

func (r *registry) newPaymentRequestCreator() service.PaymentRequestCreator {
	return watchsrv.NewPaymentRequestCreate(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newPaymentRequestRepo(),
		r.walletType,
	)
}

func (r *registry) newETHPaymentRequestCreator() service.PaymentRequestCreator {
	return ethsrv.NewPaymentRequestCreate(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newPaymentRequestRepo(),
		r.walletType,
	)
}

func (r *registry) newRPCClient() *rpcclient.Client {
	if r.rpcClient == nil {
		var err error
		r.rpcClient, err = btcgrp.NewRPCClient(&r.conf.Bitcoin)
		if err != nil {
			panic(err)
		}
	}
	return r.rpcClient
}

func (r *registry) newEthRPCClient() *ethrpc.Client {
	if r.rpcEthClient == nil {
		var err error
		r.rpcEthClient, err = ethgrp.NewRPCClient(&r.conf.Ethereum)
		if err != nil {
			panic(err)
		}
	}
	return r.rpcEthClient
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

func (r *registry) newETH() ethgrp.Ethereumer {
	if r.eth == nil {
		var err error
		r.eth, err = ethgrp.NewEthereum(
			r.newEthRPCClient(),
			&r.conf.Ethereum,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.eth
}

func (r *registry) newXRP() xrpgrp.Rippler {
	if r.xrp == nil {
		var err error
		r.xrp, err = xrpgrp.NewRipple(
			&r.conf.Ripple,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.xrp
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

func (r *registry) newTxRepo() watchrepo.BTCTxRepositorier {
	return watchrepo.NewBTCTxRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxInputRepo() watchrepo.TxInputRepositorier {
	return watchrepo.NewTxInputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxOutputRepo() watchrepo.TxOutputRepositorier {
	return watchrepo.NewTxOutputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newETHTxRepo() watchrepo.ETHTxRepositorier {
	return watchrepo.NewETHTxRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newETHTxDetailRepo() watchrepo.EthDetailTxRepositorier {
	return watchrepo.NewEthDetailTxInputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newPaymentRequestRepo() watchrepo.PaymentRequestRepositorier {
	return watchrepo.NewPaymentRequestRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAddressRepo() watchrepo.AddressRepositorier {
	return watchrepo.NewAddressRepository(
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

func (r *registry) newAddressFileRepo() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.FilePath.FullPubKey,
		r.newLogger(),
	)
}

func (r *registry) newTxFileRepo() tx.FileRepositorier {
	return tx.NewFileRepository(
		r.conf.FilePath.Tx,
		r.newLogger(),
	)
}
