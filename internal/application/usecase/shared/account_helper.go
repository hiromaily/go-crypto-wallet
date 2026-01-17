package shared

import (
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// InferSenderAccount infers the sender account from the transaction action type.
// This is a pragmatic approach since transaction files don't always encode the account concept.
//
// This helper is shared between BTC and BCH use cases for:
//   - keygen wallet: signing transactions
//   - sign wallet: signing transactions
//
// Transaction flow:
//   - [actionType:deposit]  [from] client [to] deposit (not multisig addr)
//   - [actionType:payment]  [from] payment [to] unknown (multisig addr)
//   - [actionType:transfer] [from] account [to] account (multisig addr)
func InferSenderAccount(actionType domainTx.ActionType) (domainAccount.AccountType, error) {
	switch actionType {
	case domainTx.ActionTypeDeposit:
		return domainAccount.AccountTypeClient, nil
	case domainTx.ActionTypePayment:
		return domainAccount.AccountTypePayment, nil
	case domainTx.ActionTypeTransfer:
		// Transfer could be from various accounts
		// Default to payment for now
		return domainAccount.AccountTypePayment, nil
	default:
		return "", fmt.Errorf("unsupported action type: %s", actionType)
	}
}
