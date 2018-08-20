package enum

//TxType トランザクション種別
type TxType string

// tx_type
const (
	TxTypeUnsigned TxType = "unsigned"
	TxTypeSigned   TxType = "signed"
	TxTypeSent     TxType = "sent"
	TxTypeDone     TxType = "done"
	TxTypeCancel   TxType = "cancel"
)

//TxTypeValue tx_typeの値
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned: 1,
	TxTypeSigned:   2,
	TxTypeSent:     3,
	TxTypeDone:     4,
	TxTypeCancel:   5,
}

//Action 入金/出金
type Action uint8

// action
const (
	ActionReceipt Action = iota + 1
	ActionPayment
)
