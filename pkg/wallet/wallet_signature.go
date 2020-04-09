package wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// Signer is for signature wallet service interface
type Signer interface {
	SignatureFromFile(filePath string) (string, bool, string, error)
	GenerateSeed() ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, coinType enum.CoinType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivateKey(accountType account.AccountType) error
	ImportPublicKeyForColdWallet2(fileName string, accountType account.AccountType) error
	AddMultisigAddressByAuthorization(accountType account.AccountType, addressType enum.AddressType) error
	ExportAddedPubkeyHistory(accountType account.AccountType) (string, error)
	Done()
	GetDB() *model.DB
	GetBTC() api.Bitcoiner
	GetType() WalletType
}
