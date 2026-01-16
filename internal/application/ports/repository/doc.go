// Package repository defines interfaces for database persistence operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/repository/).
//
// # Cold Wallet Repositories (keygen/sign wallets)
//
//   - SeedRepositorier: HD wallet seed storage
//   - BTCAccountKeyRepositorier: BTC/BCH account key management
//   - ETHAccountKeyRepositorier: ETH account key management
//   - XRPAccountKeyRepositorier: XRP account key management
//   - AuthFullPubkeyRepositorier: Authorization public key storage
//   - AuthAccountKeyRepositorier: Authorization account key storage
//   - HDWalletRepo: Generic HD wallet key storage abstraction
//
// # Watch Wallet Repositories
//
//   - AddressRepositorier: Address management and allocation
//   - BTCTxRepositorier: BTC transaction record management
//   - TxInputRepositorier: Transaction input records
//   - TxOutputRepositorier: Transaction output records
//   - TxRepositorier: Generic transaction records
//   - PaymentRequestRepositorier: Payment request management
//   - ETHDetailTXRepositorier: ETH transaction detail management
//   - XRPDetailTXRepositorier: XRP transaction detail management
//
// # Usage
//
// Use cases depend on these interfaces, not concrete implementations:
//
//	type myUseCase struct {
//	    addressRepo repository.AddressRepositorier
//	    txRepo      repository.BTCTxRepositorier
//	}
//
// # Related Packages
//
//   - internal/infrastructure/repository/cold/: Cold wallet implementations
//   - internal/infrastructure/repository/watch/: Watch wallet implementations
//   - internal/domain/: Domain entities used in interface methods
package repository
