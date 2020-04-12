package coldwallet

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// About structure
// Wallets/wallet
//        /coldwallet ... has any func for both keygen and signature
//        /keygen     ... has only keygen interface
//        /signature  ... has only signature interface

//// Coldwalleter may be Not used anywhere
//type Coldwalleter interface {
//	KeySigner
//	keygen.KeygenExclusiver
//	signature.SignatureExclusiver
//
//	Done()
//	GetDB() rdb.ColdStorager
//	GetBTC() api.Bitcoiner
//	GetType() types.WalletType
//}
//
//// common interface for keygen/signature
//type KeySigner interface {
//	SignatureFromFile(filePath string) (string, bool, string, error)
//	GenerateSeed() ([]byte, error)
//	GenerateAccountKey(accountType account.AccountType, coinType enum.CoinType, seed []byte, count uint32) ([]key.WalletKey, error)
//	ImportPrivateKey(accountType account.AccountType) error
//}

// TODO: this object has to includes any func in structure
// ColdWallet coldwallet for keygen/signature object
type ColdWallet struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	storager rdb.ColdStorager
	wtype    types.WalletType
}

func NewColdWalet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.ColdStorager,
	wtype types.WalletType) *ColdWallet {

	return &ColdWallet{
		btc:      btc,
		logger:   logger,
		tracer:   tracer,
		storager: storager,
		wtype:    wtype,
	}
}

// Done should be called before exit
func (w *ColdWallet) Done() {
	w.storager.Close()
	w.btc.Close()
}

func (w *ColdWallet) GetDB() rdb.ColdStorager {
	return w.storager
}

func (w *ColdWallet) GetBTC() api.Bitcoiner {
	return w.btc
}

func (w *ColdWallet) GetType() types.WalletType {
	return w.wtype
}
