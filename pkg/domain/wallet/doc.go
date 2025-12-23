// Package wallet provides domain entities and types for wallet management.
//
// This package contains pure business logic related to wallet types and operations.
// It defines the three wallet types used in the system:
//   - Watch-only wallet: Online, holds public keys only
//   - Keygen wallet: Offline, generates keys and first multisig signature
//   - Sign wallet: Offline, provides additional multisig signatures
//
// This package has no infrastructure dependencies and can be tested in isolation.
package wallet
