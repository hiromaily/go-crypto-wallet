package walletrepo

import (
	"database/sql"

	"go.uber.org/zap"
)

// WalletRepository is wallet repository interface
type WalletRepository interface {
	Close() error
	BeginTx() (*sql.Tx, error)
	Tx() TxRepository
	TxInput() TxInputRepository
	TxOutput() TxOutputRepository
	PayReq() PaymentRequestRepository
	Addr() AddressRepository
}

// walletRepository is repository for wallet
type walletRepository struct {
	db     *sql.DB
	logger *zap.Logger
	TxRepository
	TxInputRepository
	TxOutputRepository
	PaymentRequestRepository
	AddressRepository
}

// NewWalletRepository returns WalletRepository
// nolint:golint
func NewWalletRepository(
	db *sql.DB,
	logger *zap.Logger,
	txRepo TxRepository,
	txInRepo TxInputRepository,
	txOutRepo TxOutputRepository,
	payReqRepo PaymentRequestRepository,
	addrRepo AddressRepository) *walletRepository {

	return &walletRepository{
		db:                       db,
		logger:                   logger,
		TxRepository:             txRepo,
		TxInputRepository:        txInRepo,
		TxOutputRepository:       txOutRepo,
		PaymentRequestRepository: payReqRepo,
		AddressRepository:        addrRepo,
	}
}

// Close close db connection
func (r *walletRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// BeginTx starts transaction
func (r *walletRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *walletRepository) Tx() TxRepository {
	return r.TxRepository
}

func (r *walletRepository) TxInput() TxInputRepository {
	return r.TxInputRepository
}

func (r *walletRepository) TxOutput() TxOutputRepository {
	return r.TxOutputRepository
}

func (r *walletRepository) PayReq() PaymentRequestRepository {
	return r.PaymentRequestRepository
}

func (r *walletRepository) Addr() AddressRepository {
	return r.AddressRepository
}
