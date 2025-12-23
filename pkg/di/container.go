package di

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/contract"
	"github.com/hiromaily/go-crypto-wallet/pkg/converter"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	wtype "github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/erc20"
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

// Container is for DI container interface
type Container interface {
	NewWalleter() wallets.Watcher
}

type container struct {
	// config
	conf        *config.WalletRoot
	accountConf *account.AccountRoot
	// wallet
	walletType wtype.WalletType
	btc        btcgrp.Bitcoiner
	eth        ethgrp.Ethereumer
	erc20      ethgrp.ERC20er
	xrp        xrpgrp.Rippler
	// client
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *ws.WS
	wsXrpAdmin   *ws.WS
	grpcConn     *grpc.ClientConn
	rippleAPI    *xrp.RippleAPI
	// db
	mysqlClient *sql.DB
	// utility
	logger      logger.Logger
	uuidHandler uuid.UUIDHandler
}

// NewContainer is to create container interface
func NewContainer(conf *config.WalletRoot, accountConf *account.AccountRoot, walletType wtype.WalletType) Container {
	return &container{
		conf:        conf,
		accountConf: accountConf,
		walletType:  walletType,
	}
}

//
// Wallet
//

// NewWalleter is to register for walleter interface
func (c *container) NewWalleter() wallets.Watcher {
	switch {
	case coin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWalleter()
	case coin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWalleter()
	case c.conf.CoinTypeCode == coin.XRP:
		return c.newXRPWalleter()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) newBTCWalleter() wallets.Watcher {
	return btcwallet.NewBTCWatch(
		c.newBTC(),
		c.newMySQLClient(),
		c.newLogger(),
		c.conf.AddressType,
		c.newBTCAddressImporter(),
		c.newBTCTxCreator(),
		c.newBTCTxSender(),
		c.newBTCTxMonitorer(),
		c.newPaymentRequestCreator(),
		c.walletType,
	)
}

func (c *container) newETHWalleter() wallets.Watcher {
	return ethwallet.NewETHWatch(
		c.newETH(),
		c.newMySQLClient(),
		c.newLogger(),
		c.newCommonAddressImporter(),
		c.newETHTxCreator(),
		c.newETHTxSender(),
		c.newETHTxMonitorer(),
		c.newPaymentRequestCreator(),
		c.walletType,
	)
}

func (c *container) newXRPWalleter() wallets.Watcher {
	return xrpwallet.NewXRPWatch(
		c.newXRP(),
		c.newMySQLClient(),
		c.newLogger(),
		c.newCommonAddressImporter(),
		c.newXRPTxCreator(),
		c.newXRPTxSender(),
		c.newXRPTxMonitorer(),
		c.newPaymentRequestCreator(),
		c.walletType,
	)
}

//
// Wallet Service
//

func (c *container) newBTCAddressImporter() service.AddressImporter {
	return btcsrv.NewAddressImport(
		c.newBTC(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
		c.walletType,
	)
}

func (c *container) newCommonAddressImporter() watchsrv.AddressImporter {
	return watchsrv.NewAddressImport(
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
		c.walletType,
	)
}

func (c *container) newBTCTxCreator() service.TxCreator {
	return btcsrv.NewTxCreate(
		c.newBTC(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newBTCTxRepo(),
		c.newBTCTxInputRepo(),
		c.newBTCTxOutputRepo(),
		c.newPaymentRequestRepo(),
		c.newTxFileRepo(),
		c.newDepositAccount(),
		c.newPaymentAccount(),
		c.walletType,
	)
}

func (c *container) newETHTxCreator() ethsrv.TxCreator {
	var targetEthAPI ethgrp.EtherTxCreator
	if coin.IsERC20Token(c.conf.CoinTypeCode.String()) {
		targetEthAPI = c.newERC20()
	} else {
		targetEthAPI = c.newETH()
	}

	return ethsrv.NewTxCreate(
		targetEthAPI,
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newTxRepo(),
		c.newETHTxDetailRepo(),
		c.newPaymentRequestRepo(),
		c.newTxFileRepo(),
		c.newDepositAccount(),
		c.newPaymentAccount(),
		c.walletType,
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPTxCreator() xrpsrv.TxCreator {
	return xrpsrv.NewTxCreate(
		c.newXRP(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newUUIDHandler(),
		c.newAddressRepo(),
		c.newTxRepo(),
		c.newXRPTxDetailRepo(),
		c.newPaymentRequestRepo(),
		c.newTxFileRepo(),
		c.newDepositAccount(),
		c.newPaymentAccount(),
		c.walletType,
	)
}

func (c *container) newBTCTxSender() service.TxSender {
	return btcsrv.NewTxSend(
		c.newBTC(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newBTCTxRepo(),
		c.newBTCTxOutputRepo(),
		c.newTxFileRepo(),
		c.walletType,
	)
}

func (c *container) newETHTxSender() service.TxSender {
	return ethsrv.NewTxSend(
		c.newETH(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newTxRepo(),
		c.newETHTxDetailRepo(),
		c.newTxFileRepo(),
		c.walletType,
	)
}

func (c *container) newXRPTxSender() service.TxSender {
	return xrpsrv.NewTxSend(
		c.newXRP(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newTxRepo(),
		c.newXRPTxDetailRepo(),
		c.newTxFileRepo(),
		c.walletType,
	)
}

func (c *container) newBTCTxMonitorer() service.TxMonitorer {
	return btcsrv.NewTxMonitor(
		c.newBTC(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newBTCTxRepo(),
		c.newBTCTxInputRepo(),
		c.newPaymentRequestRepo(),
		c.walletType,
	)
}

func (c *container) newETHTxMonitorer() service.TxMonitorer {
	if c.conf.Ethereum.ConfirmationNum == 0 {
		panic("confirmation_num of ethereum in config is required")
	}

	return ethsrv.NewTxMonitor(
		c.newETH(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newETHTxDetailRepo(),
		c.conf.Ethereum.ConfirmationNum,
		c.walletType,
	)
}

func (c *container) newXRPTxMonitorer() service.TxMonitorer {
	return xrpsrv.NewTxMonitor(
		c.newXRP(),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newXRPTxDetailRepo(),
		c.walletType,
	)
}

func (c *container) newPaymentRequestCreator() service.PaymentRequestCreator {
	return watchsrv.NewPaymentRequestCreate(
		c.newConverter(c.conf.CoinTypeCode),
		c.newLogger(),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newPaymentRequestRepo(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newConverter(coinTypeCode coin.CoinTypeCode) converter.Converter {
	switch coinTypeCode {
	case coin.BTC:
		return c.newBTC()
	case coin.BCH, coin.LTC, coin.ETH, coin.XRP, coin.ERC20, coin.HYC:
		return converter.NewConverter()
	default:
		return converter.NewConverter()
	}
}

//
// RPC Client
//

func (c *container) newRPCClient() *rpcclient.Client {
	if c.rpcClient == nil {
		var err error
		c.rpcClient, err = btcgrp.NewRPCClient(&c.conf.Bitcoin)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcClient
}

func (c *container) newEthRPCClient() *ethrpc.Client {
	if c.rpcEthClient == nil {
		var err error
		c.rpcEthClient, err = ethgrp.NewRPCClient(&c.conf.Ethereum)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcEthClient
}

func (c *container) newXRPWSClient() (*ws.WS, *ws.WS) {
	if c.wsXrpPublic == nil {
		var err error
		c.wsXrpPublic, c.wsXrpAdmin, err = xrpgrp.NewWSClient(&c.conf.Ripple)
		if err != nil {
			panic(err)
		}
	}
	return c.wsXrpPublic, c.wsXrpAdmin
}

func (c *container) newGRPCConn() *grpc.ClientConn {
	if c.grpcConn == nil {
		var err error
		c.grpcConn, err = xrpgrp.NewGRPCClient(&c.conf.Ripple.API)
		if err != nil {
			panic(err)
		}
	}
	return c.grpcConn
}

//
// Wallet API
//

func (c *container) newBTC() btcgrp.Bitcoiner {
	if c.btc == nil {
		var err error
		c.btc, err = btcgrp.NewBitcoin(
			c.newRPCClient(),
			&c.conf.Bitcoin,
			c.newLogger(),
			c.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return c.btc
}

func (c *container) newETH() ethgrp.Ethereumer {
	if c.eth == nil {
		var err error
		c.eth, err = ethgrp.NewEthereum(
			c.newEthRPCClient(),
			&c.conf.Ethereum,
			c.newLogger(),
			c.conf.CoinTypeCode,
			c.newUUIDHandler(),
		)
		if err != nil {
			panic(err)
		}
	}
	return c.eth
}

func (c *container) newERC20() ethgrp.ERC20er {
	if c.erc20 == nil {
		var err error
		client := ethclient.NewClient(c.newEthRPCClient())
		conf := c.conf.Ethereum
		tokenClient, err := contract.NewContractToken(
			conf.ERC20s[conf.ERC20Token].ContractAddress,
			client,
		)
		if err != nil {
			panic(err)
		}
		c.erc20 = erc20.NewERC20(
			client,
			tokenClient,
			conf.ERC20Token,
			c.newUUIDHandler(),
			conf.ERC20s[conf.ERC20Token].Name,
			conf.ERC20s[conf.ERC20Token].ContractAddress,
			conf.ERC20s[conf.ERC20Token].MasterAddress,
			conf.ERC20s[conf.ERC20Token].Decimals,
			c.newLogger(),
		)
	}
	return c.erc20
}

func (c *container) newXRP() xrpgrp.Rippler {
	if c.xrp == nil {
		var err error
		wsPublic, wsAdmin := c.newXRPWSClient()
		c.xrp, err = xrpgrp.NewRipple(
			wsPublic,
			wsAdmin,
			c.newRippleAPI(),
			&c.conf.Ripple,
			c.newLogger(),
			c.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return c.xrp
}

func (c *container) newRippleAPI() *xrp.RippleAPI {
	if c.rippleAPI == nil {
		c.rippleAPI = xrp.NewRippleAPI(c.newGRPCConn(), c.newLogger())
	}
	return c.rippleAPI
}

//
// Logger
//

func (c *container) newLogger() logger.Logger {
	if c.logger == nil {
		c.logger = logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service)
	}
	return c.logger
}

//
// UUID
//

func (c *container) newUUIDHandler() uuid.UUIDHandler {
	if c.uuidHandler == nil {
		c.uuidHandler = uuid.NewGoogleUUIDHandler()
	}
	return c.uuidHandler
}

//
// DB
//

func (c *container) newMySQLClient() *sql.DB {
	if c.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&c.conf.MySQL)
		if err != nil {
			panic(err)
		}
		c.mysqlClient = dbConn
	}
	return c.mysqlClient
}

//
// Repository
//

func (c *container) newBTCTxRepo() watchrepo.BTCTxRepositorier {
	return watchrepo.NewBTCTxRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newBTCTxInputRepo() watchrepo.TxInputRepositorier {
	return watchrepo.NewBTCTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newBTCTxOutputRepo() watchrepo.TxOutputRepositorier {
	return watchrepo.NewBTCTxOutputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newTxRepo() watchrepo.TxRepositorier {
	return watchrepo.NewTxRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newETHTxDetailRepo() watchrepo.EthDetailTxRepositorier {
	return watchrepo.NewEthDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newXRPTxDetailRepo() watchrepo.XrpDetailTxRepositorier {
	return watchrepo.NewXrpDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newPaymentRequestRepo() watchrepo.PaymentRequestRepositorier {
	return watchrepo.NewPaymentRequestRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newAddressRepo() watchrepo.AddressRepositorier {
	return watchrepo.NewAddressRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newLogger(),
	)
}

func (c *container) newAddressFileRepo() address.FileRepositorier {
	return address.NewFileRepository(
		c.conf.FilePath.FullPubKey,
		c.newLogger(),
	)
}

func (c *container) newTxFileRepo() tx.FileRepositorier {
	return tx.NewFileRepository(
		c.conf.FilePath.Tx,
		c.newLogger(),
	)
}

//
// Account
//

func (c *container) newDepositAccount() account.AccountType {
	if c.accountConf == nil || c.accountConf.DepositReceiver == "" {
		return account.AccountTypeDeposit
	}
	return c.accountConf.DepositReceiver
}

func (c *container) newPaymentAccount() account.AccountType {
	if c.accountConf == nil || c.accountConf.PaymentSender == "" {
		return account.AccountTypePayment
	}
	return c.accountConf.PaymentSender
}
