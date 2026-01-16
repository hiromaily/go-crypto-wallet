// Package xrp defines interfaces for Ripple/XRP blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/xrp/).
//
// # Interfaces
//
//   - Rippler: Main interface combining admin, public, and API operations
//   - RippleAPIer: Interface for Ripple API operations (account, address, transaction)
//   - RipplePublicer: Interface for public node operations
//   - RippleAdminer: Interface for admin node operations
//
// # Capabilities
//
// The Rippler interface provides:
//   - Balance queries (GetBalance, GetTotalBalance)
//   - Transaction operations (CreateRawTransaction, SignTransaction, SubmitTransaction)
//   - Address generation (GenerateAddress, GenerateXAddress)
//   - Account information (AccountInfo, GetAccountInfo)
//   - Server information (ServerInfo)
//
// # Usage
//
// Use cases depend on these interfaces, not concrete implementations:
//
//	import apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
//
//	type myUseCase struct {
//	    xrp apixrp.Rippler
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/xrp/xrp/: XRP implementation
//   - internal/application/dto/ripple/: DTOs used in interface methods
package xrp
