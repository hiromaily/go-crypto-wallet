package enum

//NetworkType ネットワーク種別
type NetworkType string

// network type
const (
	NetworkTypeMainNet NetworkType = "mainnet"
	NetworkTypeTestNet3   NetworkType = "testnet3"
)

//TxType トランザクション種別
type TxType string

// tx_type
const (
	TxTypeUnsigned TxType = "unsigned"
	TxTypeSigned   TxType = "signed"
	TxTypeSent     TxType = "sent"
	TxTypeDone     TxType = "done"
	TxTypeNotified TxType = "notified"
	TxTypeCancel   TxType = "canceled"
)

//TxTypeValue tx_typeの値
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned: 1,
	TxTypeSigned:   2,
	TxTypeSent:     3,
	TxTypeDone:     4,
	TxTypeNotified: 5,
	TxTypeCancel:   6,
}

// ValidateTxType TxTypeのバリデーションを行う
func ValidateTxType(val string) bool {
	if _, ok := TxTypeValue[TxType(val)]; ok {
		return true
	}
	return false
}

//ActionType 入金/出金
type ActionType string

// action
const (
	ActionTypeReceipt ActionType = "receipt"
	ActionTypePayment ActionType = "payment"
)

//ActionTypeValue action_typeの値
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeReceipt: 1,
	ActionTypePayment: 2,
}

// ValidateActionType ActionTypeのバリデーションを行う
func ValidateActionType(val string) bool {
	if _, ok := ActionTypeValue[ActionType(val)]; ok {
		return true
	}
	return false
}

//EnvironmentType 実行環境
type EnvironmentType string

// environment
const (
	EnvDev  EnvironmentType = "dev"
	EnvProd EnvironmentType = "prod"
)
