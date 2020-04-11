package wallets

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
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
	GetDB() rdb.KeygenStorager
	GetBTC() api.Bitcoiner
	GetType() WalletType
}

// Keygen keygen wallet object
//  it is almost same to Wallet object, difference is storager interface
type Keygen struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	storager rdb.KeygenStorager
	wtype    WalletType
}

func NewKeygen(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.KeygenStorager,
	wtype WalletType) *Keygen {

	return &Keygen{
		btc:      btc,
		logger:   logger,
		tracer:   tracer,
		storager: storager,
		wtype:    wtype,
	}
}

// Done should be called before exit
func (w *Keygen) Done() {
	w.storager.Close()
	w.btc.Close()
}

func (w *Keygen) GetDB() rdb.KeygenStorager {
	return w.storager
}

func (w *Keygen) GetBTC() api.Bitcoiner {
	return w.btc
}

func (w *Keygen) GetType() WalletType {
	return w.wtype
}
