package wallet

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/repository/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	wtype "github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// Wallet watch only wallet object
type Wallet struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	tracer       opentracing.Tracer
	repo         walletrepo.WalletRepository
	addrFileRepo address.FileStorager
	txFileRepo   tx.FileStorager
	wtype        wtype.WalletType
}

// NewWallet returns Wallet object
func NewWallet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	repo walletrepo.WalletRepository,
	addrFileRepo address.FileStorager,
	txFileRepo tx.FileStorager,
	wtype wtype.WalletType) *Wallet {

	return &Wallet{
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
func (w *Wallet) Done() {
	w.repo.Close()
	w.btc.Close()
}

// GetDB gets repository
func (w *Wallet) GetDB() walletrepo.WalletRepository {
	return w.repo
}

// GetBTC gets btc
func (w *Wallet) GetBTC() api.Bitcoiner {
	return w.btc
}

// GetType gets wallet type
func (w *Wallet) GetType() wtype.WalletType {
	return w.wtype
}
