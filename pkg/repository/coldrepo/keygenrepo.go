package coldrepo

import (
	"database/sql"

	"go.uber.org/zap"
)

// KeygenRepository is keygen wallet repository interface
type KeygenRepository interface {
	Close() error
	BeginTx() (*sql.Tx, error)
	Seed() SeedRepository
	AccountKey() AccountKeyRepository
	//MultisigHistory() MultisigHistoryRepository
}

// keygenRepository is repository for keygen wallet
type keygenRepository struct {
	db     *sql.DB
	logger *zap.Logger
	SeedRepository
	AccountKeyRepository
}

// NewKeygenWalletRepository returns WalletRepository
func NewKeygenWalletRepository(
	db *sql.DB,
	logger *zap.Logger,
	seedRepo SeedRepository,
	accountKeyRepo AccountKeyRepository) *keygenRepository {

	return &keygenRepository{
		db:                   db,
		logger:               logger,
		SeedRepository:       seedRepo,
		AccountKeyRepository: accountKeyRepo,
	}
}

// Close close db connection
func (r *keygenRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// BeginTx starts transaction
func (r *keygenRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *keygenRepository) Seed() SeedRepository {
	return r.SeedRepository
}

func (r *keygenRepository) AccountKey() AccountKeyRepository {
	return r.AccountKeyRepository
}
