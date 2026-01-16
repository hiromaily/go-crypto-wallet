// Package repository is the parent directory for repository port interfaces.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/repository/).
//
// # Subdirectories
//
//   - cold/: Repository interfaces for cold wallets (keygen/sign)
//   - watch/: Repository interfaces for watch wallets (online)
//
// # Usage
//
// Import the appropriate subpackage based on wallet type:
//
//	import (
//	    repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
//	    repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
//	)
//
// # Related Packages
//
//   - internal/infrastructure/repository/cold/mysql/: Cold wallet MySQL implementations
//   - internal/infrastructure/repository/watch/mysql/: Watch wallet MySQL implementations
//   - internal/domain/: Domain entities used in interface methods
package repository
