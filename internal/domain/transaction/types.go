package transaction

// TxType represents the transaction status/lifecycle state.
//
// Transactions progress through a state machine:
// unsigned → signed → sent → done → (optional: notified or canceled)
type TxType string

// Transaction type constants representing the lifecycle states
const (
	// TxTypeUnsigned means the transaction has been created but not yet signed
	TxTypeUnsigned TxType = "unsigned"

	// TxTypeSigned means the transaction has been signed by required signatories
	TxTypeSigned TxType = "signed"

	// TxTypeSent means the transaction has been broadcast to the network
	TxTypeSent TxType = "sent"

	// TxTypeDone means the transaction has been confirmed on the blockchain
	TxTypeDone TxType = "done"

	// TxTypeNotified means notification has been sent (optional post-confirmation step)
	TxTypeNotified TxType = "notified"

	// TxTypeCancel means the transaction was canceled before being sent
	TxTypeCancel TxType = "canceled"
)

// String returns the string representation of the transaction type.
func (t TxType) String() string {
	return string(t)
}

// Int8 returns the numeric value of the transaction type as int8.
func (t TxType) Int8() int8 {
	return int8(TxTypeValue[t])
}

// Uint8 returns the numeric value of the transaction type as uint8.
func (t TxType) Uint8() uint8 {
	return TxTypeValue[t]
}

// TxTypeValue provides numeric values for transaction types.
// These values are used for database storage and state ordering.
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned: 1,
	TxTypeSigned:   2,
	TxTypeSent:     3,
	TxTypeDone:     4,
	TxTypeNotified: 5,
	TxTypeCancel:   6,
}

// ValidateTxType validates that the given string is a valid transaction type.
func ValidateTxType(val string) bool {
	_, ok := TxTypeValue[TxType(val)]
	return ok
}

// ActionType represents the operation type for a transaction.
//
// Action types define the purpose of a transaction:
//   - Deposit: Collect coins from client accounts to deposit account
//   - Payment: Send coins to external addresses
//   - Transfer: Move coins between internal accounts
type ActionType string

// Action type constants
const (
	// ActionTypeDeposit collects coins from client accounts to deposit account
	ActionTypeDeposit ActionType = "deposit"

	// ActionTypePayment sends coins to external user addresses
	ActionTypePayment ActionType = "payment"

	// ActionTypeTransfer moves coins between internal accounts
	ActionTypeTransfer ActionType = "transfer"
)

// String returns the string representation of the action type.
func (a ActionType) String() string {
	return string(a)
}

// Uint8 returns the numeric value of the action type.
func (a ActionType) Uint8() uint8 {
	return ActionTypeValue[a]
}

// ActionTypeValue provides numeric values for action types.
// These values are used for database storage.
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeDeposit:  1,
	ActionTypePayment:  2,
	ActionTypeTransfer: 3,
}

// ValidateActionType validates that the given string is a valid action type.
func ValidateActionType(val string) bool {
	_, ok := ActionTypeValue[ActionType(val)]
	return ok
}
