// Package keygen contains application services for keygen wallet operations.
//
// Keygen wallets are offline wallets that generate and manage cryptographic keys.
// In a multisig setup, the keygen wallet provides the first signature. They can:
//   - Generate HD wallet keys from seed
//   - Create and export addresses
//   - Import private keys
//   - Create multisig addresses
//   - Export full public keys for other signers
//   - Provide the first signature in multisig transactions
//
// # Organization
//
// Services are organized by coin type:
//   - shared/: Coin-agnostic services (work across all cryptocurrencies)
//   - btc/: Bitcoin (BTC) and Bitcoin Cash (BCH) specific services
//   - eth/: Ethereum (ETH) and ERC-20 token specific services
//   - xrp/: Ripple (XRP) specific services
//
// # Common Services
//
//   - HDWalleter: Generates HD wallet keys using BIP32/BIP44 standards
//   - Seeder: Generates and stores BIP39 seeds
//   - AddressExporter: Exports generated addresses to files
//   - PrivKeyer: Imports private keys
//   - FullPubKeyImporter: Imports full public keys from other signers
//   - Multisiger: Creates multisig addresses
//
// # Security
//
// Keygen wallets:
//   - Must be kept offline (air-gapped from the network)
//   - Hold private keys and are responsible for generating secure keys
//   - Provide the first signature in multisig setups
//   - Export only public information (addresses, full public keys)
package keygen
