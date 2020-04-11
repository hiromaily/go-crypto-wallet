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
	GetDB() rdb.SignatureStorager
	GetBTC() api.Bitcoiner
	GetType() WalletType
}

// Signature signature wallet object
//  it is almost same to Wallet object, difference is storager interface
type Signature struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	storager rdb.SignatureStorager
	wtype    WalletType
}

func NewSignature(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.SignatureStorager,
	wtype WalletType) *Signature {

	return &Signature{
		btc:      btc,
		logger:   logger,
		tracer:   tracer,
		storager: storager,
		wtype:    wtype,
	}
}

// Done should be called before exit
func (w *Signature) Done() {
	w.storager.Close()
	w.btc.Close()
}

func (w *Signature) GetDB() rdb.SignatureStorager {
	return w.storager
}

func (w *Signature) GetBTC() api.Bitcoiner {
	return w.btc
}

func (w *Signature) GetType() WalletType {
	return w.wtype
}
