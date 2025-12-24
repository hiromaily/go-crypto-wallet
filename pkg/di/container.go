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
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/erc20"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/database/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/network/websocket"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	btckeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/btc"
	ethkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/eth"
	keygenshared "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
	xrpkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/xrp"
	btcsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/btc"
	ethsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/eth"
	xrpsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/xrp"
	watchshared "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/shared"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"

	// Use case imports
	keygenusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen/btc"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen/eth"
	keygenusecaseshared "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen/shared"
	keygenusecasexrp "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen/xrp"
	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	signusecasebtc "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign/btc"
	signusecaseeth "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign/eth"
	signusecaseshared "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign/shared"
	signusecasexrp "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch/btc"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch/eth"
	watchusecaseshared "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch/shared"
	watchusecasexrp "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch/xrp"
)

// Container is for DI container interface
type Container interface {
	NewWalleter() wallets.Watcher
	NewKeygener() wallets.Keygener
	NewSigner(authName string) wallets.Signer

	// Watch Use Cases
	NewWatchCreateTransactionUseCase() any
	NewWatchMonitorTransactionUseCase() any
	NewWatchSendTransactionUseCase() any
	NewWatchImportAddressUseCase() watchusecase.ImportAddressUseCase
	NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase

	// Keygen Use Cases
	NewKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase
	NewKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase
	NewKeygenExportAddressUseCase() keygenusecase.ExportAddressUseCase
	NewKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase
	NewKeygenCreateMultisigAddressUseCase() keygenusecase.CreateMultisigAddressUseCase
	NewKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase
	NewKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase
	NewKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase

	// Sign Use Cases
	NewSignTransactionUseCase() signusecase.SignTransactionUseCase
	NewSignImportPrivateKeyUseCase(authType domainAccount.AuthType) signusecase.ImportPrivateKeyUseCase
	NewSignExportFullPubkeyUseCase(authType domainAccount.AuthType) signusecase.ExportFullPubkeyUseCase
	NewSignGenerateSeedUseCase() signusecase.GenerateSeedUseCase
	NewSignStoreSeedUseCase() signusecase.StoreSeedUseCase
	NewSignGenerateAuthKeyUseCase() signusecase.GenerateAuthKeyUseCase

	// Auth accessors
	AuthName() string
	AuthType() domainAccount.AuthType

	// Config accessors
	AddressType() address.AddrType
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
	walletType domainWallet.WalletType
	btc        bitcoin.Bitcoiner
	eth        ethereum.Ethereumer
	erc20      ethereum.ERC20er
	xrp        ripple.Rippler
	// client
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *websocket.WS
	wsXrpAdmin   *websocket.WS
	grpcConn     *grpc.ClientConn
	rippleAPI    *xrp.RippleAPI
	// keygen specific
	multisig account.MultisigAccounter
	// sign specific
	authName string
}

// NewContainer is to create container interface
func NewContainer(
	conf *config.WalletRoot,
	accountConf *account.AccountRoot,
	walletType domainWallet.WalletType,
) Container {
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
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCKeygener()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHKeygener()
	case c.conf.CoinTypeCode == domainCoin.XRP:
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
		c.newKeygenGenerateSeedUseCase(),
		c.newKeygenGenerateHDWalletUseCase(),
		c.newBTCKeygenImportPrivateKeyUseCase(),
		c.newBTCKeygenImportFullPubkeyUseCase(),
		c.newBTCKeygenCreateMultisigAddressUseCase(),
		c.newKeygenExportAddressUseCase(),
		c.newBTCKeygenSignTransactionUseCase(),
		c.walletType,
	)
}

func (c *container) newETHKeygener() wallets.Keygener {
	return ethwallet.NewETHKeygen(
		c.newETH(),
		c.newMySQLClient(),
		c.walletType,
		c.newKeygenGenerateSeedUseCase(),
		c.newKeygenGenerateHDWalletUseCase(),
		c.newETHKeygenImportPrivateKeyUseCase(),
		c.newKeygenExportAddressUseCase(),
		c.newETHKeygenSignTransactionUseCase(),
	)
}

func (c *container) newXRPKeygener() wallets.Keygener {
	return xrpwallet.NewXRPKeygen(
		c.newXRP(),
		c.newMySQLClient(),
		c.walletType,
		c.newKeygenGenerateSeedUseCase(),
		c.newKeygenGenerateHDWalletUseCase(),
		c.newXRPKeygenGenerateKeyUseCase(),
		c.newKeygenExportAddressUseCase(),
		c.newXRPKeygenSignTransactionUseCase(),
	)
}

// NewWalleter is to register for walleter interface
func (c *container) NewWalleter() wallets.Watcher {
	// set global logger
	logger.SetGlobal(logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service))

	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWalleter()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWalleter()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPWalleter()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

// NewSigner is to register for Signer interface
func (c *container) NewSigner(authName string) wallets.Signer {
	// validate
	if !domainAccount.ValidateAuthType(authName) {
		panic("authName is invalid. this should be embedded when building: " + authName)
	}

	// store authName for accessor methods
	c.authName = authName

	// set global logger
	logger.SetGlobal(logger.NewSlogFromConfig(c.conf.Logger.Env, c.conf.Logger.Level, c.conf.Logger.Service))

	authType := domainAccount.AuthTypeMap[authName]

	switch c.conf.CoinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		return c.newBTCSigner(authType)
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

// AuthName returns the authentication account name for sign wallet
func (c *container) AuthName() string {
	return c.authName
}

// AuthType returns the authentication account type for sign wallet
func (c *container) AuthType() domainAccount.AuthType {
	if c.authName == "" {
		return domainAccount.AuthType("")
	}
	return domainAccount.AuthTypeMap[c.authName]
}

// AddressType returns the address type from configuration
func (c *container) AddressType() address.AddrType {
	return c.conf.AddressType
}

func (c *container) newBTCSigner(authType domainAccount.AuthType) wallets.Signer {
	return btcwallet.NewBTCSign(
		c.newBTC(),
		c.newMySQLClient(),
		authType,
		c.conf.AddressType,
		c.NewSignGenerateSeedUseCase(),
		c.NewSignStoreSeedUseCase(),
		c.NewSignGenerateAuthKeyUseCase(),
		c.newBTCSignImportPrivateKeyUseCase(authType),
		c.newBTCSignExportFullPubkeyUseCase(authType),
		c.newBTCSignTransactionUseCase(),
		c.walletType,
	)
}

func (c *container) newBTCWalleter() wallets.Watcher {
	return btcwallet.NewBTCWatch(
		c.newBTC(),
		c.newMySQLClient(),
		c.conf.AddressType,
		c.newBTCWatchCreateTransactionUseCase(),
		c.newBTCWatchMonitorTransactionUseCase(),
		c.newBTCWatchSendTransactionUseCase(),
		c.newBTCWatchImportAddressUseCase(),
		c.newWatchCreatePaymentRequestUseCase(),
		c.walletType,
	)
}

func (c *container) newETHWalleter() wallets.Watcher {
	return ethwallet.NewETHWatch(
		c.newETH(),
		c.newMySQLClient(),
		c.newETHWatchCreateTransactionUseCase(),
		c.newETHWatchMonitorTransactionUseCase(),
		c.newETHWatchSendTransactionUseCase(),
		c.newWatchImportAddressUseCase(),
		c.newWatchCreatePaymentRequestUseCase(),
		c.walletType,
	)
}

func (c *container) newXRPWalleter() wallets.Watcher {
	return xrpwallet.NewXRPWatch(
		c.newXRP(),
		c.newMySQLClient(),
		c.newXRPWatchCreateTransactionUseCase(),
		c.newXRPWatchMonitorTransactionUseCase(),
		c.newXRPWatchSendTransactionUseCase(),
		c.newWatchImportAddressUseCase(),
		c.newWatchCreatePaymentRequestUseCase(),
		c.walletType,
	)
}

//
// Wallet Service
//

func (c *container) newCommonAddressImporter() watchshared.AddressImporter {
	return watchshared.NewAddressImport(
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
		c.walletType,
	)
}

func (c *container) newPaymentRequestCreator() service.PaymentRequestCreator {
	return watchshared.NewPaymentRequestCreate(
		c.newConverter(c.conf.CoinTypeCode),
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newPaymentRequestRepo(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newConverter(coinTypeCode domainCoin.CoinTypeCode) converter.Converter {
	switch coinTypeCode {
	case domainCoin.BTC:
		return c.newBTC()
	case domainCoin.BCH, domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
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
		c.rpcClient, err = bitcoin.NewRPCClient(&c.conf.Bitcoin)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcClient
}

func (c *container) newEthRPCClient() *ethrpc.Client {
	if c.rpcEthClient == nil {
		var err error
		c.rpcEthClient, err = ethereum.NewRPCClient(&c.conf.Ethereum)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcEthClient
}

func (c *container) newXRPWSClient() (*websocket.WS, *websocket.WS) {
	if c.wsXrpPublic == nil {
		var err error
		c.wsXrpPublic, c.wsXrpAdmin, err = ripple.NewWSClient(&c.conf.Ripple)
		if err != nil {
			panic(err)
		}
	}
	return c.wsXrpPublic, c.wsXrpAdmin
}

func (c *container) newGRPCConn() *grpc.ClientConn {
	if c.grpcConn == nil {
		var err error
		c.grpcConn, err = ripple.NewGRPCClient(&c.conf.Ripple.API)
		if err != nil {
			panic(err)
		}
	}
	return c.grpcConn
}

//
// Wallet API
//

func (c *container) newBTC() bitcoin.Bitcoiner {
	if c.btc == nil {
		var err error
		c.btc, err = bitcoin.NewBitcoin(
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

func (c *container) newETH() ethereum.Ethereumer {
	if c.eth == nil {
		var err error
		c.eth, err = ethereum.NewEthereum(
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

func (c *container) newERC20() ethereum.ERC20er {
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

func (c *container) newXRP() ripple.Rippler {
	if c.xrp == nil {
		var err error
		wsPublic, wsAdmin := c.newXRPWSClient()
		c.xrp, err = ripple.NewRipple(
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

func (c *container) newBTCTxRepo() watch.BTCTxRepositorier {
	return watch.NewBTCTxRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCTxInputRepo() watch.TxInputRepositorier {
	return watch.NewBTCTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCTxOutputRepo() watch.TxOutputRepositorier {
	return watch.NewBTCTxOutputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newTxRepo() watch.TxRepositorier {
	return watch.NewTxRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newETHTxDetailRepo() watch.EthDetailTxRepositorier {
	return watch.NewEthDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPTxDetailRepo() watch.XrpDetailTxRepositorier {
	return watch.NewXrpDetailTxInputRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newPaymentRequestRepo() watch.PaymentRequestRepositorier {
	return watch.NewPaymentRequestRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressRepo() watch.AddressRepositorier {
	return watch.NewAddressRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressFileRepo() file.AddressFileRepositorier {
	return file.NewAddressFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (c *container) newTxFileRepo() file.TransactionFileRepositorier {
	return file.NewTransactionFileRepository(
		c.conf.FilePath.Tx,
	)
}

//
// Account
//

func (c *container) newDepositAccount() domainAccount.AccountType {
	if c.accountConf == nil || c.accountConf.DepositReceiver == "" {
		return domainAccount.AccountTypeDeposit
	}
	return c.accountConf.DepositReceiver
}

func (c *container) newPaymentAccount() domainAccount.AccountType {
	if c.accountConf == nil || c.accountConf.PaymentSender == "" {
		return domainAccount.AccountTypePayment
	}
	return c.accountConf.PaymentSender
}

//
// Keygen Service
//

func (c *container) newSeeder() service.Seeder {
	return keygenshared.NewSeed(
		c.newSeedRepo(),
		c.walletType,
	)
}

func (c *container) newHdWallter() service.HDWalleter {
	return keygenshared.NewHDWallet(
		c.newHdWalletRepo(),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newHdWalletRepo() keygenshared.HDWalletRepo {
	return keygenshared.NewAccountHDWalletRepo(
		c.newAccountKeyRepo(),
	)
}

func (c *container) newPrivKeyer() service.PrivKeyer {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return btckeygensrv.NewPrivKey(
			c.newBTC(),
			c.newAccountKeyRepo(),
			c.walletType,
		)
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
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
	return keygenshared.NewAddressExport(
		c.newAccountKeyRepo(),
		c.newAddressFileStorager(),
		c.newMultiAccount(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newSigner() service.Signer {
	return btcsignsrv.NewSign(
		c.newBTC(),
		c.newAccountKeyRepo(),
		c.newAuthKeyRepo(),
		c.newTxFileStorager(),
		c.newMultiAccount(),
		c.walletType,
	)
}

func (c *container) newETHSigner() service.Signer {
	return ethsignsrv.NewSign(
		c.newETH(),
		c.newTxFileStorager(),
		c.walletType,
	)
}

func (c *container) newXRPSigner() service.Signer {
	return xrpsignsrv.NewSign(
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
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		chainConf = c.newBTC().GetChainConf()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		chainConf = c.newETH().GetChainConf()
	case c.conf.CoinTypeCode == domainCoin.XRP:
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

func (c *container) newSeedRepo() cold.SeedRepositorier {
	return cold.NewSeedRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAccountKeyRepo() cold.AccountKeyRepositorier {
	return cold.NewAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPAccountKeyRepo() cold.XRPAccountKeyRepositorier {
	return cold.NewXRPAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAuthFullPubKeyRepo() cold.AuthFullPubkeyRepositorier {
	return cold.NewAuthFullPubkeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAuthKeyRepo() cold.AuthAccountKeyRepositorier {
	return cold.NewAuthAccountKeyRepositorySqlc(
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

//
// Keygen File Storage
//

func (c *container) newAddressFileStorager() file.AddressFileRepositorier {
	return file.NewAddressFileRepository(
		c.conf.FilePath.Address,
	)
}

func (c *container) newPubkeyFileStorager() file.AddressFileRepositorier {
	return file.NewAddressFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (c *container) newTxFileStorager() file.TransactionFileRepositorier {
	return file.NewTransactionFileRepository(
		c.conf.FilePath.Tx,
	)
}

//
// Sign Service
//

func (c *container) newSignHdWallter(authType domainAccount.AuthType) service.HDWalleter {
	return keygenshared.NewHDWallet(
		c.newSignHdWalletRepo(authType),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
		c.walletType,
	)
}

func (c *container) newSignHdWalletRepo(authType domainAccount.AuthType) keygenshared.HDWalletRepo {
	return keygenshared.NewAuthHDWalletRepo(
		c.newAuthKeyRepo(),
		authType,
	)
}

func (c *container) newSignPrivKeyer(authType domainAccount.AuthType) btcsignsrv.PrivKeyer {
	return btcsignsrv.NewPrivKey(
		c.newBTC(),
		c.newAuthKeyRepo(),
		authType,
		c.walletType,
	)
}

func (c *container) newSignFullPubkeyExporter(authType domainAccount.AuthType) service.FullPubkeyExporter {
	return btcsignsrv.NewFullPubkeyExport(
		c.newAuthKeyRepo(),
		c.newPubkeyFileStorager(),
		c.conf.CoinTypeCode,
		authType,
		c.walletType,
	)
}

//
// Use Case Factory Methods
//

// Watch Use Cases

func (c *container) NewWatchCreateTransactionUseCase() any {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWatchCreateTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWatchCreateTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPWatchCreateTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewWatchMonitorTransactionUseCase() any {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWatchMonitorTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWatchMonitorTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPWatchMonitorTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewWatchSendTransactionUseCase() any {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWatchSendTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWatchSendTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPWatchSendTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewWatchImportAddressUseCase() watchusecase.ImportAddressUseCase {
	return c.newWatchImportAddressUseCase()
}

func (c *container) NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase {
	return c.newWatchCreatePaymentRequestUseCase()
}

// Keygen Use Cases

func (c *container) NewKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase {
	return c.newKeygenGenerateHDWalletUseCase()
}

func (c *container) NewKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase {
	return c.newKeygenGenerateSeedUseCase()
}

func (c *container) NewKeygenExportAddressUseCase() keygenusecase.ExportAddressUseCase {
	return c.newKeygenExportAddressUseCase()
}

func (c *container) NewKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCKeygenImportPrivateKeyUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHKeygenImportPrivateKeyUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewKeygenCreateMultisigAddressUseCase() keygenusecase.CreateMultisigAddressUseCase {
	return c.newBTCKeygenCreateMultisigAddressUseCase()
}

func (c *container) NewKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase {
	return c.newBTCKeygenImportFullPubkeyUseCase()
}

func (c *container) NewKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase {
	return c.newXRPKeygenGenerateKeyUseCase()
}

func (c *container) NewKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCKeygenSignTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHKeygenSignTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPKeygenSignTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

// Sign Use Cases

func (c *container) NewSignTransactionUseCase() signusecase.SignTransactionUseCase {
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCSignTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHSignTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPSignTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewSignImportPrivateKeyUseCase(
	authType domainAccount.AuthType,
) signusecase.ImportPrivateKeyUseCase {
	return c.newBTCSignImportPrivateKeyUseCase(authType)
}

func (c *container) NewSignExportFullPubkeyUseCase(
	authType domainAccount.AuthType,
) signusecase.ExportFullPubkeyUseCase {
	return c.newBTCSignExportFullPubkeyUseCase(authType)
}

func (c *container) NewSignGenerateSeedUseCase() signusecase.GenerateSeedUseCase {
	return signusecaseshared.NewGenerateSeedUseCase(c.newSeedRepo())
}

func (c *container) NewSignStoreSeedUseCase() signusecase.StoreSeedUseCase {
	return signusecaseshared.NewStoreSeedUseCase(c.newSeedRepo())
}

func (c *container) NewSignGenerateAuthKeyUseCase() signusecase.GenerateAuthKeyUseCase {
	authType := c.AuthType()
	return signusecaseshared.NewGenerateAuthKeyUseCase(
		c.newSignHdWalletRepo(authType),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
	)
}

// BTC Watch Use Cases

func (c *container) newBTCWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	return watchusecasebtc.NewCreateTransactionUseCase(
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

func (c *container) newBTCWatchMonitorTransactionUseCase() watchusecase.MonitorTransactionUseCase {
	return watchusecasebtc.NewMonitorTransactionUseCase(
		c.newBTC(),
		c.newMySQLClient(),
		c.newBTCTxRepo(),
		c.newBTCTxInputRepo(),
		c.newPaymentRequestRepo(),
	)
}

func (c *container) newBTCWatchSendTransactionUseCase() watchusecase.SendTransactionUseCase {
	return watchusecasebtc.NewSendTransactionUseCase(
		c.newBTC(),
		c.newAddressRepo(),
		c.newBTCTxRepo(),
		c.newBTCTxOutputRepo(),
		c.newTxFileRepo(),
	)
}

func (c *container) newBTCWatchImportAddressUseCase() watchusecase.ImportAddressUseCase {
	return watchusecasebtc.NewImportAddressUseCase(
		c.newBTC(),
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
	)
}

// ETH Watch Use Cases

func (c *container) newETHWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	// Determine which Ethereum API to use based on coin type
	var targetEthAPI ethereum.EtherTxCreator
	if domainCoin.IsERC20Token(c.conf.CoinTypeCode.String()) {
		targetEthAPI = c.newERC20()
	} else {
		targetEthAPI = c.newETH()
	}

	return watchusecaseeth.NewCreateTransactionUseCase(
		targetEthAPI,
		c.newMySQLClient(),
		c.newAddressRepo(),
		c.newTxRepo(),
		c.newETHTxDetailRepo(),
		c.newPaymentRequestRepo(),
		c.newTxFileRepo(),
		c.newDepositAccount(),
		c.newPaymentAccount(),
	)
}

func (c *container) newETHWatchMonitorTransactionUseCase() watchusecase.MonitorTransactionUseCase {
	if c.conf.Ethereum.ConfirmationNum == 0 {
		panic("confirmation_num of ethereum in config is required")
	}

	return watchusecaseeth.NewMonitorTransactionUseCase(
		c.newETH(),
		c.newAddressRepo(),
		c.newETHTxDetailRepo(),
		c.conf.Ethereum.ConfirmationNum,
	)
}

func (c *container) newETHWatchSendTransactionUseCase() watchusecase.SendTransactionUseCase {
	return watchusecaseeth.NewSendTransactionUseCase(
		c.newETH(),
		c.newETHTxDetailRepo(),
		c.newTxFileRepo(),
	)
}

// XRP Watch Use Cases

func (c *container) newXRPWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	return watchusecasexrp.NewCreateTransactionUseCase(
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
	)
}

func (c *container) newXRPWatchMonitorTransactionUseCase() watchusecase.MonitorTransactionUseCase {
	return watchusecasexrp.NewMonitorTransactionUseCase(
		c.newXRP(),
		c.newAddressRepo(),
	)
}

func (c *container) newXRPWatchSendTransactionUseCase() watchusecase.SendTransactionUseCase {
	return watchusecasexrp.NewSendTransactionUseCase(
		c.newXRP(),
		c.newXRPTxDetailRepo(),
		c.newTxFileRepo(),
	)
}

// Shared Watch Use Cases

func (c *container) newWatchImportAddressUseCase() watchusecase.ImportAddressUseCase {
	return watchusecaseshared.NewImportAddressUseCase(
		c.newCommonAddressImporter().(*watchshared.AddressImport),
	)
}

func (c *container) newWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase {
	return watchusecaseshared.NewCreatePaymentRequestUseCase(
		c.newPaymentRequestCreator().(*watchshared.PaymentRequestCreate),
	)
}

// Keygen Use Cases

func (c *container) newKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase {
	return keygenusecaseshared.NewGenerateHDWalletUseCase(
		c.newHdWalletRepo(),
		c.newKeyGenerator(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase {
	return keygenusecaseshared.NewGenerateSeedUseCase(
		c.newSeedRepo(),
	)
}

func (c *container) newKeygenExportAddressUseCase() keygenusecase.ExportAddressUseCase {
	return keygenusecaseshared.NewExportAddressUseCase(
		c.newAccountKeyRepo(),
		c.newAddressFileRepo(),
		c.newMultiAccount(),
		c.conf.CoinTypeCode,
	)
}

// BTC Keygen Use Cases

func (c *container) newBTCKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase {
	return keygenusecasebtc.NewImportPrivateKeyUseCase(
		c.newBTC(),
		c.newAccountKeyRepo(),
	)
}

func (c *container) newBTCKeygenCreateMultisigAddressUseCase() keygenusecase.CreateMultisigAddressUseCase {
	return keygenusecasebtc.NewCreateMultisigAddressUseCase(
		c.newBTC(),
		c.newAuthFullPubKeyRepo(),
		c.newAccountKeyRepo(),
		c.newMultiAccount(),
	)
}

func (c *container) newBTCKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase {
	return keygenusecasebtc.NewImportFullPubkeyUseCase(
		c.newBTC(),
		c.newAuthFullPubKeyRepo(),
		c.newAddressFileRepo(),
	)
}

// ETH Keygen Use Cases

func (c *container) newETHKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase {
	return keygenusecaseeth.NewImportPrivateKeyUseCase(
		c.newETH(),
		c.newAccountKeyRepo(),
	)
}

// XRP Keygen Use Cases

func (c *container) newXRPKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase {
	return keygenusecasexrp.NewGenerateKeyUseCase(
		c.newXRP(),
		c.newMySQLClient(),
		c.conf.CoinTypeCode,
		c.newAccountKeyRepo(),
		c.newXRPAccountKeyRepo(),
	)
}

// Keygen Sign Transaction Use Cases

func (c *container) newBTCKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	return keygenusecasebtc.NewSignTransactionUseCase(
		c.newBTC(),
		c.newAccountKeyRepo(),
		c.newTxFileRepo(),
		c.newMultiAccount(),
	)
}

func (c *container) newETHKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	return keygenusecaseeth.NewSignTransactionUseCase(
		c.newETH(),
		c.newTxFileRepo(),
	)
}

func (c *container) newXRPKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	return keygenusecasexrp.NewSignTransactionUseCase(
		c.newXRP(),
		c.newXRPAccountKeyRepo(),
		c.newTxFileRepo(),
	)
}

// Sign Use Cases

// BTC Sign Use Cases

func (c *container) newBTCSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecasebtc.NewSignTransactionUseCase(
		c.newSigner().(*btcsignsrv.Sign),
	)
}

func (c *container) newBTCSignImportPrivateKeyUseCase(
	authType domainAccount.AuthType,
) signusecase.ImportPrivateKeyUseCase {
	return signusecasebtc.NewImportPrivateKeyUseCase(
		c.newBTC(),
		c.newAuthKeyRepo(),
		authType,
		c.walletType,
	)
}

func (c *container) newBTCSignExportFullPubkeyUseCase(
	authType domainAccount.AuthType,
) signusecase.ExportFullPubkeyUseCase {
	return signusecasebtc.NewExportFullPubkeyUseCase(
		c.newAuthKeyRepo(),
		c.newPubkeyFileStorager(),
		c.conf.CoinTypeCode,
		authType,
		c.walletType,
	)
}

// ETH Sign Use Cases

func (c *container) newETHSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecaseeth.NewSignTransactionUseCase(
		c.newETHSigner().(*ethsignsrv.Sign),
	)
}

// XRP Sign Use Cases

func (c *container) newXRPSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecasexrp.NewSignTransactionUseCase(
		c.newXRPSigner().(*xrpsignsrv.Sign),
	)
}
