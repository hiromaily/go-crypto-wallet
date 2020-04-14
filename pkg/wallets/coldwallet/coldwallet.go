package coldwallet

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// ColdWallet coldwallet for keygen/signature object
type ColdWallet struct {
	btc              api.Bitcoiner
	logger           *zap.Logger
	tracer           opentracing.Tracer
	storager         rdb.ColdStorager
	addrFileStorager address.Storager
	wtype            types.WalletType
}

func NewColdWalet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	storager rdb.ColdStorager,
	addrFileStorager address.Storager,
	wtype types.WalletType) *ColdWallet {

	return &ColdWallet{
		btc:              btc,
		logger:           logger,
		tracer:           tracer,
		storager:         storager,
		addrFileStorager: addrFileStorager,
		wtype:            wtype,
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
