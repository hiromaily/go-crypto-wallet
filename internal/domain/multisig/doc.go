// Package multisig provides domain entities and types for multisig address management.
//
// This package contains pure business logic related to multisig (M-of-N) addresses:
//   - Multisig configuration validation (M must be <= N, practical limits)
//   - Public key validation
//   - Redeem script validation
//   - Account eligibility for multisig
//   - Authorization account validation
//
// The package enforces business rules such as:
//   - Multisig requires at least 2 total signers
//   - Required signatures cannot exceed total signatures
//   - Only certain account types can use multisig (deposit, payment, stored)
//   - Authorization accounts must be proper auth accounts
//
// This package has no infrastructure dependencies and can be tested in isolation.
// Actual multisig address generation remains in infrastructure layer as it requires
// blockchain-specific libraries.
package multisig
