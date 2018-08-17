package enum

//TxType トランザクション種別
type TxType string

const (
	TxTypeUnsigned TxType = "unsigned"
	TxTypeSigned   TxType = "signed"
	TxTypeSent     TxType = "sent"
	TxTypeDone     TxType = "done"
	TxTypeCancel   TxType = "cancel"
)

//Action 入金/出金
type Action uint8

const (
	ActionReceipt Action = iota + 1
	ActionPayment
)
