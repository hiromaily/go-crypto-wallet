// Package watch contains application services for watch wallet operations.
//
// Watch wallets are online wallets that hold only public keys. They can:
//   - Monitor addresses and balances
//   - Create unsigned transactions
//   - Send signed transactions to the blockchain
//   - Import addresses from external sources
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
//   - TxCreator: Creates deposit, payment, and transfer transactions
//   - TxMonitorer: Monitors transaction status and updates balances
//   - TxSender: Sends signed transactions to the blockchain network
//   - AddressImporter: Imports addresses from files
//   - PaymentRequestCreator: Creates payment requests
//
// # Security
//
// Watch wallets:
//   - Only have public keys (cannot sign transactions)
//   - Are designed to run online and communicate with blockchain nodes
//   - Create unsigned transactions that must be signed by offline wallets
package watch
