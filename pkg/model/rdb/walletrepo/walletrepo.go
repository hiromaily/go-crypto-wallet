package walletrepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// WalletRepository is repository for wallet
type WalletRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewWalletRepository return WalletRepository
func NewWalletRepository(db *sqlx.DB, logger *zap.Logger) *WalletRepository {
	return &WalletRepository{
		db:     db,
		logger: logger,
	}
}

// Close close db connection
// FIXME: it'd better to return db instance instead
func (r *WalletRepository) Close() error {
	return r.db.Close()
}

// MustBegin start mustBegin
func (r *WalletRepository) MustBegin() *sqlx.Tx {
	return r.db.MustBegin()
}
