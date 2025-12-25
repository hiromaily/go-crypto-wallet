// Package key provides domain entities and types for cryptographic key management.
//
// This package contains pure business logic related to HD wallet key generation
// and management, including:
//   - WalletKey value object with multiple address formats
//   - Key index validation (BIP32/BIP44 compliance)
//   - Seed validation
//   - Key range validation
//
// The package enforces business rules such as:
//   - Key indices must be within valid BIP32 ranges
//   - Seeds must meet minimum entropy requirements
//   - Generated keys must have required address formats
//
// This package has no infrastructure dependencies and can be tested in isolation.
// Actual key generation logic remains in infrastructure layer as it requires
// cryptographic libraries.
package key
