// Package bitcoin defines interfaces for Bitcoin/BitcoinCash blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/bitcoin/).
//
// # Interfaces
//
//   - Bitcoiner: Main interface for Bitcoin/BitcoinCash RPC operations
//
// # Capabilities
//
// The Bitcoiner interface provides:
//   - Address management (GetAddressInfo, ValidateAddress, DecodeAddress)
//   - Balance queries (GetBalance, GetBalanceByAccount)
//   - Transaction operations (CreateRawTransaction, SignRawTransaction, SendTransaction)
//   - PSBT support (CreatePSBT, SignPSBTWithKey, FinalizePSBT) - BIP174
//   - Descriptor wallet support (ImportDescriptors, ListDescriptors)
//   - Multisig operations (AddMultisigAddress)
//   - Wallet management (BackupWallet, LoadWallet, CreateWallet)
//
// # Usage
//
// Use cases depend on this interface, not concrete implementations:
//
//	type myUseCase struct {
//	    btc bitcoin.Bitcoiner
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/bitcoin/btc/: BTC implementation
//   - internal/infrastructure/api/bitcoin/bch/: BCH implementation
//   - internal/application/dto/btc/: DTOs used in interface methods
package bitcoin
