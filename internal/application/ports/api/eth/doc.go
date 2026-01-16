// Package eth defines interfaces for Ethereum blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/eth/).
//
// # Interfaces
//
//   - Ethereumer: Main interface for Ethereum RPC operations
//   - ERC20er: Interface for ERC-20 token operations
//   - EtherTxCreator: Type alias for ERC20er in transaction creation contexts
//   - EtherTxMonitor: Interface for monitoring transaction confirmations
//
// # Types
//
//   - TxCreateParams: DTO for transaction creation parameters
//
// # Capabilities
//
// The Ethereumer interface provides:
//   - Balance queries (GetBalance, GetTotalBalance, BalanceAt)
//   - Transaction operations (CreateRawTransaction, SignOnRawTransaction, SendRawTransaction)
//   - Account management (NewAccount, ImportRawKey, ListAccounts)
//   - Network information (Syncing, BlockNumber, GasPrice)
//   - Mining control (StartMining, StopMining)
//
// # Usage
//
// Use cases depend on these interfaces, not concrete implementations:
//
//	import apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
//
//	type myUseCase struct {
//	    eth apieth.Ethereumer
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/eth/eth/: ETH implementation
//   - internal/infrastructure/api/eth/erc20/: ERC-20 implementation
//   - internal/domain/ethereum/: Domain types used in interface methods
package eth
