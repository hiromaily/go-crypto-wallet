// Package xrp defines interfaces for XRP blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/xrp/).
//
// # Interfaces
//
//   - XRPer: Main interface combining admin, public, and API operations
//   - XRPAPIProvider: Interface for XRP API operations (account, address, transaction)
//   - XRPPublicer: Interface for public node operations
//   - XRPAdminer: Interface for admin node operations
//
// # Capabilities
//
// The XRPer interface provides:
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
//	    xrp apixrp.XRPer
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/xrp/: XRP implementation
//   - internal/application/dto/xrp/: DTOs used in interface methods
package xrp
