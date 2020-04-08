package wallet

import (
	"github.com/btcsuite/btcutil"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// Walleter is for watch only wallet service interface
type Walleter interface {
	ImportPublicKeyForWatchWallet(fileName string, accountType account.AccountType, isRescan bool) error
	DetectReceivedCoin(adjustmentFee float64) (string, string, error)
	CreateUnsignedTransactionForPayment(adjustmentFee float64) (string, string, error)
	SendToAccount(from, to account.AccountType, amount btcutil.Amount) (string, string, error)
	SendFromFile(filePath string) (string, error)
	UpdateStatus() error
	Done()
	GetDB() *model.DB
	GetBTC() api.Bitcoiner
	GetType() WalletType
}

// Wallet watch only wallet object
type Wallet struct {
	BTC  api.Bitcoiner
	DB   *model.DB //TODO:should be interface
	Type WalletType
}

func NewWallet(bit api.Bitcoiner, rds *model.DB, typ WalletType) *Wallet {
	return &Wallet{
		BTC:  bit,
		DB:   rds,
		Type: typ,
	}
}

// Done should be called before exit
func (w *Wallet) Done() {
	w.DB.RDB.Close()
	w.BTC.Close()
}

func (w *Wallet) GetDB() *model.DB {
	return w.DB
}

func (w *Wallet) GetBTC() api.Bitcoiner {
	return w.BTC
}

func (w *Wallet) GetType() WalletType {
	return w.Type
}
