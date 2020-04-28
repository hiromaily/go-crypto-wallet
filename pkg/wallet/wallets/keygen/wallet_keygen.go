package keygen

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// Keygen keygen wallet object
//  it is almost same to Wallet object, difference is storager interface
type Keygen struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	repo         coldrepo.ColdRepository
	keyGenerator key.Generator
	wtype        wallet.WalletType
}

// NewKeygen may be Not used anywhere
func NewKeygen(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	repo coldrepo.ColdRepository,
	keyGenerator key.Generator,
	wtype wallet.WalletType) *Keygen {

	return &Keygen{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		repo:         repo,
		keyGenerator: keyGenerator,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *Keygen) Done() {
	w.repo.Close()
	w.btc.Close()
}

// GetBTC gets btc
func (w *Keygen) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Keygen) GetType() wallet.WalletType {
	return w.wtype
}
