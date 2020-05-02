package wallets

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/repository/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv/keygensrv"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/coldwalletsrv/signsrv"
)

// About structure
// Wallets/wallet
//        /coldwallet ... has any func for both keygen and signature
//        /keygen     ... has only keygen interface
//        /signature  ... has only signature interface

// Watcher is for watch only wallet service interface
type Watcher interface {
	ImportAddress(fileName string, accountType account.AccountType, isRescan bool) error
	CreateDepositTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error)
	SendTx(filePath string) (string, error)
	UpdateTxStatus() error

	Done()
	GetDB() walletrepo.WalletRepository // for debug use
	GetBTC() api.Bitcoiner
	GetType() wallet.WalletType
}

// Coldwalleter may not be used anywhere
//type Coldwalleter interface {
//	KeySigner
//	KeygenExclusiver
//	SignatureExclusiver
//
//	Done()
//	GetBTC() api.Bitcoiner
//	GetType() wallet.WalletType
//}

// KeySigner is common interface for keygen/signature
//type KeySigner interface {
//	GenerateSeed() ([]byte, error)
//	StoreSeed(strSeed string) ([]byte, error)
//	GeneratePubKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
//	ImportPrivateKey(accountType account.AccountType) error
//	SignTx(filePath string) (string, bool, string, error)
//}

// Keygener is for keygen wallet service interface
type Keygener interface {
	Seed() coldwalletsrv.Seeder
	HDWallet() coldwalletsrv.HDWalleter
	PrivKey() keygensrv.PrivKeyer
	FullPubKeyImport() keygensrv.FullPubKeyImporter
	Multisig() keygensrv.Multisiger
	AddressExport() keygensrv.AddressExporter
	Sign() coldwalletsrv.Signer

	Done()
	GetBTC() api.Bitcoiner
	GetType() wallet.WalletType
}

// Signer is for signature wallet service interface
type Signer interface {
	Seed() coldwalletsrv.Seeder
	HDWallet() coldwalletsrv.HDWalleter
	PrivKey() signsrv.PrivKeyer
	FullPubkeyExport() signsrv.FullPubkeyExporter
	Sign() coldwalletsrv.Signer

	Done()
	GetBTC() api.Bitcoiner
	GetType() wallet.WalletType
	GetAuthType() account.AuthType
}
