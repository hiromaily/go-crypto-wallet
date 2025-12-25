// Package transaction provides domain entities and types for transaction management.
//
// This package contains pure business logic related to transactions including:
//   - Transaction lifecycle states (unsigned, signed, sent, done, notified, canceled)
//   - Action types (deposit, payment, transfer)
//   - Transaction validation rules
//   - State machine for transaction transitions
//   - Amount and balance validation
//
// The package enforces business rules such as:
//   - Transaction state transitions must follow the defined state machine
//   - Amounts must be positive and within available balance
//   - Sender/receiver combinations must be valid for the action type
//
// This package has no infrastructure dependencies and can be tested in isolation.
package transaction
