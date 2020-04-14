package wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// Wallet watch only wallet object
type Wallet struct {
	btc              api.Bitcoiner
	logger           *zap.Logger
	tracer           opentracing.Tracer
	storager         rdb.WalletStorager
	addrFileStorager address.Storager
	wtype            types.WalletType
}

func NewWallet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.WalletStorager,
	addrFileStorager address.Storager,
	wtype types.WalletType) *Wallet {

	return &Wallet{
		btc:              btc,
		logger:           logger,
		tracer:           tracer,
		storager:         storager,
		addrFileStorager: addrFileStorager,
		wtype:            wtype,
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

func (w *Wallet) GetType() types.WalletType {
	return w.wtype
}
