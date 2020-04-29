package action

// ActionType operation (deposit, payment, transfer)
type ActionType string

// action_type
const (
	ActionTypeDeposit  ActionType = "deposit"
	ActionTypePayment  ActionType = "payment"
	ActionTypeTransfer ActionType = "transfer"
)

// String converter
func (a ActionType) String() string {
	return string(a)
}

// ActionTypeValue value
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeDeposit:  1,
	ActionTypePayment:  2,
	ActionTypeTransfer: 3,
}

// ValidateActionType validate
func ValidateActionType(val string) bool {
	if _, ok := ActionTypeValue[ActionType(val)]; ok {
		return true
	}
	return false
}
