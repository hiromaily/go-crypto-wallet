package walletrepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type WalletRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewWalletRepository
func NewWalletRepository(db *sqlx.DB, logger *zap.Logger) *WalletRepository {
	return &WalletRepository{
		db:     db,
		logger: logger,
	}
}

func (r *WalletRepository) Close() error {
	return r.db.Close()
}

func (r *WalletRepository) MustBegin() *sqlx.Tx {
	return r.db.MustBegin()
}
