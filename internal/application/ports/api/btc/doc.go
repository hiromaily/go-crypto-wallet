// Package btc defines interfaces for Bitcoin/BitcoinCash blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/btc/).
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
//   - Wallet management (EncryptWallet, WalletLock, WalletPassphrase, WalletPassphraseChange, DumpWallet, ImportWallet)
//
// # Usage
//
// Use cases depend on this interface, not concrete implementations:
//
//	import apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
//
//	type myUseCase struct {
//	    btc apibtc.Bitcoiner
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/btc/btc/: BTC implementation
//   - internal/infrastructure/api/btc/bch/: BCH implementation
//   - internal/application/dto/btc/: DTOs used in interface methods
package btc
