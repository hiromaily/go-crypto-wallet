package enum

//TxType トランザクション種別
type TxType string

const (
	TxTypeUnsigned = "unsigned"
	TxTypeSigned   = "signed"
	TxTypeSent     = "sent"
	TxTypeDone     = "done"
	TxTypeCancel   = "cancel"
)
