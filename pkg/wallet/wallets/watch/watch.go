package watch

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/repository/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// Watch watch only wallet object
type Watch struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	repo         walletrepo.WalletRepositorier
	addrFileRepo address.FileStorager
	txFileRepo   tx.FileStorager
	wtype        wtype.WalletType
}

// NewWatch returns Watch object
func NewWatch(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	repo walletrepo.WalletRepositorier,
	addrFileRepo address.FileStorager,
	txFileRepo tx.FileStorager,
	wtype wtype.WalletType) *Watch {

	return &Watch{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		repo:         repo,
		addrFileRepo: addrFileRepo,
		txFileRepo:   txFileRepo,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *Watch) Done() {
	w.repo.Close()
	w.btc.Close()
}

// GetDB gets repository
func (w *Watch) GetDB() walletrepo.WalletRepositorier {
	return w.repo
}

// GetBTC gets btc
func (w *Watch) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Watch) GetType() wtype.WalletType {
	return w.wtype
}
