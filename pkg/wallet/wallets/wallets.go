package wallets

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// About structure
// Wallets/wallet
//        /coldwallet ... has any func for both keygen and signature
//        /keygen     ... has only keygen interface
//        /signature  ... has only signature interface

// Keygener is for keygen wallet service interface
type Keygener interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivKey(accountType account.AccountType) error
	ImportFullPubKey(fileName string) error
	CreateMultisigAddress(accountType account.AccountType, addressType address.AddrType) error
	ExportAddress(accountType account.AccountType) (string, error)
	SignTx(filePath string) (string, bool, string, error)

	Done()
	GetBTC() api.Bitcoiner
	//GetType() wallet.WalletType
}

// Signer is for signature wallet service interface
type Signer interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAuthKey(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivKey() error
	ExportFullPubkey() (string, error)
	SignTx(filePath string) (string, bool, string, error)

	Done()
	//GetBTC() api.Bitcoiner
	//GetType() wallet.WalletType
	GetAuthType() account.AuthType
}
