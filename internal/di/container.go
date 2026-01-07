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

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	portsEth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/ethereum"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsRipple "github.com/hiromaily/go-crypto-wallet/internal/application/ports/ripple"
	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	btcapi "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	ethimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum/erc20"
	rippleimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/contract"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/descriptor"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
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
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	keygenusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/shared"
	keygenusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/xrp"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	signusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/btc"
	signusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/eth"
	signusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/shared"
	signusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	watchusecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/shared"
	watchusecasexrp "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
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
	btc        portsBtc.Bitcoiner
	eth        portsEth.Ethereumer
	erc20      portsEth.ERC20er
	xrp        portsRipple.Rippler
	// client
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *websocket.WS
	wsXrpAdmin   *websocket.WS
	rippleAPI    *xrp.RippleAPI
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
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
func (c *container) AddressType() domainAddress.AddrType {
	return c.conf.AddressType
}

func (c *container) newBTCSigner(authType domainAccount.AuthType) wallets.Signer {
	return btcwallet.NewBTCSign(
		c.newBTC(),
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
			if publicURL = xrp.GetPublicWSServer(c.conf.Ripple.NetworkType).String(); publicURL == "" {
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

func (c *container) newBTC() portsBtc.Bitcoiner {
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

func (c *container) newMuSig2Service() *btcapi.MuSig2Service {
	return btcapi.NewMuSig2Service(c.pkgContainer.NewLogger())
}

func (c *container) newETH() portsEth.Ethereumer {
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

func (c *container) newERC20() portsEth.ERC20er {
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
			c.pkgContainer.NewUUIDHandler(),
			conf.ERC20s[conf.ERC20Token].Name,
			conf.ERC20s[conf.ERC20Token].ContractAddress,
			conf.ERC20s[conf.ERC20Token].MasterAddress,
			conf.ERC20s[conf.ERC20Token].Decimals,
		)
	}
	return c.erc20
}

func (c *container) newXRP() portsRipple.Rippler {
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

func (c *container) newRippleAPI() *xrp.RippleAPI {
	if c.rippleAPI == nil {
		c.rippleAPI = xrp.NewRippleAPI(c.pkgContainer.NewGRPCClient())
	}
	return c.rippleAPI
}

//
// Repository
//

func (c *container) newBTCTxRepo() persistence.BTCTxRepositorier {
	return watch.NewBTCTxRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCTxInputRepo() persistence.TxInputRepositorier {
	return watch.NewBTCTxInputRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newBTCTxOutputRepo() persistence.TxOutputRepositorier {
	return watch.NewBTCTxOutputRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newTxRepo() persistence.TxRepositorier {
	return watch.NewTxRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newETHTxDetailRepo() persistence.ETHDetailTXRepositorier {
	return watch.NewETHDetailTXInputRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPTxDetailRepo() persistence.XRPDetailTXRepositorier {
	return watch.NewXRPDetailTxInputRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newPaymentRequestRepo() persistence.PaymentRequestRepositorier {
	return watch.NewPaymentRequestRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressRepo() persistence.AddressRepositorier {
	return watch.NewAddressRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAddressFileRepo() portsStorage.AddressFileRepositorier {
	return address.NewAddressFileRepository(
		c.conf.FilePath.Address,
	)
}

func (c *container) newTxFileRepo() portsStorage.TransactionFileRepositorier {
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

func (c *container) newHdWalletRepo() persistence.HDWalletRepo {
	return cold.NewAccountHDWalletRepo(
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
	factory := infraKey.NewFactory()
	keyType := c.getKeyType() // Get from config or default to BIP44
	generator, err := factory.CreateGenerator(keyType, c.conf.CoinTypeCode, chainConf)
	if err != nil {
		panic(fmt.Sprintf("failed to create key generator: %v", err))
	}

	return generator
}

func (c *container) getKeyType() domainKey.KeyType {
	// Get from config if available, otherwise default to BIP44
	if c.conf.KeyType != "" {
		return c.conf.KeyType
	}
	return domainKey.KeyTypeBIP44
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

func (c *container) newSeedRepo() persistence.SeedRepositorier {
	return cold.NewSeedRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAccountKeyRepo() persistence.BTCAccountKeyRepositorier {
	return cold.NewAccountKeyRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newXRPAccountKeyRepo() persistence.XRPAccountKeyRepositorier {
	return cold.NewXRPAccountKeyRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newEthAccountKeyRepo() persistence.ETHAccountKeyRepositorier {
	return cold.NewETHAccountKeyRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
	)
}

func (c *container) newAuthFullPubKeyRepo() persistence.AuthFullPubkeyRepositorier {
	return cold.NewAuthFullPubkeyRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newAuthKeyRepo() persistence.AuthAccountKeyRepositorier {
	return cold.NewAuthAccountKeyRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
		c.conf.CoinTypeCode,
	)
}

func (c *container) newNonceRepo() *cold.NonceRepositorySqlc {
	return cold.NewNonceRepositorySqlc(
		c.pkgContainer.NewMySQLClient(),
	)
}

//
// Keygen File Storage
//

func (c *container) newPubkeyFileStorager() portsStorage.AddressFileRepositorier {
	return address.NewAddressFileRepository(
		c.conf.FilePath.FullPubKey,
	)
}

func (*container) newDescriptorFileWriter() portsStorage.DescriptorFileWriter {
	return descriptor.NewFileWriter()
}

//
// Sign Service
//

func (c *container) newSignHdWalletRepo(authType domainAccount.AuthType) persistence.HDWalletRepo {
	return cold.NewAuthHDWalletRepo(
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
	switch {
	case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
		return c.newBTCWatchImportAddressUseCase()
	case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
		return c.newWatchImportAddressUseCase()
	case c.conf.CoinTypeCode == domainCoin.XRP:
		return c.newWatchImportAddressUseCase()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
	}
}

func (c *container) NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase {
	return c.newWatchCreatePaymentRequestUseCase()
}

func (c *container) NewWatchImportDescriptorUseCase() watchusecase.ImportDescriptorUseCase {
	if !domainCoin.IsBTCGroup(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("descriptor import supported only for BTC group, got %s", c.conf.CoinTypeCode))
	}
	return c.newBTCWatchImportDescriptorUseCase()
}

func (c *container) NewWatchAggregateMuSig2SignaturesUseCase() watchusecase.AggregateMuSig2SignaturesUseCase {
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
	if !domainCoin.IsBTCGroup(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("descriptor generation is only supported for BTC group coins, got %s", c.conf.CoinTypeCode))
	}
	return c.newBTCKeygenGenerateDescriptorUseCase()
}

func (c *container) NewKeygenExportDescriptorUseCase() keygenusecase.ExportDescriptorUseCase {
	if !domainCoin.IsBTCGroup(c.conf.CoinTypeCode) {
		panic(fmt.Sprintf("descriptor export is only supported for BTC group coins, got %s", c.conf.CoinTypeCode))
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

func (c *container) NewKeygenGenerateMuSig2NonceUseCase() keygenusecase.GenerateMuSig2NonceUseCase {
	return c.newBTCKeygenGenerateMuSig2NonceUseCase()
}

func (c *container) NewKeygenMuSig2SignUseCase() keygenusecase.MuSig2SignUseCase {
	return c.newBTCKeygenMuSig2SignUseCase()
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

func (c *container) NewSignGenerateMuSig2NonceUseCase() signusecase.GenerateMuSig2NonceUseCase {
	return c.newBTCSignGenerateMuSig2NonceUseCase()
}

func (c *container) NewSignMuSig2SignUseCase() signusecase.MuSig2SignUseCase {
	return c.newBTCSignMuSig2SignUseCase()
}

// BTC Watch Use Cases

func (c *container) newBTCWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	return watchusecasebtc.NewCreateTransactionUseCase(
		c.newBTC(),
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
		btcapi.NewDescriptorParser(),
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

// ETH Watch Use Cases

func (c *container) newETHWatchCreateTransactionUseCase() watchusecase.CreateTransactionUseCase {
	// Determine which Ethereum API to use based on coin type
	var targetEthAPI portsEth.EtherTxCreator
	if domainCoin.IsERC20Token(c.conf.CoinTypeCode.String()) {
		targetEthAPI = c.newERC20()
	} else {
		targetEthAPI = c.newETH()
	}

	return watchusecaseeth.NewCreateTransactionUseCase(
		targetEthAPI,
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
		c.pkgContainer.NewMySQLClient(),
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
	return keygenusecasebtc.NewGenerateDescriptorUseCase(
		btcapi.NewDescriptorService(c.newBTC().GetChainConf()),
		c.newBTC().GetChainConf(),
		c.newAuthFullPubKeyRepo(),
		c.newAccountKeyRepo(),
		c.newMultiAccount(),
	)
}

func (c *container) newBTCKeygenExportDescriptorUseCase() keygenusecase.ExportDescriptorUseCase {
	return keygenusecasebtc.NewExportDescriptorUseCase(
		c.newBTCKeygenGenerateDescriptorUseCase(),
		c.newDescriptorFileWriter(),
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
		c.pkgContainer.NewMySQLClient(),
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
