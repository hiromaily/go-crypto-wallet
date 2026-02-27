package tx

import "context"

// Transaction abstracts database transaction operations.
// This interface enables database-agnostic implementations,
// allowing repositories to work with any transaction type
// (SQL, NoSQL, in-memory, etc.).
type Transaction interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback rolls back the transaction.
	Rollback() error
}

// UnitOfWork manages transaction lifecycle.
// It provides a consistent API for beginning transactions
// and executing code within a transactional context.
type UnitOfWork interface {
	// Begin starts a new transaction and returns it.
	// The caller is responsible for committing or rolling back.
	Begin(ctx context.Context) (Transaction, error)

	// RunInTransaction executes the provided function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function succeeds, the transaction is committed.
	RunInTransaction(ctx context.Context, fn func(transaction Transaction) error) error
}

// TransactionBeginner is a minimal interface for components that can begin transactions.
// This is useful when only transaction creation is needed, without the full UnitOfWork.
type TransactionBeginner interface {
	// BeginTx starts a new transaction and returns it.
	BeginTx(ctx context.Context) (Transaction, error)
}
