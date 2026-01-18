// Package persistence provides abstractions for database transaction operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining transaction interfaces in the application layer that abstract away
// the underlying database implementation details.
//
// # Purpose
//
// The interfaces in this package enable:
//   - Database-agnostic transaction operations
//   - In-memory repository implementations for testing
//   - NoSQL backend implementations
//   - Mock repositories without SQL dependencies
//
// # Usage
//
//	import (
//	    "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
//	)
//
//	func DoSomethingWithTransaction(uow persistence.UnitOfWork) error {
//	    return uow.RunInTransaction(ctx, func(tx persistence.Transaction) error {
//	        // perform operations with tx
//	        return nil
//	    })
//	}
//
// # Related Packages
//
//   - internal/infrastructure/database/mysql/: MySQL implementation
//   - internal/infrastructure/database/sqlite/: SQLite implementation
//   - internal/application/ports/repository/: Repository interfaces using Transaction
package persistence
