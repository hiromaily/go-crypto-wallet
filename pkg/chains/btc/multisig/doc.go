// Package multisig provides domain entities and types for multisig address management.
//
// This package contains pure business logic related to multisig (M-of-N) addresses
// and MuSig2 (Schnorr multisignatures):
//
// Traditional Multisig (M-of-N):
//   - Multisig configuration validation (M must be <= N, practical limits)
//   - Public key validation
//   - Redeem script validation
//   - Account eligibility for multisig
//   - Authorization account validation
//
// MuSig2 (Schnorr Multisignatures):
//   - Nonce commitment types and validation
//   - Partial signature types and validation
//   - Aggregated key types
//   - Signing session state management
//   - Signer count validation
//   - Nonce uniqueness enforcement (critical for security)
//
// The package enforces business rules such as:
//   - Multisig requires at least 2 total signers
//   - Required signatures cannot exceed total signatures
//   - Only certain account types can use multisig (deposit, payment, stored)
//   - Authorization accounts must be proper auth accounts
//   - MuSig2 nonces must be unique (nonce reuse leaks private keys)
//   - MuSig2 partial signatures must be from distinct signers
//
// This package has no infrastructure dependencies and can be tested in isolation.
// Actual multisig address generation and MuSig2 cryptographic operations remain
// in infrastructure layer as they require blockchain-specific libraries.
package multisig
