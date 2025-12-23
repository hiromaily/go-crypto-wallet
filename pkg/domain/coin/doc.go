// Package coin provides domain entities and types for cryptocurrency definitions.
//
// This package contains pure business logic related to supported cryptocurrencies:
//   - Coin types following SLIP-0044 standard for HD wallet derivation
//   - Coin type codes for human-readable identifiers
//   - ERC20 token definitions
//   - Coin grouping (BTC group, ETH group)
//
// Supported cryptocurrencies:
//   - Bitcoin (BTC)
//   - Bitcoin Cash (BCH)
//   - Litecoin (LTC)
//   - Ethereum (ETH)
//   - Ripple (XRP)
//   - ERC20 tokens (HYT, BAT, and others)
//
// This package has no infrastructure dependencies and can be tested in isolation.
package coin
