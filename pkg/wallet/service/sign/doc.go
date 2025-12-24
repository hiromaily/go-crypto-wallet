// Package sign contains application services for sign wallet operations.
//
// Sign wallets are offline wallets that provide additional signatures for multisig transactions.
// In a multisig setup, sign wallets provide the second and subsequent signatures after the keygen wallet.
// They can:
//   - Sign transactions created by the watch wallet
//   - Import private keys
//   - Export full public keys to the keygen wallet
//
// # Organization
//
// Services are organized by coin type:
//   - shared/: Coin-agnostic services (currently empty, but structure ready for future use)
//   - btc/: Bitcoin (BTC) and Bitcoin Cash (BCH) specific services
//   - eth/: Ethereum (ETH) and ERC-20 token specific services
//   - xrp/: Ripple (XRP) specific services
//
// # Common Services
//
//   - Signer: Signs transactions using stored private keys
//   - FullPubkeyExporter: Exports full public keys for keygen wallet
//   - PrivKeyImporter: Imports private keys
//
// # Security
//
// Sign wallets:
//   - Must be kept offline (air-gapped from the network)
//   - Hold private keys used for signing transactions
//   - Provide additional signatures in multisig setups (2-of-N, 3-of-N, etc.)
//   - Export only public information (full public keys)
//   - Are separated from keygen wallets to improve security through separation of duties
package sign
