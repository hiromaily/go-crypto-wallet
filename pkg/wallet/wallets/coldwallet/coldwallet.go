package coldwallet

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// ColdWallet coldwallet for keygen/signature object
type ColdWallet struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	storager     rdb.ColdStorager
	keyGenerator key.Generator
	addrFileRepo address.Storager
	txFileRepo   tx.Storager
	wtype        types.WalletType
}

// NewColdWalet returns ColdWallet instance
func NewColdWalet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.ColdStorager,
	keyGenerator key.Generator,
	addrFileRepo address.Storager,
	txFileRepo tx.Storager,
	wtype types.WalletType) *ColdWallet {

	return &ColdWallet{
		btc:          btc,
		logger:       logger,
		tracer:       tracer,
		storager:     storager,
		keyGenerator: keyGenerator,
		addrFileRepo: addrFileRepo,
		txFileRepo:   txFileRepo,
		wtype:        wtype,
	}
}

// Done should be called before exit
func (w *ColdWallet) Done() {
	w.storager.Close()
	w.btc.Close()
}

// GetDB gets storager
func (w *ColdWallet) GetDB() rdb.ColdStorager {
	return w.storager
}

// GetBTC gets btc
func (w *ColdWallet) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *ColdWallet) GetType() types.WalletType {
	return w.wtype
}
