package wallets

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// About structure
// Wallets/wallet
//        /coldwallet ... has any func for both keygen and signature
//        /keygen     ... has only keygen interface
//        /signature  ... has only signature interface

// Walleter is for watch only wallet service interface
type Walleter interface {
	ImportPubKey(fileName string, accountType account.AccountType, isRescan bool) error
	CreateReceiptTx(adjustmentFee float64) (string, string, error)
	CreatePaymentTx(adjustmentFee float64) (string, string, error)
	CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error)
	SendTx(filePath string) (string, error)
	UpdateTxStatus() error

	Done()
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// Coldwalleter may not be used anywhere
type Coldwalleter interface {
	KeySigner
	KeygenExclusiver
	SignatureExclusiver

	Done()
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// KeySigner is common interface for keygen/signature
type KeySigner interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GeneratePubKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivateKey(accountType account.AccountType) error
	SignTx(filePath string) (string, bool, string, error)
}

// Keygener is for keygen wallet service interface
type Keygener interface {
	KeySigner
	KeygenExclusiver

	Done()
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// KeygenExclusiver is for only Keygen interface
type KeygenExclusiver interface {
	ExportAccountKey(accountType account.AccountType, addrStatus address.AddrStatus) (string, error)
	ImportMultisigAddress(fileName string, accountType account.AccountType) error
}

// Signer is for signature wallet service interface
type Signer interface {
	KeySigner
	SignatureExclusiver

	Done()
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// SignatureExclusiver is for only signature interface
type SignatureExclusiver interface {
	ImportPubKey(fileName string, accountType account.AccountType) error
	AddMultisigAddress(accountType account.AccountType, addressType address.AddrType) error
	ExportAddedPubkeyHistory(accountType account.AccountType) (string, error)
}
