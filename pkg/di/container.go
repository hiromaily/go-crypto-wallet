package di

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
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
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
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
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	btccoldsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv"
	btckeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv/keygensrv"
	btcsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv/signsrv"
	btcsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/watchsrv"
	commonsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/coldsrv"
	ethkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/keygensrv"
	ethsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/watchsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watchsrv"
	xrpkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/keygensrv"
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
	NewKeygener() wallets.Keygener
	NewSigner(authName string) wallets.Signer
}

type container struct {
	// config
	conf        *config.WalletRoot
	accountConf *account.AccountRoot
	// db
	mysqlClient *sql.DB
	// utility
	uuidHandler uuid.UUIDHandler
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
	// keygen specific
	multisig account.MultisigAccounter
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

// NewKeygener is to register for keygener interface
func (c *container) NewKeygener() wallets.Keygener {
	// set global logger
	logger.SetGlobal(logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service))

	switch {
	case coin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCKeygener()
	case coin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHKeygener()
	case c.conf.CoinTypeCode == coin.XRP:
		return c.newXRPKeygener()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) newBTCKeygener() wallets.Keygener {
	return btcwallet.NewBTCKeygen(
		c.newBTC(),
		c.newMySQLClient(),
		c.conf.AddressType,
		c.newSeeder(),
		c.newHdWallter(),
		c.newPrivKeyer(),
		c.newFullPubKeyImporter(),
		c.newMultisiger(),
		c.newAddressExporter(),
		c.newSigner(),
		c.walletType,
	)
}

func (c *container) newETHKeygener() wallets.Keygener {
	return ethwallet.NewETHKeygen(
		c.newETH(),
		c.newMySQLClient(),
		c.walletType,
		c.newSeeder(),
		c.newHdWallter(),
		c.newPrivKeyer(),
		c.newAddressExporter(),
		c.newETHSigner(),
	)
}

func (c *container) newXRPKeygener() wallets.Keygener {
	return xrpwallet.NewXRPKeygen(
		c.newXRP(),
		c.newMySQLClient(),
		c.walletType,
		c.newSeeder(),
		c.newHdWallter(),
		c.newXRPKeyGenerator(),
		c.newAddressExporter(),
		c.newXRPSigner(),
	)
}

// NewWalleter is to register for walleter interface
func (c *container) NewWalleter() wallets.Watcher {
	// set global logger
	logger.SetGlobal(logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service))

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

// NewSigner is to register for Signer interface
func (c *container) NewSigner(authName string) wallets.Signer {
	// validate
	if !account.ValidateAuthType(authName) {
		panic("authName is invalid. this should be embedded when building: " + authName)
	}

	// set global logger
	logger.SetGlobal(logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service))

	authType := account.AuthTypeMap[authName]

	switch c.conf.CoinTypeCode {
	case coin.BTC, coin.BCH:
		return c.newBTCSigner(authType)
	case coin.LTC, coin.ETH, coin.XRP, coin.ERC20, coin.HYC:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) newBTCSigner(authType account.AuthType) wallets.Signer {
	return btcwallet.NewBTCSign(
		c.newBTC(),
		c.newMySQLClient(),
		authType,
		c.conf.AddressType,
		c.newSeeder(),
		c.newSignHdWallter(authType),
		c.newSignPrivKeyer(authType),
		c.newSignFullPubkeyExporter(authType),
		c.newSigner(),
		c.walletType,
	)
}

func (c *container) newBTCWalleter() wallets.Watcher {
	return btcwallet.NewBTCWatch(
		c.newBTC(),
		c.newMySQLClient(),
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
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newXRPTxDetailRepo(),
		c.walletType,
	)
}

func (c *container) newPaymentRequestCreator() service.PaymentRequestCreator {
	return watchsrv.NewPaymentRequestCreate(
		c.newConverter(c.conf.CoinTypeCode),
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
		c.rippleAPI = xrp.NewRippleAPI(c.newGRPCConn())
	}
	return c.rippleAPI
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
	)
}

func (c *container) newBTCTxInputRepo() watchrepo.TxInputRepositorier {
	return watchrepo.NewBTCTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCTxOutputRepo() watchrepo.TxOutputRepositorier {
	return watchrepo.NewBTCTxOutputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newTxRepo() watchrepo.TxRepositorier {
	return watchrepo.NewTxRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newETHTxDetailRepo() watchrepo.EthDetailTxRepositorier {
	return watchrepo.NewEthDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPTxDetailRepo() watchrepo.XrpDetailTxRepositorier {
	return watchrepo.NewXrpDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newPaymentRequestRepo() watchrepo.PaymentRequestRepositorier {
	return watchrepo.NewPaymentRequestRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressRepo() watchrepo.AddressRepositorier {
	return watchrepo.NewAddressRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressFileRepo() address.FileRepositorier {
	return address.NewFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (c *container) newTxFileRepo() tx.FileRepositorier {
	return tx.NewFileRepository(
		c.conf.FilePath.Tx,
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

//
// Keygen Service
//

func (c *container) newSeeder() service.Seeder {
	return commonsrv.NewSeed(
		c.newSeedRepo(),
		c.walletType,
	)
}

func (c *container) newHdWallter() service.HDWalleter {
	return commonsrv.NewHDWallet(
		c.newHdWalletRepo(),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newHdWalletRepo() commonsrv.HDWalletRepo {
	return commonsrv.NewAccountHDWalletRepo(
		c.newAccountKeyRepo(),
	)
}

func (c *container) newPrivKeyer() service.PrivKeyer {
	switch {
	case coin.IsBTCGroup(c.conf.CoinTypeCode):
		return btckeygensrv.NewPrivKey(
			c.newBTC(),
			c.newAccountKeyRepo(),
			c.walletType,
		)
	case coin.IsETHGroup(c.conf.CoinTypeCode):
		return ethkeygensrv.NewPrivKey(
			c.newETH(),
			c.newAccountKeyRepo(),
			c.walletType,
		)
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) newFullPubKeyImporter() service.FullPubKeyImporter {
	return btckeygensrv.NewFullPubkeyImport(
		c.newBTC(),
		c.newAuthFullPubKeyRepo(),
		c.newPubkeyFileStorager(),
		c.walletType,
	)
}

func (c *container) newMultisiger() service.Multisiger {
	return btckeygensrv.NewMultisig(
		c.newBTC(),
		c.newAuthFullPubKeyRepo(),
		c.newAccountKeyRepo(),
		c.newMultiAccount(),
		c.walletType,
	)
}

func (c *container) newAddressExporter() service.AddressExporter {
	return commonsrv.NewAddressExport(
		c.newAccountKeyRepo(),
		c.newAddressFileStorager(),
		c.newMultiAccount(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newSigner() service.Signer {
	return btccoldsrv.NewSign(
		c.newBTC(),
		c.newAccountKeyRepo(),
		c.newAuthKeyRepo(),
		c.newTxFileStorager(),
		c.newMultiAccount(),
		c.walletType,
	)
}

func (c *container) newETHSigner() service.Signer {
	return ethkeygensrv.NewSign(
		c.newETH(),
		c.newTxFileStorager(),
		c.walletType,
	)
}

func (c *container) newXRPSigner() service.Signer {
	return xrpkeygensrv.NewSign(
		c.newXRP(),
		c.newXRPAccountKeyRepo(),
		c.newTxFileStorager(),
		c.walletType,
	)
}

func (c *container) newXRPKeyGenerator() xrpkeygensrv.XRPKeyGenerator {
	return xrpkeygensrv.NewXRPKeyGenerate(
		c.newXRP(),
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.walletType,
		c.newAccountKeyRepo(),
		c.newXRPAccountKeyRepo(),
	)
}

func (c *container) newKeyGenerator() key.Generator {
	var chainConf *chaincfg.Params
	switch {
	case coin.IsBTCGroup(c.conf.CoinTypeCode):
		chainConf = c.newBTC().GetChainConf()
	case coin.IsETHGroup(c.conf.CoinTypeCode):
		chainConf = c.newETH().GetChainConf()
	case c.conf.CoinTypeCode == coin.XRP:
		chainConf = c.newXRP().GetChainConf()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}

	return key.NewHDKey(
		key.PurposeTypeBIP44,
		c.conf.CoinTypeCode,
		chainConf)
}

func (c *container) newMultiAccount() account.MultisigAccounter {
	if c.multisig == nil {
		if c.accountConf == nil || c.accountConf.Multisigs == nil {
			return account.NewMultisigAccounts(nil)
		}
		c.multisig = account.NewMultisigAccounts(c.accountConf.Multisigs)
	}
	return c.multisig
}

//
// Keygen Repository
//

func (c *container) newSeedRepo() coldrepo.SeedRepositorier {
	return coldrepo.NewSeedRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAccountKeyRepo() coldrepo.AccountKeyRepositorier {
	return coldrepo.NewAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPAccountKeyRepo() coldrepo.XRPAccountKeyRepositorier {
	return coldrepo.NewXRPAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAuthFullPubKeyRepo() coldrepo.AuthFullPubkeyRepositorier {
	return coldrepo.NewAuthFullPubkeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAuthKeyRepo() coldrepo.AuthAccountKeyRepositorier {
	return coldrepo.NewAuthAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

//
// Keygen File Storage
//

func (c *container) newAddressFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		c.conf.FilePath.Address,
	)
}

func (c *container) newPubkeyFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (c *container) newTxFileStorager() tx.FileRepositorier {
	return tx.NewFileRepository(
		c.conf.FilePath.Tx,
	)
}

//
// Sign Service
//

func (c *container) newSignHdWallter(authType account.AuthType) service.HDWalleter {
	return commonsrv.NewHDWallet(
		c.newSignHdWalletRepo(authType),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newSignHdWalletRepo(authType account.AuthType) commonsrv.HDWalletRepo {
	return commonsrv.NewAuthHDWalletRepo(
		c.newAuthKeyRepo(),
		authType,
	)
}

func (c *container) newSignPrivKeyer(authType account.AuthType) btcsignsrv.PrivKeyer {
	return btcsignsrv.NewPrivKey(
		c.newBTC(),
		c.newAuthKeyRepo(),
		authType,
		c.walletType,
	)
}

func (c *container) newSignFullPubkeyExporter(authType account.AuthType) service.FullPubkeyExporter {
	return btcsignsrv.NewFullPubkeyExport(
		c.newAuthKeyRepo(),
		c.newPubkeyFileStorager(),
		c.conf.CoinTypeCode,
		authType,
		c.walletType,
	)
}
