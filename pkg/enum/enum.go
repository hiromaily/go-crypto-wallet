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

//Action 入金/出金
type Action uint8

// action
const (
	ActionReceipt Action = iota + 1
	ActionPayment
)
