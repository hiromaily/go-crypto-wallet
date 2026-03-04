package xrp

import (
	"context"

	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
)

// AccountInfoProvider defines a minimal interface for querying XRP account information.
//
// This interface follows the Interface Segregation Principle by providing only
// the methods needed by the CreateTransactionUseCase, avoiding dependency on the
// full XRPer interface.
//
// # Responsibilities
//   - Query account information from XRP Ledger
//   - Query account balance
//
// # Constraints
//   - Read-only operations (no state modifications)
//   - Requires network connectivity (online operation)
//   - No transaction signing or submission methods
//
// # Implementation Notes
//   - Implemented by xrpscan/xrpl-go client in infrastructure layer
//   - Used by CreateTransactionUseCase to validate accounts before transaction creation
//
// # Related Interfaces
//   - TransactionSigner: For signing operations
//   - TransactionSubmitter: For submitting transactions
type AccountInfoProvider interface {
	// GetAccountInfo retrieves comprehensive account information from XRP Ledger.
	//
	// This method queries the ledger for account details including sequence number,
	// balance, flags, and multi-signature configuration (if applicable).
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - address: XRP address in classic format (rXXX) or X-address format
	//
	// Returns:
	//   - ResponseGetAccountInfo: Account details including sequence, balance, and flags
	//   - error: Returns error if:
	//     * address is invalid or malformed
	//     * account does not exist on the ledger
	//     * network communication fails
	//     * ledger is unavailable
	//
	// Preconditions:
	//   - address must be valid XRP address format (rXXX classic or XXXX X-address)
	//   - network connectivity required (online operation)
	//
	// Postconditions:
	//   - Returns current account state from ledger
	//   - No ledger state modifications
	//
	// Example Usage:
	//   accountInfo, err := provider.GetAccountInfo(ctx, "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH")
	//   if err != nil {
	//       return fmt.Errorf("failed to get account info: %w", err)
	//   }
	//   sequence := accountInfo.AccountData.Sequence
	GetAccountInfo(ctx context.Context, address string) (*xrplgo.AccountInfo, error)

	// GetBalance retrieves the current XRP balance for a given account.
	//
	// This method queries the ledger for the account's XRP balance and returns
	// it as a float64 value representing XRP (not drops).
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - addr: XRP address in classic format (rXXX) or X-address format
	//
	// Returns:
	//   - float64: Account balance in XRP (1 XRP = 1,000,000 drops)
	//   - error: Returns error if:
	//     * address is invalid or malformed
	//     * account does not exist on the ledger
	//     * network communication fails
	//     * ledger is unavailable
	//
	// Preconditions:
	//   - addr must be valid XRP address format
	//   - network connectivity required (online operation)
	//
	// Postconditions:
	//   - Returns balance as XRP (float64), not drops (integer)
	//   - Balance >= 0 (accounts cannot have negative balance)
	//   - No ledger state modifications
	//
	// Example Usage:
	//   balance, err := provider.GetBalance(ctx, "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH")
	//   if err != nil {
	//       return fmt.Errorf("failed to get balance: %w", err)
	//   }
	//   if balance < requiredAmount {
	//       return errors.New("insufficient balance")
	//   }
	GetBalance(ctx context.Context, addr string) (float64, error)
}
