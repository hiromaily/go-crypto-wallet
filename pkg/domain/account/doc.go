// Package account provides domain entities and types for account management.
//
// This package contains pure business logic related to account types and validation.
// Accounts represent different utilization purposes of addresses:
//   - Client: User-created addresses
//   - Deposit: Aggregation addresses
//   - Payment: Payment sender addresses
//   - Stored: Cold storage addresses
//   - Authorization (auth1-auth15): Multisig signer accounts
//   - Anonymous: External receiver addresses
//
// The package includes business rules for:
//   - Account type validation
//   - Transfer eligibility checks
//   - Multisig account validation
//   - Authorization account identification
//
// This package has no infrastructure dependencies and can be tested in isolation.
package account
