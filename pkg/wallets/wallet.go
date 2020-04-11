package wallets

import (
	"github.com/btcsuite/btcutil"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
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
	GetDB() rdb.WalletStorager
	GetBTC() api.Bitcoiner
	GetType() WalletType
}

// Wallet watch only wallet object
type Wallet struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	storager rdb.WalletStorager
	wtype    WalletType
}

func NewWallet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.WalletStorager,
	wtype WalletType) *Wallet {

	return &Wallet{
		btc:      btc,
		logger:   logger,
		tracer:   tracer,
		storager: storager,
		wtype:    wtype,
	}
}

// Done should be called before exit
func (w *Wallet) Done() {
	w.storager.Close()
	w.btc.Close()
}

func (w *Wallet) GetDB() rdb.WalletStorager {
	return w.storager
}

func (w *Wallet) GetBTC() api.Bitcoiner {
	return w.btc
}

func (w *Wallet) GetType() WalletType {
	return w.wtype
}
