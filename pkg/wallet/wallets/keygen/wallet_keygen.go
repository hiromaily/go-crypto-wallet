package keygen

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

//// Keygener is for keygen wallet service interface
//type Keygener interface {
//	coldwallet.KeySigner
//
//	KeygenExclusiver
//
//	Done()
//	GetDB() rdb.KeygenStorager
//	GetBTC() api.Bitcoiner
//	GetType() types.WalletType
//}
//
//type KeygenExclusiver interface {
//	ExportAccountKey(accountType account.AccountType, addrStatus keystatus.AddrStatus) (string, error)
//	ImportMultisigAddress(fileName string, accountType account.AccountType) error
//}

// Keygen keygen wallet object
//  it is almost same to Wallet object, difference is storager interface
type Keygen struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	storager rdb.ColdStorager
	wtype    types.WalletType
}

// NewKeygen may be Not used anywhere
func NewKeygen(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.ColdStorager,
	wtype types.WalletType) *Keygen {

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

func (w *Keygen) GetDB() rdb.ColdStorager {
	return w.storager
}

func (w *Keygen) GetBTC() api.Bitcoiner {
	return w.btc
}

func (w *Keygen) GetType() types.WalletType {
	return w.wtype
}
