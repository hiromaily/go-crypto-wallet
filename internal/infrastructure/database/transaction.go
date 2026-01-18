package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
)

// Compile-time interface compliance check
var (
	_ persistence.Transaction         = (*SQLTransaction)(nil)
	_ persistence.UnitOfWork          = (*SQLUnitOfWork)(nil)
	_ persistence.TransactionBeginner = (*SQLUnitOfWork)(nil)
)

// SQLTransaction wraps *sql.Tx to implement persistence.Transaction.
// This allows repository implementations to work with the abstract
// Transaction interface while still using SQL transactions internally.
type SQLTransaction struct {
	tx *sql.Tx
}

// NewSQLTransaction creates a new SQLTransaction wrapping the given *sql.Tx.
func NewSQLTransaction(tx *sql.Tx) *SQLTransaction {
	return &SQLTransaction{tx: tx}
}

// Commit commits the transaction.
func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction.
func (t *SQLTransaction) Rollback() error {
	return t.tx.Rollback()
}

// UnwrapTx returns the underlying *sql.Tx.
// This is used by repository implementations that need access to the raw SQL transaction.
func (t *SQLTransaction) UnwrapTx() *sql.Tx {
	return t.tx
}

// SQLUnitOfWork provides transaction management for SQL databases.
// It implements both persistence.UnitOfWork and persistence.TransactionBeginner.
type SQLUnitOfWork struct {
	db *sql.DB
}

// NewSQLUnitOfWork creates a new SQLUnitOfWork with the given database connection.
func NewSQLUnitOfWork(db *sql.DB) *SQLUnitOfWork {
	return &SQLUnitOfWork{db: db}
}

// Begin starts a new transaction and returns it.
func (u *SQLUnitOfWork) Begin(ctx context.Context) (persistence.Transaction, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return NewSQLTransaction(tx), nil
}

// BeginTx starts a new transaction and returns it (implements TransactionBeginner).
func (u *SQLUnitOfWork) BeginTx(ctx context.Context) (persistence.Transaction, error) {
	return u.Begin(ctx)
}

// RunInTransaction executes the provided function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function succeeds, the transaction is committed.
func (u *SQLUnitOfWork) RunInTransaction(ctx context.Context, fn func(tx persistence.Transaction) error) error {
	tx, err := u.Begin(ctx)
	if err != nil {
		return err
	}

	// Ensure rollback on panic or error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}

// UnwrapSQLTx extracts the underlying *sql.Tx from a persistence.Transaction.
// Returns nil if the transaction is not a SQLTransaction.
// This is a helper function for repository implementations.
func UnwrapSQLTx(tx persistence.Transaction) *sql.Tx {
	if sqlTx, ok := tx.(*SQLTransaction); ok {
		return sqlTx.UnwrapTx()
	}
	return nil
}
