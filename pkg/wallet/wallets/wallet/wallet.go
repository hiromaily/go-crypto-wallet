package wallet

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// Wallet watch only wallet object
type Wallet struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	storager     rdb.WalletStorager
	addrFileRepo address.Storager
	txFileRepo   tx.Storager
	wtype        types.WalletType
}

// NewWallet returns Wallet object
func NewWallet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.WalletStorager,
	addrFileRepo address.Storager,
	txFileRepo tx.Storager,
	wtype types.WalletType) *Wallet {

	return &Wallet{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		storager:     storager,
		addrFileRepo: addrFileRepo,
		txFileRepo:   txFileRepo,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *Wallet) Done() {
	w.storager.Close()
	w.btc.Close()
}

// GetDB gets storager
func (w *Wallet) GetDB() rdb.WalletStorager {
	return w.storager
}

// GetBTC gets btc
func (w *Wallet) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Wallet) GetType() types.WalletType {
	return w.wtype
}