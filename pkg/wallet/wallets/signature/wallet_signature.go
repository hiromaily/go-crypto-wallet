package signature

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// Signature signature wallet object
//  it is almost same to Wallet object, difference is storager interface
type Signature struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	repo         rdb.ColdStorager
	keyGenerator key.Generator
	wtype        types.WalletType
}

// NewSignature may be Not used anywhere
func NewSignature(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	repo rdb.ColdStorager,
	keyGenerator key.Generator,
	wtype types.WalletType) *Signature {

	return &Signature{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		repo:         repo,
		keyGenerator: keyGenerator,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *Signature) Done() {
	w.repo.Close()
	w.btc.Close()
}

// GetDB gets storager
func (w *Signature) GetDB() rdb.ColdStorager {
	return w.repo
}

// GetBTC gets btc
func (w *Signature) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Signature) GetType() types.WalletType {
	return w.wtype
}
