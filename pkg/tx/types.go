package tx

//----------------------------------------------------
// TxType
//----------------------------------------------------

// TxType transaction status
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

// String converter
func (t TxType) String() string {
	return string(t)
}

// Int8 converter
func (t TxType) Int8() int8 {
	return int8(TxTypeValue[t])
}

// TxTypeValue value
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned: 1,
	TxTypeSigned:   2,
	TxTypeSent:     3,
	TxTypeDone:     4,
	TxTypeNotified: 5,
	TxTypeCancel:   6,
}

// ValidateTxType validate string
func ValidateTxType(val string) bool {
	if _, ok := TxTypeValue[TxType(val)]; ok {
		return true
	}
	return false
}
