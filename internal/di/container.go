// Package di provides dependency injection for application-specific components.
//
// This package is responsible for creating and wiring instances from internal/
// directory (domain, application, infrastructure layers).
//
// For reusable infrastructure components (logger, config, etc.), see pkg/di.
package di

import (
	"context"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	bitcoin "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
	ethimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth"
	apierc20impl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/erc20"
	rippleimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp"
	apixrpimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/contract"
	coldmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mysql"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
	watchmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mysql"
	watchsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/sqlite"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/descriptor"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	infraKeyGen "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key/generator"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
	xrpwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency"
	pkgdi "github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"

	// Use case imports
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebch "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/bch"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	keygenusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/shared"
	keygenusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/xrp"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	signusecasebch "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/bch"
	signusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/btc"
	signusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/eth"
	signusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/shared"
	signusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebch "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/bch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	watchusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/shared"
	watchusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
)

// Container is for DI container interface
type Container interface {
	NewWalleter() wallets.Watcher
	NewKeygener() wallets.Keygener
	NewSigner(authName string) (wallets.Signer, error)

	// Watch Use Cases
	NewWatchCreateTransactionUseCase() any
	NewWatchMonitorTransactionUseCase() any
	NewWatchSendTransactionUseCase() any
	NewWatchImportAddressUseCase() watchusecase.ImportAddressUseCase
	NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase
	NewWatchImportDescriptorUseCase() watchusecase.ImportDescriptorUseCase
	NewWatchAggregateMuSig2SignaturesUseCase() watchusecase.AggregateMuSig2SignaturesUseCase

	// Keygen Use Cases
	NewKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase
	NewKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase
	NewKeygenGenerateDescriptorUseCase() keygenusecase.GenerateDescriptorUseCase
	NewKeygenExportDescriptorUseCase() keygenusecase.ExportDescriptorUseCase
	NewKeygenExportAddressUseCase() keygenusecase.ExportAddressUseCase
	NewKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase
	NewKeygenCreateMultisigAddressUseCase() keygenusecase.CreateMultisigAddressUseCase
	NewKeygenCreateMuSig2AddressUseCase() keygenusecase.CreateMuSig2AddressUseCase
	NewKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase
	NewKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase
	NewKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase
	NewKeygenGenerateMuSig2NonceUseCase() keygenusecase.GenerateMuSig2NonceUseCase
	NewKeygenMuSig2SignUseCase() keygenusecase.MuSig2SignUseCase

	// Sign Use Cases
	NewSignTransactionUseCase() signusecase.SignTransactionUseCase
	NewSignImportPrivateKeyUseCase(authType domainAccount.AuthType) signusecase.ImportPrivateKeyUseCase
	NewSignExportFullPubkeyUseCase(authType domainAccount.AuthType) signusecase.ExportFullPubkeyUseCase
	NewSignGenerateSeedUseCase() signusecase.GenerateSeedUseCase
	NewSignStoreSeedUseCase() signusecase.StoreSeedUseCase
	NewSignGenerateAuthKeyUseCase() signusecase.GenerateAuthKeyUseCase
	NewSignGenerateMuSig2NonceUseCase() signusecase.GenerateMuSig2NonceUseCase
	NewSignMuSig2SignUseCase() signusecase.MuSig2SignUseCase

	// Auth accessors
	AuthName() string
	AuthType() domainAccount.AuthType

	// Config accessors
	AddressType() domainAddress.AddrType
}

var _ Container = (*container)(nil)

type container struct {
	// pkg DI
	pkgContainer pkgdi.PkgContainer
	// config
	conf        *config.WalletRoot
	accountConf *config.AccountRoot
	// wallet
	walletType domainWallet.WalletType
	btc        apibtc.Bitcoiner
	eth        apieth.Ethereumer
	erc20      apieth.ERC20er
	xrp        apixrp.Rippler
	// client
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *websocket.WS
	wsXrpAdmin   *websocket.WS
	rippleAPI    *apixrpimpl.RippleAPI
	// keygen specific
	multisig *domainAccount.MultisigConfig
	// sign specific
	authName string
}

// NewContainer is to create container interface
func NewContainer(
	conf *config.WalletRoot,
	accountConf *config.AccountRoot,
	walletType domainWallet.WalletType,
) Container {
	return &container{
		pkgContainer: pkgdi.NewPkgContainer(conf, accountConf),
		conf:         conf,
		accountConf:  accountConf,
		walletType:   walletType,
	}
}

//
// Wallet
//

// NewKeygener is to register for keygener interface
func (c *container) NewKeygener() wallets.Keygener {
	// set global logger
	logger.SetGlobal(logger.NewSlogLogger(c.conf.Logger.Format, c.conf.Logger.Level, c.conf.Logger.Service, ""))

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
	// Select BTC or BCH specific sign transaction use case
	var signTxUseCase keygenusecase.SignTransactionUseCase
	if c.conf.CoinTypeCode == domainCoin.BCH {
		signTxUseCase = c.newBCHKeygenSignTransactionUseCase()
	} else {
		signTxUseCase = c.newBTCKeygenSignTransactionUseCase()
	}

	return btcwallet.NewBTCKeygen(
		c.newBTC(),
		c.pkgContainer.NewDatabaseClient(),
		c.conf.AddressType,
		c.newKeygenGenerateSeedUseCase(),
		c.newKeygenGenerateHDWalletUseCase(),
		c.newBTCKeygenImportPrivateKeyUseCase(),
		c.newBTCKeygenImportFullPubkeyUseCase(),
		c.newBTCKeygenCreateMultisigAddressUseCase(),
		c.newKeygenExportAddressUseCase(),
		signTxUseCase,
		c.walletType,
	)
}

func (c *container) newETHKeygener() wallets.Keygener {
	return ethwallet.NewETHKeygen(
		c.newETH(),
		c.pkgContainer.NewDatabaseClient(),
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
		c.pkgContainer.NewDatabaseClient(),
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
	logger.SetGlobal(logger.NewSlogLogger(c.conf.Logger.Format, c.conf.Logger.Level, c.conf.Logger.Service, ""))

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
func (c *container) NewSigner(authName string) (wallets.Signer, error) {
	// validate
	if !domainAccount.ValidateAuthType(authName) {
		if authName == "" {
			return nil, errors.New("authName is empty. Sign wallet binary must be built with ldflags.\n" +
				"Use 'make build-sign' or 'make build-all' to create sign1/sign2 binaries with proper authName.\n" +
				"Example: go build -ldflags \"-X main.authName=auth1\" -o sign1 ./cmd/sign/")
		}
		return nil, fmt.Errorf("authName '%s' is invalid. Valid values: auth1-auth15. "+
			"Sign wallet binary must be built with a valid authName via ldflags", authName)
	}

	// store authName for accessor methods
	c.authName = authName

	// set global logger
	logger.SetGlobal(logger.NewSlogLogger(c.conf.Logger.Format, c.conf.Logger.Level, c.conf.Logger.Service, ""))

	authType := domainAccount.AuthTypeMap[authName]

	switch c.conf.CoinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		return c.newBTCSigner(authType), nil
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		return nil, fmt.Errorf("coinType[%s] is not implemented yet", c.conf.CoinTypeCode)
	default:
		return nil, fmt.Errorf("coinType[%s] is not implemented yet", c.conf.CoinTypeCode)
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
func (c *container) AddressType() domainAddress.AddrType {
	return c.conf.AddressType
}

func (c *container) newBTCSigner(authType domainAccount.AuthType) wallets.Signer {
	// Select BTC or BCH specific sign transaction use case
	var signTxUseCase signusecase.SignTransactionUseCase
	if c.conf.CoinTypeCode == domainCoin.BCH {
		signTxUseCase = c.newBCHSignTransactionUseCase()
	} else {
		signTxUseCase = c.newBTCSignTransactionUseCase()
	}

	return btcwallet.NewBTCSign(
		c.newBTC(),
		c.pkgContainer.NewDatabaseClient(),
		authType,
		c.conf.AddressType,
		c.NewSignGenerateSeedUseCase(),
		c.NewSignStoreSeedUseCase(),
		c.NewSignGenerateAuthKeyUseCase(),
		c.newBTCSignImportPrivateKeyUseCase(authType),
		c.newBTCSignExportFullPubkeyUseCase(authType),
		signTxUseCase,
		c.walletType,
	)
}

func (c *container) newBTCWalleter() wallets.Watcher {
	// Select BTC or BCH specific use cases
	var createTxUseCase watchusecase.CreateTransactionUseCase
	var sendTxUseCase watchusecase.SendTransactionUseCase
	if c.conf.CoinTypeCode == domainCoin.BCH {
		createTxUseCase = c.newBCHWatchCreateTransactionUseCase()
		sendTxUseCase = c.newBCHWatchSendTransactionUseCase()
	} else {
		createTxUseCase = c.newBTCWatchCreateTransactionUseCase()
		sendTxUseCase = c.newBTCWatchSendTransactionUseCase()
	}

	return btcwallet.NewBTCWatch(
		c.newBTC(),
		c.pkgContainer.NewDatabaseClient(),
		c.conf.AddressType,
		createTxUseCase,
		c.newBTCWatchMonitorTransactionUseCase(),
		sendTxUseCase,
		c.newBTCWatchImportAddressUseCase(),
		c.newWatchCreatePaymentRequestUseCase(),
		c.walletType,
	)
}

func (c *container) newETHWalleter() wallets.Watcher {
	return ethwallet.NewETHWatch(
		c.newETH(),
		c.pkgContainer.NewDatabaseClient(),
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
		c.pkgContainer.NewDatabaseClient(),
		c.newXRPWatchCreateTransactionUseCase(),
		c.newXRPWatchMonitorTransactionUseCase(),
		c.newXRPWatchSendTransactionUseCase(),
		c.newWatchImportAddressUseCase(),
		c.newWatchCreatePaymentRequestUseCase(),
		c.walletType,
	)
}

//
// RPC Client
//

func (c *container) newRPCClient() *rpcclient.Client {
	if c.rpcClient == nil {
		var err error
		c.rpcClient, err = cryptocurrency.NewBitcoinRPCClient(&c.conf.Bitcoin)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcClient
}

func (c *container) newEthRPCClient() *ethrpc.Client {
	if c.rpcEthClient == nil {
		var err error
		c.rpcEthClient, err = cryptocurrency.NewEthereumRPCClient(&c.conf.Ethereum)
		if err != nil {
			panic(err)
		}
	}
	return c.rpcEthClient
}

func (c *container) newXRPWSClient() (*websocket.WS, *websocket.WS) {
	if c.wsXrpPublic == nil {
		var err error
		// public client
		publicURL := c.conf.Ripple.WebsocketPublicURL
		if publicURL == "" {
			if publicURL = apixrpimpl.GetPublicWSServer(c.conf.Ripple.NetworkType).String(); publicURL == "" {
				panic(errors.New("websocket URL is not found"))
			}
		}
		c.wsXrpPublic, err = websocket.New(context.Background(), publicURL)
		if err != nil {
			panic(err)
		}

		// admin client
		c.wsXrpAdmin, err = websocket.New(context.Background(), c.conf.Ripple.WebsocketAdminURL)
		if err != nil {
			panic(
				fmt.Errorf(
					"fail to call websocket.New() for admin API: %s: %w",
					c.conf.Ripple.WebsocketAdminURL, err),
			)
		}
	}
	return c.wsXrpPublic, c.wsXrpAdmin
}

//
// Wallet API
//

func (c *container) newBTC() apibtc.Bitcoiner {
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

// TODO: define interface for MuSig2Service
func (*container) newMuSig2Service() *apibtcimpl.MuSig2Service {
	return apibtcimpl.NewMuSig2Service()
}

func (c *container) newETH() apieth.Ethereumer {
	if c.eth == nil {
		var err error
		c.eth, err = ethimpl.NewEthereum(
			c.newEthRPCClient(),
			&c.conf.Ethereum,
			c.conf.CoinTypeCode,
			c.pkgContainer.NewUUIDHandler(),
		)
		if err != nil {
			panic(err)
		}
	}
	return c.eth
}

func (c *container) newERC20() apieth.ERC20er {
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
		c.erc20 = apierc20impl.NewERC20(
			client,
			tokenClient,
			conf.ERC20Token,
			c.pkgContainer.NewUUIDHandler(),
			conf.ERC20s[conf.ERC20Token].Name,
			conf.ERC20s[conf.ERC20Token].ContractAddress,
			conf.ERC20s[conf.ERC20Token].MasterAddress,
			conf.ERC20s[conf.ERC20Token].Decimals,
		)
	}
	return c.erc20
}

func (c *container) newXRP() apixrp.Rippler {
	if c.xrp == nil {
		var err error
		wsPublic, wsAdmin := c.newXRPWSClient()
		c.xrp, err = rippleimpl.NewRipple(
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

func (c *container) newRippleAPI() *apixrpimpl.RippleAPI {
	if c.rippleAPI == nil {
		c.rippleAPI = apixrpimpl.NewRippleAPI(c.pkgContainer.NewGRPCClient())
	}
	return c.rippleAPI
}

//
// Repository
//

func (c *container) newBTCTxRepo() repowatch.BTCTxRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewBTCTxRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewBTCTxRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newBTCTxInputRepo() repowatch.TxInputRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewBTCTxInputRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewBTCTxInputRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newBTCTxOutputRepo() repowatch.TxOutputRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewBTCTxOutputRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewBTCTxOutputRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newTxRepo() repowatch.TxRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewTxRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewTxRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newETHTxDetailRepo() repowatch.ETHDetailTXRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewETHDetailTXInputRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewETHDetailTXInputRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newXRPTxDetailRepo() repowatch.XRPDetailTXRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewXRPDetailTxInputRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewXRPDetailTxInputRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newPaymentRequestRepo() repowatch.PaymentRequestRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewPaymentRequestRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewPaymentRequestRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newAddressRepo() repowatch.AddressRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return watchmysql.NewAddressRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return watchsqlite.NewAddressRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newAddressFileRepo() file.AddressFileRepositorier {
	return address.NewAddressFileRepository(
		c.conf.FilePath.Address,
	)
}

func (c *container) newTxFileRepo() file.TransactionFileRepositorier {
	return transaction.NewTransactionFileRepository(
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

func (c *container) newHdWalletRepo() repocold.HDWalletRepo {
	return coldmysql.NewAccountHDWalletRepo(
		c.newAccountKeyRepo(),
	)
}

func (c *container) newKeyGenerator() portsWallet.Generator {
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

	// Use factory to create generator based on key type
	factory := infraKeyGen.NewFactory()
	keyType := c.getKeyType() // Get from config or default to BIP44
	generator, err := factory.CreateGenerator(keyType, c.conf.CoinTypeCode, chainConf)
	if err != nil {
		panic(fmt.Sprintf("failed to create key generator: %v", err))
	}

	return generator
}

func (c *container) getKeyType() domainKey.KeyType {
	// Derive KeyType from AddressType to ensure consistency
	// address_type and key_type have a 1:1 mapping:
	//   - legacy      -> bip44
	//   - p2sh-segwit -> bip49
	//   - bech32      -> bip84
	//   - taproot     -> bip86
	keyType, err := c.conf.AddressType.ToKeyType()
	if err != nil {
		// Unsupported address type is a critical configuration error.
		// Panic to prevent silent failures that could lead to incorrect
		// key generation and potential asset loss.
		panic(fmt.Sprintf("failed to derive key type from address type %q: %v", c.conf.AddressType, err))
	}
	return keyType
}

func (c *container) newMultiAccount() *domainAccount.MultisigConfig {
	if c.multisig == nil {
		if c.accountConf == nil || c.accountConf.Multisigs == nil {
			c.multisig = config.NewMultisigConfig(nil)
		} else {
			c.multisig = config.NewMultisigConfig(c.accountConf.Multisigs)
		}
	}
	return c.multisig
}

//
// Keygen Repository
//

func (c *container) newSeedRepo() repocold.SeedRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewSeedRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return coldsqlite.NewSeedRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newAccountKeyRepo() repocold.BTCAccountKeyRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewBTCAccountKeyRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return coldsqlite.NewBTCAccountKeyRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newXRPAccountKeyRepo() repocold.XRPAccountKeyRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewXRPAccountKeyRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return coldsqlite.NewXRPAccountKeyRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newEthAccountKeyRepo() repocold.ETHAccountKeyRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewETHAccountKeyRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
		)
	case "sqlite":
		return coldsqlite.NewETHAccountKeyRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newAuthFullPubKeyRepo() repocold.AuthFullPubkeyRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewAuthFullPubkeyRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return coldsqlite.NewAuthFullPubkeyRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newAuthKeyRepo() repocold.AuthAccountKeyRepositorier {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewAuthAccountKeyRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
			c.conf.CoinTypeCode,
		)
	case "sqlite":
		return coldsqlite.NewAuthAccountKeyRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
			c.conf.CoinTypeCode,
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

func (c *container) newNonceRepo() multisig.NonceRepository {
	switch c.conf.Database.Type {
	case "mysql":
		return coldmysql.NewNonceRepositorySqlc(
			c.pkgContainer.NewDatabaseClient(),
		)
	case "sqlite":
		return coldsqlite.NewNonceRepositorySqlc(
			c.pkgContainer.NewSQLiteClient(),
		)
	default:
		panic("unsupported database type: " + c.conf.Database.Type)
	}
}

//
// Keygen File Storage
//

func (c *container) newPubkeyFileStorager() file.AddressFileRepositorier {
	return address.NewAddressFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (*container) newDescriptorFileWriter() file.DescriptorFileWriter {
	return descriptor.NewFileWriter()
}

//
// Sign Service
//

func (c *container) newSignHdWalletRepo(authType domainAccount.AuthType) repocold.HDWalletRepo {
	return coldmysql.NewAuthHDWalletRepo(
		c.newAuthKeyRepo(),
		authType,
	)
}

//
// Use Case Factory Methods
//

// Watch Use Cases

func (c *container) NewWatchCreateTransactionUseCase() any {
	switch {
	case domainCoin.IsBTCOnly(c.conf.CoinTypeCode):
		return c.newBTCWatchCreateTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.BCH:
		return c.newBCHWatchCreateTransactionUseCase()
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
	case domainCoin.IsBTCOnly(c.conf.CoinTypeCode):
		return c.newBTCWatchSendTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.BCH:
		return c.newBCHWatchSendTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHWatchSendTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPWatchSendTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewWatchImportAddressUseCase() watchusecase.ImportAddressUseCase {
	// Use coin-specific implementations for BTC and BCH
	if c.conf.CoinTypeCode == domainCoin.BTC {
		return c.newBTCWatchImportAddressUseCase()
	}
	if c.conf.CoinTypeCode == domainCoin.BCH {
		return c.newBCHWatchImportAddressUseCase()
	}
	// Use shared implementation for other coins (ETH, XRP, etc.)
	return c.newWatchImportAddressUseCase()
}

func (c *container) NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase {
	return c.newWatchCreatePaymentRequestUseCase()
}

func (c *container) NewWatchImportDescriptorUseCase() watchusecase.ImportDescriptorUseCase {
	// Descriptors are BTC-only (BCH does not support descriptors)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf(
			"descriptor import is BTC-only, not supported for %s. Use ImportAddress instead.",
			c.conf.CoinTypeCode))
	}
	return c.newBTCWatchImportDescriptorUseCase()
}

func (c *container) NewWatchAggregateMuSig2SignaturesUseCase() watchusecase.AggregateMuSig2SignaturesUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s. Use P2SH multisig instead.", c.conf.CoinTypeCode))
	}
	return c.newBTCWatchAggregateMuSig2SignaturesUseCase()
}

// Keygen Use Cases

func (c *container) NewKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase {
	return c.newKeygenGenerateHDWalletUseCase()
}

func (c *container) NewKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase {
	return c.newKeygenGenerateSeedUseCase()
}

func (c *container) NewKeygenGenerateDescriptorUseCase() keygenusecase.GenerateDescriptorUseCase {
	// Descriptors are BTC-only (BCH does not support descriptors)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("descriptor generation is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenGenerateDescriptorUseCase()
}

func (c *container) NewKeygenExportDescriptorUseCase() keygenusecase.ExportDescriptorUseCase {
	// Descriptors are BTC-only (BCH does not support descriptors)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("descriptor export is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenExportDescriptorUseCase()
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

func (c *container) NewKeygenCreateMuSig2AddressUseCase() keygenusecase.CreateMuSig2AddressUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s. Use P2SH multisig instead.", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenCreateMuSig2AddressUseCase()
}

func (c *container) NewKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase {
	return c.newBTCKeygenImportFullPubkeyUseCase()
}

func (c *container) NewKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase {
	return c.newXRPKeygenGenerateKeyUseCase()
}

func (c *container) NewKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	switch {
	case domainCoin.IsBTCOnly(c.conf.CoinTypeCode):
		return c.newBTCKeygenSignTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.BCH:
		return c.newBCHKeygenSignTransactionUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newETHKeygenSignTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newXRPKeygenSignTransactionUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewKeygenGenerateMuSig2NonceUseCase() keygenusecase.GenerateMuSig2NonceUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenGenerateMuSig2NonceUseCase()
}

func (c *container) NewKeygenMuSig2SignUseCase() keygenusecase.MuSig2SignUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenMuSig2SignUseCase()
}

// Sign Use Cases

func (c *container) NewSignTransactionUseCase() signusecase.SignTransactionUseCase {
	switch {
	case domainCoin.IsBTCOnly(c.conf.CoinTypeCode):
		return c.newBTCSignTransactionUseCase()
	case c.conf.CoinTypeCode == domainCoin.BCH:
		return c.newBCHSignTransactionUseCase()
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

func (c *container) NewSignGenerateMuSig2NonceUseCase() signusecase.GenerateMuSig2NonceUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCSignGenerateMuSig2NonceUseCase()
}

func (c *container) NewSignMuSig2SignUseCase() signusecase.MuSig2SignUseCase {
	// MuSig2 is BTC-only (BCH does not support Schnorr signatures)
	if !domainCoin.IsBTCOnly(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("MuSig2 is BTC-only, not supported for %s", c.conf.CoinTypeCode))
	}
	return c.newBTCSignMuSig2SignUseCase()
}

// BTC Watch Use Cases

func (c *container) newBTCWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	return watchusecasebtc.NewCreateTransactionUseCase(
		c.newBTC(),
		c.pkgContainer.NewDatabaseClient(),
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
		c.pkgContainer.NewDatabaseClient(),
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

func (c *container) newBTCWatchImportDescriptorUseCase() watchusecase.ImportDescriptorUseCase {
	return watchusecasebtc.NewImportDescriptorUseCase(
		c.newBTC(),
		apibtcimpl.NewDescriptorParser(),
		c.newBTC().GetChainConf(),
		c.newAddressRepo(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCWatchAggregateMuSig2SignaturesUseCase() watchusecase.AggregateMuSig2SignaturesUseCase {
	return watchusecasebtc.NewAggregateMuSig2SignaturesUseCase(
		c.newMuSig2Service(),
		c.newBTC(),
	)
}

// BCH Watch Use Cases

func (c *container) newBCHWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	return watchusecasebch.NewCreateBCHTransactionUseCase(
		c.newBTC(), // BCHer interface is satisfied by Bitcoiner
		c.pkgContainer.NewDatabaseClient(),
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

func (c *container) newBCHWatchSendTransactionUseCase() watchusecase.SendTransactionUseCase {
	return watchusecasebch.NewBCHSendTransactionUseCase(
		c.newBTC(), // BCHer interface is satisfied by Bitcoiner
		c.newAddressRepo(),
		c.newBTCTxRepo(),
		c.newBTCTxOutputRepo(),
		c.newTxFileRepo(),
	)
}

func (c *container) newBCHWatchImportAddressUseCase() watchusecase.ImportAddressUseCase {
	return watchusecasebch.NewImportBCHAddressUseCase(
		c.newBTC(), // BCH BitcoinCash struct embeds Bitcoin, providing all methods
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
	)
}

// ETH Watch Use Cases

func (c *container) newETHWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	// Determine which Ethereum API to use based on coin type
	var targetEthAPI apieth.EtherTxCreator
	if domainCoin.IsERC20Token(c.conf.CoinTypeCode.String()) {
		targetEthAPI = c.newERC20()
	} else {
		targetEthAPI = c.newETH()
	}

	return watchusecaseeth.NewCreateTransactionUseCase(
		targetEthAPI,
		c.pkgContainer.NewDatabaseClient(),
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
		c.pkgContainer.NewDatabaseClient(),
		c.pkgContainer.NewUUIDHandler(),
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
		c.newAddressRepo(),
		c.newAddressFileRepo(),
		c.conf.CoinTypeCode,
		c.conf.AddressType,
		c.walletType,
	)
}

func (c *container) newWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase {
	return watchusecaseshared.NewCreatePaymentRequestUseCase(
		c.pkgContainer.NewDatabaseClient(),
		c.newAddressRepo(),
		c.newPaymentRequestRepo(),
		c.conf.CoinTypeCode,
		c.walletType,
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
	return keygenusecasebtc.NewExportAddressUseCase(
		c.newAccountKeyRepo(),
		c.newAddressFileRepo(),
		c.newMultiAccount(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCKeygenGenerateDescriptorUseCase() keygenusecase.GenerateDescriptorUseCase {
	// Create coin strategy once at DI initialization
	chainConf := c.newBTC().GetChainConf()
	coinStrategy, err := infraKey.CreateCoinKeyStrategy(c.conf.CoinTypeCode, chainConf)
	if err != nil {
		panic(fmt.Sprintf("failed to create coin strategy for descriptor use case: %v", err))
	}

	return keygenusecasebtc.NewGenerateDescriptorUseCase(
		apibtcimpl.NewDescriptorService(chainConf),
		chainConf,
		c.newAuthFullPubKeyRepo(),
		c.newAccountKeyRepo(),
		c.newSeedRepo(),
		c.conf.CoinTypeCode,
		c.newMultiAccount(),
		coinStrategy,
	)
}

func (c *container) newBTCKeygenExportDescriptorUseCase() keygenusecase.ExportDescriptorUseCase {
	return keygenusecasebtc.NewExportDescriptorUseCase(
		c.newBTCKeygenGenerateDescriptorUseCase(),
		c.newDescriptorFileWriter(),
		c.newAccountKeyRepo(),
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

func (c *container) newBTCKeygenCreateMuSig2AddressUseCase() keygenusecase.CreateMuSig2AddressUseCase {
	return keygenusecasebtc.NewCreateMuSig2AddressUseCase(
		c.newMuSig2Service(),
		c.newBTC().GetChainConf(),
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
		c.newEthAccountKeyRepo(),
	)
}

// XRP Keygen Use Cases

func (c *container) newXRPKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase {
	return keygenusecasexrp.NewGenerateKeyUseCase(
		c.newXRP(),
		c.pkgContainer.NewDatabaseClient(),
		c.conf.CoinTypeCode,
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

func (c *container) newBCHKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase {
	return keygenusecasebch.NewBCHSignTransactionUseCase(
		c.newBTC(), // BCHer interface is satisfied by Bitcoiner
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

// Keygen MuSig2 Use Cases

func (c *container) newBTCKeygenGenerateMuSig2NonceUseCase() keygenusecase.GenerateMuSig2NonceUseCase {
	return keygenusecasebtc.NewGenerateMuSig2NonceUseCase(
		c.newMuSig2Service(),
		c.newNonceRepo(),
	)
}

func (c *container) newBTCKeygenMuSig2SignUseCase() keygenusecase.MuSig2SignUseCase {
	return keygenusecasebtc.NewMuSig2SignUseCase(
		c.newMuSig2Service(),
		c.newNonceRepo(),
	)
}

// Sign Use Cases

// BTC Sign Use Cases

func (c *container) newBTCSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecasebtc.NewSignTransactionUseCase(
		c.newBTC(),
		c.newAccountKeyRepo(),
		c.newAuthKeyRepo(),
		c.newTxFileRepo(),
		c.newMultiAccount(),
		c.walletType,
		c.AuthType(),
	)
}

func (c *container) newBCHSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecasebch.NewBCHSignTransactionUseCase(
		c.newBTC(), // BCHer interface is satisfied by Bitcoiner
		c.newAccountKeyRepo(),
		c.newAuthKeyRepo(),
		c.newTxFileRepo(),
		c.newMultiAccount(),
		c.walletType,
		c.AuthType(),
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
		c.newSeedRepo(),
		c.newPubkeyFileStorager(),
		c.conf.CoinTypeCode,
		authType,
		c.walletType,
		c.newBTC().GetChainConf(),
	)
}

func (c *container) newBTCSignGenerateMuSig2NonceUseCase() signusecase.GenerateMuSig2NonceUseCase {
	return signusecasebtc.NewGenerateMuSig2NonceUseCase(
		c.newMuSig2Service(),
		c.newNonceRepo(),
		c.newAuthKeyRepo(),
	)
}

func (c *container) newBTCSignMuSig2SignUseCase() signusecase.MuSig2SignUseCase {
	return signusecasebtc.NewMuSig2SignUseCase(
		c.newMuSig2Service(),
		c.newNonceRepo(),
		c.newAuthKeyRepo(),
	)
}

// ETH Sign Use Cases

func (c *container) newETHSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecaseeth.NewSignTransactionUseCase(
		c.newETH(),
		c.newTxFileRepo(),
		c.walletType,
	)
}

// XRP Sign Use Cases

func (c *container) newXRPSignTransactionUseCase() signusecase.SignTransactionUseCase {
	return signusecasexrp.NewSignTransactionUseCase(
		c.newXRP(),
		c.newXRPAccountKeyRepo(),
		c.newTxFileRepo(),
		c.walletType,
	)
}
