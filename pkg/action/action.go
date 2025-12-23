package action

import (
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction instead.
// This package provides backward compatibility aliases.

// ActionType operation (deposit, payment, transfer)
// Deprecated: Use domain/transaction.ActionType
type ActionType = domainTx.ActionType

// action_type
// Deprecated: Use constants from domain/transaction package
const (
	ActionTypeDeposit  = domainTx.ActionTypeDeposit
	ActionTypePayment  = domainTx.ActionTypePayment
	ActionTypeTransfer = domainTx.ActionTypeTransfer
)

// ActionTypeValue value
// Deprecated: Use domain/transaction.ActionTypeValue
var ActionTypeValue = domainTx.ActionTypeValue

// ValidateActionType validate
// Deprecated: Use domain/transaction.ValidateActionType
func ValidateActionType(val string) bool {
	return domainTx.ValidateActionType(val)
}
