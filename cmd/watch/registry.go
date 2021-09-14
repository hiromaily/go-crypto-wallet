package main

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/opentracing/opentracing-go"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/converter"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tracer"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	btcsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/watchsrv"
	ethsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/watchsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watchsrv"
	xrpsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/watchsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/ws"
)

// Registry is for registry interface
type Registry interface {
	NewWalleter() wallets.Watcher
}

type registry struct {
	conf         *config.WalletRoot
	accountConf  *account.AccountRoot
	walletType   wtype.WalletType
	logger       *zap.Logger
	btc          btcgrp.Bitcoiner
	eth          ethgrp.Ethereumer
	xrp          xrpgrp.Rippler
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *ws.WS
	wsXrpAdmin   *ws.WS
	grpcConn     *grpc.ClientConn
	rippleAPI    *xrp.RippleAPI
	mysqlClient  *sql.DB
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.WalletRoot, accountConf *account.AccountRoot, walletType wtype.WalletType) Registry {
	return &registry{
		conf:        conf,
		accountConf: accountConf,
		walletType:  walletType,
	}
}

// NewWalleter is to register for walleter interface
func (r *registry) NewWalleter() wallets.Watcher {
	switch r.conf.CoinTypeCode {
	case coin.BTC, coin.BCH:
		return r.newBTCWalleter()
	case coin.ETH, coin.ERC20:
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
		r.newBTCTxCreator(),
		r.newBTCTxSender(),
		r.newBTCTxMonitorer(),
		r.newPaymentRequestCreator(),
		r.walletType,
	)
}

func (r *registry) newETHWalleter() wallets.Watcher {
	return ethwallet.NewETHWatch(
		r.newETH(),
		r.newMySQLClient(),
		r.newLogger(),
		r.newCommonAddressImporter(),
		r.newETHTxCreator(),
		r.newETHTxSender(),
		r.newETHTxMonitorer(),
		r.newPaymentRequestCreator(),
		r.walletType,
	)
}

func (r *registry) newXRPWalleter() wallets.Watcher {
	return xrpwallet.NewXRPWatch(
		r.newXRP(),
		r.newMySQLClient(),
		r.newLogger(),
		r.newCommonAddressImporter(),
		r.newXRPTxCreator(),
		r.newXRPTxSender(),
		r.newXRPTxMonitorer(),
		r.newPaymentRequestCreator(),
		r.walletType,
	)
}

func (r *registry) newBTCAddressImporter() service.AddressImporter {
	return btcsrv.NewAddressImport(
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

func (r *registry) newCommonAddressImporter() watchsrv.AddressImporter {
	return watchsrv.NewAddressImport(
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newAddressFileRepo(),
		r.conf.CoinTypeCode,
		r.conf.AddressType,
		r.walletType,
	)
}

func (r *registry) newBTCTxCreator() service.TxCreator {
	return btcsrv.NewTxCreate(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newBTCTxRepo(),
		r.newBTCTxInputRepo(),
		r.newBTCTxOutputRepo(),
		r.newPaymentRequestRepo(),
		r.newTxFileRepo(),
		r.newDepositAccount(),
		r.newPaymentAccount(),
		r.walletType,
	)
}

func (r *registry) newETHTxCreator() ethsrv.TxCreator {
	return ethsrv.NewTxCreate(
		r.newETH(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newTxRepo(),
		r.newETHTxDetailRepo(),
		r.newPaymentRequestRepo(),
		r.newTxFileRepo(),
		r.newDepositAccount(),
		r.newPaymentAccount(),
		r.walletType,
	)
}

func (r *registry) newXRPTxCreator() xrpsrv.TxCreator {
	return xrpsrv.NewTxCreate(
		r.newXRP(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newTxRepo(),
		r.newXRPTxDetailRepo(),
		r.newPaymentRequestRepo(),
		r.newTxFileRepo(),
		r.newDepositAccount(),
		r.newPaymentAccount(),
		r.walletType,
	)
}

func (r *registry) newBTCTxSender() service.TxSender {
	return btcsrv.NewTxSend(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newBTCTxRepo(),
		r.newBTCTxOutputRepo(),
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
		r.newTxRepo(),
		r.newETHTxDetailRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newXRPTxSender() service.TxSender {
	return xrpsrv.NewTxSend(
		r.newXRP(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newTxRepo(),
		r.newXRPTxDetailRepo(),
		r.newTxFileRepo(),
		r.walletType,
	)
}

func (r *registry) newBTCTxMonitorer() service.TxMonitorer {
	return btcsrv.NewTxMonitor(
		r.newBTC(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newBTCTxRepo(),
		r.newBTCTxInputRepo(),
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

func (r *registry) newXRPTxMonitorer() service.TxMonitorer {
	return xrpsrv.NewTxMonitor(
		r.newXRP(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newXRPTxDetailRepo(),
		r.walletType,
	)
}

func (r *registry) newPaymentRequestCreator() service.PaymentRequestCreator {
	return watchsrv.NewPaymentRequestCreate(
		r.newConverter(r.conf.CoinTypeCode),
		r.newLogger(),
		r.newMySQLClient(),
		r.newAddressRepo(),
		r.newPaymentRequestRepo(),
		r.conf.CoinTypeCode,
		r.walletType,
	)
}

func (r *registry) newConverter(coinTypeCode coin.CoinTypeCode) converter.Converter {
	switch coinTypeCode {
	case coin.BTC:
		return r.newBTC()
	default:
		return converter.NewConverter()
	}
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

func (r *registry) newXRPWSClient() (*ws.WS, *ws.WS) {
	if r.wsXrpPublic == nil {
		var err error
		r.wsXrpPublic, r.wsXrpAdmin, err = xrpgrp.NewWSClient(&r.conf.Ripple)
		if err != nil {
			panic(err)
		}
	}
	return r.wsXrpPublic, r.wsXrpAdmin
}

func (r *registry) newRippleAPI() *xrp.RippleAPI {
	if r.rippleAPI == nil {
		r.rippleAPI = xrp.NewRippleAPI(r.newGRPCConn(), r.newLogger())
	}
	return r.rippleAPI
}

func (r *registry) newGRPCConn() *grpc.ClientConn {
	if r.grpcConn == nil {
		var err error
		r.grpcConn, err = xrpgrp.NewGRPCClient(&r.conf.Ripple.API)
		if err != nil {
			panic(err)
		}
	}
	return r.grpcConn
}

func (r *registry) newXRP() xrpgrp.Rippler {
	if r.xrp == nil {
		var err error
		wsPublic, wsAdmin := r.newXRPWSClient()
		r.xrp, err = xrpgrp.NewRipple(
			wsPublic,
			wsAdmin,
			r.newRippleAPI(),
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

func (r *registry) newBTCTxRepo() watchrepo.BTCTxRepositorier {
	return watchrepo.NewBTCTxRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newBTCTxInputRepo() watchrepo.TxInputRepositorier {
	return watchrepo.NewBTCTxInputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newBTCTxOutputRepo() watchrepo.TxOutputRepositorier {
	return watchrepo.NewBTCTxOutputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxRepo() watchrepo.TxRepositorier {
	return watchrepo.NewTxRepository(
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

func (r *registry) newXRPTxDetailRepo() watchrepo.XrpDetailTxRepositorier {
	return watchrepo.NewXrpDetailTxInputRepository(
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

func (r *registry) newDepositAccount() account.AccountType {
	if r.accountConf == nil || r.accountConf.DepositReceiver == "" {
		return account.AccountTypeDeposit
	}
	return r.accountConf.DepositReceiver
}

func (r *registry) newPaymentAccount() account.AccountType {
	if r.accountConf == nil || r.accountConf.PaymentSender == "" {
		return account.AccountTypePayment
	}
	return r.accountConf.PaymentSender
}
