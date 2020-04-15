package action

//ActionType operation (receipt, payment, transfer)
type ActionType string

// action
const (
	ActionTypeReceipt  ActionType = "receipt"
	ActionTypePayment  ActionType = "payment"
	ActionTypeTransfer ActionType = "transfer"
)

func (a ActionType) String() string {
	return string(a)
}

//ActionTypeValue value
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeReceipt:  1,
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
