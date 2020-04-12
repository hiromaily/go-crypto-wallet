package action

//ActionType 入金/出金
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

//ActionTypeValue action_typeの値
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeReceipt:  1,
	ActionTypePayment:  2,
	ActionTypeTransfer: 3,
}

// ValidateActionType ActionTypeのバリデーションを行う
func ValidateActionType(val string) bool {
	if _, ok := ActionTypeValue[ActionType(val)]; ok {
		return true
	}
	return false
}

//ActionToAccountMap ActionTypeをAccountTypeでマッピングする
//var ActionToAccountMap = map[ActionType]AccountType{
//	ActionTypeReceipt: AccountTypeClient,
//	ActionTypePayment: AccountTypePayment,
//}
