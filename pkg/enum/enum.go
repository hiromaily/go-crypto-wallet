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
