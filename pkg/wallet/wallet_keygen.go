package wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// Keygener is for keygen wallet service interface
type Keygener interface {
	SignatureFromFile(filePath string) (string, bool, string, error)
	GenerateSeed() ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, coinType enum.CoinType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivateKey(accountType account.AccountType) error
	ExportAccountKey(accountType account.AccountType, keyStatus enum.KeyStatus) (string, error)
	ImportMultisigAddrForColdWallet1(fileName string, accountType account.AccountType) error
	Done()
	GetDB() *model.DB
	GetBTC() api.Bitcoiner
	GetType() WalletType
}
