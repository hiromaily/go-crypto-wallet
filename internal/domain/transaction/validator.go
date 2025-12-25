package transaction

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// ValidateAmount validates that the transaction amount is valid.
// Amount must be positive and balance must be sufficient.
func ValidateAmount(amount, balance float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive: got %f", amount)
	}

	if amount > balance {
		return fmt.Errorf("insufficient balance: required %f, available %f", amount, balance)
	}

	return nil
}

// ValidateSenderReceiver validates the sender and receiver account types for a transaction.
func ValidateSenderReceiver(sender, receiver account.AccountType, actionType ActionType) error {
	switch actionType {
	case ActionTypeDeposit:
		// Deposit: from client to deposit account
		if sender != account.AccountTypeClient {
			return fmt.Errorf("deposit transactions must have client as sender, got %s", sender)
		}
		if receiver != account.AccountTypeDeposit {
			return fmt.Errorf("deposit transactions must have deposit as receiver, got %s", receiver)
		}

	case ActionTypePayment:
		// Payment: from payment account to external addresses
		if sender != account.AccountTypePayment {
			return fmt.Errorf("payment transactions must have payment as sender, got %s", sender)
		}
		if receiver != account.AccountTypeAnonymous {
			return fmt.Errorf("payment transactions must have anonymous as receiver, got %s", receiver)
		}

	case ActionTypeTransfer:
		// Transfer: between internal accounts (validated by account validator)
		return account.ValidateTransferAccounts(sender, receiver)

	default:
		return fmt.Errorf("invalid action type: %s", actionType)
	}

	return nil
}

// ValidateTransactionType validates that the transaction type is valid.
func ValidateTransactionType(txType TxType) error {
	if !ValidateTxType(string(txType)) {
		return fmt.Errorf("invalid transaction type: %s", txType)
	}
	return nil
}

// CanTransitionTo checks if a transaction can transition from one state to another.
// This enforces the transaction state machine:
// unsigned → signed → sent → done → (optional: notified)
// Cancellation is only allowed before the transaction is confirmed (done)
func CanTransitionTo(from, to TxType) bool {
	// Define valid transitions
	validTransitions := map[TxType][]TxType{
		TxTypeUnsigned: {TxTypeSigned, TxTypeCancel},
		TxTypeSigned:   {TxTypeSent, TxTypeCancel},
		TxTypeSent:     {TxTypeDone, TxTypeCancel},
		TxTypeDone:     {TxTypeNotified},
		TxTypeNotified: {}, // Terminal state
		TxTypeCancel:   {}, // Terminal state
	}

	allowedTransitions, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, allowed := range allowedTransitions {
		if to == allowed {
			return true
		}
	}

	return false
}

// ValidateTransition validates a transaction state transition.
func ValidateTransition(from, to TxType) error {
	if !CanTransitionTo(from, to) {
		return fmt.Errorf("invalid state transition from %s to %s", from, to)
	}
	return nil
}
