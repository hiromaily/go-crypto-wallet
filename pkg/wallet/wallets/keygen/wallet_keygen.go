package keygen

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// Keygen keygen wallet object
//  it is almost same to Wallet object, difference is storager interface
type Keygen struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	storager     rdb.ColdStorager
	keyGenerator key.Generator
	wtype        types.WalletType
}

// NewKeygen may be Not used anywhere
func NewKeygen(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.ColdStorager,
	keyGenerator key.Generator,
	wtype types.WalletType) *Keygen {

	return &Keygen{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		storager:     storager,
		keyGenerator: keyGenerator,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *Keygen) Done() {
	w.storager.Close()
	w.btc.Close()
}

// GetDB gets storager
func (w *Keygen) GetDB() rdb.ColdStorager {
	return w.storager
}

// GetBTC gets btc
func (w *Keygen) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Keygen) GetType() types.WalletType {
	return w.wtype
}