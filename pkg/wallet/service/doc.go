// Package service contains the application layer of the cryptocurrency wallet system.
//
// The application layer is organized by wallet type following Clean Architecture principles:
//
// Organization Pattern: wallet_type → coin → use_case
//
// # Wallet Types
//
// The application layer is divided into three wallet types based on their operational characteristics:
//
//   - watch/: Watch wallet services (online, public keys only)
//     Creates and sends transactions, monitors transaction status, imports addresses from external sources.
//
//   - keygen/: Keygen wallet services (offline, first signature)
//     Generates HD wallet keys, imports private keys, creates multisig addresses, exports addresses and public keys.
//
//   - sign/: Sign wallet services (offline, subsequent signatures)
//     Signs transactions, imports private keys, exports public keys.
//
// # Directory Structure
//
// Each wallet type follows a consistent structure:
//
//	wallet/service/
//	├── watch/                      # Watch wallet services
//	│   ├── interfaces.go           # Watch wallet interfaces
//	│   ├── shared/                 # Coin-agnostic services
//	│   │   ├── address_importer.go
//	│   │   └── payment_request_creator.go
//	│   ├── btc/                    # BTC/BCH watch wallet services
//	│   │   ├── tx_creator.go
//	│   │   ├── tx_monitor.go
//	│   │   └── tx_sender.go
//	│   ├── eth/                    # ETH watch wallet services
//	│   └── xrp/                    # XRP watch wallet services
//	├── keygen/                     # Keygen wallet services
//	│   ├── interfaces.go           # Keygen wallet interfaces
//	│   ├── shared/                 # Coin-agnostic services
//	│   │   ├── hd_walleter.go
//	│   │   ├── seeder.go
//	│   │   └── address_exporter.go
//	│   ├── btc/                    # BTC/BCH keygen services
//	│   ├── eth/                    # ETH keygen services
//	│   └── xrp/                    # XRP keygen services
//	└── sign/                       # Sign wallet services
//	    ├── interfaces.go           # Sign wallet interfaces
//	    ├── shared/                 # Coin-agnostic services (currently empty)
//	    ├── btc/                    # BTC/BCH sign services
//	    ├── eth/                    # ETH sign services
//	    └── xrp/                    # XRP sign services
//
// # Common Use Cases
//
//   - Transaction Creation (tx_creator.go): Creates deposit, payment, or transfer transactions
//   - Transaction Monitoring (tx_monitor.go): Monitors transaction status and balances
//   - Transaction Sending (tx_sender.go): Sends signed transactions to the network
//   - Key Generation (hd_walleter.go, key_generator.go): Generates cryptographic keys
//   - Key Import/Export (privkey_importer.go, fullpubkey_importer.go, etc.): Imports/exports keys
//   - Address Management (address_importer.go, address_exporter.go): Manages addresses
//   - Multisig Operations (multisigaddress.go): Creates and manages multisig addresses
//   - Signing (signer.go): Signs transactions
//
// # Backward Compatibility
//
// Type aliases are provided in the root service package for backward compatibility:
//
//   - watch_interface.go: Type aliases for watch wallet interfaces
//   - cold_interface.go: Type aliases for keygen and sign wallet interfaces
//
// # Architecture Guidelines
//
//   - Application layer depends on domain layer (pkg/domain/)
//   - Application layer uses infrastructure through domain interfaces
//   - Services should be focused on single use cases
//   - Use dependency injection for all dependencies
//   - Interfaces are co-located with implementations in each wallet type directory
package service
