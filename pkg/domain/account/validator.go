package account

import (
	"errors"
	"fmt"
)

// CanTransferFrom returns true if the account type can be used as a sender in transfers.
// Client and Authorization accounts cannot be senders in transfer operations.
func CanTransferFrom(accountType AccountType) bool {
	// Client accounts are for users and should not be used as senders
	// Authorization accounts are for multisig signing only
	return accountType != AccountTypeClient && accountType != AccountTypeAuthorization
}

// CanReceiveTo returns true if the account type can be used as a receiver in transfers.
// Client and Authorization accounts cannot be receivers in transfer operations.
func CanReceiveTo(accountType AccountType) bool {
	// Client accounts are for users and should not be used as receivers in internal transfers
	// Authorization accounts are for multisig signing only
	return accountType != AccountTypeClient && accountType != AccountTypeAuthorization
}

// ValidateTransferAccounts validates sender and receiver account types for transfer operations.
// Returns an error if the accounts are invalid for transfers.
func ValidateTransferAccounts(sender, receiver AccountType) error {
	if !CanTransferFrom(sender) {
		return fmt.Errorf("account type %s cannot be used as sender in transfer", sender)
	}

	if !CanReceiveTo(receiver) {
		return fmt.Errorf("account type %s cannot be used as receiver in transfer", receiver)
	}

	if sender == receiver {
		return errors.New("sender and receiver must be different accounts")
	}

	return nil
}

// IsAuthorizationAccount returns true if the account type is an authorization account.
func IsAuthorizationAccount(accountType AccountType) bool {
	return accountType == AccountTypeAuthorization ||
		accountType == AccountTypeAuth1 ||
		accountType == AccountTypeAuth2 ||
		accountType == AccountTypeAuth3 ||
		accountType == AccountTypeAuth4 ||
		accountType == AccountTypeAuth5 ||
		accountType == AccountTypeAuth6 ||
		accountType == AccountTypeAuth7 ||
		accountType == AccountTypeAuth8 ||
		accountType == AccountTypeAuth9 ||
		accountType == AccountTypeAuth10 ||
		accountType == AccountTypeAuth11 ||
		accountType == AccountTypeAuth12 ||
		accountType == AccountTypeAuth13 ||
		accountType == AccountTypeAuth14 ||
		accountType == AccountTypeAuth15
}

// IsMultisigEligible returns true if the account type can have multisig addresses.
// Only Deposit, Payment, and Stored accounts support multisig.
func IsMultisigEligible(accountType AccountType) bool {
	return accountType == AccountTypeDeposit ||
		accountType == AccountTypePayment ||
		accountType == AccountTypeStored
}
