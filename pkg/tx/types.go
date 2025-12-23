package tx

import (
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction instead.
// This package provides backward compatibility aliases.

//----------------------------------------------------
// TxType
//----------------------------------------------------

// TxType transaction status
//
// Deprecated: Use domain/transaction.TxType
type TxType = domainTx.TxType

// tx_type
//
// Deprecated: Use constants from domain/transaction package
const (
	TxTypeUnsigned = domainTx.TxTypeUnsigned
	TxTypeSigned   = domainTx.TxTypeSigned
	TxTypeSent     = domainTx.TxTypeSent
	TxTypeDone     = domainTx.TxTypeDone
	TxTypeNotified = domainTx.TxTypeNotified
	TxTypeCancel   = domainTx.TxTypeCancel
)

// TxTypeValue value
//
// Deprecated: Use domain/transaction.TxTypeValue
var TxTypeValue = domainTx.TxTypeValue

// ValidateTxType validate string
//
// Deprecated: Use domain/transaction.ValidateTxType
func ValidateTxType(val string) bool {
	return domainTx.ValidateTxType(val)
}
